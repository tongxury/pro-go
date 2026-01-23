package product_qualityDetail_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/product_qualityDetail/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type ProductQualityDetailRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *ProductQualityDetailParam
}

func (c *ProductQualityDetailRequest) GetUrlPath() string {
	return "/product/qualityDetail"
}

func New() *ProductQualityDetailRequest {
	request := &ProductQualityDetailRequest{
		Param: &ProductQualityDetailParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *ProductQualityDetailRequest) Execute(accessToken *doudian_sdk.AccessToken) (*product_qualityDetail_response.ProductQualityDetailResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &product_qualityDetail_response.ProductQualityDetailResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *ProductQualityDetailRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *ProductQualityDetailRequest) GetParams() *ProductQualityDetailParam {
	return c.Param
}

type ProductQualityDetailParam struct {
	// 商品ID
	ProductId int64 `json:"product_id"`
	// 质量分版本
	QualityScoreVersion string `json:"quality_score_version"`
}
