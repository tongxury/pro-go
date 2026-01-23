package freightTemplate_update_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type FreightTemplateUpdateResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *FreightTemplateUpdateData `json:"data"`
}
type FreightTemplateUpdateData struct {
	// 运费模板id
	TemplateId int64 `json:"template_id"`
}
