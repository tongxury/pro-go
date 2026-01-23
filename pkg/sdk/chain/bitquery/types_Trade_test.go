package bitquery

import (
	"fmt"
	"testing"
)

func TestTrade_Parse(t *testing.T) {

	//c := Trade{
	//	Buy: TradeMetadata{
	//		Amount:      "124075.036887",
	//		AmountInUSD: "0",
	//		Currency: Currency{
	//			MintAddress: "2imW7gSCKYJBiBb6kKaN7Ru1pD1T6RVpzzFQhag2pump",
	//			Symbol:      "X",
	//		},
	//		Price:      0.00000009671567545798008,
	//		PriceInUSD: 0.00001830861413356116,
	//	},
	//	Sell: TradeMetadata{
	//		Amount:      "0.012000001",
	//		AmountInUSD: "2.271641973971451",
	//		Currency: Currency{
	//			Native: true,
	//		},
	//		Price:      10339585.545617871,
	//		PriceInUSD: 0,
	//	},
	//}

	c := Trade{
		Buy: TradeMetadata{
			Amount:      "124075.036887",
			AmountInUSD: "0",
			Currency: Currency{
				MintAddress: "xxxx",
				Symbol:      "xxxx",
			},
			Price:      0.00000009671567545798008,
			PriceInUSD: 0,
		},
		Sell: TradeMetadata{
			Amount:      "0.012000001",
			AmountInUSD: "0",
			Currency: Currency{
				Symbol: "USDC",
			},
			Price:      10339585.545617871,
			PriceInUSD: 0,
		},
	}

	v := c.Parse()

	fmt.Println(v)
}
