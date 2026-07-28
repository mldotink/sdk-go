package ink

import "time"

const (
	SandboxStatePending     = "PENDING"
	SandboxStateStarting    = "STARTING"
	SandboxStateRunning     = "RUNNING"
	SandboxStateTerminating = "TERMINATING"
	SandboxStateTerminated  = "TERMINATED"
	SandboxStateFailed      = "FAILED"

	SandboxExecutionStateQueued    = "QUEUED"
	SandboxExecutionStateRunning   = "RUNNING"
	SandboxExecutionStateSucceeded = "SUCCEEDED"
	SandboxExecutionStateFailed    = "FAILED"
	SandboxExecutionStateCancelled = "CANCELLED"
	SandboxExecutionStateTimedOut  = "TIMED_OUT"
)

type SandboxResources struct {
	VCPUs    float64 `json:"vcpus"`
	MemoryGB float64 `json:"memoryGb"`
}

type SandboxNetwork struct {
	Egress       string   `json:"egress"`
	AllowDomains []string `json:"allowDomains"`
	Inbound      string   `json:"inbound"`
}

type SandboxStorage struct {
	Root string `json:"root"`
}

type SandboxError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SandboxMetadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Sandbox struct {
	ID                    string            `json:"id"`
	WorkspaceID           string            `json:"workspaceId"`
	ProjectID             string            `json:"projectId"`
	Name                  *string           `json:"name"`
	State                 string            `json:"state"`
	Image                 string            `json:"image"`
	ImageDigest           *string           `json:"imageDigest"`
	Resources             SandboxResources  `json:"resources"`
	Network               SandboxNetwork    `json:"network"`
	Storage               SandboxStorage    `json:"storage"`
	Capabilities          []string          `json:"capabilities"`
	EnvKeys               []string          `json:"envKeys"`
	Metadata              []SandboxMetadata `json:"metadata"`
	CreatedAt             string            `json:"createdAt"`
	ReadyAt               *string           `json:"readyAt"`
	DestroyTimeoutSeconds int               `json:"destroyTimeoutSeconds"`
	DestroyAt             *string           `json:"destroyAt"`
	TerminatedAt          *string           `json:"terminatedAt"`
	TerminationReason     *string           `json:"terminationReason"`
	Error                 *SandboxError     `json:"error"`

	client        *Client
	workspaceSlug string
}

type SandboxExecution struct {
	ID              string        `json:"id"`
	SandboxID       string        `json:"sandboxId"`
	State           string        `json:"state"`
	Argv            []string      `json:"argv"`
	CWD             *string       `json:"cwd"`
	EnvKeys         []string      `json:"envKeys"`
	Stdout          string        `json:"stdout"`
	Stderr          string        `json:"stderr"`
	OutputTruncated bool          `json:"outputTruncated"`
	OutputExpired   bool          `json:"outputExpired"`
	ExitCode        *int          `json:"exitCode"`
	CreatedAt       string        `json:"createdAt"`
	StartedAt       *string       `json:"startedAt"`
	FinishedAt      *string       `json:"finishedAt"`
	Error           *SandboxError `json:"error"`
}

type SandboxFileAccess struct {
	BaseURL   string `json:"baseUrl"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

type SandboxPortSession struct {
	ID        string `json:"id"`
	SandboxID string `json:"sandboxId"`
	Port      int    `json:"port"`
	URL       string `json:"url"`
	ExpiresAt string `json:"expiresAt"`
}

type SandboxMetrics struct {
	CPUUsage                   MetricSeries `json:"cpuUsage"`
	MemoryUsageMB              MetricSeries `json:"memoryUsageMB"`
	NetworkReceiveBytesPerSec  MetricSeries `json:"networkReceiveBytesPerSec"`
	NetworkTransmitBytesPerSec MetricSeries `json:"networkTransmitBytesPerSec"`
	MemoryLimitMB              float64      `json:"memoryLimitMB"`
	CPULimitVCPUs              float64      `json:"cpuLimitVCPUs"`
}

type SandboxCreateOptions struct {
	IdempotencyKey        string
	Name                  string
	WorkspaceSlug         string
	Project               string
	Image                 string
	VCPUs                 float64
	MemoryGB              float64
	Env                   map[string]string
	NetworkEgress         SandboxNetworkEgress
	AllowDomains          []string
	DestroyTimeoutSeconds int
	Metadata              map[string]string
	PollInterval          time.Duration
}

type SandboxRunOptions struct {
	IdempotencyKey string
	CWD            string
	Env            map[string]string
	TimeoutSeconds int
	MaxOutputBytes int
	PollInterval   time.Duration
}

type SandboxPortOptions struct {
	IdempotencyKey   string
	ExpiresInSeconds int
}

type SandboxFileInfo struct {
	Path          string  `json:"path"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Size          int64   `json:"size"`
	Mode          int64   `json:"mode"`
	ModifiedAt    string  `json:"modified_at"`
	SymlinkTarget *string `json:"symlink_target,omitempty"`
}
