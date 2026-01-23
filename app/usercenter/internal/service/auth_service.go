package service

import (
	ucpb "store/api/usercenter"
	"store/app/usercenter/internal/biz"
	"store/app/usercenter/internal/data"
)

type AuthService struct {
	ucpb.UnimplementedAuthServiceServer
	data *data.Data
	item *biz.ItemBiz
}

func NewAuthService(data *data.Data, item *biz.ItemBiz) *AuthService {
	return &AuthService{
		data: data,
		item: item,
	}
}
