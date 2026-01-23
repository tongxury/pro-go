package service

import (
	"context"
	ucpb "store/api/usercenter"
	"store/pkg/krathelper"

	"github.com/go-kratos/kratos/v2/log"
)

func (t *UserService) GetUser(ctx context.Context, request *ucpb.GetUserRequest) (*ucpb.User, error) {

	userId := krathelper.FindUserId(ctx)

	return t.XGetUser(ctx, &ucpb.XGetUserRequest{
		Id: userId,
	})
}

func (t *UserService) XGetUser(ctx context.Context, request *ucpb.XGetUserRequest) (*ucpb.User, error) {

	userId := request.Id
	if userId == "" {
		return nil, nil
	}

	user, err := t.data.Mongo.User.FindByID(ctx, userId)
	if err != nil {
		log.Errorw("GetById err", err, "userId", userId)
		return nil, err
	}

	return user, nil
}
