package shopVideo_queryUploadTask_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_queryUploadTask/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoQueryUploadTaskRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoQueryUploadTaskParam
}

func (c *ShopVideoQueryUploadTaskRequest) GetUrlPath() string {
	return "/shopVideo/queryUploadTask"
}

func New() *ShopVideoQueryUploadTaskRequest {
	request := &ShopVideoQueryUploadTaskRequest{
		Param: &ShopVideoQueryUploadTaskParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoQueryUploadTaskRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_queryUploadTask_response.ShopVideoQueryUploadTaskResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_queryUploadTask_response.ShopVideoQueryUploadTaskResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoQueryUploadTaskRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoQueryUploadTaskRequest) GetParams() *ShopVideoQueryUploadTaskParam {
	return c.Param
}

type ShopVideoQueryUploadTaskParam struct {
	// uploadVideo接口得到的task_id
	TaskId string `json:"task_id"`
}
