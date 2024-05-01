// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: resources/common/i18n.proto

package common

import (
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

type TranslateItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key        string            `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Parameters map[string]string `protobuf:"bytes,2,rep,name=parameters,proto3" json:"parameters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *TranslateItem) Reset() {
	*x = TranslateItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_common_i18n_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TranslateItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TranslateItem) ProtoMessage() {}

func (x *TranslateItem) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_i18n_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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

var file_resources_common_i18n_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x69, 0x31, 0x38, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x22,
	0xb1, 0x01, 0x0a, 0x0d, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x4f, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x73, 0x1a, 0x3d, 0x0a, 0x0f, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69,
	0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_resources_common_i18n_proto_rawDescOnce sync.Once
	file_resources_common_i18n_proto_rawDescData = file_resources_common_i18n_proto_rawDesc
)

func file_resources_common_i18n_proto_rawDescGZIP() []byte {
	file_resources_common_i18n_proto_rawDescOnce.Do(func() {
		file_resources_common_i18n_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_common_i18n_proto_rawDescData)
	})
	return file_resources_common_i18n_proto_rawDescData
}

var file_resources_common_i18n_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_resources_common_i18n_proto_goTypes = []interface{}{
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
	if !protoimpl.UnsafeEnabled {
		file_resources_common_i18n_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TranslateItem); i {
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
			RawDescriptor: file_resources_common_i18n_proto_rawDesc,
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
	file_resources_common_i18n_proto_rawDesc = nil
	file_resources_common_i18n_proto_goTypes = nil
	file_resources_common_i18n_proto_depIdxs = nil
}
