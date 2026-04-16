package ink

import "context"

// CreateRepo creates an internal git repository for source deployments.
func (c *Client) CreateRepo(ctx context.Context, input CreateRepoInput) (*CreateRepoResult, error) {
	resp, err := createRepo(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.RepoCreate, nil
}

// GetRepoToken obtains a short-lived push token for an internal git repository.
func (c *Client) GetRepoToken(ctx context.Context, input GetRepoTokenInput) (*GetRepoTokenResult, error) {
	resp, err := getRepoToken(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.RepoGetToken, nil
}
