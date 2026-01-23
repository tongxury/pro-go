package shopVideo_searchItem_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoSearchItemResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoSearchItemData `json:"data"`
}
type HighQualityTag struct {
	// 是否优质
	HighQuality bool `json:"high_quality"`
	// 流量扶持
	BoostVV int64 `json:"boost_v_v"`
}
type AccountInfo struct {
	// 账号id
	UserId string `json:"user_id"`
	// 账号名称
	Name string `json:"name"`
}
type Data struct {
	// 视频列表
	Items []ItemsItem `json:"items"`
}
type Stats struct {
	// 成交金额,单位分
	OrderAmount int32 `json:"order_amount"`
	// 播放量
	PlayCount int32 `json:"play_count"`
	// 分享数
	ShareCount int32 `json:"share_count"`
	// 点赞数
	LikeCount int32 `json:"like_count"`
	// 评论数
	CommentCount int32 `json:"comment_count"`
}
type BindProductsItem struct {
	// 商品列表
	Products []ProductsItem `json:"products"`
	// 绑定关系,0-未绑定，1-挂车，2-种草，3-看后搜
	BindType int32 `json:"bind_type"`
}
type ItemsItem struct {
	// 店铺信息
	ShopInfo *ShopInfo `json:"shop_info"`
	// 视频信息
	ItemInfo *ItemInfo `json:"item_info"`
	// 绑定关系，0-未绑定，1-挂车，2-种草
	BindTypes []int64 `json:"bind_types"`
	// 短视频表现
	Stats *Stats `json:"stats"`
	// 绑定商品信息。BindType为1、2、3时，取绑定商品 0时取推荐商品
	Products *Products `json:"products"`
	// 优质标签
	HighQualityTag *HighQualityTag `json:"high_quality_tag"`
	// 账号信息
	AccountInfo *AccountInfo `json:"account_info"`
}
type ShopVideoSearchItemData struct {
	// 返回数据
	Data *Data `json:"data"`
	// 总数
	Total int64 `json:"total"`
}
type ShopInfo struct {
	// 店铺ID
	ShopId string `json:"shop_id"`
}
type ItemInfo struct {
	// 时长
	Duration string `json:"duration"`
	// item状态，22-用户删除，23-审核删除，102-公开，140-支持运营、用户并行审核流程，141-待审/自见，142-需要编辑终审，143-好友可见，144-审核自见
	Status int32 `json:"status"`
	// 是否审核过
	Reviewed bool `json:"reviewed"`
	// 发布时间
	PublishTime int64 `json:"publish_time"`
	// 描述
	Caption string `json:"caption"`
	// 视频id
	ItemId string `json:"item_id"`
	// 视频名称
	Name string `json:"name"`
	// 封面
	Cover string `json:"cover"`
}
type ProductsItem struct {
	// 商品id
	Id string `json:"id"`
	// 产品名称
	Name string `json:"name"`
	// 封面url
	Img string `json:"img"`
}
type SuggestProductsItem struct {
	// 商品id
	Id string `json:"id"`
	// 产品名称
	Name string `json:"name"`
	// 封面url
	Img string `json:"img"`
}
type Products struct {
	// 关联商品列表
	BindProducts []BindProductsItem `json:"bind_products"`
	// 推荐商品列表
	SuggestProducts []SuggestProductsItem `json:"suggest_products"`
}
