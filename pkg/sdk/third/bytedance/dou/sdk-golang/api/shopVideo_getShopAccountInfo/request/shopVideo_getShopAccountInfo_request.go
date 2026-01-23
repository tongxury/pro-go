package shopVideo_getShopAccountInfo_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_getShopAccountInfo/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGetShopAccountInfoRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoGetShopAccountInfoParam
}

func (c *ShopVideoGetShopAccountInfoRequest) GetUrlPath() string {
	return "/shopVideo/getShopAccountInfo"
}

func New() *ShopVideoGetShopAccountInfoRequest {
	request := &ShopVideoGetShopAccountInfoRequest{
		Param: &ShopVideoGetShopAccountInfoParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoGetShopAccountInfoRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_getShopAccountInfo_response.ShopVideoGetShopAccountInfoResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_getShopAccountInfo_response.ShopVideoGetShopAccountInfoResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoGetShopAccountInfoRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoGetShopAccountInfoRequest) GetParams() *ShopVideoGetShopAccountInfoParam {
	return c.Param
}

type ShopVideoGetShopAccountInfoParam struct {
}
