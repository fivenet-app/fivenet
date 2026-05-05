package perms

import "slices"

// Attrs

func (ps *Perms) lookupAttributeByID(id int64) (*cacheAttr, bool) {
	return ps.attrsMap.Load(id)
}

func (ps *Perms) lookupAttributeByPermID(id int64, key Key) (*cacheAttr, bool) {
	as, ok := ps.attrsPermsMap.Load(id)
	if !ok {
		return nil, false
	}

	aId, ok := as.Load(string(key))
	if !ok {
		return nil, false
	}

	return ps.lookupAttributeByID(aId)
}

func (ps *Perms) lookupRoleAttribute(roleId int64, attrId int64) (*cacheRoleAttr, bool) {
	attrMap, ok := ps.attrsRoleMap.Load(roleId)
	if !ok {
		return nil, false
	}

	return attrMap.Load(attrId)
}

// Roles

func (ps *Perms) lookupJobForRoleID(roleId int64) (string, bool) {
	return ps.roleIDToJobMap.Load(roleId)
}

func (ps *Perms) lookupRoleIDForJobAndGrade(job string, grade int32) (int64, bool) {
	roles, ok := ps.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok || len(roles) == 0 {
		return 0, false
	}

	return roles[len(roles)-1], true
}

// Lookup all role IDs for a job up to a certain grade.
func (ps *Perms) lookupRoleIDsForJobUpToGrade(job string, grade int32) ([]int64, bool) {
	gradesMap, ok := ps.permsJobsRoleMap.Load(job)
	if !ok {
		return nil, false
	}

	grades := []int32{}
	for g := range gradesMap.All() {
		if g > grade {
			continue
		}

		grades = append(grades, g)
	}

	slices.Sort(grades)

	gradeList := []int64{}
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

func (ps *Perms) lookupPermIDByGuard(guard string) (int64, bool) {
	return ps.permsGuardToIDMap.Load(guard)
}

func (ps *Perms) lookupPermByID(id int64) (*cachePerm, bool) {
	return ps.permsMap.Load(id)
}
