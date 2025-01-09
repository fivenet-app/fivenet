// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v3.20.3
// source: resources/sync/data.proto

package sync

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	vehicles "github.com/fivenet-app/fivenet/gen/go/proto/resources/vehicles"
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

type DataStatus struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Count         int64                  `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DataStatus) Reset() {
	*x = DataStatus{}
	mi := &file_resources_sync_data_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataStatus) ProtoMessage() {}

func (x *DataStatus) ProtoReflect() protoreflect.Message {
	mi := &file_resources_sync_data_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataStatus.ProtoReflect.Descriptor instead.
func (*DataStatus) Descriptor() ([]byte, []int) {
	return file_resources_sync_data_proto_rawDescGZIP(), []int{0}
}

func (x *DataStatus) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type DataJobs struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Jobs          []*users.Job           `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DataJobs) Reset() {
	*x = DataJobs{}
	mi := &file_resources_sync_data_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataJobs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataJobs) ProtoMessage() {}

func (x *DataJobs) ProtoReflect() protoreflect.Message {
	mi := &file_resources_sync_data_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataJobs.ProtoReflect.Descriptor instead.
func (*DataJobs) Descriptor() ([]byte, []int) {
	return file_resources_sync_data_proto_rawDescGZIP(), []int{1}
}

func (x *DataJobs) GetJobs() []*users.Job {
	if x != nil {
		return x.Jobs
	}
	return nil
}

type DataUsers struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Users         []*users.User          `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DataUsers) Reset() {
	*x = DataUsers{}
	mi := &file_resources_sync_data_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataUsers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUsers) ProtoMessage() {}

func (x *DataUsers) ProtoReflect() protoreflect.Message {
	mi := &file_resources_sync_data_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUsers.ProtoReflect.Descriptor instead.
func (*DataUsers) Descriptor() ([]byte, []int) {
	return file_resources_sync_data_proto_rawDescGZIP(), []int{2}
}

func (x *DataUsers) GetUsers() []*users.User {
	if x != nil {
		return x.Users
	}
	return nil
}

type DataVehicles struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Vehicles      []*vehicles.Vehicle    `protobuf:"bytes,1,rep,name=vehicles,proto3" json:"vehicles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DataVehicles) Reset() {
	*x = DataVehicles{}
	mi := &file_resources_sync_data_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataVehicles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataVehicles) ProtoMessage() {}

func (x *DataVehicles) ProtoReflect() protoreflect.Message {
	mi := &file_resources_sync_data_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataVehicles.ProtoReflect.Descriptor instead.
func (*DataVehicles) Descriptor() ([]byte, []int) {
	return file_resources_sync_data_proto_rawDescGZIP(), []int{3}
}

func (x *DataVehicles) GetVehicles() []*vehicles.Vehicle {
	if x != nil {
		return x.Vehicles
	}
	return nil
}

type DataLicenses struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Licenses      []*users.License       `protobuf:"bytes,1,rep,name=licenses,proto3" json:"licenses,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DataLicenses) Reset() {
	*x = DataLicenses{}
	mi := &file_resources_sync_data_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataLicenses) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataLicenses) ProtoMessage() {}

func (x *DataLicenses) ProtoReflect() protoreflect.Message {
	mi := &file_resources_sync_data_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataLicenses.ProtoReflect.Descriptor instead.
func (*DataLicenses) Descriptor() ([]byte, []int) {
	return file_resources_sync_data_proto_rawDescGZIP(), []int{4}
}

func (x *DataLicenses) GetLicenses() []*users.License {
	if x != nil {
		return x.Licenses
	}
	return nil
}

var File_resources_sync_data_proto protoreflect.FileDescriptor

var file_resources_sync_data_proto_rawDesc = []byte{
	0x0a, 0x19, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x73, 0x79, 0x6e, 0x63,
	0x2f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x1a, 0x1a, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x6a, 0x6f, 0x62,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f,
	0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x2f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x22, 0x0a, 0x0a, 0x44, 0x61, 0x74, 0x61, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3f, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x4a, 0x6f, 0x62, 0x73,
	0x12, 0x33, 0x0a, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x2e, 0x4a, 0x6f, 0x62, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x92, 0x01, 0x03, 0x10, 0xc8, 0x01, 0x52,
	0x04, 0x6a, 0x6f, 0x62, 0x73, 0x22, 0x43, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x55, 0x73, 0x65,
	0x72, 0x73, 0x12, 0x36, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x92, 0x01, 0x03,
	0x10, 0xf4, 0x03, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x52, 0x0a, 0x0c, 0x44, 0x61,
	0x74, 0x61, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x12, 0x42, 0x0a, 0x08, 0x76, 0x65,
	0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65,
	0x73, 0x2e, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x92, 0x01,
	0x03, 0x10, 0xe8, 0x07, 0x52, 0x08, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x22, 0x4f,
	0x0a, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x12, 0x3f,
	0x0a, 0x08, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x4c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x92,
	0x01, 0x03, 0x10, 0xc8, 0x01, 0x52, 0x08, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x42,
	0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69,
	0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65,
	0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x3b, 0x73, 0x79,
	0x6e, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_sync_data_proto_rawDescOnce sync.Once
	file_resources_sync_data_proto_rawDescData = file_resources_sync_data_proto_rawDesc
)

func file_resources_sync_data_proto_rawDescGZIP() []byte {
	file_resources_sync_data_proto_rawDescOnce.Do(func() {
		file_resources_sync_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_sync_data_proto_rawDescData)
	})
	return file_resources_sync_data_proto_rawDescData
}

var file_resources_sync_data_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_resources_sync_data_proto_goTypes = []any{
	(*DataStatus)(nil),       // 0: resources.sync.DataStatus
	(*DataJobs)(nil),         // 1: resources.sync.DataJobs
	(*DataUsers)(nil),        // 2: resources.sync.DataUsers
	(*DataVehicles)(nil),     // 3: resources.sync.DataVehicles
	(*DataLicenses)(nil),     // 4: resources.sync.DataLicenses
	(*users.Job)(nil),        // 5: resources.users.Job
	(*users.User)(nil),       // 6: resources.users.User
	(*vehicles.Vehicle)(nil), // 7: resources.vehicles.Vehicle
	(*users.License)(nil),    // 8: resources.users.License
}
var file_resources_sync_data_proto_depIdxs = []int32{
	5, // 0: resources.sync.DataJobs.jobs:type_name -> resources.users.Job
	6, // 1: resources.sync.DataUsers.users:type_name -> resources.users.User
	7, // 2: resources.sync.DataVehicles.vehicles:type_name -> resources.vehicles.Vehicle
	8, // 3: resources.sync.DataLicenses.licenses:type_name -> resources.users.License
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_resources_sync_data_proto_init() }
func file_resources_sync_data_proto_init() {
	if File_resources_sync_data_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_sync_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_sync_data_proto_goTypes,
		DependencyIndexes: file_resources_sync_data_proto_depIdxs,
		MessageInfos:      file_resources_sync_data_proto_msgTypes,
	}.Build()
	File_resources_sync_data_proto = out.File
	file_resources_sync_data_proto_rawDesc = nil
	file_resources_sync_data_proto_goTypes = nil
	file_resources_sync_data_proto_depIdxs = nil
}
