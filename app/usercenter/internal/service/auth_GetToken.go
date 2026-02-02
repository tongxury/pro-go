package service

import (
	"context"
	"fmt"
	ucpb "store/api/usercenter"
	"store/pkg/clients/mgz"
	"store/pkg/events"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
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

	msg := kafka.Message{
		Value: conv.S2B(events.AuthEvent{
			UserID:     userId,
			LoginBy:    "phone",
			TS:         time.Now().Unix(),
			IsRegister: isNew,
		}),
	}

	err = t.data.KafkaClient.W().Write(ctx, events.Topic_AuthLogin, msg)
	if err != nil {
		log.Errorw("KafkaClient.W().Write err", err, "msg", msg)
	}

	return &ucpb.Token{
		Token:  token,
		IsNew:  isNew,
		UserId: userId,
	}, nil
}
