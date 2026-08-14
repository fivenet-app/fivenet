package errorscalendar

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrFailedQuery.content"},
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrFailedQuery.title"},
	)
	ErrNoPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrNoPerms.content"},
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrNoPerms.title"},
	)
	ErrOnePrivateCal = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrOnePrivateCal.content"},
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrOnePrivateCal.title"},
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

	ErrNoDiscordGuildID = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrNoDiscordGuildID.content"},
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrNoDiscordGuildID.title"},
	)
	ErrInvalidDiscordChannel = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrInvalidDiscordChannel.content"},
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrInvalidDiscordChannel.title"},
	)
	ErrInvalidReminderStep = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrInvalidReminderStep.content"},
		&common.I18NItem{Key: "errors.calendar.CalendarService.ErrInvalidReminderStep.title"},
	)
)
