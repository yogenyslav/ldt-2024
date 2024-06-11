// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: api/prompter.proto

package pb

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type QueryType int32

const (
	QueryType_UNDEFINED  QueryType = 0
	QueryType_PREDICTION QueryType = 1
	QueryType_STOCK      QueryType = 2
)

// Enum value maps for QueryType.
var (
	QueryType_name = map[int32]string{
		0: "UNDEFINED",
		1: "PREDICTION",
		2: "STOCK",
	}
	QueryType_value = map[string]int32{
		"UNDEFINED":  0,
		"PREDICTION": 1,
		"STOCK":      2,
	}
)

func (x QueryType) Enum() *QueryType {
	p := new(QueryType)
	*p = x
	return p
}

func (x QueryType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QueryType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_prompter_proto_enumTypes[0].Descriptor()
}

func (QueryType) Type() protoreflect.EnumType {
	return &file_api_prompter_proto_enumTypes[0]
}

func (x QueryType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QueryType.Descriptor instead.
func (QueryType) EnumDescriptor() ([]byte, []int) {
	return file_api_prompter_proto_rawDescGZIP(), []int{0}
}

type ExtractReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prompt string `protobuf:"bytes,1,opt,name=prompt,proto3" json:"prompt,omitempty"`
}

func (x *ExtractReq) Reset() {
	*x = ExtractReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prompter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExtractReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractReq) ProtoMessage() {}

func (x *ExtractReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_prompter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractReq.ProtoReflect.Descriptor instead.
func (*ExtractReq) Descriptor() ([]byte, []int) {
	return file_api_prompter_proto_rawDescGZIP(), []int{0}
}

func (x *ExtractReq) GetPrompt() string {
	if x != nil {
		return x.Prompt
	}
	return ""
}

type ExtractedPrompt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    QueryType `protobuf:"varint,1,opt,name=type,proto3,enum=api.QueryType" json:"type,omitempty"`
	Product string    `protobuf:"bytes,2,opt,name=product,proto3" json:"product,omitempty"`
}

func (x *ExtractedPrompt) Reset() {
	*x = ExtractedPrompt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prompter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExtractedPrompt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractedPrompt) ProtoMessage() {}

func (x *ExtractedPrompt) ProtoReflect() protoreflect.Message {
	mi := &file_api_prompter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractedPrompt.ProtoReflect.Descriptor instead.
func (*ExtractedPrompt) Descriptor() ([]byte, []int) {
	return file_api_prompter_proto_rawDescGZIP(), []int{1}
}

func (x *ExtractedPrompt) GetType() QueryType {
	if x != nil {
		return x.Type
	}
	return QueryType_UNDEFINED
}

func (x *ExtractedPrompt) GetProduct() string {
	if x != nil {
		return x.Product
	}
	return ""
}

var File_api_prompter_proto protoreflect.FileDescriptor

var file_api_prompter_proto_rawDesc = []byte{
	0x0a, 0x12, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d,
	0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd2, 0x01, 0x0a, 0x0a, 0x45, 0x78, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x3a, 0xab,
	0x01, 0x92, 0x41, 0xa7, 0x01, 0x0a, 0x49, 0x2a, 0x16, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x20, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x20, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x32,
	0x26, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x20, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72,
	0x65, 0x64, 0x20, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x20, 0x66, 0x72, 0x6f, 0x6d,
	0x20, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0xd2, 0x01, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x32, 0x5a, 0x7b, 0x22, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x22, 0x3a, 0x20, 0x22, 0xd1, 0x8f,
	0x20, 0xd1, 0x85, 0xd0, 0xbe, 0xd1, 0x87, 0xd1, 0x83, 0x20, 0xd0, 0xba, 0xd1, 0x83, 0xd0, 0xbf,
	0xd0, 0xb8, 0xd1, 0x82, 0xd1, 0x8c, 0x20, 0xd0, 0xb1, 0xd1, 0x83, 0xd0, 0xbc, 0xd0, 0xb0, 0xd0,
	0xb3, 0xd1, 0x83, 0x20, 0x41, 0x34, 0x20, 0xd0, 0xbd, 0xd0, 0xb0, 0x20, 0x31, 0x38, 0x20, 0xd0,
	0xbc, 0xd0, 0xb5, 0xd1, 0x81, 0xd1, 0x8f, 0xd1, 0x86, 0xd0, 0xb5, 0xd0, 0xb2, 0x20, 0xd0, 0xb2,
	0xd0, 0xbf, 0xd0, 0xb5, 0xd1, 0x80, 0xd0, 0xb5, 0xd0, 0xb4, 0x22, 0x7d, 0x22, 0xe3, 0x01, 0x0a,
	0x0f, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x12, 0x22, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x3a, 0x91,
	0x01, 0x92, 0x41, 0x8d, 0x01, 0x0a, 0x60, 0x2a, 0x17, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x20, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x20, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x32, 0x34, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20,
	0x72, 0x61, 0x77, 0x20, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x20, 0x74, 0x68, 0x61, 0x74, 0x20,
	0x69, 0x73, 0x20, 0x70, 0x61, 0x73, 0x73, 0x65, 0x64, 0x20, 0x74, 0x6f, 0x20, 0x70, 0x72, 0x65,
	0x64, 0x69, 0x63, 0x74, 0x6f, 0x72, 0xd2, 0x01, 0x04, 0x74, 0x79, 0x70, 0x65, 0xd2, 0x01, 0x07,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x32, 0x29, 0x7b, 0x22, 0x74, 0x79, 0x70, 0x65, 0x22,
	0x3a, 0x20, 0x31, 0x2c, 0x20, 0x22, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x22, 0x3a, 0x20,
	0x22, 0xd0, 0xb1, 0xd1, 0x83, 0xd0, 0xbc, 0xd0, 0xb0, 0xd0, 0xb3, 0xd0, 0xb0, 0x20, 0x41, 0x34,
	0x22, 0x7d, 0x2a, 0x35, 0x0a, 0x09, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0d, 0x0a, 0x09, 0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e,
	0x0a, 0x0a, 0x50, 0x52, 0x45, 0x44, 0x49, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x09,
	0x0a, 0x05, 0x53, 0x54, 0x4f, 0x43, 0x4b, 0x10, 0x02, 0x32, 0x79, 0x0a, 0x08, 0x50, 0x72, 0x6f,
	0x6d, 0x70, 0x74, 0x65, 0x72, 0x12, 0x6d, 0x0a, 0x07, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x12, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x65,
	0x71, 0x1a, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x65,
	0x64, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x22, 0x3b, 0x92, 0x41, 0x15, 0x62, 0x13, 0x0a, 0x11,
	0x0a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x3a, 0x01, 0x2a, 0x22, 0x18, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x65, 0x72, 0x2f, 0x65, 0x78, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x42, 0xe9, 0x01, 0x92, 0x41, 0xd4, 0x01, 0x12, 0xa9, 0x01, 0x0a, 0x03,
	0x41, 0x50, 0x49, 0x12, 0x48, 0xd0, 0x94, 0xd0, 0xbe, 0xd0, 0xba, 0xd1, 0x83, 0xd0, 0xbc, 0xd0,
	0xb5, 0xd0, 0xbd, 0xd1, 0x82, 0xd0, 0xb0, 0xd1, 0x86, 0xd0, 0xb8, 0xd1, 0x8f, 0x20, 0xd0, 0xba,
	0x20, 0x41, 0x50, 0x49, 0x2d, 0xd1, 0x81, 0xd0, 0xb5, 0xd1, 0x80, 0xd0, 0xb2, 0xd0, 0xb8, 0xd1,
	0x81, 0xd1, 0x83, 0x20, 0xd0, 0xba, 0xd0, 0xbe, 0xd0, 0xbc, 0xd0, 0xb0, 0xd0, 0xbd, 0xd0, 0xb4,
	0xd1, 0x8b, 0x20, 0x6d, 0x69, 0x73, 0x69, 0x73, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2a, 0x58, 0x0a,
	0x14, 0x42, 0x53, 0x44, 0x20, 0x33, 0x2d, 0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x20, 0x4c, 0x69,
	0x63, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x65,
	0x63, 0x6f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2f,
	0x4c, 0x49, 0x43, 0x45, 0x4e, 0x53, 0x45, 0x5a, 0x26, 0x0a, 0x24, 0x0a, 0x0d, 0x41, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x13, 0x08, 0x02, 0x1a, 0x0d,
	0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x02, 0x5a,
	0x0f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_prompter_proto_rawDescOnce sync.Once
	file_api_prompter_proto_rawDescData = file_api_prompter_proto_rawDesc
)

func file_api_prompter_proto_rawDescGZIP() []byte {
	file_api_prompter_proto_rawDescOnce.Do(func() {
		file_api_prompter_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_prompter_proto_rawDescData)
	})
	return file_api_prompter_proto_rawDescData
}

var file_api_prompter_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_prompter_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_prompter_proto_goTypes = []interface{}{
	(QueryType)(0),          // 0: api.QueryType
	(*ExtractReq)(nil),      // 1: api.ExtractReq
	(*ExtractedPrompt)(nil), // 2: api.ExtractedPrompt
}
var file_api_prompter_proto_depIdxs = []int32{
	0, // 0: api.ExtractedPrompt.type:type_name -> api.QueryType
	1, // 1: api.Prompter.Extract:input_type -> api.ExtractReq
	2, // 2: api.Prompter.Extract:output_type -> api.ExtractedPrompt
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_prompter_proto_init() }
func file_api_prompter_proto_init() {
	if File_api_prompter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_prompter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExtractReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_prompter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExtractedPrompt); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_prompter_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_prompter_proto_goTypes,
		DependencyIndexes: file_api_prompter_proto_depIdxs,
		EnumInfos:         file_api_prompter_proto_enumTypes,
		MessageInfos:      file_api_prompter_proto_msgTypes,
	}.Build()
	File_api_prompter_proto = out.File
	file_api_prompter_proto_rawDesc = nil
	file_api_prompter_proto_goTypes = nil
	file_api_prompter_proto_depIdxs = nil
}
