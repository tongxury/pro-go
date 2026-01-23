package coingecko

type Price struct {
	Address string
	Ts      int64
	Price   float64
}

type Prices []*Price

func (ps Prices) Avg() float64 {

	var total float64

	for _, p := range ps {
		total += p.Price
	}

	return total / float64(len(ps))

}

type Ohlc struct {
	Id    string
	Time  int64
	Open  float64
	High  float64
	Low   float64
	Close float64
}

type Ohlcs []*Ohlc

type Coin struct {
	Id     string
	Symbol string
	Name   string
}

type Coins []*Coin
