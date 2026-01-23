package wxdev

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func (t *Client) GetOpenIdByCode(ctx context.Context, code string) (string, error) {

	get, err := resty.New().R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"grant_type": "client_credential",
			"appid":      "wxb0008ff0fb65a736",
			"secret":     "04d8ed66b3a8db49550f1927501863dd",
			"js_code":    code,
		}).
		Post(t.endpoint + "/sns/jscode2session")

	if err != nil {
		return "", err
	}

	fmt.Println(get.String())

	var result OpenIdResult
	err = json.Unmarshal(get.Body(), &result)
	if err != nil {
		return "", err
	}

	return result.Openid, nil
}

type OpenIdResult struct {
	SessionKey string `json:"session_key"`
	Openid     string `json:"openid"`
}
