package geminiai

import "google.golang.org/genai"

type Client struct {
	c *genai.Client
}

func (t *Client) C() *genai.Client {
	return t.c
}
