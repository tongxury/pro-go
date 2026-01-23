package rights_deductUsageCount_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type RightsDeductUsageCountResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *RightsDeductUsageCountData `json:"data"`
}
type RightsDeductUsageCountData struct {
}
