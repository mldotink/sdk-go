package ink

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	defaultSandboxPollInterval = 2 * time.Second
	maxSandboxFileSize         = 16 << 20
)

type SandboxExecutionError struct {
	Execution *SandboxExecution
}

func (e *SandboxExecutionError) Error() string {
	if e.Execution.Error != nil {
		return fmt.Sprintf(
			"ink: sandbox execution %s %s: %s",
			e.Execution.ID,
			strings.ToLower(e.Execution.State),
			e.Execution.Error.Message,
		)
	}
	if e.Execution.ExitCode != nil {
		return fmt.Sprintf(
			"ink: sandbox execution %s %s with exit code %d",
			e.Execution.ID,
			strings.ToLower(e.Execution.State),
			*e.Execution.ExitCode,
		)
	}
	return fmt.Sprintf("ink: sandbox execution %s %s", e.Execution.ID, strings.ToLower(e.Execution.State))
}

type SandboxDataPlaneError struct {
	StatusCode int
	Operation  string
}

func (e *SandboxDataPlaneError) Error() string {
	return fmt.Sprintf("ink: sandbox file %s failed with HTTP %d", e.Operation, e.StatusCode)
}

func (c *Client) CreateSandbox(ctx context.Context, options SandboxCreateOptions) (*Sandbox, error) {
	if strings.TrimSpace(options.Image) == "" {
		return nil, errors.New("ink: sandbox image is required")
	}
	if options.VCPUs <= 0 || options.MemoryGB <= 0 {
		return nil, errors.New("ink: sandbox vcpus and memoryGB must be greater than zero")
	}
	idempotencyKey := options.IdempotencyKey
	if idempotencyKey == "" {
		var err error
		idempotencyKey, err = sandboxIdempotencyKey("create")
		if err != nil {
			return nil, err
		}
	}
	input := CreateSandboxInput{
		IdempotencyKey: idempotencyKey,
		WorkspaceSlug:  optStr(options.WorkspaceSlug),
		Project:        optStr(options.Project),
		Name:           optStr(options.Name),
		Image:          options.Image,
		Resources: SandboxResourcesInput{
			Vcpus: options.VCPUs, MemoryGb: options.MemoryGB,
		},
		Env:                   sandboxEnvInput(options.Env),
		DestroyTimeoutSeconds: optInt(options.DestroyTimeoutSeconds),
		Metadata:              sandboxMetadataInput(options.Metadata),
	}
	if options.NetworkEgress != "" || len(options.AllowDomains) > 0 {
		input.Network = &SandboxNetworkInput{
			Egress:       sandboxEgressPointer(options.NetworkEgress),
			AllowDomains: options.AllowDomains,
		}
	}
	response, err := createSandboxAPI(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	sandbox := &response.SandboxCreate
	sandbox.attach(c, options.WorkspaceSlug)
	return sandbox.waitForState(ctx, SandboxStateRunning, options.PollInterval)
}

func (c *Client) GetSandbox(ctx context.Context, id, workspaceSlug string) (*Sandbox, error) {
	response, err := getSandboxAPI(ctx, c.gql, id, optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	if response.SandboxGet == nil {
		return nil, fmt.Errorf("ink: sandbox %q not found", id)
	}
	response.SandboxGet.attach(c, workspaceSlug)
	return response.SandboxGet, nil
}

func (sandbox *Sandbox) Refresh(ctx context.Context) error {
	if err := sandbox.ensureAttached(); err != nil {
		return err
	}
	latest, err := sandbox.client.GetSandbox(ctx, sandbox.ID, sandbox.workspaceSlug)
	if err != nil {
		return err
	}
	sandbox.update(latest)
	return nil
}

func (sandbox *Sandbox) Terminate(ctx context.Context) error {
	if err := sandbox.ensureAttached(); err != nil {
		return err
	}
	response, err := terminateSandboxAPI(ctx, sandbox.client.gql, sandbox.ID, optStr(sandbox.workspaceSlug))
	if err != nil {
		return err
	}
	terminated := response.SandboxTerminate
	terminated.attach(sandbox.client, sandbox.workspaceSlug)
	sandbox.update(&terminated)
	_, err = sandbox.waitForState(ctx, SandboxStateTerminated, defaultSandboxPollInterval)
	return err
}

func (sandbox *Sandbox) Run(
	ctx context.Context,
	argv []string,
	options SandboxRunOptions,
) (*SandboxExecution, error) {
	if err := sandbox.ensureAttached(); err != nil {
		return nil, err
	}
	if len(argv) == 0 {
		return nil, errors.New("ink: sandbox execution argv must not be empty")
	}
	idempotencyKey := options.IdempotencyKey
	if idempotencyKey == "" {
		var err error
		idempotencyKey, err = sandboxIdempotencyKey("exec")
		if err != nil {
			return nil, err
		}
	}
	response, err := createSandboxExecutionAPI(ctx, sandbox.client.gql, CreateSandboxExecutionInput{
		IdempotencyKey: idempotencyKey,
		SandboxId:      sandbox.ID,
		WorkspaceSlug:  optStr(sandbox.workspaceSlug),
		Argv:           append([]string(nil), argv...),
		Cwd:            optStr(options.CWD),
		Env:            sandboxEnvInput(options.Env),
		TimeoutSeconds: optInt(options.TimeoutSeconds),
		MaxOutputBytes: optInt(options.MaxOutputBytes),
	})
	if err != nil {
		return nil, err
	}
	execution := &response.SandboxExecutionCreate
	interval := options.PollInterval
	if interval <= 0 {
		interval = defaultSandboxPollInterval
	}
	return sandbox.waitForExecution(ctx, execution, interval)
}

func (sandbox *Sandbox) RunCommand(
	ctx context.Context,
	command string,
	options SandboxRunOptions,
) (*SandboxExecution, error) {
	if strings.TrimSpace(command) == "" {
		return nil, errors.New("ink: sandbox command must not be empty")
	}
	return sandbox.Run(ctx, []string{"/bin/sh", "-lc", command}, options)
}

func (sandbox *Sandbox) OpenPort(
	ctx context.Context,
	port int,
	options SandboxPortOptions,
) (*SandboxPortSession, error) {
	if err := sandbox.ensureAttached(); err != nil {
		return nil, err
	}
	idempotencyKey := options.IdempotencyKey
	if idempotencyKey == "" {
		var err error
		idempotencyKey, err = sandboxIdempotencyKey("port")
		if err != nil {
			return nil, err
		}
	}
	response, err := createSandboxPortSessionAPI(ctx, sandbox.client.gql, CreateSandboxPortSessionInput{
		IdempotencyKey:   idempotencyKey,
		SandboxId:        sandbox.ID,
		WorkspaceSlug:    optStr(sandbox.workspaceSlug),
		Port:             port,
		ExpiresInSeconds: optInt(options.ExpiresInSeconds),
	})
	if err != nil {
		return nil, err
	}
	return &response.SandboxPortSessionCreate, nil
}

func (sandbox *Sandbox) Metrics(
	ctx context.Context,
	timeRange MetricTimeRange,
	maxDataPoints int,
) (*SandboxMetrics, error) {
	if err := sandbox.ensureAttached(); err != nil {
		return nil, err
	}
	if timeRange == "" {
		timeRange = MetricTimeRangeOneHour
	}
	response, err := getSandboxMetricsAPI(
		ctx,
		sandbox.client.gql,
		sandbox.ID,
		timeRange,
		optInt(maxDataPoints),
		optStr(sandbox.workspaceSlug),
	)
	if err != nil {
		return nil, err
	}
	return &response.SandboxMetrics, nil
}

func (sandbox *Sandbox) StatFile(ctx context.Context, filePath string) (*SandboxFileInfo, error) {
	var result SandboxFileInfo
	if err := sandbox.doFileRequest(ctx, http.MethodGet, "stat", filePath, false, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (sandbox *Sandbox) ListFiles(ctx context.Context, filePath string) ([]SandboxFileInfo, error) {
	var result struct {
		Data []SandboxFileInfo `json:"data"`
	}
	if err := sandbox.doFileRequest(ctx, http.MethodGet, "list", filePath, false, nil, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (sandbox *Sandbox) ReadFile(ctx context.Context, filePath string) ([]byte, error) {
	var result []byte
	if err := sandbox.doFileRequest(ctx, http.MethodGet, "content", filePath, false, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (sandbox *Sandbox) WriteFile(
	ctx context.Context,
	filePath string,
	content []byte,
) (*SandboxFileInfo, error) {
	if len(content) > maxSandboxFileSize {
		return nil, errors.New("ink: sandbox file exceeds the 16 MiB limit")
	}
	var result SandboxFileInfo
	if err := sandbox.doFileRequest(
		ctx, http.MethodPut, "content", filePath, false, bytes.NewReader(content), &result,
	); err != nil {
		return nil, err
	}
	return &result, nil
}

func (sandbox *Sandbox) DeleteFile(ctx context.Context, filePath string, recursive bool) error {
	return sandbox.doFileRequest(ctx, http.MethodDelete, "content", filePath, recursive, nil, nil)
}

func (sandbox *Sandbox) waitForState(
	ctx context.Context,
	wanted string,
	pollInterval time.Duration,
) (*Sandbox, error) {
	if pollInterval <= 0 {
		pollInterval = defaultSandboxPollInterval
	}
	for {
		if sandbox.State == wanted {
			return sandbox, nil
		}
		if sandbox.State == SandboxStateFailed ||
			(sandbox.State == SandboxStateTerminated && wanted != SandboxStateTerminated) {
			if sandbox.Error != nil {
				return sandbox, fmt.Errorf(
					"ink: sandbox %s %s: %s",
					sandbox.ID,
					strings.ToLower(sandbox.State),
					sandbox.Error.Message,
				)
			}
			return sandbox, fmt.Errorf("ink: sandbox %s reached %s", sandbox.ID, sandbox.State)
		}
		if err := waitSandboxPoll(ctx, pollInterval); err != nil {
			return sandbox, err
		}
		if err := sandbox.Refresh(ctx); err != nil {
			return sandbox, err
		}
	}
}

func (sandbox *Sandbox) waitForExecution(
	ctx context.Context,
	execution *SandboxExecution,
	pollInterval time.Duration,
) (*SandboxExecution, error) {
	for {
		switch execution.State {
		case SandboxExecutionStateSucceeded:
			return execution, nil
		case SandboxExecutionStateFailed,
			SandboxExecutionStateCancelled,
			SandboxExecutionStateTimedOut:
			return execution, &SandboxExecutionError{Execution: execution}
		}
		if err := waitSandboxPoll(ctx, pollInterval); err != nil {
			return execution, err
		}
		response, err := getSandboxExecutionAPI(
			ctx,
			sandbox.client.gql,
			sandbox.ID,
			execution.ID,
			optStr(sandbox.workspaceSlug),
		)
		if err != nil {
			return execution, err
		}
		if response.SandboxExecutionGet == nil {
			return execution, fmt.Errorf("ink: sandbox execution %q not found", execution.ID)
		}
		execution = response.SandboxExecutionGet
	}
}

func (sandbox *Sandbox) doFileRequest(
	ctx context.Context,
	method, operation, filePath string,
	recursive bool,
	body io.Reader,
	output any,
) error {
	if err := sandbox.ensureAttached(); err != nil {
		return err
	}
	if !strings.HasPrefix(filePath, "/") {
		return errors.New("ink: sandbox file path must be absolute")
	}
	accessResponse, err := createSandboxFileAccessAPI(
		ctx,
		sandbox.client.gql,
		sandbox.ID,
		optStr(sandbox.workspaceSlug),
	)
	if err != nil {
		return err
	}
	access := accessResponse.SandboxFileAccessCreate
	endpoint, err := url.Parse(strings.TrimRight(access.BaseURL, "/") + "/" + operation)
	if err != nil {
		return fmt.Errorf("ink: invalid sandbox file endpoint: %w", err)
	}
	query := endpoint.Query()
	query.Set("path", filePath)
	if method == http.MethodDelete {
		query.Set("recursive", fmt.Sprintf("%t", recursive))
	}
	endpoint.RawQuery = query.Encode()
	request, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+access.Token)
	if method == http.MethodPut {
		request.Header.Set("Content-Type", "application/octet-stream")
	}
	response, err := sandbox.client.dataPlaneClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 64<<10))
		return &SandboxDataPlaneError{StatusCode: response.StatusCode, Operation: operation}
	}
	if output == nil {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 64<<10))
		return nil
	}
	if bytesOutput, ok := output.(*[]byte); ok {
		content, err := io.ReadAll(io.LimitReader(response.Body, maxSandboxFileSize+1))
		if err != nil {
			return err
		}
		if len(content) > maxSandboxFileSize {
			return errors.New("ink: sandbox file exceeds the 16 MiB limit")
		}
		*bytesOutput = content
		return nil
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 4<<20)).Decode(output); err != nil {
		return fmt.Errorf("ink: decode sandbox file response: %w", err)
	}
	return nil
}

func (sandbox *Sandbox) attach(client *Client, workspaceSlug string) {
	sandbox.client = client
	sandbox.workspaceSlug = workspaceSlug
}

func (sandbox *Sandbox) update(latest *Sandbox) {
	client, workspaceSlug := sandbox.client, sandbox.workspaceSlug
	*sandbox = *latest
	sandbox.attach(client, workspaceSlug)
}

func (sandbox *Sandbox) ensureAttached() error {
	if sandbox == nil || sandbox.client == nil {
		return errors.New("ink: sandbox is not attached to a client")
	}
	return nil
}

func sandboxIdempotencyKey(operation string) (string, error) {
	random := make([]byte, 18)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("ink: create sandbox %s idempotency key: %w", operation, err)
	}
	return "sdk-" + operation + "-" + base64.RawURLEncoding.EncodeToString(random), nil
}

func sandboxEgressPointer(value SandboxNetworkEgress) *SandboxNetworkEgress {
	if value == "" {
		return nil
	}
	return &value
}

func sandboxEnvInput(values map[string]string) []EnvVarInput {
	keys := sortedSandboxKeys(values)
	result := make([]EnvVarInput, 0, len(keys))
	for _, key := range keys {
		result = append(result, EnvVarInput{Key: key, Value: values[key]})
	}
	return result
}

func sandboxMetadataInput(values map[string]string) []SandboxMetadataInput {
	keys := sortedSandboxKeys(values)
	result := make([]SandboxMetadataInput, 0, len(keys))
	for _, key := range keys {
		result = append(result, SandboxMetadataInput{Key: key, Value: values[key]})
	}
	return result
}

func sortedSandboxKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func waitSandboxPoll(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
