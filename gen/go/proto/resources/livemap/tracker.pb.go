// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.20.3
// source: resources/livemap/tracker.proto

package livemap

import (
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

type UsersUpdateEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Added   []*UserMarker `protobuf:"bytes,1,rep,name=added,proto3" json:"added,omitempty"`
	Removed []*UserMarker `protobuf:"bytes,2,rep,name=removed,proto3" json:"removed,omitempty"`
}

func (x *UsersUpdateEvent) Reset() {
	*x = UsersUpdateEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resources_livemap_tracker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsersUpdateEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsersUpdateEvent) ProtoMessage() {}

func (x *UsersUpdateEvent) ProtoReflect() protoreflect.Message {
	mi := &file_resources_livemap_tracker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsersUpdateEvent.ProtoReflect.Descriptor instead.
func (*UsersUpdateEvent) Descriptor() ([]byte, []int) {
	return file_resources_livemap_tracker_proto_rawDescGZIP(), []int{0}
}

func (x *UsersUpdateEvent) GetAdded() []*UserMarker {
	if x != nil {
		return x.Added
	}
	return nil
}

func (x *UsersUpdateEvent) GetRemoved() []*UserMarker {
	if x != nil {
		return x.Removed
	}
	return nil
}

var File_resources_livemap_tracker_proto protoreflect.FileDescriptor

var file_resources_livemap_tracker_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6c, 0x69, 0x76, 0x65,
	0x6d, 0x61, 0x70, 0x2f, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x11, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76,
	0x65, 0x6d, 0x61, 0x70, 0x1a, 0x1f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f,
	0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x80, 0x01, 0x0a, 0x10, 0x55, 0x73, 0x65, 0x72, 0x73, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x33, 0x0a, 0x05, 0x61, 0x64,
	0x64, 0x65, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x05, 0x61, 0x64, 0x64, 0x65, 0x64, 0x12,
	0x37, 0x0a, 0x07, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x6c, 0x69, 0x76,
	0x65, 0x6d, 0x61, 0x70, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x72, 0x52,
	0x07, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x42, 0x47, 0x5a, 0x45, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61,
	0x70, 0x70, 0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61, 0x70, 0x3b, 0x6c, 0x69, 0x76, 0x65, 0x6d, 0x61,
	0x70, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resources_livemap_tracker_proto_rawDescOnce sync.Once
	file_resources_livemap_tracker_proto_rawDescData = file_resources_livemap_tracker_proto_rawDesc
)

func file_resources_livemap_tracker_proto_rawDescGZIP() []byte {
	file_resources_livemap_tracker_proto_rawDescOnce.Do(func() {
		file_resources_livemap_tracker_proto_rawDescData = protoimpl.X.CompressGZIP(file_resources_livemap_tracker_proto_rawDescData)
	})
	return file_resources_livemap_tracker_proto_rawDescData
}

var file_resources_livemap_tracker_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_livemap_tracker_proto_goTypes = []interface{}{
	(*UsersUpdateEvent)(nil), // 0: resources.livemap.UsersUpdateEvent
	(*UserMarker)(nil),       // 1: resources.livemap.UserMarker
}
var file_resources_livemap_tracker_proto_depIdxs = []int32{
	1, // 0: resources.livemap.UsersUpdateEvent.added:type_name -> resources.livemap.UserMarker
	1, // 1: resources.livemap.UsersUpdateEvent.removed:type_name -> resources.livemap.UserMarker
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_resources_livemap_tracker_proto_init() }
func file_resources_livemap_tracker_proto_init() {
	if File_resources_livemap_tracker_proto != nil {
		return
	}
	file_resources_livemap_livemap_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_resources_livemap_tracker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsersUpdateEvent); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resources_livemap_tracker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_livemap_tracker_proto_goTypes,
		DependencyIndexes: file_resources_livemap_tracker_proto_depIdxs,
		MessageInfos:      file_resources_livemap_tracker_proto_msgTypes,
	}.Build()
	File_resources_livemap_tracker_proto = out.File
	file_resources_livemap_tracker_proto_rawDesc = nil
	file_resources_livemap_tracker_proto_goTypes = nil
	file_resources_livemap_tracker_proto_depIdxs = nil
}
