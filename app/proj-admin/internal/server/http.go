package server

import (
	projpb "store/api/proj"
	"store/app/proj-admin/internal/service"
	"store/pkg/confcenter"
	"store/pkg/krathelper"
	"store/pkg/middlewares/encoder"
	helpers "store/pkg/sdk/helper"
	"store/pkg/sdk/helper/crond"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/robfig/cron/v3"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c confcenter.Server, service *service.ProjAdminService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.ErrorEncoder(encoder.ErrorEncoder),
		http.ResponseEncoder(encoder.ResponseEncoder),
		http.Middleware(
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

		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	service.Mi6()
		//}))

		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
			service.Segment()
		}))
		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	service.SegmentResult()
		//}))
		//
		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	service.ExtractSegment()
		//}))
		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	service.AnalyzeSegment()
		//}))
		//
		////_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		////	service.SegmentAnalyze()
		////}))
		//
		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	service.AnalyzeSubSegment()
		//}))

		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
			service.SyncTemplateSegmentsToVikingDB()
		}))
		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
			service.SyncTemplatesToVikingDB()
		}))

		cr.Run()
	}()

	projpb.RegisterProjAdminServiceHTTPServer(srv, service)
	return srv
}
