// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v3.20.3
// source: services/completor/completor.proto

package completor

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	documents "github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	laws "github.com/fivenet-app/fivenet/gen/go/proto/resources/laws"
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

type CompleteCitizensRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Search        string                 `protobuf:"bytes,1,opt,name=search,proto3" json:"search,omitempty"`
	CurrentJob    *bool                  `protobuf:"varint,2,opt,name=current_job,json=currentJob,proto3,oneof" json:"current_job,omitempty"`
	OnDuty        *bool                  `protobuf:"varint,3,opt,name=on_duty,json=onDuty,proto3,oneof" json:"on_duty,omitempty"`
	UserId        *int32                 `protobuf:"varint,4,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompleteCitizensRequest) Reset() {
	*x = CompleteCitizensRequest{}
	mi := &file_services_completor_completor_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteCitizensRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteCitizensRequest) ProtoMessage() {}

func (x *CompleteCitizensRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteCitizensRequest.ProtoReflect.Descriptor instead.
func (*CompleteCitizensRequest) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{0}
}

func (x *CompleteCitizensRequest) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}

func (x *CompleteCitizensRequest) GetCurrentJob() bool {
	if x != nil && x.CurrentJob != nil {
		return *x.CurrentJob
	}
	return false
}

func (x *CompleteCitizensRequest) GetOnDuty() bool {
	if x != nil && x.OnDuty != nil {
		return *x.OnDuty
	}
	return false
}

func (x *CompleteCitizensRequest) GetUserId() int32 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

type CompleteCitizensRespoonse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Users         []*users.UserShort     `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty" alias:"user"` // @gotags: alias:"user"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompleteCitizensRespoonse) Reset() {
	*x = CompleteCitizensRespoonse{}
	mi := &file_services_completor_completor_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteCitizensRespoonse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteCitizensRespoonse) ProtoMessage() {}

func (x *CompleteCitizensRespoonse) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteCitizensRespoonse.ProtoReflect.Descriptor instead.
func (*CompleteCitizensRespoonse) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{1}
}

func (x *CompleteCitizensRespoonse) GetUsers() []*users.UserShort {
	if x != nil {
		return x.Users
	}
	return nil
}

type CompleteJobsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Search        *string                `protobuf:"bytes,1,opt,name=search,proto3,oneof" json:"search,omitempty"`
	ExactMatch    *bool                  `protobuf:"varint,2,opt,name=exact_match,json=exactMatch,proto3,oneof" json:"exact_match,omitempty"`
	CurrentJob    *bool                  `protobuf:"varint,3,opt,name=current_job,json=currentJob,proto3,oneof" json:"current_job,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompleteJobsRequest) Reset() {
	*x = CompleteJobsRequest{}
	mi := &file_services_completor_completor_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteJobsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteJobsRequest) ProtoMessage() {}

func (x *CompleteJobsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteJobsRequest.ProtoReflect.Descriptor instead.
func (*CompleteJobsRequest) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{2}
}

func (x *CompleteJobsRequest) GetSearch() string {
	if x != nil && x.Search != nil {
		return *x.Search
	}
	return ""
}

func (x *CompleteJobsRequest) GetExactMatch() bool {
	if x != nil && x.ExactMatch != nil {
		return *x.ExactMatch
	}
	return false
}

func (x *CompleteJobsRequest) GetCurrentJob() bool {
	if x != nil && x.CurrentJob != nil {
		return *x.CurrentJob
	}
	return false
}

type CompleteJobsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Jobs          []*users.Job           `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompleteJobsResponse) Reset() {
	*x = CompleteJobsResponse{}
	mi := &file_services_completor_completor_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteJobsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteJobsResponse) ProtoMessage() {}

func (x *CompleteJobsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteJobsResponse.ProtoReflect.Descriptor instead.
func (*CompleteJobsResponse) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{3}
}

func (x *CompleteJobsResponse) GetJobs() []*users.Job {
	if x != nil {
		return x.Jobs
	}
	return nil
}

type CompleteDocumentCategoriesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Search        string                 `protobuf:"bytes,1,opt,name=search,proto3" json:"search,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompleteDocumentCategoriesRequest) Reset() {
	*x = CompleteDocumentCategoriesRequest{}
	mi := &file_services_completor_completor_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteDocumentCategoriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteDocumentCategoriesRequest) ProtoMessage() {}

func (x *CompleteDocumentCategoriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteDocumentCategoriesRequest.ProtoReflect.Descriptor instead.
func (*CompleteDocumentCategoriesRequest) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{4}
}

func (x *CompleteDocumentCategoriesRequest) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}

type CompleteDocumentCategoriesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Categories    []*documents.Category  `protobuf:"bytes,1,rep,name=categories,proto3" json:"categories,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompleteDocumentCategoriesResponse) Reset() {
	*x = CompleteDocumentCategoriesResponse{}
	mi := &file_services_completor_completor_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteDocumentCategoriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteDocumentCategoriesResponse) ProtoMessage() {}

func (x *CompleteDocumentCategoriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteDocumentCategoriesResponse.ProtoReflect.Descriptor instead.
func (*CompleteDocumentCategoriesResponse) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{5}
}

func (x *CompleteDocumentCategoriesResponse) GetCategories() []*documents.Category {
	if x != nil {
		return x.Categories
	}
	return nil
}

type ListLawBooksRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListLawBooksRequest) Reset() {
	*x = ListLawBooksRequest{}
	mi := &file_services_completor_completor_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListLawBooksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLawBooksRequest) ProtoMessage() {}

func (x *ListLawBooksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLawBooksRequest.ProtoReflect.Descriptor instead.
func (*ListLawBooksRequest) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{6}
}

type ListLawBooksResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Books         []*laws.LawBook        `protobuf:"bytes,1,rep,name=books,proto3" json:"books,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListLawBooksResponse) Reset() {
	*x = ListLawBooksResponse{}
	mi := &file_services_completor_completor_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListLawBooksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLawBooksResponse) ProtoMessage() {}

func (x *ListLawBooksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLawBooksResponse.ProtoReflect.Descriptor instead.
func (*ListLawBooksResponse) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{7}
}

func (x *ListLawBooksResponse) GetBooks() []*laws.LawBook {
	if x != nil {
		return x.Books
	}
	return nil
}

type CompleteCitizenLabelsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Search        string                 `protobuf:"bytes,1,opt,name=search,proto3" json:"search,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompleteCitizenLabelsRequest) Reset() {
	*x = CompleteCitizenLabelsRequest{}
	mi := &file_services_completor_completor_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteCitizenLabelsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteCitizenLabelsRequest) ProtoMessage() {}

func (x *CompleteCitizenLabelsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteCitizenLabelsRequest.ProtoReflect.Descriptor instead.
func (*CompleteCitizenLabelsRequest) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{8}
}

func (x *CompleteCitizenLabelsRequest) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}

type CompleteCitizenLabelsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Labels        []*users.CitizenLabel  `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompleteCitizenLabelsResponse) Reset() {
	*x = CompleteCitizenLabelsResponse{}
	mi := &file_services_completor_completor_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteCitizenLabelsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteCitizenLabelsResponse) ProtoMessage() {}

func (x *CompleteCitizenLabelsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_completor_completor_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteCitizenLabelsResponse.ProtoReflect.Descriptor instead.
func (*CompleteCitizenLabelsResponse) Descriptor() ([]byte, []int) {
	return file_services_completor_completor_proto_rawDescGZIP(), []int{9}
}

func (x *CompleteCitizenLabelsResponse) GetLabels() []*users.CitizenLabel {
	if x != nil {
		return x.Labels
	}
	return nil
}

var File_services_completor_completor_proto protoreflect.FileDescriptor

var file_services_completor_completor_proto_rawDesc = []byte{
	0x0a, 0x22, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x6f, 0x72, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x1a, 0x22, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6c, 0x61, 0x77, 0x73, 0x2f, 0x6c, 0x61, 0x77,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcd, 0x01, 0x0a, 0x17, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x40, 0x52, 0x06, 0x73, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x12, 0x24, 0x0a, 0x0b, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f,
	0x6a, 0x6f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x0a, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x4a, 0x6f, 0x62, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x07, 0x6f, 0x6e,
	0x5f, 0x64, 0x75, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x06, 0x6f,
	0x6e, 0x44, 0x75, 0x74, 0x79, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02,
	0x20, 0x00, 0x48, 0x02, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42,
	0x0e, 0x0a, 0x0c, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6a, 0x6f, 0x62, 0x42,
	0x0a, 0x0a, 0x08, 0x5f, 0x6f, 0x6e, 0x5f, 0x64, 0x75, 0x74, 0x79, 0x42, 0x0a, 0x0a, 0x08, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x22, 0x4d, 0x0a, 0x19, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52,
	0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0xb2, 0x01, 0x0a, 0x13, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24,
	0x0a, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x40, 0x48, 0x00, 0x52, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x88, 0x01, 0x01, 0x12, 0x24, 0x0a, 0x0b, 0x65, 0x78, 0x61, 0x63, 0x74, 0x5f, 0x6d, 0x61,
	0x74, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x0a, 0x65, 0x78, 0x61,
	0x63, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x88, 0x01, 0x01, 0x12, 0x24, 0x0a, 0x0b, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6a, 0x6f, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48,
	0x02, 0x52, 0x0a, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4a, 0x6f, 0x62, 0x88, 0x01, 0x01,
	0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x42, 0x0e, 0x0a, 0x0c, 0x5f,
	0x65, 0x78, 0x61, 0x63, 0x74, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x42, 0x0e, 0x0a, 0x0c, 0x5f,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6a, 0x6f, 0x62, 0x22, 0x40, 0x0a, 0x14, 0x43,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x52, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x22, 0x44, 0x0a,
	0x21, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x40, 0x52, 0x06, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x22, 0x63, 0x0a, 0x22, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x44,
	0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x0a, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x22, 0x15, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74,
	0x4c, 0x61, 0x77, 0x42, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x45, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x61, 0x77, 0x42, 0x6f, 0x6f, 0x6b, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x05, 0x62, 0x6f, 0x6f, 0x6b, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x6c, 0x61, 0x77, 0x73, 0x2e, 0x4c, 0x61, 0x77, 0x42, 0x6f, 0x6f, 0x6b, 0x52,
	0x05, 0x62, 0x6f, 0x6f, 0x6b, 0x73, 0x22, 0x3f, 0x0a, 0x1c, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65,
	0x74, 0x65, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x40, 0x52,
	0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x22, 0x56, 0x0a, 0x1d, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x69, 0x74, 0x69, 0x7a,
	0x65, 0x6e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x32,
	0xd4, 0x04, 0x0a, 0x10, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x6e, 0x0a, 0x10, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x12, 0x2b, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6f,
	0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x61, 0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x4a, 0x6f, 0x62, 0x73, 0x12, 0x27, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65,
	0x74, 0x65, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x6f, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x4a, 0x6f, 0x62, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x8b, 0x01, 0x0a, 0x1a, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x35, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x6f, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x6f, 0x63, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x61, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x61, 0x77,
	0x42, 0x6f, 0x6f, 0x6b, 0x73, 0x12, 0x27, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c,
	0x61, 0x77, 0x42, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65,
	0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x61, 0x77, 0x42, 0x6f, 0x6f, 0x6b, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x7c, 0x0a, 0x15, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x12, 0x30, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x43,
	0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x63,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x65, 0x43, 0x69, 0x74, 0x69, 0x7a, 0x65, 0x6e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4a, 0x5a, 0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70,
	0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x6f, 0x72, 0x3b, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_completor_completor_proto_rawDescOnce sync.Once
	file_services_completor_completor_proto_rawDescData = file_services_completor_completor_proto_rawDesc
)

func file_services_completor_completor_proto_rawDescGZIP() []byte {
	file_services_completor_completor_proto_rawDescOnce.Do(func() {
		file_services_completor_completor_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_completor_completor_proto_rawDescData)
	})
	return file_services_completor_completor_proto_rawDescData
}

var file_services_completor_completor_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_services_completor_completor_proto_goTypes = []any{
	(*CompleteCitizensRequest)(nil),            // 0: services.completor.CompleteCitizensRequest
	(*CompleteCitizensRespoonse)(nil),          // 1: services.completor.CompleteCitizensRespoonse
	(*CompleteJobsRequest)(nil),                // 2: services.completor.CompleteJobsRequest
	(*CompleteJobsResponse)(nil),               // 3: services.completor.CompleteJobsResponse
	(*CompleteDocumentCategoriesRequest)(nil),  // 4: services.completor.CompleteDocumentCategoriesRequest
	(*CompleteDocumentCategoriesResponse)(nil), // 5: services.completor.CompleteDocumentCategoriesResponse
	(*ListLawBooksRequest)(nil),                // 6: services.completor.ListLawBooksRequest
	(*ListLawBooksResponse)(nil),               // 7: services.completor.ListLawBooksResponse
	(*CompleteCitizenLabelsRequest)(nil),       // 8: services.completor.CompleteCitizenLabelsRequest
	(*CompleteCitizenLabelsResponse)(nil),      // 9: services.completor.CompleteCitizenLabelsResponse
	(*users.UserShort)(nil),                    // 10: resources.users.UserShort
	(*users.Job)(nil),                          // 11: resources.users.Job
	(*documents.Category)(nil),                 // 12: resources.documents.Category
	(*laws.LawBook)(nil),                       // 13: resources.laws.LawBook
	(*users.CitizenLabel)(nil),                 // 14: resources.users.CitizenLabel
}
var file_services_completor_completor_proto_depIdxs = []int32{
	10, // 0: services.completor.CompleteCitizensRespoonse.users:type_name -> resources.users.UserShort
	11, // 1: services.completor.CompleteJobsResponse.jobs:type_name -> resources.users.Job
	12, // 2: services.completor.CompleteDocumentCategoriesResponse.categories:type_name -> resources.documents.Category
	13, // 3: services.completor.ListLawBooksResponse.books:type_name -> resources.laws.LawBook
	14, // 4: services.completor.CompleteCitizenLabelsResponse.labels:type_name -> resources.users.CitizenLabel
	0,  // 5: services.completor.CompletorService.CompleteCitizens:input_type -> services.completor.CompleteCitizensRequest
	2,  // 6: services.completor.CompletorService.CompleteJobs:input_type -> services.completor.CompleteJobsRequest
	4,  // 7: services.completor.CompletorService.CompleteDocumentCategories:input_type -> services.completor.CompleteDocumentCategoriesRequest
	6,  // 8: services.completor.CompletorService.ListLawBooks:input_type -> services.completor.ListLawBooksRequest
	8,  // 9: services.completor.CompletorService.CompleteCitizenLabels:input_type -> services.completor.CompleteCitizenLabelsRequest
	1,  // 10: services.completor.CompletorService.CompleteCitizens:output_type -> services.completor.CompleteCitizensRespoonse
	3,  // 11: services.completor.CompletorService.CompleteJobs:output_type -> services.completor.CompleteJobsResponse
	5,  // 12: services.completor.CompletorService.CompleteDocumentCategories:output_type -> services.completor.CompleteDocumentCategoriesResponse
	7,  // 13: services.completor.CompletorService.ListLawBooks:output_type -> services.completor.ListLawBooksResponse
	9,  // 14: services.completor.CompletorService.CompleteCitizenLabels:output_type -> services.completor.CompleteCitizenLabelsResponse
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_services_completor_completor_proto_init() }
func file_services_completor_completor_proto_init() {
	if File_services_completor_completor_proto != nil {
		return
	}
	file_services_completor_completor_proto_msgTypes[0].OneofWrappers = []any{}
	file_services_completor_completor_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_completor_completor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_completor_completor_proto_goTypes,
		DependencyIndexes: file_services_completor_completor_proto_depIdxs,
		MessageInfos:      file_services_completor_completor_proto_msgTypes,
	}.Build()
	File_services_completor_completor_proto = out.File
	file_services_completor_completor_proto_rawDesc = nil
	file_services_completor_completor_proto_goTypes = nil
	file_services_completor_completor_proto_depIdxs = nil
}
