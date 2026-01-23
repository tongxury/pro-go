package shopVideo_queryUploadTask_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoQueryUploadTaskResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoQueryUploadTaskData `json:"data"`
}
type ShopVideoQueryUploadTaskData struct {
	// media_id，v开头的字符串
	Vid string `json:"vid"`
	// 短视频上传状态，0=初始化，1=上传中，2=上传成功，3=上传失败
	Status int64 `json:"status"`
}
