package dbtypes

import "time"

type State struct {
	Symbol                   string
	Timestamp                time.Time
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	OpenTime                 time.Time
	CloseTime                time.Time
	QuoteAssetVolume         float64
	TradeNum                 int64
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
	Rsi                      float64
	Macd12269                float64
	Macdh12269               float64
	Macds12269               float64
}
