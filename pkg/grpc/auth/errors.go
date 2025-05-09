package auth

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrNoToken          = common.I18nErr(codes.Unauthenticated, &common.TranslateItem{Key: "errors.pkg-auth.ErrNoToken"}, nil)
	ErrInvalidToken     = common.I18nErr(codes.Unauthenticated, &common.TranslateItem{Key: "errors.pkg-auth.ErrInvalidToken"}, nil)
	ErrCheckToken       = common.I18nErr(codes.Unauthenticated, &common.TranslateItem{Key: "errors.pkg-auth.ErrCheckToken"}, nil)
	ErrUserNoPerms      = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.pkg-auth.ErrUserNoPerms"}, nil)
	ErrNoUserInfo       = common.I18nErr(codes.Unauthenticated, &common.TranslateItem{Key: "errors.pkg-auth.ErrNoUserInfo"}, nil)
	ErrPermissionDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.pkg-auth.ErrPermissionDenied"}, nil)
	ErrCharLock         = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.AuthService.ErrCharLock.content"}, &common.TranslateItem{Key: "errors.AuthService.ErrCharLock.title"}) // Copied from the auth service
)
