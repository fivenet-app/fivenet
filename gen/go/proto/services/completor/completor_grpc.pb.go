// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
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

const (
	CompletorService_CompleteCitizens_FullMethodName           = "/services.completor.CompletorService/CompleteCitizens"
	CompletorService_CompleteJobs_FullMethodName               = "/services.completor.CompletorService/CompleteJobs"
	CompletorService_CompleteDocumentCategories_FullMethodName = "/services.completor.CompletorService/CompleteDocumentCategories"
	CompletorService_ListLawBooks_FullMethodName               = "/services.completor.CompletorService/ListLawBooks"
	CompletorService_CompleteCitizenAttributes_FullMethodName  = "/services.completor.CompletorService/CompleteCitizenAttributes"
)

// CompletorServiceClient is the client API for CompletorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CompletorServiceClient interface {
	// @perm
	CompleteCitizens(ctx context.Context, in *CompleteCitizensRequest, opts ...grpc.CallOption) (*CompleteCitizensRespoonse, error)
	// @perm
	CompleteJobs(ctx context.Context, in *CompleteJobsRequest, opts ...grpc.CallOption) (*CompleteJobsResponse, error)
	// @perm: Attrs=Jobs/JobList
	CompleteDocumentCategories(ctx context.Context, in *CompleteDocumentCategoriesRequest, opts ...grpc.CallOption) (*CompleteDocumentCategoriesResponse, error)
	// @perm: Name=Any
	ListLawBooks(ctx context.Context, in *ListLawBooksRequest, opts ...grpc.CallOption) (*ListLawBooksResponse, error)
	// @perm: Attrs=Jobs/JobList
	CompleteCitizenAttributes(ctx context.Context, in *CompleteCitizenAttributesRequest, opts ...grpc.CallOption) (*CompleteCitizenAttributesResponse, error)
}

type completorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCompletorServiceClient(cc grpc.ClientConnInterface) CompletorServiceClient {
	return &completorServiceClient{cc}
}

func (c *completorServiceClient) CompleteCitizens(ctx context.Context, in *CompleteCitizensRequest, opts ...grpc.CallOption) (*CompleteCitizensRespoonse, error) {
	out := new(CompleteCitizensRespoonse)
	err := c.cc.Invoke(ctx, CompletorService_CompleteCitizens_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *completorServiceClient) CompleteJobs(ctx context.Context, in *CompleteJobsRequest, opts ...grpc.CallOption) (*CompleteJobsResponse, error) {
	out := new(CompleteJobsResponse)
	err := c.cc.Invoke(ctx, CompletorService_CompleteJobs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *completorServiceClient) CompleteDocumentCategories(ctx context.Context, in *CompleteDocumentCategoriesRequest, opts ...grpc.CallOption) (*CompleteDocumentCategoriesResponse, error) {
	out := new(CompleteDocumentCategoriesResponse)
	err := c.cc.Invoke(ctx, CompletorService_CompleteDocumentCategories_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *completorServiceClient) ListLawBooks(ctx context.Context, in *ListLawBooksRequest, opts ...grpc.CallOption) (*ListLawBooksResponse, error) {
	out := new(ListLawBooksResponse)
	err := c.cc.Invoke(ctx, CompletorService_ListLawBooks_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *completorServiceClient) CompleteCitizenAttributes(ctx context.Context, in *CompleteCitizenAttributesRequest, opts ...grpc.CallOption) (*CompleteCitizenAttributesResponse, error) {
	out := new(CompleteCitizenAttributesResponse)
	err := c.cc.Invoke(ctx, CompletorService_CompleteCitizenAttributes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompletorServiceServer is the server API for CompletorService service.
// All implementations must embed UnimplementedCompletorServiceServer
// for forward compatibility
type CompletorServiceServer interface {
	// @perm
	CompleteCitizens(context.Context, *CompleteCitizensRequest) (*CompleteCitizensRespoonse, error)
	// @perm
	CompleteJobs(context.Context, *CompleteJobsRequest) (*CompleteJobsResponse, error)
	// @perm: Attrs=Jobs/JobList
	CompleteDocumentCategories(context.Context, *CompleteDocumentCategoriesRequest) (*CompleteDocumentCategoriesResponse, error)
	// @perm: Name=Any
	ListLawBooks(context.Context, *ListLawBooksRequest) (*ListLawBooksResponse, error)
	// @perm: Attrs=Jobs/JobList
	CompleteCitizenAttributes(context.Context, *CompleteCitizenAttributesRequest) (*CompleteCitizenAttributesResponse, error)
	mustEmbedUnimplementedCompletorServiceServer()
}

// UnimplementedCompletorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCompletorServiceServer struct {
}

func (UnimplementedCompletorServiceServer) CompleteCitizens(context.Context, *CompleteCitizensRequest) (*CompleteCitizensRespoonse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteCitizens not implemented")
}
func (UnimplementedCompletorServiceServer) CompleteJobs(context.Context, *CompleteJobsRequest) (*CompleteJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteJobs not implemented")
}
func (UnimplementedCompletorServiceServer) CompleteDocumentCategories(context.Context, *CompleteDocumentCategoriesRequest) (*CompleteDocumentCategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteDocumentCategories not implemented")
}
func (UnimplementedCompletorServiceServer) ListLawBooks(context.Context, *ListLawBooksRequest) (*ListLawBooksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLawBooks not implemented")
}
func (UnimplementedCompletorServiceServer) CompleteCitizenAttributes(context.Context, *CompleteCitizenAttributesRequest) (*CompleteCitizenAttributesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteCitizenAttributes not implemented")
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

func _CompletorService_CompleteCitizens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteCitizensRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompletorServiceServer).CompleteCitizens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CompletorService_CompleteCitizens_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompletorServiceServer).CompleteCitizens(ctx, req.(*CompleteCitizensRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CompletorService_CompleteJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteJobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompletorServiceServer).CompleteJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CompletorService_CompleteJobs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompletorServiceServer).CompleteJobs(ctx, req.(*CompleteJobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CompletorService_CompleteDocumentCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteDocumentCategoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompletorServiceServer).CompleteDocumentCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CompletorService_CompleteDocumentCategories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompletorServiceServer).CompleteDocumentCategories(ctx, req.(*CompleteDocumentCategoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CompletorService_ListLawBooks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLawBooksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompletorServiceServer).ListLawBooks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CompletorService_ListLawBooks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompletorServiceServer).ListLawBooks(ctx, req.(*ListLawBooksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CompletorService_CompleteCitizenAttributes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteCitizenAttributesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompletorServiceServer).CompleteCitizenAttributes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CompletorService_CompleteCitizenAttributes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompletorServiceServer).CompleteCitizenAttributes(ctx, req.(*CompleteCitizenAttributesRequest))
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
			MethodName: "CompleteCitizens",
			Handler:    _CompletorService_CompleteCitizens_Handler,
		},
		{
			MethodName: "CompleteJobs",
			Handler:    _CompletorService_CompleteJobs_Handler,
		},
		{
			MethodName: "CompleteDocumentCategories",
			Handler:    _CompletorService_CompleteDocumentCategories_Handler,
		},
		{
			MethodName: "ListLawBooks",
			Handler:    _CompletorService_ListLawBooks_Handler,
		},
		{
			MethodName: "CompleteCitizenAttributes",
			Handler:    _CompletorService_CompleteCitizenAttributes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/completor/completor.proto",
}
