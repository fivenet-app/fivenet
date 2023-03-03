// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: dispatches/dispatches.proto

package dispatches

import (
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DispatchesServiceClient is the client API for DispatchesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DispatchesServiceClient interface {
}

type dispatchesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDispatchesServiceClient(cc grpc.ClientConnInterface) DispatchesServiceClient {
	return &dispatchesServiceClient{cc}
}

// DispatchesServiceServer is the server API for DispatchesService service.
// All implementations must embed UnimplementedDispatchesServiceServer
// for forward compatibility
type DispatchesServiceServer interface {
	mustEmbedUnimplementedDispatchesServiceServer()
}

// UnimplementedDispatchesServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDispatchesServiceServer struct {
}

func (UnimplementedDispatchesServiceServer) mustEmbedUnimplementedDispatchesServiceServer() {}

// UnsafeDispatchesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DispatchesServiceServer will
// result in compilation errors.
type UnsafeDispatchesServiceServer interface {
	mustEmbedUnimplementedDispatchesServiceServer()
}

func RegisterDispatchesServiceServer(s grpc.ServiceRegistrar, srv DispatchesServiceServer) {
	s.RegisterService(&DispatchesService_ServiceDesc, srv)
}

// DispatchesService_ServiceDesc is the grpc.ServiceDesc for DispatchesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DispatchesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gen.dispatches.DispatchesService",
	HandlerType: (*DispatchesServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "dispatches/dispatches.proto",
}
