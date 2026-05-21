package ink

import "context"

func (c *Client) GetUsageBillBreakdown(ctx context.Context, workspaceSlug string) (*UsageBillBreakdown, error) {
	resp, err := getUsageBillBreakdown(ctx, c.gql, optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.UsageBillBreakdown, nil
}
