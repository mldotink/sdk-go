//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestServiceUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	svcName := name("upd")
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
	requireNoError(t, err, "create service")
	cleanupService(t, svcName, project)
	waitForActive(t, result.ServiceID, 3*time.Minute)
	upd, err := client.UpdateService(ctx, ink.UpdateServiceInput{
		ServiceID:     result.ServiceID,
		WorkspaceSlug: workspace,
		Memory:        ptr("512Mi"),
		VCPUs:         ptr("0.5"),
	})
	requireNoError(t, err, "update service")
	t.Logf("updated service, status: %s", upd.Status)
	svc := waitForActive(t, result.ServiceID, 3*time.Minute)
	requireEqual(t, svc.Memory, "512Mi", "service memory")
	requireEqual(t, svc.VCPUs, "0.5", "service vcpus")
	t.Logf("service updated: memory=%s vcpus=%s", svc.Memory, svc.VCPUs)
}
