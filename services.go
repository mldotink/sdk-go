package ink

import (
	"context"
	"fmt"
)

// CreateService deploys a new service. The returned status is typically "queued"
// — poll GetService to track build/deploy progress.
func (c *Client) CreateService(ctx context.Context, input CreateServiceInput) (*CreateServiceResult, error) {
	resp, err := createService(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.ServiceCreate, nil
}

// GetService returns full details for a single service by ID.
func (c *Client) GetService(ctx context.Context, id string) (*Service, error) {
	resp, err := getService(ctx, c.gql, id)
	if err != nil {
		return nil, err
	}
	if resp.ServiceGet == nil {
		return nil, fmt.Errorf("ink: service %q not found", id)
	}
	return resp.ServiceGet, nil
}

// ListServices returns all services in a workspace, optionally filtered by project slug.
func (c *Client) ListServices(ctx context.Context, workspaceSlug, projectSlug string) ([]Service, error) {
	resp, err := listServices(ctx, c.gql, optStr(workspaceSlug), optStr(projectSlug))
	if err != nil {
		return nil, err
	}
	return resp.ServiceList.Nodes, nil
}

// UpdateService reconfigures a service and triggers a redeployment.
func (c *Client) UpdateService(ctx context.Context, input UpdateServiceInput) (*UpdateServiceResult, error) {
	resp, err := updateService(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.ServiceUpdate, nil
}

// DeleteService permanently removes a service and tears down its resources.
func (c *Client) DeleteService(ctx context.Context, input DeleteServiceInput) (*DeleteServiceResult, error) {
	resp, err := deleteService(ctx, c.gql, optStr(input.Name), optStr(input.ServiceID), optStr(input.Project), optStr(input.ProjectID), optStr(input.WorkspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.ServiceDelete, nil
}
