package wxpay

import (
	"context"
	"testing"
)

func TestClient_GetTradeByOutTradeNo(t *testing.T) {

	c := NewClient()

	c.GetTradeByOutTradeNo(context.Background(), "68a694539324c2b91cdf9377")
}
