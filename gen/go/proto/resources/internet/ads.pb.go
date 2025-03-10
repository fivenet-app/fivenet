// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: resources/internet/ads.proto

package internet

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	filestore "github.com/fivenet-app/fivenet/gen/go/proto/resources/filestore"
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

type AdType int32

const (
	AdType_AD_TYPE_UNSPECIFIED   AdType = 0
	AdType_AD_TYPE_SPONSORED     AdType = 1
	AdType_AD_TYPE_SEARCH_RESULT AdType = 2
	AdType_AD_TYPE_CONTENT_MAIN  AdType = 3
	AdType_AD_TYPE_CONTENT_ASIDE AdType = 4
)

// Enum value maps for AdType.
var (
	AdType_name = map[int32]string{
		0: "AD_TYPE_UNSPECIFIED",
		1: "AD_TYPE_SPONSORED",
		2: "AD_TYPE_SEARCH_RESULT",
		3: "AD_TYPE_CONTENT_MAIN",
		4: "AD_TYPE_CONTENT_ASIDE",
	}
	AdType_value = map[string]int32{
		"AD_TYPE_UNSPECIFIED":   0,
		"AD_TYPE_SPONSORED":     1,
		"AD_TYPE_SEARCH_RESULT": 2,
		"AD_TYPE_CONTENT_MAIN":  3,
		"AD_TYPE_CONTENT_ASIDE": 4,
	}
)

func (x AdType) Enum() *AdType {
	p := new(AdType)
	*p = x
	return p
}

func (x AdType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AdType) Descriptor() protoreflect.EnumDescriptor {
	return file_resources_internet_ads_proto_enumTypes[0].Descriptor()
}

func (AdType) Type() protoreflect.EnumType {
	return &file_resources_internet_ads_proto_enumTypes[0]
}

func (x AdType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AdType.Descriptor instead.
func (AdType) EnumDescriptor() ([]byte, []int) {
	return file_resources_internet_ads_proto_rawDescGZIP(), []int{0}
}

type Ad struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" sql:"primary_key" alias:"id"` // @gotags: sql:"primary_key" alias:"id"
	CreatedAt *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp   `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3,oneof" json:"updated_at,omitempty"`
	DeletedAt *timestamp.Timestamp   `protobuf:"bytes,4,opt,name=deleted_at,json=deletedAt,proto3,oneof" json:"deleted_at,omitempty"`
	Disabled  bool                   `protobuf:"varint,5,opt,name=disabled,proto3" json:"disabled,omitempty"`
	AdType    AdType                 `protobuf:"varint,6,opt,name=ad_type,json=adType,proto3,enum=resources.internet.AdType" json:"ad_type,omitempty"`
	StartsAt  *timestamp.Timestamp   `protobuf:"bytes,7,opt,name=starts_at,json=startsAt,proto3,oneof" json:"starts_at,omitempty"`
	EndsAt    *timestamp.Timestamp   `protobuf:"bytes,8,opt,name=ends_at,json=endsAt,proto3,oneof" json:"ends_at,omitempty"`
	// @sanitize: method=StripTags
	Title string `protobuf:"bytes,9,opt,name=title,proto3" json:"title,omitempty"`
	// @sanitize: method=StripTags
	Description   string          `protobuf:"bytes,10,opt,name=description,proto3" json:"description,omitempty"`
	Image         *filestore.File `protobuf:"bytes,11,opt,name=image,proto3,oneof" json:"image,omitempty"`
	ApproverId    *int32          `protobuf:"varint,12,opt,name=approver_id,json=approverId,proto3,oneof" json:"approver_id,omitempty"`
	ApproverJob   *string         `protobuf:"bytes,13,opt,name=approver_job,json=approverJob,proto3,oneof" json:"approver_job,omitempty"`
	CreatorId     *int32          `protobuf:"varint,14,opt,name=creator_id,json=creatorId,proto3,oneof" json:"creator_id,omitempty"`
	CreatorJob    *string         `protobuf:"bytes,15,opt,name=creator_job,json=creatorJob,proto3,oneof" json:"creator_job,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Ad) Reset() {
	*x = Ad{}
	mi := &file_resources_internet_ads_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Ad) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ad) ProtoMessage() {}

func (x *Ad) ProtoReflect() protoreflect.Message {
	mi := &file_resources_internet_ads_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ad.ProtoReflect.Descriptor instead.
func (*Ad) Descriptor() ([]byte, []int) {
	return file_resources_internet_ads_proto_rawDescGZIP(), []int{0}
}

func (x *Ad) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Ad) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Ad) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Ad) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

func (x *Ad) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *Ad) GetAdType() AdType {
	if x != nil {
		return x.AdType
	}
	return AdType_AD_TYPE_UNSPECIFIED
}

func (x *Ad) GetStartsAt() *timestamp.Timestamp {
	if x != nil {
		return x.StartsAt
	}
	return nil
}

func (x *Ad) GetEndsAt() *timestamp.Timestamp {
	if x != nil {
		return x.EndsAt
	}
	return nil
}

func (x *Ad) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Ad) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Ad) GetImage() *filestore.File {
	if x != nil {
		return x.Image
	}
	return nil
}

func (x *Ad) GetApproverId() int32 {
	if x != nil && x.ApproverId != nil {
		return *x.ApproverId
	}
	return 0
}

func (x *Ad) GetApproverJob() string {
	if x != nil && x.ApproverJob != nil {
		return *x.ApproverJob
	}
	return ""
}

func (x *Ad) GetCreatorId() int32 {
	if x != nil && x.CreatorId != nil {
		return *x.CreatorId
	}
	return 0
}

func (x *Ad) GetCreatorJob() string {
	if x != nil && x.CreatorJob != nil {
		return *x.CreatorJob
	}
	return ""
}

var File_resources_internet_ads_proto protoreflect.FileDescriptor

var file_resources_internet_ads_proto_rawDesc = string([]byte{
	0x0a, 0x1c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x65, 0x74, 0x2f, 0x61, 0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x65, 0x74, 0x1a, 0x1e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xd6, 0x06, 0x0a, 0x02, 0x41, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x42, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x42, 0x0a, 0x0a, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01,
	0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1a,
	0x0a, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x3d, 0x0a, 0x07, 0x61, 0x64,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74,
	0x2e, 0x41, 0x64, 0x54, 0x79, 0x70, 0x65, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10,
	0x01, 0x52, 0x06, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x40, 0x0a, 0x09, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x73, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x02, 0x52, 0x08,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x73, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x3c, 0x0a, 0x07, 0x65,
	0x6e, 0x64, 0x73, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x03, 0x52, 0x06,
	0x65, 0x6e, 0x64, 0x73, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10,
	0x03, 0x18, 0xff, 0x01, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x2c, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x03, 0x18, 0x80, 0x08, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x34, 0x0a, 0x05, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x48, 0x04, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x24, 0x0a, 0x0b, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x05, 0x52, 0x0a, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72,
	0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0c, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65,
	0x72, 0x5f, 0x6a, 0x6f, 0x62, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x48, 0x06, 0x52, 0x0b, 0x61,
	0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x4a, 0x6f, 0x62, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x05, 0x48, 0x07, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x24, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x6a, 0x6f, 0x62,
	0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x48, 0x08, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x4a, 0x6f, 0x62, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x73,
	0x5f, 0x61, 0x74, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x65, 0x6e, 0x64, 0x73, 0x5f, 0x61, 0x74, 0x42,
	0x08, 0x0a, 0x06, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x61, 0x70,
	0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x61, 0x70,
	0x70, 0x72, 0x6f, 0x76, 0x65, 0x72, 0x5f, 0x6a, 0x6f, 0x62, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x6a, 0x6f, 0x62, 0x2a, 0x88, 0x01, 0x0a, 0x06, 0x41, 0x64,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x41, 0x44, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a,
	0x11, 0x41, 0x44, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x4f, 0x52,
	0x45, 0x44, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x41, 0x44, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x53, 0x45, 0x41, 0x52, 0x43, 0x48, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x10, 0x02, 0x12,
	0x18, 0x0a, 0x14, 0x41, 0x44, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x45,
	0x4e, 0x54, 0x5f, 0x4d, 0x41, 0x49, 0x4e, 0x10, 0x03, 0x12, 0x19, 0x0a, 0x15, 0x41, 0x44, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x41, 0x53, 0x49,
	0x44, 0x45, 0x10, 0x04, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66,
	0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x3b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_internet_ads_proto_rawDescOnce sync.Once
	file_resources_internet_ads_proto_rawDescData []byte
)

func file_resources_internet_ads_proto_rawDescGZIP() []byte {
	file_resources_internet_ads_proto_rawDescOnce.Do(func() {
		file_resources_internet_ads_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_internet_ads_proto_rawDesc), len(file_resources_internet_ads_proto_rawDesc)))
	})
	return file_resources_internet_ads_proto_rawDescData
}

var file_resources_internet_ads_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_resources_internet_ads_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_internet_ads_proto_goTypes = []any{
	(AdType)(0),                 // 0: resources.internet.AdType
	(*Ad)(nil),                  // 1: resources.internet.Ad
	(*timestamp.Timestamp)(nil), // 2: resources.timestamp.Timestamp
	(*filestore.File)(nil),      // 3: resources.filestore.File
}
var file_resources_internet_ads_proto_depIdxs = []int32{
	2, // 0: resources.internet.Ad.created_at:type_name -> resources.timestamp.Timestamp
	2, // 1: resources.internet.Ad.updated_at:type_name -> resources.timestamp.Timestamp
	2, // 2: resources.internet.Ad.deleted_at:type_name -> resources.timestamp.Timestamp
	0, // 3: resources.internet.Ad.ad_type:type_name -> resources.internet.AdType
	2, // 4: resources.internet.Ad.starts_at:type_name -> resources.timestamp.Timestamp
	2, // 5: resources.internet.Ad.ends_at:type_name -> resources.timestamp.Timestamp
	3, // 6: resources.internet.Ad.image:type_name -> resources.filestore.File
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_resources_internet_ads_proto_init() }
func file_resources_internet_ads_proto_init() {
	if File_resources_internet_ads_proto != nil {
		return
	}
	file_resources_internet_ads_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_internet_ads_proto_rawDesc), len(file_resources_internet_ads_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_internet_ads_proto_goTypes,
		DependencyIndexes: file_resources_internet_ads_proto_depIdxs,
		EnumInfos:         file_resources_internet_ads_proto_enumTypes,
		MessageInfos:      file_resources_internet_ads_proto_msgTypes,
	}.Build()
	File_resources_internet_ads_proto = out.File
	file_resources_internet_ads_proto_goTypes = nil
	file_resources_internet_ads_proto_depIdxs = nil
}
