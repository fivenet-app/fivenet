// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: resources/jobs/timeclock.proto

package jobs

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

type TimeclockEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job       string               `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
	Date      *timestamp.Timestamp `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	UserId    int32                `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	User      *users.UserShort     `protobuf:"bytes,4,opt,name=user,proto3,oneof" json:"user,omitempty"`
	StartTime *timestamp.Timestamp `protobuf:"bytes,5,opt,name=start_time,json=startTime,proto3,oneof" json:"start_time,omitempty"`
	EndTime   *timestamp.Timestamp `protobuf:"bytes,6,opt,name=end_time,json=endTime,proto3,oneof" json:"end_time,omitempty"`
	SpentTime float32              `protobuf:"fixed32,7,opt,name=spent_time,json=spentTime,proto3" json:"spent_time,omitempty"`
}

func (x *TimeclockEntry) Reset() {
	*x = TimeclockEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_jobs_timeclock_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeclockEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeclockEntry) ProtoMessage() {}

func (x *TimeclockEntry) ProtoReflect() protoreflect.Message {
	mi := &file_resources_jobs_timeclock_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeclockEntry.ProtoReflect.Descriptor instead.
func (*TimeclockEntry) Descriptor() ([]byte, []int) {
	return file_resources_jobs_timeclock_proto_rawDescGZIP(), []int{0}
}

func (x *TimeclockEntry) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *TimeclockEntry) GetDate() *timestamp.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *TimeclockEntry) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *TimeclockEntry) GetUser() *users.UserShort {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *TimeclockEntry) GetStartTime() *timestamp.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *TimeclockEntry) GetEndTime() *timestamp.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *TimeclockEntry) GetSpentTime() float32 {
	if x != nil {
		return x.SpentTime
	}
	return 0
}

type TimeclockStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job          string  `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
	SpentTimeSum float32 `protobuf:"fixed32,2,opt,name=spent_time_sum,json=spentTimeSum,proto3" json:"spent_time_sum,omitempty"`
	SpentTimeAvg float32 `protobuf:"fixed32,3,opt,name=spent_time_avg,json=spentTimeAvg,proto3" json:"spent_time_avg,omitempty"`
	SpentTimeMax float32 `protobuf:"fixed32,4,opt,name=spent_time_max,json=spentTimeMax,proto3" json:"spent_time_max,omitempty"`
}

func (x *TimeclockStats) Reset() {
	*x = TimeclockStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_jobs_timeclock_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeclockStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeclockStats) ProtoMessage() {}

func (x *TimeclockStats) ProtoReflect() protoreflect.Message {
	mi := &file_resources_jobs_timeclock_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeclockStats.ProtoReflect.Descriptor instead.
func (*TimeclockStats) Descriptor() ([]byte, []int) {
	return file_resources_jobs_timeclock_proto_rawDescGZIP(), []int{1}
}

func (x *TimeclockStats) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *TimeclockStats) GetSpentTimeSum() float32 {
	if x != nil {
		return x.SpentTimeSum
	}
	return 0
}

func (x *TimeclockStats) GetSpentTimeAvg() float32 {
	if x != nil {
		return x.SpentTimeAvg
	}
	return 0
}

func (x *TimeclockStats) GetSpentTimeMax() float32 {
	if x != nil {
		return x.SpentTimeMax
	}
	return 0
}

type TimeclockWeeklyStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Date string  `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	Sum  float32 `protobuf:"fixed32,2,opt,name=sum,proto3" json:"sum,omitempty"`
	Avg  float32 `protobuf:"fixed32,3,opt,name=avg,proto3" json:"avg,omitempty"`
	Max  float32 `protobuf:"fixed32,4,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *TimeclockWeeklyStats) Reset() {
	*x = TimeclockWeeklyStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_jobs_timeclock_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeclockWeeklyStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeclockWeeklyStats) ProtoMessage() {}

func (x *TimeclockWeeklyStats) ProtoReflect() protoreflect.Message {
	mi := &file_resources_jobs_timeclock_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeclockWeeklyStats.ProtoReflect.Descriptor instead.
func (*TimeclockWeeklyStats) Descriptor() ([]byte, []int) {
	return file_resources_jobs_timeclock_proto_rawDescGZIP(), []int{2}
}

func (x *TimeclockWeeklyStats) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *TimeclockWeeklyStats) GetSum() float32 {
	if x != nil {
		return x.Sum
	}
	return 0
}

func (x *TimeclockWeeklyStats) GetAvg() float32 {
	if x != nil {
		return x.Avg
	}
	return 0
}

func (x *TimeclockWeeklyStats) GetMax() float32 {
	if x != nil {
		return x.Max
	}
	return 0
}

var File_resources_jobs_timeclock_proto protoreflect.FileDescriptor

var file_resources_jobs_timeclock_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73,
	0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf5, 0x02, 0x0a, 0x0e,
	0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x19,
	0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x18, 0x14, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x32, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x48, 0x00, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48,
	0x01, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x3e, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x48, 0x02, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x1d, 0x0a, 0x0a, 0x73, 0x70, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x09, 0x73, 0x70, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x07,
	0x0a, 0x05, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x22, 0x9d, 0x01, 0x0a, 0x0e, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63,
	0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52, 0x03, 0x6a, 0x6f,
	0x62, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x70, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f,
	0x73, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0c, 0x73, 0x70, 0x65, 0x6e, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x53, 0x75, 0x6d, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x70, 0x65, 0x6e, 0x74,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x61, 0x76, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x0c, 0x73, 0x70, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x41, 0x76, 0x67, 0x12, 0x24, 0x0a,
	0x0e, 0x73, 0x70, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x6d, 0x61, 0x78, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0c, 0x73, 0x70, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x4d, 0x61, 0x78, 0x22, 0x60, 0x0a, 0x14, 0x54, 0x69, 0x6d, 0x65, 0x63, 0x6c, 0x6f, 0x63, 0x6b,
	0x57, 0x65, 0x65, 0x6b, 0x6c, 0x79, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x73, 0x75,
	0x6d, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x76, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03,
	0x61, 0x76, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x03, 0x6d, 0x61, 0x78, 0x42, 0x3d, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x66, 0x69, 0x76, 0x65,
	0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x3b,
	0x6a, 0x6f, 0x62, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_jobs_timeclock_proto_rawDescOnce sync.Once
	file_resources_jobs_timeclock_proto_rawDescData = file_resources_jobs_timeclock_proto_rawDesc
)

func file_resources_jobs_timeclock_proto_rawDescGZIP() []byte {
	file_resources_jobs_timeclock_proto_rawDescOnce.Do(func() {
		file_resources_jobs_timeclock_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_jobs_timeclock_proto_rawDescData)
	})
	return file_resources_jobs_timeclock_proto_rawDescData
}

var file_resources_jobs_timeclock_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_resources_jobs_timeclock_proto_goTypes = []interface{}{
	(*TimeclockEntry)(nil),       // 0: resources.jobs.TimeclockEntry
	(*TimeclockStats)(nil),       // 1: resources.jobs.TimeclockStats
	(*TimeclockWeeklyStats)(nil), // 2: resources.jobs.TimeclockWeeklyStats
	(*timestamp.Timestamp)(nil),  // 3: resources.timestamp.Timestamp
	(*users.UserShort)(nil),      // 4: resources.users.UserShort
}
var file_resources_jobs_timeclock_proto_depIdxs = []int32{
	3, // 0: resources.jobs.TimeclockEntry.date:type_name -> resources.timestamp.Timestamp
	4, // 1: resources.jobs.TimeclockEntry.user:type_name -> resources.users.UserShort
	3, // 2: resources.jobs.TimeclockEntry.start_time:type_name -> resources.timestamp.Timestamp
	3, // 3: resources.jobs.TimeclockEntry.end_time:type_name -> resources.timestamp.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_resources_jobs_timeclock_proto_init() }
func file_resources_jobs_timeclock_proto_init() {
	if File_resources_jobs_timeclock_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_jobs_timeclock_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeclockEntry); i {
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
		file_resources_jobs_timeclock_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeclockStats); i {
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
		file_resources_jobs_timeclock_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeclockWeeklyStats); i {
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
	file_resources_jobs_timeclock_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_jobs_timeclock_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_jobs_timeclock_proto_goTypes,
		DependencyIndexes: file_resources_jobs_timeclock_proto_depIdxs,
		MessageInfos:      file_resources_jobs_timeclock_proto_msgTypes,
	}.Build()
	File_resources_jobs_timeclock_proto = out.File
	file_resources_jobs_timeclock_proto_rawDesc = nil
	file_resources_jobs_timeclock_proto_goTypes = nil
	file_resources_jobs_timeclock_proto_depIdxs = nil
}
