// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/settings/accounts.proto
// source: services/settings/config.proto
// source: services/settings/cron.proto
// source: services/settings/laws.proto
// source: services/settings/settings.proto

package permssettings

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
)

const (
	AccountsServicePerm perms.Category = "settings.AccountsService"
	ConfigServicePerm   perms.Category = "settings.ConfigService"
	CronServicePerm     perms.Category = "settings.CronService"
	LawsServicePerm     perms.Category = "settings.LawsService"
	SettingsServicePerm perms.Category = "settings.SettingsService"

	LawsServiceCreateOrUpdateLawBookPerm perms.Name = "CreateOrUpdateLawBook"
	LawsServiceDeleteLawBookPerm         perms.Name = "DeleteLawBook"
	SettingsServiceCreateRolePerm        perms.Name = "CreateRole"
	SettingsServiceDeleteRolePerm        perms.Name = "DeleteRole"
	SettingsServiceGetJobPropsPerm       perms.Name = "GetJobProps"
	SettingsServiceGetRolesPerm          perms.Name = "GetRoles"
	SettingsServiceSetJobPropsPerm       perms.Name = "SetJobProps"
	SettingsServiceUpdateRolePermsPerm   perms.Name = "UpdateRolePerms"
	SettingsServiceViewAuditLogPerm      perms.Name = "ViewAuditLog"
)
