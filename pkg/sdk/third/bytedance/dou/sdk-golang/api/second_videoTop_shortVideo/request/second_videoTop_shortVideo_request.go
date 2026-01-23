package second_videoTop_shortVideo_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/second_videoTop_shortVideo/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondVideoTopShortVideoRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *SecondVideoTopShortVideoParam
}

func (c *SecondVideoTopShortVideoRequest) GetUrlPath() string {
	return "/second/videoTop/shortVideo"
}

func New() *SecondVideoTopShortVideoRequest {
	request := &SecondVideoTopShortVideoRequest{
		Param: &SecondVideoTopShortVideoParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *SecondVideoTopShortVideoRequest) Execute(accessToken *doudian_sdk.AccessToken) (*second_videoTop_shortVideo_response.SecondVideoTopShortVideoResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &second_videoTop_shortVideo_response.SecondVideoTopShortVideoResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *SecondVideoTopShortVideoRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *SecondVideoTopShortVideoRequest) GetParams() *SecondVideoTopShortVideoParam {
	return c.Param
}

type Sort struct {
	// 排序方式 asc - 升序， desc - 降序
	Order *string `json:"order"`
	// 排序字段，可使用字段参考https://bytedance.sg.larkoffice.com/docx/CdM7dD8FVoub3VxpnFdlPk7zglh
	Field *string `json:"field"`
}
type SecondVideoTopShortVideoParam struct {
	// 搜索条件
	Likes []LikesItem `json:"likes"`
	// 页码
	Page *int32 `json:"page"`
	// 条数
	PageSize *int32 `json:"page_size"`
	// 过滤条件
	Filters []FiltersItem `json:"filters"`
	// 排序条件
	Sort *Sort `json:"sort"`
}
type LikesItem struct {
	// 搜索方式
	LikeValue *string `json:"like_value"`
	// 搜索字段，使用方式参考https://bytedance.sg.larkoffice.com/docx/CdM7dD8FVoub3VxpnFdlPk7zglh
	Field *string `json:"field"`
}
type FiltersItem struct {
	// 过滤字段，可使用字段参考https://bytedance.sg.larkoffice.com/docx/CdM7dD8FVoub3VxpnFdlPk7zglh
	Field *string `json:"field"`
	// 过滤方式，使用方式参考https://bytedance.sg.larkoffice.com/docx/CdM7dD8FVoub3VxpnFdlPk7zglh
	Option *string `json:"option"`
}
