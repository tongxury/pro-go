package shopVideo_publishVideoNew_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_publishVideoNew/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoPublishVideoNewRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoPublishVideoNewParam
}

func (c *ShopVideoPublishVideoNewRequest) GetUrlPath() string {
	return "/shopVideo/publishVideoNew"
}

func New() *ShopVideoPublishVideoNewRequest {
	request := &ShopVideoPublishVideoNewRequest{
		Param: &ShopVideoPublishVideoNewParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoPublishVideoNewRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_publishVideoNew_response.ShopVideoPublishVideoNewResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_publishVideoNew_response.ShopVideoPublishVideoNewResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoPublishVideoNewRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoPublishVideoNewRequest) GetParams() *ShopVideoPublishVideoNewParam {
	return c.Param
}

type TopicsItem struct {
	// 活动名称
	TopicName *string `json:"topic_name"`
	// 话题id
	TopicId *string `json:"topic_id"`
}
type RelatedInfo struct {
	// 话题信息
	Topics []TopicsItem `json:"topics"`
	// 热点词
	HotSentence *string `json:"hot_sentence"`
}
type AppInfo struct {
	// app_key
	AppId int64 `json:"app_id"`
}
type ShopVideoPublishVideoNewParam struct {
	// 发布信息
	PublishInfo *PublishInfo `json:"publish_info"`
	// 视频信息
	VideoInfo *VideoInfo `json:"video_info"`
	// 锚点信息
	AnchorInfo *AnchorInfo `json:"anchor_info"`
	// 相关信息
	RelatedInfo *RelatedInfo `json:"related_info"`
	// 应用信息
	AppInfo *AppInfo `json:"app_info"`
}
type PublishInfo struct {
	// 来自open端默认传3
	Source int32 `json:"source"`
	// 用户id，不传则为人店一体账号
	UserId *string `json:"user_id"`
	// 定时发布时间
	SchedulePublishTime *int64 `json:"schedule_publish_time"`
}
type VideoInfo struct {
	// 视频描述
	Desc string `json:"desc"`
	// 封面地址。
	CoverImageUri *string `json:"cover_image_uri"`
	// 视频id，格式如：v0d27cg10001csve2qfog65o283c1h0g
	MediaId string `json:"media_id"`
	// 视频标题
	Title string `json:"title"`
}
type AnchorsItem struct {
	// 锚点标题，例如小黄车的标题
	Keyword *string `json:"keyword"`
	// 商品id，挂车和非挂车均需要
	ProductIds []string `json:"product_ids"`
	// 锚点草稿id。发布挂车视频时必传
	AnchorDraftId *string `json:"anchor_draft_id"`
	// 锚点文案
	Content *string `json:"content"`
	// 锚点类型，3-挂车，7-非挂车
	Type *int32 `json:"type"`
}
type AnchorInfo struct {
	// 锚点列表
	Anchors []AnchorsItem `json:"anchors"`
}
