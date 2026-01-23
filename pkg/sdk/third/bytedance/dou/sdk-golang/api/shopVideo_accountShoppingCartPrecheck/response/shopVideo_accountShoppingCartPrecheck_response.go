package shopVideo_accountShoppingCartPrecheck_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoAccountShoppingCartPrecheckResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *Data `json:"data"`
}
type Data struct {
	// 是否可绑定挂车视频
	CanBind bool `json:"can_bind"`
	// 剩余可绑定挂车视频数量
	RemainBindNum int64 `json:"remain_bind_num"`
	// 保证金是否缴纳
	DepositPaid bool `json:"deposit_paid"`
}
type ShopVideoAccountShoppingCartPrecheckData struct {
	// data
	Data *Data `json:"data"`
}
