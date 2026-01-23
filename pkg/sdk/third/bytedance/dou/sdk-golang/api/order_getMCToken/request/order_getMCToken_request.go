package order_getMCToken_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/order_getMCToken/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OrderGetMCTokenRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *OrderGetMCTokenParam
}

func (c *OrderGetMCTokenRequest) GetUrlPath() string {
	return "/order/getMCToken"
}

func New() *OrderGetMCTokenRequest {
	request := &OrderGetMCTokenRequest{
		Param: &OrderGetMCTokenParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *OrderGetMCTokenRequest) Execute(accessToken *doudian_sdk.AccessToken) (*order_getMCToken_response.OrderGetMCTokenResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &order_getMCToken_response.OrderGetMCTokenResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *OrderGetMCTokenRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *OrderGetMCTokenRequest) GetParams() *OrderGetMCTokenParam {
	return c.Param
}

type OrderGetMCTokenParam struct {
}
