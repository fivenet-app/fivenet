// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.20.3
// source: resources/internet/internet.proto

package internet

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

type Domain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Domain string `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
}

func (x *Domain) Reset() {
	*x = Domain{}
	mi := &file_resources_internet_internet_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Domain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Domain) ProtoMessage() {}

func (x *Domain) ProtoReflect() protoreflect.Message {
	mi := &file_resources_internet_internet_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Domain.ProtoReflect.Descriptor instead.
func (*Domain) Descriptor() ([]byte, []int) {
	return file_resources_internet_internet_proto_rawDescGZIP(), []int{0}
}

func (x *Domain) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Domain) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

type WebPage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DomainId uint64 `protobuf:"varint,2,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	Url      string `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *WebPage) Reset() {
	*x = WebPage{}
	mi := &file_resources_internet_internet_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WebPage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebPage) ProtoMessage() {}

func (x *WebPage) ProtoReflect() protoreflect.Message {
	mi := &file_resources_internet_internet_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebPage.ProtoReflect.Descriptor instead.
func (*WebPage) Descriptor() ([]byte, []int) {
	return file_resources_internet_internet_proto_rawDescGZIP(), []int{1}
}

func (x *WebPage) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *WebPage) GetDomainId() uint64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *WebPage) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_resources_internet_internet_proto protoreflect.FileDescriptor

var file_resources_internet_internet_proto_rawDesc = []byte{
	0x0a, 0x21, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x65, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x12, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x22, 0x34, 0x0a, 0x06, 0x44, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x12, 0x12, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x30,
	0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x50, 0x0a,
	0x07, 0x57, 0x65, 0x62, 0x50, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x30, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x09,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x42,
	0x02, 0x30, 0x01, 0x52, 0x08, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x42,
	0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69,
	0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65,
	0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65,
	0x74, 0x3b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_resources_internet_internet_proto_rawDescOnce sync.Once
	file_resources_internet_internet_proto_rawDescData = file_resources_internet_internet_proto_rawDesc
)

func file_resources_internet_internet_proto_rawDescGZIP() []byte {
	file_resources_internet_internet_proto_rawDescOnce.Do(func() {
		file_resources_internet_internet_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_internet_internet_proto_rawDescData)
	})
	return file_resources_internet_internet_proto_rawDescData
}

var file_resources_internet_internet_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_resources_internet_internet_proto_goTypes = []any{
	(*Domain)(nil),  // 0: resources.internet.Domain
	(*WebPage)(nil), // 1: resources.internet.WebPage
}
var file_resources_internet_internet_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_resources_internet_internet_proto_init() }
func file_resources_internet_internet_proto_init() {
	if File_resources_internet_internet_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_internet_internet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_internet_internet_proto_goTypes,
		DependencyIndexes: file_resources_internet_internet_proto_depIdxs,
		MessageInfos:      file_resources_internet_internet_proto_msgTypes,
	}.Build()
	File_resources_internet_internet_proto = out.File
	file_resources_internet_internet_proto_rawDesc = nil
	file_resources_internet_internet_proto_goTypes = nil
	file_resources_internet_internet_proto_depIdxs = nil
}