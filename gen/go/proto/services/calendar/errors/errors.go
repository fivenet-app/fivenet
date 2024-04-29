package errorscalendar

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery = status.Error(codes.Internal, "errors.CalendarService.ErrFailedQuery")
	ErrNoPerms     = status.Error(codes.InvalidArgument, "errors.CalendarService.ErrNoPerms")
)
