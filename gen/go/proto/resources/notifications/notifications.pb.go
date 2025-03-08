// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: resources/notifications/notifications.proto

package notifications

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	common "github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
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

type NotificationType int32

const (
	NotificationType_NOTIFICATION_TYPE_UNSPECIFIED NotificationType = 0
	NotificationType_NOTIFICATION_TYPE_ERROR       NotificationType = 1
	NotificationType_NOTIFICATION_TYPE_WARNING     NotificationType = 2
	NotificationType_NOTIFICATION_TYPE_INFO        NotificationType = 3
	NotificationType_NOTIFICATION_TYPE_SUCCESS     NotificationType = 4
)

// Enum value maps for NotificationType.
var (
	NotificationType_name = map[int32]string{
		0: "NOTIFICATION_TYPE_UNSPECIFIED",
		1: "NOTIFICATION_TYPE_ERROR",
		2: "NOTIFICATION_TYPE_WARNING",
		3: "NOTIFICATION_TYPE_INFO",
		4: "NOTIFICATION_TYPE_SUCCESS",
	}
	NotificationType_value = map[string]int32{
		"NOTIFICATION_TYPE_UNSPECIFIED": 0,
		"NOTIFICATION_TYPE_ERROR":       1,
		"NOTIFICATION_TYPE_WARNING":     2,
		"NOTIFICATION_TYPE_INFO":        3,
		"NOTIFICATION_TYPE_SUCCESS":     4,
	}
)

func (x NotificationType) Enum() *NotificationType {
	p := new(NotificationType)
	*p = x
	return p
}

func (x NotificationType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NotificationType) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_notifications_notifications_proto_enumTypes[0].Descriptor()
}

func (NotificationType) Type() protoreflect.EnumType {
	return &file_resources_notifications_notifications_proto_enumTypes[0]
}

func (x NotificationType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NotificationType.Descriptor instead.
func (NotificationType) EnumDescriptor() ([]byte, []int) {
	return file_resources_notifications_notifications_proto_rawDescGZIP(), []int{0}
}

type NotificationCategory int32

const (
	NotificationCategory_NOTIFICATION_CATEGORY_UNSPECIFIED NotificationCategory = 0
	NotificationCategory_NOTIFICATION_CATEGORY_GENERAL     NotificationCategory = 1
	NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT    NotificationCategory = 2
	NotificationCategory_NOTIFICATION_CATEGORY_CALENDAR    NotificationCategory = 3
)

// Enum value maps for NotificationCategory.
var (
	NotificationCategory_name = map[int32]string{
		0: "NOTIFICATION_CATEGORY_UNSPECIFIED",
		1: "NOTIFICATION_CATEGORY_GENERAL",
		2: "NOTIFICATION_CATEGORY_DOCUMENT",
		3: "NOTIFICATION_CATEGORY_CALENDAR",
	}
	NotificationCategory_value = map[string]int32{
		"NOTIFICATION_CATEGORY_UNSPECIFIED": 0,
		"NOTIFICATION_CATEGORY_GENERAL":     1,
		"NOTIFICATION_CATEGORY_DOCUMENT":    2,
		"NOTIFICATION_CATEGORY_CALENDAR":    3,
	}
)

func (x NotificationCategory) Enum() *NotificationCategory {
	p := new(NotificationCategory)
	*p = x
	return p
}

func (x NotificationCategory) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NotificationCategory) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_notifications_notifications_proto_enumTypes[1].Descriptor()
}

func (NotificationCategory) Type() protoreflect.EnumType {
	return &file_resources_notifications_notifications_proto_enumTypes[1]
}

func (x NotificationCategory) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NotificationCategory.Descriptor instead.
func (NotificationCategory) EnumDescriptor() ([]byte, []int) {
	return file_resources_notifications_notifications_proto_rawDescGZIP(), []int{1}
}

type Notification struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	ReadAt    *timestamp.Timestamp   `protobuf:"bytes,3,opt,name=read_at,json=readAt,proto3" json:"read_at,omitempty"`
	UserId    int32                  `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// @sanitize
	Title *common.TranslateItem `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	Type  NotificationType      `protobuf:"varint,6,opt,name=type,proto3,enum=resources.notifications.NotificationType" json:"type,omitempty"`
	// @sanitize
	Content       *common.TranslateItem `protobuf:"bytes,7,opt,name=content,proto3" json:"content,omitempty"`
	Category      NotificationCategory  `protobuf:"varint,8,opt,name=category,proto3,enum=resources.notifications.NotificationCategory" json:"category,omitempty"`
	Data          *Data                 `protobuf:"bytes,9,opt,name=data,proto3,oneof" json:"data,omitempty"`
	Starred       *bool                 `protobuf:"varint,10,opt,name=starred,proto3,oneof" json:"starred,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Notification) Reset() {
	*x = Notification{}
	mi := &file_resources_notifications_notifications_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification) ProtoMessage() {}

func (x *Notification) ProtoReflect() protoreflect.Message {
	mi := &file_resources_notifications_notifications_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notification.ProtoReflect.Descriptor instead.
func (*Notification) Descriptor() ([]byte, []int) {
	return file_resources_notifications_notifications_proto_rawDescGZIP(), []int{0}
}

func (x *Notification) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Notification) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Notification) GetReadAt() *timestamp.Timestamp {
	if x != nil {
		return x.ReadAt
	}
	return nil
}

func (x *Notification) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Notification) GetTitle() *common.TranslateItem {
	if x != nil {
		return x.Title
	}
	return nil
}

func (x *Notification) GetType() NotificationType {
	if x != nil {
		return x.Type
	}
	return NotificationType_NOTIFICATION_TYPE_UNSPECIFIED
}

func (x *Notification) GetContent() *common.TranslateItem {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *Notification) GetCategory() NotificationCategory {
	if x != nil {
		return x.Category
	}
	return NotificationCategory_NOTIFICATION_CATEGORY_UNSPECIFIED
}

func (x *Notification) GetData() *Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Notification) GetStarred() bool {
	if x != nil && x.Starred != nil {
		return *x.Starred
	}
	return false
}

type Data struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Link          *Link                  `protobuf:"bytes,1,opt,name=link,proto3,oneof" json:"link,omitempty"`
	CausedBy      *users.UserShort       `protobuf:"bytes,2,opt,name=caused_by,json=causedBy,proto3,oneof" json:"caused_by,omitempty"`
	Calendar      *CalendarData          `protobuf:"bytes,3,opt,name=calendar,proto3,oneof" json:"calendar,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Data) Reset() {
	*x = Data{}
	mi := &file_resources_notifications_notifications_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_resources_notifications_notifications_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_resources_notifications_notifications_proto_rawDescGZIP(), []int{1}
}

func (x *Data) GetLink() *Link {
	if x != nil {
		return x.Link
	}
	return nil
}

func (x *Data) GetCausedBy() *users.UserShort {
	if x != nil {
		return x.CausedBy
	}
	return nil
}

func (x *Data) GetCalendar() *CalendarData {
	if x != nil {
		return x.Calendar
	}
	return nil
}

type Link struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	To            string                 `protobuf:"bytes,1,opt,name=to,proto3" json:"to,omitempty"`
	Title         *string                `protobuf:"bytes,2,opt,name=title,proto3,oneof" json:"title,omitempty"`
	External      *bool                  `protobuf:"varint,3,opt,name=external,proto3,oneof" json:"external,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Link) Reset() {
	*x = Link{}
	mi := &file_resources_notifications_notifications_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Link) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Link) ProtoMessage() {}

func (x *Link) ProtoReflect() protoreflect.Message {
	mi := &file_resources_notifications_notifications_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Link.ProtoReflect.Descriptor instead.
func (*Link) Descriptor() ([]byte, []int) {
	return file_resources_notifications_notifications_proto_rawDescGZIP(), []int{2}
}

func (x *Link) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *Link) GetTitle() string {
	if x != nil && x.Title != nil {
		return *x.Title
	}
	return ""
}

func (x *Link) GetExternal() bool {
	if x != nil && x.External != nil {
		return *x.External
	}
	return false
}

type CalendarData struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	CalendarId      *uint64                `protobuf:"varint,1,opt,name=calendar_id,json=calendarId,proto3,oneof" json:"calendar_id,omitempty"`
	CalendarEntryId *uint64                `protobuf:"varint,2,opt,name=calendar_entry_id,json=calendarEntryId,proto3,oneof" json:"calendar_entry_id,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *CalendarData) Reset() {
	*x = CalendarData{}
	mi := &file_resources_notifications_notifications_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CalendarData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CalendarData) ProtoMessage() {}

func (x *CalendarData) ProtoReflect() protoreflect.Message {
	mi := &file_resources_notifications_notifications_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CalendarData.ProtoReflect.Descriptor instead.
func (*CalendarData) Descriptor() ([]byte, []int) {
	return file_resources_notifications_notifications_proto_rawDescGZIP(), []int{3}
}

func (x *CalendarData) GetCalendarId() uint64 {
	if x != nil && x.CalendarId != nil {
		return *x.CalendarId
	}
	return 0
}

func (x *CalendarData) GetCalendarEntryId() uint64 {
	if x != nil && x.CalendarEntryId != nil {
		return *x.CalendarEntryId
	}
	return 0
}

var File_resources_notifications_notifications_proto protoreflect.FileDescriptor

var file_resources_notifications_notifications_proto_rawDesc = string([]byte{
	0x0a, 0x2b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x69, 0x31, 0x38, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb4,
	0x04, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x3d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x37,
	0x0a, 0x07, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x06, 0x72, 0x65, 0x61, 0x64, 0x41, 0x74, 0x12, 0x20, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x28,
	0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x6c, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x47, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01,
	0x02, 0x10, 0x01, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x39, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x12, 0x53, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52,
	0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x36, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x88, 0x01,
	0x01, 0x12, 0x1d, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x72, 0x72, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x08, 0x48, 0x01, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x72, 0x65, 0x64, 0x88, 0x01, 0x01,
	0x42, 0x07, 0x0a, 0x05, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x73, 0x74,
	0x61, 0x72, 0x72, 0x65, 0x64, 0x22, 0xe8, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x36,
	0x0a, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x48, 0x00, 0x52, 0x04, 0x6c,
	0x69, 0x6e, 0x6b, 0x88, 0x01, 0x01, 0x12, 0x3c, 0x0a, 0x09, 0x63, 0x61, 0x75, 0x73, 0x65, 0x64,
	0x5f, 0x62, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x53, 0x68, 0x6f, 0x72, 0x74, 0x48, 0x01, 0x52, 0x08, 0x63, 0x61, 0x75, 0x73, 0x65, 0x64, 0x42,
	0x79, 0x88, 0x01, 0x01, 0x12, 0x46, 0x0a, 0x08, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x43, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x44, 0x61, 0x74, 0x61, 0x48, 0x02, 0x52,
	0x08, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05,
	0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x63, 0x61, 0x75, 0x73, 0x65, 0x64,
	0x5f, 0x62, 0x79, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72,
	0x22, 0x69, 0x0a, 0x04, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x19, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x08, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x42, 0x0b,
	0x0a, 0x09, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x22, 0x8b, 0x01, 0x0a, 0x0c,
	0x43, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x24, 0x0a, 0x0b,
	0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x48, 0x00, 0x52, 0x0a, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x49, 0x64, 0x88,
	0x01, 0x01, 0x12, 0x2f, 0x0a, 0x11, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x5f, 0x65,
	0x6e, 0x74, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x48, 0x01, 0x52,
	0x0f, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x49, 0x64,
	0x88, 0x01, 0x01, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72,
	0x5f, 0x69, 0x64, 0x42, 0x14, 0x0a, 0x12, 0x5f, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72,
	0x5f, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x2a, 0xac, 0x01, 0x0a, 0x10, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21,
	0x0a, 0x1d, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x1b, 0x0a, 0x17, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x01, 0x12, 0x1d,
	0x0a, 0x19, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x1a, 0x0a,
	0x16, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x03, 0x12, 0x1d, 0x0a, 0x19, 0x4e, 0x4f, 0x54,
	0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53,
	0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x04, 0x2a, 0xa8, 0x01, 0x0a, 0x14, 0x4e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x79, 0x12, 0x25, 0x0a, 0x21, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x43, 0x41, 0x54, 0x45, 0x47, 0x4f, 0x52, 0x59, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x21, 0x0a, 0x1d, 0x4e, 0x4f, 0x54, 0x49,
	0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x41, 0x54, 0x45, 0x47, 0x4f, 0x52,
	0x59, 0x5f, 0x47, 0x45, 0x4e, 0x45, 0x52, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x22, 0x0a, 0x1e, 0x4e,
	0x4f, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x41, 0x54, 0x45,
	0x47, 0x4f, 0x52, 0x59, 0x5f, 0x44, 0x4f, 0x43, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x02, 0x12,
	0x22, 0x0a, 0x1e, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x43, 0x41, 0x54, 0x45, 0x47, 0x4f, 0x52, 0x59, 0x5f, 0x43, 0x41, 0x4c, 0x45, 0x4e, 0x44, 0x41,
	0x52, 0x10, 0x03, 0x42, 0x53, 0x5a, 0x51, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69,
	0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3b, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_notifications_notifications_proto_rawDescOnce sync.Once
	file_resources_notifications_notifications_proto_rawDescData []byte
)

func file_resources_notifications_notifications_proto_rawDescGZIP() []byte {
	file_resources_notifications_notifications_proto_rawDescOnce.Do(func() {
		file_resources_notifications_notifications_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_notifications_notifications_proto_rawDesc), len(file_resources_notifications_notifications_proto_rawDesc)))
	})
	return file_resources_notifications_notifications_proto_rawDescData
}

var file_resources_notifications_notifications_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_resources_notifications_notifications_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_resources_notifications_notifications_proto_goTypes = []any{
	(NotificationType)(0),        // 0: resources.notifications.NotificationType
	(NotificationCategory)(0),    // 1: resources.notifications.NotificationCategory
	(*Notification)(nil),         // 2: resources.notifications.Notification
	(*Data)(nil),                 // 3: resources.notifications.Data
	(*Link)(nil),                 // 4: resources.notifications.Link
	(*CalendarData)(nil),         // 5: resources.notifications.CalendarData
	(*timestamp.Timestamp)(nil),  // 6: resources.timestamp.Timestamp
	(*common.TranslateItem)(nil), // 7: resources.common.TranslateItem
	(*users.UserShort)(nil),      // 8: resources.users.UserShort
}
var file_resources_notifications_notifications_proto_depIdxs = []int32{
	6,  // 0: resources.notifications.Notification.created_at:type_name -> resources.timestamp.Timestamp
	6,  // 1: resources.notifications.Notification.read_at:type_name -> resources.timestamp.Timestamp
	7,  // 2: resources.notifications.Notification.title:type_name -> resources.common.TranslateItem
	0,  // 3: resources.notifications.Notification.type:type_name -> resources.notifications.NotificationType
	7,  // 4: resources.notifications.Notification.content:type_name -> resources.common.TranslateItem
	1,  // 5: resources.notifications.Notification.category:type_name -> resources.notifications.NotificationCategory
	3,  // 6: resources.notifications.Notification.data:type_name -> resources.notifications.Data
	4,  // 7: resources.notifications.Data.link:type_name -> resources.notifications.Link
	8,  // 8: resources.notifications.Data.caused_by:type_name -> resources.users.UserShort
	5,  // 9: resources.notifications.Data.calendar:type_name -> resources.notifications.CalendarData
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_resources_notifications_notifications_proto_init() }
func file_resources_notifications_notifications_proto_init() {
	if File_resources_notifications_notifications_proto != nil {
		return
	}
	file_resources_notifications_notifications_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_notifications_notifications_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_notifications_notifications_proto_msgTypes[2].OneofWrappers = []any{}
	file_resources_notifications_notifications_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_notifications_notifications_proto_rawDesc), len(file_resources_notifications_notifications_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_notifications_notifications_proto_goTypes,
		DependencyIndexes: file_resources_notifications_notifications_proto_depIdxs,
		EnumInfos:         file_resources_notifications_notifications_proto_enumTypes,
		MessageInfos:      file_resources_notifications_notifications_proto_msgTypes,
	}.Build()
	File_resources_notifications_notifications_proto = out.File
	file_resources_notifications_notifications_proto_goTypes = nil
	file_resources_notifications_notifications_proto_depIdxs = nil
}
