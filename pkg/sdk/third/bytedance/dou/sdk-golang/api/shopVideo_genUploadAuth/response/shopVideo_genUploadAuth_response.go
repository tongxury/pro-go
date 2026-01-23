package shopVideo_genUploadAuth_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGenUploadAuthResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoGenUploadAuthData `json:"data"`
}
type ShopVideoGenUploadAuthData struct {
	// 生成授权时间
	CurrentTime string `json:"current_time"`
	// 临时的accessKeyId
	AccessKeyId string `json:"access_key_id"`
	// 临时的secretAccessKey
	SecretAccessKey string `json:"secret_access_key"`
	// 临时的会话token
	SessionToken string `json:"session_token"`
	// 空间名，上传时必须使用该空间，否则会上传失败
	SpaceName string `json:"space_name"`
	// 授权过期时间
	ExpiredTime string `json:"expired_time"`
}
