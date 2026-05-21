package ink

import (
	"context"
	"fmt"
)

func (c *Client) CreateService(ctx context.Context, input CreateServiceInput) (*CreateServiceResult, error) {
	resp, err := createService(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.ServiceCreate, nil
}
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
func (c *Client) GetServiceByName(ctx context.Context, name, workspaceSlug, projectSlug string) (*Service, error) {
	resp, err := getServiceByName(ctx, c.gql, name, optStr(workspaceSlug), optStr(projectSlug))
	if err != nil {
		return nil, err
	}
	if resp.ServiceGetByName == nil {
		return nil, fmt.Errorf("ink: service %q not found", name)
	}
	return resp.ServiceGetByName, nil
}
func (c *Client) ListServices(ctx context.Context, workspaceSlug, projectSlug string) ([]Service, error) {
	resp, err := listServices(ctx, c.gql, optStr(workspaceSlug), optStr(projectSlug))
	if err != nil {
		return nil, err
	}
	return resp.ServiceList.Nodes, nil
}
func (c *Client) UpdateService(ctx context.Context, input UpdateServiceInput) (*UpdateServiceResult, error) {
	resp, err := updateService(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.ServiceUpdate, nil
}
func (c *Client) DeleteService(ctx context.Context, input DeleteServiceInput) (*DeleteServiceResult, error) {
	resp, err := deleteService(ctx, c.gql, optStr(input.Name), optStr(input.ServiceID), optStr(input.Project), optStr(input.ProjectID), optStr(input.WorkspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.ServiceDelete, nil
}
