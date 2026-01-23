package tikhub

import (
	"context"
	"encoding/json"
	"errors"
)

func (t Client) XhsWebGetNoteInfoByIdV4(ctx context.Context, noteId string) (*Note, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("note_id", noteId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_note_info_v4")

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

type Note struct {
	Code int `json:"code"`
	Data []struct {
		TrackId   string `json:"track_id"`
		ModelType string `json:"model_type"`
		User      struct {
			GroupId                   string `json:"group_id"`
			Nickname                  string `json:"nickname"`
			ShowRedOfficialVerifyIcon bool   `json:"show_red_official_verify_icon"`
			Level                     struct {
				Image string `json:"image"`
			} `json:"level"`
			Name                  string `json:"name"`
			Image                 string `json:"image"`
			RedId                 string `json:"red_id"`
			RedOfficialVerified   bool   `json:"red_official_verified"`
			Id                    string `json:"id"`
			Followed              bool   `json:"followed"`
			Fstatus               string `json:"fstatus"`
			RedOfficialVerifyType int    `json:"red_official_verify_type"`
			Userid                string `json:"userid"`
			TrackDuration         int    `json:"track_duration"`
		} `json:"user"`
		NoteList []struct {
			HasRelatedGoods bool `json:"has_related_goods"`
			MiniProgramInfo struct {
				Thumb      string `json:"thumb"`
				ShareTitle string `json:"share_title"`
				UserName   string `json:"user_name"`
				Path       string `json:"path"`
				Title      string `json:"title"`
				Desc       string `json:"desc"`
				WebpageUrl string `json:"webpage_url"`
			} `json:"mini_program_info"`
			ViewCount       int        `json:"view_count"`
			RedEnvelopeNote bool       `json:"red_envelope_note"`
			ModelType       string     `json:"model_type"`
			Id              string     `json:"id"`
			WidgetsGroups   [][]string `json:"widgets_groups"`
			NeedNextStep    bool       `json:"need_next_step"`
			GoodsInfo       struct {
			} `json:"goods_info"`
			FootTags              []interface{} `json:"foot_tags"`
			InCensor              bool          `json:"in_censor"`
			EnableFlsRelatedCards bool          `json:"enable_fls_related_cards"`
			User                  struct {
				Image                     string `json:"image"`
				Followed                  bool   `json:"followed"`
				RedOfficialVerified       bool   `json:"red_official_verified"`
				Id                        string `json:"id"`
				Userid                    string `json:"userid"`
				TrackDuration             int    `json:"track_duration"`
				Fstatus                   string `json:"fstatus"`
				RedId                     string `json:"red_id"`
				Name                      string `json:"name"`
				GroupId                   string `json:"group_id"`
				Nickname                  string `json:"nickname"`
				RedOfficialVerifyType     int    `json:"red_official_verify_type"`
				ShowRedOfficialVerifyIcon bool   `json:"show_red_official_verify_icon"`
				Level                     struct {
					Image string `json:"image"`
				} `json:"level"`
			} `json:"user"`
			Collected      bool          `json:"collected"`
			WidgetsContext string        `json:"widgets_context"`
			CooperateBinds []interface{} `json:"cooperate_binds"`
			HasMusic       bool          `json:"has_music"`
			Time           int           `json:"time"`
			LastUpdateTime int           `json:"last_update_time"`
			Topics         []struct {
				Image          string `json:"image"`
				Link           string `json:"link"`
				ActivityOnline bool   `json:"activity_online"`
				Style          int    `json:"style"`
				DiscussNum     int    `json:"discuss_num"`
				BusinessType   int    `json:"business_type"`
				Id             string `json:"id"`
				Name           string `json:"name"`
			} `json:"topics"`
			UseWaterColor        bool          `json:"use_water_color"`
			SharedCount          int           `json:"shared_count"`
			HeadTags             []interface{} `json:"head_tags"`
			LikedUsers           []interface{} `json:"liked_users"`
			ApiUpgrade           int           `json:"api_upgrade"`
			NoteTextPressOptions []struct {
				Key   string `json:"key"`
				Extra string `json:"extra"`
			} `json:"note_text_press_options"`
			SeededCount          int           `json:"seeded_count"`
			EnableFlsBridgeCards bool          `json:"enable_fls_bridge_cards"`
			MayHaveRedPacket     bool          `json:"may_have_red_packet"`
			Ats                  []interface{} `json:"ats"`
			Liked                bool          `json:"liked"`
			EnableCoProduce      bool          `json:"enable_co_produce"`
			HasCoProduce         bool          `json:"has_co_produce"`
			Type                 string        `json:"type"`
			Title                string        `json:"title"`
			EnableBrandLottery   bool          `json:"enable_brand_lottery"`
			ShareCodeFlag        int           `json:"share_code_flag"`
			Privacy              struct {
				NickNames string `json:"nick_names"`
				Type      int    `json:"type"`
				ShowTips  bool   `json:"show_tips"`
			} `json:"privacy"`
			Seeded    bool `json:"seeded"`
			ShareInfo struct {
				ShowWechatTag   bool `json:"show_wechat_tag"`
				FunctionEntries []struct {
					Type string `json:"type"`
				} `json:"function_entries"`
				Content         string `json:"content"`
				Link            string `json:"link"`
				IsStar          bool   `json:"is_star"`
				GuideAudited    bool   `json:"guide_audited"`
				Image           string `json:"image"`
				Title           string `json:"title"`
				BlockPrivateMsg bool   `json:"block_private_msg"`
			} `json:"share_info"`
			NeedProductReview  bool `json:"need_product_review"`
			LongPressShareInfo struct {
				Content         string `json:"content"`
				Title           string `json:"title"`
				IsStar          bool   `json:"is_star"`
				BlockPrivateMsg bool   `json:"block_private_msg"`
				ShowWechatTag   bool   `json:"show_wechat_tag"`
				FunctionEntries []struct {
					Type string `json:"type"`
				} `json:"function_entries"`
				GuideAudited bool `json:"guide_audited"`
			} `json:"long_press_share_info"`
			Sticky           bool   `json:"sticky"`
			IpLocation       string `json:"ip_location"`
			Countdown        int    `json:"countdown"`
			TextLanguageCode string `json:"text_language_code"`
			Video            struct {
				Height      int    `json:"height"`
				Width       int    `json:"width"`
				Url         string `json:"url"`
				AvgBitrate  int    `json:"avg_bitrate"`
				UrlInfoList []struct {
					Width      int    `json:"width"`
					Height     int    `json:"height"`
					Vmaf       int    `json:"vmaf"`
					Desc       string `json:"desc"`
					Url        string `json:"url"`
					AvgBitrate int    `json:"avg_bitrate"`
				} `json:"url_info_list"`
				FirstFrame              string        `json:"first_frame"`
				Id                      string        `json:"id"`
				PreloadSize             int           `json:"preload_size"`
				PlayedCount             int           `json:"played_count"`
				IsUpload                bool          `json:"is_upload"`
				Volume                  int           `json:"volume"`
				Vmaf                    int           `json:"vmaf"`
				ThumbnailDim            string        `json:"thumbnail_dim"`
				AdaptiveStreamingUrlSet []interface{} `json:"adaptive_streaming_url_set"`
				FrameTs                 int           `json:"frame_ts"`
				IsUserSelect            bool          `json:"is_user_select"`
				Thumbnail               string        `json:"thumbnail"`
				CanSuperResolution      bool          `json:"can_super_resolution"`
				Duration                int           `json:"duration"`
			} `json:"video"`
			CollectedCount int `json:"collected_count"`
			ImagesList     []struct {
				Width         int    `json:"width"`
				Original      string `json:"original"`
				Latitude      int    `json:"latitude"`
				UrlMultiLevel struct {
					High   string `json:"high"`
					Low    string `json:"low"`
					Medium string `json:"medium"`
				} `json:"url_multi_level"`
				Index                 int    `json:"index"`
				Longitude             int    `json:"longitude"`
				TraceId               string `json:"trace_id"`
				NeedLoadOriginalImage bool   `json:"need_load_original_image"`
				Fileid                string `json:"fileid"`
				Height                int    `json:"height"`
				Url                   string `json:"url"`
				ScaleToLarge          int    `json:"scale_to_large"`
			} `json:"images_list"`
			LikedCount    int    `json:"liked_count"`
			CommentsCount int    `json:"comments_count"`
			Desc          string `json:"desc"`
			HashTag       []struct {
				Link         string `json:"link"`
				RecordCount  int    `json:"record_count"`
				RecordUnit   string `json:"record_unit"`
				CurrentScore int    `json:"current_score"`
				BizId        string `json:"bizId"`
				Id           string `json:"id"`
				Name         string `json:"name"`
				Type         string `json:"type"`
				RecordEmoji  string `json:"record_emoji"`
				TagHint      string `json:"tag_hint"`
			} `json:"hash_tag"`
			MediaSaveConfig struct {
				DisableSave       bool `json:"disable_save"`
				DisableWatermark  bool `json:"disable_watermark"`
				DisableWeiboCover bool `json:"disable_weibo_cover"`
			} `json:"media_save_config"`
		} `json:"note_list"`
		CommentList []interface{} `json:"comment_list"`
	} `json:"data"`
	Message    interface{} `json:"message"`
	RecordTime string      `json:"recordTime"`
}
