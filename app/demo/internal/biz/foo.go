package biz

import "store/app/demo/internal/data"

type FooBiz struct {
	data *data.Data
}

func NewFooBiz(data *data.Data) *FooBiz {
	return &FooBiz{
		data: data,
	}
}
