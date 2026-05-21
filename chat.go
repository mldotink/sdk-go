package ink

import "context"

func (c *Client) SendChatMessage(ctx context.Context, workspaceSlug, channel, content string) (*SendChatResult, error) {
	resp, err := sendChatMessage(ctx, c.gql, workspaceSlug, optStr(channel), content)
	if err != nil {
		return nil, err
	}
	return &resp.ChatSend, nil
}
func (c *Client) ReadChat(ctx context.Context, workspaceSlug, channel string, cursor, limit int) (*ReadChatResult, error) {
	resp, err := readChat(ctx, c.gql, workspaceSlug, optStr(channel), optInt(cursor), optInt(limit))
	if err != nil {
		return nil, err
	}
	return &resp.ChatRead, nil
}
