// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: users/users.proto

package users

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	common "github.com/galexrt/arpanet/proto/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FindUsersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset    int64             `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	OrderBy   []*common.OrderBy `protobuf:"bytes,2,rep,name=orderBy,proto3" json:"orderBy,omitempty"`
	Firstname string            `protobuf:"bytes,3,opt,name=firstname,proto3" json:"firstname,omitempty"`
	Lastname  string            `protobuf:"bytes,4,opt,name=lastname,proto3" json:"lastname,omitempty"`
	Wanted    bool              `protobuf:"varint,5,opt,name=wanted,proto3" json:"wanted,omitempty"`
}

func (x *FindUsersRequest) Reset() {
	*x = FindUsersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_users_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindUsersRequest) ProtoMessage() {}

func (x *FindUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindUsersRequest.ProtoReflect.Descriptor instead.
func (*FindUsersRequest) Descriptor() ([]byte, []int) {
	return file_users_users_proto_rawDescGZIP(), []int{0}
}

func (x *FindUsersRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *FindUsersRequest) GetOrderBy() []*common.OrderBy {
	if x != nil {
		return x.OrderBy
	}
	return nil
}

func (x *FindUsersRequest) GetFirstname() string {
	if x != nil {
		return x.Firstname
	}
	return ""
}

func (x *FindUsersRequest) GetLastname() string {
	if x != nil {
		return x.Lastname
	}
	return ""
}

func (x *FindUsersRequest) GetWanted() bool {
	if x != nil {
		return x.Wanted
	}
	return false
}

type FindUsersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalCount int64          `protobuf:"varint,1,opt,name=totalCount,proto3" json:"totalCount,omitempty"`
	Offset     int64          `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	End        int64          `protobuf:"varint,3,opt,name=end,proto3" json:"end,omitempty"`
	Users      []*common.User `protobuf:"bytes,4,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *FindUsersResponse) Reset() {
	*x = FindUsersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_users_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindUsersResponse) ProtoMessage() {}

func (x *FindUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindUsersResponse.ProtoReflect.Descriptor instead.
func (*FindUsersResponse) Descriptor() ([]byte, []int) {
	return file_users_users_proto_rawDescGZIP(), []int{1}
}

func (x *FindUsersResponse) GetTotalCount() int64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *FindUsersResponse) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *FindUsersResponse) GetEnd() int64 {
	if x != nil {
		return x.End
	}
	return 0
}

func (x *FindUsersResponse) GetUsers() []*common.User {
	if x != nil {
		return x.Users
	}
	return nil
}

type GetUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID int32 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_users_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[2]
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
	return file_users_users_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserRequest) GetUserID() int32 {
	if x != nil {
		return x.UserID
	}
	return 0
}

type GetUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *common.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *GetUserResponse) Reset() {
	*x = GetUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_users_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserResponse) ProtoMessage() {}

func (x *GetUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[3]
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
	return file_users_users_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserResponse) GetUser() *common.User {
	if x != nil {
		return x.User
	}
	return nil
}

type GetUserActivityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID int32 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (x *GetUserActivityRequest) Reset() {
	*x = GetUserActivityRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_users_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserActivityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserActivityRequest) ProtoMessage() {}

func (x *GetUserActivityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserActivityRequest.ProtoReflect.Descriptor instead.
func (*GetUserActivityRequest) Descriptor() ([]byte, []int) {
	return file_users_users_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserActivityRequest) GetUserID() int32 {
	if x != nil {
		return x.UserID
	}
	return 0
}

type GetUserActivityResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Activity []*UserActivity `protobuf:"bytes,1,rep,name=activity,proto3" json:"activity,omitempty"`
}

func (x *GetUserActivityResponse) Reset() {
	*x = GetUserActivityResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_users_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserActivityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserActivityResponse) ProtoMessage() {}

func (x *GetUserActivityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserActivityResponse.ProtoReflect.Descriptor instead.
func (*GetUserActivityResponse) Descriptor() ([]byte, []int) {
	return file_users_users_proto_rawDescGZIP(), []int{5}
}

func (x *GetUserActivityResponse) GetActivity() []*UserActivity {
	if x != nil {
		return x.Activity
	}
	return nil
}

type SetUserPropsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID int32 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
	Wanted *bool `protobuf:"varint,2,opt,name=wanted,proto3,oneof" json:"wanted,omitempty"`
}

func (x *SetUserPropsRequest) Reset() {
	*x = SetUserPropsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_users_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetUserPropsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUserPropsRequest) ProtoMessage() {}

func (x *SetUserPropsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[6]
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
	return file_users_users_proto_rawDescGZIP(), []int{6}
}

func (x *SetUserPropsRequest) GetUserID() int32 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *SetUserPropsRequest) GetWanted() bool {
	if x != nil && x.Wanted != nil {
		return *x.Wanted
	}
	return false
}

type SetUserPropsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SetUserPropsResponse) Reset() {
	*x = SetUserPropsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_users_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetUserPropsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetUserPropsResponse) ProtoMessage() {}

func (x *SetUserPropsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[7]
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
	return file_users_users_proto_rawDescGZIP(), []int{7}
}

type UserActivity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"arpanet_user_activity.id"`                // @gotags: sql:"primary_key" alias:"arpanet_user_activity.id"
	Type       string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty" alias:"arpanet_user_activity.type"`             // @gotags: alias:"arpanet_user_activity.type"
	CreatedAt  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=createdAt,proto3" json:"createdAt,omitempty" alias:"arpanet_user_activity.created_at"`   // @gotags: alias:"arpanet_user_activity.created_at"
	TargetUser *common.ShortUser      `protobuf:"bytes,4,opt,name=targetUser,proto3" json:"targetUser,omitempty" alias:"target_user"` // @gotags: alias:"target_user"
	CauseUser  *common.ShortUser      `protobuf:"bytes,5,opt,name=causeUser,proto3" json:"causeUser,omitempty" alias:"cause_user"`   // @gotags: alias:"cause_user"
	Key        string                 `protobuf:"bytes,6,opt,name=key,proto3" json:"key,omitempty" alias:"arpanet_user_activity.key"`               // @gotags: alias:"arpanet_user_activity.key"
	OldValue   string                 `protobuf:"bytes,7,opt,name=oldValue,proto3" json:"oldValue,omitempty" alias:"arpanet_user_activity.old_value"`     // @gotags: alias:"arpanet_user_activity.old_value"
	NewValue   string                 `protobuf:"bytes,8,opt,name=newValue,proto3" json:"newValue,omitempty" alias:"arpanet_user_activity.new_value"`     // @gotags: alias:"arpanet_user_activity.new_value"
	Reason     string                 `protobuf:"bytes,9,opt,name=reason,proto3" json:"reason,omitempty" alias:"arpanet_user_activity.reason"`         // @gotags: alias:"arpanet_user_activity.reason"
}

func (x *UserActivity) Reset() {
	*x = UserActivity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_users_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserActivity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserActivity) ProtoMessage() {}

func (x *UserActivity) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserActivity.ProtoReflect.Descriptor instead.
func (*UserActivity) Descriptor() ([]byte, []int) {
	return file_users_users_proto_rawDescGZIP(), []int{8}
}

func (x *UserActivity) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserActivity) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *UserActivity) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *UserActivity) GetTargetUser() *common.ShortUser {
	if x != nil {
		return x.TargetUser
	}
	return nil
}

func (x *UserActivity) GetCauseUser() *common.ShortUser {
	if x != nil {
		return x.CauseUser
	}
	return nil
}

func (x *UserActivity) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *UserActivity) GetOldValue() string {
	if x != nil {
		return x.OldValue
	}
	return ""
}

func (x *UserActivity) GetNewValue() string {
	if x != nil {
		return x.NewValue
	}
	return ""
}

func (x *UserActivity) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

var File_users_users_proto protoreflect.FileDescriptor

var file_users_users_proto_rawDesc = []byte{
	0x0a, 0x11, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x67, 0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x17,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xab, 0x01, 0x0a, 0x10, 0x46, 0x69, 0x6e, 0x64, 0x55,
	0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f,
	0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6f, 0x66, 0x66,
	0x73, 0x65, 0x74, 0x12, 0x2d, 0x0a, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x42, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x77, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x77, 0x61,
	0x6e, 0x74, 0x65, 0x64, 0x22, 0x85, 0x01, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x64, 0x55, 0x73, 0x65,
	0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x03, 0x65, 0x6e, 0x64, 0x12, 0x26, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x31, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f,
	0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x22,
	0x37, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x39, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x22, 0x4e, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63,
	0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33,
	0x0a, 0x08, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x08, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x69, 0x74, 0x79, 0x22, 0x5e, 0x0a, 0x13, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72,
	0x6f, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a,
	0x02, 0x20, 0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1b, 0x0a, 0x06, 0x77,
	0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x06, 0x77,
	0x61, 0x6e, 0x74, 0x65, 0x64, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x77, 0x61, 0x6e,
	0x74, 0x65, 0x64, 0x22, 0x16, 0x0a, 0x14, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72,
	0x6f, 0x70, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xc3, 0x02, 0x0a, 0x0c,
	0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x12, 0x17, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x32, 0x02, 0x20,
	0x00, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x35, 0x0a, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x0a,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x09, 0x63, 0x61,
	0x75, 0x73, 0x65, 0x55, 0x73, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x67, 0x65, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x09, 0x63, 0x61, 0x75, 0x73, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x6c, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x6c, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x6e, 0x65, 0x77, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6e, 0x65, 0x77, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x32, 0xc3, 0x02, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x46, 0x0a, 0x09, 0x46, 0x69, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12,
	0x1b, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x46, 0x69, 0x6e, 0x64,
	0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x67,
	0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x55, 0x73, 0x65,
	0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x07, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1a, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x58, 0x0a, 0x0f,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x12,
	0x21, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x22, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0c, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x1e, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x67, 0x65, 0x6e, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x61, 0x72,
	0x70, 0x61, 0x6e, 0x65, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_users_users_proto_rawDescOnce sync.Once
	file_users_users_proto_rawDescData = file_users_users_proto_rawDesc
)

func file_users_users_proto_rawDescGZIP() []byte {
	file_users_users_proto_rawDescOnce.Do(func() {
		file_users_users_proto_rawDescData = protoimpl.X.CompressGZIP(file_users_users_proto_rawDescData)
	})
	return file_users_users_proto_rawDescData
}

var file_users_users_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_users_users_proto_goTypes = []interface{}{
	(*FindUsersRequest)(nil),        // 0: gen.users.FindUsersRequest
	(*FindUsersResponse)(nil),       // 1: gen.users.FindUsersResponse
	(*GetUserRequest)(nil),          // 2: gen.users.GetUserRequest
	(*GetUserResponse)(nil),         // 3: gen.users.GetUserResponse
	(*GetUserActivityRequest)(nil),  // 4: gen.users.GetUserActivityRequest
	(*GetUserActivityResponse)(nil), // 5: gen.users.GetUserActivityResponse
	(*SetUserPropsRequest)(nil),     // 6: gen.users.SetUserPropsRequest
	(*SetUserPropsResponse)(nil),    // 7: gen.users.SetUserPropsResponse
	(*UserActivity)(nil),            // 8: gen.users.UserActivity
	(*common.OrderBy)(nil),          // 9: gen.common.OrderBy
	(*common.User)(nil),             // 10: gen.common.User
	(*timestamppb.Timestamp)(nil),   // 11: google.protobuf.Timestamp
	(*common.ShortUser)(nil),        // 12: gen.common.ShortUser
}
var file_users_users_proto_depIdxs = []int32{
	9,  // 0: gen.users.FindUsersRequest.orderBy:type_name -> gen.common.OrderBy
	10, // 1: gen.users.FindUsersResponse.users:type_name -> gen.common.User
	10, // 2: gen.users.GetUserResponse.user:type_name -> gen.common.User
	8,  // 3: gen.users.GetUserActivityResponse.activity:type_name -> gen.users.UserActivity
	11, // 4: gen.users.UserActivity.createdAt:type_name -> google.protobuf.Timestamp
	12, // 5: gen.users.UserActivity.targetUser:type_name -> gen.common.ShortUser
	12, // 6: gen.users.UserActivity.causeUser:type_name -> gen.common.ShortUser
	0,  // 7: gen.users.UsersService.FindUsers:input_type -> gen.users.FindUsersRequest
	2,  // 8: gen.users.UsersService.GetUser:input_type -> gen.users.GetUserRequest
	4,  // 9: gen.users.UsersService.GetUserActivity:input_type -> gen.users.GetUserActivityRequest
	6,  // 10: gen.users.UsersService.SetUserProps:input_type -> gen.users.SetUserPropsRequest
	1,  // 11: gen.users.UsersService.FindUsers:output_type -> gen.users.FindUsersResponse
	3,  // 12: gen.users.UsersService.GetUser:output_type -> gen.users.GetUserResponse
	5,  // 13: gen.users.UsersService.GetUserActivity:output_type -> gen.users.GetUserActivityResponse
	7,  // 14: gen.users.UsersService.SetUserProps:output_type -> gen.users.SetUserPropsResponse
	11, // [11:15] is the sub-list for method output_type
	7,  // [7:11] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_users_users_proto_init() }
func file_users_users_proto_init() {
	if File_users_users_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_users_users_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindUsersRequest); i {
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
		file_users_users_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindUsersResponse); i {
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
		file_users_users_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_users_users_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_users_users_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserActivityRequest); i {
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
		file_users_users_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserActivityResponse); i {
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
		file_users_users_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
		file_users_users_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
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
		file_users_users_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserActivity); i {
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
	file_users_users_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_users_users_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_users_users_proto_goTypes,
		DependencyIndexes: file_users_users_proto_depIdxs,
		MessageInfos:      file_users_users_proto_msgTypes,
	}.Build()
	File_users_users_proto = out.File
	file_users_users_proto_rawDesc = nil
	file_users_users_proto_goTypes = nil
	file_users_users_proto_depIdxs = nil
}
