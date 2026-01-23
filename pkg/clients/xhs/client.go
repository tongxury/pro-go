package xhs

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/context"
	netUrl "net/url"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/goqueryz"
	"store/pkg/sdk/helper"
	"strings"
	"time"
)

type Client struct {
	headers map[string]string
}

func NewClient() *Client {
	//xmns := "awzSYIh3W5ebFaLMeLyBHxpQKugcSljKHlLK0+FwgpCY5gFRkCMK93Zy/jxcitJyaJFfxcpmKxybbxcyRxWpSM6BkzlFXm86F2cjhhM0MEFdEfM5aXSjNkJOxzj2cMXGhzGnBtpTd5QePhg8DX4iEWFxBGLdCYf2IXDQh8DjhYo6bv3d+mQcRXLzEL2GImN9d9l7Klzy3kN72MwoDR689cFCIlJF1ZH8jh0oLZ1t1FuaP4WwP7Wnw2SfpFa6fh6Dea8Kh+WufyeIcJI+ftBI17FT9xF35oFb9HN7C/C0Sw9cc7o4iHJTxt7YoXWe5dPkG7CB8b9n/kTY4ItNFft0vjgx4fvYIY1kLvQPcQ8oFN/Mwj3vgnNLj57D35gyOeil5FXMgcn8c6pDiI6cFjYy6XLgF6QtTEhg6b6uunOzxG3wP3JQ3unJwJgPm24La0mWmbthJ9ISd/3YRM+6HlPEwIuiGfoithvogmxLhPn3C5fWZIbjGi/eKK6QwoKW"
	//xs := "XYW_eyJzaWduU3ZuIjoiNTYiLCJzaWduVHlwZSI6IngyIiwiYXBwSWQiOiJ4aHMtcGMtd2ViIiwic2lnblZlcnNpb24iOiIxIiwicGF5bG9hZCI6IjNmNTdjMGQ0ZDI1M2ZmOWRhZWZmNDNlODA5M2E4MjY2MDEyYjg0ZmFlYTU3MmM3NmJhNjdjOTAyYzVmZTEyMGIzNjViMmM3OGRjMTkzYzQzMWMyNmM1MWMyZjg1ODYxNzE0MThhYzVjOTZmYmNmNjljMGU2MzJlMzI5ZjIzZTk5YzY5NWRlMzg3MDE1MTEwZjJkMGM3MzgxNjVkY2I1YjdkZmJjMzY3MWM4MDE3ZGQ3NDBiNjFiNzQ3ZjAwM2I5ZjAzYmM0YjZmNTBiZWY4MGJiMjMyMWQxMTI2Zjk2MGVlNGYxYmM4MjM3Y2M3YWUxMDM5MzM1MmQ2ZjQ4MTVjYWUzMmQ1ZDA0YzU1MjAyM2U3MTgzMDg1ZDM1ZWZkYjM4ZTViNDkxYjkwZDhmOWZjNjdlMDVhODMwODMxOTg1ZDgzYjMwNGYyYmY1OTdmNzZiNWI4ODMzZWRmNGYyNzAyZWY5M2FmYTJhMzIzNWRmMWZkZDQ0Y2Y4MWZkNzRlN2UwOGM0ZTcyYjQxOGMyYWQ3MzA1OTRhNTNhMjIxNGI1YmU3In0="
	//xscommon := "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"
	//cookie := "abRequestId=ffa7ca92-1757-5b56-8a73-5570627c1d52; a1=194e40ecd54z2gvo7ea38d9i1fnhxdr5sqnssy5ha30000336056; webId=7486224fc1ddf358b401f55d4a34195f; gid=yj4d48ddiKMDyj4d48dSfAv224uJkhlWd0qYIx0fj3yiIMq8xCfF2M888qqK82K8JJDSijYi; x-user-id-creator.xiaohongshu.com=5cdf90a10000000011023c23; customerClientId=295944940634645; web_session=0400698dabd3e9146c759b4d92354bacfd281e; access-token-creator.xiaohongshu.com=customer.creator.AT-68c517485603822268730262c8yqfbuibzpctnyb; galaxy_creator_session_id=Nej4MWscg4Thmk2rnrccDmHcBiozpJO7tQu7; galaxy.creator.beaker.session.id=1742877956312063517434; webBuild=4.61.1; unread={%22ub%22:%2267e93aaa000000001d0222c0%22%2C%22ue%22:%2267e7b8ad0000000007035832%22%2C%22uc%22:25}; acw_tc=0a4a063b17434253798992085e3c7958056e9d882a0809453dc928876ed75d; xsecappid=xhs-pc-web; loadts=1743426639546; websectiga=7750c37de43b7be9de8ed9ff8ea0e576519e8cd2157322eb972ecb429a7735d4;sec_poison_id=28d99d4a-12ad-4dc2-b533-347d2d361068"
	//
	//xs = `XYW_eyJzaWduU3ZuIjoiNTYiLCJzaWduVHlwZSI6IngyIiwiYXBwSWQiOiJ4aHMtcGMtd2ViIiwic2lnblZlcnNpb24iOiIxIiwicGF5bG9hZCI6ImNhYTg3NjIyOTk3NDQ0OTM0ZWExYjE0MGM3MDQ1MjZiYjRkZmJiNzIxYmMxZDJlMTM0OGE3NmY2YzgxMjk3OWMzZWQxZGYyNTQwY2RkNGRmYTJkYTRiMjM5ODQwMzJiY2Y0YTg0NTUzNThmZmFkNDQ2OTE3ODhhZGM3ZTYyZTdiMmYzN2ZkYzNlODA5ZDk0M2Y5NGQzYmEzNWM2ZDczMTgyMDU5YjcxNjJkZTRiODFjNzIwMjY2ZDIzYTMyZjUwZGQyM2I0ZDkyYzFlZTE0MWMxNTIyZTgzNGNjODc2ZWNiMmZjYWY4MmEyYzViMmU3MGJjZDJjMmY5ZmE2NzBiM2FkMDU0YjJhZWY3YmJkNGJkZDAxNmM2MjdhNGQyNDAzMjQ0ZTYxNzI3YjM2OTRlZTIwYzE1MjUxZjZkY2RlYjE3MThmMjdmNDJmMjRkZDdmMTVmMjgyODcxZjk0NzkxMTQwNDljNWVlM2NiOWE2MzNiYTc4ZDk5ZDFlMDk5NDY2MWU4YWI4M2ZkMzk4Mzc4MGQ4YTJjMmU1MmRjZjAyNGEyIn0=`
	//xscommon = `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`
	//cookie = `abRequestId=ffa7ca92-1757-5b56-8a73-5570627c1d52; a1=194e40ecd54z2gvo7ea38d9i1fnhxdr5sqnssy5ha30000336056; webId=7486224fc1ddf358b401f55d4a34195f; gid=yj4d48ddiKMDyj4d48dSfAv224uJkhlWd0qYIx0fj3yiIMq8xCfF2M888qqK82K8JJDSijYi; x-user-id-creator.xiaohongshu.com=5cdf90a10000000011023c23; customerClientId=295944940634645; customer-sso-sid=68c5174909416462787924056daksn2qia25hxfn; access-token-creator.xiaohongshu.com=customer.creator.AT-68c5174909416462787304625ty0jxzgyey8ntzs; galaxy_creator_session_id=krTT90MN2muz5EEgc3oDz1YPc41QrZk1q7BI; galaxy.creator.beaker.session.id=1744120765156098347520; xsecappid=xhs-pc-web; unread={%22ub%22:%2267e3db670000000006029421%22%2C%22ue%22:%2267efcaa0000000001202d2dd%22%2C%22uc%22:29}; acw_tc=0ad595c917443420085662838e3cc2382fc297ce77b1b9fcae8e47292988ee; webBuild=4.62.3; websectiga=f3d8eaee8a8c63016320d94a1bd00562d516a5417bc43a032a80cbf70f07d5c0; sec_poison_id=38e2575b-d2a5-4f7e-99d8-4f685593741e; web_session=0400698dabd3e9146c7504d9cd354b823c36d4; loadts=1744343447679`
	//
	//xmns = "unload"
	//
	//xmns = `awpwGvvuMQ4OyuPna5LPbJ5faG+DkeyMyE8j6FJvGaW1ughPzfO3du0/6n7eGOkfK6FtKu4YzKkW0bedEpyydP99bOy7HTa5lG1wjdWO0QeSeZYvDxMn/QebS/We6lCCopmiePFuktmY3v32L0aufv4OILGEEoBNRTweF/oQOSbXEehMyu2uhKlXpXMenR+htn8FNYTO048MB9wevO2+x3leLKFTmBdbcbOfymm2zlLhldKB02bcjpCQj/Jnb1olxTk94Tm8809TBOELf/hQxo6RE0FRo2I7vit+xb5ZDJeL5jSfekzbX2gM1jv21xpYFdfFvcSGRbaeHMCgpkN+0lE4w5MxuYwKFJtoTJhB6zNgBKB164Yj4bpyJG37gQTmN5wMt6cIMNpPK5oxFCw5gNld+xSFbkYPLHlb2F+h8IY2otv6apBXeOda+HoRod8Cdg8NT9hLSikGDw2cdFiSQyw9KwwjjXmFDH+m44EbFMIn1PYtPFk3G01eNFuB`
	//xs = `XYW_eyJzaWduU3ZuIjoiNTYiLCJzaWduVHlwZSI6IngyIiwiYXBwSWQiOiJ4aHMtcGMtd2ViIiwic2lnblZlcnNpb24iOiIxIiwicGF5bG9hZCI6ImJkMzZjZjEwMjY5Zjc1YjA1M2VlOGNmNzdkOTcyMzZiM2JkNzhkY2YzMWY2ODBiZWQyNzA5ZGQyZTM4YzI5NmZhMmNjMDM0OWM4NTg1MDU1MDEzYTRiYTkxMmZjZjU1MWY0MDI3Njc1OTUxNDI4MWI4NGFkNmZiYzgxNzhhZDk0NTI5NDQzMzRkYTFjNGZlMjU3ZjdkMjg2M2Q3NjczMzljOTBiZTY5N2QxOWM2ZWEzZmY1MjYwMThjMzc3YzMwYTc0ZGQ0ZTQ1MjM3YTk5MjEzMmYyZTgwMjA0MDdjY2Q4NzQyOTMxYWEwMGI2YzE4ZmI2ZmU0MjBmNjRiY2IzYTYyOGU3MmVmMjg0YTY3MjVlNGRjYjVhZDk0M2MwYmUwNjBlODYzMmY2ZDY2YTJlMzQ3NWI1ZWQyMWZjNzlkN2NjZDY2ZWJhYzJhNGFkNzA3NmI0YTk2OWRiODU4N2E2MWVhNzY0ZDI5MWE5N2NhYzA0ZmYwZTBhNTQ1NjI3YzYzMTkwNjAzYjAyNzFhZjc5MTgxNDc0Y2E4YmU3OTE3MWIzIn0=`
	//xscommon = `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`
	//cookie = `abRequestId=ffa7ca92-1757-5b56-8a73-5570627c1d52; a1=194e40ecd54z2gvo7ea38d9i1fnhxdr5sqnssy5ha30000336056; webId=7486224fc1ddf358b401f55d4a34195f; gid=yj4d48ddiKMDyj4d48dSfAv224uJkhlWd0qYIx0fj3yiIMq8xCfF2M888qqK82K8JJDSijYi; x-user-id-creator.xiaohongshu.com=5cdf90a10000000011023c23; customerClientId=295944940634645; access-token-creator.xiaohongshu.com=customer.creator.AT-68c5174909416462787304625ty0jxzgyey8ntzs; galaxy_creator_session_id=krTT90MN2muz5EEgc3oDz1YPc41QrZk1q7BI; galaxy.creator.beaker.session.id=1744120765156098347520; webBuild=4.62.3; web_session=0400698dabd3e9146c75752d333a4b48bb6c4a; unread={%22ub%22:%2267ece601000000001d01fcb5%22%2C%22ue%22:%226804eadb000000001201dcc1%22%2C%22uc%22:28}; xsecappid=xhs-pc-web; loadts=1745426949355; acw_tc=0a50892f17454578549161900e319560e854e298bf53914c3b5ec76cf77bfa; websectiga=10f9a40ba454a07755a08f27ef8194c53637eba4551cf9751c009d9afb564467; sec_poison_id=96fefc3d-bf0f-4b3b-93e5-f9833d9348c4`

	return &Client{
		//headers: map[string]string{
		//	"x-mns":      xmns,
		//	"x-s":        xs,
		//	"x-s-common": xscommon,
		//	"cookie":     cookie,
		//	"x-t":        conv.Str(time.Now().UnixMilli()),
		//	"origin":     "https://www.xiaohongshu.com",
		//},
	}
}

func (t *Client) SetAuth(headers map[string]string) *Client {
	t.headers = map[string]string{
		"x-mns":      helper.OrString(headers["xmns"], "unload"),
		"x-s-common": headers["xscommon"],
		"cookie":     headers["cookie"],
		"x-s":        headers["xs"],
		"x-t":        conv.Str(time.Now().UnixMilli()),
		"origin":     "https://www.xiaohongshu.com",
	}

	return t
}

type Notes []*Note

type Filter struct {
	Tags []string `json:"tags"`
	Type string   `json:"type"`
}
type ListNotesParams struct {
	Keyword      string        `json:"keyword"`
	Page         int           `json:"page"`
	PageSize     int           `json:"page_size"`
	SearchId     string        `json:"search_id"`
	Sort         string        `json:"sort"`
	NoteType     int           `json:"note_type"`
	ExtFlags     []interface{} `json:"ext_flags"`
	Filters      []Filter      `json:"filters"`
	Geo          string        `json:"geo"`
	ImageFormats []string      `json:"image_formats"`
}

type Note struct {
	XsecToken string `json:"xsec_token"`
	Id        string `json:"id"`
	ModelType string `json:"model_type"`
	NoteCard  struct {
		User struct {
			Avatar    string `json:"avatar"`
			UserId    string `json:"user_id"`
			Nickname  string `json:"nickname"`
			XsecToken string `json:"xsec_token"`
			NickName  string `json:"nick_name"`
		} `json:"user"`
		InteractInfo struct {
			CollectedCount string `json:"collected_count"`
			CommentCount   string `json:"comment_count"`
			SharedCount    string `json:"shared_count"`
			Liked          bool   `json:"liked"`
			LikedCount     string `json:"liked_count"`
			Collected      bool   `json:"collected"`
		} `json:"interact_info"`
		Cover struct {
			Width      int    `json:"width"`
			UrlDefault string `json:"url_default"`
			UrlPre     string `json:"url_pre"`
			Height     int    `json:"height"`
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
		Type         string `json:"type"`
		DisplayTitle string `json:"display_title,omitempty"`
	} `json:"note_card,omitempty"`
	RecQuery struct {
		Queries []struct {
			Id         string `json:"id"`
			Name       string `json:"name"`
			SearchWord string `json:"search_word"`
		} `json:"queries"`
		Title         string `json:"title"`
		Source        int    `json:"source"`
		WordRequestId string `json:"word_request_id"`
	} `json:"rec_query,omitempty"`
}

type Response struct {
	Data struct {
		HasMore bool  `json:"has_more"`
		Items   Notes `json:"items"`
	} `json:"data"`
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

func (t *Client) ListNotes(ctx context.Context, params ListNotesParams) (Notes, bool, error) {

	//body := map[string]any{
	//	"keyword": "热点",
	//	"page": 3,
	//	"page_size": 20,
	//	"search_id": "2epl6zbwf8mfl5tx1p4pt@2epl71sktwb9rmle9yplk",
	//	"sort": "popularity_descending",
	//	"note_type": 1,
	//	"ext_flags":[]string{},
	//	"filters": [{"tags":["general"], "type":"sort_type"},{"tags":["不限"], "type":"filter_note_type"},{"tags":["不限"], "type":"filter_note_time"},{"tags":["不限"], "type":"filter_note_range"},{"tags":["不限"], "type":"filter_pos_distance"}], "geo":"", "image_formats":["jpg", "webp", "avif"]}

	post, err := resty.New().R().SetContext(ctx).
		SetHeaders(t.headers).
		SetBody(params).
		Post("https://edith.xiaohongshu.com/api/sns/web/v1/search/notes")
	if err != nil {
		return nil, false, err
	}

	var res Response
	if err := json.Unmarshal(post.Body(), &res); err != nil {
		return nil, false, err
	}

	if !res.Success {
		return nil, false, errors.New(res.Msg)
	}

	return res.Data.Items, res.Data.HasMore, nil
}

type NoteMetadata struct {
	VideoUrl   string
	ProfileUrl string
	Title      string
	Desc       string
}

// "https://www.xiaohongshu.com/search_result/%s?xsec_token=%s&xsec_source=pc_search", x.Raw["id"], x.Raw["xsecToken"]
func (t *Client) GetNoteMetadataByUrl(ctx context.Context, url string) (*NoteMetadata, error) {

	htmlDoc, err := t.GetHtmlDoc(ctx, url)
	if err != nil {
		return nil, err
	}

	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlDoc))
	if err != nil {
		return nil, err
	}

	videoUrl := goqueryz.FindMetaContent(reader, "og:video")
	//if videoUrl == "" {
	//	return nil, errors.New("videoNotFound")
	//}

	var profileUrl string
	var title string
	var desc string
	reader.Find(".interaction-container").Each(func(i int, s *goquery.Selection) {
		s.Find(".author-wrapper .info a").Each(func(i int, s *goquery.Selection) {
			profileUrl = s.AttrOr("href", "")
		})

		s.Find(".note-scroller").Each(func(i int, s *goquery.Selection) {
			s.Find(".note-content").Each(func(i int, s *goquery.Selection) {
				s.Find("div[id=detail-title]").Each(func(i int, s *goquery.Selection) {
					title = s.Text()
				})
				s.Find("div[id=detail-desc]").Each(func(i int, s *goquery.Selection) {
					desc = s.Text()
				})
			})
		})

	})

	if profileUrl == "" {
		return nil, errors.New("profileUrlNotFound")
	}

	return &NoteMetadata{
		VideoUrl:   videoUrl,
		ProfileUrl: "https://www.xiaohongshu.com" + profileUrl,
		Title:      title,
		Desc:       desc,
	}, nil

}

func (t *Client) GetHtmlDoc(ctx context.Context, url string) ([]byte, error) {

	rsp, err := resty.New().R().
		SetHeaders(t.headers).
		SetContext(ctx).
		//Get("https://edith.xiaohongshu.com/api/sns/web/v1/search/recommend?keyword=675330826")
		Get(url)
	if err != nil {
		return nil, err
	}

	return rsp.Body(), nil
}

func (t *Client) GetHtmlDocV2(ctx context.Context, url string) ([]byte, string, error) {

	rsp, err := resty.New().R().
		SetHeaders(t.headers).
		SetContext(ctx).
		//Get("https://edith.xiaohongshu.com/api/sns/web/v1/search/recommend?keyword=675330826")
		Get(url)
	if err != nil {
		return nil, "", err
	}

	if strings.Contains(rsp.RawResponse.Request.URL.Path, "/login") {
		return nil, "", errors.New("need to login")
	}

	// 从url中获取 笔记id，调用tikhub获取内容  作为信息补充
	noteId := strings.ReplaceAll(rsp.RawResponse.Request.URL.Path, "/discovery/item/", "")

	if noteId == "" {
		//rsp.RawResponse.Request.URL.RawQuery = ""
		urlObj := &netUrl.URL{}

		// 设置查询字符串
		urlObj.RawQuery = rsp.RawResponse.Request.URL.RawQuery
		// 解析查询参数
		queryParams := urlObj.Query()
		// 获取 noteId
		noteId = queryParams.Get("noteId")

	}

	return rsp.Body(), noteId, nil
}

type Profile struct {
	Avatar         string   `json:"avatar"`
	Username       string   `json:"username"`
	Id             string   `json:"id"`
	IpAddress      string   `json:"ipAddress"`
	Sign           string   `json:"sign"`
	Tags           []string `json:"tags"`
	FollowingCount string   `json:"followingCount"`
	FollowerCount  string   `json:"followerCount"`
	LikedCount     string   `json:"likedCount"`
	NoteCount      string   `json:"noteCount"`
}

func (t *Client) GetProfileByLink(ctx context.Context, link string) (*Profile, error) {

	htmlDoc, err := t.GetHtmlDoc(ctx, link)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlDoc))
	if err != nil {
		return nil, err
	}

	var avatar string
	var username string
	var id string
	var ipAddress string
	var sign string
	var tags []string
	var followingCount string
	var followerCount string
	var likedCount string

	doc.Find(".info-part .info .basic-info").Each(func(i int, x *goquery.Selection) {
		x.Find(".avatar .avatar-wrapper .user-image").Each(func(i int, x *goquery.Selection) {
			avatar, _ = x.Attr("src")
		})

		x.Find(".user-basic .user-nickname .user-name").Each(func(i int, x *goquery.Selection) {
			username = x.Text()
		})

		x.Find(".user-basic .user-content .user-redId").Each(func(i int, x *goquery.Selection) {

			parts := strings.Split(x.Text(), "：")
			if len(parts) == 2 {
				id = parts[1]
			}
		})
		x.Find(".user-basic .user-content .user-IP").Each(func(i int, x *goquery.Selection) {
			ipAddress = x.Text()
		})
	})

	doc.Find(".info-part .info .user-desc").Each(func(i int, x *goquery.Selection) {
		sign = x.Text()
	})
	doc.Find(".info-part .info .user-tags .tag-item").Each(func(i int, x *goquery.Selection) {
		if x.Text() != "" {
			tags = append(tags, x.Text())
		}
	})
	doc.Find(".info-part .info .data-info .user-interactions .count").
		Each(func(i int, x *goquery.Selection) {
			if i == 0 {
				followingCount = x.Text()
			}

			if i == 1 {
				followerCount = x.Text()
			}
			if i == 2 {
				likedCount = x.Text()
			}
		})

	noteCount, err := t.getNoteCount(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Profile{
		Avatar:         avatar,
		Username:       username,
		Id:             id,
		IpAddress:      ipAddress,
		Sign:           sign,
		Tags:           tags,
		FollowingCount: followingCount,
		FollowerCount:  followerCount,
		LikedCount:     likedCount,
		NoteCount:      noteCount,
	}, nil
}

func (t *Client) getNoteCount(ctx context.Context, id string) (string, error) {

	rsp2, err := resty.New().R().
		SetHeaders(t.headers).
		SetContext(ctx).
		Get("https://edith.xiaohongshu.com/api/sns/web/v1/search/recommend?keyword=" + id)
	if err != nil {
		return "", err
	}

	var result Result

	err = json.Unmarshal(rsp2.Body(), &result)
	if err != nil {
		return "", err
	}

	var noteCount string
	if len(result.Data.SugItems) > 0 {
		noteCount = conv.Str(result.Data.SugItems[0].User.NoteCount)
	}

	return noteCount, nil
}

type Result struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    struct {
		WordRequestId string `json:"word_request_id"`
		SugItems      []struct {
			User struct {
				UpdateTime                string `json:"update_time"`
				Id                        string `json:"id"`
				Followed                  bool   `json:"followed"`
				RedOfficialVerifyType     int    `json:"red_official_verify_type"`
				TrackDuration             int    `json:"track_duration"`
				Name                      string `json:"name"`
				Image                     string `json:"image"`
				Desc                      string `json:"desc"`
				Fans                      string `json:"fans"`
				RedId                     string `json:"red_id"`
				ShowRedOfficialVerifyIcon bool   `json:"show_red_official_verify_icon"`
				IsSelf                    bool   `json:"is_self"`
				RedOfficialVerified       bool   `json:"red_official_verified"`
				NoteCount                 int    `json:"note_count"`
			} `json:"user,omitempty"`
			Text           string `json:"text"`
			SearchType     string `json:"search_type"`
			Type           string `json:"type"`
			JumpType       string `json:"jump_type,omitempty"`
			HighlightFlags []bool `json:"highlight_flags,omitempty"`
		} `json:"sug_items"`
		SearchCplId string `json:"search_cpl_id"`
	} `json:"data"`
}
