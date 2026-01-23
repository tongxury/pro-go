package product_createComponentTemplateV2_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ProductCreateComponentTemplateV2Response struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ProductCreateComponentTemplateV2Data `json:"data"`
}
type ProductCreateComponentTemplateV2Data struct {
	// 模板ID
	TemplateId int64 `json:"template_id"`
}
