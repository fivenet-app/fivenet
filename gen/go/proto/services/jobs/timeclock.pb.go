// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: services/jobs/timeclock.proto

package jobs

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
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

type ListTimeclockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *database.PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	// Search
	UserIds []int32              `protobuf:"varint,2,rep,packed,name=user_ids,json=userIds,proto3" json:"user_ids,omitempty"`
	From    *timestamp.Timestamp `protobuf:"bytes,3,opt,name=from,proto3,oneof" json:"from,omitempty"`
	To      *timestamp.Timestamp `protobuf:"bytes,4,opt,name=to,proto3,oneof" json:"to,omitempty"`
	PerDay  *bool                `protobuf:"varint,5,opt,name=per_day,json=perDay,proto3,oneof" json:"per_day,omitempty"`
}

func (x *ListTimeclockRequest) Reset() {
	*x = ListTimeclockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_jobs_timeclock_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListTimeclockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTimeclockRequest) ProtoMessage() {}

func (x *ListTimeclockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_jobs_timeclock_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTimeclockRequest.ProtoReflect.Descriptor instead.
func (*ListTimeclockRequest) Descriptor() ([]byte, []int) {
	return file_services_jobs_timeclock_proto_rawDescGZIP(), []int{0}
}

func (x *ListTimeclockRequest) GetPagination() *database.PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListTimeclockRequest) GetUserIds() []int32 {
	if x != nil {
		return x.UserIds
	}
	return nil
}

func (x *ListTimeclockRequest) GetFrom() *timestamp.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *ListTimeclockRequest) GetTo() *timestamp.Timestamp {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *ListTimeclockRequest) GetPerDay() bool {
	if x != nil && x.PerDay != nil {
		return *x.PerDay
	}
	return false
}

type ListTimeclockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *database.PaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Entries    []*jobs.TimeclockEntry       `protobuf:"bytes,2,rep,name=entries,proto3" json:"entries,omitempty"`
	Stats      *jobs.TimeclockStats         `protobuf:"bytes,3,opt,name=stats,proto3" json:"stats,omitempty"`
	Weekly     []*jobs.TimeclockWeeklyStats `protobuf:"bytes,4,rep,name=weekly,proto3" json:"weekly,omitempty"`
}

func (x *ListTimeclockResponse) Reset() {
	*x = ListTimeclockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_jobs_timeclock_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListTimeclockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTimeclockResponse) ProtoMessage() {}

func (x *ListTimeclockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_jobs_timeclock_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTimeclockResponse.ProtoReflect.Descriptor instead.
func (*ListTimeclockResponse) Descriptor() ([]byte, []int) {
	return file_services_jobs_timeclock_proto_rawDescGZIP(), []int{1}
}

func (x *ListTimeclockResponse) GetPagination() *database.PaginationResponse {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListTimeclockResponse) GetEntries() []*jobs.TimeclockEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

func (x *ListTimeclockResponse) GetStats() *jobs.TimeclockStats {
	if x != nil {
		return x.Stats
	}
	return nil
}

func (x *ListTimeclockResponse) GetWeekly() []*jobs.TimeclockWeeklyStats {
	if x != nil {
		return x.Weekly
	}
	return nil
}

type GetTimeclockStatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId *int32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
}

func (x *GetTimeclockStatsRequest) Reset() {
	*x = GetTimeclockStatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_jobs_timeclock_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTimeclockStatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTimeclockStatsRequest) ProtoMessage() {}

func (x *GetTimeclockStatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_jobs_timeclock_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTimeclockStatsRequest.ProtoReflect.Descriptor instead.
func (*GetTimeclockStatsRequest) Descriptor() ([]byte, []int) {
	return file_services_jobs_timeclock_proto_rawDescGZIP(), []int{2}
}

func (x *GetTimeclockStatsRequest) GetUserId() int32 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

type GetTimeclockStatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stats  *jobs.TimeclockStats         `protobuf:"bytes,1,opt,name=stats,proto3" json:"stats,omitempty"`
	Weekly []*jobs.TimeclockWeeklyStats `protobuf:"bytes,2,rep,name=weekly,proto3" json:"weekly,omitempty"`
}

func (x *GetTimeclockStatsResponse) Reset() {
	*x = GetTimeclockStatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_jobs_timeclock_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTimeclockStatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTimeclockStatsResponse) ProtoMessage() {}

func (x *GetTimeclockStatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_jobs_timeclock_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTimeclockStatsResponse.ProtoReflect.Descriptor instead.
func (*GetTimeclockStatsResponse) Descriptor() ([]byte, []int) {
	return file_services_jobs_timeclock_proto_rawDescGZIP(), []int{3}
}

func (x *GetTimeclockStatsResponse) GetStats() *jobs.TimeclockStats {
	if x != nil {
		return x.Stats
	}
	return nil
}

func (x *GetTimeclockStatsResponse) GetWeekly() []*jobs.TimeclockWeeklyStats {
	if x != nil {
		return x.Weekly
	}
	return nil
}

type ListInactiveEmployeesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *database.PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Days       int32                       `protobuf:"varint,2,opt,name=days,proto3" json:"days,omitempty"`
}

func (x *ListInactiveEmployeesRequest) Reset() {
	*x = ListInactiveEmployeesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_jobs_timeclock_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListInactiveEmployeesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListInactiveEmployeesRequest) ProtoMessage() {}

func (x *ListInactiveEmployeesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_jobs_timeclock_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListInactiveEmployeesRequest.ProtoReflect.Descriptor instead.
func (*ListInactiveEmployeesRequest) Descriptor() ([]byte, []int) {
	return file_services_jobs_timeclock_proto_rawDescGZIP(), []int{4}
}

func (x *ListInactiveEmployeesRequest) GetPagination() *database.PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListInactiveEmployeesRequest) GetDays() int32 {
	if x != nil {
		return x.Days
	}
	return 0
}

type ListInactiveEmployeesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *database.PaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Colleagues []*jobs.Colleague            `protobuf:"bytes,2,rep,name=colleagues,proto3" json:"colleagues,omitempty"`
}

func (x *ListInactiveEmployeesResponse) Reset() {
	*x = ListInactiveEmployeesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_jobs_timeclock_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListInactiveEmployeesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListInactiveEmployeesResponse) ProtoMessage() {}

func (x *ListInactiveEmployeesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_jobs_timeclock_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListInactiveEmployeesResponse.ProtoReflect.Descriptor instead.
func (*ListInactiveEmployeesResponse) Descriptor() ([]byte, []int) {
	return file_services_jobs_timeclock_proto_rawDescGZIP(), []int{5}
}

func (x *ListInactiveEmployeesResponse) GetPagination() *database.PaginationResponse {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListInactiveEmployeesResponse) GetColleagues() []*jobs.Colleague {
	if x != nil {
		return x.Colleagues
	}
	return nil
}

var File_services_jobs_timeclock_proto protoreflect.FileDescriptor

var file_services_jobs_timeclock_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x1a, 0x28,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x61, 0x67,
	0x75, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x63, 0x6c,
	0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb1, 0x02, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x56, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x70, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x05, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x73, 0x12, 0x37, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x48, 0x00, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a, 0x02,
	0x74, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52, 0x02, 0x74, 0x6f, 0x88, 0x01,
	0x01, 0x12, 0x1c, 0x0a, 0x07, 0x70, 0x65, 0x72, 0x5f, 0x64, 0x61, 0x79, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x48, 0x02, 0x52, 0x06, 0x70, 0x65, 0x72, 0x44, 0x61, 0x79, 0x88, 0x01, 0x01, 0x42,
	0x07, 0x0a, 0x05, 0x5f, 0x66, 0x72, 0x6f, 0x6d, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x74, 0x6f, 0x42,
	0x0a, 0x0a, 0x08, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x64, 0x61, 0x79, 0x22, 0x94, 0x02, 0x0a, 0x15,
	0x4c, 0x69, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x34,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x73, 0x12, 0x3c, 0x0a, 0x06, 0x77, 0x65, 0x65, 0x6b, 0x6c, 0x79, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x57,
	0x65, 0x65, 0x6b, 0x6c, 0x79, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x06, 0x77, 0x65, 0x65, 0x6b,
	0x6c, 0x79, 0x22, 0x44, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f,
	0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x48,
	0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x22, 0x8f, 0x01, 0x0a, 0x19, 0x47, 0x65, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x12, 0x3c, 0x0a, 0x06,
	0x77, 0x65, 0x65, 0x6b, 0x6c, 0x79, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x57, 0x65, 0x65, 0x6b, 0x6c, 0x79, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x52, 0x06, 0x77, 0x65, 0x65, 0x6b, 0x6c, 0x79, 0x22, 0x93, 0x01, 0x0a, 0x1c, 0x4c,
	0x69, 0x73, 0x74, 0x49, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x45, 0x6d, 0x70, 0x6c, 0x6f,
	0x79, 0x65, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x56, 0x0a, 0x0a, 0x70,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x04, 0x64, 0x61, 0x79, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x28, 0x01, 0x52, 0x04, 0x64, 0x61, 0x79, 0x73,
	0x22, 0xa9, 0x01, 0x0a, 0x1d, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x61, 0x67, 0x75, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x61, 0x67, 0x75, 0x65,
	0x52, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x61, 0x67, 0x75, 0x65, 0x73, 0x32, 0xce, 0x02, 0x0a,
	0x14, 0x4a, 0x6f, 0x62, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x23, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x63,
	0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x66, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63,
	0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x27, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c,
	0x6f, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x28, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x72, 0x0a, 0x15, 0x4c, 0x69, 0x73,
	0x74, 0x49, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65,
	0x65, 0x73, 0x12, 0x2b, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f,
	0x62, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x45,
	0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x45, 0x6d, 0x70, 0x6c,
	0x6f, 0x79, 0x65, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x40, 0x5a,
	0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65,
	0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x3b, 0x6a, 0x6f, 0x62, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_jobs_timeclock_proto_rawDescOnce sync.Once
	file_services_jobs_timeclock_proto_rawDescData = file_services_jobs_timeclock_proto_rawDesc
)

func file_services_jobs_timeclock_proto_rawDescGZIP() []byte {
	file_services_jobs_timeclock_proto_rawDescOnce.Do(func() {
		file_services_jobs_timeclock_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_jobs_timeclock_proto_rawDescData)
	})
	return file_services_jobs_timeclock_proto_rawDescData
}

var file_services_jobs_timeclock_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_services_jobs_timeclock_proto_goTypes = []interface{}{
	(*ListTimeclockRequest)(nil),          // 0: services.jobs.ListTimeclockRequest
	(*ListTimeclockResponse)(nil),         // 1: services.jobs.ListTimeclockResponse
	(*GetTimeclockStatsRequest)(nil),      // 2: services.jobs.GetTimeclockStatsRequest
	(*GetTimeclockStatsResponse)(nil),     // 3: services.jobs.GetTimeclockStatsResponse
	(*ListInactiveEmployeesRequest)(nil),  // 4: services.jobs.ListInactiveEmployeesRequest
	(*ListInactiveEmployeesResponse)(nil), // 5: services.jobs.ListInactiveEmployeesResponse
	(*database.PaginationRequest)(nil),    // 6: resources.common.database.PaginationRequest
	(*timestamp.Timestamp)(nil),           // 7: resources.timestamp.Timestamp
	(*database.PaginationResponse)(nil),   // 8: resources.common.database.PaginationResponse
	(*jobs.TimeclockEntry)(nil),           // 9: resources.jobs.TimeclockEntry
	(*jobs.TimeclockStats)(nil),           // 10: resources.jobs.TimeclockStats
	(*jobs.TimeclockWeeklyStats)(nil),     // 11: resources.jobs.TimeclockWeeklyStats
	(*jobs.Colleague)(nil),                // 12: resources.jobs.Colleague
}
var file_services_jobs_timeclock_proto_depIdxs = []int32{
	6,  // 0: services.jobs.ListTimeclockRequest.pagination:type_name -> resources.common.database.PaginationRequest
	7,  // 1: services.jobs.ListTimeclockRequest.from:type_name -> resources.timestamp.Timestamp
	7,  // 2: services.jobs.ListTimeclockRequest.to:type_name -> resources.timestamp.Timestamp
	8,  // 3: services.jobs.ListTimeclockResponse.pagination:type_name -> resources.common.database.PaginationResponse
	9,  // 4: services.jobs.ListTimeclockResponse.entries:type_name -> resources.jobs.TimeclockEntry
	10, // 5: services.jobs.ListTimeclockResponse.stats:type_name -> resources.jobs.TimeclockStats
	11, // 6: services.jobs.ListTimeclockResponse.weekly:type_name -> resources.jobs.TimeclockWeeklyStats
	10, // 7: services.jobs.GetTimeclockStatsResponse.stats:type_name -> resources.jobs.TimeclockStats
	11, // 8: services.jobs.GetTimeclockStatsResponse.weekly:type_name -> resources.jobs.TimeclockWeeklyStats
	6,  // 9: services.jobs.ListInactiveEmployeesRequest.pagination:type_name -> resources.common.database.PaginationRequest
	8,  // 10: services.jobs.ListInactiveEmployeesResponse.pagination:type_name -> resources.common.database.PaginationResponse
	12, // 11: services.jobs.ListInactiveEmployeesResponse.colleagues:type_name -> resources.jobs.Colleague
	0,  // 12: services.jobs.JobsTimeclockService.ListTimeclock:input_type -> services.jobs.ListTimeclockRequest
	2,  // 13: services.jobs.JobsTimeclockService.GetTimeclockStats:input_type -> services.jobs.GetTimeclockStatsRequest
	4,  // 14: services.jobs.JobsTimeclockService.ListInactiveEmployees:input_type -> services.jobs.ListInactiveEmployeesRequest
	1,  // 15: services.jobs.JobsTimeclockService.ListTimeclock:output_type -> services.jobs.ListTimeclockResponse
	3,  // 16: services.jobs.JobsTimeclockService.GetTimeclockStats:output_type -> services.jobs.GetTimeclockStatsResponse
	5,  // 17: services.jobs.JobsTimeclockService.ListInactiveEmployees:output_type -> services.jobs.ListInactiveEmployeesResponse
	15, // [15:18] is the sub-list for method output_type
	12, // [12:15] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_services_jobs_timeclock_proto_init() }
func file_services_jobs_timeclock_proto_init() {
	if File_services_jobs_timeclock_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_jobs_timeclock_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListTimeclockRequest); i {
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
		file_services_jobs_timeclock_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListTimeclockResponse); i {
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
		file_services_jobs_timeclock_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTimeclockStatsRequest); i {
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
		file_services_jobs_timeclock_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTimeclockStatsResponse); i {
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
		file_services_jobs_timeclock_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListInactiveEmployeesRequest); i {
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
		file_services_jobs_timeclock_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListInactiveEmployeesResponse); i {
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
	file_services_jobs_timeclock_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_services_jobs_timeclock_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_jobs_timeclock_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_jobs_timeclock_proto_goTypes,
		DependencyIndexes: file_services_jobs_timeclock_proto_depIdxs,
		MessageInfos:      file_services_jobs_timeclock_proto_msgTypes,
	}.Build()
	File_services_jobs_timeclock_proto = out.File
	file_services_jobs_timeclock_proto_rawDesc = nil
	file_services_jobs_timeclock_proto_goTypes = nil
	file_services_jobs_timeclock_proto_depIdxs = nil
}
