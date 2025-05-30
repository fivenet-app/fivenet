// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: services/citizens/citizens.proto

package citizens

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	filestore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/filestore"
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

type ListCitizensRequest struct {
	state      protoimpl.MessageState      `protogen:"open.v1"`
	Pagination *database.PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Sort       *database.Sort              `protobuf:"bytes,2,opt,name=sort,proto3,oneof" json:"sort,omitempty"`
	// Search params
	Search                  string  `protobuf:"bytes,3,opt,name=search,proto3" json:"search,omitempty"`
	Wanted                  *bool   `protobuf:"varint,4,opt,name=wanted,proto3,oneof" json:"wanted,omitempty"`
	PhoneNumber             *string `protobuf:"bytes,5,opt,name=phone_number,json=phoneNumber,proto3,oneof" json:"phone_number,omitempty"`
	TrafficInfractionPoints *uint32 `protobuf:"varint,6,opt,name=traffic_infraction_points,json=trafficInfractionPoints,proto3,oneof" json:"traffic_infraction_points,omitempty"`
	Dateofbirth             *string `protobuf:"bytes,7,opt,name=dateofbirth,proto3,oneof" json:"dateofbirth,omitempty"`
	OpenFines               *uint64 `protobuf:"varint,8,opt,name=open_fines,json=openFines,proto3,oneof" json:"open_fines,omitempty"`
	unknownFields           protoimpl.UnknownFields
	sizeCache               protoimpl.SizeCache
}

func (x *ListCitizensRequest) Reset() {
	*x = ListCitizensRequest{}
	mi := &file_services_citizens_citizens_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCitizensRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCitizensRequest) ProtoMessage() {}

func (x *ListCitizensRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[0]
	if x != nil {
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
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{0}
}

func (x *ListCitizensRequest) GetPagination() *database.PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListCitizensRequest) GetSort() *database.Sort {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListCitizensRequest) GetSearch() string {
	if x != nil {
		return x.Search
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

func (x *ListCitizensRequest) GetTrafficInfractionPoints() uint32 {
	if x != nil && x.TrafficInfractionPoints != nil {
		return *x.TrafficInfractionPoints
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
	state         protoimpl.MessageState       `protogen:"open.v1"`
	Pagination    *database.PaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Users         []*users.User                `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCitizensResponse) Reset() {
	*x = ListCitizensResponse{}
	mi := &file_services_citizens_citizens_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCitizensResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCitizensResponse) ProtoMessage() {}

func (x *ListCitizensResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[1]
	if x != nil {
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
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{1}
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int32                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	InfoOnly      *bool                  `protobuf:"varint,2,opt,name=info_only,json=infoOnly,proto3,oneof" json:"info_only,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	mi := &file_services_citizens_citizens_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[2]
	if x != nil {
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
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetUserRequest) GetInfoOnly() bool {
	if x != nil && x.InfoOnly != nil {
		return *x.InfoOnly
	}
	return false
}

type GetUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *users.User            `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserResponse) Reset() {
	*x = GetUserResponse{}
	mi := &file_services_citizens_citizens_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserResponse) ProtoMessage() {}

func (x *GetUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[3]
	if x != nil {
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
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserResponse) GetUser() *users.User {
	if x != nil {
		return x.User
	}
	return nil
}

type ListUserActivityRequest struct {
	state      protoimpl.MessageState      `protogen:"open.v1"`
	Pagination *database.PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Sort       *database.Sort              `protobuf:"bytes,2,opt,name=sort,proto3,oneof" json:"sort,omitempty"`
	// Search params
	UserId        int32                    `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Types         []users.UserActivityType `protobuf:"varint,4,rep,packed,name=types,proto3,enum=resources.users.UserActivityType" json:"types,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListUserActivityRequest) Reset() {
	*x = ListUserActivityRequest{}
	mi := &file_services_citizens_citizens_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListUserActivityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserActivityRequest) ProtoMessage() {}

func (x *ListUserActivityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[4]
	if x != nil {
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
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{4}
}

func (x *ListUserActivityRequest) GetPagination() *database.PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListUserActivityRequest) GetSort() *database.Sort {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListUserActivityRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ListUserActivityRequest) GetTypes() []users.UserActivityType {
	if x != nil {
		return x.Types
	}
	return nil
}

type ListUserActivityResponse struct {
	state         protoimpl.MessageState       `protogen:"open.v1"`
	Pagination    *database.PaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Activity      []*users.UserActivity        `protobuf:"bytes,2,rep,name=activity,proto3" json:"activity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListUserActivityResponse) Reset() {
	*x = ListUserActivityResponse{}
	mi := &file_services_citizens_citizens_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListUserActivityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserActivityResponse) ProtoMessage() {}

func (x *ListUserActivityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[5]
	if x != nil {
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
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{5}
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
	state protoimpl.MessageState `protogen:"open.v1"`
	Props *users.UserProps       `protobuf:"bytes,1,opt,name=props,proto3" json:"props,omitempty"`
	// @sanitize
	Reason        string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetUserPropsRequest) Reset() {
	*x = SetUserPropsRequest{}
	mi := &file_services_citizens_citizens_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetUserPropsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUserPropsRequest) ProtoMessage() {}

func (x *SetUserPropsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[6]
	if x != nil {
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
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{6}
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	Props         *users.UserProps       `protobuf:"bytes,1,opt,name=props,proto3" json:"props,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetUserPropsResponse) Reset() {
	*x = SetUserPropsResponse{}
	mi := &file_services_citizens_citizens_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetUserPropsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUserPropsResponse) ProtoMessage() {}

func (x *SetUserPropsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[7]
	if x != nil {
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
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{7}
}

func (x *SetUserPropsResponse) GetProps() *users.UserProps {
	if x != nil {
		return x.Props
	}
	return nil
}

type SetProfilePictureRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Avatar        *filestore.File        `protobuf:"bytes,1,opt,name=avatar,proto3" json:"avatar,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetProfilePictureRequest) Reset() {
	*x = SetProfilePictureRequest{}
	mi := &file_services_citizens_citizens_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetProfilePictureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetProfilePictureRequest) ProtoMessage() {}

func (x *SetProfilePictureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetProfilePictureRequest.ProtoReflect.Descriptor instead.
func (*SetProfilePictureRequest) Descriptor() ([]byte, []int) {
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{8}
}

func (x *SetProfilePictureRequest) GetAvatar() *filestore.File {
	if x != nil {
		return x.Avatar
	}
	return nil
}

type SetProfilePictureResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Avatar        *filestore.File        `protobuf:"bytes,1,opt,name=avatar,proto3" json:"avatar,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetProfilePictureResponse) Reset() {
	*x = SetProfilePictureResponse{}
	mi := &file_services_citizens_citizens_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetProfilePictureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetProfilePictureResponse) ProtoMessage() {}

func (x *SetProfilePictureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetProfilePictureResponse.ProtoReflect.Descriptor instead.
func (*SetProfilePictureResponse) Descriptor() ([]byte, []int) {
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{9}
}

func (x *SetProfilePictureResponse) GetAvatar() *filestore.File {
	if x != nil {
		return x.Avatar
	}
	return nil
}

type ManageLabelsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Labels        []*users.Label         `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ManageLabelsRequest) Reset() {
	*x = ManageLabelsRequest{}
	mi := &file_services_citizens_citizens_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ManageLabelsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManageLabelsRequest) ProtoMessage() {}

func (x *ManageLabelsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManageLabelsRequest.ProtoReflect.Descriptor instead.
func (*ManageLabelsRequest) Descriptor() ([]byte, []int) {
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{10}
}

func (x *ManageLabelsRequest) GetLabels() []*users.Label {
	if x != nil {
		return x.Labels
	}
	return nil
}

type ManageLabelsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Labels        []*users.Label         `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ManageLabelsResponse) Reset() {
	*x = ManageLabelsResponse{}
	mi := &file_services_citizens_citizens_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ManageLabelsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManageLabelsResponse) ProtoMessage() {}

func (x *ManageLabelsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_citizens_citizens_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManageLabelsResponse.ProtoReflect.Descriptor instead.
func (*ManageLabelsResponse) Descriptor() ([]byte, []int) {
	return file_services_citizens_citizens_proto_rawDescGZIP(), []int{11}
}

func (x *ManageLabelsResponse) GetLabels() []*users.Label {
	if x != nil {
		return x.Labels
	}
	return nil
}

var File_services_citizens_citizens_proto protoreflect.FileDescriptor

const file_services_citizens_citizens_proto_rawDesc = "" +
	"\n" +
	" services/citizens/citizens.proto\x12\x11services.citizens\x1a(resources/common/database/database.proto\x1a\x1eresources/filestore/file.proto\x1a\x1eresources/users/activity.proto\x1a\x1cresources/users/labels.proto\x1a\x1bresources/users/users.proto\x1a\x1bresources/users/props.proto\x1a\x17validate/validate.proto\"\x8d\x04\n" +
	"\x13ListCitizensRequest\x12V\n" +
	"\n" +
	"pagination\x18\x01 \x01(\v2,.resources.common.database.PaginationRequestB\b\xfaB\x05\x8a\x01\x02\x10\x01R\n" +
	"pagination\x128\n" +
	"\x04sort\x18\x02 \x01(\v2\x1f.resources.common.database.SortH\x00R\x04sort\x88\x01\x01\x12\x1f\n" +
	"\x06search\x18\x03 \x01(\tB\a\xfaB\x04r\x02\x18@R\x06search\x12\x1b\n" +
	"\x06wanted\x18\x04 \x01(\bH\x01R\x06wanted\x88\x01\x01\x12/\n" +
	"\fphone_number\x18\x05 \x01(\tB\a\xfaB\x04r\x02\x18\x14H\x02R\vphoneNumber\x88\x01\x01\x12?\n" +
	"\x19traffic_infraction_points\x18\x06 \x01(\rH\x03R\x17trafficInfractionPoints\x88\x01\x01\x12.\n" +
	"\vdateofbirth\x18\a \x01(\tB\a\xfaB\x04r\x02\x18\n" +
	"H\x04R\vdateofbirth\x88\x01\x01\x12\"\n" +
	"\n" +
	"open_fines\x18\b \x01(\x04H\x05R\topenFines\x88\x01\x01B\a\n" +
	"\x05_sortB\t\n" +
	"\a_wantedB\x0f\n" +
	"\r_phone_numberB\x1c\n" +
	"\x1a_traffic_infraction_pointsB\x0e\n" +
	"\f_dateofbirthB\r\n" +
	"\v_open_fines\"\x92\x01\n" +
	"\x14ListCitizensResponse\x12M\n" +
	"\n" +
	"pagination\x18\x01 \x01(\v2-.resources.common.database.PaginationResponseR\n" +
	"pagination\x12+\n" +
	"\x05users\x18\x02 \x03(\v2\x15.resources.users.UserR\x05users\"b\n" +
	"\x0eGetUserRequest\x12 \n" +
	"\auser_id\x18\x01 \x01(\x05B\a\xfaB\x04\x1a\x02 \x00R\x06userId\x12 \n" +
	"\tinfo_only\x18\x02 \x01(\bH\x00R\binfoOnly\x88\x01\x01B\f\n" +
	"\n" +
	"_info_only\"<\n" +
	"\x0fGetUserResponse\x12)\n" +
	"\x04user\x18\x01 \x01(\v2\x15.resources.users.UserR\x04user\"\x99\x02\n" +
	"\x17ListUserActivityRequest\x12V\n" +
	"\n" +
	"pagination\x18\x01 \x01(\v2,.resources.common.database.PaginationRequestB\b\xfaB\x05\x8a\x01\x02\x10\x01R\n" +
	"pagination\x128\n" +
	"\x04sort\x18\x02 \x01(\v2\x1f.resources.common.database.SortH\x00R\x04sort\x88\x01\x01\x12 \n" +
	"\auser_id\x18\x03 \x01(\x05B\a\xfaB\x04\x1a\x02 \x00R\x06userId\x12A\n" +
	"\x05types\x18\x04 \x03(\x0e2!.resources.users.UserActivityTypeB\b\xfaB\x05\x92\x01\x02\x10\x14R\x05typesB\a\n" +
	"\x05_sort\"\xa4\x01\n" +
	"\x18ListUserActivityResponse\x12M\n" +
	"\n" +
	"pagination\x18\x01 \x01(\v2-.resources.common.database.PaginationResponseR\n" +
	"pagination\x129\n" +
	"\bactivity\x18\x02 \x03(\v2\x1d.resources.users.UserActivityR\bactivity\"x\n" +
	"\x13SetUserPropsRequest\x12:\n" +
	"\x05props\x18\x01 \x01(\v2\x1a.resources.users.UserPropsB\b\xfaB\x05\x8a\x01\x02\x10\x01R\x05props\x12%\n" +
	"\x06reason\x18\x02 \x01(\tB\r\xfaB\n" +
	"r\b\x10\x03\x18\xff\x01\xd0\x01\x01R\x06reason\"H\n" +
	"\x14SetUserPropsResponse\x120\n" +
	"\x05props\x18\x01 \x01(\v2\x1a.resources.users.UserPropsR\x05props\"M\n" +
	"\x18SetProfilePictureRequest\x121\n" +
	"\x06avatar\x18\x01 \x01(\v2\x19.resources.filestore.FileR\x06avatar\"N\n" +
	"\x19SetProfilePictureResponse\x121\n" +
	"\x06avatar\x18\x01 \x01(\v2\x19.resources.filestore.FileR\x06avatar\"E\n" +
	"\x13ManageLabelsRequest\x12.\n" +
	"\x06labels\x18\x01 \x03(\v2\x16.resources.users.LabelR\x06labels\"F\n" +
	"\x14ManageLabelsResponse\x12.\n" +
	"\x06labels\x18\x01 \x03(\v2\x16.resources.users.LabelR\x06labels2\xe3\x04\n" +
	"\x0fCitizensService\x12_\n" +
	"\fListCitizens\x12&.services.citizens.ListCitizensRequest\x1a'.services.citizens.ListCitizensResponse\x12P\n" +
	"\aGetUser\x12!.services.citizens.GetUserRequest\x1a\".services.citizens.GetUserResponse\x12k\n" +
	"\x10ListUserActivity\x12*.services.citizens.ListUserActivityRequest\x1a+.services.citizens.ListUserActivityResponse\x12_\n" +
	"\fSetUserProps\x12&.services.citizens.SetUserPropsRequest\x1a'.services.citizens.SetUserPropsResponse\x12n\n" +
	"\x11SetProfilePicture\x12+.services.citizens.SetProfilePictureRequest\x1a,.services.citizens.SetProfilePictureResponse\x12_\n" +
	"\fManageLabels\x12&.services.citizens.ManageLabelsRequest\x1a'.services.citizens.ManageLabelsResponseBNZLgithub.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens;citizensb\x06proto3"

var (
	file_services_citizens_citizens_proto_rawDescOnce sync.Once
	file_services_citizens_citizens_proto_rawDescData []byte
)

func file_services_citizens_citizens_proto_rawDescGZIP() []byte {
	file_services_citizens_citizens_proto_rawDescOnce.Do(func() {
		file_services_citizens_citizens_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_services_citizens_citizens_proto_rawDesc), len(file_services_citizens_citizens_proto_rawDesc)))
	})
	return file_services_citizens_citizens_proto_rawDescData
}

var file_services_citizens_citizens_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_services_citizens_citizens_proto_goTypes = []any{
	(*ListCitizensRequest)(nil),         // 0: services.citizens.ListCitizensRequest
	(*ListCitizensResponse)(nil),        // 1: services.citizens.ListCitizensResponse
	(*GetUserRequest)(nil),              // 2: services.citizens.GetUserRequest
	(*GetUserResponse)(nil),             // 3: services.citizens.GetUserResponse
	(*ListUserActivityRequest)(nil),     // 4: services.citizens.ListUserActivityRequest
	(*ListUserActivityResponse)(nil),    // 5: services.citizens.ListUserActivityResponse
	(*SetUserPropsRequest)(nil),         // 6: services.citizens.SetUserPropsRequest
	(*SetUserPropsResponse)(nil),        // 7: services.citizens.SetUserPropsResponse
	(*SetProfilePictureRequest)(nil),    // 8: services.citizens.SetProfilePictureRequest
	(*SetProfilePictureResponse)(nil),   // 9: services.citizens.SetProfilePictureResponse
	(*ManageLabelsRequest)(nil),         // 10: services.citizens.ManageLabelsRequest
	(*ManageLabelsResponse)(nil),        // 11: services.citizens.ManageLabelsResponse
	(*database.PaginationRequest)(nil),  // 12: resources.common.database.PaginationRequest
	(*database.Sort)(nil),               // 13: resources.common.database.Sort
	(*database.PaginationResponse)(nil), // 14: resources.common.database.PaginationResponse
	(*users.User)(nil),                  // 15: resources.users.User
	(users.UserActivityType)(0),         // 16: resources.users.UserActivityType
	(*users.UserActivity)(nil),          // 17: resources.users.UserActivity
	(*users.UserProps)(nil),             // 18: resources.users.UserProps
	(*filestore.File)(nil),              // 19: resources.filestore.File
	(*users.Label)(nil),                 // 20: resources.users.Label
}
var file_services_citizens_citizens_proto_depIdxs = []int32{
	12, // 0: services.citizens.ListCitizensRequest.pagination:type_name -> resources.common.database.PaginationRequest
	13, // 1: services.citizens.ListCitizensRequest.sort:type_name -> resources.common.database.Sort
	14, // 2: services.citizens.ListCitizensResponse.pagination:type_name -> resources.common.database.PaginationResponse
	15, // 3: services.citizens.ListCitizensResponse.users:type_name -> resources.users.User
	15, // 4: services.citizens.GetUserResponse.user:type_name -> resources.users.User
	12, // 5: services.citizens.ListUserActivityRequest.pagination:type_name -> resources.common.database.PaginationRequest
	13, // 6: services.citizens.ListUserActivityRequest.sort:type_name -> resources.common.database.Sort
	16, // 7: services.citizens.ListUserActivityRequest.types:type_name -> resources.users.UserActivityType
	14, // 8: services.citizens.ListUserActivityResponse.pagination:type_name -> resources.common.database.PaginationResponse
	17, // 9: services.citizens.ListUserActivityResponse.activity:type_name -> resources.users.UserActivity
	18, // 10: services.citizens.SetUserPropsRequest.props:type_name -> resources.users.UserProps
	18, // 11: services.citizens.SetUserPropsResponse.props:type_name -> resources.users.UserProps
	19, // 12: services.citizens.SetProfilePictureRequest.avatar:type_name -> resources.filestore.File
	19, // 13: services.citizens.SetProfilePictureResponse.avatar:type_name -> resources.filestore.File
	20, // 14: services.citizens.ManageLabelsRequest.labels:type_name -> resources.users.Label
	20, // 15: services.citizens.ManageLabelsResponse.labels:type_name -> resources.users.Label
	0,  // 16: services.citizens.CitizensService.ListCitizens:input_type -> services.citizens.ListCitizensRequest
	2,  // 17: services.citizens.CitizensService.GetUser:input_type -> services.citizens.GetUserRequest
	4,  // 18: services.citizens.CitizensService.ListUserActivity:input_type -> services.citizens.ListUserActivityRequest
	6,  // 19: services.citizens.CitizensService.SetUserProps:input_type -> services.citizens.SetUserPropsRequest
	8,  // 20: services.citizens.CitizensService.SetProfilePicture:input_type -> services.citizens.SetProfilePictureRequest
	10, // 21: services.citizens.CitizensService.ManageLabels:input_type -> services.citizens.ManageLabelsRequest
	1,  // 22: services.citizens.CitizensService.ListCitizens:output_type -> services.citizens.ListCitizensResponse
	3,  // 23: services.citizens.CitizensService.GetUser:output_type -> services.citizens.GetUserResponse
	5,  // 24: services.citizens.CitizensService.ListUserActivity:output_type -> services.citizens.ListUserActivityResponse
	7,  // 25: services.citizens.CitizensService.SetUserProps:output_type -> services.citizens.SetUserPropsResponse
	9,  // 26: services.citizens.CitizensService.SetProfilePicture:output_type -> services.citizens.SetProfilePictureResponse
	11, // 27: services.citizens.CitizensService.ManageLabels:output_type -> services.citizens.ManageLabelsResponse
	22, // [22:28] is the sub-list for method output_type
	16, // [16:22] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_services_citizens_citizens_proto_init() }
func file_services_citizens_citizens_proto_init() {
	if File_services_citizens_citizens_proto != nil {
		return
	}
	file_services_citizens_citizens_proto_msgTypes[0].OneofWrappers = []any{}
	file_services_citizens_citizens_proto_msgTypes[2].OneofWrappers = []any{}
	file_services_citizens_citizens_proto_msgTypes[4].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_services_citizens_citizens_proto_rawDesc), len(file_services_citizens_citizens_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_citizens_citizens_proto_goTypes,
		DependencyIndexes: file_services_citizens_citizens_proto_depIdxs,
		MessageInfos:      file_services_citizens_citizens_proto_msgTypes,
	}.Build()
	File_services_citizens_citizens_proto = out.File
	file_services_citizens_citizens_proto_goTypes = nil
	file_services_citizens_citizens_proto_depIdxs = nil
}
