// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: resources/common/i18n.proto

package common

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Wrapped translated message for the client
// @dbscanner: json
type TranslateItem struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// @sanitize: method=StripTags
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// @sanitize: method=StripTags
	Parameters    map[string]string `protobuf:"bytes,2,rep,name=parameters,proto3" json:"parameters,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TranslateItem) Reset() {
	*x = TranslateItem{}
	mi := &file_resources_common_i18n_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TranslateItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TranslateItem) ProtoMessage() {}

func (x *TranslateItem) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_i18n_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TranslateItem.ProtoReflect.Descriptor instead.
func (*TranslateItem) Descriptor() ([]byte, []int) {
	return file_resources_common_i18n_proto_rawDescGZIP(), []int{0}
}

func (x *TranslateItem) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *TranslateItem) GetParameters() map[string]string {
	if x != nil {
		return x.Parameters
	}
	return nil
}

var File_resources_common_i18n_proto protoreflect.FileDescriptor

const file_resources_common_i18n_proto_rawDesc = "" +
	"\n" +
	"\x1bresources/common/i18n.proto\x12\x10resources.common\"\xb1\x01\n" +
	"\rTranslateItem\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12O\n" +
	"\n" +
	"parameters\x18\x02 \x03(\v2/.resources.common.TranslateItem.ParametersEntryR\n" +
	"parameters\x1a=\n" +
	"\x0fParametersEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01BEZCgithub.com/fivenet-app/fivenet/gen/go/proto/resources/common;commonb\x06proto3"

var (
	file_resources_common_i18n_proto_rawDescOnce sync.Once
	file_resources_common_i18n_proto_rawDescData []byte
)

func file_resources_common_i18n_proto_rawDescGZIP() []byte {
	file_resources_common_i18n_proto_rawDescOnce.Do(func() {
		file_resources_common_i18n_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_common_i18n_proto_rawDesc), len(file_resources_common_i18n_proto_rawDesc)))
	})
	return file_resources_common_i18n_proto_rawDescData
}

var file_resources_common_i18n_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_resources_common_i18n_proto_goTypes = []any{
	(*TranslateItem)(nil), // 0: resources.common.TranslateItem
	nil,                   // 1: resources.common.TranslateItem.ParametersEntry
}
var file_resources_common_i18n_proto_depIdxs = []int32{
	1, // 0: resources.common.TranslateItem.parameters:type_name -> resources.common.TranslateItem.ParametersEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_resources_common_i18n_proto_init() }
func file_resources_common_i18n_proto_init() {
	if File_resources_common_i18n_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_common_i18n_proto_rawDesc), len(file_resources_common_i18n_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_common_i18n_proto_goTypes,
		DependencyIndexes: file_resources_common_i18n_proto_depIdxs,
		MessageInfos:      file_resources_common_i18n_proto_msgTypes,
	}.Build()
	File_resources_common_i18n_proto = out.File
	file_resources_common_i18n_proto_goTypes = nil
	file_resources_common_i18n_proto_depIdxs = nil
}
