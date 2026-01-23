package product_isv_createProductFromSupplyPlatform_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ProductIsvCreateProductFromSupplyPlatformResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ProductIsvCreateProductFromSupplyPlatformData `json:"data"`
}
type ProductIsvCreateProductFromSupplyPlatformData struct {
}
