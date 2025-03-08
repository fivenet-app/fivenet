// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: resources/notifications/events.proto

package notifications

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type UserEvent struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*UserEvent_RefreshToken
	//	*UserEvent_Notification
	//	*UserEvent_NotificationsReadCount
	Data          isUserEvent_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserEvent) Reset() {
	*x = UserEvent{}
	mi := &file_resources_notifications_events_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserEvent) ProtoMessage() {}

func (x *UserEvent) ProtoReflect() protoreflect.Message {
	mi := &file_resources_notifications_events_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserEvent.ProtoReflect.Descriptor instead.
func (*UserEvent) Descriptor() ([]byte, []int) {
	return file_resources_notifications_events_proto_rawDescGZIP(), []int{0}
}

func (x *UserEvent) GetData() isUserEvent_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *UserEvent) GetRefreshToken() bool {
	if x != nil {
		if x, ok := x.Data.(*UserEvent_RefreshToken); ok {
			return x.RefreshToken
		}
	}
	return false
}

func (x *UserEvent) GetNotification() *Notification {
	if x != nil {
		if x, ok := x.Data.(*UserEvent_Notification); ok {
			return x.Notification
		}
	}
	return nil
}

func (x *UserEvent) GetNotificationsReadCount() int32 {
	if x != nil {
		if x, ok := x.Data.(*UserEvent_NotificationsReadCount); ok {
			return x.NotificationsReadCount
		}
	}
	return 0
}

type isUserEvent_Data interface {
	isUserEvent_Data()
}

type UserEvent_RefreshToken struct {
	RefreshToken bool `protobuf:"varint,1,opt,name=refresh_token,json=refreshToken,proto3,oneof"`
}

type UserEvent_Notification struct {
	// Notifications
	Notification *Notification `protobuf:"bytes,2,opt,name=notification,proto3,oneof"`
}

type UserEvent_NotificationsReadCount struct {
	NotificationsReadCount int32 `protobuf:"varint,3,opt,name=notifications_read_count,json=notificationsReadCount,proto3,oneof"`
}

func (*UserEvent_RefreshToken) isUserEvent_Data() {}

func (*UserEvent_Notification) isUserEvent_Data() {}

func (*UserEvent_NotificationsReadCount) isUserEvent_Data() {}

type JobEvent struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*JobEvent_JobProps
	Data          isJobEvent_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JobEvent) Reset() {
	*x = JobEvent{}
	mi := &file_resources_notifications_events_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobEvent) ProtoMessage() {}

func (x *JobEvent) ProtoReflect() protoreflect.Message {
	mi := &file_resources_notifications_events_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobEvent.ProtoReflect.Descriptor instead.
func (*JobEvent) Descriptor() ([]byte, []int) {
	return file_resources_notifications_events_proto_rawDescGZIP(), []int{1}
}

func (x *JobEvent) GetData() isJobEvent_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *JobEvent) GetJobProps() *users.JobProps {
	if x != nil {
		if x, ok := x.Data.(*JobEvent_JobProps); ok {
			return x.JobProps
		}
	}
	return nil
}

type isJobEvent_Data interface {
	isJobEvent_Data()
}

type JobEvent_JobProps struct {
	JobProps *users.JobProps `protobuf:"bytes,1,opt,name=job_props,json=jobProps,proto3,oneof"`
}

func (*JobEvent_JobProps) isJobEvent_Data() {}

type JobGradeEvent struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*JobGradeEvent_RefreshToken
	Data          isJobGradeEvent_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JobGradeEvent) Reset() {
	*x = JobGradeEvent{}
	mi := &file_resources_notifications_events_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobGradeEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobGradeEvent) ProtoMessage() {}

func (x *JobGradeEvent) ProtoReflect() protoreflect.Message {
	mi := &file_resources_notifications_events_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobGradeEvent.ProtoReflect.Descriptor instead.
func (*JobGradeEvent) Descriptor() ([]byte, []int) {
	return file_resources_notifications_events_proto_rawDescGZIP(), []int{2}
}

func (x *JobGradeEvent) GetData() isJobGradeEvent_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *JobGradeEvent) GetRefreshToken() bool {
	if x != nil {
		if x, ok := x.Data.(*JobGradeEvent_RefreshToken); ok {
			return x.RefreshToken
		}
	}
	return false
}

type isJobGradeEvent_Data interface {
	isJobGradeEvent_Data()
}

type JobGradeEvent_RefreshToken struct {
	RefreshToken bool `protobuf:"varint,1,opt,name=refresh_token,json=refreshToken,proto3,oneof"`
}

func (*JobGradeEvent_RefreshToken) isJobGradeEvent_Data() {}

type SystemEvent struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*SystemEvent_Ping
	Data          isSystemEvent_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SystemEvent) Reset() {
	*x = SystemEvent{}
	mi := &file_resources_notifications_events_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SystemEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SystemEvent) ProtoMessage() {}

func (x *SystemEvent) ProtoReflect() protoreflect.Message {
	mi := &file_resources_notifications_events_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SystemEvent.ProtoReflect.Descriptor instead.
func (*SystemEvent) Descriptor() ([]byte, []int) {
	return file_resources_notifications_events_proto_rawDescGZIP(), []int{3}
}

func (x *SystemEvent) GetData() isSystemEvent_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *SystemEvent) GetPing() bool {
	if x != nil {
		if x, ok := x.Data.(*SystemEvent_Ping); ok {
			return x.Ping
		}
	}
	return false
}

type isSystemEvent_Data interface {
	isSystemEvent_Data()
}

type SystemEvent_Ping struct {
	Ping bool `protobuf:"varint,1,opt,name=ping,proto3,oneof"`
}

func (*SystemEvent_Ping) isSystemEvent_Data() {}

var File_resources_notifications_events_proto protoreflect.FileDescriptor

var file_resources_notifications_events_proto_rawDesc = string([]byte{
	0x0a, 0x24, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a,
	0x2b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x6a, 0x6f,
	0x62, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc8, 0x01, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x0c, 0x72,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x4b, 0x0a, 0x0c, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x25, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0c, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x18, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x16, 0x6e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x61, 0x64, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x42, 0x0b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x03, 0xf8, 0x42,
	0x01, 0x22, 0x51, 0x0a, 0x08, 0x4a, 0x6f, 0x62, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x38, 0x0a,
	0x09, 0x6a, 0x6f, 0x62, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x48, 0x00, 0x52, 0x08, 0x6a,
	0x6f, 0x62, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x42, 0x0b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x03, 0xf8, 0x42, 0x01, 0x22, 0x43, 0x0a, 0x0d, 0x4a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64, 0x65,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x0c,
	0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x0b, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x03, 0xf8, 0x42, 0x01, 0x22, 0x30, 0x0a, 0x0b, 0x53, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x04, 0x70, 0x69, 0x6e, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x04, 0x70, 0x69, 0x6e, 0x67, 0x42, 0x0b,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x03, 0xf8, 0x42, 0x01, 0x42, 0x53, 0x5a, 0x51, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65,
	0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x3b, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_notifications_events_proto_rawDescOnce sync.Once
	file_resources_notifications_events_proto_rawDescData []byte
)

func file_resources_notifications_events_proto_rawDescGZIP() []byte {
	file_resources_notifications_events_proto_rawDescOnce.Do(func() {
		file_resources_notifications_events_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_notifications_events_proto_rawDesc), len(file_resources_notifications_events_proto_rawDesc)))
	})
	return file_resources_notifications_events_proto_rawDescData
}

var file_resources_notifications_events_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_resources_notifications_events_proto_goTypes = []any{
	(*UserEvent)(nil),      // 0: resources.notifications.UserEvent
	(*JobEvent)(nil),       // 1: resources.notifications.JobEvent
	(*JobGradeEvent)(nil),  // 2: resources.notifications.JobGradeEvent
	(*SystemEvent)(nil),    // 3: resources.notifications.SystemEvent
	(*Notification)(nil),   // 4: resources.notifications.Notification
	(*users.JobProps)(nil), // 5: resources.users.JobProps
}
var file_resources_notifications_events_proto_depIdxs = []int32{
	4, // 0: resources.notifications.UserEvent.notification:type_name -> resources.notifications.Notification
	5, // 1: resources.notifications.JobEvent.job_props:type_name -> resources.users.JobProps
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_resources_notifications_events_proto_init() }
func file_resources_notifications_events_proto_init() {
	if File_resources_notifications_events_proto != nil {
		return
	}
	file_resources_notifications_notifications_proto_init()
	file_resources_notifications_events_proto_msgTypes[0].OneofWrappers = []any{
		(*UserEvent_RefreshToken)(nil),
		(*UserEvent_Notification)(nil),
		(*UserEvent_NotificationsReadCount)(nil),
	}
	file_resources_notifications_events_proto_msgTypes[1].OneofWrappers = []any{
		(*JobEvent_JobProps)(nil),
	}
	file_resources_notifications_events_proto_msgTypes[2].OneofWrappers = []any{
		(*JobGradeEvent_RefreshToken)(nil),
	}
	file_resources_notifications_events_proto_msgTypes[3].OneofWrappers = []any{
		(*SystemEvent_Ping)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_notifications_events_proto_rawDesc), len(file_resources_notifications_events_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_notifications_events_proto_goTypes,
		DependencyIndexes: file_resources_notifications_events_proto_depIdxs,
		MessageInfos:      file_resources_notifications_events_proto_msgTypes,
	}.Build()
	File_resources_notifications_events_proto = out.File
	file_resources_notifications_events_proto_goTypes = nil
	file_resources_notifications_events_proto_depIdxs = nil
}
