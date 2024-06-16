
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource, SERVICE_NAME
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from functools import wraps

def get_tracer(endpoint="localhost:4317"):
    resource = Resource(attributes={
    SERVICE_NAME: "predictor"
    })

    trace.set_tracer_provider(TracerProvider(resource=resource))
    otlp_exporter = OTLPSpanExporter(
        endpoint=endpoint,  
        insecure=True
    )
    span_processor = BatchSpanProcessor(otlp_exporter)
    trace.get_tracer_provider().add_span_processor(span_processor)
    tracer = trace.get_tracer('predictor.traces')

    return tracer

def trace_function(tracer, span_name=None):
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            name = span_name or func.name
            with tracer.start_as_current_span(name):
                return func(*args, **kwargs)
        return wrapper
    return decorator