// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v3.20.3
// source: resources/centrum/settings.proto

package centrum

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type CentrumMode int32

const (
	CentrumMode_CENTRUM_MODE_UNSPECIFIED      CentrumMode = 0
	CentrumMode_CENTRUM_MODE_MANUAL           CentrumMode = 1
	CentrumMode_CENTRUM_MODE_CENTRAL_COMMAND  CentrumMode = 2
	CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN CentrumMode = 3
	CentrumMode_CENTRUM_MODE_SIMPLIFIED       CentrumMode = 4
)

// Enum value maps for CentrumMode.
var (
	CentrumMode_name = map[int32]string{
		0: "CENTRUM_MODE_UNSPECIFIED",
		1: "CENTRUM_MODE_MANUAL",
		2: "CENTRUM_MODE_CENTRAL_COMMAND",
		3: "CENTRUM_MODE_AUTO_ROUND_ROBIN",
		4: "CENTRUM_MODE_SIMPLIFIED",
	}
	CentrumMode_value = map[string]int32{
		"CENTRUM_MODE_UNSPECIFIED":      0,
		"CENTRUM_MODE_MANUAL":           1,
		"CENTRUM_MODE_CENTRAL_COMMAND":  2,
		"CENTRUM_MODE_AUTO_ROUND_ROBIN": 3,
		"CENTRUM_MODE_SIMPLIFIED":       4,
	}
)

func (x CentrumMode) Enum() *CentrumMode {
	p := new(CentrumMode)
	*p = x
	return p
}

func (x CentrumMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CentrumMode) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_centrum_settings_proto_enumTypes[0].Descriptor()
}

func (CentrumMode) Type() protoreflect.EnumType {
	return &file_resources_centrum_settings_proto_enumTypes[0]
}

func (x CentrumMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CentrumMode.Descriptor instead.
func (CentrumMode) EnumDescriptor() ([]byte, []int) {
	return file_resources_centrum_settings_proto_rawDescGZIP(), []int{0}
}

type Settings struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	Job              string                 `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
	Enabled          bool                   `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Mode             CentrumMode            `protobuf:"varint,3,opt,name=mode,proto3,enum=resources.centrum.CentrumMode" json:"mode,omitempty"`
	FallbackMode     CentrumMode            `protobuf:"varint,4,opt,name=fallback_mode,json=fallbackMode,proto3,enum=resources.centrum.CentrumMode" json:"fallback_mode,omitempty"`
	PredefinedStatus *PredefinedStatus      `protobuf:"bytes,5,opt,name=predefined_status,json=predefinedStatus,proto3,oneof" json:"predefined_status,omitempty"`
	Timings          *Timings               `protobuf:"bytes,6,opt,name=timings,proto3" json:"timings,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *Settings) Reset() {
	*x = Settings{}
	mi := &file_resources_centrum_settings_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Settings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Settings) ProtoMessage() {}

func (x *Settings) ProtoReflect() protoreflect.Message {
	mi := &file_resources_centrum_settings_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Settings.ProtoReflect.Descriptor instead.
func (*Settings) Descriptor() ([]byte, []int) {
	return file_resources_centrum_settings_proto_rawDescGZIP(), []int{0}
}

func (x *Settings) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *Settings) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *Settings) GetMode() CentrumMode {
	if x != nil {
		return x.Mode
	}
	return CentrumMode_CENTRUM_MODE_UNSPECIFIED
}

func (x *Settings) GetFallbackMode() CentrumMode {
	if x != nil {
		return x.FallbackMode
	}
	return CentrumMode_CENTRUM_MODE_UNSPECIFIED
}

func (x *Settings) GetPredefinedStatus() *PredefinedStatus {
	if x != nil {
		return x.PredefinedStatus
	}
	return nil
}

func (x *Settings) GetTimings() *Timings {
	if x != nil {
		return x.Timings
	}
	return nil
}

type PredefinedStatus struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// @sanitize: method=StripTags
	UnitStatus []string `protobuf:"bytes,1,rep,name=unit_status,json=unitStatus,proto3" json:"unit_status,omitempty"`
	// @sanitize: method=StripTags
	DispatchStatus []string `protobuf:"bytes,2,rep,name=dispatch_status,json=dispatchStatus,proto3" json:"dispatch_status,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *PredefinedStatus) Reset() {
	*x = PredefinedStatus{}
	mi := &file_resources_centrum_settings_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PredefinedStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredefinedStatus) ProtoMessage() {}

func (x *PredefinedStatus) ProtoReflect() protoreflect.Message {
	mi := &file_resources_centrum_settings_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredefinedStatus.ProtoReflect.Descriptor instead.
func (*PredefinedStatus) Descriptor() ([]byte, []int) {
	return file_resources_centrum_settings_proto_rawDescGZIP(), []int{1}
}

func (x *PredefinedStatus) GetUnitStatus() []string {
	if x != nil {
		return x.UnitStatus
	}
	return nil
}

func (x *PredefinedStatus) GetDispatchStatus() []string {
	if x != nil {
		return x.DispatchStatus
	}
	return nil
}

type Timings struct {
	state                      protoimpl.MessageState `protogen:"open.v1"`
	DispatchMaxWait            int64                  `protobuf:"varint,1,opt,name=dispatch_max_wait,json=dispatchMaxWait,proto3" json:"dispatch_max_wait,omitempty"`
	RequireUnit                bool                   `protobuf:"varint,2,opt,name=require_unit,json=requireUnit,proto3" json:"require_unit,omitempty"`
	RequireUnitReminderSeconds int64                  `protobuf:"varint,3,opt,name=require_unit_reminder_seconds,json=requireUnitReminderSeconds,proto3" json:"require_unit_reminder_seconds,omitempty"`
	unknownFields              protoimpl.UnknownFields
	sizeCache                  protoimpl.SizeCache
}

func (x *Timings) Reset() {
	*x = Timings{}
	mi := &file_resources_centrum_settings_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Timings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Timings) ProtoMessage() {}

func (x *Timings) ProtoReflect() protoreflect.Message {
	mi := &file_resources_centrum_settings_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Timings.ProtoReflect.Descriptor instead.
func (*Timings) Descriptor() ([]byte, []int) {
	return file_resources_centrum_settings_proto_rawDescGZIP(), []int{2}
}

func (x *Timings) GetDispatchMaxWait() int64 {
	if x != nil {
		return x.DispatchMaxWait
	}
	return 0
}

func (x *Timings) GetRequireUnit() bool {
	if x != nil {
		return x.RequireUnit
	}
	return false
}

func (x *Timings) GetRequireUnitReminderSeconds() int64 {
	if x != nil {
		return x.RequireUnitReminderSeconds
	}
	return 0
}

var File_resources_centrum_settings_proto protoreflect.FileDescriptor

var file_resources_centrum_settings_proto_rawDesc = []byte{
	0x0a, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x65, 0x6e, 0x74,
	0x72, 0x75, 0x6d, 0x2f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x11, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x65,
	0x6e, 0x74, 0x72, 0x75, 0x6d, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xef,
	0x02, 0x0a, 0x08, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x19, 0x0a, 0x03, 0x6a,
	0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18,
	0x14, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64,
	0x12, 0x3c, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72,
	0x75, 0x6d, 0x2e, 0x43, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x4d, 0x6f, 0x64, 0x65, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x12, 0x4d,
	0x0a, 0x0d, 0x66, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x2e, 0x43, 0x65, 0x6e, 0x74, 0x72, 0x75,
	0x6d, 0x4d, 0x6f, 0x64, 0x65, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52,
	0x0c, 0x66, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x55, 0x0a,
	0x11, 0x70, 0x72, 0x65, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x2e, 0x50, 0x72, 0x65,
	0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x48, 0x00, 0x52,
	0x10, 0x70, 0x72, 0x65, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x88, 0x01, 0x01, 0x12, 0x34, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x69, 0x6e, 0x67, 0x73, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x2e, 0x54, 0x69, 0x6d, 0x69, 0x6e, 0x67,
	0x73, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x69, 0x6e, 0x67, 0x73, 0x42, 0x14, 0x0a, 0x12, 0x5f, 0x70,
	0x72, 0x65, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x22, 0x7c, 0x0a, 0x10, 0x50, 0x72, 0x65, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x0b, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x42, 0x0e, 0xfa, 0x42, 0x0b, 0x92, 0x01,
	0x08, 0x10, 0x05, 0x22, 0x04, 0x72, 0x02, 0x18, 0x40, 0x52, 0x0a, 0x75, 0x6e, 0x69, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x37, 0x0a, 0x0f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x42, 0x0e,
	0xfa, 0x42, 0x0b, 0x92, 0x01, 0x08, 0x10, 0x05, 0x22, 0x04, 0x72, 0x02, 0x18, 0x40, 0x52, 0x0e,
	0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xb3,
	0x01, 0x0a, 0x07, 0x54, 0x69, 0x6d, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x36, 0x0a, 0x11, 0x64, 0x69,
	0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x77, 0x61, 0x69, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x22, 0x05, 0x10, 0xf0, 0x2e, 0x20,
	0x1e, 0x52, 0x0f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x4d, 0x61, 0x78, 0x57, 0x61,
	0x69, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x5f, 0x75, 0x6e,
	0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72,
	0x65, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x4d, 0x0a, 0x1d, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x5f, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x72, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x73,
	0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x42, 0x0a, 0xfa, 0x42,
	0x07, 0x22, 0x05, 0x10, 0xf0, 0x2e, 0x20, 0x1e, 0x52, 0x1a, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72,
	0x65, 0x55, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x53, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x73, 0x2a, 0xa6, 0x01, 0x0a, 0x0b, 0x43, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d,
	0x4d, 0x6f, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x18, 0x43, 0x45, 0x4e, 0x54, 0x52, 0x55, 0x4d, 0x5f,
	0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x45, 0x4e, 0x54, 0x52, 0x55, 0x4d, 0x5f, 0x4d, 0x4f,
	0x44, 0x45, 0x5f, 0x4d, 0x41, 0x4e, 0x55, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x20, 0x0a, 0x1c, 0x43,
	0x45, 0x4e, 0x54, 0x52, 0x55, 0x4d, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x43, 0x45, 0x4e, 0x54,
	0x52, 0x41, 0x4c, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x10, 0x02, 0x12, 0x21, 0x0a,
	0x1d, 0x43, 0x45, 0x4e, 0x54, 0x52, 0x55, 0x4d, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x41, 0x55,
	0x54, 0x4f, 0x5f, 0x52, 0x4f, 0x55, 0x4e, 0x44, 0x5f, 0x52, 0x4f, 0x42, 0x49, 0x4e, 0x10, 0x03,
	0x12, 0x1b, 0x0a, 0x17, 0x43, 0x45, 0x4e, 0x54, 0x52, 0x55, 0x4d, 0x5f, 0x4d, 0x4f, 0x44, 0x45,
	0x5f, 0x53, 0x49, 0x4d, 0x50, 0x4c, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x04, 0x42, 0x47, 0x5a,
	0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65,
	0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x3b, 0x63,
	0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_centrum_settings_proto_rawDescOnce sync.Once
	file_resources_centrum_settings_proto_rawDescData = file_resources_centrum_settings_proto_rawDesc
)

func file_resources_centrum_settings_proto_rawDescGZIP() []byte {
	file_resources_centrum_settings_proto_rawDescOnce.Do(func() {
		file_resources_centrum_settings_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_centrum_settings_proto_rawDescData)
	})
	return file_resources_centrum_settings_proto_rawDescData
}

var file_resources_centrum_settings_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_centrum_settings_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_resources_centrum_settings_proto_goTypes = []any{
	(CentrumMode)(0),         // 0: resources.centrum.CentrumMode
	(*Settings)(nil),         // 1: resources.centrum.Settings
	(*PredefinedStatus)(nil), // 2: resources.centrum.PredefinedStatus
	(*Timings)(nil),          // 3: resources.centrum.Timings
}
var file_resources_centrum_settings_proto_depIdxs = []int32{
	0, // 0: resources.centrum.Settings.mode:type_name -> resources.centrum.CentrumMode
	0, // 1: resources.centrum.Settings.fallback_mode:type_name -> resources.centrum.CentrumMode
	2, // 2: resources.centrum.Settings.predefined_status:type_name -> resources.centrum.PredefinedStatus
	3, // 3: resources.centrum.Settings.timings:type_name -> resources.centrum.Timings
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_resources_centrum_settings_proto_init() }
func file_resources_centrum_settings_proto_init() {
	if File_resources_centrum_settings_proto != nil {
		return
	}
	file_resources_centrum_settings_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_centrum_settings_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_centrum_settings_proto_goTypes,
		DependencyIndexes: file_resources_centrum_settings_proto_depIdxs,
		EnumInfos:         file_resources_centrum_settings_proto_enumTypes,
		MessageInfos:      file_resources_centrum_settings_proto_msgTypes,
	}.Build()
	File_resources_centrum_settings_proto = out.File
	file_resources_centrum_settings_proto_rawDesc = nil
	file_resources_centrum_settings_proto_goTypes = nil
	file_resources_centrum_settings_proto_depIdxs = nil
}
