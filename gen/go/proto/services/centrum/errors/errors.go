package errorscentrum

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery       = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.CentrumService.ErrFailedQuery"}, nil)
	ErrNotPartOfDispatch = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CentrumService.ErrNotPartOfDispatch"}, nil)
	ErrNotPartOfUnit     = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CentrumService.ErrNotPartOfUnit"}, nil)
	ErrNotOnDuty         = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CentrumService.ErrNotOnDuty.content"}, &common.TranslateItem{Key: "errors.CentrumService.ErrNotOnDuty.title"})
	ErrStaticUnit        = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CentrumService.ErrStaticUnit"}, nil)

	ErrModeForbidsAction        = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CentrumService.ErrModeForbidsAction.content"}, &common.TranslateItem{Key: "errors.CentrumService.ErrModeForbidsAction.title"})
	ErrDispatchAlreadyCompleted = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CentrumService.ErrDispatchAlreadyCompleted.content"}, &common.TranslateItem{Key: "errors.CentrumService.ErrDispatchAlreadyCompleted.title"})
)
