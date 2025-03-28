// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: resources/wiki/access.proto

package wiki

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type AccessLevel int32

const (
	AccessLevel_ACCESS_LEVEL_UNSPECIFIED AccessLevel = 0
	AccessLevel_ACCESS_LEVEL_BLOCKED     AccessLevel = 1
	AccessLevel_ACCESS_LEVEL_VIEW        AccessLevel = 2
	AccessLevel_ACCESS_LEVEL_ACCESS      AccessLevel = 3
	AccessLevel_ACCESS_LEVEL_EDIT        AccessLevel = 4
)

// Enum value maps for AccessLevel.
var (
	AccessLevel_name = map[int32]string{
		0: "ACCESS_LEVEL_UNSPECIFIED",
		1: "ACCESS_LEVEL_BLOCKED",
		2: "ACCESS_LEVEL_VIEW",
		3: "ACCESS_LEVEL_ACCESS",
		4: "ACCESS_LEVEL_EDIT",
	}
	AccessLevel_value = map[string]int32{
		"ACCESS_LEVEL_UNSPECIFIED": 0,
		"ACCESS_LEVEL_BLOCKED":     1,
		"ACCESS_LEVEL_VIEW":        2,
		"ACCESS_LEVEL_ACCESS":      3,
		"ACCESS_LEVEL_EDIT":        4,
	}
)

func (x AccessLevel) Enum() *AccessLevel {
	p := new(AccessLevel)
	*p = x
	return p
}

func (x AccessLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccessLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_wiki_access_proto_enumTypes[0].Descriptor()
}

func (AccessLevel) Type() protoreflect.EnumType {
	return &file_resources_wiki_access_proto_enumTypes[0]
}

func (x AccessLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccessLevel.Descriptor instead.
func (AccessLevel) EnumDescriptor() ([]byte, []int) {
	return file_resources_wiki_access_proto_rawDescGZIP(), []int{0}
}

type PageAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Jobs          []*PageJobAccess       `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty" alias:"job_access"`   // @gotags: alias:"job_access"
	Users         []*PageUserAccess      `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty" alias:"user_access"` // @gotags: alias:"user_access"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PageAccess) Reset() {
	*x = PageAccess{}
	mi := &file_resources_wiki_access_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PageAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageAccess) ProtoMessage() {}

func (x *PageAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_wiki_access_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageAccess.ProtoReflect.Descriptor instead.
func (*PageAccess) Descriptor() ([]byte, []int) {
	return file_resources_wiki_access_proto_rawDescGZIP(), []int{0}
}

func (x *PageAccess) GetJobs() []*PageJobAccess {
	if x != nil {
		return x.Jobs
	}
	return nil
}

func (x *PageAccess) GetUsers() []*PageUserAccess {
	if x != nil {
		return x.Users
	}
	return nil
}

type PageJobAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	TargetId      uint64                 `protobuf:"varint,3,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`
	Job           string                 `protobuf:"bytes,4,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel      *string                `protobuf:"bytes,5,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty"`
	MinimumGrade  int32                  `protobuf:"varint,6,opt,name=minimum_grade,json=minimumGrade,proto3" json:"minimum_grade,omitempty"`
	JobGradeLabel *string                `protobuf:"bytes,7,opt,name=job_grade_label,json=jobGradeLabel,proto3,oneof" json:"job_grade_label,omitempty"`
	Access        AccessLevel            `protobuf:"varint,8,opt,name=access,proto3,enum=resources.wiki.AccessLevel" json:"access,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PageJobAccess) Reset() {
	*x = PageJobAccess{}
	mi := &file_resources_wiki_access_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PageJobAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageJobAccess) ProtoMessage() {}

func (x *PageJobAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_wiki_access_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageJobAccess.ProtoReflect.Descriptor instead.
func (*PageJobAccess) Descriptor() ([]byte, []int) {
	return file_resources_wiki_access_proto_rawDescGZIP(), []int{1}
}

func (x *PageJobAccess) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PageJobAccess) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *PageJobAccess) GetTargetId() uint64 {
	if x != nil {
		return x.TargetId
	}
	return 0
}

func (x *PageJobAccess) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *PageJobAccess) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
	}
	return ""
}

func (x *PageJobAccess) GetMinimumGrade() int32 {
	if x != nil {
		return x.MinimumGrade
	}
	return 0
}

func (x *PageJobAccess) GetJobGradeLabel() string {
	if x != nil && x.JobGradeLabel != nil {
		return *x.JobGradeLabel
	}
	return ""
}

func (x *PageJobAccess) GetAccess() AccessLevel {
	if x != nil {
		return x.Access
	}
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

type PageUserAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	TargetId      uint64                 `protobuf:"varint,3,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`
	UserId        int32                  `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	User          *users.UserShort       `protobuf:"bytes,5,opt,name=user,proto3,oneof" json:"user,omitempty"`
	Access        AccessLevel            `protobuf:"varint,6,opt,name=access,proto3,enum=resources.wiki.AccessLevel" json:"access,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PageUserAccess) Reset() {
	*x = PageUserAccess{}
	mi := &file_resources_wiki_access_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PageUserAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageUserAccess) ProtoMessage() {}

func (x *PageUserAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_wiki_access_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageUserAccess.ProtoReflect.Descriptor instead.
func (*PageUserAccess) Descriptor() ([]byte, []int) {
	return file_resources_wiki_access_proto_rawDescGZIP(), []int{2}
}

func (x *PageUserAccess) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PageUserAccess) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *PageUserAccess) GetTargetId() uint64 {
	if x != nil {
		return x.TargetId
	}
	return 0
}

func (x *PageUserAccess) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *PageUserAccess) GetUser() *users.UserShort {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *PageUserAccess) GetAccess() AccessLevel {
	if x != nil {
		return x.Access
	}
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

var File_resources_wiki_access_proto protoreflect.FileDescriptor

var file_resources_wiki_access_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x77, 0x69, 0x6b, 0x69,
	0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x77, 0x69, 0x6b, 0x69, 0x1a, 0x23, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x89, 0x01, 0x0a, 0x0a, 0x50, 0x61, 0x67,
	0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x3b, 0x0a, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x77, 0x69, 0x6b, 0x69, 0x2e, 0x50, 0x61, 0x67, 0x65, 0x4a, 0x6f, 0x62, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x92, 0x01, 0x02, 0x10, 0x14, 0x52, 0x04,
	0x6a, 0x6f, 0x62, 0x73, 0x12, 0x3e, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x77, 0x69, 0x6b, 0x69, 0x2e, 0x50, 0x61, 0x67, 0x65, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x92, 0x01, 0x02, 0x10, 0x14, 0x52, 0x05, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x22, 0x9a, 0x03, 0x0a, 0x0d, 0x50, 0x61, 0x67, 0x65, 0x4a, 0x6f, 0x62,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52, 0x03, 0x6a,
	0x6f, 0x62, 0x12, 0x29, 0x0a, 0x09, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x01,
	0x52, 0x08, 0x6a, 0x6f, 0x62, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x2c, 0x0a,
	0x0d, 0x6d, 0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x28, 0x00, 0x52, 0x0c, 0x6d,
	0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x47, 0x72, 0x61, 0x64, 0x65, 0x12, 0x34, 0x0a, 0x0f, 0x6a,
	0x6f, 0x62, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x02, 0x52,
	0x0d, 0x6a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64, 0x65, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x88, 0x01,
	0x01, 0x12, 0x3d, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1b, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x77, 0x69,
	0x6b, 0x69, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42,
	0x0c, 0x0a, 0x0a, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x42, 0x12, 0x0a,
	0x10, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x22, 0xaf, 0x02, 0x0a, 0x0e, 0x50, 0x61, 0x67, 0x65, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72,
	0x74, 0x48, 0x01, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x3d, 0x0a, 0x06,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x77, 0x69, 0x6b, 0x69, 0x2e, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01,
	0x02, 0x10, 0x01, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x2a, 0x8c, 0x01, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x65,
	0x76, 0x65, 0x6c, 0x12, 0x1c, 0x0a, 0x18, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45,
	0x56, 0x45, 0x4c, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56, 0x45,
	0x4c, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x45, 0x44, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x41,
	0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x56, 0x49, 0x45, 0x57,
	0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56,
	0x45, 0x4c, 0x5f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x03, 0x12, 0x15, 0x0a, 0x11, 0x41,
	0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x45, 0x44, 0x49, 0x54,
	0x10, 0x04, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76,
	0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x77, 0x69, 0x6b, 0x69,
	0x3b, 0x77, 0x69, 0x6b, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_wiki_access_proto_rawDescOnce sync.Once
	file_resources_wiki_access_proto_rawDescData []byte
)

func file_resources_wiki_access_proto_rawDescGZIP() []byte {
	file_resources_wiki_access_proto_rawDescOnce.Do(func() {
		file_resources_wiki_access_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_wiki_access_proto_rawDesc), len(file_resources_wiki_access_proto_rawDesc)))
	})
	return file_resources_wiki_access_proto_rawDescData
}

var file_resources_wiki_access_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_wiki_access_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_resources_wiki_access_proto_goTypes = []any{
	(AccessLevel)(0),            // 0: resources.wiki.AccessLevel
	(*PageAccess)(nil),          // 1: resources.wiki.PageAccess
	(*PageJobAccess)(nil),       // 2: resources.wiki.PageJobAccess
	(*PageUserAccess)(nil),      // 3: resources.wiki.PageUserAccess
	(*timestamp.Timestamp)(nil), // 4: resources.timestamp.Timestamp
	(*users.UserShort)(nil),     // 5: resources.users.UserShort
}
var file_resources_wiki_access_proto_depIdxs = []int32{
	2, // 0: resources.wiki.PageAccess.jobs:type_name -> resources.wiki.PageJobAccess
	3, // 1: resources.wiki.PageAccess.users:type_name -> resources.wiki.PageUserAccess
	4, // 2: resources.wiki.PageJobAccess.created_at:type_name -> resources.timestamp.Timestamp
	0, // 3: resources.wiki.PageJobAccess.access:type_name -> resources.wiki.AccessLevel
	4, // 4: resources.wiki.PageUserAccess.created_at:type_name -> resources.timestamp.Timestamp
	5, // 5: resources.wiki.PageUserAccess.user:type_name -> resources.users.UserShort
	0, // 6: resources.wiki.PageUserAccess.access:type_name -> resources.wiki.AccessLevel
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_resources_wiki_access_proto_init() }
func file_resources_wiki_access_proto_init() {
	if File_resources_wiki_access_proto != nil {
		return
	}
	file_resources_wiki_access_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_wiki_access_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_wiki_access_proto_rawDesc), len(file_resources_wiki_access_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_wiki_access_proto_goTypes,
		DependencyIndexes: file_resources_wiki_access_proto_depIdxs,
		EnumInfos:         file_resources_wiki_access_proto_enumTypes,
		MessageInfos:      file_resources_wiki_access_proto_msgTypes,
	}.Build()
	File_resources_wiki_access_proto = out.File
	file_resources_wiki_access_proto_goTypes = nil
	file_resources_wiki_access_proto_depIdxs = nil
}
