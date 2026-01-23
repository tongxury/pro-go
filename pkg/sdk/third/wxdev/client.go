package wxdev

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	endpoint string
}

func NewClient() *Client {
	return &Client{
		endpoint: "https://api.weixin.qq.com",
	}
}

func (t *Client) GetAccessToken() (*Token, error) {

	get, err := resty.New().R().
		SetQueryParams(map[string]string{
			"grant_type": "client_credential",
			"appid":      "wxb0008ff0fb65a736",
			"secret":     "04d8ed66b3a8db49550f1927501863dd",
		}).
		Get(t.endpoint + "/cgi-bin/token")

	if err != nil {
		return nil, err
	}

	var token Token

	err = json.Unmarshal(get.Body(), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
