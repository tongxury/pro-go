package service

import (
	"context"
	ucpb "store/api/usercenter"
	"store/pkg/clients/mgz"

	"github.com/go-kratos/kratos/v2/log"
)

func (t *UserService) XGetUserByPhone(ctx context.Context, request *ucpb.XGetUserByPhoneRequest) (*ucpb.User, error) {

	phone := request.Phone
	if phone == "" {
		return nil, nil
	}

	user, err := t.data.Mongo.User.FindOne(ctx, mgz.Filter().EQ("phone", phone).B())
	if err != nil {
		log.Errorw("GetById err", err, "phone", phone)
		return nil, err
	}

	return user, nil
}
