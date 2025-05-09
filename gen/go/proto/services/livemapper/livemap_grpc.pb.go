// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: services/livemapper/livemap.proto

package livemapper

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
	LivemapperService_Stream_FullMethodName               = "/services.livemapper.LivemapperService/Stream"
	LivemapperService_CreateOrUpdateMarker_FullMethodName = "/services.livemapper.LivemapperService/CreateOrUpdateMarker"
	LivemapperService_DeleteMarker_FullMethodName         = "/services.livemapper.LivemapperService/DeleteMarker"
)

// LivemapperServiceClient is the client API for LivemapperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LivemapperServiceClient interface {
	// @perm: Attrs=Markers/JobList|Players/JobGradeList
	Stream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StreamResponse], error)
	// @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}
	CreateOrUpdateMarker(ctx context.Context, in *CreateOrUpdateMarkerRequest, opts ...grpc.CallOption) (*CreateOrUpdateMarkerResponse, error)
	// @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}
	DeleteMarker(ctx context.Context, in *DeleteMarkerRequest, opts ...grpc.CallOption) (*DeleteMarkerResponse, error)
}

type livemapperServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLivemapperServiceClient(cc grpc.ClientConnInterface) LivemapperServiceClient {
	return &livemapperServiceClient{cc}
}

func (c *livemapperServiceClient) Stream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StreamResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &LivemapperService_ServiceDesc.Streams[0], LivemapperService_Stream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StreamRequest, StreamResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LivemapperService_StreamClient = grpc.ServerStreamingClient[StreamResponse]

func (c *livemapperServiceClient) CreateOrUpdateMarker(ctx context.Context, in *CreateOrUpdateMarkerRequest, opts ...grpc.CallOption) (*CreateOrUpdateMarkerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateOrUpdateMarkerResponse)
	err := c.cc.Invoke(ctx, LivemapperService_CreateOrUpdateMarker_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *livemapperServiceClient) DeleteMarker(ctx context.Context, in *DeleteMarkerRequest, opts ...grpc.CallOption) (*DeleteMarkerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMarkerResponse)
	err := c.cc.Invoke(ctx, LivemapperService_DeleteMarker_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LivemapperServiceServer is the server API for LivemapperService service.
// All implementations must embed UnimplementedLivemapperServiceServer
// for forward compatibility.
type LivemapperServiceServer interface {
	// @perm: Attrs=Markers/JobList|Players/JobGradeList
	Stream(*StreamRequest, grpc.ServerStreamingServer[StreamResponse]) error
	// @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}
	CreateOrUpdateMarker(context.Context, *CreateOrUpdateMarkerRequest) (*CreateOrUpdateMarkerResponse, error)
	// @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}
	DeleteMarker(context.Context, *DeleteMarkerRequest) (*DeleteMarkerResponse, error)
	mustEmbedUnimplementedLivemapperServiceServer()
}

// UnimplementedLivemapperServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLivemapperServiceServer struct{}

func (UnimplementedLivemapperServiceServer) Stream(*StreamRequest, grpc.ServerStreamingServer[StreamResponse]) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedLivemapperServiceServer) CreateOrUpdateMarker(context.Context, *CreateOrUpdateMarkerRequest) (*CreateOrUpdateMarkerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateMarker not implemented")
}
func (UnimplementedLivemapperServiceServer) DeleteMarker(context.Context, *DeleteMarkerRequest) (*DeleteMarkerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMarker not implemented")
}
func (UnimplementedLivemapperServiceServer) mustEmbedUnimplementedLivemapperServiceServer() {}
func (UnimplementedLivemapperServiceServer) testEmbeddedByValue()                           {}

// UnsafeLivemapperServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LivemapperServiceServer will
// result in compilation errors.
type UnsafeLivemapperServiceServer interface {
	mustEmbedUnimplementedLivemapperServiceServer()
}

func RegisterLivemapperServiceServer(s grpc.ServiceRegistrar, srv LivemapperServiceServer) {
	// If the following call pancis, it indicates UnimplementedLivemapperServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LivemapperService_ServiceDesc, srv)
}

func _LivemapperService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LivemapperServiceServer).Stream(m, &grpc.GenericServerStream[StreamRequest, StreamResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LivemapperService_StreamServer = grpc.ServerStreamingServer[StreamResponse]

func _LivemapperService_CreateOrUpdateMarker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrUpdateMarkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LivemapperServiceServer).CreateOrUpdateMarker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LivemapperService_CreateOrUpdateMarker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LivemapperServiceServer).CreateOrUpdateMarker(ctx, req.(*CreateOrUpdateMarkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LivemapperService_DeleteMarker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMarkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LivemapperServiceServer).DeleteMarker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LivemapperService_DeleteMarker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LivemapperServiceServer).DeleteMarker(ctx, req.(*DeleteMarkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LivemapperService_ServiceDesc is the grpc.ServiceDesc for LivemapperService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LivemapperService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.livemapper.LivemapperService",
	HandlerType: (*LivemapperServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrUpdateMarker",
			Handler:    _LivemapperService_CreateOrUpdateMarker_Handler,
		},
		{
			MethodName: "DeleteMarker",
			Handler:    _LivemapperService_DeleteMarker_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _LivemapperService_Stream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "services/livemapper/livemap.proto",
}
