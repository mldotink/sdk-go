package ink

import "context"

// SetSecrets sets environment variables on a service. Existing vars with the
// same key are overwritten; other vars are preserved unless Replace is true.
// Triggers a redeployment.
func (c *Client) SetSecrets(ctx context.Context, input SetSecretsInput) error {
	_, err := setSecrets(ctx, c.gql, input)
	return err
}

// DeleteSecrets removes the specified environment variable keys from a service.
// Triggers a redeployment.
func (c *Client) DeleteSecrets(ctx context.Context, input DeleteSecretsInput) error {
	_, err := deleteSecrets(ctx, c.gql, input)
	return err
}
