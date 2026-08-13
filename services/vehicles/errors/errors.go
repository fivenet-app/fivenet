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
	ErrPropsWantedDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrPropsWantedDenied.content"},
		&common.I18NItem{Key: "errors.vehicles.VehiclesService.ErrPropsWantedDenied.title"},
	)
)
