package perms

import (
	permissionspermissions "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/permissions"
)

const (
	PermInternalNamespace Namespace = "internal"

	PermSuperuserService Service = "Superuser"
)

var (
	PermAnyRef = NewPermissionRef(PermInternalNamespace, "Any", "Any")
	PermAny    = &permissionspermissions.Permission{
		Namespace: string(PermInternalNamespace),
		Service:   "Any",
		Name:      "Any",
		GuardName: "internal-any-any",
	}

	PermCanBeSuperuser = &permissionspermissions.Permission{
		Namespace: string(PermInternalNamespace),
		Service:   string(PermSuperuserService),
		Name:      "CanBeSuperuser",
		GuardName: BuildGuard(PermInternalNamespace, PermSuperuserService, "CanBeSuperuser"),
	}
	PermCanBeSuperuserRef = NewPermissionRef(
		PermInternalNamespace,
		PermSuperuserService,
		"CanBeSuperuser",
	)

	PermConfigAdminRef = NewPermissionRef(
		PermInternalNamespace,
		PermSuperuserService,
		"ConfigAdmin",
	)
	PermConfigAdmin = &permissionspermissions.Permission{
		Namespace: string(PermInternalNamespace),
		Service:   string(PermSuperuserService),
		Name:      "ConfigAdmin",
		GuardName: BuildGuard(PermInternalNamespace, PermSuperuserService, "ConfigAdmin"),
	}

	PermSuperuserRef = NewPermissionRef(PermInternalNamespace, PermSuperuserService, "JobAdmin")
	PermSuperuser    = &permissionspermissions.Permission{
		Namespace: string(PermInternalNamespace),
		Service:   string(PermSuperuserService),
		Name:      "JobAdmin",
		GuardName: BuildGuard(PermInternalNamespace, PermSuperuserService, "JobAdmin"),
	}

	PermJobAdminRef = PermSuperuserRef
	PermJobAdmin    = PermSuperuser
)
