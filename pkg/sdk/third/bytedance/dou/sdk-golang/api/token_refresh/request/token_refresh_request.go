package token_refresh_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/token_refresh/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type TokenRefreshRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *TokenRefreshParam
}

func (c *TokenRefreshRequest) GetUrlPath() string {
	return "/token/refresh"
}

func New() *TokenRefreshRequest {
	request := &TokenRefreshRequest{
		Param: &TokenRefreshParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *TokenRefreshRequest) Execute(accessToken *doudian_sdk.AccessToken) (*token_refresh_response.TokenRefreshResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &token_refresh_response.TokenRefreshResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *TokenRefreshRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *TokenRefreshRequest) GetParams() *TokenRefreshParam {
	return c.Param
}

type TokenRefreshParam struct {
	// 用于刷新access_token的刷新令牌；有效期：14 天；
	RefreshToken string `json:"refresh_token"`
	// 授权类型；请传入默认值：refresh_token
	GrantType string `json:"grant_type"`
}
