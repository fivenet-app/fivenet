// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.20.3
// source: resources/documents/requests.proto

package documents

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

type DocRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt       *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt       *timestamp.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DocumentId      uint64               `protobuf:"varint,4,opt,name=document_id,json=documentId,proto3" json:"document_id,omitempty"`
	RequestType     DocActivityType      `protobuf:"varint,5,opt,name=request_type,json=requestType,proto3,enum=resources.documents.DocActivityType" json:"request_type,omitempty"`
	CreatorId       *int32               `protobuf:"varint,6,opt,name=creator_id,json=creatorId,proto3,oneof" json:"creator_id,omitempty"`
	Creator         *users.UserShort     `protobuf:"bytes,7,opt,name=creator,proto3,oneof" json:"creator,omitempty" alias:"creator"` // @gotags: alias:"creator"
	CreatorJob      string               `protobuf:"bytes,8,opt,name=creator_job,json=creatorJob,proto3" json:"creator_job,omitempty"`
	CreatorJobLabel *string              `protobuf:"bytes,9,opt,name=creator_job_label,json=creatorJobLabel,proto3,oneof" json:"creator_job_label,omitempty"`
	Reason          *string              `protobuf:"bytes,10,opt,name=reason,proto3,oneof" json:"reason,omitempty"`
	Data            *DocActivityData     `protobuf:"bytes,11,opt,name=data,proto3" json:"data,omitempty"`
	Accepted        *bool                `protobuf:"varint,12,opt,name=accepted,proto3,oneof" json:"accepted,omitempty"`
}

func (x *DocRequest) Reset() {
	*x = DocRequest{}
	mi := &file_resources_documents_requests_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DocRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocRequest) ProtoMessage() {}

func (x *DocRequest) ProtoReflect() protoreflect.Message {
	mi := &file_resources_documents_requests_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocRequest.ProtoReflect.Descriptor instead.
func (*DocRequest) Descriptor() ([]byte, []int) {
	return file_resources_documents_requests_proto_rawDescGZIP(), []int{0}
}

func (x *DocRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DocRequest) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *DocRequest) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *DocRequest) GetDocumentId() uint64 {
	if x != nil {
		return x.DocumentId
	}
	return 0
}

func (x *DocRequest) GetRequestType() DocActivityType {
	if x != nil {
		return x.RequestType
	}
	return DocActivityType_DOC_ACTIVITY_TYPE_UNSPECIFIED
}

func (x *DocRequest) GetCreatorId() int32 {
	if x != nil && x.CreatorId != nil {
		return *x.CreatorId
	}
	return 0
}

func (x *DocRequest) GetCreator() *users.UserShort {
	if x != nil {
		return x.Creator
	}
	return nil
}

func (x *DocRequest) GetCreatorJob() string {
	if x != nil {
		return x.CreatorJob
	}
	return ""
}

func (x *DocRequest) GetCreatorJobLabel() string {
	if x != nil && x.CreatorJobLabel != nil {
		return *x.CreatorJobLabel
	}
	return ""
}

func (x *DocRequest) GetReason() string {
	if x != nil && x.Reason != nil {
		return *x.Reason
	}
	return ""
}

func (x *DocRequest) GetData() *DocActivityData {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DocRequest) GetAccepted() bool {
	if x != nil && x.Accepted != nil {
		return *x.Accepted
	}
	return false
}

var File_resources_documents_requests_proto protoreflect.FileDescriptor

var file_resources_documents_requests_proto_rawDesc = []byte{
	0x0a, 0x22, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x22, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x61,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xae, 0x05, 0x0a, 0x0a, 0x44, 0x6f, 0x63,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x42, 0x02, 0x30, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3d, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x3d, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x23, 0x0a, 0x0b, 0x64, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02,
	0x30, 0x01, 0x52, 0x0a, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x5b,
	0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x44, 0x6f, 0x63, 0x41, 0x63,
	0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x42, 0x12, 0xfa, 0x42, 0x0f, 0x82,
	0x01, 0x0c, 0x18, 0x0d, 0x18, 0x0e, 0x18, 0x0f, 0x18, 0x10, 0x18, 0x11, 0x18, 0x12, 0x52, 0x0b,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x48,
	0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12,
	0x39, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x48, 0x01, 0x52, 0x07,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x12, 0x28, 0x0a, 0x0b, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x6a, 0x6f, 0x62, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x14, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x4a, 0x6f, 0x62, 0x12, 0x38, 0x0a, 0x11, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f,
	0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x32, 0x48, 0x02, 0x52, 0x0f, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x6f, 0x72, 0x4a, 0x6f, 0x62, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x25,
	0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01, 0x48, 0x03, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x38, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x44, 0x6f, 0x63, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x69, 0x74, 0x79, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x1f, 0x0a, 0x08, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x08, 0x48, 0x04, 0x52, 0x08, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x88, 0x01, 0x01,
	0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x42,
	0x0a, 0x0a, 0x08, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x42, 0x14, 0x0a, 0x12, 0x5f,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x42, 0x0b, 0x0a, 0x09,
	0x5f, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x42, 0x4b, 0x5a, 0x49, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d,
	0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x3b, 0x64, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_documents_requests_proto_rawDescOnce sync.Once
	file_resources_documents_requests_proto_rawDescData = file_resources_documents_requests_proto_rawDesc
)

func file_resources_documents_requests_proto_rawDescGZIP() []byte {
	file_resources_documents_requests_proto_rawDescOnce.Do(func() {
		file_resources_documents_requests_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_documents_requests_proto_rawDescData)
	})
	return file_resources_documents_requests_proto_rawDescData
}

var file_resources_documents_requests_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_documents_requests_proto_goTypes = []any{
	(*DocRequest)(nil),          // 0: resources.documents.DocRequest
	(*timestamp.Timestamp)(nil), // 1: resources.timestamp.Timestamp
	(DocActivityType)(0),        // 2: resources.documents.DocActivityType
	(*users.UserShort)(nil),     // 3: resources.users.UserShort
	(*DocActivityData)(nil),     // 4: resources.documents.DocActivityData
}
var file_resources_documents_requests_proto_depIdxs = []int32{
	1, // 0: resources.documents.DocRequest.created_at:type_name -> resources.timestamp.Timestamp
	1, // 1: resources.documents.DocRequest.updated_at:type_name -> resources.timestamp.Timestamp
	2, // 2: resources.documents.DocRequest.request_type:type_name -> resources.documents.DocActivityType
	3, // 3: resources.documents.DocRequest.creator:type_name -> resources.users.UserShort
	4, // 4: resources.documents.DocRequest.data:type_name -> resources.documents.DocActivityData
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_resources_documents_requests_proto_init() }
func file_resources_documents_requests_proto_init() {
	if File_resources_documents_requests_proto != nil {
		return
	}
	file_resources_documents_activity_proto_init()
	file_resources_documents_requests_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_documents_requests_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_documents_requests_proto_goTypes,
		DependencyIndexes: file_resources_documents_requests_proto_depIdxs,
		MessageInfos:      file_resources_documents_requests_proto_msgTypes,
	}.Build()
	File_resources_documents_requests_proto = out.File
	file_resources_documents_requests_proto_rawDesc = nil
	file_resources_documents_requests_proto_goTypes = nil
	file_resources_documents_requests_proto_depIdxs = nil
}
