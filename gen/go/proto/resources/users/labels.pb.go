// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: resources/users/labels.proto

package users

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

// @dbscanner: json
type CitizenLabels struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	List          []*CitizenLabel        `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CitizenLabels) Reset() {
	*x = CitizenLabels{}
	mi := &file_resources_users_labels_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CitizenLabels) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CitizenLabels) ProtoMessage() {}

func (x *CitizenLabels) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_labels_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CitizenLabels.ProtoReflect.Descriptor instead.
func (*CitizenLabels) Descriptor() ([]byte, []int) {
	return file_resources_users_labels_proto_rawDescGZIP(), []int{0}
}

func (x *CitizenLabels) GetList() []*CitizenLabel {
	if x != nil {
		return x.List
	}
	return nil
}

type CitizenLabel struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	Job   *string                `protobuf:"bytes,2,opt,name=job,proto3,oneof" json:"job,omitempty"`
	Name  string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// @sanitize: method=StripTags
	Color         string `protobuf:"bytes,4,opt,name=color,proto3" json:"color,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CitizenLabel) Reset() {
	*x = CitizenLabel{}
	mi := &file_resources_users_labels_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CitizenLabel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CitizenLabel) ProtoMessage() {}

func (x *CitizenLabel) ProtoReflect() protoreflect.Message {
	mi := &file_resources_users_labels_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CitizenLabel.ProtoReflect.Descriptor instead.
func (*CitizenLabel) Descriptor() ([]byte, []int) {
	return file_resources_users_labels_proto_rawDescGZIP(), []int{1}
}

func (x *CitizenLabel) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CitizenLabel) GetJob() string {
	if x != nil && x.Job != nil {
		return *x.Job
	}
	return ""
}

func (x *CitizenLabel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CitizenLabel) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

var File_resources_users_labels_proto protoreflect.FileDescriptor

const file_resources_users_labels_proto_rawDesc = "" +
	"\n" +
	"\x1cresources/users/labels.proto\x12\x0fresources.users\x1a\x17validate/validate.proto\"L\n" +
	"\rCitizenLabels\x12;\n" +
	"\x04list\x18\x01 \x03(\v2\x1d.resources.users.CitizenLabelB\b\xfaB\x05\x92\x01\x02\x10\n" +
	"R\x04list\"\x96\x01\n" +
	"\fCitizenLabel\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x1e\n" +
	"\x03job\x18\x02 \x01(\tB\a\xfaB\x04r\x02\x18\x14H\x00R\x03job\x88\x01\x01\x12\x1b\n" +
	"\x04name\x18\x03 \x01(\tB\a\xfaB\x04r\x02\x180R\x04name\x121\n" +
	"\x05color\x18\x04 \x01(\tB\x1b\xfaB\x18r\x162\x11^#[A-Fa-f0-9]{6}$\x98\x01\aR\x05colorB\x06\n" +
	"\x04_jobBCZAgithub.com/fivenet-app/fivenet/gen/go/proto/resources/users;usersb\x06proto3"

var (
	file_resources_users_labels_proto_rawDescOnce sync.Once
	file_resources_users_labels_proto_rawDescData []byte
)

func file_resources_users_labels_proto_rawDescGZIP() []byte {
	file_resources_users_labels_proto_rawDescOnce.Do(func() {
		file_resources_users_labels_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_users_labels_proto_rawDesc), len(file_resources_users_labels_proto_rawDesc)))
	})
	return file_resources_users_labels_proto_rawDescData
}

var file_resources_users_labels_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_resources_users_labels_proto_goTypes = []any{
	(*CitizenLabels)(nil), // 0: resources.users.CitizenLabels
	(*CitizenLabel)(nil),  // 1: resources.users.CitizenLabel
}
var file_resources_users_labels_proto_depIdxs = []int32{
	1, // 0: resources.users.CitizenLabels.list:type_name -> resources.users.CitizenLabel
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_resources_users_labels_proto_init() }
func file_resources_users_labels_proto_init() {
	if File_resources_users_labels_proto != nil {
		return
	}
	file_resources_users_labels_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_users_labels_proto_rawDesc), len(file_resources_users_labels_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_users_labels_proto_goTypes,
		DependencyIndexes: file_resources_users_labels_proto_depIdxs,
		MessageInfos:      file_resources_users_labels_proto_msgTypes,
	}.Build()
	File_resources_users_labels_proto = out.File
	file_resources_users_labels_proto_goTypes = nil
	file_resources_users_labels_proto_depIdxs = nil
}
