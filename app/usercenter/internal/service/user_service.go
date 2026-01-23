package service

import (
	ucpb "store/api/usercenter"
	"store/app/usercenter/internal/biz"
	"store/app/usercenter/internal/data"
)

type UserService struct {
	ucpb.UnimplementedUserServiceServer
	data *data.Data
	item *biz.ItemBiz
}

func NewUserService(data *data.Data, item *biz.ItemBiz) *UserService {
	return &UserService{
		data: data,
		item: item,
	}
}
