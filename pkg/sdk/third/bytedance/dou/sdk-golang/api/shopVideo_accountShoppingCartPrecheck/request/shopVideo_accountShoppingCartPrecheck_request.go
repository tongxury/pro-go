package shopVideo_accountShoppingCartPrecheck_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_accountShoppingCartPrecheck/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoAccountShoppingCartPrecheckRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoAccountShoppingCartPrecheckParam
}

func (c *ShopVideoAccountShoppingCartPrecheckRequest) GetUrlPath() string {
	return "/shopVideo/accountShoppingCartPrecheck"
}

func New() *ShopVideoAccountShoppingCartPrecheckRequest {
	request := &ShopVideoAccountShoppingCartPrecheckRequest{
		Param: &ShopVideoAccountShoppingCartPrecheckParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoAccountShoppingCartPrecheckRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_accountShoppingCartPrecheck_response.ShopVideoAccountShoppingCartPrecheckResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_accountShoppingCartPrecheck_response.ShopVideoAccountShoppingCartPrecheckResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoAccountShoppingCartPrecheckRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoAccountShoppingCartPrecheckRequest) GetParams() *ShopVideoAccountShoppingCartPrecheckParam {
	return c.Param
}

type ShopVideoAccountShoppingCartPrecheckParam struct {
	// 账号ID，支持获取授权号信息
	UserId *string `json:"user_id"`
}
