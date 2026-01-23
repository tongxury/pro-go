package second_ecomInsight_deleteTask_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/second_ecomInsight_deleteTask/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightDeleteTaskRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *SecondEcomInsightDeleteTaskParam
}

func (c *SecondEcomInsightDeleteTaskRequest) GetUrlPath() string {
	return "/second/ecomInsight/deleteTask"
}

func New() *SecondEcomInsightDeleteTaskRequest {
	request := &SecondEcomInsightDeleteTaskRequest{
		Param: &SecondEcomInsightDeleteTaskParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *SecondEcomInsightDeleteTaskRequest) Execute(accessToken *doudian_sdk.AccessToken) (*second_ecomInsight_deleteTask_response.SecondEcomInsightDeleteTaskResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &second_ecomInsight_deleteTask_response.SecondEcomInsightDeleteTaskResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *SecondEcomInsightDeleteTaskRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *SecondEcomInsightDeleteTaskRequest) GetParams() *SecondEcomInsightDeleteTaskParam {
	return c.Param
}

type SecondEcomInsightDeleteTaskParam struct {
	// 任务id
	Id *int64 `json:"id"`
}
