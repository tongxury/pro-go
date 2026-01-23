package jupiter

import (
	"context"
	"github.com/go-resty/resty/v2"
	"store/pkg/sdk/helper/restyd"
	"strings"
	"time"
)

//type TokenMetadata struct {
//	Address         string      `json:"address"`
//	Name            string      `json:"name"`
//	Symbol          string      `json:"symbol"`
//	Decimals        int         `json:"decimals"`
//	LogoURI         string      `json:"logoURI"`
//	Tags            []string    `json:"tags"`
//	DailyVolume     float64     `json:"daily_volume"`
//	FreezeAuthority interface{} `json:"freeze_authority"`
//	MintAuthority   string      `json:"mint_authority"`
//}

type TokenMetadata struct {
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
	DailyVolume float64   `json:"daily_volume"`
	Decimals    int       `json:"decimals"`
	Extensions  struct {
		CoingeckoId string `json:"coingeckoId"`
	} `json:"extensions"`
	FreezeAuthority   interface{} `json:"freeze_authority"`
	LogoURI           string      `json:"logoURI"`
	MintAuthority     interface{} `json:"mint_authority"`
	MintedAt          time.Time   `json:"minted_at"`
	Name              string      `json:"name"`
	PermanentDelegate interface{} `json:"permanent_delegate"`
	Symbol            string      `json:"symbol"`
	Tags              []string    `json:"tags"`
}

type TokenMetadatas []*TokenMetadata

func (t *JupiterClient) GetTokenMetadata(ctx context.Context, tokenID string) (*TokenMetadata, error) {

	url := t.tokenEndpoint + "/token/" + tokenID

	var quote TokenMetadata

	rsp, err := resty.New().R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	err = restyd.ParseResult(rsp, &quote)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

type ListTokensParams struct {
	Tags []string
}

func (t *JupiterClient) ListTokens(ctx context.Context, params ListTokensParams) (TokenMetadatas, error) {

	url := t.tokenEndpoint + "/tokens"

	var quote TokenMetadatas

	rsp, err := resty.New().R().
		SetQueryParam("tags", strings.Join(params.Tags, ",")).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	err = restyd.ParseResult(rsp, &quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

// large response
func (t *JupiterClient) ListTradableTokens(ctx context.Context) error {

	url := t.tokenEndpoint + "/tokens_with_markets?toke=10"

	// 结构没确认
	var quote TokenMetadata

	rsp, err := resty.New().R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return err
	}

	err = restyd.ParseResult(rsp, &quote)
	if err != nil {
		return err
	}

	return nil
}
