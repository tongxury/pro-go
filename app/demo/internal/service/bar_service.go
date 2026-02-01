package service

import (
	demopb "store/api/demo"
	"store/app/demo/internal/data"
)

type BarService struct {
	demopb.UnimplementedBarServiceServer
	data *data.Data
}

func NewBarService(data *data.Data) *BarService {
	return &BarService{
		data: data,
	}
}
