// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.20.3
// source: resources/documents/access.proto

package documents

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
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

type AccessLevelUpdateMode int32

const (
	AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED AccessLevelUpdateMode = 0
	AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE      AccessLevelUpdateMode = 1
	AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_DELETE      AccessLevelUpdateMode = 2
	AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_CLEAR       AccessLevelUpdateMode = 3
)

// Enum value maps for AccessLevelUpdateMode.
var (
	AccessLevelUpdateMode_name = map[int32]string{
		0: "ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED",
		1: "ACCESS_LEVEL_UPDATE_MODE_UPDATE",
		2: "ACCESS_LEVEL_UPDATE_MODE_DELETE",
		3: "ACCESS_LEVEL_UPDATE_MODE_CLEAR",
	}
	AccessLevelUpdateMode_value = map[string]int32{
		"ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED": 0,
		"ACCESS_LEVEL_UPDATE_MODE_UPDATE":      1,
		"ACCESS_LEVEL_UPDATE_MODE_DELETE":      2,
		"ACCESS_LEVEL_UPDATE_MODE_CLEAR":       3,
	}
)

func (x AccessLevelUpdateMode) Enum() *AccessLevelUpdateMode {
	p := new(AccessLevelUpdateMode)
	*p = x
	return p
}

func (x AccessLevelUpdateMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccessLevelUpdateMode) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_documents_access_proto_enumTypes[0].Descriptor()
}

func (AccessLevelUpdateMode) Type() protoreflect.EnumType {
	return &file_resources_documents_access_proto_enumTypes[0]
}

func (x AccessLevelUpdateMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccessLevelUpdateMode.Descriptor instead.
func (AccessLevelUpdateMode) EnumDescriptor() ([]byte, []int) {
	return file_resources_documents_access_proto_rawDescGZIP(), []int{0}
}

type AccessLevel int32

const (
	AccessLevel_ACCESS_LEVEL_UNSPECIFIED AccessLevel = 0
	AccessLevel_ACCESS_LEVEL_BLOCKED     AccessLevel = 1
	AccessLevel_ACCESS_LEVEL_VIEW        AccessLevel = 2
	AccessLevel_ACCESS_LEVEL_COMMENT     AccessLevel = 3
	AccessLevel_ACCESS_LEVEL_STATUS      AccessLevel = 4
	AccessLevel_ACCESS_LEVEL_ACCESS      AccessLevel = 5
	AccessLevel_ACCESS_LEVEL_EDIT        AccessLevel = 6
)

// Enum value maps for AccessLevel.
var (
	AccessLevel_name = map[int32]string{
		0: "ACCESS_LEVEL_UNSPECIFIED",
		1: "ACCESS_LEVEL_BLOCKED",
		2: "ACCESS_LEVEL_VIEW",
		3: "ACCESS_LEVEL_COMMENT",
		4: "ACCESS_LEVEL_STATUS",
		5: "ACCESS_LEVEL_ACCESS",
		6: "ACCESS_LEVEL_EDIT",
	}
	AccessLevel_value = map[string]int32{
		"ACCESS_LEVEL_UNSPECIFIED": 0,
		"ACCESS_LEVEL_BLOCKED":     1,
		"ACCESS_LEVEL_VIEW":        2,
		"ACCESS_LEVEL_COMMENT":     3,
		"ACCESS_LEVEL_STATUS":      4,
		"ACCESS_LEVEL_ACCESS":      5,
		"ACCESS_LEVEL_EDIT":        6,
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
	return file_resources_documents_access_proto_enumTypes[1].Descriptor()
}

func (AccessLevel) Type() protoreflect.EnumType {
	return &file_resources_documents_access_proto_enumTypes[1]
}

func (x AccessLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccessLevel.Descriptor instead.
func (AccessLevel) EnumDescriptor() ([]byte, []int) {
	return file_resources_documents_access_proto_rawDescGZIP(), []int{1}
}

type DocumentAccess struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jobs  []*DocumentJobAccess  `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty" alias:"job_access"`   // @gotags: alias:"job_access"
	Users []*DocumentUserAccess `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty" alias:"user_access"` // @gotags: alias:"user_access"
}

func (x *DocumentAccess) Reset() {
	*x = DocumentAccess{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_documents_access_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocumentAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentAccess) ProtoMessage() {}

func (x *DocumentAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_access_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentAccess.ProtoReflect.Descriptor instead.
func (*DocumentAccess) Descriptor() ([]byte, []int) {
	return file_resources_documents_access_proto_rawDescGZIP(), []int{0}
}

func (x *DocumentAccess) GetJobs() []*DocumentJobAccess {
	if x != nil {
		return x.Jobs
	}
	return nil
}

func (x *DocumentAccess) GetUsers() []*DocumentUserAccess {
	if x != nil {
		return x.Users
	}
	return nil
}

type DocumentJobAccess struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	DocumentId    uint64               `protobuf:"varint,3,opt,name=document_id,json=documentId,proto3" json:"document_id,omitempty"`
	Job           string               `protobuf:"bytes,4,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel      *string              `protobuf:"bytes,5,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty" alias:"job_label"` // @gotags: alias:"job_label"
	MinimumGrade  int32                `protobuf:"varint,6,opt,name=minimum_grade,json=minimumGrade,proto3" json:"minimum_grade,omitempty"`
	JobGradeLabel *string              `protobuf:"bytes,7,opt,name=job_grade_label,json=jobGradeLabel,proto3,oneof" json:"job_grade_label,omitempty" alias:"job_grade_label"` // @gotags: alias:"job_grade_label"
	Access        AccessLevel          `protobuf:"varint,8,opt,name=access,proto3,enum=resources.documents.AccessLevel" json:"access,omitempty"`
	Required      *bool                `protobuf:"varint,9,opt,name=required,proto3,oneof" json:"required,omitempty" alias:"required"` // @gotags: alias:"required"
}

func (x *DocumentJobAccess) Reset() {
	*x = DocumentJobAccess{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_documents_access_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocumentJobAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentJobAccess) ProtoMessage() {}

func (x *DocumentJobAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_access_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentJobAccess.ProtoReflect.Descriptor instead.
func (*DocumentJobAccess) Descriptor() ([]byte, []int) {
	return file_resources_documents_access_proto_rawDescGZIP(), []int{1}
}

func (x *DocumentJobAccess) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DocumentJobAccess) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *DocumentJobAccess) GetDocumentId() uint64 {
	if x != nil {
		return x.DocumentId
	}
	return 0
}

func (x *DocumentJobAccess) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *DocumentJobAccess) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
	}
	return ""
}

func (x *DocumentJobAccess) GetMinimumGrade() int32 {
	if x != nil {
		return x.MinimumGrade
	}
	return 0
}

func (x *DocumentJobAccess) GetJobGradeLabel() string {
	if x != nil && x.JobGradeLabel != nil {
		return *x.JobGradeLabel
	}
	return ""
}

func (x *DocumentJobAccess) GetAccess() AccessLevel {
	if x != nil {
		return x.Access
	}
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

func (x *DocumentJobAccess) GetRequired() bool {
	if x != nil && x.Required != nil {
		return *x.Required
	}
	return false
}

type DocumentUserAccess struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt  *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	DocumentId uint64               `protobuf:"varint,3,opt,name=document_id,json=documentId,proto3" json:"document_id,omitempty"`
	UserId     int32                `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	User       *users.UserShort     `protobuf:"bytes,5,opt,name=user,proto3,oneof" json:"user,omitempty"`
	Access     AccessLevel          `protobuf:"varint,6,opt,name=access,proto3,enum=resources.documents.AccessLevel" json:"access,omitempty"`
	Required   *bool                `protobuf:"varint,7,opt,name=required,proto3,oneof" json:"required,omitempty" alias:"required"` // @gotags: alias:"required"
}

func (x *DocumentUserAccess) Reset() {
	*x = DocumentUserAccess{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_documents_access_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocumentUserAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentUserAccess) ProtoMessage() {}

func (x *DocumentUserAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_access_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentUserAccess.ProtoReflect.Descriptor instead.
func (*DocumentUserAccess) Descriptor() ([]byte, []int) {
	return file_resources_documents_access_proto_rawDescGZIP(), []int{2}
}

func (x *DocumentUserAccess) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DocumentUserAccess) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *DocumentUserAccess) GetDocumentId() uint64 {
	if x != nil {
		return x.DocumentId
	}
	return 0
}

func (x *DocumentUserAccess) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *DocumentUserAccess) GetUser() *users.UserShort {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *DocumentUserAccess) GetAccess() AccessLevel {
	if x != nil {
		return x.Access
	}
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

func (x *DocumentUserAccess) GetRequired() bool {
	if x != nil && x.Required != nil {
		return *x.Required
	}
	return false
}

var File_resources_documents_access_proto protoreflect.FileDescriptor

var file_resources_documents_access_proto_rawDesc = []byte{
	0x0a, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x13, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x9f, 0x01, 0x0a, 0x0e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x44, 0x0a, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x4a, 0x6f, 0x62, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42, 0x08, 0xfa, 0x42, 0x05,
	0x92, 0x01, 0x02, 0x10, 0x14, 0x52, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x12, 0x47, 0x0a, 0x05, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x92, 0x01, 0x02, 0x10, 0x14, 0x52, 0x05, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x22, 0xdd, 0x03, 0x0a, 0x11, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x4a, 0x6f, 0x62, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x30, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x23, 0x0a, 0x0b, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x30, 0x01, 0x52, 0x0a, 0x64, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52, 0x03, 0x6a,
	0x6f, 0x62, 0x12, 0x29, 0x0a, 0x09, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x01,
	0x52, 0x08, 0x6a, 0x6f, 0x62, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x2c, 0x0a,
	0x0d, 0x6d, 0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x0c, 0x6d,
	0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x47, 0x72, 0x61, 0x64, 0x65, 0x12, 0x34, 0x0a, 0x0f, 0x6a,
	0x6f, 0x62, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x02, 0x52,
	0x0d, 0x6a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64, 0x65, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x88, 0x01,
	0x01, 0x12, 0x42, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x20, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x65,
	0x76, 0x65, 0x6c, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x61,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1f, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x48, 0x03, 0x52, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x67, 0x72, 0x61, 0x64,
	0x65, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x72, 0x65, 0x71, 0x75,
	0x69, 0x72, 0x65, 0x64, 0x22, 0xf2, 0x02, 0x0a, 0x12, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x30, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x0b, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x30, 0x01, 0x52, 0x0a, 0x64, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02,
	0x20, 0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x48, 0x01, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12,
	0x42, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x20, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65,
	0x6c, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x1f, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x08, 0x48, 0x02, 0x52, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x64, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x0b, 0x0a, 0x09,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2a, 0xaf, 0x01, 0x0a, 0x15, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d,
	0x6f, 0x64, 0x65, 0x12, 0x28, 0x0a, 0x24, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45,
	0x56, 0x45, 0x4c, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x23, 0x0a,
	0x1f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x55, 0x50,
	0x44, 0x41, 0x54, 0x45, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45,
	0x10, 0x01, 0x12, 0x23, 0x0a, 0x1f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56,
	0x45, 0x4c, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x44,
	0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x02, 0x12, 0x22, 0x0a, 0x1e, 0x41, 0x43, 0x43, 0x45, 0x53,
	0x53, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x4d,
	0x4f, 0x44, 0x45, 0x5f, 0x43, 0x4c, 0x45, 0x41, 0x52, 0x10, 0x03, 0x2a, 0xbf, 0x01, 0x0a, 0x0b,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x1c, 0x0a, 0x18, 0x41,
	0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x43, 0x43,
	0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x45,
	0x44, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45,
	0x56, 0x45, 0x4c, 0x5f, 0x56, 0x49, 0x45, 0x57, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x43,
	0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x45,
	0x4e, 0x54, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c,
	0x45, 0x56, 0x45, 0x4c, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x10, 0x04, 0x12, 0x17, 0x0a,
	0x13, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x41, 0x43,
	0x43, 0x45, 0x53, 0x53, 0x10, 0x05, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53,
	0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x45, 0x44, 0x49, 0x54, 0x10, 0x06, 0x42, 0x4b, 0x5a,
	0x49, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65,
	0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x3b, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_resources_documents_access_proto_rawDescOnce sync.Once
	file_resources_documents_access_proto_rawDescData = file_resources_documents_access_proto_rawDesc
)

func file_resources_documents_access_proto_rawDescGZIP() []byte {
	file_resources_documents_access_proto_rawDescOnce.Do(func() {
		file_resources_documents_access_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_documents_access_proto_rawDescData)
	})
	return file_resources_documents_access_proto_rawDescData
}

var file_resources_documents_access_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_resources_documents_access_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_resources_documents_access_proto_goTypes = []any{
	(AccessLevelUpdateMode)(0),  // 0: resources.documents.AccessLevelUpdateMode
	(AccessLevel)(0),            // 1: resources.documents.AccessLevel
	(*DocumentAccess)(nil),      // 2: resources.documents.DocumentAccess
	(*DocumentJobAccess)(nil),   // 3: resources.documents.DocumentJobAccess
	(*DocumentUserAccess)(nil),  // 4: resources.documents.DocumentUserAccess
	(*timestamp.Timestamp)(nil), // 5: resources.timestamp.Timestamp
	(*users.UserShort)(nil),     // 6: resources.users.UserShort
}
var file_resources_documents_access_proto_depIdxs = []int32{
	3, // 0: resources.documents.DocumentAccess.jobs:type_name -> resources.documents.DocumentJobAccess
	4, // 1: resources.documents.DocumentAccess.users:type_name -> resources.documents.DocumentUserAccess
	5, // 2: resources.documents.DocumentJobAccess.created_at:type_name -> resources.timestamp.Timestamp
	1, // 3: resources.documents.DocumentJobAccess.access:type_name -> resources.documents.AccessLevel
	5, // 4: resources.documents.DocumentUserAccess.created_at:type_name -> resources.timestamp.Timestamp
	6, // 5: resources.documents.DocumentUserAccess.user:type_name -> resources.users.UserShort
	1, // 6: resources.documents.DocumentUserAccess.access:type_name -> resources.documents.AccessLevel
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_resources_documents_access_proto_init() }
func file_resources_documents_access_proto_init() {
	if File_resources_documents_access_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_documents_access_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*DocumentAccess); i {
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
		file_resources_documents_access_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*DocumentJobAccess); i {
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
		file_resources_documents_access_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*DocumentUserAccess); i {
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
	file_resources_documents_access_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_documents_access_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_documents_access_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_documents_access_proto_goTypes,
		DependencyIndexes: file_resources_documents_access_proto_depIdxs,
		EnumInfos:         file_resources_documents_access_proto_enumTypes,
		MessageInfos:      file_resources_documents_access_proto_msgTypes,
	}.Build()
	File_resources_documents_access_proto = out.File
	file_resources_documents_access_proto_rawDesc = nil
	file_resources_documents_access_proto_goTypes = nil
	file_resources_documents_access_proto_depIdxs = nil
}
