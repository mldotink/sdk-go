//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

const harborSandboxImage = "public.ecr.aws/s2s3s5b7/ink/ink-graphql:physicslab-harbor-e2e-30353546022"

type sandboxResult struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Image string `json:"image"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type sandboxExecutionResult struct {
	ID       string `json:"id"`
	State    string `json:"state"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode *int   `json:"exitCode"`
	Error    *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestSandboxHarborEnvironment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()

	var created struct {
		Sandbox sandboxResult `json:"sandboxCreate"`
	}
	sandboxName := name("harbor")
	requireNoError(t, sandboxGraphQL(ctx, `
		mutation CreateHarborSandbox($input: CreateSandboxInput!) {
			sandboxCreate(input: $input) {
				id
				state
				image
				error { code message }
			}
		}`,
		map[string]any{"input": map[string]any{
			"idempotencyKey":        "harbor-create-" + runID,
			"workspaceSlug":         workspace,
			"project":               project,
			"name":                  sandboxName,
			"image":                 harborSandboxImage,
			"resources":             map[string]any{"vcpus": 3, "memoryGb": 2},
			"network":               map[string]any{"egress": "NONE"},
			"destroyTimeoutSeconds": 900,
			"metadata": []map[string]string{
				{"key": "e2e", "value": "harbor-environment"},
			},
		}}, &created), "create Harbor sandbox")
	requireNotEmpty(t, created.Sandbox.ID, "sandbox id is empty")
	requireEqual(t, created.Sandbox.Image, harborSandboxImage, "sandbox image")
	t.Logf("created sandbox %s (id: %s)", sandboxName, created.Sandbox.ID)

	terminated := false
	t.Cleanup(func() {
		if terminated {
			return
		}
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cleanupCancel()
		if err := terminateSandbox(cleanupCtx, created.Sandbox.ID); err != nil {
			t.Logf("cleanup: terminate sandbox %s: %v", created.Sandbox.ID, err)
		}
	})

	running := waitForSandboxState(t, ctx, created.Sandbox.ID, "RUNNING")
	if running.Error != nil {
		t.Fatalf("sandbox reached RUNNING with error %s: %s", running.Error.Code, running.Error.Message)
	}

	checkScript := `set -eu
test "$(uname -m)" = "x86_64"
test "$HOME" = "/workspace"
test -w /workspace
test -w /app/submission
touch /workspace/.ink-e2e-write
touch /app/submission/.ink-e2e-write
test ! -e /root/.codex
test ! -e /root/.claude
test ! -e /root/.claude.json
test ! -e /root/.grok
test ! -e /root/.config/opencode
test -z "$(find /app/data -mindepth 1 -print -quit)"
test "$(codex --version)" = "codex-cli 0.144.6"
claude --version | grep -F "2.1.216"
grok --version | grep -F "grok 0.2.106"
test "$(opencode --version)" = "1.18.4"
test "$(node --version)" = "v22.17.0"
/usr/local/bin/harbor-offline-shell --self-test
/usr/local/bin/harbor-grok-launch --self-test
python - <<'PY'
import matplotlib, mpmath, numpy, scipy, sympy
assert matplotlib.__version__ == "3.10.3"
assert mpmath.__version__ == "1.3.0"
assert numpy.__version__ == "2.2.6"
assert scipy.__version__ == "1.15.3"
assert sympy.__version__ == "1.14.0"
print("scientific-stack-ok")
PY
echo harbor-runtime-ok`

	var executionCreated struct {
		Execution sandboxExecutionResult `json:"sandboxExecutionCreate"`
	}
	requireNoError(t, sandboxGraphQL(ctx, `
		mutation CheckHarborEnvironment($input: CreateSandboxExecutionInput!) {
			sandboxExecutionCreate(input: $input) {
				id
				state
				stdout
				stderr
				exitCode
				error { code message }
			}
		}`,
		map[string]any{"input": map[string]any{
			"idempotencyKey": "harbor-exec-" + runID,
			"sandboxId":      created.Sandbox.ID,
			"workspaceSlug":  workspace,
			"argv":           []string{"/bin/sh", "-lc", checkScript},
			"cwd":            "/app",
			"timeoutSeconds": 180,
			"maxOutputBytes": 1 << 20,
		}}, &executionCreated), "create Harbor validation execution")
	requireNotEmpty(t, executionCreated.Execution.ID, "execution id is empty")

	execution := waitForSandboxExecution(t, ctx, created.Sandbox.ID, executionCreated.Execution.ID)
	if execution.State != "SUCCEEDED" {
		t.Fatalf("Harbor validation execution state %s (exit=%v, error=%v)\nstdout:\n%s\nstderr:\n%s",
			execution.State, execution.ExitCode, execution.Error, execution.Stdout, execution.Stderr)
	}
	if execution.ExitCode == nil || *execution.ExitCode != 0 {
		t.Fatalf("Harbor validation exit code = %v\nstdout:\n%s\nstderr:\n%s",
			execution.ExitCode, execution.Stdout, execution.Stderr)
	}
	for _, marker := range []string{"scientific-stack-ok", "harbor-runtime-ok"} {
		if !strings.Contains(execution.Stdout, marker) {
			t.Fatalf("Harbor validation output missing %q:\n%s", marker, execution.Stdout)
		}
	}
	t.Logf("Harbor environment validation succeeded:\n%s", execution.Stdout)

	requireNoError(t, terminateSandbox(ctx, created.Sandbox.ID), "terminate Harbor sandbox")
	waitForSandboxState(t, ctx, created.Sandbox.ID, "TERMINATED")
	terminated = true
}

func waitForSandboxState(t *testing.T, ctx context.Context, sandboxID, wanted string) sandboxResult {
	t.Helper()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		var result struct {
			Sandbox *sandboxResult `json:"sandboxGet"`
		}
		err := sandboxGraphQL(ctx, `
			query GetSandbox($id: ID!, $workspace: String) {
				sandboxGet(id: $id, workspaceSlug: $workspace) {
					id
					state
					image
					error { code message }
				}
			}`,
			map[string]any{"id": sandboxID, "workspace": workspace}, &result)
		if err == nil && result.Sandbox != nil {
			t.Logf("sandbox %s state: %s", sandboxID, result.Sandbox.State)
			if result.Sandbox.State == wanted {
				return *result.Sandbox
			}
			if result.Sandbox.State == "FAILED" {
				t.Fatalf("sandbox %s failed: %+v", sandboxID, result.Sandbox.Error)
			}
		} else if err != nil {
			t.Logf("sandbox %s poll error: %v", sandboxID, err)
		}

		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for sandbox %s to reach %s: %v", sandboxID, wanted, ctx.Err())
		case <-ticker.C:
		}
	}
}

func waitForSandboxExecution(t *testing.T, ctx context.Context, sandboxID, executionID string) sandboxExecutionResult {
	t.Helper()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		var result struct {
			Execution *sandboxExecutionResult `json:"sandboxExecutionGet"`
		}
		err := sandboxGraphQL(ctx, `
			query GetSandboxExecution($sandboxId: ID!, $id: ID!, $workspace: String) {
				sandboxExecutionGet(sandboxId: $sandboxId, id: $id, workspaceSlug: $workspace) {
					id
					state
					stdout
					stderr
					exitCode
					error { code message }
				}
			}`,
			map[string]any{"sandboxId": sandboxID, "id": executionID, "workspace": workspace}, &result)
		if err == nil && result.Execution != nil {
			t.Logf("execution %s state: %s", executionID, result.Execution.State)
			switch result.Execution.State {
			case "SUCCEEDED", "FAILED", "CANCELLED", "TIMED_OUT":
				return *result.Execution
			}
		} else if err != nil {
			t.Logf("execution %s poll error: %v", executionID, err)
		}

		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for execution %s: %v", executionID, ctx.Err())
		case <-ticker.C:
		}
	}
}

func terminateSandbox(ctx context.Context, sandboxID string) error {
	var result struct {
		Sandbox sandboxResult `json:"sandboxTerminate"`
	}
	return sandboxGraphQL(ctx, `
		mutation TerminateSandbox($id: ID!, $workspace: String) {
			sandboxTerminate(id: $id, workspaceSlug: $workspace) {
				id
				state
				error { code message }
			}
		}`,
		map[string]any{"id": sandboxID, "workspace": workspace}, &result)
}

func sandboxGraphQL(ctx context.Context, query string, variables map[string]any, out any) error {
	requestBody, err := json.Marshal(map[string]any{"query": query, "variables": variables})
	if err != nil {
		return fmt.Errorf("marshal GraphQL request: %w", err)
	}

	apiURL := envOr("INK_API_URL", ink.DefaultBaseURL)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("create GraphQL request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+os.Getenv("INK_API_KEY"))
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("execute GraphQL request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if err != nil {
		return fmt.Errorf("read GraphQL response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("GraphQL HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var envelope struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message    string         `json:"message"`
			Extensions map[string]any `json:"extensions"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return fmt.Errorf("decode GraphQL response: %w", err)
	}
	if len(envelope.Errors) > 0 {
		return fmt.Errorf("GraphQL error: %s (%v)", envelope.Errors[0].Message, envelope.Errors[0].Extensions)
	}
	if err := json.Unmarshal(envelope.Data, out); err != nil {
		return fmt.Errorf("decode GraphQL data: %w", err)
	}
	return nil
}
