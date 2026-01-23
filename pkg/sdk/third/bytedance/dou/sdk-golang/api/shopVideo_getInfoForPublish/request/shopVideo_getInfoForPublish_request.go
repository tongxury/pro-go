package shopVideo_getInfoForPublish_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_getInfoForPublish/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGetInfoForPublishRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoGetInfoForPublishParam
}

func (c *ShopVideoGetInfoForPublishRequest) GetUrlPath() string {
	return "/shopVideo/getInfoForPublish"
}

func New() *ShopVideoGetInfoForPublishRequest {
	request := &ShopVideoGetInfoForPublishRequest{
		Param: &ShopVideoGetInfoForPublishParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoGetInfoForPublishRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_getInfoForPublish_response.ShopVideoGetInfoForPublishResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_getInfoForPublish_response.ShopVideoGetInfoForPublishResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoGetInfoForPublishRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoGetInfoForPublishRequest) GetParams() *ShopVideoGetInfoForPublishParam {
	return c.Param
}

type ShopVideoGetInfoForPublishParam struct {
	// 是否需要热点数据
	NeedHotspot bool `json:"need_hotspot"`
	// 是否查询创作者活动
	NeedActivity bool `json:"need_activity"`
	// 是否需要检查挂车次数限制
	NeedAllianceLimit bool `json:"need_alliance_limit"`
}
