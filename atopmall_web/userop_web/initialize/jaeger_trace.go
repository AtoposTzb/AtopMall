package initialize

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"atopmall_web/userop_web/global"
)

// TracerInit 全局追踪器初始化
// 返回 *TracerProvider，用于程序退出时优雅刷新剩余数据
func JaegerTracerInit() (*sdktrace.TracerProvider, error) {
	// 从配置中获取 Jaeger 配置
	jaegerConfig := global.ServerConfig.JaegerInfo
	strConfig := fmt.Sprintf("%s:%d", jaegerConfig.Host, jaegerConfig.Port)
	fmt.Println("检查Jaeger配置+++++++++", strConfig)
	// 1. 创建 OTLP 导出器：负责把 Span 数据发送到 Jaeger
	exp, err := otlptracegrpc.New(
		context.Background(),
		// Jaeger 的 OTLP gRPC 地址
		otlptracegrpc.WithEndpoint(strConfig),
		// 开发环境无 TLS 加密，生产环境建议开启
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// 2. 配置资源信息：标记这条链路属于哪个服务
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			// 服务名，Jaeger UI 中会显示，必须设置
			semconv.ServiceNameKey.String(jaegerConfig.Name),
			// 可附加自定义属性，比如环境、版本号
			// semconv.DeploymentEnvironmentKey.String("dev"),
			// semconv.ServiceVersionKey.String("v1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	// 3. 创建 TracerProvider：追踪器的核心管理器
	tp := sdktrace.NewTracerProvider(
		// 批量处理器：攒一批 Span 再统一上报，减少网络请求
		sdktrace.WithBatcher(exp),
		// 绑定服务资源信息
		sdktrace.WithResource(res),
		// 默认全量采样，生产环境可配置采样率
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// 4. 设置为全局默认追踪器，后续代码直接用 otel.Tracer() 即可获取
	//TracerProvider：管理 Span 的"生老病死"
	otel.SetTracerProvider(tp)
	// 5. 设置全局默认的 Propagator，后续代码直接用 otel.GetPropagator() 即可获取
	//Propagator：管理 trace context 的"跨进程传递"
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}
