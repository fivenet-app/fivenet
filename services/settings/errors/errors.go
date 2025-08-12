package errorssettings

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.SettingsService.ErrFailedQuery"},
		nil,
	)
	ErrInvalidRequest = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.SettingsService.ErrInvalidRequest"},
		nil,
	)
	ErrNoPermission = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.SettingsService.ErrNoPermission"},
		nil,
	)
	ErrRoleAlreadyExists = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.SettingsService.ErrRoleAlreadyExists"},
		nil,
	)
	ErrOwnRoleDeletion = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.SettingsService.ErrOwnRoleDeletion"},
		nil,
	)
	ErrInvalidAttrs = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.SettingsService.ErrInvalidAttrs"},
		nil,
	)
	ErrInvalidPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.SettingsService.ErrInvalidPerms"},
		nil,
	)

	ErrDiscordNotEnabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.SettingsService.ErrDiscordNotEnabled"},
		nil,
	)
	ErrDiscordConnectRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.SettingsService.ErrDiscordConnectRequired.content"},
		&common.I18NItem{Key: "errors.SettingsService.ErrDiscordConnectRequired.title"},
	)
)
