// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: services/centrum/centrum.proto

package centrum

import (
	context "context"
	common "github.com/galexrt/fivenet/gen/go/proto/resources/common"
	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	CentrumService_GetSettings_FullMethodName          = "/services.centrum.CentrumService/GetSettings"
	CentrumService_UpdateSettings_FullMethodName       = "/services.centrum.CentrumService/UpdateSettings"
	CentrumService_CreateOrUpdateUnit_FullMethodName   = "/services.centrum.CentrumService/CreateOrUpdateUnit"
	CentrumService_DeleteUnit_FullMethodName           = "/services.centrum.CentrumService/DeleteUnit"
	CentrumService_TakeControl_FullMethodName          = "/services.centrum.CentrumService/TakeControl"
	CentrumService_AssignUnit_FullMethodName           = "/services.centrum.CentrumService/AssignUnit"
	CentrumService_AssignDispatch_FullMethodName       = "/services.centrum.CentrumService/AssignDispatch"
	CentrumService_CreateDispatch_FullMethodName       = "/services.centrum.CentrumService/CreateDispatch"
	CentrumService_UpdateDispatch_FullMethodName       = "/services.centrum.CentrumService/UpdateDispatch"
	CentrumService_TakeDispatch_FullMethodName         = "/services.centrum.CentrumService/TakeDispatch"
	CentrumService_UpdateDispatchStatus_FullMethodName = "/services.centrum.CentrumService/UpdateDispatchStatus"
	CentrumService_UpdateUnitStatus_FullMethodName     = "/services.centrum.CentrumService/UpdateUnitStatus"
	CentrumService_ListDispatches_FullMethodName       = "/services.centrum.CentrumService/ListDispatches"
	CentrumService_ListUnits_FullMethodName            = "/services.centrum.CentrumService/ListUnits"
	CentrumService_ListUnitActivity_FullMethodName     = "/services.centrum.CentrumService/ListUnitActivity"
	CentrumService_ListDispatchActivity_FullMethodName = "/services.centrum.CentrumService/ListDispatchActivity"
	CentrumService_Stream_FullMethodName               = "/services.centrum.CentrumService/Stream"
)

// CentrumServiceClient is the client API for CentrumService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CentrumServiceClient interface {
	// @perm: Name=Stream
	GetSettings(ctx context.Context, in *common.EmptyRequest, opts ...grpc.CallOption) (*dispatch.Settings, error)
	// @perm
	UpdateSettings(ctx context.Context, in *dispatch.Settings, opts ...grpc.CallOption) (*dispatch.Settings, error)
	// @perm
	CreateOrUpdateUnit(ctx context.Context, in *CreateOrUpdateUnitRequest, opts ...grpc.CallOption) (*CreateOrUpdateUnitResponse, error)
	// @perm
	DeleteUnit(ctx context.Context, in *DeleteUnitRequest, opts ...grpc.CallOption) (*DeleteUnitResponse, error)
	// @perm
	TakeControl(ctx context.Context, in *TakeControlRequest, opts ...grpc.CallOption) (*TakeControlResponse, error)
	// @perm: Name=TakeControl
	AssignUnit(ctx context.Context, in *AssignUnitRequest, opts ...grpc.CallOption) (*AssignUnitResponse, error)
	// @perm: Name=TakeControl
	AssignDispatch(ctx context.Context, in *AssignDispatchRequest, opts ...grpc.CallOption) (*AssignDispatchResponse, error)
	// @perm
	CreateDispatch(ctx context.Context, in *CreateDispatchRequest, opts ...grpc.CallOption) (*CreateDispatchResponse, error)
	// @perm
	UpdateDispatch(ctx context.Context, in *UpdateDispatchRequest, opts ...grpc.CallOption) (*UpdateDispatchResponse, error)
	// @perm
	TakeDispatch(ctx context.Context, in *TakeDispatchRequest, opts ...grpc.CallOption) (*TakeDispatchResponse, error)
	// @perm: Name=TakeDispatch
	UpdateDispatchStatus(ctx context.Context, in *UpdateDispatchStatusRequest, opts ...grpc.CallOption) (*UpdateDispatchStatusResponse, error)
	// @perm: Name=TakeDispatch
	UpdateUnitStatus(ctx context.Context, in *UpdateUnitStatusRequest, opts ...grpc.CallOption) (*UpdateUnitStatusResponse, error)
	// @perm: Name=Stream
	ListDispatches(ctx context.Context, in *ListDispatchesRequest, opts ...grpc.CallOption) (*ListDispatchesResponse, error)
	// @perm: Name=Stream
	ListUnits(ctx context.Context, in *ListUnitsRequest, opts ...grpc.CallOption) (*ListUnitsResponse, error)
	// @perm: Name=Stream
	ListUnitActivity(ctx context.Context, in *ListActivityRequest, opts ...grpc.CallOption) (*ListUnitActivityResponse, error)
	// @perm: Name=Stream
	ListDispatchActivity(ctx context.Context, in *ListActivityRequest, opts ...grpc.CallOption) (*ListDispatchActivityResponse, error)
	// @perm
	Stream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (CentrumService_StreamClient, error)
}

type centrumServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCentrumServiceClient(cc grpc.ClientConnInterface) CentrumServiceClient {
	return &centrumServiceClient{cc}
}

func (c *centrumServiceClient) GetSettings(ctx context.Context, in *common.EmptyRequest, opts ...grpc.CallOption) (*dispatch.Settings, error) {
	out := new(dispatch.Settings)
	err := c.cc.Invoke(ctx, CentrumService_GetSettings_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) UpdateSettings(ctx context.Context, in *dispatch.Settings, opts ...grpc.CallOption) (*dispatch.Settings, error) {
	out := new(dispatch.Settings)
	err := c.cc.Invoke(ctx, CentrumService_UpdateSettings_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) CreateOrUpdateUnit(ctx context.Context, in *CreateOrUpdateUnitRequest, opts ...grpc.CallOption) (*CreateOrUpdateUnitResponse, error) {
	out := new(CreateOrUpdateUnitResponse)
	err := c.cc.Invoke(ctx, CentrumService_CreateOrUpdateUnit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) DeleteUnit(ctx context.Context, in *DeleteUnitRequest, opts ...grpc.CallOption) (*DeleteUnitResponse, error) {
	out := new(DeleteUnitResponse)
	err := c.cc.Invoke(ctx, CentrumService_DeleteUnit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) TakeControl(ctx context.Context, in *TakeControlRequest, opts ...grpc.CallOption) (*TakeControlResponse, error) {
	out := new(TakeControlResponse)
	err := c.cc.Invoke(ctx, CentrumService_TakeControl_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) AssignUnit(ctx context.Context, in *AssignUnitRequest, opts ...grpc.CallOption) (*AssignUnitResponse, error) {
	out := new(AssignUnitResponse)
	err := c.cc.Invoke(ctx, CentrumService_AssignUnit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) AssignDispatch(ctx context.Context, in *AssignDispatchRequest, opts ...grpc.CallOption) (*AssignDispatchResponse, error) {
	out := new(AssignDispatchResponse)
	err := c.cc.Invoke(ctx, CentrumService_AssignDispatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) CreateDispatch(ctx context.Context, in *CreateDispatchRequest, opts ...grpc.CallOption) (*CreateDispatchResponse, error) {
	out := new(CreateDispatchResponse)
	err := c.cc.Invoke(ctx, CentrumService_CreateDispatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) UpdateDispatch(ctx context.Context, in *UpdateDispatchRequest, opts ...grpc.CallOption) (*UpdateDispatchResponse, error) {
	out := new(UpdateDispatchResponse)
	err := c.cc.Invoke(ctx, CentrumService_UpdateDispatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) TakeDispatch(ctx context.Context, in *TakeDispatchRequest, opts ...grpc.CallOption) (*TakeDispatchResponse, error) {
	out := new(TakeDispatchResponse)
	err := c.cc.Invoke(ctx, CentrumService_TakeDispatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) UpdateDispatchStatus(ctx context.Context, in *UpdateDispatchStatusRequest, opts ...grpc.CallOption) (*UpdateDispatchStatusResponse, error) {
	out := new(UpdateDispatchStatusResponse)
	err := c.cc.Invoke(ctx, CentrumService_UpdateDispatchStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) UpdateUnitStatus(ctx context.Context, in *UpdateUnitStatusRequest, opts ...grpc.CallOption) (*UpdateUnitStatusResponse, error) {
	out := new(UpdateUnitStatusResponse)
	err := c.cc.Invoke(ctx, CentrumService_UpdateUnitStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) ListDispatches(ctx context.Context, in *ListDispatchesRequest, opts ...grpc.CallOption) (*ListDispatchesResponse, error) {
	out := new(ListDispatchesResponse)
	err := c.cc.Invoke(ctx, CentrumService_ListDispatches_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) ListUnits(ctx context.Context, in *ListUnitsRequest, opts ...grpc.CallOption) (*ListUnitsResponse, error) {
	out := new(ListUnitsResponse)
	err := c.cc.Invoke(ctx, CentrumService_ListUnits_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) ListUnitActivity(ctx context.Context, in *ListActivityRequest, opts ...grpc.CallOption) (*ListUnitActivityResponse, error) {
	out := new(ListUnitActivityResponse)
	err := c.cc.Invoke(ctx, CentrumService_ListUnitActivity_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) ListDispatchActivity(ctx context.Context, in *ListActivityRequest, opts ...grpc.CallOption) (*ListDispatchActivityResponse, error) {
	out := new(ListDispatchActivityResponse)
	err := c.cc.Invoke(ctx, CentrumService_ListDispatchActivity_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centrumServiceClient) Stream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (CentrumService_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &CentrumService_ServiceDesc.Streams[0], CentrumService_Stream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &centrumServiceStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CentrumService_StreamClient interface {
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type centrumServiceStreamClient struct {
	grpc.ClientStream
}

func (x *centrumServiceStreamClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CentrumServiceServer is the server API for CentrumService service.
// All implementations must embed UnimplementedCentrumServiceServer
// for forward compatibility
type CentrumServiceServer interface {
	// @perm: Name=Stream
	GetSettings(context.Context, *common.EmptyRequest) (*dispatch.Settings, error)
	// @perm
	UpdateSettings(context.Context, *dispatch.Settings) (*dispatch.Settings, error)
	// @perm
	CreateOrUpdateUnit(context.Context, *CreateOrUpdateUnitRequest) (*CreateOrUpdateUnitResponse, error)
	// @perm
	DeleteUnit(context.Context, *DeleteUnitRequest) (*DeleteUnitResponse, error)
	// @perm
	TakeControl(context.Context, *TakeControlRequest) (*TakeControlResponse, error)
	// @perm: Name=TakeControl
	AssignUnit(context.Context, *AssignUnitRequest) (*AssignUnitResponse, error)
	// @perm: Name=TakeControl
	AssignDispatch(context.Context, *AssignDispatchRequest) (*AssignDispatchResponse, error)
	// @perm
	CreateDispatch(context.Context, *CreateDispatchRequest) (*CreateDispatchResponse, error)
	// @perm
	UpdateDispatch(context.Context, *UpdateDispatchRequest) (*UpdateDispatchResponse, error)
	// @perm
	TakeDispatch(context.Context, *TakeDispatchRequest) (*TakeDispatchResponse, error)
	// @perm: Name=TakeDispatch
	UpdateDispatchStatus(context.Context, *UpdateDispatchStatusRequest) (*UpdateDispatchStatusResponse, error)
	// @perm: Name=TakeDispatch
	UpdateUnitStatus(context.Context, *UpdateUnitStatusRequest) (*UpdateUnitStatusResponse, error)
	// @perm: Name=Stream
	ListDispatches(context.Context, *ListDispatchesRequest) (*ListDispatchesResponse, error)
	// @perm: Name=Stream
	ListUnits(context.Context, *ListUnitsRequest) (*ListUnitsResponse, error)
	// @perm: Name=Stream
	ListUnitActivity(context.Context, *ListActivityRequest) (*ListUnitActivityResponse, error)
	// @perm: Name=Stream
	ListDispatchActivity(context.Context, *ListActivityRequest) (*ListDispatchActivityResponse, error)
	// @perm
	Stream(*StreamRequest, CentrumService_StreamServer) error
	mustEmbedUnimplementedCentrumServiceServer()
}

// UnimplementedCentrumServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCentrumServiceServer struct {
}

func (UnimplementedCentrumServiceServer) GetSettings(context.Context, *common.EmptyRequest) (*dispatch.Settings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSettings not implemented")
}
func (UnimplementedCentrumServiceServer) UpdateSettings(context.Context, *dispatch.Settings) (*dispatch.Settings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSettings not implemented")
}
func (UnimplementedCentrumServiceServer) CreateOrUpdateUnit(context.Context, *CreateOrUpdateUnitRequest) (*CreateOrUpdateUnitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateUnit not implemented")
}
func (UnimplementedCentrumServiceServer) DeleteUnit(context.Context, *DeleteUnitRequest) (*DeleteUnitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUnit not implemented")
}
func (UnimplementedCentrumServiceServer) TakeControl(context.Context, *TakeControlRequest) (*TakeControlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TakeControl not implemented")
}
func (UnimplementedCentrumServiceServer) AssignUnit(context.Context, *AssignUnitRequest) (*AssignUnitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignUnit not implemented")
}
func (UnimplementedCentrumServiceServer) AssignDispatch(context.Context, *AssignDispatchRequest) (*AssignDispatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignDispatch not implemented")
}
func (UnimplementedCentrumServiceServer) CreateDispatch(context.Context, *CreateDispatchRequest) (*CreateDispatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDispatch not implemented")
}
func (UnimplementedCentrumServiceServer) UpdateDispatch(context.Context, *UpdateDispatchRequest) (*UpdateDispatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDispatch not implemented")
}
func (UnimplementedCentrumServiceServer) TakeDispatch(context.Context, *TakeDispatchRequest) (*TakeDispatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TakeDispatch not implemented")
}
func (UnimplementedCentrumServiceServer) UpdateDispatchStatus(context.Context, *UpdateDispatchStatusRequest) (*UpdateDispatchStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDispatchStatus not implemented")
}
func (UnimplementedCentrumServiceServer) UpdateUnitStatus(context.Context, *UpdateUnitStatusRequest) (*UpdateUnitStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUnitStatus not implemented")
}
func (UnimplementedCentrumServiceServer) ListDispatches(context.Context, *ListDispatchesRequest) (*ListDispatchesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDispatches not implemented")
}
func (UnimplementedCentrumServiceServer) ListUnits(context.Context, *ListUnitsRequest) (*ListUnitsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUnits not implemented")
}
func (UnimplementedCentrumServiceServer) ListUnitActivity(context.Context, *ListActivityRequest) (*ListUnitActivityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUnitActivity not implemented")
}
func (UnimplementedCentrumServiceServer) ListDispatchActivity(context.Context, *ListActivityRequest) (*ListDispatchActivityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDispatchActivity not implemented")
}
func (UnimplementedCentrumServiceServer) Stream(*StreamRequest, CentrumService_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedCentrumServiceServer) mustEmbedUnimplementedCentrumServiceServer() {}

// UnsafeCentrumServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CentrumServiceServer will
// result in compilation errors.
type UnsafeCentrumServiceServer interface {
	mustEmbedUnimplementedCentrumServiceServer()
}

func RegisterCentrumServiceServer(s grpc.ServiceRegistrar, srv CentrumServiceServer) {
	s.RegisterService(&CentrumService_ServiceDesc, srv)
}

func _CentrumService_GetSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).GetSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_GetSettings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).GetSettings(ctx, req.(*common.EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_UpdateSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(dispatch.Settings)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).UpdateSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_UpdateSettings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).UpdateSettings(ctx, req.(*dispatch.Settings))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_CreateOrUpdateUnit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrUpdateUnitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).CreateOrUpdateUnit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_CreateOrUpdateUnit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).CreateOrUpdateUnit(ctx, req.(*CreateOrUpdateUnitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_DeleteUnit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUnitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).DeleteUnit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_DeleteUnit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).DeleteUnit(ctx, req.(*DeleteUnitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_TakeControl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TakeControlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).TakeControl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_TakeControl_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).TakeControl(ctx, req.(*TakeControlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_AssignUnit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignUnitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).AssignUnit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_AssignUnit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).AssignUnit(ctx, req.(*AssignUnitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_AssignDispatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignDispatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).AssignDispatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_AssignDispatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).AssignDispatch(ctx, req.(*AssignDispatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_CreateDispatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDispatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).CreateDispatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_CreateDispatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).CreateDispatch(ctx, req.(*CreateDispatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_UpdateDispatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDispatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).UpdateDispatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_UpdateDispatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).UpdateDispatch(ctx, req.(*UpdateDispatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_TakeDispatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TakeDispatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).TakeDispatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_TakeDispatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).TakeDispatch(ctx, req.(*TakeDispatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_UpdateDispatchStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDispatchStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).UpdateDispatchStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_UpdateDispatchStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).UpdateDispatchStatus(ctx, req.(*UpdateDispatchStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_UpdateUnitStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUnitStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).UpdateUnitStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_UpdateUnitStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).UpdateUnitStatus(ctx, req.(*UpdateUnitStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_ListDispatches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDispatchesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).ListDispatches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_ListDispatches_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).ListDispatches(ctx, req.(*ListDispatchesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_ListUnits_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUnitsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).ListUnits(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_ListUnits_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).ListUnits(ctx, req.(*ListUnitsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_ListUnitActivity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListActivityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).ListUnitActivity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_ListUnitActivity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).ListUnitActivity(ctx, req.(*ListActivityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_ListDispatchActivity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListActivityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentrumServiceServer).ListDispatchActivity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentrumService_ListDispatchActivity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentrumServiceServer).ListDispatchActivity(ctx, req.(*ListActivityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentrumService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CentrumServiceServer).Stream(m, &centrumServiceStreamServer{stream})
}

type CentrumService_StreamServer interface {
	Send(*StreamResponse) error
	grpc.ServerStream
}

type centrumServiceStreamServer struct {
	grpc.ServerStream
}

func (x *centrumServiceStreamServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// CentrumService_ServiceDesc is the grpc.ServiceDesc for CentrumService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CentrumService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.centrum.CentrumService",
	HandlerType: (*CentrumServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSettings",
			Handler:    _CentrumService_GetSettings_Handler,
		},
		{
			MethodName: "UpdateSettings",
			Handler:    _CentrumService_UpdateSettings_Handler,
		},
		{
			MethodName: "CreateOrUpdateUnit",
			Handler:    _CentrumService_CreateOrUpdateUnit_Handler,
		},
		{
			MethodName: "DeleteUnit",
			Handler:    _CentrumService_DeleteUnit_Handler,
		},
		{
			MethodName: "TakeControl",
			Handler:    _CentrumService_TakeControl_Handler,
		},
		{
			MethodName: "AssignUnit",
			Handler:    _CentrumService_AssignUnit_Handler,
		},
		{
			MethodName: "AssignDispatch",
			Handler:    _CentrumService_AssignDispatch_Handler,
		},
		{
			MethodName: "CreateDispatch",
			Handler:    _CentrumService_CreateDispatch_Handler,
		},
		{
			MethodName: "UpdateDispatch",
			Handler:    _CentrumService_UpdateDispatch_Handler,
		},
		{
			MethodName: "TakeDispatch",
			Handler:    _CentrumService_TakeDispatch_Handler,
		},
		{
			MethodName: "UpdateDispatchStatus",
			Handler:    _CentrumService_UpdateDispatchStatus_Handler,
		},
		{
			MethodName: "UpdateUnitStatus",
			Handler:    _CentrumService_UpdateUnitStatus_Handler,
		},
		{
			MethodName: "ListDispatches",
			Handler:    _CentrumService_ListDispatches_Handler,
		},
		{
			MethodName: "ListUnits",
			Handler:    _CentrumService_ListUnits_Handler,
		},
		{
			MethodName: "ListUnitActivity",
			Handler:    _CentrumService_ListUnitActivity_Handler,
		},
		{
			MethodName: "ListDispatchActivity",
			Handler:    _CentrumService_ListDispatchActivity_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _CentrumService_Stream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "services/centrum/centrum.proto",
}
