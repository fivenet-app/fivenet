// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: services/wiki/collab.proto

package wiki

import (
	collab "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/collab"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_services_wiki_collab_proto protoreflect.FileDescriptor

const file_services_wiki_collab_proto_rawDesc = "" +
	"\n" +
	"\x1aservices/wiki/collab.proto\x12\rservices.wiki\x1a\x1dresources/collab/collab.proto2_\n" +
	"\rCollabService\x12N\n" +
	"\bJoinRoom\x12\x1e.resources.collab.ClientPacket\x1a\x1e.resources.collab.ServerPacket(\x010\x01BFZDgithub.com/fivenet-app/fivenet/v2025/gen/go/proto/services/wiki;wikib\x06proto3"

var file_services_wiki_collab_proto_goTypes = []any{
	(*collab.ClientPacket)(nil), // 0: resources.collab.ClientPacket
	(*collab.ServerPacket)(nil), // 1: resources.collab.ServerPacket
}
var file_services_wiki_collab_proto_depIdxs = []int32{
	0, // 0: services.wiki.CollabService.JoinRoom:input_type -> resources.collab.ClientPacket
	1, // 1: services.wiki.CollabService.JoinRoom:output_type -> resources.collab.ServerPacket
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_services_wiki_collab_proto_init() }
func file_services_wiki_collab_proto_init() {
	if File_services_wiki_collab_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_services_wiki_collab_proto_rawDesc), len(file_services_wiki_collab_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_wiki_collab_proto_goTypes,
		DependencyIndexes: file_services_wiki_collab_proto_depIdxs,
	}.Build()
	File_services_wiki_collab_proto = out.File
	file_services_wiki_collab_proto_goTypes = nil
	file_services_wiki_collab_proto_depIdxs = nil
}
