package ink

import "context"

func (c *Client) GetLogs(ctx context.Context, input LogsInput) (*LogsResult, error) {
	resp, err := getServiceLogs(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.ServiceLogs, nil
}
