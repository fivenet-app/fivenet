// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.20.3
// source: services/wiki/wiki.proto

package wiki

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	WikiService_ListPages_FullMethodName          = "/services.wiki.WikiService/ListPages"
	WikiService_GetPage_FullMethodName            = "/services.wiki.WikiService/GetPage"
	WikiService_CreateOrUpdatePage_FullMethodName = "/services.wiki.WikiService/CreateOrUpdatePage"
	WikiService_DeletePage_FullMethodName         = "/services.wiki.WikiService/DeletePage"
	WikiService_GetPageHistory_FullMethodName     = "/services.wiki.WikiService/GetPageHistory"
)

// WikiServiceClient is the client API for WikiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WikiServiceClient interface {
	// @perm
	ListPages(ctx context.Context, in *ListPagesRequest, opts ...grpc.CallOption) (*ListPagesResponse, error)
	// @perm
	GetPage(ctx context.Context, in *GetPageRequest, opts ...grpc.CallOption) (*GetPageResponse, error)
	// @perm
	CreateOrUpdatePage(ctx context.Context, in *CreateOrUpdatePageRequest, opts ...grpc.CallOption) (*CreateOrUpdatePageResponse, error)
	// @perm
	DeletePage(ctx context.Context, in *DeletePageRequest, opts ...grpc.CallOption) (*DeletePageResponse, error)
	// @perm
	GetPageHistory(ctx context.Context, in *GetPageHistoryRequest, opts ...grpc.CallOption) (*GetPageHistoryResponse, error)
}

type wikiServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWikiServiceClient(cc grpc.ClientConnInterface) WikiServiceClient {
	return &wikiServiceClient{cc}
}

func (c *wikiServiceClient) ListPages(ctx context.Context, in *ListPagesRequest, opts ...grpc.CallOption) (*ListPagesResponse, error) {
	out := new(ListPagesResponse)
	err := c.cc.Invoke(ctx, WikiService_ListPages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wikiServiceClient) GetPage(ctx context.Context, in *GetPageRequest, opts ...grpc.CallOption) (*GetPageResponse, error) {
	out := new(GetPageResponse)
	err := c.cc.Invoke(ctx, WikiService_GetPage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wikiServiceClient) CreateOrUpdatePage(ctx context.Context, in *CreateOrUpdatePageRequest, opts ...grpc.CallOption) (*CreateOrUpdatePageResponse, error) {
	out := new(CreateOrUpdatePageResponse)
	err := c.cc.Invoke(ctx, WikiService_CreateOrUpdatePage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wikiServiceClient) DeletePage(ctx context.Context, in *DeletePageRequest, opts ...grpc.CallOption) (*DeletePageResponse, error) {
	out := new(DeletePageResponse)
	err := c.cc.Invoke(ctx, WikiService_DeletePage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wikiServiceClient) GetPageHistory(ctx context.Context, in *GetPageHistoryRequest, opts ...grpc.CallOption) (*GetPageHistoryResponse, error) {
	out := new(GetPageHistoryResponse)
	err := c.cc.Invoke(ctx, WikiService_GetPageHistory_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WikiServiceServer is the server API for WikiService service.
// All implementations must embed UnimplementedWikiServiceServer
// for forward compatibility
type WikiServiceServer interface {
	// @perm
	ListPages(context.Context, *ListPagesRequest) (*ListPagesResponse, error)
	// @perm
	GetPage(context.Context, *GetPageRequest) (*GetPageResponse, error)
	// @perm
	CreateOrUpdatePage(context.Context, *CreateOrUpdatePageRequest) (*CreateOrUpdatePageResponse, error)
	// @perm
	DeletePage(context.Context, *DeletePageRequest) (*DeletePageResponse, error)
	// @perm
	GetPageHistory(context.Context, *GetPageHistoryRequest) (*GetPageHistoryResponse, error)
	mustEmbedUnimplementedWikiServiceServer()
}

// UnimplementedWikiServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWikiServiceServer struct {
}

func (UnimplementedWikiServiceServer) ListPages(context.Context, *ListPagesRequest) (*ListPagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPages not implemented")
}
func (UnimplementedWikiServiceServer) GetPage(context.Context, *GetPageRequest) (*GetPageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPage not implemented")
}
func (UnimplementedWikiServiceServer) CreateOrUpdatePage(context.Context, *CreateOrUpdatePageRequest) (*CreateOrUpdatePageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdatePage not implemented")
}
func (UnimplementedWikiServiceServer) DeletePage(context.Context, *DeletePageRequest) (*DeletePageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePage not implemented")
}
func (UnimplementedWikiServiceServer) GetPageHistory(context.Context, *GetPageHistoryRequest) (*GetPageHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPageHistory not implemented")
}
func (UnimplementedWikiServiceServer) mustEmbedUnimplementedWikiServiceServer() {}

// UnsafeWikiServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WikiServiceServer will
// result in compilation errors.
type UnsafeWikiServiceServer interface {
	mustEmbedUnimplementedWikiServiceServer()
}

func RegisterWikiServiceServer(s grpc.ServiceRegistrar, srv WikiServiceServer) {
	s.RegisterService(&WikiService_ServiceDesc, srv)
}

func _WikiService_ListPages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WikiServiceServer).ListPages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WikiService_ListPages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WikiServiceServer).ListPages(ctx, req.(*ListPagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WikiService_GetPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WikiServiceServer).GetPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WikiService_GetPage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WikiServiceServer).GetPage(ctx, req.(*GetPageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WikiService_CreateOrUpdatePage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrUpdatePageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WikiServiceServer).CreateOrUpdatePage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WikiService_CreateOrUpdatePage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WikiServiceServer).CreateOrUpdatePage(ctx, req.(*CreateOrUpdatePageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WikiService_DeletePage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WikiServiceServer).DeletePage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WikiService_DeletePage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WikiServiceServer).DeletePage(ctx, req.(*DeletePageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WikiService_GetPageHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPageHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WikiServiceServer).GetPageHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WikiService_GetPageHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WikiServiceServer).GetPageHistory(ctx, req.(*GetPageHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WikiService_ServiceDesc is the grpc.ServiceDesc for WikiService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WikiService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.wiki.WikiService",
	HandlerType: (*WikiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListPages",
			Handler:    _WikiService_ListPages_Handler,
		},
		{
			MethodName: "GetPage",
			Handler:    _WikiService_GetPage_Handler,
		},
		{
			MethodName: "CreateOrUpdatePage",
			Handler:    _WikiService_CreateOrUpdatePage_Handler,
		},
		{
			MethodName: "DeletePage",
			Handler:    _WikiService_DeletePage_Handler,
		},
		{
			MethodName: "GetPageHistory",
			Handler:    _WikiService_GetPageHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/wiki/wiki.proto",
}
