package server

import (
	demopb "store/api/demo"
	"store/app/demo/internal/service"
	"store/pkg/confcenter"
	"store/pkg/middlewares/encoder"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c confcenter.Server,
	foo *service.FooService,
	bar *service.BarService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.ErrorEncoder(encoder.ErrorEncoder),
		http.ResponseEncoder(encoder.ResponseEncoder),
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout > 0 {
		opts = append(opts, http.Timeout(c.Http.Timeout))
	}
	srv := http.NewServer(opts...)

	demopb.RegisterFooServiceHTTPServer(srv, foo)
	demopb.RegisterBarServiceHTTPServer(srv, bar)
	return srv
}
