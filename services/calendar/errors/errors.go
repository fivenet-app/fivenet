package errorscalendar

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery    = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.CalendarService.ErrFailedQuery"}, nil)
	ErrNoPerms        = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CalendarService.ErrNoPerms"}, nil)
	ErrOnePrivateCal  = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CalendarService.ErrOnePrivateCal"}, nil)
	ErrCalendarClosed = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CalendarService.ErrCalendarClosed.content"}, &common.TranslateItem{Key: "errors.CalendarService.ErrCalendarClosed.title"})
	ErrEntryClosed    = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CalendarService.ErrEntryClosed.content"}, &common.TranslateItem{Key: "errors.CalendarService.ErrEntryClosed.title"})
)
