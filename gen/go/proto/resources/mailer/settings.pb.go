// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: resources/mailer/settings.proto

package mailer

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

type EmailSettings struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	EmailId uint64                 `protobuf:"varint,1,opt,name=email_id,json=emailId,proto3" json:"email_id,omitempty"`
	// @sanitize
	Signature *string `protobuf:"bytes,2,opt,name=signature,proto3,oneof" json:"signature,omitempty"`
	// @sanitize: method=StripTags
	BlockedEmails []string `protobuf:"bytes,3,rep,name=blocked_emails,json=blockedEmails,proto3" json:"blocked_emails,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmailSettings) Reset() {
	*x = EmailSettings{}
	mi := &file_resources_mailer_settings_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailSettings) ProtoMessage() {}

func (x *EmailSettings) ProtoReflect() protoreflect.Message {
	mi := &file_resources_mailer_settings_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailSettings.ProtoReflect.Descriptor instead.
func (*EmailSettings) Descriptor() ([]byte, []int) {
	return file_resources_mailer_settings_proto_rawDescGZIP(), []int{0}
}

func (x *EmailSettings) GetEmailId() uint64 {
	if x != nil {
		return x.EmailId
	}
	return 0
}

func (x *EmailSettings) GetSignature() string {
	if x != nil && x.Signature != nil {
		return *x.Signature
	}
	return ""
}

func (x *EmailSettings) GetBlockedEmails() []string {
	if x != nil {
		return x.BlockedEmails
	}
	return nil
}

var File_resources_mailer_settings_proto protoreflect.FileDescriptor

const file_resources_mailer_settings_proto_rawDesc = "" +
	"\n" +
	"\x1fresources/mailer/settings.proto\x12\x10resources.mailer\x1a\x17validate/validate.proto\"\x96\x01\n" +
	"\rEmailSettings\x12\x19\n" +
	"\bemail_id\x18\x01 \x01(\x04R\aemailId\x12+\n" +
	"\tsignature\x18\x02 \x01(\tB\b\xfaB\x05r\x03\x18\x80\bH\x00R\tsignature\x88\x01\x01\x12/\n" +
	"\x0eblocked_emails\x18\x03 \x03(\tB\b\xfaB\x05\x92\x01\x02\x10\x19R\rblockedEmailsB\f\n" +
	"\n" +
	"_signatureBEZCgithub.com/fivenet-app/fivenet/gen/go/proto/resources/mailer;mailerb\x06proto3"

var (
	file_resources_mailer_settings_proto_rawDescOnce sync.Once
	file_resources_mailer_settings_proto_rawDescData []byte
)

func file_resources_mailer_settings_proto_rawDescGZIP() []byte {
	file_resources_mailer_settings_proto_rawDescOnce.Do(func() {
		file_resources_mailer_settings_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_mailer_settings_proto_rawDesc), len(file_resources_mailer_settings_proto_rawDesc)))
	})
	return file_resources_mailer_settings_proto_rawDescData
}

var file_resources_mailer_settings_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_mailer_settings_proto_goTypes = []any{
	(*EmailSettings)(nil), // 0: resources.mailer.EmailSettings
}
var file_resources_mailer_settings_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_resources_mailer_settings_proto_init() }
func file_resources_mailer_settings_proto_init() {
	if File_resources_mailer_settings_proto != nil {
		return
	}
	file_resources_mailer_settings_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_mailer_settings_proto_rawDesc), len(file_resources_mailer_settings_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_mailer_settings_proto_goTypes,
		DependencyIndexes: file_resources_mailer_settings_proto_depIdxs,
		MessageInfos:      file_resources_mailer_settings_proto_msgTypes,
	}.Build()
	File_resources_mailer_settings_proto = out.File
	file_resources_mailer_settings_proto_goTypes = nil
	file_resources_mailer_settings_proto_depIdxs = nil
}
