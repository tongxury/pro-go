package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	adminpb "store/api/admin"
	userpb "store/api/user"
	"store/app/user/internal/service"
	"store/pkg/confcenter"
	pkgConf "store/pkg/confcenter"
	"store/pkg/krathelper"
	"store/pkg/middlewares/encoder"
	"store/pkg/middlewares/prometircs"
	helpers "store/pkg/sdk/helper"
	"strings"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c confcenter.Server,
	user *service.UserService,
	//member *service.MemberService,
	auth *service.AuthService,
	admin *service.AdminService,
	logger log.Logger) *http.Server {

	var opts = []http.ServerOption{
		http.ErrorEncoder(encoder.ErrorEncoder),
		http.ResponseEncoder(encoder.ResponseEncoder),
		http.RequestDecoder(encoder.RequestDecoder),
		http.Middleware(
			//logging.Server(logger),
			krathelper.NormalizeContentType(),
			krathelper.InternalAPIAuth(),
			krathelper.NormalizeAuthorization(krathelper.SecretSignKey),
			krathelper.FromCookie(pkgConf.AuthCookieName),
			selector.Server(
				jwt.Server(
					func(token *jwtv4.Token) (interface{}, error) {
						return []byte(krathelper.SecretSignKey), nil
					},
				)).Match(func(ctx context.Context, operation string) bool {

				req, ok := http.RequestFromServerContext(ctx)
				if ok {
					if strings.Contains(req.RequestURI, "/internal-api/") {
						return false
					}
				}

				return !helpers.InSlice(operation, []string{
					userpb.Auth_GetPhoneAuthCode_FullMethodName,
					userpb.Auth_GetPhoneAuthToken_FullMethodName,
					userpb.Auth_GetAuthCode_FullMethodName,
					userpb.Auth_GetAuthToken_FullMethodName,
					userpb.Auth_GetAuthTokenV2_FullMethodName,

					userpb.Auth_AddRegisterUser_FullMethodName,
					userpb.Auth_GetVerifyCode_FullMethodName,
					userpb.Auth_SubmitVerifyCode_FullMethodName,
					userpb.Auth_GetAuthTokenByEmail_FullMethodName,
					userpb.Auth_ResetAuthPassword_FullMethodName,

					userpb.Member_GetMemberMetadata_FullMethodName,
					userpb.User_GetSettings_FullMethodName,
					adminpb.Admin_GetNotifications_FullMethodName,

					userpb.User_GetUser_FullMethodName,
					userpb.User_CallbackPayment_FullMethodName,
					userpb.Auth_GetPhoneAuthToken_FullMethodName,
					userpb.Auth_GetPhoneAuthToken_FullMethodName,
					userpb.Auth_GetEmailAuthCode_FullMethodName,
					userpb.Auth_GetEmailAuthToken_FullMethodName,
					userpb.Auth_GetAppleAuthToken_FullMethodName,
					userpb.Auth_GetWxAuthToken_FullMethodName,

					userpb.User_GetAppSettings_FullMethodName,
					userpb.User_GetAppVersion_FullMethodName,

					userpb.User_GetServerMeta_FullMethodName,
				})
			}).Build(),
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

	httpSrv.Handle("/metrics", promhttp.Handler())
	httpSrv.HandlePrefix("/q/", openapiv2.NewHandler())

	httpSrv.HandleFunc("/api/v1/payment/callback", func(w http.ResponseWriter, r *http.Request) {

		log.Debugw("r.RequestURI", r.RequestURI)

		//pid=1215&trade_no=2025041809413132472&out_trade_no=12345678&type=wxpay&name=product&money=0.5&trade_status=TRADE_SUCCESS&api_trade_no=2025041809415239739&timestamp=1744940539&sign_type=RSA&sign=dT%2B7ldn0PNNatyDoMm8WkMkkiGWbbxaKrrZFNcqrTKm45EYHH1xFw15IfnInhVpSsCnz%2F7qmCfD5huyojCVwfSddMvq1yN910DNsE7mm3BYK%2B4u2ZdpDSJDSqyHBSu%2FjEekO00r3TeZ1GfXY%2FKda0f6lyXL%2FktX0b7v1mNUMgiRsYVPyujLHBfMowoEWiogSUl4WF%2BtT2TfUS1hnhYnadV51fMyoz2d7xiHF7%2FdFNVA99fkk0sWIbuXxeorpzGBVnYUcLhSXnuTF022D7jx00btt2Gm6rFI7XTnq69Psl%2F%2FmYR8ybv9wLyhO3NGO46zLL3H%2FNXFW54kGRf2S%2Fqm8qA%3D%3D

		all, _ := io.ReadAll(r.Body)

		log.Debugw("r.Body", string(all))

	})

	//go func() {
	//	for {
	//		time.Sleep(time.Second)
	//
	//		log.Errorw("test", "every 1s")
	//	}
	//
	//}()

	userpb.RegisterUserHTTPServer(httpSrv, user)
	userpb.RegisterAuthHTTPServer(httpSrv, auth)
	//userpb.RegisterMemberHTTPServer(httpSrv, member)
	adminpb.RegisterAdminHTTPServer(httpSrv, admin)
	return httpSrv
}
