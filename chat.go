package ink

import "context"

// SendChatMessage posts a message to a workspace chat channel.
// channel defaults to the workspace general channel if empty.
func (c *Client) SendChatMessage(ctx context.Context, workspaceSlug, channel, content string) (*SendChatResult, error) {
	resp, err := sendChatMessage(ctx, c.gql, workspaceSlug, optStr(channel), content)
	if err != nil {
		return nil, err
	}
	return &resp.ChatSend, nil
}

// ReadChat reads messages from a workspace chat channel.
// cursor=0 reads from the beginning; use ReadChatResult.NextCursor to paginate.
// limit=0 uses the server default.
func (c *Client) ReadChat(ctx context.Context, workspaceSlug, channel string, cursor, limit int) (*ReadChatResult, error) {
	resp, err := readChat(ctx, c.gql, workspaceSlug, optStr(channel), optInt(cursor), optInt(limit))
	if err != nil {
		return nil, err
	}
	return &resp.ChatRead, nil
}
