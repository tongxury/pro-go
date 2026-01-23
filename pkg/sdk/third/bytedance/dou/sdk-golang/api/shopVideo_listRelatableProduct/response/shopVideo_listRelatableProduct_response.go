package shopVideo_listRelatableProduct_response

import (
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoListRelatableProductResponse struct {
	doudian_sdk.BaseDoudianOpApiResponse
	Data *ShopVideoListRelatableProductData `json:"data"`
}
type ShopVideoListRelatableProductData struct {
	// 商品列表
	ProductList []ProductListItem `json:"product_list"`
	// 总数
	Total int64 `json:"total"`
	// 当前页
	Page int32 `json:"page"`
	// 每页数量
	PageSize int32 `json:"page_size"`
}
type ProductListItem struct {
	// 库存数量
	StockNum int64 `json:"stock_num"`
	// 原价，单位分
	MarketPrice int64 `json:"market_price"`
	// 商品列表图片
	Img string `json:"img"`
	// 商品状态，0-上架中，1-下架中，2-删除中
	Status int64 `json:"status"`
	// 草稿审核状态，1-待提交，2-审核中，3-审核通过，4-审核被驳回
	DraftStatus int64 `json:"draft_status"`
	// 小店商品审核状态，1-待提交，2-审核中，3-审核通过，4-审核被驳回，5-被封禁
	CheckStatus int64 `json:"check_status"`
	// 商品id
	ProductId int64 `json:"product_id"`
	// 商品名称
	Name string `json:"name"`
}
