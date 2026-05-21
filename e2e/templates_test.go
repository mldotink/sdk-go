//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestTemplates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
	defer cancel()
	templates, err := client.ListTemplates(ctx, "")
	requireNoError(t, err, "list templates")
	if len(templates) == 0 {
		t.Fatal("template catalog should not be empty")
	}
	t.Logf("found %d templates", len(templates))
	var tmpl ink.Template
	for _, tpl := range templates {
		if tpl.Slug == "redis" {
			tmpl = tpl
			break
		}
	}
	if tmpl.Slug == "" {
		tmpl = templates[0]
	}
	t.Logf("deploying template: %s (%s)", tmpl.Name, tmpl.Slug)
	deployName := name("tpl")
	result, err := client.DeployTemplate(ctx, ink.TemplateDeployInput{
		Template:      tmpl.Slug,
		Name:          deployName,
		Project:       project,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "deploy template")
	requireNotEmpty(t, result.TemplateInstanceID, "template instance id is empty")
	if len(result.Services) == 0 {
		t.Fatal("template deployment returned no services")
	}
	t.Logf("deployed template instance %s with %d services", result.TemplateInstanceID, len(result.Services))
	for _, svc := range result.Services {
		cleanupService(t, svc.Name, project)
	}
	for _, svc := range result.Services {
		t.Logf("waiting for template service %s...", svc.Name)
		waitForActive(t, svc.ServiceID, 5*time.Minute)
	}
	instances, err := client.ListTemplateInstances(ctx, project, "", workspace)
	requireNoError(t, err, "list template instances")
	found := false
	for _, inst := range instances {
		if inst.ID == result.TemplateInstanceID {
			found = true
			break
		}
	}
	requireTrue(t, found, "template instance %s not found in list", result.TemplateInstanceID)
}
