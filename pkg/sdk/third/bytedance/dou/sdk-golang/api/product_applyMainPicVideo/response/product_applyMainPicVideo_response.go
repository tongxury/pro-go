package product_applyMainPicVideo_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ProductApplyMainPicVideoResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ProductApplyMainPicVideoData `json:"data"`
}
type ProductApplyMainPicVideoData struct {
	// 库中素材ID
	MaterialId int64 `json:"material_id"`
}
