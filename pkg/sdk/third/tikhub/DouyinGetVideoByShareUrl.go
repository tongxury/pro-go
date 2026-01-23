package tikhub

import (
	"context"
	"encoding/json"
	"errors"
)

func (t Client) DouyinGetVideoByShareUrl(ctx context.Context, shareUrl string) (*Video, error) {

	url2, err := t.GetVideoByShareUrl2(ctx, shareUrl)
	if err == nil {
		return url2, err
	}

	return t.GetVideoByShareUrl(ctx, shareUrl)
}

func (t Client) GetVideoByShareUrl2(ctx context.Context, url string) (*Video, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("share_url", url).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/douyin/web/fetch_one_video_by_share_url")
	//Get(t.baseUrl + "/api/v1/douyin/web/fetch_one_video_by_share_url")

	if err != nil {
		return nil, err
	}

	var resp Response[Video]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	urls := resp.Data.GetPlayAddr()

	if len(urls) == 0 {
		return nil, errors.New("playAddr url list is empty")
	}

	return &resp.Data, nil
}

func (t Client) GetVideoByShareUrl(ctx context.Context, url string) (*Video, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("share_url", url).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/douyin/app/v3/fetch_one_video_by_share_url")
	//Get(t.baseUrl + "/api/v1/douyin/web/fetch_one_video_by_share_url")

	if err != nil {
		return nil, err
	}

	var resp Response[Video]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	urls := resp.Data.GetPlayAddr()

	if len(urls) == 0 {
		return nil, errors.New("playAddr url list is empty")
	}
	return &resp.Data, nil
}

func (t Video) GetPlayAddr() []string {

	urls := t.AwemeDetail.Video.PlayAddrH264.UrlList
	if len(urls) == 0 {
		urls = t.AwemeDetail.Video.PlayAddr265.UrlList
	}

	if len(urls) == 0 {
		for _, x := range t.AwemeDetail.Video.BitRate {
			if len(urls) == 0 {
				urls = x.PlayAddr.UrlList
			}
		}
	}

	return urls
}

type Video struct {
	AwemeDetail struct {
		//Is24Story    int    `json:"is_24_story"`
		Desc string `json:"desc"`
		//ImageComment struct {
		//} `json:"image_comment"`
		//MediaType               int         `json:"media_type"`
		//CommentList             interface{} `json:"comment_list"`
		//DuetAggregateInMusicTab bool        `json:"duet_aggregate_in_music_tab"`
		//AuthenticationToken     string      `json:"authentication_token"`
		//SeriesPaidInfo          struct {
		//	SeriesPaidStatus int `json:"series_paid_status"`
		//	ItemPrice        int `json:"item_price"`
		//} `json:"series_paid_info"`
		//AwemeListenStruct struct {
		//	TraceInfo string `json:"trace_info"`
		//} `json:"aweme_listen_struct"`
		//StoryTtl         int  `json:"story_ttl"`
		//HaveDashboard    bool `json:"have_dashboard"`
		//DistributeCircle struct {
		//	DistributeType         int  `json:"distribute_type"`
		//	IsCampus               bool `json:"is_campus"`
		//	CampusBlockInteraction bool `json:"campus_block_interaction"`
		//} `json:"distribute_circle"`
		//ImageInfos    interface{} `json:"image_infos"`
		//XiguaBaseInfo struct {
		//	StarAltarType    int `json:"star_altar_type"`
		//	ItemId           int `json:"item_id"`
		//	Status           int `json:"status"`
		//	StarAltarOrderId int `json:"star_altar_order_id"`
		//} `json:"xigua_base_info"`
		//IsLifeItem        bool        `json:"is_life_item"`
		//PackedClips       interface{} `json:"packed_clips"`
		//IsKaraoke         bool        `json:"is_karaoke"`
		//DistributeType    int         `json:"distribute_type"`
		//ChallengePosition interface{} `json:"challenge_position"`
		//City              string      `json:"city"`
		//ItemReact         int         `json:"item_react"`
		//IsVr              bool        `json:"is_vr"`
		//HybridLabel       interface{} `json:"hybrid_label"`
		//Position          interface{} `json:"position"`
		//ShareUrl          string      `json:"share_url"`
		//FallCardStruct    struct {
		//	RecommendReasonV2 string `json:"recommend_reason_v2"`
		//} `json:"fall_card_struct"`
		//ItemDuet          int `json:"item_duet"`
		//FeedCommentConfig struct {
		//	InputConfigText   string `json:"input_config_text"`
		//	AuthorAuditStatus int    `json:"author_audit_status"`
		//	ImageChallenge    struct {
		//		InnerPublishText string `json:"inner_publish_text"`
		//	} `json:"image_challenge"`
		//} `json:"feed_comment_config"`
		//IsMomentStory        int `json:"is_moment_story"`
		//VideoShareEditStatus int `json:"video_share_edit_status"`
		//ComponentControl     struct {
		//} `json:"component_control"`
		//Anchors          interface{} `json:"anchors"`
		//IsPgcshow        bool        `json:"is_pgcshow"`
		//ShowFollowButton struct {
		//} `json:"show_follow_button"`
		//EntertainmentProductInfo struct {
		//	MarketInfo struct {
		//		LimitFree struct {
		//			InFree bool `json:"in_free"`
		//		} `json:"limit_free"`
		//	} `json:"market_info"`
		//} `json:"entertainment_product_info"`
		//CanCacheToLocal           bool   `json:"can_cache_to_local"`
		//Caption                   string `json:"caption"`
		//AwemeTypeTags             string `json:"aweme_type_tags"`
		//EntertainmentVideoPaidWay struct {
		//	PaidWays            []interface{} `json:"paid_ways"`
		//	PaidType            int           `json:"paid_type"`
		//	EnableUseNewEntData bool          `json:"enable_use_new_ent_data"`
		//} `json:"entertainment_video_paid_way"`
		//ShootWay     string `json:"shoot_way"`
		//IsFromAdAuth bool   `json:"is_from_ad_auth"`
		//Statistics   struct {
		//	CommentCount       int    `json:"comment_count"`
		//	DiggCount          int    `json:"digg_count"`
		//	LoseCount          int    `json:"lose_count"`
		//	AwemeId            string `json:"aweme_id"`
		//	LiveWatchCount     int    `json:"live_watch_count"`
		//	Digest             string `json:"digest"`
		//	LoseCommentCount   int    `json:"lose_comment_count"`
		//	WhatsappShareCount int    `json:"whatsapp_share_count"`
		//	CollectCount       int    `json:"collect_count"`
		//	ShareCount         int    `json:"share_count"`
		//	ForwardCount       int    `json:"forward_count"`
		//	PlayCount          int    `json:"play_count"`
		//	ExposureCount      int    `json:"exposure_count"`
		//	AdmireCount        int    `json:"admire_count"`
		//	DownloadCount      int    `json:"download_count"`
		//} `json:"statistics"`
		//ImageAlbumMusicInfo struct {
		//	EndTime   int `json:"end_time"`
		//	Volume    int `json:"volume"`
		//	BeginTime int `json:"begin_time"`
		//} `json:"image_album_music_info"`
		//PressPanelInfo   string `json:"press_panel_info"`
		//MiscInfo         string `json:"misc_info"`
		//WithoutWatermark bool   `json:"without_watermark"`
		//DescLanguage     string `json:"desc_language"`
		//LibfinsertTaskId string `json:"libfinsert_task_id"`
		//Original         int    `json:"original"`
		//ItemTitle        string `json:"item_title"`
		//VideoTag         []struct {
		//	TagId   int    `json:"tag_id"`
		//	TagName string `json:"tag_name"`
		//	Level   int    `json:"level"`
		//} `json:"video_tag"`
		Author struct {
			//	Geofencing       []interface{} `json:"geofencing"`
			//	IsBindedWeibo    bool          `json:"is_binded_weibo"`
			//	AppleAccount     int           `json:"apple_account"`
			//	ShowImageBubble  bool          `json:"show_image_bubble"`
			//	NewStoryCover    interface{}   `json:"new_story_cover"`
			//	TypeLabel        interface{}   `json:"type_label"`
			//	Birthday         string        `json:"birthday"`
			//	IsBlockingV2     bool          `json:"is_blocking_v2"`
			//	AccountRegion    string        `json:"account_region"`
			//	DownloadPromptTs int           `json:"download_prompt_ts"`
			//	DisplayInfo      interface{}   `json:"display_info"`
			//	CardSortPriority interface{}   `json:"card_sort_priority"`
			//	Location         string        `json:"location"`
			//	HasYoutubeToken  bool          `json:"has_youtube_token"`
			//	GoogleAccount    string        `json:"google_account"`
			//	TwExpireTime     int           `json:"tw_expire_time"`
			//	CreateTime       int           `json:"create_time"`
			//	Avatar168X168    struct {
			//		Uri     string   `json:"uri"`
			//		Height  int      `json:"height"`
			//		UrlList []string `json:"url_list"`
			//		Width   int      `json:"width"`
			//	} `json:"avatar_168x168"`
			//	LiveAgreement       int         `json:"live_agreement"`
			//	DuetSetting         int         `json:"duet_setting"`
			//	TotalFavorited      int         `json:"total_favorited"`
			//	IsVerified          bool        `json:"is_verified"`
			//	TextExtra           interface{} `json:"text_extra"`
			//	WhiteCoverUrl       interface{} `json:"white_cover_url"`
			//	LinkItemList        interface{} `json:"link_item_list"`
			//	UserTags            interface{} `json:"user_tags"`
			//	UserRate            int         `json:"user_rate"`
			//	VerifyInfo          string      `json:"verify_info"`
			//	SignatureExtra      interface{} `json:"signature_extra"`
			//	OfflineInfoList     interface{} `json:"offline_info_list"`
			//	AuthorityStatus     int         `json:"authority_status"`
			//	SpecialFollowStatus int         `json:"special_follow_status"`
			//	EnableNearbyVisible bool        `json:"enable_nearby_visible"`
			//	WeiboVerify         string      `json:"weibo_verify"`
			//	HasFacebookToken    bool        `json:"has_facebook_token"`
			//	MaxFollowerCount    int         `json:"max_follower_count"`
			//	ContrailList        interface{} `json:"contrail_list"`
			//	ChaList             interface{} `json:"cha_list"`
			//	IsStar              bool        `json:"is_star"`
			//	SpecialLock         int         `json:"special_lock"`
			//	CommentSetting      int         `json:"comment_setting"`
			//	YoutubeChannelId    string      `json:"youtube_channel_id"`
			//	DataLabelList       interface{} `json:"data_label_list"`
			//	AdCoverUrl          interface{} `json:"ad_cover_url"`
			//	ShortId             string      `json:"short_id"`
			//	UserNotSee          int         `json:"user_not_see"`
			AwemeCount int `json:"aweme_count"`
			//	LiveVerify          int         `json:"live_verify"`
			//	ReflowPageGid       int         `json:"reflow_page_gid"`
			//	CustomVerify        string      `json:"custom_verify"`
			//	CfList              interface{} `json:"cf_list"`
			//	Uid                 string      `json:"uid"`
			Avatar300X300 struct {
				UrlList []string `json:"url_list"`
				//Width   int      `json:"width"`
				//Uri     string   `json:"uri"`
				//Height  int      `json:"height"`
			} `json:"avatar_300x300"`
			//	WithCommerceEntry     bool        `json:"with_commerce_entry"`
			//	MateAddPermission     int         `json:"mate_add_permission"`
			//	AwemehtsGreetInfo     string      `json:"awemehts_greet_info"`
			//	PersonalTagList       interface{} `json:"personal_tag_list"`
			//	CommerceUserLevel     int         `json:"commerce_user_level"`
			//	ImRoleIds             interface{} `json:"im_role_ids"`
			//	SpecialPeopleLabels   interface{} `json:"special_people_labels"`
			//	IsGovMediaVip         bool        `json:"is_gov_media_vip"`
			//	ShowNearbyActive      bool        `json:"show_nearby_active"`
			//	ShieldFollowNotice    int         `json:"shield_follow_notice"`
			//	YoutubeExpireTime     int         `json:"youtube_expire_time"`
			//	HasEmail              bool        `json:"has_email"`
			//	KyOnlyPredict         int         `json:"ky_only_predict"`
			//	CloseFriendType       int         `json:"close_friend_type"`
			//	AvatarUri             string      `json:"avatar_uri"`
			//	SignatureDisplayLines int         `json:"signature_display_lines"`
			//	AvatarThumb           struct {
			//		UrlList []string `json:"url_list"`
			//		Height  int      `json:"height"`
			//		Width   int      `json:"width"`
			//		Uri     string   `json:"uri"`
			//	} `json:"avatar_thumb"`
			//	FollowerCount                          int           `json:"follower_count"`
			//	Region                                 string        `json:"region"`
			//	SecUid                                 string        `json:"sec_uid"`
			//	UniqueIdModifyTime                     int           `json:"unique_id_modify_time"`
			//	Language                               string        `json:"language"`
			//	BanUserFunctions                       []interface{} `json:"ban_user_functions"`
			//	WithFusionShopEntry                    bool          `json:"with_fusion_shop_entry"`
			//	SchoolCategory                         int           `json:"school_category"`
			//	FollowerListSecondaryInformationStruct interface{}   `json:"follower_list_secondary_information_struct"`
			//	SyncToToutiao                          int           `json:"sync_to_toutiao"`
			//	PreventDownload                        bool          `json:"prevent_download"`
			//	ReactSetting                           int           `json:"react_setting"`
			//	Gender                                 int           `json:"gender"`
			//	FollowerStatus                         int           `json:"follower_status"`
			//	CoverUrl                               []struct {
			//		Uri     string   `json:"uri"`
			//		UrlList []string `json:"url_list"`
			//		Height  int      `json:"height"`
			//		Width   int      `json:"width"`
			//	} `json:"cover_url"`
			//	IsAdFake            bool        `json:"is_ad_fake"`
			//	AwemeHotsoonAuth    int         `json:"aweme_hotsoon_auth"`
			//	WithDouEntry        bool        `json:"with_dou_entry"`
			//	UserAge             int         `json:"user_age"`
			//	WeiboName           string      `json:"weibo_name"`
			//	TwitterId           string      `json:"twitter_id"`
			//	WeiboUrl            string      `json:"weibo_url"`
			//	AccountCertInfo     string      `json:"account_cert_info"`
			//	IsCf                int         `json:"is_cf"`
			//	ItemList            interface{} `json:"item_list"`
			//	Status              int         `json:"status"`
			Nickname string `json:"nickname"`
			//	ShieldDiggNotice    int         `json:"shield_digg_notice"`
			//	WeiboSchema         string      `json:"weibo_schema"`
			//	InterestTags        interface{} `json:"interest_tags"`
			//	UserPeriod          int         `json:"user_period"`
			//	SchoolName          string      `json:"school_name"`
			//	YoutubeChannelTitle string      `json:"youtube_channel_title"`
			//	StitchSetting       int         `json:"stitch_setting"`
			//	ReflowPageUid       int         `json:"reflow_page_uid"`
			//	SchoolId            string      `json:"school_id"`
			//	FollowersDetail     interface{} `json:"followers_detail"`
			//	FbExpireTime        int         `json:"fb_expire_time"`
			//	VerificationType    int         `json:"verification_type"`
			//	EndorsementInfoList interface{} `json:"endorsement_info_list"`
			//	CanSetGeofencing    interface{} `json:"can_set_geofencing"`
			//	LiveHighValue       int         `json:"live_high_value"`
			//	UserCanceled        bool        `json:"user_canceled"`
			//	CvLevel             string      `json:"cv_level"`
			//	NeiguangShield      int         `json:"neiguang_shield"`
			//	CardEntries         interface{} `json:"card_entries"`
			//	StoryCount          int         `json:"story_count"`
			//	BindPhone           string      `json:"bind_phone"`
			//	AvatarMedium        struct {
			//		UrlList []string `json:"url_list"`
			//		Height  int      `json:"height"`
			//		Width   int      `json:"width"`
			//		Uri     string   `json:"uri"`
			//	} `json:"avatar_medium"`
			//	CommentFilterStatus int    `json:"comment_filter_status"`
			InsId string `json:"ins_id"`
			//	ShareQrcodeUri      string `json:"share_qrcode_uri"`
			//	IsPhoneBinded       bool   `json:"is_phone_binded"`
			//	RiskNoticeText      string `json:"risk_notice_text"`
			//	VideoIcon           struct {
			//		Width   int           `json:"width"`
			//		UrlList []interface{} `json:"url_list"`
			//		Height  int           `json:"height"`
			//		Uri     string        `json:"uri"`
			//	} `json:"video_icon"`
			//	HasTwitterToken  bool        `json:"has_twitter_token"`
			//	UserNotShow      int         `json:"user_not_show"`
			//	DownloadSetting  int         `json:"download_setting"`
			//	NeedPoints       interface{} `json:"need_points"`
			//	PlatformSyncInfo interface{} `json:"platform_sync_info"`
			//	SearchImpr       struct {
			//		EntityId string `json:"entity_id"`
			//	} `json:"search_impr"`
			//	IsNotShow           bool        `json:"is_not_show"`
			//	HideSearch          bool        `json:"hide_search"`
			//	ContactsStatus      int         `json:"contacts_status"`
			//	FavoritingCount     int         `json:"favoriting_count"`
			//	WithShopEntry       bool        `json:"with_shop_entry"`
			//	HasInsights         bool        `json:"has_insights"`
			//	NeedRecommend       int         `json:"need_recommend"`
			//	UserMode            int         `json:"user_mode"`
			//	LiveStatus          int         `json:"live_status"`
			//	HasUnreadStory      bool        `json:"has_unread_story"`
			//	HideLocation        bool        `json:"hide_location"`
			//	HomepageBottomToast interface{} `json:"homepage_bottom_toast"`
			//	AvatarLarger        struct {
			//		Uri     string   `json:"uri"`
			//		UrlList []string `json:"url_list"`
			//		Width   int      `json:"width"`
			//		Height  int      `json:"height"`
			//	} `json:"avatar_larger"`
			//	AwemeControl struct {
			//		CanComment     bool `json:"can_comment"`
			//		CanShowComment bool `json:"can_show_comment"`
			//		CanForward     bool `json:"can_forward"`
			//		CanShare       bool `json:"can_share"`
			//	} `json:"aweme_control"`
			//	RelativeUsers          interface{} `json:"relative_users"`
			Signature string `json:"signature"`
			//	LiveCommerce           bool        `json:"live_commerce"`
			//	ShieldCommentNotice    int         `json:"shield_comment_notice"`
			//	FollowingCount         int         `json:"following_count"`
			//	UniqueId               string      `json:"unique_id"`
			//	IsBlockedV2            bool        `json:"is_blocked_v2"`
			//	LiveAgreementTime      int         `json:"live_agreement_time"`
			//	AcceptPrivatePolicy    bool        `json:"accept_private_policy"`
			//	RoomId                 int         `json:"room_id"`
			//	EnterpriseVerifyReason string      `json:"enterprise_verify_reason"`
			//	UserPermissions        interface{} `json:"user_permissions"`
			//	IsDisciplineMember     bool        `json:"is_discipline_member"`
			//	TwitterName            string      `json:"twitter_name"`
			//	Constellation          int         `json:"constellation"`
			//	ShareInfo              struct {
			//		ShareUrl       string `json:"share_url"`
			//		ShareTitle     string `json:"share_title"`
			//		ShareQrcodeUrl struct {
			//			Uri     string   `json:"uri"`
			//			Width   int      `json:"width"`
			//			UrlList []string `json:"url_list"`
			//			Height  int      `json:"height"`
			//		} `json:"share_qrcode_url"`
			//		ShareTitleOther  string `json:"share_title_other"`
			//		ShareDescInfo    string `json:"share_desc_info"`
			//		ShareWeiboDesc   string `json:"share_weibo_desc"`
			//		ShareDesc        string `json:"share_desc"`
			//		ShareTitleMyself string `json:"share_title_myself"`
			//	} `json:"share_info"`
			//	IsBlock               bool        `json:"is_block"`
			//	IsMixUser             bool        `json:"is_mix_user"`
			//	StoryOpen             bool        `json:"story_open"`
			//	SchoolPoiId           string      `json:"school_poi_id"`
			//	SchoolType            int         `json:"school_type"`
			//	CardEntriesNotDisplay interface{} `json:"card_entries_not_display"`
			//	FollowStatus          int         `json:"follow_status"`
			//	FollowerRequestStatus int         `json:"follower_request_status"`
			//	HasOrders             bool        `json:"has_orders"`
			//	Secret                int         `json:"secret"`
		} `json:"author"`
		//AwemeControl struct {
		//	CanComment     bool `json:"can_comment"`
		//	CanShowComment bool `json:"can_show_comment"`
		//	CanForward     bool `json:"can_forward"`
		//	CanShare       bool `json:"can_share"`
		//} `json:"aweme_control"`
		//IsCollectsSelected int `json:"is_collects_selected"`
		//ChaList            []struct {
		//	ShowItems  interface{} `json:"show_items"`
		//	Cid        string      `json:"cid"`
		//	Schema     string      `json:"schema"`
		//	SearchImpr struct {
		//		EntityId string `json:"entity_id"`
		//	} `json:"search_impr"`
		//	ShareInfo struct {
		//		ShareSignatureUrl  string `json:"share_signature_url"`
		//		ShareUrl           string `json:"share_url"`
		//		ShareQuote         string `json:"share_quote"`
		//		ShareTitleMyself   string `json:"share_title_myself"`
		//		ShareDescInfo      string `json:"share_desc_info"`
		//		ShareWeiboDesc     string `json:"share_weibo_desc"`
		//		ShareTitle         string `json:"share_title"`
		//		ShareSignatureDesc string `json:"share_signature_desc"`
		//		ShareTitleOther    string `json:"share_title_other"`
		//		ShareDesc          string `json:"share_desc"`
		//		BoolPersist        int    `json:"bool_persist"`
		//	} `json:"share_info"`
		//	CollectStat int `json:"collect_stat"`
		//	ExtraAttr   struct {
		//		IsLive bool `json:"is_live"`
		//	} `json:"extra_attr"`
		//	BannerList  interface{} `json:"banner_list"`
		//	ChaAttrs    interface{} `json:"cha_attrs"`
		//	Desc        string      `json:"desc"`
		//	Type        int         `json:"type"`
		//	IsPgcshow   bool        `json:"is_pgcshow"`
		//	IsChallenge int         `json:"is_challenge"`
		//	IsCommerce  bool        `json:"is_commerce"`
		//	SubType     int         `json:"sub_type"`
		//	ChaName     string      `json:"cha_name"`
		//	ViewCount   int         `json:"view_count"`
		//	Author      struct {
		//		DataLabelList   interface{}   `json:"data_label_list"`
		//		Status          int           `json:"status"`
		//		HasEmail        bool          `json:"has_email"`
		//		OfflineInfoList interface{}   `json:"offline_info_list"`
		//		AvatarUri       string        `json:"avatar_uri"`
		//		CoverUrl        []interface{} `json:"cover_url"`
		//		ItemList        interface{}   `json:"item_list"`
		//		Nickname        string        `json:"nickname"`
		//		SearchImpr      struct {
		//			EntityId string `json:"entity_id"`
		//		} `json:"search_impr"`
		//		Avatar300X300 struct {
		//			UrlList []interface{} `json:"url_list"`
		//			Height  int           `json:"height"`
		//			Uri     string        `json:"uri"`
		//			Width   int           `json:"width"`
		//		} `json:"avatar_300x300"`
		//		PlatformSyncInfo    interface{} `json:"platform_sync_info"`
		//		UserPermissions     interface{} `json:"user_permissions"`
		//		RelativeUsers       interface{} `json:"relative_users"`
		//		DisplayInfo         interface{} `json:"display_info"`
		//		EndorsementInfoList interface{} `json:"endorsement_info_list"`
		//		ContrailList        interface{} `json:"contrail_list"`
		//		HomepageBottomToast interface{} `json:"homepage_bottom_toast"`
		//		LinkItemList        interface{} `json:"link_item_list"`
		//		CreateTime          int         `json:"create_time"`
		//		BanUserFunctions    interface{} `json:"ban_user_functions"`
		//		IsPhoneBinded       bool        `json:"is_phone_binded"`
		//		TypeLabel           interface{} `json:"type_label"`
		//		AwemeControl        struct {
		//		} `json:"aweme_control"`
		//		CfList       interface{} `json:"cf_list"`
		//		FollowStatus int         `json:"follow_status"`
		//		IsBlock      bool        `json:"is_block"`
		//		SecUid       string      `json:"sec_uid"`
		//		VideoIcon    struct {
		//			UrlList []interface{} `json:"url_list"`
		//			Width   int           `json:"width"`
		//			Uri     string        `json:"uri"`
		//			Height  int           `json:"height"`
		//		} `json:"video_icon"`
		//		UniqueIdModifyTime int         `json:"unique_id_modify_time"`
		//		Constellation      int         `json:"constellation"`
		//		Uid                string      `json:"uid"`
		//		ChaList            interface{} `json:"cha_list"`
		//		AdCoverUrl         interface{} `json:"ad_cover_url"`
		//		FollowersDetail    interface{} `json:"followers_detail"`
		//		UserTags           interface{} `json:"user_tags"`
		//		WithDouEntry       bool        `json:"with_dou_entry"`
		//		ShowImageBubble    bool        `json:"show_image_bubble"`
		//		BindPhone          string      `json:"bind_phone"`
		//		InterestTags       interface{} `json:"interest_tags"`
		//		AvatarLarger       struct {
		//			Uri     string        `json:"uri"`
		//			Height  int           `json:"height"`
		//			UrlList []interface{} `json:"url_list"`
		//			Width   int           `json:"width"`
		//		} `json:"avatar_larger"`
		//		FollowerListSecondaryInformationStruct interface{} `json:"follower_list_secondary_information_struct"`
		//		WhiteCoverUrl                          interface{} `json:"white_cover_url"`
		//		Region                                 string      `json:"region"`
		//		Language                               string      `json:"language"`
		//		AvatarMedium                           struct {
		//			Uri     string        `json:"uri"`
		//			UrlList []interface{} `json:"url_list"`
		//			Width   int           `json:"width"`
		//			Height  int           `json:"height"`
		//		} `json:"avatar_medium"`
		//		Avatar168X168 struct {
		//			Uri     string        `json:"uri"`
		//			UrlList []interface{} `json:"url_list"`
		//			Height  int           `json:"height"`
		//			Width   int           `json:"width"`
		//		} `json:"avatar_168x168"`
		//		ImRoleIds             interface{} `json:"im_role_ids"`
		//		TextExtra             interface{} `json:"text_extra"`
		//		CardEntriesNotDisplay interface{} `json:"card_entries_not_display"`
		//		NewStoryCover         interface{} `json:"new_story_cover"`
		//		UniqueId              string      `json:"unique_id"`
		//		AvatarThumb           struct {
		//			Uri     string        `json:"uri"`
		//			Height  int           `json:"height"`
		//			Width   int           `json:"width"`
		//			UrlList []interface{} `json:"url_list"`
		//		} `json:"avatar_thumb"`
		//		CardEntries         interface{} `json:"card_entries"`
		//		CardSortPriority    interface{} `json:"card_sort_priority"`
		//		ShortId             string      `json:"short_id"`
		//		CanSetGeofencing    interface{} `json:"can_set_geofencing"`
		//		SignatureExtra      interface{} `json:"signature_extra"`
		//		PersonalTagList     interface{} `json:"personal_tag_list"`
		//		SpecialPeopleLabels interface{} `json:"special_people_labels"`
		//		Signature           string      `json:"signature"`
		//		Birthday            string      `json:"birthday"`
		//		Gender              int         `json:"gender"`
		//		NeedPoints          interface{} `json:"need_points"`
		//	} `json:"author"`
		//	HashtagProfile string        `json:"hashtag_profile"`
		//	UserCount      int           `json:"user_count"`
		//	ConnectMusic   []interface{} `json:"connect_music"`
		//} `json:"cha_list"`
		//GuideBtnType        int         `json:"guide_btn_type"`
		//Rate                int         `json:"rate"`
		//CoverLabels         interface{} `json:"cover_labels"`
		//FriendRecommendInfo struct {
		//	FriendRecommendSource int `json:"friend_recommend_source"`
		//} `json:"friend_recommend_info"`
		//VisualSearchInfo struct {
		//	IsShowImgEntrance  bool `json:"is_show_img_entrance"`
		//	IsEcomImg          bool `json:"is_ecom_img"`
		//	IsHighAccuracyEcom bool `json:"is_high_accuracy_ecom"`
		//	IsHighRecallEcom   bool `json:"is_high_recall_ecom"`
		//} `json:"visual_search_info"`
		//IsMomentHistory int `json:"is_moment_history"`
		//BodydanceScore  int `json:"bodydance_score"`
		//PoiPatchInfo    struct {
		//	ItemPatchPoiPromptMark int    `json:"item_patch_poi_prompt_mark"`
		//	Extra                  string `json:"extra"`
		//} `json:"poi_patch_info"`
		//WithPromotionalMusic bool        `json:"with_promotional_music"`
		//SortLabel            string      `json:"sort_label"`
		//CommerceConfigData   interface{} `json:"commerce_config_data"`
		//ItemWarnNotification struct {
		//	Content string `json:"content"`
		//	Show    bool   `json:"show"`
		//	Type    int    `json:"type"`
		//} `json:"item_warn_notification"`
		//IsFirstVideo     bool        `json:"is_first_video"`
		//OriginCommentIds interface{} `json:"origin_comment_ids"`
		//DanmakuControl   struct {
		//	Activities []struct {
		//		Id   int `json:"id"`
		//		Type int `json:"type"`
		//	} `json:"activities"`
		//	PassThroughParams  string `json:"pass_through_params"`
		//	SmartModeDecision  int    `json:"smart_mode_decision"`
		//	PostPrivilegeLevel int    `json:"post_privilege_level"`
		//	PostDeniedReason   string `json:"post_denied_reason"`
		//	DanmakuCnt         int    `json:"danmaku_cnt"`
		//	EnableDanmaku      bool   `json:"enable_danmaku"`
		//	IsPostDenied       bool   `json:"is_post_denied"`
		//	SkipDanmaku        bool   `json:"skip_danmaku"`
		//} `json:"danmaku_control"`
		//GameTagInfo struct {
		//	IsGame bool `json:"is_game"`
		//} `json:"game_tag_info"`
		//ImageList           interface{} `json:"image_list"`
		//UserRecommendStatus int         `json:"user_recommend_status"`
		//NicknamePosition    interface{} `json:"nickname_position"`
		//Music               struct {
		//	Album             string `json:"album"`
		//	BindedChallengeId int    `json:"binded_challenge_id"`
		//	CoverHd           struct {
		//		Uri     string   `json:"uri"`
		//		Width   int      `json:"width"`
		//		UrlList []string `json:"url_list"`
		//		Height  int      `json:"height"`
		//	} `json:"cover_hd"`
		//	PgcMusicType    int  `json:"pgc_music_type"`
		//	IsCommerceMusic bool `json:"is_commerce_music"`
		//	AvatarThumb     struct {
		//		Height  int      `json:"height"`
		//		Width   int      `json:"width"`
		//		Uri     string   `json:"uri"`
		//		UrlList []string `json:"url_list"`
		//	} `json:"avatar_thumb"`
		//	LyricShortPosition interface{} `json:"lyric_short_position"`
		//	Id                 int64       `json:"id"`
		//	DspStatus          int         `json:"dsp_status"`
		//	SearchImpr         struct {
		//		EntityId string `json:"entity_id"`
		//	} `json:"search_impr"`
		//	AuditionDuration int `json:"audition_duration"`
		//	EndTime          int `json:"end_time"`
		//	AvatarLarge      struct {
		//		Width   int      `json:"width"`
		//		UrlList []string `json:"url_list"`
		//		Height  int      `json:"height"`
		//		Uri     string   `json:"uri"`
		//	} `json:"avatar_large"`
		//	OwnerHandle       string        `json:"owner_handle"`
		//	Artists           []interface{} `json:"artists"`
		//	IsMatchedMetadata bool          `json:"is_matched_metadata"`
		//	AuthorPosition    interface{}   `json:"author_position"`
		//	UnshelveCountries interface{}   `json:"unshelve_countries"`
		//	Extra             string        `json:"extra"`
		//	IsOriginalSound   bool          `json:"is_original_sound"`
		//	ExternalSongInfo  []interface{} `json:"external_song_info"`
		//	OwnerId           string        `json:"owner_id"`
		//	SourcePlatform    int           `json:"source_platform"`
		//	CoverMedium       struct {
		//		Uri     string   `json:"uri"`
		//		UrlList []string `json:"url_list"`
		//		Height  int      `json:"height"`
		//		Width   int      `json:"width"`
		//	} `json:"cover_medium"`
		//	CanBackgroundPlay bool        `json:"can_background_play"`
		//	MusicianUserInfos interface{} `json:"musician_user_infos"`
		//	MusicChartRanks   interface{} `json:"music_chart_ranks"`
		//	VideoDuration     int         `json:"video_duration"`
		//	ShootDuration     int         `json:"shoot_duration"`
		//	IsDelVideo        bool        `json:"is_del_video"`
		//	ReasonType        int         `json:"reason_type"`
		//	Title             string      `json:"title"`
		//	PlayUrl           struct {
		//		UrlList []string `json:"url_list"`
		//		UrlKey  string   `json:"url_key"`
		//		Height  int      `json:"height"`
		//		Uri     string   `json:"uri"`
		//		Width   int      `json:"width"`
		//	} `json:"play_url"`
		//	AvatarMedium struct {
		//		UrlList []string `json:"url_list"`
		//		Height  int      `json:"height"`
		//		Uri     string   `json:"uri"`
		//		Width   int      `json:"width"`
		//	} `json:"avatar_medium"`
		//	TagList              interface{} `json:"tag_list"`
		//	Position             interface{} `json:"position"`
		//	Mid                  string      `json:"mid"`
		//	UserCount            int         `json:"user_count"`
		//	IsAudioUrlWithCookie bool        `json:"is_audio_url_with_cookie"`
		//	Duration             int         `json:"duration"`
		//	CoverThumb           struct {
		//		Height  int      `json:"height"`
		//		UrlList []string `json:"url_list"`
		//		Width   int      `json:"width"`
		//		Uri     string   `json:"uri"`
		//	} `json:"cover_thumb"`
		//	MuteShare                      bool   `json:"mute_share"`
		//	IsRestricted                   bool   `json:"is_restricted"`
		//	PreventDownload                bool   `json:"prevent_download"`
		//	OwnerNickname                  string `json:"owner_nickname"`
		//	Author                         string `json:"author"`
		//	PreviewEndTime                 int    `json:"preview_end_time"`
		//	StartTime                      int    `json:"start_time"`
		//	PreventItemDownloadStatus      int    `json:"prevent_item_download_status"`
		//	IsVideoSelfSee                 bool   `json:"is_video_self_see"`
		//	SecUid                         string `json:"sec_uid"`
		//	IdStr                          string `json:"id_str"`
		//	MusicCoverAtmosphereColorValue string `json:"music_cover_atmosphere_color_value"`
		//	StrongBeatUrl                  struct {
		//		Width   int      `json:"width"`
		//		Uri     string   `json:"uri"`
		//		UrlList []string `json:"url_list"`
		//		Height  int      `json:"height"`
		//	} `json:"strong_beat_url"`
		//	MusicStatus int `json:"music_status"`
		//	CoverLarge  struct {
		//		Uri     string   `json:"uri"`
		//		UrlList []string `json:"url_list"`
		//		Height  int      `json:"height"`
		//		Width   int      `json:"width"`
		//	} `json:"cover_large"`
		//	ArtistUserInfos interface{} `json:"artist_user_infos"`
		//	Status          int         `json:"status"`
		//	AuthorDeleted   bool        `json:"author_deleted"`
		//	IsPgc           bool        `json:"is_pgc"`
		//	AuthorStatus    int         `json:"author_status"`
		//	MusicImageBeats struct {
		//		MusicImageBeatsUrl struct {
		//			Height  int      `json:"height"`
		//			Uri     string   `json:"uri"`
		//			UrlList []string `json:"url_list"`
		//			Width   int      `json:"width"`
		//		} `json:"music_image_beats_url"`
		//	} `json:"music_image_beats"`
		//	MusicCollectCount int    `json:"music_collect_count"`
		//	PreviewStartTime  int    `json:"preview_start_time"`
		//	DmvAutoShow       bool   `json:"dmv_auto_show"`
		//	Redirect          bool   `json:"redirect"`
		//	IsOriginal        bool   `json:"is_original"`
		//	OfflineDesc       string `json:"offline_desc"`
		//	CollectStat       int    `json:"collect_stat"`
		//	SchemaUrl         string `json:"schema_url"`
		//} `json:"music"`
		//PreviewVideoStatus     int  `json:"preview_video_status"`
		//IsUseMusic             bool `json:"is_use_music"`
		//DislikeDimensionListV2 []struct {
		//	Icon    string `json:"icon"`
		//	Text    string `json:"text"`
		//	Entitys []struct {
		//		DimensionId int    `json:"dimension_id"`
		//		PreText     string `json:"pre_text"`
		//		SubText     string `json:"sub_text,omitempty"`
		//		ServerExtra string `json:"server_extra"`
		//		EnName      string `json:"en_name"`
		//	} `json:"entitys"`
		//} `json:"dislike_dimension_list_v2"`
		//AdmireAuth struct {
		//	IsClickAdmireIconRecently     int `json:"is_click_admire_icon_recently"`
		//	ExitAdmireInAwemePost         int `json:"exit_admire_in_aweme_post"`
		//	IsShowAdmireTab               int `json:"is_show_admire_tab"`
		//	IsShowAdmireButton            int `json:"is_show_admire_button"`
		//	AdmireButton                  int `json:"admire_button"`
		//	IsAdmire                      int `json:"is_admire"`
		//	IsFiftyAdmireAuthorStableFans int `json:"is_fifty_admire_author_stable_fans"`
		//	AuthorCanAdmire               int `json:"author_can_admire"`
		//	IsIronFansInAwemePost         int `json:"is_iron_fans_in_aweme_post"`
		//} `json:"admire_auth"`
		//CollectionCornerMark  int           `json:"collection_corner_mark"`
		//Images                interface{}   `json:"images"`
		//OriginDuetResourceUri string        `json:"origin_duet_resource_uri"`
		//ComponentInfoV2       string        `json:"component_info_v2"`
		//OriginTextExtra       []interface{} `json:"origin_text_extra"`
		//IsRelieve             bool          `json:"is_relieve"`
		//CommentWordsRecommend struct {
		//} `json:"comment_words_recommend"`
		//AwemeType                  int           `json:"aweme_type"`
		//ShareRecExtra              string        `json:"share_rec_extra"`
		//Promotions                 []interface{} `json:"promotions"`
		//VideoGameDataChannelConfig struct {
		//} `json:"video_game_data_channel_config"`
		//GeofencingRegions         interface{} `json:"geofencing_regions"`
		//RelationLabels            interface{} `json:"relation_labels"`
		//EcomCommentAtmosphereType int         `json:"ecom_comment_atmosphere_type"`
		//EntVideoDetail            struct {
		//	IsEntEncryptVideo bool `json:"is_ent_encrypt_video"`
		//} `json:"ent_video_detail"`
		//ImpressionData struct {
		//	SimilarIdListB interface{}   `json:"similar_id_list_b"`
		//	GroupIdListC   []interface{} `json:"group_id_list_c"`
		//	GroupIdListA   []interface{} `json:"group_id_list_a"`
		//	GroupIdListB   []interface{} `json:"group_id_list_b"`
		//	SimilarIdListA interface{}   `json:"similar_id_list_a"`
		//} `json:"impression_data"`
		//IsImageBeat                     bool          `json:"is_image_beat"`
		//IsInScope                       bool          `json:"is_in_scope"`
		//IsDuetSing                      bool          `json:"is_duet_sing"`
		//NearbyLevel                     int           `json:"nearby_level"`
		//LongVideo                       interface{}   `json:"long_video"`
		//Geofencing                      []interface{} `json:"geofencing"`
		//PersonalPageBottonDiagnoseStyle int           `json:"personal_page_botton_diagnose_style"`
		//ReportAction                    bool          `json:"report_action"`
		//Duration                        int           `json:"duration"`
		//XiguaTask                       struct {
		//	IsXiguaTask bool `json:"is_xigua_task"`
		//} `json:"xigua_task"`
		//NeedVsEntry    bool        `json:"need_vs_entry"`
		//ImageCropCtrl  int         `json:"image_crop_ctrl"`
		//OriginalImages interface{} `json:"original_images"`
		//UniqidPosition interface{} `json:"uniqid_position"`
		//IsTop          int         `json:"is_top"`
		//CreateTime     int         `json:"create_time"`
		//GuideSceneInfo struct {
		//} `json:"guide_scene_info"`
		//TrendsEventTrack string `json:"trends_event_track"`
		//BoostStatus      int    `json:"boost_status"`
		//CommerceInfo     struct {
		//	IsAd   bool `json:"is_ad"`
		//	AdType int  `json:"ad_type"`
		//} `json:"commerce_info"`
		//FlashMobTrends    int    `json:"flash_mob_trends"`
		//ActivityVideoType int    `json:"activity_video_type"`
		//AwemeId           string `json:"aweme_id"`
		Video struct {
			//	Width         int `json:"width"`
			//	CdnUrlExpired int `json:"cdn_url_expired"`
			//	IsH265        int `json:"is_h265"`
			//	IsSourceHDR   int `json:"is_source_HDR"`
			//	Audio         struct {
			//	} `json:"audio"`
			//	Meta     string      `json:"meta"`
			//	Tags     interface{} `json:"tags"`
			//	Height   int         `json:"height"`
			//	Duration int         `json:"duration"`
			BitRate []struct {
				GearName    string `json:"gear_name"`
				BitRate     int    `json:"bit_rate"`
				IsH265      int    `json:"is_h265"`
				IsBytevc1   int    `json:"is_bytevc1"`
				HDRType     string `json:"HDR_type"`
				FPS         int    `json:"FPS"`
				VideoExtra  string `json:"video_extra"`
				QualityType int    `json:"quality_type"`
				HDRBit      string `json:"HDR_bit"`
				PlayAddr    struct {
					Uri      string   `json:"uri"`
					Height   int      `json:"height"`
					Width    int      `json:"width"`
					DataSize int      `json:"data_size"`
					UrlList  []string `json:"url_list"`
					UrlKey   string   `json:"url_key"`
					FileHash string   `json:"file_hash"`
					FileCs   string   `json:"file_cs"`
				} `json:"play_addr"`
				Format string `json:"format"`
			} `json:"bit_rate"`
			//	CoverOriginalScale struct {
			//		Uri     string   `json:"uri"`
			//		Width   int      `json:"width"`
			//		UrlList []string `json:"url_list"`
			//		Height  int      `json:"height"`
			//	} `json:"cover_original_scale"`
			//	HasWatermark bool `json:"has_watermark"`
			//	Cover        struct {
			//		UrlList []string `json:"url_list"`
			//		Width   int      `json:"width"`
			//		Uri     string   `json:"uri"`
			//		Height  int      `json:"height"`
			//	} `json:"cover"`
			//	AnimatedCover struct {
			//		Uri     string   `json:"uri"`
			//		UrlList []string `json:"url_list"`
			//	} `json:"animated_cover"`
			PlayAddr265 struct {
				Uri      string   `json:"uri"`
				FileHash string   `json:"file_hash"`
				Height   int      `json:"height"`
				Width    int      `json:"width"`
				DataSize int      `json:"data_size"`
				FileCs   string   `json:"file_cs"`
				UrlList  []string `json:"url_list"`
				UrlKey   string   `json:"url_key"`
			} `json:"play_addr_265"`
			//	Format      string        `json:"format"`
			//	BigThumbs   []interface{} `json:"big_thumbs"`
			//	IsBytevc1   int           `json:"is_bytevc1"`
			//	OriginCover struct {
			//		UrlList []string `json:"url_list"`
			//		Uri     string   `json:"uri"`
			//		Width   int      `json:"width"`
			//		Height  int      `json:"height"`
			//	} `json:"origin_cover"`
			DownloadAddr struct {
				UrlList  []string `json:"url_list"`
				Height   int      `json:"height"`
				Width    int      `json:"width"`
				Uri      string   `json:"uri"`
				DataSize int      `json:"data_size"`
				FileCs   string   `json:"file_cs"`
			} `json:"download_addr"`
			//	Ratio    string `json:"ratio"`
			//	PlayAddr struct {
			//		UrlList  []string `json:"url_list"`
			//		FileHash string   `json:"file_hash"`
			//		DataSize int      `json:"data_size"`
			//		FileCs   string   `json:"file_cs"`
			//		Uri      string   `json:"uri"`
			//		UrlKey   string   `json:"url_key"`
			//		Height   int      `json:"height"`
			//		Width    int      `json:"width"`
			//	} `json:"play_addr"`
			//	IsCallback   bool `json:"is_callback"`
			//	NeedSetToken bool `json:"need_set_token"`
			//	DynamicCover struct {
			//		UrlList []string `json:"url_list"`
			//		Uri     string   `json:"uri"`
			//		Height  int      `json:"height"`
			//		Width   int      `json:"width"`
			//	} `json:"dynamic_cover"`
			PlayAddrH264 struct {
				Uri      string   `json:"uri"`
				Height   int      `json:"height"`
				Width    int      `json:"width"`
				DataSize int      `json:"data_size"`
				FileHash string   `json:"file_hash"`
				UrlList  []string `json:"url_list"`
				UrlKey   string   `json:"url_key"`
				FileCs   string   `json:"file_cs"`
			} `json:"play_addr_h264"`
			//	BitRateAudio interface{} `json:"bit_rate_audio"`
		} `json:"video"`
		//LabelTopText         interface{} `json:"label_top_text"`
		//DislikeDimensionList interface{} `json:"dislike_dimension_list"`
		//AuthorUserId         int64       `json:"author_user_id"`
		//CfAssetsType         int         `json:"cf_assets_type"`
		//IsFantasy            bool        `json:"is_fantasy"`
		//SocialTagList        interface{} `json:"social_tag_list"`
		//PhotoSearchEntrance  struct {
		//	EcomType int `json:"ecom_type"`
		//} `json:"photo_search_entrance"`
		//IsPreview              int    `json:"is_preview"`
		//Region                 string `json:"region"`
		//HasVsEntry             bool   `json:"has_vs_entry"`
		//PreventDownload        bool   `json:"prevent_download"`
		//EntertainmentVideoType int    `json:"entertainment_video_type"`
		//RiskInfos              struct {
		//	Vote     bool   `json:"vote"`
		//	Warn     bool   `json:"warn"`
		//	RiskSink bool   `json:"risk_sink"`
		//	Type     int    `json:"type"`
		//	Content  string `json:"content"`
		//} `json:"risk_infos"`
		//VideoLabels interface{} `json:"video_labels"`
		//ImgBitrate  interface{} `json:"img_bitrate"`
		//CommentGid  int64       `json:"comment_gid"`
		//ShareInfo   struct {
		//	ShareSignatureDesc string `json:"share_signature_desc"`
		//	ShareTitleMyself   string `json:"share_title_myself"`
		//	ShareTitleOther    string `json:"share_title_other"`
		//	ShareDescInfo      string `json:"share_desc_info"`
		//	ShareUrl           string `json:"share_url"`
		//	ShareQuote         string `json:"share_quote"`
		//	ShareSignatureUrl  string `json:"share_signature_url"`
		//	ShareWeiboDesc     string `json:"share_weibo_desc"`
		//	ShareTitle         string `json:"share_title"`
		//	ShareDesc          string `json:"share_desc"`
		//	BoolPersist        int    `json:"bool_persist"`
		//	ShareLinkDesc      string `json:"share_link_desc"`
		//} `json:"share_info"`
		//ChapterList   interface{} `json:"chapter_list"`
		//IsSharePost   bool        `json:"is_share_post"`
		//CfRecheckTs   int         `json:"cf_recheck_ts"`
		//GroupId       string      `json:"group_id"`
		//AuthorMaskTag int         `json:"author_mask_tag"`
		//CmtSwt        bool        `json:"cmt_swt"`
		//UserDigged    int         `json:"user_digged"`
		//IsAds         bool        `json:"is_ads"`
		//VideoControl  struct {
		//	TimerStatus           int    `json:"timer_status"`
		//	AllowDynamicWallpaper bool   `json:"allow_dynamic_wallpaper"`
		//	AllowShare            bool   `json:"allow_share"`
		//	ShareGrayed           bool   `json:"share_grayed"`
		//	ShareType             int    `json:"share_type"`
		//	PreventDownloadType   int    `json:"prevent_download_type"`
		//	DisableRecordReason   string `json:"disable_record_reason"`
		//	TimerInfo             struct {
		//	} `json:"timer_info"`
		//	DuetInfo struct {
		//		Level int `json:"level"`
		//	} `json:"duet_info"`
		//	AllowMusic               bool `json:"allow_music"`
		//	AllowDouplus             bool `json:"allow_douplus"`
		//	AllowStitch              bool `json:"allow_stitch"`
		//	ShareIgnoreVisibility    bool `json:"share_ignore_visibility"`
		//	ShowProgressBar          int  `json:"show_progress_bar"`
		//	AllowRecord              bool `json:"allow_record"`
		//	DraftProgressBar         int  `json:"draft_progress_bar"`
		//	AllowReact               bool `json:"allow_react"`
		//	AllowDuet                bool `json:"allow_duet"`
		//	DownloadIgnoreVisibility bool `json:"download_ignore_visibility"`
		//	DuetIgnoreVisibility     bool `json:"duet_ignore_visibility"`
		//	DownloadInfo             struct {
		//		Level int `json:"level"`
		//	} `json:"download_info"`
		//	AllowDownload bool `json:"allow_download"`
		//} `json:"video_control"`
		//Distance           string `json:"distance"`
		//IsStory            int    `json:"is_story"`
		//IsHashTag          int    `json:"is_hash_tag"`
		//DisableRelationBar int    `json:"disable_relation_bar"`
		//Status             struct {
		//	IsPrivate                  bool `json:"is_private"`
		//	AllowSelfRecommendToFriend bool `json:"allow_self_recommend_to_friend"`
		//	PartSee                    int  `json:"part_see"`
		//	Reviewed                   int  `json:"reviewed"`
		//	AllowFriendRecommendGuide  bool `json:"allow_friend_recommend_guide"`
		//	AllowShare                 bool `json:"allow_share"`
		//	AllowFriendRecommend       bool `json:"allow_friend_recommend"`
		//	IsProhibited               bool `json:"is_prohibited"`
		//	ListenVideoStatus          int  `json:"listen_video_status"`
		//	ReviewResult               struct {
		//		ReviewStatus int `json:"review_status"`
		//	} `json:"review_result"`
		//	DontShareStatus       int    `json:"dont_share_status"`
		//	NotAllowSoftDelReason string `json:"not_allow_soft_del_reason"`
		//	WithGoods             bool   `json:"with_goods"`
		//	SelfSee               bool   `json:"self_see"`
		//	IsDelete              bool   `json:"is_delete"`
		//	AllowComment          bool   `json:"allow_comment"`
		//	AwemeEditInfo         struct {
		//		ButtonToast    string `json:"button_toast"`
		//		EditStatus     int    `json:"edit_status"`
		//		HasModifiedAll bool   `json:"has_modified_all"`
		//		ButtonStatus   int    `json:"button_status"`
		//	} `json:"aweme_edit_info"`
		//	WithFusionGoods  bool   `json:"with_fusion_goods"`
		//	EnableSoftDelete int    `json:"enable_soft_delete"`
		//	PrivateStatus    int    `json:"private_status"`
		//	DownloadStatus   int    `json:"download_status"`
		//	InReviewing      bool   `json:"in_reviewing"`
		//	AwemeId          string `json:"aweme_id"`
		//	VideoHideSearch  int    `json:"video_hide_search"`
		//} `json:"status"`
		//EnableCommentStickerRec bool `json:"enable_comment_sticker_rec"`
		//ItemStitch              int  `json:"item_stitch"`
		//PublishPlusAlienation   struct {
		//	AlienationType int `json:"alienation_type"`
		//} `json:"publish_plus_alienation"`
		//VrType int `json:"vr_type"`
		//PoiBiz struct {
		//} `json:"poi_biz"`
		//InteractionStickers  interface{} `json:"interaction_stickers"`
		//ShouldOpenAdReport   bool        `json:"should_open_ad_report"`
		//MarkLargelyFollowing bool        `json:"mark_largely_following"`
		//TextExtra            []struct {
		//	Start        int    `json:"start"`
		//	End          int    `json:"end"`
		//	Type         int    `json:"type"`
		//	CaptionEnd   int    `json:"caption_end"`
		//	CaptionStart int    `json:"caption_start"`
		//	HashtagName  string `json:"hashtag_name"`
		//	HashtagId    string `json:"hashtag_id"`
		//	IsCommerce   bool   `json:"is_commerce"`
		//} `json:"text_extra"`
		//ItemCommentSettings   int           `json:"item_comment_settings"`
		PreviewTitle string `json:"preview_title"`
		//ItemShare             int           `json:"item_share"`
		//CollectStat           int           `json:"collect_stat"`
		//VideoText             []interface{} `json:"video_text"`
		//CommentPermissionInfo struct {
		//	CanComment              bool `json:"can_comment"`
		//	ItemDetailEntry         bool `json:"item_detail_entry"`
		//	PressEntry              bool `json:"press_entry"`
		//	ToastGuide              bool `json:"toast_guide"`
		//	CommentPermissionStatus int  `json:"comment_permission_status"`
		//} `json:"comment_permission_info"`
		//IncentiveItemType int `json:"incentive_item_type"`
	} `json:"aweme_detail"`
}
