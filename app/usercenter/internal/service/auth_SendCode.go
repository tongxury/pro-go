package service

import (
	"context"
	"fmt"
	ucpb "store/api/usercenter"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper/mathz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (t AuthService) SendCode(ctx context.Context, request *ucpb.SendCodeRequest) (*emptypb.Empty, error) {

	redisKey := fmt.Sprintf("authCode:%s", request.Phone)

	code := mathz.RandNumber(100000, 999999)

	log.Debugw("验证码", "", "phone", request.Phone, "code", code)

	if t.data.VolcSmsClient != nil {
		err := t.data.VolcSmsClient.SendSmsCode(ctx, request.Phone, conv.Str(code))
		if err != nil {
			return nil, err
		}
	} else if t.data.Alisms != nil {
		_, err := t.data.Alisms.Send(ctx, []string{request.Phone}, conv.Str(code))
		if err != nil {
			return nil, err
		}
	}

	err := t.data.Redis.Set(ctx, redisKey, code, 5*time.Minute).Err()
	if err != nil {
		log.Errorw("RedisClient set error", err, "params", request)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
