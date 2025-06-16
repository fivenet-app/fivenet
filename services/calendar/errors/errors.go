package errorscalendar

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery    = common.NewI18nErr(codes.Internal, &common.I18NItem{Key: "errors.CalendarService.ErrFailedQuery"}, nil)
	ErrNoPerms        = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.CalendarService.ErrNoPerms"}, nil)
	ErrOnePrivateCal  = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.CalendarService.ErrOnePrivateCal"}, nil)
	ErrCalendarClosed = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.CalendarService.ErrCalendarClosed.content"}, &common.I18NItem{Key: "errors.CalendarService.ErrCalendarClosed.title"})
	ErrEntryClosed    = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.CalendarService.ErrEntryClosed.content"}, &common.I18NItem{Key: "errors.CalendarService.ErrEntryClosed.title"})
)
