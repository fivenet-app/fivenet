package errorscalendar

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery    = status.Error(codes.Internal, "errors.CalendarService.ErrFailedQuery")
	ErrNoPerms        = status.Error(codes.InvalidArgument, "errors.CalendarService.ErrNoPerms")
	ErrOnePrivateCal  = status.Error(codes.InvalidArgument, "errors.CalendarService.ErrOnePrivateCal")
	ErrCalendarClosed = status.Error(codes.InvalidArgument, "errors.CalendarService.ErrCalendarClosed.title;errors.CalendarService.ErrCalendarClosed.content")
	ErrEntryClosed    = status.Error(codes.InvalidArgument, "errors.CalendarService.ErrEntryClosed.title;errors.CalendarService.ErrEntryClosed.content")
)
