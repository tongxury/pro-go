package open_getProductVideoCheckResult_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenGetProductVideoCheckResultResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *OpenGetProductVideoCheckResultData `json:"data"`
}
type OpenGetProductVideoCheckResultData struct {
	// 任务状态 0执行成功 1执行中 2执行失败
	TaskStatus int64 `json:"task_status"`
	// 主图视频是否合格 0不合格 1合格
	IsPass int64 `json:"is_pass"`
}
