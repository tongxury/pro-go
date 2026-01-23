package ape

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
	"time"
)

// https://api.ape.pro/api/v1/pools?createdAt=2024-11-14T16%3A45%3A10.540Z&sortBy=listedTime&sortDir=desc&limit=50&offset=0&notPumpfunToken=false
// https://api.ape.pro/api/v1/pools?limit=10&offset=0&searchText=rG19j9jSirQx4oSzvgnhWoSbBuzQEPcVu31w2jFpump
func (t *Client) GetPoolsByToken(ctx context.Context, token string, cache bool) ([]Pool, error) {

	if cache {
		if cached, found := t.tokenPoolsCache[token]; found {
			return cached, nil
		}
	}

	pools, err := t.ListPools(ctx, ListPoolsParams{
		SearchText: token,
	})
	if err != nil {
		return nil, err
	}

	if len(pools) == 0 {
		return nil, fmt.Errorf("no pools found for token: %s", token)
	}

	t.tokenPoolsCache[token] = pools

	return pools, nil
}

//curl 'https://api.ape.pro/api/v1/pools?createdAt=2024-11-21T02%3A56%3A40.058Z&sortBy=listedTime&sortDir=desc&limit=50&offset=0&notPumpfunToken=false' \
//-H 'accept: */*' \
//-H 'accept-language: zh-CN,zh;q=0.9,ar;q=0.8' \
//-H 'cache-control: no-cache' \
//-H 'origin: https://ape.pro' \
//-H 'pragma: no-cache' \
//-H 'priority: u=1, i' \
//-H 'referer: https://ape.pro/' \
//-H 'sec-ch-ua: "Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"' \
//-H 'sec-ch-ua-mobile: ?0' \
//-H 'sec-ch-ua-platform: "macOS"' \
//-H 'sec-fetch-dest: empty' \
//-H 'sec-fetch-mode: cors' \
//-H 'sec-fetch-site: same-site' \
//-H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'

//curl 'https://api.ape.pro/api/v1/pools?createdAt=2024-11-21T03%3A03%3A33.564Z&sortBy=listedTime&sortDir=desc&limit=50&offset=0&notPumpfunToken=false' \
//-H 'accept: */*' \
//-H 'accept-language: zh-CN,zh;q=0.9,ar;q=0.8' \
//-H 'cache-control: no-cache' \
//-H 'origin: https://ape.pro' \
//-H 'pragma: no-cache' \
//-H 'priority: u=1, i' \
//-H 'referer: https://ape.pro/' \
//-H 'sec-ch-ua: "Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"' \
//-H 'sec-ch-ua-mobile: ?0' \
//-H 'sec-ch-ua-platform: "macOS"' \
//-H 'sec-fetch-dest: empty' \
//-H 'sec-fetch-mode: cors' \
//-H 'sec-fetch-site: same-site' \
//-H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'

func (t *Client) ListPools(ctx context.Context, params ListPoolsParams) ([]Pool, error) {

	url := "https://api.ape.pro/api/v1/pools"

	reqParams := map[string]string{
		//"createdAt":       "2024-11-21T03:03:33.564Z", // todo
		"createdAt":       time.Now().Format("2006-01-02T15:04:05.999Z"), // todo
		"sortBy":          "listedTime",
		"sortDir":         "desc",
		"notPumpfunToken": "false",
		"limit":           "50",
		"offset":          "0",
	}

	//
	//if params.CreatedAt != "" {
	//	reqParams["createdAt"] = params.CreatedAt
	//}

	if params.SearchText != "" {
		reqParams["searchText"] = params.SearchText
		delete(reqParams, "createdAt")
	}

	if len(params.AssetIds) > 0 {
		reqParams["assetIds"] = strings.Join(params.AssetIds, ",")
	}

	//if params.Limit > 0 {
	//	reqParams["limit"] = fmt.Sprintf("%d", params.Limit)
	//}

	if params.Offset > 0 {
		reqParams["offset"] = fmt.Sprintf("%d", params.Offset)
	}

	result, err := resty.New().R().
		SetContext(ctx).
		SetQueryParams(reqParams).
		SetHeaders(map[string]string{
			//"accept":             "*/*",
			//"accept-language":    "zh-CN,zh;q=0.9,ar;q=0.8",
			//"cache-control":      "no-cache",
			//"origin": "https://ape.pro",
			//"pragma":             "no-cache",
			//"priority":           "u=1, i",
			//"referer":            "https://ape.pro/",
			//"sec-ch-ua":          `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
			//"sec-ch-ua-mobile":   "?0",
			//"sec-ch-ua-platform": `"macOS"`,
			//"sec-fetch-dest":     "empty",
			//"sec-fetch-mode":     "cors",
			//"sec-fetch-site":     "same-site",
			//"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		}).
		Get(url)

	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, fmt.Errorf(result.String())
	}

	var pools PoolList
	if err := json.Unmarshal(result.Body(), &pools); err != nil {
		return nil, err
	}

	return pools.Pools, nil
}

type ListPoolsParams struct {
	SearchText string
	AssetIds   []string
	CreatedAt  string
	Limit      int
	Offset     int
}

type Asset struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	Icon         string  `json:"icon"`
	Decimals     int     `json:"decimals"`
	Dev          string  `json:"dev"`
	UsdPrice     float64 `json:"usdPrice"`
	NativePrice  float64 `json:"nativePrice"`
	PoolAmount   float64 `json:"poolAmount"`
	CircSupply   float64 `json:"circSupply"`
	TotalSupply  float64 `json:"totalSupply"`
	Fdv          float64 `json:"fdv"`
	Mcap         float64 `json:"mcap"`
	Launchpad    string  `json:"launchpad"`
	TokenProgram string  `json:"tokenProgram"`
	DevMintCount int     `json:"devMintCount"`
	Twitter      string  `json:"twitter,omitempty"`
	Website      string  `json:"website,omitempty"`
	Telegram     string  `json:"telegram,omitempty"`
}

type Stats struct {
	PriceChange float64 `json:"priceChange"`
	BuyVolume   float64 `json:"buyVolume"`
	SellVolume  float64 `json:"sellVolume,omitempty"`
	NumBuys     int     `json:"numBuys"`
	NumSells    int     `json:"numSells,omitempty"`
	NumTraders  int     `json:"numTraders"`
	NumBuyers   int     `json:"numBuyers"`
	NumSellers  int     `json:"numSellers,omitempty"`
}

type Audit struct {
	MintAuthorityDisabled   bool    `json:"mintAuthorityDisabled"`
	FreezeAuthorityDisabled bool    `json:"freezeAuthorityDisabled"`
	TopHoldersPercentage    float64 `json:"topHoldersPercentage"`
	LpBurnedPercentage      float64 `json:"lpBurnedPercentage"`
}

type Pool struct {
	Id           string    `json:"id"`
	Chain        string    `json:"chain"`
	Dex          string    `json:"dex"`
	Type         string    `json:"type"`
	BaseAsset    Asset     `json:"baseAsset"`
	QuoteAsset   Asset     `json:"quoteAsset"`
	Audit        Audit     `json:"audit"`
	CreatedAt    time.Time `json:"createdAt"`
	Liquidity    float64   `json:"liquidity"`
	Stats5M      Stats     `json:"stats5m"`
	Stats1H      Stats     `json:"stats1h"`
	Stats6H      Stats     `json:"stats6h"`
	Stats24H     Stats     `json:"stats24h"`
	BondingCurve float64   `json:"bondingCurve"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PoolList struct {
	Pools []Pool `json:"pools"`
	Total int    `json:"total"`
	Next  int    `json:"next"`
}
