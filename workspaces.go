package ink

import "context"

// ListWorkspaces returns all workspaces the authenticated user belongs to.
func (c *Client) ListWorkspaces(ctx context.Context) ([]Workspace, error) {
	resp, err := listWorkspaces(ctx, c.gql)
	if err != nil {
		return nil, err
	}
	return resp.WorkspaceList, nil
}

// CreateWorkspace creates a new workspace (team). Slug must be unique and URL-safe.
func (c *Client) CreateWorkspace(ctx context.Context, name, slug, description string) (*Workspace, error) {
	resp, err := createWorkspace(ctx, c.gql, name, slug, optStr(description))
	if err != nil {
		return nil, err
	}
	return &resp.WorkspaceCreate, nil
}

// DeleteWorkspace permanently deletes a workspace and all its resources.
func (c *Client) DeleteWorkspace(ctx context.Context, id string) error {
	_, err := deleteWorkspace(ctx, c.gql, id)
	return err
}

// ListWorkspaceMembers returns all members of a workspace.
func (c *Client) ListWorkspaceMembers(ctx context.Context, workspaceSlug string) ([]WorkspaceMember, error) {
	resp, err := listWorkspaceMembers(ctx, c.gql, workspaceSlug)
	if err != nil {
		return nil, err
	}
	return resp.WorkspaceListMembers, nil
}

// InviteToWorkspace invites a user (by email or username) to a workspace.
// Role defaults to "member" if empty.
func (c *Client) InviteToWorkspace(ctx context.Context, workspaceID, user, role string) (*WorkspaceInvite, error) {
	resp, err := inviteToWorkspace(ctx, c.gql, workspaceID, user, optStr(role))
	if err != nil {
		return nil, err
	}
	return &resp.WorkspaceInvite, nil
}

// RemoveWorkspaceMember removes a member from a workspace.
func (c *Client) RemoveWorkspaceMember(ctx context.Context, workspaceID, userID string) error {
	_, err := removeWorkspaceMember(ctx, c.gql, workspaceID, userID)
	return err
}

// ListMyInvites returns all pending workspace invitations for the authenticated user.
func (c *Client) ListMyInvites(ctx context.Context) ([]WorkspaceInvite, error) {
	resp, err := listMyInvites(ctx, c.gql)
	if err != nil {
		return nil, err
	}
	return resp.WorkspaceListMyInvites, nil
}

// ListWorkspaceInvites returns all pending invites for a workspace.
func (c *Client) ListWorkspaceInvites(ctx context.Context, workspaceSlug string) ([]WorkspaceInvite, error) {
	resp, err := listWorkspaceInvites(ctx, c.gql, workspaceSlug)
	if err != nil {
		return nil, err
	}
	return resp.WorkspaceListInvites, nil
}

// AcceptInvite accepts a workspace invitation.
func (c *Client) AcceptInvite(ctx context.Context, inviteID string) error {
	_, err := acceptInvite(ctx, c.gql, inviteID)
	return err
}

// DeclineInvite declines a workspace invitation.
func (c *Client) DeclineInvite(ctx context.Context, inviteID string) error {
	_, err := declineInvite(ctx, c.gql, inviteID)
	return err
}

// RevokeInvite cancels a pending workspace invitation.
func (c *Client) RevokeInvite(ctx context.Context, inviteID string) error {
	_, err := revokeInvite(ctx, c.gql, inviteID)
	return err
}
