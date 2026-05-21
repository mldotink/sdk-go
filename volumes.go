package ink

import "context"

func (c *Client) ListVolumes(ctx context.Context, workspaceSlug, projectSlug string) ([]VolumeInfo, error) {
	resp, err := listVolumes(ctx, c.gql, optStr(workspaceSlug), optStr(projectSlug))
	if err != nil {
		return nil, err
	}
	return resp.VolumeList, nil
}

func (c *Client) DeleteVolume(ctx context.Context, name, projectSlug, workspaceSlug string) error {
	_, err := deleteVolume(ctx, c.gql, name, optStr(projectSlug), optStr(workspaceSlug))
	return err
}
