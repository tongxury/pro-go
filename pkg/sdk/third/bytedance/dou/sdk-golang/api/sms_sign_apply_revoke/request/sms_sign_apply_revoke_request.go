package sms_sign_apply_revoke_request

import (
	"encoding/json"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/api/sms_sign_apply_revoke/response"
	"store/pkg/sdk/third/bytedance/dou/sdk-golang/core"
)

type SmsSignApplyRevokeRequest struct {
	doudian_sdk.BaseDoudianOpApiRequest
	Param *SmsSignApplyRevokeParam
}

func (c *SmsSignApplyRevokeRequest) GetUrlPath() string {
	return "/sms/sign/apply/revoke"
}

func New() *SmsSignApplyRevokeRequest {
	request := &SmsSignApplyRevokeRequest{
		Param: &SmsSignApplyRevokeParam{},
	}
	request.SetConfig(doudian_sdk.GlobalConfig)
	request.SetClient(doudian_sdk.DefaultDoudianOpApiClient)
	return request

}

func (c *SmsSignApplyRevokeRequest) Execute(accessToken *doudian_sdk.AccessToken) (*sms_sign_apply_revoke_response.SmsSignApplyRevokeResponse, error) {
	responseJson, err := c.GetClient().Request(c, accessToken)
	if err != nil {
		return nil, err
	}
	response := &sms_sign_apply_revoke_response.SmsSignApplyRevokeResponse{}
	_ = json.Unmarshal([]byte(responseJson), response)
	return response, nil

}

func (c *SmsSignApplyRevokeRequest) GetParamObject() interface{} {
	return c.Param
}

func (c *SmsSignApplyRevokeRequest) GetParams() *SmsSignApplyRevokeParam {
	return c.Param
}

type SmsSignApplyRevokeParam struct {
	// 短信发送渠道，主要做资源隔离
	SmsAccount string `json:"sms_account"`
	// 申请单id
	SmsSignApplyId string `json:"sms_sign_apply_id"`
}
