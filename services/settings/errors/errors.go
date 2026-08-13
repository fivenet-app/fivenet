package errorssettings

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrFailedQuery.content"},
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrFailedQuery.title"},
	)
	ErrInvalidRequest = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrInvalidRequest.content"},
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrInvalidRequest.title"},
	)
	ErrNoPermission = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrNoPermission.content"},
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrNoPermission.title"},
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
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrInvalidAttrs.content"},
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrInvalidAttrs.title"},
	)
	ErrInvalidPerms = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrInvalidPerms.content"},
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrInvalidPerms.title"},
	)

	ErrDiscordNotEnabled = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrDiscordNotEnabled.content"},
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrDiscordNotEnabled.title"},
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
	ErrDiscordTokenExpired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrDiscordTokenExpired.content"},
		&common.I18NItem{Key: "errors.settings.SettingsService.ErrDiscordTokenExpired.title"},
	)
)
