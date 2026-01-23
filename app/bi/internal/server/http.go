package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	bipb "store/api/bi"
	"store/app/bi/internal/service"
	"store/pkg/confcenter"
	pkgConf "store/pkg/confcenter"
	"store/pkg/krathelper"
	"store/pkg/middlewares/encoder"
	"store/pkg/middlewares/prometircs"
	helpers "store/pkg/sdk/helper"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c confcenter.Server,
	bi *service.BiService,
	logger log.Logger) *http.Server {

	var opts = []http.ServerOption{
		http.ErrorEncoder(encoder.ErrorEncoder),
		http.ResponseEncoder(encoder.ResponseEncoder),
		http.Middleware(
			//logging.Server(logger),
			krathelper.NormalizeAuthorization(krathelper.SecretSignKey),
			krathelper.FromCookie(pkgConf.AuthCookieName),
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					return nil
				}),
			),
			validate.Validator(),
			ratelimit.Server(),
			prometircs.Metrics(),
		),
		krathelper.DefaultCorsConfig,
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

	httpSrv := http.NewServer(opts...)

	go func() {
		helpers.DeferFunc()
		bi.StartConsumer()
	}()

	httpSrv.Handle("/metrics", promhttp.Handler())
	httpSrv.HandlePrefix("/q/", openapiv2.NewHandler())

	bipb.RegisterBiHTTPServer(httpSrv, bi)
	return httpSrv
}
