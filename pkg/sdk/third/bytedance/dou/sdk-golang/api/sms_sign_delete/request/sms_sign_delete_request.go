package sms_sign_delete_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/sms_sign_delete/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SmsSignDeleteRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *SmsSignDeleteParam
}

func (c *SmsSignDeleteRequest) GetUrlPath() string {
	return "/sms/sign/delete"
}

func New() *SmsSignDeleteRequest {
	request := &SmsSignDeleteRequest{
		Param: &SmsSignDeleteParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *SmsSignDeleteRequest) Execute(accessToken *doudian_sdk.AccessToken) (*sms_sign_delete_response.SmsSignDeleteResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &sms_sign_delete_response.SmsSignDeleteResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *SmsSignDeleteRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *SmsSignDeleteRequest) GetParams() *SmsSignDeleteParam {
	return c.Param
}

type SmsSignDeleteParam struct {
	// 短信发送渠道，主要做资源隔离
	SmsAccount string `json:"sms_account"`
	// 签名
	Sign string `json:"sign"`
}
