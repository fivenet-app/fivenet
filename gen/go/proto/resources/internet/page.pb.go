// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
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

const file_resources_internet_page_proto_rawDesc = "" +
	"\n" +
	"\x1dresources/internet/page.proto\x12\x12resources.internet\x1a&resources/common/content/content.proto\x1a#resources/timestamp/timestamp.proto\x1a\x17validate/validate.proto\"\xab\x04\n" +
	"\x04Page\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12=\n" +
	"\n" +
	"created_at\x18\x02 \x01(\v2\x1e.resources.timestamp.TimestampR\tcreatedAt\x12B\n" +
	"\n" +
	"updated_at\x18\x03 \x01(\v2\x1e.resources.timestamp.TimestampH\x00R\tupdatedAt\x88\x01\x01\x12B\n" +
	"\n" +
	"deleted_at\x18\x04 \x01(\v2\x1e.resources.timestamp.TimestampH\x01R\tdeletedAt\x88\x01\x01\x12\x1b\n" +
	"\tdomain_id\x18\x05 \x01(\x04R\bdomainId\x12\x1c\n" +
	"\x04path\x18\x06 \x01(\tB\b\xfaB\x05r\x03\x18\x80\x01R\x04path\x12 \n" +
	"\x05title\x18\a \x01(\tB\n" +
	"\xfaB\ar\x05\x10\x01\x18\xff\x01R\x05title\x12,\n" +
	"\vdescription\x18\b \x01(\tB\n" +
	"\xfaB\ar\x05\x10\x03\x18\x80\x04R\vdescription\x12:\n" +
	"\x04data\x18\t \x01(\v2\x1c.resources.internet.PageDataB\b\xfaB\x05\x8a\x01\x02\x10\x01R\x04data\x12$\n" +
	"\vcreator_job\x18\n" +
	" \x01(\tH\x02R\n" +
	"creatorJob\x88\x01\x01\x12\"\n" +
	"\n" +
	"creator_id\x18\v \x01(\x05H\x03R\tcreatorId\x88\x01\x01B\r\n" +
	"\v_updated_atB\r\n" +
	"\v_deleted_atB\x0e\n" +
	"\f_creator_jobB\r\n" +
	"\v_creator_id\"\x92\x01\n" +
	"\bPageData\x12C\n" +
	"\vlayout_type\x18\x01 \x01(\x0e2\".resources.internet.PageLayoutTypeR\n" +
	"layoutType\x128\n" +
	"\x04node\x18\x02 \x01(\v2\x1f.resources.internet.ContentNodeH\x00R\x04node\x88\x01\x01B\a\n" +
	"\x05_node\"\x8d\x03\n" +
	"\vContentNode\x12@\n" +
	"\x04type\x18\x01 \x01(\x0e2\".resources.common.content.NodeTypeB\b\xfaB\x05\x82\x01\x02\x10\x01R\x04type\x12\x13\n" +
	"\x02id\x18\x02 \x01(\tH\x00R\x02id\x88\x01\x01\x12\x10\n" +
	"\x03tag\x18\x03 \x01(\tR\x03tag\x12@\n" +
	"\x05attrs\x18\x04 \x03(\v2*.resources.internet.ContentNode.AttrsEntryR\x05attrs\x12\x17\n" +
	"\x04text\x18\x05 \x01(\tH\x01R\x04text\x88\x01\x01\x129\n" +
	"\acontent\x18\x06 \x03(\v2\x1f.resources.internet.ContentNodeR\acontent\x125\n" +
	"\x05slots\x18\a \x03(\v2\x1f.resources.internet.ContentNodeR\x05slots\x1a8\n" +
	"\n" +
	"AttrsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01B\x05\n" +
	"\x03_idB\a\n" +
	"\x05_text*v\n" +
	"\x0ePageLayoutType\x12 \n" +
	"\x1cPAGE_LAYOUT_TYPE_UNSPECIFIED\x10\x00\x12\x1f\n" +
	"\x1bPAGE_LAYOUT_TYPE_BASIC_PAGE\x10\x01\x12!\n" +
	"\x1dPAGE_LAYOUT_TYPE_LANDING_PAGE\x10\x02BIZGgithub.com/fivenet-app/fivenet/gen/go/proto/resources/internet;internetb\x06proto3"

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
