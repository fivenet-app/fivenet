// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: resources/permissions/permissions.proto

package permissions

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

type Permission struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" alias:"id"`                                     // @gotags: alias:"id"
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty" alias:"created_at"` // @gotags: alias:"created_at"
	Category  string               `protobuf:"bytes,3,opt,name=category,proto3" json:"category,omitempty" alias:"category"`                          // @gotags: alias:"category"
	Name      string               `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty" alias:"name"`                                  // @gotags: alias:"name"
	GuardName string               `protobuf:"bytes,5,opt,name=guard_name,json=guardName,proto3" json:"guard_name,omitempty" alias:"guard_name"`       // @gotags: alias:"guard_name"
	Val       int32                `protobuf:"varint,6,opt,name=val,proto3" json:"val,omitempty" alias:"val"`                                   // @gotags: alias:"val"
}

func (x *Permission) Reset() {
	*x = Permission{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_permissions_permissions_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Permission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Permission) ProtoMessage() {}

func (x *Permission) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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

func (x *Permission) GetVal() int32 {
	if x != nil {
		return x.Val
	}
	return 0
}

type Role struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" alias:"id"`                                             // @gotags: alias:"id"
	CreatedAt     *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty" alias:"created_at"`         // @gotags: alias:"created_at"
	Job           string               `protobuf:"bytes,3,opt,name=job,proto3" json:"job,omitempty" alias:"job"`                                            // @gotags: alias:"job"
	JobLabel      string               `protobuf:"bytes,4,opt,name=job_label,json=jobLabel,proto3" json:"job_label,omitempty" alias:"job_label"`                  // @gotags: alias:"job_label"
	Grade         int32                `protobuf:"varint,5,opt,name=grade,proto3" json:"grade,omitempty" alias:"grade"`                                       // @gotags: alias:"grade"
	JobGradeLabel string               `protobuf:"bytes,6,opt,name=job_grade_label,json=jobGradeLabel,proto3" json:"job_grade_label,omitempty" alias:"job_grade_label"` // @gotags: alias:"job_grade_label"
	Permissions   []*Permission        `protobuf:"bytes,7,rep,name=permissions,proto3" json:"permissions,omitempty"`
	Attributes    []*RoleAttribute     `protobuf:"bytes,8,rep,name=attributes,proto3" json:"attributes,omitempty"`
}

func (x *Role) Reset() {
	*x = Role{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_permissions_permissions_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Role) ProtoMessage() {}

func (x *Role) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
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
	if x != nil {
		return x.JobLabel
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
	if x != nil {
		return x.JobGradeLabel
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

type Attribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" alias:"id"`                                         // @gotags: alias:"id"
	CreatedAt    *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty" alias:"created_at"`     // @gotags: alias:"created_at"
	PermissionId uint64               `protobuf:"varint,3,opt,name=permission_id,json=permissionId,proto3" json:"permission_id,omitempty" alias:"permission_id"` // @gotags: alias:"permission_id"
	Key          string               `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty" alias:"key"`                                        // @gotags: alias:"key"
	Type         string               `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty" alias:"type"`                                      // @gotags: alias:"type"
	Value        string               `protobuf:"bytes,6,opt,name=value,proto3" json:"value,omitempty" alias:"value"`                                    // @gotags: alias:"value"
	ValidValues  []string             `protobuf:"bytes,7,rep,name=validValues,proto3" json:"validValues,omitempty"`
}

func (x *Attribute) Reset() {
	*x = Attribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_permissions_permissions_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attribute) ProtoMessage() {}

func (x *Attribute) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attribute.ProtoReflect.Descriptor instead.
func (*Attribute) Descriptor() ([]byte, []int) {
	return file_resources_permissions_permissions_proto_rawDescGZIP(), []int{2}
}

func (x *Attribute) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Attribute) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Attribute) GetPermissionId() uint64 {
	if x != nil {
		return x.PermissionId
	}
	return 0
}

func (x *Attribute) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Attribute) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Attribute) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Attribute) GetValidValues() []string {
	if x != nil {
		return x.ValidValues
	}
	return nil
}

type RoleAttribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId       uint64               `protobuf:"varint,1,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty" alias:"role_id"`                   // @gotags: alias:"role_id"
	CreatedAt    *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty" alias:"created_at"`     // @gotags: alias:"created_at"
	AttrId       uint64               `protobuf:"varint,3,opt,name=attr_id,json=attrId,proto3" json:"attr_id,omitempty" alias:"attr_id"`                   // @gotags: alias:"attr_id"
	PermissionId uint64               `protobuf:"varint,4,opt,name=permission_id,json=permissionId,proto3" json:"permission_id,omitempty" alias:"permission_id"` // @gotags: alias:"permission_id"
	Category     string               `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty" alias:"category"`                              // @gotags: alias:"category"
	Name         string               `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty" alias:"name"`                                      // @gotags: alias:"name"
	Key          string               `protobuf:"bytes,7,opt,name=key,proto3" json:"key,omitempty" alias:"key"`                                        // @gotags: alias:"key"
	Type         string               `protobuf:"bytes,8,opt,name=type,proto3" json:"type,omitempty" alias:"type"`                                      // @gotags: alias:"type"
	Value        string               `protobuf:"bytes,9,opt,name=value,proto3" json:"value,omitempty" alias:"value"`                                    // @gotags: alias:"value"
	ValidValues  string               `protobuf:"bytes,10,opt,name=valid_values,json=validValues,proto3" json:"valid_values,omitempty" alias:"valid_values"`    // @gotags: alias:"valid_values"
}

func (x *RoleAttribute) Reset() {
	*x = RoleAttribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_permissions_permissions_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleAttribute) ProtoMessage() {}

func (x *RoleAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_resources_permissions_permissions_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
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

func (x *RoleAttribute) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *RoleAttribute) GetValidValues() string {
	if x != nil {
		return x.ValidValues
	}
	return ""
}

var File_resources_permissions_permissions_proto protoreflect.FileDescriptor

var file_resources_permissions_permissions_proto_rawDesc = []byte{
	0x0a, 0x27, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x82,
	0x02, 0x0a, 0x0a, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01,
	0x01, 0x12, 0x24, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x01, 0x52, 0x08, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x1c, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0a, 0x67, 0x75, 0x61, 0x72, 0x64, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03,
	0x18, 0xff, 0x01, 0x52, 0x09, 0x67, 0x75, 0x61, 0x72, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24,
	0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x42, 0x12, 0xfa, 0x42, 0x0f,
	0x1a, 0x0d, 0x18, 0x01, 0x28, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x52,
	0x03, 0x76, 0x61, 0x6c, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x22, 0x86, 0x03, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01,
	0x12, 0x1a, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x24, 0x0a, 0x09,
	0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x52, 0x08, 0x6a, 0x6f, 0x62, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x12, 0x1d, 0x0a, 0x05, 0x67, 0x72, 0x61, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x05, 0x67, 0x72, 0x61, 0x64,
	0x65, 0x12, 0x2f, 0x0a, 0x0f, 0x6a, 0x6f, 0x62, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72,
	0x02, 0x18, 0x32, 0x52, 0x0d, 0x6a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64, 0x65, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x12, 0x43, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x44, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x42, 0x0d, 0x0a,
	0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x22, 0x90, 0x02, 0x0a,
	0x09, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x23,
	0x0a, 0x0d, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x1c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xfa, 0x42,
	0x06, 0x72, 0x04, 0x18, 0xff, 0xff, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x07, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0b, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x22,
	0x86, 0x03, 0x0a, 0x0d, 0x52, 0x6f, 0x6c, 0x65, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x42, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x17,
	0x0a, 0x07, 0x61, 0x74, 0x74, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x61, 0x74, 0x74, 0x72, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x65, 0x72, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c,
	0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x08,
	0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x01, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x12, 0x1c, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1a, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72,
	0x03, 0x18, 0xff, 0x01, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04,
	0x18, 0xff, 0xff, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2c, 0x0a, 0x0c, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x18, 0xff, 0xff, 0x03, 0x52, 0x0b, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x4b, 0x5a, 0x49, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x66,
	0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x70, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x3b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_permissions_permissions_proto_rawDescOnce sync.Once
	file_resources_permissions_permissions_proto_rawDescData = file_resources_permissions_permissions_proto_rawDesc
)

func file_resources_permissions_permissions_proto_rawDescGZIP() []byte {
	file_resources_permissions_permissions_proto_rawDescOnce.Do(func() {
		file_resources_permissions_permissions_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_permissions_permissions_proto_rawDescData)
	})
	return file_resources_permissions_permissions_proto_rawDescData
}

var file_resources_permissions_permissions_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_resources_permissions_permissions_proto_goTypes = []interface{}{
	(*Permission)(nil),          // 0: resources.permissions.Permission
	(*Role)(nil),                // 1: resources.permissions.Role
	(*Attribute)(nil),           // 2: resources.permissions.Attribute
	(*RoleAttribute)(nil),       // 3: resources.permissions.RoleAttribute
	(*timestamp.Timestamp)(nil), // 4: resources.timestamp.Timestamp
}
var file_resources_permissions_permissions_proto_depIdxs = []int32{
	4, // 0: resources.permissions.Permission.created_at:type_name -> resources.timestamp.Timestamp
	4, // 1: resources.permissions.Role.created_at:type_name -> resources.timestamp.Timestamp
	0, // 2: resources.permissions.Role.permissions:type_name -> resources.permissions.Permission
	3, // 3: resources.permissions.Role.attributes:type_name -> resources.permissions.RoleAttribute
	4, // 4: resources.permissions.Attribute.created_at:type_name -> resources.timestamp.Timestamp
	4, // 5: resources.permissions.RoleAttribute.created_at:type_name -> resources.timestamp.Timestamp
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_resources_permissions_permissions_proto_init() }
func file_resources_permissions_permissions_proto_init() {
	if File_resources_permissions_permissions_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_permissions_permissions_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Permission); i {
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
		file_resources_permissions_permissions_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Role); i {
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
		file_resources_permissions_permissions_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attribute); i {
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
		file_resources_permissions_permissions_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoleAttribute); i {
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
	file_resources_permissions_permissions_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_resources_permissions_permissions_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_resources_permissions_permissions_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_resources_permissions_permissions_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_permissions_permissions_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_permissions_permissions_proto_goTypes,
		DependencyIndexes: file_resources_permissions_permissions_proto_depIdxs,
		MessageInfos:      file_resources_permissions_permissions_proto_msgTypes,
	}.Build()
	File_resources_permissions_permissions_proto = out.File
	file_resources_permissions_permissions_proto_rawDesc = nil
	file_resources_permissions_permissions_proto_goTypes = nil
	file_resources_permissions_permissions_proto_depIdxs = nil
}
