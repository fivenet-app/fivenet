// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: services/completor/completor.proto

package completor

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

// CompletorServiceClient is the client API for CompletorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CompletorServiceClient interface {
	// @permission
	CompleteCharNames(ctx context.Context, in *CompleteCharNamesRequest, opts ...grpc.CallOption) (*CompleteCharNamesRespoonse, error)
	// @permission
	CompleteJobNames(ctx context.Context, in *CompleteJobNamesRequest, opts ...grpc.CallOption) (*CompleteJobNamesResponse, error)
	// @permission
	CompleteJobGrades(ctx context.Context, in *CompleteJobGradesRequest, opts ...grpc.CallOption) (*CompleteJobGradesResponse, error)
	// @permission: PerJob=true
	CompleteDocumentCategory(ctx context.Context, in *CompleteDocumentCategoryRequest, opts ...grpc.CallOption) (*CompleteDocumentCategoryResponse, error)
}

type completorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCompletorServiceClient(cc grpc.ClientConnInterface) CompletorServiceClient {
	return &completorServiceClient{cc}
}

func (c *completorServiceClient) CompleteCharNames(ctx context.Context, in *CompleteCharNamesRequest, opts ...grpc.CallOption) (*CompleteCharNamesRespoonse, error) {
	out := new(CompleteCharNamesRespoonse)
	err := c.cc.Invoke(ctx, "/services.completor.CompletorService/CompleteCharNames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *completorServiceClient) CompleteJobNames(ctx context.Context, in *CompleteJobNamesRequest, opts ...grpc.CallOption) (*CompleteJobNamesResponse, error) {
	out := new(CompleteJobNamesResponse)
	err := c.cc.Invoke(ctx, "/services.completor.CompletorService/CompleteJobNames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *completorServiceClient) CompleteJobGrades(ctx context.Context, in *CompleteJobGradesRequest, opts ...grpc.CallOption) (*CompleteJobGradesResponse, error) {
	out := new(CompleteJobGradesResponse)
	err := c.cc.Invoke(ctx, "/services.completor.CompletorService/CompleteJobGrades", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *completorServiceClient) CompleteDocumentCategory(ctx context.Context, in *CompleteDocumentCategoryRequest, opts ...grpc.CallOption) (*CompleteDocumentCategoryResponse, error) {
	out := new(CompleteDocumentCategoryResponse)
	err := c.cc.Invoke(ctx, "/services.completor.CompletorService/CompleteDocumentCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompletorServiceServer is the server API for CompletorService service.
// All implementations must embed UnimplementedCompletorServiceServer
// for forward compatibility
type CompletorServiceServer interface {
	// @permission
	CompleteCharNames(context.Context, *CompleteCharNamesRequest) (*CompleteCharNamesRespoonse, error)
	// @permission
	CompleteJobNames(context.Context, *CompleteJobNamesRequest) (*CompleteJobNamesResponse, error)
	// @permission
	CompleteJobGrades(context.Context, *CompleteJobGradesRequest) (*CompleteJobGradesResponse, error)
	// @permission: PerJob=true
	CompleteDocumentCategory(context.Context, *CompleteDocumentCategoryRequest) (*CompleteDocumentCategoryResponse, error)
	mustEmbedUnimplementedCompletorServiceServer()
}

// UnimplementedCompletorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCompletorServiceServer struct {
}

func (UnimplementedCompletorServiceServer) CompleteCharNames(context.Context, *CompleteCharNamesRequest) (*CompleteCharNamesRespoonse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteCharNames not implemented")
}
func (UnimplementedCompletorServiceServer) CompleteJobNames(context.Context, *CompleteJobNamesRequest) (*CompleteJobNamesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteJobNames not implemented")
}
func (UnimplementedCompletorServiceServer) CompleteJobGrades(context.Context, *CompleteJobGradesRequest) (*CompleteJobGradesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteJobGrades not implemented")
}
func (UnimplementedCompletorServiceServer) CompleteDocumentCategory(context.Context, *CompleteDocumentCategoryRequest) (*CompleteDocumentCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteDocumentCategory not implemented")
}
func (UnimplementedCompletorServiceServer) mustEmbedUnimplementedCompletorServiceServer() {}

// UnsafeCompletorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CompletorServiceServer will
// result in compilation errors.
type UnsafeCompletorServiceServer interface {
	mustEmbedUnimplementedCompletorServiceServer()
}

func RegisterCompletorServiceServer(s grpc.ServiceRegistrar, srv CompletorServiceServer) {
	s.RegisterService(&CompletorService_ServiceDesc, srv)
}

func _CompletorService_CompleteCharNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteCharNamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompletorServiceServer).CompleteCharNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.completor.CompletorService/CompleteCharNames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompletorServiceServer).CompleteCharNames(ctx, req.(*CompleteCharNamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CompletorService_CompleteJobNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteJobNamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompletorServiceServer).CompleteJobNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.completor.CompletorService/CompleteJobNames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompletorServiceServer).CompleteJobNames(ctx, req.(*CompleteJobNamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CompletorService_CompleteJobGrades_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteJobGradesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompletorServiceServer).CompleteJobGrades(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.completor.CompletorService/CompleteJobGrades",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompletorServiceServer).CompleteJobGrades(ctx, req.(*CompleteJobGradesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CompletorService_CompleteDocumentCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteDocumentCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompletorServiceServer).CompleteDocumentCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.completor.CompletorService/CompleteDocumentCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompletorServiceServer).CompleteDocumentCategory(ctx, req.(*CompleteDocumentCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CompletorService_ServiceDesc is the grpc.ServiceDesc for CompletorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CompletorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.completor.CompletorService",
	HandlerType: (*CompletorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CompleteCharNames",
			Handler:    _CompletorService_CompleteCharNames_Handler,
		},
		{
			MethodName: "CompleteJobNames",
			Handler:    _CompletorService_CompleteJobNames_Handler,
		},
		{
			MethodName: "CompleteJobGrades",
			Handler:    _CompletorService_CompleteJobGrades_Handler,
		},
		{
			MethodName: "CompleteDocumentCategory",
			Handler:    _CompletorService_CompleteDocumentCategory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/completor/completor.proto",
}
