package errorsjobs

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrFailedQuery"},
		nil,
	)
	ErrPropsAbsenceDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrPropsAbsenceDenied"},
		nil,
	)
	ErrPropsNoteDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrPropsNoteDenied"},
		nil,
	)
	ErrPropsLabelsDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrPropsLabelsDenied"},
		nil,
	)
	ErrPropsNameDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrPropsNameDenied"},
		nil,
	)
	ErrReasonRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrReasonRequired"},
		nil,
	)
	ErrNotFoundOrNoPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrNotFoundOrNoPerms"},
		nil,
	)

	ErrAbsenceBeginOutOfRange = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrAbsenceBeginOutOfRange"},
		nil,
	)
	ErrAbsenceEndOutOfRange = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrAbsenceEndOutOfRange"},
		nil,
	)

	ErrTimeclockOutOfRange = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrTimeclockOutOfRange"},
		nil,
	)

	ErrLabelsNoPerms = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrLabelsNoPerms.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrLabelsNoPerms.title"},
	)
)
