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

	PermSuperuserRef = NewPermissionRef(PermInternalNamespace, PermSuperuserService, "Superuser")
	PermSuperuser    = &permissionspermissions.Permission{
		Namespace: string(PermInternalNamespace),
		Service:   string(PermSuperuserService),
		Name:      "Superuser",
		GuardName: BuildGuard(PermInternalNamespace, PermSuperuserService, "Superuser"),
	}
)
