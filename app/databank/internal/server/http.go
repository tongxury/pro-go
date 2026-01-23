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
	databankpb "store/api/databank"
	"store/app/databank/internal/service"
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
	databank *service.DatabankService,
	logger log.Logger) *http.Server {

	var opts = []http.ServerOption{
		http.ErrorEncoder(encoder.ErrorEncoder),
		http.ResponseEncoder(encoder.ResponseEncoder),
		http.RequestDecoder(encoder.RequestDecoder),
		http.Middleware(
			logging.Server(logger),
			//kratosutil.NormalizeContentType(),
			krathelper.NormalizeAuthorization(krathelper.SecretSignKey),
			krathelper.InternalAPIAuth(),
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

				return !helpers.InSlice(operation, []string{})
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
	//
	//httpSrv.HandleFunc("/api/databank/files", func(writer nethttp.ResponseWriter, request *nethttp.Request) {
	//
	//	ctx := request.Context()
	//
	//	log.Debugw("userID", userID, "err", err)
	//
	//	err = request.ParseMultipartForm(32 << 20)
	//	if err != nil {
	//		writer.WriteHeader(nethttp.StatusBadRequest) // Return a 400 error on a bad signature
	//		log.Errorw("ParseMultipartForm err", err)
	//		return
	//	}
	//	files := request.MultipartForm.File["files"]
	//
	//	if len(files) == 0 {
	//		writer.WriteHeader(nethttp.StatusBadRequest)
	//		return
	//	}
	//
	//	var assets []*typepb.File
	//	for _, x := range files {
	//
	//		fileBody, _ := x.Open()
	//
	//		fileBytes, err := io.ReadAll(fileBody)
	//		if err != nil {
	//			return
	//		}
	//
	//		assets = append(assets, &typepb.File{
	//			Body: fileBytes,
	//			Name: x.Filename,
	//			Field:  helpers.MD5(fileBytes),
	//		})
	//	}
	//
	//	uploadedFiles, err := databank.AddFiles(ctx, &databankpb.AddFilesParams{Files: assets})
	//	if err != nil {
	//		writer.WriteHeader(nethttp.StatusInternalServerError) // Return a 400 error on a bad signature
	//		log.Errorw("AddFiles err", err)
	//		return
	//	}
	//
	//	encoder.WriteResponse(writer, request, nethttp.StatusOK, uploadedFiles)
	//})

	databankpb.RegisterDatabankHTTPServer(httpSrv, databank)
	return httpSrv
}
