package shopVideo_deleteShopVideo_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_deleteShopVideo/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoDeleteShopVideoRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoDeleteShopVideoParam
}

func (c *ShopVideoDeleteShopVideoRequest) GetUrlPath() string {
	return "/shopVideo/deleteShopVideo"
}

func New() *ShopVideoDeleteShopVideoRequest {
	request := &ShopVideoDeleteShopVideoRequest{
		Param: &ShopVideoDeleteShopVideoParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoDeleteShopVideoRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_deleteShopVideo_response.ShopVideoDeleteShopVideoResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_deleteShopVideo_response.ShopVideoDeleteShopVideoResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoDeleteShopVideoRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoDeleteShopVideoRequest) GetParams() *ShopVideoDeleteShopVideoParam {
	return c.Param
}

type ShopVideoDeleteShopVideoParam struct {
	// 视频ID
	VideoId string `json:"video_id"`
	// user_id，不传则默认人店一体账号，支持渠道号
	UserId *string `json:"user_id"`
}
