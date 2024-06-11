from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.api import annotations_pb2 as _annotations_pb2
from protoc_gen_openapiv2.options import annotations_pb2 as _annotations_pb2_1
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PrepareDataReq(_message.Message):
    __slots__ = ("sources",)
    SOURCES_FIELD_NUMBER: _ClassVar[int]
    sources: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, sources: _Optional[_Iterable[str]] = ...) -> None: ...

class PredictReq(_message.Message):
    __slots__ = ("ts", "months_count", "segment")
    TS_FIELD_NUMBER: _ClassVar[int]
    MONTHS_COUNT_FIELD_NUMBER: _ClassVar[int]
    SEGMENT_FIELD_NUMBER: _ClassVar[int]
    ts: _timestamp_pb2.Timestamp
    months_count: int
    segment: str
    def __init__(self, ts: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., months_count: _Optional[int] = ..., segment: _Optional[str] = ...) -> None: ...

class PredictResp(_message.Message):
    __slots__ = ("predicts",)
    PREDICTS_FIELD_NUMBER: _ClassVar[int]
    predicts: _containers.RepeatedCompositeFieldContainer[Predict]
    def __init__(self, predicts: _Optional[_Iterable[_Union[Predict, _Mapping]]] = ...) -> None: ...

class Predict(_message.Message):
    __slots__ = ("ts", "price")
    TS_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    ts: _timestamp_pb2.Timestamp
    price: float
    def __init__(self, ts: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., price: _Optional[float] = ...) -> None: ...

class ClientIdentifier(_message.Message):
    __slots__ = ("value",)
    VALUE_FIELD_NUMBER: _ClassVar[int]
    value: str
    def __init__(self, value: _Optional[str] = ...) -> None: ...

class UniqueCodesResp(_message.Message):
    __slots__ = ("codes",)
    CODES_FIELD_NUMBER: _ClassVar[int]
    codes: _containers.RepeatedCompositeFieldContainer[UniqueCode]
    def __init__(self, codes: _Optional[_Iterable[_Union[UniqueCode, _Mapping]]] = ...) -> None: ...

class UniqueCode(_message.Message):
    __slots__ = ("segment", "regular")
    SEGMENT_FIELD_NUMBER: _ClassVar[int]
    REGULAR_FIELD_NUMBER: _ClassVar[int]
    segment: str
    regular: bool
    def __init__(self, segment: _Optional[str] = ..., regular: bool = ...) -> None: ...
