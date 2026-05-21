//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestServiceImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	svcName := name("nginx")
	result, err := client.CreateService(ctx, ink.CreateServiceInput{
		Name:          svcName,
		Source:        "image",
		Image:         "nginx:latest",
		Memory:        "256Mi",
		VCPUs:         "0.25",
		Ports:         publicHTTPPort(80),
		Project:       project,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "create image service")
	requireNotEmpty(t, result.ServiceID, "service id is empty")
	requireEqual(t, result.Name, svcName, "service name")
	t.Logf("created service %s (id: %s)", result.Name, result.ServiceID)
	cleanupService(t, svcName, project)
	svc := waitForActive(t, result.ServiceID, 3*time.Minute)
	requireEqual(t, svc.Status, "active", "service status")
	requireEqual(t, svc.Memory, "256Mi", "service memory")
	requireEqual(t, svc.VCPUs, "0.25", "service vcpus")
	requireNotEmpty(t, svc.Subdomain, "service subdomain is empty")
	url := publicHTTPURL(t, svc)
	status, _ := httpGet(t, url)
	requireEqual(t, status, 200, "expected HTTP 200")
	t.Logf("service live at %s", url)
}
