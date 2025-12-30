package errorsauth

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrGenericAccount = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.AuthService.ErrGenericAccount"},
		nil,
	)
	ErrAccountCreateFailed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrAccountCreateFailed"},
		nil,
	)
	ErrAccountExistsFailed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrAccountExistsFailed"},
		nil,
	)
	ErrInvalidLogin = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrInvalidLogin"},
		nil,
	)
	ErrNoAccount = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrNoAccount"},
		nil,
	)
	ErrNoCharFound = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.AuthService.ErrNoCharFound"},
		nil,
	)
	ErrGenericLogin = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.AuthService.ErrGenericLogin"},
		nil,
	)
	ErrUnableToChooseChar = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.AuthService.ErrUnableToChooseChar"},
		nil,
	)
	ErrUpdateAccount = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrUpdateAccount"},
		nil,
	)
	ErrChangePassword = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrChangePassword"},
		nil,
	)
	ErrForgotPassword = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrForgotPassword"},
		nil,
	)
	ErrSignupDisabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrSignupDisabled"},
		nil,
	)
	ErrAccountDuplicate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrAccountDuplicate"},
		nil,
	)
	ErrChangeUsername = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrChangeUsername"},
		nil,
	)
	ErrBadUsername = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrBadUsername"},
		nil,
	)
	ErrNotSuperuser = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.AuthService.ErrNotSuperuser"},
		nil,
	)
)
