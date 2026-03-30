package errorscitizens

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrFailedQuery"},
		nil,
	)
	ErrJobGradeNoPermission = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrJobGradeNoPermission"},
		nil,
	)
	ErrReasonRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrReasonRequired"},
		nil,
	)
	ErrPropsWantedDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsWantedDenied"},
		nil,
	)
	ErrPropsJobDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsJobDenied"},
		nil,
	)
	ErrPropsJobPublic = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsJobPublic"},
		nil,
	)
	ErrPropsJobInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsJobInvalid"},
		nil,
	)
	ErrPropsTrafficPointsDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsTrafficPointsDenied"},
		nil,
	)
	ErrPropsMugshotDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsMugshotDenied"},
		nil,
	)
	ErrPropsLabelsDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsLabelsDenied"},
		nil,
	)
	ErrCitizenNotFound = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrCitizenNotFound"},
		nil,
	)
	ErrLabelNotFound = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.citizens.LabelsService.ErrLabelNotFound"},
		nil,
	)
)
