package shopVideo_listRelatableProduct_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/shopVideo_listRelatableProduct/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ShopVideoListRelatableProductRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ShopVideoListRelatableProductParam
}

func (c *ShopVideoListRelatableProductRequest) GetUrlPath() string {
	return "/shopVideo/listRelatableProduct"
}

func New() *ShopVideoListRelatableProductRequest {
	request := &ShopVideoListRelatableProductRequest{
		Param: &ShopVideoListRelatableProductParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ShopVideoListRelatableProductRequest) Execute(accessToken *doudian_sdk.AccessToken) (*shopVideo_listRelatableProduct_response.ShopVideoListRelatableProductResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &shopVideo_listRelatableProduct_response.ShopVideoListRelatableProductResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ShopVideoListRelatableProductRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ShopVideoListRelatableProductRequest) GetParams() *ShopVideoListRelatableProductParam {
	return c.Param
}

type ShopVideoListRelatableProductParam struct {
	// 商品id列表
	ProductIdList []int64 `json:"product_id_list"`
	// 销量左区间
	SellNumStart *int64 `json:"sell_num_start"`
	// 销量右区间
	SellNumEnd *int64 `json:"sell_num_end"`
	// 页码
	Page int32 `json:"page"`
	// 每页数量
	PageSize int32 `json:"page_size"`
	// 类目id列表
	CategoryIdList []int64 `json:"category_id_list"`
	// 商品名称
	ProductName *string `json:"product_name"`
}
