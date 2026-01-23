package shopVideo_getShopAccountInfo_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoGetShopAccountInfoResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoGetShopAccountInfoData `json:"data"`
}
type AccountInfosItem struct {
	// 用户ID
	UserId int64 `json:"user_id"`
	// 1-人店一体账号，2-授权号
	Type int64 `json:"type"`
	// 用户昵称
	UserName string `json:"user_name"`
	// 用户头像链接
	AvatarUrl string `json:"avatar_url"`
}
type ShopVideoGetShopAccountInfoData struct {
	// 账户信息
	AccountInfos []AccountInfosItem `json:"account_infos"`
}
