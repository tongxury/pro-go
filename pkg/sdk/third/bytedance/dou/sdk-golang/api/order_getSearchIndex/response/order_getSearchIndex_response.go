package order_getSearchIndex_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OrderGetSearchIndexResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *OrderGetSearchIndexData `json:"data"`
}
type OrderGetSearchIndexData struct {
	// 索引串
	EncryptIndexText string `json:"encrypt_index_text"`
}
