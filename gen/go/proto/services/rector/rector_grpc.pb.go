// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: services/rector/rector.proto

package rector

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
	RectorService_GetJobProps_FullMethodName      = "/services.rector.RectorService/GetJobProps"
	RectorService_SetJobProps_FullMethodName      = "/services.rector.RectorService/SetJobProps"
	RectorService_GetRoles_FullMethodName         = "/services.rector.RectorService/GetRoles"
	RectorService_GetRole_FullMethodName          = "/services.rector.RectorService/GetRole"
	RectorService_CreateRole_FullMethodName       = "/services.rector.RectorService/CreateRole"
	RectorService_DeleteRole_FullMethodName       = "/services.rector.RectorService/DeleteRole"
	RectorService_UpdateRolePerms_FullMethodName  = "/services.rector.RectorService/UpdateRolePerms"
	RectorService_GetPermissions_FullMethodName   = "/services.rector.RectorService/GetPermissions"
	RectorService_ViewAuditLog_FullMethodName     = "/services.rector.RectorService/ViewAuditLog"
	RectorService_UpdateRoleLimits_FullMethodName = "/services.rector.RectorService/UpdateRoleLimits"
	RectorService_DeleteFaction_FullMethodName    = "/services.rector.RectorService/DeleteFaction"
)

// RectorServiceClient is the client API for RectorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RectorServiceClient interface {
	// @perm
	GetJobProps(ctx context.Context, in *GetJobPropsRequest, opts ...grpc.CallOption) (*GetJobPropsResponse, error)
	// @perm
	SetJobProps(ctx context.Context, in *SetJobPropsRequest, opts ...grpc.CallOption) (*SetJobPropsResponse, error)
	// @perm
	GetRoles(ctx context.Context, in *GetRolesRequest, opts ...grpc.CallOption) (*GetRolesResponse, error)
	// @perm: Name=GetRoles
	GetRole(ctx context.Context, in *GetRoleRequest, opts ...grpc.CallOption) (*GetRoleResponse, error)
	// @perm
	CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error)
	// @perm
	DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error)
	// @perm
	UpdateRolePerms(ctx context.Context, in *UpdateRolePermsRequest, opts ...grpc.CallOption) (*UpdateRolePermsResponse, error)
	// @perm: Name=GetRoles
	GetPermissions(ctx context.Context, in *GetPermissionsRequest, opts ...grpc.CallOption) (*GetPermissionsResponse, error)
	// @perm
	ViewAuditLog(ctx context.Context, in *ViewAuditLogRequest, opts ...grpc.CallOption) (*ViewAuditLogResponse, error)
	// @perm: Name=SuperUser
	UpdateRoleLimits(ctx context.Context, in *UpdateRoleLimitsRequest, opts ...grpc.CallOption) (*UpdateRoleLimitsResponse, error)
	// @perm: Name=SuperUser
	DeleteFaction(ctx context.Context, in *DeleteFactionRequest, opts ...grpc.CallOption) (*DeleteFactionResponse, error)
}

type rectorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRectorServiceClient(cc grpc.ClientConnInterface) RectorServiceClient {
	return &rectorServiceClient{cc}
}

func (c *rectorServiceClient) GetJobProps(ctx context.Context, in *GetJobPropsRequest, opts ...grpc.CallOption) (*GetJobPropsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetJobPropsResponse)
	err := c.cc.Invoke(ctx, RectorService_GetJobProps_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) SetJobProps(ctx context.Context, in *SetJobPropsRequest, opts ...grpc.CallOption) (*SetJobPropsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetJobPropsResponse)
	err := c.cc.Invoke(ctx, RectorService_SetJobProps_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) GetRoles(ctx context.Context, in *GetRolesRequest, opts ...grpc.CallOption) (*GetRolesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRolesResponse)
	err := c.cc.Invoke(ctx, RectorService_GetRoles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) GetRole(ctx context.Context, in *GetRoleRequest, opts ...grpc.CallOption) (*GetRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoleResponse)
	err := c.cc.Invoke(ctx, RectorService_GetRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateRoleResponse)
	err := c.cc.Invoke(ctx, RectorService_CreateRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteRoleResponse)
	err := c.cc.Invoke(ctx, RectorService_DeleteRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) UpdateRolePerms(ctx context.Context, in *UpdateRolePermsRequest, opts ...grpc.CallOption) (*UpdateRolePermsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateRolePermsResponse)
	err := c.cc.Invoke(ctx, RectorService_UpdateRolePerms_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) GetPermissions(ctx context.Context, in *GetPermissionsRequest, opts ...grpc.CallOption) (*GetPermissionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPermissionsResponse)
	err := c.cc.Invoke(ctx, RectorService_GetPermissions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) ViewAuditLog(ctx context.Context, in *ViewAuditLogRequest, opts ...grpc.CallOption) (*ViewAuditLogResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ViewAuditLogResponse)
	err := c.cc.Invoke(ctx, RectorService_ViewAuditLog_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) UpdateRoleLimits(ctx context.Context, in *UpdateRoleLimitsRequest, opts ...grpc.CallOption) (*UpdateRoleLimitsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateRoleLimitsResponse)
	err := c.cc.Invoke(ctx, RectorService_UpdateRoleLimits_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rectorServiceClient) DeleteFaction(ctx context.Context, in *DeleteFactionRequest, opts ...grpc.CallOption) (*DeleteFactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteFactionResponse)
	err := c.cc.Invoke(ctx, RectorService_DeleteFaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RectorServiceServer is the server API for RectorService service.
// All implementations must embed UnimplementedRectorServiceServer
// for forward compatibility.
type RectorServiceServer interface {
	// @perm
	GetJobProps(context.Context, *GetJobPropsRequest) (*GetJobPropsResponse, error)
	// @perm
	SetJobProps(context.Context, *SetJobPropsRequest) (*SetJobPropsResponse, error)
	// @perm
	GetRoles(context.Context, *GetRolesRequest) (*GetRolesResponse, error)
	// @perm: Name=GetRoles
	GetRole(context.Context, *GetRoleRequest) (*GetRoleResponse, error)
	// @perm
	CreateRole(context.Context, *CreateRoleRequest) (*CreateRoleResponse, error)
	// @perm
	DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleResponse, error)
	// @perm
	UpdateRolePerms(context.Context, *UpdateRolePermsRequest) (*UpdateRolePermsResponse, error)
	// @perm: Name=GetRoles
	GetPermissions(context.Context, *GetPermissionsRequest) (*GetPermissionsResponse, error)
	// @perm
	ViewAuditLog(context.Context, *ViewAuditLogRequest) (*ViewAuditLogResponse, error)
	// @perm: Name=SuperUser
	UpdateRoleLimits(context.Context, *UpdateRoleLimitsRequest) (*UpdateRoleLimitsResponse, error)
	// @perm: Name=SuperUser
	DeleteFaction(context.Context, *DeleteFactionRequest) (*DeleteFactionResponse, error)
	mustEmbedUnimplementedRectorServiceServer()
}

// UnimplementedRectorServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRectorServiceServer struct{}

func (UnimplementedRectorServiceServer) GetJobProps(context.Context, *GetJobPropsRequest) (*GetJobPropsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJobProps not implemented")
}
func (UnimplementedRectorServiceServer) SetJobProps(context.Context, *SetJobPropsRequest) (*SetJobPropsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetJobProps not implemented")
}
func (UnimplementedRectorServiceServer) GetRoles(context.Context, *GetRolesRequest) (*GetRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoles not implemented")
}
func (UnimplementedRectorServiceServer) GetRole(context.Context, *GetRoleRequest) (*GetRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRole not implemented")
}
func (UnimplementedRectorServiceServer) CreateRole(context.Context, *CreateRoleRequest) (*CreateRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRole not implemented")
}
func (UnimplementedRectorServiceServer) DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRole not implemented")
}
func (UnimplementedRectorServiceServer) UpdateRolePerms(context.Context, *UpdateRolePermsRequest) (*UpdateRolePermsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRolePerms not implemented")
}
func (UnimplementedRectorServiceServer) GetPermissions(context.Context, *GetPermissionsRequest) (*GetPermissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPermissions not implemented")
}
func (UnimplementedRectorServiceServer) ViewAuditLog(context.Context, *ViewAuditLogRequest) (*ViewAuditLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ViewAuditLog not implemented")
}
func (UnimplementedRectorServiceServer) UpdateRoleLimits(context.Context, *UpdateRoleLimitsRequest) (*UpdateRoleLimitsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRoleLimits not implemented")
}
func (UnimplementedRectorServiceServer) DeleteFaction(context.Context, *DeleteFactionRequest) (*DeleteFactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFaction not implemented")
}
func (UnimplementedRectorServiceServer) mustEmbedUnimplementedRectorServiceServer() {}
func (UnimplementedRectorServiceServer) testEmbeddedByValue()                       {}

// UnsafeRectorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RectorServiceServer will
// result in compilation errors.
type UnsafeRectorServiceServer interface {
	mustEmbedUnimplementedRectorServiceServer()
}

func RegisterRectorServiceServer(s grpc.ServiceRegistrar, srv RectorServiceServer) {
	// If the following call pancis, it indicates UnimplementedRectorServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RectorService_ServiceDesc, srv)
}

func _RectorService_GetJobProps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJobPropsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).GetJobProps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_GetJobProps_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).GetJobProps(ctx, req.(*GetJobPropsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_SetJobProps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetJobPropsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).SetJobProps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_SetJobProps_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).SetJobProps(ctx, req.(*SetJobPropsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_GetRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).GetRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_GetRoles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).GetRoles(ctx, req.(*GetRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_GetRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).GetRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_GetRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).GetRole(ctx, req.(*GetRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_CreateRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).CreateRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_CreateRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).CreateRole(ctx, req.(*CreateRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_DeleteRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).DeleteRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_DeleteRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).DeleteRole(ctx, req.(*DeleteRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_UpdateRolePerms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRolePermsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).UpdateRolePerms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_UpdateRolePerms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).UpdateRolePerms(ctx, req.(*UpdateRolePermsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_GetPermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPermissionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).GetPermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_GetPermissions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).GetPermissions(ctx, req.(*GetPermissionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_ViewAuditLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ViewAuditLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).ViewAuditLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_ViewAuditLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).ViewAuditLog(ctx, req.(*ViewAuditLogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_UpdateRoleLimits_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRoleLimitsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).UpdateRoleLimits(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_UpdateRoleLimits_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).UpdateRoleLimits(ctx, req.(*UpdateRoleLimitsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RectorService_DeleteFaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RectorServiceServer).DeleteFaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RectorService_DeleteFaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RectorServiceServer).DeleteFaction(ctx, req.(*DeleteFactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RectorService_ServiceDesc is the grpc.ServiceDesc for RectorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RectorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.rector.RectorService",
	HandlerType: (*RectorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetJobProps",
			Handler:    _RectorService_GetJobProps_Handler,
		},
		{
			MethodName: "SetJobProps",
			Handler:    _RectorService_SetJobProps_Handler,
		},
		{
			MethodName: "GetRoles",
			Handler:    _RectorService_GetRoles_Handler,
		},
		{
			MethodName: "GetRole",
			Handler:    _RectorService_GetRole_Handler,
		},
		{
			MethodName: "CreateRole",
			Handler:    _RectorService_CreateRole_Handler,
		},
		{
			MethodName: "DeleteRole",
			Handler:    _RectorService_DeleteRole_Handler,
		},
		{
			MethodName: "UpdateRolePerms",
			Handler:    _RectorService_UpdateRolePerms_Handler,
		},
		{
			MethodName: "GetPermissions",
			Handler:    _RectorService_GetPermissions_Handler,
		},
		{
			MethodName: "ViewAuditLog",
			Handler:    _RectorService_ViewAuditLog_Handler,
		},
		{
			MethodName: "UpdateRoleLimits",
			Handler:    _RectorService_UpdateRoleLimits_Handler,
		},
		{
			MethodName: "DeleteFaction",
			Handler:    _RectorService_DeleteFaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/rector/rector.proto",
}
