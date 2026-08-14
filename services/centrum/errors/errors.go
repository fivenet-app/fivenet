package errorscentrum

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrFailedQuery.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrFailedQuery.title"},
	)
	ErrNotPartOfDispatch = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotPartOfDispatch.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotPartOfDispatch.title"},
	)
	ErrNotPartOfUnit = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotPartOfUnit.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotPartOfUnit.title"},
	)
	ErrUnitPermDenied = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrUnitPermDenied.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrUnitPermDenied.title"},
	)
	ErrDispatchJobPermDenied = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrDispatchJobPermDenied.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrDispatchJobPermDenied.title"},
	)
	ErrNotOnDuty = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotOnDuty.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotOnDuty.title"},
	)
	ErrStaticUnit = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrStaticUnit.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrStaticUnit.title"},
	)
	ErrDisabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrDisabled.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrDisabled.title"},
	)

	ErrModeForbidsAction = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrModeForbidsAction.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrModeForbidsAction.title"},
	)
	ErrDispatchAlreadyCompleted = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrDispatchAlreadyCompleted.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrDispatchAlreadyCompleted.title"},
	)
	ErrDispatchNoJobs = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrDispatchNoJobs.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrDispatchNoJobs.title"},
	)
)
