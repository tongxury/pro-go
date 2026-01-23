package tikhub

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_GetVideoByShareUrl(t *testing.T) {

	c := NewClient()

	url, err := c.GetVideoByShareUrl(context.Background(), "https://v.douyin.com/e3x2fjE/")

	fmt.Println(url, err)
}
