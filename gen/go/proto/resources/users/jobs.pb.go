// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: resources/users/jobs.proto

package users

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	timestamp "github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
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

type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" sql:"primary_key" alias:"name"` // @gotags: sql:"primary_key" alias:"name"
	Label  string      `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	Grades []*JobGrade `protobuf:"bytes,3,rep,name=grades,proto3" json:"grades,omitempty"`
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_users_jobs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_jobs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_resources_users_jobs_proto_rawDescGZIP(), []int{0}
}

func (x *Job) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Job) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *Job) GetGrades() []*JobGrade {
	if x != nil {
		return x.Grades
	}
	return nil
}

type JobGrade struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName *string `protobuf:"bytes,1,opt,name=job_name,json=jobName,proto3,oneof" json:"job_name,omitempty"`
	Grade   int32   `protobuf:"varint,2,opt,name=grade,proto3" json:"grade,omitempty"`
	Label   string  `protobuf:"bytes,3,opt,name=label,proto3" json:"label,omitempty"`
}

func (x *JobGrade) Reset() {
	*x = JobGrade{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_users_jobs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobGrade) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobGrade) ProtoMessage() {}

func (x *JobGrade) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_jobs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobGrade.ProtoReflect.Descriptor instead.
func (*JobGrade) Descriptor() ([]byte, []int) {
	return file_resources_users_jobs_proto_rawDescGZIP(), []int{1}
}

func (x *JobGrade) GetJobName() string {
	if x != nil && x.JobName != nil {
		return *x.JobName
	}
	return ""
}

func (x *JobGrade) GetGrade() int32 {
	if x != nil {
		return x.Grade
	}
	return 0
}

func (x *JobGrade) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

type JobProps struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job                 string               `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
	Theme               string               `protobuf:"bytes,2,opt,name=theme,proto3" json:"theme,omitempty"`
	LivemapMarkerColor  string               `protobuf:"bytes,3,opt,name=livemap_marker_color,json=livemapMarkerColor,proto3" json:"livemap_marker_color,omitempty"`
	QuickButtons        *QuickButtons        `protobuf:"bytes,4,opt,name=quick_buttons,json=quickButtons,proto3" json:"quick_buttons,omitempty"`
	DiscordGuildId      *uint64              `protobuf:"varint,5,opt,name=discord_guild_id,json=discordGuildId,proto3,oneof" json:"discord_guild_id,omitempty"`
	DiscordLastSync     *timestamp.Timestamp `protobuf:"bytes,6,opt,name=discord_last_sync,json=discordLastSync,proto3,oneof" json:"discord_last_sync,omitempty"`
	DiscordSyncSettings *DiscordSyncSettings `protobuf:"bytes,7,opt,name=discord_sync_settings,json=discordSyncSettings,proto3" json:"discord_sync_settings,omitempty"`
}

func (x *JobProps) Reset() {
	*x = JobProps{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_users_jobs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobProps) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobProps) ProtoMessage() {}

func (x *JobProps) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_jobs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobProps.ProtoReflect.Descriptor instead.
func (*JobProps) Descriptor() ([]byte, []int) {
	return file_resources_users_jobs_proto_rawDescGZIP(), []int{2}
}

func (x *JobProps) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *JobProps) GetTheme() string {
	if x != nil {
		return x.Theme
	}
	return ""
}

func (x *JobProps) GetLivemapMarkerColor() string {
	if x != nil {
		return x.LivemapMarkerColor
	}
	return ""
}

func (x *JobProps) GetQuickButtons() *QuickButtons {
	if x != nil {
		return x.QuickButtons
	}
	return nil
}

func (x *JobProps) GetDiscordGuildId() uint64 {
	if x != nil && x.DiscordGuildId != nil {
		return *x.DiscordGuildId
	}
	return 0
}

func (x *JobProps) GetDiscordLastSync() *timestamp.Timestamp {
	if x != nil {
		return x.DiscordLastSync
	}
	return nil
}

func (x *JobProps) GetDiscordSyncSettings() *DiscordSyncSettings {
	if x != nil {
		return x.DiscordSyncSettings
	}
	return nil
}

type QuickButtons struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PenaltyCalculator bool `protobuf:"varint,1,opt,name=penalty_calculator,json=penaltyCalculator,proto3" json:"penalty_calculator,omitempty"`
	BodyCheckup       bool `protobuf:"varint,2,opt,name=body_checkup,json=bodyCheckup,proto3" json:"body_checkup,omitempty"`
}

func (x *QuickButtons) Reset() {
	*x = QuickButtons{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_users_jobs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuickButtons) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuickButtons) ProtoMessage() {}

func (x *QuickButtons) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_jobs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuickButtons.ProtoReflect.Descriptor instead.
func (*QuickButtons) Descriptor() ([]byte, []int) {
	return file_resources_users_jobs_proto_rawDescGZIP(), []int{3}
}

func (x *QuickButtons) GetPenaltyCalculator() bool {
	if x != nil {
		return x.PenaltyCalculator
	}
	return false
}

func (x *QuickButtons) GetBodyCheckup() bool {
	if x != nil {
		return x.BodyCheckup
	}
	return false
}

type DiscordSyncSettings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserInfoSync bool `protobuf:"varint,1,opt,name=user_info_sync,json=userInfoSync,proto3" json:"user_info_sync,omitempty"`
}

func (x *DiscordSyncSettings) Reset() {
	*x = DiscordSyncSettings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_users_jobs_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscordSyncSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscordSyncSettings) ProtoMessage() {}

func (x *DiscordSyncSettings) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_jobs_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscordSyncSettings.ProtoReflect.Descriptor instead.
func (*DiscordSyncSettings) Descriptor() ([]byte, []int) {
	return file_resources_users_jobs_proto_rawDescGZIP(), []int{4}
}

func (x *DiscordSyncSettings) GetUserInfoSync() bool {
	if x != nil {
		return x.UserInfoSync
	}
	return false
}

var File_resources_users_jobs_proto protoreflect.FileDescriptor

var file_resources_users_jobs_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x23, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x74, 0x0a, 0x03, 0x4a,
	0x6f, 0x62, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1d, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x31,
	0x0a, 0x06, 0x67, 0x72, 0x61, 0x64, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x2e, 0x4a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64, 0x65, 0x52, 0x06, 0x67, 0x72, 0x61, 0x64, 0x65,
	0x73, 0x22, 0x7e, 0x0a, 0x08, 0x4a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64, 0x65, 0x12, 0x27, 0x0a,
	0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x00, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x4e,
	0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x05, 0x67, 0x72, 0x61, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x05,
	0x67, 0x72, 0x61, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x52, 0x05, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0xdf, 0x03, 0x0a, 0x08, 0x4a, 0x6f, 0x62, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x19,
	0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x18, 0x14, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x1d, 0x0a, 0x05, 0x74, 0x68, 0x65,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18,
	0x14, 0x52, 0x05, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x4c, 0x0a, 0x14, 0x6c, 0x69, 0x76, 0x65,
	0x6d, 0x61, 0x70, 0x5f, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1a, 0xfa, 0x42, 0x17, 0x72, 0x15, 0x32, 0x10, 0x5e,
	0x5b, 0x41, 0x2d, 0x46, 0x61, 0x2d, 0x66, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x36, 0x7d, 0x24, 0x98,
	0x01, 0x06, 0x52, 0x12, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x4d, 0x61, 0x72, 0x6b, 0x65,
	0x72, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x42, 0x0a, 0x0d, 0x71, 0x75, 0x69, 0x63, 0x6b, 0x5f,
	0x62, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e,
	0x51, 0x75, 0x69, 0x63, 0x6b, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x73, 0x52, 0x0c, 0x71, 0x75,
	0x69, 0x63, 0x6b, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x73, 0x12, 0x31, 0x0a, 0x10, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x30, 0x01, 0x48, 0x00, 0x52, 0x0e, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x72, 0x64, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x4f, 0x0a,
	0x11, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x79,
	0x6e, 0x63, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52, 0x0f, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x72, 0x64, 0x4c, 0x61, 0x73, 0x74, 0x53, 0x79, 0x6e, 0x63, 0x88, 0x01, 0x01, 0x12, 0x58,
	0x0a, 0x15, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x73,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x73, 0x52, 0x13, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x79, 0x6e, 0x63,
	0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x72, 0x64, 0x5f, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x69, 0x64, 0x42, 0x14, 0x0a,
	0x12, 0x5f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73,
	0x79, 0x6e, 0x63, 0x22, 0x60, 0x0a, 0x0c, 0x51, 0x75, 0x69, 0x63, 0x6b, 0x42, 0x75, 0x74, 0x74,
	0x6f, 0x6e, 0x73, 0x12, 0x2d, 0x0a, 0x12, 0x70, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79, 0x5f, 0x63,
	0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x11, 0x70, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79, 0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74,
	0x6f, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6f, 0x64, 0x79, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b,
	0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x62, 0x6f, 0x64, 0x79, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x75, 0x70, 0x22, 0x3b, 0x0a, 0x13, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64,
	0x53, 0x79, 0x6e, 0x63, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x24, 0x0a, 0x0e,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x5f, 0x73, 0x79, 0x6e, 0x63, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x79,
	0x6e, 0x63, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x3b, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_users_jobs_proto_rawDescOnce sync.Once
	file_resources_users_jobs_proto_rawDescData = file_resources_users_jobs_proto_rawDesc
)

func file_resources_users_jobs_proto_rawDescGZIP() []byte {
	file_resources_users_jobs_proto_rawDescOnce.Do(func() {
		file_resources_users_jobs_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_users_jobs_proto_rawDescData)
	})
	return file_resources_users_jobs_proto_rawDescData
}

var file_resources_users_jobs_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_resources_users_jobs_proto_goTypes = []interface{}{
	(*Job)(nil),                 // 0: resources.users.Job
	(*JobGrade)(nil),            // 1: resources.users.JobGrade
	(*JobProps)(nil),            // 2: resources.users.JobProps
	(*QuickButtons)(nil),        // 3: resources.users.QuickButtons
	(*DiscordSyncSettings)(nil), // 4: resources.users.DiscordSyncSettings
	(*timestamp.Timestamp)(nil), // 5: resources.timestamp.Timestamp
}
var file_resources_users_jobs_proto_depIdxs = []int32{
	1, // 0: resources.users.Job.grades:type_name -> resources.users.JobGrade
	3, // 1: resources.users.JobProps.quick_buttons:type_name -> resources.users.QuickButtons
	5, // 2: resources.users.JobProps.discord_last_sync:type_name -> resources.timestamp.Timestamp
	4, // 3: resources.users.JobProps.discord_sync_settings:type_name -> resources.users.DiscordSyncSettings
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_resources_users_jobs_proto_init() }
func file_resources_users_jobs_proto_init() {
	if File_resources_users_jobs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_users_jobs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job); i {
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
		file_resources_users_jobs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobGrade); i {
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
		file_resources_users_jobs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobProps); i {
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
		file_resources_users_jobs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuickButtons); i {
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
		file_resources_users_jobs_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscordSyncSettings); i {
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
	file_resources_users_jobs_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_resources_users_jobs_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_users_jobs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_users_jobs_proto_goTypes,
		DependencyIndexes: file_resources_users_jobs_proto_depIdxs,
		MessageInfos:      file_resources_users_jobs_proto_msgTypes,
	}.Build()
	File_resources_users_jobs_proto = out.File
	file_resources_users_jobs_proto_rawDesc = nil
	file_resources_users_jobs_proto_goTypes = nil
	file_resources_users_jobs_proto_depIdxs = nil
}
