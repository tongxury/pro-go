package cartesia

import (
	"context"
	"fmt"
	"time"
)

type Voice struct {
	Id            string    `json:"id"`
	Mode          string    `json:"mode"`
	IsPublic      bool      `json:"is_public"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	Gender        string    `json:"gender"`
	Embedding     []float64 `json:"embedding"`
	Language      string    `json:"language"`
	Popularity    int       `json:"popularity"`
	PreviewFileId string    `json:"preview_file_id"`
	OwnerId       string    `json:"owner_id,omitempty"`
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
