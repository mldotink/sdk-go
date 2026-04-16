package ink

import "context"

// AddDomain attaches a custom domain to a service and begins certificate provisioning.
func (c *Client) AddDomain(ctx context.Context, serviceName, domain, project, workspaceSlug string) (*AddDomainResult, error) {
	resp, err := addDomain(ctx, c.gql, serviceName, domain, optStr(project), optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.DomainAdd, nil
}

// RemoveDomain detaches the custom domain from a service.
func (c *Client) RemoveDomain(ctx context.Context, serviceName, project, workspaceSlug string) (*RemoveDomainResult, error) {
	resp, err := removeDomain(ctx, c.gql, serviceName, optStr(project), optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.DomainRemove, nil
}
