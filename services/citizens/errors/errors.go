package errorscitizens

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrFailedQuery.content"},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrFailedQuery.title"},
	)
	ErrJobGradeNoPermission = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrJobGradeNoPermission.content"},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrJobGradeNoPermission.title"},
	)
	ErrReasonRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrReasonRequired.content"},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrReasonRequired.title"},
	)
	ErrPropsWantedDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsWantedDenied.content"},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsWantedDenied.title"},
	)
	ErrPropsJobDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsJobDenied.content"},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsJobDenied.title"},
	)
	ErrPropsJobPublic = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsJobPublic.content"},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsJobPublic.title"},
	)
	ErrPropsJobInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsJobInvalid.content"},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsJobInvalid.title"},
	)
	ErrPropsTrafficPointsDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{
			Key: "errors.citizens.CitizensService.ErrPropsTrafficPointsDenied.content",
		},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsTrafficPointsDenied.title"},
	)
	ErrPropsMugshotDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsMugshotDenied.content"},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsMugshotDenied.title"},
	)
	ErrPropsLabelsDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsLabelsDenied.content"},
		&common.I18NItem{Key: "errors.citizens.CitizensService.ErrPropsLabelsDenied.title"},
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
