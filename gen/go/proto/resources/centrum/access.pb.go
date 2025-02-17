// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v3.20.3
// source: resources/centrum/access.proto

package centrum

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	qualifications "github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
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

type UnitAccessLevel int32

const (
	UnitAccessLevel_UNIT_ACCESS_LEVEL_UNSPECIFIED UnitAccessLevel = 0
	UnitAccessLevel_UNIT_ACCESS_LEVEL_BLOCKED     UnitAccessLevel = 1
	UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN        UnitAccessLevel = 2
)

// Enum value maps for UnitAccessLevel.
var (
	UnitAccessLevel_name = map[int32]string{
		0: "UNIT_ACCESS_LEVEL_UNSPECIFIED",
		1: "UNIT_ACCESS_LEVEL_BLOCKED",
		2: "UNIT_ACCESS_LEVEL_JOIN",
	}
	UnitAccessLevel_value = map[string]int32{
		"UNIT_ACCESS_LEVEL_UNSPECIFIED": 0,
		"UNIT_ACCESS_LEVEL_BLOCKED":     1,
		"UNIT_ACCESS_LEVEL_JOIN":        2,
	}
)

func (x UnitAccessLevel) Enum() *UnitAccessLevel {
	p := new(UnitAccessLevel)
	*p = x
	return p
}

func (x UnitAccessLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UnitAccessLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_centrum_access_proto_enumTypes[0].Descriptor()
}

func (UnitAccessLevel) Type() protoreflect.EnumType {
	return &file_resources_centrum_access_proto_enumTypes[0]
}

func (x UnitAccessLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UnitAccessLevel.Descriptor instead.
func (UnitAccessLevel) EnumDescriptor() ([]byte, []int) {
	return file_resources_centrum_access_proto_rawDescGZIP(), []int{0}
}

type UnitAccess struct {
	state          protoimpl.MessageState     `protogen:"open.v1"`
	Jobs           []*UnitJobAccess           `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty" alias:"job_access"`                     // @gotags: alias:"job_access"
	Qualifications []*UnitQualificationAccess `protobuf:"bytes,3,rep,name=qualifications,proto3" json:"qualifications,omitempty" alias:"qualification_access"` // @gotags: alias:"qualification_access"
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *UnitAccess) Reset() {
	*x = UnitAccess{}
	mi := &file_resources_centrum_access_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnitAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnitAccess) ProtoMessage() {}

func (x *UnitAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_centrum_access_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnitAccess.ProtoReflect.Descriptor instead.
func (*UnitAccess) Descriptor() ([]byte, []int) {
	return file_resources_centrum_access_proto_rawDescGZIP(), []int{0}
}

func (x *UnitAccess) GetJobs() []*UnitJobAccess {
	if x != nil {
		return x.Jobs
	}
	return nil
}

func (x *UnitAccess) GetQualifications() []*UnitQualificationAccess {
	if x != nil {
		return x.Qualifications
	}
	return nil
}

type UnitJobAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	TargetId      uint64                 `protobuf:"varint,3,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty" alias:"calendar_id"` // @gotags: alias:"calendar_id"
	Job           string                 `protobuf:"bytes,4,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel      *string                `protobuf:"bytes,5,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty"`
	MinimumGrade  int32                  `protobuf:"varint,6,opt,name=minimum_grade,json=minimumGrade,proto3" json:"minimum_grade,omitempty"`
	JobGradeLabel *string                `protobuf:"bytes,7,opt,name=job_grade_label,json=jobGradeLabel,proto3,oneof" json:"job_grade_label,omitempty"`
	Access        UnitAccessLevel        `protobuf:"varint,8,opt,name=access,proto3,enum=resources.centrum.UnitAccessLevel" json:"access,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UnitJobAccess) Reset() {
	*x = UnitJobAccess{}
	mi := &file_resources_centrum_access_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnitJobAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnitJobAccess) ProtoMessage() {}

func (x *UnitJobAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_centrum_access_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnitJobAccess.ProtoReflect.Descriptor instead.
func (*UnitJobAccess) Descriptor() ([]byte, []int) {
	return file_resources_centrum_access_proto_rawDescGZIP(), []int{1}
}

func (x *UnitJobAccess) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UnitJobAccess) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *UnitJobAccess) GetTargetId() uint64 {
	if x != nil {
		return x.TargetId
	}
	return 0
}

func (x *UnitJobAccess) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *UnitJobAccess) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
	}
	return ""
}

func (x *UnitJobAccess) GetMinimumGrade() int32 {
	if x != nil {
		return x.MinimumGrade
	}
	return 0
}

func (x *UnitJobAccess) GetJobGradeLabel() string {
	if x != nil && x.JobGradeLabel != nil {
		return *x.JobGradeLabel
	}
	return ""
}

func (x *UnitJobAccess) GetAccess() UnitAccessLevel {
	if x != nil {
		return x.Access
	}
	return UnitAccessLevel_UNIT_ACCESS_LEVEL_UNSPECIFIED
}

type UnitUserAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UnitUserAccess) Reset() {
	*x = UnitUserAccess{}
	mi := &file_resources_centrum_access_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnitUserAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnitUserAccess) ProtoMessage() {}

func (x *UnitUserAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_centrum_access_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnitUserAccess.ProtoReflect.Descriptor instead.
func (*UnitUserAccess) Descriptor() ([]byte, []int) {
	return file_resources_centrum_access_proto_rawDescGZIP(), []int{2}
}

type UnitQualificationAccess struct {
	state           protoimpl.MessageState             `protogen:"open.v1"`
	Id              uint64                             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt       *timestamp.Timestamp               `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	TargetId        uint64                             `protobuf:"varint,3,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty" alias:"thread_id"` // @gotags: alias:"thread_id"
	QualificationId uint64                             `protobuf:"varint,4,opt,name=qualification_id,json=qualificationId,proto3" json:"qualification_id,omitempty"`
	Qualification   *qualifications.QualificationShort `protobuf:"bytes,5,opt,name=qualification,proto3,oneof" json:"qualification,omitempty"`
	Access          UnitAccessLevel                    `protobuf:"varint,6,opt,name=access,proto3,enum=resources.centrum.UnitAccessLevel" json:"access,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *UnitQualificationAccess) Reset() {
	*x = UnitQualificationAccess{}
	mi := &file_resources_centrum_access_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnitQualificationAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnitQualificationAccess) ProtoMessage() {}

func (x *UnitQualificationAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_centrum_access_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnitQualificationAccess.ProtoReflect.Descriptor instead.
func (*UnitQualificationAccess) Descriptor() ([]byte, []int) {
	return file_resources_centrum_access_proto_rawDescGZIP(), []int{3}
}

func (x *UnitQualificationAccess) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UnitQualificationAccess) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *UnitQualificationAccess) GetTargetId() uint64 {
	if x != nil {
		return x.TargetId
	}
	return 0
}

func (x *UnitQualificationAccess) GetQualificationId() uint64 {
	if x != nil {
		return x.QualificationId
	}
	return 0
}

func (x *UnitQualificationAccess) GetQualification() *qualifications.QualificationShort {
	if x != nil {
		return x.Qualification
	}
	return nil
}

func (x *UnitQualificationAccess) GetAccess() UnitAccessLevel {
	if x != nil {
		return x.Access
	}
	return UnitAccessLevel_UNIT_ACCESS_LEVEL_UNSPECIFIED
}

var File_resources_centrum_access_proto protoreflect.FileDescriptor

var file_resources_centrum_access_proto_rawDesc = string([]byte{
	0x0a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x65, 0x6e, 0x74,
	0x72, 0x75, 0x6d, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x11, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x65, 0x6e, 0x74,
	0x72, 0x75, 0x6d, 0x1a, 0x2d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x71,
	0x75, 0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x71, 0x75,
	0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xaa, 0x01, 0x0a, 0x0a, 0x55, 0x6e, 0x69, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x3e, 0x0a, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75,
	0x6d, 0x2e, 0x55, 0x6e, 0x69, 0x74, 0x4a, 0x6f, 0x62, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42,
	0x08, 0xfa, 0x42, 0x05, 0x92, 0x01, 0x02, 0x10, 0x14, 0x52, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x12,
	0x5c, 0x0a, 0x0e, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x2e, 0x55, 0x6e, 0x69, 0x74,
	0x51, 0x75, 0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x92, 0x01, 0x02, 0x10, 0x14, 0x52, 0x0e, 0x71,
	0x75, 0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xa1, 0x03,
	0x0a, 0x0d, 0x55, 0x6e, 0x69, 0x74, 0x4a, 0x6f, 0x62, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64,
	0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x29, 0x0a, 0x09, 0x6a,
	0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x01, 0x52, 0x08, 0x6a, 0x6f, 0x62, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x2c, 0x0a, 0x0d, 0x6d, 0x69, 0x6e, 0x69, 0x6d, 0x75,
	0x6d, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x1a, 0x02, 0x28, 0x00, 0x52, 0x0c, 0x6d, 0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x47,
	0x72, 0x61, 0x64, 0x65, 0x12, 0x34, 0x0a, 0x0f, 0x6a, 0x6f, 0x62, 0x5f, 0x67, 0x72, 0x61, 0x64,
	0x65, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x02, 0x52, 0x0d, 0x6a, 0x6f, 0x62, 0x47, 0x72, 0x61,
	0x64, 0x65, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x44, 0x0a, 0x06, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x2e, 0x55,
	0x6e, 0x69, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42,
	0x0c, 0x0a, 0x0a, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x42, 0x12, 0x0a,
	0x10, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x22, 0x10, 0x0a, 0x0e, 0x55, 0x6e, 0x69, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x22, 0xf5, 0x02, 0x0a, 0x17, 0x55, 0x6e, 0x69, 0x74, 0x51, 0x75, 0x61, 0x6c,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64,
	0x12, 0x29, 0x0a, 0x10, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x71, 0x75, 0x61, 0x6c,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x57, 0x0a, 0x0d, 0x71,
	0x75, 0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x71,
	0x75, 0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x51, 0x75,
	0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x48, 0x01, 0x52, 0x0d, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x88, 0x01, 0x01, 0x12, 0x44, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x2e, 0x55, 0x6e, 0x69, 0x74, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02,
	0x10, 0x01, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x71, 0x75,
	0x61, 0x6c, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x6f, 0x0a, 0x0f, 0x55,
	0x6e, 0x69, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x21,
	0x0a, 0x1d, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45,
	0x56, 0x45, 0x4c, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x1d, 0x0a, 0x19, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53,
	0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x45, 0x44, 0x10, 0x01,
	0x12, 0x1a, 0x0a, 0x16, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f,
	0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x4a, 0x4f, 0x49, 0x4e, 0x10, 0x02, 0x42, 0x47, 0x5a, 0x45,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e,
	0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x75, 0x6d, 0x3b, 0x63, 0x65,
	0x6e, 0x74, 0x72, 0x75, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_centrum_access_proto_rawDescOnce sync.Once
	file_resources_centrum_access_proto_rawDescData []byte
)

func file_resources_centrum_access_proto_rawDescGZIP() []byte {
	file_resources_centrum_access_proto_rawDescOnce.Do(func() {
		file_resources_centrum_access_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_centrum_access_proto_rawDesc), len(file_resources_centrum_access_proto_rawDesc)))
	})
	return file_resources_centrum_access_proto_rawDescData
}

var file_resources_centrum_access_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_centrum_access_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_resources_centrum_access_proto_goTypes = []any{
	(UnitAccessLevel)(0),                      // 0: resources.centrum.UnitAccessLevel
	(*UnitAccess)(nil),                        // 1: resources.centrum.UnitAccess
	(*UnitJobAccess)(nil),                     // 2: resources.centrum.UnitJobAccess
	(*UnitUserAccess)(nil),                    // 3: resources.centrum.UnitUserAccess
	(*UnitQualificationAccess)(nil),           // 4: resources.centrum.UnitQualificationAccess
	(*timestamp.Timestamp)(nil),               // 5: resources.timestamp.Timestamp
	(*qualifications.QualificationShort)(nil), // 6: resources.qualifications.QualificationShort
}
var file_resources_centrum_access_proto_depIdxs = []int32{
	2, // 0: resources.centrum.UnitAccess.jobs:type_name -> resources.centrum.UnitJobAccess
	4, // 1: resources.centrum.UnitAccess.qualifications:type_name -> resources.centrum.UnitQualificationAccess
	5, // 2: resources.centrum.UnitJobAccess.created_at:type_name -> resources.timestamp.Timestamp
	0, // 3: resources.centrum.UnitJobAccess.access:type_name -> resources.centrum.UnitAccessLevel
	5, // 4: resources.centrum.UnitQualificationAccess.created_at:type_name -> resources.timestamp.Timestamp
	6, // 5: resources.centrum.UnitQualificationAccess.qualification:type_name -> resources.qualifications.QualificationShort
	0, // 6: resources.centrum.UnitQualificationAccess.access:type_name -> resources.centrum.UnitAccessLevel
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_resources_centrum_access_proto_init() }
func file_resources_centrum_access_proto_init() {
	if File_resources_centrum_access_proto != nil {
		return
	}
	file_resources_centrum_access_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_centrum_access_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_centrum_access_proto_rawDesc), len(file_resources_centrum_access_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_centrum_access_proto_goTypes,
		DependencyIndexes: file_resources_centrum_access_proto_depIdxs,
		EnumInfos:         file_resources_centrum_access_proto_enumTypes,
		MessageInfos:      file_resources_centrum_access_proto_msgTypes,
	}.Build()
	File_resources_centrum_access_proto = out.File
	file_resources_centrum_access_proto_goTypes = nil
	file_resources_centrum_access_proto_depIdxs = nil
}
