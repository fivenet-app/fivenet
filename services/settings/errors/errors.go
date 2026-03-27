package errorssettings

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrFailedQuery"},
		nil,
	)
	ErrInvalidRequest = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrInvalidRequest"},
		nil,
	)
	ErrNoPermission = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrNoPermission"},
		nil,
	)
	ErrRoleAlreadyExists = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrRoleAlreadyExists"},
		nil,
	)
	ErrOwnRoleDeletion = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrOwnRoleDeletion"},
		nil,
	)
	ErrInvalidAttrs = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrInvalidAttrs"},
		nil,
	)
	ErrInvalidPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrInvalidPerms"},
		nil,
	)

	ErrDiscordNotEnabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrDiscordNotEnabled"},
		nil,
	)
	ErrDiscordConnectRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrDiscordConnectRequired.content"},
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrDiscordConnectRequired.title"},
	)
	ErrCannotDeleteOwnAccount = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrCannotDeleteOwnAccount.content"},
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrCannotDeleteOwnAccount.title"},
	)
)
