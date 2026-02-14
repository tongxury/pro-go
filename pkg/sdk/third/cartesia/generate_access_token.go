package cartesia

import (
	"context"
	"fmt"
)

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

	var res AccessTokenResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(reqBody).
		SetResult(&res).
		Post(url)

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("cartesia api error: status %d body %s", resp.StatusCode(), resp.String())
	}

	return res.Token, nil
}
