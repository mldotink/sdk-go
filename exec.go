package ink

import "context"

// ExecURL obtains a short-lived WebSocket URL and token for an interactive
// shell session (valid ~120 seconds). Used by the exec sub-package.
func (c *Client) ExecURL(ctx context.Context, serviceID string) (*ExecSession, error) {
	resp, err := getServiceExecUrl(ctx, c.gql, serviceID)
	if err != nil {
		return nil, err
	}
	return &resp.ServiceExecUrl, nil
}

// Exec runs a command in a running service container and returns the output.
// Maximum 30 second timeout and 1 MiB output. Use the exec sub-package for
// interactive shell sessions.
func (c *Client) Exec(ctx context.Context, target ExecInput, command string) (*ExecResult, error) {
	resp, err := execService(ctx, c.gql, optStr(target.ServiceID), optStr(target.Name), command, optStr(target.Project), optStr(target.WorkspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.ServiceExec, nil
}
