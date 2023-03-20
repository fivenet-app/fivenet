// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: resources/vehicles/vehicles.proto

package vehicles

import (
	users "github.com/galexrt/arpanet/proto/resources/users"
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

type Vehicle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Plate string             `protobuf:"bytes,1,opt,name=plate,proto3" json:"plate,omitempty" alias:"plate"` // @gotags: alias:"plate"
	Model string             `protobuf:"bytes,2,opt,name=model,proto3" json:"model,omitempty" alias:"model"` // @gotags: alias:"model"
	Type  string             `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty" alias:"type"`   // @gotags: alias:"type"
	Owner *users.UserShortNI `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (x *Vehicle) Reset() {
	*x = Vehicle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_vehicles_vehicles_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vehicle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vehicle) ProtoMessage() {}

func (x *Vehicle) ProtoReflect() protoreflect.Message {
	mi := &file_resources_vehicles_vehicles_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *Vehicle) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Vehicle) GetOwner() *users.UserShortNI {
	if x != nil {
		return x.Owner
	}
	return nil
}

var File_resources_vehicles_vehicles_proto protoreflect.FileDescriptor

var file_resources_vehicles_vehicles_proto_rawDesc = []byte{
	0x0a, 0x21, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x76, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x73, 0x2f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x12, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x76,
	0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7d, 0x0a, 0x07, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x32, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x4e, 0x49, 0x52, 0x05, 0x6f, 0x77,
	0x6e, 0x65, 0x72, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x61, 0x72, 0x70, 0x61, 0x6e, 0x65,
	0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x3b, 0x76, 0x65, 0x68, 0x69, 0x63,
	0x6c, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_vehicles_vehicles_proto_rawDescOnce sync.Once
	file_resources_vehicles_vehicles_proto_rawDescData = file_resources_vehicles_vehicles_proto_rawDesc
)

func file_resources_vehicles_vehicles_proto_rawDescGZIP() []byte {
	file_resources_vehicles_vehicles_proto_rawDescOnce.Do(func() {
		file_resources_vehicles_vehicles_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_vehicles_vehicles_proto_rawDescData)
	})
	return file_resources_vehicles_vehicles_proto_rawDescData
}

var file_resources_vehicles_vehicles_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_vehicles_vehicles_proto_goTypes = []interface{}{
	(*Vehicle)(nil),           // 0: resources.vehicles.Vehicle
	(*users.UserShortNI)(nil), // 1: resources.users.UserShortNI
}
var file_resources_vehicles_vehicles_proto_depIdxs = []int32{
	1, // 0: resources.vehicles.Vehicle.owner:type_name -> resources.users.UserShortNI
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_resources_vehicles_vehicles_proto_init() }
func file_resources_vehicles_vehicles_proto_init() {
	if File_resources_vehicles_vehicles_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_vehicles_vehicles_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vehicle); i {
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
			RawDescriptor: file_resources_vehicles_vehicles_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_vehicles_vehicles_proto_goTypes,
		DependencyIndexes: file_resources_vehicles_vehicles_proto_depIdxs,
		MessageInfos:      file_resources_vehicles_vehicles_proto_msgTypes,
	}.Build()
	File_resources_vehicles_vehicles_proto = out.File
	file_resources_vehicles_vehicles_proto_rawDesc = nil
	file_resources_vehicles_vehicles_proto_goTypes = nil
	file_resources_vehicles_vehicles_proto_depIdxs = nil
}
