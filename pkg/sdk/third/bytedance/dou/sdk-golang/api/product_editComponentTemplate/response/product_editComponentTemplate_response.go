package product_editComponentTemplate_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ProductEditComponentTemplateResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ProductEditComponentTemplateData `json:"data"`
}
type ProductEditComponentTemplateData struct {
}
