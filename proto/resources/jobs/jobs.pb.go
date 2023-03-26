// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: resources/jobs/jobs.proto

package jobs

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" sql:"primary_key" alias:"name"`   // @gotags: sql:"primary_key" alias:"name"
	Label  string      `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty" alias:"label"` // @gotags: alias:"label"
	Grades []*JobGrade `protobuf:"bytes,3,rep,name=grades,proto3" json:"grades,omitempty"`
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_jobs_jobs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_resources_jobs_jobs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_resources_jobs_jobs_proto_rawDescGZIP(), []int{0}
}

func (x *Job) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Job) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *Job) GetGrades() []*JobGrade {
	if x != nil {
		return x.Grades
	}
	return nil
}

type JobGrade struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName string `protobuf:"bytes,1,opt,name=job_name,json=jobName,proto3" json:"job_name,omitempty" alias:"job_name"` // @gotags: alias:"job_name"
	Grade   int32  `protobuf:"varint,2,opt,name=grade,proto3" json:"grade,omitempty" alias:"grade"`                   // @gotags: alias:"grade"
	Label   string `protobuf:"bytes,3,opt,name=label,proto3" json:"label,omitempty" alias:"label"`                    // @gotags: alias:"label"
}

func (x *JobGrade) Reset() {
	*x = JobGrade{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_jobs_jobs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobGrade) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobGrade) ProtoMessage() {}

func (x *JobGrade) ProtoReflect() protoreflect.Message {
	mi := &file_resources_jobs_jobs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobGrade.ProtoReflect.Descriptor instead.
func (*JobGrade) Descriptor() ([]byte, []int) {
	return file_resources_jobs_jobs_proto_rawDescGZIP(), []int{1}
}

func (x *JobGrade) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *JobGrade) GetGrade() int32 {
	if x != nil {
		return x.Grade
	}
	return 0
}

func (x *JobGrade) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

type JobProps struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job   string `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty" alias:"job"`     // @gotags: alias:"job"
	Theme string `protobuf:"bytes,2,opt,name=theme,proto3" json:"theme,omitempty" alias:"theme"` // @gotags: alias:"theme"
}

func (x *JobProps) Reset() {
	*x = JobProps{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_jobs_jobs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobProps) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobProps) ProtoMessage() {}

func (x *JobProps) ProtoReflect() protoreflect.Message {
	mi := &file_resources_jobs_jobs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobProps.ProtoReflect.Descriptor instead.
func (*JobProps) Descriptor() ([]byte, []int) {
	return file_resources_jobs_jobs_proto_rawDescGZIP(), []int{2}
}

func (x *JobProps) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *JobProps) GetTheme() string {
	if x != nil {
		return x.Theme
	}
	return ""
}

var File_resources_jobs_jobs_proto protoreflect.FileDescriptor

var file_resources_jobs_jobs_proto_rawDesc = []byte{
	0x0a, 0x19, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73,
	0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x1a, 0x17, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a, 0x03, 0x4a, 0x6f, 0x62, 0x12, 0x1b, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02,
	0x18, 0x32, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32,
	0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x30, 0x0a, 0x06, 0x67, 0x72, 0x61, 0x64, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x47, 0x72, 0x61, 0x64,
	0x65, 0x52, 0x06, 0x67, 0x72, 0x61, 0x64, 0x65, 0x73, 0x22, 0x63, 0x0a, 0x08, 0x4a, 0x6f, 0x62,
	0x47, 0x72, 0x61, 0x64, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x1d, 0x0a, 0x05, 0x67, 0x72, 0x61, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x52, 0x05, 0x67, 0x72, 0x61, 0x64, 0x65, 0x12,
	0x1d, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x22, 0x44,
	0x0a, 0x08, 0x4a, 0x6f, 0x62, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f,
	0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14,
	0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x1d, 0x0a, 0x05, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52, 0x05, 0x74,
	0x68, 0x65, 0x6d, 0x65, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78, 0x72, 0x74, 0x2f, 0x61, 0x72, 0x70, 0x61, 0x6e,
	0x65, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x3b, 0x6a, 0x6f, 0x62, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_jobs_jobs_proto_rawDescOnce sync.Once
	file_resources_jobs_jobs_proto_rawDescData = file_resources_jobs_jobs_proto_rawDesc
)

func file_resources_jobs_jobs_proto_rawDescGZIP() []byte {
	file_resources_jobs_jobs_proto_rawDescOnce.Do(func() {
		file_resources_jobs_jobs_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_jobs_jobs_proto_rawDescData)
	})
	return file_resources_jobs_jobs_proto_rawDescData
}

var file_resources_jobs_jobs_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_resources_jobs_jobs_proto_goTypes = []interface{}{
	(*Job)(nil),      // 0: resources.jobs.Job
	(*JobGrade)(nil), // 1: resources.jobs.JobGrade
	(*JobProps)(nil), // 2: resources.jobs.JobProps
}
var file_resources_jobs_jobs_proto_depIdxs = []int32{
	1, // 0: resources.jobs.Job.grades:type_name -> resources.jobs.JobGrade
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_resources_jobs_jobs_proto_init() }
func file_resources_jobs_jobs_proto_init() {
	if File_resources_jobs_jobs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resources_jobs_jobs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job); i {
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
		file_resources_jobs_jobs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobGrade); i {
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
		file_resources_jobs_jobs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobProps); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_jobs_jobs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_jobs_jobs_proto_goTypes,
		DependencyIndexes: file_resources_jobs_jobs_proto_depIdxs,
		MessageInfos:      file_resources_jobs_jobs_proto_msgTypes,
	}.Build()
	File_resources_jobs_jobs_proto = out.File
	file_resources_jobs_jobs_proto_rawDesc = nil
	file_resources_jobs_jobs_proto_goTypes = nil
	file_resources_jobs_jobs_proto_depIdxs = nil
}
