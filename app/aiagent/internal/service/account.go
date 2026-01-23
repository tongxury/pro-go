package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/tikhub"
	"strings"
	"time"
)

func (t TrackerService) ListAccounts(ctx context.Context, params *aiagentpb.ListAccountsParams) (*aiagentpb.AccountList, error) {
	userId := krathelper.FindUserId(ctx)

	// 避免前端调登录
	if userId == "" {
		return &aiagentpb.AccountList{}, nil
	}

	list, err := t.Data.Mongo.Account.List(ctx,
		bson.M{"userId": userId},
		options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}

	return &aiagentpb.AccountList{
		List: list,
	}, nil
}

func (t TrackerService) UpdateAccount(ctx context.Context, params *aiagentpb.UpdateAccountParams) (*aiagentpb.Account, error) {
	userId := krathelper.RequireUserId(ctx)

	account, err := t.Data.Mongo.Account.GetById(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	switch params.Action {
	case "delete":
		_, err = t.Data.Mongo.Account.DestroyById(ctx, params.GetId())
	case "update":

		switch account.Platform {
		case "xiaohongshu":

			platformId := account.Extra["redId"]
			if platformId == "" {
				return account, nil
			}

			newAccount, err := t.getXhsAccount(ctx, userId, platformId)
			if err != nil {
				return nil, err
			}

			up := bson.M{
				"lastUpdatedAt": time.Now().Unix(),
			}

			if newAccount.Posts != "" {
				up["posts"] = newAccount.Posts
			}

			if newAccount.Followers != "" {
				up["followers"] = newAccount.Followers
			}

			if newAccount.Interacts != "" {
				up["interacts"] = newAccount.Interacts
			}

			if newAccount.Posts != "" {
				up["sign"] = newAccount.Sign
			}

			_, err = t.Data.Mongo.Account.UpdateOneXX(ctx,
				bson.M{
					"_id": mgz.ObjectId(params.Id),
				},
				up,
			)
		case "douyin":

			platformId := account.Extra["douyinAccount"]
			if platformId == "" {
				return account, nil
			}

			newAccount, err := t.getDouyinAccount(ctx, userId, platformId)
			if err != nil {
				return nil, err
			}

			_, err = t.Data.Mongo.Account.UpdateOneXX(ctx,
				bson.M{
					"_id": mgz.ObjectId(params.Id),
				},
				bson.M{
					"followers":     newAccount.Followers,
					"interacts":     newAccount.Interacts,
					"posts":         newAccount.Posts,
					"sign":          newAccount.Sign,
					"lastUpdatedAt": time.Now().Unix(),
				},
			)

		}

	case "setDefault":

		if account.IsDefault {
			return account, nil
		}

		_, err = t.Data.Mongo.Account.UpdateOneXX(ctx,
			bson.M{"isDefault": true},
			bson.M{"isDefault": false},
		)

		account, err = t.Data.Mongo.Account.FindAndUpdateOne(ctx,
			bson.M{"_id": mgz.ObjectId(params.GetId())},
			bson.M{"isDefault": true},
		)
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return account, nil

}

func (t TrackerService) BindAccount(ctx context.Context, params *aiagentpb.BindAccountParams) (*aiagentpb.Account, error) {

	userId := krathelper.RequireUserId(ctx)

	filters := bson.M{
		"userId":   userId,
		"extra.id": params.Keyword,
		"platform": params.Platform,
	}

	accounts, err := t.Data.Mongo.Account.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	if len(accounts) > 0 {
		return accounts[0], nil
	}

	var newAccount *aiagentpb.Account

	switch params.Platform {
	case "xiaohongshu":

		if strings.HasPrefix(params.Keyword, "https://") {
			newAccount, err = t.getXhsAccountByImage(ctx, userId, params.Keyword)
		} else {
			newAccount, err = t.getXhsAccount(ctx, userId, params.Keyword)
		}

	case "douyin":
		newAccount, err = t.getDouyinAccount(ctx, userId, params.Keyword)
	//case "tiktok":
	//	profile, err = t.addTiktokProfile(ctx, userId, params.Link)
	default:
		newAccount, err = t.getXhsAccount(ctx, userId, params.Keyword)
	}

	if err != nil {
		//errorNotFound
		return nil, errors.New(0, "bindErr", "bindErr")
	}

	id, err := t.Data.Mongo.Account.Insert(ctx, newAccount)
	if err != nil {
		return nil, err
	}

	if id != "" {
		newAccount.XId = id

		_, _ = t.Data.Mongo.Account.UpdateOneXX(ctx, bson.M{"isDefault": true}, bson.M{"isDefault": false})
		_, _ = t.Data.Mongo.Account.UpdateOneXX(ctx, bson.M{"_id": mgz.ObjectId(id)}, bson.M{"isDefault": true})
		//_, _ = t.Data.Mongo.Account.UpdateOneXX(ctx, bson.M{"_id": id}, bson.M{"isDefault": true})
	}

	return newAccount, nil
}

func (t TrackerService) getDouyinAccount(ctx context.Context, userId, douyinCode string) (*aiagentpb.Account, error) {

	user, err := t.Data.Tikhub.DouyinGetUserProfileV2(ctx, douyinCode)
	if err != nil {
		return nil, err
	}
	//
	//profile, err := t.Data.Tikhub.DouyinGetUserProfile(ctx, user.UserInfo.SecUid)
	//if err != nil {
	//	return nil, err
	//}

	newProfile := &aiagentpb.Account{
		Platform: "douyin",
		Avatar:   user.GetAvatar(),
		Nickname: user.UserInfo.Nickname,
		//PlatformId:     profile.Id,
		//IpAddress: profile.User.IpLocation,
		Sign: user.UserInfo.Signature,
		//Tags:           profile.Tags,
		//FollowingCount: conv.Str(user.UserInfo.FollowingCount),
		Followers:     conv.Str(user.UserInfo.MplatformFollowersCount),
		Interacts:     conv.Str(user.UserInfo.TotalFavorited),
		Posts:         conv.Str(user.UserInfo.AwemeCount),
		UserId:        userId,
		LastUpdatedAt: time.Now().Unix(),
		//Raw:            conv.S2J(profile),
		Extra: map[string]string{
			"douyinUniqueId": user.UserInfo.UniqueId,
			"douyinShortId":  user.UserInfo.ShortId,
			"douyinSecUid":   user.UserInfo.SecUid,
			"douyinAccount":  douyinCode,
		},
	}

	return newProfile, nil
}

func (t TrackerService) getXhsAccountByImage(ctx context.Context, userId, url string) (*aiagentpb.Account, error) {

	genaiClient := t.Data.GenaiFactory.Get()

	// 1. 准备参数
	parts := []*genai.Part{
		gemini.NewTextPart("帮我分析图片中的博主信息"),
	}

	p, err := gemini.NewImagePart(url)
	if err != nil {
		return nil, err
	}
	parts = append(parts, p)

	// 2. 配置
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:     genai.TypeObject,
			Required: []string{"nickname", "sign", "tags", "followers", "interacts", "extra"},
			Properties: map[string]*genai.Schema{
				"nickname": {
					Type:        genai.TypeString,
					Description: "昵称",
				},
				"sign": {
					Type:        genai.TypeString,
					Description: "个性签名",
				},
				"tags": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
					Description: "标签",
				},
				"followers": {
					Type:        genai.TypeString,
					Description: "粉丝",
				},
				"interacts": {
					Type:        genai.TypeString,
					Description: "获赞与收藏",
				},
				"extra": {
					Type:     genai.TypeObject,
					Required: []string{"redId"},
					Properties: map[string]*genai.Schema{
						"redId": {
							Type:        genai.TypeString,
							Description: "小红书号",
						},
					},
				},
			},
		},
	}

	// 3. 调用
	text, err := genaiClient.GenerateContent(ctx, gemini.GenerateContentRequest{
		Model:  "gemini-2.5-pro-preview-05-06",
		Parts:  parts,
		Config: config,
	})

	if err != nil {
		log.Errorw("GenerateContent err", err)
		return nil, err
	}

	var account *aiagentpb.Account

	log.Debugw("account", account)

	err = json.Unmarshal([]byte(text), &account)
	if err != nil {
		return nil, err
	}

	if account != nil {
		account.UserId = userId
		account.Platform = "xiaohongshu"
		account.Avatar = url
	}

	return account, nil
}

func (t TrackerService) getXhsAccount(ctx context.Context, userId, keyword string) (*aiagentpb.Account, error) {

	users, err := t.Data.Tikhub.XhsSearchUsers(ctx, keyword)
	if err != nil {
		return nil, err
	}

	users = helper.Filter(users, func(param tikhub.XHSUser) bool {
		return param.RedId == keyword
	})

	if len(users) == 0 {
		return nil, fmt.Errorf("account not found by " + keyword)
	}

	user := users[0]

	user2, err := t.Data.Tikhub.XhsGetUserById(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	//id := primitive.NewObjectID().Hex()

	// todo
	notes := conv.Str(user2.Data.Notes)
	if notes == "0" {
		notes = user.GetNotes()
	}

	var tags []string
	for _, x := range user2.Data.Tags {
		tags = append(tags, x.Name)
	}
	newProfile := &aiagentpb.Account{
		//XId:       id,
		Nickname:  user.Name,
		Sign:      user2.Data.Desc,
		Tags:      tags,
		Domain:    nil,
		Posts:     notes,
		Interacts: conv.Str(user2.Data.Liked),
		Followers: conv.Str(user2.Data.Fans),
		Platform:  "xiaohongshu",
		Avatar:    user.Image,
		//Username:       user.Data.Nickname,
		//PlatformId:     user.Data.RedId,
		//IpAddress:      user.Data.IpLocation,
		//Sign:           user.Data.Desc,
		//Tags:           tags,
		//FollowingCount: conv.Str(user.Data.Follows),
		//FollowerCount:  conv.Str(user.Data.Fans),
		//LikedCount:     conv.Str(user.Data.Liked),
		//NoteCount:      conv.Str(user.Data.Notes),
		UserId:    userId,
		IsDefault: false,
		Extra: map[string]string{
			"id":    user.Id,
			"redId": user.RedId,
		},
	}

	return newProfile, nil
}

func (t TrackerService) AddAccount(ctx context.Context, params *aiagentpb.AddAccountParams) (*aiagentpb.Account, error) {

	userId := krathelper.RequireUserId(ctx)

	newAccount := &aiagentpb.Account{
		Nickname:  params.Nickname,
		Sign:      params.Sign,
		Tags:      nil,
		Domain:    params.Domain,
		Posts:     params.Posts,
		Interacts: params.Interacts,
		Followers: params.Followers,
		Platform:  params.Platform,
		UserId:    userId,
	}

	_, id, err := t.Data.Mongo.Account.InsertNX(ctx,
		bson.M{
			"userId":   userId,
			"nickname": params.Nickname,
			"platform": params.Platform,
		},
		newAccount,
	)
	if err != nil {
		return nil, err
	}
	newAccount.XId = id

	_, _ = t.Data.Mongo.Account.UpdateOneXX(ctx, bson.M{"idDefault": true}, bson.M{"isDefault": false})
	_, _ = t.Data.Mongo.Account.UpdateOneXX(ctx, bson.M{"_id": id}, bson.M{"isDefault": true})

	return newAccount, nil
}

func (t TrackerService) PutAccount(ctx context.Context, params *aiagentpb.AddAccountParams) (*aiagentpb.Account, error) {

	userId := krathelper.RequireUserId(ctx)

	newAccount := &aiagentpb.Account{

		Nickname:  params.Nickname,
		Sign:      params.Sign,
		Tags:      nil,
		Domain:    params.Domain,
		Posts:     params.Posts,
		Interacts: params.Interacts,
		Followers: params.Followers,
		Platform:  params.Platform,
		UserId:    userId,
	}

	err := t.Data.Mongo.Account.ReplaceOneXX(ctx, bson.M{"_id": mgz.ObjectId(params.XId)}, newAccount)
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}
