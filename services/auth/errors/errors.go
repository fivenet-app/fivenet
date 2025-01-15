package errorsauth

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrAccountCreateFailed = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrAccountCreateFailed"}, nil)
	ErrAccountExistsFailed = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrAccountExistsFailed"}, nil)
	ErrInvalidLogin        = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrInvalidLogin"}, nil)
	ErrNoAccount           = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrNoAccount"}, nil)
	ErrNoCharFound         = common.I18nErr(codes.NotFound, &common.TranslateItem{Key: "errors.AuthService.ErrNoCharFound"}, nil)
	ErrGenericLogin        = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.AuthService.ErrGenericLogin"}, nil)
	ErrUnableToChooseChar  = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.AuthService.ErrUnableToChooseChar"}, nil)
	ErrUpdateAccount       = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrUpdateAccount"}, nil)
	ErrChangePassword      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrChangePassword"}, nil)
	ErrForgotPassword      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrForgotPassword"}, nil)
	ErrSignupDisabled      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrSignupDisabled"}, nil)
	ErrAccountDuplicate    = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrAccountDuplicate"}, nil)
	ErrChangeUsername      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrChangeUsername"}, nil)
	ErrBadUsername         = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrBadUsername"}, nil)
	ErrCharLock            = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.AuthService.ErrCharLock.content"}, &common.TranslateItem{Key: "errors.AuthService.ErrCharLock.title"})
)
