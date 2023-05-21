package perms

import (
	"github.com/galexrt/fivenet/pkg/utils"
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
	roles, ok := p.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok || len(roles) == 0 {
		return 0, false
	}

	return roles[len(roles)-1], true
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

	if len(gradeList) == 0 {
		return nil, false
	}

	return gradeList, true
}

// Permissions
func (p *Perms) lookupPermIDByGuard(guard string) (uint64, bool) {
	return p.permsGuardToIDMap.Load(guard)
}

func (p *Perms) lookupPermByID(id uint64) (*cachePerm, bool) {
	return p.permsMap.Load(id)
}
