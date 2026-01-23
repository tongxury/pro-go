package volcengine

import (
	"context"
	"encoding/json"
	"errors"
	"store/pkg/sdk/conv"
)

// {"ResponseMetadata":{"RequestId":"2025091218261822503CC3C7D04C797E80","Action":"CreateUrlMaterial","Version":"2022-02-01","Service":"iccloud_muse","Region":"cn-north","Code":0},"Result":{"MediaId":"7549146192932864051"}}
// 7549146192932864051
func (t *Client) GetMediaInfo(ctx context.Context, params GetMediaInfoParams) (*GetMediaInfoResult, error) {

	bytes, err := t.doRequest(ctx,
		Req{
			Version: "2022-02-01",
			Action:  "GetMediaInfo",
			Method:  "POST",
			Body: conv.M2B(map[string]any{
				"MediaIds":  params.MediaIds,
				"MediaType": params.MediaType,
			}),
		},
	)
	var resp Response[GetMediaInfoResult]
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return &resp.Result, nil
}

type GetMediaInfoParams struct {
	MediaIds  []string
	MediaType int // 1是素材，2是草稿，3是成片，目前支持成片，素材；默认是查询成片
	//Title       string
	//MaterialUrl string
}

type GetMediaInfoResult struct {
	MediaInfos []MediaInfo `json:"MediaInfos"`
}
