package errorsjobs

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.JobsService.ErrFailedQuery"},
		nil,
	)
	ErrPropsAbsenceDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.JobsService.ErrPropsAbsenceDenied"},
		nil,
	)
	ErrPropsNoteDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.JobsService.ErrPropsNoteDenied"},
		nil,
	)
	ErrPropsLabelsDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.JobsService.ErrPropsLabelsDenied"},
		nil,
	)
	ErrPropsNameDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.JobsService.ErrPropsNameDenied"},
		nil,
	)
	ErrReasonRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.JobsService.ErrReasonRequired"},
		nil,
	)
	ErrNotFoundOrNoPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.JobsService.ErrNotFoundOrNoPerms"},
		nil,
	)

	ErrAbsenceBeginOutOfRange = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.JobsService.ErrAbsenceBeginOutOfRange"},
		nil,
	)
	ErrAbsenceEndOutOfRange = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.JobsService.ErrAbsenceEndOutOfRange"},
		nil,
	)

	ErrTimeclockOutOfRange = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.JobsService.ErrTimeclockOutOfRange"},
		nil,
	)

	ErrLabelsNoPerms = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.JobsService.ErrLabelsNoPerms.content"},
		&common.I18NItem{Key: "errors.JobsService.ErrLabelsNoPerms.title"},
	)
)
