package microapp_invoke_cloudApi_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/microapp_invoke_cloudApi/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type MicroappInvokeCloudApiRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *MicroappInvokeCloudApiParam
}

func (c *MicroappInvokeCloudApiRequest) GetUrlPath() string {
	return "/microapp/invoke/cloudApi"
}

func New() *MicroappInvokeCloudApiRequest {
	request := &MicroappInvokeCloudApiRequest{
		Param: &MicroappInvokeCloudApiParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *MicroappInvokeCloudApiRequest) Execute(accessToken *doudian_sdk.AccessToken) (*microapp_invoke_cloudApi_response.MicroappInvokeCloudApiResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &microapp_invoke_cloudApi_response.MicroappInvokeCloudApiResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *MicroappInvokeCloudApiRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *MicroappInvokeCloudApiRequest) GetParams() *MicroappInvokeCloudApiParam {
	return c.Param
}

type MicroappInvokeCloudApiParam struct {
	// 云内应用提供的接口名
	InterfaceName string `json:"interface_name"`
	// 入参，json字符串
	BizParamJson string `json:"biz_param_json"`
	// 是否调用至"测试环境"云内应用的接口
	UseTest int64 `json:"use_test"`
}
