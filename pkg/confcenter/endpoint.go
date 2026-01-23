package confcenter

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"store/pkg/krathelper"
	"store/pkg/middlewares/tracer"
	"store/pkg/sdk/conv"
)

type Config[Biz any] struct {
	Meta      Meta
	Server    Server
	Biz       Biz
	Database  Database
	Component Component
	Logger    log.Logger
}

func get[T any](client *NacosClient, dataId, appName string) (*T, error) {
	source, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId, Group: appName,
	})
	if err != nil {
		return nil, fmt.Errorf("%s, %s, %s", appName, dataId, err.Error())
	}

	var t T

	err = conv.Y2S([]byte(source), &t)
	if err != nil {
		return nil, fmt.Errorf("%s, %s, %s", appName, dataId, err.Error())
	}

	log.Infow("dataId", dataId, "appName", appName, "source", source, "t", t)

	return &t, nil
}

func GetMeta() (*Meta, error) {
	c := config.New(config.WithSource(env.NewSource()))
	if err := c.Load(); err != nil {
		return nil, err
	}

	namespace, err := c.Value("POD_NAMESPACE").String()
	if err != nil {
		return nil, err
	}
	appname, err := c.Value("POD_NAME").String()
	if err != nil {
		return nil, err
	}

	return &Meta{
		Namespace: namespace,
		Appname:   appname,
	}, nil
}

// todo 定制化严重

func GetConfig[Biz any](filepath ...string) (*Config[Biz], error) {

	if len(filepath) > 0 && filepath[0] != "" {
		var target Config[Biz]
		err := krathelper.ScanFile(filepath[0], &target)
		if err != nil {
			return nil, err
		}

		return &target, nil
	}

	c := config.New(config.WithSource(env.NewSource()))
	if err := c.Load(); err != nil {
		return nil, err
	}

	namespace, err := c.Value("POD_NAMESPACE").String()
	if err != nil {
		return nil, err
	}
	appname, err := c.Value("POD_NAME").String()
	if err != nil {
		return nil, err
	}

	log.Infow("namespace", namespace, "appname", appname)

	configClient, err := NewClient(namespace)
	if err != nil {
		return nil, err
	}

	// server
	server, err := get[Server](configClient, "server.yaml", "public")
	if err != nil {
		return nil, err
	}
	//level, err := zapcore.ParseLevel(helper.OrString(server.Log.Level, "info"))
	//if err != nil {
	//	return nil, err
	//}

	//encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	//encoderConfig.LineEnding = "\n\n\n"
	//encoderConfig.SkipLineEnding = false

	//// log hook
	//logHookCache := cache.New(time.Minute, time.Minute)
	////feishuClient := feishuo.NewClient(consts.FEISHU_ENDPOINT)
	//logHook := func(entry zapcore.Entry) error {
	//
	//	if server.Log.Alarm.Threshold <= 0 {
	//		return nil
	//	}
	//
	//	if entry.Level == zap.ErrorLevel {
	//
	//		cacheKey := namespace + appname
	//
	//		var times = 1
	//
	//		_, found := logHookCache.Get(cacheKey)
	//		if found {
	//			cur, err := logHookCache.IncrementInt(cacheKey, 1)
	//			if err != nil {
	//				return err
	//			}
	//			times = cur
	//
	//		} else {
	//			logHookCache.Set(cacheKey, 1, time.Minute)
	//		}
	//
	//		if times > server.Log.Alarm.Threshold {
	//			//log.Warnf(fmt.Sprintf("【%s】【%s】服务错误日志每分钟频频次 %d 次, 超过报警阈值(%d次)", namespace, appname, times, server.Log.Alarm.Threshold))
	//			//log.Warn(entry.Stack, entry.Message)
	//
	//			//err := feishuClient.SendRichText(context.Background(),
	//			//	feishuo.RichText{
	//			//		Key:   "error",
	//			//		Title: "服务错误日志报警",
	//			//		Content: []feishuo.RichTextParagraph{
	//			//			[]feishuo.RichTextContent{
	//			//				{
	//			//					Tag:  "text",
	//			//					Text: fmt.Sprintf("【%s】【%s】服务错误日志每分钟频次 %d次, 超过报警阈值(%d次)", namespace, appname, times, server.Log.Alarm.Threshold),
	//			//				},
	//			//				{
	//			//					Tag:  "a",
	//			//					Text: "请尽快查看",
	//			//					Href: fmt.Sprintf(consts.K8S_ADMIN_URL_FORMAT, namespace, appname),
	//			//				},
	//			//			},
	//			//			[]feishuo.RichTextContent{
	//			//				{
	//			//					Tag:  "text",
	//			//					Text: conv.S2J(entry),
	//			//				},
	//			//			},
	//			//			[]feishuo.RichTextContent{
	//			//				{
	//			//					Tag:  "text",
	//			//					Text: entry.Stack,
	//			//				},
	//			//			},
	//			//		},
	//			//	}, feishuo.Option{IntervalSeconds: server.Log.Alarm.Interval})
	//			//if err != nil {
	//			//	return err
	//			//}
	//		}
	//	}
	//	return nil
	//}

	//logger := kratoszap.NewLogger(
	//	zap.New(
	//		zapcore.NewCore(
	//			zapcore.NewJSONEncoder(encoderConfig),
	//			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
	//			level,
	//		),
	//		zap.AddCaller(),
	//		zap.AddCallerSkip(3),
	//		//zap.Hooks(logHook),
	//	).With(
	//		zap.String("service.name", appname),
	//		//zap.String("service.version", Version),
	//		//zap.String("trace.id", tracing.TraceID()(context.Background())),
	//		//zap.Any("span.id", tracing.SpanID()),
	//	),
	//)

	// 链路追踪
	if server.Tracer.Endpoint != "" {
		if err := tracer.SetTracerProvider(appname, &server.Tracer); err != nil {
			panic(err)
		}
	}

	// database
	//databaseConfigGroup, _ := c.Value("DATABASE_CONFIG_GROUP").String()
	//databaseConfigGroup = helper.OrString(databaseConfigGroup, "public")

	db, err := get[Database](configClient, "database.yaml", "public")
	if err != nil {
		return nil, err
	}
	// component
	cp, err := get[Component](configClient, "component.yaml", "public")
	if err != nil {
		return nil, err
	}
	// biz
	biz, err := get[Biz](configClient, "biz.yaml", appname)
	if err != nil {
		return nil, err
	}

	// 最后再设置logger 保留 config的初始格式 便于查看

	rsp := &Config[Biz]{
		Meta: Meta{
			Namespace: namespace,
			Appname:   appname,
		},
		Server:    *server,
		Database:  *db,
		Component: *cp,
		Biz:       *biz,
		Logger:    log.DefaultLogger,
	}

	// 设置logger 后会出现 go 协程无法打印日志的问题
	//log.SetLogger(logger)

	log.Infow("configs", rsp)

	return rsp, nil
}
