package fuwu_order_list_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type FuwuOrderListResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *FuwuOrderListData `json:"data"`
}
type OrdersItem struct {
	// 订单ID
	OrderId int64 `json:"order_id"`
	// 店铺ID
	ShopId int64 `json:"shop_id"`
	// AppID
	AppId int64 `json:"app_id"`
	// Sku标题
	SkuTitle string `json:"sku_title"`
	// Sku单价
	SkuPrice int64 `json:"sku_price"`
	// Sku周期数量
	SkuDuration int32 `json:"sku_duration"`
	// Sku周期单位，0-天，1-月，2-年，3-小时
	SkuDurationUnit int32 `json:"sku_duration_unit"`
	// Sku规格类型，0-普通规格，1-质检服务项目，2-微应用定制规格
	SkuSpecType int32 `json:"sku_spec_type"`
	// 规格名称
	SkuSpecValue string `json:"sku_spec_value"`
	// 用户支付类型，1-微信（2022年8月前订单会返回该字段），2-支付宝（2022年8月前订单会返回该字段），7-零元单，16-收银台付款
	PayType int32 `json:"pay_type"`
	// 支付金额，单位：分
	PayAmount int64 `json:"pay_amount"`
	// 履约开始时间，unix秒级时间戳
	ServiceStartTime int64 `json:"service_start_time"`
	// 履约结束时间，unix秒级时间戳
	ServiceEndTime int64 `json:"service_end_time"`
	// 退款状态，1-退款中，99-退款关闭，100-退款成功
	RefundStatus int32 `json:"refund_status"`
}
type FuwuOrderListData struct {
	// 订单列表
	Orders []OrdersItem `json:"orders"`
}
