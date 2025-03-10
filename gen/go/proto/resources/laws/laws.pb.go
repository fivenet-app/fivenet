// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: resources/laws/laws.proto

package laws

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	timestamp "github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
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

type LawBook struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	// @sanitize
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// @sanitize
	Description   *string `protobuf:"bytes,5,opt,name=description,proto3,oneof" json:"description,omitempty"`
	Laws          []*Law  `protobuf:"bytes,6,rep,name=laws,proto3" json:"laws,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LawBook) Reset() {
	*x = LawBook{}
	mi := &file_resources_laws_laws_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LawBook) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LawBook) ProtoMessage() {}

func (x *LawBook) ProtoReflect() protoreflect.Message {
	mi := &file_resources_laws_laws_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LawBook.ProtoReflect.Descriptor instead.
func (*LawBook) Descriptor() ([]byte, []int) {
	return file_resources_laws_laws_proto_rawDescGZIP(), []int{0}
}

func (x *LawBook) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LawBook) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *LawBook) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *LawBook) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LawBook) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *LawBook) GetLaws() []*Law {
	if x != nil {
		return x.Laws
	}
	return nil
}

type Law struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"law.id"` // @gotags: sql:"primary_key" alias:"law.id"
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	LawbookId uint64                 `protobuf:"varint,4,opt,name=lawbook_id,json=lawbookId,proto3" json:"lawbook_id,omitempty"`
	// @sanitize
	Name string `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	// @sanitize
	Description *string `protobuf:"bytes,6,opt,name=description,proto3,oneof" json:"description,omitempty"`
	// @sanitize
	Hint          *string `protobuf:"bytes,7,opt,name=hint,proto3,oneof" json:"hint,omitempty"`
	Fine          *uint32 `protobuf:"varint,8,opt,name=fine,proto3,oneof" json:"fine,omitempty"`
	DetentionTime *uint32 `protobuf:"varint,9,opt,name=detention_time,json=detentionTime,proto3,oneof" json:"detention_time,omitempty"`
	StvoPoints    *uint32 `protobuf:"varint,10,opt,name=stvo_points,json=stvoPoints,proto3,oneof" json:"stvo_points,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Law) Reset() {
	*x = Law{}
	mi := &file_resources_laws_laws_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Law) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Law) ProtoMessage() {}

func (x *Law) ProtoReflect() protoreflect.Message {
	mi := &file_resources_laws_laws_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Law.ProtoReflect.Descriptor instead.
func (*Law) Descriptor() ([]byte, []int) {
	return file_resources_laws_laws_proto_rawDescGZIP(), []int{1}
}

func (x *Law) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Law) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Law) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Law) GetLawbookId() uint64 {
	if x != nil {
		return x.LawbookId
	}
	return 0
}

func (x *Law) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Law) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Law) GetHint() string {
	if x != nil && x.Hint != nil {
		return *x.Hint
	}
	return ""
}

func (x *Law) GetFine() uint32 {
	if x != nil && x.Fine != nil {
		return *x.Fine
	}
	return 0
}

func (x *Law) GetDetentionTime() uint32 {
	if x != nil && x.DetentionTime != nil {
		return *x.DetentionTime
	}
	return 0
}

func (x *Law) GetStvoPoints() uint32 {
	if x != nil && x.StvoPoints != nil {
		return *x.StvoPoints
	}
	return 0
}

var File_resources_laws_laws_proto protoreflect.FileDescriptor

var file_resources_laws_laws_proto_rawDesc = string([]byte{
	0x0a, 0x19, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6c, 0x61, 0x77, 0x73,
	0x2f, 0x6c, 0x61, 0x77, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x61, 0x77, 0x73, 0x1a, 0x23, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc9, 0x02, 0x0a, 0x07, 0x4c, 0x61,
	0x77, 0x42, 0x6f, 0x6f, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07,
	0x72, 0x05, 0x10, 0x03, 0x18, 0x80, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2f, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x48, 0x02, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x27,
	0x0a, 0x04, 0x6c, 0x61, 0x77, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x61, 0x77, 0x73, 0x2e, 0x4c, 0x61,
	0x77, 0x52, 0x04, 0x6c, 0x61, 0x77, 0x73, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xfe, 0x03, 0x0a, 0x03, 0x4c, 0x61, 0x77, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x48, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01,
	0x01, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x61, 0x77, 0x62, 0x6f, 0x6f, 0x6b,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x6c, 0x61, 0x77, 0x62, 0x6f,
	0x6f, 0x6b, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x03, 0x18, 0x80, 0x01, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2f, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03,
	0x18, 0x80, 0x08, 0x48, 0x02, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x04, 0x68, 0x69, 0x6e, 0x74, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x04, 0x48, 0x03, 0x52,
	0x04, 0x68, 0x69, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x66, 0x69, 0x6e, 0x65,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x04, 0x52, 0x04, 0x66, 0x69, 0x6e, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x2a, 0x0a, 0x0e, 0x64, 0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x05, 0x52, 0x0d, 0x64, 0x65, 0x74,
	0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x24, 0x0a,
	0x0b, 0x73, 0x74, 0x76, 0x6f, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x0d, 0x48, 0x06, 0x52, 0x0a, 0x73, 0x74, 0x76, 0x6f, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x68, 0x69, 0x6e, 0x74, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x66,
	0x69, 0x6e, 0x65, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x64, 0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x73, 0x74, 0x76, 0x6f, 0x5f,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70,
	0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f,
	0x6c, 0x61, 0x77, 0x73, 0x3b, 0x6c, 0x61, 0x77, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_resources_laws_laws_proto_rawDescOnce sync.Once
	file_resources_laws_laws_proto_rawDescData []byte
)

func file_resources_laws_laws_proto_rawDescGZIP() []byte {
	file_resources_laws_laws_proto_rawDescOnce.Do(func() {
		file_resources_laws_laws_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_laws_laws_proto_rawDesc), len(file_resources_laws_laws_proto_rawDesc)))
	})
	return file_resources_laws_laws_proto_rawDescData
}

var file_resources_laws_laws_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_resources_laws_laws_proto_goTypes = []any{
	(*LawBook)(nil),             // 0: resources.laws.LawBook
	(*Law)(nil),                 // 1: resources.laws.Law
	(*timestamp.Timestamp)(nil), // 2: resources.timestamp.Timestamp
}
var file_resources_laws_laws_proto_depIdxs = []int32{
	2, // 0: resources.laws.LawBook.created_at:type_name -> resources.timestamp.Timestamp
	2, // 1: resources.laws.LawBook.updated_at:type_name -> resources.timestamp.Timestamp
	1, // 2: resources.laws.LawBook.laws:type_name -> resources.laws.Law
	2, // 3: resources.laws.Law.created_at:type_name -> resources.timestamp.Timestamp
	2, // 4: resources.laws.Law.updated_at:type_name -> resources.timestamp.Timestamp
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_resources_laws_laws_proto_init() }
func file_resources_laws_laws_proto_init() {
	if File_resources_laws_laws_proto != nil {
		return
	}
	file_resources_laws_laws_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_laws_laws_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_laws_laws_proto_rawDesc), len(file_resources_laws_laws_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_laws_laws_proto_goTypes,
		DependencyIndexes: file_resources_laws_laws_proto_depIdxs,
		MessageInfos:      file_resources_laws_laws_proto_msgTypes,
	}.Build()
	File_resources_laws_laws_proto = out.File
	file_resources_laws_laws_proto_goTypes = nil
	file_resources_laws_laws_proto_depIdxs = nil
}
