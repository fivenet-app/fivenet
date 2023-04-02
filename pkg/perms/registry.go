package perms

import (
	"fmt"
	"strings"
	"sync"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/proto/resources/jobs"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"golang.org/x/exp/slices"
)

var (
	list = []*Perm{}

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

func (p *Perms) Register() error {
	for _, perm := range list {
		baseKey := fmt.Sprintf("%s.%s", perm.Key, perm.Name)
		if err := p.CreatePermission(baseKey, perm.Description); err != nil {
			return err
		}

		if perm.PerJob {
			for _, job := range config.C.FiveM.PermissionRoleJobs {
				jobKey := fmt.Sprintf("%s.%s", baseKey, job)
				if err := p.CreatePermission(jobKey, perm.Description); err != nil {
					return err
				}
			}
			continue
		}

		for _, field := range perm.Fields {
			fieldKey := fmt.Sprintf("%s.%s", baseKey, field)
			if err := p.CreatePermission(fieldKey, perm.Description); err != nil {
				return err
			}
		}
	}

	return p.setupRoles()
}

func (p *Perms) setupRoles() error {
	role, err := p.CreateRole("masterofdisaster", "")
	if err != nil {
		return err
	}

	perms, _ := p.GetAllPermissions()
	// Ensure the "masterofdisaster" role always has all permissions
	if err := p.AddPermissionsToRole(role.ID, perms.IDs()); err != nil {
		return err
	}

	j := table.Jobs.AS("job")
	jg := table.JobGrades.AS("jobgrade")
	stmt := j.
		SELECT(
			j.Name,
			j.Label,
			jg.JobName.AS("jobname"),
			jg.Grade,
			jg.Name,
			jg.Label,
		).
		FROM(j.
			INNER_JOIN(jg,
				jg.JobName.EQ(j.Name),
			))

	var dest []*jobs.Job
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		return err
	}

	existingRolesList, err := p.GetRoles("job-")
	if err != nil {
		return err
	}
	existingRoles := existingRolesList.IDs()

	// Iterate over current job and job grades to find any non-existant roles in our database
	for i := 0; i < len(dest); i++ {
		for _, grade := range dest[i].Grades {
			roleName := strings.ToLower(GetRoleName(dest[i].Name, grade.Grade))

			role, err := p.GetRoleByGuardName(roleName)
			if err != nil {
				return err
			}

			index := slices.Index(existingRoles, role.ID)
			if index >= 0 {
				existingRoles = append(existingRoles[:index], existingRoles[index+1:]...)
			}
		}
	}

	// Remove all roles that shouldn't exist anymore
	for i := 0; i < len(existingRoles); i++ {
		if err := p.DeleteRole(existingRoles[i]); err != nil {
			return err
		}
	}

	return nil
}

func GetRoleName(job string, grade int32) string {
	return fmt.Sprintf("job-%s-%d", job, grade)
}
