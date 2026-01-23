package service

import (
	projpb "store/api/proj"
	"store/app/proj-pro/internal/biz"
	"store/app/proj-pro/internal/data"
)

type SessionService struct {
	projpb.UnimplementedSessionServiceServer
	data *data.Data
	item *biz.ItemBiz
}

func NewSessionService(data *data.Data, item *biz.ItemBiz) *SessionService {
	return &SessionService{
		data: data,
		item: item,
	}
}
