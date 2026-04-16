package ink

import "context"

// ListProjects returns all projects in a workspace.
func (c *Client) ListProjects(ctx context.Context, workspaceSlug string) ([]Project, error) {
	resp, err := listProjects(ctx, c.gql, optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return resp.ProjectList.Nodes, nil
}

// CreateProject creates a new project within a workspace.
func (c *Client) CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error) {
	resp, err := createProject(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.ProjectCreate, nil
}

// DeleteProject permanently deletes a project and all its services.
func (c *Client) DeleteProject(ctx context.Context, slug, workspaceSlug string) error {
	_, err := deleteProject(ctx, c.gql, slug, optStr(workspaceSlug))
	return err
}
