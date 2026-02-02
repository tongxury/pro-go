package server

import (
	"context"
	creditpb "store/api/credit"
	toolkitpb "store/api/toolkit"
	ucpb "store/api/usercenter"
	"store/app/usercenter/internal/service"
	"store/pkg/confcenter"
	pkgConf "store/pkg/confcenter"
	"store/pkg/krathelper"
	"store/pkg/middlewares/encoder"
	helpers "store/pkg/sdk/helper"
	"store/pkg/sdk/helper/crond"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/robfig/cron/v3"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c confcenter.Server,
	user *service.UserService,
	auth *service.AuthService,
	toolkit *service.ToolkitService,
	credit *service.CreditService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.ErrorEncoder(encoder.ErrorEncoder),
		http.ResponseEncoder(encoder.ResponseEncoder),

		http.Middleware(
			krathelper.NormalizeAuthorization(krathelper.SecretSignKey),
			krathelper.FromCookie(pkgConf.AuthCookieName),
			selector.Server(
				jwt.Server(
					func(token *jwtv4.Token) (interface{}, error) {
						return []byte(krathelper.SecretSignKey), nil
					},
				)).Match(func(ctx context.Context, operation string) bool {

				return !helpers.InSlice(operation, []string{
					ucpb.AuthService_GetToken_FullMethodName,
					ucpb.AuthService_SendCode_FullMethodName,
					ucpb.UserService_GetUser_FullMethodName,
					ucpb.AuthService_GetAppleAuthToken_FullMethodName,
				})
			}).Build(),
			recovery.Recovery(),
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
	srv := http.NewServer(opts...)

	go func() {
		defer helpers.DeferFunc()

		cr := cron.New()

		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		}))

		cr.Run()
	}()

	ucpb.RegisterAuthServiceHTTPServer(srv, auth)
	ucpb.RegisterUserServiceHTTPServer(srv, user)

	creditpb.RegisterCreditServiceHTTPServer(srv, credit)
	toolkitpb.RegisterToolkitServiceHTTPServer(srv, toolkit)
	return srv
}
