package product_getPerformance_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ProductGetPerformanceResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ProductGetPerformanceData `json:"data"`
}
type ProductGetPerformanceData struct {
	// 商品ID
	ProductId string `json:"product_id"`
	// 是否动销
	HasSold bool `json:"has_sold"`
	// 是否曝光
	HasVisited bool `json:"has_visited"`
}
