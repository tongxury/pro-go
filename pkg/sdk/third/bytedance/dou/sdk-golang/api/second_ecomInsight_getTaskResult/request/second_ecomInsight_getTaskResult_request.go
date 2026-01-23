package second_ecomInsight_getTaskResult_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/second_ecomInsight_getTaskResult/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightGetTaskResultRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *SecondEcomInsightGetTaskResultParam
}

func (c *SecondEcomInsightGetTaskResultRequest) GetUrlPath() string {
	return "/second/ecomInsight/getTaskResult"
}

func New() *SecondEcomInsightGetTaskResultRequest {
	request := &SecondEcomInsightGetTaskResultRequest{
		Param: &SecondEcomInsightGetTaskResultParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *SecondEcomInsightGetTaskResultRequest) Execute(accessToken *doudian_sdk.AccessToken) (*second_ecomInsight_getTaskResult_response.SecondEcomInsightGetTaskResultResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &second_ecomInsight_getTaskResult_response.SecondEcomInsightGetTaskResultResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *SecondEcomInsightGetTaskResultRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *SecondEcomInsightGetTaskResultRequest) GetParams() *SecondEcomInsightGetTaskResultParam {
	return c.Param
}

type SecondEcomInsightGetTaskResultParam struct {
	// 任务id
	Id *int64 `json:"id"`
	// 任务类型
	TaskType *string `json:"task_type"`
}
