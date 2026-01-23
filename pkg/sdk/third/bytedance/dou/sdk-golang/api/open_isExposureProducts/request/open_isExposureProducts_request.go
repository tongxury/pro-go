package open_isExposureProducts_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/open_isExposureProducts/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type OpenIsExposureProductsRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *OpenIsExposureProductsParam
}

func (c *OpenIsExposureProductsRequest) GetUrlPath() string {
	return "/open/isExposureProducts"
}

func New() *OpenIsExposureProductsRequest {
	request := &OpenIsExposureProductsRequest{
		Param: &OpenIsExposureProductsParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *OpenIsExposureProductsRequest) Execute(accessToken *doudian_sdk.AccessToken) (*open_isExposureProducts_response.OpenIsExposureProductsResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &open_isExposureProducts_response.OpenIsExposureProductsResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *OpenIsExposureProductsRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *OpenIsExposureProductsRequest) GetParams() *OpenIsExposureProductsParam {
	return c.Param
}

type OpenIsExposureProductsParam struct {
	// 商品id
	ProductIdList []int64 `json:"product_id_list"`
}
