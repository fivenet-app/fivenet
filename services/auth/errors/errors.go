package errorsauth

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrGenericAccount = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.AuthService.ErrGenericAccount.content"},
		&common.I18NItem{Key: "errors.AuthService.ErrGenericAccount.title"},
	)
	ErrAccountCreateFailed = common.NewI18nErrFunc(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrAccountCreateFailed.content"},
		&common.I18NItem{Key: "errors.AuthService.ErrAccountCreateFailed.title"},
	)
	ErrAccountExistsFailed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrAccountExistsFailed.content"},
		&common.I18NItem{Key: "errors.AuthService.ErrAccountExistsFailed.title"},
	)
	ErrInvalidLogin = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrInvalidLogin.content"},
		&common.I18NItem{Key: "errors.AuthService.ErrInvalidLogin.title"},
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
	ErrChangePassword = common.NewI18nErrFunc(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrChangePassword.content"},
		&common.I18NItem{Key: "errors.AuthService.ErrChangePassword.title"},
	)
	ErrForgotPassword = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrForgotPassword.content"},
		&common.I18NItem{Key: "errors.AuthService.ErrForgotPassword.title"},
	)
	ErrSignupDisabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrSignupDisabled.content"},
		&common.I18NItem{Key: "errors.AuthService.ErrSignupDisabled.title"},
	)
	ErrAccountDuplicate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrAccountDuplicate.content"},
		&common.I18NItem{Key: "errors.AuthService.ErrAccountDuplicate.title"},
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
	ErrImpersonateJobInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.AuthService.ErrImpersonateJobInvalid"},
		nil,
	)
)
