package tikhub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type XHSNote struct {
	Id    string
	Title string
	Desc  string
	Cover string
	Url   string

	CollectedCount int
	CommentsCount  int
	LikedCount     int
	NiceCount      int
	SharedCount    int
	CreatedAt      time.Time

	User XHSUser
}

type XHSNotes []XHSNote

func (t Client) XhsGetUserNotes(ctx context.Context, userId string) (XHSNotes, error) {

	notes, err := t.xhsWebV2FetchHomeNotesApp(ctx, userId)
	if err == nil {
		return notes, err
	}

	return nil, nil
}

func (t Client) xhsWebV2FetchHomeNotesApp(ctx context.Context, userId string) (XHSNotes, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("user_id", userId).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		//Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_user_info")
		Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_home_notes_app")

	if err != nil {
		return nil, err
	}
	fmt.Println(r)

	var resp Response[xhsWebV2FetchHomeNotesAppResult]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}
	//
	//if resp.Data.Nickname == "" {
	//	return nil, errors.New("not found: " + userId)
	//}

	var notes XHSNotes
	for _, x := range resp.Data.Notes {
		notes = append(notes, XHSNote{
			Id:             x.Id,
			Title:          x.Title,
			Desc:           x.Desc,
			CollectedCount: x.CollectedCount,
			CommentsCount:  x.CommentsCount,
			LikedCount:     x.Likes,
			NiceCount:      x.NiceCount,
			SharedCount:    x.ShareCount,
			CreatedAt:      time.Unix(int64(x.CreateTime), 0),
		})
	}

	return notes, nil
}

type xhsWebV2FetchHomeNotesAppResult struct {
	HasMore bool `json:"has_more"`
	Notes   []struct {
		//AdvancedWidgetsGroups struct {
		//	Groups []struct {
		//		FetchTypes []string `json:"fetch_types"`
		//		Mode       int      `json:"mode"`
		//	} `json:"groups"`
		//} `json:"advanced_widgets_groups"`
		//Ats            []interface{} `json:"ats"`
		CollectedCount int `json:"collected_count"`
		CommentsCount  int `json:"comments_count"`
		CreateTime     int `json:"create_time"`
		//Cursor         string        `json:"cursor"`
		Desc         string `json:"desc"`
		DisplayTitle string `json:"display_title"`
		//HasMusic       bool          `json:"has_music"`
		Id string `json:"id"`
		//ImagesList     []struct {
		//	Fileid       string `json:"fileid"`
		//	Height       int    `json:"height"`
		//	Original     string `json:"original"`
		//	TraceId      string `json:"trace_id"`
		//	Url          string `json:"url"`
		//	UrlSizeLarge string `json:"url_size_large"`
		//	Width        int    `json:"width"`
		//} `json:"images_list"`
		//Infavs         bool `json:"infavs"`
		//Inlikes        bool `json:"inlikes"`
		//IsGoodsNote    bool `json:"is_goods_note"`
		//LastUpdateTime int  `json:"last_update_time"`
		//Level          int  `json:"level"`
		Likes     int `json:"likes"`
		NiceCount int `json:"nice_count"`
		//Niced          bool `json:"niced"`
		//Price          int  `json:"price"`
		//Recommend      struct {
		//	Desc       string `json:"desc"`
		//	Icon       string `json:"icon"`
		//	TargetId   string `json:"target_id"`
		//	TargetName string `json:"target_name"`
		//	TrackId    string `json:"track_id"`
		//	Type       string `json:"type"`
		//} `json:"recommend"`
		ShareCount int `json:"share_count"`
		//Sticky     bool   `json:"sticky"`
		Title string `json:"title"`
		//Type       string `json:"type"`
		//User       struct {
		//	Followed              bool   `json:"followed"`
		//	Fstatus               string `json:"fstatus"`
		//	Images                string `json:"images"`
		//	Nickname              string `json:"nickname"`
		//	RedOfficialVerifyType int    `json:"red_official_verify_type"`
		//	Userid                string `json:"userid"`
		//} `json:"user"`
		//ViewCount int `json:"view_count"`
		//WidgetsContext string `json:"widgets_context"`
	} `json:"notes"`
	//Tags []interface{} `json:"tags"`
}
