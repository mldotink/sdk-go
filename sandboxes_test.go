package ink

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSandboxSDKHidesLifecycleAndDataPlanePlumbing(t *testing.T) {
	var mutex sync.Mutex
	graphQLOperations := make([]string, 0, 12)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/graphql" {
			if request.Header.Get("Authorization") != "Bearer capability-secret" {
				t.Errorf("data-plane authorization = %q", request.Header.Get("Authorization"))
			}
			switch request.Method {
			case http.MethodGet:
				_, _ = writer.Write([]byte("sandbox-content"))
			case http.MethodPut:
				_, _ = io.Copy(io.Discard, request.Body)
				writeTestJSON(writer, map[string]any{
					"path": "/workspace/result.txt", "name": "result.txt", "type": "file", "size": 7,
				})
			default:
				writer.WriteHeader(http.StatusNoContent)
			}
			return
		}
		if request.Header.Get("Authorization") != "Bearer api-secret" {
			t.Errorf("GraphQL authorization = %q", request.Header.Get("Authorization"))
		}
		var envelope struct {
			OperationName string         `json:"operationName"`
			Variables     map[string]any `json:"variables"`
		}
		if err := json.NewDecoder(request.Body).Decode(&envelope); err != nil {
			t.Errorf("decode GraphQL request: %v", err)
			return
		}
		mutex.Lock()
		graphQLOperations = append(graphQLOperations, envelope.OperationName)
		mutex.Unlock()
		switch envelope.OperationName {
		case "createSandboxAPI":
			input := envelope.Variables["input"].(map[string]any)
			network := input["network"].(map[string]any)
			if input["idempotencyKey"] != "create-key-0001" ||
				network["allowDomains"].([]any)[0] != "api.openai.com" {
				t.Errorf("create input = %#v", input)
			}
			writeGraphQLData(writer, "sandboxCreate", testSandboxJSON("PENDING"))
		case "getSandboxAPI":
			writeGraphQLData(writer, "sandboxGet", testSandboxJSON("RUNNING"))
		case "createSandboxExecutionAPI":
			writeGraphQLData(writer, "sandboxExecutionCreate", map[string]any{
				"id": "exe_test", "sandboxId": "sbx_test", "state": "QUEUED", "argv": []string{"printf", "ok"},
			})
		case "getSandboxExecutionAPI":
			exitCode := 0
			writeGraphQLData(writer, "sandboxExecutionGet", map[string]any{
				"id": "exe_test", "sandboxId": "sbx_test", "state": "SUCCEEDED",
				"stdout": "ok", "exitCode": exitCode,
			})
		case "createSandboxFileAccessAPI":
			writeGraphQLData(writer, "sandboxFileAccessCreate", map[string]any{
				"baseUrl": server.URL + "/sandbox-files/sbx_test",
				"token":   "capability-secret", "expiresAt": "2026-07-28T12:00:00Z",
			})
		case "createSandboxPortSessionAPI":
			writeGraphQLData(writer, "sandboxPortSessionCreate", map[string]any{
				"id": "psn_test", "sandboxId": "sbx_test", "port": 8080,
				"url": "https://sandbox.example.test", "expiresAt": "2026-07-28T12:00:00Z",
			})
		case "getSandboxMetricsAPI":
			if envelope.Variables["timeRange"] != "ONE_HOUR" {
				t.Errorf("metrics time range = %#v", envelope.Variables["timeRange"])
			}
			writeGraphQLData(writer, "sandboxMetrics", testSandboxMetricsJSON())
		case "terminateSandboxAPI":
			writeGraphQLData(writer, "sandboxTerminate", testSandboxJSON("TERMINATED"))
		default:
			t.Errorf("unexpected GraphQL operation %q", envelope.OperationName)
		}
	}))
	defer server.Close()

	client := NewClient(Config{
		APIKey: "api-secret", BaseURL: server.URL + "/graphql", HTTPClient: server.Client(),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sandbox, err := client.CreateSandbox(ctx, SandboxCreateOptions{
		IdempotencyKey: "create-key-0001",
		WorkspaceSlug:  "research",
		Image:          "python:3.12",
		VCPUs:          1.5,
		MemoryGB:       2,
		AllowDomains:   []string{"api.openai.com"},
		PollInterval:   time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if sandbox.State != SandboxStateRunning || sandbox.Network.AllowDomains[0] != "api.openai.com" {
		t.Fatalf("sandbox = %#v", sandbox)
	}

	execution, err := sandbox.Run(ctx, []string{"printf", "ok"}, SandboxRunOptions{
		IdempotencyKey: "execute-key-0001", PollInterval: time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if execution.Stdout != "ok" || execution.ExitCode == nil || *execution.ExitCode != 0 {
		t.Fatalf("execution = %#v", execution)
	}

	if _, err := sandbox.WriteFile(ctx, "/workspace/result.txt", []byte("content")); err != nil {
		t.Fatal(err)
	}
	content, err := sandbox.ReadFile(ctx, "/workspace/result.txt")
	if err != nil || string(content) != "sandbox-content" {
		t.Fatalf("file content = %q, err=%v", content, err)
	}
	portSession, err := sandbox.OpenPort(ctx, 8080, SandboxPortOptions{IdempotencyKey: "port-key-0001"})
	if err != nil || !strings.HasPrefix(portSession.URL, "https://") {
		t.Fatalf("port session = %#v, err=%v", portSession, err)
	}
	metrics, err := sandbox.Metrics(ctx, "", 100)
	if err != nil || metrics.CPULimitVCPUs != 1.5 {
		t.Fatalf("metrics = %#v, err=%v", metrics, err)
	}
	if err := sandbox.Terminate(ctx); err != nil {
		t.Fatal(err)
	}
	if sandbox.State != SandboxStateTerminated {
		t.Fatalf("terminated sandbox state = %q", sandbox.State)
	}

	mutex.Lock()
	defer mutex.Unlock()
	for _, wanted := range []string{
		"createSandboxAPI",
		"getSandboxAPI",
		"createSandboxExecutionAPI",
		"getSandboxExecutionAPI",
		"createSandboxFileAccessAPI",
		"createSandboxPortSessionAPI",
		"getSandboxMetricsAPI",
		"terminateSandboxAPI",
	} {
		if !containsString(graphQLOperations, wanted) {
			t.Fatalf("operation %q was not called: %#v", wanted, graphQLOperations)
		}
	}
}

func writeGraphQLData(writer http.ResponseWriter, field string, value any) {
	writeTestJSON(writer, map[string]any{"data": map[string]any{field: value}})
}

func writeTestJSON(writer http.ResponseWriter, value any) {
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(value)
}

func testSandboxJSON(state string) map[string]any {
	return map[string]any{
		"id": "sbx_test", "workspaceId": "ws_test", "projectId": "prj_test",
		"state": state, "image": "python:3.12",
		"resources": map[string]any{"vcpus": 1.5, "memoryGb": 2},
		"network": map[string]any{
			"egress": "ALLOWLIST", "allowDomains": []string{"api.openai.com"}, "inbound": "private_sessions",
		},
		"storage":               map[string]any{"root": "ephemeral"},
		"capabilities":          []string{"commands", "files", "private_ports"},
		"envKeys":               []string{},
		"metadata":              []any{},
		"destroyTimeoutSeconds": 0,
		"createdAt":             "2026-07-28T11:00:00Z",
	}
}

func testSandboxMetricsJSON() map[string]any {
	series := func(metric string) map[string]any {
		return map[string]any{
			"metric": metric,
			"dataPoints": []any{
				map[string]any{"timestamp": "2026-07-28T11:00:00Z", "value": 1},
			},
		}
	}
	return map[string]any{
		"cpuUsage": series("cpu_usage"), "memoryUsageMB": series("memory_usage_mb"),
		"networkReceiveBytesPerSec":  series("network_receive_bytes_per_sec"),
		"networkTransmitBytesPerSec": series("network_transmit_bytes_per_sec"),
		"memoryLimitMB":              2048, "cpuLimitVCPUs": 1.5,
	}
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
