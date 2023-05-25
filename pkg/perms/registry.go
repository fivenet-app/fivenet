package perms

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/perms/helpers"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/exp/slices"
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
	ID          uint64
	Key         Key
	Type        AttributeTypes
	ValidValues any
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

func (p *Perms) Register(defaultRolePerms []string) error {
	ctx, span := p.tracer.Start(p.ctx, "perms-register")
	defer span.End()

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
			case "config.C.Game.Livemap.Jobs":
				attr.ValidValues = config.C.Game.Livemap.Jobs
			}

			_, err := p.createOrUpdateAttribute(ctx, permId, attr.Key, attr.Type, attr.ValidValues)
			if err != nil {
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
		if !errors.Is(qrm.ErrNoRows, err) {
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

func (p *Perms) createOrUpdateAttribute(ctx context.Context, permId uint64, key Key, aType AttributeTypes, validValues any) (uint64, error) {
	attr, err := p.getAttributeFromDatabase(ctx, permId, key)
	if err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return 0, err
		}
	}

	if attr != nil {
		var validVal interface{}
		if validValues != nil {
			validVal, err = json.MarshalToString(validValues)
			if err != nil {
				return 0, err
			}
		}

		if attr.Type != string(aType) ||
			((attr.ValidValues == nil && validVal != nil) || (validVal != nil && attr.ValidValues != nil && validVal != *attr.ValidValues)) {
			return attr.ID, p.UpdateAttribute(ctx, attr.ID, permId, key, aType, validValues)
		}

		if err := p.addOrUpdateAttributeInMap(permId, attr.ID, key, aType, validValues); err != nil {
			return 0, err
		}

		return attr.ID, nil
	}

	return p.CreateAttribute(ctx, permId, key, aType, validValues)
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

	var dest []*jobs.Job
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return err
	}
	// Add default job to avoid deletion
	dest = append(dest, &jobs.Job{
		Name: DefaultRoleJob,
		Grades: []*jobs.JobGrade{
			{
				JobName: DefaultRoleJob,
				Grade:   DefaultRoleJobGrade,
			},
		},
	})

	allRoles, err := p.getRoles(ctx)
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
