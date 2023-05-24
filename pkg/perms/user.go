package perms

import (
	"fmt"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/query/fivenet/model"
)

func (p *Perms) GetPermissionsOfUser(userInfo *userinfo.UserInfo) (collections.Permissions, error) {
	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(userInfo.Job, userInfo.JobGrade)
	if !ok {
		roleId, ok := p.lookupRoleIDForJobAndGrade(DefaultRoleJob, DefaultRoleJobGrade)
		if !ok {
			return nil, fmt.Errorf("failed to fallback to default role")
		}
		roleIds = []uint64{roleId}
	}

	ps := p.getRolePermissionsFromCache(roleIds)
	perms := make(collections.Permissions, len(ps))

	if len(ps) == 0 {
		return perms, nil
	}

	for i := 0; i < len(ps); i++ {
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
	perms := map[uint64]interface{}{}
	for i := len(roleIds) - 1; i >= 0; i-- {
		permsMap, ok := p.permsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		permsMap.Range(func(key uint64, value bool) bool {
			if !value {
				return true
			}

			if _, ok := perms[key]; !ok {
				perms[key] = nil
			}

			return true
		})
	}

	ps := []*cachePerm{}
	for k := range perms {
		p, ok := p.lookupPermByID(k)
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
		result = p.checkIfCan(permId, userInfo, category, name)
	}

	p.userCanCache.Set(cacheKey, result,
		cache.WithExpiration(p.userCanCacheTTL))

	return result
}

func (p *Perms) checkIfCan(permId uint64, userInfo *userinfo.UserInfo, category Category, name Name) (result bool) {
	return p.checkRoleJob(userInfo.Job, userInfo.JobGrade, permId)
}

func (p *Perms) checkRoleJob(job string, grade int32, permId uint64) bool {
	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		// Fallback to default role
		roleIds, ok = p.lookupRoleIDsForJobUpToGrade(DefaultRoleJob, DefaultRoleJobGrade)
		if !ok {
			return false
		}
	}

	for i := len(roleIds) - 1; i >= 0; i-- {
		ps, ok := p.permsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}
		val, ok := ps.Load(permId)
		if !ok {
			continue
		}
		if val {
			return true
		}
	}

	return false
}
