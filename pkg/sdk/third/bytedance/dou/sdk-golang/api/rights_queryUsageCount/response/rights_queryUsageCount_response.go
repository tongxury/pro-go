package rights_queryUsageCount_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type RightsQueryUsageCountResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *RightsQueryUsageCountData `json:"data"`
}
type RightsQueryUsageCountData struct {
	// 可扣减次数
	CanDeductCount int64 `json:"canDeductCount"`
	// 冻结中的次数(处于退款中/换版本中等状态的次数)
	FrozenCount int64 `json:"frozenCount"`
	// 总的可使用次数
	TotalCount int64 `json:"totalCount"`
}
