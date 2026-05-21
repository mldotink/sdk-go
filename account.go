package ink

import "context"

func (c *Client) GetAccountStatus(ctx context.Context) (*AccountStatus, error) {
	resp, err := getAccountStatus(ctx, c.gql)
	if err != nil {
		return nil, err
	}
	return resp.AccountStatus, nil
}
