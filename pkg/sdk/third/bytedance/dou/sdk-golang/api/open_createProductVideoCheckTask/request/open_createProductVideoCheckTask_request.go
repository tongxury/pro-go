package open_createProductVideoCheckTask_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/open_createProductVideoCheckTask/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenCreateProductVideoCheckTaskRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *OpenCreateProductVideoCheckTaskParam
}

func (c *OpenCreateProductVideoCheckTaskRequest) GetUrlPath() string {
	return "/open/createProductVideoCheckTask"
}

func New() *OpenCreateProductVideoCheckTaskRequest {
	request := &OpenCreateProductVideoCheckTaskRequest{
		Param: &OpenCreateProductVideoCheckTaskParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *OpenCreateProductVideoCheckTaskRequest) Execute(accessToken *doudian_sdk.AccessToken) (*open_createProductVideoCheckTask_response.OpenCreateProductVideoCheckTaskResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &open_createProductVideoCheckTask_response.OpenCreateProductVideoCheckTaskResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *OpenCreateProductVideoCheckTaskRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *OpenCreateProductVideoCheckTaskRequest) GetParams() *OpenCreateProductVideoCheckTaskParam {
	return c.Param
}

type OpenCreateProductVideoCheckTaskParam struct {
	// 商品id
	ProductId int64 `json:"product_id"`
}
