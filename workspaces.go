package ink

import "context"

func (c *Client) ListWorkspaces(ctx context.Context) ([]Workspace, error) {
	resp, err := listWorkspaces(ctx, c.gql)
	if err != nil {
		return nil, err
	}
	return resp.WorkspaceList, nil
}
func (c *Client) CreateWorkspace(ctx context.Context, name, slug, description string) (*Workspace, error) {
	resp, err := createWorkspace(ctx, c.gql, name, slug, optStr(description))
	if err != nil {
		return nil, err
	}
	return &resp.WorkspaceCreate, nil
}
func (c *Client) DeleteWorkspace(ctx context.Context, id string) error {
	_, err := deleteWorkspace(ctx, c.gql, id)
	return err
}
func (c *Client) ListWorkspaceMembers(ctx context.Context, workspaceSlug string) ([]WorkspaceMember, error) {
	resp, err := listWorkspaceMembers(ctx, c.gql, workspaceSlug)
	if err != nil {
		return nil, err
	}
	return resp.WorkspaceListMembers, nil
}
func (c *Client) InviteToWorkspace(ctx context.Context, workspaceID, user, role string) (*WorkspaceInvite, error) {
	resp, err := inviteToWorkspace(ctx, c.gql, workspaceID, user, optStr(role))
	if err != nil {
		return nil, err
	}
	return &resp.WorkspaceInvite, nil
}
func (c *Client) RemoveWorkspaceMember(ctx context.Context, workspaceID, userID string) error {
	_, err := removeWorkspaceMember(ctx, c.gql, workspaceID, userID)
	return err
}
func (c *Client) ListMyInvites(ctx context.Context) ([]WorkspaceInvite, error) {
	resp, err := listMyInvites(ctx, c.gql)
	if err != nil {
		return nil, err
	}
	return resp.WorkspaceListMyInvites, nil
}
func (c *Client) ListWorkspaceInvites(ctx context.Context, workspaceSlug string) ([]WorkspaceInvite, error) {
	resp, err := listWorkspaceInvites(ctx, c.gql, workspaceSlug)
	if err != nil {
		return nil, err
	}
	return resp.WorkspaceListInvites, nil
}
func (c *Client) AcceptInvite(ctx context.Context, inviteID string) error {
	_, err := acceptInvite(ctx, c.gql, inviteID)
	return err
}
func (c *Client) DeclineInvite(ctx context.Context, inviteID string) error {
	_, err := declineInvite(ctx, c.gql, inviteID)
	return err
}
func (c *Client) RevokeInvite(ctx context.Context, inviteID string) error {
	_, err := revokeInvite(ctx, c.gql, inviteID)
	return err
}
