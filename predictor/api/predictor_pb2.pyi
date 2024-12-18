from google.api import annotations_pb2 as _annotations_pb2
from google.protobuf import empty_pb2 as _empty_pb2
from protoc_gen_openapiv2.options import annotations_pb2 as _annotations_pb2_1
from api import prompter_pb2 as _prompter_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Source(_message.Message):
    __slots__ = ("name", "path")
    NAME_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    name: str
    path: str
    def __init__(self, name: _Optional[str] = ..., path: _Optional[str] = ...) -> None: ...

class PrepareDataReq(_message.Message):
    __slots__ = ("sources", "organization")
    SOURCES_FIELD_NUMBER: _ClassVar[int]
    ORGANIZATION_FIELD_NUMBER: _ClassVar[int]
    sources: _containers.RepeatedCompositeFieldContainer[Source]
    organization: str
    def __init__(self, sources: _Optional[_Iterable[_Union[Source, _Mapping]]] = ..., organization: _Optional[str] = ...) -> None: ...

class PredictReq(_message.Message):
    __slots__ = ("type", "product", "period", "organization")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PRODUCT_FIELD_NUMBER: _ClassVar[int]
    PERIOD_FIELD_NUMBER: _ClassVar[int]
    ORGANIZATION_FIELD_NUMBER: _ClassVar[int]
    type: _prompter_pb2.QueryType
    product: str
    period: str
    organization: str
    def __init__(self, type: _Optional[_Union[_prompter_pb2.QueryType, str]] = ..., product: _Optional[str] = ..., period: _Optional[str] = ..., organization: _Optional[str] = ...) -> None: ...

class PredictResp(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...

class UniqueCodesReq(_message.Message):
    __slots__ = ("organization",)
    ORGANIZATION_FIELD_NUMBER: _ClassVar[int]
    organization: str
    def __init__(self, organization: _Optional[str] = ...) -> None: ...

class UniqueCodesResp(_message.Message):
    __slots__ = ("codes",)
    CODES_FIELD_NUMBER: _ClassVar[int]
    codes: _containers.RepeatedCompositeFieldContainer[UniqueCode]
    def __init__(self, codes: _Optional[_Iterable[_Union[UniqueCode, _Mapping]]] = ...) -> None: ...

class UniqueCode(_message.Message):
    __slots__ = ("segment", "name", "regular")
    SEGMENT_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    REGULAR_FIELD_NUMBER: _ClassVar[int]
    segment: str
    name: str
    regular: bool
    def __init__(self, segment: _Optional[str] = ..., name: _Optional[str] = ..., regular: bool = ...) -> None: ...
