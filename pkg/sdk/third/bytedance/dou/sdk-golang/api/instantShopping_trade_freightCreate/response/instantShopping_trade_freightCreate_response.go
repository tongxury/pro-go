package instantShopping_trade_freightCreate_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type InstantShoppingTradeFreightCreateResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *InstantShoppingTradeFreightCreateData `json:"data"`
}
type InstantShoppingTradeFreightCreateData struct {
	// 模板id
	TemplateId int64 `json:"template_id"`
}
