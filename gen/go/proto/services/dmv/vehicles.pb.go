// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: services/dmv/vehicles.proto

package dmv

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
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

type ListVehiclesRequest struct {
	state      protoimpl.MessageState      `protogen:"open.v1"`
	Pagination *database.PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Sort       *database.Sort              `protobuf:"bytes,2,opt,name=sort,proto3,oneof" json:"sort,omitempty"`
	// Search params
	LicensePlate  *string `protobuf:"bytes,3,opt,name=license_plate,json=licensePlate,proto3,oneof" json:"license_plate,omitempty"`
	Model         *string `protobuf:"bytes,4,opt,name=model,proto3,oneof" json:"model,omitempty"`
	UserId        *int32  `protobuf:"varint,5,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	Job           *string `protobuf:"bytes,6,opt,name=job,proto3,oneof" json:"job,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListVehiclesRequest) Reset() {
	*x = ListVehiclesRequest{}
	mi := &file_services_dmv_vehicles_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListVehiclesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListVehiclesRequest) ProtoMessage() {}

func (x *ListVehiclesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_dmv_vehicles_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListVehiclesRequest.ProtoReflect.Descriptor instead.
func (*ListVehiclesRequest) Descriptor() ([]byte, []int) {
	return file_services_dmv_vehicles_proto_rawDescGZIP(), []int{0}
}

func (x *ListVehiclesRequest) GetPagination() *database.PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListVehiclesRequest) GetSort() *database.Sort {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListVehiclesRequest) GetLicensePlate() string {
	if x != nil && x.LicensePlate != nil {
		return *x.LicensePlate
	}
	return ""
}

func (x *ListVehiclesRequest) GetModel() string {
	if x != nil && x.Model != nil {
		return *x.Model
	}
	return ""
}

func (x *ListVehiclesRequest) GetUserId() int32 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

func (x *ListVehiclesRequest) GetJob() string {
	if x != nil && x.Job != nil {
		return *x.Job
	}
	return ""
}

type ListVehiclesResponse struct {
	state         protoimpl.MessageState       `protogen:"open.v1"`
	Pagination    *database.PaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Vehicles      []*vehicles.Vehicle          `protobuf:"bytes,2,rep,name=vehicles,proto3" json:"vehicles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListVehiclesResponse) Reset() {
	*x = ListVehiclesResponse{}
	mi := &file_services_dmv_vehicles_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListVehiclesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListVehiclesResponse) ProtoMessage() {}

func (x *ListVehiclesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_dmv_vehicles_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListVehiclesResponse.ProtoReflect.Descriptor instead.
func (*ListVehiclesResponse) Descriptor() ([]byte, []int) {
	return file_services_dmv_vehicles_proto_rawDescGZIP(), []int{1}
}

func (x *ListVehiclesResponse) GetPagination() *database.PaginationResponse {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListVehiclesResponse) GetVehicles() []*vehicles.Vehicle {
	if x != nil {
		return x.Vehicles
	}
	return nil
}

var File_services_dmv_vehicles_proto protoreflect.FileDescriptor

var file_services_dmv_vehicles_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x6d, 0x76, 0x2f, 0x76,
	0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x6d, 0x76, 0x1a, 0x28, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x64, 0x61,
	0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x2f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xfe, 0x02, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x56, 0x0a, 0x0a, 0x70, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x08, 0xfa, 0x42, 0x05,
	0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x38, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x53, 0x6f, 0x72, 0x74,
	0x48, 0x00, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x88, 0x01, 0x01, 0x12, 0x31, 0x0a, 0x0d, 0x6c,
	0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x5f, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x20, 0x48, 0x01, 0x52, 0x0c, 0x6c,
	0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x50, 0x6c, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x22,
	0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x18, 0x20, 0x48, 0x02, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x88,
	0x01, 0x01, 0x12, 0x25, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x28, 0x00, 0x48, 0x03, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x03, 0x6a, 0x6f, 0x62,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x48,
	0x04, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x73, 0x6f,
	0x72, 0x74, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x5f, 0x70,
	0x6c, 0x61, 0x74, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x42, 0x0a,
	0x0a, 0x08, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6a,
	0x6f, 0x62, 0x22, 0x9e, 0x01, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x65, 0x68, 0x69, 0x63,
	0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x70,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0a,
	0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x37, 0x0a, 0x08, 0x76, 0x65,
	0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65,
	0x73, 0x2e, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x08, 0x76, 0x65, 0x68, 0x69, 0x63,
	0x6c, 0x65, 0x73, 0x32, 0x63, 0x0a, 0x0a, 0x44, 0x4d, 0x56, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x55, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65,
	0x73, 0x12, 0x21, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x6d, 0x76,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x64, 0x6d, 0x76, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61,
	0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2f, 0x64, 0x6d, 0x76, 0x3b, 0x64, 0x6d, 0x76, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_services_dmv_vehicles_proto_rawDescOnce sync.Once
	file_services_dmv_vehicles_proto_rawDescData []byte
)

func file_services_dmv_vehicles_proto_rawDescGZIP() []byte {
	file_services_dmv_vehicles_proto_rawDescOnce.Do(func() {
		file_services_dmv_vehicles_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_services_dmv_vehicles_proto_rawDesc), len(file_services_dmv_vehicles_proto_rawDesc)))
	})
	return file_services_dmv_vehicles_proto_rawDescData
}

var file_services_dmv_vehicles_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_services_dmv_vehicles_proto_goTypes = []any{
	(*ListVehiclesRequest)(nil),         // 0: services.dmv.ListVehiclesRequest
	(*ListVehiclesResponse)(nil),        // 1: services.dmv.ListVehiclesResponse
	(*database.PaginationRequest)(nil),  // 2: resources.common.database.PaginationRequest
	(*database.Sort)(nil),               // 3: resources.common.database.Sort
	(*database.PaginationResponse)(nil), // 4: resources.common.database.PaginationResponse
	(*vehicles.Vehicle)(nil),            // 5: resources.vehicles.Vehicle
}
var file_services_dmv_vehicles_proto_depIdxs = []int32{
	2, // 0: services.dmv.ListVehiclesRequest.pagination:type_name -> resources.common.database.PaginationRequest
	3, // 1: services.dmv.ListVehiclesRequest.sort:type_name -> resources.common.database.Sort
	4, // 2: services.dmv.ListVehiclesResponse.pagination:type_name -> resources.common.database.PaginationResponse
	5, // 3: services.dmv.ListVehiclesResponse.vehicles:type_name -> resources.vehicles.Vehicle
	0, // 4: services.dmv.DMVService.ListVehicles:input_type -> services.dmv.ListVehiclesRequest
	1, // 5: services.dmv.DMVService.ListVehicles:output_type -> services.dmv.ListVehiclesResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_services_dmv_vehicles_proto_init() }
func file_services_dmv_vehicles_proto_init() {
	if File_services_dmv_vehicles_proto != nil {
		return
	}
	file_services_dmv_vehicles_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_services_dmv_vehicles_proto_rawDesc), len(file_services_dmv_vehicles_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_dmv_vehicles_proto_goTypes,
		DependencyIndexes: file_services_dmv_vehicles_proto_depIdxs,
		MessageInfos:      file_services_dmv_vehicles_proto_msgTypes,
	}.Build()
	File_services_dmv_vehicles_proto = out.File
	file_services_dmv_vehicles_proto_goTypes = nil
	file_services_dmv_vehicles_proto_depIdxs = nil
}
