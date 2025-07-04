// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: resources/common/grpcws/grpcws.proto

package grpcws

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

type GrpcFrame struct {
	state    protoimpl.MessageState `protogen:"open.v1"`
	StreamId uint32                 `protobuf:"varint,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
	// Types that are valid to be assigned to Payload:
	//
	//	*GrpcFrame_Ping
	//	*GrpcFrame_Header
	//	*GrpcFrame_Body
	//	*GrpcFrame_Complete
	//	*GrpcFrame_Failure
	//	*GrpcFrame_Cancel
	Payload       isGrpcFrame_Payload `protobuf_oneof:"payload"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GrpcFrame) Reset() {
	*x = GrpcFrame{}
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GrpcFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcFrame) ProtoMessage() {}

func (x *GrpcFrame) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcFrame.ProtoReflect.Descriptor instead.
func (*GrpcFrame) Descriptor() ([]byte, []int) {
	return file_resources_common_grpcws_grpcws_proto_rawDescGZIP(), []int{0}
}

func (x *GrpcFrame) GetStreamId() uint32 {
	if x != nil {
		return x.StreamId
	}
	return 0
}

func (x *GrpcFrame) GetPayload() isGrpcFrame_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *GrpcFrame) GetPing() *Ping {
	if x != nil {
		if x, ok := x.Payload.(*GrpcFrame_Ping); ok {
			return x.Ping
		}
	}
	return nil
}

func (x *GrpcFrame) GetHeader() *Header {
	if x != nil {
		if x, ok := x.Payload.(*GrpcFrame_Header); ok {
			return x.Header
		}
	}
	return nil
}

func (x *GrpcFrame) GetBody() *Body {
	if x != nil {
		if x, ok := x.Payload.(*GrpcFrame_Body); ok {
			return x.Body
		}
	}
	return nil
}

func (x *GrpcFrame) GetComplete() *Complete {
	if x != nil {
		if x, ok := x.Payload.(*GrpcFrame_Complete); ok {
			return x.Complete
		}
	}
	return nil
}

func (x *GrpcFrame) GetFailure() *Failure {
	if x != nil {
		if x, ok := x.Payload.(*GrpcFrame_Failure); ok {
			return x.Failure
		}
	}
	return nil
}

func (x *GrpcFrame) GetCancel() *Cancel {
	if x != nil {
		if x, ok := x.Payload.(*GrpcFrame_Cancel); ok {
			return x.Cancel
		}
	}
	return nil
}

type isGrpcFrame_Payload interface {
	isGrpcFrame_Payload()
}

type GrpcFrame_Ping struct {
	Ping *Ping `protobuf:"bytes,3,opt,name=ping,proto3,oneof"`
}

type GrpcFrame_Header struct {
	Header *Header `protobuf:"bytes,4,opt,name=header,proto3,oneof"`
}

type GrpcFrame_Body struct {
	Body *Body `protobuf:"bytes,5,opt,name=body,proto3,oneof"`
}

type GrpcFrame_Complete struct {
	Complete *Complete `protobuf:"bytes,6,opt,name=complete,proto3,oneof"`
}

type GrpcFrame_Failure struct {
	Failure *Failure `protobuf:"bytes,7,opt,name=failure,proto3,oneof"`
}

type GrpcFrame_Cancel struct {
	Cancel *Cancel `protobuf:"bytes,8,opt,name=cancel,proto3,oneof"`
}

func (*GrpcFrame_Ping) isGrpcFrame_Payload() {}

func (*GrpcFrame_Header) isGrpcFrame_Payload() {}

func (*GrpcFrame_Body) isGrpcFrame_Payload() {}

func (*GrpcFrame_Complete) isGrpcFrame_Payload() {}

func (*GrpcFrame_Failure) isGrpcFrame_Payload() {}

func (*GrpcFrame_Cancel) isGrpcFrame_Payload() {}

type Ping struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Pong          bool                   `protobuf:"varint,1,opt,name=pong,proto3" json:"pong,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Ping) Reset() {
	*x = Ping{}
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_resources_common_grpcws_grpcws_proto_rawDescGZIP(), []int{1}
}

func (x *Ping) GetPong() bool {
	if x != nil {
		return x.Pong
	}
	return false
}

type Header struct {
	state         protoimpl.MessageState  `protogen:"open.v1"`
	Operation     string                  `protobuf:"bytes,1,opt,name=operation,proto3" json:"operation,omitempty"`
	Headers       map[string]*HeaderValue `protobuf:"bytes,2,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Status        int32                   `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Header) Reset() {
	*x = Header{}
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_resources_common_grpcws_grpcws_proto_rawDescGZIP(), []int{2}
}

func (x *Header) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

func (x *Header) GetHeaders() map[string]*HeaderValue {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *Header) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type HeaderValue struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Value         []string               `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HeaderValue) Reset() {
	*x = HeaderValue{}
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeaderValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeaderValue) ProtoMessage() {}

func (x *HeaderValue) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeaderValue.ProtoReflect.Descriptor instead.
func (*HeaderValue) Descriptor() ([]byte, []int) {
	return file_resources_common_grpcws_grpcws_proto_rawDescGZIP(), []int{3}
}

func (x *HeaderValue) GetValue() []string {
	if x != nil {
		return x.Value
	}
	return nil
}

type Body struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []byte                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Complete      bool                   `protobuf:"varint,2,opt,name=complete,proto3" json:"complete,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Body) Reset() {
	*x = Body{}
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Body) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Body) ProtoMessage() {}

func (x *Body) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Body.ProtoReflect.Descriptor instead.
func (*Body) Descriptor() ([]byte, []int) {
	return file_resources_common_grpcws_grpcws_proto_rawDescGZIP(), []int{4}
}

func (x *Body) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Body) GetComplete() bool {
	if x != nil {
		return x.Complete
	}
	return false
}

type Complete struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Complete) Reset() {
	*x = Complete{}
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Complete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Complete) ProtoMessage() {}

func (x *Complete) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Complete.ProtoReflect.Descriptor instead.
func (*Complete) Descriptor() ([]byte, []int) {
	return file_resources_common_grpcws_grpcws_proto_rawDescGZIP(), []int{5}
}

type Failure struct {
	state         protoimpl.MessageState  `protogen:"open.v1"`
	ErrorMessage  string                  `protobuf:"bytes,1,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	ErrorStatus   string                  `protobuf:"bytes,2,opt,name=error_status,json=errorStatus,proto3" json:"error_status,omitempty"`
	Headers       map[string]*HeaderValue `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Failure) Reset() {
	*x = Failure{}
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Failure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Failure) ProtoMessage() {}

func (x *Failure) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Failure.ProtoReflect.Descriptor instead.
func (*Failure) Descriptor() ([]byte, []int) {
	return file_resources_common_grpcws_grpcws_proto_rawDescGZIP(), []int{6}
}

func (x *Failure) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *Failure) GetErrorStatus() string {
	if x != nil {
		return x.ErrorStatus
	}
	return ""
}

func (x *Failure) GetHeaders() map[string]*HeaderValue {
	if x != nil {
		return x.Headers
	}
	return nil
}

type Cancel struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cancel) Reset() {
	*x = Cancel{}
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cancel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cancel) ProtoMessage() {}

func (x *Cancel) ProtoReflect() protoreflect.Message {
	mi := &file_resources_common_grpcws_grpcws_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cancel.ProtoReflect.Descriptor instead.
func (*Cancel) Descriptor() ([]byte, []int) {
	return file_resources_common_grpcws_grpcws_proto_rawDescGZIP(), []int{7}
}

var File_resources_common_grpcws_grpcws_proto protoreflect.FileDescriptor

const file_resources_common_grpcws_grpcws_proto_rawDesc = "" +
	"\n" +
	"$resources/common/grpcws/grpcws.proto\x12\x17resources.common.grpcws\"\x99\x03\n" +
	"\tGrpcFrame\x12\x1b\n" +
	"\tstream_id\x18\x01 \x01(\rR\bstreamId\x123\n" +
	"\x04ping\x18\x03 \x01(\v2\x1d.resources.common.grpcws.PingH\x00R\x04ping\x129\n" +
	"\x06header\x18\x04 \x01(\v2\x1f.resources.common.grpcws.HeaderH\x00R\x06header\x123\n" +
	"\x04body\x18\x05 \x01(\v2\x1d.resources.common.grpcws.BodyH\x00R\x04body\x12?\n" +
	"\bcomplete\x18\x06 \x01(\v2!.resources.common.grpcws.CompleteH\x00R\bcomplete\x12<\n" +
	"\afailure\x18\a \x01(\v2 .resources.common.grpcws.FailureH\x00R\afailure\x129\n" +
	"\x06cancel\x18\b \x01(\v2\x1f.resources.common.grpcws.CancelH\x00R\x06cancelB\x10\n" +
	"\apayload\x12\x05\xbaH\x02\b\x01\"\x1a\n" +
	"\x04Ping\x12\x12\n" +
	"\x04pong\x18\x01 \x01(\bR\x04pong\"\xe8\x01\n" +
	"\x06Header\x12\x1c\n" +
	"\toperation\x18\x01 \x01(\tR\toperation\x12F\n" +
	"\aheaders\x18\x02 \x03(\v2,.resources.common.grpcws.Header.HeadersEntryR\aheaders\x12\x16\n" +
	"\x06status\x18\x03 \x01(\x05R\x06status\x1a`\n" +
	"\fHeadersEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12:\n" +
	"\x05value\x18\x02 \x01(\v2$.resources.common.grpcws.HeaderValueR\x05value:\x028\x01\"#\n" +
	"\vHeaderValue\x12\x14\n" +
	"\x05value\x18\x01 \x03(\tR\x05value\"6\n" +
	"\x04Body\x12\x12\n" +
	"\x04data\x18\x01 \x01(\fR\x04data\x12\x1a\n" +
	"\bcomplete\x18\x02 \x01(\bR\bcomplete\"\n" +
	"\n" +
	"\bComplete\"\xfc\x01\n" +
	"\aFailure\x12#\n" +
	"\rerror_message\x18\x01 \x01(\tR\ferrorMessage\x12!\n" +
	"\ferror_status\x18\x02 \x01(\tR\verrorStatus\x12G\n" +
	"\aheaders\x18\x03 \x03(\v2-.resources.common.grpcws.Failure.HeadersEntryR\aheaders\x1a`\n" +
	"\fHeadersEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12:\n" +
	"\x05value\x18\x02 \x01(\v2$.resources.common.grpcws.HeaderValueR\x05value:\x028\x01\"\b\n" +
	"\x06CancelBRZPgithub.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/grpcws;grpcwsb\x06proto3"

var (
	file_resources_common_grpcws_grpcws_proto_rawDescOnce sync.Once
	file_resources_common_grpcws_grpcws_proto_rawDescData []byte
)

func file_resources_common_grpcws_grpcws_proto_rawDescGZIP() []byte {
	file_resources_common_grpcws_grpcws_proto_rawDescOnce.Do(func() {
		file_resources_common_grpcws_grpcws_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_common_grpcws_grpcws_proto_rawDesc), len(file_resources_common_grpcws_grpcws_proto_rawDesc)))
	})
	return file_resources_common_grpcws_grpcws_proto_rawDescData
}

var file_resources_common_grpcws_grpcws_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_resources_common_grpcws_grpcws_proto_goTypes = []any{
	(*GrpcFrame)(nil),   // 0: resources.common.grpcws.GrpcFrame
	(*Ping)(nil),        // 1: resources.common.grpcws.Ping
	(*Header)(nil),      // 2: resources.common.grpcws.Header
	(*HeaderValue)(nil), // 3: resources.common.grpcws.HeaderValue
	(*Body)(nil),        // 4: resources.common.grpcws.Body
	(*Complete)(nil),    // 5: resources.common.grpcws.Complete
	(*Failure)(nil),     // 6: resources.common.grpcws.Failure
	(*Cancel)(nil),      // 7: resources.common.grpcws.Cancel
	nil,                 // 8: resources.common.grpcws.Header.HeadersEntry
	nil,                 // 9: resources.common.grpcws.Failure.HeadersEntry
}
var file_resources_common_grpcws_grpcws_proto_depIdxs = []int32{
	1,  // 0: resources.common.grpcws.GrpcFrame.ping:type_name -> resources.common.grpcws.Ping
	2,  // 1: resources.common.grpcws.GrpcFrame.header:type_name -> resources.common.grpcws.Header
	4,  // 2: resources.common.grpcws.GrpcFrame.body:type_name -> resources.common.grpcws.Body
	5,  // 3: resources.common.grpcws.GrpcFrame.complete:type_name -> resources.common.grpcws.Complete
	6,  // 4: resources.common.grpcws.GrpcFrame.failure:type_name -> resources.common.grpcws.Failure
	7,  // 5: resources.common.grpcws.GrpcFrame.cancel:type_name -> resources.common.grpcws.Cancel
	8,  // 6: resources.common.grpcws.Header.headers:type_name -> resources.common.grpcws.Header.HeadersEntry
	9,  // 7: resources.common.grpcws.Failure.headers:type_name -> resources.common.grpcws.Failure.HeadersEntry
	3,  // 8: resources.common.grpcws.Header.HeadersEntry.value:type_name -> resources.common.grpcws.HeaderValue
	3,  // 9: resources.common.grpcws.Failure.HeadersEntry.value:type_name -> resources.common.grpcws.HeaderValue
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_resources_common_grpcws_grpcws_proto_init() }
func file_resources_common_grpcws_grpcws_proto_init() {
	if File_resources_common_grpcws_grpcws_proto != nil {
		return
	}
	file_resources_common_grpcws_grpcws_proto_msgTypes[0].OneofWrappers = []any{
		(*GrpcFrame_Ping)(nil),
		(*GrpcFrame_Header)(nil),
		(*GrpcFrame_Body)(nil),
		(*GrpcFrame_Complete)(nil),
		(*GrpcFrame_Failure)(nil),
		(*GrpcFrame_Cancel)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_common_grpcws_grpcws_proto_rawDesc), len(file_resources_common_grpcws_grpcws_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_common_grpcws_grpcws_proto_goTypes,
		DependencyIndexes: file_resources_common_grpcws_grpcws_proto_depIdxs,
		MessageInfos:      file_resources_common_grpcws_grpcws_proto_msgTypes,
	}.Build()
	File_resources_common_grpcws_grpcws_proto = out.File
	file_resources_common_grpcws_grpcws_proto_goTypes = nil
	file_resources_common_grpcws_grpcws_proto_depIdxs = nil
}
