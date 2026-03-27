package errorscalendar

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrFailedQuery"},
		nil,
	)
	ErrNoPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrNoPerms"},
		nil,
	)
	ErrOnePrivateCal = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrOnePrivateCal"},
		nil,
	)
	ErrCalendarClosed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrCalendarClosed.content"},
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrCalendarClosed.title"},
	)
	ErrEntryClosed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrEntryClosed.content"},
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrEntryClosed.title"},
	)
)
