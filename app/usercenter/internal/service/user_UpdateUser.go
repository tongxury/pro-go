package service

import (
	"context"
	ucpb "store/api/usercenter"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
)

func (t *UserService) UpdateUser(ctx context.Context, req *ucpb.UpdateUserRequest) (*ucpb.User, error) {

	userId := krathelper.FindUserId(ctx)

	switch req.Action {
	case "updateNickname":
		_, err := t.data.Mongo.User.UpdateByIDIfExists(ctx, userId,
			mgz.Op().Set("nickname", req.Params["nickname"]))
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
