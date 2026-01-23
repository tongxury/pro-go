package service

import (
	"context"
	"fmt"
	ucpb "store/api/usercenter"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

func (t AuthService) GetToken(ctx context.Context, request *ucpb.GetTokenRequest) (*ucpb.Token, error) {

	redisKey := fmt.Sprintf("authCode:%s", request.Phone)

	code, err := t.data.Redis.Get(ctx, redisKey).Result()
	if err != nil {
		log.Errorw("RedisClient get error", err, "params", request)
		return nil, err
	}

	if code != request.Code {
		return nil, errors.Unauthorized("", "验证码错误")
	}

	userId, isNew, err := t.data.Mongo.User.InsertNX(ctx,
		&ucpb.User{
			Key:       request.Phone,
			Phone:     request.Phone,
			CreatedAt: time.Now().Unix(),
		},
		mgz.Filter().EQ("key", request.Phone).
			B(),
	)

	if err != nil {
		log.Errorw("InsertNX err", err, "params", request)
		return nil, err
	}

	token, err := krathelper.GenerateTokenV2(userId)

	return &ucpb.Token{
		Token:  token,
		IsNew:  isNew,
		UserId: userId,
	}, nil
}
