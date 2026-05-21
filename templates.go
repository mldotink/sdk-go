package ink

import "context"

func (c *Client) ListTemplates(ctx context.Context, search string) ([]Template, error) {
	resp, err := listTemplates(ctx, c.gql, optStr(search))
	if err != nil {
		return nil, err
	}
	return resp.TemplateList, nil
}
func (c *Client) DeployTemplate(ctx context.Context, input TemplateDeployInput) (*TemplateDeployResult, error) {
	resp, err := deployTemplate(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.TemplateDeploy, nil
}
func (c *Client) ListTemplateInstances(ctx context.Context, project, projectID, workspaceSlug string) ([]TemplateInstance, error) {
	resp, err := listTemplateInstances(ctx, c.gql, optStr(project), optStr(projectID), optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return resp.TemplateInstanceList, nil
}
