package ink

import "context"

func (c *Client) ListDNSZones(ctx context.Context, workspaceSlug string) ([]DNSZone, error) {
	resp, err := listDNSZones(ctx, c.gql, optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return resp.DnsListZones, nil
}
func (c *Client) ListDNSRecords(ctx context.Context, zone, workspaceSlug string) ([]ZoneRecord, error) {
	resp, err := listDNSRecords(ctx, c.gql, zone, optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return resp.DnsListRecords, nil
}
func (c *Client) AddDNSRecord(ctx context.Context, zone, name, recordType, content string, ttl int, workspaceSlug string) (*ZoneRecord, error) {
	resp, err := addDNSRecord(ctx, c.gql, zone, name, recordType, content, optInt(ttl), optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.DnsAddRecord, nil
}
func (c *Client) DeleteDNSRecord(ctx context.Context, zone, recordID, workspaceSlug string) error {
	_, err := deleteDNSRecord(ctx, c.gql, zone, recordID, optStr(workspaceSlug))
	return err
}
