package sms_public_template_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SmsPublicTemplateResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *SmsPublicTemplateData `json:"data"`
}
type SmsPublicTemplateData struct {
	// 数据总量
	Total int64 `json:"total"`
	// 列表页数据
	PublicTemplateList []PublicTemplateListItem `json:"public_template_list"`
}
type PublicTemplateListItem struct {
	// 模版id
	TemplateId string `json:"template_id"`
	// 模版名称
	TemplateName string `json:"template_name"`
	// 模版内容
	TemplateContent string `json:"template_content"`
	// 模版类型
	ChannelType string `json:"channel_type"`
}
