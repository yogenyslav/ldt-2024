from google.api import annotations_pb2 as _annotations_pb2
from protoc_gen_openapiv2.options import annotations_pb2 as _annotations_pb2_1
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNDEFINED: _ClassVar[QueryType]
    PREDICTION: _ClassVar[QueryType]
    STOCK: _ClassVar[QueryType]

UNDEFINED: QueryType
PREDICTION: QueryType
STOCK: QueryType

class ExtractReq(_message.Message):
    __slots__ = ("prompt",)
    PROMPT_FIELD_NUMBER: _ClassVar[int]
    prompt: str
    def __init__(self, prompt: _Optional[str] = ...) -> None: ...

class ExtractedPrompt(_message.Message):
    __slots__ = ("type", "product", "period")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PRODUCT_FIELD_NUMBER: _ClassVar[int]
    PERIOD_FIELD_NUMBER: _ClassVar[int]
    type: QueryType
    product: str
    period: str
    def __init__(
        self,
        type: _Optional[_Union[QueryType, str]] = ...,
        product: _Optional[str] = ...,
        period: _Optional[str] = ...,
    ) -> None: ...
