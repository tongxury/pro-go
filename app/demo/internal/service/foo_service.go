package service

import (
	demopb "store/api/demo"
	"store/app/demo/internal/biz"
	"store/app/demo/internal/data"
)

type FooService struct {
	demopb.UnimplementedFooServiceServer
	data *data.Data
	foo  *biz.FooBiz
}

func NewFooService(data *data.Data, foo *biz.FooBiz) *FooService {
	return &FooService{
		data: data,
		foo:  foo,
	}
}
