// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: resources/dispatch/units.proto

package dispatch

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

type UNIT_STATUS int32

const (
	UNIT_STATUS_UNAVAILABLE UNIT_STATUS = 0
	UNIT_STATUS_AVAILABLE   UNIT_STATUS = 1
	UNIT_STATUS_ON_BREAK    UNIT_STATUS = 2
	UNIT_STATUS_BUSY        UNIT_STATUS = 3
)

// Enum value maps for UNIT_STATUS.
var (
	UNIT_STATUS_name = map[int32]string{
		0: "UNAVAILABLE",
		1: "AVAILABLE",
		2: "ON_BREAK",
		3: "BUSY",
	}
	UNIT_STATUS_value = map[string]int32{
		"UNAVAILABLE": 0,
		"AVAILABLE":   1,
		"ON_BREAK":    2,
		"BUSY":        3,
	}
)

func (x UNIT_STATUS) Enum() *UNIT_STATUS {
	p := new(UNIT_STATUS)
	*p = x
	return p
}

func (x UNIT_STATUS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UNIT_STATUS) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_dispatch_units_proto_enumTypes[0].Descriptor()
}

func (UNIT_STATUS) Type() protoreflect.EnumType {
	return &file_resources_dispatch_units_proto_enumTypes[0]
}

func (x UNIT_STATUS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UNIT_STATUS.Descriptor instead.
func (UNIT_STATUS) EnumDescriptor() ([]byte, []int) {
	return file_resources_dispatch_units_proto_rawDescGZIP(), []int{0}
}

type Unit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"`                                     // @gotags: sql:"primary_key" alias:"id"
	CreatedAt   *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty" alias:"created_at"` // @gotags: alias:"created_at"
	UpdatedAt   *timestamp.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty" alias:"updated_at"` // @gotags: alias:"updated_at"
	Name        string               `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Initials    string               `protobuf:"bytes,5,opt,name=initials,proto3" json:"initials,omitempty"`
	Color       *string              `protobuf:"bytes,6,opt,name=color,proto3,oneof" json:"color,omitempty"`
	Description *string              `protobuf:"bytes,7,opt,name=description,proto3,oneof" json:"description,omitempty"`
	Status      *UNIT_STATUS         `protobuf:"varint,8,opt,name=status,proto3,enum=resources.dispatch.UNIT_STATUS,oneof" json:"status,omitempty"`
	Users       []*users.UserShort   `protobuf:"bytes,9,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *Unit) Reset() {
	*x = Unit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_dispatch_units_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Unit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Unit) ProtoMessage() {}

func (x *Unit) ProtoReflect() protoreflect.Message {
	mi := &file_resources_dispatch_units_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Unit.ProtoReflect.Descriptor instead.
func (*Unit) Descriptor() ([]byte, []int) {
	return file_resources_dispatch_units_proto_rawDescGZIP(), []int{0}
}

func (x *Unit) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Unit) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Unit) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Unit) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Unit) GetInitials() string {
	if x != nil {
		return x.Initials
	}
	return ""
}

func (x *Unit) GetColor() string {
	if x != nil && x.Color != nil {
		return *x.Color
	}
	return ""
}

func (x *Unit) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Unit) GetStatus() UNIT_STATUS {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return UNIT_STATUS_UNAVAILABLE
}

func (x *Unit) GetUsers() []*users.UserShort {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_resources_dispatch_units_proto protoreflect.FileDescriptor

var file_resources_dispatch_units_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x69, 0x73, 0x70,
	0x61, 0x74, 0x63, 0x68, 0x2f, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x12, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x70,
	0x61, 0x74, 0x63, 0x68, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xec, 0x03, 0x0a, 0x04, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x48, 0x01, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01,
	0x12, 0x1d, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09,
	0xfa, 0x42, 0x06, 0x72, 0x04, 0x10, 0x03, 0x18, 0x18, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x25, 0x0a, 0x08, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x10, 0x02, 0x18, 0x04, 0x52, 0x08, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x12, 0x22, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x06, 0x48, 0x02,
	0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x12, 0x2f, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x48, 0x03, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x3c, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x2e, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x48, 0x04, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x88, 0x01, 0x01, 0x12, 0x30, 0x0a, 0x05, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x63, 0x6f,
	0x6c, 0x6f, 0x72, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2a, 0x45,
	0x0a, 0x0b, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x12, 0x0f, 0x0a,
	0x0b, 0x55, 0x4e, 0x41, 0x56, 0x41, 0x49, 0x4c, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x0d,
	0x0a, 0x09, 0x41, 0x56, 0x41, 0x49, 0x4c, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x01, 0x12, 0x0c, 0x0a,
	0x08, 0x4f, 0x4e, 0x5f, 0x42, 0x52, 0x45, 0x41, 0x4b, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x42,
	0x55, 0x53, 0x59, 0x10, 0x03, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x66, 0x69, 0x76, 0x65,
	0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x69, 0x73, 0x70, 0x61,
	0x74, 0x63, 0x68, 0x3b, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_dispatch_units_proto_rawDescOnce sync.Once
	file_resources_dispatch_units_proto_rawDescData = file_resources_dispatch_units_proto_rawDesc
)

func file_resources_dispatch_units_proto_rawDescGZIP() []byte {
	file_resources_dispatch_units_proto_rawDescOnce.Do(func() {
		file_resources_dispatch_units_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_dispatch_units_proto_rawDescData)
	})
	return file_resources_dispatch_units_proto_rawDescData
}

var file_resources_dispatch_units_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_dispatch_units_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_dispatch_units_proto_goTypes = []interface{}{
	(UNIT_STATUS)(0),            // 0: resources.dispatch.UNIT_STATUS
	(*Unit)(nil),                // 1: resources.dispatch.Unit
	(*timestamp.Timestamp)(nil), // 2: resources.timestamp.Timestamp
	(*users.UserShort)(nil),     // 3: resources.users.UserShort
}
var file_resources_dispatch_units_proto_depIdxs = []int32{
	2, // 0: resources.dispatch.Unit.created_at:type_name -> resources.timestamp.Timestamp
	2, // 1: resources.dispatch.Unit.updated_at:type_name -> resources.timestamp.Timestamp
	0, // 2: resources.dispatch.Unit.status:type_name -> resources.dispatch.UNIT_STATUS
	3, // 3: resources.dispatch.Unit.users:type_name -> resources.users.UserShort
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_resources_dispatch_units_proto_init() }
func file_resources_dispatch_units_proto_init() {
	if File_resources_dispatch_units_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_dispatch_units_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Unit); i {
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
	file_resources_dispatch_units_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_dispatch_units_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_dispatch_units_proto_goTypes,
		DependencyIndexes: file_resources_dispatch_units_proto_depIdxs,
		EnumInfos:         file_resources_dispatch_units_proto_enumTypes,
		MessageInfos:      file_resources_dispatch_units_proto_msgTypes,
	}.Build()
	File_resources_dispatch_units_proto = out.File
	file_resources_dispatch_units_proto_rawDesc = nil
	file_resources_dispatch_units_proto_goTypes = nil
	file_resources_dispatch_units_proto_depIdxs = nil
}
