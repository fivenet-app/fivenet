// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: services/settings/filestore.proto

package settings

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	filestore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/filestore"
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

type ListFilesRequest struct {
	state         protoimpl.MessageState      `protogen:"open.v1"`
	Pagination    *database.PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Path          *string                     `protobuf:"bytes,2,opt,name=path,proto3,oneof" json:"path,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListFilesRequest) Reset() {
	*x = ListFilesRequest{}
	mi := &file_services_settings_filestore_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListFilesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFilesRequest) ProtoMessage() {}

func (x *ListFilesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_settings_filestore_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFilesRequest.ProtoReflect.Descriptor instead.
func (*ListFilesRequest) Descriptor() ([]byte, []int) {
	return file_services_settings_filestore_proto_rawDescGZIP(), []int{0}
}

func (x *ListFilesRequest) GetPagination() *database.PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListFilesRequest) GetPath() string {
	if x != nil && x.Path != nil {
		return *x.Path
	}
	return ""
}

type ListFilesResponse struct {
	state         protoimpl.MessageState       `protogen:"open.v1"`
	Pagination    *database.PaginationResponse `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	Files         []*filestore.FileInfo        `protobuf:"bytes,2,rep,name=files,proto3" json:"files,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListFilesResponse) Reset() {
	*x = ListFilesResponse{}
	mi := &file_services_settings_filestore_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListFilesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFilesResponse) ProtoMessage() {}

func (x *ListFilesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_settings_filestore_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFilesResponse.ProtoReflect.Descriptor instead.
func (*ListFilesResponse) Descriptor() ([]byte, []int) {
	return file_services_settings_filestore_proto_rawDescGZIP(), []int{1}
}

func (x *ListFilesResponse) GetPagination() *database.PaginationResponse {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *ListFilesResponse) GetFiles() []*filestore.FileInfo {
	if x != nil {
		return x.Files
	}
	return nil
}

type UploadFileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Prefix        string                 `protobuf:"bytes,1,opt,name=prefix,proto3" json:"prefix,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	File          *filestore.File        `protobuf:"bytes,3,opt,name=file,proto3" json:"file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadFileRequest) Reset() {
	*x = UploadFileRequest{}
	mi := &file_services_settings_filestore_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileRequest) ProtoMessage() {}

func (x *UploadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_settings_filestore_proto_msgTypes[2]
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
	return file_services_settings_filestore_proto_rawDescGZIP(), []int{2}
}

func (x *UploadFileRequest) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

func (x *UploadFileRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UploadFileRequest) GetFile() *filestore.File {
	if x != nil {
		return x.File
	}
	return nil
}

type UploadFileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	File          *filestore.FileInfo    `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadFileResponse) Reset() {
	*x = UploadFileResponse{}
	mi := &file_services_settings_filestore_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileResponse) ProtoMessage() {}

func (x *UploadFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_settings_filestore_proto_msgTypes[3]
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
	return file_services_settings_filestore_proto_rawDescGZIP(), []int{3}
}

func (x *UploadFileResponse) GetFile() *filestore.FileInfo {
	if x != nil {
		return x.File
	}
	return nil
}

type DeleteFileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Path          string                 `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFileRequest) Reset() {
	*x = DeleteFileRequest{}
	mi := &file_services_settings_filestore_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileRequest) ProtoMessage() {}

func (x *DeleteFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_settings_filestore_proto_msgTypes[4]
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
	return file_services_settings_filestore_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteFileRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type DeleteFileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFileResponse) Reset() {
	*x = DeleteFileResponse{}
	mi := &file_services_settings_filestore_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileResponse) ProtoMessage() {}

func (x *DeleteFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_settings_filestore_proto_msgTypes[5]
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
	return file_services_settings_filestore_proto_rawDescGZIP(), []int{5}
}

var File_services_settings_filestore_proto protoreflect.FileDescriptor

const file_services_settings_filestore_proto_rawDesc = "" +
	"\n" +
	"!services/settings/filestore.proto\x12\x11services.settings\x1a(resources/common/database/database.proto\x1a\x1eresources/filestore/file.proto\x1a\x17validate/validate.proto\"\x96\x01\n" +
	"\x10ListFilesRequest\x12V\n" +
	"\n" +
	"pagination\x18\x01 \x01(\v2,.resources.common.database.PaginationRequestB\b\xfaB\x05\x8a\x01\x02\x10\x01R\n" +
	"pagination\x12!\n" +
	"\x04path\x18\x02 \x01(\tB\b\xfaB\x05r\x03\x18\x80\x01H\x00R\x04path\x88\x01\x01B\a\n" +
	"\x05_path\"\x97\x01\n" +
	"\x11ListFilesResponse\x12M\n" +
	"\n" +
	"pagination\x18\x01 \x01(\v2-.resources.common.database.PaginationResponseR\n" +
	"pagination\x123\n" +
	"\x05files\x18\x02 \x03(\v2\x1d.resources.filestore.FileInfoR\x05files\"x\n" +
	"\x11UploadFileRequest\x12\x16\n" +
	"\x06prefix\x18\x01 \x01(\tR\x06prefix\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x127\n" +
	"\x04file\x18\x03 \x01(\v2\x19.resources.filestore.FileB\b\xfaB\x05\x8a\x01\x02\x10\x01R\x04file\"G\n" +
	"\x12UploadFileResponse\x121\n" +
	"\x04file\x18\x01 \x01(\v2\x1d.resources.filestore.FileInfoR\x04file\"'\n" +
	"\x11DeleteFileRequest\x12\x12\n" +
	"\x04path\x18\x01 \x01(\tR\x04path\"\x14\n" +
	"\x12DeleteFileResponse2\xa0\x02\n" +
	"\x10FilestoreService\x12V\n" +
	"\tListFiles\x12#.services.settings.ListFilesRequest\x1a$.services.settings.ListFilesResponse\x12Y\n" +
	"\n" +
	"UploadFile\x12$.services.settings.UploadFileRequest\x1a%.services.settings.UploadFileResponse\x12Y\n" +
	"\n" +
	"DeleteFile\x12$.services.settings.DeleteFileRequest\x1a%.services.settings.DeleteFileResponseBNZLgithub.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings;settingsb\x06proto3"

var (
	file_services_settings_filestore_proto_rawDescOnce sync.Once
	file_services_settings_filestore_proto_rawDescData []byte
)

func file_services_settings_filestore_proto_rawDescGZIP() []byte {
	file_services_settings_filestore_proto_rawDescOnce.Do(func() {
		file_services_settings_filestore_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_services_settings_filestore_proto_rawDesc), len(file_services_settings_filestore_proto_rawDesc)))
	})
	return file_services_settings_filestore_proto_rawDescData
}

var file_services_settings_filestore_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_services_settings_filestore_proto_goTypes = []any{
	(*ListFilesRequest)(nil),            // 0: services.settings.ListFilesRequest
	(*ListFilesResponse)(nil),           // 1: services.settings.ListFilesResponse
	(*UploadFileRequest)(nil),           // 2: services.settings.UploadFileRequest
	(*UploadFileResponse)(nil),          // 3: services.settings.UploadFileResponse
	(*DeleteFileRequest)(nil),           // 4: services.settings.DeleteFileRequest
	(*DeleteFileResponse)(nil),          // 5: services.settings.DeleteFileResponse
	(*database.PaginationRequest)(nil),  // 6: resources.common.database.PaginationRequest
	(*database.PaginationResponse)(nil), // 7: resources.common.database.PaginationResponse
	(*filestore.FileInfo)(nil),          // 8: resources.filestore.FileInfo
	(*filestore.File)(nil),              // 9: resources.filestore.File
}
var file_services_settings_filestore_proto_depIdxs = []int32{
	6, // 0: services.settings.ListFilesRequest.pagination:type_name -> resources.common.database.PaginationRequest
	7, // 1: services.settings.ListFilesResponse.pagination:type_name -> resources.common.database.PaginationResponse
	8, // 2: services.settings.ListFilesResponse.files:type_name -> resources.filestore.FileInfo
	9, // 3: services.settings.UploadFileRequest.file:type_name -> resources.filestore.File
	8, // 4: services.settings.UploadFileResponse.file:type_name -> resources.filestore.FileInfo
	0, // 5: services.settings.FilestoreService.ListFiles:input_type -> services.settings.ListFilesRequest
	2, // 6: services.settings.FilestoreService.UploadFile:input_type -> services.settings.UploadFileRequest
	4, // 7: services.settings.FilestoreService.DeleteFile:input_type -> services.settings.DeleteFileRequest
	1, // 8: services.settings.FilestoreService.ListFiles:output_type -> services.settings.ListFilesResponse
	3, // 9: services.settings.FilestoreService.UploadFile:output_type -> services.settings.UploadFileResponse
	5, // 10: services.settings.FilestoreService.DeleteFile:output_type -> services.settings.DeleteFileResponse
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_services_settings_filestore_proto_init() }
func file_services_settings_filestore_proto_init() {
	if File_services_settings_filestore_proto != nil {
		return
	}
	file_services_settings_filestore_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_services_settings_filestore_proto_rawDesc), len(file_services_settings_filestore_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_settings_filestore_proto_goTypes,
		DependencyIndexes: file_services_settings_filestore_proto_depIdxs,
		MessageInfos:      file_services_settings_filestore_proto_msgTypes,
	}.Build()
	File_services_settings_filestore_proto = out.File
	file_services_settings_filestore_proto_goTypes = nil
	file_services_settings_filestore_proto_depIdxs = nil
}
