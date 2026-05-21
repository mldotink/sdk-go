//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func name(suffix string) string {
	return fmt.Sprintf("e2e-%s-%s", runID, suffix)
}

func waitForActive(t *testing.T, serviceID string, timeout time.Duration) *ink.Service {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	last := "not polled yet"
	for {
		svc, err := client.GetService(ctx, serviceID)
		if err != nil {
			last = err.Error()
			t.Logf("  service %s poll error: %v", serviceID, err)
		} else {
			t.Logf("  service %s status: %s", svc.Name, svc.Status)
			last = svc.Status
			switch svc.Status {
			case "active":
				return svc
			case "failed", "crashed", "cancelled":
				t.Fatalf("service %s reached terminal status %q: %s", svc.Name, svc.Status, svc.ErrorMessage)
			}
		}

		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for service %s to become active (last result: %s)", serviceID, last)
		case <-ticker.C:
		}
	}
}

func cleanupService(t *testing.T, serviceName, project string) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		_, err := client.DeleteService(ctx, ink.DeleteServiceInput{
			Name:          serviceName,
			Project:       project,
			WorkspaceSlug: workspace,
		})
		if err := ignoreNotFound(err); err != nil {
			t.Logf("cleanup: delete service %s: %v", serviceName, err)
		}
	})
}

func cleanupProject(t *testing.T, slug string) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := ignoreNotFound(client.DeleteProject(ctx, slug, workspace)); err != nil {
			t.Logf("cleanup: delete project %s: %v", slug, err)
		}
	})
}

func gitPush(t *testing.T, gitRemote string, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	cmds := [][]string{
		{"git", "init", "-b", "main"},
		{"git", "add", "."},
		{"git", "config", "user.email", "e2e@ml.ink"},
		{"git", "config", "user.name", "Ink E2E"},
		{"git", "commit", "-m", "e2e test"},
		{"git", "remote", "add", "ink", gitRemote},
		{"git", "push", "ink", "main"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return dir
}

func httpGet(t *testing.T, url string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	status, body, err := fetchURL(ctx, url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	return status, body
}

func fetchURL(ctx context.Context, url string) (int, string, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return resp.StatusCode, "", err
	}
	return resp.StatusCode, string(body), nil
}

func waitForHTTPContains(t *testing.T, url, needle string, timeout time.Duration) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var last string
	for {
		status, body, err := fetchURL(ctx, url)
		if err == nil && status == http.StatusOK && strings.Contains(body, needle) {
			return
		}
		if err != nil {
			last = err.Error()
		} else {
			last = fmt.Sprintf("status=%d body=%q", status, body)
		}

		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for %s to contain %q; last response: %s", url, needle, last)
		case <-ticker.C:
		}
	}
}

func waitForHTTPGone(t *testing.T, url string, timeout time.Duration) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	last := "not checked yet"
	for {
		status, body, err := fetchURL(ctx, url)
		if err != nil || status >= 400 {
			return
		}
		last = fmt.Sprintf("status=%d body=%q", status, body)

		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for %s to become unavailable after delete; last response: %s", url, last)
		case <-ticker.C:
		}
	}
}

func waitForServiceGone(t *testing.T, serviceID string, timeout time.Duration) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		svc, err := client.GetService(ctx, serviceID)
		if err != nil {
			return
		}
		if svc.Status == "removed" {
			return
		}

		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for service %s to be deleted; last status: %s", serviceID, svc.Status)
		case <-ticker.C:
		}
	}
}

func waitForVolumeStatus(t *testing.T, volumeName, project, status string, timeout time.Duration) ink.VolumeInfo {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	last := "not checked yet"
	for {
		volumes, err := client.ListVolumes(ctx, workspace, project)
		if err != nil {
			last = err.Error()
		} else if vol := findVolume(volumes, volumeName); vol != nil {
			last = vol.Status
			if vol.Status == status {
				return *vol
			}
		} else {
			last = "not found"
		}

		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for volume %s to be %s (last result: %s)", volumeName, status, last)
		case <-ticker.C:
		}
	}
}

func waitForVolumeGone(t *testing.T, volumeName, project string, timeout time.Duration) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	last := "not checked yet"
	for {
		volumes, err := client.ListVolumes(ctx, workspace, project)
		if err != nil {
			last = err.Error()
		} else if vol := findVolume(volumes, volumeName); vol == nil {
			return
		} else {
			last = vol.Status
		}

		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for volume %s to be deleted (last result: %s)", volumeName, last)
		case <-ticker.C:
		}
	}
}

func findVolume(volumes []ink.VolumeInfo, name string) *ink.VolumeInfo {
	for i := range volumes {
		if volumes[i].Name == name {
			return &volumes[i]
		}
	}
	return nil
}

func volumeNameSet(volumes []ink.VolumeInfo) map[string]bool {
	names := make(map[string]bool, len(volumes))
	for _, vol := range volumes {
		names[vol.Name] = true
	}
	return names
}

func cleanupServiceAndVolumes(t *testing.T, serviceID, serviceName, project string, volumeNames ...string) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		if serviceID != "" || serviceName != "" {
			_, err := client.DeleteService(ctx, ink.DeleteServiceInput{
				ServiceID:     serviceID,
				Name:          serviceName,
				Project:       project,
				WorkspaceSlug: workspace,
			})
			if err := ignoreNotFound(err); err != nil {
				t.Logf("cleanup: delete service %s: %v", serviceName, err)
			}
		}
		for _, volumeName := range volumeNames {
			if err := ignoreNotFound(client.DeleteVolume(ctx, volumeName, project, workspace)); err != nil {
				t.Logf("cleanup: delete volume %s: %v", volumeName, err)
			}
		}
	})
}

func publicHTTPURL(t *testing.T, svc *ink.Service) string {
	t.Helper()
	for _, port := range svc.Ports {
		if strings.EqualFold(port.Protocol, "http") && strings.EqualFold(port.Visibility, "public") && port.PublicEndpoint != "" {
			return port.PublicEndpoint
		}
	}
	t.Fatalf("service %s has no public HTTP endpoint in ports: %+v", svc.Name, svc.Ports)
	return ""
}

func publicHTTPPort(port int) []ink.ServicePortInput {
	return []ink.ServicePortInput{{
		Name:       "http",
		Port:       port,
		Protocol:   "http",
		Visibility: "public",
	}}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func ptr[T any](v T) *T { return &v }

func loadDotEnv() {
	for _, path := range []string{".env", "../.env"} {
		if err := loadDotEnvFile(path); err == nil {
			return
		}
	}
}

func loadDotEnvFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if key == "" || os.Getenv(key) != "" {
			continue
		}
		_ = os.Setenv(key, value)
	}
	return nil
}

func ignoreNotFound(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(strings.ToLower(err.Error()), "not found") {
		return nil
	}
	return err
}

func requireNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

func requireNotEmpty(t *testing.T, value, msg string) {
	t.Helper()
	if value == "" {
		t.Fatal(msg)
	}
}

func requireEqual[T comparable](t *testing.T, got, want T, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %v, want %v", msg, got, want)
	}
}

func requireTrue(t *testing.T, ok bool, msg string, args ...any) {
	t.Helper()
	if !ok {
		t.Fatalf(msg, args...)
	}
}

func requireError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatal(msg)
	}
}
