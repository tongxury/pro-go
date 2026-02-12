package wavespeed

import (
	"context"
	"testing"
)

func TestClient_Gemini3ProImage(t *testing.T) {

	c := NewClient()

	c.Gemini3ProImage(context.Background(), Gemini3ProImageRequest{
		"参考给定的图，帮我生成2张图",
		[]string{"https://yoozyres.tos-cn-shanghai.volces.com/wechat_20260105113151_58_402.png"},
		"",
		"",
		"",
		false,
		false,
	})

	c.GetResult(context.Background(), "7d8d7372150044d283b4c66285bc8154")
}
