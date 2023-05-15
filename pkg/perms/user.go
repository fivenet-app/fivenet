package perms

import (
	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (p *Perms) GetPermissionsOfUser(userInfo *userinfo.UserInfo) (collections.Permissions, error) {
	stmt := tPerms.
		SELECT(
			tPerms.AllColumns,
		).
		FROM(tPerms).
		WHERE(
			tPerms.ID.IN(
				tRoles.
					SELECT(tRolePerms.PermissionID).
					FROM(tRoles.
						INNER_JOIN(tRolePerms,
							tRolePerms.RoleID.EQ(tRoles.ID)),
					).
					WHERE(jet.AND(
						tRoles.Job.EQ(jet.String(userInfo.Job)),
						tRoles.Grade.EQ(jet.Int32(userInfo.JobGrade)),
					)),
			),
		)

	var perms collections.Permissions
	if err := stmt.QueryContext(p.ctx, p.db, &perms); err != nil {
		return nil, err
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

func (p *Perms) getRoleIDsForJobUpToGrade(job string, grade int32) ([]uint64, bool) {
	grades, ok := p.jobsToRoleIDMap.Load(job)
	if !ok {
		return nil, false
	}

	gradeList := []uint64{}
	for g, value := range grades {
		if g > grade {
			continue
		}
		gradeList = append(gradeList, value)
	}

	return gradeList, true
}

func (p *Perms) checkRoleJob(job string, grade int32, permId uint64) bool {
	roleIds, ok := p.getRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return false
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
