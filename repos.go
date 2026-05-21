package ink

import "context"

func (c *Client) CreateRepo(ctx context.Context, input CreateRepoInput) (*CreateRepoResult, error) {
	resp, err := createRepo(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.RepoCreate, nil
}
func (c *Client) GetRepoToken(ctx context.Context, input GetRepoTokenInput) (*GetRepoTokenResult, error) {
	resp, err := getRepoToken(ctx, c.gql, input)
	if err != nil {
		return nil, err
	}
	return &resp.RepoGetToken, nil
}
