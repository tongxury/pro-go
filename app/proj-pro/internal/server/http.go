package server

import (
	"context"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/service"
	"store/pkg/clients/cronexecutor"
	"store/pkg/confcenter"
	pkgConf "store/pkg/confcenter"
	"store/pkg/krathelper"
	"store/pkg/middlewares/encoder"
	helpers "store/pkg/sdk/helper"
	"store/pkg/sdk/helper/crond"
	"strings"

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
	service *service.ProjService,
	workflow *service.WorkFlowService,
	session *service.SessionService, logger log.Logger) *http.Server {
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

				req, ok := http.RequestFromServerContext(ctx)
				if ok {
					if strings.Contains(req.RequestURI, "/internal-api/") {
						return false
					}
				}

				return !helpers.InSlice(operation, []string{
					projpb.OperationProjServiceCallbackByDouyin,
					projpb.OperationProjServiceWebhookByDouyin,
					projpb.OperationProjServiceWebhookByWxpay,
					projpb.OperationProjServiceGetBindingAccount,

					projpb.OperationProjServiceGetTask,
					projpb.OperationProjServiceExportScript,

					projpb.OperationProjServiceGetPhoneAuthCode,
					projpb.OperationProjServiceGetPhoneAuthToken,
					projpb.OperationProjServiceGetWxAuthToken,

					projpb.ProjProService_CallbackByVolcengine_FullMethodName,
					projpb.OperationProjProServiceListQualityAssets,

					projpb.WorkflowService_CreateWorkflow_FullMethodName,
					projpb.WorkflowService_GetWorkflow_FullMethodName,
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

	cc := cronexecutor.NewClient()

	if err := cc.Register("@every 5s", workflow, cronexecutor.RegisterOptions{Concurrency: 100}); err != nil {
		panic(err)
	}

	if err := cc.Register("@every 1m",
		workflow.SyncRunningWorkflowIdsExecutor(),
		cronexecutor.RegisterOptions{Concurrency: 1}); err != nil {
		panic(err)
	}

	if err := cc.Register("@every 10m",
		workflow.PauseTimeoutWorkflowIdsExecutor(),
		cronexecutor.RegisterOptions{Concurrency: 1}); err != nil {
		panic(err)
	}

	cc.Start()

	go func() {
		defer helpers.DeferFunc()

		cr := cron.New()

		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
			service.AnalyzeChances()
		}))

		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
			service.SegmentResult()
		}))

		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
			service.ExtractSegment()
		}))
		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
			service.AnalyzeSegment()
		}))

		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	service.SegmentAnalyze()
		//}))

		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
			service.AnalyzeSubSegment()
		}))

		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	service.GenerateText()
		//}))

		//// session
		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	session.GenerateSubtitle()
		//}))
		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	session.GenerateNewKeyFrames()
		//}))
		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	session.Remix()
		//}))

		cr.Run()
	}()

	projpb.RegisterProjProServiceHTTPServer(srv, service)
	projpb.RegisterSessionServiceHTTPServer(srv, session)
	projpb.RegisterWorkflowServiceHTTPServer(srv, workflow)
	return srv
}
