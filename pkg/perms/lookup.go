package perms

import (
	"github.com/galexrt/fivenet/pkg/perms/helpers"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (p *Perms) lookupAttributeByID(id uint64) (*cacheAttr, bool) {
	return p.attrsMap.Load(id)
}

func (p *Perms) lookupAttributeByPermID(id uint64, key Key) (*cacheAttr, bool) {
	as, ok := p.attrsPermsMap.Load(id)
	if !ok {
		return nil, false
	}

	aId, ok := as[key]
	if !ok {
		return nil, false
	}

	return p.lookupAttributeByID(aId)
}

func (p *Perms) lookupRoleAttribute(roleId uint64, attrId uint64) (*cacheRoleAttr, bool) {
	as, ok := p.attrsRoleMap.Load(roleId)
	if !ok {
		return nil, false
	}

	a, ok := as[attrId]
	return a, ok
}

// Roles
func (p *Perms) lookupRoleIDForJobAndGrade(job string, grade int32) (uint64, bool) {
	grades, ok := p.permsJobsRoleMap.Load(job)
	if !ok {
		return 0, false
	}
	roleId, ok := grades[grade]
	if !ok {
		return 0, false
	}

	return roleId, true
}

func (p *Perms) lookupRoleIDsForJobUpToGrade(job string, grade int32) ([]uint64, bool) {
	gradesMap, ok := p.permsJobsRoleMap.Load(job)
	if !ok {
		return nil, false
	}

	grades := []int32{}
	for g := range gradesMap {
		if g > grade {
			continue
		}
		grades = append(grades, g)
	}
	utils.SortInt32Slice(grades)
	gradeList := []uint64{}
	for i := 0; i < len(grades); i++ {
		gradeList = append(gradeList, gradesMap[grades[i]])
	}

	return gradeList, true
}

// Permissions
func (p *Perms) lookupPermIDByGuard(guard string) (uint64, bool) {
	return p.permsGuardToIDMap.Load(guard)
}

func (p *Perms) lookupPermissionByGuard(name string) (*model.FivenetPermissions, error) {
	guard := helpers.Guard(name)

	stmt := tPerms.
		SELECT(
			tPerms.AllColumns,
		).
		FROM(tPerms).
		WHERE(
			tPerms.GuardName.EQ(jet.String(guard)),
		).
		LIMIT(1)

	var dest model.FivenetPermissions
	err := stmt.QueryContext(p.ctx, p.db, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (p *Perms) lookupPermByID(id uint64) (*cachePerm, bool) {
	return p.permsMap.Load(id)
}
