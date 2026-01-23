package shopVideo_uploadVideo_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_uploadVideo/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoUploadVideoRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoUploadVideoParam
}

func (c *ShopVideoUploadVideoRequest) GetUrlPath() string {
	return "/shopVideo/uploadVideo"
}

func New() *ShopVideoUploadVideoRequest {
	request := &ShopVideoUploadVideoRequest{
		Param: &ShopVideoUploadVideoParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoUploadVideoRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_uploadVideo_response.ShopVideoUploadVideoResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_uploadVideo_response.ShopVideoUploadVideoResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoUploadVideoRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoUploadVideoRequest) GetParams() *ShopVideoUploadVideoParam {
	return c.Param
}

type ShopVideoUploadVideoParam struct {
	// 视频url，目前支持mp4视频格式
	Url string `json:"url"`
}
