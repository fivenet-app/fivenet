// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.20.3
// source: resources/common/database/database.proto

package database

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

type PaginationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset   int64  `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	PageSize *int64 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3,oneof" json:"page_size,omitempty"`
}

func (x *PaginationRequest) Reset() {
	*x = PaginationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_common_database_database_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaginationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaginationRequest) ProtoMessage() {}

func (x *PaginationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_database_database_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaginationRequest.ProtoReflect.Descriptor instead.
func (*PaginationRequest) Descriptor() ([]byte, []int) {
	return file_resources_common_database_database_proto_rawDescGZIP(), []int{0}
}

func (x *PaginationRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *PaginationRequest) GetPageSize() int64 {
	if x != nil && x.PageSize != nil {
		return *x.PageSize
	}
	return 0
}

type PaginationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalCount int64 `protobuf:"varint,1,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	Offset     int64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	End        int64 `protobuf:"varint,3,opt,name=end,proto3" json:"end,omitempty"`
	PageSize   int64 `protobuf:"varint,4,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *PaginationResponse) Reset() {
	*x = PaginationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_common_database_database_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaginationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaginationResponse) ProtoMessage() {}

func (x *PaginationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_database_database_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaginationResponse.ProtoReflect.Descriptor instead.
func (*PaginationResponse) Descriptor() ([]byte, []int) {
	return file_resources_common_database_database_proto_rawDescGZIP(), []int{1}
}

func (x *PaginationResponse) GetTotalCount() int64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *PaginationResponse) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *PaginationResponse) GetEnd() int64 {
	if x != nil {
		return x.End
	}
	return 0
}

func (x *PaginationResponse) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type Sort struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Column    string `protobuf:"bytes,1,opt,name=column,proto3" json:"column,omitempty"`
	Direction string `protobuf:"bytes,2,opt,name=direction,proto3" json:"direction,omitempty"`
}

func (x *Sort) Reset() {
	*x = Sort{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_common_database_database_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sort) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sort) ProtoMessage() {}

func (x *Sort) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_database_database_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sort.ProtoReflect.Descriptor instead.
func (*Sort) Descriptor() ([]byte, []int) {
	return file_resources_common_database_database_proto_rawDescGZIP(), []int{2}
}

func (x *Sort) GetColumn() string {
	if x != nil {
		return x.Column
	}
	return ""
}

func (x *Sort) GetDirection() string {
	if x != nil {
		return x.Direction
	}
	return ""
}

var File_resources_common_database_database_proto protoreflect.FileDescriptor

var file_resources_common_database_database_proto_rawDesc = []byte{
	0x0a, 0x28, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x62, 0x61, 0x73, 0x65, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x71,
	0x0a, 0x11, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x42, 0x09, 0xfa, 0x42, 0x04, 0x22, 0x02, 0x28, 0x00, 0x30, 0x02, 0x52, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x2b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x42, 0x09, 0xfa, 0x42, 0x04, 0x22, 0x02,
	0x28, 0x00, 0x30, 0x02, 0x48, 0x00, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x88, 0x01, 0x01, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a,
	0x65, 0x22, 0x8c, 0x01, 0x0a, 0x12, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x02, 0x30,
	0x02, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1a, 0x0a,
	0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x42, 0x02, 0x30,
	0x02, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x03, 0x65, 0x6e, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x42, 0x02, 0x30, 0x02, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12,
	0x1f, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x42, 0x02, 0x30, 0x02, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x22, 0x59, 0x0a, 0x04, 0x53, 0x6f, 0x72, 0x74, 0x12, 0x21, 0x0a, 0x06, 0x63, 0x6f, 0x6c, 0x75,
	0x6d, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x10,
	0x01, 0x18, 0x40, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x12, 0x2e, 0x0a, 0x09, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x10,
	0xfa, 0x42, 0x0d, 0x72, 0x0b, 0x52, 0x03, 0x61, 0x73, 0x63, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63,
	0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x50, 0x5a, 0x4e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65,
	0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x3b, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_common_database_database_proto_rawDescOnce sync.Once
	file_resources_common_database_database_proto_rawDescData = file_resources_common_database_database_proto_rawDesc
)

func file_resources_common_database_database_proto_rawDescGZIP() []byte {
	file_resources_common_database_database_proto_rawDescOnce.Do(func() {
		file_resources_common_database_database_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_common_database_database_proto_rawDescData)
	})
	return file_resources_common_database_database_proto_rawDescData
}

var file_resources_common_database_database_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_resources_common_database_database_proto_goTypes = []any{
	(*PaginationRequest)(nil),  // 0: resources.common.database.PaginationRequest
	(*PaginationResponse)(nil), // 1: resources.common.database.PaginationResponse
	(*Sort)(nil),               // 2: resources.common.database.Sort
}
var file_resources_common_database_database_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_resources_common_database_database_proto_init() }
func file_resources_common_database_database_proto_init() {
	if File_resources_common_database_database_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_common_database_database_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*PaginationRequest); i {
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
		file_resources_common_database_database_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*PaginationResponse); i {
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
		file_resources_common_database_database_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Sort); i {
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
	file_resources_common_database_database_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_common_database_database_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_common_database_database_proto_goTypes,
		DependencyIndexes: file_resources_common_database_database_proto_depIdxs,
		MessageInfos:      file_resources_common_database_database_proto_msgTypes,
	}.Build()
	File_resources_common_database_database_proto = out.File
	file_resources_common_database_database_proto_rawDesc = nil
	file_resources_common_database_database_proto_goTypes = nil
	file_resources_common_database_database_proto_depIdxs = nil
}
