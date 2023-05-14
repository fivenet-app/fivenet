package perms

import (
	"errors"
	"fmt"
	"sync"

	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/pkg/perms/helpers"
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
	ValidValues string
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

func (p *Perms) Register() error {
	for _, perm := range permsList {
		permId, err := p.createOrUpdatePermission(perm.Category, perm.Name)
		if err != nil {
			return err
		}
		p.guardToPermIDMap.Store(BuildGuard(perm.Category, perm.Name), permId)

		for _, attr := range perm.Attrs {
			attrId, err := p.createOrUpdateAttribute(permId, attr.Key, attr.Type, attr.ValidValues)
			if err != nil {
				return err
			}
			attr.ID = attrId
			p.permIdToAttrsMap.LoadOrStore(permId, map[Key]Attr{attr.Key: attr})
		}
	}

	return p.cleanupRoles()
}

func (p *Perms) createOrUpdatePermission(category Category, name Name) (uint64, error) {
	perm, err := p.getPermissionByGuard(BuildGuard(category, name))
	if err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return 0, err
		}
	}

	if perm != nil {
		if Category(perm.Category) != category || Name(perm.Name) != name {
			return perm.ID, p.UpdatePermission(perm.ID, category, name)
		}

		return perm.ID, nil
	}

	return p.CreatePermission(category, name)
}

func (p *Perms) createOrUpdateAttribute(permId uint64, key Key, aType AttributeTypes, validValues string) (uint64, error) {
	attr, err := p.GetAttribute(permId, key)
	if err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return 0, err
		}
	}

	var validVals interface{}
	if validValues != "" {
		if err = json.UnmarshalFromString(validValues, &validVals); err != nil {
			return 0, err
		}
	}

	if attr != nil {
		if Key(attr.Key) != key ||
			((attr.ValidValues == nil && validValues != "") || (attr.ValidValues != nil && attr.ValidValues != &validValues)) {
			return attr.ID, p.UpdateAttribute(attr.ID, permId, key, aType, validVals)
		}

		return attr.ID, nil
	}

	return p.CreateAttribute(permId, key, aType, validVals)
}

func (p *Perms) cleanupRoles() error {
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

	allRoles, err := p.getRoles()
	if err != nil {
		return err
	}
	existingRoles := allRoles.IDs()

	// Iterate over current job and job grades to find any non-existant roles in our database
	for i := 0; i < len(dest); i++ {
		for _, grade := range dest[i].Grades {
			role, err := p.GetRoleByJobAndGrade(dest[i].Name, grade.Grade)
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
		if err := p.DeleteRole(existingRoles[i]); err != nil {
			return err
		}
	}

	return nil
}
