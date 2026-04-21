package dnsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CloudflareClient struct {
	APIToken string
	ZoneID   string
	BaseURL  string
}

func NewCloudflareClient(token, zoneID string) *CloudflareClient {
	return &CloudflareClient{
		APIToken: token,
		ZoneID:   zoneID,
		BaseURL:  fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID),
	}
}

// MapToTunnel creates or updates a CNAME record pointing to a specific Cloudflare Tunnel.
func (c *CloudflareClient) MapToTunnel(subdomain string, tunnelUUID string) error {
	recordID, err := c.getRecordID(subdomain)
	if err != nil {
		return err
	}

	// The target must be the Cloudflare Tunnel URL format
	target := fmt.Sprintf("%s.cfargotunnel.com", tunnelUUID)

	payload := map[string]interface{}{
		"type":    "CNAME",
		"name":    subdomain,
		"content": target,
		"ttl":     60,   // 60s allows for fast failover switches
		"proxied": true, // Must be proxied for Cloudflare tunnels to work securely
	}

	body, _ := json.Marshal(payload)

	var req *http.Request
	if recordID != "" {
		// Update existing record
		req, _ = http.NewRequest("PUT", fmt.Sprintf("%s/%s", c.BaseURL, recordID), bytes.NewBuffer(body))
	} else {
		// Create new record
		req, _ = http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(body))
	}

	req.Header.Add("Authorization", "Bearer "+c.APIToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("network error contacting cloudflare: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("cloudflare API returned status: %d", resp.StatusCode)
	}

	return nil
}

// getRecordID finds an existing CNAME record by name
func (c *CloudflareClient) getRecordID(subdomain string) (string, error) {
	url := fmt.Sprintf("%s?name=%s&type=CNAME", c.BaseURL, subdomain)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+c.APIToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var cfResp struct {
		Success bool `json:"success"`
		Result  []struct {
			ID string `json:"id"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return "", err
	}

	if cfResp.Success && len(cfResp.Result) > 0 {
		return cfResp.Result[0].ID, nil
	}

	return "", nil
}
