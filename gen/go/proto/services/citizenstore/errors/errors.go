package errorscitizenstore

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery              = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrFailedQuery"}, nil)
	ErrJobGradeNoPermission     = common.I18nErr(codes.NotFound, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrJobGradeNoPermission"}, nil)
	ErrReasonRequired           = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrReasonRequired"}, nil)
	ErrPropsWantedDenied        = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrPropsWantedDenied"}, nil)
	ErrPropsJobDenied           = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrPropsJobDenied"}, nil)
	ErrPropsJobPublic           = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrPropsJobPublic"}, nil)
	ErrPropsJobInvalid          = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrPropsJobInvalid"}, nil)
	ErrPropsTrafficPointsDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrPropsTrafficPointsDenied"}, nil)
	ErrPropsMugShotDenied       = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrPropsMugShotDenied"}, nil)
	ErrPropsLabelsDenied        = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.CitizenStoreService.ErrPropsLabelsDenied"}, nil)
)
