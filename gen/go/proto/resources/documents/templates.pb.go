// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: resources/documents/templates.proto

package documents

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	vehicles "github.com/fivenet-app/fivenet/gen/go/proto/resources/vehicles"
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

type Template struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" alias:"id"` // @gotags: alias:"id"
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	Category  *Category              `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty" alias:"category"` // @gotags: alias:"category"
	Weight    uint32                 `protobuf:"varint,5,opt,name=weight,proto3" json:"weight,omitempty"`
	// @sanitize
	Title string `protobuf:"bytes,6,opt,name=title,proto3" json:"title,omitempty"`
	// @sanitize
	Description string `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	// @sanitize: method=StripTags
	Color *string `protobuf:"bytes,8,opt,name=color,proto3,oneof" json:"color,omitempty"`
	// @sanitize: method=StripTags
	Icon            *string              `protobuf:"bytes,9,opt,name=icon,proto3,oneof" json:"icon,omitempty"`
	ContentTitle    string               `protobuf:"bytes,10,opt,name=content_title,json=contentTitle,proto3" json:"content_title,omitempty" alias:"content_title"` // @gotags: alias:"content_title"
	Content         string               `protobuf:"bytes,11,opt,name=content,proto3" json:"content,omitempty" alias:"content"`                               // @gotags: alias:"content"
	State           string               `protobuf:"bytes,12,opt,name=state,proto3" json:"state,omitempty" alias:"state"`                                   // @gotags: alias:"state"
	Schema          *TemplateSchema      `protobuf:"bytes,13,opt,name=schema,proto3" json:"schema,omitempty" alias:"schema"`                                 // @gotags: alias:"schema"
	CreatorJob      string               `protobuf:"bytes,14,opt,name=creator_job,json=creatorJob,proto3" json:"creator_job,omitempty"`
	CreatorJobLabel *string              `protobuf:"bytes,15,opt,name=creator_job_label,json=creatorJobLabel,proto3,oneof" json:"creator_job_label,omitempty"`
	JobAccess       []*TemplateJobAccess `protobuf:"bytes,16,rep,name=job_access,json=jobAccess,proto3" json:"job_access,omitempty"`
	ContentAccess   *DocumentAccess      `protobuf:"bytes,17,opt,name=content_access,json=contentAccess,proto3" json:"content_access,omitempty" alias:"access"` // @gotags: alias:"access"
	Workflow        *Workflow            `protobuf:"bytes,18,opt,name=workflow,proto3,oneof" json:"workflow,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Template) Reset() {
	*x = Template{}
	mi := &file_resources_documents_templates_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Template) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Template) ProtoMessage() {}

func (x *Template) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_templates_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Template.ProtoReflect.Descriptor instead.
func (*Template) Descriptor() ([]byte, []int) {
	return file_resources_documents_templates_proto_rawDescGZIP(), []int{0}
}

func (x *Template) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Template) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Template) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Template) GetCategory() *Category {
	if x != nil {
		return x.Category
	}
	return nil
}

func (x *Template) GetWeight() uint32 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *Template) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Template) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Template) GetColor() string {
	if x != nil && x.Color != nil {
		return *x.Color
	}
	return ""
}

func (x *Template) GetIcon() string {
	if x != nil && x.Icon != nil {
		return *x.Icon
	}
	return ""
}

func (x *Template) GetContentTitle() string {
	if x != nil {
		return x.ContentTitle
	}
	return ""
}

func (x *Template) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Template) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *Template) GetSchema() *TemplateSchema {
	if x != nil {
		return x.Schema
	}
	return nil
}

func (x *Template) GetCreatorJob() string {
	if x != nil {
		return x.CreatorJob
	}
	return ""
}

func (x *Template) GetCreatorJobLabel() string {
	if x != nil && x.CreatorJobLabel != nil {
		return *x.CreatorJobLabel
	}
	return ""
}

func (x *Template) GetJobAccess() []*TemplateJobAccess {
	if x != nil {
		return x.JobAccess
	}
	return nil
}

func (x *Template) GetContentAccess() *DocumentAccess {
	if x != nil {
		return x.ContentAccess
	}
	return nil
}

func (x *Template) GetWorkflow() *Workflow {
	if x != nil {
		return x.Workflow
	}
	return nil
}

type TemplateShort struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" alias:"id"` // @gotags: alias:"id"
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	Category  *Category              `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty" alias:"category"` // @gotags: alias:"category"
	Weight    uint32                 `protobuf:"varint,5,opt,name=weight,proto3" json:"weight,omitempty"`
	// @sanitize
	Title string `protobuf:"bytes,6,opt,name=title,proto3" json:"title,omitempty"`
	// @sanitize
	Description string `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	// @sanitize: method=StripTags
	Color *string `protobuf:"bytes,8,opt,name=color,proto3,oneof" json:"color,omitempty"`
	// @sanitize: method=StripTags
	Icon            *string         `protobuf:"bytes,9,opt,name=icon,proto3,oneof" json:"icon,omitempty"`
	Schema          *TemplateSchema `protobuf:"bytes,10,opt,name=schema,proto3" json:"schema,omitempty" alias:"schema"` // @gotags: alias:"schema"
	CreatorJob      string          `protobuf:"bytes,11,opt,name=creator_job,json=creatorJob,proto3" json:"creator_job,omitempty"`
	CreatorJobLabel *string         `protobuf:"bytes,12,opt,name=creator_job_label,json=creatorJobLabel,proto3,oneof" json:"creator_job_label,omitempty"`
	Workflow        *Workflow       `protobuf:"bytes,18,opt,name=workflow,proto3,oneof" json:"workflow,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *TemplateShort) Reset() {
	*x = TemplateShort{}
	mi := &file_resources_documents_templates_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TemplateShort) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateShort) ProtoMessage() {}

func (x *TemplateShort) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_templates_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateShort.ProtoReflect.Descriptor instead.
func (*TemplateShort) Descriptor() ([]byte, []int) {
	return file_resources_documents_templates_proto_rawDescGZIP(), []int{1}
}

func (x *TemplateShort) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *TemplateShort) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *TemplateShort) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *TemplateShort) GetCategory() *Category {
	if x != nil {
		return x.Category
	}
	return nil
}

func (x *TemplateShort) GetWeight() uint32 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *TemplateShort) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *TemplateShort) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *TemplateShort) GetColor() string {
	if x != nil && x.Color != nil {
		return *x.Color
	}
	return ""
}

func (x *TemplateShort) GetIcon() string {
	if x != nil && x.Icon != nil {
		return *x.Icon
	}
	return ""
}

func (x *TemplateShort) GetSchema() *TemplateSchema {
	if x != nil {
		return x.Schema
	}
	return nil
}

func (x *TemplateShort) GetCreatorJob() string {
	if x != nil {
		return x.CreatorJob
	}
	return ""
}

func (x *TemplateShort) GetCreatorJobLabel() string {
	if x != nil && x.CreatorJobLabel != nil {
		return *x.CreatorJobLabel
	}
	return ""
}

func (x *TemplateShort) GetWorkflow() *Workflow {
	if x != nil {
		return x.Workflow
	}
	return nil
}

// @dbscanner: json
type TemplateSchema struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Requirements  *TemplateRequirements  `protobuf:"bytes,1,opt,name=requirements,proto3" json:"requirements,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TemplateSchema) Reset() {
	*x = TemplateSchema{}
	mi := &file_resources_documents_templates_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TemplateSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateSchema) ProtoMessage() {}

func (x *TemplateSchema) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_templates_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateSchema.ProtoReflect.Descriptor instead.
func (*TemplateSchema) Descriptor() ([]byte, []int) {
	return file_resources_documents_templates_proto_rawDescGZIP(), []int{2}
}

func (x *TemplateSchema) GetRequirements() *TemplateRequirements {
	if x != nil {
		return x.Requirements
	}
	return nil
}

type TemplateRequirements struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Documents     *ObjectSpecs           `protobuf:"bytes,1,opt,name=documents,proto3,oneof" json:"documents,omitempty"`
	Users         *ObjectSpecs           `protobuf:"bytes,2,opt,name=users,proto3,oneof" json:"users,omitempty"`
	Vehicles      *ObjectSpecs           `protobuf:"bytes,3,opt,name=vehicles,proto3,oneof" json:"vehicles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TemplateRequirements) Reset() {
	*x = TemplateRequirements{}
	mi := &file_resources_documents_templates_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TemplateRequirements) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateRequirements) ProtoMessage() {}

func (x *TemplateRequirements) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_templates_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateRequirements.ProtoReflect.Descriptor instead.
func (*TemplateRequirements) Descriptor() ([]byte, []int) {
	return file_resources_documents_templates_proto_rawDescGZIP(), []int{3}
}

func (x *TemplateRequirements) GetDocuments() *ObjectSpecs {
	if x != nil {
		return x.Documents
	}
	return nil
}

func (x *TemplateRequirements) GetUsers() *ObjectSpecs {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *TemplateRequirements) GetVehicles() *ObjectSpecs {
	if x != nil {
		return x.Vehicles
	}
	return nil
}

type ObjectSpecs struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Required      *bool                  `protobuf:"varint,1,opt,name=required,proto3,oneof" json:"required,omitempty"`
	Min           *int32                 `protobuf:"varint,2,opt,name=min,proto3,oneof" json:"min,omitempty"`
	Max           *int32                 `protobuf:"varint,3,opt,name=max,proto3,oneof" json:"max,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ObjectSpecs) Reset() {
	*x = ObjectSpecs{}
	mi := &file_resources_documents_templates_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ObjectSpecs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ObjectSpecs) ProtoMessage() {}

func (x *ObjectSpecs) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_templates_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ObjectSpecs.ProtoReflect.Descriptor instead.
func (*ObjectSpecs) Descriptor() ([]byte, []int) {
	return file_resources_documents_templates_proto_rawDescGZIP(), []int{4}
}

func (x *ObjectSpecs) GetRequired() bool {
	if x != nil && x.Required != nil {
		return *x.Required
	}
	return false
}

func (x *ObjectSpecs) GetMin() int32 {
	if x != nil && x.Min != nil {
		return *x.Min
	}
	return 0
}

func (x *ObjectSpecs) GetMax() int32 {
	if x != nil && x.Max != nil {
		return *x.Max
	}
	return 0
}

type TemplateData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ActiveChar    *users.User            `protobuf:"bytes,1,opt,name=activeChar,proto3" json:"activeChar,omitempty"`
	Documents     []*DocumentShort       `protobuf:"bytes,2,rep,name=documents,proto3" json:"documents,omitempty"`
	Users         []*users.UserShort     `protobuf:"bytes,3,rep,name=users,proto3" json:"users,omitempty"`
	Vehicles      []*vehicles.Vehicle    `protobuf:"bytes,4,rep,name=vehicles,proto3" json:"vehicles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TemplateData) Reset() {
	*x = TemplateData{}
	mi := &file_resources_documents_templates_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TemplateData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateData) ProtoMessage() {}

func (x *TemplateData) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_templates_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateData.ProtoReflect.Descriptor instead.
func (*TemplateData) Descriptor() ([]byte, []int) {
	return file_resources_documents_templates_proto_rawDescGZIP(), []int{5}
}

func (x *TemplateData) GetActiveChar() *users.User {
	if x != nil {
		return x.ActiveChar
	}
	return nil
}

func (x *TemplateData) GetDocuments() []*DocumentShort {
	if x != nil {
		return x.Documents
	}
	return nil
}

func (x *TemplateData) GetUsers() []*users.UserShort {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *TemplateData) GetVehicles() []*vehicles.Vehicle {
	if x != nil {
		return x.Vehicles
	}
	return nil
}

type TemplateJobAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" alias:"id"` // @gotags: alias:"id"
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	TargetId      uint64                 `protobuf:"varint,3,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty" alias:"template_id"` // @gotags: alias:"template_id"
	Job           string                 `protobuf:"bytes,4,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel      *string                `protobuf:"bytes,5,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty"`
	MinimumGrade  int32                  `protobuf:"varint,6,opt,name=minimum_grade,json=minimumGrade,proto3" json:"minimum_grade,omitempty"`
	JobGradeLabel *string                `protobuf:"bytes,7,opt,name=job_grade_label,json=jobGradeLabel,proto3,oneof" json:"job_grade_label,omitempty"`
	Access        AccessLevel            `protobuf:"varint,8,opt,name=access,proto3,enum=resources.documents.AccessLevel" json:"access,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TemplateJobAccess) Reset() {
	*x = TemplateJobAccess{}
	mi := &file_resources_documents_templates_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TemplateJobAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateJobAccess) ProtoMessage() {}

func (x *TemplateJobAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_templates_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateJobAccess.ProtoReflect.Descriptor instead.
func (*TemplateJobAccess) Descriptor() ([]byte, []int) {
	return file_resources_documents_templates_proto_rawDescGZIP(), []int{6}
}

func (x *TemplateJobAccess) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *TemplateJobAccess) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *TemplateJobAccess) GetTargetId() uint64 {
	if x != nil {
		return x.TargetId
	}
	return 0
}

func (x *TemplateJobAccess) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *TemplateJobAccess) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
	}
	return ""
}

func (x *TemplateJobAccess) GetMinimumGrade() int32 {
	if x != nil {
		return x.MinimumGrade
	}
	return 0
}

func (x *TemplateJobAccess) GetJobGradeLabel() string {
	if x != nil && x.JobGradeLabel != nil {
		return *x.JobGradeLabel
	}
	return ""
}

func (x *TemplateJobAccess) GetAccess() AccessLevel {
	if x != nil {
		return x.Access
	}
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

// Dummy - DO NOT USE!
type TemplateUserAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TemplateUserAccess) Reset() {
	*x = TemplateUserAccess{}
	mi := &file_resources_documents_templates_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TemplateUserAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateUserAccess) ProtoMessage() {}

func (x *TemplateUserAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_templates_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateUserAccess.ProtoReflect.Descriptor instead.
func (*TemplateUserAccess) Descriptor() ([]byte, []int) {
	return file_resources_documents_templates_proto_rawDescGZIP(), []int{7}
}

var File_resources_documents_templates_proto protoreflect.FileDescriptor

const file_resources_documents_templates_proto_rawDesc = "" +
	"\n" +
	"#resources/documents/templates.proto\x12\x13resources.documents\x1a resources/documents/access.proto\x1a\"resources/documents/category.proto\x1a#resources/documents/documents.proto\x1a\"resources/documents/workflow.proto\x1a#resources/timestamp/timestamp.proto\x1a\x1bresources/users/users.proto\x1a!resources/vehicles/vehicles.proto\x1a\x17validate/validate.proto\"\xe0\a\n" +
	"\bTemplate\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12B\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12B\n" +
	"\n" +
	"updated_at\x18\x03 \x01(\v2\x1e.resources.timestamp.TimestampH\x01R\tupdatedAt\x88\x01\x01\x129\n" +
	"\bcategory\x18\x04 \x01(\v2\x1d.resources.documents.CategoryR\bcategory\x12#\n" +
	"\x06weight\x18\x05 \x01(\rB\v\xfaB\b*\x06\x10\xff\xff\xff\xff\x0fR\x06weight\x12\x1d\n" +
	"\x05title\x18\x06 \x01(\tB\a\xfaB\x04r\x02\x10\x03R\x05title\x12*\n" +
	"\vdescription\x18\a \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\vdescription\x12$\n" +
	"\x05color\x18\b \x01(\tB\t\xfaB\x06r\x04\x10\x03\x18\aH\x02R\x05color\x88\x01\x01\x12!\n" +
	"\x04icon\x18\t \x01(\tB\b\xfaB\x05r\x03\x18\x80\x01H\x03R\x04icon\x88\x01\x01\x12/\n" +
	"\rcontent_title\x18\n" +
	" \x01(\tB\n" +
	"\xfaB\ar\x05\x10\x03(\x80PR\fcontentTitle\x12%\n" +
	"\acontent\x18\v \x01(\tB\v\xfaB\br\x06\x10\x00(\x80\x89zR\acontent\x12\x1e\n" +
	"\x05state\x18\f \x01(\tB\b\xfaB\x05r\x03\x18\x80\x04R\x05state\x12;\n" +
	"\x06schema\x18\r \x01(\v2#.resources.documents.TemplateSchemaR\x06schema\x12(\n" +
	"\vcreator_job\x18\x0e \x01(\tB\a\xfaB\x04r\x02\x18\x14R\n" +
	"creatorJob\x128\n" +
	"\x11creator_job_label\x18\x0f \x01(\tB\a\xfaB\x04r\x02\x182H\x04R\x0fcreatorJobLabel\x88\x01\x01\x12O\n" +
	"\n" +
	"job_access\x18\x10 \x03(\v2&.resources.documents.TemplateJobAccessB\b\xfaB\x05\x92\x01\x02\x10\x14R\tjobAccess\x12J\n" +
	"\x0econtent_access\x18\x11 \x01(\v2#.resources.documents.DocumentAccessR\rcontentAccess\x12>\n" +
	"\bworkflow\x18\x12 \x01(\v2\x1d.resources.documents.WorkflowH\x05R\bworkflow\x88\x01\x01B\r\n" +
	"\v_created_atB\r\n" +
	"\v_updated_atB\b\n" +
	"\x06_colorB\a\n" +
	"\x05_iconB\x14\n" +
	"\x12_creator_job_labelB\v\n" +
	"\t_workflow\"\xd0\x05\n" +
	"\rTemplateShort\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12B\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12B\n" +
	"\n" +
	"updated_at\x18\x03 \x01(\v2\x1e.resources.timestamp.TimestampH\x01R\tupdatedAt\x88\x01\x01\x129\n" +
	"\bcategory\x18\x04 \x01(\v2\x1d.resources.documents.CategoryR\bcategory\x12#\n" +
	"\x06weight\x18\x05 \x01(\rB\v\xfaB\b*\x06\x10\xff\xff\xff\xff\x0fR\x06weight\x12\x1d\n" +
	"\x05title\x18\x06 \x01(\tB\a\xfaB\x04r\x02\x10\x03R\x05title\x12*\n" +
	"\vdescription\x18\a \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\vdescription\x12$\n" +
	"\x05color\x18\b \x01(\tB\t\xfaB\x06r\x04\x10\x03\x18\aH\x02R\x05color\x88\x01\x01\x12!\n" +
	"\x04icon\x18\t \x01(\tB\b\xfaB\x05r\x03\x18\x80\x01H\x03R\x04icon\x88\x01\x01\x12;\n" +
	"\x06schema\x18\n" +
	" \x01(\v2#.resources.documents.TemplateSchemaR\x06schema\x12(\n" +
	"\vcreator_job\x18\v \x01(\tB\a\xfaB\x04r\x02\x18\x14R\n" +
	"creatorJob\x128\n" +
	"\x11creator_job_label\x18\f \x01(\tB\a\xfaB\x04r\x02\x182H\x04R\x0fcreatorJobLabel\x88\x01\x01\x12>\n" +
	"\bworkflow\x18\x12 \x01(\v2\x1d.resources.documents.WorkflowH\x05R\bworkflow\x88\x01\x01B\r\n" +
	"\v_created_atB\r\n" +
	"\v_updated_atB\b\n" +
	"\x06_colorB\a\n" +
	"\x05_iconB\x14\n" +
	"\x12_creator_job_labelB\v\n" +
	"\t_workflow\"_\n" +
	"\x0eTemplateSchema\x12M\n" +
	"\frequirements\x18\x01 \x01(\v2).resources.documents.TemplateRequirementsR\frequirements\"\x80\x02\n" +
	"\x14TemplateRequirements\x12C\n" +
	"\tdocuments\x18\x01 \x01(\v2 .resources.documents.ObjectSpecsH\x00R\tdocuments\x88\x01\x01\x12;\n" +
	"\x05users\x18\x02 \x01(\v2 .resources.documents.ObjectSpecsH\x01R\x05users\x88\x01\x01\x12A\n" +
	"\bvehicles\x18\x03 \x01(\v2 .resources.documents.ObjectSpecsH\x02R\bvehicles\x88\x01\x01B\f\n" +
	"\n" +
	"_documentsB\b\n" +
	"\x06_usersB\v\n" +
	"\t_vehicles\"y\n" +
	"\vObjectSpecs\x12\x1f\n" +
	"\brequired\x18\x01 \x01(\bH\x00R\brequired\x88\x01\x01\x12\x15\n" +
	"\x03min\x18\x02 \x01(\x05H\x01R\x03min\x88\x01\x01\x12\x15\n" +
	"\x03max\x18\x03 \x01(\x05H\x02R\x03max\x88\x01\x01B\v\n" +
	"\t_requiredB\x06\n" +
	"\x04_minB\x06\n" +
	"\x04_max\"\x9a\x02\n" +
	"\fTemplateData\x12?\n" +
	"\n" +
	"activeChar\x18\x01 \x01(\v2\x15.resources.users.UserB\b\xfaB\x05\x8a\x01\x02\x10\x01R\n" +
	"activeChar\x12J\n" +
	"\tdocuments\x18\x02 \x03(\v2\".resources.documents.DocumentShortB\b\xfaB\x05\x92\x01\x02\x10\fR\tdocuments\x12:\n" +
	"\x05users\x18\x03 \x03(\v2\x1a.resources.users.UserShortB\b\xfaB\x05\x92\x01\x02\x10\fR\x05users\x12A\n" +
	"\bvehicles\x18\x04 \x03(\v2\x1b.resources.vehicles.VehicleB\b\xfaB\x05\x92\x01\x02\x10\fR\bvehicles\"\xa3\x03\n" +
	"\x11TemplateJobAccess\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12B\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12\x1b\n" +
	"\ttarget_id\x18\x03 \x01(\x04R\btargetId\x12\x19\n" +
	"\x03job\x18\x04 \x01(\tB\a\xfaB\x04r\x02\x18\x14R\x03job\x12)\n" +
	"\tjob_label\x18\x05 \x01(\tB\a\xfaB\x04r\x02\x182H\x01R\bjobLabel\x88\x01\x01\x12,\n" +
	"\rminimum_grade\x18\x06 \x01(\x05B\a\xfaB\x04\x1a\x02(\x00R\fminimumGrade\x124\n" +
	"\x0fjob_grade_label\x18\a \x01(\tB\a\xfaB\x04r\x02\x182H\x02R\rjobGradeLabel\x88\x01\x01\x12B\n" +
	"\x06access\x18\b \x01(\x0e2 .resources.documents.AccessLevelB\b\xfaB\x05\x82\x01\x02\x10\x01R\x06accessB\r\n" +
	"\v_created_atB\f\n" +
	"\n" +
	"_job_labelB\x12\n" +
	"\x10_job_grade_label\"\x14\n" +
	"\x12TemplateUserAccessBKZIgithub.com/fivenet-app/fivenet/gen/go/proto/resources/documents;documentsb\x06proto3"

var (
	file_resources_documents_templates_proto_rawDescOnce sync.Once
	file_resources_documents_templates_proto_rawDescData []byte
)

func file_resources_documents_templates_proto_rawDescGZIP() []byte {
	file_resources_documents_templates_proto_rawDescOnce.Do(func() {
		file_resources_documents_templates_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_documents_templates_proto_rawDesc), len(file_resources_documents_templates_proto_rawDesc)))
	})
	return file_resources_documents_templates_proto_rawDescData
}

var file_resources_documents_templates_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_resources_documents_templates_proto_goTypes = []any{
	(*Template)(nil),             // 0: resources.documents.Template
	(*TemplateShort)(nil),        // 1: resources.documents.TemplateShort
	(*TemplateSchema)(nil),       // 2: resources.documents.TemplateSchema
	(*TemplateRequirements)(nil), // 3: resources.documents.TemplateRequirements
	(*ObjectSpecs)(nil),          // 4: resources.documents.ObjectSpecs
	(*TemplateData)(nil),         // 5: resources.documents.TemplateData
	(*TemplateJobAccess)(nil),    // 6: resources.documents.TemplateJobAccess
	(*TemplateUserAccess)(nil),   // 7: resources.documents.TemplateUserAccess
	(*timestamp.Timestamp)(nil),  // 8: resources.timestamp.Timestamp
	(*Category)(nil),             // 9: resources.documents.Category
	(*DocumentAccess)(nil),       // 10: resources.documents.DocumentAccess
	(*Workflow)(nil),             // 11: resources.documents.Workflow
	(*users.User)(nil),           // 12: resources.users.User
	(*DocumentShort)(nil),        // 13: resources.documents.DocumentShort
	(*users.UserShort)(nil),      // 14: resources.users.UserShort
	(*vehicles.Vehicle)(nil),     // 15: resources.vehicles.Vehicle
	(AccessLevel)(0),             // 16: resources.documents.AccessLevel
}
var file_resources_documents_templates_proto_depIdxs = []int32{
	8,  // 0: resources.documents.Template.created_at:type_name -> resources.timestamp.Timestamp
	8,  // 1: resources.documents.Template.updated_at:type_name -> resources.timestamp.Timestamp
	9,  // 2: resources.documents.Template.category:type_name -> resources.documents.Category
	2,  // 3: resources.documents.Template.schema:type_name -> resources.documents.TemplateSchema
	6,  // 4: resources.documents.Template.job_access:type_name -> resources.documents.TemplateJobAccess
	10, // 5: resources.documents.Template.content_access:type_name -> resources.documents.DocumentAccess
	11, // 6: resources.documents.Template.workflow:type_name -> resources.documents.Workflow
	8,  // 7: resources.documents.TemplateShort.created_at:type_name -> resources.timestamp.Timestamp
	8,  // 8: resources.documents.TemplateShort.updated_at:type_name -> resources.timestamp.Timestamp
	9,  // 9: resources.documents.TemplateShort.category:type_name -> resources.documents.Category
	2,  // 10: resources.documents.TemplateShort.schema:type_name -> resources.documents.TemplateSchema
	11, // 11: resources.documents.TemplateShort.workflow:type_name -> resources.documents.Workflow
	3,  // 12: resources.documents.TemplateSchema.requirements:type_name -> resources.documents.TemplateRequirements
	4,  // 13: resources.documents.TemplateRequirements.documents:type_name -> resources.documents.ObjectSpecs
	4,  // 14: resources.documents.TemplateRequirements.users:type_name -> resources.documents.ObjectSpecs
	4,  // 15: resources.documents.TemplateRequirements.vehicles:type_name -> resources.documents.ObjectSpecs
	12, // 16: resources.documents.TemplateData.activeChar:type_name -> resources.users.User
	13, // 17: resources.documents.TemplateData.documents:type_name -> resources.documents.DocumentShort
	14, // 18: resources.documents.TemplateData.users:type_name -> resources.users.UserShort
	15, // 19: resources.documents.TemplateData.vehicles:type_name -> resources.vehicles.Vehicle
	8,  // 20: resources.documents.TemplateJobAccess.created_at:type_name -> resources.timestamp.Timestamp
	16, // 21: resources.documents.TemplateJobAccess.access:type_name -> resources.documents.AccessLevel
	22, // [22:22] is the sub-list for method output_type
	22, // [22:22] is the sub-list for method input_type
	22, // [22:22] is the sub-list for extension type_name
	22, // [22:22] is the sub-list for extension extendee
	0,  // [0:22] is the sub-list for field type_name
}

func init() { file_resources_documents_templates_proto_init() }
func file_resources_documents_templates_proto_init() {
	if File_resources_documents_templates_proto != nil {
		return
	}
	file_resources_documents_access_proto_init()
	file_resources_documents_category_proto_init()
	file_resources_documents_documents_proto_init()
	file_resources_documents_workflow_proto_init()
	file_resources_documents_templates_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_documents_templates_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_documents_templates_proto_msgTypes[3].OneofWrappers = []any{}
	file_resources_documents_templates_proto_msgTypes[4].OneofWrappers = []any{}
	file_resources_documents_templates_proto_msgTypes[6].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_documents_templates_proto_rawDesc), len(file_resources_documents_templates_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_documents_templates_proto_goTypes,
		DependencyIndexes: file_resources_documents_templates_proto_depIdxs,
		MessageInfos:      file_resources_documents_templates_proto_msgTypes,
	}.Build()
	File_resources_documents_templates_proto = out.File
	file_resources_documents_templates_proto_goTypes = nil
	file_resources_documents_templates_proto_depIdxs = nil
}
