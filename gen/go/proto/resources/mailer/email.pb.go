// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.20.3
// source: resources/mailer/email.proto

package mailer

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
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

type Email struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	DeletedAt *timestamp.Timestamp `protobuf:"bytes,4,opt,name=deleted_at,json=deletedAt,proto3,oneof" json:"deleted_at,omitempty"`
	Disabled  bool                 `protobuf:"varint,5,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Job       *string              `protobuf:"bytes,6,opt,name=job,proto3,oneof" json:"job,omitempty"`
	UserId    *int32               `protobuf:"varint,7,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	User      *users.UserShort     `protobuf:"bytes,8,opt,name=user,proto3,oneof" json:"user,omitempty"`
	// @sanitize: method=StripTags
	Email *string `protobuf:"bytes,9,opt,name=email,proto3,oneof" json:"email,omitempty"`
	// @sanitize: method=StripTags
	Label     string  `protobuf:"bytes,10,opt,name=label,proto3" json:"label,omitempty"`
	Internal  bool    `protobuf:"varint,11,opt,name=internal,proto3" json:"internal,omitempty"`
	Signature *string `protobuf:"bytes,12,opt,name=signature,proto3,oneof" json:"signature,omitempty"`
	Access    *Access `protobuf:"bytes,13,opt,name=access,proto3" json:"access,omitempty"`
}

func (x *Email) Reset() {
	*x = Email{}
	mi := &file_resources_mailer_email_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Email) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Email) ProtoMessage() {}

func (x *Email) ProtoReflect() protoreflect.Message {
	mi := &file_resources_mailer_email_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Email.ProtoReflect.Descriptor instead.
func (*Email) Descriptor() ([]byte, []int) {
	return file_resources_mailer_email_proto_rawDescGZIP(), []int{0}
}

func (x *Email) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Email) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Email) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Email) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *Email) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *Email) GetJob() string {
	if x != nil && x.Job != nil {
		return *x.Job
	}
	return ""
}

func (x *Email) GetUserId() int32 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

func (x *Email) GetUser() *users.UserShort {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *Email) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *Email) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *Email) GetInternal() bool {
	if x != nil {
		return x.Internal
	}
	return false
}

func (x *Email) GetSignature() string {
	if x != nil && x.Signature != nil {
		return *x.Signature
	}
	return ""
}

func (x *Email) GetAccess() *Access {
	if x != nil {
		return x.Access
	}
	return nil
}

type EmailShort struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	DeletedAt *timestamp.Timestamp `protobuf:"bytes,4,opt,name=deleted_at,json=deletedAt,proto3,oneof" json:"deleted_at,omitempty"`
	Disabled  bool                 `protobuf:"varint,5,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Job       *string              `protobuf:"bytes,6,opt,name=job,proto3,oneof" json:"job,omitempty"`
	UserId    *int32               `protobuf:"varint,7,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	User      *users.UserShort     `protobuf:"bytes,8,opt,name=user,proto3,oneof" json:"user,omitempty"`
	// @sanitize: method=StripTags
	Email *string `protobuf:"bytes,9,opt,name=email,proto3,oneof" json:"email,omitempty"`
	// @sanitize: method=StripTags
	Label    string `protobuf:"bytes,10,opt,name=label,proto3" json:"label,omitempty"`
	Internal bool   `protobuf:"varint,11,opt,name=internal,proto3" json:"internal,omitempty"`
}

func (x *EmailShort) Reset() {
	*x = EmailShort{}
	mi := &file_resources_mailer_email_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailShort) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailShort) ProtoMessage() {}

func (x *EmailShort) ProtoReflect() protoreflect.Message {
	mi := &file_resources_mailer_email_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailShort.ProtoReflect.Descriptor instead.
func (*EmailShort) Descriptor() ([]byte, []int) {
	return file_resources_mailer_email_proto_rawDescGZIP(), []int{1}
}

func (x *EmailShort) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *EmailShort) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *EmailShort) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *EmailShort) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *EmailShort) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *EmailShort) GetJob() string {
	if x != nil && x.Job != nil {
		return *x.Job
	}
	return ""
}

func (x *EmailShort) GetUserId() int32 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

func (x *EmailShort) GetUser() *users.UserShort {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *EmailShort) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *EmailShort) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *EmailShort) GetInternal() bool {
	if x != nil {
		return x.Internal
	}
	return false
}

var File_resources_mailer_email_proto protoreflect.FileDescriptor

var file_resources_mailer_email_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6d, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x2f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72,
	0x1a, 0x1d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6d, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x05, 0x0a, 0x05, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x42, 0x02, 0x30, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48,
	0x01, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12,
	0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x03, 0x6a,
	0x6f, 0x62, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18,
	0x28, 0x48, 0x02, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x1a, 0x02, 0x20, 0x00, 0x48, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88,
	0x01, 0x01, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x48, 0x04, 0x52, 0x04,
	0x75, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x24, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x10, 0x06, 0x18,
	0x50, 0x48, 0x05, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a,
	0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42,
	0x07, 0x72, 0x05, 0x10, 0x02, 0x18, 0x80, 0x01, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12,
	0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x12, 0x2b, 0x0a, 0x09, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x08, 0x48, 0x06, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x88, 0x01, 0x01, 0x12, 0x30, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6a, 0x6f, 0x62,
	0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x42, 0x07, 0x0a, 0x05,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x42,
	0x0c, 0x0a, 0x0a, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0xa8, 0x04,
	0x0a, 0x0a, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x30, 0x01, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x3d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62,
	0x6c, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62,
	0x6c, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x28, 0x48, 0x02, 0x52, 0x03, 0x6a, 0x6f, 0x62,
	0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x48, 0x03, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x48, 0x04, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12,
	0x24, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09,
	0xfa, 0x42, 0x06, 0x72, 0x04, 0x10, 0x02, 0x18, 0x28, 0x48, 0x05, 0x52, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x02, 0x18, 0x80, 0x01,
	0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x6a, 0x6f, 0x62, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x08,
	0x0a, 0x06, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61,
	0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x3b, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_mailer_email_proto_rawDescOnce sync.Once
	file_resources_mailer_email_proto_rawDescData = file_resources_mailer_email_proto_rawDesc
)

func file_resources_mailer_email_proto_rawDescGZIP() []byte {
	file_resources_mailer_email_proto_rawDescOnce.Do(func() {
		file_resources_mailer_email_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_mailer_email_proto_rawDescData)
	})
	return file_resources_mailer_email_proto_rawDescData
}

var file_resources_mailer_email_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_resources_mailer_email_proto_goTypes = []any{
	(*Email)(nil),               // 0: resources.mailer.Email
	(*EmailShort)(nil),          // 1: resources.mailer.EmailShort
	(*timestamp.Timestamp)(nil), // 2: resources.timestamp.Timestamp
	(*users.UserShort)(nil),     // 3: resources.users.UserShort
	(*Access)(nil),              // 4: resources.mailer.Access
}
var file_resources_mailer_email_proto_depIdxs = []int32{
	2, // 0: resources.mailer.Email.created_at:type_name -> resources.timestamp.Timestamp
	2, // 1: resources.mailer.Email.updated_at:type_name -> resources.timestamp.Timestamp
	2, // 2: resources.mailer.Email.deleted_at:type_name -> resources.timestamp.Timestamp
	3, // 3: resources.mailer.Email.user:type_name -> resources.users.UserShort
	4, // 4: resources.mailer.Email.access:type_name -> resources.mailer.Access
	2, // 5: resources.mailer.EmailShort.created_at:type_name -> resources.timestamp.Timestamp
	2, // 6: resources.mailer.EmailShort.updated_at:type_name -> resources.timestamp.Timestamp
	2, // 7: resources.mailer.EmailShort.deleted_at:type_name -> resources.timestamp.Timestamp
	3, // 8: resources.mailer.EmailShort.user:type_name -> resources.users.UserShort
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_resources_mailer_email_proto_init() }
func file_resources_mailer_email_proto_init() {
	if File_resources_mailer_email_proto != nil {
		return
	}
	file_resources_mailer_access_proto_init()
	file_resources_mailer_email_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_mailer_email_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_mailer_email_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_mailer_email_proto_goTypes,
		DependencyIndexes: file_resources_mailer_email_proto_depIdxs,
		MessageInfos:      file_resources_mailer_email_proto_msgTypes,
	}.Build()
	File_resources_mailer_email_proto = out.File
	file_resources_mailer_email_proto_rawDesc = nil
	file_resources_mailer_email_proto_goTypes = nil
	file_resources_mailer_email_proto_depIdxs = nil
}
