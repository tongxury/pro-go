package material_easyShuttle_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type MaterialEasyShuttleResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *MaterialEasyShuttleData `json:"data"`
}
type MaterialEasyShuttleData struct {
}
