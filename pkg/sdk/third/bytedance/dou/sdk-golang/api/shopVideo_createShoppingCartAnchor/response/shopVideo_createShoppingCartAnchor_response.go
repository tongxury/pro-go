package shopVideo_createShoppingCartAnchor_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoCreateShoppingCartAnchorResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data string `json:"data"`
}
type ShopVideoCreateShoppingCartAnchorData struct {
	// 生成的锚点ID
	Data string `json:"data"`
}
