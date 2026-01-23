package microapp_invoke_cloudApi_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type MicroappInvokeCloudApiResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *MicroappInvokeCloudApiData `json:"data"`
}
type MicroappInvokeCloudApiData struct {
	// 后端接口的返回值，json string
	BizResultJson string `json:"biz_result_json"`
}
