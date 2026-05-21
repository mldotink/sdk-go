package ink

import "context"

func (c *Client) GetMetrics(ctx context.Context, serviceID string, timeRange MetricTimeRange, maxDataPoints int) (*ServiceMetrics, error) {
	resp, err := getServiceMetrics(ctx, c.gql, serviceID, timeRange, optInt(maxDataPoints))
	if err != nil {
		return nil, err
	}
	return &resp.ServiceMetrics, nil
}
