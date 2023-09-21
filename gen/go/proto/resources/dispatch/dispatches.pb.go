// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.2
// source: resources/dispatch/dispatches.proto

package dispatch

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	timestamp "github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
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

type StatusDispatch int32

const (
	StatusDispatch_STATUS_DISPATCH_UNSPECIFIED     StatusDispatch = 0
	StatusDispatch_STATUS_DISPATCH_NEW             StatusDispatch = 1
	StatusDispatch_STATUS_DISPATCH_UNASSIGNED      StatusDispatch = 2
	StatusDispatch_STATUS_DISPATCH_UPDATED         StatusDispatch = 3
	StatusDispatch_STATUS_DISPATCH_UNIT_ASSIGNED   StatusDispatch = 4
	StatusDispatch_STATUS_DISPATCH_UNIT_UNASSIGNED StatusDispatch = 5
	StatusDispatch_STATUS_DISPATCH_EN_ROUTE        StatusDispatch = 6
	StatusDispatch_STATUS_DISPATCH_ON_SCENE        StatusDispatch = 7
	StatusDispatch_STATUS_DISPATCH_NEED_ASSISTANCE StatusDispatch = 8
	StatusDispatch_STATUS_DISPATCH_COMPLETED       StatusDispatch = 9
	StatusDispatch_STATUS_DISPATCH_CANCELLED       StatusDispatch = 10
	StatusDispatch_STATUS_DISPATCH_ARCHIVED        StatusDispatch = 11
)

// Enum value maps for StatusDispatch.
var (
	StatusDispatch_name = map[int32]string{
		0:  "STATUS_DISPATCH_UNSPECIFIED",
		1:  "STATUS_DISPATCH_NEW",
		2:  "STATUS_DISPATCH_UNASSIGNED",
		3:  "STATUS_DISPATCH_UPDATED",
		4:  "STATUS_DISPATCH_UNIT_ASSIGNED",
		5:  "STATUS_DISPATCH_UNIT_UNASSIGNED",
		6:  "STATUS_DISPATCH_EN_ROUTE",
		7:  "STATUS_DISPATCH_ON_SCENE",
		8:  "STATUS_DISPATCH_NEED_ASSISTANCE",
		9:  "STATUS_DISPATCH_COMPLETED",
		10: "STATUS_DISPATCH_CANCELLED",
		11: "STATUS_DISPATCH_ARCHIVED",
	}
	StatusDispatch_value = map[string]int32{
		"STATUS_DISPATCH_UNSPECIFIED":     0,
		"STATUS_DISPATCH_NEW":             1,
		"STATUS_DISPATCH_UNASSIGNED":      2,
		"STATUS_DISPATCH_UPDATED":         3,
		"STATUS_DISPATCH_UNIT_ASSIGNED":   4,
		"STATUS_DISPATCH_UNIT_UNASSIGNED": 5,
		"STATUS_DISPATCH_EN_ROUTE":        6,
		"STATUS_DISPATCH_ON_SCENE":        7,
		"STATUS_DISPATCH_NEED_ASSISTANCE": 8,
		"STATUS_DISPATCH_COMPLETED":       9,
		"STATUS_DISPATCH_CANCELLED":       10,
		"STATUS_DISPATCH_ARCHIVED":        11,
	}
)

func (x StatusDispatch) Enum() *StatusDispatch {
	p := new(StatusDispatch)
	*p = x
	return p
}

func (x StatusDispatch) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StatusDispatch) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_dispatch_dispatches_proto_enumTypes[0].Descriptor()
}

func (StatusDispatch) Type() protoreflect.EnumType {
	return &file_resources_dispatch_dispatches_proto_enumTypes[0]
}

func (x StatusDispatch) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StatusDispatch.Descriptor instead.
func (StatusDispatch) EnumDescriptor() ([]byte, []int) {
	return file_resources_dispatch_dispatches_proto_rawDescGZIP(), []int{0}
}

type Dispatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	Job       string               `protobuf:"bytes,4,opt,name=job,proto3" json:"job,omitempty"`
	Status    *DispatchStatus      `protobuf:"bytes,5,opt,name=status,proto3,oneof" json:"status,omitempty"`
	// @sanitize
	Message string `protobuf:"bytes,7,opt,name=message,proto3" json:"message,omitempty"`
	// @sanitize
	Description *string     `protobuf:"bytes,8,opt,name=description,proto3,oneof" json:"description,omitempty"`
	Attributes  *Attributes `protobuf:"bytes,9,opt,name=attributes,proto3,oneof" json:"attributes,omitempty"`
	X           float64     `protobuf:"fixed64,10,opt,name=x,proto3" json:"x,omitempty"`
	Y           float64     `protobuf:"fixed64,11,opt,name=y,proto3" json:"y,omitempty"`
	// @sanitize
	Postal *string               `protobuf:"bytes,12,opt,name=postal,proto3,oneof" json:"postal,omitempty"`
	Anon   bool                  `protobuf:"varint,13,opt,name=anon,proto3" json:"anon,omitempty"`
	UserId *int32                `protobuf:"varint,14,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	User   *users.User           `protobuf:"bytes,15,opt,name=user,proto3,oneof" json:"user,omitempty"`
	Units  []*DispatchAssignment `protobuf:"bytes,16,rep,name=units,proto3" json:"units,omitempty"`
}

func (x *Dispatch) Reset() {
	*x = Dispatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_dispatch_dispatches_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Dispatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dispatch) ProtoMessage() {}

func (x *Dispatch) ProtoReflect() protoreflect.Message {
	mi := &file_resources_dispatch_dispatches_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dispatch.ProtoReflect.Descriptor instead.
func (*Dispatch) Descriptor() ([]byte, []int) {
	return file_resources_dispatch_dispatches_proto_rawDescGZIP(), []int{0}
}

func (x *Dispatch) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Dispatch) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Dispatch) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Dispatch) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *Dispatch) GetStatus() *DispatchStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *Dispatch) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Dispatch) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Dispatch) GetAttributes() *Attributes {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Dispatch) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Dispatch) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *Dispatch) GetPostal() string {
	if x != nil && x.Postal != nil {
		return *x.Postal
	}
	return ""
}

func (x *Dispatch) GetAnon() bool {
	if x != nil {
		return x.Anon
	}
	return false
}

func (x *Dispatch) GetUserId() int32 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

func (x *Dispatch) GetUser() *users.User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *Dispatch) GetUnits() []*DispatchAssignment {
	if x != nil {
		return x.Units
	}
	return nil
}

type Attributes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []string `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *Attributes) Reset() {
	*x = Attributes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_dispatch_dispatches_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attributes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attributes) ProtoMessage() {}

func (x *Attributes) ProtoReflect() protoreflect.Message {
	mi := &file_resources_dispatch_dispatches_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attributes.ProtoReflect.Descriptor instead.
func (*Attributes) Descriptor() ([]byte, []int) {
	return file_resources_dispatch_dispatches_proto_rawDescGZIP(), []int{1}
}

func (x *Attributes) GetList() []string {
	if x != nil {
		return x.List
	}
	return nil
}

type DispatchAssignments struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DispatchId uint64                `protobuf:"varint,1,opt,name=dispatch_id,json=dispatchId,proto3" json:"dispatch_id,omitempty"`
	Job        string                `protobuf:"bytes,2,opt,name=job,proto3" json:"job,omitempty"`
	Units      []*DispatchAssignment `protobuf:"bytes,3,rep,name=units,proto3" json:"units,omitempty"`
}

func (x *DispatchAssignments) Reset() {
	*x = DispatchAssignments{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_dispatch_dispatches_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DispatchAssignments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DispatchAssignments) ProtoMessage() {}

func (x *DispatchAssignments) ProtoReflect() protoreflect.Message {
	mi := &file_resources_dispatch_dispatches_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DispatchAssignments.ProtoReflect.Descriptor instead.
func (*DispatchAssignments) Descriptor() ([]byte, []int) {
	return file_resources_dispatch_dispatches_proto_rawDescGZIP(), []int{2}
}

func (x *DispatchAssignments) GetDispatchId() uint64 {
	if x != nil {
		return x.DispatchId
	}
	return 0
}

func (x *DispatchAssignments) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *DispatchAssignments) GetUnits() []*DispatchAssignment {
	if x != nil {
		return x.Units
	}
	return nil
}

type DispatchAssignment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DispatchId uint64               `protobuf:"varint,1,opt,name=dispatch_id,json=dispatchId,proto3" json:"dispatch_id,omitempty" sql:"primary_key" alias:"dispatch_id"` // @gotags: sql:"primary_key" alias:"dispatch_id"
	UnitId     uint64               `protobuf:"varint,2,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty" sql:"primary_key" alias:"unit_id"`             // @gotags: sql:"primary_key" alias:"unit_id"
	Unit       *Unit                `protobuf:"bytes,3,opt,name=unit,proto3,oneof" json:"unit,omitempty"`
	CreatedAt  *timestamp.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	ExpiresAt  *timestamp.Timestamp `protobuf:"bytes,5,opt,name=expires_at,json=expiresAt,proto3,oneof" json:"expires_at,omitempty"`
}

func (x *DispatchAssignment) Reset() {
	*x = DispatchAssignment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_dispatch_dispatches_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DispatchAssignment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DispatchAssignment) ProtoMessage() {}

func (x *DispatchAssignment) ProtoReflect() protoreflect.Message {
	mi := &file_resources_dispatch_dispatches_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DispatchAssignment.ProtoReflect.Descriptor instead.
func (*DispatchAssignment) Descriptor() ([]byte, []int) {
	return file_resources_dispatch_dispatches_proto_rawDescGZIP(), []int{3}
}

func (x *DispatchAssignment) GetDispatchId() uint64 {
	if x != nil {
		return x.DispatchId
	}
	return 0
}

func (x *DispatchAssignment) GetUnitId() uint64 {
	if x != nil {
		return x.UnitId
	}
	return 0
}

func (x *DispatchAssignment) GetUnit() *Unit {
	if x != nil {
		return x.Unit
	}
	return nil
}

func (x *DispatchAssignment) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *DispatchAssignment) GetExpiresAt() *timestamp.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

type DispatchStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	CreatedAt  *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	DispatchId uint64               `protobuf:"varint,3,opt,name=dispatch_id,json=dispatchId,proto3" json:"dispatch_id,omitempty"`
	UnitId     *uint64              `protobuf:"varint,4,opt,name=unit_id,json=unitId,proto3,oneof" json:"unit_id,omitempty"`
	Unit       *Unit                `protobuf:"bytes,5,opt,name=unit,proto3,oneof" json:"unit,omitempty"`
	Status     StatusDispatch       `protobuf:"varint,6,opt,name=status,proto3,enum=resources.dispatch.StatusDispatch" json:"status,omitempty"`
	// @sanitize
	Reason *string `protobuf:"bytes,7,opt,name=reason,proto3,oneof" json:"reason,omitempty"`
	// @sanitize
	Code   *string          `protobuf:"bytes,8,opt,name=code,proto3,oneof" json:"code,omitempty"`
	UserId *int32           `protobuf:"varint,9,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	User   *users.UserShort `protobuf:"bytes,10,opt,name=user,proto3,oneof" json:"user,omitempty"`
	X      *float64         `protobuf:"fixed64,11,opt,name=x,proto3,oneof" json:"x,omitempty"`
	Y      *float64         `protobuf:"fixed64,12,opt,name=y,proto3,oneof" json:"y,omitempty"`
	// @sanitize
	Postal *string `protobuf:"bytes,13,opt,name=postal,proto3,oneof" json:"postal,omitempty"`
}

func (x *DispatchStatus) Reset() {
	*x = DispatchStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_dispatch_dispatches_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DispatchStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DispatchStatus) ProtoMessage() {}

func (x *DispatchStatus) ProtoReflect() protoreflect.Message {
	mi := &file_resources_dispatch_dispatches_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DispatchStatus.ProtoReflect.Descriptor instead.
func (*DispatchStatus) Descriptor() ([]byte, []int) {
	return file_resources_dispatch_dispatches_proto_rawDescGZIP(), []int{4}
}

func (x *DispatchStatus) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DispatchStatus) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *DispatchStatus) GetDispatchId() uint64 {
	if x != nil {
		return x.DispatchId
	}
	return 0
}

func (x *DispatchStatus) GetUnitId() uint64 {
	if x != nil && x.UnitId != nil {
		return *x.UnitId
	}
	return 0
}

func (x *DispatchStatus) GetUnit() *Unit {
	if x != nil {
		return x.Unit
	}
	return nil
}

func (x *DispatchStatus) GetStatus() StatusDispatch {
	if x != nil {
		return x.Status
	}
	return StatusDispatch_STATUS_DISPATCH_UNSPECIFIED
}

func (x *DispatchStatus) GetReason() string {
	if x != nil && x.Reason != nil {
		return *x.Reason
	}
	return ""
}

func (x *DispatchStatus) GetCode() string {
	if x != nil && x.Code != nil {
		return *x.Code
	}
	return ""
}

func (x *DispatchStatus) GetUserId() int32 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

func (x *DispatchStatus) GetUser() *users.UserShort {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *DispatchStatus) GetX() float64 {
	if x != nil && x.X != nil {
		return *x.X
	}
	return 0
}

func (x *DispatchStatus) GetY() float64 {
	if x != nil && x.Y != nil {
		return *x.Y
	}
	return 0
}

func (x *DispatchStatus) GetPostal() string {
	if x != nil && x.Postal != nil {
		return *x.Postal
	}
	return ""
}

var File_resources_dispatch_dispatches_proto protoreflect.FileDescriptor

var file_resources_dispatch_dispatches_proto_rawDesc = []byte{
	0x0a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x69, 0x73, 0x70,
	0x61, 0x74, 0x63, 0x68, 0x2f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x1a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x75, 0x6e,
	0x69, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xeb, 0x05, 0x0a, 0x08, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f, 0x62,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52,
	0x03, 0x6a, 0x6f, 0x62, 0x12, 0x3f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74,
	0x63, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x48, 0x02, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2f, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x08, 0x48, 0x03, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x43, 0x0a, 0x0a, 0x61, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x70, 0x61,
	0x74, 0x63, 0x68, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x48, 0x04,
	0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x88, 0x01, 0x01, 0x12,
	0x0c, 0x0a, 0x01, 0x78, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a,
	0x01, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x79, 0x12, 0x24, 0x0a, 0x06, 0x70,
	0x6f, 0x73, 0x74, 0x61, 0x6c, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x18, 0x30, 0x48, 0x05, 0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x61, 0x6c, 0x88, 0x01,
	0x01, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x6e, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x04, 0x61, 0x6e, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x48,
	0x06, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x2e, 0x0a, 0x04,
	0x75, 0x73, 0x65, 0x72, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x48, 0x07, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x3c, 0x0a, 0x05,
	0x75, 0x6e, 0x69, 0x74, 0x73, 0x18, 0x10, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x2e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x05, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x61, 0x6c, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x22, 0x20, 0x0a, 0x0a, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04,
	0x6c, 0x69, 0x73, 0x74, 0x22, 0x8f, 0x01, 0x0a, 0x13, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1f, 0x0a, 0x0b,
	0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x12, 0x19, 0x0a,
	0x03, 0x6a, 0x6f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72,
	0x02, 0x18, 0x14, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x3c, 0x0a, 0x05, 0x75, 0x6e, 0x69, 0x74,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x44, 0x69, 0x73,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x05, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x22, 0xb0, 0x02, 0x0a, 0x12, 0x44, 0x69, 0x73, 0x70, 0x61,
	0x74, 0x63, 0x68, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x75, 0x6e, 0x69, 0x74, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x55, 0x6e, 0x69, 0x74, 0x48,
	0x00, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42,
	0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x48, 0x02, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x88,
	0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x6e, 0x69, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x65,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61, 0x74, 0x22, 0xf1, 0x04, 0x0a, 0x0e, 0x44, 0x69,
	0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01,
	0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x49,
	0x64, 0x12, 0x1c, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x04, 0x48, 0x01, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12,
	0x31, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74,
	0x63, 0x68, 0x2e, 0x55, 0x6e, 0x69, 0x74, 0x48, 0x02, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x44, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64,
	0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x44, 0x69,
	0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18,
	0xff, 0x01, 0x48, 0x03, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12,
	0x20, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x48, 0x04, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x25, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x48, 0x05, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f,
	0x72, 0x74, 0x48, 0x06, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x11, 0x0a,
	0x01, 0x78, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x48, 0x07, 0x52, 0x01, 0x78, 0x88, 0x01, 0x01,
	0x12, 0x11, 0x0a, 0x01, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x01, 0x48, 0x08, 0x52, 0x01, 0x79,
	0x88, 0x01, 0x01, 0x12, 0x24, 0x0a, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x61, 0x6c, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x30, 0x48, 0x09, 0x52, 0x06,
	0x70, 0x6f, 0x73, 0x74, 0x61, 0x6c, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x75, 0x6e, 0x69,
	0x74, 0x5f, 0x69, 0x64, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x6e, 0x69, 0x74, 0x42, 0x09, 0x0a,
	0x07, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x63, 0x6f, 0x64,
	0x65, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x42, 0x07, 0x0a,
	0x05, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x04, 0x0a, 0x02, 0x5f, 0x78, 0x42, 0x04, 0x0a, 0x02,
	0x5f, 0x79, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x61, 0x6c, 0x2a, 0x8c, 0x03,
	0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x12, 0x1f, 0x0a, 0x1b, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41,
	0x54, 0x43, 0x48, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50,
	0x41, 0x54, 0x43, 0x48, 0x5f, 0x4e, 0x45, 0x57, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54, 0x43, 0x48, 0x5f, 0x55, 0x4e,
	0x41, 0x53, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x44, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54, 0x43, 0x48, 0x5f, 0x55, 0x50,
	0x44, 0x41, 0x54, 0x45, 0x44, 0x10, 0x03, 0x12, 0x21, 0x0a, 0x1d, 0x53, 0x54, 0x41, 0x54, 0x55,
	0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54, 0x43, 0x48, 0x5f, 0x55, 0x4e, 0x49, 0x54, 0x5f,
	0x41, 0x53, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x44, 0x10, 0x04, 0x12, 0x23, 0x0a, 0x1f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54, 0x43, 0x48, 0x5f, 0x55, 0x4e,
	0x49, 0x54, 0x5f, 0x55, 0x4e, 0x41, 0x53, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x44, 0x10, 0x05, 0x12,
	0x1c, 0x0a, 0x18, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54,
	0x43, 0x48, 0x5f, 0x45, 0x4e, 0x5f, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x10, 0x06, 0x12, 0x1c, 0x0a,
	0x18, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54, 0x43, 0x48,
	0x5f, 0x4f, 0x4e, 0x5f, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x10, 0x07, 0x12, 0x23, 0x0a, 0x1f, 0x53,
	0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54, 0x43, 0x48, 0x5f, 0x4e,
	0x45, 0x45, 0x44, 0x5f, 0x41, 0x53, 0x53, 0x49, 0x53, 0x54, 0x41, 0x4e, 0x43, 0x45, 0x10, 0x08,
	0x12, 0x1d, 0x0a, 0x19, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41,
	0x54, 0x43, 0x48, 0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x09, 0x12,
	0x1d, 0x0a, 0x19, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54,
	0x43, 0x48, 0x5f, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x4c, 0x45, 0x44, 0x10, 0x0a, 0x12, 0x1c,
	0x0a, 0x18, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54, 0x43,
	0x48, 0x5f, 0x41, 0x52, 0x43, 0x48, 0x49, 0x56, 0x45, 0x44, 0x10, 0x0b, 0x42, 0x45, 0x5a, 0x43,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78,
	0x72, 0x74, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x3b, 0x64, 0x69, 0x73, 0x70, 0x61,
	0x74, 0x63, 0x68, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_dispatch_dispatches_proto_rawDescOnce sync.Once
	file_resources_dispatch_dispatches_proto_rawDescData = file_resources_dispatch_dispatches_proto_rawDesc
)

func file_resources_dispatch_dispatches_proto_rawDescGZIP() []byte {
	file_resources_dispatch_dispatches_proto_rawDescOnce.Do(func() {
		file_resources_dispatch_dispatches_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_dispatch_dispatches_proto_rawDescData)
	})
	return file_resources_dispatch_dispatches_proto_rawDescData
}

var file_resources_dispatch_dispatches_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_dispatch_dispatches_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_resources_dispatch_dispatches_proto_goTypes = []interface{}{
	(StatusDispatch)(0),         // 0: resources.dispatch.StatusDispatch
	(*Dispatch)(nil),            // 1: resources.dispatch.Dispatch
	(*Attributes)(nil),          // 2: resources.dispatch.Attributes
	(*DispatchAssignments)(nil), // 3: resources.dispatch.DispatchAssignments
	(*DispatchAssignment)(nil),  // 4: resources.dispatch.DispatchAssignment
	(*DispatchStatus)(nil),      // 5: resources.dispatch.DispatchStatus
	(*timestamp.Timestamp)(nil), // 6: resources.timestamp.Timestamp
	(*users.User)(nil),          // 7: resources.users.User
	(*Unit)(nil),                // 8: resources.dispatch.Unit
	(*users.UserShort)(nil),     // 9: resources.users.UserShort
}
var file_resources_dispatch_dispatches_proto_depIdxs = []int32{
	6,  // 0: resources.dispatch.Dispatch.created_at:type_name -> resources.timestamp.Timestamp
	6,  // 1: resources.dispatch.Dispatch.updated_at:type_name -> resources.timestamp.Timestamp
	5,  // 2: resources.dispatch.Dispatch.status:type_name -> resources.dispatch.DispatchStatus
	2,  // 3: resources.dispatch.Dispatch.attributes:type_name -> resources.dispatch.Attributes
	7,  // 4: resources.dispatch.Dispatch.user:type_name -> resources.users.User
	4,  // 5: resources.dispatch.Dispatch.units:type_name -> resources.dispatch.DispatchAssignment
	4,  // 6: resources.dispatch.DispatchAssignments.units:type_name -> resources.dispatch.DispatchAssignment
	8,  // 7: resources.dispatch.DispatchAssignment.unit:type_name -> resources.dispatch.Unit
	6,  // 8: resources.dispatch.DispatchAssignment.created_at:type_name -> resources.timestamp.Timestamp
	6,  // 9: resources.dispatch.DispatchAssignment.expires_at:type_name -> resources.timestamp.Timestamp
	6,  // 10: resources.dispatch.DispatchStatus.created_at:type_name -> resources.timestamp.Timestamp
	8,  // 11: resources.dispatch.DispatchStatus.unit:type_name -> resources.dispatch.Unit
	0,  // 12: resources.dispatch.DispatchStatus.status:type_name -> resources.dispatch.StatusDispatch
	9,  // 13: resources.dispatch.DispatchStatus.user:type_name -> resources.users.UserShort
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_resources_dispatch_dispatches_proto_init() }
func file_resources_dispatch_dispatches_proto_init() {
	if File_resources_dispatch_dispatches_proto != nil {
		return
	}
	file_resources_dispatch_units_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_resources_dispatch_dispatches_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Dispatch); i {
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
		file_resources_dispatch_dispatches_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attributes); i {
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
		file_resources_dispatch_dispatches_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DispatchAssignments); i {
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
		file_resources_dispatch_dispatches_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DispatchAssignment); i {
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
		file_resources_dispatch_dispatches_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DispatchStatus); i {
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
	file_resources_dispatch_dispatches_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_resources_dispatch_dispatches_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_resources_dispatch_dispatches_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_dispatch_dispatches_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_dispatch_dispatches_proto_goTypes,
		DependencyIndexes: file_resources_dispatch_dispatches_proto_depIdxs,
		EnumInfos:         file_resources_dispatch_dispatches_proto_enumTypes,
		MessageInfos:      file_resources_dispatch_dispatches_proto_msgTypes,
	}.Build()
	File_resources_dispatch_dispatches_proto = out.File
	file_resources_dispatch_dispatches_proto_rawDesc = nil
	file_resources_dispatch_dispatches_proto_goTypes = nil
	file_resources_dispatch_dispatches_proto_depIdxs = nil
}
