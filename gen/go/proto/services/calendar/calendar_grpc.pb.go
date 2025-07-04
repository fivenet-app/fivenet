// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: services/calendar/calendar.proto

package calendar

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
	CalendarService_ListCalendars_FullMethodName               = "/services.calendar.CalendarService/ListCalendars"
	CalendarService_GetCalendar_FullMethodName                 = "/services.calendar.CalendarService/GetCalendar"
	CalendarService_CreateCalendar_FullMethodName              = "/services.calendar.CalendarService/CreateCalendar"
	CalendarService_UpdateCalendar_FullMethodName              = "/services.calendar.CalendarService/UpdateCalendar"
	CalendarService_DeleteCalendar_FullMethodName              = "/services.calendar.CalendarService/DeleteCalendar"
	CalendarService_ListCalendarEntries_FullMethodName         = "/services.calendar.CalendarService/ListCalendarEntries"
	CalendarService_GetUpcomingEntries_FullMethodName          = "/services.calendar.CalendarService/GetUpcomingEntries"
	CalendarService_GetCalendarEntry_FullMethodName            = "/services.calendar.CalendarService/GetCalendarEntry"
	CalendarService_CreateOrUpdateCalendarEntry_FullMethodName = "/services.calendar.CalendarService/CreateOrUpdateCalendarEntry"
	CalendarService_DeleteCalendarEntry_FullMethodName         = "/services.calendar.CalendarService/DeleteCalendarEntry"
	CalendarService_ShareCalendarEntry_FullMethodName          = "/services.calendar.CalendarService/ShareCalendarEntry"
	CalendarService_ListCalendarEntryRSVP_FullMethodName       = "/services.calendar.CalendarService/ListCalendarEntryRSVP"
	CalendarService_RSVPCalendarEntry_FullMethodName           = "/services.calendar.CalendarService/RSVPCalendarEntry"
	CalendarService_ListSubscriptions_FullMethodName           = "/services.calendar.CalendarService/ListSubscriptions"
	CalendarService_SubscribeToCalendar_FullMethodName         = "/services.calendar.CalendarService/SubscribeToCalendar"
)

// CalendarServiceClient is the client API for CalendarService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalendarServiceClient interface {
	// @perm: Name=Any
	ListCalendars(ctx context.Context, in *ListCalendarsRequest, opts ...grpc.CallOption) (*ListCalendarsResponse, error)
	// @perm: Name=Any
	GetCalendar(ctx context.Context, in *GetCalendarRequest, opts ...grpc.CallOption) (*GetCalendarResponse, error)
	// @perm: Attrs=Fields/StringList:[]string{"Job", "Public"}
	CreateCalendar(ctx context.Context, in *CreateCalendarRequest, opts ...grpc.CallOption) (*CreateCalendarResponse, error)
	// @perm: Name=Any
	UpdateCalendar(ctx context.Context, in *UpdateCalendarRequest, opts ...grpc.CallOption) (*UpdateCalendarResponse, error)
	// @perm: Name=Any
	DeleteCalendar(ctx context.Context, in *DeleteCalendarRequest, opts ...grpc.CallOption) (*DeleteCalendarResponse, error)
	// @perm: Name=Any
	ListCalendarEntries(ctx context.Context, in *ListCalendarEntriesRequest, opts ...grpc.CallOption) (*ListCalendarEntriesResponse, error)
	// @perm: Name=Any
	GetUpcomingEntries(ctx context.Context, in *GetUpcomingEntriesRequest, opts ...grpc.CallOption) (*GetUpcomingEntriesResponse, error)
	// @perm: Name=Any
	GetCalendarEntry(ctx context.Context, in *GetCalendarEntryRequest, opts ...grpc.CallOption) (*GetCalendarEntryResponse, error)
	// @perm: Name=Any
	CreateOrUpdateCalendarEntry(ctx context.Context, in *CreateOrUpdateCalendarEntryRequest, opts ...grpc.CallOption) (*CreateOrUpdateCalendarEntryResponse, error)
	// @perm: Name=Any
	DeleteCalendarEntry(ctx context.Context, in *DeleteCalendarEntryRequest, opts ...grpc.CallOption) (*DeleteCalendarEntryResponse, error)
	// @perm: Name=Any
	ShareCalendarEntry(ctx context.Context, in *ShareCalendarEntryRequest, opts ...grpc.CallOption) (*ShareCalendarEntryResponse, error)
	// @perm: Name=Any
	ListCalendarEntryRSVP(ctx context.Context, in *ListCalendarEntryRSVPRequest, opts ...grpc.CallOption) (*ListCalendarEntryRSVPResponse, error)
	// @perm: Name=Any
	RSVPCalendarEntry(ctx context.Context, in *RSVPCalendarEntryRequest, opts ...grpc.CallOption) (*RSVPCalendarEntryResponse, error)
	// @perm: Name=Any
	ListSubscriptions(ctx context.Context, in *ListSubscriptionsRequest, opts ...grpc.CallOption) (*ListSubscriptionsResponse, error)
	// @perm: Name=Any
	SubscribeToCalendar(ctx context.Context, in *SubscribeToCalendarRequest, opts ...grpc.CallOption) (*SubscribeToCalendarResponse, error)
}

type calendarServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCalendarServiceClient(cc grpc.ClientConnInterface) CalendarServiceClient {
	return &calendarServiceClient{cc}
}

func (c *calendarServiceClient) ListCalendars(ctx context.Context, in *ListCalendarsRequest, opts ...grpc.CallOption) (*ListCalendarsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCalendarsResponse)
	err := c.cc.Invoke(ctx, CalendarService_ListCalendars_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) GetCalendar(ctx context.Context, in *GetCalendarRequest, opts ...grpc.CallOption) (*GetCalendarResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCalendarResponse)
	err := c.cc.Invoke(ctx, CalendarService_GetCalendar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) CreateCalendar(ctx context.Context, in *CreateCalendarRequest, opts ...grpc.CallOption) (*CreateCalendarResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCalendarResponse)
	err := c.cc.Invoke(ctx, CalendarService_CreateCalendar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) UpdateCalendar(ctx context.Context, in *UpdateCalendarRequest, opts ...grpc.CallOption) (*UpdateCalendarResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateCalendarResponse)
	err := c.cc.Invoke(ctx, CalendarService_UpdateCalendar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) DeleteCalendar(ctx context.Context, in *DeleteCalendarRequest, opts ...grpc.CallOption) (*DeleteCalendarResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteCalendarResponse)
	err := c.cc.Invoke(ctx, CalendarService_DeleteCalendar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) ListCalendarEntries(ctx context.Context, in *ListCalendarEntriesRequest, opts ...grpc.CallOption) (*ListCalendarEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCalendarEntriesResponse)
	err := c.cc.Invoke(ctx, CalendarService_ListCalendarEntries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) GetUpcomingEntries(ctx context.Context, in *GetUpcomingEntriesRequest, opts ...grpc.CallOption) (*GetUpcomingEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUpcomingEntriesResponse)
	err := c.cc.Invoke(ctx, CalendarService_GetUpcomingEntries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) GetCalendarEntry(ctx context.Context, in *GetCalendarEntryRequest, opts ...grpc.CallOption) (*GetCalendarEntryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCalendarEntryResponse)
	err := c.cc.Invoke(ctx, CalendarService_GetCalendarEntry_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) CreateOrUpdateCalendarEntry(ctx context.Context, in *CreateOrUpdateCalendarEntryRequest, opts ...grpc.CallOption) (*CreateOrUpdateCalendarEntryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateOrUpdateCalendarEntryResponse)
	err := c.cc.Invoke(ctx, CalendarService_CreateOrUpdateCalendarEntry_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) DeleteCalendarEntry(ctx context.Context, in *DeleteCalendarEntryRequest, opts ...grpc.CallOption) (*DeleteCalendarEntryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteCalendarEntryResponse)
	err := c.cc.Invoke(ctx, CalendarService_DeleteCalendarEntry_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) ShareCalendarEntry(ctx context.Context, in *ShareCalendarEntryRequest, opts ...grpc.CallOption) (*ShareCalendarEntryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShareCalendarEntryResponse)
	err := c.cc.Invoke(ctx, CalendarService_ShareCalendarEntry_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) ListCalendarEntryRSVP(ctx context.Context, in *ListCalendarEntryRSVPRequest, opts ...grpc.CallOption) (*ListCalendarEntryRSVPResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCalendarEntryRSVPResponse)
	err := c.cc.Invoke(ctx, CalendarService_ListCalendarEntryRSVP_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) RSVPCalendarEntry(ctx context.Context, in *RSVPCalendarEntryRequest, opts ...grpc.CallOption) (*RSVPCalendarEntryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RSVPCalendarEntryResponse)
	err := c.cc.Invoke(ctx, CalendarService_RSVPCalendarEntry_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) ListSubscriptions(ctx context.Context, in *ListSubscriptionsRequest, opts ...grpc.CallOption) (*ListSubscriptionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListSubscriptionsResponse)
	err := c.cc.Invoke(ctx, CalendarService_ListSubscriptions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarServiceClient) SubscribeToCalendar(ctx context.Context, in *SubscribeToCalendarRequest, opts ...grpc.CallOption) (*SubscribeToCalendarResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SubscribeToCalendarResponse)
	err := c.cc.Invoke(ctx, CalendarService_SubscribeToCalendar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarServiceServer is the server API for CalendarService service.
// All implementations must embed UnimplementedCalendarServiceServer
// for forward compatibility.
type CalendarServiceServer interface {
	// @perm: Name=Any
	ListCalendars(context.Context, *ListCalendarsRequest) (*ListCalendarsResponse, error)
	// @perm: Name=Any
	GetCalendar(context.Context, *GetCalendarRequest) (*GetCalendarResponse, error)
	// @perm: Attrs=Fields/StringList:[]string{"Job", "Public"}
	CreateCalendar(context.Context, *CreateCalendarRequest) (*CreateCalendarResponse, error)
	// @perm: Name=Any
	UpdateCalendar(context.Context, *UpdateCalendarRequest) (*UpdateCalendarResponse, error)
	// @perm: Name=Any
	DeleteCalendar(context.Context, *DeleteCalendarRequest) (*DeleteCalendarResponse, error)
	// @perm: Name=Any
	ListCalendarEntries(context.Context, *ListCalendarEntriesRequest) (*ListCalendarEntriesResponse, error)
	// @perm: Name=Any
	GetUpcomingEntries(context.Context, *GetUpcomingEntriesRequest) (*GetUpcomingEntriesResponse, error)
	// @perm: Name=Any
	GetCalendarEntry(context.Context, *GetCalendarEntryRequest) (*GetCalendarEntryResponse, error)
	// @perm: Name=Any
	CreateOrUpdateCalendarEntry(context.Context, *CreateOrUpdateCalendarEntryRequest) (*CreateOrUpdateCalendarEntryResponse, error)
	// @perm: Name=Any
	DeleteCalendarEntry(context.Context, *DeleteCalendarEntryRequest) (*DeleteCalendarEntryResponse, error)
	// @perm: Name=Any
	ShareCalendarEntry(context.Context, *ShareCalendarEntryRequest) (*ShareCalendarEntryResponse, error)
	// @perm: Name=Any
	ListCalendarEntryRSVP(context.Context, *ListCalendarEntryRSVPRequest) (*ListCalendarEntryRSVPResponse, error)
	// @perm: Name=Any
	RSVPCalendarEntry(context.Context, *RSVPCalendarEntryRequest) (*RSVPCalendarEntryResponse, error)
	// @perm: Name=Any
	ListSubscriptions(context.Context, *ListSubscriptionsRequest) (*ListSubscriptionsResponse, error)
	// @perm: Name=Any
	SubscribeToCalendar(context.Context, *SubscribeToCalendarRequest) (*SubscribeToCalendarResponse, error)
	mustEmbedUnimplementedCalendarServiceServer()
}

// UnimplementedCalendarServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCalendarServiceServer struct{}

func (UnimplementedCalendarServiceServer) ListCalendars(context.Context, *ListCalendarsRequest) (*ListCalendarsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCalendars not implemented")
}
func (UnimplementedCalendarServiceServer) GetCalendar(context.Context, *GetCalendarRequest) (*GetCalendarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCalendar not implemented")
}
func (UnimplementedCalendarServiceServer) CreateCalendar(context.Context, *CreateCalendarRequest) (*CreateCalendarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCalendar not implemented")
}
func (UnimplementedCalendarServiceServer) UpdateCalendar(context.Context, *UpdateCalendarRequest) (*UpdateCalendarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCalendar not implemented")
}
func (UnimplementedCalendarServiceServer) DeleteCalendar(context.Context, *DeleteCalendarRequest) (*DeleteCalendarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCalendar not implemented")
}
func (UnimplementedCalendarServiceServer) ListCalendarEntries(context.Context, *ListCalendarEntriesRequest) (*ListCalendarEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCalendarEntries not implemented")
}
func (UnimplementedCalendarServiceServer) GetUpcomingEntries(context.Context, *GetUpcomingEntriesRequest) (*GetUpcomingEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUpcomingEntries not implemented")
}
func (UnimplementedCalendarServiceServer) GetCalendarEntry(context.Context, *GetCalendarEntryRequest) (*GetCalendarEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCalendarEntry not implemented")
}
func (UnimplementedCalendarServiceServer) CreateOrUpdateCalendarEntry(context.Context, *CreateOrUpdateCalendarEntryRequest) (*CreateOrUpdateCalendarEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateCalendarEntry not implemented")
}
func (UnimplementedCalendarServiceServer) DeleteCalendarEntry(context.Context, *DeleteCalendarEntryRequest) (*DeleteCalendarEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCalendarEntry not implemented")
}
func (UnimplementedCalendarServiceServer) ShareCalendarEntry(context.Context, *ShareCalendarEntryRequest) (*ShareCalendarEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShareCalendarEntry not implemented")
}
func (UnimplementedCalendarServiceServer) ListCalendarEntryRSVP(context.Context, *ListCalendarEntryRSVPRequest) (*ListCalendarEntryRSVPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCalendarEntryRSVP not implemented")
}
func (UnimplementedCalendarServiceServer) RSVPCalendarEntry(context.Context, *RSVPCalendarEntryRequest) (*RSVPCalendarEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RSVPCalendarEntry not implemented")
}
func (UnimplementedCalendarServiceServer) ListSubscriptions(context.Context, *ListSubscriptionsRequest) (*ListSubscriptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSubscriptions not implemented")
}
func (UnimplementedCalendarServiceServer) SubscribeToCalendar(context.Context, *SubscribeToCalendarRequest) (*SubscribeToCalendarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeToCalendar not implemented")
}
func (UnimplementedCalendarServiceServer) mustEmbedUnimplementedCalendarServiceServer() {}
func (UnimplementedCalendarServiceServer) testEmbeddedByValue()                         {}

// UnsafeCalendarServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalendarServiceServer will
// result in compilation errors.
type UnsafeCalendarServiceServer interface {
	mustEmbedUnimplementedCalendarServiceServer()
}

func RegisterCalendarServiceServer(s grpc.ServiceRegistrar, srv CalendarServiceServer) {
	// If the following call pancis, it indicates UnimplementedCalendarServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CalendarService_ServiceDesc, srv)
}

func _CalendarService_ListCalendars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCalendarsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).ListCalendars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_ListCalendars_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).ListCalendars(ctx, req.(*ListCalendarsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_GetCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).GetCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_GetCalendar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).GetCalendar(ctx, req.(*GetCalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_CreateCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).CreateCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_CreateCalendar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).CreateCalendar(ctx, req.(*CreateCalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_UpdateCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).UpdateCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_UpdateCalendar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).UpdateCalendar(ctx, req.(*UpdateCalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_DeleteCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).DeleteCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_DeleteCalendar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).DeleteCalendar(ctx, req.(*DeleteCalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_ListCalendarEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCalendarEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).ListCalendarEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_ListCalendarEntries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).ListCalendarEntries(ctx, req.(*ListCalendarEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_GetUpcomingEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUpcomingEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).GetUpcomingEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_GetUpcomingEntries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).GetUpcomingEntries(ctx, req.(*GetUpcomingEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_GetCalendarEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCalendarEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).GetCalendarEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_GetCalendarEntry_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).GetCalendarEntry(ctx, req.(*GetCalendarEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_CreateOrUpdateCalendarEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrUpdateCalendarEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).CreateOrUpdateCalendarEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_CreateOrUpdateCalendarEntry_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).CreateOrUpdateCalendarEntry(ctx, req.(*CreateOrUpdateCalendarEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_DeleteCalendarEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCalendarEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).DeleteCalendarEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_DeleteCalendarEntry_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).DeleteCalendarEntry(ctx, req.(*DeleteCalendarEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_ShareCalendarEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareCalendarEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).ShareCalendarEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_ShareCalendarEntry_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).ShareCalendarEntry(ctx, req.(*ShareCalendarEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_ListCalendarEntryRSVP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCalendarEntryRSVPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).ListCalendarEntryRSVP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_ListCalendarEntryRSVP_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).ListCalendarEntryRSVP(ctx, req.(*ListCalendarEntryRSVPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_RSVPCalendarEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RSVPCalendarEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).RSVPCalendarEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_RSVPCalendarEntry_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).RSVPCalendarEntry(ctx, req.(*RSVPCalendarEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_ListSubscriptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSubscriptionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).ListSubscriptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_ListSubscriptions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).ListSubscriptions(ctx, req.(*ListSubscriptionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalendarService_SubscribeToCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeToCalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServiceServer).SubscribeToCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CalendarService_SubscribeToCalendar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServiceServer).SubscribeToCalendar(ctx, req.(*SubscribeToCalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CalendarService_ServiceDesc is the grpc.ServiceDesc for CalendarService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CalendarService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.calendar.CalendarService",
	HandlerType: (*CalendarServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListCalendars",
			Handler:    _CalendarService_ListCalendars_Handler,
		},
		{
			MethodName: "GetCalendar",
			Handler:    _CalendarService_GetCalendar_Handler,
		},
		{
			MethodName: "CreateCalendar",
			Handler:    _CalendarService_CreateCalendar_Handler,
		},
		{
			MethodName: "UpdateCalendar",
			Handler:    _CalendarService_UpdateCalendar_Handler,
		},
		{
			MethodName: "DeleteCalendar",
			Handler:    _CalendarService_DeleteCalendar_Handler,
		},
		{
			MethodName: "ListCalendarEntries",
			Handler:    _CalendarService_ListCalendarEntries_Handler,
		},
		{
			MethodName: "GetUpcomingEntries",
			Handler:    _CalendarService_GetUpcomingEntries_Handler,
		},
		{
			MethodName: "GetCalendarEntry",
			Handler:    _CalendarService_GetCalendarEntry_Handler,
		},
		{
			MethodName: "CreateOrUpdateCalendarEntry",
			Handler:    _CalendarService_CreateOrUpdateCalendarEntry_Handler,
		},
		{
			MethodName: "DeleteCalendarEntry",
			Handler:    _CalendarService_DeleteCalendarEntry_Handler,
		},
		{
			MethodName: "ShareCalendarEntry",
			Handler:    _CalendarService_ShareCalendarEntry_Handler,
		},
		{
			MethodName: "ListCalendarEntryRSVP",
			Handler:    _CalendarService_ListCalendarEntryRSVP_Handler,
		},
		{
			MethodName: "RSVPCalendarEntry",
			Handler:    _CalendarService_RSVPCalendarEntry_Handler,
		},
		{
			MethodName: "ListSubscriptions",
			Handler:    _CalendarService_ListSubscriptions_Handler,
		},
		{
			MethodName: "SubscribeToCalendar",
			Handler:    _CalendarService_SubscribeToCalendar_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/calendar/calendar.proto",
}
