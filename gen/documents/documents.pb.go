// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: documents/documents.proto

package documents

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

type GetDocumentsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetDocumentsResponse) Reset() {
	*x = GetDocumentsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documents_documents_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDocumentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDocumentsResponse) ProtoMessage() {}

func (x *GetDocumentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_documents_documents_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDocumentsResponse.ProtoReflect.Descriptor instead.
func (*GetDocumentsResponse) Descriptor() ([]byte, []int) {
	return file_documents_documents_proto_rawDescGZIP(), []int{0}
}

type GetDocumentsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetDocumentsRequest) Reset() {
	*x = GetDocumentsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documents_documents_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDocumentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDocumentsRequest) ProtoMessage() {}

func (x *GetDocumentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_documents_documents_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDocumentsRequest.ProtoReflect.Descriptor instead.
func (*GetDocumentsRequest) Descriptor() ([]byte, []int) {
	return file_documents_documents_proto_rawDescGZIP(), []int{1}
}

type Document struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Document) Reset() {
	*x = Document{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documents_documents_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Document) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Document) ProtoMessage() {}

func (x *Document) ProtoReflect() protoreflect.Message {
	mi := &file_documents_documents_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Document.ProtoReflect.Descriptor instead.
func (*Document) Descriptor() ([]byte, []int) {
	return file_documents_documents_proto_rawDescGZIP(), []int{2}
}

var File_documents_documents_proto protoreflect.FileDescriptor

var file_documents_documents_proto_rawDesc = []byte{
	0x0a, 0x19, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x67, 0x65, 0x6e,
	0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x16, 0x0a, 0x14, 0x47, 0x65,
	0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x15, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0a, 0x0a, 0x08, 0x44, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x32, 0x6b, 0x0a, 0x10, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x0c, 0x47, 0x65, 0x74,
	0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x22, 0x2e, 0x67, 0x65, 0x6e, 0x2e,
	0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e,
	0x67, 0x65, 0x6e, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x47, 0x65,
	0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x61, 0x72, 0x70, 0x61, 0x6e, 0x65, 0x74,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_documents_documents_proto_rawDescOnce sync.Once
	file_documents_documents_proto_rawDescData = file_documents_documents_proto_rawDesc
)

func file_documents_documents_proto_rawDescGZIP() []byte {
	file_documents_documents_proto_rawDescOnce.Do(func() {
		file_documents_documents_proto_rawDescData = protoimpl.X.CompressGZIP(file_documents_documents_proto_rawDescData)
	})
	return file_documents_documents_proto_rawDescData
}

var file_documents_documents_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_documents_documents_proto_goTypes = []interface{}{
	(*GetDocumentsResponse)(nil), // 0: gen.documents.GetDocumentsResponse
	(*GetDocumentsRequest)(nil),  // 1: gen.documents.GetDocumentsRequest
	(*Document)(nil),             // 2: gen.documents.Document
}
var file_documents_documents_proto_depIdxs = []int32{
	1, // 0: gen.documents.DocumentsService.GetDocuments:input_type -> gen.documents.GetDocumentsRequest
	0, // 1: gen.documents.DocumentsService.GetDocuments:output_type -> gen.documents.GetDocumentsResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_documents_documents_proto_init() }
func file_documents_documents_proto_init() {
	if File_documents_documents_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_documents_documents_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDocumentsResponse); i {
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
		file_documents_documents_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDocumentsRequest); i {
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
		file_documents_documents_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Document); i {
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
			RawDescriptor: file_documents_documents_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_documents_documents_proto_goTypes,
		DependencyIndexes: file_documents_documents_proto_depIdxs,
		MessageInfos:      file_documents_documents_proto_msgTypes,
	}.Build()
	File_documents_documents_proto = out.File
	file_documents_documents_proto_rawDesc = nil
	file_documents_documents_proto_goTypes = nil
	file_documents_documents_proto_depIdxs = nil
}
