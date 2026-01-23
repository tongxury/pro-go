package tikhub

import (
	"context"
	"encoding/json"
)

func (t Client) XhsGetNoteByShareUrl(ctx context.Context, url string) (*NoteMeta, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("share_link", url).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/app/extract_share_info")

	if err != nil {
		return nil, err
	}

	var resp Response[NoteMeta]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}
	//https: //www.xiaohongshu.com/explore/680e0c3b000000000b0174e0?app_platform=ios&app_version=8.80.3&share_from_user_hidden=true&xsec_source=app_share&type=video&xsec_token=CBwF3yJI_lLcIBZxrd9OiGtX8YTSeYH-aOVOVqFrIBtwU&author_share=1&xhsshare=CopyLink&shareRedId=N0lGRjw1RTo2NzUyOTgwNjc0OTk3SDg8&apptime=1751704306&share_id=6fc617647ede48e49c3651eedf1f6903
	return &resp.Data, nil
}

type NoteMeta struct {
	NoteId    string `json:"note_id"`
	XsecToken string `json:"xsec_token"`
}
