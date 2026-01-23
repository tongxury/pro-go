package dou

type Client struct {
	conf Config
}

type Config struct {
	AppKey    string
	AppSecret string
}

func NewClient() *Client {
	return &Client{
		conf: Config{
			AppKey:    "7582142872741217828",
			AppSecret: "e27d9d6a-35b7-42db-8859-d6fb0f0c99d7",
		},
	}
}
