package biz

import "store/app/proj-pro/internal/data"

type ItemBiz struct {
	data *data.Data
}

func NewItemBiz(data *data.Data) *ItemBiz {
	return &ItemBiz{
		data: data,
	}
}
