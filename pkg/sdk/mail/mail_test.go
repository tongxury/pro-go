package mail

import (
	"fmt"
	"testing"
)

func TestPoster_Send(t *testing.T) {

	poster := NewClient(Config{})

	err := poster.Send(Params{
		//Subject: "Subject",
		//To:      "tongxurt@gmail.com",
		//Content: "content content",
		Subject: "【Veogo】登录验证",
		To:      "tongxurt@gmail.com",
		Content: fmt.Sprintf("<div>【Veogo】您的验证码为: <strong>%d</strong>, 5分钟之内有效。</div>", 123421),
	})

	if err != nil {
		panic(err)
	}
}
