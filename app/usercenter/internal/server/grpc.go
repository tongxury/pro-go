package server

import (
	creditpb "store/api/credit"
	toolkitpb "store/api/toolkit"
	ucpb "store/api/usercenter"
	"store/app/usercenter/internal/service"
	"store/pkg/confcenter"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c confcenter.Server,
	user *service.UserService,
	auth *service.AuthService,
	toolkit *service.ToolkitService,
	credit *service.CreditService,
	logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout > 0 {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout))
	}
	srv := grpc.NewServer(opts...)

	ucpb.RegisterAuthServiceServer(srv, auth)
	ucpb.RegisterUserServiceServer(srv, user)
	creditpb.RegisterCreditServiceServer(srv, credit)
	toolkitpb.RegisterToolkitServiceServer(srv, toolkit)
	return srv
}
