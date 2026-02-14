package cartesia

import (
	"context"
	"fmt"
)

type Voice struct {
	Id          string `json:"id"`
	IsOwner     bool   `json:"is_owner"`
	IsPublic    bool   `json:"is_public"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	CreatedAt   string `json:"created_at"`
}

type ListVoicesRequest struct {
	Language      string `json:"language"`
	StartingAfter string `json:"starting_after"`
}

type ListVoicesResponse struct {
	Data     []Voice `json:"data"`
	HasMore  bool    `json:"has_more"`
	NextPage string  `json:"next_page"`
}

func (c *Client) ListVoices(ctx context.Context, in *ListVoicesRequest) (*ListVoicesResponse, error) {
	url := fmt.Sprintf("%s/voices", BaseURL)

	var res ListVoicesResponse
	req := c.client.R().
		SetContext(ctx).
		SetResult(&res)

	if in != nil {
		if in.Language != "" {
			req.SetQueryParam("language", in.Language)
		}
		if in.StartingAfter != "" {
			req.SetQueryParam("starting_after", in.StartingAfter)
		}
	}

	resp, err := req.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("cartesia api error: status %d body %s", resp.StatusCode(), resp.String())
	}

	return &res, nil
}
