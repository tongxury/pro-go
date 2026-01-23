package arkr

import (
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
)

type Client struct {
	c *arkruntime.Client
}

func NewClient() *Client {

	client := arkruntime.NewClientWithApiKey(
		"c6d7518e-8174-4c94-96ca-e64c032bba07",
		//arkruntime.WithTimeout(10*time.Minute),
		//arkruntime.WithHTTPClient(&http.Client{
		//	Timeout: 10 * time.Minute,
		//}),
	)

	return &Client{
		c: client,
	}
}

func (t *Client) C() *arkruntime.Client {
	return t.c
}
