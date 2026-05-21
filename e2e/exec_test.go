//go:build e2e

package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestExec(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	svcName := name("exec")
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

	t.Run("echo", func(t *testing.T) {
		execResult, err := client.Exec(ctx, ink.ExecInput{
			ServiceID:     result.ServiceID,
			WorkspaceSlug: workspace,
		}, "echo hello")
		requireNoError(t, err, "exec echo")
		requireEqual(t, execResult.ExitCode, 0, "echo exit code")
		requireTrue(t, strings.Contains(execResult.Stdout, "hello"), "expected exec stdout to contain hello, got %q", execResult.Stdout)
	})

	t.Run("nonzero_exit", func(t *testing.T) {
		execResult, err := client.Exec(ctx, ink.ExecInput{
			ServiceID:     result.ServiceID,
			WorkspaceSlug: workspace,
		}, "false")
		requireNoError(t, err, "exec false")
		if execResult.ExitCode == 0 {
			t.Fatal("expected non-zero exit code")
		}
	})

	t.Run("read_file", func(t *testing.T) {
		execResult, err := client.Exec(ctx, ink.ExecInput{
			ServiceID:     result.ServiceID,
			WorkspaceSlug: workspace,
		}, "cat /etc/os-release")
		requireNoError(t, err, "exec read file")
		requireEqual(t, execResult.ExitCode, 0, "read file exit code")
		requireNotEmpty(t, execResult.Stdout, "os-release output is empty")
	})
}
