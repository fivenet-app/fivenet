package errorscentrum

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrFailedQuery"},
		nil,
	)
	ErrNotPartOfDispatch = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotPartOfDispatch"},
		nil,
	)
	ErrNotPartOfUnit = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotPartOfUnit"},
		nil,
	)
	ErrUnitPermDenied = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrUnitPermDenied"},
		nil,
	)
	ErrNotOnDuty = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotOnDuty.content"},
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrNotOnDuty.title"},
	)
	ErrStaticUnit = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrStaticUnit"},
		nil,
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
		&common.I18NItem{Key: "errors.centrum.CentrumService.ErrDispatchNoJobs"},
		nil,
	)
)
