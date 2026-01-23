package shopVideo_submitHiddenWaterMarkTask_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoSubmitHiddenWaterMarkTaskResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoSubmitHiddenWaterMarkTaskData `json:"data"`
}
type ShopVideoSubmitHiddenWaterMarkTaskData struct {
	// 任务id
	TaskId string `json:"task_id"`
}
