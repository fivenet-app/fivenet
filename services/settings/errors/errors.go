package errorssettings

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery       = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.SettingsService.ErrFailedQuery"}, nil)
	ErrInvalidRequest    = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.SettingsService.ErrInvalidRequest"}, nil)
	ErrNoPermission      = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.SettingsService.ErrNoPermission"}, nil)
	ErrRoleAlreadyExists = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.SettingsService.ErrRoleAlreadyExists"}, nil)
	ErrOwnRoleDeletion   = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.SettingsService.ErrOwnRoleDeletion"}, nil)
	ErrInvalidAttrs      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.SettingsService.ErrInvalidAttrs"}, nil)
	ErrInvalidPerms      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.SettingsService.ErrInvalidPerms"}, nil)

	ErrDiscordNotEnabled      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.SettingsService.ErrDiscordNotEnabled"}, nil)
	ErrDiscordConnectRequired = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.SettingsService.ErrDiscordConnectRequired.content"}, &common.TranslateItem{Key: "errors.SettingsService.ErrDiscordConnectRequired.title"})
)
