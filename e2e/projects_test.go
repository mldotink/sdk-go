//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestProjectCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	projName := name("proj")
	proj, err := client.CreateProject(ctx, ink.CreateProjectInput{
		Name:          projName,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "create project")
	requireEqual(t, proj.Name, projName, "project name")
	requireNotEmpty(t, proj.ID, "project id is empty")
	requireNotEmpty(t, proj.Slug, "project slug is empty")
	t.Logf("created project %s (slug: %s)", proj.Name, proj.Slug)
	cleanupProject(t, proj.Slug)
	projects, err := client.ListProjects(ctx, workspace)
	requireNoError(t, err, "list projects")
	found := false
	for _, p := range projects {
		if p.ID == proj.ID {
			found = true
			break
		}
	}
	requireTrue(t, found, "project %s not found in list", proj.ID)
	err = client.DeleteProject(ctx, proj.Slug, workspace)
	requireNoError(t, err, "delete project")
	projects, err = client.ListProjects(ctx, workspace)
	requireNoError(t, err, "list projects after delete")
	for _, p := range projects {
		if p.ID == proj.ID {
			t.Fatalf("project %s should be deleted", proj.ID)
		}
	}
}
