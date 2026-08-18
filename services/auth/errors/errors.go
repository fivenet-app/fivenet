package errorsauth

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrGenericAccount = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrGenericAccount.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrGenericAccount.title"},
	)
	ErrInvalidLogin = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrInvalidLogin.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrInvalidLogin.title"},
	)
	ErrNoCharFound = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrNoCharFound.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrNoCharFound.title"},
	)
	ErrGenericLogin = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrGenericLogin.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrGenericLogin.title"},
	)
	ErrUnableToChooseChar = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrUnableToChooseChar.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrUnableToChooseChar.title"},
	)
	ErrChangePassword = common.NewI18nErrFunc(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrChangePassword.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrChangePassword.title"},
	)
	ErrForgotPassword = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrForgotPassword.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrForgotPassword.title"},
	)
	ErrSignupDisabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrSignupDisabled.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrSignupDisabled.title"},
	)
	ErrChangeUsername = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrChangeUsername.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrChangeUsername.title"},
	)
	ErrNotSuperuser = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrNotSuperuser.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrNotSuperuser.title"},
	)
	ErrImpersonateJobInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrImpersonateJobInvalid.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrImpersonateJobInvalid.title"},
	)
	ErrReauthRequired = common.NewI18nErr(
		codes.Unauthenticated,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrReauthRequired.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrReauthRequired.title"},
	)
)
