package errorsjobs

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrFailedQuery.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrFailedQuery.title"},
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
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrReasonRequired.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrReasonRequired.title"},
	)
	ErrNotFoundOrNoPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrNotFoundOrNoPerms.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrNotFoundOrNoPerms.title"},
	)
	ErrGroupMemberRulesRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrGroupMemberRulesRequired.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrGroupMemberRulesRequired.title"},
	)
	ErrGroupPolicyViolation = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrGroupPolicyViolation.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrGroupPolicyViolation.title"},
	)

	ErrAbsenceBeginOutOfRange = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrAbsenceBeginOutOfRange.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrAbsenceBeginOutOfRange.title"},
	)
	ErrAbsenceEndOutOfRange = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrAbsenceEndOutOfRange.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrAbsenceEndOutOfRange.title"},
	)

	ErrTimeclockOutOfRange = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrTimeclockOutOfRange.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrTimeclockOutOfRange.title"},
	)

	ErrLabelsNoPerms = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrLabelsNoPerms.content"},
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrLabelsNoPerms.title"},
	)
	ErrLabelNotFound = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.jobs.JobsService.ErrLabelNotFound"},
		nil,
	)
)
