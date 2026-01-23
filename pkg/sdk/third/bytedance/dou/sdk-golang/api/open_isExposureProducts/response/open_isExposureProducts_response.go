package open_isExposureProducts_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenIsExposureProductsResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *OpenIsExposureProductsData `json:"data"`
}
type OpenIsExposureProductsData struct {
	// 结果集，key商品id；value结果 0-没有曝光 1-有曝光 2-高曝光
	IsExposureMap map[int64]int16 `json:"is_exposure_map"`
}
