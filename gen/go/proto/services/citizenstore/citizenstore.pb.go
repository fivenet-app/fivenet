// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.2
// source: services/citizenstore/citizenstore.proto

package citizenstore

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
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

type ListCitizensRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *database.PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	// Search params
	SearchName    string  `protobuf:"bytes,2,opt,name=search_name,json=searchName,proto3" json:"search_name,omitempty"`
	Wanted        *bool   `protobuf:"varint,3,opt,name=wanted,proto3,oneof" json:"wanted,omitempty"`
	PhoneNumber   *string `protobuf:"bytes,4,opt,name=phone_number,json=phoneNumber,proto3,oneof" json:"phone_number,omitempty"`
	TrafficPoints *uint64 `protobuf:"varint,5,opt,name=traffic_points,json=trafficPoints,proto3,oneof" json:"traffic_points,omitempty"`
	Dateofbirth   *string `protobuf:"bytes,6,opt,name=dateofbirth,proto3,oneof" json:"dateofbirth,omitempty"`
	OpenFines     *uint64 `protobuf:"varint,7,opt,name=open_fines,json=openFines,proto3,oneof" json:"open_fines,omitempty"`
}

func (x *ListCitizensRequest) Reset() {
	*x = ListCitizensRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_citizenstore_citizenstore_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCitizensRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCitizensRequest) ProtoMessage() {}

func (x *ListCitizensRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizenstore_citizenstore_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCitizensRequest.ProtoReflect.Descriptor instead.
func (*ListCitizensRequest) Descriptor() ([]byte, []int) {
	return file_services_citizenstore_citizenstore_proto_rawDescGZIP(), []int{0}
}

func (x *ListCitizensRequest) GetPagination() *database.PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListCitizensRequest) GetSearchName() string {
	if x != nil {
		return x.SearchName
	}
	return ""
}

func (x *ListCitizensRequest) GetWanted() bool {
	if x != nil && x.Wanted != nil {
		return *x.Wanted
	}
	return false
}

func (x *ListCitizensRequest) GetPhoneNumber() string {
	if x != nil && x.PhoneNumber != nil {
		return *x.PhoneNumber
	}
	return ""
}

func (x *ListCitizensRequest) GetTrafficPoints() uint64 {
	if x != nil && x.TrafficPoints != nil {
		return *x.TrafficPoints
	}
	return 0
}

func (x *ListCitizensRequest) GetDateofbirth() string {
	if x != nil && x.Dateofbirth != nil {
		return *x.Dateofbirth
	}
	return ""
}

func (x *ListCitizensRequest) GetOpenFines() uint64 {
	if x != nil && x.OpenFines != nil {
		return *x.OpenFines
	}
	return 0
}

type ListCitizensResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *database.PaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Users      []*users.User                `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *ListCitizensResponse) Reset() {
	*x = ListCitizensResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_citizenstore_citizenstore_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCitizensResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCitizensResponse) ProtoMessage() {}

func (x *ListCitizensResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizenstore_citizenstore_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCitizensResponse.ProtoReflect.Descriptor instead.
func (*ListCitizensResponse) Descriptor() ([]byte, []int) {
	return file_services_citizenstore_citizenstore_proto_rawDescGZIP(), []int{1}
}

func (x *ListCitizensResponse) GetPagination() *database.PaginationResponse {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListCitizensResponse) GetUsers() []*users.User {
	if x != nil {
		return x.Users
	}
	return nil
}

type GetUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_citizenstore_citizenstore_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizenstore_citizenstore_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserRequest.ProtoReflect.Descriptor instead.
func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return file_services_citizenstore_citizenstore_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type GetUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *users.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *GetUserResponse) Reset() {
	*x = GetUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_citizenstore_citizenstore_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserResponse) ProtoMessage() {}

func (x *GetUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizenstore_citizenstore_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserResponse.ProtoReflect.Descriptor instead.
func (*GetUserResponse) Descriptor() ([]byte, []int) {
	return file_services_citizenstore_citizenstore_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserResponse) GetUser() *users.User {
	if x != nil {
		return x.User
	}
	return nil
}

type ListUserActivityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *database.PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	UserId     int32                       `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *ListUserActivityRequest) Reset() {
	*x = ListUserActivityRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_citizenstore_citizenstore_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUserActivityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserActivityRequest) ProtoMessage() {}

func (x *ListUserActivityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizenstore_citizenstore_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUserActivityRequest.ProtoReflect.Descriptor instead.
func (*ListUserActivityRequest) Descriptor() ([]byte, []int) {
	return file_services_citizenstore_citizenstore_proto_rawDescGZIP(), []int{4}
}

func (x *ListUserActivityRequest) GetPagination() *database.PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListUserActivityRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ListUserActivityResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *database.PaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Activity   []*users.UserActivity        `protobuf:"bytes,2,rep,name=activity,proto3" json:"activity,omitempty"`
}

func (x *ListUserActivityResponse) Reset() {
	*x = ListUserActivityResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_citizenstore_citizenstore_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUserActivityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserActivityResponse) ProtoMessage() {}

func (x *ListUserActivityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizenstore_citizenstore_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUserActivityResponse.ProtoReflect.Descriptor instead.
func (*ListUserActivityResponse) Descriptor() ([]byte, []int) {
	return file_services_citizenstore_citizenstore_proto_rawDescGZIP(), []int{5}
}

func (x *ListUserActivityResponse) GetPagination() *database.PaginationResponse {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListUserActivityResponse) GetActivity() []*users.UserActivity {
	if x != nil {
		return x.Activity
	}
	return nil
}

type SetUserPropsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Props *users.UserProps `protobuf:"bytes,1,opt,name=props,proto3" json:"props,omitempty"`
	// @sanitize
	Reason string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *SetUserPropsRequest) Reset() {
	*x = SetUserPropsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_citizenstore_citizenstore_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetUserPropsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUserPropsRequest) ProtoMessage() {}

func (x *SetUserPropsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizenstore_citizenstore_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetUserPropsRequest.ProtoReflect.Descriptor instead.
func (*SetUserPropsRequest) Descriptor() ([]byte, []int) {
	return file_services_citizenstore_citizenstore_proto_rawDescGZIP(), []int{6}
}

func (x *SetUserPropsRequest) GetProps() *users.UserProps {
	if x != nil {
		return x.Props
	}
	return nil
}

func (x *SetUserPropsRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type SetUserPropsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Props *users.UserProps `protobuf:"bytes,1,opt,name=props,proto3" json:"props,omitempty"`
}

func (x *SetUserPropsResponse) Reset() {
	*x = SetUserPropsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_citizenstore_citizenstore_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetUserPropsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUserPropsResponse) ProtoMessage() {}

func (x *SetUserPropsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizenstore_citizenstore_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetUserPropsResponse.ProtoReflect.Descriptor instead.
func (*SetUserPropsResponse) Descriptor() ([]byte, []int) {
	return file_services_citizenstore_citizenstore_proto_rawDescGZIP(), []int{7}
}

func (x *SetUserPropsResponse) GetProps() *users.UserProps {
	if x != nil {
		return x.Props
	}
	return nil
}

var File_services_citizenstore_citizenstore_proto protoreflect.FileDescriptor

var file_services_citizenstore_citizenstore_proto_rawDesc = []byte{
	0x0a, 0x28, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x69, 0x74, 0x69, 0x7a,
	0x65, 0x6e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x1a, 0x28, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x64, 0x61, 0x74,
	0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xb3, 0x03, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65,
	0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x56, 0x0a, 0x0a, 0x70, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x08, 0xfa, 0x42, 0x05,
	0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x28, 0x0a, 0x0b, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x52,
	0x0a, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x06, 0x77,
	0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x06, 0x77,
	0x61, 0x6e, 0x74, 0x65, 0x64, 0x88, 0x01, 0x01, 0x12, 0x2f, 0x0a, 0x0c, 0x70, 0x68, 0x6f, 0x6e,
	0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x48, 0x01, 0x52, 0x0b, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x2a, 0x0a, 0x0e, 0x74, 0x72, 0x61,
	0x66, 0x66, 0x69, 0x63, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x48, 0x02, 0x52, 0x0d, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x50, 0x6f, 0x69, 0x6e,
	0x74, 0x73, 0x88, 0x01, 0x01, 0x12, 0x2e, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x65, 0x6f, 0x66, 0x62,
	0x69, 0x72, 0x74, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72,
	0x02, 0x18, 0x0a, 0x48, 0x03, 0x52, 0x0b, 0x64, 0x61, 0x74, 0x65, 0x6f, 0x66, 0x62, 0x69, 0x72,
	0x74, 0x68, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x6f, 0x70, 0x65, 0x6e, 0x5f, 0x66, 0x69,
	0x6e, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x48, 0x04, 0x52, 0x09, 0x6f, 0x70, 0x65,
	0x6e, 0x46, 0x69, 0x6e, 0x65, 0x73, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x77, 0x61,
	0x6e, 0x74, 0x65, 0x64, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x5f, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69,
	0x63, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x61, 0x74,
	0x65, 0x6f, 0x66, 0x62, 0x69, 0x72, 0x74, 0x68, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x6f, 0x70, 0x65,
	0x6e, 0x5f, 0x66, 0x69, 0x6e, 0x65, 0x73, 0x22, 0x92, 0x01, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74,
	0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x4d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x2b, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x32, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x3c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x93,
	0x01, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76,
	0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x56, 0x0a, 0x0a, 0x70, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x08, 0xfa, 0x42,
	0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x22, 0xa4, 0x01, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73,
	0x65, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x39, 0x0a, 0x08, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74,
	0x79, 0x52, 0x08, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x22, 0x78, 0x0a, 0x13, 0x53,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x3a, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x25,
	0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d,
	0xfa, 0x42, 0x0a, 0x72, 0x08, 0x10, 0x03, 0x18, 0xff, 0x01, 0xd0, 0x01, 0x01, 0x52, 0x06, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x48, 0x0a, 0x14, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x50, 0x72, 0x6f, 0x70, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a,
	0x05, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x32,
	0xb6, 0x03, 0x0a, 0x13, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x67, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x43,
	0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x12, 0x2a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x63, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63,
	0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x58, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x25, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x26, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x69,
	0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x73, 0x0a, 0x10, 0x4c, 0x69,
	0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x12, 0x2e,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x69, 0x74, 0x69, 0x7a, 0x65,
	0x6e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x69, 0x74, 0x69, 0x7a, 0x65,
	0x6e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x67, 0x0a, 0x0c, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x12,
	0x2a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x69, 0x74, 0x69, 0x7a,
	0x65, 0x6e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x50,
	0x72, 0x6f, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4c, 0x5a, 0x4a, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x66,
	0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x69, 0x74,
	0x69, 0x7a, 0x65, 0x6e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x3b, 0x63, 0x69, 0x74, 0x69, 0x7a, 0x65,
	0x6e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_citizenstore_citizenstore_proto_rawDescOnce sync.Once
	file_services_citizenstore_citizenstore_proto_rawDescData = file_services_citizenstore_citizenstore_proto_rawDesc
)

func file_services_citizenstore_citizenstore_proto_rawDescGZIP() []byte {
	file_services_citizenstore_citizenstore_proto_rawDescOnce.Do(func() {
		file_services_citizenstore_citizenstore_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_citizenstore_citizenstore_proto_rawDescData)
	})
	return file_services_citizenstore_citizenstore_proto_rawDescData
}

var file_services_citizenstore_citizenstore_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_services_citizenstore_citizenstore_proto_goTypes = []interface{}{
	(*ListCitizensRequest)(nil),         // 0: services.citizenstore.ListCitizensRequest
	(*ListCitizensResponse)(nil),        // 1: services.citizenstore.ListCitizensResponse
	(*GetUserRequest)(nil),              // 2: services.citizenstore.GetUserRequest
	(*GetUserResponse)(nil),             // 3: services.citizenstore.GetUserResponse
	(*ListUserActivityRequest)(nil),     // 4: services.citizenstore.ListUserActivityRequest
	(*ListUserActivityResponse)(nil),    // 5: services.citizenstore.ListUserActivityResponse
	(*SetUserPropsRequest)(nil),         // 6: services.citizenstore.SetUserPropsRequest
	(*SetUserPropsResponse)(nil),        // 7: services.citizenstore.SetUserPropsResponse
	(*database.PaginationRequest)(nil),  // 8: resources.common.database.PaginationRequest
	(*database.PaginationResponse)(nil), // 9: resources.common.database.PaginationResponse
	(*users.User)(nil),                  // 10: resources.users.User
	(*users.UserActivity)(nil),          // 11: resources.users.UserActivity
	(*users.UserProps)(nil),             // 12: resources.users.UserProps
}
var file_services_citizenstore_citizenstore_proto_depIdxs = []int32{
	8,  // 0: services.citizenstore.ListCitizensRequest.pagination:type_name -> resources.common.database.PaginationRequest
	9,  // 1: services.citizenstore.ListCitizensResponse.pagination:type_name -> resources.common.database.PaginationResponse
	10, // 2: services.citizenstore.ListCitizensResponse.users:type_name -> resources.users.User
	10, // 3: services.citizenstore.GetUserResponse.user:type_name -> resources.users.User
	8,  // 4: services.citizenstore.ListUserActivityRequest.pagination:type_name -> resources.common.database.PaginationRequest
	9,  // 5: services.citizenstore.ListUserActivityResponse.pagination:type_name -> resources.common.database.PaginationResponse
	11, // 6: services.citizenstore.ListUserActivityResponse.activity:type_name -> resources.users.UserActivity
	12, // 7: services.citizenstore.SetUserPropsRequest.props:type_name -> resources.users.UserProps
	12, // 8: services.citizenstore.SetUserPropsResponse.props:type_name -> resources.users.UserProps
	0,  // 9: services.citizenstore.CitizenStoreService.ListCitizens:input_type -> services.citizenstore.ListCitizensRequest
	2,  // 10: services.citizenstore.CitizenStoreService.GetUser:input_type -> services.citizenstore.GetUserRequest
	4,  // 11: services.citizenstore.CitizenStoreService.ListUserActivity:input_type -> services.citizenstore.ListUserActivityRequest
	6,  // 12: services.citizenstore.CitizenStoreService.SetUserProps:input_type -> services.citizenstore.SetUserPropsRequest
	1,  // 13: services.citizenstore.CitizenStoreService.ListCitizens:output_type -> services.citizenstore.ListCitizensResponse
	3,  // 14: services.citizenstore.CitizenStoreService.GetUser:output_type -> services.citizenstore.GetUserResponse
	5,  // 15: services.citizenstore.CitizenStoreService.ListUserActivity:output_type -> services.citizenstore.ListUserActivityResponse
	7,  // 16: services.citizenstore.CitizenStoreService.SetUserProps:output_type -> services.citizenstore.SetUserPropsResponse
	13, // [13:17] is the sub-list for method output_type
	9,  // [9:13] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_services_citizenstore_citizenstore_proto_init() }
func file_services_citizenstore_citizenstore_proto_init() {
	if File_services_citizenstore_citizenstore_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_citizenstore_citizenstore_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCitizensRequest); i {
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
		file_services_citizenstore_citizenstore_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCitizensResponse); i {
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
		file_services_citizenstore_citizenstore_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserRequest); i {
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
		file_services_citizenstore_citizenstore_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserResponse); i {
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
		file_services_citizenstore_citizenstore_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUserActivityRequest); i {
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
		file_services_citizenstore_citizenstore_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUserActivityResponse); i {
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
		file_services_citizenstore_citizenstore_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetUserPropsRequest); i {
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
		file_services_citizenstore_citizenstore_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetUserPropsResponse); i {
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
	file_services_citizenstore_citizenstore_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_citizenstore_citizenstore_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_citizenstore_citizenstore_proto_goTypes,
		DependencyIndexes: file_services_citizenstore_citizenstore_proto_depIdxs,
		MessageInfos:      file_services_citizenstore_citizenstore_proto_msgTypes,
	}.Build()
	File_services_citizenstore_citizenstore_proto = out.File
	file_services_citizenstore_citizenstore_proto_rawDesc = nil
	file_services_citizenstore_citizenstore_proto_goTypes = nil
	file_services_citizenstore_citizenstore_proto_depIdxs = nil
}
