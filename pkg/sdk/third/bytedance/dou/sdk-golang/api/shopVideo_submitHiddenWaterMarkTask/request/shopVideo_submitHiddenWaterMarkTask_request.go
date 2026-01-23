package shopVideo_submitHiddenWaterMarkTask_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_submitHiddenWaterMarkTask/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoSubmitHiddenWaterMarkTaskRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoSubmitHiddenWaterMarkTaskParam
}

func (c *ShopVideoSubmitHiddenWaterMarkTaskRequest) GetUrlPath() string {
	return "/shopVideo/submitHiddenWaterMarkTask"
}

func New() *ShopVideoSubmitHiddenWaterMarkTaskRequest {
	request := &ShopVideoSubmitHiddenWaterMarkTaskRequest{
		Param: &ShopVideoSubmitHiddenWaterMarkTaskParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoSubmitHiddenWaterMarkTaskRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_submitHiddenWaterMarkTask_response.ShopVideoSubmitHiddenWaterMarkTaskResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_submitHiddenWaterMarkTask_response.ShopVideoSubmitHiddenWaterMarkTaskResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoSubmitHiddenWaterMarkTaskRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoSubmitHiddenWaterMarkTaskRequest) GetParams() *ShopVideoSubmitHiddenWaterMarkTaskParam {
	return c.Param
}

type ShopVideoSubmitHiddenWaterMarkTaskParam struct {
	// 通过sdk上传后的获取到的mediaId，url和mediaId必须传一个
	MediaId *string `json:"mediaId"`
	// 视频url，必须公网可访问，url和mediaId必须传一个，优先mediaId
	Url *string `json:"url"`
}
