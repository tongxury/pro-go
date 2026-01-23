package second_ecomInsight_taskList_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/second_ecomInsight_taskList/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightTaskListRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *SecondEcomInsightTaskListParam
}

func (c *SecondEcomInsightTaskListRequest) GetUrlPath() string {
	return "/second/ecomInsight/taskList"
}

func New() *SecondEcomInsightTaskListRequest {
	request := &SecondEcomInsightTaskListRequest{
		Param: &SecondEcomInsightTaskListParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *SecondEcomInsightTaskListRequest) Execute(accessToken *doudian_sdk.AccessToken) (*second_ecomInsight_taskList_response.SecondEcomInsightTaskListResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &second_ecomInsight_taskList_response.SecondEcomInsightTaskListResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *SecondEcomInsightTaskListRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *SecondEcomInsightTaskListRequest) GetParams() *SecondEcomInsightTaskListParam {
	return c.Param
}

type SecondEcomInsightTaskListParam struct {
	// 条数
	PageSize *int32 `json:"page_size"`
	// 任务类型
	Type *string `json:"type"`
	// 任务状态：0 - 失败，1 - 创建，2 - 运行中，3 - 成功，4 - 取消
	Status *int32 `json:"status"`
	// 任务名
	Kw *string `json:"kw"`
	// 页码
	Page *int32 `json:"page"`
}
