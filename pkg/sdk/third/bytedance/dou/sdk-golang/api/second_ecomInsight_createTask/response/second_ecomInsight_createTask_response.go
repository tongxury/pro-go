package second_ecomInsight_createTask_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightCreateTaskResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *SecondEcomInsightCreateTaskData `json:"data"`
}
type SecondEcomInsightCreateTaskData struct {
	// 返回结构体
	Data *Data `json:"data"`
	// 消息
	Msg string `json:"msg"`
	// 状态码
	Code int32 `json:"code"`
}
type Data struct {
	// 更新时间
	UpdateTime string `json:"update_time"`
	// 消息
	Msg string `json:"msg"`
	// 任务名
	Name string `json:"name"`
	// 创建结构体
	CreateInfo string `json:"create_info"`
	// 模型taskId
	ModelTaskId string `json:"model_task_id"`
	// pilotTaskId
	PilotTaskId string `json:"pilot_task_id"`
	// 任务id
	Id int64 `json:"id"`
	// 状态：0-失败、1-创建、2-运行中、3-成功、4-取消
	Status int32 `json:"status"`
	// 任务类型
	TaskType string `json:"task_type"`
	// 任务返回值
	TaskResult string `json:"task_result"`
	// 创建时间
	CreateTime string `json:"create_time"`
}
