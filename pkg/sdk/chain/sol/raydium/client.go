package raydium

type Client struct {
	endpoint string
}

func NewClient() *Client {
	return &Client{
		endpoint: "https://transaction-v1.raydium.io",
	}
}
