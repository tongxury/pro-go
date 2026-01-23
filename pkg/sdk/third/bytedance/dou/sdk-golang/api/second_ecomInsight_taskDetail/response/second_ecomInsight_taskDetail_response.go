package second_ecomInsight_taskDetail_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightTaskDetailResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *SecondEcomInsightTaskDetailData `json:"data"`
}
type Data struct {
	// 任务id
	Id int64 `json:"id"`
	// 任务类型
	TaskType string `json:"task_type"`
	// 模型taskId
	ModelTaskId string `json:"model_task_id"`
	// pilotTaskId
	PilotTaskId string `json:"pilot_task_id"`
	// 创建时间
	CreateTime string `json:"create_time"`
	// 更新时间
	UpdateTime string `json:"update_time"`
	// 消息
	Msg string `json:"msg"`
	// 任务名
	Name string `json:"name"`
	// 任务状态：0-失败、1-创建、2-运行中、3-成功
	Status int32 `json:"status"`
	// JSON string
	CreateInfo string `json:"create_info"`
	// JSON string
	TaskResult string `json:"task_result"`
}
type SecondEcomInsightTaskDetailData struct {
	// 消息
	Msg string `json:"msg"`
	// 状态码
	Code int32 `json:"code"`
	// 任务结构体
	Data *Data `json:"data"`
}
