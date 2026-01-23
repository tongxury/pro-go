package biz

import "store/app/usercenter/internal/data"

type ItemBiz struct {
	data *data.Data
}

func NewItemBiz(data *data.Data) *ItemBiz {
	return &ItemBiz{
		data: data,
	}
}
