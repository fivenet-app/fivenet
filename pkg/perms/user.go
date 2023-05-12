package perms

import (
	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var tUserPerms = table.FivenetUserPermissions

func (p *Perms) GetPermissionsOfUser(userId int32, job string, grade int32) (collections.Permissions, error) {
	stmt := tPerms.
		SELECT(
			tPerms.AllColumns,
		).
		FROM(tPerms).
		WHERE(
			tPerms.ID.IN(
				tUserPerms.
					SELECT(
						tUserPerms.PermissionID,
					).
					FROM(
						tUserPerms,
					).
					WHERE(
						tUserPerms.UserID.EQ(jet.Int32(userId)),
					).
					UNION(
						tRoles.
							SELECT(tRolePerms.PermissionID).
							FROM(tRoles.
								INNER_JOIN(tRolePerms,
									tRolePerms.RoleID.EQ(tRoles.ID)),
							).
							WHERE(jet.AND(
								tRoles.Job.EQ(jet.String(job)),
								tRoles.Grade.EQ(jet.Int32(grade)),
							)),
					),
			),
		)

	var perms collections.Permissions
	if err := stmt.QueryContext(p.ctx, p.db, &perms); err != nil {
		return nil, err
	}

	return perms, nil
}

func (p *Perms) Can(userId int32, job string, grade int32, category Category, name Name) bool {
	permId, ok := p.lookupPermIDByGuard(BuildGuard(category, name))
	if !ok {
		return false
	}

	cached, ok := p.userCanCache.Get(userId)
	if ok {
		if result, found := cached[permId]; found {
			return result
		}
	}

	result := p.checkIfCan(permId, userId, job, grade, category, name)

	if !result {
		result = p.checkIfCan(permId, userId, job, grade, common.SuperuserCategoryPerm, common.SuperuserAnyAccessName)
	}

	if cached == nil {
		cached = map[uint64]bool{}
	}
	cached[permId] = result

	p.userCanCache.Set(userId, cached,
		cache.WithExpiration(p.userCanCacheTTL))

	return result
}

func (p *Perms) checkIfCan(permId uint64, userId int32, job string, grade int32, category Category, name Name) (result bool) {
	if p.checkRoleJob(job, grade, permId) {
		return true
	}

	return p.checkIfUserCan(userId, permId)
}

func (p *Perms) checkIfUserCan(userId int32, permId uint64) bool {
	stmt :=
		tUserPerms.
			SELECT(
				tUserPerms.PermissionID.AS("id"),
			).
			FROM(tUserPerms).
			WHERE(jet.AND(
				tUserPerms.UserID.EQ(jet.Int32(userId)),
				tUserPerms.PermissionID.EQ(jet.Uint64(permId)),
			)).
			LIMIT(1)

	var dest struct {
		ID int32
	}
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		return false
	}

	return dest.ID > 0
}

func (p *Perms) lookupPermIDByGuard(guard string) (uint64, bool) {
	return p.guardToPermIDMap.Load(guard)
}

func (p *Perms) getRoleIDForJobAndGrade(job string, grade int32) (uint64, bool) {
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
		if val == 1 {
			return true
		}
	}

	return false
}
