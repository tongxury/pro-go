package shopVideo_publishVideoNew_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoPublishVideoNewResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *Data `json:"data"`
}
type Data struct {
	// 视频id
	ItemId string `json:"item_id"`
}
type ShopVideoPublishVideoNewData struct {
	// 响应参数
	Data *Data `json:"data"`
}
