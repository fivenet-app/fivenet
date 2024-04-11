// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: services/livemapper/livemap.proto

package livemapper

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	livemap "github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
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

type StreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_livemapper_livemap_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_livemapper_livemap_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRequest.ProtoReflect.Descriptor instead.
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return file_services_livemapper_livemap_proto_rawDescGZIP(), []int{0}
}

type StreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*StreamResponse_Jobs
	//	*StreamResponse_Markers
	//	*StreamResponse_Users
	Data isStreamResponse_Data `protobuf_oneof:"data"`
}

func (x *StreamResponse) Reset() {
	*x = StreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_livemapper_livemap_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResponse) ProtoMessage() {}

func (x *StreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_livemapper_livemap_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamResponse.ProtoReflect.Descriptor instead.
func (*StreamResponse) Descriptor() ([]byte, []int) {
	return file_services_livemapper_livemap_proto_rawDescGZIP(), []int{1}
}

func (m *StreamResponse) GetData() isStreamResponse_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *StreamResponse) GetJobs() *JobsList {
	if x, ok := x.GetData().(*StreamResponse_Jobs); ok {
		return x.Jobs
	}
	return nil
}

func (x *StreamResponse) GetMarkers() *MarkerMarkersUpdates {
	if x, ok := x.GetData().(*StreamResponse_Markers); ok {
		return x.Markers
	}
	return nil
}

func (x *StreamResponse) GetUsers() *UserMarkersUpdates {
	if x, ok := x.GetData().(*StreamResponse_Users); ok {
		return x.Users
	}
	return nil
}

type isStreamResponse_Data interface {
	isStreamResponse_Data()
}

type StreamResponse_Jobs struct {
	Jobs *JobsList `protobuf:"bytes,1,opt,name=jobs,proto3,oneof"`
}

type StreamResponse_Markers struct {
	Markers *MarkerMarkersUpdates `protobuf:"bytes,2,opt,name=markers,proto3,oneof"`
}

type StreamResponse_Users struct {
	Users *UserMarkersUpdates `protobuf:"bytes,3,opt,name=users,proto3,oneof"`
}

func (*StreamResponse_Jobs) isStreamResponse_Data() {}

func (*StreamResponse_Markers) isStreamResponse_Data() {}

func (*StreamResponse_Users) isStreamResponse_Data() {}

type JobsList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users   []*users.Job `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	Markers []*users.Job `protobuf:"bytes,2,rep,name=markers,proto3" json:"markers,omitempty"`
}

func (x *JobsList) Reset() {
	*x = JobsList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_livemapper_livemap_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobsList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobsList) ProtoMessage() {}

func (x *JobsList) ProtoReflect() protoreflect.Message {
	mi := &file_services_livemapper_livemap_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobsList.ProtoReflect.Descriptor instead.
func (*JobsList) Descriptor() ([]byte, []int) {
	return file_services_livemapper_livemap_proto_rawDescGZIP(), []int{2}
}

func (x *JobsList) GetUsers() []*users.Job {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *JobsList) GetMarkers() []*users.Job {
	if x != nil {
		return x.Markers
	}
	return nil
}

type MarkerMarkersUpdates struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Markers []*livemap.MarkerMarker `protobuf:"bytes,1,rep,name=markers,proto3" json:"markers,omitempty"`
}

func (x *MarkerMarkersUpdates) Reset() {
	*x = MarkerMarkersUpdates{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_livemapper_livemap_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarkerMarkersUpdates) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkerMarkersUpdates) ProtoMessage() {}

func (x *MarkerMarkersUpdates) ProtoReflect() protoreflect.Message {
	mi := &file_services_livemapper_livemap_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkerMarkersUpdates.ProtoReflect.Descriptor instead.
func (*MarkerMarkersUpdates) Descriptor() ([]byte, []int) {
	return file_services_livemapper_livemap_proto_rawDescGZIP(), []int{3}
}

func (x *MarkerMarkersUpdates) GetMarkers() []*livemap.MarkerMarker {
	if x != nil {
		return x.Markers
	}
	return nil
}

type UserMarkersUpdates struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*livemap.UserMarker `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	Part  int32                 `protobuf:"varint,2,opt,name=part,proto3" json:"part,omitempty"`
}

func (x *UserMarkersUpdates) Reset() {
	*x = UserMarkersUpdates{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_livemapper_livemap_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserMarkersUpdates) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserMarkersUpdates) ProtoMessage() {}

func (x *UserMarkersUpdates) ProtoReflect() protoreflect.Message {
	mi := &file_services_livemapper_livemap_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserMarkersUpdates.ProtoReflect.Descriptor instead.
func (*UserMarkersUpdates) Descriptor() ([]byte, []int) {
	return file_services_livemapper_livemap_proto_rawDescGZIP(), []int{4}
}

func (x *UserMarkersUpdates) GetUsers() []*livemap.UserMarker {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *UserMarkersUpdates) GetPart() int32 {
	if x != nil {
		return x.Part
	}
	return 0
}

type CreateOrUpdateMarkerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Marker *livemap.MarkerMarker `protobuf:"bytes,1,opt,name=marker,proto3" json:"marker,omitempty"`
}

func (x *CreateOrUpdateMarkerRequest) Reset() {
	*x = CreateOrUpdateMarkerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_livemapper_livemap_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateOrUpdateMarkerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrUpdateMarkerRequest) ProtoMessage() {}

func (x *CreateOrUpdateMarkerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_livemapper_livemap_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrUpdateMarkerRequest.ProtoReflect.Descriptor instead.
func (*CreateOrUpdateMarkerRequest) Descriptor() ([]byte, []int) {
	return file_services_livemapper_livemap_proto_rawDescGZIP(), []int{5}
}

func (x *CreateOrUpdateMarkerRequest) GetMarker() *livemap.MarkerMarker {
	if x != nil {
		return x.Marker
	}
	return nil
}

type CreateOrUpdateMarkerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Marker *livemap.MarkerMarker `protobuf:"bytes,1,opt,name=marker,proto3" json:"marker,omitempty"`
}

func (x *CreateOrUpdateMarkerResponse) Reset() {
	*x = CreateOrUpdateMarkerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_livemapper_livemap_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateOrUpdateMarkerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrUpdateMarkerResponse) ProtoMessage() {}

func (x *CreateOrUpdateMarkerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_livemapper_livemap_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrUpdateMarkerResponse.ProtoReflect.Descriptor instead.
func (*CreateOrUpdateMarkerResponse) Descriptor() ([]byte, []int) {
	return file_services_livemapper_livemap_proto_rawDescGZIP(), []int{6}
}

func (x *CreateOrUpdateMarkerResponse) GetMarker() *livemap.MarkerMarker {
	if x != nil {
		return x.Marker
	}
	return nil
}

type DeleteMarkerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteMarkerRequest) Reset() {
	*x = DeleteMarkerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_livemapper_livemap_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteMarkerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMarkerRequest) ProtoMessage() {}

func (x *DeleteMarkerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_livemapper_livemap_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMarkerRequest.ProtoReflect.Descriptor instead.
func (*DeleteMarkerRequest) Descriptor() ([]byte, []int) {
	return file_services_livemapper_livemap_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteMarkerRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DeleteMarkerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteMarkerResponse) Reset() {
	*x = DeleteMarkerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_livemapper_livemap_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteMarkerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMarkerResponse) ProtoMessage() {}

func (x *DeleteMarkerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_livemapper_livemap_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMarkerResponse.ProtoReflect.Descriptor instead.
func (*DeleteMarkerResponse) Descriptor() ([]byte, []int) {
	return file_services_livemapper_livemap_proto_rawDescGZIP(), []int{8}
}

var File_services_livemapper_livemap_proto protoreflect.FileDescriptor

var file_services_livemapper_livemap_proto_rawDesc = []byte{
	0x0a, 0x21, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6d,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x13, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69,
	0x76, 0x65, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x1a, 0x1f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x2f, 0x6c, 0x69, 0x76, 0x65,
	0x6d, 0x61, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0f,
	0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0xd5, 0x01, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x33, 0x0a, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65,
	0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x4a, 0x6f, 0x62, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x48,
	0x00, 0x52, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x12, 0x45, 0x0a, 0x07, 0x6d, 0x61, 0x72, 0x6b, 0x65,
	0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x4d,
	0x61, 0x72, 0x6b, 0x65, 0x72, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x73, 0x48, 0x00, 0x52, 0x07, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x12, 0x3f,
	0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x48, 0x00, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x42,
	0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x66, 0x0a, 0x08, 0x4a, 0x6f, 0x62, 0x73, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12,
	0x2e, 0x0a, 0x07, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x4a, 0x6f, 0x62, 0x52, 0x07, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x22,
	0x51, 0x0a, 0x14, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x73,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x12, 0x39, 0x0a, 0x07, 0x6d, 0x61, 0x72, 0x6b, 0x65,
	0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x2e, 0x4d, 0x61, 0x72,
	0x6b, 0x65, 0x72, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x07, 0x6d, 0x61, 0x72, 0x6b, 0x65,
	0x72, 0x73, 0x22, 0x5d, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72,
	0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x12, 0x33, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x72,
	0x74, 0x22, 0x60, 0x0a, 0x1b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x41, 0x0a, 0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76,
	0x65, 0x6d, 0x61, 0x70, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x4d, 0x61, 0x72, 0x6b, 0x65,
	0x72, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x6d, 0x61, 0x72,
	0x6b, 0x65, 0x72, 0x22, 0x57, 0x0a, 0x1c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e,
	0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x4d, 0x61,
	0x72, 0x6b, 0x65, 0x72, 0x52, 0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x22, 0x29, 0x0a, 0x13,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42,
	0x02, 0x30, 0x01, 0x52, 0x02, 0x69, 0x64, 0x22, 0x16, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32,
	0xca, 0x02, 0x0a, 0x11, 0x4c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x53, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12,
	0x22, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65, 0x6d,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6c,
	0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x7b, 0x0a, 0x14, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b,
	0x65, 0x72, 0x12, 0x30, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69,
	0x76, 0x65, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f,
	0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4f, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x63, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x12, 0x28, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x29, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76,
	0x65, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x61,
	0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x48, 0x5a, 0x46,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x6c, 0x65, 0x78,
	0x72, 0x74, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x3b, 0x6c, 0x69, 0x76, 0x65,
	0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_livemapper_livemap_proto_rawDescOnce sync.Once
	file_services_livemapper_livemap_proto_rawDescData = file_services_livemapper_livemap_proto_rawDesc
)

func file_services_livemapper_livemap_proto_rawDescGZIP() []byte {
	file_services_livemapper_livemap_proto_rawDescOnce.Do(func() {
		file_services_livemapper_livemap_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_livemapper_livemap_proto_rawDescData)
	})
	return file_services_livemapper_livemap_proto_rawDescData
}

var file_services_livemapper_livemap_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_services_livemapper_livemap_proto_goTypes = []interface{}{
	(*StreamRequest)(nil),                // 0: services.livemapper.StreamRequest
	(*StreamResponse)(nil),               // 1: services.livemapper.StreamResponse
	(*JobsList)(nil),                     // 2: services.livemapper.JobsList
	(*MarkerMarkersUpdates)(nil),         // 3: services.livemapper.MarkerMarkersUpdates
	(*UserMarkersUpdates)(nil),           // 4: services.livemapper.UserMarkersUpdates
	(*CreateOrUpdateMarkerRequest)(nil),  // 5: services.livemapper.CreateOrUpdateMarkerRequest
	(*CreateOrUpdateMarkerResponse)(nil), // 6: services.livemapper.CreateOrUpdateMarkerResponse
	(*DeleteMarkerRequest)(nil),          // 7: services.livemapper.DeleteMarkerRequest
	(*DeleteMarkerResponse)(nil),         // 8: services.livemapper.DeleteMarkerResponse
	(*users.Job)(nil),                    // 9: resources.users.Job
	(*livemap.MarkerMarker)(nil),         // 10: resources.livemap.MarkerMarker
	(*livemap.UserMarker)(nil),           // 11: resources.livemap.UserMarker
}
var file_services_livemapper_livemap_proto_depIdxs = []int32{
	2,  // 0: services.livemapper.StreamResponse.jobs:type_name -> services.livemapper.JobsList
	3,  // 1: services.livemapper.StreamResponse.markers:type_name -> services.livemapper.MarkerMarkersUpdates
	4,  // 2: services.livemapper.StreamResponse.users:type_name -> services.livemapper.UserMarkersUpdates
	9,  // 3: services.livemapper.JobsList.users:type_name -> resources.users.Job
	9,  // 4: services.livemapper.JobsList.markers:type_name -> resources.users.Job
	10, // 5: services.livemapper.MarkerMarkersUpdates.markers:type_name -> resources.livemap.MarkerMarker
	11, // 6: services.livemapper.UserMarkersUpdates.users:type_name -> resources.livemap.UserMarker
	10, // 7: services.livemapper.CreateOrUpdateMarkerRequest.marker:type_name -> resources.livemap.MarkerMarker
	10, // 8: services.livemapper.CreateOrUpdateMarkerResponse.marker:type_name -> resources.livemap.MarkerMarker
	0,  // 9: services.livemapper.LivemapperService.Stream:input_type -> services.livemapper.StreamRequest
	5,  // 10: services.livemapper.LivemapperService.CreateOrUpdateMarker:input_type -> services.livemapper.CreateOrUpdateMarkerRequest
	7,  // 11: services.livemapper.LivemapperService.DeleteMarker:input_type -> services.livemapper.DeleteMarkerRequest
	1,  // 12: services.livemapper.LivemapperService.Stream:output_type -> services.livemapper.StreamResponse
	6,  // 13: services.livemapper.LivemapperService.CreateOrUpdateMarker:output_type -> services.livemapper.CreateOrUpdateMarkerResponse
	8,  // 14: services.livemapper.LivemapperService.DeleteMarker:output_type -> services.livemapper.DeleteMarkerResponse
	12, // [12:15] is the sub-list for method output_type
	9,  // [9:12] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_services_livemapper_livemap_proto_init() }
func file_services_livemapper_livemap_proto_init() {
	if File_services_livemapper_livemap_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_livemapper_livemap_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_livemapper_livemap_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_livemapper_livemap_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobsList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_livemapper_livemap_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarkerMarkersUpdates); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_livemapper_livemap_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserMarkersUpdates); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_livemapper_livemap_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateOrUpdateMarkerRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_livemapper_livemap_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateOrUpdateMarkerResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_livemapper_livemap_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteMarkerRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_livemapper_livemap_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteMarkerResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_services_livemapper_livemap_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*StreamResponse_Jobs)(nil),
		(*StreamResponse_Markers)(nil),
		(*StreamResponse_Users)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_livemapper_livemap_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_livemapper_livemap_proto_goTypes,
		DependencyIndexes: file_services_livemapper_livemap_proto_depIdxs,
		MessageInfos:      file_services_livemapper_livemap_proto_msgTypes,
	}.Build()
	File_services_livemapper_livemap_proto = out.File
	file_services_livemapper_livemap_proto_rawDesc = nil
	file_services_livemapper_livemap_proto_goTypes = nil
	file_services_livemapper_livemap_proto_depIdxs = nil
}