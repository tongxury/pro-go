package ape

type Client struct {
	tokenPoolsCache map[string][]Pool
}

func NewClient() *Client {
	return &Client{
		tokenPoolsCache: make(map[string][]Pool),
	}
}
