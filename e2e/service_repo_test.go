//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestServiceFromPrivateGitRepoFlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()

	repoName := name("hello-repo")
	svcName := name("hello-svc")
	expected := fmt.Sprintf("hello from ink e2e %s", runID)

	repo, err := client.CreateRepo(ctx, ink.CreateRepoInput{
		Name:          repoName,
		Project:       project,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "create repo")
	requireNotEmpty(t, repo.GitRemote, "repo git remote is empty")
	t.Logf("created repo %s", repo.Name)

	gitPush(t, repo.GitRemote, map[string]string{
		"package.json": `{"scripts":{"start":"node server.js"},"engines":{"node":">=20"}}` + "\n",
		"server.js": fmt.Sprintf(`const http = require("http");

const body = %q;
const server = http.createServer((req, res) => {
  res.writeHead(200, {"content-type": "text/plain"});
  res.end(body);
});

server.listen(process.env.PORT || 3000, "0.0.0.0");
`, expected),
	})

	result, err := client.CreateService(ctx, ink.CreateServiceInput{
		Name:      svcName,
		Source:    "repo",
		Repo:      repoName,
		Branch:    "main",
		Project:   project,
		BuildPack: "railpack",
		Memory:    "256Mi",
		VCPUs:     "0.25",
		Ports: []ink.ServicePortInput{{
			Name:       "http",
			Port:       3000,
			Protocol:   "http",
			Visibility: "public",
		}},
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "create service")
	t.Logf("created service %s (id: %s)", result.Name, result.ServiceID)
	cleanupService(t, svcName, project)

	svc := waitForActive(t, result.ServiceID, 8*time.Minute)
	requireEqual(t, svc.Status, "active", "service status")

	url := publicHTTPURL(t, svc)
	waitForHTTPContains(t, url, expected, 2*time.Minute)
	t.Logf("service live at %s", url)

	_, err = client.DeleteService(ctx, ink.DeleteServiceInput{
		ServiceID:     result.ServiceID,
		Project:       project,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "delete service")

	waitForServiceGone(t, result.ServiceID, 2*time.Minute)
	waitForHTTPGone(t, url, 2*time.Minute)
}
