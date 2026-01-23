package wxdev

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func (t *Client) GetUserPhoneNumber(ctx context.Context, code string) (string, error) {

	accessToken, err := t.GetAccessToken()
	if err != nil {
		return "", err
	}

	get, err := resty.New().R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"access_token": accessToken.AccessToken,
		}).
		SetBody(map[string]string{
			"code": code,
		}).
		Post(t.endpoint + "/wxa/business/getuserphonenumber")

	if err != nil {
		return "", err
	}

	fmt.Println(get.String())

	var token Result

	err = json.Unmarshal(get.Body(), &token)
	if err != nil {
		return "", err
	}

	return token.PhoneInfo.PhoneNumber, nil
}

type Result struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int    `json:"timestamp"`
			Appid     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}
