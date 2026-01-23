package tikhub

import (
	"context"
	"encoding/json"
	"fmt"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"time"
)

type XhsSearchNotesParams struct {
	Keyword string
	//general：综合
	//time_descending：最新
	//popularity_descending：最热
	Sort string
	//0：全部
	//1：视频
	//2：图文
	NoteType string
}

func (t Client) XhsSearchNotes(ctx context.Context, params XhsSearchNotesParams) ([]XHSNote, error) {

	notes, err := t.xhsWebV2FetchSearchNotes(ctx, params)
	if err == nil {
		return notes, nil
	}

	return nil, nil
}

func (t Client) xhsWebV2FetchSearchNotes(ctx context.Context, params XhsSearchNotesParams) ([]XHSNote, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("keywords", params.Keyword).
		SetQueryParam("sort_type", helper.OrString(params.Sort, "popularity_descending")).
		SetQueryParam("note_type", helper.OrString(params.NoteType, "0")).
		SetQueryParam("page", "1").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_search_notes")
	//Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_user_info_v3")

	if err != nil {
		return nil, err
	}

	fmt.Println(r.String())

	var resp Response[xhsWebV2FetchSearchNotesResult]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	var notes []XHSNote
	for _, x := range resp.Data.Items {
		notes = append(notes, XHSNote{
			Id:             x.Id,
			Title:          x.NoteCard.DisplayTitle,
			Cover:          x.NoteCard.Cover.UrlPre,
			Desc:           "",
			CollectedCount: conv.Int(x.NoteCard.InteractInfo.CollectedCount),
			CommentsCount:  conv.Int(x.NoteCard.InteractInfo.CommentCount),
			LikedCount:     conv.Int(x.NoteCard.InteractInfo.LikedCount),
			SharedCount:    conv.Int(x.NoteCard.InteractInfo.SharedCount),
			CreatedAt:      time.Time{},
			User: XHSUser{
				Id:     x.NoteCard.User.UserId,
				Avatar: x.NoteCard.User.Avatar,

				ShowRedOfficialVerifyIcon: false,
				RedOfficialVerified:       false,
				Self:                      false,
				Reason:                    "",
				RedOfficialVerifyType:     0,
				Name:                      x.NoteCard.User.NickName,
				Desc:                      "",
				Image:                     "",
				Followed:                  false,
				RedId:                     "",
				Link:                      "",
				SubTitle:                  "",
			},
		})
	}

	return notes, nil
}

type xhsWebV2FetchSearchNotesResult struct {
	HasMore bool `json:"has_more"`
	Items   []struct {
		NoteCard struct {
			Type         string `json:"type"`
			DisplayTitle string `json:"display_title"`
			User         struct {
				Nickname  string `json:"nickname"`
				XsecToken string `json:"xsec_token"`
				NickName  string `json:"nick_name"`
				Avatar    string `json:"avatar"`
				UserId    string `json:"user_id"`
			} `json:"user"`
			InteractInfo struct {
				SharedCount    string `json:"shared_count"`
				Liked          bool   `json:"liked"`
				LikedCount     string `json:"liked_count"`
				Collected      bool   `json:"collected"`
				CollectedCount string `json:"collected_count"`
				CommentCount   string `json:"comment_count"`
			} `json:"interact_info"`
			Cover struct {
				Height     int    `json:"height"`
				Width      int    `json:"width"`
				UrlDefault string `json:"url_default"`
				UrlPre     string `json:"url_pre"`
			} `json:"cover"`
			ImageList []struct {
				Height   int `json:"height"`
				Width    int `json:"width"`
				InfoList []struct {
					ImageScene string `json:"image_scene"`
					Url        string `json:"url"`
				} `json:"info_list"`
			} `json:"image_list"`
			CornerTagInfo []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"corner_tag_info"`
		} `json:"note_card"`
		XsecToken string `json:"xsec_token"`
		Id        string `json:"id"`
		ModelType string `json:"model_type"`
	} `json:"items"`
}
