package sms_template_revoke_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SmsTemplateRevokeResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *SmsTemplateRevokeData `json:"data"`
}
type SmsTemplateRevokeData struct {
	// 是否成功 0表示成功
	Code int64 `json:"code"`
	// 说明
	Message string `json:"message"`
}
