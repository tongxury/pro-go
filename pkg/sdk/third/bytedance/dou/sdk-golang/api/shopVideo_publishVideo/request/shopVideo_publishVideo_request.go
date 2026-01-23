package shopVideo_publishVideo_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_publishVideo/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoPublishVideoRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoPublishVideoParam
}

func (c *ShopVideoPublishVideoRequest) GetUrlPath() string {
	return "/shopVideo/publishVideo"
}

func New() *ShopVideoPublishVideoRequest {
	request := &ShopVideoPublishVideoRequest{
		Param: &ShopVideoPublishVideoParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoPublishVideoRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_publishVideo_response.ShopVideoPublishVideoResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_publishVideo_response.ShopVideoPublishVideoResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoPublishVideoRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoPublishVideoRequest) GetParams() *ShopVideoPublishVideoParam {
	return c.Param
}

type AnchorListItem struct {
	// 锚点类型：3-国内电商（挂车），7-电商种草
	Type *int32 `json:"type"`
	// 锚点关联商品id，当前只支持关联一个商品
	ProductIds []int64 `json:"product_ids"`
}
type ActivityListItem struct {
	// 活动id
	ActivityId *string `json:"activity_id"`
	// 活动名称
	ActivityName *string `json:"activity_name"`
}
type BasicInfo struct {
	// 视频id
	MediaId string `json:"media_id"`
	// 描述
	Desc *string `json:"desc"`
	// 标题
	Title *string `json:"title"`
	// 活动列表
	ActivityList []ActivityListItem `json:"activity_list"`
	// 封面uri，通过sdk上传封面素材后获取。上传方式见：https://bytedance.larkoffice.com/docx/doxcnyCpq3C9gik15T2dtKK0Tlc?source_type=message&from=message
	CoverImageUri *string `json:"cover_image_uri"`
	// 热点
	Hotspot *string `json:"hotspot"`
}
type ShopVideoPublishVideoParam struct {
	// 短视频锚点（关联商品）
	AnchorList []AnchorListItem `json:"anchor_list"`
	// 短视频基础信息
	BasicInfo *BasicInfo `json:"basic_info"`
}
