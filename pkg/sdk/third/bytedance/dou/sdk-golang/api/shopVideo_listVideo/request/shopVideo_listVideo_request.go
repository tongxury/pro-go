package shopVideo_listVideo_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_listVideo/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoListVideoRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoListVideoParam
}

func (c *ShopVideoListVideoRequest) GetUrlPath() string {
	return "/shopVideo/listVideo"
}

func New() *ShopVideoListVideoRequest {
	request := &ShopVideoListVideoRequest{
		Param: &ShopVideoListVideoParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoListVideoRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_listVideo_response.ShopVideoListVideoResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_listVideo_response.ShopVideoListVideoResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoListVideoRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoListVideoRequest) GetParams() *ShopVideoListVideoParam {
	return c.Param
}

type ShopVideoListVideoParam struct {
	// 视频上传开始时间戳，支持查询90天内视频
	CreateTimeStart *int64 `json:"create_time_start"`
	// 页码，>1
	Page int32 `json:"page"`
	// 每页数量，1至50
	Size int32 `json:"size"`
	// 排序类型：0-降序，1-升序
	OrderByType *int32 `json:"order_by_type"`
	// 抖音主端短视频ID
	ItemId *int64 `json:"item_id"`
	// 排序字段，1-创建时间，2-累计成交量，3-播放量
	OrderByField *int32 `json:"order_by_field"`
	// 关联的商品id
	BindProductId *int64 `json:"bind_product_id"`
	// 绑定类型：0-未绑定，1-挂车，2-种草
	BindTypeList []int32 `json:"bind_type_list"`
}
