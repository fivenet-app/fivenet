// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v3.20.3
// source: resources/jobs/conduct.proto

package jobs

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type ConductType int32

const (
	ConductType_CONDUCT_TYPE_UNSPECIFIED ConductType = 0
	ConductType_CONDUCT_TYPE_NEUTRAL     ConductType = 1
	ConductType_CONDUCT_TYPE_POSITIVE    ConductType = 2
	ConductType_CONDUCT_TYPE_NEGATIVE    ConductType = 3
	ConductType_CONDUCT_TYPE_WARNING     ConductType = 4
	ConductType_CONDUCT_TYPE_SUSPENSION  ConductType = 5
	ConductType_CONDUCT_TYPE_NOTE        ConductType = 6
)

// Enum value maps for ConductType.
var (
	ConductType_name = map[int32]string{
		0: "CONDUCT_TYPE_UNSPECIFIED",
		1: "CONDUCT_TYPE_NEUTRAL",
		2: "CONDUCT_TYPE_POSITIVE",
		3: "CONDUCT_TYPE_NEGATIVE",
		4: "CONDUCT_TYPE_WARNING",
		5: "CONDUCT_TYPE_SUSPENSION",
		6: "CONDUCT_TYPE_NOTE",
	}
	ConductType_value = map[string]int32{
		"CONDUCT_TYPE_UNSPECIFIED": 0,
		"CONDUCT_TYPE_NEUTRAL":     1,
		"CONDUCT_TYPE_POSITIVE":    2,
		"CONDUCT_TYPE_NEGATIVE":    3,
		"CONDUCT_TYPE_WARNING":     4,
		"CONDUCT_TYPE_SUSPENSION":  5,
		"CONDUCT_TYPE_NOTE":        6,
	}
)

func (x ConductType) Enum() *ConductType {
	p := new(ConductType)
	*p = x
	return p
}

func (x ConductType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConductType) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_jobs_conduct_proto_enumTypes[0].Descriptor()
}

func (ConductType) Type() protoreflect.EnumType {
	return &file_resources_jobs_conduct_proto_enumTypes[0]
}

func (x ConductType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConductType.Descriptor instead.
func (ConductType) EnumDescriptor() ([]byte, []int) {
	return file_resources_jobs_conduct_proto_rawDescGZIP(), []int{0}
}

type ConductEntry struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	Job       string                 `protobuf:"bytes,4,opt,name=job,proto3" json:"job,omitempty"`
	Type      ConductType            `protobuf:"varint,5,opt,name=type,proto3,enum=resources.jobs.ConductType" json:"type,omitempty"`
	// @sanitize
	Message       string               `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
	ExpiresAt     *timestamp.Timestamp `protobuf:"bytes,7,opt,name=expires_at,json=expiresAt,proto3,oneof" json:"expires_at,omitempty"`
	TargetUserId  int32                `protobuf:"varint,8,opt,name=target_user_id,json=targetUserId,proto3" json:"target_user_id,omitempty"`
	TargetUser    *Colleague           `protobuf:"bytes,9,opt,name=target_user,json=targetUser,proto3,oneof" json:"target_user,omitempty" alias:"target_user"` // @gotags: alias:"target_user"
	CreatorId     int32                `protobuf:"varint,10,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
	Creator       *Colleague           `protobuf:"bytes,11,opt,name=creator,proto3,oneof" json:"creator,omitempty" alias:"creator"` // @gotags: alias:"creator"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConductEntry) Reset() {
	*x = ConductEntry{}
	mi := &file_resources_jobs_conduct_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConductEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConductEntry) ProtoMessage() {}

func (x *ConductEntry) ProtoReflect() protoreflect.Message {
	mi := &file_resources_jobs_conduct_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConductEntry.ProtoReflect.Descriptor instead.
func (*ConductEntry) Descriptor() ([]byte, []int) {
	return file_resources_jobs_conduct_proto_rawDescGZIP(), []int{0}
}

func (x *ConductEntry) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ConductEntry) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *ConductEntry) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *ConductEntry) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *ConductEntry) GetType() ConductType {
	if x != nil {
		return x.Type
	}
	return ConductType_CONDUCT_TYPE_UNSPECIFIED
}

func (x *ConductEntry) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ConductEntry) GetExpiresAt() *timestamp.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

func (x *ConductEntry) GetTargetUserId() int32 {
	if x != nil {
		return x.TargetUserId
	}
	return 0
}

func (x *ConductEntry) GetTargetUser() *Colleague {
	if x != nil {
		return x.TargetUser
	}
	return nil
}

func (x *ConductEntry) GetCreatorId() int32 {
	if x != nil {
		return x.CreatorId
	}
	return 0
}

func (x *ConductEntry) GetCreator() *Colleague {
	if x != nil {
		return x.Creator
	}
	return nil
}

var File_resources_jobs_conduct_proto protoreflect.FileDescriptor

var file_resources_jobs_conduct_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73,
	0x2f, 0x63, 0x6f, 0x6e, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x1a, 0x1f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2f, 0x63,
	0x6f, 0x6c, 0x6c, 0x65, 0x61, 0x67, 0x75, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x81, 0x05,
	0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x64, 0x75, 0x63, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52, 0x03, 0x6a, 0x6f,
	0x62, 0x12, 0x39, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1b, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73,
	0x2e, 0x43, 0x6f, 0x6e, 0x64, 0x75, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x42, 0x08, 0xfa, 0x42,
	0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x24, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa,
	0x42, 0x07, 0x72, 0x05, 0x10, 0x03, 0x18, 0x80, 0x10, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x42, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x02, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65,
	0x73, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x3f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x43, 0x6f, 0x6c, 0x6c,
	0x65, 0x61, 0x67, 0x75, 0x65, 0x48, 0x03, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a,
	0x02, 0x20, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x38,
	0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73,
	0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x61, 0x67, 0x75, 0x65, 0x48, 0x04, 0x52, 0x07, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x65, 0x73, 0x5f, 0x61, 0x74, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x2a, 0xc9, 0x01, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x64, 0x75, 0x63, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1c, 0x0a, 0x18, 0x43, 0x4f, 0x4e, 0x44, 0x55, 0x43, 0x54, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x18, 0x0a, 0x14, 0x43, 0x4f, 0x4e, 0x44, 0x55, 0x43, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x4e, 0x45, 0x55, 0x54, 0x52, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x43, 0x4f, 0x4e,
	0x44, 0x55, 0x43, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x4f, 0x53, 0x49, 0x54, 0x49,
	0x56, 0x45, 0x10, 0x02, 0x12, 0x19, 0x0a, 0x15, 0x43, 0x4f, 0x4e, 0x44, 0x55, 0x43, 0x54, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x45, 0x47, 0x41, 0x54, 0x49, 0x56, 0x45, 0x10, 0x03, 0x12,
	0x18, 0x0a, 0x14, 0x43, 0x4f, 0x4e, 0x44, 0x55, 0x43, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x12, 0x1b, 0x0a, 0x17, 0x43, 0x4f, 0x4e,
	0x44, 0x55, 0x43, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x55, 0x53, 0x50, 0x45, 0x4e,
	0x53, 0x49, 0x4f, 0x4e, 0x10, 0x05, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x4f, 0x4e, 0x44, 0x55, 0x43,
	0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x54, 0x45, 0x10, 0x06, 0x42, 0x41, 0x5a,
	0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65,
	0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x3b, 0x6a, 0x6f, 0x62, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_jobs_conduct_proto_rawDescOnce sync.Once
	file_resources_jobs_conduct_proto_rawDescData = file_resources_jobs_conduct_proto_rawDesc
)

func file_resources_jobs_conduct_proto_rawDescGZIP() []byte {
	file_resources_jobs_conduct_proto_rawDescOnce.Do(func() {
		file_resources_jobs_conduct_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_jobs_conduct_proto_rawDescData)
	})
	return file_resources_jobs_conduct_proto_rawDescData
}

var file_resources_jobs_conduct_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_jobs_conduct_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_jobs_conduct_proto_goTypes = []any{
	(ConductType)(0),            // 0: resources.jobs.ConductType
	(*ConductEntry)(nil),        // 1: resources.jobs.ConductEntry
	(*timestamp.Timestamp)(nil), // 2: resources.timestamp.Timestamp
	(*Colleague)(nil),           // 3: resources.jobs.Colleague
}
var file_resources_jobs_conduct_proto_depIdxs = []int32{
	2, // 0: resources.jobs.ConductEntry.created_at:type_name -> resources.timestamp.Timestamp
	2, // 1: resources.jobs.ConductEntry.updated_at:type_name -> resources.timestamp.Timestamp
	0, // 2: resources.jobs.ConductEntry.type:type_name -> resources.jobs.ConductType
	2, // 3: resources.jobs.ConductEntry.expires_at:type_name -> resources.timestamp.Timestamp
	3, // 4: resources.jobs.ConductEntry.target_user:type_name -> resources.jobs.Colleague
	3, // 5: resources.jobs.ConductEntry.creator:type_name -> resources.jobs.Colleague
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_resources_jobs_conduct_proto_init() }
func file_resources_jobs_conduct_proto_init() {
	if File_resources_jobs_conduct_proto != nil {
		return
	}
	file_resources_jobs_colleagues_proto_init()
	file_resources_jobs_conduct_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_jobs_conduct_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_jobs_conduct_proto_goTypes,
		DependencyIndexes: file_resources_jobs_conduct_proto_depIdxs,
		EnumInfos:         file_resources_jobs_conduct_proto_enumTypes,
		MessageInfos:      file_resources_jobs_conduct_proto_msgTypes,
	}.Build()
	File_resources_jobs_conduct_proto = out.File
	file_resources_jobs_conduct_proto_rawDesc = nil
	file_resources_jobs_conduct_proto_goTypes = nil
	file_resources_jobs_conduct_proto_depIdxs = nil
}
