package shopVideo_searchItem_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_searchItem/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoSearchItemRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoSearchItemParam
}

func (c *ShopVideoSearchItemRequest) GetUrlPath() string {
	return "/shopVideo/searchItem"
}

func New() *ShopVideoSearchItemRequest {
	request := &ShopVideoSearchItemRequest{
		Param: &ShopVideoSearchItemParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoSearchItemRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_searchItem_response.ShopVideoSearchItemResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_searchItem_response.ShopVideoSearchItemResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoSearchItemRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoSearchItemRequest) GetParams() *ShopVideoSearchItemParam {
	return c.Param
}

type ShopVideoSearchItemParam struct {
	// 页码
	Page *int32 `json:"page"`
	// 绑定场景，1-只查挂车短视频，2-查询所有非挂车短视频，3-查询所有未绑定商品的短视频
	BindScene *int32 `json:"bind_scene"`
	// 可选参数: 是否获取相关数据
	OptionalParam *OptionalParam `json:"optional_param"`
	// 每页数量，默认10
	PageSize *int32 `json:"page_size"`
	// 作者账号id，不传默认人店一体账号，支持渠道号
	UserId *string `json:"user_id"`
	// 短视频id
	ItemId *string `json:"item_id"`
	// 商品参数
	Product *Product `json:"product"`
}
type OptionalParam struct {
	// 是否需要短视频表现
	NeedStats *bool `json:"need_stats"`
	// 是否需要优质标签
	NeedHighQualityTag *bool `json:"need_high_quality_tag"`
}
type Product struct {
	// 商品id
	ProductId *string `json:"product_id"`
}
