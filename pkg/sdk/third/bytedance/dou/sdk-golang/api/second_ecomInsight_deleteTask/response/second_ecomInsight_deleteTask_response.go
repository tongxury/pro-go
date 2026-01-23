package second_ecomInsight_deleteTask_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightDeleteTaskResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *SecondEcomInsightDeleteTaskData `json:"data"`
}
type Data struct {
	// pilot任务id
	PilotTaskId string `json:"pilot_task_id"`
	// 更新时间
	UpdateTime string `json:"update_time"`
	// 任务id
	Id int64 `json:"id"`
	// 任务名
	Name string `json:"name"`
	// 任务状态
	Status int32 `json:"status"`
	// 任务类型
	TaskType string `json:"task_type"`
	// JSON string
	CreateInfo string `json:"create_info"`
	// JSON string
	TaskResult string `json:"task_result"`
	// 创建时间
	CreateTime string `json:"create_time"`
	// 消息
	Msg string `json:"msg"`
	// 模型任务id
	ModelTaskId string `json:"model_task_id"`
}
type SecondEcomInsightDeleteTaskData struct {
	// 消息
	Msg string `json:"msg"`
	// 状态码
	Code int32 `json:"code"`
	// 结构体
	Data *Data `json:"data"`
}
