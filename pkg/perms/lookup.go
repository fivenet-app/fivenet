package perms

import "slices"

// Attrs

func (p *Perms) LookupAttributeByID(id uint64) (*cacheAttr, bool) {
	return p.attrsMap.Load(id)
}

func (p *Perms) lookupAttributeByPermID(id uint64, key Key) (*cacheAttr, bool) {
	as, ok := p.attrsPermsMap.Load(id)
	if !ok {
		return nil, false
	}

	aId, ok := as.Load(string(key))
	if !ok {
		return nil, false
	}

	return p.LookupAttributeByID(aId)
}

func (p *Perms) lookupRoleAttribute(roleId uint64, attrId uint64) (*cacheRoleAttr, bool) {
	as, ok := p.attrsRoleMap.Load(roleId)
	if !ok {
		return nil, false
	}

	return as.Load(attrId)
}

// Roles

func (p *Perms) lookupJobForRoleID(roleId uint64) (string, bool) {
	return p.roleIDToJobMap.Load(roleId)
}

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
	gradesMap.Range(func(g int32, _ uint64) bool {
		if g > grade {
			return true
		}

		grades = append(grades, g)
		return true
	})

	slices.Sort(grades)

	gradeList := []uint64{}
	for i := range grades {
		grade, ok := gradesMap.Load(grades[i])
		if !ok {
			return nil, false
		}

		gradeList = append(gradeList, grade)
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
