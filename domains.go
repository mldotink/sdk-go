package ink

import "context"

func (c *Client) AddDomain(ctx context.Context, serviceName, domain, project, workspaceSlug string) (*AddDomainResult, error) {
	resp, err := addDomain(ctx, c.gql, serviceName, domain, optStr(project), optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.DomainAdd, nil
}
func (c *Client) RemoveDomain(ctx context.Context, serviceName, project, workspaceSlug string) (*RemoveDomainResult, error) {
	resp, err := removeDomain(ctx, c.gql, serviceName, optStr(project), optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.DomainRemove, nil
}
