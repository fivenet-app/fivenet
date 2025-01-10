// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v3.20.3
// source: services/sync/sync.proto

package sync

import (
	jobs "github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
	sync "github.com/fivenet-app/fivenet/gen/go/proto/resources/sync"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync1 "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetStatusRequest) Reset() {
	*x = GetStatusRequest{}
	mi := &file_services_sync_sync_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatusRequest) ProtoMessage() {}

func (x *GetStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_sync_sync_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatusRequest.ProtoReflect.Descriptor instead.
func (*GetStatusRequest) Descriptor() ([]byte, []int) {
	return file_services_sync_sync_proto_rawDescGZIP(), []int{0}
}

type GetStatusResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Jobs          *sync.DataStatus       `protobuf:"bytes,1,opt,name=jobs,proto3" json:"jobs,omitempty"`
	Licenses      *sync.DataStatus       `protobuf:"bytes,2,opt,name=licenses,proto3" json:"licenses,omitempty"`
	Users         *sync.DataStatus       `protobuf:"bytes,3,opt,name=users,proto3" json:"users,omitempty"`
	Vehicles      *sync.DataStatus       `protobuf:"bytes,4,opt,name=vehicles,proto3" json:"vehicles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetStatusResponse) Reset() {
	*x = GetStatusResponse{}
	mi := &file_services_sync_sync_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatusResponse) ProtoMessage() {}

func (x *GetStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_sync_sync_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatusResponse.ProtoReflect.Descriptor instead.
func (*GetStatusResponse) Descriptor() ([]byte, []int) {
	return file_services_sync_sync_proto_rawDescGZIP(), []int{1}
}

func (x *GetStatusResponse) GetJobs() *sync.DataStatus {
	if x != nil {
		return x.Jobs
	}
	return nil
}

func (x *GetStatusResponse) GetLicenses() *sync.DataStatus {
	if x != nil {
		return x.Licenses
	}
	return nil
}

func (x *GetStatusResponse) GetUsers() *sync.DataStatus {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *GetStatusResponse) GetVehicles() *sync.DataStatus {
	if x != nil {
		return x.Vehicles
	}
	return nil
}

type AddActivityRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Activity:
	//
	//	*AddActivityRequest_UserOauth2
	//	*AddActivityRequest_UserActivity
	//	*AddActivityRequest_UserProps
	//	*AddActivityRequest_JobsUserActivity
	//	*AddActivityRequest_JobsUserProps
	//	*AddActivityRequest_JobsTimeclock
	Activity      isAddActivityRequest_Activity `protobuf_oneof:"activity"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddActivityRequest) Reset() {
	*x = AddActivityRequest{}
	mi := &file_services_sync_sync_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddActivityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddActivityRequest) ProtoMessage() {}

func (x *AddActivityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_sync_sync_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddActivityRequest.ProtoReflect.Descriptor instead.
func (*AddActivityRequest) Descriptor() ([]byte, []int) {
	return file_services_sync_sync_proto_rawDescGZIP(), []int{2}
}

func (x *AddActivityRequest) GetActivity() isAddActivityRequest_Activity {
	if x != nil {
		return x.Activity
	}
	return nil
}

func (x *AddActivityRequest) GetUserOauth2() *sync.UserOAuth2Conn {
	if x != nil {
		if x, ok := x.Activity.(*AddActivityRequest_UserOauth2); ok {
			return x.UserOauth2
		}
	}
	return nil
}

func (x *AddActivityRequest) GetUserActivity() *users.UserActivity {
	if x != nil {
		if x, ok := x.Activity.(*AddActivityRequest_UserActivity); ok {
			return x.UserActivity
		}
	}
	return nil
}

func (x *AddActivityRequest) GetUserProps() *sync.UserProps {
	if x != nil {
		if x, ok := x.Activity.(*AddActivityRequest_UserProps); ok {
			return x.UserProps
		}
	}
	return nil
}

func (x *AddActivityRequest) GetJobsUserActivity() *jobs.JobsUserActivity {
	if x != nil {
		if x, ok := x.Activity.(*AddActivityRequest_JobsUserActivity); ok {
			return x.JobsUserActivity
		}
	}
	return nil
}

func (x *AddActivityRequest) GetJobsUserProps() *sync.JobsUserProps {
	if x != nil {
		if x, ok := x.Activity.(*AddActivityRequest_JobsUserProps); ok {
			return x.JobsUserProps
		}
	}
	return nil
}

func (x *AddActivityRequest) GetJobsTimeclock() *jobs.TimeclockEntry {
	if x != nil {
		if x, ok := x.Activity.(*AddActivityRequest_JobsTimeclock); ok {
			return x.JobsTimeclock
		}
	}
	return nil
}

type isAddActivityRequest_Activity interface {
	isAddActivityRequest_Activity()
}

type AddActivityRequest_UserOauth2 struct {
	UserOauth2 *sync.UserOAuth2Conn `protobuf:"bytes,1,opt,name=user_oauth2,json=userOauth2,proto3,oneof"`
}

type AddActivityRequest_UserActivity struct {
	// User activity
	UserActivity *users.UserActivity `protobuf:"bytes,2,opt,name=user_activity,json=userActivity,proto3,oneof"`
}

type AddActivityRequest_UserProps struct {
	// Setting props will cause activity to be created automtically
	UserProps *sync.UserProps `protobuf:"bytes,3,opt,name=user_props,json=userProps,proto3,oneof"`
}

type AddActivityRequest_JobsUserActivity struct {
	// Jobs user activity
	JobsUserActivity *jobs.JobsUserActivity `protobuf:"bytes,4,opt,name=jobs_user_activity,json=jobsUserActivity,proto3,oneof"`
}

type AddActivityRequest_JobsUserProps struct {
	// Setting props will cause activity to be created automtically
	JobsUserProps *sync.JobsUserProps `protobuf:"bytes,5,opt,name=jobs_user_props,json=jobsUserProps,proto3,oneof"`
}

type AddActivityRequest_JobsTimeclock struct {
	// Timeclock user entry
	JobsTimeclock *jobs.TimeclockEntry `protobuf:"bytes,6,opt,name=jobs_timeclock,json=jobsTimeclock,proto3,oneof"`
}

func (*AddActivityRequest_UserOauth2) isAddActivityRequest_Activity() {}

func (*AddActivityRequest_UserActivity) isAddActivityRequest_Activity() {}

func (*AddActivityRequest_UserProps) isAddActivityRequest_Activity() {}

func (*AddActivityRequest_JobsUserActivity) isAddActivityRequest_Activity() {}

func (*AddActivityRequest_JobsUserProps) isAddActivityRequest_Activity() {}

func (*AddActivityRequest_JobsTimeclock) isAddActivityRequest_Activity() {}

type AddActivityResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddActivityResponse) Reset() {
	*x = AddActivityResponse{}
	mi := &file_services_sync_sync_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddActivityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddActivityResponse) ProtoMessage() {}

func (x *AddActivityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_sync_sync_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddActivityResponse.ProtoReflect.Descriptor instead.
func (*AddActivityResponse) Descriptor() ([]byte, []int) {
	return file_services_sync_sync_proto_rawDescGZIP(), []int{3}
}

func (x *AddActivityResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type SendDataRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*SendDataRequest_Jobs
	//	*SendDataRequest_Licenses
	//	*SendDataRequest_Users
	//	*SendDataRequest_Vehicles
	Data          isSendDataRequest_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendDataRequest) Reset() {
	*x = SendDataRequest{}
	mi := &file_services_sync_sync_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendDataRequest) ProtoMessage() {}

func (x *SendDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_sync_sync_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendDataRequest.ProtoReflect.Descriptor instead.
func (*SendDataRequest) Descriptor() ([]byte, []int) {
	return file_services_sync_sync_proto_rawDescGZIP(), []int{4}
}

func (x *SendDataRequest) GetData() isSendDataRequest_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *SendDataRequest) GetJobs() *sync.DataJobs {
	if x != nil {
		if x, ok := x.Data.(*SendDataRequest_Jobs); ok {
			return x.Jobs
		}
	}
	return nil
}

func (x *SendDataRequest) GetLicenses() *sync.DataLicenses {
	if x != nil {
		if x, ok := x.Data.(*SendDataRequest_Licenses); ok {
			return x.Licenses
		}
	}
	return nil
}

func (x *SendDataRequest) GetUsers() *sync.DataUsers {
	if x != nil {
		if x, ok := x.Data.(*SendDataRequest_Users); ok {
			return x.Users
		}
	}
	return nil
}

func (x *SendDataRequest) GetVehicles() *sync.DataVehicles {
	if x != nil {
		if x, ok := x.Data.(*SendDataRequest_Vehicles); ok {
			return x.Vehicles
		}
	}
	return nil
}

type isSendDataRequest_Data interface {
	isSendDataRequest_Data()
}

type SendDataRequest_Jobs struct {
	Jobs *sync.DataJobs `protobuf:"bytes,1,opt,name=jobs,proto3,oneof"`
}

type SendDataRequest_Licenses struct {
	Licenses *sync.DataLicenses `protobuf:"bytes,2,opt,name=licenses,proto3,oneof"`
}

type SendDataRequest_Users struct {
	Users *sync.DataUsers `protobuf:"bytes,3,opt,name=users,proto3,oneof"`
}

type SendDataRequest_Vehicles struct {
	Vehicles *sync.DataVehicles `protobuf:"bytes,4,opt,name=vehicles,proto3,oneof"`
}

func (*SendDataRequest_Jobs) isSendDataRequest_Data() {}

func (*SendDataRequest_Licenses) isSendDataRequest_Data() {}

func (*SendDataRequest_Users) isSendDataRequest_Data() {}

func (*SendDataRequest_Vehicles) isSendDataRequest_Data() {}

type SendDataResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AffectedRows  int64                  `protobuf:"varint,1,opt,name=affected_rows,json=affectedRows,proto3" json:"affected_rows,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendDataResponse) Reset() {
	*x = SendDataResponse{}
	mi := &file_services_sync_sync_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendDataResponse) ProtoMessage() {}

func (x *SendDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_sync_sync_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendDataResponse.ProtoReflect.Descriptor instead.
func (*SendDataResponse) Descriptor() ([]byte, []int) {
	return file_services_sync_sync_proto_rawDescGZIP(), []int{5}
}

func (x *SendDataResponse) GetAffectedRows() int64 {
	if x != nil {
		return x.AffectedRows
	}
	return 0
}

type StreamRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	mi := &file_services_sync_sync_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_sync_sync_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRequest.ProtoReflect.Descriptor instead.
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return file_services_sync_sync_proto_rawDescGZIP(), []int{6}
}

type StreamResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int32                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamResponse) Reset() {
	*x = StreamResponse{}
	mi := &file_services_sync_sync_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResponse) ProtoMessage() {}

func (x *StreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_sync_sync_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamResponse.ProtoReflect.Descriptor instead.
func (*StreamResponse) Descriptor() ([]byte, []int) {
	return file_services_sync_sync_proto_rawDescGZIP(), []int{7}
}

func (x *StreamResponse) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

var File_services_sync_sync_proto protoreflect.FileDescriptor

var file_services_sync_sync_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f,
	0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x1a, 0x1d, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f,
	0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x12, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xe5, 0x01, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x04,
	0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x12, 0x36, 0x0a, 0x08,
	0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e,
	0x44, 0x61, 0x74, 0x61, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x08, 0x6c, 0x69, 0x63, 0x65,
	0x6e, 0x73, 0x65, 0x73, 0x12, 0x30, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x73, 0x79, 0x6e, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x36, 0x0a, 0x08, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c,
	0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x08, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x22, 0xc9,
	0x03, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x41, 0x0a, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6f, 0x61,
	0x75, 0x74, 0x68, 0x32, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x43, 0x6f, 0x6e, 0x6e, 0x48, 0x00, 0x52, 0x0a, 0x75, 0x73,
	0x65, 0x72, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x12, 0x44, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x48, 0x00,
	0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x12, 0x3a,
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x73,
	0x79, 0x6e, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x48, 0x00, 0x52,
	0x09, 0x75, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x50, 0x0a, 0x12, 0x6a, 0x6f,
	0x62, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x73, 0x55, 0x73, 0x65, 0x72,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x48, 0x00, 0x52, 0x10, 0x6a, 0x6f, 0x62, 0x73,
	0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x12, 0x47, 0x0a, 0x0f,
	0x6a, 0x6f, 0x62, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x4a, 0x6f, 0x62, 0x73, 0x55, 0x73, 0x65, 0x72, 0x50,
	0x72, 0x6f, 0x70, 0x73, 0x48, 0x00, 0x52, 0x0d, 0x6a, 0x6f, 0x62, 0x73, 0x55, 0x73, 0x65, 0x72,
	0x50, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x47, 0x0a, 0x0e, 0x6a, 0x6f, 0x62, 0x73, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x48, 0x00, 0x52,
	0x0d, 0x6a, 0x6f, 0x62, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x42, 0x0a,
	0x0a, 0x08, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x22, 0x25, 0x0a, 0x13, 0x41, 0x64,
	0x64, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69,
	0x64, 0x22, 0xf4, 0x01, 0x0a, 0x0f, 0x53, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x73, 0x79, 0x6e, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x4a, 0x6f, 0x62, 0x73, 0x48, 0x00, 0x52,
	0x04, 0x6a, 0x6f, 0x62, 0x73, 0x12, 0x3a, 0x0a, 0x08, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x63,
	0x65, 0x6e, 0x73, 0x65, 0x73, 0x48, 0x00, 0x52, 0x08, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65,
	0x73, 0x12, 0x31, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e,
	0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x55, 0x73, 0x65, 0x72, 0x73, 0x48, 0x00, 0x52, 0x05, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x12, 0x3a, 0x0a, 0x08, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x56, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x73, 0x48, 0x00, 0x52, 0x08, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73,
	0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x37, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d,
	0x61, 0x66, 0x66, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x6f, 0x77, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0c, 0x61, 0x66, 0x66, 0x65, 0x63, 0x74, 0x65, 0x64, 0x52, 0x6f, 0x77,
	0x73, 0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x29, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x32, 0xc9, 0x02,
	0x0a, 0x0b, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a,
	0x09, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a,
	0x0b, 0x41, 0x64, 0x64, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x12, 0x21, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x41, 0x64, 0x64,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x22, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e,
	0x41, 0x64, 0x64, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x4b, 0x0a, 0x08, 0x53, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x1e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x47, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1c, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d,
	0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x3b, 0x73, 0x79, 0x6e, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_services_sync_sync_proto_rawDescOnce sync1.Once
	file_services_sync_sync_proto_rawDescData = file_services_sync_sync_proto_rawDesc
)

func file_services_sync_sync_proto_rawDescGZIP() []byte {
	file_services_sync_sync_proto_rawDescOnce.Do(func() {
		file_services_sync_sync_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_sync_sync_proto_rawDescData)
	})
	return file_services_sync_sync_proto_rawDescData
}

var file_services_sync_sync_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_services_sync_sync_proto_goTypes = []any{
	(*GetStatusRequest)(nil),      // 0: services.sync.GetStatusRequest
	(*GetStatusResponse)(nil),     // 1: services.sync.GetStatusResponse
	(*AddActivityRequest)(nil),    // 2: services.sync.AddActivityRequest
	(*AddActivityResponse)(nil),   // 3: services.sync.AddActivityResponse
	(*SendDataRequest)(nil),       // 4: services.sync.SendDataRequest
	(*SendDataResponse)(nil),      // 5: services.sync.SendDataResponse
	(*StreamRequest)(nil),         // 6: services.sync.StreamRequest
	(*StreamResponse)(nil),        // 7: services.sync.StreamResponse
	(*sync.DataStatus)(nil),       // 8: resources.sync.DataStatus
	(*sync.UserOAuth2Conn)(nil),   // 9: resources.sync.UserOAuth2Conn
	(*users.UserActivity)(nil),    // 10: resources.users.UserActivity
	(*sync.UserProps)(nil),        // 11: resources.sync.UserProps
	(*jobs.JobsUserActivity)(nil), // 12: resources.jobs.JobsUserActivity
	(*sync.JobsUserProps)(nil),    // 13: resources.sync.JobsUserProps
	(*jobs.TimeclockEntry)(nil),   // 14: resources.jobs.TimeclockEntry
	(*sync.DataJobs)(nil),         // 15: resources.sync.DataJobs
	(*sync.DataLicenses)(nil),     // 16: resources.sync.DataLicenses
	(*sync.DataUsers)(nil),        // 17: resources.sync.DataUsers
	(*sync.DataVehicles)(nil),     // 18: resources.sync.DataVehicles
}
var file_services_sync_sync_proto_depIdxs = []int32{
	8,  // 0: services.sync.GetStatusResponse.jobs:type_name -> resources.sync.DataStatus
	8,  // 1: services.sync.GetStatusResponse.licenses:type_name -> resources.sync.DataStatus
	8,  // 2: services.sync.GetStatusResponse.users:type_name -> resources.sync.DataStatus
	8,  // 3: services.sync.GetStatusResponse.vehicles:type_name -> resources.sync.DataStatus
	9,  // 4: services.sync.AddActivityRequest.user_oauth2:type_name -> resources.sync.UserOAuth2Conn
	10, // 5: services.sync.AddActivityRequest.user_activity:type_name -> resources.users.UserActivity
	11, // 6: services.sync.AddActivityRequest.user_props:type_name -> resources.sync.UserProps
	12, // 7: services.sync.AddActivityRequest.jobs_user_activity:type_name -> resources.jobs.JobsUserActivity
	13, // 8: services.sync.AddActivityRequest.jobs_user_props:type_name -> resources.sync.JobsUserProps
	14, // 9: services.sync.AddActivityRequest.jobs_timeclock:type_name -> resources.jobs.TimeclockEntry
	15, // 10: services.sync.SendDataRequest.jobs:type_name -> resources.sync.DataJobs
	16, // 11: services.sync.SendDataRequest.licenses:type_name -> resources.sync.DataLicenses
	17, // 12: services.sync.SendDataRequest.users:type_name -> resources.sync.DataUsers
	18, // 13: services.sync.SendDataRequest.vehicles:type_name -> resources.sync.DataVehicles
	0,  // 14: services.sync.SyncService.GetStatus:input_type -> services.sync.GetStatusRequest
	2,  // 15: services.sync.SyncService.AddActivity:input_type -> services.sync.AddActivityRequest
	4,  // 16: services.sync.SyncService.SendData:input_type -> services.sync.SendDataRequest
	6,  // 17: services.sync.SyncService.Stream:input_type -> services.sync.StreamRequest
	1,  // 18: services.sync.SyncService.GetStatus:output_type -> services.sync.GetStatusResponse
	3,  // 19: services.sync.SyncService.AddActivity:output_type -> services.sync.AddActivityResponse
	5,  // 20: services.sync.SyncService.SendData:output_type -> services.sync.SendDataResponse
	7,  // 21: services.sync.SyncService.Stream:output_type -> services.sync.StreamResponse
	18, // [18:22] is the sub-list for method output_type
	14, // [14:18] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_services_sync_sync_proto_init() }
func file_services_sync_sync_proto_init() {
	if File_services_sync_sync_proto != nil {
		return
	}
	file_services_sync_sync_proto_msgTypes[2].OneofWrappers = []any{
		(*AddActivityRequest_UserOauth2)(nil),
		(*AddActivityRequest_UserActivity)(nil),
		(*AddActivityRequest_UserProps)(nil),
		(*AddActivityRequest_JobsUserActivity)(nil),
		(*AddActivityRequest_JobsUserProps)(nil),
		(*AddActivityRequest_JobsTimeclock)(nil),
	}
	file_services_sync_sync_proto_msgTypes[4].OneofWrappers = []any{
		(*SendDataRequest_Jobs)(nil),
		(*SendDataRequest_Licenses)(nil),
		(*SendDataRequest_Users)(nil),
		(*SendDataRequest_Vehicles)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_sync_sync_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_sync_sync_proto_goTypes,
		DependencyIndexes: file_services_sync_sync_proto_depIdxs,
		MessageInfos:      file_services_sync_sync_proto_msgTypes,
	}.Build()
	File_services_sync_sync_proto = out.File
	file_services_sync_sync_proto_rawDesc = nil
	file_services_sync_sync_proto_goTypes = nil
	file_services_sync_sync_proto_depIdxs = nil
}
