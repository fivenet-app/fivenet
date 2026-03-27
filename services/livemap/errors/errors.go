package errorslivemap

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrStreamFailed = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.livemap.LivemapService.ErrStreamFailed"},
		nil,
	)
	ErrMarkerFailed = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.livemap.LivemapService.ErrMarkerFailed"},
		nil,
	)
	ErrMarkerDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.livemap.LivemapService.ErrMarkerDenied"},
		nil,
	)
)
