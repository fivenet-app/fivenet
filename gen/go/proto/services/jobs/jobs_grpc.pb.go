// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: services/jobs/jobs.proto

package jobs

import (
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const ()

// JobsServiceClient is the client API for JobsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JobsServiceClient interface {
}

type jobsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewJobsServiceClient(cc grpc.ClientConnInterface) JobsServiceClient {
	return &jobsServiceClient{cc}
}

// JobsServiceServer is the server API for JobsService service.
// All implementations must embed UnimplementedJobsServiceServer
// for forward compatibility
type JobsServiceServer interface {
	mustEmbedUnimplementedJobsServiceServer()
}

// UnimplementedJobsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedJobsServiceServer struct {
}

func (UnimplementedJobsServiceServer) mustEmbedUnimplementedJobsServiceServer() {}

// UnsafeJobsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JobsServiceServer will
// result in compilation errors.
type UnsafeJobsServiceServer interface {
	mustEmbedUnimplementedJobsServiceServer()
}

func RegisterJobsServiceServer(s grpc.ServiceRegistrar, srv JobsServiceServer) {
	s.RegisterService(&JobsService_ServiceDesc, srv)
}

// JobsService_ServiceDesc is the grpc.ServiceDesc for JobsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var JobsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.jobs.JobsService",
	HandlerType: (*JobsServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "services/jobs/jobs.proto",
}
