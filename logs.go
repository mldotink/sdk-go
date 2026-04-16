package ink

import "context"

// GetLogs fetches log entries from a service.
func (c *Client) GetLogs(ctx context.Context, input LogsInput) (*LogsResult, error) {
	resp, err := getServiceLogs(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.ServiceLogs, nil
}
