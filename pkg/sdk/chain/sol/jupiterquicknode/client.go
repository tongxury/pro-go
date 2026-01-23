package jupiterquicknode

import (
	"context"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	endpoint string
}

func NewClient() *Client {
	return &Client{
		endpoint: "https://public.jupiterapi.com",
	}
}

func (t *Client) GetTokenMetadata(ctx context.Context) error {

	resty.New().R().SetContext(ctx).Get(t.endpoint + "/tokens")
	return nil
}
