package dnsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// handles interactions with the Cloudflare DNS API
type CloudflareClient struct {
	APIToken string
	ZoneID   string // the id of root domain
	BaseURL  string
}

func NewCloudflareClient(token, zoneID string) *CloudflareClient {
	return &CloudflareClient{
		APIToken: token,
		ZoneID:   zoneID,
		BaseURL:  fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID),
	}
}

// Internal Structs for parsing Cloudflare JSON responses
type cfResponse struct {
	Success bool       `json:"success"`
	Result  []cfRecord `json:"result"`
}

type cfRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (c *CloudflareClient) UpsertRecord(subdomain, ipAddress string) error {
	// Full domain name (app-name.domain.com)
	// For the MVP, we assume the subdomain string passed already includes the root domain

	// check if the record already exists
	recordID, err := c.getRecordID(subdomain)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"type":    "A",
		"name":    subdomain,
		"content": ipAddress,
		"ttl":     60,    // 60 seconds is the minimum TTL, crucial for fast failover
		"proxied": false, // TODO: setup tls on server and cloudflare both ends
	}

	body, _ := json.Marshal(payload)

	var req *http.Request
	if recordID != "" {
		// Record exists: Send PUT request to update it
		req, _ = http.NewRequest("PUT", fmt.Sprintf("%s/%s", c.BaseURL, recordID), bytes.NewBuffer(body))
	} else {
		// Record does not exist: Send POST request to create it
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

// queries Cloudflare to find an existing DNS record by name
func (c *CloudflareClient) getRecordID(subdomain string) (string, error) {
	// Query param to filter by name
	url := fmt.Sprintf("%s?name=%s&type=A", c.BaseURL, subdomain)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+c.APIToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var cfResp cfResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return "", err
	}

	// If a record is found, return its ID
	if cfResp.Success && len(cfResp.Result) > 0 {
		return cfResp.Result[0].ID, nil
	}

	// Return empty string if no record exists
	return "", nil
}
