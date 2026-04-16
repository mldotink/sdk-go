# Ink Go SDK

Go client library for the [Ink](https://ml.ink) cloud platform.

## Installation

```bash
go get github.com/mldotink/sdk-go
```

## Usage

### Create a client

```go
import "github.com/mldotink/sdk-go"

client := ink.NewClient(ink.Config{
    APIKey: "dk_live_...", // https://ml.ink/account/api-keys
})
```

### Deploy a service

```go
result, err := client.CreateService(ctx, ink.CreateServiceInput{
    Name:   "my-api",
    Source: "image",
    Image:  "nginx:latest",
    Memory: "256Mi",
    VCPUs:  "0.25",
})
fmt.Println(result.ServiceID, result.Status)
```

### Get service status

```go
svc, err := client.GetService(ctx, serviceID)
fmt.Println(svc.Name, svc.Status)
```

### List services in a workspace

```go
services, err := client.ListServices(ctx, "my-workspace", "")
for _, s := range services {
    fmt.Println(s.Name, s.Status)
}
```

### Update and redeploy a service

```go
newImage := "nginx:1.27"
result, err := client.UpdateService(ctx, ink.UpdateServiceInput{
    ServiceID: serviceID,
    Image:     &newImage,
})
```

### Set environment variables

```go
err := client.SetSecrets(ctx, ink.SetSecretsInput{
    ServiceID: serviceID,
    EnvVars: []ink.EnvVar{
        {Key: "DATABASE_URL", Value: "postgres://..."},
        {Key: "API_KEY", Value: "secret"},
    },
})
```

### Run a one-shot command

```go
result, err := client.Exec(ctx, ink.ExecInput{ServiceID: serviceID}, "ls -la /app")
fmt.Printf("exit %d\n%s", result.ExitCode, result.Stdout)
```

### Interactive shell session

```go
import "github.com/mldotink/sdk-go/exec"

session, err := exec.Dial(ctx, client, serviceID)
if err != nil {
    log.Fatal(err)
}
defer session.Close()

go io.Copy(os.Stdout, session.Stdout())
go io.Copy(os.Stdout, session.Stderr())
io.Copy(session.Stdin(), os.Stdin)
session.Wait()
```

### Deploy a template

```go
result, err := client.DeployTemplate(ctx, ink.TemplateDeployInput{
    Template:      "postgres",
    WorkspaceSlug: "my-workspace",
    Variables: []ink.TemplateVariableValue{
        {Key: "POSTGRES_DB", Value: "mydb"},
    },
})
fmt.Println(result.TemplateInstanceID)
```

### DNS management

```go
zones, err := client.ListDNSZones(ctx, "my-workspace")
record, err := client.AddDNSRecord(ctx, "example.com", "api", "A", "1.2.3.4", 300, "my-workspace")
err = client.DeleteDNSRecord(ctx, "example.com", record.ID, "my-workspace")
```

### Billing

```go
breakdown, err := client.GetUsageBillBreakdown(ctx, "my-workspace")
fmt.Printf("Current bill: $%.2f\n", float64(breakdown.CurrentBillCents)/100)
```

## API Reference

### Services
| Method | Description |
|--------|-------------|
| `CreateService(ctx, CreateServiceInput)` | Deploy a new service |
| `GetService(ctx, id)` | Get full service details |
| `ListServices(ctx, workspaceSlug, projectSlug)` | List services |
| `UpdateService(ctx, UpdateServiceInput)` | Reconfigure and redeploy |
| `DeleteService(ctx, DeleteServiceInput)` | Permanently delete a service |

### Secrets
| Method | Description |
|--------|-------------|
| `SetSecrets(ctx, SetSecretsInput)` | Set environment variables |
| `DeleteSecrets(ctx, DeleteSecretsInput)` | Remove environment variables |

### Exec
| Method | Description |
|--------|-------------|
| `Exec(ctx, ExecInput, command)` | Run a one-shot command (30s timeout) |
| `ExecURL(ctx, serviceID)` | Get WebSocket URL for interactive shell |
| `exec.Dial(ctx, client, serviceID)` | Open an interactive shell session |

### Workspaces
| Method | Description |
|--------|-------------|
| `ListWorkspaces(ctx)` | List workspaces the user belongs to |
| `CreateWorkspace(ctx, name, slug, description)` | Create a workspace |
| `DeleteWorkspace(ctx, id)` | Delete a workspace |
| `ListWorkspaceMembers(ctx, workspaceSlug)` | List members |
| `InviteToWorkspace(ctx, workspaceID, user, role)` | Invite a user |
| `RemoveWorkspaceMember(ctx, workspaceID, userID)` | Remove a member |
| `ListMyInvites(ctx)` | List pending invitations for the current user |
| `AcceptInvite(ctx, inviteID)` | Accept an invitation |
| `DeclineInvite(ctx, inviteID)` | Decline an invitation |
| `RevokeInvite(ctx, inviteID)` | Cancel a pending invitation |

### Projects
| Method | Description |
|--------|-------------|
| `ListProjects(ctx, workspaceSlug)` | List projects |
| `CreateProject(ctx, CreateProjectInput)` | Create a project |
| `DeleteProject(ctx, slug, workspaceSlug)` | Delete a project |

### DNS
| Method | Description |
|--------|-------------|
| `ListDNSZones(ctx, workspaceSlug)` | List DNS zones |
| `ListDNSRecords(ctx, zone, workspaceSlug)` | List records in a zone |
| `AddDNSRecord(ctx, zone, name, type, content, ttl, workspaceSlug)` | Create a record |
| `DeleteDNSRecord(ctx, zone, recordID, workspaceSlug)` | Delete a record |

### Domains
| Method | Description |
|--------|-------------|
| `AddDomain(ctx, serviceName, domain, project, workspaceSlug)` | Attach a custom domain |
| `RemoveDomain(ctx, serviceName, project, workspaceSlug)` | Detach a custom domain |

### Templates
| Method | Description |
|--------|-------------|
| `ListTemplates(ctx, search)` | Browse available templates |
| `DeployTemplate(ctx, TemplateDeployInput)` | Deploy a template |
| `ListTemplateInstances(ctx, project, projectID, workspaceSlug)` | List deployed instances |

### Account & Billing
| Method | Description |
|--------|-------------|
| `GetAccountStatus(ctx)` | Current user details |
| `GetUsageBillBreakdown(ctx, workspaceSlug)` | Current billing period breakdown |

### Observability
| Method | Description |
|--------|-------------|
| `GetLogs(ctx, LogsInput)` | Fetch service log entries |
| `GetMetrics(ctx, serviceID, timeRange, maxDataPoints)` | CPU/memory/network metrics |

### Repos
| Method | Description |
|--------|-------------|
| `CreateRepo(ctx, CreateRepoInput)` | Create an internal git repo |
| `GetRepoToken(ctx, GetRepoTokenInput)` | Get a short-lived push token |

### Chat
| Method | Description |
|--------|-------------|
| `SendChatMessage(ctx, workspaceSlug, channel, content)` | Post a message |
| `ReadChat(ctx, workspaceSlug, channel, cursor, limit)` | Read messages |

## exec.Session methods

| Method | Description |
|--------|-------------|
| `Stdin()` | Writer for shell input |
| `Stdout()` | Reader for shell stdout |
| `Stderr()` | Reader for shell stderr |
| `Resize(w, h)` | Send terminal resize |
| `Wait()` | Block until session ends |
| `Close()` | Terminate the session |
