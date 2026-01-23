package tikhub

import (
	"context"
	"encoding/json"
)

func (t Client) GetTiktokUserProfile(ctx context.Context, uniqueId string) (*UserProfile, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("unique_id", uniqueId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/tiktok/app/v3/handler_user_profile")

	if err != nil {
		return nil, err
	}

	var resp Response[UserProfile]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
