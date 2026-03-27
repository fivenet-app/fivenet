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
	ErrAccountCreateFailed = common.NewI18nErrFunc(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrAccountCreateFailed.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrAccountCreateFailed.title"},
	)
	ErrAccountExistsFailed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrAccountExistsFailed.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrAccountExistsFailed.title"},
	)
	ErrInvalidLogin = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrInvalidLogin.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrInvalidLogin.title"},
	)
	ErrNoAccount = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrNoAccount"},
		nil,
	)
	ErrNoCharFound = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrNoCharFound"},
		nil,
	)
	ErrGenericLogin = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrGenericLogin"},
		nil,
	)
	ErrUnableToChooseChar = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrUnableToChooseChar"},
		nil,
	)
	ErrUpdateAccount = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrUpdateAccount"},
		nil,
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
	ErrAccountDuplicate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrAccountDuplicate.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrAccountDuplicate.title"},
	)
	ErrChangeUsername = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrChangeUsername"},
		nil,
	)
	ErrBadUsername = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrBadUsername"},
		nil,
	)
	ErrNotSuperuser = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrNotSuperuser"},
		nil,
	)
	ErrImpersonateJobInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrImpersonateJobInvalid"},
		nil,
	)
)
