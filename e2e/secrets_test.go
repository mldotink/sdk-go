//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestSecrets(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
	defer cancel()

	svcName := name("sec")
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
	err = client.SetSecrets(ctx, ink.SetSecretsInput{
		ServiceID:     result.ServiceID,
		WorkspaceSlug: workspace,
		EnvVars: []ink.EnvVar{
			{Key: "FOO", Value: "bar"},
			{Key: "BAZ", Value: "qux"},
		},
	})
	requireNoError(t, err, "set secrets")
	waitForActive(t, result.ServiceID, 3*time.Minute)

	svc, err := client.GetService(ctx, result.ServiceID)
	requireNoError(t, err, "get service after set secrets")
	envMap := envToMap(svc.EnvVars)
	requireEqual(t, envMap["FOO"], "bar", "FOO value")
	requireEqual(t, envMap["BAZ"], "qux", "BAZ value")
	err = client.SetSecrets(ctx, ink.SetSecretsInput{
		ServiceID:     result.ServiceID,
		WorkspaceSlug: workspace,
		EnvVars: []ink.EnvVar{
			{Key: "FOO", Value: "updated"},
		},
	})
	requireNoError(t, err, "update secrets")
	waitForActive(t, result.ServiceID, 3*time.Minute)

	svc, err = client.GetService(ctx, result.ServiceID)
	requireNoError(t, err, "get service after update secrets")
	envMap = envToMap(svc.EnvVars)
	requireEqual(t, envMap["FOO"], "updated", "FOO updated value")
	requireEqual(t, envMap["BAZ"], "qux", "BAZ should be preserved after merge")
	err = client.DeleteSecrets(ctx, ink.DeleteSecretsInput{
		ServiceID:     result.ServiceID,
		WorkspaceSlug: workspace,
		Keys:          []string{"BAZ"},
	})
	requireNoError(t, err, "delete secret")
	waitForActive(t, result.ServiceID, 3*time.Minute)

	svc, err = client.GetService(ctx, result.ServiceID)
	requireNoError(t, err, "get service after delete secret")
	envMap = envToMap(svc.EnvVars)
	requireEqual(t, envMap["FOO"], "updated", "FOO value after delete")
	_, hasBaz := envMap["BAZ"]
	if hasBaz {
		t.Fatal("BAZ should be deleted")
	}
}

func envToMap(vars []ink.EnvVar) map[string]string {
	m := make(map[string]string, len(vars))
	for _, v := range vars {
		m[v.Key] = v.Value
	}
	return m
}
