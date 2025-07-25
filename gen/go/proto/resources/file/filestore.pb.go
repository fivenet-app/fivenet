// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: resources/file/filestore.proto

package file

import (
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

type UploadFileRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Payload:
	//
	//	*UploadFileRequest_Meta
	//	*UploadFileRequest_Data
	Payload       isUploadFileRequest_Payload `protobuf_oneof:"payload"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadFileRequest) Reset() {
	*x = UploadFileRequest{}
	mi := &file_resources_file_filestore_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileRequest) ProtoMessage() {}

func (x *UploadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_resources_file_filestore_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileRequest.ProtoReflect.Descriptor instead.
func (*UploadFileRequest) Descriptor() ([]byte, []int) {
	return file_resources_file_filestore_proto_rawDescGZIP(), []int{0}
}

func (x *UploadFileRequest) GetPayload() isUploadFileRequest_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *UploadFileRequest) GetMeta() *UploadMeta {
	if x != nil {
		if x, ok := x.Payload.(*UploadFileRequest_Meta); ok {
			return x.Meta
		}
	}
	return nil
}

func (x *UploadFileRequest) GetData() []byte {
	if x != nil {
		if x, ok := x.Payload.(*UploadFileRequest_Data); ok {
			return x.Data
		}
	}
	return nil
}

type isUploadFileRequest_Payload interface {
	isUploadFileRequest_Payload()
}

type UploadFileRequest_Meta struct {
	Meta *UploadMeta `protobuf:"bytes,1,opt,name=meta,proto3,oneof"`
}

type UploadFileRequest_Data struct {
	// Raw bytes <= 128 KiB each, browsers should only read 64 KiB at a time, but this is a buffer just in case
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3,oneof"`
}

func (*UploadFileRequest_Meta) isUploadFileRequest_Payload() {}

func (*UploadFileRequest_Data) isUploadFileRequest_Payload() {}

type UploadMeta struct {
	state        protoimpl.MessageState `protogen:"open.v1"`
	ParentId     uint64                 `protobuf:"varint,1,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`
	Namespace    string                 `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"` // "documents", "wiki", …
	OriginalName string                 `protobuf:"bytes,3,opt,name=original_name,json=originalName,proto3" json:"original_name,omitempty"`
	ContentType  string                 `protobuf:"bytes,4,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"` // optional – server re-validates
	Size         int64                  `protobuf:"varint,5,opt,name=size,proto3" json:"size,omitempty"`                                 // Size in bytes
	// @sanitize
	Reason        string `protobuf:"bytes,6,opt,name=reason,proto3" json:"reason,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadMeta) Reset() {
	*x = UploadMeta{}
	mi := &file_resources_file_filestore_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadMeta) ProtoMessage() {}

func (x *UploadMeta) ProtoReflect() protoreflect.Message {
	mi := &file_resources_file_filestore_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadMeta.ProtoReflect.Descriptor instead.
func (*UploadMeta) Descriptor() ([]byte, []int) {
	return file_resources_file_filestore_proto_rawDescGZIP(), []int{1}
}

func (x *UploadMeta) GetParentId() uint64 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

func (x *UploadMeta) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *UploadMeta) GetOriginalName() string {
	if x != nil {
		return x.OriginalName
	}
	return ""
}

func (x *UploadMeta) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *UploadMeta) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *UploadMeta) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type UploadFileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`    // Unique ID for the uploaded file
	Url           string                 `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`   // URL to the uploaded file
	File          *File                  `protobuf:"bytes,3,opt,name=file,proto3" json:"file,omitempty"` // File info
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadFileResponse) Reset() {
	*x = UploadFileResponse{}
	mi := &file_resources_file_filestore_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileResponse) ProtoMessage() {}

func (x *UploadFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_resources_file_filestore_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileResponse.ProtoReflect.Descriptor instead.
func (*UploadFileResponse) Descriptor() ([]byte, []int) {
	return file_resources_file_filestore_proto_rawDescGZIP(), []int{2}
}

func (x *UploadFileResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UploadFileResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *UploadFileResponse) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

type DeleteFileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ParentId      uint64                 `protobuf:"varint,1,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`
	FileId        uint64                 `protobuf:"varint,2,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFileRequest) Reset() {
	*x = DeleteFileRequest{}
	mi := &file_resources_file_filestore_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileRequest) ProtoMessage() {}

func (x *DeleteFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_resources_file_filestore_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileRequest.ProtoReflect.Descriptor instead.
func (*DeleteFileRequest) Descriptor() ([]byte, []int) {
	return file_resources_file_filestore_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteFileRequest) GetParentId() uint64 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

func (x *DeleteFileRequest) GetFileId() uint64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

type DeleteFileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFileResponse) Reset() {
	*x = DeleteFileResponse{}
	mi := &file_resources_file_filestore_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileResponse) ProtoMessage() {}

func (x *DeleteFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_resources_file_filestore_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileResponse.ProtoReflect.Descriptor instead.
func (*DeleteFileResponse) Descriptor() ([]byte, []int) {
	return file_resources_file_filestore_proto_rawDescGZIP(), []int{4}
}

var File_resources_file_filestore_proto protoreflect.FileDescriptor

const file_resources_file_filestore_proto_rawDesc = "" +
	"\n" +
	"\x1eresources/file/filestore.proto\x12\x0eresources.file\x1a\x19resources/file/file.proto\"q\n" +
	"\x11UploadFileRequest\x120\n" +
	"\x04meta\x18\x01 \x01(\v2\x1a.resources.file.UploadMetaH\x00R\x04meta\x12\x1f\n" +
	"\x04data\x18\x02 \x01(\fB\t\xbaH\x06z\x04\x18\x80\x80\bH\x00R\x04dataB\t\n" +
	"\apayload\"\xdd\x01\n" +
	"\n" +
	"UploadMeta\x12\x1b\n" +
	"\tparent_id\x18\x01 \x01(\x04R\bparentId\x12\x1c\n" +
	"\tnamespace\x18\x02 \x01(\tR\tnamespace\x12-\n" +
	"\roriginal_name\x18\x03 \x01(\tB\b\xbaH\x05r\x03\x18\xff\x01R\foriginalName\x12!\n" +
	"\fcontent_type\x18\x04 \x01(\tR\vcontentType\x12\x1b\n" +
	"\x04size\x18\x05 \x01(\x03B\a\xbaH\x04\"\x02 \x00R\x04size\x12%\n" +
	"\x06reason\x18\x06 \x01(\tB\r\xbaH\n" +
	"\xd8\x01\x01r\x05\x10\x03\x18\xff\x01R\x06reason\"`\n" +
	"\x12UploadFileResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x10\n" +
	"\x03url\x18\x02 \x01(\tR\x03url\x12(\n" +
	"\x04file\x18\x03 \x01(\v2\x14.resources.file.FileR\x04file\"[\n" +
	"\x11DeleteFileRequest\x12$\n" +
	"\tparent_id\x18\x01 \x01(\x04B\a\xbaH\x042\x02 \x00R\bparentId\x12 \n" +
	"\afile_id\x18\x02 \x01(\x04B\a\xbaH\x042\x02 \x00R\x06fileId\"\x14\n" +
	"\x12DeleteFileResponseBGZEgithub.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file;fileb\x06proto3"

var (
	file_resources_file_filestore_proto_rawDescOnce sync.Once
	file_resources_file_filestore_proto_rawDescData []byte
)

func file_resources_file_filestore_proto_rawDescGZIP() []byte {
	file_resources_file_filestore_proto_rawDescOnce.Do(func() {
		file_resources_file_filestore_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_file_filestore_proto_rawDesc), len(file_resources_file_filestore_proto_rawDesc)))
	})
	return file_resources_file_filestore_proto_rawDescData
}

var file_resources_file_filestore_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_resources_file_filestore_proto_goTypes = []any{
	(*UploadFileRequest)(nil),  // 0: resources.file.UploadFileRequest
	(*UploadMeta)(nil),         // 1: resources.file.UploadMeta
	(*UploadFileResponse)(nil), // 2: resources.file.UploadFileResponse
	(*DeleteFileRequest)(nil),  // 3: resources.file.DeleteFileRequest
	(*DeleteFileResponse)(nil), // 4: resources.file.DeleteFileResponse
	(*File)(nil),               // 5: resources.file.File
}
var file_resources_file_filestore_proto_depIdxs = []int32{
	1, // 0: resources.file.UploadFileRequest.meta:type_name -> resources.file.UploadMeta
	5, // 1: resources.file.UploadFileResponse.file:type_name -> resources.file.File
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_resources_file_filestore_proto_init() }
func file_resources_file_filestore_proto_init() {
	if File_resources_file_filestore_proto != nil {
		return
	}
	file_resources_file_file_proto_init()
	file_resources_file_filestore_proto_msgTypes[0].OneofWrappers = []any{
		(*UploadFileRequest_Meta)(nil),
		(*UploadFileRequest_Data)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_file_filestore_proto_rawDesc), len(file_resources_file_filestore_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_file_filestore_proto_goTypes,
		DependencyIndexes: file_resources_file_filestore_proto_depIdxs,
		MessageInfos:      file_resources_file_filestore_proto_msgTypes,
	}.Build()
	File_resources_file_filestore_proto = out.File
	file_resources_file_filestore_proto_goTypes = nil
	file_resources_file_filestore_proto_depIdxs = nil
}
