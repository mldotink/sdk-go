//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestDomains(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	zone := envOr("INK_E2E_DOMAIN_ZONE", "")
	if zone == "" {
		t.Skip("INK_E2E_DOMAIN_ZONE is required for custom domain e2e")
	}

	svcName := name("dom")
	domain := fmt.Sprintf("%s.%s", name("d"), zone)
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
	domResult, err := client.AddDomain(ctx, svcName, domain, project, workspace)
	requireNoError(t, err, "add domain")
	requireNotEmpty(t, domResult.ServiceID, "domain result service id is empty")
	t.Logf("added domain %s (status: %s)", domResult.Domain, domResult.Status)
	t.Cleanup(func() {
		rctx, rcancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer rcancel()
		_, err := client.RemoveDomain(rctx, svcName, project, workspace)
		if err != nil {
			t.Logf("cleanup: remove domain: %v", err)
		}
	})
	svc, err := client.GetService(ctx, result.ServiceID)
	requireNoError(t, err, "get service after add domain")
	requireEqual(t, svc.CustomDomain, domain, "custom domain")
	records, err := client.ListDNSRecords(ctx, zone, workspace)
	requireNoError(t, err, "list DNS records")
	foundRecord := false
	for _, r := range records {
		if r.Managed && r.Name == domain {
			foundRecord = true
			break
		}
	}
	requireTrue(t, foundRecord, "managed DNS record for %s not found", domain)
	_, err = client.RemoveDomain(ctx, svcName, project, workspace)
	requireNoError(t, err, "remove domain")
	svc, err = client.GetService(ctx, result.ServiceID)
	requireNoError(t, err, "get service after remove domain")
	if svc.CustomDomain != "" {
		t.Fatalf("custom domain should be cleared, got %q", svc.CustomDomain)
	}
}
