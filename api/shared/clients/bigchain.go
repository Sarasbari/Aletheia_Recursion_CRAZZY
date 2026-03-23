package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"aletheia-api/shared/models"
)

type BigchainClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewBigchainClient(baseURL, apiKey string) *BigchainClient {
	return &BigchainClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *BigchainClient) Save(ctx context.Context, record models.ProvenRecord) error {
	if c.baseURL == "" {
		return nil
	}
	body, err := json.Marshal(record)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/metadata", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bigchaindb save failed: %d", resp.StatusCode)
	}
	return nil
}

func (c *BigchainClient) FindBySHA(ctx context.Context, sha string) (*models.ProvenRecord, error) {
	if c.baseURL == "" {
		return nil, nil
	}
	url := fmt.Sprintf("%s/metadata/%s", c.baseURL, sha)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bigchaindb find failed: %d", resp.StatusCode)
	}
	var out models.ProvenRecord
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
