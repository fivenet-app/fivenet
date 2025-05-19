package perms

import (
	"fmt"
	"slices"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms/collections"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
)

func (p *Perms) GetPermissionsOfUser(userInfo *userinfo.UserInfo) (collections.Permissions, error) {
	defaultRoleId, ok := p.lookupRoleIDForJobAndGrade(DefaultRoleJob, p.startJobGrade)
	if !ok {
		return nil, fmt.Errorf("failed to fallback to default role")
	}

	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(userInfo.Job, userInfo.JobGrade)
	if !ok {
		// Fallback to default role
		roleIds = []uint64{defaultRoleId}
	} else {
		// Prepend default role to default perms
		roleIds = append([]uint64{defaultRoleId}, roleIds...)
	}

	ps := p.getRolePermissionsFromCache(roleIds)
	if len(ps) == 0 {
		return nil, nil
	}

	perms := make(collections.Permissions, len(ps))
	for i := range ps {
		perms[i] = &model.FivenetPermissions{
			ID:        ps[i].ID,
			Category:  string(ps[i].Category),
			Name:      string(ps[i].Name),
			GuardName: ps[i].GuardName,
		}
	}

	return perms, nil
}

func (p *Perms) getRolePermissionsFromCache(roleIds []uint64) []*cachePerm {
	perms := map[uint64]bool{}
	for i := range slices.Backward(roleIds) {
		permsRoleMap, ok := p.permsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		permsRoleMap.Range(func(key uint64, value bool) bool {
			// Only allow the perm "value" to be set once (because that's how role perms inheritance works)
			if _, ok := perms[key]; !ok {
				perms[key] = value
			}

			return true
		})
	}

	ps := []*cachePerm{}
	for i, v := range perms {
		if !v {
			continue
		}

		p, ok := p.lookupPermByID(i)
		if !ok {
			continue
		}

		ps = append(ps, p)
	}

	return ps
}

func (p *Perms) Can(userInfo *userinfo.UserInfo, category Category, name Name) bool {
	permId, ok := p.lookupPermIDByGuard(BuildGuard(category, name))
	if !ok {
		return false
	}

	cacheKey := userCacheKey{
		userId: userInfo.UserId,
		permId: permId,
	}
	result, ok := p.userCanCache.Get(cacheKey)
	if ok {
		return result
	}

	if userInfo.SuperUser {
		result = true
	} else {
		result = p.checkIfCan(permId, userInfo)
	}

	p.userCanCache.Set(cacheKey, result,
		cache.WithExpiration(p.userCanCacheTTL))

	return result
}

func (p *Perms) checkIfCan(permId uint64, userInfo *userinfo.UserInfo) bool {
	if check, ok := p.checkRoleJob(userInfo.Job, userInfo.JobGrade, permId); ok {
		return check
	}

	// Check default role perms
	check, _ := p.checkRoleJob(DefaultRoleJob, p.startJobGrade, permId)
	return check
}

func (p *Perms) checkRoleJob(job string, grade int32, permId uint64) (bool, bool) {
	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return false, false
	}

	for i := range slices.Backward(roleIds) {
		ps, ok := p.permsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		val, ok := ps.Load(permId)
		if !ok {
			continue
		}

		return val, true
	}

	return false, false
}
