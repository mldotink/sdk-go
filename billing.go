package ink

import "context"

// GetUsageBillBreakdown returns the current billing period cost breakdown for a workspace.
func (c *Client) GetUsageBillBreakdown(ctx context.Context, workspaceSlug string) (*UsageBillBreakdown, error) {
	resp, err := getUsageBillBreakdown(ctx, c.gql, optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.UsageBillBreakdown, nil
}
