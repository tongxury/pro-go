package open_createProductVideoCheckTask_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenCreateProductVideoCheckTaskResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *OpenCreateProductVideoCheckTaskData `json:"data"`
}
type OpenCreateProductVideoCheckTaskData struct {
	// 任务id
	TaskId string `json:"task_id"`
}
