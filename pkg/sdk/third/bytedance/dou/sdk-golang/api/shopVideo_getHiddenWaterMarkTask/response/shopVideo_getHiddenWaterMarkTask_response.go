package shopVideo_getHiddenWaterMarkTask_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGetHiddenWaterMarkTaskResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoGetHiddenWaterMarkTaskData `json:"data"`
}
type ShopVideoGetHiddenWaterMarkTaskData struct {
	// 加完暗水印后的资源url，在status=3时有值，注意：该url有效期10分钟，禁止对外暴露，请在自己系统内转储后给商家播放/下载
	WatermarkUrl string `json:"watermark_url"`
	// 任务id
	TaskId string `json:"task_id"`
	// 任务状态：1：上传视频中，2-加暗水印中，3：任务成功，4：任务失败
	DisplayStatus int32 `json:"display_status"`
	// 失败原因，在task_status=4时有值
	FailedReason string `json:"failed_reason"`
}
