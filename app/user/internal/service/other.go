package service

import (
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	userpb "store/api/user"
	"store/pkg/krathelper"
)

func (t *UserService) GetServerMeta(ctx context.Context, _ *empty.Empty) (*userpb.ServerMeta, error) {
	return &userpb.ServerMeta{
		CountryCode: krathelper.ClientCountryCode(ctx),
	}, nil

}
