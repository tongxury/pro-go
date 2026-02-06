package gemini

import (
	"net/http"
	"store/pkg/sdk"

	"google.golang.org/genai"
)

type Client struct {
	c          *genai.Client
	httpClient *http.Client
	cache      sdk.ICache
}

func (t *Client) C() *genai.Client {
	return t.c
}

func (t *Client) HTTPClient() *http.Client {
	return t.httpClient
}
