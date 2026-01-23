package second_videoTop_shortVideo_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SecondVideoTopShortVideoResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *SecondVideoTopShortVideoData `json:"data"`
}
type RecordsItem struct {
	// 唯一id
	Id int64 `json:"id"`
	// 商品名
	ProductName string `json:"product_name"`
	// 商品二级类目
	ProductSecondCatName string `json:"product_second_cat_name"`
	// 视频url
	VideoUrl string `json:"video_url"`
	// 视频id
	VideoId string `json:"video_id"`
	// 店铺名
	ShopName string `json:"shop_name"`
	// 商品一级类目
	ProductFirstCatName string `json:"product_first_cat_name"`
	// 商品叶子类目
	ProductLeafCatName string `json:"product_leaf_cat_name"`
	// 视频创建日期
	VideoCreateDate string `json:"video_create_date"`
	// 视频标题
	VideoTitle string `json:"video_title"`
}
type Data struct {
	// 返回总条数
	Total int64 `json:"total"`
	// 页码
	Page int64 `json:"page"`
	// 条数
	PageSize int64 `json:"page_size"`
	// 短视频列表
	Records []RecordsItem `json:"records"`
}
type SecondVideoTopShortVideoData struct {
	// 数据
	Data *Data `json:"data"`
	// 消息
	Msg string `json:"msg"`
	// 状态码
	Code int32 `json:"code"`
}
