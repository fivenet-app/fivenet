// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v3.20.3
// source: resources/documents/workflow.proto

package documents

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Workflow struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	Reminder          bool                   `protobuf:"varint,1,opt,name=reminder,proto3" json:"reminder,omitempty"`
	ReminderSettings  *ReminderSettings      `protobuf:"bytes,2,opt,name=reminder_settings,json=reminderSettings,proto3" json:"reminder_settings,omitempty"`
	AutoClose         bool                   `protobuf:"varint,3,opt,name=auto_close,json=autoClose,proto3" json:"auto_close,omitempty"`
	AutoCloseSettings *AutoCloseSettings     `protobuf:"bytes,4,opt,name=auto_close_settings,json=autoCloseSettings,proto3" json:"auto_close_settings,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *Workflow) Reset() {
	*x = Workflow{}
	mi := &file_resources_documents_workflow_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Workflow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Workflow) ProtoMessage() {}

func (x *Workflow) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_workflow_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Workflow.ProtoReflect.Descriptor instead.
func (*Workflow) Descriptor() ([]byte, []int) {
	return file_resources_documents_workflow_proto_rawDescGZIP(), []int{0}
}

func (x *Workflow) GetReminder() bool {
	if x != nil {
		return x.Reminder
	}
	return false
}

func (x *Workflow) GetReminderSettings() *ReminderSettings {
	if x != nil {
		return x.ReminderSettings
	}
	return nil
}

func (x *Workflow) GetAutoClose() bool {
	if x != nil {
		return x.AutoClose
	}
	return false
}

func (x *Workflow) GetAutoCloseSettings() *AutoCloseSettings {
	if x != nil {
		return x.AutoCloseSettings
	}
	return nil
}

type ReminderSettings struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Reminders     []*Reminder            `protobuf:"bytes,1,rep,name=reminders,proto3" json:"reminders,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReminderSettings) Reset() {
	*x = ReminderSettings{}
	mi := &file_resources_documents_workflow_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReminderSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReminderSettings) ProtoMessage() {}

func (x *ReminderSettings) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_workflow_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReminderSettings.ProtoReflect.Descriptor instead.
func (*ReminderSettings) Descriptor() ([]byte, []int) {
	return file_resources_documents_workflow_proto_rawDescGZIP(), []int{1}
}

func (x *ReminderSettings) GetReminders() []*Reminder {
	if x != nil {
		return x.Reminders
	}
	return nil
}

type Reminder struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Duration      *durationpb.Duration   `protobuf:"bytes,1,opt,name=duration,proto3" json:"duration,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Reminder) Reset() {
	*x = Reminder{}
	mi := &file_resources_documents_workflow_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Reminder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reminder) ProtoMessage() {}

func (x *Reminder) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_workflow_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reminder.ProtoReflect.Descriptor instead.
func (*Reminder) Descriptor() ([]byte, []int) {
	return file_resources_documents_workflow_proto_rawDescGZIP(), []int{2}
}

func (x *Reminder) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *Reminder) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type AutoCloseSettings struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Duration      *durationpb.Duration   `protobuf:"bytes,1,opt,name=duration,proto3" json:"duration,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AutoCloseSettings) Reset() {
	*x = AutoCloseSettings{}
	mi := &file_resources_documents_workflow_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AutoCloseSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AutoCloseSettings) ProtoMessage() {}

func (x *AutoCloseSettings) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_workflow_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AutoCloseSettings.ProtoReflect.Descriptor instead.
func (*AutoCloseSettings) Descriptor() ([]byte, []int) {
	return file_resources_documents_workflow_proto_rawDescGZIP(), []int{3}
}

func (x *AutoCloseSettings) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *AutoCloseSettings) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type WorkflowCronData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LastDocId     uint64                 `protobuf:"varint,1,opt,name=last_doc_id,json=lastDocId,proto3" json:"last_doc_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *WorkflowCronData) Reset() {
	*x = WorkflowCronData{}
	mi := &file_resources_documents_workflow_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WorkflowCronData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkflowCronData) ProtoMessage() {}

func (x *WorkflowCronData) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_workflow_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkflowCronData.ProtoReflect.Descriptor instead.
func (*WorkflowCronData) Descriptor() ([]byte, []int) {
	return file_resources_documents_workflow_proto_rawDescGZIP(), []int{4}
}

func (x *WorkflowCronData) GetLastDocId() uint64 {
	if x != nil {
		return x.LastDocId
	}
	return 0
}

var File_resources_documents_workflow_proto protoreflect.FileDescriptor

var file_resources_documents_workflow_proto_rawDesc = []byte{
	0x0a, 0x22, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xf1, 0x01, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x12,
	0x1a, 0x0a, 0x08, 0x72, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x08, 0x72, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x52, 0x0a, 0x11, 0x72,
	0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x52, 0x65, 0x6d,
	0x69, 0x6e, 0x64, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x10, 0x72,
	0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12,
	0x1d, 0x0a, 0x0a, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x61, 0x75, 0x74, 0x6f, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x56,
	0x0a, 0x13, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x73, 0x65, 0x74,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x41, 0x75, 0x74, 0x6f, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x53, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x73, 0x52, 0x11, 0x61, 0x75, 0x74, 0x6f, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0x59, 0x0a, 0x10, 0x52, 0x65, 0x6d, 0x69, 0x6e, 0x64,
	0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x45, 0x0a, 0x09, 0x72, 0x65,
	0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x52, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x42, 0x08, 0xfa, 0x42,
	0x05, 0x92, 0x01, 0x02, 0x10, 0x03, 0x52, 0x09, 0x72, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72,
	0x73, 0x22, 0x7c, 0x0a, 0x08, 0x52, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x4c, 0x0a,
	0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x15, 0xfa, 0x42, 0x12, 0xaa,
	0x01, 0x0f, 0x08, 0x01, 0x1a, 0x05, 0x08, 0x80, 0xce, 0xda, 0x03, 0x32, 0x04, 0x08, 0x80, 0xa3,
	0x05, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42,
	0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x85, 0x01, 0x0a, 0x11, 0x41, 0x75, 0x74, 0x6f, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x53, 0x65, 0x74,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x4c, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x15, 0xfa, 0x42, 0x12, 0xaa, 0x01, 0x0f, 0x08, 0x01, 0x1a, 0x05, 0x08, 0x80,
	0xce, 0xda, 0x03, 0x32, 0x04, 0x08, 0x80, 0xa3, 0x05, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x32, 0x0a, 0x10, 0x57, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x43, 0x72, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1e, 0x0a, 0x0b, 0x6c,
	0x61, 0x73, 0x74, 0x5f, 0x64, 0x6f, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x44, 0x6f, 0x63, 0x49, 0x64, 0x42, 0x4b, 0x5a, 0x49, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65,
	0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x3b, 0x64,
	0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_documents_workflow_proto_rawDescOnce sync.Once
	file_resources_documents_workflow_proto_rawDescData = file_resources_documents_workflow_proto_rawDesc
)

func file_resources_documents_workflow_proto_rawDescGZIP() []byte {
	file_resources_documents_workflow_proto_rawDescOnce.Do(func() {
		file_resources_documents_workflow_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_documents_workflow_proto_rawDescData)
	})
	return file_resources_documents_workflow_proto_rawDescData
}

var file_resources_documents_workflow_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_resources_documents_workflow_proto_goTypes = []any{
	(*Workflow)(nil),            // 0: resources.documents.Workflow
	(*ReminderSettings)(nil),    // 1: resources.documents.ReminderSettings
	(*Reminder)(nil),            // 2: resources.documents.Reminder
	(*AutoCloseSettings)(nil),   // 3: resources.documents.AutoCloseSettings
	(*WorkflowCronData)(nil),    // 4: resources.documents.WorkflowCronData
	(*durationpb.Duration)(nil), // 5: google.protobuf.Duration
}
var file_resources_documents_workflow_proto_depIdxs = []int32{
	1, // 0: resources.documents.Workflow.reminder_settings:type_name -> resources.documents.ReminderSettings
	3, // 1: resources.documents.Workflow.auto_close_settings:type_name -> resources.documents.AutoCloseSettings
	2, // 2: resources.documents.ReminderSettings.reminders:type_name -> resources.documents.Reminder
	5, // 3: resources.documents.Reminder.duration:type_name -> google.protobuf.Duration
	5, // 4: resources.documents.AutoCloseSettings.duration:type_name -> google.protobuf.Duration
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_resources_documents_workflow_proto_init() }
func file_resources_documents_workflow_proto_init() {
	if File_resources_documents_workflow_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_documents_workflow_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_documents_workflow_proto_goTypes,
		DependencyIndexes: file_resources_documents_workflow_proto_depIdxs,
		MessageInfos:      file_resources_documents_workflow_proto_msgTypes,
	}.Build()
	File_resources_documents_workflow_proto = out.File
	file_resources_documents_workflow_proto_rawDesc = nil
	file_resources_documents_workflow_proto_goTypes = nil
	file_resources_documents_workflow_proto_depIdxs = nil
}
