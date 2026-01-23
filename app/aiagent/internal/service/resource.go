package service

import (
	"bytes"
	"context"
	errors2 "errors"
	"fmt"
	"io"
	trackerpb "store/api/aiagent"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/goqueryz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/filed"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/tikhub"
	"store/pkg/sdk/urlz"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (t TrackerService) GetAliOssSignedUrl(ctx context.Context, params *trackerpb.GetAliOssSignedUrlParams) (*trackerpb.GetAliOssSignedUrlResult, error) {

	//url, err := t.Data.Alioss.GetSignedUrl(ctx, params.Bucket, params.FileKey, params.ContentType)
	//if err != nil {
	//	return nil, err
	//}

	return &trackerpb.GetAliOssSignedUrlResult{
		//Url: url,
	}, nil
}

func (t TrackerService) GetQiniuUploadToken(ctx context.Context, params *trackerpb.GetQiniuUploadTokenParams) (*trackerpb.GetQiniuUploadTokenResult, error) {

	//url, err := t.Data.Qiniu.GetUploadToken(ctx, params.Bucket)
	//if err != nil {
	//	return nil, err
	//}

	return &trackerpb.GetQiniuUploadTokenResult{
		//Token: url,
	}, nil
}

func (t TrackerService) addTiktokResource(ctx context.Context, url string) (*trackerpb.ResourceV2, error) {

	video, err := t.Data.Tikhub.TiktokGetVideoByShareUrl(ctx, url)
	if err != nil {
		return nil, err
	}

	if len(video.AwemeDetails) == 0 {
		return nil, errors.New(0, "noVideoFound", "")
	}

	x := video.AwemeDetails[0]

	if len(x.Video.DownloadAddr.UrlList) == 0 {
		return nil, errors.New(0, "no url found", "")
	}

	videoUrl := x.Video.DownloadAddr.UrlList[0]

	authorAvatar := x.Author.Avatar300X300.UrlList[0]

	return &trackerpb.ResourceV2{
		Url:      "",
		Category: "",
		//GenaiUri: "",
		Name: "",
		Profile: &trackerpb.Profile{
			XId: "",
			//LastUpdatedAt:  0,
			Content:        nil,
			Avatar:         authorAvatar,
			Username:       x.Author.Nickname,
			PlatformId:     x.Author.InsId,
			IpAddress:      "",
			Sign:           x.Author.Signature,
			Tags:           nil,
			FollowingCount: "",
			FollowerCount:  "",
			LikedCount:     "",
			NoteCount:      conv.Str(x.Author.AwemeCount),
			//Followers:      0,
			Platform: "tiktok",
			Raw:      "",
			UserId:   "",
		},
		PlatformUrl: videoUrl,
		//Title:       x,
		Desc: x.Desc,
		//UploadUserId: userId,
		CreatedAt: time.Now().Unix(),
		MimeType:  "video/mp4",
	}, nil

	//return nil, nil
}

func (t TrackerService) addDouyinResource(ctx context.Context, url string) (*trackerpb.ResourceV2, error) {

	video, err := t.Data.Tikhub.GetVideoByShareUrl(ctx, url)
	if err != nil {
		return nil, err
	}

	if len(video.AwemeDetail.Video.DownloadAddr.UrlList) == 0 {
		return nil, errors.New(0, "no url found", "")
	}

	videoUrl := video.AwemeDetail.Video.DownloadAddr.UrlList[0]

	authorAvatar := video.AwemeDetail.Author.Avatar300X300.UrlList[0]

	return &trackerpb.ResourceV2{
		Url:      "",
		Category: "",
		//GenaiUri: "",
		Name: "",
		Profile: &trackerpb.Profile{
			XId: "",
			//LastUpdatedAt:  0,
			Content:        nil,
			Avatar:         authorAvatar,
			Username:       video.AwemeDetail.Author.Nickname,
			PlatformId:     video.AwemeDetail.Author.InsId,
			IpAddress:      "",
			Sign:           video.AwemeDetail.Author.Signature,
			Tags:           nil,
			FollowingCount: "",
			FollowerCount:  "",
			LikedCount:     "",
			NoteCount:      conv.Str(video.AwemeDetail.Author.AwemeCount),
			//Followers:      0,
			Platform: "douyin",
			Raw:      "",
			UserId:   "",
		},
		PlatformUrl: videoUrl,
		Title:       video.AwemeDetail.PreviewTitle,
		Desc:        video.AwemeDetail.Desc,
		//UploadUserId: userId,
		CreatedAt: time.Now().Unix(),
		MimeType:  "video/mp4",
	}, nil
}

func (t TrackerService) addXhsResource(ctx context.Context, sessionId, userId, url string) (*trackerpb.ResourceV2, error) {

	settings, err := t.Data.Mongo.Settings.Get(ctx)
	if err != nil {
		return nil, err
	}

	client := t.Data.XhsClient.SetAuth(settings.XhsCookies)

	doc, err := client.GetHtmlDoc(ctx, url)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(doc))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//html, _ := reader.Html()
	//log.Debugw("NewDocumentFromReader response", html)

	videoUrl := goqueryz.FindMetaContent(reader, "og:video")
	if videoUrl == "" {
		return nil, errors.New(10400, "videoNotFound", "")
	}

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
		log.Error("Parse HtmlDoc err", "", "profileUrl", profileUrl, "videoUrl", videoUrl)
		return nil, errors.New(10500, "", "profile url not found")
	}

	profile, err := client.GetProfileByLink(ctx, "https://www.xiaohongshu.com"+profileUrl)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	resourceId := primitive.NewObjectID().Hex()

	newResource := &trackerpb.ResourceV2{
		XId:         resourceId,
		SessionId:   sessionId,
		PlatformUrl: videoUrl,
		Title:       title,
		Desc:        desc,
		MimeType:    "video/mp4",
		Profile: &trackerpb.Profile{
			Avatar:         profile.Avatar,
			Username:       profile.Username,
			PlatformId:     profile.Id,
			IpAddress:      profile.IpAddress,
			Sign:           profile.Sign,
			Tags:           profile.Tags,
			FollowingCount: profile.FollowingCount,
			FollowerCount:  profile.FollowerCount,
			LikedCount:     profile.LikedCount,
			NoteCount:      profile.NoteCount,
		},
		UploadUserId: userId,
	}

	_, err = t.Data.Mongo.Resource.InsertOne(ctx, newResource)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	ctx = context.Background()

	//geminiUrl, ossUrl, err := t.upload(ctx, t.Data.GenaiClient, videoUrl)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	ossUrl, err := t.uploadToS3(ctx, videoUrl, "video/mp4")
	if err != nil {
		return nil, err
	}

	_, err = t.Data.Mongo.Resource.UpdateFieldsById(ctx, resourceId,
		bson.M{"url": ossUrl}, //"genaiUri": geminiUrl,

	)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//newResource.GenaiUri = geminiUrl
	newResource.Url = ossUrl

	// 冗余到session

	_, err = t.Data.Mongo.Session.UpdateOne(ctx,
		bson.M{"_id": sessionId},
		bson.M{
			"$push": bson.M{"resources": newResource},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	return newResource, nil
}

func (t TrackerService) convToShareUrl(url string) string {

	//https://www.xiaohongshu.com/explore/68b6fc54000000001d00bc4d?xsec_token=ABTLXAmhdRDxIY5R5_fzU31JPCeZRzy5xqOGvzvX9DJJ0=&xsec_source=
	//https://www.xiaohongshu.com/discovery/item/68b6fc54000000001d00bc4d?source=webshare&xhsshare=pc_web&xsec_token=ABTLXAmhdRDxIY5R5_fzU31JPCeZRzy5xqOGvzvX9DJJ0=&xsec_source=pc_share
	urlInfo, _ := urlz.ParseURL(url)

	noteId := urlInfo.PathSegments[len(urlInfo.PathSegments)-1]

	xsecToken := urlInfo.QueryParams["xsec_token"]

	return fmt.Sprintf("https://www.xiaohongshu.com/discovery/item/%s?source=webshare&xhsshare=pc_web&xsec_token=%s&xsec_source=pc_share", noteId, xsecToken)

}

func (t TrackerService) addXhsResourceV3(ctx context.Context, url string) ([]*trackerpb.ResourceV2, error) {

	// 短连接(app复制的)
	if strings.HasPrefix(url, "http://xhslink.com") || strings.HasPrefix(url, "https://xhslink.com") || strings.HasPrefix(url, "xhslink.com") {
		noteMeta, err := t.Data.Tikhub.XhsGetNoteByShareUrl(ctx, url)
		if err != nil {
			return nil, err
		}

		//// 手机上复制的分享链接 http://xhslink.com/a/8bW9FkbVb4qgb 需要登录才行，电脑端的可以直接用
		//// 用tikhub转换下 转成电脑端的 https://www.xiaohongshu.com/discovery/item/6847cd61000000002202b553?source=webshare&xhsshare=pc_web&xsec_token=CBxSy-OEn1Kn2lTZtq215_kbkGpi12eHxJK3w4fi6coHg=&xsec_source=pc_share
		//url = fmt.Sprintf("https://www.xiaohongshu.com/discovery/item/%s?source=webshare&xhsshare=pc_web&xsec_token=%s&xsec_source=pc_share", noteMeta.NoteId, noteMeta.XsecToken)
		//
		//note, err := t.addXhsByRawHtml(ctx, url)
		//if err == nil {
		//	return note, err
		//}
		//
		//log.Errorw("addXhsByRawHtml err", err, "url", url)

		return t.addXhsByNoteIdViaTikhub(ctx, noteMeta.NoteId)
	} else {

		if strings.Contains(url, "www.xiaohongshu.com/explore/") {
			url = t.convToShareUrl(url)
		}

		note, err := t.addXhsByRawHtml(ctx, url)
		if err == nil {
			return note, err
		}

		log.Errorw("addXhsByRawHtml err", err, "url", url)

		urlInfo, _ := urlz.ParseURL(url)

		noteId := urlInfo.PathSegments[len(urlInfo.PathSegments)-1]

		note, err = t.addXhsByNoteIdViaTikhub(ctx, noteId)
		if err == nil {
			return note, err
		}
		log.Errorw("addXhsByNoteIdViaTikhub err", err, "noteId", noteId, "url", url)

		return nil, errors2.New("resource not found")
	}
}

func (t TrackerService) addXhsByShareUrlViaTikhub(ctx context.Context, url string) ([]*trackerpb.ResourceV2, error) {

	noteMeta, err := t.Data.Tikhub.XhsGetNoteByShareUrl(ctx, url)
	if err != nil {
		return nil, err
	}

	note, err := t.Data.Tikhub.XhsWebGetNoteInfoById(ctx, noteMeta.NoteId)
	if err == nil {
		return t.pack(note)
	}
	note, err = t.Data.Tikhub.XhsWebGetNoteInfoByIdV2(ctx, noteMeta.NoteId)
	if err == nil {
		return t.pack(note)
	}

	note, err = t.Data.Tikhub.XhsWebGetNoteInfoByIdV4(ctx, noteMeta.NoteId)
	if err == nil {
		return t.pack(note)
	}

	fmt.Println(note)

	return []*trackerpb.ResourceV2{}, nil
}

func (t TrackerService) addXhsByNoteIdViaTikhub(ctx context.Context, noteId string) ([]*trackerpb.ResourceV2, error) {

	//note, err := t.Data.Tikhub.XhsWebGetNoteInfoById(ctx, noteId)
	//if err == nil {
	//	return t.pack(note)
	//}
	//note, err = t.Data.Tikhub.XhsWebGetNoteInfoByIdV2(ctx, noteId)
	//if err == nil {
	//	return t.pack(note)
	//}
	//
	//note, err = t.Data.Tikhub.XhsWebGetNoteInfoByIdV4(ctx, noteId)
	//if err == nil {
	//	return t.pack(note)
	//}

	note, err := t.Data.Tikhub.XhsGetNoteInfoById(ctx, noteId)
	if err != nil {
		return nil, err
	}

	var resources []*trackerpb.ResourceV2

	if note.Title != "" {
		resources = append(resources, &trackerpb.ResourceV2{
			//XId:      primitive.NewObjectID().Hex(),
			MimeType: "text/plain",
			Category: "title",
			Content:  note.Title,
		})
	}

	if note.Content != "" {
		resources = append(resources, &trackerpb.ResourceV2{
			//XId:      primitive.NewObjectID().Hex(),
			MimeType: "text/plain",
			Category: "description",
			Content:  note.Content,
		})
	}

	resources = append(resources, &trackerpb.ResourceV2{
		//XId:      primitive.NewObjectID().Hex(),
		MimeType: "text/plain",
		Category: "interactionMetrics",
		Meta: map[string]string{
			//"keywords":       keywords,
			//"description":    description,
			//"title":          title,
			//"commentCount":   commentCount,
			//"likedCount":     likedCount,
			"sharedCount": conv.Str(note.ShareCount),
			//"collectedCount": collectedCount,
		},
		Content: conv.M2J(map[string]interface{}{
			//"keywords":       keywords,
			//"description":    description,
			//"title":          title,
			//"commentCount":   commentCount,
			//"likedCount":     likedCount,
			//"collectedCount": collectedCount,
			"sharedCount": conv.Str(note.ShareCount),
		}),
	})

	if note.Url != "" {

		var coverUrl string
		if len(note.Images) > 0 {
			coverUrl = note.Images[0].Url
		}

		resources = append(resources, &trackerpb.ResourceV2{
			//XId:         primitive.NewObjectID().Hex(),
			PlatformUrl: note.Url,
			CoverUrl:    coverUrl,
			MimeType:    "video/mp4",
		})
	}

	if len(note.Images) > 0 {
		for _, x := range note.Images {
			resources = append(resources, &trackerpb.ResourceV2{
				PlatformUrl: x.Url,
				MimeType:    "image/jpeg",
			})
		}

	}

	profile_ := &trackerpb.Profile{
		Avatar:     note.User.Avatar,
		Username:   note.User.Name,
		PlatformId: note.User.Id,
		//IpAddress:  profile.IpAddress,
		Sign: note.User.Desc,
		//Tags: note.User.Tags,
		//FollowingCount: profile.FollowingCount,
		//FollowerCount:  profile.FollowerCount,
		//LikedCount:     profile.LikedCount,
		//NoteCount:      profile.NoteCount,
		Platform: "xiaohongshu",
	}

	resources = append(resources, &trackerpb.ResourceV2{
		//XId:      primitive.NewObjectID().Hex(),
		MimeType: "text/plain",
		Category: "authorProfile",
		Content:  conv.S2J(profile_),
		Meta: map[string]string{
			"avatar":     note.User.Avatar,
			"username":   note.User.Name,
			"platformId": note.User.Id,
			//"ipAddress":  profile.IpAddress,
			"sign": note.User.Desc,
			//"tags":       strings.Join(note.Tags, ","),
			//"followingCount": profile.FollowingCount,
			//"followerCount":  profile.FollowerCount,
			//"likedCount":     profile.LikedCount,
			//"noteCount":      profile.NoteCount,
			"platform": "xiaohongshu",
		},
	})

	return resources, nil
}

func (t TrackerService) pack(note *tikhub.Note) ([]*trackerpb.ResourceV2, error) {
	return nil, nil
}

func (t TrackerService) addXhsByRawHtml(ctx context.Context, url string) ([]*trackerpb.ResourceV2, error) {

	settings, err := t.Data.Mongo.Settings.Get(ctx)
	if err != nil {
		return nil, err
	}

	client := t.Data.XhsClient.SetAuth(settings.XhsCookies)

	//if !strings.HasPrefix(url, "https") {
	//	url = strings.Replace(url, "http", "https", 1)
	//}

	doc, noteId, err := client.GetHtmlDocV2(ctx, url)
	if err != nil {
		log.Errorw("GetHtmlDocV2", err, "url", url)
		return nil, err
	}

	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(doc))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var resources []*trackerpb.ResourceV2

	keywords := goqueryz.FindMetaContent(reader, "keywords")
	description := goqueryz.FindMetaContent(reader, "description")
	//commentCount := goqueryz.FindMetaContent(reader, "og:xhs:note_comment")
	//likedCount := goqueryz.FindMetaContent(reader, "og:xhs:note_like")
	//sharedCount := ""
	//collectedCount := goqueryz.FindMetaContent(reader, "og:xhs:note_collect")

	var commentCount string
	var likedCount string
	var sharedCount string
	var collectedCount string

	title := goqueryz.FindMetaContent(reader, "og:title")
	title = strings.Split(title, "-")[0]

	// 从url中获取 笔记id，调用tikhub获取内容  作为信息补充
	//noteId := strings.ReplaceAll(path, "/discovery/item/", "")

	var images []string

	note, err := t.Data.Tikhub.XhsWebGetNoteInfoByIdV4(ctx, noteId)
	if err != nil {
		log.Error(err)
		//return nil, err
	}

	if note != nil && len(note.Data) > 0 && len(note.Data[0].NoteList) > 0 {
		x := note.Data[0].NoteList[0]

		commentCount = conv.Str(x.CommentsCount)
		likedCount = conv.Str(x.LikedCount)
		sharedCount = conv.Str(x.SharedCount)
		collectedCount = conv.Str(x.CollectedCount)

		description = x.Desc

		var keywordsList []string
		for _, xx := range x.HashTag {
			keywordsList = append(keywordsList, xx.Name)
		}
		if len(keywordsList) > 0 {
			keywords = strings.Join(keywordsList, ",")
		}

		title = x.Title

		for _, xx := range x.ImagesList {
			images = append(images, xx.Original)

			//resources = append(resources, &trackerpb.ResourceV2{
			//	PlatformUrl: x.Url,
			//	MimeType:    "image/jpeg",
			//})
		}
	}

	// 元数据

	if title != "" {
		resources = append(resources, &trackerpb.ResourceV2{
			//XId:      primitive.NewObjectID().Hex(),
			MimeType: "text/plain",
			Category: "title",
			Content:  title,
		})
	}

	if description != "" {
		resources = append(resources, &trackerpb.ResourceV2{
			//XId:      primitive.NewObjectID().Hex(),
			MimeType: "text/plain",
			Category: "description",
			Content:  description,
		})
	}

	resources = append(resources, &trackerpb.ResourceV2{
		//XId:      primitive.NewObjectID().Hex(),
		MimeType: "text/plain",
		Category: "interactionMetrics",
		Meta: map[string]string{
			"keywords":       keywords,
			"description":    description,
			"title":          title,
			"commentCount":   commentCount,
			"likedCount":     likedCount,
			"sharedCount":    sharedCount,
			"collectedCount": collectedCount,
		},
		Content: conv.M2J(map[string]interface{}{
			"keywords":       keywords,
			"description":    description,
			"title":          title,
			"commentCount":   commentCount,
			"likedCount":     likedCount,
			"collectedCount": collectedCount,
		}),
	})

	// 资源
	ogType := goqueryz.FindMetaContent(reader, "og:type")
	if ogType == "video" {
		videoUrl := goqueryz.FindMetaContent(reader, "og:video")
		if videoUrl == "" {
			return nil, errors.New(10400, "videoNotFound", "")
		}
		resources = append(resources, &trackerpb.ResourceV2{
			//XId:         primitive.NewObjectID().Hex(),
			PlatformUrl: videoUrl,
			CoverUrl:    goqueryz.FindMetaContent(reader, "og:image"),
			MimeType:    "video/mp4",
		})

	} else {
		// og:type article

		if len(images) == 0 {
			images = goqueryz.FindMetaContents(reader, "og:image")
			//resources = append(resources, helper.Mapping(images, func(x string) *trackerpb.ResourceV2 {
			//
			//	if !strings.HasPrefix(x, "http") {
			//
			//	}
			//
			//	return &trackerpb.ResourceV2{
			//		//XId:         primitive.NewObjectID().Hex(),
			//		PlatformUrl: x,
			//		MimeType:    "image/jpeg",
			//	}
			//})...)
		}

		for _, x := range images {

			if !strings.HasPrefix(x, "http") {
				continue
			}
			resources = append(resources, &trackerpb.ResourceV2{
				PlatformUrl: x,
				MimeType:    "image/jpeg",
			})
		}

	}

	// 作者信息
	var profileUrl string
	reader.Find(".interaction-container").Each(func(i int, s *goquery.Selection) {
		s.Find(".author-wrapper .info a").Each(func(i int, s *goquery.Selection) {
			profileUrl = s.AttrOr("href", "")
		})
	})

	if profileUrl != "" {
		//log.Error("Parse HtmlDoc err", "", "profileUrl", profileUrl)
		//return nil, errors.New(10500, "", "profile url not found")
		profile, err := client.GetProfileByLink(ctx, "https://www.xiaohongshu.com"+profileUrl)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		profile_ := &trackerpb.Profile{
			Avatar:     profile.Avatar,
			Username:   profile.Username,
			PlatformId: profile.Id,
			IpAddress:  profile.IpAddress,
			Sign:       profile.Sign,
			Tags:       profile.Tags,
			//FollowingCount: profile.FollowingCount,
			//FollowerCount:  profile.FollowerCount,
			//LikedCount:     profile.LikedCount,
			//NoteCount:      profile.NoteCount,
			Platform: "xiaohongshu",
		}

		resources = append(resources, &trackerpb.ResourceV2{
			//XId:      primitive.NewObjectID().Hex(),
			MimeType: "text/plain",
			Category: "authorProfile",
			Content:  conv.S2J(profile_),
			Meta: map[string]string{
				"avatar":     profile.Avatar,
				"username":   profile.Username,
				"platformId": profile.Id,
				"ipAddress":  profile.IpAddress,
				"sign":       profile.Sign,
				"tags":       strings.Join(profile.Tags, ","),
				//"followingCount": profile.FollowingCount,
				//"followerCount":  profile.FollowerCount,
				//"likedCount":     profile.LikedCount,
				//"noteCount":      profile.NoteCount,
				"platform": "xiaohongshu",
			},
		})
	}

	return resources, nil

}

func (t TrackerService) addXhsResourceV2(ctx context.Context, url string) (*trackerpb.ResourceV2, error) {

	settings, err := t.Data.Mongo.Settings.Get(ctx)
	if err != nil {
		return nil, err
	}

	client := t.Data.XhsClient.SetAuth(settings.XhsCookies)

	doc, err := client.GetHtmlDoc(ctx, url)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(doc))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//html, _ := reader.Html()
	//log.Debugw("NewDocumentFromReader response", html)

	videoUrl := goqueryz.FindMetaContent(reader, "og:video")
	if videoUrl == "" {
		return nil, errors.New(10400, "videoNotFound", "")
	}

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
		log.Error("Parse HtmlDoc err", "", "profileUrl", profileUrl, "videoUrl", videoUrl)
		return nil, errors.New(10500, "", "profile url not found")
	}

	profile, err := client.GetProfileByLink(ctx, "https://www.xiaohongshu.com"+profileUrl)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	resourceId := primitive.NewObjectID().Hex()

	newResource := &trackerpb.ResourceV2{
		XId: resourceId,
		//SessionId:   sessionId,
		PlatformUrl: videoUrl,
		Title:       title,
		Desc:        desc,
		MimeType:    "video/mp4",
		Profile: &trackerpb.Profile{
			Avatar:         profile.Avatar,
			Username:       profile.Username,
			PlatformId:     profile.Id,
			IpAddress:      profile.IpAddress,
			Sign:           profile.Sign,
			Tags:           profile.Tags,
			FollowingCount: profile.FollowingCount,
			FollowerCount:  profile.FollowerCount,
			LikedCount:     profile.LikedCount,
			NoteCount:      profile.NoteCount,
		},
		//UploadUserId: userId,
	}

	return newResource, nil
}

func (t TrackerService) AddResourceByLink(ctx context.Context, params *trackerpb.AddResourceByLinkParams) (*trackerpb.ResourceV2, error) {
	userId := krathelper.RequireUserId(ctx)

	var tmpResources []*trackerpb.ResourceV2
	var tmpResource *trackerpb.ResourceV2
	var err error

	if strings.Contains(params.Link, "v.douyin.com") {
		tmpResource, err = t.addDouyinResource(ctx, params.Link)
	} else if strings.Contains(params.Link, "tiktok.com") {
		tmpResource, err = t.addTiktokResource(ctx, params.Link)

	} else {
		tmpResources, err = t.addXhsResourceV3(ctx, params.Link)
	}

	if err != nil {
		return nil, err
	}

	if tmpResource != nil {
		tmpResources = append(tmpResources, tmpResource)
	}

	for i := range tmpResources {

		x := tmpResources[i]

		if x.IsVideo() || x.IsImage() {

			ossUrl, err := t.uploadToS3(ctx, x.PlatformUrl, x.MimeType)
			if err != nil {
				return nil, err
			}
			x.Url = ossUrl
		}

		x.XId = primitive.NewObjectID().Hex()
		//tmpResource.GenaiUri = geminiUrl
		x.SessionId = params.SessionId
		x.UploadUserId = userId
		x.UserId = userId

		// 冗余到session
		_, err = t.Data.Mongo.Session.UpdateOne(ctx,
			bson.M{"_id": params.SessionId},
			bson.M{
				"$push": bson.M{"resources": x},
			},
			options.Update().SetUpsert(true))

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
func (t TrackerService) AddResourceByLinkV2(ctx context.Context, params *trackerpb.AddResourceByLinkV2Params) (*trackerpb.AddResourceByLinkV2Result, error) {
	userId := krathelper.RequireUserId(ctx)

	var tmpResources []*trackerpb.ResourceV2
	var tmpResource *trackerpb.ResourceV2
	var err error

	if strings.Contains(params.Link, "v.douyin.com") {
		tmpResource, err = t.addDouyinResource(ctx, params.Link)
	} else if strings.Contains(params.Link, "tiktok.com") {
		tmpResource, err = t.addTiktokResource(ctx, params.Link)

	} else {
		tmpResources, err = t.addXhsResourceV3(ctx, params.Link)
	}

	if err != nil {
		return nil, err
	}

	if tmpResource != nil {
		tmpResources = append(tmpResources, tmpResource)
	}

	for i := range tmpResources {

		x := tmpResources[i]

		if x.IsVideo() || x.IsImage() {

			ossUrl, err := t.uploadToS3(ctx, x.PlatformUrl, x.MimeType)
			if err != nil {
				return nil, err
			}
			x.Url = ossUrl
		}

		x.XId = primitive.NewObjectID().Hex()
		//tmpResource.GenaiUri = geminiUrl
		//x.SessionId = params.SessionId
		x.UploadUserId = userId
		x.UserId = userId
	}

	return &trackerpb.AddResourceByLinkV2Result{
		Resources: tmpResources,
	}, nil
}

func (t TrackerService) ExtractResourceByLink(ctx context.Context, params *trackerpb.ExtractResourceByLinkParams) (*trackerpb.ExtractResourceByLinkResult, error) {
	//userId := krathelper.RequireUserId(ctx)

	var tmpResources []*trackerpb.ResourceV2
	var tmpResource *trackerpb.ResourceV2
	var err error

	if strings.Contains(params.Link, "v.douyin.com") {
		tmpResource, err = t.addDouyinResource(ctx, params.Link)
	} else if strings.Contains(params.Link, "tiktok.com") {
		tmpResource, err = t.addTiktokResource(ctx, params.Link)

	} else {
		tmpResources, err = t.addXhsResourceV3(ctx, params.Link)
	}

	if err != nil {
		return nil, err
	}

	if tmpResource != nil {
		tmpResources = append(tmpResources, tmpResource)
	}

	for i := range tmpResources {

		x := tmpResources[i]

		if x.IsVideo() || x.IsImage() {

			//ossUrl, err := t.uploadToS3(ctx, x.PlatformUrl, x.MimeType)
			//if err != nil {
			//	return nil, err
			//}
			x.Url = x.PlatformUrl

			if x.CoverUrl != "" {
				cover, err := t.upload(ctx, x.CoverUrl, "image/jpeg")
				if err != nil {
					return nil, err
				}
				x.CoverUrl = cover
			}
		}

		//x.XId = primitive.NewObjectID().Hex()
		//tmpResource.GenaiUri = geminiUrl
		//x.SessionId = params.SessionId
		//x.UploadUserId = userId
		//x.UserId = userId
	}

	return &trackerpb.ExtractResourceByLinkResult{
		Resources: tmpResources,
	}, nil
}

func (t TrackerService) AddResourceV2(ctx context.Context, params *trackerpb.AddResourceV2Params) (*trackerpb.ResourceV2, error) {
	userId := krathelper.RequireUserId(ctx)

	var tmpResource *trackerpb.ResourceV2
	var err error

	if strings.Contains(params.Link, "v.douyin.com") {
		tmpResource, err = t.addDouyinResource(ctx, params.Link)
	} else if strings.Contains(params.Link, "tiktok.com") {
		tmpResource, err = t.addTiktokResource(ctx, params.Link)
	} else {
		tmpResource, err = t.addXhsResourceV2(ctx, params.Link)
	}

	if err != nil {
		return nil, err
	}

	//geminiUrl, ossUrl, err := t.upload(ctx, t.Data.GenaiClient, tmpResource.PlatformUrl)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	ossUrl, err := t.uploadToS3(ctx, tmpResource.PlatformUrl, "video/mp4")
	if err != nil {
		return nil, err
	}

	tmpResource.XId = primitive.NewObjectID().Hex()
	//tmpResource.GenaiUri = geminiUrl
	tmpResource.SessionId = params.SessionId
	tmpResource.Url = ossUrl
	tmpResource.UploadUserId = userId

	// 冗余到session
	_, err = t.Data.Mongo.Session.UpdateOne(ctx,
		bson.M{"_id": params.SessionId},
		bson.M{
			"$push": bson.M{"resources": tmpResource},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	return tmpResource, nil
}

func (t TrackerService) upload(ctx context.Context, url, mimeType string) (string, error) {

	//if t.Data.Conf.Biz.Storage == "s3" {
	//	return t.uploadToS3(ctx, url, mimeType)
	//} else {
	//	return t.uploadToAliOss(ctx, url, mimeType)
	//}

	return t.uploadToS3(ctx, url, mimeType)
}

func (t TrackerService) uploadToS3(ctx context.Context, url, mimeType string) (string, error) {
	result, err := resty.New().R().Get(url)
	if err != nil {
		return "", err
	}

	fileBody := bytes.NewReader(result.Body())

	fileBytes, err := io.ReadAll(fileBody)
	if err != nil {
		return "", err
	}

	md5 := helper.MD5(fileBytes)
	//category := filed.FindSuffix(url)

	category := "." + strings.Split(mimeType, "/")[1]

	ctx = context.Background()

	//ossUrl, err := t.Data.S3Client.Upload(ctx, mimeType, md5+category, fileBytes)
	//if err != nil {
	//	return "", err
	//}

	ossUrl, err := t.Data.TosClient.Put(ctx, tos.PutRequest{
		Content: fileBytes,
		Key:     md5 + category,
	})
	if err != nil {
		return "", err
	}

	//fmt.Println(s3Url)

	//ossUrl, err := t.Data.Uploader.Client().Upload(ctx, "oscar-res", md5+category, fileBytes)
	//if err != nil {
	//	log.Errorw("Upload to oss err", err)
	//	return "", err
	//}

	return ossUrl, nil
}

//func (t TrackerService) uploadToAliOss(ctx context.Context, url, mimeType string) (string, error) {
//	result, err := resty.New().R().Get(url)
//	if err != nil {
//		return "", err
//	}
//
//	fileBody := bytes.NewReader(result.Body())
//
//	fileBytes, err := io.ReadAll(fileBody)
//	if err != nil {
//		return "", err
//	}
//
//	md5 := helper.MD5(fileBytes)
//	//category := filed.FindSuffix(url)
//
//	category := "." + strings.Split(mimeType, "/")[1]
//
//	ctx = context.Background()
//
//	ossUrl, err := t.Data.Alioss.Upload(ctx, "veogocn", md5+category, fileBytes)
//	if err != nil {
//		return "", err
//	}
//
//	//fmt.Println(s3Url)
//
//	//ossUrl, err := t.Data.Uploader.Client().Upload(ctx, "oscar-res", md5+category, fileBytes)
//	//if err != nil {
//	//	log.Errorw("Upload to oss err", err)
//	//	return "", err
//	//}
//
//	return ossUrl, nil
//}

func (t TrackerService) mapToMIMEType(name string) string {

	//video/mp4
	//video/mpeg
	//video/mov
	//video/avi
	//video/x-flv
	//video/mpg
	//video/webm
	//video/wmv
	//video/3gpp

	name = strings.ToLower(name)

	if strings.HasSuffix(name, ".mp4") {
		return "video/mp4"
	} else if strings.HasSuffix(name, ".webm") {
		return "video/webm"
	} else if strings.HasSuffix(name, ".avi") {
		return "video/avi"
	} else if strings.HasSuffix(name, ".mpg") {
		return "video/mpg"
	} else if strings.HasSuffix(name, ".mov") {
		return "video/mov"
	} else if strings.HasSuffix(name, ".png") {
		return "image/png"
	} else if strings.HasSuffix(name, ".jpg") {
		return "image/jpeg"
	} else if strings.HasSuffix(name, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(name, ".gif") {
		return "image/gif"
	}

	return ""

}

func (t TrackerService) AddResourceV5(ctx context.Context, params *trackerpb.AddResourceV5Params) (*trackerpb.ResourceV2, error) {
	userId := krathelper.RequireUserId(ctx)

	request, _ := http.RequestFromServerContext(ctx)
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Errorw("ParseMultipartForm err", err)
		return nil, err
	}

	mimeType := request.MultipartForm.Value["mimeType"][0]
	sessionId := request.MultipartForm.Value["sessionId"][0]

	var category string
	if len(request.MultipartForm.Value["category"]) > 0 {
		category = request.MultipartForm.Value["category"][0]
	}

	log.Debugw("AddResourceV5 ", "", "mimeType", mimeType, "sessionId", sessionId, "category", category)

	var newResource *trackerpb.ResourceV2
	resourceId := primitive.NewObjectID().Hex()

	switch mimeType {
	case "text/plain":
		newResource = &trackerpb.ResourceV2{
			XId: resourceId,
			//GenaiUri:     genaiResult.URI,
			UploadUserId: userId,
			CreatedAt:    time.Now().Unix(),
			MimeType:     mimeType,
			SessionId:    sessionId,
			Category:     category,
			Content:      request.MultipartForm.Value["value"][0],
		}
	default:
		f := request.MultipartForm.File["value"][0]

		fileBody, _ := f.Open()

		fileBytes, err := io.ReadAll(fileBody)
		if err != nil {
			return nil, err
		}

		md5 := helper.MD5(fileBytes)
		suffix := filed.FindSuffix(f.Filename)

		ctx = context.Background()

		//ossUrl, err := t.Data.Uploader.Client().Upload(ctx, "oscar-res", md5+category, fileBytes)
		//if err != nil {
		//	log.Errorw("Upload to oss err", err)
		//	return nil, err
		//}
		//
		//ossUrl, err := t.Data.S3Client.Upload(ctx, mimeType, md5+suffix, fileBytes)
		//if err != nil {
		//	log.Errorw("Upload to oss err", err)
		//	return nil, err
		//}

		ossUrl, err := t.Data.TosClient.Put(ctx, tos.PutRequest{
			Content: fileBytes,
			Key:     md5 + suffix,
		})
		if err != nil {
			return nil, err
		}

		newResource = &trackerpb.ResourceV2{
			XId: resourceId,
			//GenaiUri:     genaiResult.URI,
			UploadUserId: userId,
			CreatedAt:    time.Now().Unix(),
			MimeType:     mimeType,
			SessionId:    sessionId,
			Category:     category,
			Name:         f.Filename,
			Url:          ossUrl,
		}
	}

	_, err = t.Data.Mongo.Resource.InsertOne(ctx, newResource)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// 冗余到session
	_, err = t.Data.Mongo.Session.UpdateOne(ctx,
		bson.M{"_id": sessionId},
		bson.M{
			"$push": bson.M{"resources": newResource},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	return newResource, nil

}

func (t TrackerService) AddResourceV3(ctx context.Context, params *trackerpb.AddResourceParams) (*trackerpb.ResourceV2, error) {

	t0 := time.Now()
	log.Debugw("AddResource params", params, "t0", t0)

	userId := krathelper.RequireUserId(ctx)

	request, _ := http.RequestFromServerContext(ctx)

	// google
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Errorw("ParseMultipartForm err", err)
		return nil, err
	}
	files := request.MultipartForm.File["files"]

	sessionId := request.MultipartForm.Value["sessionId"][0]

	if len(files) == 0 {
		return nil, err
	}
	log.Debugw("MultipartForm time", time.Since(t0))

	f := files[0]

	fileBody, _ := f.Open()

	//t1 := time.Now()
	//log.Debugw("start to upload file to genai time", t1)
	//
	//ctx = context.Background()
	//
	mimeType := t.mapToMIMEType(f.Filename)

	//opts := genai.UploadFileOptions{DisplayName: f.Filename, MIMEType: mimeType}
	//genaiResult, err := t.Data.GenaiClient.UploadFile(ctx, "", fileBody, &opts)
	//if err != nil {
	//	log.Errorw("upload file error", err)
	//	return nil, err
	//}
	//
	//log.Debugw("end upload file to genai time", time.Since(t1))
	//t2 := time.Now()

	// oss
	//_, err = fileBody.Seek(0, 0)
	//if err != nil {
	//	return nil, err
	//}

	fileBytes, err := io.ReadAll(fileBody)
	if err != nil {
		return nil, err
	}

	md5 := helper.MD5(fileBytes)
	category := filed.FindSuffix(f.Filename)

	ctx = context.Background()

	//ossUrl, err := t.Data.Uploader.Client().Upload(ctx, "oscar-res", md5+category, fileBytes)
	//if err != nil {
	//	log.Errorw("Upload to oss err", err)
	//	return nil, err
	//}
	//
	//ossUrl, err := t.Data.S3Client.Upload(ctx, "video/mp4", md5+category, fileBytes)
	//if err != nil {
	//	log.Errorw("Upload to oss err", err)
	//	return nil, err
	//}

	ossUrl, err := t.Data.TosClient.Put(ctx, tos.PutRequest{
		Content: fileBytes,
		Key:     md5 + category,
	})
	if err != nil {
		return nil, err
	}

	//log.Debugw("upload file to oss time", time.Since(t2))

	//// 校验
	//for genaiResult.State == genai.FileStateProcessing {
	//	fmt.Print(".")
	//	// Sleep for 10 seconds
	//	time.Sleep(1 * time.Second)
	//	// Fetch the file from the API again.
	//	genaiResult, err = t.Data.GenaiClient.GetFile(ctx, genaiResult.Name)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	log.Debugw("GenaiClient.GetFile response", genaiResult)
	//}

	//
	resourceId := primitive.NewObjectID().Hex()

	newResource := &trackerpb.ResourceV2{
		XId: resourceId,
		//GenaiUri:     genaiResult.URI,
		Name:         f.Filename,
		Url:          ossUrl,
		UploadUserId: userId,
		CreatedAt:    time.Now().Unix(),
		MimeType:     mimeType,
		SessionId:    sessionId,
	}

	_, err = t.Data.Mongo.Resource.InsertOne(ctx, newResource)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// 冗余到session
	_, err = t.Data.Mongo.Session.UpdateOne(ctx,
		bson.M{"_id": sessionId},
		bson.M{
			"$push": bson.M{"resources": newResource},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	return newResource, nil
}

func (t TrackerService) ListResources(ctx context.Context, params *trackerpb.ListResourcesParams) (*trackerpb.ResourceList, error) {
	return &trackerpb.ResourceList{}, nil
}

func (t TrackerService) AddResource(ctx context.Context, params *trackerpb.AddResourceParams) (*trackerpb.Resource, error) {

	//t0 := time.Now()
	//log.Debugw("AddResource params", params, "t0", t0)
	//
	//userId := krathelper.RequireUserId(ctx)
	//
	//request, _ := http.RequestFromServerContext(ctx)
	//
	//// google
	//err := request.ParseMultipartForm(32 << 20)
	//if err != nil {
	//	log.Errorw("ParseMultipartForm err", err)
	//	return nil, err
	//}
	//files := request.MultipartForm.File["files"]
	//
	//if len(files) == 0 {
	//	return nil, err
	//}
	//log.Debugw("MultipartForm time", time.Since(t0))
	//
	//f := files[0]
	//
	//fileBody, _ := f.Open()
	//
	//t1 := time.Now()
	//log.Debugw("start to upload file to genai time", t1)
	//
	//ctx = context.Background()
	//
	//opts := genai.UploadFileOptions{DisplayName: f.Filename, MIMEType: "video/mp4"}
	//response, err := t.Data.GenaiClient.UploadFile(ctx, "", fileBody, &opts)
	//if err != nil {
	//	log.Errorw("upload file error", err)
	//	return nil, err
	//}
	//
	//log.Debugw("end upload file to genai time", time.Since(t1))
	//t2 := time.Now()
	//
	////for response.State == genai.FileStateProcessing {
	////	fmt.Print(".")
	////	// Sleep for 10 seconds
	////	time.Sleep(1 * time.Second)
	////	// Fetch the file from the API again.
	////	response, err = t.Data.GenaiClient.GetFile(ctx, response.Name)
	////	if err != nil {
	////		log.Fatal(err)
	////	}
	////
	////	log.Debugw("GenaiClient.GetFile response", response)
	////}
	//
	//// oss
	//_, err = fileBody.Seek(0, 0)
	//if err != nil {
	//	return nil, err
	//}
	//
	//fileBytes, err := io.ReadAll(fileBody)
	//if err != nil {
	//	return nil, err
	//}
	//
	//md5 := helper.MD5(fileBytes)
	//category := filed.FindSuffix(f.Filename)
	//
	//ctx = context.Background()
	//
	//ossUrl, err := t.Data.Uploader.Client().Upload(ctx, "oscar-res", md5+category, fileBytes)
	//if err != nil {
	//	log.Errorw("Upload to oss err", err)
	//	return nil, err
	//}
	//
	//log.Debugw("upload file to oss time", time.Since(t2))
	//t3 := time.Now()
	//
	//id, err := t.Data.EntClient.Resource.Create().
	//	SetKey(md5).
	//	SetUserID(userId).
	//	SetCategory(category).
	//	SetURL(ossUrl).
	//	SetName(f.Filename).
	//	SetLastUpdatedAt(time.Now()).
	//	SetExtra(entz.Extra{
	//		"genaiUri": response.URI,
	//	}).
	//	OnConflictColumns(resource.FieldUserID, resource.FieldKey).
	//	UpdateFields(func(upsert *ent.ResourceUpsert) {
	//		upsert.UpdateExtra()
	//		upsert.UpdateLastUpdatedAt()
	//	}).
	//	ID(ctx)
	//if err != nil {
	//	log.Errorw("create resource err", err)
	//	return nil, err
	//}
	//log.Debugw("EntClient time", time.Since(t3))

	return &trackerpb.Resource{
		//Id:       conv.Str(id),
		//Url:      ossUrl,
		//Category: category,
		//GenaiUri: response.URI,
		//Name:     f.Filename,
	}, nil
}
