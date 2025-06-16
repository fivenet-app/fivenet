package errorscitizens

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery              = common.NewI18nErr(codes.Internal, &common.I18NItem{Key: "errors.CitizensService.ErrFailedQuery"}, nil)
	ErrJobGradeNoPermission     = common.NewI18nErr(codes.NotFound, &common.I18NItem{Key: "errors.CitizensService.ErrJobGradeNoPermission"}, nil)
	ErrReasonRequired           = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.CitizensService.ErrReasonRequired"}, nil)
	ErrPropsWantedDenied        = common.NewI18nErr(codes.PermissionDenied, &common.I18NItem{Key: "errors.CitizensService.ErrPropsWantedDenied"}, nil)
	ErrPropsJobDenied           = common.NewI18nErr(codes.PermissionDenied, &common.I18NItem{Key: "errors.CitizensService.ErrPropsJobDenied"}, nil)
	ErrPropsJobPublic           = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.CitizensService.ErrPropsJobPublic"}, nil)
	ErrPropsJobInvalid          = common.NewI18nErr(codes.InvalidArgument, &common.I18NItem{Key: "errors.CitizensService.ErrPropsJobInvalid"}, nil)
	ErrPropsTrafficPointsDenied = common.NewI18nErr(codes.PermissionDenied, &common.I18NItem{Key: "errors.CitizensService.ErrPropsTrafficPointsDenied"}, nil)
	ErrPropsMugshotDenied       = common.NewI18nErr(codes.PermissionDenied, &common.I18NItem{Key: "errors.CitizensService.ErrPropsMugshotDenied"}, nil)
	ErrPropsLabelsDenied        = common.NewI18nErr(codes.PermissionDenied, &common.I18NItem{Key: "errors.CitizensService.ErrPropsLabelsDenied"}, nil)
)
