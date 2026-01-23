package rights_deductUsageCount_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/rights_deductUsageCount/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type RightsDeductUsageCountRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *RightsDeductUsageCountParam
}

func (c *RightsDeductUsageCountRequest) GetUrlPath() string {
	return "/rights/deductUsageCount"
}

func New() *RightsDeductUsageCountRequest {
	request := &RightsDeductUsageCountRequest{
		Param: &RightsDeductUsageCountParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *RightsDeductUsageCountRequest) Execute(accessToken *doudian_sdk.AccessToken) (*rights_deductUsageCount_response.RightsDeductUsageCountResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &rights_deductUsageCount_response.RightsDeductUsageCountResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *RightsDeductUsageCountRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *RightsDeductUsageCountRequest) GetParams() *RightsDeductUsageCountParam {
	return c.Param
}

type RightsDeductUsageCountParam struct {
	// 外部业务ID(模板市场为模板code)，非必传，不可与service_id同传
	OuterBizId string `json:"outer_biz_id"`
	// 抖店服务市场服务ID，非必传，不可与outer_biz_id同传
	ServiceId int64 `json:"service_id"`
	// 用户使用软件时, 调用抖店开放平台接口的日志id
	LogId string `json:"log_id"`
	// 请求id, 用作幂等键
	ReqId string `json:"req_id"`
	// 扣减次数
	DeductCnt int64 `json:"deduct_cnt"`
}
