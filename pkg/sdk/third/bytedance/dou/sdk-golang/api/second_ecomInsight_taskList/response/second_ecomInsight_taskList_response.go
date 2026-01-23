package second_ecomInsight_taskList_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondEcomInsightTaskListResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *SecondEcomInsightTaskListData `json:"data"`
}
type RecordsItem struct {
	// 任务名
	Name string `json:"name"`
	// 任务状态
	Status int32 `json:"status"`
	// 创建任务信息
	CreateInfo string `json:"create_info"`
	// model任务id
	ModelTaskId string `json:"model_task_id"`
	// 更新时间
	UpdateTime string `json:"update_time"`
	// 消息
	Msg string `json:"msg"`
	// 任务类型
	TaskType string `json:"task_type"`
	// 任务结果
	TaskResult string `json:"task_result"`
	// pilot任务id
	PilotTaskId string `json:"pilot_task_id"`
	// 创建时间
	CreateTime string `json:"create_time"`
	// 唯一id
	Id int64 `json:"id"`
}
type Data struct {
	// 任务列表
	Records []RecordsItem `json:"records"`
	// 总条数
	Total int64 `json:"total"`
	// 页码
	Page int64 `json:"page"`
	// 条数
	PageSize int64 `json:"page_size"`
}
type SecondEcomInsightTaskListData struct {
	// 消息
	Msg string `json:"msg"`
	// 状态码
	Code int32 `json:"code"`
	// 任务信息
	Data *Data `json:"data"`
}
