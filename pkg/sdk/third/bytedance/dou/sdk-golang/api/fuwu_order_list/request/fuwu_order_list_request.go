package fuwu_order_list_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/fuwu_order_list/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type FuwuOrderListRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *FuwuOrderListParam
}

func (c *FuwuOrderListRequest) GetUrlPath() string {
	return "/fuwu/order/list"
}

func New() *FuwuOrderListRequest {
	request := &FuwuOrderListRequest{
		Param: &FuwuOrderListParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *FuwuOrderListRequest) Execute(accessToken *doudian_sdk.AccessToken) (*fuwu_order_list_response.FuwuOrderListResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &fuwu_order_list_response.FuwuOrderListResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *FuwuOrderListRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *FuwuOrderListRequest) GetParams() *FuwuOrderListParam {
	return c.Param
}

type FuwuOrderListParam struct {
	// 要搜索的订单创建时间起，unix秒级时间戳，必须大于2019年1月1日
	OrderCreateTimeStart int64 `json:"order_create_time_start"`
	// 要搜索的订单创建时间止，unix秒级时间戳，需小于次月1日
	OrderCreateTimeEnd int64 `json:"order_create_time_end"`
	// 页码，从1起
	PageIndex int64 `json:"page_index"`
	// 每页大小，最小1，最大100。
	PageSize int64 `json:"page_size"`
}
