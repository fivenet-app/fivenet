// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: resources/wiki/page.proto

package wiki

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	content "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/content"
	filestore "github.com/fivenet-app/fivenet/gen/go/proto/resources/filestore"
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

type Page struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	// @sanitize: method=StripTags
	Job           string           `protobuf:"bytes,2,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel      *string          `protobuf:"bytes,3,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty"`
	ParentId      *uint64          `protobuf:"varint,4,opt,name=parent_id,json=parentId,proto3,oneof" json:"parent_id,omitempty"`
	Meta          *PageMeta        `protobuf:"bytes,5,opt,name=meta,proto3" json:"meta,omitempty"`
	Content       *content.Content `protobuf:"bytes,6,opt,name=content,proto3" json:"content,omitempty"`
	Access        *PageAccess      `protobuf:"bytes,7,opt,name=access,proto3" json:"access,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Page) Reset() {
	*x = Page{}
	mi := &file_resources_wiki_page_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Page) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Page) ProtoMessage() {}

func (x *Page) ProtoReflect() protoreflect.Message {
	mi := &file_resources_wiki_page_proto_msgTypes[0]
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
	return file_resources_wiki_page_proto_rawDescGZIP(), []int{0}
}

func (x *Page) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Page) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *Page) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
	}
	return ""
}

func (x *Page) GetParentId() uint64 {
	if x != nil && x.ParentId != nil {
		return *x.ParentId
	}
	return 0
}

func (x *Page) GetMeta() *PageMeta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *Page) GetContent() *content.Content {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *Page) GetAccess() *PageAccess {
	if x != nil {
		return x.Access
	}
	return nil
}

type PageMeta struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,1,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	DeletedAt *timestamp.Timestamp   `protobuf:"bytes,3,opt,name=deleted_at,json=deletedAt,proto3,oneof" json:"deleted_at,omitempty"`
	// @sanitize: method=StripTags
	Slug *string `protobuf:"bytes,4,opt,name=slug,proto3,oneof" json:"slug,omitempty"`
	// @sanitize
	Title string `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	// @sanitize: method=StripTags
	Description string              `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	CreatorId   *int32              `protobuf:"varint,7,opt,name=creator_id,json=creatorId,proto3,oneof" json:"creator_id,omitempty"`
	Creator     *users.UserShort    `protobuf:"bytes,8,opt,name=creator,proto3,oneof" json:"creator,omitempty" alias:"creator"` // @gotags: alias:"creator"
	ContentType content.ContentType `protobuf:"varint,9,opt,name=content_type,json=contentType,proto3,enum=resources.common.content.ContentType" json:"content_type,omitempty"`
	// @sanitize: method=StripTags
	Tags          []string `protobuf:"bytes,10,rep,name=tags,proto3" json:"tags,omitempty"`
	Toc           *bool    `protobuf:"varint,11,opt,name=toc,proto3,oneof" json:"toc,omitempty"`
	Public        bool     `protobuf:"varint,12,opt,name=public,proto3" json:"public,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PageMeta) Reset() {
	*x = PageMeta{}
	mi := &file_resources_wiki_page_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PageMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageMeta) ProtoMessage() {}

func (x *PageMeta) ProtoReflect() protoreflect.Message {
	mi := &file_resources_wiki_page_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageMeta.ProtoReflect.Descriptor instead.
func (*PageMeta) Descriptor() ([]byte, []int) {
	return file_resources_wiki_page_proto_rawDescGZIP(), []int{1}
}

func (x *PageMeta) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *PageMeta) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *PageMeta) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *PageMeta) GetSlug() string {
	if x != nil && x.Slug != nil {
		return *x.Slug
	}
	return ""
}

func (x *PageMeta) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *PageMeta) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PageMeta) GetCreatorId() int32 {
	if x != nil && x.CreatorId != nil {
		return *x.CreatorId
	}
	return 0
}

func (x *PageMeta) GetCreator() *users.UserShort {
	if x != nil {
		return x.Creator
	}
	return nil
}

func (x *PageMeta) GetContentType() content.ContentType {
	if x != nil {
		return x.ContentType
	}
	return content.ContentType(0)
}

func (x *PageMeta) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *PageMeta) GetToc() bool {
	if x != nil && x.Toc != nil {
		return *x.Toc
	}
	return false
}

func (x *PageMeta) GetPublic() bool {
	if x != nil {
		return x.Public
	}
	return false
}

type PageShort struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	Job       string                 `protobuf:"bytes,2,opt,name=job,proto3" json:"job,omitempty"`
	JobLabel  *string                `protobuf:"bytes,3,opt,name=job_label,json=jobLabel,proto3,oneof" json:"job_label,omitempty"`
	ParentId  *uint64                `protobuf:"varint,4,opt,name=parent_id,json=parentId,proto3,oneof" json:"parent_id,omitempty"`
	DeletedAt *timestamp.Timestamp   `protobuf:"bytes,5,opt,name=deleted_at,json=deletedAt,proto3,oneof" json:"deleted_at,omitempty"`
	// @sanitize: method=StripTags
	Slug          *string       `protobuf:"bytes,6,opt,name=slug,proto3,oneof" json:"slug,omitempty"`
	Title         string        `protobuf:"bytes,7,opt,name=title,proto3" json:"title,omitempty"`
	Description   string        `protobuf:"bytes,8,opt,name=description,proto3" json:"description,omitempty"`
	Children      []*PageShort  `protobuf:"bytes,9,rep,name=children,proto3" json:"children,omitempty"`
	RootInfo      *PageRootInfo `protobuf:"bytes,10,opt,name=root_info,json=rootInfo,proto3,oneof" json:"root_info,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PageShort) Reset() {
	*x = PageShort{}
	mi := &file_resources_wiki_page_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PageShort) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageShort) ProtoMessage() {}

func (x *PageShort) ProtoReflect() protoreflect.Message {
	mi := &file_resources_wiki_page_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageShort.ProtoReflect.Descriptor instead.
func (*PageShort) Descriptor() ([]byte, []int) {
	return file_resources_wiki_page_proto_rawDescGZIP(), []int{2}
}

func (x *PageShort) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PageShort) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *PageShort) GetJobLabel() string {
	if x != nil && x.JobLabel != nil {
		return *x.JobLabel
	}
	return ""
}

func (x *PageShort) GetParentId() uint64 {
	if x != nil && x.ParentId != nil {
		return *x.ParentId
	}
	return 0
}

func (x *PageShort) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *PageShort) GetSlug() string {
	if x != nil && x.Slug != nil {
		return *x.Slug
	}
	return ""
}

func (x *PageShort) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *PageShort) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PageShort) GetChildren() []*PageShort {
	if x != nil {
		return x.Children
	}
	return nil
}

func (x *PageShort) GetRootInfo() *PageRootInfo {
	if x != nil {
		return x.RootInfo
	}
	return nil
}

type PageRootInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Logo          *filestore.File        `protobuf:"bytes,1,opt,name=logo,proto3,oneof" json:"logo,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PageRootInfo) Reset() {
	*x = PageRootInfo{}
	mi := &file_resources_wiki_page_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PageRootInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageRootInfo) ProtoMessage() {}

func (x *PageRootInfo) ProtoReflect() protoreflect.Message {
	mi := &file_resources_wiki_page_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageRootInfo.ProtoReflect.Descriptor instead.
func (*PageRootInfo) Descriptor() ([]byte, []int) {
	return file_resources_wiki_page_proto_rawDescGZIP(), []int{3}
}

func (x *PageRootInfo) GetLogo() *filestore.File {
	if x != nil {
		return x.Logo
	}
	return nil
}

var File_resources_wiki_page_proto protoreflect.FileDescriptor

var file_resources_wiki_page_proto_rawDesc = string([]byte{
	0x0a, 0x19, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x77, 0x69, 0x6b, 0x69,
	0x2f, 0x70, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x77, 0x69, 0x6b, 0x69, 0x1a, 0x26, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2f, 0x77, 0x69, 0x6b, 0x69, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcd, 0x02, 0x0a, 0x04,
	0x50, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12,
	0x29, 0x0a, 0x09, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x00, 0x52, 0x08, 0x6a,
	0x6f, 0x62, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x48, 0x01, 0x52,
	0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x36, 0x0a, 0x04,
	0x6d, 0x65, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x77, 0x69, 0x6b, 0x69, 0x2e, 0x50, 0x61, 0x67, 0x65,
	0x4d, 0x65, 0x74, 0x61, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04,
	0x6d, 0x65, 0x74, 0x61, 0x12, 0x3b, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x12, 0x3c, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x77, 0x69,
	0x6b, 0x69, 0x2e, 0x50, 0x61, 0x67, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x42,
	0x0c, 0x0a, 0x0a, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x42, 0x0c, 0x0a,
	0x0a, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x22, 0x8a, 0x05, 0x0a, 0x08,
	0x50, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x3d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48,
	0x01, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12,
	0x20, 0x0a, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x18, 0x64, 0x48, 0x02, 0x52, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x88, 0x01,
	0x01, 0x12, 0x20, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x03, 0x18, 0x80, 0x08, 0x52, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x12, 0x2a, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18,
	0x80, 0x01, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x2b, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x20, 0x00, 0x48, 0x03, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x39, 0x0a, 0x07,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x48, 0x04, 0x52, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x12, 0x52, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0b,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x61, 0x67, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12,
	0x15, 0x0a, 0x03, 0x74, 0x6f, 0x63, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x48, 0x05, 0x52, 0x03,
	0x74, 0x6f, 0x63, 0x88, 0x01, 0x01, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x42, 0x0d,
	0x0a, 0x0b, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a,
	0x0b, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x07, 0x0a, 0x05,
	0x5f, 0x73, 0x6c, 0x75, 0x67, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x5f, 0x69, 0x64, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72,
	0x42, 0x06, 0x0a, 0x04, 0x5f, 0x74, 0x6f, 0x63, 0x22, 0xda, 0x03, 0x0a, 0x09, 0x50, 0x61, 0x67,
	0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x52, 0x03, 0x6a, 0x6f,
	0x62, 0x12, 0x29, 0x0a, 0x09, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x00, 0x52,
	0x08, 0x6a, 0x6f, 0x62, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09,
	0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x48,
	0x01, 0x52, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x42,
	0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x48, 0x02, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x20, 0x0a, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x64, 0x48, 0x03, 0x52, 0x04, 0x73, 0x6c, 0x75,
	0x67, 0x88, 0x01, 0x01, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x35, 0x0a, 0x08,
	0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x77, 0x69, 0x6b, 0x69, 0x2e,
	0x50, 0x61, 0x67, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64,
	0x72, 0x65, 0x6e, 0x12, 0x3e, 0x0a, 0x09, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x77, 0x69, 0x6b, 0x69, 0x2e, 0x50, 0x61, 0x67, 0x65, 0x52, 0x6f, 0x6f, 0x74,
	0x49, 0x6e, 0x66, 0x6f, 0x48, 0x04, 0x52, 0x08, 0x72, 0x6f, 0x6f, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x88, 0x01, 0x01, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x42,
	0x0d, 0x0a, 0x0b, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x07,
	0x0a, 0x05, 0x5f, 0x73, 0x6c, 0x75, 0x67, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x72, 0x6f, 0x6f, 0x74,
	0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x4b, 0x0a, 0x0c, 0x50, 0x61, 0x67, 0x65, 0x52, 0x6f, 0x6f,
	0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x32, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x00,
	0x52, 0x04, 0x6c, 0x6f, 0x67, 0x6f, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6c, 0x6f,
	0x67, 0x6f, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76,
	0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x77, 0x69, 0x6b, 0x69,
	0x3b, 0x77, 0x69, 0x6b, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_wiki_page_proto_rawDescOnce sync.Once
	file_resources_wiki_page_proto_rawDescData []byte
)

func file_resources_wiki_page_proto_rawDescGZIP() []byte {
	file_resources_wiki_page_proto_rawDescOnce.Do(func() {
		file_resources_wiki_page_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_wiki_page_proto_rawDesc), len(file_resources_wiki_page_proto_rawDesc)))
	})
	return file_resources_wiki_page_proto_rawDescData
}

var file_resources_wiki_page_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_resources_wiki_page_proto_goTypes = []any{
	(*Page)(nil),                // 0: resources.wiki.Page
	(*PageMeta)(nil),            // 1: resources.wiki.PageMeta
	(*PageShort)(nil),           // 2: resources.wiki.PageShort
	(*PageRootInfo)(nil),        // 3: resources.wiki.PageRootInfo
	(*content.Content)(nil),     // 4: resources.common.content.Content
	(*PageAccess)(nil),          // 5: resources.wiki.PageAccess
	(*timestamp.Timestamp)(nil), // 6: resources.timestamp.Timestamp
	(*users.UserShort)(nil),     // 7: resources.users.UserShort
	(content.ContentType)(0),    // 8: resources.common.content.ContentType
	(*filestore.File)(nil),      // 9: resources.filestore.File
}
var file_resources_wiki_page_proto_depIdxs = []int32{
	1,  // 0: resources.wiki.Page.meta:type_name -> resources.wiki.PageMeta
	4,  // 1: resources.wiki.Page.content:type_name -> resources.common.content.Content
	5,  // 2: resources.wiki.Page.access:type_name -> resources.wiki.PageAccess
	6,  // 3: resources.wiki.PageMeta.created_at:type_name -> resources.timestamp.Timestamp
	6,  // 4: resources.wiki.PageMeta.updated_at:type_name -> resources.timestamp.Timestamp
	6,  // 5: resources.wiki.PageMeta.deleted_at:type_name -> resources.timestamp.Timestamp
	7,  // 6: resources.wiki.PageMeta.creator:type_name -> resources.users.UserShort
	8,  // 7: resources.wiki.PageMeta.content_type:type_name -> resources.common.content.ContentType
	6,  // 8: resources.wiki.PageShort.deleted_at:type_name -> resources.timestamp.Timestamp
	2,  // 9: resources.wiki.PageShort.children:type_name -> resources.wiki.PageShort
	3,  // 10: resources.wiki.PageShort.root_info:type_name -> resources.wiki.PageRootInfo
	9,  // 11: resources.wiki.PageRootInfo.logo:type_name -> resources.filestore.File
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_resources_wiki_page_proto_init() }
func file_resources_wiki_page_proto_init() {
	if File_resources_wiki_page_proto != nil {
		return
	}
	file_resources_wiki_access_proto_init()
	file_resources_wiki_page_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_wiki_page_proto_msgTypes[1].OneofWrappers = []any{}
	file_resources_wiki_page_proto_msgTypes[2].OneofWrappers = []any{}
	file_resources_wiki_page_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_wiki_page_proto_rawDesc), len(file_resources_wiki_page_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_wiki_page_proto_goTypes,
		DependencyIndexes: file_resources_wiki_page_proto_depIdxs,
		MessageInfos:      file_resources_wiki_page_proto_msgTypes,
	}.Build()
	File_resources_wiki_page_proto = out.File
	file_resources_wiki_page_proto_goTypes = nil
	file_resources_wiki_page_proto_depIdxs = nil
}
