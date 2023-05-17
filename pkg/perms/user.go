package perms

import (
	"fmt"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
)

func (p *Perms) GetPermissionsOfUser(userInfo *userinfo.UserInfo) (collections.Permissions, error) {
	roleId, ok := p.lookupRoleIDForJobAndGrade(userInfo.Job, userInfo.JobGrade)
	if !ok {
		roleId, ok = p.lookupRoleIDForJobAndGrade(DefaultRoleJob, DefaultRoleJobGrade)
		if !ok {
			return nil, fmt.Errorf("failed to fallback to default role")
		}
	}

	ps, err := p.GetRolePermissions(roleId)
	if err != nil {
		return nil, err
	}

	perms := make(collections.Permissions, len(ps))
	for i := 0; i < len(ps); i++ {
		var createdAt *time.Time
		if ps[i].CreatedAt != nil {
			createdAtTime := ps[i].CreatedAt.AsTime()
			createdAt = &createdAtTime
		}

		perms[i] = &model.FivenetPermissions{
			ID:        ps[i].Id,
			Category:  ps[i].Category,
			CreatedAt: createdAt,
			Name:      ps[i].Name,
			GuardName: ps[i].GuardName,
		}
	}

	return perms, nil
}

func (p *Perms) Can(userInfo *userinfo.UserInfo, category Category, name Name) bool {
	permId, ok := p.lookupPermIDByGuard(BuildGuard(category, name))
	if !ok {
		return false
	}

	cached, ok := p.userCanCache.Get(userInfo.UserId)
	if ok {
		if result, found := cached[permId]; found {
			return result
		}
	}

	var result bool
	if userInfo.SuperUser {
		result = true
	} else {
		result = p.checkIfCan(permId, userInfo, category, name)
	}

	if cached == nil {
		cached = map[uint64]bool{}
	}
	cached[permId] = result

	p.userCanCache.Set(userInfo.UserId, cached,
		cache.WithExpiration(p.userCanCacheTTL))

	return result
}

func (p *Perms) checkIfCan(permId uint64, userInfo *userinfo.UserInfo, category Category, name Name) (result bool) {
	return p.checkRoleJob(userInfo.Job, userInfo.JobGrade, permId)
}

func (p *Perms) lookupPermIDByGuard(guard string) (uint64, bool) {
	return p.guardToPermIDMap.Load(guard)
}

func (p *Perms) lookupRoleIDForJobAndGrade(job string, grade int32) (uint64, bool) {
	grades, ok := p.jobsToRoleIDMap.Load(job)
	if !ok {
		return 0, false
	}
	roleId, ok := grades[grade]
	if !ok {
		return 0, false
	}

	return roleId, true
}

func (p *Perms) getRoleIDsForJobUpToGrade(job string, grade int32) ([]uint64, bool) {
	gradesMap, ok := p.jobsToRoleIDMap.Load(job)
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

func (p *Perms) checkRoleJob(job string, grade int32, permId uint64) bool {
	roleIds, ok := p.getRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		// Fallback to default role
		roleIds, ok = p.getRoleIDsForJobUpToGrade(DefaultRoleJob, DefaultRoleJobGrade)
		if !ok {
			return false
		}
	}

	for i := len(roleIds) - 1; i >= 0; i-- {
		ps, ok := p.rolePermsMap.Load(roleIds[i])
		if !ok {
			continue
		}
		val, ok := ps[permId]
		if !ok {
			continue
		}
		if val {
			return true
		}
	}

	return false
}
