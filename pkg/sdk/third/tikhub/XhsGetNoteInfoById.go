package tikhub

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"github.com/go-kratos/kratos/v2/log"
)

type NoteInfo struct {
	Id      string
	Title   string
	Content string
	Url     string
	Images  []Image
	User    XHSUser
	Tags    []Tag

	ShareCount int
}

type Image struct {
	Url string
}

type Tag struct {
	Name string
}

func (t Client) XhsGetNoteInfoById(ctx context.Context, noteId string) (*NoteInfo, error) {

	note, err := t.xhsAppFetchFeedNotesV2(ctx, noteId)
	if err == nil {
		return note, nil
	}

	note, err = t.xhsAppGetNoteInfoV2(ctx, noteId)
	if err == nil {
		return note, nil
	}

	return nil, errors2.New("note not found: " + noteId)
}

func (t Client) xhsAppFetchFeedNotesV2(ctx context.Context, noteId string) (*NoteInfo, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("note_id", noteId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_feed_notes_v2")

	if err != nil {
		return nil, err
	}

	var resp Response[AppFetchFeedNotesV2Result]

	log.Debugw("r", r.String())

	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	if len(resp.Data.NoteList) == 0 {
		return nil, errors2.New("note not found: " + noteId)
	}

	d := resp.Data.NoteList[0]

	var images []Image
	for _, x := range d.ImagesList {
		images = append(images, Image{
			Url: x.Url,
		})
	}

	var tags []Tag
	for _, x := range d.HashTag {
		tags = append(tags, Tag{
			Name: x.Name,
		})
	}

	return &NoteInfo{
		Id:      d.Id,
		Title:   d.Title,
		Content: d.Desc,
		Url:     d.Video.Url,
		Images:  images,
		User: XHSUser{
			Id:     d.User.Id,
			Avatar: d.User.Image,
			Name:   d.User.Nickname,
		},
		Tags:       tags,
		ShareCount: d.SharedCount,
	}, nil
}

type AppFetchFeedNotesV2Result struct {
	CommentList []interface{} `json:"comment_list"`
	ModelType   string        `json:"model_type"`
	NoteList    []struct {
		ApiUpgrade       int           `json:"api_upgrade"`
		Ats              []interface{} `json:"ats"`
		Collected        bool          `json:"collected"`
		CollectedCount   int           `json:"collected_count"`
		CommentsCount    int           `json:"comments_count"`
		ContentTransInfo struct {
			Strategy int `json:"strategy"`
		} `json:"content_trans_info"`
		CooperateBinds        []interface{} `json:"cooperate_binds"`
		Countdown             int           `json:"countdown"`
		Desc                  string        `json:"desc"`
		EnableBrandLottery    bool          `json:"enable_brand_lottery"`
		EnableCoProduce       bool          `json:"enable_co_produce"`
		EnableFlsBridgeCards  bool          `json:"enable_fls_bridge_cards"`
		EnableFlsRelatedCards bool          `json:"enable_fls_related_cards"`
		FeedbackInfo          struct {
			DislikeStatus int `json:"dislike_status"`
		} `json:"feedback_info"`
		FootTags  []interface{} `json:"foot_tags"`
		GoodsInfo struct {
		} `json:"goods_info"`
		HasCoProduce    bool `json:"has_co_produce"`
		HasMusic        bool `json:"has_music"`
		HasRelatedGoods bool `json:"has_related_goods"`
		HashTag         []struct {
			BizId        string `json:"bizId"`
			CurrentScore int    `json:"current_score"`
			Id           string `json:"id"`
			Link         string `json:"link"`
			Name         string `json:"name"`
			RecordCount  int    `json:"record_count"`
			RecordEmoji  string `json:"record_emoji"`
			RecordUnit   string `json:"record_unit"`
			TagHint      string `json:"tag_hint"`
			Type         string `json:"type"`
		} `json:"hash_tag"`
		HeadTags   []interface{} `json:"head_tags"`
		Id         string        `json:"id"`
		ImagesList []struct {
			Fileid                string `json:"fileid"`
			Height                int    `json:"height"`
			Index                 int    `json:"index"`
			Latitude              int    `json:"latitude"`
			Longitude             int    `json:"longitude"`
			NeedLoadOriginalImage bool   `json:"need_load_original_image"`
			Original              string `json:"original"`
			ScaleToLarge          int    `json:"scale_to_large"`
			TextIntention         int    `json:"text_intention"`
			TraceId               string `json:"trace_id"`
			Url                   string `json:"url"`
			UrlMultiLevel         struct {
				High   string `json:"high"`
				Low    string `json:"low"`
				Medium string `json:"medium"`
			} `json:"url_multi_level"`
			Width int `json:"width"`
		} `json:"images_list"`
		InCensor           bool          `json:"in_censor"`
		IpLocation         string        `json:"ip_location"`
		LastUpdateTime     int           `json:"last_update_time"`
		Liked              bool          `json:"liked"`
		LikedCount         int           `json:"liked_count"`
		LikedUsers         []interface{} `json:"liked_users"`
		LongPressShareInfo struct {
			BlockPrivateMsg bool   `json:"block_private_msg"`
			Content         string `json:"content"`
			FunctionEntries []struct {
				Type string `json:"type"`
			} `json:"function_entries"`
			GuideAudited  bool   `json:"guide_audited"`
			IsStar        bool   `json:"is_star"`
			ShowWechatTag bool   `json:"show_wechat_tag"`
			Title         string `json:"title"`
		} `json:"long_press_share_info"`
		MayHaveRedPacket bool `json:"may_have_red_packet"`
		MediaSaveConfig  struct {
			DisableSave       bool `json:"disable_save"`
			DisableWatermark  bool `json:"disable_watermark"`
			DisableWeiboCover bool `json:"disable_weibo_cover"`
		} `json:"media_save_config"`
		MiniProgramInfo struct {
			Desc       string `json:"desc"`
			Path       string `json:"path"`
			ShareTitle string `json:"share_title"`
			Thumb      string `json:"thumb"`
			Title      string `json:"title"`
			UserName   string `json:"user_name"`
			WebpageUrl string `json:"webpage_url"`
		} `json:"mini_program_info"`
		ModelType       string `json:"model_type"`
		NativeVoiceInfo struct {
			Cover                 string `json:"cover"`
			Desc                  string `json:"desc"`
			Duration              int    `json:"duration"`
			Md5Sum                string `json:"md5sum"`
			Name                  string `json:"name"`
			OptimizeExpInUseCount int    `json:"optimize_exp_in_use_count"`
			SoundBgMusicType      int    `json:"sound_bg_music_type"`
			SoundId               string `json:"sound_id"`
			Subtitle              string `json:"subtitle"`
			Url                   string `json:"url"`
			UseCount              int    `json:"use_count"`
		} `json:"native_voice_info"`
		NeedNextStep         bool `json:"need_next_step"`
		NeedProductReview    bool `json:"need_product_review"`
		NoteTextPressOptions []struct {
			Extra string `json:"extra"`
			Key   string `json:"key"`
		} `json:"note_text_press_options"`
		Privacy struct {
			NickNames string `json:"nick_names"`
			ShowTips  bool   `json:"show_tips"`
			Type      int    `json:"type"`
		} `json:"privacy"`
		QqMiniProgramInfo struct {
			Desc       string `json:"desc"`
			Path       string `json:"path"`
			ShareTitle string `json:"share_title"`
			Thumb      string `json:"thumb"`
			Title      string `json:"title"`
			UserName   string `json:"user_name"`
			WebpageUrl string `json:"webpage_url"`
		} `json:"qq_mini_program_info"`
		RedEnvelopeNote bool `json:"red_envelope_note"`
		Seeded          bool `json:"seeded"`
		SeededCount     int  `json:"seeded_count"`
		ShareCodeFlag   int  `json:"share_code_flag"`
		ShareInfo       struct {
			BlockPrivateMsg bool   `json:"block_private_msg"`
			Content         string `json:"content"`
			FunctionEntries []struct {
				Type string `json:"type"`
			} `json:"function_entries"`
			GuideAudited  bool   `json:"guide_audited"`
			Image         string `json:"image"`
			IsStar        bool   `json:"is_star"`
			Link          string `json:"link"`
			ShowWechatTag bool   `json:"show_wechat_tag"`
			Title         string `json:"title"`
		} `json:"share_info"`
		SharedCount      int    `json:"shared_count"`
		Sticky           bool   `json:"sticky"`
		TextLanguageCode string `json:"text_language_code"`
		Time             int    `json:"time"`
		Title            string `json:"title"`
		Topics           []struct {
			ActivityOnline bool   `json:"activity_online"`
			BusinessType   int    `json:"business_type"`
			DiscussNum     int    `json:"discuss_num"`
			Id             string `json:"id"`
			Image          string `json:"image"`
			Link           string `json:"link"`
			Name           string `json:"name"`
			Style          int    `json:"style"`
		} `json:"topics"`
		Type          string `json:"type"`
		UseWaterColor bool   `json:"use_water_color"`
		User          struct {
			Followed bool   `json:"followed"`
			Fstatus  string `json:"fstatus"`
			Id       string `json:"id"`
			Image    string `json:"image"`
			Level    struct {
				Image string `json:"image"`
			} `json:"level"`
			Name                      string `json:"name"`
			Nickname                  string `json:"nickname"`
			RedId                     string `json:"red_id"`
			RedOfficialVerified       bool   `json:"red_official_verified"`
			RedOfficialVerifyType     int    `json:"red_official_verify_type"`
			ShowRedOfficialVerifyIcon bool   `json:"show_red_official_verify_icon"`
			TrackDuration             int    `json:"track_duration"`
			Userid                    string `json:"userid"`
		} `json:"user"`
		Video struct {
			AdaptiveStreamingUrlSet []interface{} `json:"adaptive_streaming_url_set"`
			AvgBitrate              int           `json:"avg_bitrate"`
			CanSuperResolution      bool          `json:"can_super_resolution"`
			Duration                int           `json:"duration"`
			FirstFrame              string        `json:"first_frame"`
			FrameTs                 int           `json:"frame_ts"`
			Height                  int           `json:"height"`
			Id                      string        `json:"id"`
			IsUpload                bool          `json:"is_upload"`
			IsUserSelect            bool          `json:"is_user_select"`
			PlayedCount             int           `json:"played_count"`
			PreloadSize             int           `json:"preload_size"`
			Thumbnail               string        `json:"thumbnail"`
			ThumbnailDim            string        `json:"thumbnail_dim"`
			Url                     string        `json:"url"`
			UrlInfoList             []struct {
				AvgBitrate int    `json:"avg_bitrate"`
				Desc       string `json:"desc"`
				Height     int    `json:"height"`
				Url        string `json:"url"`
				Vmaf       int    `json:"vmaf"`
				Width      int    `json:"width"`
			} `json:"url_info_list"`
			Vmaf   int `json:"vmaf"`
			Volume int `json:"volume"`
			Width  int `json:"width"`
		} `json:"video"`
		ViewCount      int        `json:"view_count"`
		WidgetsContext string     `json:"widgets_context"`
		WidgetsGroups  [][]string `json:"widgets_groups"`
	} `json:"note_list"`
	TrackId string `json:"track_id"`
	User    struct {
		Followed bool   `json:"followed"`
		Fstatus  string `json:"fstatus"`
		Id       string `json:"id"`
		Image    string `json:"image"`
		Level    struct {
			Image string `json:"image"`
		} `json:"level"`
		Name                      string `json:"name"`
		Nickname                  string `json:"nickname"`
		RedId                     string `json:"red_id"`
		RedOfficialVerified       bool   `json:"red_official_verified"`
		RedOfficialVerifyType     int    `json:"red_official_verify_type"`
		ShowRedOfficialVerifyIcon bool   `json:"show_red_official_verify_icon"`
		TrackDuration             int    `json:"track_duration"`
		Userid                    string `json:"userid"`
	} `json:"user"`
}

func (t Client) xhsAppGetNoteInfoV2(ctx context.Context, noteId string) (*NoteInfo, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("note_id", noteId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/app/get_note_info_v2")

	if err != nil {
		return nil, err
	}

	var resp Response[AppGetNoteInfoV2Result]

	log.Debugw("r", r.String())

	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	if resp.Data.Data.NoteId == "" {
		return nil, errors2.New("note not found: " + resp.Data.Data.NoteId)
	}

	d := resp.Data.Data

	var images []Image
	for _, x := range d.ImagesList {
		images = append(images, Image{
			Url: x.Url,
		})
	}

	return &NoteInfo{
		Id:      d.NoteId,
		Title:   d.Title,
		Content: d.Content,
		Url:     d.VideoInfo.VideoUrl,
		Images:  images,
		User: XHSUser{
			Id:     d.UserId,
			Avatar: d.UserInfo.Avatar,
			Name:   d.UserInfo.NickName,
			Fans:   d.FavNum,
		},
	}, nil
}

type AppGetNoteInfoV2Result struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Guid    interface{} `json:"guid"`
	Success bool        `json:"success"`
	Data    struct {
		NoteId     string      `json:"noteId"`
		NoteLink   string      `json:"noteLink"`
		UserId     string      `json:"userId"`
		HeadPhoto  interface{} `json:"headPhoto"`
		Name       interface{} `json:"name"`
		RedId      interface{} `json:"redId"`
		Type       int         `json:"type"`
		AtUserList interface{} `json:"atUserList"`
		Title      string      `json:"title"`
		Content    string      `json:"content"`
		ImagesList []struct {
			FileId    string      `json:"fileId"`
			Url       string      `json:"url"`
			Original  string      `json:"original"`
			Width     int         `json:"width"`
			Height    int         `json:"height"`
			Latitude  interface{} `json:"latitude"`
			Longitude interface{} `json:"longitude"`
			TraceId   string      `json:"traceId"`
			Sticker   interface{} `json:"sticker"`
		} `json:"imagesList"`
		VideoInfo struct {
			Id             string `json:"id"`
			VideoKey       string `json:"videoKey"`
			OriginVideoKey string `json:"originVideoKey"`
			Meta           struct {
				Width    int `json:"width"`
				Height   int `json:"height"`
				Duration int `json:"duration"`
			} `json:"meta"`
			GifKey       string        `json:"gifKey"`
			VideoUrl     string        `json:"videoUrl"`
			GifUrl       string        `json:"gifUrl"`
			VideoKeyList []interface{} `json:"videoKeyList"`
			HasFragments bool          `json:"hasFragments"`
			Thumbnail    string        `json:"thumbnail"`
			FirstFrame   string        `json:"firstFrame"`
			Volume       float64       `json:"volume"`
			Chapters     interface{}   `json:"chapters"`
		} `json:"videoInfo"`
		Time struct {
			CreateTime     int64 `json:"createTime"`
			UpdateTime     int64 `json:"updateTime"`
			UserUpdateTime int64 `json:"userUpdateTime"`
		} `json:"time"`
		CreateTime        string      `json:"createTime"`
		ImpNum            int         `json:"impNum"`
		LikeNum           int         `json:"likeNum"`
		FavNum            int         `json:"favNum"`
		CmtNum            int         `json:"cmtNum"`
		ReadNum           int         `json:"readNum"`
		ShareNum          int         `json:"shareNum"`
		FollowCnt         int         `json:"followCnt"`
		ReportBrandUserId interface{} `json:"reportBrandUserId"`
		ReportBrandName   interface{} `json:"reportBrandName"`
		FeatureTags       interface{} `json:"featureTags"`
		UserInfo          struct {
			NickName       string        `json:"nickName"`
			Avatar         string        `json:"avatar"`
			UserId         string        `json:"userId"`
			AdvertiserId   interface{}   `json:"advertiserId"`
			FansNum        int           `json:"fansNum"`
			CooperType     int           `json:"cooperType"`
			PriceState     interface{}   `json:"priceState"`
			PictureState   interface{}   `json:"pictureState"`
			PicturePrice   interface{}   `json:"picturePrice"`
			VideoState     interface{}   `json:"videoState"`
			VideoPrice     interface{}   `json:"videoPrice"`
			UserType       int           `json:"userType"`
			OperateState   interface{}   `json:"operateState"`
			CurrentLevel   interface{}   `json:"currentLevel"`
			Location       string        `json:"location"`
			ContentTags    []interface{} `json:"contentTags"`
			FeatureTags    []interface{} `json:"featureTags"`
			PersonalTags   []interface{} `json:"personalTags"`
			Gender         string        `json:"gender"`
			IsCollect      bool          `json:"isCollect"`
			ClickMidNum    int           `json:"clickMidNum"`
			InterMidNum    int           `json:"interMidNum"`
			PictureInCart  interface{}   `json:"pictureInCart"`
			VideoInCart    interface{}   `json:"videoInCart"`
			KolType        interface{}   `json:"kolType"`
			MEngagementNum int           `json:"mEngagementNum"`
		} `json:"userInfo"`
		CompClickData interface{} `json:"compClickData"`
	} `json:"data"`
}
