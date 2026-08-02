from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor, SpanExportResult
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource, SERVICE_NAME

# 健康检查的 span 名称后缀，Consul 每5秒触发一次，会产生大量无意义 trace
_HEALTH_CHECK_SPAN_NAME = '/grpc.health.v1.Health/Check'


class _FilterHealthExporter:
    """包装 OTLP 导出器，过滤掉健康检查产生的 span，避免 Jaeger UI 被淹没"""

    def __init__(self, real_exporter):
        self._exporter = real_exporter

    def export(self, spans):
        # print(f"[DEBUG] _FilterHealthExporter.export called, span count={len(spans)}")
        # for s in spans:
        #     print(f"  span name={s.name}, service={s.resource.attributes.get('service.name', 'N/A')}")
        filtered = [s for s in spans if s.name != _HEALTH_CHECK_SPAN_NAME]
        if not filtered:
            return SpanExportResult.SUCCESS
        return self._exporter.export(filtered)

    def shutdown(self):
        self._exporter.shutdown()

    def force_flush(self, timeout_millis=None):
        self._exporter.force_flush(timeout_millis)


def init_tracer(service_name: str, jaeger_endpoint: str):
    """
    公共链路初始化
    :param service_name: 当前微服务名称
    :param jaeger_endpoint: jaeger otlp上报地址，格式 ip:4317
    """
    # 构建OTLP导出器
    raw_exporter = OTLPSpanExporter(
        endpoint=jaeger_endpoint,
        insecure=True
    )
    # 过滤健康检查 span，避免 Jaeger 被 Consul 健康检查淹没
    exporter = _FilterHealthExporter(raw_exporter)

    # 定义服务标识资源
    resource = Resource(attributes={
        SERVICE_NAME: service_name
    })

    # 构建追踪器
    provider = TracerProvider(resource=resource)
    provider.add_span_processor(BatchSpanProcessor(exporter))
    trace.set_tracer_provider(provider)