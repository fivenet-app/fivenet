// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: resources/vehicles/vehicles.proto

package vehicles

import (
	timestamp "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
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

type Vehicle struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Plate           string                 `protobuf:"bytes,1,opt,name=plate,proto3" json:"plate,omitempty"`
	Model           *string                `protobuf:"bytes,2,opt,name=model,proto3,oneof" json:"model,omitempty"`
	Type            string                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	OwnerId         *int32                 `protobuf:"varint,4,opt,name=owner_id,json=ownerId,proto3,oneof" json:"owner_id,omitempty"`
	OwnerIdentifier *string                `protobuf:"bytes,6,opt,name=owner_identifier,json=ownerIdentifier,proto3,oneof" json:"owner_identifier,omitempty"`
	Owner           *users.UserShort       `protobuf:"bytes,5,opt,name=owner,proto3,oneof" json:"owner,omitempty"`
	Job             *string                `protobuf:"bytes,7,opt,name=job,proto3,oneof" json:"job,omitempty"`
	JobLabel        *string                `protobuf:"bytes,8,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Vehicle) Reset() {
	*x = Vehicle{}
	mi := &file_resources_vehicles_vehicles_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Vehicle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vehicle) ProtoMessage() {}

func (x *Vehicle) ProtoReflect() protoreflect.Message {
	mi := &file_resources_vehicles_vehicles_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vehicle.ProtoReflect.Descriptor instead.
func (*Vehicle) Descriptor() ([]byte, []int) {
	return file_resources_vehicles_vehicles_proto_rawDescGZIP(), []int{0}
}

func (x *Vehicle) GetPlate() string {
	if x != nil {
		return x.Plate
	}
	return ""
}

func (x *Vehicle) GetModel() string {
	if x != nil && x.Model != nil {
		return *x.Model
	}
	return ""
}

func (x *Vehicle) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Vehicle) GetOwnerId() int32 {
	if x != nil && x.OwnerId != nil {
		return *x.OwnerId
	}
	return 0
}

func (x *Vehicle) GetOwnerIdentifier() string {
	if x != nil && x.OwnerIdentifier != nil {
		return *x.OwnerIdentifier
	}
	return ""
}

func (x *Vehicle) GetOwner() *users.UserShort {
	if x != nil {
		return x.Owner
	}
	return nil
}

func (x *Vehicle) GetJob() string {
	if x != nil && x.Job != nil {
		return *x.Job
	}
	return ""
}

func (x *Vehicle) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
	}
	return ""
}

type VehicleProps struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Plate         string                 `protobuf:"bytes,1,opt,name=plate,proto3" json:"plate,omitempty"`
	UpdatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	Wanted        *bool                  `protobuf:"varint,3,opt,name=wanted,proto3,oneof" json:"wanted,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VehicleProps) Reset() {
	*x = VehicleProps{}
	mi := &file_resources_vehicles_vehicles_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VehicleProps) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VehicleProps) ProtoMessage() {}

func (x *VehicleProps) ProtoReflect() protoreflect.Message {
	mi := &file_resources_vehicles_vehicles_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VehicleProps.ProtoReflect.Descriptor instead.
func (*VehicleProps) Descriptor() ([]byte, []int) {
	return file_resources_vehicles_vehicles_proto_rawDescGZIP(), []int{1}
}

func (x *VehicleProps) GetPlate() string {
	if x != nil {
		return x.Plate
	}
	return ""
}

func (x *VehicleProps) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *VehicleProps) GetWanted() bool {
	if x != nil && x.Wanted != nil {
		return *x.Wanted
	}
	return false
}

var File_resources_vehicles_vehicles_proto protoreflect.FileDescriptor

const file_resources_vehicles_vehicles_proto_rawDesc = "" +
	"\n" +
	"!resources/vehicles/vehicles.proto\x12\x12resources.vehicles\x1a#resources/timestamp/timestamp.proto\x1a\x1bresources/users/users.proto\"\x90\x03\n" +
	"\aVehicle\x12\x1d\n" +
	"\x05plate\x18\x01 \x01(\tB\a\xbaH\x04r\x02\x18 R\x05plate\x12\"\n" +
	"\x05model\x18\x02 \x01(\tB\a\xbaH\x04r\x02\x18@H\x00R\x05model\x88\x01\x01\x12\x1b\n" +
	"\x04type\x18\x03 \x01(\tB\a\xbaH\x04r\x02\x18 R\x04type\x12\x1e\n" +
	"\bowner_id\x18\x04 \x01(\x05H\x01R\aownerId\x88\x01\x01\x127\n" +
	"\x10owner_identifier\x18\x06 \x01(\tB\a\xbaH\x04r\x02\x18@H\x02R\x0fownerIdentifier\x88\x01\x01\x125\n" +
	"\x05owner\x18\x05 \x01(\v2\x1a.resources.users.UserShortH\x03R\x05owner\x88\x01\x01\x12\x1e\n" +
	"\x03job\x18\a \x01(\tB\a\xbaH\x04r\x02\x18\x14H\x04R\x03job\x88\x01\x01\x12)\n" +
	"\tjob_label\x18\b \x01(\tB\a\xbaH\x04r\x02\x182H\x05R\bjobLabel\x88\x01\x01B\b\n" +
	"\x06_modelB\v\n" +
	"\t_owner_idB\x13\n" +
	"\x11_owner_identifierB\b\n" +
	"\x06_ownerB\x06\n" +
	"\x04_jobB\f\n" +
	"\n" +
	"_job_label\"\xa8\x01\n" +
	"\fVehicleProps\x12\x1d\n" +
	"\x05plate\x18\x01 \x01(\tB\a\xbaH\x04r\x02\x18 R\x05plate\x12B\n" +
	"\n" +
	"updated_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tupdatedAt\x88\x01\x01\x12\x1b\n" +
	"\x06wanted\x18\x03 \x01(\bH\x01R\x06wanted\x88\x01\x01B\r\n" +
	"\v_updated_atB\t\n" +
	"\a_wantedBOZMgithub.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/vehicles;vehiclesb\x06proto3"

var (
	file_resources_vehicles_vehicles_proto_rawDescOnce sync.Once
	file_resources_vehicles_vehicles_proto_rawDescData []byte
)

func file_resources_vehicles_vehicles_proto_rawDescGZIP() []byte {
	file_resources_vehicles_vehicles_proto_rawDescOnce.Do(func() {
		file_resources_vehicles_vehicles_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_vehicles_vehicles_proto_rawDesc), len(file_resources_vehicles_vehicles_proto_rawDesc)))
	})
	return file_resources_vehicles_vehicles_proto_rawDescData
}

var file_resources_vehicles_vehicles_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_resources_vehicles_vehicles_proto_goTypes = []any{
	(*Vehicle)(nil),             // 0: resources.vehicles.Vehicle
	(*VehicleProps)(nil),        // 1: resources.vehicles.VehicleProps
	(*users.UserShort)(nil),     // 2: resources.users.UserShort
	(*timestamp.Timestamp)(nil), // 3: resources.timestamp.Timestamp
}
var file_resources_vehicles_vehicles_proto_depIdxs = []int32{
	2, // 0: resources.vehicles.Vehicle.owner:type_name -> resources.users.UserShort
	3, // 1: resources.vehicles.VehicleProps.updated_at:type_name -> resources.timestamp.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_resources_vehicles_vehicles_proto_init() }
func file_resources_vehicles_vehicles_proto_init() {
	if File_resources_vehicles_vehicles_proto != nil {
		return
	}
	file_resources_vehicles_vehicles_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_vehicles_vehicles_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_vehicles_vehicles_proto_rawDesc), len(file_resources_vehicles_vehicles_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_vehicles_vehicles_proto_goTypes,
		DependencyIndexes: file_resources_vehicles_vehicles_proto_depIdxs,
		MessageInfos:      file_resources_vehicles_vehicles_proto_msgTypes,
	}.Build()
	File_resources_vehicles_vehicles_proto = out.File
	file_resources_vehicles_vehicles_proto_goTypes = nil
	file_resources_vehicles_vehicles_proto_depIdxs = nil
}
