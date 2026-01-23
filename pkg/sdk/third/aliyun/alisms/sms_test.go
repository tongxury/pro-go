package alisms

import (
	"context"
	"fmt"
	"store/confs"
	"testing"
)

func TestSMSClient_SendSMS(t *testing.T) {
	var code = "2322"
	println(fmt.Sprintf(`{"code":"%s"}`, code))

	c, _ := NewClient(Config{
		AccessKey:    confs.AliyunAccessKey,
		AccessSecret: confs.AliyunAccessSecret,

		Sign:         "唯构科技深圳",
		TemplateCode: "SMS_316000047",
	})

	c.Send(context.Background(), []string{"15701348086"}, "123456")
}
