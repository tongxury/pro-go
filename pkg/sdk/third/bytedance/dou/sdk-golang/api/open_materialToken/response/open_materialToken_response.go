package open_materialToken_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenMaterialTokenResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *OpenMaterialTokenData `json:"data"`
}
type OpenMaterialTokenData struct {
	// 获取下载地址query
	AuthQuery string `json:"auth_query"`
}
