// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: resources/dispatch/dispatch.proto

package dispatch

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	livemap "github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
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

type DISPATCH_STATUS int32

const (
	DISPATCH_STATUS_DECLINED        DISPATCH_STATUS = 0
	DISPATCH_STATUS_UNASSIGNED      DISPATCH_STATUS = 1
	DISPATCH_STATUS_EN_ROUTE        DISPATCH_STATUS = 2
	DISPATCH_STATUS_AT_SCENE        DISPATCH_STATUS = 3
	DISPATCH_STATUS_NEED_ASSISTANCE DISPATCH_STATUS = 4
	DISPATCH_STATUS_COMPLETED       DISPATCH_STATUS = 5
)

// Enum value maps for DISPATCH_STATUS.
var (
	DISPATCH_STATUS_name = map[int32]string{
		0: "DECLINED",
		1: "UNASSIGNED",
		2: "EN_ROUTE",
		3: "AT_SCENE",
		4: "NEED_ASSISTANCE",
		5: "COMPLETED",
	}
	DISPATCH_STATUS_value = map[string]int32{
		"DECLINED":        0,
		"UNASSIGNED":      1,
		"EN_ROUTE":        2,
		"AT_SCENE":        3,
		"NEED_ASSISTANCE": 4,
		"COMPLETED":       5,
	}
)

func (x DISPATCH_STATUS) Enum() *DISPATCH_STATUS {
	p := new(DISPATCH_STATUS)
	*p = x
	return p
}

func (x DISPATCH_STATUS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DISPATCH_STATUS) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_dispatch_dispatch_proto_enumTypes[0].Descriptor()
}

func (DISPATCH_STATUS) Type() protoreflect.EnumType {
	return &file_resources_dispatch_dispatch_proto_enumTypes[0]
}

func (x DISPATCH_STATUS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DISPATCH_STATUS.Descriptor instead.
func (DISPATCH_STATUS) EnumDescriptor() ([]byte, []int) {
	return file_resources_dispatch_dispatch_proto_rawDescGZIP(), []int{0}
}

type Dispatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	CreatedAt   *timestamp.Timestamp    `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	UpdatedAt   *timestamp.Timestamp    `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	Job         *string                 `protobuf:"bytes,4,opt,name=job,proto3,oneof" json:"job,omitempty"`
	Status      *DispatchStatus         `protobuf:"bytes,5,opt,name=status,proto3,oneof" json:"status,omitempty"`
	Message     string                  `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
	Description *string                 `protobuf:"bytes,7,opt,name=description,proto3,oneof" json:"description,omitempty"`
	Attributes  map[string]string       `protobuf:"bytes,8,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Marker      *livemap.DispatchMarker `protobuf:"bytes,9,opt,name=marker,proto3" json:"marker,omitempty"`
	Anon        *bool                   `protobuf:"varint,10,opt,name=anon,proto3,oneof" json:"anon,omitempty"`
	UserId      *bool                   `protobuf:"varint,11,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	User        *users.UserShort        `protobuf:"bytes,12,opt,name=user,proto3,oneof" json:"user,omitempty"`
	Units       []*DispatchAssignment   `protobuf:"bytes,13,rep,name=units,proto3" json:"units,omitempty"`
}

func (x *Dispatch) Reset() {
	*x = Dispatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_dispatch_dispatch_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Dispatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dispatch) ProtoMessage() {}

func (x *Dispatch) ProtoReflect() protoreflect.Message {
	mi := &file_resources_dispatch_dispatch_proto_msgTypes[0]
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
	return file_resources_dispatch_dispatch_proto_rawDescGZIP(), []int{0}
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
	if x != nil && x.Job != nil {
		return *x.Job
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

func (x *Dispatch) GetAttributes() map[string]string {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Dispatch) GetMarker() *livemap.DispatchMarker {
	if x != nil {
		return x.Marker
	}
	return nil
}

func (x *Dispatch) GetAnon() bool {
	if x != nil && x.Anon != nil {
		return *x.Anon
	}
	return false
}

func (x *Dispatch) GetUserId() bool {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return false
}

func (x *Dispatch) GetUser() *users.UserShort {
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

type DispatchAssignment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DispatchId uint64 `protobuf:"varint,1,opt,name=dispatch_id,json=dispatchId,proto3" json:"dispatch_id,omitempty"`
	UnitId     uint64 `protobuf:"varint,2,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
	Unit       *Unit  `protobuf:"bytes,3,opt,name=unit,proto3,oneof" json:"unit,omitempty"`
}

func (x *DispatchAssignment) Reset() {
	*x = DispatchAssignment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_dispatch_dispatch_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DispatchAssignment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DispatchAssignment) ProtoMessage() {}

func (x *DispatchAssignment) ProtoReflect() protoreflect.Message {
	mi := &file_resources_dispatch_dispatch_proto_msgTypes[1]
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
	return file_resources_dispatch_dispatch_proto_rawDescGZIP(), []int{1}
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

type DispatchStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	CreatedAt  *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	DispatchId uint64               `protobuf:"varint,3,opt,name=dispatch_id,json=dispatchId,proto3" json:"dispatch_id,omitempty"`
	UnitId     uint64               `protobuf:"varint,4,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
	Status     DISPATCH_STATUS      `protobuf:"varint,5,opt,name=status,proto3,enum=resources.dispatch.DISPATCH_STATUS" json:"status,omitempty"`
	Reason     *string              `protobuf:"bytes,6,opt,name=reason,proto3,oneof" json:"reason,omitempty"`
	Code       *string              `protobuf:"bytes,7,opt,name=code,proto3,oneof" json:"code,omitempty"`
	UserId     int32                `protobuf:"varint,8,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	User       *users.UserShort     `protobuf:"bytes,9,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *DispatchStatus) Reset() {
	*x = DispatchStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_dispatch_dispatch_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DispatchStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DispatchStatus) ProtoMessage() {}

func (x *DispatchStatus) ProtoReflect() protoreflect.Message {
	mi := &file_resources_dispatch_dispatch_proto_msgTypes[2]
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
	return file_resources_dispatch_dispatch_proto_rawDescGZIP(), []int{2}
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
	if x != nil {
		return x.UnitId
	}
	return 0
}

func (x *DispatchStatus) GetStatus() DISPATCH_STATUS {
	if x != nil {
		return x.Status
	}
	return DISPATCH_STATUS_DECLINED
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
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *DispatchStatus) GetUser() *users.UserShort {
	if x != nil {
		return x.User
	}
	return nil
}

var File_resources_dispatch_dispatch_proto protoreflect.FileDescriptor

var file_resources_dispatch_dispatch_proto_rawDesc = []byte{
	0x0a, 0x21, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x69, 0x73, 0x70,
	0x61, 0x74, 0x63, 0x68, 0x2f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x12, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64,
	0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x1a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x75, 0x6e, 0x69, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6d,
	0x61, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xa9, 0x06, 0x0a, 0x08, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x02,
	0x52, 0x03, 0x6a, 0x6f, 0x62, 0x88, 0x01, 0x01, 0x12, 0x3f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x44, 0x69,
	0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x48, 0x03, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72,
	0x03, 0x18, 0xff, 0x01, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2f, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x08, 0x48, 0x04, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x4c,
	0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x08, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64,
	0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x39, 0x0a, 0x06,
	0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70,
	0x2e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52,
	0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x17, 0x0a, 0x04, 0x61, 0x6e, 0x6f, 0x6e, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x08, 0x48, 0x05, 0x52, 0x04, 0x61, 0x6e, 0x6f, 0x6e, 0x88, 0x01, 0x01,
	0x12, 0x1c, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x08, 0x48, 0x06, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x33,
	0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x48, 0x07, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72,
	0x88, 0x01, 0x01, 0x12, 0x3c, 0x0a, 0x05, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x18, 0x0d, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x26, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64,
	0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68,
	0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x75, 0x6e, 0x69, 0x74,
	0x73, 0x1a, 0x3d, 0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42,
	0x0d, 0x0a, 0x0b, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x06,
	0x0a, 0x04, 0x5f, 0x6a, 0x6f, 0x62, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x61, 0x6e, 0x6f, 0x6e, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x22,
	0x8a, 0x01, 0x0a, 0x12, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x41, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74,
	0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x64, 0x69, 0x73,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x74, 0x49, 0x64,
	0x12, 0x31, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x70, 0x61,
	0x74, 0x63, 0x68, 0x2e, 0x55, 0x6e, 0x69, 0x74, 0x48, 0x00, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74,
	0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x6e, 0x69, 0x74, 0x22, 0xa3, 0x03, 0x0a,
	0x0e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74,
	0x63, 0x68, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x74, 0x49, 0x64, 0x12, 0x45, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74,
	0x63, 0x68, 0x2e, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54, 0x43, 0x48, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x48, 0x01,
	0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02,
	0x18, 0x14, 0x48, 0x02, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x2e, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x42,
	0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x09,
	0x0a, 0x07, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x63, 0x6f,
	0x64, 0x65, 0x2a, 0x6f, 0x0a, 0x0f, 0x44, 0x49, 0x53, 0x50, 0x41, 0x54, 0x43, 0x48, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x55, 0x53, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x45, 0x43, 0x4c, 0x49, 0x4e, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x55, 0x4e, 0x41, 0x53, 0x53, 0x49, 0x47, 0x4e, 0x45,
	0x44, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x4e, 0x5f, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x10,
	0x02, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x54, 0x5f, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x10, 0x03, 0x12,
	0x13, 0x0a, 0x0f, 0x4e, 0x45, 0x45, 0x44, 0x5f, 0x41, 0x53, 0x53, 0x49, 0x53, 0x54, 0x41, 0x4e,
	0x43, 0x45, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45,
	0x44, 0x10, 0x05, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65,
	0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x3b, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_resources_dispatch_dispatch_proto_rawDescOnce sync.Once
	file_resources_dispatch_dispatch_proto_rawDescData = file_resources_dispatch_dispatch_proto_rawDesc
)

func file_resources_dispatch_dispatch_proto_rawDescGZIP() []byte {
	file_resources_dispatch_dispatch_proto_rawDescOnce.Do(func() {
		file_resources_dispatch_dispatch_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_dispatch_dispatch_proto_rawDescData)
	})
	return file_resources_dispatch_dispatch_proto_rawDescData
}

var file_resources_dispatch_dispatch_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_dispatch_dispatch_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_resources_dispatch_dispatch_proto_goTypes = []interface{}{
	(DISPATCH_STATUS)(0),           // 0: resources.dispatch.DISPATCH_STATUS
	(*Dispatch)(nil),               // 1: resources.dispatch.Dispatch
	(*DispatchAssignment)(nil),     // 2: resources.dispatch.DispatchAssignment
	(*DispatchStatus)(nil),         // 3: resources.dispatch.DispatchStatus
	nil,                            // 4: resources.dispatch.Dispatch.AttributesEntry
	(*timestamp.Timestamp)(nil),    // 5: resources.timestamp.Timestamp
	(*livemap.DispatchMarker)(nil), // 6: resources.livemap.DispatchMarker
	(*users.UserShort)(nil),        // 7: resources.users.UserShort
	(*Unit)(nil),                   // 8: resources.dispatch.Unit
}
var file_resources_dispatch_dispatch_proto_depIdxs = []int32{
	5,  // 0: resources.dispatch.Dispatch.created_at:type_name -> resources.timestamp.Timestamp
	5,  // 1: resources.dispatch.Dispatch.updated_at:type_name -> resources.timestamp.Timestamp
	3,  // 2: resources.dispatch.Dispatch.status:type_name -> resources.dispatch.DispatchStatus
	4,  // 3: resources.dispatch.Dispatch.attributes:type_name -> resources.dispatch.Dispatch.AttributesEntry
	6,  // 4: resources.dispatch.Dispatch.marker:type_name -> resources.livemap.DispatchMarker
	7,  // 5: resources.dispatch.Dispatch.user:type_name -> resources.users.UserShort
	2,  // 6: resources.dispatch.Dispatch.units:type_name -> resources.dispatch.DispatchAssignment
	8,  // 7: resources.dispatch.DispatchAssignment.unit:type_name -> resources.dispatch.Unit
	5,  // 8: resources.dispatch.DispatchStatus.created_at:type_name -> resources.timestamp.Timestamp
	0,  // 9: resources.dispatch.DispatchStatus.status:type_name -> resources.dispatch.DISPATCH_STATUS
	7,  // 10: resources.dispatch.DispatchStatus.user:type_name -> resources.users.UserShort
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_resources_dispatch_dispatch_proto_init() }
func file_resources_dispatch_dispatch_proto_init() {
	if File_resources_dispatch_dispatch_proto != nil {
		return
	}
	file_resources_dispatch_units_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_resources_dispatch_dispatch_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_resources_dispatch_dispatch_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_resources_dispatch_dispatch_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
	file_resources_dispatch_dispatch_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_resources_dispatch_dispatch_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_resources_dispatch_dispatch_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_dispatch_dispatch_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_dispatch_dispatch_proto_goTypes,
		DependencyIndexes: file_resources_dispatch_dispatch_proto_depIdxs,
		EnumInfos:         file_resources_dispatch_dispatch_proto_enumTypes,
		MessageInfos:      file_resources_dispatch_dispatch_proto_msgTypes,
	}.Build()
	File_resources_dispatch_dispatch_proto = out.File
	file_resources_dispatch_dispatch_proto_rawDesc = nil
	file_resources_dispatch_dispatch_proto_goTypes = nil
	file_resources_dispatch_dispatch_proto_depIdxs = nil
}
