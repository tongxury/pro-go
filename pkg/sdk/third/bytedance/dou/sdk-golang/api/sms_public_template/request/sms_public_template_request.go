package sms_public_template_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/sms_public_template/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SmsPublicTemplateRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *SmsPublicTemplateParam
}

func (c *SmsPublicTemplateRequest) GetUrlPath() string {
	return "/sms/public/template"
}

func New() *SmsPublicTemplateRequest {
	request := &SmsPublicTemplateRequest{
		Param: &SmsPublicTemplateParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *SmsPublicTemplateRequest) Execute(accessToken *doudian_sdk.AccessToken) (*sms_public_template_response.SmsPublicTemplateResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &sms_public_template_response.SmsPublicTemplateResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *SmsPublicTemplateRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *SmsPublicTemplateRequest) GetParams() *SmsPublicTemplateParam {
	return c.Param
}

type SmsPublicTemplateParam struct {
	// 每页数据大小
	Size *int64 `json:"size"`
	// 第几页，从0开始
	Page *int64 `json:"page"`
	// 模版id
	TemplateId *string `json:"template_id"`
}
