package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	nethttp "net/http"
	paymentpb "store/api/payment"
	"store/app/payment/internal/conf"
	"store/app/payment/internal/service"
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
	conf confcenter.Config[bizconf.BizConfig],
	payment *service.PaymentService,
	logger log.Logger) *http.Server {

	var opts = []http.ServerOption{

		http.ErrorEncoder(encoder.ErrorEncoder),
		http.ResponseEncoder(encoder.ResponseEncoder),
		krathelper.CorsConfig(),
		http.Middleware(
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
					paymentpb.Payment_CallbackPayment_FullMethodName,
					paymentpb.Payment_CallbackByAirWallex_FullMethodName,
					paymentpb.Payment_CallbackByAppleV2_FullMethodName,
					paymentpb.Payment_MockCredit_FullMethodName,
				})
			}).Build(),
			krathelper.NormalizeContentType(),
			logging.Server(logger),
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					return nil
				}),
			),
			validate.Validator(),
			ratelimit.Server(),
			prometircs.Metrics(),
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

	httpSrv := http.NewServer(opts...)

	httpSrv.Handle("/metrics", promhttp.Handler())
	httpSrv.HandlePrefix("/q/", openapiv2.NewHandler())

	//cr := cron.New()
	//
	//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
	//	payment.GeneratePayment(context.Background(), &paymentpb.GeneratePaymentParams{})
	//}))
	//
	//cr.Start()

	//httpSrv.HandleFunc(bizconf.StripeCallbackURL_ProdV1, func(w http.ResponseWriter, req *http.Request) {
	//
	//	ctx := req.Context()
	//
	//	req.Body = nethttp.MaxBytesReader(w, req.Body, int64(65536))
	//
	//	body, err := io.ReadAll(req.Body)
	//	if err != nil {
	//		w.WriteHeader(nethttp.StatusBadRequest)
	//		return
	//	}
	//
	//	err = payment.OnEvent(ctx, body, req.Header.Get("Stripe-Signature"), bizconf.StripeCallbackURL_ProdV1)
	//	if err != nil {
	//		w.WriteHeader(nethttp.StatusBadRequest) // Return a 400 error on a bad signature
	//		log.Errorw("stipe-callback CallbackPayment err", err)
	//		return
	//	}
	//
	//	w.WriteHeader(200)
	//})

	httpSrv.HandleFunc(bizconf.StripeCallbackURL_Prod, func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		req.Body = nethttp.MaxBytesReader(w, req.Body, int64(65536))

		body, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)
			return
		}

		err = payment.OnEvent(ctx, body, req.Header.Get("Stripe-Signature"), bizconf.StripeCallbackURL_Prod)
		if err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)
			log.Errorw("stipe-callback CallbackPayment err", err)
			return
		}

		w.WriteHeader(200)
	})

	httpSrv.HandleFunc(bizconf.StripeCallbackURL_Test, func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		req.Body = nethttp.MaxBytesReader(w, req.Body, int64(65536))

		body, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)
			return
		}

		err = payment.OnEvent(ctx, body, req.Header.Get("Stripe-Signature"), bizconf.StripeCallbackURL_Test)
		if err != nil {
			w.WriteHeader(nethttp.StatusBadRequest) // Return a 400 error on a bad signature
			log.Errorw("stipe-callback CallbackPayment err", err)
			return
		}

		w.WriteHeader(200)
	})

	paymentpb.RegisterPaymentHTTPServer(httpSrv, payment)
	return httpSrv
}
