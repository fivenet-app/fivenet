// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: resources/centrum/user_unit.proto

package centrum

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
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

type UserUnitMapping struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UnitId        uint64                 `protobuf:"varint,1,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
	Job           string                 `protobuf:"bytes,2,opt,name=job,proto3" json:"job,omitempty"`
	UserId        int32                  `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserUnitMapping) Reset() {
	*x = UserUnitMapping{}
	mi := &file_resources_centrum_user_unit_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserUnitMapping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserUnitMapping) ProtoMessage() {}

func (x *UserUnitMapping) ProtoReflect() protoreflect.Message {
	mi := &file_resources_centrum_user_unit_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserUnitMapping.ProtoReflect.Descriptor instead.
func (*UserUnitMapping) Descriptor() ([]byte, []int) {
	return file_resources_centrum_user_unit_proto_rawDescGZIP(), []int{0}
}

func (x *UserUnitMapping) GetUnitId() uint64 {
	if x != nil {
		return x.UnitId
	}
	return 0
}

func (x *UserUnitMapping) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *UserUnitMapping) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserUnitMapping) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

var File_resources_centrum_user_unit_proto protoreflect.FileDescriptor

var file_resources_centrum_user_unit_proto_rawDesc = string([]byte{
	0x0a, 0x21, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x65, 0x6e, 0x74,
	0x72, 0x75, 0x6d, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x75, 0x6e, 0x69, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x11, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63,
	0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9d, 0x01, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x55, 0x6e, 0x69,
	0x74, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x74, 0x49,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6a, 0x6f, 0x62, 0x12, 0x20, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x28, 0x00, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x3d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x42, 0x47, 0x5a, 0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66,
	0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x65,
	0x6e, 0x74, 0x72, 0x75, 0x6d, 0x3b, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_centrum_user_unit_proto_rawDescOnce sync.Once
	file_resources_centrum_user_unit_proto_rawDescData []byte
)

func file_resources_centrum_user_unit_proto_rawDescGZIP() []byte {
	file_resources_centrum_user_unit_proto_rawDescOnce.Do(func() {
		file_resources_centrum_user_unit_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_centrum_user_unit_proto_rawDesc), len(file_resources_centrum_user_unit_proto_rawDesc)))
	})
	return file_resources_centrum_user_unit_proto_rawDescData
}

var file_resources_centrum_user_unit_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_centrum_user_unit_proto_goTypes = []any{
	(*UserUnitMapping)(nil),     // 0: resources.centrum.UserUnitMapping
	(*timestamp.Timestamp)(nil), // 1: resources.timestamp.Timestamp
}
var file_resources_centrum_user_unit_proto_depIdxs = []int32{
	1, // 0: resources.centrum.UserUnitMapping.created_at:type_name -> resources.timestamp.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_resources_centrum_user_unit_proto_init() }
func file_resources_centrum_user_unit_proto_init() {
	if File_resources_centrum_user_unit_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_centrum_user_unit_proto_rawDesc), len(file_resources_centrum_user_unit_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_centrum_user_unit_proto_goTypes,
		DependencyIndexes: file_resources_centrum_user_unit_proto_depIdxs,
		MessageInfos:      file_resources_centrum_user_unit_proto_msgTypes,
	}.Build()
	File_resources_centrum_user_unit_proto = out.File
	file_resources_centrum_user_unit_proto_goTypes = nil
	file_resources_centrum_user_unit_proto_depIdxs = nil
}
