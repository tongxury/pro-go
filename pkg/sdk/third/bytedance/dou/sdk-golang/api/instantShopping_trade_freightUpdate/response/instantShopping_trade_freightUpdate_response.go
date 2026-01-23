package instantShopping_trade_freightUpdate_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type InstantShoppingTradeFreightUpdateResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *InstantShoppingTradeFreightUpdateData `json:"data"`
}
type InstantShoppingTradeFreightUpdateData struct {
	// 运费模板id
	TemplateId int64 `json:"template_id"`
}
