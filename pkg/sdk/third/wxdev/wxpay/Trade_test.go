package wxpay

import (
	"context"
	"log"
	"testing"
)

func TestClient_Trade(t *testing.T) {

	c := NewClient()

	trade, err := c.Trade(context.Background(), TradeParams{
		Description: "11111111111",
		OutTradeNo:  "111111111111",
		NotifyUrl:   "https://api.mch.weixin.qq.com/pay/unifiedorder",
		TotalAmount: 100,
		OpenId:      "oyetIvj0rDA6unzW4BnlTOPYjpnA",
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(trade)

	//c.GetUserPhoneNumber(context.Background(), "05a282ef27b0ad7ce70e14d10112398f33db4c181c1fe59f23353ce3b9c62edd")

}
