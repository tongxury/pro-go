package gmgnai

type Client struct {
	endpoint string
}

func NewClient() *Client {
	return &Client{}
}

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
