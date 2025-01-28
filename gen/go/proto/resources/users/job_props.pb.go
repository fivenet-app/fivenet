// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v3.20.3
// source: resources/users/job_props.proto

package users

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	filestore "github.com/fivenet-app/fivenet/gen/go/proto/resources/filestore"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
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

type JobProps struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	Job                 string                 `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel            *string                `protobuf:"bytes,2,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty"`
	Theme               string                 `protobuf:"bytes,3,opt,name=theme,proto3" json:"theme,omitempty"`
	LivemapMarkerColor  string                 `protobuf:"bytes,4,opt,name=livemap_marker_color,json=livemapMarkerColor,proto3" json:"livemap_marker_color,omitempty"`
	QuickButtons        *QuickButtons          `protobuf:"bytes,5,opt,name=quick_buttons,json=quickButtons,proto3" json:"quick_buttons,omitempty"`
	RadioFrequency      *string                `protobuf:"bytes,6,opt,name=radio_frequency,json=radioFrequency,proto3,oneof" json:"radio_frequency,omitempty"`
	DiscordGuildId      *string                `protobuf:"bytes,7,opt,name=discord_guild_id,json=discordGuildId,proto3,oneof" json:"discord_guild_id,omitempty"`
	DiscordLastSync     *timestamp.Timestamp   `protobuf:"bytes,8,opt,name=discord_last_sync,json=discordLastSync,proto3,oneof" json:"discord_last_sync,omitempty"`
	DiscordSyncSettings *DiscordSyncSettings   `protobuf:"bytes,9,opt,name=discord_sync_settings,json=discordSyncSettings,proto3" json:"discord_sync_settings,omitempty"`
	DiscordSyncChanges  *DiscordSyncChanges    `protobuf:"bytes,10,opt,name=discord_sync_changes,json=discordSyncChanges,proto3,oneof" json:"discord_sync_changes,omitempty"`
	Motd                *string                `protobuf:"bytes,11,opt,name=motd,proto3,oneof" json:"motd,omitempty"`
	LogoUrl             *filestore.File        `protobuf:"bytes,12,opt,name=logo_url,json=logoUrl,proto3,oneof" json:"logo_url,omitempty"`
	Settings            *JobSettings           `protobuf:"bytes,13,opt,name=settings,proto3" json:"settings,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *JobProps) Reset() {
	*x = JobProps{}
	mi := &file_resources_users_job_props_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobProps) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobProps) ProtoMessage() {}

func (x *JobProps) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_job_props_proto_msgTypes[0]
	if x != nil {
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
	return file_resources_users_job_props_proto_rawDescGZIP(), []int{0}
}

func (x *JobProps) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *JobProps) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
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

func (x *JobProps) GetRadioFrequency() string {
	if x != nil && x.RadioFrequency != nil {
		return *x.RadioFrequency
	}
	return ""
}

func (x *JobProps) GetDiscordGuildId() string {
	if x != nil && x.DiscordGuildId != nil {
		return *x.DiscordGuildId
	}
	return ""
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

func (x *JobProps) GetDiscordSyncChanges() *DiscordSyncChanges {
	if x != nil {
		return x.DiscordSyncChanges
	}
	return nil
}

func (x *JobProps) GetMotd() string {
	if x != nil && x.Motd != nil {
		return *x.Motd
	}
	return ""
}

func (x *JobProps) GetLogoUrl() *filestore.File {
	if x != nil {
		return x.LogoUrl
	}
	return nil
}

func (x *JobProps) GetSettings() *JobSettings {
	if x != nil {
		return x.Settings
	}
	return nil
}

type QuickButtons struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	PenaltyCalculator bool                   `protobuf:"varint,1,opt,name=penalty_calculator,json=penaltyCalculator,proto3" json:"penalty_calculator,omitempty"`
	BodyCheckup       bool                   `protobuf:"varint,2,opt,name=body_checkup,json=bodyCheckup,proto3" json:"body_checkup,omitempty"`
	MathCalculator    bool                   `protobuf:"varint,3,opt,name=math_calculator,json=mathCalculator,proto3" json:"math_calculator,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *QuickButtons) Reset() {
	*x = QuickButtons{}
	mi := &file_resources_users_job_props_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QuickButtons) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuickButtons) ProtoMessage() {}

func (x *QuickButtons) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_job_props_proto_msgTypes[1]
	if x != nil {
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
	return file_resources_users_job_props_proto_rawDescGZIP(), []int{1}
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

func (x *QuickButtons) GetMathCalculator() bool {
	if x != nil {
		return x.MathCalculator
	}
	return false
}

var File_resources_users_job_props_proto protoreflect.FileDescriptor

var file_resources_users_job_props_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x1a, 0x22, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x66, 0x69, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x83, 0x07, 0x0a, 0x08, 0x4a, 0x6f, 0x62, 0x50, 0x72, 0x6f, 0x70,
	0x73, 0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x29, 0x0a, 0x09,
	0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x00, 0x52, 0x08, 0x6a, 0x6f, 0x62, 0x4c,
	0x61, 0x62, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x05, 0x74, 0x68, 0x65, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52,
	0x05, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x4d, 0x0a, 0x14, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61,
	0x70, 0x5f, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x1b, 0xfa, 0x42, 0x18, 0x72, 0x16, 0x32, 0x11, 0x5e, 0x23, 0x5b,
	0x41, 0x2d, 0x46, 0x61, 0x2d, 0x66, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x36, 0x7d, 0x24, 0x98, 0x01,
	0x07, 0x52, 0x12, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72,
	0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x42, 0x0a, 0x0d, 0x71, 0x75, 0x69, 0x63, 0x6b, 0x5f, 0x62,
	0x75, 0x74, 0x74, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x51,
	0x75, 0x69, 0x63, 0x6b, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x73, 0x52, 0x0c, 0x71, 0x75, 0x69,
	0x63, 0x6b, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x73, 0x12, 0x35, 0x0a, 0x0f, 0x72, 0x61, 0x64,
	0x69, 0x6f, 0x5f, 0x66, 0x72, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x18, 0x48, 0x01, 0x52, 0x0e, 0x72,
	0x61, 0x64, 0x69, 0x6f, 0x46, 0x72, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x79, 0x88, 0x01, 0x01,
	0x12, 0x2d, 0x0a, 0x10, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x67, 0x75, 0x69, 0x6c,
	0x64, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x0e, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x72, 0x64, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12,
	0x4f, 0x0a, 0x11, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f,
	0x73, 0x79, 0x6e, 0x63, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x03, 0x52, 0x0f, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x72, 0x64, 0x4c, 0x61, 0x73, 0x74, 0x53, 0x79, 0x6e, 0x63, 0x88, 0x01, 0x01,
	0x12, 0x58, 0x0a, 0x15, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x79, 0x6e, 0x63,
	0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x24, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x65, 0x74,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x13, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x79,
	0x6e, 0x63, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x5a, 0x0a, 0x14, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f,
	0x72, 0x64, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x48, 0x04, 0x52,
	0x12, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x73, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x04, 0x6d, 0x6f, 0x74, 0x64, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x08, 0x48, 0x05,
	0x52, 0x04, 0x6d, 0x6f, 0x74, 0x64, 0x88, 0x01, 0x01, 0x12, 0x39, 0x0a, 0x08, 0x6c, 0x6f, 0x67,
	0x6f, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x06, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x6f, 0x55, 0x72,
	0x6c, 0x88, 0x01, 0x01, 0x12, 0x38, 0x0a, 0x08, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x53, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x73, 0x52, 0x08, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x42, 0x0c,
	0x0a, 0x0a, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x42, 0x12, 0x0a, 0x10,
	0x5f, 0x72, 0x61, 0x64, 0x69, 0x6f, 0x5f, 0x66, 0x72, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x79,
	0x42, 0x13, 0x0a, 0x11, 0x5f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x67, 0x75, 0x69,
	0x6c, 0x64, 0x5f, 0x69, 0x64, 0x42, 0x14, 0x0a, 0x12, 0x5f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72,
	0x64, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x79, 0x6e, 0x63, 0x42, 0x17, 0x0a, 0x15, 0x5f,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x73, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6d, 0x6f, 0x74, 0x64, 0x42, 0x0b, 0x0a,
	0x09, 0x5f, 0x6c, 0x6f, 0x67, 0x6f, 0x5f, 0x75, 0x72, 0x6c, 0x22, 0x89, 0x01, 0x0a, 0x0c, 0x51,
	0x75, 0x69, 0x63, 0x6b, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x73, 0x12, 0x2d, 0x0a, 0x12, 0x70,
	0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79, 0x5f, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x70, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79,
	0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6f,
	0x64, 0x79, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0b, 0x62, 0x6f, 0x64, 0x79, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x75, 0x70, 0x12, 0x27, 0x0a,
	0x0f, 0x6d, 0x61, 0x74, 0x68, 0x5f, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x6d, 0x61, 0x74, 0x68, 0x43, 0x61, 0x6c, 0x63,
	0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70,
	0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x3b, 0x75, 0x73, 0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_resources_users_job_props_proto_rawDescOnce sync.Once
	file_resources_users_job_props_proto_rawDescData = file_resources_users_job_props_proto_rawDesc
)

func file_resources_users_job_props_proto_rawDescGZIP() []byte {
	file_resources_users_job_props_proto_rawDescOnce.Do(func() {
		file_resources_users_job_props_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_users_job_props_proto_rawDescData)
	})
	return file_resources_users_job_props_proto_rawDescData
}

var file_resources_users_job_props_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_resources_users_job_props_proto_goTypes = []any{
	(*JobProps)(nil),            // 0: resources.users.JobProps
	(*QuickButtons)(nil),        // 1: resources.users.QuickButtons
	(*timestamp.Timestamp)(nil), // 2: resources.timestamp.Timestamp
	(*DiscordSyncSettings)(nil), // 3: resources.users.DiscordSyncSettings
	(*DiscordSyncChanges)(nil),  // 4: resources.users.DiscordSyncChanges
	(*filestore.File)(nil),      // 5: resources.filestore.File
	(*JobSettings)(nil),         // 6: resources.users.JobSettings
}
var file_resources_users_job_props_proto_depIdxs = []int32{
	1, // 0: resources.users.JobProps.quick_buttons:type_name -> resources.users.QuickButtons
	2, // 1: resources.users.JobProps.discord_last_sync:type_name -> resources.timestamp.Timestamp
	3, // 2: resources.users.JobProps.discord_sync_settings:type_name -> resources.users.DiscordSyncSettings
	4, // 3: resources.users.JobProps.discord_sync_changes:type_name -> resources.users.DiscordSyncChanges
	5, // 4: resources.users.JobProps.logo_url:type_name -> resources.filestore.File
	6, // 5: resources.users.JobProps.settings:type_name -> resources.users.JobSettings
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_resources_users_job_props_proto_init() }
func file_resources_users_job_props_proto_init() {
	if File_resources_users_job_props_proto != nil {
		return
	}
	file_resources_users_job_settings_proto_init()
	file_resources_users_job_props_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_users_job_props_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_users_job_props_proto_goTypes,
		DependencyIndexes: file_resources_users_job_props_proto_depIdxs,
		MessageInfos:      file_resources_users_job_props_proto_msgTypes,
	}.Build()
	File_resources_users_job_props_proto = out.File
	file_resources_users_job_props_proto_rawDesc = nil
	file_resources_users_job_props_proto_goTypes = nil
	file_resources_users_job_props_proto_depIdxs = nil
}
