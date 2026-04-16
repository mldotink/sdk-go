package ink

import "context"

// GetMetrics returns CPU, memory, network, and disk metrics for a service.
// maxDataPoints=0 uses the server default (~120 points).
func (c *Client) GetMetrics(ctx context.Context, serviceID string, timeRange MetricTimeRange, maxDataPoints int) (*ServiceMetrics, error) {
	resp, err := getServiceMetrics(ctx, c.gql, serviceID, timeRange, optInt(maxDataPoints))
	if err != nil {
		return nil, err
	}
	return &resp.ServiceMetrics, nil
}
