package shopVideo_createShoppingCartAnchor_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_createShoppingCartAnchor/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoCreateShoppingCartAnchorRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoCreateShoppingCartAnchorParam
}

func (c *ShopVideoCreateShoppingCartAnchorRequest) GetUrlPath() string {
	return "/shopVideo/createShoppingCartAnchor"
}

func New() *ShopVideoCreateShoppingCartAnchorRequest {
	request := &ShopVideoCreateShoppingCartAnchorRequest{
		Param: &ShopVideoCreateShoppingCartAnchorParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoCreateShoppingCartAnchorRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_createShoppingCartAnchor_response.ShopVideoCreateShoppingCartAnchorResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_createShoppingCartAnchor_response.ShopVideoCreateShoppingCartAnchorResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoCreateShoppingCartAnchorRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoCreateShoppingCartAnchorRequest) GetParams() *ShopVideoCreateShoppingCartAnchorParam {
	return c.Param
}

type ShopVideoCreateShoppingCartAnchorParam struct {
	// 商品ID
	ProductId string `json:"product_id"`
	// 不传默认人店一体账号，支持渠道号
	UserId *string `json:"user_id"`
	// 小黄车标题
	Title *string `json:"title"`
}
