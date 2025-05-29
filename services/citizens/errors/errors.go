package errorscitizens

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery              = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.CitizensService.ErrFailedQuery"}, nil)
	ErrJobGradeNoPermission     = common.I18nErr(codes.NotFound, &common.TranslateItem{Key: "errors.CitizensService.ErrJobGradeNoPermission"}, nil)
	ErrReasonRequired           = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CitizensService.ErrReasonRequired"}, nil)
	ErrPropsWantedDenied        = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizensService.ErrPropsWantedDenied"}, nil)
	ErrPropsJobDenied           = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizensService.ErrPropsJobDenied"}, nil)
	ErrPropsJobPublic           = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CitizensService.ErrPropsJobPublic"}, nil)
	ErrPropsJobInvalid          = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CitizensService.ErrPropsJobInvalid"}, nil)
	ErrPropsTrafficPointsDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizensService.ErrPropsTrafficPointsDenied"}, nil)
	ErrPropsMugshotDenied       = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizensService.ErrPropsMugshotDenied"}, nil)
	ErrPropsLabelsDenied        = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizensService.ErrPropsLabelsDenied"}, nil)
)
