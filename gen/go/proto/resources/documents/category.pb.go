// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.20.3
// source: resources/documents/category.proto

package documents

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type Category struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// @sanitize
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// @sanitize
	Description *string `protobuf:"bytes,3,opt,name=description,proto3,oneof" json:"description,omitempty"`
	Job         *string `protobuf:"bytes,4,opt,name=job,proto3,oneof" json:"job,omitempty"`
	// @sanitize: method=StripTags
	Color string `protobuf:"bytes,5,opt,name=color,proto3" json:"color,omitempty"`
	// @sanitize: method=StripTags
	Icon *string `protobuf:"bytes,6,opt,name=icon,proto3,oneof" json:"icon,omitempty"`
}

func (x *Category) Reset() {
	*x = Category{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_documents_category_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Category) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Category) ProtoMessage() {}

func (x *Category) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_category_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Category.ProtoReflect.Descriptor instead.
func (*Category) Descriptor() ([]byte, []int) {
	return file_resources_documents_category_proto_rawDescGZIP(), []int{0}
}

func (x *Category) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Category) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Category) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Category) GetJob() string {
	if x != nil && x.Job != nil {
		return *x.Job
	}
	return ""
}

func (x *Category) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *Category) GetIcon() string {
	if x != nil && x.Icon != nil {
		return *x.Icon
	}
	return ""
}

var File_resources_documents_category_proto protoreflect.FileDescriptor

var file_resources_documents_category_proto_rawDesc = []byte{
	0x0a, 0x22, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xf4, 0x01, 0x0a, 0x08, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12,
	0x12, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x30, 0x01, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x03, 0x18, 0x80, 0x01, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x2f, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18,
	0xff, 0x01, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x48, 0x01, 0x52, 0x03, 0x6a, 0x6f,
	0x62, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x10, 0x03, 0x18, 0x07, 0x52, 0x05,
	0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x21, 0x0a, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x01, 0x48, 0x02, 0x52,
	0x04, 0x69, 0x63, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6a, 0x6f, 0x62,
	0x42, 0x07, 0x0a, 0x05, 0x5f, 0x69, 0x63, 0x6f, 0x6e, 0x42, 0x4b, 0x5a, 0x49, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d,
	0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x3b, 0x64, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_documents_category_proto_rawDescOnce sync.Once
	file_resources_documents_category_proto_rawDescData = file_resources_documents_category_proto_rawDesc
)

func file_resources_documents_category_proto_rawDescGZIP() []byte {
	file_resources_documents_category_proto_rawDescOnce.Do(func() {
		file_resources_documents_category_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_documents_category_proto_rawDescData)
	})
	return file_resources_documents_category_proto_rawDescData
}

var file_resources_documents_category_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_documents_category_proto_goTypes = []any{
	(*Category)(nil), // 0: resources.documents.Category
}
var file_resources_documents_category_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_resources_documents_category_proto_init() }
func file_resources_documents_category_proto_init() {
	if File_resources_documents_category_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_documents_category_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Category); i {
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
	file_resources_documents_category_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_documents_category_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_documents_category_proto_goTypes,
		DependencyIndexes: file_resources_documents_category_proto_depIdxs,
		MessageInfos:      file_resources_documents_category_proto_msgTypes,
	}.Build()
	File_resources_documents_category_proto = out.File
	file_resources_documents_category_proto_rawDesc = nil
	file_resources_documents_category_proto_goTypes = nil
	file_resources_documents_category_proto_depIdxs = nil
}
