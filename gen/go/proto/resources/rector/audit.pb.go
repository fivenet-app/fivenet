// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.2
// source: resources/rector/audit.proto

package rector

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	timestamp "github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
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

type EVENT_TYPE int32

const (
	EVENT_TYPE_UNSPECIFIED EVENT_TYPE = 0
	EVENT_TYPE_ERRORED     EVENT_TYPE = 1
	EVENT_TYPE_VIEWED      EVENT_TYPE = 2
	EVENT_TYPE_CREATED     EVENT_TYPE = 3
	EVENT_TYPE_UPDATED     EVENT_TYPE = 4
	EVENT_TYPE_DELETED     EVENT_TYPE = 5
)

// Enum value maps for EVENT_TYPE.
var (
	EVENT_TYPE_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "ERRORED",
		2: "VIEWED",
		3: "CREATED",
		4: "UPDATED",
		5: "DELETED",
	}
	EVENT_TYPE_value = map[string]int32{
		"UNSPECIFIED": 0,
		"ERRORED":     1,
		"VIEWED":      2,
		"CREATED":     3,
		"UPDATED":     4,
		"DELETED":     5,
	}
)

func (x EVENT_TYPE) Enum() *EVENT_TYPE {
	p := new(EVENT_TYPE)
	*p = x
	return p
}

func (x EVENT_TYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EVENT_TYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_rector_audit_proto_enumTypes[0].Descriptor()
}

func (EVENT_TYPE) Type() protoreflect.EnumType {
	return &file_resources_rector_audit_proto_enumTypes[0]
}

func (x EVENT_TYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EVENT_TYPE.Descriptor instead.
func (EVENT_TYPE) EnumDescriptor() ([]byte, []int) {
	return file_resources_rector_audit_proto_rawDescGZIP(), []int{0}
}

type AuditEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" alias:"id"` // @gotags: alias:"id"
	CreatedAt    *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UserId       uint64               `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" alias:"user_id"` // @gotags: alias:"user_id"
	User         *users.UserShort     `protobuf:"bytes,4,opt,name=user,proto3,oneof" json:"user,omitempty"`
	UserJob      string               `protobuf:"bytes,5,opt,name=user_job,json=userJob,proto3" json:"user_job,omitempty" alias:"user_job"`                        // @gotags: alias:"user_job"
	TargetUserId *string              `protobuf:"bytes,6,opt,name=target_user_id,json=targetUserId,proto3,oneof" json:"target_user_id,omitempty" alias:"target_user_id"` // @gotags: alias:"target_user_id"
	TargetUser   *users.UserShort     `protobuf:"bytes,7,opt,name=target_user,json=targetUser,proto3,oneof" json:"target_user,omitempty"`
	Service      string               `protobuf:"bytes,8,opt,name=service,proto3" json:"service,omitempty" alias:"service"`                                // @gotags: alias:"service"
	Method       string               `protobuf:"bytes,9,opt,name=method,proto3" json:"method,omitempty" alias:"method"`                                  // @gotags: alias:"method"
	State        EVENT_TYPE           `protobuf:"varint,10,opt,name=state,proto3,enum=resources.rector.EVENT_TYPE" json:"state,omitempty" alias:"state"` // @gotags: alias:"state"
	Data         *string              `protobuf:"bytes,11,opt,name=data,proto3,oneof" json:"data,omitempty" alias:"data"`                               // @gotags: alias:"data"
}

func (x *AuditEntry) Reset() {
	*x = AuditEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_rector_audit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuditEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuditEntry) ProtoMessage() {}

func (x *AuditEntry) ProtoReflect() protoreflect.Message {
	mi := &file_resources_rector_audit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuditEntry.ProtoReflect.Descriptor instead.
func (*AuditEntry) Descriptor() ([]byte, []int) {
	return file_resources_rector_audit_proto_rawDescGZIP(), []int{0}
}

func (x *AuditEntry) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AuditEntry) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *AuditEntry) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *AuditEntry) GetUser() *users.UserShort {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *AuditEntry) GetUserJob() string {
	if x != nil {
		return x.UserJob
	}
	return ""
}

func (x *AuditEntry) GetTargetUserId() string {
	if x != nil && x.TargetUserId != nil {
		return *x.TargetUserId
	}
	return ""
}

func (x *AuditEntry) GetTargetUser() *users.UserShort {
	if x != nil {
		return x.TargetUser
	}
	return nil
}

func (x *AuditEntry) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *AuditEntry) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *AuditEntry) GetState() EVENT_TYPE {
	if x != nil {
		return x.State
	}
	return EVENT_TYPE_UNSPECIFIED
}

func (x *AuditEntry) GetData() string {
	if x != nil && x.Data != nil {
		return *x.Data
	}
	return ""
}

var File_resources_rector_audit_proto protoreflect.FileDescriptor

var file_resources_rector_audit_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x72, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xef, 0x03, 0x0a, 0x0a,
	0x41, 0x75, 0x64, 0x69, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3d, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x48, 0x00, 0x52, 0x04,
	0x75, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x6a, 0x6f, 0x62, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x4a,
	0x6f, 0x62, 0x12, 0x29, 0x0a, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0c, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x40, 0x0a,
	0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x48, 0x02,
	0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x12, 0x3c, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x2e, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x17, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x5d, 0x0a,
	0x0a, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x12, 0x0f, 0x0a, 0x0b, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x56, 0x49, 0x45,
	0x57, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x44,
	0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x44, 0x10, 0x04, 0x12,
	0x0b, 0x0a, 0x07, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x05, 0x42, 0x41, 0x5a, 0x3f,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78,
	0x72, 0x74, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x3b, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_rector_audit_proto_rawDescOnce sync.Once
	file_resources_rector_audit_proto_rawDescData = file_resources_rector_audit_proto_rawDesc
)

func file_resources_rector_audit_proto_rawDescGZIP() []byte {
	file_resources_rector_audit_proto_rawDescOnce.Do(func() {
		file_resources_rector_audit_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_rector_audit_proto_rawDescData)
	})
	return file_resources_rector_audit_proto_rawDescData
}

var file_resources_rector_audit_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_rector_audit_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_rector_audit_proto_goTypes = []interface{}{
	(EVENT_TYPE)(0),             // 0: resources.rector.EVENT_TYPE
	(*AuditEntry)(nil),          // 1: resources.rector.AuditEntry
	(*timestamp.Timestamp)(nil), // 2: resources.timestamp.Timestamp
	(*users.UserShort)(nil),     // 3: resources.users.UserShort
}
var file_resources_rector_audit_proto_depIdxs = []int32{
	2, // 0: resources.rector.AuditEntry.created_at:type_name -> resources.timestamp.Timestamp
	3, // 1: resources.rector.AuditEntry.user:type_name -> resources.users.UserShort
	3, // 2: resources.rector.AuditEntry.target_user:type_name -> resources.users.UserShort
	0, // 3: resources.rector.AuditEntry.state:type_name -> resources.rector.EVENT_TYPE
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_resources_rector_audit_proto_init() }
func file_resources_rector_audit_proto_init() {
	if File_resources_rector_audit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_rector_audit_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuditEntry); i {
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
	file_resources_rector_audit_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_rector_audit_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_rector_audit_proto_goTypes,
		DependencyIndexes: file_resources_rector_audit_proto_depIdxs,
		EnumInfos:         file_resources_rector_audit_proto_enumTypes,
		MessageInfos:      file_resources_rector_audit_proto_msgTypes,
	}.Build()
	File_resources_rector_audit_proto = out.File
	file_resources_rector_audit_proto_rawDesc = nil
	file_resources_rector_audit_proto_goTypes = nil
	file_resources_rector_audit_proto_depIdxs = nil
}
