package errorsgrpcauth

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrNoToken = common.NewI18nErr(
		codes.Unauthenticated,
		&common.I18NItem{Key: "errors.pkg-auth.ErrNoToken.content"},
		&common.I18NItem{Key: "errors.pkg-auth.ErrNoToken.title"},
	)
	ErrInvalidToken = common.NewI18nErr(
		codes.Unauthenticated,
		&common.I18NItem{Key: "errors.pkg-auth.ErrInvalidToken.content"},
		&common.I18NItem{Key: "errors.pkg-auth.ErrInvalidToken.title"},
	)
	ErrCheckToken = common.NewI18nErr(
		codes.Unauthenticated,
		&common.I18NItem{Key: "errors.pkg-auth.ErrCheckToken.content"},
		&common.I18NItem{Key: "errors.pkg-auth.ErrCheckToken.title"},
	)
	ErrUserNoPerms = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.pkg-auth.ErrUserNoPerms.content"},
		&common.I18NItem{Key: "errors.pkg-auth.ErrUserNoPerms.title"},
	)
	ErrNoUserInfo = common.NewI18nErr(
		codes.Unauthenticated,
		&common.I18NItem{Key: "errors.pkg-auth.ErrNoUserInfo.content"},
		&common.I18NItem{Key: "errors.pkg-auth.ErrNoUserInfo.title"},
	)
	ErrPermissionDenied = common.NewI18nErrFunc(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.pkg-auth.ErrPermissionDenied.content"},
		&common.I18NItem{Key: "errors.pkg-auth.ErrPermissionDenied.title"},
	)
	ErrCharLock = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.auth.AuthService.ErrCharLock.content"},
		&common.I18NItem{Key: "errors.auth.AuthService.ErrCharLock.title"},
	)
)
