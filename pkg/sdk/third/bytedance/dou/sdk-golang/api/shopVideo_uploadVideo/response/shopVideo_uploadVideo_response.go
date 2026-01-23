package shopVideo_uploadVideo_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoUploadVideoResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoUploadVideoData `json:"data"`
}
type ShopVideoUploadVideoData struct {
	// 任务ID，根据任务ID查询上传任务状态
	TaskId string `json:"task_id"`
}
