package tikhub

import (
	"context"
	"encoding/json"
	"errors"
)

type SearchDouyinUserParams struct {
	Keyword string
	Cursor  int
}

func (t Client) SearchDouyinUser(ctx context.Context, params SearchDouyinUserParams) ([]*User, error) {

	r, err := t.c.R().SetContext(ctx).
		SetBody(map[string]interface{}{
			"keyword": params.Keyword,
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Post(t.baseUrl + "/api/v1/douyin/search/fetch_user_search")

	if err != nil {
		return nil, err
	}

	var tmp Response[SearchDouyinUserResult]
	err = json.Unmarshal(r.Body(), &tmp)
	if err != nil {
		return nil, err
	}

	if tmp.Code != 200 {
		return nil, errors.New("invalid search result: " + r.String())
	}

	var users []*User

	for _, x := range tmp.Data.UserList {
		var y User
		err = json.Unmarshal([]byte(x.DynamicPatch.RawData), &y)
		if err != nil {
			continue
		}

		users = append(users, &y)
	}

	return users, nil
}

type User struct {
	IsNicknameRecall     interface{} `json:"is_nickname_recall"`
	IsRedUniqueid        bool        `json:"is_red_uniqueid"`
	IsPrivateLetter      interface{} `json:"is_private_letter"`
	IsRedPhoneNumber     interface{} `json:"is_red_phone_number"`
	PhoneNumberEncrypted interface{} `json:"phone_number_encrypted"`
	UserInfo             struct {
		Uid          string `json:"uid"`
		Nickname     string `json:"nickname"`
		Gender       int    `json:"gender"`
		AvatarLarger struct {
			Uri     string   `json:"uri"`
			UrlList []string `json:"url_list"`
			Width   int      `json:"width"`
			Height  int      `json:"height"`
		} `json:"avatar_larger"`
		AvatarThumb struct {
			Uri     string   `json:"uri"`
			UrlList []string `json:"url_list"`
			Width   int      `json:"width"`
			Height  int      `json:"height"`
		} `json:"avatar_thumb"`
		AvatarMedium struct {
			Uri     string   `json:"uri"`
			UrlList []string `json:"url_list"`
			Width   int      `json:"width"`
			Height  int      `json:"height"`
		} `json:"avatar_medium"`
		FollowStatus           int         `json:"follow_status"`
		FollowerCount          int         `json:"follower_count"`
		WeiboVerify            string      `json:"weibo_verify"`
		CustomVerify           string      `json:"custom_verify"`
		UniqueId               string      `json:"unique_id"`
		VerificationType       int         `json:"verification_type"`
		EnterpriseVerifyReason string      `json:"enterprise_verify_reason"`
		FollowersDetail        interface{} `json:"followers_detail"`
		PlatformSyncInfo       interface{} `json:"platform_sync_info"`
		Secret                 int         `json:"secret"`
		Geofencing             interface{} `json:"geofencing"`
		FollowerStatus         int         `json:"follower_status"`
		CoverUrl               interface{} `json:"cover_url"`
		ItemList               interface{} `json:"item_list"`
		NewStoryCover          interface{} `json:"new_story_cover"`
		IsStar                 bool        `json:"is_star"`
		TypeLabel              interface{} `json:"type_label"`
		AdCoverUrl             interface{} `json:"ad_cover_url"`
		RelativeUsers          interface{} `json:"relative_users"`
		ChaList                interface{} `json:"cha_list"`
		SecUid                 string      `json:"sec_uid"`
		NeedPoints             interface{} `json:"need_points"`
		HomepageBottomToast    interface{} `json:"homepage_bottom_toast"`
		CanSetGeofencing       interface{} `json:"can_set_geofencing"`
		WhiteCoverUrl          interface{} `json:"white_cover_url"`
		UserTags               []struct {
			Description string `json:"description"`
			IconUrl     string `json:"icon_url"`
			Type        string `json:"type"`
		} `json:"user_tags"`
		BanUserFunctions interface{} `json:"ban_user_functions"`
		VersatileDisplay string      `json:"versatile_display"`
		CardEntries      interface{} `json:"card_entries"`
		DisplayInfo      []struct {
			Title string `json:"title,omitempty"`
			Tag   []struct {
				Description string `json:"description"`
				IconUrl     string `json:"icon_url"`
				Type        string `json:"type"`
			} `json:"tag,omitempty"`
			Text string `json:"text,omitempty"`
		} `json:"display_info"`
		LiveStatus                             int         `json:"live_status"`
		CardEntriesNotDisplay                  interface{} `json:"card_entries_not_display"`
		CardSortPriority                       interface{} `json:"card_sort_priority"`
		InterestTags                           interface{} `json:"interest_tags"`
		LinkItemList                           interface{} `json:"link_item_list"`
		UserPermissions                        interface{} `json:"user_permissions"`
		OfflineInfoList                        interface{} `json:"offline_info_list"`
		SignatureExtra                         interface{} `json:"signature_extra"`
		PersonalTagList                        interface{} `json:"personal_tag_list"`
		CfList                                 interface{} `json:"cf_list"`
		ImRoleIds                              interface{} `json:"im_role_ids"`
		NotSeenItemIdList                      interface{} `json:"not_seen_item_id_list"`
		FollowerListSecondaryInformationStruct interface{} `json:"follower_list_secondary_information_struct"`
		EndorsementInfoList                    interface{} `json:"endorsement_info_list"`
		TextExtra                              interface{} `json:"text_extra"`
		ContrailList                           interface{} `json:"contrail_list"`
		DataLabelList                          interface{} `json:"data_label_list"`
		NotSeenItemIdListV2                    interface{} `json:"not_seen_item_id_list_v2"`
		SpecialPeopleLabels                    interface{} `json:"special_people_labels"`
		FamiliarVisitorUser                    interface{} `json:"familiar_visitor_user"`
		AvatarSchemaList                       interface{} `json:"avatar_schema_list"`
		ProfileMobParams                       interface{} `json:"profile_mob_params"`
		VerificationPermissionIds              interface{} `json:"verification_permission_ids"`
		BatchUnfollowRelationDesc              interface{} `json:"batch_unfollow_relation_desc"`
		BatchUnfollowContainTabs               interface{} `json:"batch_unfollow_contain_tabs"`
		CreatorTagList                         interface{} `json:"creator_tag_list"`
		PrivateRelationList                    interface{} `json:"private_relation_list"`
		FollowerCountStr                       string      `json:"follower_count_str"`
	} `json:"user_info"`
}
type SearchDouyinUserResult struct {
	Type     int `json:"type"`
	UserList []struct {
		Position       interface{} `json:"position"`
		UniqidPosition interface{} `json:"uniqid_position"`
		Effects        interface{} `json:"effects"`
		Musics         interface{} `json:"musics"`
		Items          interface{} `json:"items"`
		MixList        interface{} `json:"mix_list"`
		Challenges     interface{} `json:"challenges"`
		DynamicPatch   struct {
			Height     int    `json:"height"`
			Schema     string `json:"schema"`
			OriginType int    `json:"origin_type"`
			RawData    string `json:"raw_data"`
		} `json:"dynamic_patch"`
		ProductInfo     interface{} `json:"product_info"`
		ProductList     interface{} `json:"product_list"`
		Baikes          interface{} `json:"baikes"`
		UserSubLightApp interface{} `json:"userSubLightApp"`
		ShopProductInfo interface{} `json:"shop_product_info"`
		UserServiceInfo interface{} `json:"user_service_info"`
	} `json:"user_list"`
	ChallengeList interface{} `json:"challenge_list"`
	MusicList     interface{} `json:"music_list"`
	Cursor        int         `json:"cursor"`
	HasMore       int         `json:"has_more"`
	StatusCode    int         `json:"status_code"`
	Qc            string      `json:"qc"`
	MyselfUserId  string      `json:"myself_user_id"`
	Rid           string      `json:"rid"`
	LogPb         struct {
		ImprId    string `json:"impr_id"`
		StabExtra struct {
			TrafficIdentification string `json:"traffic_identification"`
		} `json:"stab_extra"`
	} `json:"log_pb"`
	Extra struct {
		Now             int64         `json:"now"`
		Logid           string        `json:"logid"`
		FatalItemIds    []interface{} `json:"fatal_item_ids"`
		SearchRequestId string        `json:"search_request_id"`
		StartMs         int64         `json:"start_ms"`
	} `json:"extra"`
	InputKeyword       string `json:"input_keyword"`
	GlobalDoodleConfig struct {
		Keyword        string `json:"keyword"`
		FilterShowDot  int    `json:"filter_show_dot"`
		FilterSettings []struct {
			Title        string `json:"title"`
			Name         string `json:"name"`
			DefaultIndex int    `json:"default_index"`
			LogName      string `json:"log_name"`
			Items        []struct {
				Title    string `json:"title"`
				Value    string `json:"value"`
				LogValue string `json:"log_value"`
			} `json:"items"`
		} `json:"filter_settings"`
	} `json:"global_doodle_config"`
}
