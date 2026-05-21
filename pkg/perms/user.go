package perms

import (
	"slices"

	permissionspermissions "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/permissions"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/pkg/errors"
)

func (p *Perms) GetPermissionsOfUser(
	userInfo *userinfo.UserInfo,
) ([]*permissionspermissions.Permission, error) {
	defaultRoleId, ok := p.lookupRoleIDForJobAndGrade(DefaultRoleJob, p.startJobGrade)
	if !ok {
		return nil, errors.New("failed to fallback to default role")
	}

	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(userInfo.GetJob(), userInfo.GetJobGrade())
	if !ok {
		// Fallback to default role
		roleIds = []int64{defaultRoleId}
	} else {
		// Prepend default role to default perms
		roleIds = append([]int64{defaultRoleId}, roleIds...)
	}

	ps := p.getRolePermissionsFromCache(roleIds)
	if len(ps) == 0 {
		return nil, nil
	}

	perms := make([]*permissionspermissions.Permission, len(ps))
	for i := range ps {
		perms[i] = &permissionspermissions.Permission{
			Id:        ps[i].ID,
			Namespace: string(ps[i].Namespace),
			Service:   string(ps[i].Service),
			Name:      string(ps[i].Name),
			GuardName: ps[i].GuardName,
		}
	}

	return perms, nil
}

func (p *Perms) getRolePermissionsFromCache(roleIds []int64) []*cachePerm {
	perms := map[int64]bool{}
	for i := range slices.Backward(roleIds) {
		permsRoleMap, ok := p.permsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		for key, value := range permsRoleMap.All() {
			// Only allow the perm "value" to be set once (because that's how role perms inheritance works)
			if _, ok := perms[key]; !ok {
				perms[key] = value
			}
		}
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

func (ps *Perms) can(
	userInfo *userinfo.UserInfo,
	namespace Namespace,
	service Service,
	name Name,
) bool {
	permId, ok := ps.lookupPermIDByGuard(BuildGuard(namespace, service, name))
	if !ok {
		return false
	}

	// Don't check permissions for superusers and don't cache the result
	if userInfo.GetSuperuser() {
		return true
	}

	cacheKey := userCacheKey{
		userId: userInfo.GetUserId(),
		job:    userInfo.GetJob(),
		grade:  userInfo.GetJobGrade(),
		permId: permId,
	}
	result, ok := ps.userCanCache.Get(cacheKey)
	if ok {
		return result
	}

	result = ps.checkIfCan(permId, userInfo)

	ps.userCanCache.Put(cacheKey, result, ps.userCanCacheTTL)

	return result
}

func (ps *Perms) checkIfCan(permId int64, userInfo *userinfo.UserInfo) bool {
	if check, ok := ps.checkRoleJob(userInfo.GetJob(), userInfo.GetJobGrade(), permId); ok {
		return check
	}

	// Check default role perms
	check, _ := ps.checkRoleJob(DefaultRoleJob, ps.startJobGrade, permId)
	return check
}

func (ps *Perms) checkRoleJob(job string, grade int32, permId int64) (bool, bool) {
	roleIds, ok := ps.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return false, false
	}

	for i := range slices.Backward(roleIds) {
		ps, ok := ps.permsRoleMap.Load(roleIds[i])
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

func (ps *Perms) clearUserCanCache() {
	if ps.userCanCache != nil {
		ps.userCanCache.Clear()
	}
}
