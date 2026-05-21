//go:build e2e

package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	ink "github.com/mldotink/sdk-go"
)

func TestPostgresTemplateEndpoints(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	beforeVolumes, err := client.ListVolumes(ctx, workspace, project)
	requireNoError(t, err, "list volumes before postgres deploy")
	beforeVolumeNames := volumeNameSet(beforeVolumes)
	var serviceIDs []string

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 4*time.Minute)
		defer cleanupCancel()
		for _, serviceID := range serviceIDs {
			_, err := client.DeleteService(cleanupCtx, ink.DeleteServiceInput{
				ServiceID:     serviceID,
				Project:       project,
				WorkspaceSlug: workspace,
			})
			if err := ignoreNotFound(err); err != nil {
				t.Logf("cleanup: delete service %s: %v", serviceID, err)
			}
		}
		volumes, err := client.ListVolumes(cleanupCtx, workspace, project)
		if err != nil {
			t.Logf("cleanup: list volumes: %v", err)
			return
		}
		for _, vol := range volumes {
			if beforeVolumeNames[vol.Name] {
				continue
			}
			if err := ignoreNotFound(client.DeleteVolume(cleanupCtx, vol.Name, project, workspace)); err != nil {
				t.Logf("cleanup: delete volume %s: %v", vol.Name, err)
			}
		}
	})

	deployName := name("pg")
	database := "e2e"
	username := "e2e"
	result, err := client.DeployTemplate(ctx, ink.TemplateDeployInput{
		Template:      "postgres",
		Name:          deployName,
		Project:       project,
		WorkspaceSlug: workspace,
		Variables: []ink.TemplateVariableValue{
			{Key: "database_name", Value: database},
			{Key: "username", Value: username},
			{Key: "storage_gi", Value: "1"},
		},
	})
	requireNoError(t, err, "deploy postgres template")
	requireNotEmpty(t, result.TemplateInstanceID, "template instance id is empty")
	for _, svc := range result.Services {
		serviceIDs = append(serviceIDs, svc.ServiceID)
	}
	t.Logf("deployed postgres template instance %s", result.TemplateInstanceID)

	var endpoint ink.TemplateServiceEndpoint
	for _, svc := range result.Services {
		waitForActive(t, svc.ServiceID, 5*time.Minute)
		for _, candidate := range svc.Endpoints {
			if candidate.Name == "postgres" && candidate.Protocol == "tcp" {
				endpoint = candidate
			}
		}
	}
	requireNotEmpty(t, endpoint.InternalEndpoint, "postgres internal endpoint is empty")
	requireNotEmpty(t, endpoint.PublicEndpoint, "postgres public endpoint is empty")

	publicConnectionString := templateOutputValue(result.Outputs, "connection_string")
	requireTrue(t, strings.Contains(publicConnectionString, endpoint.PublicEndpoint), "connection_string should contain public endpoint")
	password := templateOutputValue(result.Outputs, "password")
	requireNotEmpty(t, password, "postgres password output is empty")
	internalHost, internalPort, err := net.SplitHostPort(endpoint.InternalEndpoint)
	requireNoError(t, err, "parse internal postgres endpoint")

	clientName := name("pgclient")
	clientResult, err := client.CreateService(ctx, ink.CreateServiceInput{
		Name:   clientName,
		Source: "image",
		Image:  "postgres:17",
		Memory: "512Mi",
		VCPUs:  "0.25",
		Ports: []ink.ServicePortInput{{
			Name:       "postgres",
			Port:       5432,
			Protocol:   "tcp",
			Visibility: "internal",
		}},
		EnvVars: []ink.EnvVar{
			{Key: "POSTGRES_PASSWORD", Value: "e2e"},
			{Key: "POSTGRES_USER", Value: "e2e"},
			{Key: "POSTGRES_DB", Value: "e2e"},
			{Key: "TARGET_PGHOST", Value: internalHost},
			{Key: "TARGET_PGPORT", Value: internalPort},
			{Key: "TARGET_PGUSER", Value: username},
			{Key: "TARGET_PGPASSWORD", Value: password},
			{Key: "TARGET_PGDATABASE", Value: database},
		},
		Project:       project,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "create postgres client service")
	serviceIDs = append(serviceIDs, clientResult.ServiceID)
	waitForActive(t, clientResult.ServiceID, 3*time.Minute)

	err = execPSQLInService(ctx, clientResult.ServiceID)
	requireNoError(t, err, "connect to postgres internal endpoint from service network")
	t.Logf("internal postgres endpoint connected: %s", endpoint.InternalEndpoint)

	err = queryPostgresEndpoint(ctx, endpoint.PublicEndpoint, username, password, database)
	requireNoError(t, err, "connect to postgres public endpoint from local network")
	t.Logf("public postgres endpoint connected: %s", endpoint.PublicEndpoint)
}

func templateOutputValue(outputs []ink.TemplateDeployedOutput, key string) string {
	for _, output := range outputs {
		if output.Key == key {
			return output.Value
		}
	}
	return ""
}

func execPSQLInService(ctx context.Context, serviceID string) error {
	command := "PGHOST=\"$TARGET_PGHOST\" PGPORT=\"$TARGET_PGPORT\" PGUSER=\"$TARGET_PGUSER\" PGPASSWORD=\"$TARGET_PGPASSWORD\" PGDATABASE=\"$TARGET_PGDATABASE\" PGSSLMODE=disable PGCONNECT_TIMEOUT=10 psql -v ON_ERROR_STOP=1 -Atc 'SELECT 1'"
	result, err := client.Exec(ctx, ink.ExecInput{
		ServiceID:     serviceID,
		WorkspaceSlug: workspace,
	}, command)
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("psql exit code %d: %s", result.ExitCode, strings.TrimSpace(result.Stderr))
	}
	if strings.TrimSpace(result.Stdout) != "1" {
		return fmt.Errorf("unexpected psql output: %q", result.Stdout)
	}
	return nil
}

func queryPostgresEndpoint(ctx context.Context, endpoint, username, password, database string) error {
	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		return err
	}
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("invalid postgres port %q: %w", port, err)
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   net.JoinHostPort(host, port),
		Path:   database,
	}
	query := u.Query()
	query.Set("sslmode", "disable")
	query.Set("connect_timeout", "10")
	u.RawQuery = query.Encode()

	queryCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", u.String())
	if err != nil {
		return err
	}
	defer db.Close()

	var one int
	if err := db.QueryRowContext(queryCtx, "SELECT 1").Scan(&one); err != nil {
		return err
	}
	if one != 1 {
		return fmt.Errorf("unexpected postgres query result: %d", one)
	}
	return nil
}
