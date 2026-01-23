package server

import (
	"context"
	trackerpb "store/api/aiagent"
	"store/app/aiagent/internal/service"
	"store/pkg/confcenter"
	pkgConf "store/pkg/confcenter"
	"store/pkg/krathelper"
	"store/pkg/middlewares/encoder"
	"store/pkg/middlewares/eventsource"
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
func NewHTTPServer(c confcenter.Server, service *service.TrackerService, logger log.Logger) *http.Server {
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

				return !helpers.InSlice(operation, []string{
					trackerpb.AIAgentService_ListPublicRecords_FullMethodName,
					trackerpb.AIAgentService_GetPublicRecord_FullMethodName,
					trackerpb.AIAgentService_Debug_FullMethodName,
					trackerpb.AIAgentService_ListItems_FullMethodName,
					trackerpb.AIAgentService_ListAccounts_FullMethodName,
					trackerpb.AIAgentService_SubmitSurvey_FullMethodName,
					trackerpb.AIAgentService_ListQuestions_FullMethodName,
					trackerpb.AIAgentService_ListSessions_FullMethodName,
					trackerpb.AIAgentService_GetAliOssSignedUrl_FullMethodName,
					trackerpb.AIAgentService_GetQiniuUploadToken_FullMethodName,
				})
			}).Build(),
			//logging.Server(logger),
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

	go func() {
		defer helpers.DeferFunc()

		//service.Data.Mongo.Settings.UpdateFieldsById(context.Background(), "1",
		//	bson.M{"xhsCookies": map[string]string{
		//		"xs":       "XYW_eyJzaWduU3ZuIjoiNTYiLCJzaWduVHlwZSI6IngyIiwiYXBwSWQiOiJ4aHMtcGMtd2ViIiwic2lnblZlcnNpb24iOiIxIiwicGF5bG9hZCI6ImQyMzQ3MzUwM2E4OWE2NzY3YjFiNTMzNTliOGNiMmY0MzM5NDU4NjMxZWMyMzljMTY4ZDQxNThiNWE3NTE4MGIwMzI5Yjc5Yzc1NWQ3NDUzOTk2NDdjZTc4MDMxNzcxNmVhMWQxYzU2NDljNzI2YjQ2NDMyYWVmZmFlNjU3NzBmMzEzMTkwNDllMDUyOGFhYWVkYzE3NDJkOWVmZDcyNjJlMzliY2UxNWQ3MmQ4NDEwZjRiY2NhOGQ3Y2I5NGIzM2YzZTVkMjllMjQwMTAwODYxNTVkYmQ4NjAwMWUzZDEyMmY2OGMyYTA0N2RmMWRjZGFlZTUwNWEwMTQ2NzQyMzZlZmU0MWI3YmI5NTVkYTA4NDQ5ODJiMDYyZWIwMmY0NzE0NmRjYWIyYjBhM2MwNmMxMGVkZDMzZmJmOTgzODg0MDAyYjA5MzBlNTQ5ZTE1ODMyMjQ4YjRjOGNkOTAzM2RkMzhmZTkyZjJjNTBjNTRjMzk4YjMyOGUwMjFiOGVlNTk0YjljYjhiNTJmOTQ2NjZlYTIyYjEwZTFhZDRjZjhhIn0=\n",
		//		"xscommon": "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\n",
		//		"cookie":   "abRequestId=ffa7ca92-1757-5b56-8a73-5570627c1d52; a1=194e40ecd54z2gvo7ea38d9i1fnhxdr5sqnssy5ha30000336056; webId=7486224fc1ddf358b401f55d4a34195f; gid=yj4d48ddiKMDyj4d48dSfAv224uJkhlWd0qYIx0fj3yiIMq8xCfF2M888qqK82K8JJDSijYi; x-user-id-creator.xiaohongshu.com=5cdf90a10000000011023c23; customerClientId=295944940634645; access-token-creator.xiaohongshu.com=customer.creator.AT-68c5174909416462787304625ty0jxzgyey8ntzs; galaxy_creator_session_id=krTT90MN2muz5EEgc3oDz1YPc41QrZk1q7BI; galaxy.creator.beaker.session.id=1744120765156098347520; webBuild=4.62.3; xsecappid=xhs-pc-web; unread={%22ub%22:%2267eb62ec000000001c029cb4%22%2C%22ue%22:%2267ecb1180000000009017b3a%22%2C%22uc%22:29}; web_session=0400698dabd3e9146c75e0783e3a4b213d896a; acw_tc=0ad5963417455718240348419e0acfcacae01ed0b093f9dbcc15869083549a; websectiga=984412fef754c018e472127b8effd174be8a5d51061c991aadd200c69a2801d6; sec_poison_id=ae08eb96-7c41-4317-a4a8-f581be18c18a; loadts=1745572703799",
		//	}},
		//)

		cr := cron.New()

		_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
			//service.UpdateCategory(context.Background())
			//service.UpdateItemId(context.Background())
		}))

		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	service.UpdateQuestionStatus(context.Background(), &empty.Empty{})
		//}))

		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	_ = service.AnswerQuestion(context.Background())
		//}))

		//_, _ = cr.AddJob("@every 1s", crond.NewJobWrapper(func() {
		//	_ = service.GenerateSurveyResult(context.Background())
		//}))

		cr.Start()
	}()

	srv.HandleFunc("/api/ag/v1/answer-chunks-stream", eventsource.Handle[*trackerpb.GetAnswerChunksStreamParams](
		func(ctx eventsource.Ctx, params *trackerpb.GetAnswerChunksStreamParams) (int, error) {

			code, err := service.GetAnswerChunksStream(ctx, params)
			if err != nil {
				//ctx.(eventsource.Ctx).Abort(code, err.Error())
				log.Errorw("GetAnswerChunksStream err", err)
				return code, err
			}
			return 0, nil
		}),
	)

	srv.HandleFunc("/api/ag/v1/analyse", eventsource.Handle[*trackerpb.CreateQuestionParams](
		func(ctx eventsource.Ctx, params *trackerpb.CreateQuestionParams) (int, error) {

			code, err := service.CreateQuestion(ctx, params)
			if err != nil {

				//ctx.(eventsource.Ctx).Abort(code, err.Error())
				log.Errorw("ChatCompletions err", err)
				return code, err
			}
			return 0, nil
		}),
	)

	srv.HandleFunc("/api/ag/v1/analysis", eventsource.Handle[*trackerpb.CreateQuestionParams](
		func(ctx eventsource.Ctx, params *trackerpb.CreateQuestionParams) (int, error) {

			log.Debugw("/api/ag/v1/analysis", "", "params", params)

			code, err := service.CreateQuestionV3(ctx, params)
			if err != nil {

				//ctx.(eventsource.Ctx).Abort(code, err.Error())
				log.Errorw("ChatCompletions err", err)
				return code, err
			}
			return 0, nil
		}),
	)

	srv.HandleFunc("/api/ag/v2/analysis", eventsource.Handle[*trackerpb.CreateQuestionParams](
		func(ctx eventsource.Ctx, params *trackerpb.CreateQuestionParams) (int, error) {

			code, err := service.CreateQuestionV5(ctx, params)
			if err != nil {

				//ctx.(eventsource.Ctx).Abort(code, err.Error())
				log.Errorw("ChatCompletions err", err)
				return code, err
			}
			return 0, nil
		}),
	)

	trackerpb.RegisterAIAgentServiceHTTPServer(srv, service)
	return srv
}
