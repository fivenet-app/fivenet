package permissions

import (
	"fmt"
	"sync"

	"github.com/galexrt/arpanet/pkg/config"
	"github.com/galexrt/arpanet/pkg/permify/options"
	"github.com/galexrt/arpanet/query"
)

var (
	Perms = []*Perm{}

	mu sync.Mutex
)

type Perm struct {
	Key         string
	Name        string
	Description string
	Fields      []string
	PerJob      bool
}

func RegisterPerms(perms []*Perm) {
	mu.Lock()
	defer mu.Unlock()

	Perms = append(Perms, perms...)
}

func createPermission(key string, description string) error {
	return query.Perms.CreatePermission(key, description)
}

func Register() {
	for _, perm := range Perms {
		baseKey := fmt.Sprintf("%s.%s", perm.Key, perm.Name)
		createPermission(baseKey, perm.Description)

		if perm.PerJob {
			for _, job := range config.C.FiveM.PermissionRoleJobs {
				jobKey := fmt.Sprintf("%s.%s", baseKey, job)
				createPermission(jobKey, perm.Description)
			}
			continue
		}

		for _, field := range perm.Fields {
			fieldKey := fmt.Sprintf("%s.%s", baseKey, field)
			createPermission(fieldKey, perm.Description)
		}
	}

	setupRoles()
}

func setupRoles() {
	query.Perms.CreateRole("masterofdisaster", "")
	perms, _, _ := query.Perms.GetAllPermissions(options.PermissionOption{})
	// Ensure the "masterofdisaster" role always has all permissions
	query.Perms.AddPermissionsToRole("masterofdisaster", perms.IDs())
}
