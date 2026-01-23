package service

import (
	creditpb "store/api/credit"
	"store/app/usercenter/internal/biz"
	"store/app/usercenter/internal/data"
)

type CreditService struct {
	creditpb.UnimplementedCreditServiceServer
	data *data.Data
	item *biz.ItemBiz
}

func NewCreditService(data *data.Data, item *biz.ItemBiz) *CreditService {
	return &CreditService{
		data: data,
		item: item,
	}
}
