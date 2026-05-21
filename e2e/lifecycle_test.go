//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	projName := name("lc")
	svcName := name("lcsvc")
	proj, err := client.CreateProject(ctx, ink.CreateProjectInput{
		Name:          projName,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "create project")
	cleanupProject(t, proj.Slug)
	t.Logf("created project %s", proj.Slug)
	result, err := client.CreateService(ctx, ink.CreateServiceInput{
		Name:          svcName,
		Source:        "image",
		Image:         "nginx:latest",
		Memory:        "256Mi",
		VCPUs:         "0.25",
		Ports:         publicHTTPPort(80),
		Project:       proj.Slug,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "create service")
	t.Logf("created service %s in project %s", result.Name, proj.Slug)
	waitForActive(t, result.ServiceID, 3*time.Minute)
	err = client.DeleteProject(ctx, proj.Slug, workspace)
	requireNoError(t, err, "delete project")
	t.Logf("deleted project %s", proj.Slug)
	_, err = client.GetService(ctx, result.ServiceID)
	requireError(t, err, "service should be gone after project deletion")
	projects, err := client.ListProjects(ctx, workspace)
	requireNoError(t, err, "list projects")
	for _, p := range projects {
		if p.ID == proj.ID {
			t.Fatalf("project %s should be deleted", proj.ID)
		}
	}
}
