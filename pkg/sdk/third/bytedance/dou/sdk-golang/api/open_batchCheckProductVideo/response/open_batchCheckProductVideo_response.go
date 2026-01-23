package open_batchCheckProductVideo_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenBatchCheckProductVideoResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *OpenBatchCheckProductVideoData `json:"data"`
}
type OpenBatchCheckProductVideoData struct {
	// 出参
	ProductVideoCheckResultList []ProductVideoCheckResultListItem `json:"product_video_check_result_list"`
}
type ProductVideoCheckResultListItem struct {
	// 是否合格 1合格 0不合格
	IsPass int64 `json:"is_pass"`
	// 错误码
	Code int64 `json:"code"`
	// 错误信息
	Msg string `json:"msg"`
}
