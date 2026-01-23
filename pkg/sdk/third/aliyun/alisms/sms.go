package alisms

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"strings"
)

type Config struct {
	AccessKey    string
	AccessSecret string
	Sign         string
	TemplateCode string
}

func NewClient(conf Config) (*Client, error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", conf.AccessKey, conf.AccessSecret)
	if err != nil {
		return nil, err
	}

	return &Client{client: client, conf: conf}, nil
}

type Client struct {
	client *dysmsapi.Client
	//sign   string
	conf Config
}

func (s *Client) Send(ctx context.Context, phoneNos []string, code string) (*dysmsapi.SendSmsResponse, error) {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = strings.Join(phoneNos, ",")
	request.SignName = s.conf.Sign
	request.TemplateCode = s.conf.TemplateCode
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)

	response, err := s.client.SendSms(request)

	if err != nil {
		return nil, err
	}
	if response.Code != "OK" {
		return nil, errors.New("aliyun.oss: " + response.Message)
	}

	return response, nil
}
