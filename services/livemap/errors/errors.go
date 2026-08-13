package errorslivemap

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrStreamFailed = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.livemap.LivemapService.ErrStreamFailed.content"},
		&common.I18NItem{Key: "errors.livemap.LivemapService.ErrStreamFailed.title"},
	)
	ErrMarkerFailed = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.livemap.LivemapService.ErrMarkerFailed.content"},
		&common.I18NItem{Key: "errors.livemap.LivemapService.ErrMarkerFailed.title"},
	)
	ErrMarkerDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.livemap.LivemapService.ErrMarkerDenied.content"},
		&common.I18NItem{Key: "errors.livemap.LivemapService.ErrMarkerDenied.title"},
	)
)
