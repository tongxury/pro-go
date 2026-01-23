package sms

import (
	"context"
	"store/pkg/sdk/conv"

	"store/confs"

	"github.com/volcengine/volc-sdk-golang/service/sms"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SendSmsCode(context context.Context, phoneNum, code string) error {

	sms.DefaultInstance.Client.SetAccessKey(confs.BytedanceSmsAccessKey)
	sms.DefaultInstance.Client.SetSecretKey(confs.BytedanceSmsSecretKey)
	req := &sms.SmsRequest{
		SmsAccount:    "87cc2ff3",
		Sign:          "广州予之文化",
		TemplateID:    "S1T_1y2p19gp2rpj6",
		TemplateParam: conv.M2J(map[string]string{"code": code}),
		PhoneNumbers:  phoneNum,
		Tag:           "",
	}
	result, statusCode, err := sms.DefaultInstance.Send(req)
	if err != nil {
		return err
	}

	_, _ = result, statusCode

	return nil
}
