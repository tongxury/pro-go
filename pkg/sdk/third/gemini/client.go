package gemini

import (
	"store/pkg/sdk"

	"google.golang.org/genai"
)

type Client struct {
	c     *genai.Client
	cache sdk.ICache
}

func (t *Client) C() *genai.Client {
	return t.c
}
