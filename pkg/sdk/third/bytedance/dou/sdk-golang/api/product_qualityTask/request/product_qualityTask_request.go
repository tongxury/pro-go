package product_qualityTask_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/product_qualityTask/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ProductQualityTaskRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ProductQualityTaskParam
}

func (c *ProductQualityTaskRequest) GetUrlPath() string {
	return "/product/qualityTask"
}

func New() *ProductQualityTaskRequest {
	request := &ProductQualityTaskRequest{
		Param: &ProductQualityTaskParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ProductQualityTaskRequest) Execute(accessToken *doudian_sdk.AccessToken) (*product_qualityTask_response.ProductQualityTaskResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &product_qualityTask_response.ProductQualityTaskResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ProductQualityTaskRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ProductQualityTaskRequest) GetParams() *ProductQualityTaskParam {
	return c.Param
}

type ProductQualityTaskParam struct {
	// 是否只返回简要信息，不写默认false
	BriefOnly *bool `json:"brief_only"`
}
