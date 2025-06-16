package errorslivemap

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrStreamFailed = common.NewI18nErr(codes.Internal, &common.I18NItem{Key: "errors.LivemapService.ErrStreamFailed"}, nil)
	ErrMarkerFailed = common.NewI18nErr(codes.Internal, &common.I18NItem{Key: "errors.LivemapService.ErrMarkerFailed"}, nil)
	ErrMarkerDenied = common.NewI18nErr(codes.PermissionDenied, &common.I18NItem{Key: "errors.LivemapService.ErrMarkerDenied"}, nil)
)
