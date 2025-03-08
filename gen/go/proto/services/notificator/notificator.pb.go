// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: services/notificator/notificator.proto

package notificator

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	mailer "github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	notifications "github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
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

type GetNotificationsRequest struct {
	state         protoimpl.MessageState               `protogen:"open.v1"`
	Pagination    *database.PaginationRequest          `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	IncludeRead   *bool                                `protobuf:"varint,2,opt,name=include_read,json=includeRead,proto3,oneof" json:"include_read,omitempty"`
	Categories    []notifications.NotificationCategory `protobuf:"varint,3,rep,packed,name=categories,proto3,enum=resources.notifications.NotificationCategory" json:"categories,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetNotificationsRequest) Reset() {
	*x = GetNotificationsRequest{}
	mi := &file_services_notificator_notificator_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetNotificationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotificationsRequest) ProtoMessage() {}

func (x *GetNotificationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[0]
	if x != nil {
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

func (x *GetNotificationsRequest) GetCategories() []notifications.NotificationCategory {
	if x != nil {
		return x.Categories
	}
	return nil
}

type GetNotificationsResponse struct {
	state         protoimpl.MessageState        `protogen:"open.v1"`
	Pagination    *database.PaginationResponse  `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Notifications []*notifications.Notification `protobuf:"bytes,2,rep,name=notifications,proto3" json:"notifications,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetNotificationsResponse) Reset() {
	*x = GetNotificationsResponse{}
	mi := &file_services_notificator_notificator_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetNotificationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotificationsResponse) ProtoMessage() {}

func (x *GetNotificationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[1]
	if x != nil {
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

type MarkNotificationsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Unread        bool                   `protobuf:"varint,1,opt,name=unread,proto3" json:"unread,omitempty"`
	Ids           []uint64               `protobuf:"varint,2,rep,packed,name=ids,proto3" json:"ids,omitempty"`
	All           *bool                  `protobuf:"varint,3,opt,name=all,proto3,oneof" json:"all,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MarkNotificationsRequest) Reset() {
	*x = MarkNotificationsRequest{}
	mi := &file_services_notificator_notificator_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MarkNotificationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkNotificationsRequest) ProtoMessage() {}

func (x *MarkNotificationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkNotificationsRequest.ProtoReflect.Descriptor instead.
func (*MarkNotificationsRequest) Descriptor() ([]byte, []int) {
	return file_services_notificator_notificator_proto_rawDescGZIP(), []int{2}
}

func (x *MarkNotificationsRequest) GetUnread() bool {
	if x != nil {
		return x.Unread
	}
	return false
}

func (x *MarkNotificationsRequest) GetIds() []uint64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *MarkNotificationsRequest) GetAll() bool {
	if x != nil && x.All != nil {
		return *x.All
	}
	return false
}

type MarkNotificationsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Updated       uint64                 `protobuf:"varint,1,opt,name=updated,proto3" json:"updated,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MarkNotificationsResponse) Reset() {
	*x = MarkNotificationsResponse{}
	mi := &file_services_notificator_notificator_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MarkNotificationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkNotificationsResponse) ProtoMessage() {}

func (x *MarkNotificationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkNotificationsResponse.ProtoReflect.Descriptor instead.
func (*MarkNotificationsResponse) Descriptor() ([]byte, []int) {
	return file_services_notificator_notificator_proto_rawDescGZIP(), []int{3}
}

func (x *MarkNotificationsResponse) GetUpdated() uint64 {
	if x != nil {
		return x.Updated
	}
	return 0
}

type StreamRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	mi := &file_services_notificator_notificator_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[4]
	if x != nil {
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

type StreamResponse struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	NotificationCount int32                  `protobuf:"varint,1,opt,name=notification_count,json=notificationCount,proto3" json:"notification_count,omitempty"`
	Restart           *bool                  `protobuf:"varint,2,opt,name=restart,proto3,oneof" json:"restart,omitempty"`
	// Types that are valid to be assigned to Data:
	//
	//	*StreamResponse_UserEvent
	//	*StreamResponse_JobEvent
	//	*StreamResponse_JobGradeEvent
	//	*StreamResponse_SystemEvent
	//	*StreamResponse_MailerEvent
	Data          isStreamResponse_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamResponse) Reset() {
	*x = StreamResponse{}
	mi := &file_services_notificator_notificator_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResponse) ProtoMessage() {}

func (x *StreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_notificator_notificator_proto_msgTypes[5]
	if x != nil {
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

func (x *StreamResponse) GetNotificationCount() int32 {
	if x != nil {
		return x.NotificationCount
	}
	return 0
}

func (x *StreamResponse) GetRestart() bool {
	if x != nil && x.Restart != nil {
		return *x.Restart
	}
	return false
}

func (x *StreamResponse) GetData() isStreamResponse_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *StreamResponse) GetUserEvent() *notifications.UserEvent {
	if x != nil {
		if x, ok := x.Data.(*StreamResponse_UserEvent); ok {
			return x.UserEvent
		}
	}
	return nil
}

func (x *StreamResponse) GetJobEvent() *notifications.JobEvent {
	if x != nil {
		if x, ok := x.Data.(*StreamResponse_JobEvent); ok {
			return x.JobEvent
		}
	}
	return nil
}

func (x *StreamResponse) GetJobGradeEvent() *notifications.JobGradeEvent {
	if x != nil {
		if x, ok := x.Data.(*StreamResponse_JobGradeEvent); ok {
			return x.JobGradeEvent
		}
	}
	return nil
}

func (x *StreamResponse) GetSystemEvent() *notifications.SystemEvent {
	if x != nil {
		if x, ok := x.Data.(*StreamResponse_SystemEvent); ok {
			return x.SystemEvent
		}
	}
	return nil
}

func (x *StreamResponse) GetMailerEvent() *mailer.MailerEvent {
	if x != nil {
		if x, ok := x.Data.(*StreamResponse_MailerEvent); ok {
			return x.MailerEvent
		}
	}
	return nil
}

type isStreamResponse_Data interface {
	isStreamResponse_Data()
}

type StreamResponse_UserEvent struct {
	UserEvent *notifications.UserEvent `protobuf:"bytes,3,opt,name=user_event,json=userEvent,proto3,oneof"`
}

type StreamResponse_JobEvent struct {
	JobEvent *notifications.JobEvent `protobuf:"bytes,4,opt,name=job_event,json=jobEvent,proto3,oneof"`
}

type StreamResponse_JobGradeEvent struct {
	JobGradeEvent *notifications.JobGradeEvent `protobuf:"bytes,7,opt,name=job_grade_event,json=jobGradeEvent,proto3,oneof"`
}

type StreamResponse_SystemEvent struct {
	SystemEvent *notifications.SystemEvent `protobuf:"bytes,5,opt,name=system_event,json=systemEvent,proto3,oneof"`
}

type StreamResponse_MailerEvent struct {
	MailerEvent *mailer.MailerEvent `protobuf:"bytes,6,opt,name=mailer_event,json=mailerEvent,proto3,oneof"`
}

func (*StreamResponse_UserEvent) isStreamResponse_Data() {}

func (*StreamResponse_JobEvent) isStreamResponse_Data() {}

func (*StreamResponse_JobGradeEvent) isStreamResponse_Data() {}

func (*StreamResponse_SystemEvent) isStreamResponse_Data() {}

func (*StreamResponse_MailerEvent) isStreamResponse_Data() {}

var File_services_notificator_notificator_proto protoreflect.FileDescriptor

var file_services_notificator_notificator_proto_rawDesc = string([]byte{
	0x0a, 0x26, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x1a, 0x28,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2b, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x8a, 0x02, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x56, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e,
	0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x0c, 0x69, 0x6e, 0x63, 0x6c, 0x75,
	0x64, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52,
	0x0b, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x52, 0x65, 0x61, 0x64, 0x88, 0x01, 0x01, 0x12,
	0x5e, 0x0a, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0e, 0x32, 0x2d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x42, 0x0f, 0xfa, 0x42, 0x0c, 0x92, 0x01, 0x09, 0x10, 0x04, 0x22, 0x05, 0x82, 0x01,
	0x02, 0x10, 0x01, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x42,
	0x0f, 0x0a, 0x0d, 0x5f, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x64,
	0x22, 0xb6, 0x01, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a,
	0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4b, 0x0a, 0x0d,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x71, 0x0a, 0x18, 0x4d, 0x61, 0x72,
	0x6b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x75, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x12, 0x1e, 0x0a,
	0x03, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x04, 0x42, 0x0c, 0xfa, 0x42, 0x09, 0x92,
	0x01, 0x06, 0x08, 0x01, 0x10, 0x14, 0x28, 0x01, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x15, 0x0a,
	0x03, 0x61, 0x6c, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x03, 0x61, 0x6c,
	0x6c, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x61, 0x6c, 0x6c, 0x22, 0x35, 0x0a, 0x19,
	0x4d, 0x61, 0x72, 0x6b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0xdf, 0x03, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x12, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x11, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x07, 0x72, 0x65, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x88, 0x01, 0x01, 0x12, 0x43, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52,
	0x09, 0x75, 0x73, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x40, 0x0a, 0x09, 0x6a, 0x6f,
	0x62, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x48, 0x00, 0x52, 0x08, 0x6a, 0x6f, 0x62, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x50, 0x0a, 0x0f,
	0x6a, 0x6f, 0x62, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x4a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52,
	0x0d, 0x6a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x49,
	0x0a, 0x0c, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x42, 0x0a, 0x0c, 0x6d, 0x61, 0x69,
	0x6c, 0x65, 0x72, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6d, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00,
	0x52, 0x0b, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x0b, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x03, 0xf8, 0x42, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x72,
	0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x32, 0xd4, 0x02, 0x0a, 0x12, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x71, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x2d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x74, 0x0a, 0x11, 0x4d, 0x61, 0x72, 0x6b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4d, 0x61, 0x72,
	0x6b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4d, 0x61, 0x72,
	0x6b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x12, 0x23, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x4e, 0x5a,
	0x4c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65,
	0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f,
	0x72, 0x3b, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_services_notificator_notificator_proto_rawDescOnce sync.Once
	file_services_notificator_notificator_proto_rawDescData []byte
)

func file_services_notificator_notificator_proto_rawDescGZIP() []byte {
	file_services_notificator_notificator_proto_rawDescOnce.Do(func() {
		file_services_notificator_notificator_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_services_notificator_notificator_proto_rawDesc), len(file_services_notificator_notificator_proto_rawDesc)))
	})
	return file_services_notificator_notificator_proto_rawDescData
}

var file_services_notificator_notificator_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_services_notificator_notificator_proto_goTypes = []any{
	(*GetNotificationsRequest)(nil),         // 0: services.notificator.GetNotificationsRequest
	(*GetNotificationsResponse)(nil),        // 1: services.notificator.GetNotificationsResponse
	(*MarkNotificationsRequest)(nil),        // 2: services.notificator.MarkNotificationsRequest
	(*MarkNotificationsResponse)(nil),       // 3: services.notificator.MarkNotificationsResponse
	(*StreamRequest)(nil),                   // 4: services.notificator.StreamRequest
	(*StreamResponse)(nil),                  // 5: services.notificator.StreamResponse
	(*database.PaginationRequest)(nil),      // 6: resources.common.database.PaginationRequest
	(notifications.NotificationCategory)(0), // 7: resources.notifications.NotificationCategory
	(*database.PaginationResponse)(nil),     // 8: resources.common.database.PaginationResponse
	(*notifications.Notification)(nil),      // 9: resources.notifications.Notification
	(*notifications.UserEvent)(nil),         // 10: resources.notifications.UserEvent
	(*notifications.JobEvent)(nil),          // 11: resources.notifications.JobEvent
	(*notifications.JobGradeEvent)(nil),     // 12: resources.notifications.JobGradeEvent
	(*notifications.SystemEvent)(nil),       // 13: resources.notifications.SystemEvent
	(*mailer.MailerEvent)(nil),              // 14: resources.mailer.MailerEvent
}
var file_services_notificator_notificator_proto_depIdxs = []int32{
	6,  // 0: services.notificator.GetNotificationsRequest.pagination:type_name -> resources.common.database.PaginationRequest
	7,  // 1: services.notificator.GetNotificationsRequest.categories:type_name -> resources.notifications.NotificationCategory
	8,  // 2: services.notificator.GetNotificationsResponse.pagination:type_name -> resources.common.database.PaginationResponse
	9,  // 3: services.notificator.GetNotificationsResponse.notifications:type_name -> resources.notifications.Notification
	10, // 4: services.notificator.StreamResponse.user_event:type_name -> resources.notifications.UserEvent
	11, // 5: services.notificator.StreamResponse.job_event:type_name -> resources.notifications.JobEvent
	12, // 6: services.notificator.StreamResponse.job_grade_event:type_name -> resources.notifications.JobGradeEvent
	13, // 7: services.notificator.StreamResponse.system_event:type_name -> resources.notifications.SystemEvent
	14, // 8: services.notificator.StreamResponse.mailer_event:type_name -> resources.mailer.MailerEvent
	0,  // 9: services.notificator.NotificatorService.GetNotifications:input_type -> services.notificator.GetNotificationsRequest
	2,  // 10: services.notificator.NotificatorService.MarkNotifications:input_type -> services.notificator.MarkNotificationsRequest
	4,  // 11: services.notificator.NotificatorService.Stream:input_type -> services.notificator.StreamRequest
	1,  // 12: services.notificator.NotificatorService.GetNotifications:output_type -> services.notificator.GetNotificationsResponse
	3,  // 13: services.notificator.NotificatorService.MarkNotifications:output_type -> services.notificator.MarkNotificationsResponse
	5,  // 14: services.notificator.NotificatorService.Stream:output_type -> services.notificator.StreamResponse
	12, // [12:15] is the sub-list for method output_type
	9,  // [9:12] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_services_notificator_notificator_proto_init() }
func file_services_notificator_notificator_proto_init() {
	if File_services_notificator_notificator_proto != nil {
		return
	}
	file_services_notificator_notificator_proto_msgTypes[0].OneofWrappers = []any{}
	file_services_notificator_notificator_proto_msgTypes[2].OneofWrappers = []any{}
	file_services_notificator_notificator_proto_msgTypes[5].OneofWrappers = []any{
		(*StreamResponse_UserEvent)(nil),
		(*StreamResponse_JobEvent)(nil),
		(*StreamResponse_JobGradeEvent)(nil),
		(*StreamResponse_SystemEvent)(nil),
		(*StreamResponse_MailerEvent)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_services_notificator_notificator_proto_rawDesc), len(file_services_notificator_notificator_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_notificator_notificator_proto_goTypes,
		DependencyIndexes: file_services_notificator_notificator_proto_depIdxs,
		MessageInfos:      file_services_notificator_notificator_proto_msgTypes,
	}.Build()
	File_services_notificator_notificator_proto = out.File
	file_services_notificator_notificator_proto_goTypes = nil
	file_services_notificator_notificator_proto_depIdxs = nil
}
