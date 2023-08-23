// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: services/notificator/notificator.proto

package notificator

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	notifications "github.com/galexrt/fivenet/gen/go/proto/resources/notifications"
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

type GetNotificationsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination  *database.PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	IncludeRead *bool                       `protobuf:"varint,2,opt,name=include_read,json=includeRead,proto3,oneof" json:"include_read,omitempty"`
}

func (x *GetNotificationsRequest) Reset() {
	*x = GetNotificationsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_notificator_notificator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotificationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotificationsRequest) ProtoMessage() {}

func (x *GetNotificationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotificationsRequest.ProtoReflect.Descriptor instead.
func (*GetNotificationsRequest) Descriptor() ([]byte, []int) {
	return file_services_notificator_notificator_proto_rawDescGZIP(), []int{0}
}

func (x *GetNotificationsRequest) GetPagination() *database.PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *GetNotificationsRequest) GetIncludeRead() bool {
	if x != nil && x.IncludeRead != nil {
		return *x.IncludeRead
	}
	return false
}

type GetNotificationsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination    *database.PaginationResponse  `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Notifications []*notifications.Notification `protobuf:"bytes,2,rep,name=notifications,proto3" json:"notifications,omitempty"`
}

func (x *GetNotificationsResponse) Reset() {
	*x = GetNotificationsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_notificator_notificator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotificationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotificationsResponse) ProtoMessage() {}

func (x *GetNotificationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotificationsResponse.ProtoReflect.Descriptor instead.
func (*GetNotificationsResponse) Descriptor() ([]byte, []int) {
	return file_services_notificator_notificator_proto_rawDescGZIP(), []int{1}
}

func (x *GetNotificationsResponse) GetPagination() *database.PaginationResponse {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *GetNotificationsResponse) GetNotifications() []*notifications.Notification {
	if x != nil {
		return x.Notifications
	}
	return nil
}

type ReadNotificationsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []uint64 `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"`
	All *bool    `protobuf:"varint,2,opt,name=all,proto3,oneof" json:"all,omitempty"`
}

func (x *ReadNotificationsRequest) Reset() {
	*x = ReadNotificationsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_notificator_notificator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadNotificationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadNotificationsRequest) ProtoMessage() {}

func (x *ReadNotificationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadNotificationsRequest.ProtoReflect.Descriptor instead.
func (*ReadNotificationsRequest) Descriptor() ([]byte, []int) {
	return file_services_notificator_notificator_proto_rawDescGZIP(), []int{2}
}

func (x *ReadNotificationsRequest) GetIds() []uint64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ReadNotificationsRequest) GetAll() bool {
	if x != nil && x.All != nil {
		return *x.All
	}
	return false
}

type ReadNotificationsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Updated uint64 `protobuf:"varint,1,opt,name=updated,proto3" json:"updated,omitempty"`
}

func (x *ReadNotificationsResponse) Reset() {
	*x = ReadNotificationsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_notificator_notificator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadNotificationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadNotificationsResponse) ProtoMessage() {}

func (x *ReadNotificationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadNotificationsResponse.ProtoReflect.Descriptor instead.
func (*ReadNotificationsResponse) Descriptor() ([]byte, []int) {
	return file_services_notificator_notificator_proto_rawDescGZIP(), []int{3}
}

func (x *ReadNotificationsResponse) GetUpdated() uint64 {
	if x != nil {
		return x.Updated
	}
	return 0
}

type StreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastId uint64 `protobuf:"varint,1,opt,name=last_id,json=lastId,proto3" json:"last_id,omitempty"`
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_notificator_notificator_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRequest.ProtoReflect.Descriptor instead.
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return file_services_notificator_notificator_proto_rawDescGZIP(), []int{4}
}

func (x *StreamRequest) GetLastId() uint64 {
	if x != nil {
		return x.LastId
	}
	return 0
}

type StreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastId        uint64                        `protobuf:"varint,1,opt,name=last_id,json=lastId,proto3" json:"last_id,omitempty"`
	RestartStream bool                          `protobuf:"varint,2,opt,name=restart_stream,json=restartStream,proto3" json:"restart_stream,omitempty"`
	Notifications []*notifications.Notification `protobuf:"bytes,3,rep,name=notifications,proto3" json:"notifications,omitempty"`
	Token         *TokenUpdate                  `protobuf:"bytes,4,opt,name=token,proto3,oneof" json:"token,omitempty"`
}

func (x *StreamResponse) Reset() {
	*x = StreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_notificator_notificator_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResponse) ProtoMessage() {}

func (x *StreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamResponse.ProtoReflect.Descriptor instead.
func (*StreamResponse) Descriptor() ([]byte, []int) {
	return file_services_notificator_notificator_proto_rawDescGZIP(), []int{5}
}

func (x *StreamResponse) GetLastId() uint64 {
	if x != nil {
		return x.LastId
	}
	return 0
}

func (x *StreamResponse) GetRestartStream() bool {
	if x != nil {
		return x.RestartStream
	}
	return false
}

func (x *StreamResponse) GetNotifications() []*notifications.Notification {
	if x != nil {
		return x.Notifications
	}
	return nil
}

func (x *StreamResponse) GetToken() *TokenUpdate {
	if x != nil {
		return x.Token
	}
	return nil
}

type TokenUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewToken    *string              `protobuf:"bytes,1,opt,name=new_token,json=newToken,proto3,oneof" json:"new_token,omitempty"`
	Expires     *timestamp.Timestamp `protobuf:"bytes,2,opt,name=expires,proto3" json:"expires,omitempty"`
	Permissions []string             `protobuf:"bytes,3,rep,name=permissions,proto3" json:"permissions,omitempty"`
	UserInfo    *users.User          `protobuf:"bytes,4,opt,name=user_info,json=userInfo,proto3,oneof" json:"user_info,omitempty"`
	JobProps    *users.JobProps      `protobuf:"bytes,5,opt,name=job_props,json=jobProps,proto3,oneof" json:"job_props,omitempty"`
}

func (x *TokenUpdate) Reset() {
	*x = TokenUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_notificator_notificator_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenUpdate) ProtoMessage() {}

func (x *TokenUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenUpdate.ProtoReflect.Descriptor instead.
func (*TokenUpdate) Descriptor() ([]byte, []int) {
	return file_services_notificator_notificator_proto_rawDescGZIP(), []int{6}
}

func (x *TokenUpdate) GetNewToken() string {
	if x != nil && x.NewToken != nil {
		return *x.NewToken
	}
	return ""
}

func (x *TokenUpdate) GetExpires() *timestamp.Timestamp {
	if x != nil {
		return x.Expires
	}
	return nil
}

func (x *TokenUpdate) GetPermissions() []string {
	if x != nil {
		return x.Permissions
	}
	return nil
}

func (x *TokenUpdate) GetUserInfo() *users.User {
	if x != nil {
		return x.UserInfo
	}
	return nil
}

func (x *TokenUpdate) GetJobProps() *users.JobProps {
	if x != nil {
		return x.JobProps
	}
	return nil
}

var File_services_notificator_notificator_proto protoreflect.FileDescriptor

var file_services_notificator_notificator_proto_rawDesc = []byte{
	0x0a, 0x26, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x1a, 0x28,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaa, 0x01, 0x0a,
	0x17, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x56, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a,
	0x01, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x26, 0x0a, 0x0c, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x0b, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64,
	0x65, 0x52, 0x65, 0x61, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x69, 0x6e, 0x63,
	0x6c, 0x75, 0x64, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x22, 0xb6, 0x01, 0x0a, 0x18, 0x47, 0x65,
	0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61,
	0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4b, 0x0a, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x22, 0x59, 0x0a, 0x18, 0x52, 0x65, 0x61, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e,
	0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x42, 0x0c, 0xfa, 0x42, 0x09,
	0x92, 0x01, 0x06, 0x08, 0x01, 0x10, 0x14, 0x28, 0x01, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x15,
	0x0a, 0x03, 0x61, 0x6c, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x03, 0x61,
	0x6c, 0x6c, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x61, 0x6c, 0x6c, 0x22, 0x35, 0x0a,
	0x19, 0x52, 0x65, 0x61, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x22, 0x28, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x64, 0x22, 0xe5,
	0x01, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x17, 0x0a, 0x07, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x12, 0x4b, 0x0a, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3c,
	0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x6f, 0x72, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x48, 0x00, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06,
	0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xab, 0x02, 0x0a, 0x0b, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x09, 0x6e, 0x65, 0x77, 0x5f, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x6e, 0x65, 0x77,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x38, 0x0a, 0x07, 0x65, 0x78, 0x70, 0x69,
	0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x65, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x37, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x48, 0x01,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x88, 0x01, 0x01, 0x12, 0x3b, 0x0a,
	0x09, 0x6a, 0x6f, 0x62, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x48, 0x02, 0x52, 0x08, 0x6a,
	0x6f, 0x62, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x88, 0x01, 0x01, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6e,
	0x65, 0x77, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x70,
	0x72, 0x6f, 0x70, 0x73, 0x32, 0xd4, 0x02, 0x0a, 0x12, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x71, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x2d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x74,
	0x0a, 0x11, 0x52, 0x65, 0x61, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x2e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x23,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x4a, 0x5a, 0x48, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72,
	0x74, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x3b, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_notificator_notificator_proto_rawDescOnce sync.Once
	file_services_notificator_notificator_proto_rawDescData = file_services_notificator_notificator_proto_rawDesc
)

func file_services_notificator_notificator_proto_rawDescGZIP() []byte {
	file_services_notificator_notificator_proto_rawDescOnce.Do(func() {
		file_services_notificator_notificator_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_notificator_notificator_proto_rawDescData)
	})
	return file_services_notificator_notificator_proto_rawDescData
}

var file_services_notificator_notificator_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_services_notificator_notificator_proto_goTypes = []interface{}{
	(*GetNotificationsRequest)(nil),     // 0: services.notificator.GetNotificationsRequest
	(*GetNotificationsResponse)(nil),    // 1: services.notificator.GetNotificationsResponse
	(*ReadNotificationsRequest)(nil),    // 2: services.notificator.ReadNotificationsRequest
	(*ReadNotificationsResponse)(nil),   // 3: services.notificator.ReadNotificationsResponse
	(*StreamRequest)(nil),               // 4: services.notificator.StreamRequest
	(*StreamResponse)(nil),              // 5: services.notificator.StreamResponse
	(*TokenUpdate)(nil),                 // 6: services.notificator.TokenUpdate
	(*database.PaginationRequest)(nil),  // 7: resources.common.database.PaginationRequest
	(*database.PaginationResponse)(nil), // 8: resources.common.database.PaginationResponse
	(*notifications.Notification)(nil),  // 9: resources.notifications.Notification
	(*timestamp.Timestamp)(nil),         // 10: resources.timestamp.Timestamp
	(*users.User)(nil),                  // 11: resources.users.User
	(*users.JobProps)(nil),              // 12: resources.users.JobProps
}
var file_services_notificator_notificator_proto_depIdxs = []int32{
	7,  // 0: services.notificator.GetNotificationsRequest.pagination:type_name -> resources.common.database.PaginationRequest
	8,  // 1: services.notificator.GetNotificationsResponse.pagination:type_name -> resources.common.database.PaginationResponse
	9,  // 2: services.notificator.GetNotificationsResponse.notifications:type_name -> resources.notifications.Notification
	9,  // 3: services.notificator.StreamResponse.notifications:type_name -> resources.notifications.Notification
	6,  // 4: services.notificator.StreamResponse.token:type_name -> services.notificator.TokenUpdate
	10, // 5: services.notificator.TokenUpdate.expires:type_name -> resources.timestamp.Timestamp
	11, // 6: services.notificator.TokenUpdate.user_info:type_name -> resources.users.User
	12, // 7: services.notificator.TokenUpdate.job_props:type_name -> resources.users.JobProps
	0,  // 8: services.notificator.NotificatorService.GetNotifications:input_type -> services.notificator.GetNotificationsRequest
	2,  // 9: services.notificator.NotificatorService.ReadNotifications:input_type -> services.notificator.ReadNotificationsRequest
	4,  // 10: services.notificator.NotificatorService.Stream:input_type -> services.notificator.StreamRequest
	1,  // 11: services.notificator.NotificatorService.GetNotifications:output_type -> services.notificator.GetNotificationsResponse
	3,  // 12: services.notificator.NotificatorService.ReadNotifications:output_type -> services.notificator.ReadNotificationsResponse
	5,  // 13: services.notificator.NotificatorService.Stream:output_type -> services.notificator.StreamResponse
	11, // [11:14] is the sub-list for method output_type
	8,  // [8:11] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_services_notificator_notificator_proto_init() }
func file_services_notificator_notificator_proto_init() {
	if File_services_notificator_notificator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_notificator_notificator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotificationsRequest); i {
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
		file_services_notificator_notificator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotificationsResponse); i {
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
		file_services_notificator_notificator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadNotificationsRequest); i {
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
		file_services_notificator_notificator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadNotificationsResponse); i {
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
		file_services_notificator_notificator_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamRequest); i {
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
		file_services_notificator_notificator_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamResponse); i {
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
		file_services_notificator_notificator_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenUpdate); i {
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
	file_services_notificator_notificator_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_services_notificator_notificator_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_services_notificator_notificator_proto_msgTypes[5].OneofWrappers = []interface{}{}
	file_services_notificator_notificator_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_notificator_notificator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_notificator_notificator_proto_goTypes,
		DependencyIndexes: file_services_notificator_notificator_proto_depIdxs,
		MessageInfos:      file_services_notificator_notificator_proto_msgTypes,
	}.Build()
	File_services_notificator_notificator_proto = out.File
	file_services_notificator_notificator_proto_rawDesc = nil
	file_services_notificator_notificator_proto_goTypes = nil
	file_services_notificator_notificator_proto_depIdxs = nil
}
