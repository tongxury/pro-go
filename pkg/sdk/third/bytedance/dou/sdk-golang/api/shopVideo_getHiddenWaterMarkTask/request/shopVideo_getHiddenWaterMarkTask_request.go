package shopVideo_getHiddenWaterMarkTask_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_getHiddenWaterMarkTask/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGetHiddenWaterMarkTaskRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoGetHiddenWaterMarkTaskParam
}

func (c *ShopVideoGetHiddenWaterMarkTaskRequest) GetUrlPath() string {
	return "/shopVideo/getHiddenWaterMarkTask"
}

func New() *ShopVideoGetHiddenWaterMarkTaskRequest {
	request := &ShopVideoGetHiddenWaterMarkTaskRequest{
		Param: &ShopVideoGetHiddenWaterMarkTaskParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoGetHiddenWaterMarkTaskRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_getHiddenWaterMarkTask_response.ShopVideoGetHiddenWaterMarkTaskResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_getHiddenWaterMarkTask_response.ShopVideoGetHiddenWaterMarkTaskResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoGetHiddenWaterMarkTaskRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoGetHiddenWaterMarkTaskRequest) GetParams() *ShopVideoGetHiddenWaterMarkTaskParam {
	return c.Param
}

type ShopVideoGetHiddenWaterMarkTaskParam struct {
	// 任务id
	TaskId string `json:"task_id"`
}
