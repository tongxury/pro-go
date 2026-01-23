package pricing

import (
	"context"
	"errors"
	"store/pkg/sdk/chain/gmgnai"
	"strconv"
	"time"
)

type Gmgnai struct {
}

func (t Gmgnai) GetPrice(ctx context.Context, params *GetPriceParams) (*GetPriceResult, error) {

	cs, err := gmgnai.NewClient().GetOHLCs(ctx, gmgnai.GetOHLCsParams{
		Token:      params.Token,
		Resolution: "1m",
		FromTs:     time.Now().Add(10 * time.Second).Unix(),
		ToTs:       time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	if len(cs) == 0 {
		return nil, errors.New("no prices found: " + params.Token)
	}

	result, _ := strconv.ParseFloat(cs[len(cs)-1].Close, 64)

	return &GetPriceResult{
		Value: result,
	}, nil
}

func NewGmgnai() IPricing {
	return &Gmgnai{}

}
