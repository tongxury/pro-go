package tracer

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// 设置全局trace
func SetTracerProvider(name string, conf *Config) error {
	// 创建 Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.Endpoint)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		tracesdk.WithBatcher(exp),
		// 在资源中记录有关此应用程序的信息
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(name),
			attribute.String("attr1", "attr1v"),
			attribute.Float64("attr2", 312.23),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

//func Server(serverName string) middleware.Middleware {
//	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("")))
//	if err != nil {
//		panic(err)
//	}
//	tp := tracesdk.NewTracerProvider(
//		tracesdk.WithBatcher(exp),
//		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
//		tracesdk.WithResource(resource.NewSchemaless(
//			semconv.ServiceNameKey.String(serverName),
//			attribute.String("exporter", "jaeger"),
//			attribute.Float64("float", 312.23),
//		)),
//	)
//	otel.SetTracerProvider(tp)
//	return tracing.Server(
//	//tracing.WithTracerProvider(tp),
//	)
//}
