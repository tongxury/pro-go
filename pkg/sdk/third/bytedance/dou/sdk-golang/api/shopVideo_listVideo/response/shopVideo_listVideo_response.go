package shopVideo_listVideo_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoListVideoResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoListVideoData `json:"data"`
}
type ShopVideoListVideoData struct {
	// 视频列表
	VideoList []VideoListItem `json:"video_list"`
	// 总数
	Total int32 `json:"total"`
	// 当前页码
	CurrentPage int32 `json:"current_page"`
}
type VideoListItem struct {
	// 店铺id
	ShopId int64 `json:"shop_id"`
	// 抖音主端视频id
	ItemId int64 `json:"item_id"`
	// 抖音官方账号id
	UserId int64 `json:"user_id"`
	// 视频状态：-1-删除，0-有效
	Status int32 `json:"status"`
	// 播放数量
	PlayCount int32 `json:"play_count"`
	// 挂车商品id列表
	AllianceProductList []int64 `json:"alliance_product_list"`
	// 创建时间
	CreateTime int64 `json:"create_time"`
	// 更新时间
	ModifyTime int64 `json:"modify_time"`
	// 成交订单金额，单位分
	OrderAmount int32 `json:"order_amount"`
	// 绑定类型：0-未绑定，1-挂车，2-种草
	BindType int32 `json:"bind_type"`
	// 种草商品id列表
	SeedingProductList []int64 `json:"seeding_product_list"`
}
