// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: resources/calendar/access.proto

package calendar

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
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

type AccessLevel int32

const (
	AccessLevel_ACCESS_LEVEL_UNSPECIFIED AccessLevel = 0
	AccessLevel_ACCESS_LEVEL_BLOCKED     AccessLevel = 1
	AccessLevel_ACCESS_LEVEL_VIEW        AccessLevel = 2
	AccessLevel_ACCESS_LEVEL_SHARE       AccessLevel = 3
	AccessLevel_ACCESS_LEVEL_EDIT        AccessLevel = 4
	AccessLevel_ACCESS_LEVEL_MANAGE      AccessLevel = 5
)

// Enum value maps for AccessLevel.
var (
	AccessLevel_name = map[int32]string{
		0: "ACCESS_LEVEL_UNSPECIFIED",
		1: "ACCESS_LEVEL_BLOCKED",
		2: "ACCESS_LEVEL_VIEW",
		3: "ACCESS_LEVEL_SHARE",
		4: "ACCESS_LEVEL_EDIT",
		5: "ACCESS_LEVEL_MANAGE",
	}
	AccessLevel_value = map[string]int32{
		"ACCESS_LEVEL_UNSPECIFIED": 0,
		"ACCESS_LEVEL_BLOCKED":     1,
		"ACCESS_LEVEL_VIEW":        2,
		"ACCESS_LEVEL_SHARE":       3,
		"ACCESS_LEVEL_EDIT":        4,
		"ACCESS_LEVEL_MANAGE":      5,
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
	return file_resources_calendar_access_proto_enumTypes[0].Descriptor()
}

func (AccessLevel) Type() protoreflect.EnumType {
	return &file_resources_calendar_access_proto_enumTypes[0]
}

func (x AccessLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccessLevel.Descriptor instead.
func (AccessLevel) EnumDescriptor() ([]byte, []int) {
	return file_resources_calendar_access_proto_rawDescGZIP(), []int{0}
}

type CalendarAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Jobs          []*CalendarJobAccess   `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty" alias:"job_access"`   // @gotags: alias:"job_access"
	Users         []*CalendarUserAccess  `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty" alias:"user_access"` // @gotags: alias:"user_access"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CalendarAccess) Reset() {
	*x = CalendarAccess{}
	mi := &file_resources_calendar_access_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CalendarAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CalendarAccess) ProtoMessage() {}

func (x *CalendarAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_calendar_access_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CalendarAccess.ProtoReflect.Descriptor instead.
func (*CalendarAccess) Descriptor() ([]byte, []int) {
	return file_resources_calendar_access_proto_rawDescGZIP(), []int{0}
}

func (x *CalendarAccess) GetJobs() []*CalendarJobAccess {
	if x != nil {
		return x.Jobs
	}
	return nil
}

func (x *CalendarAccess) GetUsers() []*CalendarUserAccess {
	if x != nil {
		return x.Users
	}
	return nil
}

type CalendarJobAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	TargetId      uint64                 `protobuf:"varint,3,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`
	Job           string                 `protobuf:"bytes,4,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel      *string                `protobuf:"bytes,5,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty"`
	MinimumGrade  int32                  `protobuf:"varint,6,opt,name=minimum_grade,json=minimumGrade,proto3" json:"minimum_grade,omitempty"`
	JobGradeLabel *string                `protobuf:"bytes,7,opt,name=job_grade_label,json=jobGradeLabel,proto3,oneof" json:"job_grade_label,omitempty"`
	Access        AccessLevel            `protobuf:"varint,8,opt,name=access,proto3,enum=resources.calendar.AccessLevel" json:"access,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CalendarJobAccess) Reset() {
	*x = CalendarJobAccess{}
	mi := &file_resources_calendar_access_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CalendarJobAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CalendarJobAccess) ProtoMessage() {}

func (x *CalendarJobAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_calendar_access_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CalendarJobAccess.ProtoReflect.Descriptor instead.
func (*CalendarJobAccess) Descriptor() ([]byte, []int) {
	return file_resources_calendar_access_proto_rawDescGZIP(), []int{1}
}

func (x *CalendarJobAccess) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CalendarJobAccess) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *CalendarJobAccess) GetTargetId() uint64 {
	if x != nil {
		return x.TargetId
	}
	return 0
}

func (x *CalendarJobAccess) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *CalendarJobAccess) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
	}
	return ""
}

func (x *CalendarJobAccess) GetMinimumGrade() int32 {
	if x != nil {
		return x.MinimumGrade
	}
	return 0
}

func (x *CalendarJobAccess) GetJobGradeLabel() string {
	if x != nil && x.JobGradeLabel != nil {
		return *x.JobGradeLabel
	}
	return ""
}

func (x *CalendarJobAccess) GetAccess() AccessLevel {
	if x != nil {
		return x.Access
	}
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

type CalendarUserAccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	TargetId      uint64                 `protobuf:"varint,3,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`
	UserId        int32                  `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	User          *users.UserShort       `protobuf:"bytes,5,opt,name=user,proto3,oneof" json:"user,omitempty"`
	Access        AccessLevel            `protobuf:"varint,6,opt,name=access,proto3,enum=resources.calendar.AccessLevel" json:"access,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CalendarUserAccess) Reset() {
	*x = CalendarUserAccess{}
	mi := &file_resources_calendar_access_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CalendarUserAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CalendarUserAccess) ProtoMessage() {}

func (x *CalendarUserAccess) ProtoReflect() protoreflect.Message {
	mi := &file_resources_calendar_access_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CalendarUserAccess.ProtoReflect.Descriptor instead.
func (*CalendarUserAccess) Descriptor() ([]byte, []int) {
	return file_resources_calendar_access_proto_rawDescGZIP(), []int{2}
}

func (x *CalendarUserAccess) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CalendarUserAccess) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *CalendarUserAccess) GetTargetId() uint64 {
	if x != nil {
		return x.TargetId
	}
	return 0
}

func (x *CalendarUserAccess) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CalendarUserAccess) GetUser() *users.UserShort {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *CalendarUserAccess) GetAccess() AccessLevel {
	if x != nil {
		return x.Access
	}
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

var File_resources_calendar_access_proto protoreflect.FileDescriptor

const file_resources_calendar_access_proto_rawDesc = "" +
	"\n" +
	"\x1fresources/calendar/access.proto\x12\x12resources.calendar\x1a#resources/timestamp/timestamp.proto\x1a\x1bresources/users/users.proto\x1a\x17validate/validate.proto\"\x9d\x01\n" +
	"\x0eCalendarAccess\x12C\n" +
	"\x04jobs\x18\x01 \x03(\v2%.resources.calendar.CalendarJobAccessB\b\xfaB\x05\x92\x01\x02\x10\x14R\x04jobs\x12F\n" +
	"\x05users\x18\x02 \x03(\v2&.resources.calendar.CalendarUserAccessB\b\xfaB\x05\x92\x01\x02\x10\x14R\x05users\"\xa2\x03\n" +
	"\x11CalendarJobAccess\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12B\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12\x1b\n" +
	"\ttarget_id\x18\x03 \x01(\x04R\btargetId\x12\x19\n" +
	"\x03job\x18\x04 \x01(\tB\a\xfaB\x04r\x02\x18\x14R\x03job\x12)\n" +
	"\tjob_label\x18\x05 \x01(\tB\a\xfaB\x04r\x02\x182H\x01R\bjobLabel\x88\x01\x01\x12,\n" +
	"\rminimum_grade\x18\x06 \x01(\x05B\a\xfaB\x04\x1a\x02(\x00R\fminimumGrade\x124\n" +
	"\x0fjob_grade_label\x18\a \x01(\tB\a\xfaB\x04r\x02\x182H\x02R\rjobGradeLabel\x88\x01\x01\x12A\n" +
	"\x06access\x18\b \x01(\x0e2\x1f.resources.calendar.AccessLevelB\b\xfaB\x05\x82\x01\x02\x10\x01R\x06accessB\r\n" +
	"\v_created_atB\f\n" +
	"\n" +
	"_job_labelB\x12\n" +
	"\x10_job_grade_label\"\xb7\x02\n" +
	"\x12CalendarUserAccess\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12B\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tcreatedAt\x88\x01\x01\x12\x1b\n" +
	"\ttarget_id\x18\x03 \x01(\x04R\btargetId\x12 \n" +
	"\auser_id\x18\x04 \x01(\x05B\a\xfaB\x04\x1a\x02(\x00R\x06userId\x123\n" +
	"\x04user\x18\x05 \x01(\v2\x1a.resources.users.UserShortH\x01R\x04user\x88\x01\x01\x12A\n" +
	"\x06access\x18\x06 \x01(\x0e2\x1f.resources.calendar.AccessLevelB\b\xfaB\x05\x82\x01\x02\x10\x01R\x06accessB\r\n" +
	"\v_created_atB\a\n" +
	"\x05_user*\xa4\x01\n" +
	"\vAccessLevel\x12\x1c\n" +
	"\x18ACCESS_LEVEL_UNSPECIFIED\x10\x00\x12\x18\n" +
	"\x14ACCESS_LEVEL_BLOCKED\x10\x01\x12\x15\n" +
	"\x11ACCESS_LEVEL_VIEW\x10\x02\x12\x16\n" +
	"\x12ACCESS_LEVEL_SHARE\x10\x03\x12\x15\n" +
	"\x11ACCESS_LEVEL_EDIT\x10\x04\x12\x17\n" +
	"\x13ACCESS_LEVEL_MANAGE\x10\x05BIZGgithub.com/fivenet-app/fivenet/gen/go/proto/resources/calendar;calendarb\x06proto3"

var (
	file_resources_calendar_access_proto_rawDescOnce sync.Once
	file_resources_calendar_access_proto_rawDescData []byte
)

func file_resources_calendar_access_proto_rawDescGZIP() []byte {
	file_resources_calendar_access_proto_rawDescOnce.Do(func() {
		file_resources_calendar_access_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_calendar_access_proto_rawDesc), len(file_resources_calendar_access_proto_rawDesc)))
	})
	return file_resources_calendar_access_proto_rawDescData
}

var file_resources_calendar_access_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_calendar_access_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_resources_calendar_access_proto_goTypes = []any{
	(AccessLevel)(0),            // 0: resources.calendar.AccessLevel
	(*CalendarAccess)(nil),      // 1: resources.calendar.CalendarAccess
	(*CalendarJobAccess)(nil),   // 2: resources.calendar.CalendarJobAccess
	(*CalendarUserAccess)(nil),  // 3: resources.calendar.CalendarUserAccess
	(*timestamp.Timestamp)(nil), // 4: resources.timestamp.Timestamp
	(*users.UserShort)(nil),     // 5: resources.users.UserShort
}
var file_resources_calendar_access_proto_depIdxs = []int32{
	2, // 0: resources.calendar.CalendarAccess.jobs:type_name -> resources.calendar.CalendarJobAccess
	3, // 1: resources.calendar.CalendarAccess.users:type_name -> resources.calendar.CalendarUserAccess
	4, // 2: resources.calendar.CalendarJobAccess.created_at:type_name -> resources.timestamp.Timestamp
	0, // 3: resources.calendar.CalendarJobAccess.access:type_name -> resources.calendar.AccessLevel
	4, // 4: resources.calendar.CalendarUserAccess.created_at:type_name -> resources.timestamp.Timestamp
	5, // 5: resources.calendar.CalendarUserAccess.user:type_name -> resources.users.UserShort
	0, // 6: resources.calendar.CalendarUserAccess.access:type_name -> resources.calendar.AccessLevel
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_resources_calendar_access_proto_init() }
func file_resources_calendar_access_proto_init() {
	if File_resources_calendar_access_proto != nil {
		return
	}
	file_resources_calendar_access_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_calendar_access_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_calendar_access_proto_rawDesc), len(file_resources_calendar_access_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_calendar_access_proto_goTypes,
		DependencyIndexes: file_resources_calendar_access_proto_depIdxs,
		EnumInfos:         file_resources_calendar_access_proto_enumTypes,
		MessageInfos:      file_resources_calendar_access_proto_msgTypes,
	}.Build()
	File_resources_calendar_access_proto = out.File
	file_resources_calendar_access_proto_goTypes = nil
	file_resources_calendar_access_proto_depIdxs = nil
}
