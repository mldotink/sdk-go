package ink

import "context"

func (c *Client) ExecURL(ctx context.Context, serviceID string) (*ExecSession, error) {
	resp, err := getServiceExecUrl(ctx, c.gql, serviceID)
	if err != nil {
		return nil, err
	}
	return &resp.ServiceExecUrl, nil
}
func (c *Client) Exec(ctx context.Context, target ExecInput, command string) (*ExecResult, error) {
	resp, err := execService(ctx, c.gql, optStr(target.ServiceID), optStr(target.Name), command, optStr(target.Project), optStr(target.WorkspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.ServiceExec, nil
}
