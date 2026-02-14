package cartesia

import (
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	BaseURL = "https://api.cartesia.ai"
)

type Client struct {
	apiKey string
	client *resty.Client
}

func NewClient(apiKey string) *Client {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetHeader("X-API-Key", apiKey)
	client.SetHeader("Cartesia-Version", "2025-04-16")

	return &Client{
		apiKey: apiKey,
		client: client,
	}
}
