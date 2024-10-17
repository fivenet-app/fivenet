// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.20.3
// source: resources/common/cron/cron.proto

package cron

import (
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CronjobState int32

const (
	CronjobState_CRONJOB_STATE_UNSPECIFIED CronjobState = 0
	CronjobState_CRONJOB_STATE_WAITING     CronjobState = 1
	CronjobState_CRONJOB_STATE_PENDING     CronjobState = 2
	CronjobState_CRONJOB_STATE_RUNNING     CronjobState = 3
)

// Enum value maps for CronjobState.
var (
	CronjobState_name = map[int32]string{
		0: "CRONJOB_STATE_UNSPECIFIED",
		1: "CRONJOB_STATE_WAITING",
		2: "CRONJOB_STATE_PENDING",
		3: "CRONJOB_STATE_RUNNING",
	}
	CronjobState_value = map[string]int32{
		"CRONJOB_STATE_UNSPECIFIED": 0,
		"CRONJOB_STATE_WAITING":     1,
		"CRONJOB_STATE_PENDING":     2,
		"CRONJOB_STATE_RUNNING":     3,
	}
)

func (x CronjobState) Enum() *CronjobState {
	p := new(CronjobState)
	*p = x
	return p
}

func (x CronjobState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CronjobState) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_common_cron_cron_proto_enumTypes[0].Descriptor()
}

func (CronjobState) Type() protoreflect.EnumType {
	return &file_resources_common_cron_cron_proto_enumTypes[0]
}

func (x CronjobState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CronjobState.Descriptor instead.
func (CronjobState) EnumDescriptor() ([]byte, []int) {
	return file_resources_common_cron_cron_proto_rawDescGZIP(), []int{0}
}

type Cronjob struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name             string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Schedule         string               `protobuf:"bytes,2,opt,name=schedule,proto3" json:"schedule,omitempty"`
	State            CronjobState         `protobuf:"varint,3,opt,name=state,proto3,enum=resources.common.cron.CronjobState" json:"state,omitempty"`
	NextScheduleTime *timestamp.Timestamp `protobuf:"bytes,4,opt,name=next_schedule_time,json=nextScheduleTime,proto3" json:"next_schedule_time,omitempty"`
	LastAttemptTime  *timestamp.Timestamp `protobuf:"bytes,5,opt,name=last_attempt_time,json=lastAttemptTime,proto3" json:"last_attempt_time,omitempty"`
	Data             *CronjobData         `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Cronjob) Reset() {
	*x = Cronjob{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_common_cron_cron_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cronjob) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cronjob) ProtoMessage() {}

func (x *Cronjob) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_cron_cron_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cronjob.ProtoReflect.Descriptor instead.
func (*Cronjob) Descriptor() ([]byte, []int) {
	return file_resources_common_cron_cron_proto_rawDescGZIP(), []int{0}
}

func (x *Cronjob) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Cronjob) GetSchedule() string {
	if x != nil {
		return x.Schedule
	}
	return ""
}

func (x *Cronjob) GetState() CronjobState {
	if x != nil {
		return x.State
	}
	return CronjobState_CRONJOB_STATE_UNSPECIFIED
}

func (x *Cronjob) GetNextScheduleTime() *timestamp.Timestamp {
	if x != nil {
		return x.NextScheduleTime
	}
	return nil
}

func (x *Cronjob) GetLastAttemptTime() *timestamp.Timestamp {
	if x != nil {
		return x.LastAttemptTime
	}
	return nil
}

func (x *Cronjob) GetData() *CronjobData {
	if x != nil {
		return x.Data
	}
	return nil
}

type CronjobData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,1,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Data      *anypb.Any           `protobuf:"bytes,2,opt,name=data,proto3,oneof" json:"data,omitempty"`
}

func (x *CronjobData) Reset() {
	*x = CronjobData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_common_cron_cron_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CronjobData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CronjobData) ProtoMessage() {}

func (x *CronjobData) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_cron_cron_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CronjobData.ProtoReflect.Descriptor instead.
func (*CronjobData) Descriptor() ([]byte, []int) {
	return file_resources_common_cron_cron_proto_rawDescGZIP(), []int{1}
}

func (x *CronjobData) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *CronjobData) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

type CronjobLockOwnerState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hostname  string               `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *CronjobLockOwnerState) Reset() {
	*x = CronjobLockOwnerState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_common_cron_cron_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CronjobLockOwnerState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CronjobLockOwnerState) ProtoMessage() {}

func (x *CronjobLockOwnerState) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_cron_cron_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CronjobLockOwnerState.ProtoReflect.Descriptor instead.
func (*CronjobLockOwnerState) Descriptor() ([]byte, []int) {
	return file_resources_common_cron_cron_proto_rawDescGZIP(), []int{2}
}

func (x *CronjobLockOwnerState) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *CronjobLockOwnerState) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type CronjobSchedulerEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cronjob *Cronjob `protobuf:"bytes,1,opt,name=cronjob,proto3" json:"cronjob,omitempty"`
}

func (x *CronjobSchedulerEvent) Reset() {
	*x = CronjobSchedulerEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_common_cron_cron_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CronjobSchedulerEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CronjobSchedulerEvent) ProtoMessage() {}

func (x *CronjobSchedulerEvent) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_cron_cron_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CronjobSchedulerEvent.ProtoReflect.Descriptor instead.
func (*CronjobSchedulerEvent) Descriptor() ([]byte, []int) {
	return file_resources_common_cron_cron_proto_rawDescGZIP(), []int{3}
}

func (x *CronjobSchedulerEvent) GetCronjob() *Cronjob {
	if x != nil {
		return x.Cronjob
	}
	return nil
}

type CronjobCompletedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Sucess  bool                 `protobuf:"varint,2,opt,name=sucess,proto3" json:"sucess,omitempty"`
	EndDate *timestamp.Timestamp `protobuf:"bytes,3,opt,name=endDate,proto3" json:"endDate,omitempty"`
	Elapsed *durationpb.Duration `protobuf:"bytes,4,opt,name=elapsed,proto3" json:"elapsed,omitempty"`
	Data    *CronjobData         `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *CronjobCompletedEvent) Reset() {
	*x = CronjobCompletedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_common_cron_cron_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CronjobCompletedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CronjobCompletedEvent) ProtoMessage() {}

func (x *CronjobCompletedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_cron_cron_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CronjobCompletedEvent.ProtoReflect.Descriptor instead.
func (*CronjobCompletedEvent) Descriptor() ([]byte, []int) {
	return file_resources_common_cron_cron_proto_rawDescGZIP(), []int{4}
}

func (x *CronjobCompletedEvent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CronjobCompletedEvent) GetSucess() bool {
	if x != nil {
		return x.Sucess
	}
	return false
}

func (x *CronjobCompletedEvent) GetEndDate() *timestamp.Timestamp {
	if x != nil {
		return x.EndDate
	}
	return nil
}

func (x *CronjobCompletedEvent) GetElapsed() *durationpb.Duration {
	if x != nil {
		return x.Elapsed
	}
	return nil
}

func (x *CronjobCompletedEvent) GetData() *CronjobData {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_resources_common_cron_cron_proto protoreflect.FileDescriptor

var file_resources_common_cron_cron_proto_rawDesc = []byte{
	0x0a, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2f, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x15, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc6, 0x02, 0x0a, 0x07, 0x43, 0x72,
	0x6f, 0x6e, 0x6a, 0x6f, 0x62, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x39, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x43, 0x72, 0x6f,
	0x6e, 0x6a, 0x6f, 0x62, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x4c, 0x0a, 0x12, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x10, 0x6e, 0x65,
	0x78, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x4a,
	0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x41,
	0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x36, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x72, 0x6f, 0x6e,
	0x2e, 0x43, 0x72, 0x6f, 0x6e, 0x6a, 0x6f, 0x62, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x84, 0x01, 0x0a, 0x0b, 0x43, 0x72, 0x6f, 0x6e, 0x6a, 0x6f, 0x62, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x3d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x2d, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x41, 0x6e, 0x79, 0x48, 0x00, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x88, 0x01, 0x01,
	0x42, 0x07, 0x0a, 0x05, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x22, 0x72, 0x0a, 0x15, 0x43, 0x72, 0x6f,
	0x6e, 0x6a, 0x6f, 0x62, 0x4c, 0x6f, 0x63, 0x6b, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3d,
	0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x51, 0x0a,
	0x15, 0x43, 0x72, 0x6f, 0x6e, 0x6a, 0x6f, 0x62, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x07, 0x63, 0x72, 0x6f, 0x6e, 0x6a, 0x6f,
	0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e,
	0x43, 0x72, 0x6f, 0x6e, 0x6a, 0x6f, 0x62, 0x52, 0x07, 0x63, 0x72, 0x6f, 0x6e, 0x6a, 0x6f, 0x62,
	0x22, 0xea, 0x01, 0x0a, 0x15, 0x43, 0x72, 0x6f, 0x6e, 0x6a, 0x6f, 0x62, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x75, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x73, 0x75, 0x63, 0x65, 0x73, 0x73, 0x12, 0x38, 0x0a, 0x07, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x65,
	0x12, 0x33, 0x0a, 0x07, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x65, 0x6c,
	0x61, 0x70, 0x73, 0x65, 0x64, 0x12, 0x36, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x72, 0x6f, 0x6e, 0x2e, 0x43, 0x72, 0x6f, 0x6e,
	0x6a, 0x6f, 0x62, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x7e, 0x0a,
	0x0c, 0x43, 0x72, 0x6f, 0x6e, 0x6a, 0x6f, 0x62, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x0a,
	0x19, 0x43, 0x52, 0x4f, 0x4e, 0x4a, 0x4f, 0x42, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x19, 0x0a, 0x15,
	0x43, 0x52, 0x4f, 0x4e, 0x4a, 0x4f, 0x42, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x57, 0x41,
	0x49, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x43, 0x52, 0x4f, 0x4e, 0x4a,
	0x4f, 0x42, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47,
	0x10, 0x02, 0x12, 0x19, 0x0a, 0x15, 0x43, 0x52, 0x4f, 0x4e, 0x4a, 0x4f, 0x42, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x45, 0x5f, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x42, 0x48, 0x5a,
	0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65,
	0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x72,
	0x6f, 0x6e, 0x3b, 0x63, 0x72, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_common_cron_cron_proto_rawDescOnce sync.Once
	file_resources_common_cron_cron_proto_rawDescData = file_resources_common_cron_cron_proto_rawDesc
)

func file_resources_common_cron_cron_proto_rawDescGZIP() []byte {
	file_resources_common_cron_cron_proto_rawDescOnce.Do(func() {
		file_resources_common_cron_cron_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_common_cron_cron_proto_rawDescData)
	})
	return file_resources_common_cron_cron_proto_rawDescData
}

var file_resources_common_cron_cron_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_common_cron_cron_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_resources_common_cron_cron_proto_goTypes = []any{
	(CronjobState)(0),             // 0: resources.common.cron.CronjobState
	(*Cronjob)(nil),               // 1: resources.common.cron.Cronjob
	(*CronjobData)(nil),           // 2: resources.common.cron.CronjobData
	(*CronjobLockOwnerState)(nil), // 3: resources.common.cron.CronjobLockOwnerState
	(*CronjobSchedulerEvent)(nil), // 4: resources.common.cron.CronjobSchedulerEvent
	(*CronjobCompletedEvent)(nil), // 5: resources.common.cron.CronjobCompletedEvent
	(*timestamp.Timestamp)(nil),   // 6: resources.timestamp.Timestamp
	(*anypb.Any)(nil),             // 7: google.protobuf.Any
	(*durationpb.Duration)(nil),   // 8: google.protobuf.Duration
}
var file_resources_common_cron_cron_proto_depIdxs = []int32{
	0,  // 0: resources.common.cron.Cronjob.state:type_name -> resources.common.cron.CronjobState
	6,  // 1: resources.common.cron.Cronjob.next_schedule_time:type_name -> resources.timestamp.Timestamp
	6,  // 2: resources.common.cron.Cronjob.last_attempt_time:type_name -> resources.timestamp.Timestamp
	2,  // 3: resources.common.cron.Cronjob.data:type_name -> resources.common.cron.CronjobData
	6,  // 4: resources.common.cron.CronjobData.updated_at:type_name -> resources.timestamp.Timestamp
	7,  // 5: resources.common.cron.CronjobData.data:type_name -> google.protobuf.Any
	6,  // 6: resources.common.cron.CronjobLockOwnerState.updated_at:type_name -> resources.timestamp.Timestamp
	1,  // 7: resources.common.cron.CronjobSchedulerEvent.cronjob:type_name -> resources.common.cron.Cronjob
	6,  // 8: resources.common.cron.CronjobCompletedEvent.endDate:type_name -> resources.timestamp.Timestamp
	8,  // 9: resources.common.cron.CronjobCompletedEvent.elapsed:type_name -> google.protobuf.Duration
	2,  // 10: resources.common.cron.CronjobCompletedEvent.data:type_name -> resources.common.cron.CronjobData
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_resources_common_cron_cron_proto_init() }
func file_resources_common_cron_cron_proto_init() {
	if File_resources_common_cron_cron_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_common_cron_cron_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Cronjob); i {
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
		file_resources_common_cron_cron_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CronjobData); i {
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
		file_resources_common_cron_cron_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*CronjobLockOwnerState); i {
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
		file_resources_common_cron_cron_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*CronjobSchedulerEvent); i {
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
		file_resources_common_cron_cron_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*CronjobCompletedEvent); i {
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
	file_resources_common_cron_cron_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_common_cron_cron_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_common_cron_cron_proto_goTypes,
		DependencyIndexes: file_resources_common_cron_cron_proto_depIdxs,
		EnumInfos:         file_resources_common_cron_cron_proto_enumTypes,
		MessageInfos:      file_resources_common_cron_cron_proto_msgTypes,
	}.Build()
	File_resources_common_cron_cron_proto = out.File
	file_resources_common_cron_cron_proto_rawDesc = nil
	file_resources_common_cron_cron_proto_goTypes = nil
	file_resources_common_cron_cron_proto_depIdxs = nil
}
