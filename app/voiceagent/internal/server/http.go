package server

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/app/voiceagent/internal/service"
	"store/pkg/confcenter"
	pkgConf "store/pkg/confcenter"
	"store/pkg/krathelper"
	"store/pkg/middlewares/encoder"
	helpers "store/pkg/sdk/helper"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c confcenter.Server, service *service.VoiceAgentService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.ErrorEncoder(encoder.ErrorEncoder),
		http.ResponseEncoder(encoder.ResponseEncoder),
		http.RequestDecoder(encoder.RequestDecoder),

		http.Middleware(
			krathelper.NormalizeAuthorization(krathelper.SecretSignKey),
			krathelper.FromCookie(pkgConf.AuthCookieName),
			selector.Server(
				jwt.Server(
					func(token *jwtv4.Token) (interface{}, error) {
						return []byte(krathelper.SecretSignKey), nil
					},
				)).Match(func(ctx context.Context, operation string) bool {
				// 可以根据需要排除不需要鉴权的接口
				return !helpers.InSlice(operation, []string{})
			}).Build(),
			recovery.Recovery(),
		),
		krathelper.CorsConfig(),
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
	voiceagent.RegisterVoiceAgentServiceHTTPServer(srv, service)
	voiceagent.RegisterMotivationServiceHTTPServer(srv, service)

	go func() {
		helpers.DeferFunc()
		service.StartConsumer()
	}()

	return srv
}
