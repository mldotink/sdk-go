package ink

import "context"

// ListTemplates returns available deployment templates, optionally filtered by search term.
func (c *Client) ListTemplates(ctx context.Context, search string) ([]Template, error) {
	resp, err := listTemplates(ctx, c.gql, optStr(search))
	if err != nil {
		return nil, err
	}
	return resp.TemplateList, nil
}

// DeployTemplate deploys a template into a workspace/project.
func (c *Client) DeployTemplate(ctx context.Context, input TemplateDeployInput) (*TemplateDeployResult, error) {
	resp, err := deployTemplate(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.TemplateDeploy, nil
}

// ListTemplateInstances returns all template instances in a project.
func (c *Client) ListTemplateInstances(ctx context.Context, project, projectID, workspaceSlug string) ([]TemplateInstance, error) {
	resp, err := listTemplateInstances(ctx, c.gql, optStr(project), optStr(projectID), optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return resp.TemplateInstanceList, nil
}
