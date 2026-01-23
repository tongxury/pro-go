package second_videoTop_mainPicVideo_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondVideoTopMainPicVideoResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *SecondVideoTopMainPicVideoData `json:"data"`
}
type RecordsItem struct {
	// 唯一id
	Id int64 `json:"id"`
	// 商品描述
	ProdDescBackup string `json:"prod_desc_backup"`
	// 截图
	Screenshot string `json:"screenshot"`
	// 视频id
	VideoId int64 `json:"video_id"`
	// 商品名称
	ProdName string `json:"prod_name"`
	// 商品id
	ProdId int64 `json:"prod_id"`
	// 店铺名
	ShopName string `json:"shop_name"`
	// 商品主图
	ProdMainImg string `json:"prod_main_img"`
	// 商品链接
	ProdUrl string `json:"prod_url"`
}
type Data struct {
	// 页码
	Page int64 `json:"page"`
	// 页码条数
	PageSize int64 `json:"page_size"`
	// 主图视频列表
	Records []RecordsItem `json:"records"`
	// 总条数
	Total int64 `json:"total"`
}
type SecondVideoTopMainPicVideoData struct {
	// 返回值
	Data *Data `json:"data"`
	// 返回消息
	Msg string `json:"msg"`
	// 状态码
	Code int32 `json:"code"`
}
