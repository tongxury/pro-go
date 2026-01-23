package tikhub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

func (t Client) XhsSearchUsers(ctx context.Context, keyword string) ([]XHSUser, error) {

	users, err := t.xhsWebSearchUsers(ctx, keyword)
	if err == nil {
		return users, nil
	}

	users, err = t.xhsAppSearchUsers(ctx, keyword)
	if err == nil {
		return users, nil
	}

	users, err = t.xhsWebV2SearchUsers(ctx, keyword)
	if err == nil {
		return users, nil
	}

	return nil, err
}

func (t Client) xhsAppSearchUsers(ctx context.Context, keyword string) ([]XHSUser, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("keyword", keyword).
		SetQueryParam("page", "1").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/app/search_users")
	//Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_user_info_v3")

	if err != nil {
		return nil, err
	}

	fmt.Println(r.String())

	var resp Response[XhsAppSearchUsersResult]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	//x := resp.Data

	if len(resp.Data.Data.Users) == 0 {
		return nil, errors.New("not found: " + keyword)
	}

	var users []XHSUser
	for _, x := range resp.Data.Data.Users {
		users = append(users, x)
	}

	return users, nil
}

type XhsAppSearchUsersResult struct {
	Code    int  `json:"code"`
	Success bool `json:"success"`
	Data    struct {
		Users   []XHSUser `json:"users"`
		Filters []struct {
			Name          string `json:"name"`
			Id            string `json:"id"`
			GroupShowType int    `json:"group_show_type"`
			FilterTags    []struct {
				Id                   string      `json:"id"`
				OriginText           string      `json:"origin_text"`
				IconUrlNight         string      `json:"icon_url_night"`
				IconTailUrl          string      `json:"icon_tail_url"`
				IconTailUrlNight     string      `json:"icon_tail_url_night"`
				SubFilters           interface{} `json:"sub_filters"`
				Name                 string      `json:"name"`
				IconUrl              string      `json:"icon_url"`
				NeedLocationInfo     bool        `json:"need_location_info"`
				SubFiltersSelectType string      `json:"sub_filters_select_type"`
			} `json:"filter_tags"`
			WordRequestId string `json:"word_request_id"`
			Invisible     bool   `json:"invisible"`
			Type          string `json:"type"`
		} `json:"filters"`
	} `json:"data"`
	SearchId string `json:"searchId"`
}

func (t Client) xhsWebV2SearchUsers(ctx context.Context, keyword string) ([]XHSUser, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("keywords", keyword).
		SetQueryParam("page", "1").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/web_v2/fetch_search_users")
	//Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_user_info_v3")

	if err != nil {
		return nil, err
	}

	fmt.Println(r.String())

	var resp Response[XhsWebV2SearchUsersResult]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	//x := resp.Data

	if len(resp.Data.Users) == 0 {
		return nil, errors.New("not found: " + keyword)
	}

	var users []XHSUser
	for _, x := range resp.Data.Users {
		users = append(users, XHSUser{
			ShowRedOfficialVerifyIcon: x.ShowRedOfficialVerifyIcon,
			RedOfficialVerified:       x.RedOfficialVerified,
			Self:                      x.IsSelf,
			//Reason:                    "",
			RedOfficialVerifyType: x.RedOfficialVerifyType,
			Name:                  x.Name,
			Desc:                  x.Profession,
			Image:                 x.Image,
			Followed:              x.Followed,
			RedId:                 x.RedId,
			Link:                  x.Link,
			SubTitle:              x.SubTitle,
			Id:                    x.Id,
		})
	}

	return users, nil
}

type XhsWebV2SearchUsersResult struct {
	Users []struct {
		Fans                      string `json:"fans"`
		UpdateTime                string `json:"update_time,omitempty"`
		XsecToken                 string `json:"xsec_token"`
		Id                        string `json:"id"`
		RedId                     string `json:"red_id"`
		RedOfficialVerifyType     int    `json:"red_official_verify_type"`
		ShowRedOfficialVerifyIcon bool   `json:"show_red_official_verify_icon"`
		RedOfficialVerified       bool   `json:"red_official_verified"`
		IsSelf                    bool   `json:"is_self"`
		Vshow                     int    `json:"vshow"`
		Profession                string `json:"profession,omitempty"`
		Name                      string `json:"name"`
		Image                     string `json:"image"`
		SubTitle                  string `json:"sub_title"`
		Link                      string `json:"link"`
		Followed                  bool   `json:"followed"`
		NoteCount                 int    `json:"note_count"`
	} `json:"users"`
	HasMore bool `json:"has_more"`
	Result  struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
}

func (t Client) xhsWebSearchUsers(ctx context.Context, keyword string) ([]XHSUser, error) {

	r, err := t.c.R().SetContext(ctx).
		SetQueryParam("keyword", keyword).
		SetQueryParam("page", "1").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+t.apiKey).
		Get(t.baseUrl + "/api/v1/xiaohongshu/web/search_users")
	//Get(t.baseUrl + "/api/v1/xiaohongshu/web/get_user_info_v3")

	if err != nil {
		return nil, err
	}

	fmt.Println(r.String())

	var resp Response[XhsWebSearchUsersResult]
	if err = json.Unmarshal(r.Body(), &resp); err != nil {
		return nil, err
	}

	//x := resp.Data

	if len(resp.Data.Data.Users) == 0 {
		return nil, errors.New("not found: " + keyword)
	}

	var users []XHSUser
	for _, x := range resp.Data.Data.Users {
		users = append(users, XHSUser{
			ShowRedOfficialVerifyIcon: x.ShowRedOfficialVerifyIcon,
			RedOfficialVerified:       x.RedOfficialVerified,
			Self:                      x.Self,
			//Reason:                    "",
			RedOfficialVerifyType: x.RedOfficialVerifyType,
			Name:                  x.Name,
			Desc:                  x.Desc,
			Image:                 x.Image,
			Followed:              x.Followed,
			RedId:                 x.RedId,
			Link:                  x.Link,
			SubTitle:              x.SubTitle,
			Id:                    x.Id,
		})
	}

	return users, nil
}

type XhsWebSearchUsersResult struct {
	Code int `json:"code"`
	Data struct {
		Users []struct {
			RedId                     string `json:"red_id"`
			RedOfficialVerifyType     int    `json:"red_official_verify_type"`
			ShowRedOfficialVerifyIcon bool   `json:"show_red_official_verify_icon"`
			RedOfficialVerified       bool   `json:"red_official_verified"`
			Self                      bool   `json:"self"`
			Name                      string `json:"name"`
			Image                     string `json:"image"`
			Followed                  bool   `json:"followed"`
			SubTitle                  string `json:"sub_title"`
			Reason                    string `json:"reason"`
			Id                        string `json:"id"`
			Desc                      string `json:"desc"`
			Link                      string `json:"link"`
		} `json:"users"`
		Filters []struct {
			WordRequestId string `json:"word_request_id"`
			Invisible     bool   `json:"invisible"`
			Type          string `json:"type"`
			Name          string `json:"name"`
			Id            string `json:"id"`
			GroupShowType int    `json:"group_show_type"`
			FilterTags    []struct {
				IconUrl              string      `json:"icon_url"`
				IconUrlNight         string      `json:"icon_url_night"`
				IconTailUrlNight     string      `json:"icon_tail_url_night"`
				NeedLocationInfo     bool        `json:"need_location_info"`
				SubFiltersSelectType string      `json:"sub_filters_select_type"`
				SubFilters           interface{} `json:"sub_filters"`
				Id                   string      `json:"id"`
				Name                 string      `json:"name"`
				OriginText           string      `json:"origin_text"`
				IconTailUrl          string      `json:"icon_tail_url"`
			} `json:"filter_tags"`
		} `json:"filters"`
	} `json:"data"`
	Message    interface{} `json:"message"`
	RecordTime string      `json:"recordTime"`
}
