package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	trackerpb "store/api/aiagent"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/urlz"
	"strings"
)

func (t TrackerService) addXhsProfile(ctx context.Context, userId, link string) (*trackerpb.Profile, error) {

	urlInfo, err := urlz.ParseURL(link)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(urlInfo.Path, "/user/profile") {
		return nil, errors.New(0, "invalidLink", "")
	}

	//profile, err := t.Data.XhsClient.GetProfileByLink(ctx, link)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if profile.Username == "" {
	//	return nil, errors.New(0, "invalidLink", "")
	//}

	user, err := t.Data.Tikhub.XhsGetUserById(ctx, urlInfo.PathSegments[len(urlInfo.PathSegments)-1])
	if err != nil {
		return nil, errors.New(0, "invalidLink", "")
	}

	id := primitive.NewObjectID().Hex()

	var tags []string
	for _, x := range user.Data.Tags {
		tags = append(tags, x.Name)
	}
	newProfile := &trackerpb.Profile{
		XId:      id,
		Platform: "xiaohongshu",
		//Avatar:         user.Data.Image,
		Username:       user.Data.Nickname,
		PlatformId:     user.Data.RedId,
		IpAddress:      user.Data.IpLocation,
		Sign:           user.Data.Desc,
		Tags:           tags,
		FollowingCount: conv.Str(user.Data.Follows),
		FollowerCount:  conv.Str(user.Data.Fans),
		LikedCount:     conv.Str(user.Data.Liked),
		//NoteCount:      conv.Str(user.Data.Notes),
		UserId: userId,
		//LastUpdatedAt: time.Now().Unix(),
	}

	err = t.Data.Mongo.Profile.Insert(ctx, newProfile)
	if err != nil {
		return nil, err
	}

	return newProfile, nil
}

func (t TrackerService) addTiktokProfile(ctx context.Context, userId, douyinCode string) (*trackerpb.Profile, error) {

	profile, err := t.Data.Tikhub.GetTiktokUserProfile(ctx, douyinCode)
	if err != nil {
		return nil, err
	}

	id := primitive.NewObjectID().Hex()
	newProfile := &trackerpb.Profile{
		XId:      id,
		Platform: "tiktok",
		//Avatar:         profile.User,
		Username: profile.User.Nickname,
		//PlatformId:     profile.Id,
		IpAddress: profile.User.IpLocation,
		Sign:      profile.User.Signature,
		//Tags:           profile.Tags,
		FollowingCount: conv.Str(profile.User.FollowingCount),
		FollowerCount:  conv.Str(profile.User.FollowerCount),
		LikedCount:     conv.Str(profile.User.TotalFavorited),
		NoteCount:      conv.Str(profile.User.AwemeCount),
		UserId:         userId,
		//LastUpdatedAt:  time.Now().Unix(),
		//Raw:            conv.S2J(profile),
	}

	err = t.Data.Mongo.Profile.Insert(ctx, newProfile)
	if err != nil {
		return nil, err
	}

	return newProfile, nil
}

func (t TrackerService) addDouyinProfile(ctx context.Context, userId, douyinCode string) (*trackerpb.Profile, error) {

	user, err := t.Data.Tikhub.DouyinGetUserProfileV2(ctx, douyinCode)
	if err != nil {
		return nil, err
	}
	//
	//profile, err := t.Data.Tikhub.DouyinGetUserProfile(ctx, user.UserInfo.SecUid)
	//if err != nil {
	//	return nil, err
	//}

	id := primitive.NewObjectID().Hex()
	newProfile := &trackerpb.Profile{
		XId:      id,
		Platform: "douyin",
		//Avatar:         profile.User,
		Username: user.UserInfo.Nickname,
		//PlatformId:     profile.Id,
		//IpAddress: profile.User.IpLocation,
		Sign: user.UserInfo.Signature,
		//Tags:           profile.Tags,
		FollowingCount: conv.Str(user.UserInfo.FollowingCount),
		FollowerCount:  conv.Str(user.UserInfo.MplatformFollowersCount),
		LikedCount:     conv.Str(user.UserInfo.TotalFavorited),
		NoteCount:      conv.Str(user.UserInfo.AwemeCount),
		UserId:         userId,
		//LastUpdatedAt:  time.Now().Unix(),
		//Raw:            conv.S2J(profile),
	}

	err = t.Data.Mongo.Profile.Insert(ctx, newProfile)
	if err != nil {
		return nil, err
	}

	return newProfile, nil
}

func (t TrackerService) AddProfile(ctx context.Context, params *trackerpb.AddProfileParams) (*trackerpb.Profile, error) {
	userId := krathelper.RequireUserId(ctx)

	var profile *trackerpb.Profile
	var err error

	switch params.Platform {
	case "xiaohongshu":
		profile, err = t.addXhsProfile(ctx, userId, params.Link)
	case "douyin":
		profile, err = t.addDouyinProfile(ctx, userId, params.Keyword)
	case "tiktok":
		profile, err = t.addTiktokProfile(ctx, userId, params.Link)
	default:
		profile, err = t.addXhsProfile(ctx, userId, params.Link)
	}

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (t TrackerService) GetProfile(ctx context.Context, params *trackerpb.GetProfileParams) (*trackerpb.Profile, error) {
	//userId := krathelper.RequireUserId(ctx)

	profile, err := t.Data.Mongo.Profile.FindById(ctx, params.Id)
	if err != nil {
		log.Errorw("get profile error", err)
		return nil, err
	}

	return profile, nil
}

func (t TrackerService) ListProfiles(ctx context.Context, params *trackerpb.ListProfilesParams) (*trackerpb.ProfileList, error) {
	userId := krathelper.RequireUserId(ctx)

	filters := bson.M{
		"userId": userId,
	}

	if params.Platform != "" {
		filters["platform"] = params.Platform
	}

	list, err := t.Data.Mongo.Profile.List(ctx, filters,
		options.Find().SetSort(bson.M{"lastUpdatedAt": -1}),
	)
	if err != nil {
		return nil, err
	}

	return &trackerpb.ProfileList{
		List: list,
	}, nil
}

type XhsResult struct {
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

func (t TrackerService) DeleteProfile(ctx context.Context, params *trackerpb.DeleteProfileParams) (*empty.Empty, error) {
	//userId := krathelper.RequireUserId(ctx)

	err := t.Data.Mongo.Profile.DestroyById(ctx, params.Id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &empty.Empty{}, nil
}
