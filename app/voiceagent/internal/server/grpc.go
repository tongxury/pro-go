package server

import (
	voiceagent "store/api/voiceagent"
	"store/app/voiceagent/internal/service"
	"store/pkg/confcenter"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c confcenter.Server, ss *service.VoiceAgentService, livekit *service.LiveKitService, logger log.Logger) *grpc.Server {
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
	voiceagent.RegisterVoiceAgentServiceServer(srv, ss)
	voiceagent.RegisterMotivationServiceServer(srv, ss)
	voiceagent.RegisterLiveKitServiceServer(srv, livekit)
	return srv
}
