package second_ecomInsight_taskDetail_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/second_ecomInsight_taskDetail/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightTaskDetailRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *SecondEcomInsightTaskDetailParam
}

func (c *SecondEcomInsightTaskDetailRequest) GetUrlPath() string {
	return "/second/ecomInsight/taskDetail"
}

func New() *SecondEcomInsightTaskDetailRequest {
	request := &SecondEcomInsightTaskDetailRequest{
		Param: &SecondEcomInsightTaskDetailParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *SecondEcomInsightTaskDetailRequest) Execute(accessToken *doudian_sdk.AccessToken) (*second_ecomInsight_taskDetail_response.SecondEcomInsightTaskDetailResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &second_ecomInsight_taskDetail_response.SecondEcomInsightTaskDetailResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *SecondEcomInsightTaskDetailRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *SecondEcomInsightTaskDetailRequest) GetParams() *SecondEcomInsightTaskDetailParam {
	return c.Param
}

type SecondEcomInsightTaskDetailParam struct {
	// 任务id
	Id *int64 `json:"id"`
}
