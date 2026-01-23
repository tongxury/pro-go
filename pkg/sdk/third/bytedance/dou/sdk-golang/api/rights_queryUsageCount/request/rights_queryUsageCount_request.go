package rights_queryUsageCount_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/rights_queryUsageCount/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type RightsQueryUsageCountRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *RightsQueryUsageCountParam
}

func (c *RightsQueryUsageCountRequest) GetUrlPath() string {
	return "/rights/queryUsageCount"
}

func New() *RightsQueryUsageCountRequest {
	request := &RightsQueryUsageCountRequest{
		Param: &RightsQueryUsageCountParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *RightsQueryUsageCountRequest) Execute(accessToken *doudian_sdk.AccessToken) (*rights_queryUsageCount_response.RightsQueryUsageCountResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &rights_queryUsageCount_response.RightsQueryUsageCountResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *RightsQueryUsageCountRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *RightsQueryUsageCountRequest) GetParams() *RightsQueryUsageCountParam {
	return c.Param
}

type RightsQueryUsageCountParam struct {
	// 外部业务ID(模板市场为模板code)，非必传，不可与service_id同传
	OuterBizId string `json:"outer_biz_id"`
	// 抖店服务市场服务ID，非必传，不可与outer_biz_id同传
	ServiceId int64 `json:"service_id"`
}
