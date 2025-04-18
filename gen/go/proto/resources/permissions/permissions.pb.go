// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: resources/permissions/permissions.proto

package permissions

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type Permission struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	Category      string                 `protobuf:"bytes,3,opt,name=category,proto3" json:"category,omitempty"`
	Name          string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	GuardName     string                 `protobuf:"bytes,5,opt,name=guard_name,json=guardName,proto3" json:"guard_name,omitempty"`
	Val           bool                   `protobuf:"varint,6,opt,name=val,proto3" json:"val,omitempty"`
	Order         *int32                 `protobuf:"varint,7,opt,name=order,proto3,oneof" json:"order,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Permission) Reset() {
	*x = Permission{}
	mi := &file_resources_permissions_permissions_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Permission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Permission) ProtoMessage() {}

func (x *Permission) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Permission.ProtoReflect.Descriptor instead.
func (*Permission) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{0}
}

func (x *Permission) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Permission) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Permission) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Permission) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Permission) GetGuardName() string {
	if x != nil {
		return x.GuardName
	}
	return ""
}

func (x *Permission) GetVal() bool {
	if x != nil {
		return x.Val
	}
	return false
}

func (x *Permission) GetOrder() int32 {
	if x != nil && x.Order != nil {
		return *x.Order
	}
	return 0
}

type Role struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	Job           string                 `protobuf:"bytes,3,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel      *string                `protobuf:"bytes,4,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty"`
	Grade         int32                  `protobuf:"varint,5,opt,name=grade,proto3" json:"grade,omitempty"`
	JobGradeLabel *string                `protobuf:"bytes,6,opt,name=job_grade_label,json=jobGradeLabel,proto3,oneof" json:"job_grade_label,omitempty"`
	Permissions   []*Permission          `protobuf:"bytes,7,rep,name=permissions,proto3" json:"permissions,omitempty"`
	Attributes    []*RoleAttribute       `protobuf:"bytes,8,rep,name=attributes,proto3" json:"attributes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Role) Reset() {
	*x = Role{}
	mi := &file_resources_permissions_permissions_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Role) ProtoMessage() {}

func (x *Role) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Role.ProtoReflect.Descriptor instead.
func (*Role) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{1}
}

func (x *Role) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Role) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Role) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *Role) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
	}
	return ""
}

func (x *Role) GetGrade() int32 {
	if x != nil {
		return x.Grade
	}
	return 0
}

func (x *Role) GetJobGradeLabel() string {
	if x != nil && x.JobGradeLabel != nil {
		return *x.JobGradeLabel
	}
	return ""
}

func (x *Role) GetPermissions() []*Permission {
	if x != nil {
		return x.Permissions
	}
	return nil
}

func (x *Role) GetAttributes() []*RoleAttribute {
	if x != nil {
		return x.Attributes
	}
	return nil
}

type RawRoleAttribute struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RoleId        uint64                 `protobuf:"varint,1,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	AttrId        uint64                 `protobuf:"varint,3,opt,name=attr_id,json=attrId,proto3" json:"attr_id,omitempty"`
	PermissionId  uint64                 `protobuf:"varint,4,opt,name=permission_id,json=permissionId,proto3" json:"permission_id,omitempty"`
	Category      string                 `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty"`
	Name          string                 `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	Key           string                 `protobuf:"bytes,7,opt,name=key,proto3" json:"key,omitempty"`
	Type          string                 `protobuf:"bytes,8,opt,name=type,proto3" json:"type,omitempty"`
	ValidValues   *AttributeValues       `protobuf:"bytes,9,opt,name=valid_values,json=validValues,proto3" json:"valid_values,omitempty"`
	Value         *AttributeValues       `protobuf:"bytes,10,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawRoleAttribute) Reset() {
	*x = RawRoleAttribute{}
	mi := &file_resources_permissions_permissions_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawRoleAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawRoleAttribute) ProtoMessage() {}

func (x *RawRoleAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawRoleAttribute.ProtoReflect.Descriptor instead.
func (*RawRoleAttribute) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{2}
}

func (x *RawRoleAttribute) GetRoleId() uint64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

func (x *RawRoleAttribute) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *RawRoleAttribute) GetAttrId() uint64 {
	if x != nil {
		return x.AttrId
	}
	return 0
}

func (x *RawRoleAttribute) GetPermissionId() uint64 {
	if x != nil {
		return x.PermissionId
	}
	return 0
}

func (x *RawRoleAttribute) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *RawRoleAttribute) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RawRoleAttribute) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *RawRoleAttribute) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *RawRoleAttribute) GetValidValues() *AttributeValues {
	if x != nil {
		return x.ValidValues
	}
	return nil
}

func (x *RawRoleAttribute) GetValue() *AttributeValues {
	if x != nil {
		return x.Value
	}
	return nil
}

type RoleAttribute struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RoleId        uint64                 `protobuf:"varint,1,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	AttrId        uint64                 `protobuf:"varint,3,opt,name=attr_id,json=attrId,proto3" json:"attr_id,omitempty"`
	PermissionId  uint64                 `protobuf:"varint,4,opt,name=permission_id,json=permissionId,proto3" json:"permission_id,omitempty"`
	Category      string                 `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty"`
	Name          string                 `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	Key           string                 `protobuf:"bytes,7,opt,name=key,proto3" json:"key,omitempty"`
	Type          string                 `protobuf:"bytes,8,opt,name=type,proto3" json:"type,omitempty"`
	ValidValues   *AttributeValues       `protobuf:"bytes,9,opt,name=valid_values,json=validValues,proto3" json:"valid_values,omitempty"`
	Value         *AttributeValues       `protobuf:"bytes,10,opt,name=value,proto3" json:"value,omitempty"`
	MaxValues     *AttributeValues       `protobuf:"bytes,11,opt,name=max_values,json=maxValues,proto3,oneof" json:"max_values,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RoleAttribute) Reset() {
	*x = RoleAttribute{}
	mi := &file_resources_permissions_permissions_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RoleAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleAttribute) ProtoMessage() {}

func (x *RoleAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleAttribute.ProtoReflect.Descriptor instead.
func (*RoleAttribute) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{3}
}

func (x *RoleAttribute) GetRoleId() uint64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

func (x *RoleAttribute) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *RoleAttribute) GetAttrId() uint64 {
	if x != nil {
		return x.AttrId
	}
	return 0
}

func (x *RoleAttribute) GetPermissionId() uint64 {
	if x != nil {
		return x.PermissionId
	}
	return 0
}

func (x *RoleAttribute) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *RoleAttribute) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RoleAttribute) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *RoleAttribute) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *RoleAttribute) GetValidValues() *AttributeValues {
	if x != nil {
		return x.ValidValues
	}
	return nil
}

func (x *RoleAttribute) GetValue() *AttributeValues {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *RoleAttribute) GetMaxValues() *AttributeValues {
	if x != nil {
		return x.MaxValues
	}
	return nil
}

// @dbscanner: json
type AttributeValues struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to ValidValues:
	//
	//	*AttributeValues_StringList
	//	*AttributeValues_JobList
	//	*AttributeValues_JobGradeList
	//	*AttributeValues_JobGradeMap
	ValidValues   isAttributeValues_ValidValues `protobuf_oneof:"valid_values"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AttributeValues) Reset() {
	*x = AttributeValues{}
	mi := &file_resources_permissions_permissions_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AttributeValues) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttributeValues) ProtoMessage() {}

func (x *AttributeValues) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttributeValues.ProtoReflect.Descriptor instead.
func (*AttributeValues) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{4}
}

func (x *AttributeValues) GetValidValues() isAttributeValues_ValidValues {
	if x != nil {
		return x.ValidValues
	}
	return nil
}

func (x *AttributeValues) GetStringList() *StringList {
	if x != nil {
		if x, ok := x.ValidValues.(*AttributeValues_StringList); ok {
			return x.StringList
		}
	}
	return nil
}

func (x *AttributeValues) GetJobList() *StringList {
	if x != nil {
		if x, ok := x.ValidValues.(*AttributeValues_JobList); ok {
			return x.JobList
		}
	}
	return nil
}

func (x *AttributeValues) GetJobGradeList() *JobGradeList {
	if x != nil {
		if x, ok := x.ValidValues.(*AttributeValues_JobGradeList); ok {
			return x.JobGradeList
		}
	}
	return nil
}

func (x *AttributeValues) GetJobGradeMap() *JobGradeMap {
	if x != nil {
		if x, ok := x.ValidValues.(*AttributeValues_JobGradeMap); ok {
			return x.JobGradeMap
		}
	}
	return nil
}

type isAttributeValues_ValidValues interface {
	isAttributeValues_ValidValues()
}

type AttributeValues_StringList struct {
	StringList *StringList `protobuf:"bytes,1,opt,name=string_list,json=stringList,proto3,oneof"`
}

type AttributeValues_JobList struct {
	JobList *StringList `protobuf:"bytes,2,opt,name=job_list,json=jobList,proto3,oneof"`
}

type AttributeValues_JobGradeList struct {
	JobGradeList *JobGradeList `protobuf:"bytes,3,opt,name=job_grade_list,json=jobGradeList,proto3,oneof"`
}

type AttributeValues_JobGradeMap struct {
	JobGradeMap *JobGradeMap `protobuf:"bytes,4,opt,name=job_grade_map,json=jobGradeMap,proto3,oneof"`
}

func (*AttributeValues_StringList) isAttributeValues_ValidValues() {}

func (*AttributeValues_JobList) isAttributeValues_ValidValues() {}

func (*AttributeValues_JobGradeList) isAttributeValues_ValidValues() {}

func (*AttributeValues_JobGradeMap) isAttributeValues_ValidValues() {}

type StringList struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// @sanitize: method=StripTags
	Strings       []string `protobuf:"bytes,1,rep,name=strings,proto3" json:"strings,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StringList) Reset() {
	*x = StringList{}
	mi := &file_resources_permissions_permissions_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StringList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringList) ProtoMessage() {}

func (x *StringList) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringList.ProtoReflect.Descriptor instead.
func (*StringList) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{5}
}

func (x *StringList) GetStrings() []string {
	if x != nil {
		return x.Strings
	}
	return nil
}

type JobGradeList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Jobs          map[string]int32       `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JobGradeList) Reset() {
	*x = JobGradeList{}
	mi := &file_resources_permissions_permissions_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobGradeList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobGradeList) ProtoMessage() {}

func (x *JobGradeList) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobGradeList.ProtoReflect.Descriptor instead.
func (*JobGradeList) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{6}
}

func (x *JobGradeList) GetJobs() map[string]int32 {
	if x != nil {
		return x.Jobs
	}
	return nil
}

type JobGradeMap struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Jobs          map[string]*JobGrades  `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JobGradeMap) Reset() {
	*x = JobGradeMap{}
	mi := &file_resources_permissions_permissions_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobGradeMap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobGradeMap) ProtoMessage() {}

func (x *JobGradeMap) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobGradeMap.ProtoReflect.Descriptor instead.
func (*JobGradeMap) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{7}
}

func (x *JobGradeMap) GetJobs() map[string]*JobGrades {
	if x != nil {
		return x.Jobs
	}
	return nil
}

type JobGrades struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Grades        []int32                `protobuf:"varint,1,rep,packed,name=grades,proto3" json:"grades,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JobGrades) Reset() {
	*x = JobGrades{}
	mi := &file_resources_permissions_permissions_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobGrades) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobGrades) ProtoMessage() {}

func (x *JobGrades) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobGrades.ProtoReflect.Descriptor instead.
func (*JobGrades) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{8}
}

func (x *JobGrades) GetGrades() []int32 {
	if x != nil {
		return x.Grades
	}
	return nil
}

var File_resources_permissions_permissions_proto protoreflect.FileDescriptor

const file_resources_permissions_permissions_proto_rawDesc = "" +
	"\n" +
	"'resources/permissions/permissions.proto\x12\x15resources.permissions\x1a#resources/timestamp/timestamp.proto\x1a\x17validate/validate.proto\"\x9c\x02\n" +
	"\n" +
	"Permission\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12B\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12$\n" +
	"\bcategory\x18\x03 \x01(\tB\b\xfaB\x05r\x03\x18\x80\x01R\bcategory\x12\x1c\n" +
	"\x04name\x18\x04 \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\x04name\x12'\n" +
	"\n" +
	"guard_name\x18\x05 \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\tguardName\x12\x10\n" +
	"\x03val\x18\x06 \x01(\bR\x03val\x12\"\n" +
	"\x05order\x18\a \x01(\x05B\a\xfaB\x04\x1a\x02(\x00H\x01R\x05order\x88\x01\x01B\r\n" +
	"\v_created_atB\b\n" +
	"\x06_order\"\xb2\x03\n" +
	"\x04Role\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12B\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12\x1a\n" +
	"\x03job\x18\x03 \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\x03job\x12)\n" +
	"\tjob_label\x18\x04 \x01(\tB\a\xfaB\x04r\x02\x182H\x01R\bjobLabel\x88\x01\x01\x12\x1d\n" +
	"\x05grade\x18\x05 \x01(\x05B\a\xfaB\x04\x1a\x02(\x00R\x05grade\x124\n" +
	"\x0fjob_grade_label\x18\x06 \x01(\tB\a\xfaB\x04r\x02\x182H\x02R\rjobGradeLabel\x88\x01\x01\x12C\n" +
	"\vpermissions\x18\a \x03(\v2!.resources.permissions.PermissionR\vpermissions\x12D\n" +
	"\n" +
	"attributes\x18\b \x03(\v2$.resources.permissions.RoleAttributeR\n" +
	"attributesB\r\n" +
	"\v_created_atB\f\n" +
	"\n" +
	"_job_labelB\x12\n" +
	"\x10_job_grade_label\"\xc3\x03\n" +
	"\x10RawRoleAttribute\x12\x17\n" +
	"\arole_id\x18\x01 \x01(\x04R\x06roleId\x12B\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12\x17\n" +
	"\aattr_id\x18\x03 \x01(\x04R\x06attrId\x12#\n" +
	"\rpermission_id\x18\x04 \x01(\x04R\fpermissionId\x12$\n" +
	"\bcategory\x18\x05 \x01(\tB\b\xfaB\x05r\x03\x18\x80\x01R\bcategory\x12\x1c\n" +
	"\x04name\x18\x06 \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\x04name\x12\x1a\n" +
	"\x03key\x18\a \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\x03key\x12\x1c\n" +
	"\x04type\x18\b \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\x04type\x12I\n" +
	"\fvalid_values\x18\t \x01(\v2&.resources.permissions.AttributeValuesR\vvalidValues\x12<\n" +
	"\x05value\x18\n" +
	" \x01(\v2&.resources.permissions.AttributeValuesR\x05valueB\r\n" +
	"\v_created_at\"\x9b\x04\n" +
	"\rRoleAttribute\x12\x17\n" +
	"\arole_id\x18\x01 \x01(\x04R\x06roleId\x12B\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12\x17\n" +
	"\aattr_id\x18\x03 \x01(\x04R\x06attrId\x12#\n" +
	"\rpermission_id\x18\x04 \x01(\x04R\fpermissionId\x12$\n" +
	"\bcategory\x18\x05 \x01(\tB\b\xfaB\x05r\x03\x18\x80\x01R\bcategory\x12\x1c\n" +
	"\x04name\x18\x06 \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\x04name\x12\x1a\n" +
	"\x03key\x18\a \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\x03key\x12\x1c\n" +
	"\x04type\x18\b \x01(\tB\b\xfaB\x05r\x03\x18\xff\x01R\x04type\x12I\n" +
	"\fvalid_values\x18\t \x01(\v2&.resources.permissions.AttributeValuesR\vvalidValues\x12<\n" +
	"\x05value\x18\n" +
	" \x01(\v2&.resources.permissions.AttributeValuesR\x05value\x12J\n" +
	"\n" +
	"max_values\x18\v \x01(\v2&.resources.permissions.AttributeValuesH\x01R\tmaxValues\x88\x01\x01B\r\n" +
	"\v_created_atB\r\n" +
	"\v_max_values\"\xc3\x02\n" +
	"\x0fAttributeValues\x12D\n" +
	"\vstring_list\x18\x01 \x01(\v2!.resources.permissions.StringListH\x00R\n" +
	"stringList\x12>\n" +
	"\bjob_list\x18\x02 \x01(\v2!.resources.permissions.StringListH\x00R\ajobList\x12K\n" +
	"\x0ejob_grade_list\x18\x03 \x01(\v2#.resources.permissions.JobGradeListH\x00R\fjobGradeList\x12H\n" +
	"\rjob_grade_map\x18\x04 \x01(\v2\".resources.permissions.JobGradeMapH\x00R\vjobGradeMapB\x13\n" +
	"\fvalid_values\x12\x03\xf8B\x01\"&\n" +
	"\n" +
	"StringList\x12\x18\n" +
	"\astrings\x18\x01 \x03(\tR\astrings\"\x8a\x01\n" +
	"\fJobGradeList\x12A\n" +
	"\x04jobs\x18\x01 \x03(\v2-.resources.permissions.JobGradeList.JobsEntryR\x04jobs\x1a7\n" +
	"\tJobsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\x05R\x05value:\x028\x01\"\xaa\x01\n" +
	"\vJobGradeMap\x12@\n" +
	"\x04jobs\x18\x01 \x03(\v2,.resources.permissions.JobGradeMap.JobsEntryR\x04jobs\x1aY\n" +
	"\tJobsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x126\n" +
	"\x05value\x18\x02 \x01(\v2 .resources.permissions.JobGradesR\x05value:\x028\x01\"#\n" +
	"\tJobGrades\x12\x16\n" +
	"\x06grades\x18\x01 \x03(\x05R\x06gradesBOZMgithub.com/fivenet-app/fivenet/gen/go/proto/resources/permissions;permissionsb\x06proto3"

var (
	file_resources_permissions_permissions_proto_rawDescOnce sync.Once
	file_resources_permissions_permissions_proto_rawDescData []byte
)

func file_resources_permissions_permissions_proto_rawDescGZIP() []byte {
	file_resources_permissions_permissions_proto_rawDescOnce.Do(func() {
		file_resources_permissions_permissions_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_permissions_permissions_proto_rawDesc), len(file_resources_permissions_permissions_proto_rawDesc)))
	})
	return file_resources_permissions_permissions_proto_rawDescData
}

var file_resources_permissions_permissions_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_resources_permissions_permissions_proto_goTypes = []any{
	(*Permission)(nil),          // 0: resources.permissions.Permission
	(*Role)(nil),                // 1: resources.permissions.Role
	(*RawRoleAttribute)(nil),    // 2: resources.permissions.RawRoleAttribute
	(*RoleAttribute)(nil),       // 3: resources.permissions.RoleAttribute
	(*AttributeValues)(nil),     // 4: resources.permissions.AttributeValues
	(*StringList)(nil),          // 5: resources.permissions.StringList
	(*JobGradeList)(nil),        // 6: resources.permissions.JobGradeList
	(*JobGradeMap)(nil),         // 7: resources.permissions.JobGradeMap
	(*JobGrades)(nil),           // 8: resources.permissions.JobGrades
	nil,                         // 9: resources.permissions.JobGradeList.JobsEntry
	nil,                         // 10: resources.permissions.JobGradeMap.JobsEntry
	(*timestamp.Timestamp)(nil), // 11: resources.timestamp.Timestamp
}
var file_resources_permissions_permissions_proto_depIdxs = []int32{
	11, // 0: resources.permissions.Permission.created_at:type_name -> resources.timestamp.Timestamp
	11, // 1: resources.permissions.Role.created_at:type_name -> resources.timestamp.Timestamp
	0,  // 2: resources.permissions.Role.permissions:type_name -> resources.permissions.Permission
	3,  // 3: resources.permissions.Role.attributes:type_name -> resources.permissions.RoleAttribute
	11, // 4: resources.permissions.RawRoleAttribute.created_at:type_name -> resources.timestamp.Timestamp
	4,  // 5: resources.permissions.RawRoleAttribute.valid_values:type_name -> resources.permissions.AttributeValues
	4,  // 6: resources.permissions.RawRoleAttribute.value:type_name -> resources.permissions.AttributeValues
	11, // 7: resources.permissions.RoleAttribute.created_at:type_name -> resources.timestamp.Timestamp
	4,  // 8: resources.permissions.RoleAttribute.valid_values:type_name -> resources.permissions.AttributeValues
	4,  // 9: resources.permissions.RoleAttribute.value:type_name -> resources.permissions.AttributeValues
	4,  // 10: resources.permissions.RoleAttribute.max_values:type_name -> resources.permissions.AttributeValues
	5,  // 11: resources.permissions.AttributeValues.string_list:type_name -> resources.permissions.StringList
	5,  // 12: resources.permissions.AttributeValues.job_list:type_name -> resources.permissions.StringList
	6,  // 13: resources.permissions.AttributeValues.job_grade_list:type_name -> resources.permissions.JobGradeList
	7,  // 14: resources.permissions.AttributeValues.job_grade_map:type_name -> resources.permissions.JobGradeMap
	9,  // 15: resources.permissions.JobGradeList.jobs:type_name -> resources.permissions.JobGradeList.JobsEntry
	10, // 16: resources.permissions.JobGradeMap.jobs:type_name -> resources.permissions.JobGradeMap.JobsEntry
	8,  // 17: resources.permissions.JobGradeMap.JobsEntry.value:type_name -> resources.permissions.JobGrades
	18, // [18:18] is the sub-list for method output_type
	18, // [18:18] is the sub-list for method input_type
	18, // [18:18] is the sub-list for extension type_name
	18, // [18:18] is the sub-list for extension extendee
	0,  // [0:18] is the sub-list for field type_name
}

func init() { file_resources_permissions_permissions_proto_init() }
func file_resources_permissions_permissions_proto_init() {
	if File_resources_permissions_permissions_proto != nil {
		return
	}
	file_resources_permissions_permissions_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_permissions_permissions_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_permissions_permissions_proto_msgTypes[2].OneofWrappers = []any{}
	file_resources_permissions_permissions_proto_msgTypes[3].OneofWrappers = []any{}
	file_resources_permissions_permissions_proto_msgTypes[4].OneofWrappers = []any{
		(*AttributeValues_StringList)(nil),
		(*AttributeValues_JobList)(nil),
		(*AttributeValues_JobGradeList)(nil),
		(*AttributeValues_JobGradeMap)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_permissions_permissions_proto_rawDesc), len(file_resources_permissions_permissions_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_permissions_permissions_proto_goTypes,
		DependencyIndexes: file_resources_permissions_permissions_proto_depIdxs,
		MessageInfos:      file_resources_permissions_permissions_proto_msgTypes,
	}.Build()
	File_resources_permissions_permissions_proto = out.File
	file_resources_permissions_permissions_proto_goTypes = nil
	file_resources_permissions_permissions_proto_depIdxs = nil
}
