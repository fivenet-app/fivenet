// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: services/internet/internet.proto

package internet

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	InternetService_Search_FullMethodName  = "/services.internet.InternetService/Search"
	InternetService_GetPage_FullMethodName = "/services.internet.InternetService/GetPage"
)

// InternetServiceClient is the client API for InternetService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InternetServiceClient interface {
	// @perm: Name=Any
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
	// @perm: Name=Any
	GetPage(ctx context.Context, in *GetPageRequest, opts ...grpc.CallOption) (*GetPageResponse, error)
}

type internetServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInternetServiceClient(cc grpc.ClientConnInterface) InternetServiceClient {
	return &internetServiceClient{cc}
}

func (c *internetServiceClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, InternetService_Search_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internetServiceClient) GetPage(ctx context.Context, in *GetPageRequest, opts ...grpc.CallOption) (*GetPageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPageResponse)
	err := c.cc.Invoke(ctx, InternetService_GetPage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InternetServiceServer is the server API for InternetService service.
// All implementations must embed UnimplementedInternetServiceServer
// for forward compatibility.
type InternetServiceServer interface {
	// @perm: Name=Any
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	// @perm: Name=Any
	GetPage(context.Context, *GetPageRequest) (*GetPageResponse, error)
	mustEmbedUnimplementedInternetServiceServer()
}

// UnimplementedInternetServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedInternetServiceServer struct{}

func (UnimplementedInternetServiceServer) Search(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedInternetServiceServer) GetPage(context.Context, *GetPageRequest) (*GetPageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPage not implemented")
}
func (UnimplementedInternetServiceServer) mustEmbedUnimplementedInternetServiceServer() {}
func (UnimplementedInternetServiceServer) testEmbeddedByValue()                         {}

// UnsafeInternetServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InternetServiceServer will
// result in compilation errors.
type UnsafeInternetServiceServer interface {
	mustEmbedUnimplementedInternetServiceServer()
}

func RegisterInternetServiceServer(s grpc.ServiceRegistrar, srv InternetServiceServer) {
	// If the following call pancis, it indicates UnimplementedInternetServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&InternetService_ServiceDesc, srv)
}

func _InternetService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternetServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InternetService_Search_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternetServiceServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InternetService_GetPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternetServiceServer).GetPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InternetService_GetPage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternetServiceServer).GetPage(ctx, req.(*GetPageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InternetService_ServiceDesc is the grpc.ServiceDesc for InternetService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InternetService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.internet.InternetService",
	HandlerType: (*InternetServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _InternetService_Search_Handler,
		},
		{
			MethodName: "GetPage",
			Handler:    _InternetService_GetPage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/internet/internet.proto",
}
