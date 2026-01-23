package tikhub

import (
	"context"
	"encoding/json"
)

func (t Client) TiktokGetVideoByShareUrl(ctx context.Context, url string) (*TkVideo, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("share_url", url).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/tiktok/app/v3/fetch_one_video_by_share_url")

	if err != nil {
		return nil, err
	}

	var resp Response[TkVideo]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

type TkVideo struct {
	AwemeDetails []struct {
		AddedSoundMusicInfo struct {
			Album                         string        `json:"album"`
			AllowOfflineMusicToDetailPage bool          `json:"allow_offline_music_to_detail_page"`
			Artists                       []interface{} `json:"artists"`
			AuditionDuration              int           `json:"audition_duration"`
			Author                        string        `json:"author"`
			AuthorDeleted                 bool          `json:"author_deleted"`
			AuthorPosition                interface{}   `json:"author_position"`
			AvatarMedium                  struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"avatar_medium"`
			AvatarThumb struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"avatar_thumb"`
			BindedChallengeId   int  `json:"binded_challenge_id"`
			CanBeStitched       bool `json:"can_be_stitched"`
			CanNotReuse         bool `json:"can_not_reuse"`
			CollectStat         int  `json:"collect_stat"`
			CommercialRightType int  `json:"commercial_right_type"`
			CoverLarge          struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"cover_large"`
			CoverMedium struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"cover_medium"`
			CoverThumb struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"cover_thumb"`
			CreateTime            int  `json:"create_time"`
			DmvAutoShow           bool `json:"dmv_auto_show"`
			Duration              int  `json:"duration"`
			DurationHighPrecision struct {
				AuditionDurationPrecision float64 `json:"audition_duration_precision"`
				DurationPrecision         float64 `json:"duration_precision"`
				ShootDurationPrecision    float64 `json:"shoot_duration_precision"`
				VideoDurationPrecision    float64 `json:"video_duration_precision"`
			} `json:"duration_high_precision"`
			ExternalSongInfo     []interface{} `json:"external_song_info"`
			Extra                string        `json:"extra"`
			HasCommerceRight     bool          `json:"has_commerce_right"`
			Id                   int64         `json:"id"`
			IdStr                string        `json:"id_str"`
			IsAudioUrlWithCookie bool          `json:"is_audio_url_with_cookie"`
			IsAuthorArtist       bool          `json:"is_author_artist"`
			IsCommerceMusic      bool          `json:"is_commerce_music"`
			IsMatchedMetadata    bool          `json:"is_matched_metadata"`
			IsOriginal           bool          `json:"is_original"`
			IsOriginalSound      bool          `json:"is_original_sound"`
			IsPgc                bool          `json:"is_pgc"`
			IsPlayMusic          bool          `json:"is_play_music"`
			IsShootingAllow      bool          `json:"is_shooting_allow"`
			LogExtra             string        `json:"log_extra"`
			LyricShortPosition   interface{}   `json:"lyric_short_position"`
			MemeSongInfo         struct {
			} `json:"meme_song_info"`
			Mid                  string      `json:"mid"`
			MultiBitRatePlayInfo interface{} `json:"multi_bit_rate_play_info"`
			MuteShare            bool        `json:"mute_share"`
			OfflineDesc          string      `json:"offline_desc"`
			OwnerHandle          string      `json:"owner_handle"`
			OwnerId              string      `json:"owner_id"`
			OwnerNickname        string      `json:"owner_nickname"`
			PlayUrl              struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"play_url"`
			Position         interface{} `json:"position"`
			PreventDownload  bool        `json:"prevent_download"`
			PreviewEndTime   int         `json:"preview_end_time"`
			PreviewStartTime int         `json:"preview_start_time"`
			RecommendStatus  int         `json:"recommend_status"`
			SearchHighlight  interface{} `json:"search_highlight"`
			SecUid           string      `json:"sec_uid"`
			ShootDuration    int         `json:"shoot_duration"`
			SourcePlatform   int         `json:"source_platform"`
			Status           int         `json:"status"`
			StrongBeatUrl    struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"strong_beat_url"`
			TagList          interface{} `json:"tag_list"`
			Title            string      `json:"title"`
			TtToDspSongInfos interface{} `json:"tt_to_dsp_song_infos"`
			UncertArtists    interface{} `json:"uncert_artists"`
			UserCount        int         `json:"user_count"`
			VideoDuration    int         `json:"video_duration"`
		} `json:"added_sound_music_info"`
		AigcInfo struct {
			AigcLabelType int  `json:"aigc_label_type"`
			CreatedByAi   bool `json:"created_by_ai"`
		} `json:"aigc_info"`
		Anchors           interface{} `json:"anchors"`
		AnchorsExtras     string      `json:"anchors_extras"`
		AnimatedImageInfo struct {
			Effect int `json:"effect"`
			Type   int `json:"type"`
		} `json:"animated_image_info"`
		Author struct {
			AcceptPrivatePolicy     bool        `json:"accept_private_policy"`
			AccountLabels           interface{} `json:"account_labels"`
			AccountRegion           string      `json:"account_region"`
			AdCoverUrl              interface{} `json:"ad_cover_url"`
			AdvanceFeatureItemOrder interface{} `json:"advance_feature_item_order"`
			AdvancedFeatureInfo     interface{} `json:"advanced_feature_info"`
			AppleAccount            int         `json:"apple_account"`
			AuthorityStatus         int         `json:"authority_status"`
			Avatar168X168           struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"avatar_168x168"`
			Avatar300X300 struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"avatar_300x300"`
			AvatarLarger struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"avatar_larger"`
			AvatarMedium struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"avatar_medium"`
			AvatarThumb struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"avatar_thumb"`
			AvatarUri                  string        `json:"avatar_uri"`
			AwemeCount                 int           `json:"aweme_count"`
			BindPhone                  string        `json:"bind_phone"`
			BoldFields                 interface{}   `json:"bold_fields"`
			CanMessageFollowStatusList interface{}   `json:"can_message_follow_status_list"`
			CanSetGeofencing           interface{}   `json:"can_set_geofencing"`
			ChaList                    interface{}   `json:"cha_list"`
			CommentFilterStatus        int           `json:"comment_filter_status"`
			CommentSetting             int           `json:"comment_setting"`
			CommerceUserLevel          int           `json:"commerce_user_level"`
			CoverUrl                   []interface{} `json:"cover_url"`
			CreateTime                 int           `json:"create_time"`
			CustomVerify               string        `json:"custom_verify"`
			CvLevel                    string        `json:"cv_level"`
			DownloadPromptTs           int           `json:"download_prompt_ts"`
			DownloadSetting            int           `json:"download_setting"`
			DuetSetting                int           `json:"duet_setting"`
			EnabledFilterAllComments   bool          `json:"enabled_filter_all_comments"`
			EnterpriseVerifyReason     string        `json:"enterprise_verify_reason"`
			Events                     interface{}   `json:"events"`
			FakeDataInfo               struct {
			} `json:"fake_data_info"`
			FavoritingCount       int         `json:"favoriting_count"`
			FbExpireTime          int         `json:"fb_expire_time"`
			FollowStatus          int         `json:"follow_status"`
			FollowerCount         int         `json:"follower_count"`
			FollowerStatus        int         `json:"follower_status"`
			FollowersDetail       interface{} `json:"followers_detail"`
			FollowingCount        int         `json:"following_count"`
			FriendsStatus         int         `json:"friends_status"`
			Geofencing            interface{} `json:"geofencing"`
			GoogleAccount         string      `json:"google_account"`
			HasEmail              bool        `json:"has_email"`
			HasFacebookToken      bool        `json:"has_facebook_token"`
			HasInsights           bool        `json:"has_insights"`
			HasOrders             bool        `json:"has_orders"`
			HasTwitterToken       bool        `json:"has_twitter_token"`
			HasYoutubeToken       bool        `json:"has_youtube_token"`
			HideSearch            bool        `json:"hide_search"`
			HomepageBottomToast   interface{} `json:"homepage_bottom_toast"`
			InsId                 string      `json:"ins_id"`
			IsAdFake              bool        `json:"is_ad_fake"`
			IsBlock               bool        `json:"is_block"`
			IsDisciplineMember    bool        `json:"is_discipline_member"`
			IsMute                int         `json:"is_mute"`
			IsMuteLives           int         `json:"is_mute_lives"`
			IsMuteNonStoryPost    int         `json:"is_mute_non_story_post"`
			IsMuteStory           int         `json:"is_mute_story"`
			IsPhoneBinded         bool        `json:"is_phone_binded"`
			IsStar                bool        `json:"is_star"`
			ItemList              interface{} `json:"item_list"`
			Language              string      `json:"language"`
			LiveAgreement         int         `json:"live_agreement"`
			LiveCommerce          bool        `json:"live_commerce"`
			LiveVerify            int         `json:"live_verify"`
			MentionStatus         int         `json:"mention_status"`
			MutualRelationAvatars interface{} `json:"mutual_relation_avatars"`
			NeedPoints            interface{} `json:"need_points"`
			NeedRecommend         int         `json:"need_recommend"`
			Nickname              string      `json:"nickname"`
			PlatformSyncInfo      interface{} `json:"platform_sync_info"`
			PreventDownload       bool        `json:"prevent_download"`
			ReactSetting          int         `json:"react_setting"`
			Region                string      `json:"region"`
			RelativeUsers         interface{} `json:"relative_users"`
			ReplyWithVideoFlag    int         `json:"reply_with_video_flag"`
			RoomId                int         `json:"room_id"`
			SearchHighlight       interface{} `json:"search_highlight"`
			SecUid                string      `json:"sec_uid"`
			Secret                int         `json:"secret"`
			ShareInfo             struct {
				NowInvitationCardImageUrls interface{} `json:"now_invitation_card_image_urls"`
				ShareDesc                  string      `json:"share_desc"`
				ShareDescInfo              string      `json:"share_desc_info"`
				ShareQrcodeUrl             struct {
					Height    int           `json:"height"`
					Uri       string        `json:"uri"`
					UrlList   []interface{} `json:"url_list"`
					UrlPrefix interface{}   `json:"url_prefix"`
					Width     int           `json:"width"`
				} `json:"share_qrcode_url"`
				ShareTitle       string `json:"share_title"`
				ShareTitleMyself string `json:"share_title_myself"`
				ShareTitleOther  string `json:"share_title_other"`
				ShareUrl         string `json:"share_url"`
			} `json:"share_info"`
			ShareQrcodeUri      string      `json:"share_qrcode_uri"`
			ShieldCommentNotice int         `json:"shield_comment_notice"`
			ShieldDiggNotice    int         `json:"shield_digg_notice"`
			ShieldEditFieldInfo interface{} `json:"shield_edit_field_info"`
			ShieldFollowNotice  int         `json:"shield_follow_notice"`
			ShortId             string      `json:"short_id"`
			ShowImageBubble     bool        `json:"show_image_bubble"`
			Signature           string      `json:"signature"`
			SpecialAccount      struct {
				SpecialAccountList interface{} `json:"special_account_list"`
			} `json:"special_account"`
			SpecialLock        int         `json:"special_lock"`
			Status             int         `json:"status"`
			StitchSetting      int         `json:"stitch_setting"`
			TotalFavorited     int         `json:"total_favorited"`
			TwExpireTime       int         `json:"tw_expire_time"`
			TwitterId          string      `json:"twitter_id"`
			TwitterName        string      `json:"twitter_name"`
			TypeLabel          interface{} `json:"type_label"`
			Uid                string      `json:"uid"`
			UniqueId           string      `json:"unique_id"`
			UniqueIdModifyTime int         `json:"unique_id_modify_time"`
			UserCanceled       bool        `json:"user_canceled"`
			UserMode           int         `json:"user_mode"`
			UserPeriod         int         `json:"user_period"`
			UserProfileGuide   interface{} `json:"user_profile_guide"`
			UserRate           int         `json:"user_rate"`
			UserSparkInfo      struct {
			} `json:"user_spark_info"`
			UserTags         interface{} `json:"user_tags"`
			VerificationType int         `json:"verification_type"`
			VerifyInfo       string      `json:"verify_info"`
			VideoIcon        struct {
				Height    int           `json:"height"`
				Uri       string        `json:"uri"`
				UrlList   []interface{} `json:"url_list"`
				UrlPrefix interface{}   `json:"url_prefix"`
				Width     int           `json:"width"`
			} `json:"video_icon"`
			WhiteCoverUrl       interface{} `json:"white_cover_url"`
			WithCommerceEntry   bool        `json:"with_commerce_entry"`
			WithShopEntry       bool        `json:"with_shop_entry"`
			YoutubeChannelId    string      `json:"youtube_channel_id"`
			YoutubeChannelTitle string      `json:"youtube_channel_title"`
			YoutubeExpireTime   int         `json:"youtube_expire_time"`
		} `json:"author"`
		AuthorUserId int64 `json:"author_user_id"`
		AwemeAcl     struct {
			DownloadGeneral struct {
				Code      int  `json:"code"`
				Mute      bool `json:"mute"`
				ShowType  int  `json:"show_type"`
				Transcode int  `json:"transcode"`
			} `json:"download_general"`
			DownloadMaskPanel struct {
				Code      int  `json:"code"`
				Mute      bool `json:"mute"`
				ShowType  int  `json:"show_type"`
				Transcode int  `json:"transcode"`
			} `json:"download_mask_panel"`
			PlatformList    interface{} `json:"platform_list"`
			PressActionList interface{} `json:"press_action_list"`
			ShareActionList interface{} `json:"share_action_list"`
			ShareGeneral    struct {
				Code      int  `json:"code"`
				Mute      bool `json:"mute"`
				ShowType  int  `json:"show_type"`
				Transcode int  `json:"transcode"`
			} `json:"share_general"`
			ShareListStatus int `json:"share_list_status"`
		} `json:"aweme_acl"`
		AwemeId                    string      `json:"aweme_id"`
		AwemeType                  int         `json:"aweme_type"`
		Banners                    interface{} `json:"banners"`
		BehindTheSongMusicIds      interface{} `json:"behind_the_song_music_ids"`
		BehindTheSongVideoMusicIds interface{} `json:"behind_the_song_video_music_ids"`
		BodydanceScore             int         `json:"bodydance_score"`
		BrandedContentAccounts     interface{} `json:"branded_content_accounts"`
		CcTemplateInfo             struct {
			AuthorName           string `json:"author_name"`
			ClipCount            int    `json:"clip_count"`
			Desc                 string `json:"desc"`
			DurationMilliseconds int    `json:"duration_milliseconds"`
			RelatedMusicId       string `json:"related_music_id"`
			TemplateId           string `json:"template_id"`
		} `json:"cc_template_info"`
		ChaList           interface{} `json:"cha_list"`
		ChallengePosition interface{} `json:"challenge_position"`
		CmtSwt            bool        `json:"cmt_swt"`
		CollectStat       int         `json:"collect_stat"`
		CommentConfig     struct {
			CommentPanelShowTabConfig struct {
				CommentTabInfoConfig []struct {
					Priority int    `json:"priority"`
					TabId    int    `json:"tab_id"`
					TabName  string `json:"tab_name"`
				} `json:"comment_tab_info_config"`
				MaxTabCount int `json:"max_tab_count"`
			} `json:"comment_panel_show_tab_config"`
			EmojiRecommendList     interface{} `json:"emoji_recommend_list"`
			LongPressRecommendList interface{} `json:"long_press_recommend_list"`
			Preload                struct {
				Preds string `json:"preds"`
			} `json:"preload"`
			QuickComment struct {
				Enabled bool `json:"enabled"`
			} `json:"quick_comment"`
			QuickCommentEmojiRecommendList interface{} `json:"quick_comment_emoji_recommend_list"`
		} `json:"comment_config"`
		CommentTopbarInfo  interface{} `json:"comment_topbar_info"`
		CommerceConfigData interface{} `json:"commerce_config_data"`
		CommerceInfo       struct {
			AdvPromotable          bool   `json:"adv_promotable"`
			AuctionAdInvited       bool   `json:"auction_ad_invited"`
			BrandedContentType     int    `json:"branded_content_type"`
			OrganicLogExtra        string `json:"organic_log_extra"`
			WithCommentFilterWords bool   `json:"with_comment_filter_words"`
		} `json:"commerce_info"`
		ContentDesc      string        `json:"content_desc"`
		ContentDescExtra []interface{} `json:"content_desc_extra"`
		ContentLevel     int           `json:"content_level"`
		ContentModel     struct {
			CustomBiz struct {
				AwemeTrace string `json:"aweme_trace"`
			} `json:"custom_biz"`
			StandardBiz struct {
				CreatorAnalytics struct {
					CreatorAnalyticsEntranceStatus int `json:"creator_analytics_entrance_status"`
				} `json:"creator_analytics"`
				ECommerce struct {
					TtecContentTag struct {
						RecommendationTagConsumerStr string `json:"recommendation_tag_consumer_str"`
						RecommendationTagCreatorStr  string `json:"recommendation_tag_creator_str"`
					} `json:"ttec_content_tag"`
				} `json:"e_commerce"`
				LocalAllianceInfo struct {
					AllianceItemLabelText string `json:"alliance_item_label_text"`
					AllianceItemLabelType int    `json:"alliance_item_label_type"`
				} `json:"local_alliance_info"`
				TtsVoiceInfo struct {
					TtsVoiceAttr        string `json:"tts_voice_attr"`
					TtsVoiceReuseParams string `json:"tts_voice_reuse_params"`
				} `json:"tts_voice_info"`
				VcFilterInfo struct {
					VcFilterAttr string `json:"vc_filter_attr"`
				} `json:"vc_filter_info"`
			} `json:"standard_biz"`
		} `json:"content_model"`
		ContentOriginalType int         `json:"content_original_type"`
		ContentSizeType     int         `json:"content_size_type"`
		ContentType         string      `json:"content_type"`
		CoverLabels         interface{} `json:"cover_labels"`
		CreateTime          int         `json:"create_time"`
		CreationInfo        struct {
			CreationUsedFunctions []string `json:"creation_used_functions"`
		} `json:"creation_info"`
		Desc                     string      `json:"desc"`
		DescLanguage             string      `json:"desc_language"`
		DisableSearchTrendingBar bool        `json:"disable_search_trending_bar"`
		Distance                 string      `json:"distance"`
		DistributeType           int         `json:"distribute_type"`
		FollowUpPublishFromId    int         `json:"follow_up_publish_from_id"`
		Geofencing               interface{} `json:"geofencing"`
		GeofencingRegions        interface{} `json:"geofencing_regions"`
		GreenScreenMaterials     interface{} `json:"green_screen_materials"`
		GroupId                  string      `json:"group_id"`
		GroupIdList              struct {
			GroupdIdList0 interface{} `json:"GroupdIdList0"`
			GroupdIdList1 []int64     `json:"GroupdIdList1"`
		} `json:"group_id_list"`
		HasDanmaku         bool        `json:"has_danmaku"`
		HasVsEntry         bool        `json:"has_vs_entry"`
		HaveDashboard      bool        `json:"have_dashboard"`
		HybridLabel        interface{} `json:"hybrid_label"`
		ImageInfos         interface{} `json:"image_infos"`
		InteractPermission struct {
			AllowAddingAsPost struct {
				Status int `json:"status"`
			} `json:"allow_adding_as_post"`
			AllowAddingToStory int `json:"allow_adding_to_story"`
			AllowCreateSticker struct {
				Status int `json:"status"`
			} `json:"allow_create_sticker"`
			AllowStorySwitchToPost struct {
			} `json:"allow_story_switch_to_post"`
			Duet                 int `json:"duet"`
			DuetPrivacySetting   int `json:"duet_privacy_setting"`
			Stitch               int `json:"stitch"`
			StitchPrivacySetting int `json:"stitch_privacy_setting"`
			Upvote               int `json:"upvote"`
		} `json:"interact_permission"`
		InteractionStickers       interface{} `json:"interaction_stickers"`
		IsAds                     bool        `json:"is_ads"`
		IsDescriptionTranslatable bool        `json:"is_description_translatable"`
		IsHashTag                 int         `json:"is_hash_tag"`
		IsNffOrNr                 bool        `json:"is_nff_or_nr"`
		IsOnThisDay               int         `json:"is_on_this_day"`
		IsPgcshow                 bool        `json:"is_pgcshow"`
		IsPreview                 int         `json:"is_preview"`
		IsRelieve                 bool        `json:"is_relieve"`
		IsTextStickerTranslatable bool        `json:"is_text_sticker_translatable"`
		IsTitleTranslatable       bool        `json:"is_title_translatable"`
		IsTop                     int         `json:"is_top"`
		IsVr                      bool        `json:"is_vr"`
		ItemCommentSettings       int         `json:"item_comment_settings"`
		ItemDuet                  int         `json:"item_duet"`
		ItemReact                 int         `json:"item_react"`
		ItemStitch                int         `json:"item_stitch"`
		LabelTop                  struct {
			Height    int         `json:"height"`
			Uri       string      `json:"uri"`
			UrlList   []string    `json:"url_list"`
			UrlPrefix interface{} `json:"url_prefix"`
			Width     int         `json:"width"`
		} `json:"label_top"`
		LabelTopText   interface{}   `json:"label_top_text"`
		LongVideo      interface{}   `json:"long_video"`
		MainArchCommon string        `json:"main_arch_common"`
		MaskInfos      []interface{} `json:"mask_infos"`
		MemeRegInfo    struct {
		} `json:"meme_reg_info"`
		MiscInfo         string      `json:"misc_info"`
		MufCommentInfoV2 interface{} `json:"muf_comment_info_v2"`
		Music            struct {
			Album                         string        `json:"album"`
			AllowOfflineMusicToDetailPage bool          `json:"allow_offline_music_to_detail_page"`
			Artists                       []interface{} `json:"artists"`
			AuditionDuration              int           `json:"audition_duration"`
			Author                        string        `json:"author"`
			AuthorDeleted                 bool          `json:"author_deleted"`
			AuthorPosition                interface{}   `json:"author_position"`
			AvatarMedium                  struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"avatar_medium"`
			AvatarThumb struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"avatar_thumb"`
			BindedChallengeId   int  `json:"binded_challenge_id"`
			CanBeStitched       bool `json:"can_be_stitched"`
			CanNotReuse         bool `json:"can_not_reuse"`
			CollectStat         int  `json:"collect_stat"`
			CommercialRightType int  `json:"commercial_right_type"`
			CoverLarge          struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"cover_large"`
			CoverMedium struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"cover_medium"`
			CoverThumb struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"cover_thumb"`
			CreateTime            int  `json:"create_time"`
			DmvAutoShow           bool `json:"dmv_auto_show"`
			Duration              int  `json:"duration"`
			DurationHighPrecision struct {
				AuditionDurationPrecision float64 `json:"audition_duration_precision"`
				DurationPrecision         float64 `json:"duration_precision"`
				ShootDurationPrecision    float64 `json:"shoot_duration_precision"`
				VideoDurationPrecision    float64 `json:"video_duration_precision"`
			} `json:"duration_high_precision"`
			ExternalSongInfo     []interface{} `json:"external_song_info"`
			Extra                string        `json:"extra"`
			HasCommerceRight     bool          `json:"has_commerce_right"`
			Id                   int64         `json:"id"`
			IdStr                string        `json:"id_str"`
			IsAudioUrlWithCookie bool          `json:"is_audio_url_with_cookie"`
			IsAuthorArtist       bool          `json:"is_author_artist"`
			IsCommerceMusic      bool          `json:"is_commerce_music"`
			IsMatchedMetadata    bool          `json:"is_matched_metadata"`
			IsOriginal           bool          `json:"is_original"`
			IsOriginalSound      bool          `json:"is_original_sound"`
			IsPgc                bool          `json:"is_pgc"`
			IsPlayMusic          bool          `json:"is_play_music"`
			IsShootingAllow      bool          `json:"is_shooting_allow"`
			LogExtra             string        `json:"log_extra"`
			LyricShortPosition   interface{}   `json:"lyric_short_position"`
			MemeSongInfo         struct {
			} `json:"meme_song_info"`
			Mid                  string      `json:"mid"`
			MultiBitRatePlayInfo interface{} `json:"multi_bit_rate_play_info"`
			MuteShare            bool        `json:"mute_share"`
			OfflineDesc          string      `json:"offline_desc"`
			OwnerHandle          string      `json:"owner_handle"`
			OwnerId              string      `json:"owner_id"`
			OwnerNickname        string      `json:"owner_nickname"`
			PlayUrl              struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"play_url"`
			Position         interface{} `json:"position"`
			PreventDownload  bool        `json:"prevent_download"`
			PreviewEndTime   int         `json:"preview_end_time"`
			PreviewStartTime int         `json:"preview_start_time"`
			RecommendStatus  int         `json:"recommend_status"`
			SearchHighlight  interface{} `json:"search_highlight"`
			SecUid           string      `json:"sec_uid"`
			ShootDuration    int         `json:"shoot_duration"`
			SourcePlatform   int         `json:"source_platform"`
			Status           int         `json:"status"`
			StrongBeatUrl    struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"strong_beat_url"`
			TagList          interface{} `json:"tag_list"`
			Title            string      `json:"title"`
			TtToDspSongInfos interface{} `json:"tt_to_dsp_song_infos"`
			UncertArtists    interface{} `json:"uncert_artists"`
			UserCount        int         `json:"user_count"`
			VideoDuration    int         `json:"video_duration"`
		} `json:"music"`
		MusicBeginTimeInMs int         `json:"music_begin_time_in_ms"`
		MusicEndTimeInMs   int         `json:"music_end_time_in_ms"`
		MusicSelectedFrom  string      `json:"music_selected_from"`
		MusicTitleStyle    int         `json:"music_title_style"`
		MusicVolume        string      `json:"music_volume"`
		NeedTrimStep       bool        `json:"need_trim_step"`
		NeedVsEntry        bool        `json:"need_vs_entry"`
		NicknamePosition   interface{} `json:"nickname_position"`
		NoSelectedMusic    bool        `json:"no_selected_music"`
		OperatorBoostInfo  interface{} `json:"operator_boost_info"`
		OriginCommentIds   interface{} `json:"origin_comment_ids"`
		OriginVolume       string      `json:"origin_volume"`
		OriginalClientText struct {
			MarkupText string      `json:"markup_text"`
			TextExtra  interface{} `json:"text_extra"`
		} `json:"original_client_text"`
		PickedUsers             []interface{} `json:"picked_users"`
		PlaylistBlocked         bool          `json:"playlist_blocked"`
		PoiReTagSignal          int           `json:"poi_re_tag_signal"`
		Position                interface{}   `json:"position"`
		PreventDownload         bool          `json:"prevent_download"`
		ProductsInfo            interface{}   `json:"products_info"`
		QuestionList            interface{}   `json:"question_list"`
		QuickReplyEmojis        []string      `json:"quick_reply_emojis"`
		Rate                    int           `json:"rate"`
		ReferenceTtsVoiceIds    interface{}   `json:"reference_tts_voice_ids"`
		ReferenceVoiceFilterIds interface{}   `json:"reference_voice_filter_ids"`
		Region                  string        `json:"region"`
		RetryType               int           `json:"retry_type"`
		RiskInfos               struct {
			Content  string `json:"content"`
			RiskSink bool   `json:"risk_sink"`
			Type     int    `json:"type"`
			Vote     bool   `json:"vote"`
			Warn     bool   `json:"warn"`
		} `json:"risk_infos"`
		SearchHighlight interface{} `json:"search_highlight"`
		ShareInfo       struct {
			BoolPersist                int         `json:"bool_persist"`
			NowInvitationCardImageUrls interface{} `json:"now_invitation_card_image_urls"`
			ShareDesc                  string      `json:"share_desc"`
			ShareDescInfo              string      `json:"share_desc_info"`
			ShareLinkDesc              string      `json:"share_link_desc"`
			ShareQuote                 string      `json:"share_quote"`
			ShareSignatureDesc         string      `json:"share_signature_desc"`
			ShareSignatureUrl          string      `json:"share_signature_url"`
			ShareTitle                 string      `json:"share_title"`
			ShareTitleMyself           string      `json:"share_title_myself"`
			ShareTitleOther            string      `json:"share_title_other"`
			ShareUrl                   string      `json:"share_url"`
			WhatsappDesc               string      `json:"whatsapp_desc"`
		} `json:"share_info"`
		ShareUrl              string `json:"share_url"`
		ShootTabName          string `json:"shoot_tab_name"`
		SocialInteractionBlob struct {
			AuxiliaryModelContent string `json:"auxiliary_model_content"`
		} `json:"social_interaction_blob"`
		SolariaProfile struct {
		} `json:"solaria_profile"`
		SortLabel  string `json:"sort_label"`
		Statistics struct {
			AwemeId            string `json:"aweme_id"`
			CollectCount       int    `json:"collect_count"`
			CommentCount       int    `json:"comment_count"`
			DiggCount          int    `json:"digg_count"`
			DownloadCount      int    `json:"download_count"`
			ForwardCount       int    `json:"forward_count"`
			LoseCommentCount   int    `json:"lose_comment_count"`
			LoseCount          int    `json:"lose_count"`
			PlayCount          int    `json:"play_count"`
			RepostCount        int    `json:"repost_count"`
			ShareCount         int    `json:"share_count"`
			WhatsappShareCount int    `json:"whatsapp_share_count"`
		} `json:"statistics"`
		Status struct {
			AllowComment   bool   `json:"allow_comment"`
			AllowShare     bool   `json:"allow_share"`
			AwemeId        string `json:"aweme_id"`
			DownloadStatus int    `json:"download_status"`
			InReviewing    bool   `json:"in_reviewing"`
			IsDelete       bool   `json:"is_delete"`
			IsProhibited   bool   `json:"is_prohibited"`
			PrivateStatus  int    `json:"private_status"`
			ReviewResult   struct {
				ReviewStatus int `json:"review_status"`
			} `json:"review_result"`
			Reviewed int  `json:"reviewed"`
			SelfSee  bool `json:"self_see"`
		} `json:"status"`
		SupportDanmaku       bool          `json:"support_danmaku"`
		TextExtra            []interface{} `json:"text_extra"`
		TextStickerMajorLang string        `json:"text_sticker_major_lang"`
		TitleLanguage        string        `json:"title_language"`
		TtecSuggestWords     struct {
			TtecSuggestWords interface{} `json:"ttec_suggest_words"`
		} `json:"ttec_suggest_words"`
		TtsVoiceIds          interface{} `json:"tts_voice_ids"`
		TttProductRecallType int         `json:"ttt_product_recall_type"`
		UniqidPosition       interface{} `json:"uniqid_position"`
		UsedFullSong         bool        `json:"used_full_song"`
		UserDigged           int         `json:"user_digged"`
		Video                struct {
			CoverTsp       float64 `json:"CoverTsp"`
			AiDynamicCover struct {
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
			} `json:"ai_dynamic_cover"`
			AiDynamicCoverBak struct {
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
			} `json:"ai_dynamic_cover_bak"`
			AnimatedCover struct {
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
			} `json:"animated_cover"`
			BigThumbs []interface{} `json:"big_thumbs"`
			BitRate   []struct {
				HDRBit           string      `json:"HDR_bit"`
				HDRType          string      `json:"HDR_type"`
				BitRate          int         `json:"bit_rate"`
				DubInfos         interface{} `json:"dub_infos"`
				FidProfileLabels string      `json:"fid_profile_labels"`
				Fps              int         `json:"fps"`
				GearName         string      `json:"gear_name"`
				IsBytevc1        int         `json:"is_bytevc1"`
				PlayAddr         struct {
					DataSize  int         `json:"data_size"`
					FileCs    string      `json:"file_cs"`
					FileHash  string      `json:"file_hash"`
					Height    int         `json:"height"`
					Uri       string      `json:"uri"`
					UrlKey    string      `json:"url_key"`
					UrlList   []string    `json:"url_list"`
					UrlPrefix interface{} `json:"url_prefix"`
					Width     int         `json:"width"`
				} `json:"play_addr"`
				QualityType int    `json:"quality_type"`
				VideoExtra  string `json:"video_extra"`
			} `json:"bit_rate"`
			BitRateAudio  []interface{} `json:"bit_rate_audio"`
			CdnUrlExpired int           `json:"cdn_url_expired"`
			ClaInfo       struct {
				CaptionInfos []struct {
					CaptionFormat     string   `json:"caption_format"`
					CaptionLength     int      `json:"caption_length"`
					ClaSubtitleId     int64    `json:"cla_subtitle_id"`
					ComplaintId       int64    `json:"complaint_id"`
					Expire            int      `json:"expire"`
					IsAutoGenerated   bool     `json:"is_auto_generated"`
					IsOriginalCaption bool     `json:"is_original_caption"`
					Lang              string   `json:"lang"`
					LanguageCode      string   `json:"language_code"`
					LanguageId        int      `json:"language_id"`
					SourceTag         string   `json:"source_tag"`
					SubId             int      `json:"sub_id"`
					SubVersion        string   `json:"sub_version"`
					SubtitleType      int      `json:"subtitle_type"`
					TranslationType   int      `json:"translation_type"`
					TranslatorId      int      `json:"translator_id"`
					Url               string   `json:"url"`
					UrlList           []string `json:"url_list"`
					Variant           string   `json:"variant"`
				} `json:"caption_infos"`
				CaptionsType           int  `json:"captions_type"`
				CreatorEditedCaptionId int  `json:"creator_edited_caption_id"`
				EnableAutoCaption      int  `json:"enable_auto_caption"`
				HasOriginalAudio       int  `json:"has_original_audio"`
				HideOriginalCaption    bool `json:"hide_original_caption"`
				NoCaptionReason        int  `json:"no_caption_reason"`
				OriginalLanguageInfo   struct {
					CanTranslateRealtime                         bool   `json:"can_translate_realtime"`
					CanTranslateRealtimeSkipTranslationLangCheck bool   `json:"can_translate_realtime_skip_translation_lang_check"`
					FirstSubtitleTime                            int    `json:"first_subtitle_time"`
					IsBurninCaption                              bool   `json:"is_burnin_caption"`
					Lang                                         string `json:"lang"`
					LanguageCode                                 string `json:"language_code"`
					LanguageId                                   int    `json:"language_id"`
					OriginalCaptionType                          int    `json:"original_caption_type"`
				} `json:"original_language_info"`
				VerticalPositions interface{} `json:"vertical_positions"`
			} `json:"cla_info"`
			Cover struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"cover"`
			CoverIsCustom    bool   `json:"cover_is_custom"`
			DidProfileLabels string `json:"did_profile_labels"`
			DownloadAddr     struct {
				DataSize  int         `json:"data_size"`
				FileCs    string      `json:"file_cs"`
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"download_addr"`
			DownloadNoWatermarkAddr struct {
				DataSize  int         `json:"data_size"`
				FileCs    string      `json:"file_cs"`
				FileHash  string      `json:"file_hash"`
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlKey    string      `json:"url_key"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"download_no_watermark_addr"`
			Duration     int `json:"duration"`
			DynamicCover struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"dynamic_cover"`
			HasWatermark      bool   `json:"has_watermark"`
			Height            int    `json:"height"`
			IsBytevc1         int    `json:"is_bytevc1"`
			IsCallback        bool   `json:"is_callback"`
			Meta              string `json:"meta"`
			MiscDownloadAddrs string `json:"misc_download_addrs"`
			NeedSetToken      bool   `json:"need_set_token"`
			OriginCover       struct {
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"origin_cover"`
			PlayAddr struct {
				DataSize  int         `json:"data_size"`
				FileCs    string      `json:"file_cs"`
				FileHash  string      `json:"file_hash"`
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlKey    string      `json:"url_key"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"play_addr"`
			PlayAddrBytevc1 struct {
				DataSize  int         `json:"data_size"`
				FileCs    string      `json:"file_cs"`
				FileHash  string      `json:"file_hash"`
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlKey    string      `json:"url_key"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"play_addr_bytevc1"`
			PlayAddrH264 struct {
				DataSize  int         `json:"data_size"`
				FileCs    string      `json:"file_cs"`
				FileHash  string      `json:"file_hash"`
				Height    int         `json:"height"`
				Uri       string      `json:"uri"`
				UrlKey    string      `json:"url_key"`
				UrlList   []string    `json:"url_list"`
				UrlPrefix interface{} `json:"url_prefix"`
				Width     int         `json:"width"`
			} `json:"play_addr_h264"`
			Ratio            string      `json:"ratio"`
			SourceHDRType    int         `json:"source_HDR_type"`
			Tags             interface{} `json:"tags"`
			VidProfileLabels string      `json:"vid_profile_labels"`
			Width            int         `json:"width"`
		} `json:"video"`
		VideoControl struct {
			AllowDownload         bool `json:"allow_download"`
			AllowDuet             bool `json:"allow_duet"`
			AllowDynamicWallpaper bool `json:"allow_dynamic_wallpaper"`
			AllowMusic            bool `json:"allow_music"`
			AllowReact            bool `json:"allow_react"`
			AllowStitch           bool `json:"allow_stitch"`
			DraftProgressBar      int  `json:"draft_progress_bar"`
			PreventDownloadType   int  `json:"prevent_download_type"`
			ShareType             int  `json:"share_type"`
			ShowProgressBar       int  `json:"show_progress_bar"`
			TimerStatus           int  `json:"timer_status"`
		} `json:"video_control"`
		VideoLabels      []interface{} `json:"video_labels"`
		VideoText        []interface{} `json:"video_text"`
		VisualSearchInfo struct {
		} `json:"visual_search_info"`
		VoiceFilterIds       interface{} `json:"voice_filter_ids"`
		WithPromotionalMusic bool        `json:"with_promotional_music"`
		WithoutWatermark     bool        `json:"without_watermark"`
	} `json:"aweme_details"`
	AwemeStatus interface{} `json:"aweme_status"`
	Extra       struct {
		FatalItemIds []interface{} `json:"fatal_item_ids"`
		Logid        string        `json:"logid"`
		Now          int64         `json:"now"`
	} `json:"extra"`
	LogPb struct {
		ImprId string `json:"impr_id"`
	} `json:"log_pb"`
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
