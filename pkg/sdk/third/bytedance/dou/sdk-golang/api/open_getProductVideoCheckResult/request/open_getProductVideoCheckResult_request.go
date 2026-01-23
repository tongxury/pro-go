package open_getProductVideoCheckResult_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/open_getProductVideoCheckResult/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenGetProductVideoCheckResultRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *OpenGetProductVideoCheckResultParam
}

func (c *OpenGetProductVideoCheckResultRequest) GetUrlPath() string {
	return "/open/getProductVideoCheckResult"
}

func New() *OpenGetProductVideoCheckResultRequest {
	request := &OpenGetProductVideoCheckResultRequest{
		Param: &OpenGetProductVideoCheckResultParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *OpenGetProductVideoCheckResultRequest) Execute(accessToken *doudian_sdk.AccessToken) (*open_getProductVideoCheckResult_response.OpenGetProductVideoCheckResultResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &open_getProductVideoCheckResult_response.OpenGetProductVideoCheckResultResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *OpenGetProductVideoCheckResultRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *OpenGetProductVideoCheckResultRequest) GetParams() *OpenGetProductVideoCheckResultParam {
	return c.Param
}

type OpenGetProductVideoCheckResultParam struct {
	// 任务id
	TaskId string `json:"task_id"`
}
