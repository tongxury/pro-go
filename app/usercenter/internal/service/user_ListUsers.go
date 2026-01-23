package service

import (
	"context"
	ucpb "store/api/usercenter"
	"store/pkg/clients/mgz"

	"github.com/go-kratos/kratos/v2/log"
)

func (t *UserService) XListUsers(ctx context.Context, request *ucpb.XListUsersRequest) (*ucpb.UserList, error) {

	filter := mgz.Filter()

	if request.Keyword != "" {
		filter = filter.HasPrefix("phone", request.Keyword)
	}

	list, total, err := t.data.Mongo.User.ListAndCount(ctx,
		filter.B(),
		mgz.Find().
			PageSize(request.Page, request.Size).
			SetSort("createdAt", -1).
			B(),
	)

	if err != nil {
		log.Errorw("ListAndCount err", err)
		return nil, err
	}

	return &ucpb.UserList{
		List:  list,
		Total: total,
	}, nil
}
