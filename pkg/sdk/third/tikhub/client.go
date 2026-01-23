package tikhub

import "github.com/go-resty/resty/v2"

type Client struct {
	apiKey  string
	baseUrl string
	c       *resty.Client
}

func NewClient() *Client {
	return &Client{
		//baseUrl: "https://api.tikhub.io",
		baseUrl: "https://api.tikhub.dev",
		//apiKey:  "zIEVE3F3/XQwrywxjEmviVtt9D0RtM3UoLJ48WzokGxo1HdSfqdcmiPvFA==",
		apiKey: "h8I4PZMa4X+srzTMNtiOn0jxDR24CqJJwCfvtme+aNO11g6+OkSYaCOSKw==",
		c:      resty.New(),
		//SetProxy("http://proxy:strOngPAssWOrd@aa404deaba3e54490a68959040f22566-685580509.ap-southeast-1.elb.amazonaws.com:6060"),
	}
}

type Response[T any] struct {
	Code   int    `json:"code"`
	Router string `json:"router"`
	Params struct {
	} `json:"params"`
	Data T `json:"data"`
}
