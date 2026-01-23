package product_isv_saveGoodsSupplyStatus_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ProductIsvSaveGoodsSupplyStatusResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ProductIsvSaveGoodsSupplyStatusData `json:"data"`
}
type ProductIsvSaveGoodsSupplyStatusData struct {
}
