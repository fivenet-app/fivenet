// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v3.20.3
// source: resources/users/props.proto

package users

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	filestore "github.com/fivenet-app/fivenet/gen/go/proto/resources/filestore"
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

type UserProps struct {
	state                   protoimpl.MessageState `protogen:"open.v1"`
	UserId                  int32                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UpdatedAt               *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	Wanted                  *bool                  `protobuf:"varint,3,opt,name=wanted,proto3,oneof" json:"wanted,omitempty"`
	JobName                 *string                `protobuf:"bytes,4,opt,name=job_name,json=jobName,proto3,oneof" json:"job_name,omitempty" alias:"job"` // @gotags: alias:"job"
	Job                     *Job                   `protobuf:"bytes,5,opt,name=job,proto3,oneof" json:"job,omitempty"`
	JobGradeNumber          *int32                 `protobuf:"varint,6,opt,name=job_grade_number,json=jobGradeNumber,proto3,oneof" json:"job_grade_number,omitempty" alias:"job_grade"` // @gotags: alias:"job_grade"
	JobGrade                *JobGrade              `protobuf:"bytes,7,opt,name=job_grade,json=jobGrade,proto3,oneof" json:"job_grade,omitempty"`
	TrafficInfractionPoints *uint32                `protobuf:"varint,8,opt,name=traffic_infraction_points,json=trafficInfractionPoints,proto3,oneof" json:"traffic_infraction_points,omitempty"`
	OpenFines               *int64                 `protobuf:"varint,9,opt,name=open_fines,json=openFines,proto3,oneof" json:"open_fines,omitempty"`
	BloodType               *string                `protobuf:"bytes,10,opt,name=blood_type,json=bloodType,proto3,oneof" json:"blood_type,omitempty"`
	MugShot                 *filestore.File        `protobuf:"bytes,11,opt,name=mug_shot,json=mugShot,proto3,oneof" json:"mug_shot,omitempty"`
	Attributes              *CitizenAttributes     `protobuf:"bytes,12,opt,name=attributes,proto3,oneof" json:"attributes,omitempty"`
	// @sanitize: method=StripTags
	Email         *string `protobuf:"bytes,19,opt,name=email,proto3,oneof" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserProps) Reset() {
	*x = UserProps{}
	mi := &file_resources_users_props_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserProps) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserProps) ProtoMessage() {}

func (x *UserProps) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_props_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserProps.ProtoReflect.Descriptor instead.
func (*UserProps) Descriptor() ([]byte, []int) {
	return file_resources_users_props_proto_rawDescGZIP(), []int{0}
}

func (x *UserProps) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserProps) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *UserProps) GetWanted() bool {
	if x != nil && x.Wanted != nil {
		return *x.Wanted
	}
	return false
}

func (x *UserProps) GetJobName() string {
	if x != nil && x.JobName != nil {
		return *x.JobName
	}
	return ""
}

func (x *UserProps) GetJob() *Job {
	if x != nil {
		return x.Job
	}
	return nil
}

func (x *UserProps) GetJobGradeNumber() int32 {
	if x != nil && x.JobGradeNumber != nil {
		return *x.JobGradeNumber
	}
	return 0
}

func (x *UserProps) GetJobGrade() *JobGrade {
	if x != nil {
		return x.JobGrade
	}
	return nil
}

func (x *UserProps) GetTrafficInfractionPoints() uint32 {
	if x != nil && x.TrafficInfractionPoints != nil {
		return *x.TrafficInfractionPoints
	}
	return 0
}

func (x *UserProps) GetOpenFines() int64 {
	if x != nil && x.OpenFines != nil {
		return *x.OpenFines
	}
	return 0
}

func (x *UserProps) GetBloodType() string {
	if x != nil && x.BloodType != nil {
		return *x.BloodType
	}
	return ""
}

func (x *UserProps) GetMugShot() *filestore.File {
	if x != nil {
		return x.MugShot
	}
	return nil
}

func (x *UserProps) GetAttributes() *CitizenAttributes {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *UserProps) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

var File_resources_users_props_proto protoreflect.FileDescriptor

var file_resources_users_props_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x2f, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x1e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xae, 0x06, 0x0a, 0x09, 0x55,
	0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x20, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02,
	0x20, 0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00,
	0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1b,
	0x0a, 0x06, 0x77, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01,
	0x52, 0x06, 0x77, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x08, 0x6a,
	0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52,
	0x07, 0x6a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x2b, 0x0a, 0x03, 0x6a,
	0x6f, 0x62, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x48, 0x03,
	0x52, 0x03, 0x6a, 0x6f, 0x62, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x10, 0x6a, 0x6f, 0x62, 0x5f,
	0x67, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x05, 0x48, 0x04, 0x52, 0x0e, 0x6a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64, 0x65, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x3b, 0x0a, 0x09, 0x6a, 0x6f, 0x62, 0x5f, 0x67,
	0x72, 0x61, 0x64, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x4a, 0x6f, 0x62,
	0x47, 0x72, 0x61, 0x64, 0x65, 0x48, 0x05, 0x52, 0x08, 0x6a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x3f, 0x0a, 0x19, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x5f,
	0x69, 0x6e, 0x66, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x06, 0x52, 0x17, 0x74, 0x72, 0x61, 0x66, 0x66,
	0x69, 0x63, 0x49, 0x6e, 0x66, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x69, 0x6e,
	0x74, 0x73, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x6f, 0x70, 0x65, 0x6e, 0x5f, 0x66, 0x69,
	0x6e, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x48, 0x07, 0x52, 0x09, 0x6f, 0x70, 0x65,
	0x6e, 0x46, 0x69, 0x6e, 0x65, 0x73, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x62, 0x6c, 0x6f,
	0x6f, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x48, 0x08, 0x52,
	0x09, 0x62, 0x6c, 0x6f, 0x6f, 0x64, 0x54, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x12, 0x39, 0x0a,
	0x08, 0x6d, 0x75, 0x67, 0x5f, 0x73, 0x68, 0x6f, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x09, 0x52, 0x07, 0x6d, 0x75,
	0x67, 0x53, 0x68, 0x6f, 0x74, 0x88, 0x01, 0x01, 0x12, 0x47, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x43,
	0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x48, 0x0a, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x88, 0x01,
	0x01, 0x12, 0x24, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x10, 0x06, 0x18, 0x50, 0x48, 0x0b, 0x52, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x77, 0x61, 0x6e, 0x74, 0x65,
	0x64, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x06,
	0x0a, 0x04, 0x5f, 0x6a, 0x6f, 0x62, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x67,
	0x72, 0x61, 0x64, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x6a, 0x6f, 0x62, 0x5f, 0x67, 0x72, 0x61, 0x64, 0x65, 0x42, 0x1c, 0x0a, 0x1a, 0x5f, 0x74, 0x72,
	0x61, 0x66, 0x66, 0x69, 0x63, 0x5f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x6f, 0x70, 0x65, 0x6e,
	0x5f, 0x66, 0x69, 0x6e, 0x65, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x62, 0x6c, 0x6f, 0x6f, 0x64,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6d, 0x75, 0x67, 0x5f, 0x73, 0x68,
	0x6f, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x43, 0x5a, 0x41, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65,
	0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x3b, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_users_props_proto_rawDescOnce sync.Once
	file_resources_users_props_proto_rawDescData = file_resources_users_props_proto_rawDesc
)

func file_resources_users_props_proto_rawDescGZIP() []byte {
	file_resources_users_props_proto_rawDescOnce.Do(func() {
		file_resources_users_props_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_users_props_proto_rawDescData)
	})
	return file_resources_users_props_proto_rawDescData
}

var file_resources_users_props_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_users_props_proto_goTypes = []any{
	(*UserProps)(nil),           // 0: resources.users.UserProps
	(*timestamp.Timestamp)(nil), // 1: resources.timestamp.Timestamp
	(*Job)(nil),                 // 2: resources.users.Job
	(*JobGrade)(nil),            // 3: resources.users.JobGrade
	(*filestore.File)(nil),      // 4: resources.filestore.File
	(*CitizenAttributes)(nil),   // 5: resources.users.CitizenAttributes
}
var file_resources_users_props_proto_depIdxs = []int32{
	1, // 0: resources.users.UserProps.updated_at:type_name -> resources.timestamp.Timestamp
	2, // 1: resources.users.UserProps.job:type_name -> resources.users.Job
	3, // 2: resources.users.UserProps.job_grade:type_name -> resources.users.JobGrade
	4, // 3: resources.users.UserProps.mug_shot:type_name -> resources.filestore.File
	5, // 4: resources.users.UserProps.attributes:type_name -> resources.users.CitizenAttributes
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_resources_users_props_proto_init() }
func file_resources_users_props_proto_init() {
	if File_resources_users_props_proto != nil {
		return
	}
	file_resources_users_attributes_proto_init()
	file_resources_users_jobs_proto_init()
	file_resources_users_props_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_users_props_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_users_props_proto_goTypes,
		DependencyIndexes: file_resources_users_props_proto_depIdxs,
		MessageInfos:      file_resources_users_props_proto_msgTypes,
	}.Build()
	File_resources_users_props_proto = out.File
	file_resources_users_props_proto_rawDesc = nil
	file_resources_users_props_proto_goTypes = nil
	file_resources_users_props_proto_depIdxs = nil
}
