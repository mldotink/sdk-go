# E2E Workflows

## Commands

```sh
make e2e
make e2e E2E_RUN='^TestServiceFromPrivateGitRepoFlow$'
make e2e-list
```

## `TestSmoke`

1. Get account status.
2. Check account id.
3. Check account email.
4. Check default workspace.

## `TestProjectCRUD`

1. Create project.
2. List projects.
3. Check created project is listed.
4. Delete project.
5. List projects.
6. Check deleted project is not listed.

## `TestServiceFromPrivateGitRepoFlow`

1. Create private Ink git repo.
2. Create hello-world app in a temp git repo.
3. Push app to Ink git repo.
4. Create service from repo.
5. Wait for service to become active.
6. Get public service URL.
7. Visit public service URL.
8. Check response body matches expected app content.
9. Delete service.
10. Wait for service to be deleted.
11. Check public service URL no longer works.

## `TestServiceImage`

1. Create service from public Docker image.
2. Wait for service to become active.
3. Check service memory.
4. Check service CPU.
5. Get public service URL.
6. Visit public service URL.
7. Check HTTP 200.

## `TestServiceUpdate`

1. Create service from public Docker image.
2. Wait for service to become active.
3. Update service memory.
4. Update service CPU.
5. Wait for redeploy to become active.
6. Check updated memory.
7. Check updated CPU.

## `TestServiceVolume`

1. Create service from public Docker image with a mounted volume.
2. Wait for service to become active.
3. Check volume is provisioned.
4. Delete service.
5. Check volume still exists as detached.
6. Delete volume.
7. Check volume is gone.

## `TestSecrets`

1. Create service from public Docker image.
2. Wait for service to become active.
3. Set two env vars.
4. Wait for redeploy to become active.
5. Check both env vars exist.
6. Update one env var.
7. Wait for redeploy to become active.
8. Check updated env var changed.
9. Check untouched env var remains.
10. Delete one env var.
11. Wait for redeploy to become active.
12. Check deleted env var is gone.

## `TestExec`

1. Create service from public Docker image.
2. Wait for service to become active.
3. Exec `echo hello`.
4. Check exit code is 0.
5. Check stdout contains `hello`.
6. Exec `false`.
7. Check exit code is non-zero.
8. Exec `cat /etc/os-release`.
9. Check stdout is non-empty.

## `TestTemplates`

1. List templates.
2. Select template.
3. Deploy template.
4. Check template instance id exists.
5. Check template created services.
6. Wait for created services to become active.
7. List template instances.
8. Check deployed template instance is listed.

## `TestPostgresTemplateEndpoints`

1. Deploy Postgres template.
2. Wait for Postgres service to become active.
3. Check template output connection string uses the public TCP endpoint.
4. Create Postgres client service.
5. Connect to the internal endpoint from the client service.
6. Connect to the public endpoint from the local network.
7. Delete created services.
8. Delete created volumes.

## `TestLifecycle`

1. Create project.
2. Create service in project.
3. Wait for service to become active.
4. Delete project.
5. Check service is gone.
6. Check project is not listed.

## `TestDomains`

1. Create service from public Docker image.
2. Wait for service to become active.
3. Attach custom domain.
4. Check service custom domain.
5. List DNS records.
6. Check managed DNS record exists.
7. Remove custom domain.
8. Check service custom domain is cleared.
