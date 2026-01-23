package shopVideo_getSugWords_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGetSugWordsResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data []DataItem `json:"data"`
}
type DataItem struct {
	// 话题名称
	TopicName string `json:"topic_name"`
	// 话题ID
	TopicId string `json:"topic_id"`
	// 话题热度
	PlayCount int64 `json:"play_count"`
}
type ShopVideoGetSugWordsData struct {
	// 响应参数
	Data []DataItem `json:"data"`
}
