package tikhub

import (
	"context"
	"encoding/json"
	"errors"
)

func (t Client) XhsWebGetNoteInfoByIdV2(ctx context.Context, noteId string) (*Note, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("note_id", noteId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_note_info_v2")

	if err != nil {
		return nil, err
	}

	var resp Response[Note]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	if len(resp.Data.Data) == 0 {
		return nil, errors.New("not found")
	}

	return &resp.Data, nil
}
