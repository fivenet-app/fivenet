package perms

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/perms/helpers"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	muPerms   sync.Mutex
	permsList = []*Perm{}
)

type Category string
type Name string

type Perm struct {
	Category Category
	Name     Name
	Attrs    []Attr
}

type Attr struct {
	ID            uint64
	Key           Key
	Type          permissions.AttributeTypes
	ValidValues   any
	DefaultValues any
}

func AddPermsToList(perms []*Perm) {
	muPerms.Lock()
	defer muPerms.Unlock()

	permsList = append(permsList, perms...)
}

func BuildGuard(category Category, name Name) string {
	return helpers.Guard(fmt.Sprintf("%s.%s", category, name))
}

func BuildGuardWithKey(category Category, name Name, key Key) string {
	return helpers.Guard(fmt.Sprintf("%s.%s.%s", category, name, key))
}

func (p *Perms) register(ctx context.Context, defaultRolePerms []string) error {
	if err := p.cleanupRoles(ctx); err != nil {
		return err
	}

	defaultRole, err := p.CreateRole(ctx, DefaultRoleJob, DefaultRoleJobGrade)
	if err != nil {
		return err
	}

	for _, perm := range permsList {
		permId, err := p.createOrUpdatePermission(ctx, perm.Category, perm.Name)
		if err != nil {
			return err
		}
		p.permsMap.Store(permId, &cachePerm{
			ID:        permId,
			Category:  perm.Category,
			Name:      perm.Name,
			GuardName: BuildGuard(perm.Category, perm.Name),
		})
		p.permsGuardToIDMap.Store(BuildGuard(perm.Category, perm.Name), permId)

		for _, attr := range perm.Attrs {
			switch attr.ValidValues {
			case "config.Game.Livemap.Jobs":
				attr.ValidValues = p.cfg.Game.Livemap.Jobs
			}

			if _, err := p.createOrUpdateAttribute(ctx, permId, attr.Key, attr.Type, attr.ValidValues, attr.DefaultValues); err != nil {
				return err
			}
		}
	}

	if err := p.setupDefaultRolePerms(ctx, defaultRole, defaultRolePerms); err != nil {
		return err
	}

	return nil
}

func (p *Perms) setupDefaultRolePerms(ctx context.Context, role *model.FivenetRoles, defaultPerms []string) error {
	if len(defaultPerms) == 0 {
		return nil
	}

	addPerms := make([]AddPerm, len(defaultPerms))
	for i, perm := range defaultPerms {
		permId, ok := p.permsGuardToIDMap.Load(perm)
		if !ok {
			return fmt.Errorf("permission by guard %s not found", perm)
		}

		addPerms[i] = AddPerm{
			Id:  permId,
			Val: true,
		}
	}

	if err := p.UpdateRolePermissions(ctx, role.ID, addPerms...); err != nil {
		return err
	}

	return nil
}

func (p *Perms) createOrUpdatePermission(ctx context.Context, category Category, name Name) (uint64, error) {
	perm, err := p.loadPermissionFromDatabaseByGuard(ctx, BuildGuard(category, name))
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	if perm != nil {
		if Category(perm.Category) != category || Name(perm.Name) != name {
			return perm.ID, p.UpdatePermission(ctx, perm.ID, category, name)
		}

		return perm.ID, nil
	}

	return p.CreatePermission(ctx, category, name)
}

func (p *Perms) createOrUpdateAttribute(ctx context.Context, permId uint64, key Key, aType permissions.AttributeTypes, validValues any, defaultValues any) (uint64, error) {
	attr, err := p.getAttributeFromDatabase(ctx, permId, key)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	if attr != nil && attr.ID > 0 {
		var validVal interface{}
		if validValues != nil {
			validVal, err = json.MarshalToString(validValues)
			if err != nil {
				return 0, err
			}
		}

		if err := p.addOrUpdateAttributeInMap(permId, attr.ID, key, aType, validValues, defaultValues); err != nil {
			return 0, err
		}

		if attr.Type != string(aType) || (attr.ValidValues == nil || validVal != *attr.ValidValues) || (attr.DefaultValues == nil || defaultValues != *attr.DefaultValues) {
			return attr.ID, p.UpdateAttribute(ctx, attr.ID, permId, key, aType, validValues, defaultValues)
		}

		return attr.ID, nil
	}

	attrId, err := p.CreateAttribute(ctx, permId, key, aType, validValues, defaultValues)
	if err != nil {
		return 0, err
	}

	if err := p.addOrUpdateAttributeInMap(permId, attrId, key, aType, validValues, defaultValues); err != nil {
		return 0, err
	}

	return attrId, nil
}

func (p *Perms) cleanupRoles(ctx context.Context) error {
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

	var dest []*users.Job
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}
	jobName := DefaultRoleJob
	jobGrade := DefaultRoleJobGrade
	// Add default job to avoid it being deleted
	dest = append(dest, &users.Job{
		Name: DefaultRoleJob,
		Grades: []*users.JobGrade{
			{
				JobName: &jobName,
				Grade:   jobGrade,
			},
		},
	})

	allRoles, err := p.GetRoles(ctx, false)
	if err != nil {
		return err
	}
	existingRoles := allRoles.IDs()

	// Iterate over current job and job grades to find any non-existant roles in our database
	for i := 0; i < len(dest); i++ {
		for _, grade := range dest[i].Grades {
			role, err := p.GetRoleByJobAndGrade(ctx, dest[i].Name, grade.Grade)
			if err != nil {
				return err
			}
			if role == nil {
				continue
			}

			index := slices.Index(existingRoles, role.ID)
			if index >= 0 {
				existingRoles = append(existingRoles[:index], existingRoles[index+1:]...)
			}
		}
	}

	// Remove all roles that shouldn't exist anymore
	for i := 0; i < len(existingRoles); i++ {
		if err := p.DeleteRole(ctx, existingRoles[i]); err != nil {
			return err
		}
	}

	return nil
}

func (p *Perms) getActiveJobs(ctx context.Context) ([]string, error) {
	stmt := tRoles.
		SELECT(
			tRoles.Job,
		).
		FROM(tRoles).
		WHERE(
			tRoles.Job.NOT_EQ(jet.String(DefaultRoleJob)),
		).
		GROUP_BY(tRoles.Job)

	var dest []string
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (p *Perms) ApplyJobPermissions(ctx context.Context, job string) error {
	jobs := []string{}
	if job != "" {
		jobs = append(jobs, job)
	} else {
		var err error
		jobs, err = p.getActiveJobs(ctx)
		if err != nil {
			return err
		}
	}

	for _, job := range jobs {
		if err := p.applyJobPermissions(ctx, job); err != nil {
			return err
		}
	}

	return nil
}

func (p *Perms) applyJobPermissions(ctx context.Context, job string) error {
	roles, err := p.GetJobRoles(ctx, job)
	if err != nil {
		return err
	}

	if len(roles) == 0 {
		return nil
	}

	jps, err := p.GetJobPermissions(ctx, job)
	if err != nil {
		return err
	}

	for _, role := range roles {
		ps, err := p.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return err
		}

		if len(ps) == 0 {
			continue
		}

		if len(jps) == 0 {
			pIds := []uint64{}
			for _, p := range ps {
				pIds = append(pIds, p.Id)
			}

			if err := p.RemovePermissionsFromRole(ctx, role.ID, pIds...); err != nil {
				return err
			}
			continue
		}

		toRemove := []uint64{}
		for _, p := range ps {
			if !slices.ContainsFunc(jps, func(in *permissions.Permission) bool {
				return in.Id == p.Id && in.Val
			}) {
				toRemove = append(toRemove, p.Id)
			}
		}

		if len(toRemove) > 0 {
			if p.RemovePermissionsFromRole(ctx, role.ID, toRemove...); err != nil {
				return err
			}
		}
	}

	return p.applyJobPermissionsToAttrs(ctx, job)
}

func (p *Perms) applyJobPermissionsToAttrs(ctx context.Context, job string) error {
	roles, err := p.GetJobRoles(ctx, job)
	if err != nil {
		return err
	}

	if len(roles) == 0 {
		return nil
	}

	jps, err := p.GetJobPermissions(ctx, job)
	if err != nil {
		return err
	}

	for _, role := range roles {
		attrs, err := p.GetRoleAttributes(role.Job, role.Grade)
		if err != nil {
			return err
		}

		if len(attrs) == 0 {
			continue
		}

		if len(jps) == 0 {
			if err := p.RemoveAttributesFromRole(ctx, role.ID, attrs...); err != nil {
				return err
			}
			continue
		}

		toRemove := []*permissions.RoleAttribute{}
		toUpdate := []*permissions.RoleAttribute{}
		for _, attr := range attrs {
			maxValues, _ := p.GetClosestRoleAttrMaxVals(role.Job, role.Grade, attr.PermissionId, Key(attr.Key))
			if slices.ContainsFunc(jps, func(in *permissions.Permission) bool {
				return in.Id == attr.PermissionId
			}) {
				if _, changed := attr.Value.Check(permissions.AttributeTypes(attr.Type), attr.ValidValues, maxValues); changed {
					toUpdate = append(toUpdate, attr)
				}
			} else {
				toRemove = append(toRemove, attr)
			}
		}

		if len(toRemove) > 0 {
			if p.RemoveAttributesFromRole(ctx, role.ID, toRemove...); err != nil {
				return err
			}
		}

		if len(toUpdate) > 0 {
			if p.AddOrUpdateAttributesToRole(ctx, role.Job, role.Grade, role.ID, toUpdate...); err != nil {
				return err
			}
		}
	}

	return nil
}