package perms

import (
	"fmt"
	"sync"

	"github.com/galexrt/arpanet/pkg/config"
)

var (
	list = []*Perm{
		{Key: "overview", Name: "View"},
	}

	mu sync.Mutex
)

type Perm struct {
	Key         string
	Name        string
	Description string
	Fields      []string
	PerJob      bool
}

func AddPermsToList(perms []*Perm) {
	mu.Lock()
	defer mu.Unlock()

	list = append(list, perms...)
}

func Register() {
	for _, perm := range list {
		baseKey := fmt.Sprintf("%s.%s", perm.Key, perm.Name)
		P.CreatePermission(baseKey, perm.Description)

		if perm.PerJob {
			for _, job := range config.C.FiveM.PermissionRoleJobs {
				jobKey := fmt.Sprintf("%s.%s", baseKey, job)
				P.CreatePermission(jobKey, perm.Description)
			}
			continue
		}

		for _, field := range perm.Fields {
			fieldKey := fmt.Sprintf("%s.%s", baseKey, field)
			P.CreatePermission(fieldKey, perm.Description)
		}
	}

	setupRoles()
}

func setupRoles() {
	P.CreateRole("masterofdisaster", "")
	perms, _ := P.GetAllPermissions()
	// Ensure the "masterofdisaster" role always has all permissions
	P.AddPermissionsToRole("masterofdisaster", perms)
}
