// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/rector/config.proto
// source: services/rector/filestore.proto
// source: services/rector/laws.proto
// source: services/rector/rector.proto
// source: services/rector/sync.proto

package permsrector

import (
	"github.com/fivenet-app/fivenet/pkg/perms"
)

const (
	RectorConfigServicePerm    perms.Category = "RectorConfigService"
	RectorFilestoreServicePerm perms.Category = "RectorFilestoreService"
	RectorLawsServicePerm      perms.Category = "RectorLawsService"
	RectorServicePerm          perms.Category = "RectorService"
	SyncServicePerm            perms.Category = "SyncService"

	RectorServiceCreateRolePerm      perms.Name = "CreateRole"
	RectorServiceDeleteRolePerm      perms.Name = "DeleteRole"
	RectorServiceGetJobPropsPerm     perms.Name = "GetJobProps"
	RectorServiceGetRolesPerm        perms.Name = "GetRoles"
	RectorServiceSetJobPropsPerm     perms.Name = "SetJobProps"
	RectorServiceUpdateRolePermsPerm perms.Name = "UpdateRolePerms"
	RectorServiceViewAuditLogPerm    perms.Name = "ViewAuditLog"
)
