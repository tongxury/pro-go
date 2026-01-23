package shopVideo_getInfoForPublish_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGetInfoForPublishResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoGetInfoForPublishData `json:"data"`
}
type HotspotListItem struct {
	// 热点名称
	Name string `json:"name"`
}
type HotspotInfo struct {
	// 热点列表
	HotspotList []HotspotListItem `json:"hotspot_list"`
}
type ActivityListItem struct {
	// 活动头图
	HeadImage string `json:"head_image"`
	// 活动名称
	ActivityName string `json:"activity_name"`
	// 活动ID
	ActivityId string `json:"activity_id"`
	// 详情跳转地址
	JumpUrl string `json:"jump_url"`
}
type AllianceLimitInfo struct {
	// 剩余可挂车数量
	RemainBindNum int64 `json:"remain_bind_num"`
	// 是否可挂车
	CanBind bool `json:"can_bind"`
	// 当前挂车数量
	CurrentBindNum int64 `json:"current_bind_num"`
}
type ShopVideoGetInfoForPublishData struct {
	// 热点信息
	HotspotInfo *HotspotInfo `json:"hotspot_info"`
	// 活动列表
	ActivityList []ActivityListItem `json:"activity_list"`
	// 挂车次数限制信息
	AllianceLimitInfo *AllianceLimitInfo `json:"alliance_limit_info"`
}
