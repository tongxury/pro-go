package cartesia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.cartesia.ai"
)

type Client struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type AccessTokenRequest struct {
	Grants    map[string]bool `json:"grants"`
	ExpiresIn int             `json:"expires_in"`
}

type AccessTokenResponse struct {
	Token string `json:"token"`
}

func (c *Client) GenerateAccessToken(ctx context.Context, grants map[string]bool, expiresIn int) (string, error) {
	url := fmt.Sprintf("%s/access-token", BaseURL)

	reqBody := AccessTokenRequest{
		Grants:    grants,
		ExpiresIn: expiresIn,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cartesia-Version", "2024-06-10")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cartesia api error: status %d", resp.StatusCode)
	}

	var res AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Token, nil
}
