//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestServiceStaticPublishDirectory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	repoName := name("static-repo")
	svcName := name("static-svc")
	expected := fmt.Sprintf("static publish directory from ink e2e %s", runID)

	repo, err := client.CreateRepo(ctx, ink.CreateRepoInput{
		Name:          repoName,
		Project:       project,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "create repo")
	requireNotEmpty(t, repo.GitRemote, "repo git remote is empty")
	t.Logf("created repo %s", repo.Name)

	gitPush(t, repo.GitRemote, map[string]string{
		"dist/index.html": fmt.Sprintf(`<!doctype html>
<html>
<head><meta charset="utf-8"><title>Ink static e2e</title></head>
<body>%s</body>
</html>
`, expected),
		"ignored.html": "this file should not be served when publishDirectory=dist\n",
	})

	result, err := client.CreateService(ctx, ink.CreateServiceInput{
		Name:             svcName,
		Source:           "repo",
		Repo:             repoName,
		Branch:           "main",
		Project:          project,
		BuildPack:        "static",
		PublishDirectory: "dist",
		Memory:           "256Mi",
		VCPUs:            "0.25",
		WorkspaceSlug:    workspace,
	})
	requireNoError(t, err, "create static service")
	t.Logf("created static service %s (id: %s)", result.Name, result.ServiceID)
	cleanupService(t, svcName, project)

	svc := waitForActive(t, result.ServiceID, 8*time.Minute)
	requireEqual(t, svc.Status, "active", "service status")
	requireEqual(t, svc.BuildPack, "static", "service build pack")
	requireEqual(t, svc.PublishDirectory, "dist", "service publish directory")

	url := publicHTTPURL(t, svc)
	waitForHTTPContains(t, url, expected, 2*time.Minute)
	t.Logf("static service live at %s", url)
}
