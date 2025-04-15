// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: resources/sync/data.proto

package sync

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	livemap "github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	vehicles "github.com/fivenet-app/fivenet/gen/go/proto/resources/vehicles"
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

type DataUserLocations struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Users         []*UserLocation        `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	ClearAll      *bool                  `protobuf:"varint,2,opt,name=clear_all,json=clearAll,proto3,oneof" json:"clear_all,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DataUserLocations) Reset() {
	*x = DataUserLocations{}
	mi := &file_resources_sync_data_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataUserLocations) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUserLocations) ProtoMessage() {}

func (x *DataUserLocations) ProtoReflect() protoreflect.Message {
	mi := &file_resources_sync_data_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUserLocations.ProtoReflect.Descriptor instead.
func (*DataUserLocations) Descriptor() ([]byte, []int) {
	return file_resources_sync_data_proto_rawDescGZIP(), []int{5}
}

func (x *DataUserLocations) GetUsers() []*UserLocation {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *DataUserLocations) GetClearAll() bool {
	if x != nil && x.ClearAll != nil {
		return *x.ClearAll
	}
	return false
}

type UserLocation struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Identifier    string                 `protobuf:"bytes,1,opt,name=identifier,proto3" json:"identifier,omitempty"`
	Job           string                 `protobuf:"bytes,2,opt,name=job,proto3" json:"job,omitempty"`
	Coords        *livemap.Coords        `protobuf:"bytes,3,opt,name=coords,proto3" json:"coords,omitempty"`
	Hidden        bool                   `protobuf:"varint,4,opt,name=hidden,proto3" json:"hidden,omitempty"`
	Remove        bool                   `protobuf:"varint,5,opt,name=remove,proto3" json:"remove,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserLocation) Reset() {
	*x = UserLocation{}
	mi := &file_resources_sync_data_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserLocation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserLocation) ProtoMessage() {}

func (x *UserLocation) ProtoReflect() protoreflect.Message {
	mi := &file_resources_sync_data_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserLocation.ProtoReflect.Descriptor instead.
func (*UserLocation) Descriptor() ([]byte, []int) {
	return file_resources_sync_data_proto_rawDescGZIP(), []int{6}
}

func (x *UserLocation) GetIdentifier() string {
	if x != nil {
		return x.Identifier
	}
	return ""
}

func (x *UserLocation) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *UserLocation) GetCoords() *livemap.Coords {
	if x != nil {
		return x.Coords
	}
	return nil
}

func (x *UserLocation) GetHidden() bool {
	if x != nil {
		return x.Hidden
	}
	return false
}

func (x *UserLocation) GetRemove() bool {
	if x != nil {
		return x.Remove
	}
	return false
}

type DeleteUsers struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserIds       []int32                `protobuf:"varint,1,rep,packed,name=user_ids,json=userIds,proto3" json:"user_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteUsers) Reset() {
	*x = DeleteUsers{}
	mi := &file_resources_sync_data_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUsers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUsers) ProtoMessage() {}

func (x *DeleteUsers) ProtoReflect() protoreflect.Message {
	mi := &file_resources_sync_data_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUsers.ProtoReflect.Descriptor instead.
func (*DeleteUsers) Descriptor() ([]byte, []int) {
	return file_resources_sync_data_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteUsers) GetUserIds() []int32 {
	if x != nil {
		return x.UserIds
	}
	return nil
}

type DeleteVehicles struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Plates        []string               `protobuf:"bytes,1,rep,name=plates,proto3" json:"plates,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteVehicles) Reset() {
	*x = DeleteVehicles{}
	mi := &file_resources_sync_data_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteVehicles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteVehicles) ProtoMessage() {}

func (x *DeleteVehicles) ProtoReflect() protoreflect.Message {
	mi := &file_resources_sync_data_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteVehicles.ProtoReflect.Descriptor instead.
func (*DeleteVehicles) Descriptor() ([]byte, []int) {
	return file_resources_sync_data_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteVehicles) GetPlates() []string {
	if x != nil {
		return x.Plates
	}
	return nil
}

var File_resources_sync_data_proto protoreflect.FileDescriptor

const file_resources_sync_data_proto_rawDesc = "" +
	"\n" +
	"\x19resources/sync/data.proto\x12\x0eresources.sync\x1a\x1fresources/livemap/livemap.proto\x1a\x1aresources/users/jobs.proto\x1a\x1eresources/users/licenses.proto\x1a\x1bresources/users/users.proto\x1a!resources/vehicles/vehicles.proto\x1a\x17validate/validate.proto\"\"\n" +
	"\n" +
	"DataStatus\x12\x14\n" +
	"\x05count\x18\x01 \x01(\x03R\x05count\"?\n" +
	"\bDataJobs\x123\n" +
	"\x04jobs\x18\x01 \x03(\v2\x14.resources.users.JobB\t\xfaB\x06\x92\x01\x03\x10\xc8\x01R\x04jobs\"C\n" +
	"\tDataUsers\x126\n" +
	"\x05users\x18\x01 \x03(\v2\x15.resources.users.UserB\t\xfaB\x06\x92\x01\x03\x10\xf4\x03R\x05users\"R\n" +
	"\fDataVehicles\x12B\n" +
	"\bvehicles\x18\x01 \x03(\v2\x1b.resources.vehicles.VehicleB\t\xfaB\x06\x92\x01\x03\x10\xe8\aR\bvehicles\"O\n" +
	"\fDataLicenses\x12?\n" +
	"\blicenses\x18\x01 \x03(\v2\x18.resources.users.LicenseB\t\xfaB\x06\x92\x01\x03\x10\xc8\x01R\blicenses\"\x82\x01\n" +
	"\x11DataUserLocations\x12=\n" +
	"\x05users\x18\x01 \x03(\v2\x1c.resources.sync.UserLocationB\t\xfaB\x06\x92\x01\x03\x10\xd0\x0fR\x05users\x12 \n" +
	"\tclear_all\x18\x02 \x01(\bH\x00R\bclearAll\x88\x01\x01B\f\n" +
	"\n" +
	"_clear_all\"\xbf\x01\n" +
	"\fUserLocation\x12'\n" +
	"\n" +
	"identifier\x18\x01 \x01(\tB\a\xfaB\x04r\x02\x18@R\n" +
	"identifier\x12\x19\n" +
	"\x03job\x18\x02 \x01(\tB\a\xfaB\x04r\x02\x18\x14R\x03job\x12;\n" +
	"\x06coords\x18\x03 \x01(\v2\x19.resources.livemap.CoordsB\b\xfaB\x05\x8a\x01\x02\x10\x01R\x06coords\x12\x16\n" +
	"\x06hidden\x18\x04 \x01(\bR\x06hidden\x12\x16\n" +
	"\x06remove\x18\x05 \x01(\bR\x06remove\"2\n" +
	"\vDeleteUsers\x12#\n" +
	"\buser_ids\x18\x01 \x03(\x05B\b\xfaB\x05\x92\x01\x02\x10dR\auserIds\"2\n" +
	"\x0eDeleteVehicles\x12 \n" +
	"\x06plates\x18\x01 \x03(\tB\b\xfaB\x05\x92\x01\x02\x10dR\x06platesBAZ?github.com/fivenet-app/fivenet/gen/go/proto/resources/sync;syncb\x06proto3"

var (
	file_resources_sync_data_proto_rawDescOnce sync.Once
	file_resources_sync_data_proto_rawDescData []byte
)

func file_resources_sync_data_proto_rawDescGZIP() []byte {
	file_resources_sync_data_proto_rawDescOnce.Do(func() {
		file_resources_sync_data_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_sync_data_proto_rawDesc), len(file_resources_sync_data_proto_rawDesc)))
	})
	return file_resources_sync_data_proto_rawDescData
}

var file_resources_sync_data_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_resources_sync_data_proto_goTypes = []any{
	(*DataStatus)(nil),        // 0: resources.sync.DataStatus
	(*DataJobs)(nil),          // 1: resources.sync.DataJobs
	(*DataUsers)(nil),         // 2: resources.sync.DataUsers
	(*DataVehicles)(nil),      // 3: resources.sync.DataVehicles
	(*DataLicenses)(nil),      // 4: resources.sync.DataLicenses
	(*DataUserLocations)(nil), // 5: resources.sync.DataUserLocations
	(*UserLocation)(nil),      // 6: resources.sync.UserLocation
	(*DeleteUsers)(nil),       // 7: resources.sync.DeleteUsers
	(*DeleteVehicles)(nil),    // 8: resources.sync.DeleteVehicles
	(*users.Job)(nil),         // 9: resources.users.Job
	(*users.User)(nil),        // 10: resources.users.User
	(*vehicles.Vehicle)(nil),  // 11: resources.vehicles.Vehicle
	(*users.License)(nil),     // 12: resources.users.License
	(*livemap.Coords)(nil),    // 13: resources.livemap.Coords
}
var file_resources_sync_data_proto_depIdxs = []int32{
	9,  // 0: resources.sync.DataJobs.jobs:type_name -> resources.users.Job
	10, // 1: resources.sync.DataUsers.users:type_name -> resources.users.User
	11, // 2: resources.sync.DataVehicles.vehicles:type_name -> resources.vehicles.Vehicle
	12, // 3: resources.sync.DataLicenses.licenses:type_name -> resources.users.License
	6,  // 4: resources.sync.DataUserLocations.users:type_name -> resources.sync.UserLocation
	13, // 5: resources.sync.UserLocation.coords:type_name -> resources.livemap.Coords
	6,  // [6:6] is the sub-list for method output_type
	6,  // [6:6] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_resources_sync_data_proto_init() }
func file_resources_sync_data_proto_init() {
	if File_resources_sync_data_proto != nil {
		return
	}
	file_resources_sync_data_proto_msgTypes[5].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_sync_data_proto_rawDesc), len(file_resources_sync_data_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_sync_data_proto_goTypes,
		DependencyIndexes: file_resources_sync_data_proto_depIdxs,
		MessageInfos:      file_resources_sync_data_proto_msgTypes,
	}.Build()
	File_resources_sync_data_proto = out.File
	file_resources_sync_data_proto_goTypes = nil
	file_resources_sync_data_proto_depIdxs = nil
}
