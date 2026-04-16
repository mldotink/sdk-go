package ink

// ── Services ─────────────────────────────────────────────────────────────────

// Service represents a deployed service on Ink.
type Service struct {
	ID                    string        `json:"id"`
	ProjectID             string        `json:"projectId"`
	Name                  string        `json:"name"`
	Subdomain             string        `json:"subdomain"`
	Source                string        `json:"source"`
	Repo                  string        `json:"repo"`
	Image                 string        `json:"image"`
	Branch                string        `json:"branch"`
	Status                string        `json:"status"`
	ErrorMessage          string        `json:"errorMessage"`
	EnvVars               []EnvVar      `json:"envVars"`
	Ports                 []ServicePort `json:"ports"`
	GitProvider           string        `json:"gitProvider"`
	CommitHash            string        `json:"commitHash"`
	Memory                string        `json:"memory"`
	VCPUs                 string        `json:"vcpus"`
	CustomDomain          string        `json:"customDomain"`
	CustomDomainStatus    string        `json:"customDomainStatus"`
	BuildPack             string        `json:"buildPack"`
	BuildCommand          string        `json:"buildCommand"`
	StartCommand          string        `json:"startCommand"`
	PublishDirectory      string        `json:"publishDirectory"`
	RootDirectory         string        `json:"rootDirectory"`
	DockerfilePath        string        `json:"dockerfilePath"`
	TeardownEnabled       bool          `json:"teardownEnabled"`
	TeardownOverlapSeconds  *int        `json:"teardownOverlapSeconds"`
	TeardownDrainingSeconds *int        `json:"teardownDrainingSeconds"`
	DestroyTimeoutSeconds int           `json:"destroyTimeoutSeconds"`
	CreatedAt             string        `json:"createdAt"`
	UpdatedAt             string        `json:"updatedAt"`
}

// EnvVar is a key-value environment variable on a service.
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ServicePort describes a network port exposed by a service.
type ServicePort struct {
	Name             string `json:"name"`
	Port             string `json:"port"`
	Protocol         string `json:"protocol"`
	Visibility       string `json:"visibility"`
	InternalEndpoint string `json:"internalEndpoint"`
	PublicEndpoint   string `json:"publicEndpoint"`
}

// VolumeSpec defines a persistent volume to attach to a service.
// Corresponds to the GraphQL VolumeInput type.
type VolumeSpec struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	SizeGi    int    `json:"sizeGi,omitempty"`
}

// ServicePortInput defines a port for service creation/update.
type ServicePortInput struct {
	Name       string `json:"name"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	Visibility string `json:"visibility"`
}

// BucketMountInput configures a storage bucket mount on a service.
type BucketMountInput struct {
	Name         string `json:"name"`
	MountPath    string `json:"mountPath,omitempty"`
	Mode         string `json:"mode,omitempty"`
	Prefix       string `json:"prefix,omitempty"`
	SyncInterval int    `json:"syncInterval,omitempty"`
}

// CreateServiceInput defines the parameters for creating a new service.
type CreateServiceInput struct {
	Name                  string             `json:"name,omitempty"`
	Subdomain             string             `json:"subdomain,omitempty"`
	Source                string             `json:"source,omitempty"`
	Repo                  string             `json:"repo,omitempty"`
	Image                 string             `json:"image,omitempty"`
	Host                  string             `json:"host,omitempty"`
	Branch                string             `json:"branch,omitempty"`
	Project               string             `json:"project,omitempty"`
	WorkspaceSlug         string             `json:"workspaceSlug,omitempty"`
	BuildPack             string             `json:"buildPack,omitempty"`
	Ports                 []ServicePortInput `json:"ports,omitempty"`
	EnvVars               []EnvVar           `json:"envVars,omitempty"`
	Memory                string             `json:"memory,omitempty"`
	VCPUs                 string             `json:"vcpus,omitempty"`
	BuildCommand          string             `json:"buildCommand,omitempty"`
	StartCommand          string             `json:"startCommand,omitempty"`
	PublishDirectory      string             `json:"publishDirectory,omitempty"`
	RootDirectory         string             `json:"rootDirectory,omitempty"`
	DockerfilePath        string             `json:"dockerfilePath,omitempty"`
	Regions               []string           `json:"regions,omitempty"`
	Volumes               []VolumeSpec       `json:"volumes,omitempty"`
	Bucket                *BucketMountInput  `json:"bucket,omitempty"`
	DestroyTimeoutSeconds int                `json:"destroyTimeoutSeconds,omitempty"`
	TeardownEnabled         bool `json:"teardownEnabled,omitempty"`
	TeardownOverlapSeconds  *int `json:"teardownOverlapSeconds,omitempty"`
	TeardownDrainingSeconds *int `json:"teardownDrainingSeconds,omitempty"`
}

// UpdateServiceInput defines parameters for updating an existing service.
// Only set fields are sent; pointer fields let you explicitly clear a value.
type UpdateServiceInput struct {
	Name                  string             `json:"name,omitempty"`
	ServiceID             string             `json:"serviceId,omitempty"`
	Project               string             `json:"project,omitempty"`
	ProjectID             string             `json:"projectId,omitempty"`
	WorkspaceSlug         string             `json:"workspaceSlug,omitempty"`
	Source                *string            `json:"source,omitempty"`
	Image                 *string            `json:"image,omitempty"`
	Repo                  *string            `json:"repo,omitempty"`
	Host                  *string            `json:"host,omitempty"`
	Branch                *string            `json:"branch,omitempty"`
	BuildPack             *string            `json:"buildPack,omitempty"`
	Memory                *string            `json:"memory,omitempty"`
	VCPUs                 *string            `json:"vcpus,omitempty"`
	Ports                 []ServicePortInput `json:"ports,omitempty"`
	EnvVars               []EnvVar           `json:"envVars,omitempty"`
	BuildCommand          *string            `json:"buildCommand,omitempty"`
	StartCommand          *string            `json:"startCommand,omitempty"`
	PublishDirectory      *string            `json:"publishDirectory,omitempty"`
	RootDirectory         *string            `json:"rootDirectory,omitempty"`
	DockerfilePath        *string            `json:"dockerfilePath,omitempty"`
	Volumes               []VolumeSpec       `json:"volumes,omitempty"`
	Bucket                *BucketMountInput  `json:"bucket,omitempty"`
	DestroyTimeoutSeconds *int               `json:"destroyTimeoutSeconds,omitempty"`
	TeardownEnabled       *bool              `json:"teardownEnabled,omitempty"`
	TeardownOverlapSeconds  *int             `json:"teardownOverlapSeconds,omitempty"`
	TeardownDrainingSeconds *int             `json:"teardownDrainingSeconds,omitempty"`
}

// CreateServiceResult is the result of creating a service.
type CreateServiceResult struct {
	ServiceID string        `json:"serviceId"`
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Repo      string        `json:"repo"`
	Ports     []ServicePort `json:"ports"`
}

// UpdateServiceResult is the result of updating a service.
type UpdateServiceResult struct {
	ServiceID string `json:"serviceId"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}

// DeleteServiceInput identifies a service to delete.
type DeleteServiceInput struct {
	Name          string `json:"name,omitempty"`
	ServiceID     string `json:"serviceId,omitempty"`
	Project       string `json:"project,omitempty"`
	ProjectID     string `json:"projectId,omitempty"`
	WorkspaceSlug string `json:"workspaceSlug,omitempty"`
}

// DeleteServiceResult is the result of deleting a service.
type DeleteServiceResult struct {
	ServiceID string `json:"serviceId"`
	Name      string `json:"name"`
	Message   string `json:"message"`
}

// SetSecretsResult is the result of setting or deleting env vars.
type SetSecretsResult struct {
	ServiceID string `json:"serviceId"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}

// SetSecretsInput defines env vars to set on a service.
type SetSecretsInput struct {
	Name          string   `json:"name,omitempty"`
	ServiceID     string   `json:"serviceId,omitempty"`
	Project       string   `json:"project,omitempty"`
	ProjectID     string   `json:"projectId,omitempty"`
	WorkspaceSlug string   `json:"workspaceSlug,omitempty"`
	EnvVars       []EnvVar `json:"envVars"`
	Replace       bool     `json:"replace,omitempty"`
}

// DeleteSecretsInput defines env var keys to remove from a service.
type DeleteSecretsInput struct {
	Name          string   `json:"name,omitempty"`
	ServiceID     string   `json:"serviceId,omitempty"`
	Project       string   `json:"project,omitempty"`
	ProjectID     string   `json:"projectId,omitempty"`
	WorkspaceSlug string   `json:"workspaceSlug,omitempty"`
	Keys          []string `json:"keys"`
}

// ExecSession holds connection details for an interactive shell session.
type ExecSession struct {
	URL       string `json:"url"`
	Token     string `json:"token"`
	ServiceID string `json:"serviceId"`
}

// ExecResult is the output of a one-shot command execution.
type ExecResult struct {
	ExitCode int    `json:"exitCode"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

// ExecInput identifies a service for command execution.
// Provide either ServiceID or Name (+ optional Project/WorkspaceSlug).
type ExecInput struct {
	ServiceID     string
	Name          string
	Project       string
	WorkspaceSlug string
}

// ── Workspaces ────────────────────────────────────────────────────────────────

// Workspace is a team/organisation on Ink.
type Workspace struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsDefault bool   `json:"isDefault"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

// WorkspaceMember is a member of a workspace.
type WorkspaceMember struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
	Role        string `json:"role"`
	JoinedAt    string `json:"joinedAt"`
}

// WorkspaceInvite is a workspace invitation.
type WorkspaceInvite struct {
	ID                   string `json:"id"`
	WorkspaceID          string `json:"workspaceId"`
	WorkspaceName        string `json:"workspaceName"`
	WorkspaceSlug        string `json:"workspaceSlug"`
	InviterDisplayName   string `json:"inviterDisplayName"`
	InviteeDisplayName   string `json:"inviteeDisplayName"`
	Role                 string `json:"role"`
	Status               string `json:"status"`
	CreatedAt            string `json:"createdAt"`
}

// ── Projects ─────────────────────────────────────────────────────────────────

// Project is a logical grouping of services within a workspace.
type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CreateProjectInput defines the parameters for creating a project.
type CreateProjectInput struct {
	Name          string `json:"name"`
	WorkspaceSlug string `json:"workspaceSlug,omitempty"`
}

// ── DNS ──────────────────────────────────────────────────────────────────────

// DNSZone is a DNS zone managed by Ink. Corresponds to GraphQL HostedZone.
type DNSZone struct {
	ID        string `json:"id"`
	Zone      string `json:"zone"`
	Status    string `json:"status"`
	Error     string `json:"error"`
	CreatedAt string `json:"createdAt"`
}

// ZoneRecord is a DNS record within a hosted zone. Corresponds to GraphQL ZoneRecord.
type ZoneRecord struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	TTL       int    `json:"ttl"`
	Managed   bool   `json:"managed"`
	CreatedAt string `json:"createdAt"`
}

// ── Domains ──────────────────────────────────────────────────────────────────

// AddDomainResult is the result of attaching a custom domain. Corresponds to GraphQL DomainAddResult.
type AddDomainResult struct {
	ServiceID string `json:"serviceId"`
	Domain    string `json:"domain"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// RemoveDomainResult is the result of detaching a custom domain. Corresponds to GraphQL DomainRemoveResult.
type RemoveDomainResult struct {
	ServiceID string `json:"serviceId"`
	Message   string `json:"message"`
}

// ── Templates ────────────────────────────────────────────────────────────────

// Template is a reusable deployment blueprint. Corresponds to GraphQL ServiceTemplate.
type Template struct {
	Slug        string             `json:"slug"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Tags        []string           `json:"tags"`
	Icon        string             `json:"icon"`
	Variables   []TemplateVariable `json:"variables"`
	Services    []TemplateService  `json:"services"`
	Outputs     []TemplateOutput   `json:"outputs"`
}

// TemplateVariable is a user-configurable input for a template. Corresponds to GraphQL TemplateVariableDef.
type TemplateVariable struct {
	Key          string   `json:"key"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Required     bool     `json:"required"`
	Sensitive    bool     `json:"sensitive"`
	DefaultValue string   `json:"defaultValue"`
	Options      []string `json:"options"`
}

// TemplateService describes a service blueprint within a template. Corresponds to GraphQL TemplateServiceDef.
type TemplateService struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Source string `json:"source"`
	Image  string `json:"image"`
	Memory string `json:"memory"`
	VCPUs  string `json:"vcpus"`
}

// TemplateOutput is a metadata descriptor for a template output. Corresponds to GraphQL TemplateOutputDef.
type TemplateOutput struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Sensitive   bool   `json:"sensitive"`
}

// TemplateVariableValue is a key-value pair used when deploying a template.
type TemplateVariableValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TemplateDeployInput defines the parameters for deploying a template.
type TemplateDeployInput struct {
	// Template is the slug of the template to deploy.
	Template      string                  `json:"template"`
	Name          string                  `json:"name"`
	WorkspaceSlug string                  `json:"workspaceSlug,omitempty"`
	Project       string                  `json:"project,omitempty"`
	Variables     []TemplateVariableValue `json:"variables,omitempty"`
}

// TemplateDeployResult is the result of deploying a template.
type TemplateDeployResult struct {
	TemplateInstanceID string                    `json:"templateInstanceId"`
	ProjectID          string                    `json:"projectId"`
	Services           []TemplateDeployedService `json:"services"`
	Outputs            []TemplateDeployedOutput  `json:"outputs"`
}

// TemplateDeployedService is a service created during a template deployment.
type TemplateDeployedService struct {
	ServiceID string                    `json:"serviceId"`
	Key       string                    `json:"key"`
	Name      string                    `json:"name"`
	Status    string                    `json:"status"`
	Endpoints []TemplateServiceEndpoint `json:"endpoints"`
}

// TemplateServiceEndpoint is an endpoint exposed by a deployed template service.
type TemplateServiceEndpoint struct {
	Name             string `json:"name"`
	Port             string `json:"port"`
	Protocol         string `json:"protocol"`
	Visibility       string `json:"visibility"`
	InternalEndpoint string `json:"internalEndpoint"`
	PublicEndpoint   string `json:"publicEndpoint"`
}

// TemplateDeployedOutput is a resolved output value after a template deployment.
type TemplateDeployedOutput struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Sensitive   bool   `json:"sensitive"`
	Value       string `json:"value"`
}

// TemplateInstance is a deployed instance of a template.
type TemplateInstance struct {
	ID           string                    `json:"id"`
	TemplateSlug string                    `json:"templateSlug"`
	ProjectID    string                    `json:"projectId"`
	Name         string                    `json:"name"`
	Status       string                    `json:"status"`
	Services     []TemplateDeployedService `json:"services"`
	Outputs      []TemplateDeployedOutput  `json:"outputs"`
	CreatedAt    string                    `json:"createdAt"`
}

// ── Account ──────────────────────────────────────────────────────────────────

// AccountStatus is the current authenticated user's account details.
// Corresponds to the GraphQL User type.
type AccountStatus struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	DisplayName      string `json:"displayName"`
	Username         string `json:"username"`
	GitHubUsername   string `json:"githubUsername"`
	HasGitHubOAuth   bool   `json:"hasGitHubOAuth"`
	HasGitHubApp     bool   `json:"hasGitHubApp"`
	DefaultWorkspace string `json:"defaultWorkspace"`
	SubscriptionTier string `json:"subscriptionTier"`
}

// ── Billing ──────────────────────────────────────────────────────────────────

// UsageBillBreakdown is a line-item cost breakdown for the current billing period.
type UsageBillBreakdown struct {
	Memory             UsageLineItem `json:"memory"`
	CPU                UsageLineItem `json:"cpu"`
	Egress             UsageLineItem `json:"egress"`
	Subtotal           string        `json:"subtotal"`
	IncludedUsageCents int           `json:"includedUsageCents"`
	PlanFeeCents       int           `json:"planFeeCents"`
	CurrentBillCents   int           `json:"currentBillCents"`
	PeriodStart        string        `json:"periodStart"`
	PeriodEnd          string        `json:"periodEnd"`
}

// UsageLineItem is a single cost component in the bill breakdown.
type UsageLineItem struct {
	Quantity   string `json:"quantity"`
	UnitPrice  string `json:"unitPrice"`
	Unit       string `json:"unit"`
	TotalCents int    `json:"totalCents"`
}

// ── Logs ─────────────────────────────────────────────────────────────────────

// LogType specifies which log stream to query.
type LogType = string

const (
	LogTypeBuild   LogType = "BUILD"
	LogTypeRuntime LogType = "RUNTIME"
)

// LogsInput defines the parameters for querying service logs.
type LogsInput struct {
	ServiceID string  `json:"serviceId"`
	LogType   LogType `json:"logType"`
	StartTime string  `json:"startTime,omitempty"`
	EndTime   string  `json:"endTime,omitempty"`
	Query     string  `json:"query,omitempty"`
	Limit     int     `json:"limit,omitempty"`
}

// LogEntry is a single log line from a service.
type LogEntry struct {
	Timestamp  string `json:"timestamp"`
	Level      string `json:"level"`
	Message    string `json:"message"`
	Attributes string `json:"attributes"`
}

// LogsResult is the result of a logs query.
type LogsResult struct {
	Entries []LogEntry `json:"entries"`
	HasMore bool       `json:"hasMore"`
}

// ── Metrics ──────────────────────────────────────────────────────────────────

// MetricDataPoint is a single point in a metric time series.
type MetricDataPoint struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

// MetricSeries is a named time series.
type MetricSeries struct {
	Metric     string            `json:"metric"`
	DataPoints []MetricDataPoint `json:"dataPoints"`
}

// ServiceMetrics contains CPU, memory, and network metrics for a service.
type ServiceMetrics struct {
	CPUUsage                   MetricSeries `json:"cpuUsage"`
	MemoryUsageMB              MetricSeries `json:"memoryUsageMB"`
	NetworkReceiveBytesPerSec  MetricSeries `json:"networkReceiveBytesPerSec"`
	NetworkTransmitBytesPerSec MetricSeries `json:"networkTransmitBytesPerSec"`
	MemoryLimitMB              float64      `json:"memoryLimitMB"`
	CPULimitVCPUs              float64      `json:"cpuLimitVCPUs"`
	DiskUsageMB                MetricSeries `json:"diskUsageMB"`
	VolumeSizeGi               int          `json:"volumeSizeGi"`
}

// ── Repos ─────────────────────────────────────────────────────────────────────

// CreateRepoInput defines the parameters for creating an internal git repo.
// Corresponds to GraphQL RepoCreateInput.
type CreateRepoInput struct {
	Name          string `json:"name"`
	Host          string `json:"host,omitempty"`
	Description   string `json:"description,omitempty"`
	Project       string `json:"project,omitempty"`
	WorkspaceSlug string `json:"workspaceSlug,omitempty"`
}

// CreateRepoResult is the result of creating a git repo.
// Corresponds to GraphQL RepoCreateResult.
type CreateRepoResult struct {
	Name      string `json:"name"`
	GitRemote string `json:"gitRemote"`
	ExpiresAt string `json:"expiresAt"`
	Message   string `json:"message"`
}

// GetRepoTokenInput identifies a repo to get an access token for.
// Corresponds to GraphQL RepoGetTokenInput.
type GetRepoTokenInput struct {
	Name          string `json:"name"`
	Host          string `json:"host,omitempty"`
	WorkspaceSlug string `json:"workspaceSlug,omitempty"`
}

// GetRepoTokenResult is the result of getting a short-lived push token.
// Corresponds to GraphQL RepoGetTokenResult.
type GetRepoTokenResult struct {
	GitRemote string `json:"gitRemote"`
	ExpiresAt string `json:"expiresAt"`
}

// ── Chat ─────────────────────────────────────────────────────────────────────

// ChatMessage is a message in a workspace chat channel.
type ChatMessage struct {
	Seq        int    `json:"seq"`
	MessageID  string `json:"messageId"`
	SenderID   string `json:"senderId"`
	SenderName string `json:"senderName"`
	Channel    string `json:"channel"`
	Content    string `json:"content"`
	Metadata   string `json:"metadata"`
	CreatedAt  string `json:"createdAt"`
}

// ReadChatResult is the result of reading chat messages.
// Corresponds to GraphQL ChatReadResult.
type ReadChatResult struct {
	Messages   []ChatMessage `json:"messages"`
	NextCursor int           `json:"nextCursor"`
	HasMore    bool          `json:"hasMore"`
}

// SendChatResult is the result of sending a chat message.
// Corresponds to GraphQL ChatSendResult.
type SendChatResult struct {
	Seq       int    `json:"seq"`
	MessageID string `json:"messageId"`
}
