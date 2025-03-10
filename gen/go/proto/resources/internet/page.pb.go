// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: resources/internet/page.proto

package internet

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	content "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/content"
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

type PageLayoutType int32

const (
	PageLayoutType_PAGE_LAYOUT_TYPE_UNSPECIFIED  PageLayoutType = 0
	PageLayoutType_PAGE_LAYOUT_TYPE_BASIC_PAGE   PageLayoutType = 1
	PageLayoutType_PAGE_LAYOUT_TYPE_LANDING_PAGE PageLayoutType = 2
)

// Enum value maps for PageLayoutType.
var (
	PageLayoutType_name = map[int32]string{
		0: "PAGE_LAYOUT_TYPE_UNSPECIFIED",
		1: "PAGE_LAYOUT_TYPE_BASIC_PAGE",
		2: "PAGE_LAYOUT_TYPE_LANDING_PAGE",
	}
	PageLayoutType_value = map[string]int32{
		"PAGE_LAYOUT_TYPE_UNSPECIFIED":  0,
		"PAGE_LAYOUT_TYPE_BASIC_PAGE":   1,
		"PAGE_LAYOUT_TYPE_LANDING_PAGE": 2,
	}
)

func (x PageLayoutType) Enum() *PageLayoutType {
	p := new(PageLayoutType)
	*p = x
	return p
}

func (x PageLayoutType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PageLayoutType) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_internet_page_proto_enumTypes[0].Descriptor()
}

func (PageLayoutType) Type() protoreflect.EnumType {
	return &file_resources_internet_page_proto_enumTypes[0]
}

func (x PageLayoutType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PageLayoutType.Descriptor instead.
func (PageLayoutType) EnumDescriptor() ([]byte, []int) {
	return file_resources_internet_page_proto_rawDescGZIP(), []int{0}
}

type Page struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	DeletedAt *timestamp.Timestamp   `protobuf:"bytes,4,opt,name=deleted_at,json=deletedAt,proto3,oneof" json:"deleted_at,omitempty"`
	DomainId  uint64                 `protobuf:"varint,5,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	// @sanitize: method=StripTags
	Path string `protobuf:"bytes,6,opt,name=path,proto3" json:"path,omitempty"`
	// @sanitize: method=StripTags
	Title string `protobuf:"bytes,7,opt,name=title,proto3" json:"title,omitempty"`
	// @sanitize: method=StripTags
	Description   string    `protobuf:"bytes,8,opt,name=description,proto3" json:"description,omitempty"`
	Data          *PageData `protobuf:"bytes,9,opt,name=data,proto3" json:"data,omitempty"`
	CreatorJob    *string   `protobuf:"bytes,10,opt,name=creator_job,json=creatorJob,proto3,oneof" json:"creator_job,omitempty"`
	CreatorId     *int32    `protobuf:"varint,11,opt,name=creator_id,json=creatorId,proto3,oneof" json:"creator_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Page) Reset() {
	*x = Page{}
	mi := &file_resources_internet_page_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Page) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Page) ProtoMessage() {}

func (x *Page) ProtoReflect() protoreflect.Message {
	mi := &file_resources_internet_page_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Page.ProtoReflect.Descriptor instead.
func (*Page) Descriptor() ([]byte, []int) {
	return file_resources_internet_page_proto_rawDescGZIP(), []int{0}
}

func (x *Page) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Page) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Page) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Page) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *Page) GetDomainId() uint64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *Page) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Page) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Page) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Page) GetData() *PageData {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Page) GetCreatorJob() string {
	if x != nil && x.CreatorJob != nil {
		return *x.CreatorJob
	}
	return ""
}

func (x *Page) GetCreatorId() int32 {
	if x != nil && x.CreatorId != nil {
		return *x.CreatorId
	}
	return 0
}

// @dbscanner: json
type PageData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LayoutType    PageLayoutType         `protobuf:"varint,1,opt,name=layout_type,json=layoutType,proto3,enum=resources.internet.PageLayoutType" json:"layout_type,omitempty"`
	Node          *ContentNode           `protobuf:"bytes,2,opt,name=node,proto3,oneof" json:"node,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PageData) Reset() {
	*x = PageData{}
	mi := &file_resources_internet_page_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PageData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageData) ProtoMessage() {}

func (x *PageData) ProtoReflect() protoreflect.Message {
	mi := &file_resources_internet_page_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageData.ProtoReflect.Descriptor instead.
func (*PageData) Descriptor() ([]byte, []int) {
	return file_resources_internet_page_proto_rawDescGZIP(), []int{1}
}

func (x *PageData) GetLayoutType() PageLayoutType {
	if x != nil {
		return x.LayoutType
	}
	return PageLayoutType_PAGE_LAYOUT_TYPE_UNSPECIFIED
}

func (x *PageData) GetNode() *ContentNode {
	if x != nil {
		return x.Node
	}
	return nil
}

type ContentNode struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Type  content.NodeType       `protobuf:"varint,1,opt,name=type,proto3,enum=resources.common.content.NodeType" json:"type,omitempty"`
	// @sanitize: method=StripTags
	Id *string `protobuf:"bytes,2,opt,name=id,proto3,oneof" json:"id,omitempty"`
	// @sanitize: method=StripTags
	Tag string `protobuf:"bytes,3,opt,name=tag,proto3" json:"tag,omitempty"`
	// @sanitize: method=StripTags
	Attrs map[string]string `protobuf:"bytes,4,rep,name=attrs,proto3" json:"attrs,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// @sanitize: method=StripTags
	Text          *string        `protobuf:"bytes,5,opt,name=text,proto3,oneof" json:"text,omitempty"`
	Content       []*ContentNode `protobuf:"bytes,6,rep,name=content,proto3" json:"content,omitempty"`
	Slots         []*ContentNode `protobuf:"bytes,7,rep,name=slots,proto3" json:"slots,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ContentNode) Reset() {
	*x = ContentNode{}
	mi := &file_resources_internet_page_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContentNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContentNode) ProtoMessage() {}

func (x *ContentNode) ProtoReflect() protoreflect.Message {
	mi := &file_resources_internet_page_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContentNode.ProtoReflect.Descriptor instead.
func (*ContentNode) Descriptor() ([]byte, []int) {
	return file_resources_internet_page_proto_rawDescGZIP(), []int{2}
}

func (x *ContentNode) GetType() content.NodeType {
	if x != nil {
		return x.Type
	}
	return content.NodeType(0)
}

func (x *ContentNode) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *ContentNode) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *ContentNode) GetAttrs() map[string]string {
	if x != nil {
		return x.Attrs
	}
	return nil
}

func (x *ContentNode) GetText() string {
	if x != nil && x.Text != nil {
		return *x.Text
	}
	return ""
}

func (x *ContentNode) GetContent() []*ContentNode {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *ContentNode) GetSlots() []*ContentNode {
	if x != nil {
		return x.Slots
	}
	return nil
}

var File_resources_internet_page_proto protoreflect.FileDescriptor

var file_resources_internet_page_proto_rawDesc = string([]byte{
	0x0a, 0x1d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x65, 0x74, 0x2f, 0x70, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x12, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x65, 0x74, 0x1a, 0x26, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xab, 0x04, 0x0a, 0x04, 0x50, 0x61,
	0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x3d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52, 0x09, 0x64, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x01, 0x52, 0x04,
	0x70, 0x61, 0x74, 0x68, 0x12, 0x20, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x01, 0x18, 0xff, 0x01, 0x52,
	0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x2c, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07,
	0x72, 0x05, 0x10, 0x03, 0x18, 0x80, 0x04, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x50, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x24, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x6a, 0x6f, 0x62, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72,
	0x4a, 0x6f, 0x62, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x48, 0x03, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x5f, 0x6a, 0x6f, 0x62, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x22, 0x92, 0x01, 0x0a, 0x08, 0x50, 0x61, 0x67, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x43, 0x0a, 0x0b, 0x6c, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x50,
	0x61, 0x67, 0x65, 0x4c, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x6c,
	0x61, 0x79, 0x6f, 0x75, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x04, 0x6e, 0x6f, 0x64,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65,
	0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x8d, 0x03, 0x0a,
	0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x40, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x13,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x02, 0x69, 0x64,
	0x88, 0x01, 0x01, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x40, 0x0a, 0x05, 0x61, 0x74, 0x74, 0x72, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x4e, 0x6f, 0x64, 0x65, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x05, 0x61, 0x74, 0x74, 0x72, 0x73, 0x12, 0x17, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x88, 0x01, 0x01,
	0x12, 0x39, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4e, 0x6f,
	0x64, 0x65, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x35, 0x0a, 0x05, 0x73,
	0x6c, 0x6f, 0x74, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x73, 0x6c, 0x6f,
	0x74, 0x73, 0x1a, 0x38, 0x0a, 0x0a, 0x41, 0x74, 0x74, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x05, 0x0a, 0x03,
	0x5f, 0x69, 0x64, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x2a, 0x76, 0x0a, 0x0e,
	0x50, 0x61, 0x67, 0x65, 0x4c, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x20,
	0x0a, 0x1c, 0x50, 0x41, 0x47, 0x45, 0x5f, 0x4c, 0x41, 0x59, 0x4f, 0x55, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x1f, 0x0a, 0x1b, 0x50, 0x41, 0x47, 0x45, 0x5f, 0x4c, 0x41, 0x59, 0x4f, 0x55, 0x54, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x41, 0x53, 0x49, 0x43, 0x5f, 0x50, 0x41, 0x47, 0x45, 0x10,
	0x01, 0x12, 0x21, 0x0a, 0x1d, 0x50, 0x41, 0x47, 0x45, 0x5f, 0x4c, 0x41, 0x59, 0x4f, 0x55, 0x54,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x41, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x5f, 0x50, 0x41,
	0x47, 0x45, 0x10, 0x02, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66,
	0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x3b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_internet_page_proto_rawDescOnce sync.Once
	file_resources_internet_page_proto_rawDescData []byte
)

func file_resources_internet_page_proto_rawDescGZIP() []byte {
	file_resources_internet_page_proto_rawDescOnce.Do(func() {
		file_resources_internet_page_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_internet_page_proto_rawDesc), len(file_resources_internet_page_proto_rawDesc)))
	})
	return file_resources_internet_page_proto_rawDescData
}

var file_resources_internet_page_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_internet_page_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_resources_internet_page_proto_goTypes = []any{
	(PageLayoutType)(0),         // 0: resources.internet.PageLayoutType
	(*Page)(nil),                // 1: resources.internet.Page
	(*PageData)(nil),            // 2: resources.internet.PageData
	(*ContentNode)(nil),         // 3: resources.internet.ContentNode
	nil,                         // 4: resources.internet.ContentNode.AttrsEntry
	(*timestamp.Timestamp)(nil), // 5: resources.timestamp.Timestamp
	(content.NodeType)(0),       // 6: resources.common.content.NodeType
}
var file_resources_internet_page_proto_depIdxs = []int32{
	5,  // 0: resources.internet.Page.created_at:type_name -> resources.timestamp.Timestamp
	5,  // 1: resources.internet.Page.updated_at:type_name -> resources.timestamp.Timestamp
	5,  // 2: resources.internet.Page.deleted_at:type_name -> resources.timestamp.Timestamp
	2,  // 3: resources.internet.Page.data:type_name -> resources.internet.PageData
	0,  // 4: resources.internet.PageData.layout_type:type_name -> resources.internet.PageLayoutType
	3,  // 5: resources.internet.PageData.node:type_name -> resources.internet.ContentNode
	6,  // 6: resources.internet.ContentNode.type:type_name -> resources.common.content.NodeType
	4,  // 7: resources.internet.ContentNode.attrs:type_name -> resources.internet.ContentNode.AttrsEntry
	3,  // 8: resources.internet.ContentNode.content:type_name -> resources.internet.ContentNode
	3,  // 9: resources.internet.ContentNode.slots:type_name -> resources.internet.ContentNode
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_resources_internet_page_proto_init() }
func file_resources_internet_page_proto_init() {
	if File_resources_internet_page_proto != nil {
		return
	}
	file_resources_internet_page_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_internet_page_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_internet_page_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_internet_page_proto_rawDesc), len(file_resources_internet_page_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_internet_page_proto_goTypes,
		DependencyIndexes: file_resources_internet_page_proto_depIdxs,
		EnumInfos:         file_resources_internet_page_proto_enumTypes,
		MessageInfos:      file_resources_internet_page_proto_msgTypes,
	}.Build()
	File_resources_internet_page_proto = out.File
	file_resources_internet_page_proto_goTypes = nil
	file_resources_internet_page_proto_depIdxs = nil
}
