package ape

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type GemsFilters struct {
	MinVolume24h int64 `json:"minVolume24h,omitempty"`
}

type ListGemsParams struct {
	AboutToGraduate GemsFilters `json:"aboutToGraduate,omitempty"`
	New             GemsFilters `json:"new,omitempty"`
	Graduated       GemsFilters `json:"graduated,omitempty"`
}

// 修改参数 注意修改  	return results.Graduated, nil
func (t *Client) ListGems(ctx context.Context, params ListGemsParams) (*PoolList, error) {

	url := "https://api.ape.pro/api/v1/gems"

	//params := ListGemsParams{
	//	Graduated: GemsFilters{
	//		MinVolume24h: 7000000,
	//	},
	//}

	result, err := resty.New().R().
		SetContext(ctx).
		SetBody(params).
		SetHeaders(map[string]string{
			"accept":             "*/*",
			"accept-language":    "zh-CN,zh;q=0.9,ar;q=0.8",
			"cache-control":      "no-cache",
			"origin":             "https://ape.pro",
			"pragma":             "no-cache",
			"priority":           "u=1, i",
			"referer":            "https://ape.pro/",
			"sec-ch-ua":          `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": `"macOS"`,
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-site",
			"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		}).
		Post(url)

	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, fmt.Errorf(result.String())
	}

	var results struct {
		New             *PoolList `json:"new,omitempty"`
		AboutToGraduate *PoolList `json:"aboutToGraduate,omitempty"`
		Graduated       *PoolList `json:"graduated,omitempty"`
	}
	if err := json.Unmarshal(result.Body(), &results); err != nil {
		return nil, err
	}

	return results.Graduated, nil
}
