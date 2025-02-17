// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v3.20.3
// source: resources/common/uuid.proto

package common

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type UUID struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Uuid          string                 `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UUID) Reset() {
	*x = UUID{}
	mi := &file_resources_common_uuid_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UUID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UUID) ProtoMessage() {}

func (x *UUID) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_uuid_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UUID.ProtoReflect.Descriptor instead.
func (*UUID) Descriptor() ([]byte, []int) {
	return file_resources_common_uuid_proto_rawDescGZIP(), []int{0}
}

func (x *UUID) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

var File_resources_common_uuid_proto protoreflect.FileDescriptor

var file_resources_common_uuid_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x75, 0x75, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x1a,
	0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x24, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44,
	0x12, 0x1c, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x42, 0x45,
	0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76,
	0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3b, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_common_uuid_proto_rawDescOnce sync.Once
	file_resources_common_uuid_proto_rawDescData []byte
)

func file_resources_common_uuid_proto_rawDescGZIP() []byte {
	file_resources_common_uuid_proto_rawDescOnce.Do(func() {
		file_resources_common_uuid_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_common_uuid_proto_rawDesc), len(file_resources_common_uuid_proto_rawDesc)))
	})
	return file_resources_common_uuid_proto_rawDescData
}

var file_resources_common_uuid_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_common_uuid_proto_goTypes = []any{
	(*UUID)(nil), // 0: resources.common.UUID
}
var file_resources_common_uuid_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_resources_common_uuid_proto_init() }
func file_resources_common_uuid_proto_init() {
	if File_resources_common_uuid_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_common_uuid_proto_rawDesc), len(file_resources_common_uuid_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_common_uuid_proto_goTypes,
		DependencyIndexes: file_resources_common_uuid_proto_depIdxs,
		MessageInfos:      file_resources_common_uuid_proto_msgTypes,
	}.Build()
	File_resources_common_uuid_proto = out.File
	file_resources_common_uuid_proto_goTypes = nil
	file_resources_common_uuid_proto_depIdxs = nil
}
