# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: api/prompter.proto
# Protobuf Python Version: 5.26.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder

# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from protoc_gen_openapiv2.options import (
    annotations_pb2 as protoc__gen__openapiv2_dot_options_dot_annotations__pb2,
)


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x12\x61pi/prompter.proto\x12\x03\x61pi\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\"\x1b\n\tStreamReq\x12\x0e\n\x06prompt\x18\x01 \x01(\x0c\"\x1b\n\nStreamResp\x12\r\n\x05\x63hunk\x18\x01 \x01(\t\"\xb3\x02\n\nExtractReq\x12\x0e\n\x06prompt\x18\x01 \x01(\t:\x94\x02\x92\x41\x90\x02\n\xb1\x01*9\xd0\x9f\xd0\xbe\xd0\xbb\xd1\x83\xd1\x87\xd0\xb8\xd1\x82\xd1\x8c \xd0\xbc\xd0\xb5\xd1\x82\xd0\xb0\xd0\xb4\xd0\xb0\xd0\xbd\xd0\xbd\xd1\x8b\xd0\xb5 \xd0\xb8\xd0\xb7 \xd0\xb7\xd0\xb0\xd0\xbf\xd1\x80\xd0\xbe\xd1\x81\xd0\xb0\x32k\xd0\x98\xd0\xb7\xd0\xb2\xd0\xbb\xd0\xb5\xd1\x87\xd0\xb5\xd0\xbd\xd0\xb8\xd0\xb5 \xd0\xbc\xd0\xb5\xd1\x82\xd0\xb0\xd0\xb4\xd0\xb0\xd0\xbd\xd0\xbd\xd1\x8b\xd1\x85 \xd0\xb8\xd0\xb7 \xd0\xb7\xd0\xb0\xd0\xbf\xd1\x80\xd0\xbe\xd1\x81\xd0\xb0 \xd0\xb4\xd0\xbb\xd1\x8f \xd0\xbf\xd0\xb5\xd1\x80\xd0\xb5\xd0\xb4\xd0\xb0\xd1\x87\xd0\xb8 \xd0\xb2 \xd0\xbf\xd1\x80\xd0\xb5\xd0\xb4\xd0\xb8\xd0\xba\xd1\x82\xd0\xbe\xd1\x80\xd2\x01\x06prompt2Z{\"prompt\": \"\xd1\x8f \xd1\x85\xd0\xbe\xd1\x87\xd1\x83 \xd0\xba\xd1\x83\xd0\xbf\xd0\xb8\xd1\x82\xd1\x8c \xd0\xb1\xd1\x83\xd0\xbc\xd0\xb0\xd0\xb3\xd1\x83 A4 \xd0\xbd\xd0\xb0 18 \xd0\xbc\xd0\xb5\xd1\x81\xd1\x8f\xd1\x86\xd0\xb5\xd0\xb2 \xd0\xb2\xd0\xbf\xd0\xb5\xd1\x80\xd0\xb5\xd0\xb4\"}\"\xc0\x02\n\x0f\x45xtractedPrompt\x12\x1c\n\x04type\x18\x01 \x01(\x0e\x32\x0e.api.QueryType\x12\x0f\n\x07product\x18\x02 \x01(\t\x12\x0e\n\x06period\x18\x03 \x01(\t:\xed\x01\x92\x41\xe9\x01\n\x9c\x01*(\xd0\x9c\xd0\xb5\xd1\x82\xd0\xb0\xd0\xb4\xd0\xb0\xd0\xbd\xd0\xbd\xd1\x8b\xd0\xb5 \xd0\xb8\xd0\xb7 \xd0\xb7\xd0\xb0\xd0\xbf\xd1\x80\xd0\xbe\xd1\x81\xd0\xb0\x32V\xd0\x9c\xd0\xb5\xd1\x82\xd0\xb0\xd0\xb4\xd0\xb0\xd0\xbd\xd0\xbd\xd1\x8b\xd0\xb5 \xd0\xb8\xd0\xb7 \xd0\xb7\xd0\xb0\xd0\xbf\xd1\x80\xd0\xbe\xd1\x81\xd0\xb0 \xd0\xb4\xd0\xbb\xd1\x8f \xd0\xbf\xd0\xb5\xd1\x80\xd0\xb5\xd0\xb4\xd0\xb0\xd1\x87\xd0\xb8 \xd0\xb2 \xd0\xbf\xd1\x80\xd0\xb5\xd0\xb4\xd0\xb8\xd0\xba\xd1\x82\xd0\xbe\xd1\x80\xd2\x01\x04type\xd2\x01\x07product\xd2\x01\x06period2H{\"type\": 1, \"product\": \"\xd0\xb1\xd1\x83\xd0\xbc\xd0\xb0\xd0\xb3\xd0\xb0 A4\", \"period\": \"18 \xd0\xbc\xd0\xb5\xd1\x81\xd1\x8f\xd1\x86\xd0\xb5\xd0\xb2\"}*5\n\tQueryType\x12\r\n\tUNDEFINED\x10\x00\x12\x0e\n\nPREDICTION\x10\x01\x12\t\n\x05STOCK\x10\x02\x32\xf5\x03\n\x08Prompter\x12\xb2\x03\n\x07\x45xtract\x12\x0f.api.ExtractReq\x1a\x14.api.ExtractedPrompt\"\xff\x02\x92\x41\xd8\x02\n\x08prompter\x12\x37\xd0\x98\xd0\xb7\xd0\xb2\xd0\xbb\xd0\xb5\xd1\x87\xd1\x8c \xd0\xbc\xd0\xb5\xd1\x82\xd0\xb0\xd0\xb4\xd0\xb0\xd0\xbd\xd0\xbd\xd1\x8b\xd0\xb5 \xd0\xb8\xd0\xb7 \xd0\xb7\xd0\xb0\xd0\xbf\xd1\x80\xd0\xbe\xd1\x81\xd0\xb0\x1ak\xd0\x98\xd0\xb7\xd0\xb2\xd0\xbb\xd0\xb5\xd1\x87\xd0\xb5\xd0\xbd\xd0\xb8\xd0\xb5 \xd0\xbc\xd0\xb5\xd1\x82\xd0\xb0\xd0\xb4\xd0\xb0\xd0\xbd\xd0\xbd\xd1\x8b\xd1\x85 \xd0\xb8\xd0\xb7 \xd0\xb7\xd0\xb0\xd0\xbf\xd1\x80\xd0\xbe\xd1\x81\xd0\xb0 \xd0\xb4\xd0\xbb\xd1\x8f \xd0\xbf\xd0\xb5\xd1\x80\xd0\xb5\xd0\xb4\xd0\xb0\xd1\x87\xd0\xb8 \xd0\xb2 \xd0\xbf\xd1\x80\xd0\xb5\xd0\xb4\xd0\xb8\xd0\xba\xd1\x82\xd0\xbe\xd1\x80J\x90\x01\n\x03\x32\x30\x30\x12\x88\x01\n(\xd0\x9c\xd0\xb5\xd1\x82\xd0\xb0\xd0\xb4\xd0\xb0\xd0\xbd\xd0\xbd\xd1\x8b\xd0\xb5 \xd0\xb8\xd0\xb7 \xd0\xb7\xd0\xb0\xd0\xbf\xd1\x80\xd0\xbe\xd1\x81\xd0\xb0\"\\\n\x10\x61pplication/json\x12H{\"type\": 1, \"product\": \"\xd0\xb1\xd1\x83\xd0\xbc\xd0\xb0\xd0\xb3\xd0\xb0 A4\", \"period\": \"18 \xd0\xbc\xd0\xb5\xd1\x81\xd1\x8f\xd1\x86\xd0\xb5\xd0\xb2\"}b\x13\n\x11\n\rAuthorization\x12\x00\x82\xd3\xe4\x93\x02\x1d\"\x18/api/v1/prompter/extract:\x01*\x12\x34\n\rRespondStream\x12\x0e.api.StreamReq\x1a\x0f.api.StreamResp\"\x00\x30\x01\x42\xe9\x01Z\x0finternal/api/pb\x92\x41\xd4\x01\x12\xa9\x01\n\x03\x41PI\x12H\xd0\x94\xd0\xbe\xd0\xba\xd1\x83\xd0\xbc\xd0\xb5\xd0\xbd\xd1\x82\xd0\xb0\xd1\x86\xd0\xb8\xd1\x8f \xd0\xba API-\xd1\x81\xd0\xb5\xd1\x80\xd0\xb2\xd0\xb8\xd1\x81\xd1\x83 \xd0\xba\xd0\xbe\xd0\xbc\xd0\xb0\xd0\xbd\xd0\xb4\xd1\x8b misis.tech*X\n\x14\x42SD 3-Clause License\x12@https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSEZ&\n$\n\rAuthorization\x12\x13\x08\x02\x1a\rAuthorization \x02\x62\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, "api.prompter_pb2", _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\017internal/api/pb\222A\324\001\022\251\001\n\003API\022H\320\224\320\276\320\272\321\203\320\274\320\265\320\275\321\202\320\260\321\206\320\270\321\217 \320\272 API-\321\201\320\265\321\200\320\262\320\270\321\201\321\203 \320\272\320\276\320\274\320\260\320\275\320\264\321\213 misis.tech*X\n\024BSD 3-Clause License\022@https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSEZ&\n$\n\rAuthorization\022\023\010\002\032\rAuthorization \002'
  _globals['_EXTRACTREQ']._loaded_options = None
  _globals['_EXTRACTREQ']._serialized_options = b'\222A\220\002\n\261\001*9\320\237\320\276\320\273\321\203\321\207\320\270\321\202\321\214 \320\274\320\265\321\202\320\260\320\264\320\260\320\275\320\275\321\213\320\265 \320\270\320\267 \320\267\320\260\320\277\321\200\320\276\321\201\320\2602k\320\230\320\267\320\262\320\273\320\265\321\207\320\265\320\275\320\270\320\265 \320\274\320\265\321\202\320\260\320\264\320\260\320\275\320\275\321\213\321\205 \320\270\320\267 \320\267\320\260\320\277\321\200\320\276\321\201\320\260 \320\264\320\273\321\217 \320\277\320\265\321\200\320\265\320\264\320\260\321\207\320\270 \320\262 \320\277\321\200\320\265\320\264\320\270\320\272\321\202\320\276\321\200\322\001\006prompt2Z{\"prompt\": \"\321\217 \321\205\320\276\321\207\321\203 \320\272\321\203\320\277\320\270\321\202\321\214 \320\261\321\203\320\274\320\260\320\263\321\203 A4 \320\275\320\260 18 \320\274\320\265\321\201\321\217\321\206\320\265\320\262 \320\262\320\277\320\265\321\200\320\265\320\264\"}'
  _globals['_EXTRACTEDPROMPT']._loaded_options = None
  _globals['_EXTRACTEDPROMPT']._serialized_options = b'\222A\351\001\n\234\001*(\320\234\320\265\321\202\320\260\320\264\320\260\320\275\320\275\321\213\320\265 \320\270\320\267 \320\267\320\260\320\277\321\200\320\276\321\201\320\2602V\320\234\320\265\321\202\320\260\320\264\320\260\320\275\320\275\321\213\320\265 \320\270\320\267 \320\267\320\260\320\277\321\200\320\276\321\201\320\260 \320\264\320\273\321\217 \320\277\320\265\321\200\320\265\320\264\320\260\321\207\320\270 \320\262 \320\277\321\200\320\265\320\264\320\270\320\272\321\202\320\276\321\200\322\001\004type\322\001\007product\322\001\006period2H{\"type\": 1, \"product\": \"\320\261\321\203\320\274\320\260\320\263\320\260 A4\", \"period\": \"18 \320\274\320\265\321\201\321\217\321\206\320\265\320\262\"}'
  _globals['_PROMPTER'].methods_by_name['Extract']._loaded_options = None
  _globals['_PROMPTER'].methods_by_name['Extract']._serialized_options = b'\222A\330\002\n\010prompter\0227\320\230\320\267\320\262\320\273\320\265\321\207\321\214 \320\274\320\265\321\202\320\260\320\264\320\260\320\275\320\275\321\213\320\265 \320\270\320\267 \320\267\320\260\320\277\321\200\320\276\321\201\320\260\032k\320\230\320\267\320\262\320\273\320\265\321\207\320\265\320\275\320\270\320\265 \320\274\320\265\321\202\320\260\320\264\320\260\320\275\320\275\321\213\321\205 \320\270\320\267 \320\267\320\260\320\277\321\200\320\276\321\201\320\260 \320\264\320\273\321\217 \320\277\320\265\321\200\320\265\320\264\320\260\321\207\320\270 \320\262 \320\277\321\200\320\265\320\264\320\270\320\272\321\202\320\276\321\200J\220\001\n\003200\022\210\001\n(\320\234\320\265\321\202\320\260\320\264\320\260\320\275\320\275\321\213\320\265 \320\270\320\267 \320\267\320\260\320\277\321\200\320\276\321\201\320\260\"\\\n\020application/json\022H{\"type\": 1, \"product\": \"\320\261\321\203\320\274\320\260\320\263\320\260 A4\", \"period\": \"18 \320\274\320\265\321\201\321\217\321\206\320\265\320\262\"}b\023\n\021\n\rAuthorization\022\000\202\323\344\223\002\035\"\030/api/v1/prompter/extract:\001*'
  _globals['_QUERYTYPE']._serialized_start=796
  _globals['_QUERYTYPE']._serialized_end=849
  _globals['_STREAMREQ']._serialized_start=105
  _globals['_STREAMREQ']._serialized_end=132
  _globals['_STREAMRESP']._serialized_start=134
  _globals['_STREAMRESP']._serialized_end=161
  _globals['_EXTRACTREQ']._serialized_start=164
  _globals['_EXTRACTREQ']._serialized_end=471
  _globals['_EXTRACTEDPROMPT']._serialized_start=474
  _globals['_EXTRACTEDPROMPT']._serialized_end=794
  _globals['_PROMPTER']._serialized_start=852
  _globals['_PROMPTER']._serialized_end=1353
# @@protoc_insertion_point(module_scope)
