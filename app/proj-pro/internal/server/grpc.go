package server

import (
	projpb "store/api/proj"
	"store/app/proj-pro/internal/service"
	"store/pkg/confcenter"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c confcenter.Server,
	greeter *service.ProjService,
	session *service.SessionService,
	workflow *service.WorkFlowService,
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
	projpb.RegisterProjProServiceServer(srv, greeter)
	projpb.RegisterSessionServiceServer(srv, session)
	projpb.RegisterWorkflowServiceServer(srv, workflow)
	return srv
}
