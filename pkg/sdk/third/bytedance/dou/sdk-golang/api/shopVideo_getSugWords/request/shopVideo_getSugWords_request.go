package shopVideo_getSugWords_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_getSugWords/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGetSugWordsRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoGetSugWordsParam
}

func (c *ShopVideoGetSugWordsRequest) GetUrlPath() string {
	return "/shopVideo/getSugWords"
}

func New() *ShopVideoGetSugWordsRequest {
	request := &ShopVideoGetSugWordsRequest{
		Param: &ShopVideoGetSugWordsParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoGetSugWordsRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_getSugWords_response.ShopVideoGetSugWordsResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_getSugWords_response.ShopVideoGetSugWordsResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoGetSugWordsRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoGetSugWordsRequest) GetParams() *ShopVideoGetSugWordsParam {
	return c.Param
}

type ShopVideoGetSugWordsParam struct {
	// 话题词
	Keyword string `json:"keyword"`
}
