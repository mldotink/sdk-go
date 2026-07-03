package ink

type Service struct {
	ID                      string        `json:"id"`
	ProjectID               string        `json:"projectId"`
	Name                    string        `json:"name"`
	Subdomain               string        `json:"subdomain"`
	Source                  string        `json:"source"`
	Repo                    string        `json:"repo"`
	Image                   string        `json:"image"`
	Branch                  string        `json:"branch"`
	Status                  string        `json:"status"`
	ErrorMessage            string        `json:"errorMessage"`
	EnvVars                 []EnvVar      `json:"envVars"`
	Ports                   []ServicePort `json:"ports"`
	GitProvider             string        `json:"gitProvider"`
	CommitHash              string        `json:"commitHash"`
	Memory                  string        `json:"memory"`
	VCPUs                   string        `json:"vcpus"`
	CustomDomain            string        `json:"customDomain"`
	CustomDomainStatus      string        `json:"customDomainStatus"`
	BuildPack               string        `json:"buildPack"`
	BuildCommand            string        `json:"buildCommand"`
	StartCommand            string        `json:"startCommand"`
	PublishDirectory        string        `json:"publishDirectory"`
	RootDirectory           string        `json:"rootDirectory"`
	DockerfilePath          string        `json:"dockerfilePath"`
	TeardownEnabled         bool          `json:"teardownEnabled"`
	TeardownOverlapSeconds  *int          `json:"teardownOverlapSeconds"`
	TeardownDrainingSeconds *int          `json:"teardownDrainingSeconds"`
	DestroyTimeoutSeconds   int           `json:"destroyTimeoutSeconds"`
	CreatedAt               string        `json:"createdAt"`
	UpdatedAt               string        `json:"updatedAt"`
}
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type ServicePort struct {
	Name             string `json:"name"`
	Port             string `json:"port"`
	Protocol         string `json:"protocol"`
	Visibility       string `json:"visibility"`
	AuthPolicy       string `json:"authPolicy"`
	InternalEndpoint string `json:"internalEndpoint"`
	PublicEndpoint   string `json:"publicEndpoint"`
}
type VolumeSpec struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	SizeGi    int    `json:"sizeGi,omitempty"`
}
type ServicePortInput struct {
	Name       string `json:"name"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	Visibility string `json:"visibility"`
	// AuthPolicy is the public HTTP auth policy: public, org_sso, or deployer_sso.
	// Empty uses the installation default.
	AuthPolicy string `json:"authPolicy,omitempty"`
}
type BucketMountInput struct {
	Name         string `json:"name"`
	MountPath    string `json:"mountPath,omitempty"`
	Mode         string `json:"mode,omitempty"`
	Prefix       string `json:"prefix,omitempty"`
	SyncInterval int    `json:"syncInterval,omitempty"`
}
type CreateServiceInput struct {
	Name          string             `json:"name,omitempty"`
	Subdomain     string             `json:"subdomain,omitempty"`
	Source        string             `json:"source,omitempty"`
	Repo          string             `json:"repo,omitempty"`
	Image         string             `json:"image,omitempty"`
	Host          string             `json:"host,omitempty"`
	Branch        string             `json:"branch,omitempty"`
	Project       string             `json:"project,omitempty"`
	WorkspaceSlug string             `json:"workspaceSlug,omitempty"`
	BuildPack     string             `json:"buildPack,omitempty"`
	Ports         []ServicePortInput `json:"ports,omitempty"`
	// AuthPolicy sets the auth policy of the public HTTP endpoint (public,
	// org_sso, or deployer_sso) without spelling out Ports.
	AuthPolicy              string            `json:"authPolicy,omitempty"`
	EnvVars                 []EnvVar          `json:"envVars,omitempty"`
	Memory                  string            `json:"memory,omitempty"`
	VCPUs                   string            `json:"vcpus,omitempty"`
	BuildCommand            string            `json:"buildCommand,omitempty"`
	StartCommand            string            `json:"startCommand,omitempty"`
	PublishDirectory        string            `json:"publishDirectory,omitempty"`
	RootDirectory           string            `json:"rootDirectory,omitempty"`
	DockerfilePath          string            `json:"dockerfilePath,omitempty"`
	Regions                 []string          `json:"regions,omitempty"`
	Volumes                 []VolumeSpec      `json:"volumes,omitempty"`
	Bucket                  *BucketMountInput `json:"bucket,omitempty"`
	DestroyTimeoutSeconds   int               `json:"destroyTimeoutSeconds,omitempty"`
	TeardownEnabled         bool              `json:"teardownEnabled,omitempty"`
	TeardownOverlapSeconds  *int              `json:"teardownOverlapSeconds,omitempty"`
	TeardownDrainingSeconds *int              `json:"teardownDrainingSeconds,omitempty"`
}
type UpdateServiceInput struct {
	Name          string             `json:"name,omitempty"`
	ServiceID     string             `json:"serviceId,omitempty"`
	Project       string             `json:"project,omitempty"`
	ProjectID     string             `json:"projectId,omitempty"`
	WorkspaceSlug string             `json:"workspaceSlug,omitempty"`
	Source        *string            `json:"source,omitempty"`
	Image         *string            `json:"image,omitempty"`
	Repo          *string            `json:"repo,omitempty"`
	Host          *string            `json:"host,omitempty"`
	Branch        *string            `json:"branch,omitempty"`
	BuildPack     *string            `json:"buildPack,omitempty"`
	Memory        *string            `json:"memory,omitempty"`
	VCPUs         *string            `json:"vcpus,omitempty"`
	Ports         []ServicePortInput `json:"ports,omitempty"`
	// AuthPolicy changes the auth policy of the public HTTP endpoint (public,
	// org_sso, or deployer_sso) without resending Ports.
	AuthPolicy              string            `json:"authPolicy,omitempty"`
	EnvVars                 []EnvVar          `json:"envVars,omitempty"`
	BuildCommand            *string           `json:"buildCommand,omitempty"`
	StartCommand            *string           `json:"startCommand,omitempty"`
	PublishDirectory        *string           `json:"publishDirectory,omitempty"`
	RootDirectory           *string           `json:"rootDirectory,omitempty"`
	DockerfilePath          *string           `json:"dockerfilePath,omitempty"`
	Volumes                 []VolumeSpec      `json:"volumes,omitempty"`
	Bucket                  *BucketMountInput `json:"bucket,omitempty"`
	DestroyTimeoutSeconds   *int              `json:"destroyTimeoutSeconds,omitempty"`
	TeardownEnabled         *bool             `json:"teardownEnabled,omitempty"`
	TeardownOverlapSeconds  *int              `json:"teardownOverlapSeconds,omitempty"`
	TeardownDrainingSeconds *int              `json:"teardownDrainingSeconds,omitempty"`
}
type CreateServiceResult struct {
	ServiceID string        `json:"serviceId"`
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Repo      string        `json:"repo"`
	Ports     []ServicePort `json:"ports"`
}
type UpdateServiceResult struct {
	ServiceID string `json:"serviceId"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}
type DeleteServiceInput struct {
	Name          string `json:"name,omitempty"`
	ServiceID     string `json:"serviceId,omitempty"`
	Project       string `json:"project,omitempty"`
	ProjectID     string `json:"projectId,omitempty"`
	WorkspaceSlug string `json:"workspaceSlug,omitempty"`
}
type DeleteServiceResult struct {
	ServiceID string `json:"serviceId"`
	Name      string `json:"name"`
	Message   string `json:"message"`
}
type SetSecretsResult struct {
	ServiceID string `json:"serviceId"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}
type SetSecretsInput struct {
	Name          string   `json:"name,omitempty"`
	ServiceID     string   `json:"serviceId,omitempty"`
	Project       string   `json:"project,omitempty"`
	ProjectID     string   `json:"projectId,omitempty"`
	WorkspaceSlug string   `json:"workspaceSlug,omitempty"`
	EnvVars       []EnvVar `json:"envVars"`
	Replace       bool     `json:"replace,omitempty"`
}
type DeleteSecretsInput struct {
	Name          string   `json:"name,omitempty"`
	ServiceID     string   `json:"serviceId,omitempty"`
	Project       string   `json:"project,omitempty"`
	ProjectID     string   `json:"projectId,omitempty"`
	WorkspaceSlug string   `json:"workspaceSlug,omitempty"`
	Keys          []string `json:"keys"`
}
type ExecSession struct {
	URL       string `json:"url"`
	Token     string `json:"token"`
	ServiceID string `json:"serviceId"`
}
type ExecResult struct {
	ExitCode int    `json:"exitCode"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}
type ExecInput struct {
	ServiceID     string
	Name          string
	Project       string
	WorkspaceSlug string
}
type Workspace struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsDefault bool   `json:"isDefault"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}
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
type WorkspaceInvite struct {
	ID                 string `json:"id"`
	WorkspaceID        string `json:"workspaceId"`
	WorkspaceName      string `json:"workspaceName"`
	WorkspaceSlug      string `json:"workspaceSlug"`
	InviterDisplayName string `json:"inviterDisplayName"`
	InviteeDisplayName string `json:"inviteeDisplayName"`
	Role               string `json:"role"`
	Status             string `json:"status"`
	CreatedAt          string `json:"createdAt"`
}
type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
type CreateProjectInput struct {
	Name          string `json:"name"`
	WorkspaceSlug string `json:"workspaceSlug,omitempty"`
}
type DNSZone struct {
	ID        string `json:"id"`
	Zone      string `json:"zone"`
	Status    string `json:"status"`
	Error     string `json:"error"`
	CreatedAt string `json:"createdAt"`
}
type ZoneRecord struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	TTL       int    `json:"ttl"`
	Managed   bool   `json:"managed"`
	CreatedAt string `json:"createdAt"`
}
type AddDomainResult struct {
	ServiceID string `json:"serviceId"`
	Domain    string `json:"domain"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
type RemoveDomainResult struct {
	ServiceID string `json:"serviceId"`
	Message   string `json:"message"`
}
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
type TemplateService struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Source string `json:"source"`
	Image  string `json:"image"`
	Memory string `json:"memory"`
	VCPUs  string `json:"vcpus"`
}
type TemplateOutput struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Sensitive   bool   `json:"sensitive"`
}
type TemplateVariableValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type TemplateDeployInput struct {
	Template      string                  `json:"template"`
	Name          string                  `json:"name"`
	WorkspaceSlug string                  `json:"workspaceSlug,omitempty"`
	Project       string                  `json:"project,omitempty"`
	Regions       []string                `json:"regions,omitempty"`
	Variables     []TemplateVariableValue `json:"variables,omitempty"`
}
type TemplateDeployResult struct {
	TemplateInstanceID string                    `json:"templateInstanceId"`
	ProjectID          string                    `json:"projectId"`
	Services           []TemplateDeployedService `json:"services"`
	Outputs            []TemplateDeployedOutput  `json:"outputs"`
}
type TemplateDeployedService struct {
	ServiceID string                    `json:"serviceId"`
	Key       string                    `json:"key"`
	Name      string                    `json:"name"`
	Status    string                    `json:"status"`
	Endpoints []TemplateServiceEndpoint `json:"endpoints"`
}
type TemplateServiceEndpoint struct {
	Name             string `json:"name"`
	Port             string `json:"port"`
	Protocol         string `json:"protocol"`
	Visibility       string `json:"visibility"`
	InternalEndpoint string `json:"internalEndpoint"`
	PublicEndpoint   string `json:"publicEndpoint"`
}
type TemplateDeployedOutput struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Sensitive   bool   `json:"sensitive"`
	Value       string `json:"value"`
}
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
type AccountStatus struct {
	ID               string   `json:"id"`
	Email            string   `json:"email"`
	DisplayName      string   `json:"displayName"`
	Username         string   `json:"username"`
	GitHubUsername   string   `json:"githubUsername"`
	HasGitHubOAuth   bool     `json:"hasGitHubOAuth"`
	HasGitHubApp     bool     `json:"hasGitHubApp"`
	DefaultWorkspace string   `json:"defaultWorkspace"`
	SubscriptionTier string   `json:"subscriptionTier"`
	GitHubScopes     []string `json:"githubScopes"`
}
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
type UsageLineItem struct {
	Quantity   string `json:"quantity"`
	UnitPrice  string `json:"unitPrice"`
	Unit       string `json:"unit"`
	TotalCents int    `json:"totalCents"`
}
type LogType = string

const (
	LogTypeBuild   LogType = "BUILD"
	LogTypeRuntime LogType = "RUNTIME"
)

type LogsInput struct {
	ServiceID string  `json:"serviceId"`
	LogType   LogType `json:"logType"`
	StartTime string  `json:"startTime,omitempty"`
	EndTime   string  `json:"endTime,omitempty"`
	Query     string  `json:"query,omitempty"`
	Limit     int     `json:"limit,omitempty"`
}
type LogEntry struct {
	Timestamp  string `json:"timestamp"`
	Level      string `json:"level"`
	Message    string `json:"message"`
	Attributes string `json:"attributes"`
}
type LogsResult struct {
	Entries []LogEntry `json:"entries"`
	HasMore bool       `json:"hasMore"`
}
type MetricDataPoint struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}
type MetricSeries struct {
	Metric     string            `json:"metric"`
	DataPoints []MetricDataPoint `json:"dataPoints"`
}
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
type VolumeInfo struct {
	Name        string `json:"name"`
	MountPath   string `json:"mountPath"`
	SizeGi      int    `json:"sizeGi"`
	Status      string `json:"status"`
	DeleteAfter string `json:"deleteAfter"`
}
type CreateRepoInput struct {
	Name          string `json:"name"`
	Host          string `json:"host,omitempty"`
	Description   string `json:"description,omitempty"`
	Project       string `json:"project,omitempty"`
	WorkspaceSlug string `json:"workspaceSlug,omitempty"`
}
type CreateRepoResult struct {
	Name      string `json:"name"`
	GitRemote string `json:"gitRemote"`
	ExpiresAt string `json:"expiresAt"`
	Message   string `json:"message"`
}
type GetRepoTokenInput struct {
	Name          string `json:"name"`
	Host          string `json:"host,omitempty"`
	WorkspaceSlug string `json:"workspaceSlug,omitempty"`
}
type GetRepoTokenResult struct {
	GitRemote string `json:"gitRemote"`
	ExpiresAt string `json:"expiresAt"`
}
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
type ReadChatResult struct {
	Messages   []ChatMessage `json:"messages"`
	NextCursor int           `json:"nextCursor"`
	HasMore    bool          `json:"hasMore"`
}
type SendChatResult struct {
	Seq       int    `json:"seq"`
	MessageID string `json:"messageId"`
}
