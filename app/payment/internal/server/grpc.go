package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	paymentpb "store/api/payment"
	"store/app/payment/internal/service"
	"store/pkg/confcenter"
	"store/pkg/middlewares/prometircs"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(conf confcenter.Server, payment *service.PaymentService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			logging.Server(logger),
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					return nil
				}),
			),
			validate.Validator(),
			tracing.Server(),
			ratelimit.Server(),
			prometircs.Metrics(),
		),
	}
	if conf.Grpc.Network != "" {
		opts = append(opts, grpc.Network(conf.Grpc.Network))
	}
	if conf.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(conf.Grpc.Addr))
	}
	if conf.Grpc.Timeout > 0 {
		opts = append(opts, grpc.Timeout(conf.Grpc.Timeout))
	}
	grpcSrv := grpc.NewServer(opts...)
	paymentpb.RegisterPaymentServer(grpcSrv, payment)
	return grpcSrv
}
