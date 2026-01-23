package shopVideo_genUploadAuth_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_genUploadAuth/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGenUploadAuthRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoGenUploadAuthParam
}

func (c *ShopVideoGenUploadAuthRequest) GetUrlPath() string {
	return "/shopVideo/genUploadAuth"
}

func New() *ShopVideoGenUploadAuthRequest {
	request := &ShopVideoGenUploadAuthRequest{
		Param: &ShopVideoGenUploadAuthParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoGenUploadAuthRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_genUploadAuth_response.ShopVideoGenUploadAuthResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_genUploadAuth_response.ShopVideoGenUploadAuthResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoGenUploadAuthRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoGenUploadAuthRequest) GetParams() *ShopVideoGenUploadAuthParam {
	return c.Param
}

type ShopVideoGenUploadAuthParam struct {
	// 上传场景：1-短视频发稿
	UploadScene int32 `json:"upload_scene"`
}
