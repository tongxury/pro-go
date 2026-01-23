package shopVideo_publishVideo_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoPublishVideoResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoPublishVideoData `json:"data"`
}
type ShopVideoPublishVideoData struct {
	// 抖音主端短视频ID
	ItemId int64 `json:"item_id"`
}
