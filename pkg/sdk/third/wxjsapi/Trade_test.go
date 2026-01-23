package wxjsapi

import (
	"context"
	"testing"
)

func TestClient_Trade(t *testing.T) {

	c, _ := NewClient()

	c.Trade(context.Background(), TradeParams{})

	//c.GetUserPhoneNumber(context.Background(), "05a282ef27b0ad7ce70e14d10112398f33db4c181c1fe59f23353ce3b9c62edd")

}
