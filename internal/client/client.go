package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Skyvko6607/fastcrud/cli/internal/schema"
)

type Client struct {
	BaseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) Authenticate(accessKeyID string) (string, error) {
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/authenticate/crud/%s", c.BaseURL, accessKeyID),
		"application/json",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("authentication failed (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %w", err)
	}

	return result.AccessToken, nil
}

func (c *Client) FetchSchema(token string) ([]schema.Table, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/schema", c.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema request: %w", err)
	}
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("schema request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("schema fetch failed (%d): %s", resp.StatusCode, string(body))
	}

	var tables []schema.Table
	if err := json.NewDecoder(resp.Body).Decode(&tables); err != nil {
		return nil, fmt.Errorf("failed to parse schema response: %w", err)
	}

	return tables, nil
}
