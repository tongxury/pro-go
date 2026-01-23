package dypenapi

import (
	"fmt"
	"github.com/bytedance/douyin-openapi-sdk-go/client"
)

func (t *Client) GetUserInfo(code string) (*client.OauthUserinfoResponseData, error) {

	at, err := t.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	fmt.Println("at", at)

	userinfo, err := t.c.OauthUserinfo(&client.OauthUserinfoRequest{
		AccessToken: at.AccessToken,
		OpenId:      at.OpenId,
	})

	fmt.Println("userinfo", userinfo, "err", err)
	if err != nil {
		return nil, err
	}

	return userinfo.Data, nil
}
