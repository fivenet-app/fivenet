// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v3.20.3
// source: resources/accounts/oauth2.proto

package accounts

import (
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

type OAuth2Account struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccountId     uint64                 `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	CreatedAt     *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	ProviderName  string                 `protobuf:"bytes,3,opt,name=provider_name,json=providerName,proto3" json:"provider_name,omitempty"`
	Provider      *OAuth2Provider        `protobuf:"bytes,4,opt,name=provider,proto3" json:"provider,omitempty"`
	ExternalId    string                 `protobuf:"bytes,5,opt,name=external_id,json=externalId,proto3" json:"external_id,omitempty"`
	Username      string                 `protobuf:"bytes,6,opt,name=username,proto3" json:"username,omitempty"`
	Avatar        string                 `protobuf:"bytes,7,opt,name=avatar,proto3" json:"avatar,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OAuth2Account) Reset() {
	*x = OAuth2Account{}
	mi := &file_resources_accounts_oauth2_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OAuth2Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OAuth2Account) ProtoMessage() {}

func (x *OAuth2Account) ProtoReflect() protoreflect.Message {
	mi := &file_resources_accounts_oauth2_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OAuth2Account.ProtoReflect.Descriptor instead.
func (*OAuth2Account) Descriptor() ([]byte, []int) {
	return file_resources_accounts_oauth2_proto_rawDescGZIP(), []int{0}
}

func (x *OAuth2Account) GetAccountId() uint64 {
	if x != nil {
		return x.AccountId
	}
	return 0
}

func (x *OAuth2Account) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *OAuth2Account) GetProviderName() string {
	if x != nil {
		return x.ProviderName
	}
	return ""
}

func (x *OAuth2Account) GetProvider() *OAuth2Provider {
	if x != nil {
		return x.Provider
	}
	return nil
}

func (x *OAuth2Account) GetExternalId() string {
	if x != nil {
		return x.ExternalId
	}
	return ""
}

func (x *OAuth2Account) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *OAuth2Account) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

type OAuth2Provider struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Label         string                 `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	Homepage      string                 `protobuf:"bytes,3,opt,name=homepage,proto3" json:"homepage,omitempty"`
	Icon          *string                `protobuf:"bytes,4,opt,name=icon,proto3,oneof" json:"icon,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OAuth2Provider) Reset() {
	*x = OAuth2Provider{}
	mi := &file_resources_accounts_oauth2_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OAuth2Provider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OAuth2Provider) ProtoMessage() {}

func (x *OAuth2Provider) ProtoReflect() protoreflect.Message {
	mi := &file_resources_accounts_oauth2_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OAuth2Provider.ProtoReflect.Descriptor instead.
func (*OAuth2Provider) Descriptor() ([]byte, []int) {
	return file_resources_accounts_oauth2_proto_rawDescGZIP(), []int{1}
}

func (x *OAuth2Provider) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *OAuth2Provider) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *OAuth2Provider) GetHomepage() string {
	if x != nil {
		return x.Homepage
	}
	return ""
}

func (x *OAuth2Provider) GetIcon() string {
	if x != nil && x.Icon != nil {
		return *x.Icon
	}
	return ""
}

var File_resources_accounts_oauth2_proto protoreflect.FileDescriptor

var file_resources_accounts_oauth2_proto_rawDesc = string([]byte{
	0x0a, 0x1f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x73, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x12, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x73, 0x1a, 0x23, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbb, 0x02, 0x0a, 0x0d, 0x4f,
	0x41, 0x75, 0x74, 0x68, 0x32, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x42, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48,
	0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x88, 0x01, 0x01, 0x12,
	0x23, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x3e, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x4f, 0x41, 0x75, 0x74,
	0x68, 0x32, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x22, 0x78, 0x0a, 0x0e, 0x4f, 0x41, 0x75, 0x74,
	0x68, 0x32, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x6d, 0x65, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x6d, 0x65, 0x70, 0x61, 0x67, 0x65,
	0x12, 0x17, 0x0a, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x69, 0x63,
	0x6f, 0x6e, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x66, 0x69, 0x76, 0x65, 0x6e, 0x65, 0x74, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x66, 0x69, 0x76,
	0x65, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x73, 0x3b, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_resources_accounts_oauth2_proto_rawDescOnce sync.Once
	file_resources_accounts_oauth2_proto_rawDescData []byte
)

func file_resources_accounts_oauth2_proto_rawDescGZIP() []byte {
	file_resources_accounts_oauth2_proto_rawDescOnce.Do(func() {
		file_resources_accounts_oauth2_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resources_accounts_oauth2_proto_rawDesc), len(file_resources_accounts_oauth2_proto_rawDesc)))
	})
	return file_resources_accounts_oauth2_proto_rawDescData
}

var file_resources_accounts_oauth2_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_resources_accounts_oauth2_proto_goTypes = []any{
	(*OAuth2Account)(nil),       // 0: resources.accounts.OAuth2Account
	(*OAuth2Provider)(nil),      // 1: resources.accounts.OAuth2Provider
	(*timestamp.Timestamp)(nil), // 2: resources.timestamp.Timestamp
}
var file_resources_accounts_oauth2_proto_depIdxs = []int32{
	2, // 0: resources.accounts.OAuth2Account.created_at:type_name -> resources.timestamp.Timestamp
	1, // 1: resources.accounts.OAuth2Account.provider:type_name -> resources.accounts.OAuth2Provider
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_resources_accounts_oauth2_proto_init() }
func file_resources_accounts_oauth2_proto_init() {
	if File_resources_accounts_oauth2_proto != nil {
		return
	}
	file_resources_accounts_oauth2_proto_msgTypes[0].OneofWrappers = []any{}
	file_resources_accounts_oauth2_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resources_accounts_oauth2_proto_rawDesc), len(file_resources_accounts_oauth2_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resources_accounts_oauth2_proto_goTypes,
		DependencyIndexes: file_resources_accounts_oauth2_proto_depIdxs,
		MessageInfos:      file_resources_accounts_oauth2_proto_msgTypes,
	}.Build()
	File_resources_accounts_oauth2_proto = out.File
	file_resources_accounts_oauth2_proto_goTypes = nil
	file_resources_accounts_oauth2_proto_depIdxs = nil
}
