package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	userpb "store/api/user"
	typepb "store/api/user/types"
	"store/app/user/internal/data"
	"store/app/user/internal/data/repo"
	"store/app/user/internal/data/repo/ent"
	"store/app/user/internal/data/repo/ent/user"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/timed"
	"time"
)

type UserBiz struct {
	data *data.Data
}

func NewUserBiz(data *data.Data) *UserBiz {
	return &UserBiz{data: data}
}

type GetUserExtra struct {
	Bonus int64
	Limit int64
}

func (t *UserBiz) UpdateUser(ctx context.Context, params *userpb.UpdateUserParams) (*emptypb.Empty, error) {

	var err error
	switch params.Action {
	case userpb.UpdateUserParams_addBonus:
		//amount := params.Values.Amount
		//err = t.data.Repos.EntClient.User.UpdateFields().
		//	Where(user.ID(conv.Int64(params.Id))).
		//	AddBonus(amount).Exec(ctx)

	}

	if err != nil {
		log.Errorw("UpdateUser err")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (t *UserBiz) Delete(ctx context.Context, params *userpb.DeleteUserParams) (*emptypb.Empty, error) {

	_, err := t.data.Repos.EntClient.User.Delete().
		Where(user.Or(
			user.Email(params.IdOrEmail),
			user.ID(conv.Int64(params.IdOrEmail)),
		)).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

type ListUsersParams struct {
	UserID, Keyword string
	Emails          []string
	Page, Size      int64
}

func (t *UserBiz) GetUsersSummary(ctx context.Context) (*userpb.UserSummary, error) {

	var totalCounts []struct {
		Count int64
	}
	err := t.data.Repos.EntClient.User.Query().
		Aggregate(ent.As(ent.Count(), "count")).
		Scan(ctx, &totalCounts)
	if err != nil {
		return nil, err
	}

	var todayCounts []struct {
		Count int64
	}

	err = t.data.Repos.EntClient.User.Query().
		Where(user.CreatedAtGTE(timed.DateStart(time.Now()))).
		Aggregate(ent.As(ent.Count(), "count")).
		Scan(ctx, &todayCounts)
	if err != nil {
		return nil, err
	}

	return &userpb.UserSummary{
		TotalCount: totalCounts[0].Count,
		TodayCount: todayCounts[0].Count,
	}, nil

}

func (t *UserBiz) ListUsers(ctx context.Context, params ListUsersParams) (*userpb.ListUsersResult, error) {

	users, total, err := t.data.Repos.User.List(ctx, repo.ListParams{
		Keyword: params.Keyword,
		UserID:  params.UserID,
		Emails:  params.Emails,
		Page:    params.Page,
		Size:    params.Size,
	})
	if err != nil {
		return nil, err
	}

	return &userpb.ListUsersResult{
		Page:  params.Page,
		Size:  params.Size,
		Total: total,
		List:  t.asPbUsers(users),
	}, nil
}

func (t *UserBiz) asPbUsers(users repo.Users) []*typepb.User {

	var rsp []*typepb.User
	for _, x := range users {

		y := &typepb.User{
			//Id:          conv.String(x.ID),
			Username: x.Nickname,
			//UserAvatar:  x.Avatar,
			Email:     x.Email,
			CreatedAt: x.CreatedAt.Unix(),
			//LastLoginAt: x.LastLoginAt.Unix(),
			//Bonus:       x.Bonus,
		}

		rsp = append(rsp, y)
	}

	return rsp
}

func (t *UserBiz) GetUserDetail(ctx context.Context, userId, authPlatform string) (*typepb.User, error) {

	dbUser, err := t.data.Repos.User.GetById(ctx, userId, true)

	if err != nil {
		log.Errorw("GetUserDetail err", err, "userId", userId)
		return nil, err
	}

	return &typepb.User{
		Id:         conv.String(dbUser.ID),
		Username:   helper.OrString(dbUser.Nickname, fmt.Sprintf("V%d", dbUser.ID)),
		UserAvatar: "",
		//UserAvatar: helper.OrString(confcenter.DefaultUserAvatar, conv.String(dbUser.GoogleAuth["picture"])),
		Email:       dbUser.Email,
		CreatedAt:   0,
		LastLoginAt: 0,
		Bonus:       0,
		Phone:       dbUser.Phone,
	}, err
}
