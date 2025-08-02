package errorsgrpcauth

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrNoToken          = common.NewI18nErr(codes.Unauthenticated, &common.I18NItem{Key: "errors.pkg-auth.ErrNoToken"}, nil)
	ErrInvalidToken     = common.NewI18nErr(codes.Unauthenticated, &common.I18NItem{Key: "errors.pkg-auth.ErrInvalidToken"}, nil)
	ErrCheckToken       = common.NewI18nErr(codes.Unauthenticated, &common.I18NItem{Key: "errors.pkg-auth.ErrCheckToken"}, nil)
	ErrUserNoPerms      = common.NewI18nErr(codes.PermissionDenied, &common.I18NItem{Key: "errors.pkg-auth.ErrUserNoPerms"}, nil)
	ErrNoUserInfo       = common.NewI18nErr(codes.Unauthenticated, &common.I18NItem{Key: "errors.pkg-auth.ErrNoUserInfo"}, nil)
	ErrPermissionDenied = common.NewI18nErr(codes.PermissionDenied, &common.I18NItem{Key: "errors.pkg-auth.ErrPermissionDenied"}, nil)
	ErrCharLock         = common.NewI18nErr(codes.PermissionDenied, &common.I18NItem{Key: "errors.AuthService.ErrCharLock.content"}, &common.I18NItem{Key: "errors.AuthService.ErrCharLock.title"})
)
