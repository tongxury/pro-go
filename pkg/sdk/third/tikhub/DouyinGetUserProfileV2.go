package tikhub

import (
	"context"
	"encoding/json"
)

func (t Client) DouyinGetUserProfileV2(ctx context.Context, uniqueId string) (*UserProfileV2, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("unique_id", uniqueId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/douyin/web/handler_user_profile_v2")

	if err != nil {
		return nil, err
	}

	var resp Response[UserProfileV2]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (t *UserProfileV2) GetAvatar() string {

	if t == nil {
		return ""
	}

	if len(t.UserInfo.AvatarThumb.UrlList) > 0 {
		return t.UserInfo.AvatarThumb.UrlList[0]
	}

	return ""
}

type UserProfileV2 struct {
	Extra struct {
		Logid string `json:"logid"`
		Now   int64  `json:"now"`
	} `json:"extra"`
	IsOversea  int `json:"is_oversea"`
	StatusCode int `json:"status_code"`
	UserInfo   struct {
		ShortId     string `json:"short_id"`
		Nickname    string `json:"nickname"`
		Signature   string `json:"signature"`
		AvatarThumb struct {
			Uri     string   `json:"uri"`
			UrlList []string `json:"url_list"`
		} `json:"avatar_thumb"`
		AvatarMedium struct {
			Uri     string   `json:"uri"`
			UrlList []string `json:"url_list"`
		} `json:"avatar_medium"`
		FollowStatus     int    `json:"follow_status"`
		AwemeCount       int    `json:"aweme_count"`
		FollowingCount   int    `json:"following_count"`
		FavoritingCount  int    `json:"favoriting_count"`
		TotalFavorited   string `json:"total_favorited"`
		CustomVerify     string `json:"custom_verify"`
		UniqueId         string `json:"unique_id"`
		VerificationType int    `json:"verification_type"`
		OriginalMusician struct {
			MusicCount     int `json:"music_count"`
			MusicUsedCount int `json:"music_used_count"`
			DiggCount      int `json:"digg_count"`
		} `json:"original_musician"`
		EnterpriseVerifyReason  string        `json:"enterprise_verify_reason"`
		MplatformFollowersCount int           `json:"mplatform_followers_count"`
		FollowersDetail         interface{}   `json:"followers_detail"`
		PlatformSyncInfo        interface{}   `json:"platform_sync_info"`
		Geofencing              interface{}   `json:"geofencing"`
		PolicyVersion           interface{}   `json:"policy_version"`
		SecUid                  string        `json:"sec_uid"`
		TypeLabel               interface{}   `json:"type_label"`
		ShowFavoriteList        bool          `json:"show_favorite_list"`
		CardEntries             []interface{} `json:"card_entries"`
		MixInfo                 interface{}   `json:"mix_info"`
		AccountCertInfo         string        `json:"account_cert_info"`
	} `json:"user_info"`
}
