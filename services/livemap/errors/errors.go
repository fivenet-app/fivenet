package errorslivemap

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrStreamFailed = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.LivemapService.ErrStreamFailed"}, nil)
	ErrMarkerFailed = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.LivemapService.ErrMarkerFailed"}, nil)
	ErrMarkerDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.LivemapService.ErrMarkerDenied"}, nil)
)
