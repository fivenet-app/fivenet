package errorsvehicles

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrFailedQuery.content"},
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrFailedQuery.title"},
	)
	ErrReasonRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrReasonRequired.content"},
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrReasonRequired.title"},
	)
	ErrPropsWantedDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrPropsWantedDenied.content"},
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrPropsWantedDenied.title"},
	)
	ErrPropsPropertiesDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrPropsPropertiesDenied.content"},
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrPropsPropertiesDenied.title"},
	)
	ErrInvalidProperty = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrInvalidProperty.content"},
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrInvalidProperty.title"},
	)
	ErrSystemPropertyDisabled = common.NewI18nErr(
		codes.FailedPrecondition,
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrSystemPropertyDisabled.content"},
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrSystemPropertyDisabled.title"},
	)
)
