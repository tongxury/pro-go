package dypenapi

import (
	"errors"
	credential "github.com/bytedance/douyin-openapi-credential-go/client"
	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
)

type Client struct {
	endpoint string
	c        *openApiSdkClient.Client
	conf     Config
}

type Config struct {
	ClientKey    string
	ClientSecret string
}

func NewClient(conf Config) (*Client, error) {

	opt := new(credential.Config).
		SetClientKey(conf.ClientKey).
		SetClientSecret(conf.ClientSecret)

	sdkClient, err := openApiSdkClient.NewClient(opt)
	if err != nil {
		return nil, err
	}

	return &Client{
		endpoint: "https://open.douyin.com",
		c:        sdkClient,
		conf:     conf,
	}, nil
}

func (t *Client) getAccessToken(code string) (*openApiSdkClient.OauthAccessTokenResponseData, error) {

	sdkRequest := &openApiSdkClient.OauthAccessTokenRequest{}
	sdkRequest.SetClientKey(t.conf.ClientKey)
	sdkRequest.SetClientSecret(t.conf.ClientSecret)
	sdkRequest.SetCode(code)
	sdkRequest.SetGrantType("authorization_code")
	// sdk调用
	sdkResponse, err := t.c.OauthAccessToken(sdkRequest)
	if err != nil {
		return nil, err
	}

	if sdkResponse.Data.ErrorCode != nil && *sdkResponse.Data.ErrorCode > 0 {
		return nil, errors.New(*sdkResponse.Data.Description)
	}

	return sdkResponse.Data, nil
}

//{
//"message": "success",
//"data": {
//"error_code": 0,
//"description": "",
//"expires_in": 1296000,
//"open_id": "_000bzTjnd9AQ8_t2bPsS4HURrMvM_6gm7Tq",
//"refresh_token": "rft.2b2e1d6ae8d2660f7ab21e219b5ee4d47Mm3S7bQSpfO7nEH1cbrurTbj7Sb_hl",
//"log_id": "2025081711553849F0E7DCB39A22629E93",
//"scope": "trial.whitelist",
//"refresh_expires_in": 2592000,
//"access_token": "act.3.qppASmiIOCv7pC_yHt009k9YnMGsSWna7BlAE_8ckWfj0vtQAAqUj3zb1YSe0XKArEsGuTkJKDWe9cgeRjMTE_LCfY_-RJ01-hdS5y2-kSYEFCMYXous-pfcq-ErAbrmlsA-wSbIh5C4oK3to-urcgUozSl18nj4G_BN3g==_hl"
//}
//}
