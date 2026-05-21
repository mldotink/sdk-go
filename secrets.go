package ink

import "context"

func (c *Client) SetSecrets(ctx context.Context, input SetSecretsInput) error {
	_, err := setSecrets(ctx, c.gql, input)
	return err
}
func (c *Client) DeleteSecrets(ctx context.Context, input DeleteSecretsInput) error {
	_, err := deleteSecrets(ctx, c.gql, input)
	return err
}
