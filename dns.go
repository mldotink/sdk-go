package ink

import "context"

// ListDNSZones returns all DNS zones managed in a workspace.
func (c *Client) ListDNSZones(ctx context.Context, workspaceSlug string) ([]DNSZone, error) {
	resp, err := listDNSZones(ctx, c.gql, optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return resp.DnsListZones, nil
}

// ListDNSRecords returns all records in a DNS zone.
func (c *Client) ListDNSRecords(ctx context.Context, zone, workspaceSlug string) ([]ZoneRecord, error) {
	resp, err := listDNSRecords(ctx, c.gql, zone, optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return resp.DnsListRecords, nil
}

// AddDNSRecord creates a DNS record in a zone. ttl=0 uses the server default.
func (c *Client) AddDNSRecord(ctx context.Context, zone, name, recordType, content string, ttl int, workspaceSlug string) (*ZoneRecord, error) {
	resp, err := addDNSRecord(ctx, c.gql, zone, name, recordType, content, optInt(ttl), optStr(workspaceSlug))
	if err != nil {
		return nil, err
	}
	return &resp.DnsAddRecord, nil
}

// DeleteDNSRecord removes a DNS record from a zone by record ID.
func (c *Client) DeleteDNSRecord(ctx context.Context, zone, recordID, workspaceSlug string) error {
	_, err := deleteDNSRecord(ctx, c.gql, zone, recordID, optStr(workspaceSlug))
	return err
}
