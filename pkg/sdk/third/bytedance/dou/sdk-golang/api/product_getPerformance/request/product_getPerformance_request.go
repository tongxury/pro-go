package product_getPerformance_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/product_getPerformance/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ProductGetPerformanceRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ProductGetPerformanceParam
}

func (c *ProductGetPerformanceRequest) GetUrlPath() string {
	return "/product/getPerformance"
}

func New() *ProductGetPerformanceRequest {
	request := &ProductGetPerformanceRequest{
		Param: &ProductGetPerformanceParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ProductGetPerformanceRequest) Execute(accessToken *doudian_sdk.AccessToken) (*product_getPerformance_response.ProductGetPerformanceResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &product_getPerformance_response.ProductGetPerformanceResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ProductGetPerformanceRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ProductGetPerformanceRequest) GetParams() *ProductGetPerformanceParam {
	return c.Param
}

type ProductGetPerformanceParam struct {
	// 商品ID
	ProductId string `json:"product_id"`
}
