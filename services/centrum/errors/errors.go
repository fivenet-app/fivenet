package errorscentrum

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.CentrumService.ErrFailedQuery"},
		nil,
	)
	ErrNotPartOfDispatch = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.CentrumService.ErrNotPartOfDispatch"},
		nil,
	)
	ErrNotPartOfUnit = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.CentrumService.ErrNotPartOfUnit"},
		nil,
	)
	ErrUnitPermDenied = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.CentrumService.ErrUnitPermDenied"},
		nil,
	)
	ErrNotOnDuty = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.CentrumService.ErrNotOnDuty.content"},
		&common.I18NItem{Key: "errors.CentrumService.ErrNotOnDuty.title"},
	)
	ErrStaticUnit = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.CentrumService.ErrStaticUnit"},
		nil,
	)
	ErrDisabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.CentrumService.ErrDisabled.content"},
		&common.I18NItem{Key: "errors.CentrumService.ErrDisabled.title"},
	)

	ErrModeForbidsAction = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.CentrumService.ErrModeForbidsAction.content"},
		&common.I18NItem{Key: "errors.CentrumService.ErrModeForbidsAction.title"},
	)
	ErrDispatchAlreadyCompleted = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.CentrumService.ErrDispatchAlreadyCompleted.content"},
		&common.I18NItem{Key: "errors.CentrumService.ErrDispatchAlreadyCompleted.title"},
	)
	ErrDispatchNoJobs = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.CentrumService.ErrDispatchNoJobs"},
		nil,
	)
)
