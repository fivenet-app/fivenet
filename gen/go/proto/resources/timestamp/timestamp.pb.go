// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: resources/timestamp/timestamp.proto

package timestamp

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

// Timestamp for storage messages.  We've defined a new local type wrapper
// of google.protobuf.Timestamp so we can implement sql.Scanner and sql.Valuer
// interfaces.  See:
// https://golang.org/pkg/database/sql/#Scanner
// https://golang.org/pkg/database/sql/driver/#Valuer
type Timestamp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Timestamp) Reset() {
	*x = Timestamp{}
	mi := &file_resources_timestamp_timestamp_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Timestamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Timestamp) ProtoMessage() {}

func (x *Timestamp) ProtoReflect() protoreflect.Message {
	mi := &file_resources_timestamp_timestamp_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Timestamp.ProtoReflect.Descriptor instead.
func (*Timestamp) Descriptor() ([]byte, []int) {
	return file_resources_timestamp_timestamp_proto_rawDescGZIP(), []int{0}
}

func (x *Timestamp) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

var File_resources_timestamp_timestamp_proto protoreflect.FileDescriptor

const file_resources_timestamp_timestamp_proto_rawDesc = "" +
	"\n" +
	"#resources/timestamp/timestamp.proto\x12\x13resources.timestamp\x1a\x1fgoogle/protobuf/timestamp.proto\"E\n" +
	"\tTimestamp\x128\n" +
	"\ttimestamp\x18\x01 \x01(\v2\x1a.google.protobuf.TimestampR\ttimestampBKZIgithub.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp;timestampb\x06proto3"

var (
	file_resources_timestamp_timestamp_proto_rawDescOnce sync.Once
	file_resources_timestamp_timestamp_proto_rawDescData []byte
)

func file_resources_timestamp_timestamp_proto_rawDescGZIP() []byte {
	file_resources_timestamp_timestamp_proto_rawDescOnce.Do(func() {
		file_resources_timestamp_timestamp_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_timestamp_timestamp_proto_rawDesc), len(file_resources_timestamp_timestamp_proto_rawDesc)))
	})
	return file_resources_timestamp_timestamp_proto_rawDescData
}

var file_resources_timestamp_timestamp_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resources_timestamp_timestamp_proto_goTypes = []any{
	(*Timestamp)(nil),             // 0: resources.timestamp.Timestamp
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_resources_timestamp_timestamp_proto_depIdxs = []int32{
	1, // 0: resources.timestamp.Timestamp.timestamp:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_resources_timestamp_timestamp_proto_init() }
func file_resources_timestamp_timestamp_proto_init() {
	if File_resources_timestamp_timestamp_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_timestamp_timestamp_proto_rawDesc), len(file_resources_timestamp_timestamp_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_timestamp_timestamp_proto_goTypes,
		DependencyIndexes: file_resources_timestamp_timestamp_proto_depIdxs,
		MessageInfos:      file_resources_timestamp_timestamp_proto_msgTypes,
	}.Build()
	File_resources_timestamp_timestamp_proto = out.File
	file_resources_timestamp_timestamp_proto_goTypes = nil
	file_resources_timestamp_timestamp_proto_depIdxs = nil
}
