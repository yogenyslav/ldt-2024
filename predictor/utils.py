from datetime import datetime
from functools import wraps
import math

from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor


def get_tracer(endpoint="localhost:4317"):
    if not isinstance(trace.get_tracer_provider(), TracerProvider):
        resource = Resource(attributes={SERVICE_NAME: "predictor"})

        trace.set_tracer_provider(TracerProvider(resource=resource))
        otlp_exporter = OTLPSpanExporter(endpoint=endpoint, insecure=True)
        span_processor = BatchSpanProcessor(otlp_exporter)
        trace.get_tracer_provider().add_span_processor(span_processor)

    tracer = trace.get_tracer("predictor.traces")

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


def convert_to_datetime(iso_str):
    iso_str = iso_str.replace("Z", "+00:00")
    dt = datetime.strptime(iso_str, "%Y-%m-%dT%H:%M:%S.%f%z")
    return dt


def convert_datetime_to_str(obj):
    if isinstance(obj, dict):
        return {k: convert_datetime_to_str(v) for k, v in obj.items()}
    elif isinstance(obj, list):
        return [convert_datetime_to_str(i) for i in obj]
    elif isinstance(obj, datetime):
        return obj.strftime("%Y-%m-%d")
    else:
        return obj

def convert_float_nan_to_none(obj):
    if isinstance(obj, dict):
        return {k: convert_float_nan_to_none(v) for k, v in obj.items()}
    elif isinstance(obj, list):
        return [convert_float_nan_to_none(i) for i in obj]
    elif isinstance(obj, float):
        if math.isnan(obj):
            return None
        return obj
    else:
        return obj

def mdb_instert_many(data, mdb, collection_name, drop_exists=True):
    if collection_name in mdb.list_collection_names() and drop_exists:
        mdb[collection_name].drop()
        
    collection = mdb[collection_name]
    collection.insert_many(data)