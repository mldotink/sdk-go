package ink

import "context"

// GetAccountStatus returns the authenticated user's account details.
func (c *Client) GetAccountStatus(ctx context.Context) (*AccountStatus, error) {
	resp, err := getAccountStatus(ctx, c.gql)
	if err != nil {
		return nil, err
	}
	return resp.AccountStatus, nil
}
