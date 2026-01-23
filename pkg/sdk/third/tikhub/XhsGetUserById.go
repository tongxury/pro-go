package tikhub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"store/pkg/sdk/conv"
)

func (t Client) XhsGetUserById(ctx context.Context, userId string) (*XhsUserV2, error) {

	user, err := t.xhsWebGetUserInfo(ctx, userId)
	if err == nil {
		return user, nil
	}

	user, err = t.xhsWebGetUserInfoV3(ctx, userId)
	if err == nil {
		return user, nil
	}

	user, err = t.xhsWebV2FetchUserInfoApp(ctx, userId)
	if err == nil {
		return user, nil
	}

	user, err = t.xhsWebV2FetchUserInfo(ctx, userId)
	if err == nil {
		return user, nil
	}

	return nil, errors.New("user not found: " + userId)
}

func (t Client) xhsWebV2FetchUserInfo(ctx context.Context, userId string) (*XhsUserV2, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("user_id", userId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_user_info")
	//Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_user_info_v3")

	if err != nil {
		return nil, err
	}

	fmt.Println(r)

	var resp Response[XhsUser]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	x := resp.Data

	if x.BasicInfo.Nickname == "" {
		return nil, errors.New("not found: " + userId)
	}

	return &XhsUserV2{
		Code: 0,
		Data: XhsUserV2Data{
			Desc:    x.BasicInfo.Desc,
			Fans:    x.GetFollowerCount(),
			Follows: x.GetFollowingCount(),
			//IpLocation: "",
			Liked:    x.GetLikedCount(),
			Nickname: x.BasicInfo.Nickname,
			Notes:    x.GetNoteCount(),
			RedId:    x.BasicInfo.RedId,
			Image:    x.BasicInfo.Images,
			//Tags:       x.Tags,
		},
		Message:    nil,
		RecordTime: "",
	}, nil
}

func (t Client) xhsWebV2FetchUserInfoApp(ctx context.Context, userId string) (*XhsUserV2, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("user_id", userId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_user_info_app")
	//Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_user_info_v3")

	if err != nil {
		return nil, err
	}

	fmt.Println(r)

	var resp Response[XhsWebV2FetchUserInfoAppResult]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	x := resp.Data

	if x.Nickname == "" {
		return nil, errors.New("not found: " + userId)
	}

	return &XhsUserV2{
		Code: 0,
		Data: XhsUserV2Data{
			Desc:    x.Desc,
			Fans:    x.Fans,
			Follows: x.Follows,
			//IpLocation: "",
			Liked:    x.Liked,
			Nickname: x.Nickname,
			Notes:    x.NoteNumStat.Posted,
			RedId:    x.RedId,
			Image:    x.Images,
			//Tags:       x.Tags,
		},
		Message:    nil,
		RecordTime: "",
	}, nil
}

type XhsWebV2FetchUserInfoAppResult struct {
	AvatarLikeStatus bool `json:"avatar_like_status"`
	AvatarPendant    struct {
		CurrentUserPendant bool `json:"current_user_pendant"`
		CurrentUserPet     bool `json:"current_user_pet"`
	} `json:"avatar_pendant"`
	BannerInfo struct {
		BgColor    string `json:"bg_color"`
		Image      string `json:"image"`
		LikeStatus bool   `json:"like_status"`
	} `json:"banner_info"`
	Blocked          bool `json:"blocked"`
	Blocking         bool `json:"blocking"`
	BrandAccountInfo struct {
		Conversions []struct {
			Icon      string `json:"icon"`
			Id        string `json:"id"`
			IsRedShop bool   `json:"is_red_shop"`
			IsShop    bool   `json:"is_shop"`
			Link      string `json:"link"`
			Name      string `json:"name"`
			SubTitle  string `json:"sub_title"`
		} `json:"conversions"`
	} `json:"brand_account_info"`
	Collected            int           `json:"collected"`
	CollectedBookNum     int           `json:"collected_book_num"`
	CollectedBrandNum    int           `json:"collected_brand_num"`
	CollectedMovieNum    int           `json:"collected_movie_num"`
	CollectedNotesNum    int           `json:"collected_notes_num"`
	CollectedPoiNum      int           `json:"collected_poi_num"`
	CollectedProductNum  int           `json:"collected_product_num"`
	CollectedTagsNum     int           `json:"collected_tags_num"`
	CommunityRuleUrl     string        `json:"community_rule_url"`
	DefaultCollectionTab string        `json:"default_collection_tab"`
	Desc                 string        `json:"desc"`
	DescAtUsers          []interface{} `json:"desc_at_users"`
	Fans                 int           `json:"fans"`
	Follows              int           `json:"follows"`
	Fstatus              string        `json:"fstatus"`
	Gender               int           `json:"gender"`
	HulaTabs             struct {
		AllShowTabConfig []struct {
			TabId   string `json:"tab_id"`
			TabName string `json:"tab_name"`
		} `json:"all_show_tab_config"`
		TabIdSelected string `json:"tab_id_selected"`
	} `json:"hula_tabs"`
	IdentityDeeplink      string `json:"identity_deeplink"`
	IdentityLabelMigrated bool   `json:"identity_label_migrated"`
	Imageb                string `json:"imageb"`
	Images                string `json:"images"`
	Interactions          []struct {
		Count     int    `json:"count"`
		IsPrivate bool   `json:"is_private"`
		Name      string `json:"name"`
		Toast     string `json:"toast"`
		Type      string `json:"type"`
	} `json:"interactions"`
	IpLocation              string `json:"ip_location"`
	IsRecommendLevelIllegal bool   `json:"is_recommend_level_illegal"`
	Level                   struct {
		ImageLink string `json:"image_link"`
		Number    int    `json:"number"`
	} `json:"level"`
	Liked        int    `json:"liked"`
	Location     string `json:"location"`
	LocationJump bool   `json:"location_jump"`
	Nboards      int    `json:"nboards"`
	Ndiscovery   int    `json:"ndiscovery"`
	Nickname     string `json:"nickname"`
	NoteNumStat  struct {
		Collected int `json:"collected"`
		Liked     int `json:"liked"`
		Posted    int `json:"posted"`
	} `json:"note_num_stat"`
	RecommendInfo     string `json:"recommend_info"`
	RecommendInfoIcon string `json:"recommend_info_icon"`
	RedClubInfo       struct {
		RedClub      bool   `json:"red_club"`
		RedClubLevel int    `json:"red_club_level"`
		RedClubUrl   string `json:"red_club_url"`
		Redclubscore int    `json:"redclubscore"`
	} `json:"red_club_info"`
	RedId                     string `json:"red_id"`
	RedOfficialVerified       bool   `json:"red_official_verified"`
	RedOfficialVerifyBaseInfo string `json:"red_official_verify_base_info"`
	RedOfficialVerifyContent  string `json:"red_official_verify_content"`
	RedOfficialVerifyType     int    `json:"red_official_verify_type"`
	RemarkName                string `json:"remark_name"`
	Result                    struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Success bool   `json:"success"`
	} `json:"result"`
	SellerInfo struct {
		IsTabGoodsFirst    bool          `json:"is_tab_goods_first"`
		TabCodeNames       []interface{} `json:"tab_code_names"`
		TabGoodsApiVersion int           `json:"tab_goods_api_version"`
	} `json:"seller_info"`
	ShareInfo struct {
		Content string `json:"content"`
		Title   string `json:"title"`
	} `json:"share_info"`
	ShareInfoV2 struct {
		Content string `json:"content"`
		Title   string `json:"title"`
	} `json:"share_info_v2"`
	ShareLink           string `json:"share_link"`
	ShowExtraInfoButton bool   `json:"show_extra_info_button"`
	TabPublic           struct {
		Collection      bool `json:"collection"`
		CollectionBoard bool `json:"collection_board"`
		CollectionNote  bool `json:"collection_note"`
		Seed            bool `json:"seed"`
	} `json:"tab_public"`
	TabVisible struct {
		Collect bool `json:"collect"`
		Like    bool `json:"like"`
		Note    bool `json:"note"`
		Seed    bool `json:"seed"`
	} `json:"tab_visible"`
	Tags         []interface{} `json:"tags"`
	UserDescInfo struct {
		Desc        string        `json:"desc"`
		DescAtUsers []interface{} `json:"desc_at_users"`
	} `json:"user_desc_info"`
	UserRoleType     int    `json:"user_role_type"`
	UserWidgetSwitch bool   `json:"user_widget_switch"`
	Userid           string `json:"userid"`
	ZhongTongBarInfo struct {
		Conversions []struct {
			Icon     string `json:"icon"`
			Id       string `json:"id"`
			Link     string `json:"link"`
			Name     string `json:"name"`
			SubTitle string `json:"sub_title"`
		} `json:"conversions"`
	} `json:"zhong_tong_bar_info"`
}

func (t Client) xhsWebGetUserInfoV2(ctx context.Context, userId string) (*XhsUserV2, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("user_id", userId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		//Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_user_info")
		Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_user_info_v2")

	if err != nil {
		return nil, err
	}
	fmt.Println(r)
	var resp Response[XhsUserV2]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	if resp.Data.Data.Nickname == "" {
		return nil, errors.New("not found: " + userId)
	}

	return &resp.Data, nil
}

func (t Client) xhsWebGetUserInfo(ctx context.Context, userId string) (*XhsUserV2, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("user_id", userId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		//Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_user_info")
		Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_user_info")

	if err != nil {
		return nil, err
	}
	fmt.Println(r)
	var resp Response[XhsUserV2]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	if resp.Data.Data.Nickname == "" {
		return nil, errors.New("not found: " + userId)
	}

	return &resp.Data, nil
}

func (t Client) xhsWebGetUserInfoV3(ctx context.Context, userId string) (*XhsUserV2, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("user_id", userId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		//Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_user_info")
		Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_user_info_v3")

	if err != nil {
		return nil, err
	}
	fmt.Println(r)
	var resp Response[XhsUserV2]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	if resp.Data.Data.Nickname == "" {
		return nil, errors.New("not found: " + userId)
	}

	return &resp.Data, nil
}

func (t *XhsUser) GetFollowingCount() int {
	if t == nil {
		return 0
	}

	for _, x := range t.Interactions {
		if x.Type == "follows" {
			return conv.Int(x.Count)
		}
	}

	return 0
}

func (t *XhsUser) GetFollowerCount() int {
	if t == nil {
		return 0
	}

	for _, x := range t.Interactions {
		if x.Type == "fans" {
			return conv.Int(x.Count)
		}
	}

	return 0
}

func (t *XhsUser) GetLikedCount() int {
	if t == nil {
		return 0
	}

	for _, x := range t.Interactions {
		if x.Type == "interaction" {
			return conv.Int(x.Count)
		}
	}

	return 0
}

func (t *XhsUser) GetNoteCount() int {
	if t == nil {
		return 0
	}

	return t.TabPublic.CollectionNote.Count
}

type XhsUserV2 struct {
	Code       int           `json:"code"`
	Data       XhsUserV2Data `json:"data"`
	Message    interface{}   `json:"message"`
	RecordTime string        `json:"recordTime"`
}

type XhsUserV2Data struct {
	//BannerImage      string `json:"bannerImage"`
	//Boards           int    `json:"boards"`
	//BrandAccountInfo struct {
	//	BannerImage string `json:"bannerImage"`
	//} `json:"brandAccountInfo"`
	//Collected  int    `json:"collected"`
	Desc    string `json:"desc"`
	Fans    int    `json:"fans"`
	Follows int    `json:"follows"`
	//Fstatus    string `json:"fstatus"`
	//Gender     int    `json:"gender"`
	//Id         string `json:"id"`
	Image      string `json:"image"`
	IpLocation string `json:"ip_location"`
	//Level      struct {
	//	ImageLink string `json:"image_link"`
	//	Number    int    `json:"number"`
	//} `json:"level"`
	Liked    int    `json:"liked"`
	Location string `json:"location"`
	Nickname string `json:"nickname"`
	//NoteCollectionIsPublic    string `json:"noteCollectionIsPublic"`
	Notes int `json:"notes"`
	//OfficialVerified          bool   `json:"officialVerified"`
	//OfficialVerifyIcon        string `json:"officialVerifyIcon"`
	//OfficialVerifyName        string `json:"officialVerifyName"`
	//RedOfficialVerifyIconType int    `json:"redOfficialVerifyIconType"`
	//RedOfficialVerifyShowIcon string `json:"redOfficialVerifyShowIcon"`
	//RedOfficialVerifyType     string `json:"redOfficialVerifyType"`
	RedId string `json:"red_id"`
	Tags  []struct {
		Icon    string `json:"icon"`
		Name    string `json:"name"`
		TagType string `json:"tag_type"`
	} `json:"tags"`
	//VerifyContent string `json:"verifyContent"`
}

type XhsUser struct {
	Result struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	BasicInfo struct {
		Images     string `json:"images"`
		RedId      string `json:"red_id"`
		Gender     int    `json:"gender"`
		IpLocation string `json:"ip_location"`
		Desc       string `json:"desc"`
		Imageb     string `json:"imageb"`
		Nickname   string `json:"nickname"`
	} `json:"basic_info"`
	Interactions []struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Count string `json:"count"`
	} `json:"interactions"`
	Tags []struct {
		Icon    string `json:"icon,omitempty"`
		TagType string `json:"tagType"`
		Name    string `json:"name,omitempty"`
	} `json:"tags"`
	TabPublic struct {
		Collection     bool `json:"collection"`
		CollectionNote struct {
			Display bool `json:"display"`
			Lock    bool `json:"lock"`
			Count   int  `json:"count"`
		} `json:"collectionNote"`
		CollectionBoard struct {
			Count   int  `json:"count"`
			Display bool `json:"display"`
			Lock    bool `json:"lock"`
		} `json:"collectionBoard"`
	} `json:"tab_public"`
	ExtraInfo struct {
		Fstatus   string `json:"fstatus"`
		BlockType string `json:"blockType"`
	} `json:"extra_info"`
}
