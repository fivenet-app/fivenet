package perms

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms/collections"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/puzpuzpuz/xsync/v4"
)

const (
	DefaultRoleJob = "__default__"
)

var (
	tRoles     = table.FivenetRbacRoles
	tRolePerms = table.FivenetRbacRolesPermissions
	tJobPerms  = table.FivenetRbacJobPermissions
)

type AddPerm struct {
	Id  uint64
	Val bool
}

func (p *Perms) GetJobRoles(ctx context.Context, job string) (collections.Roles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.ID,
			tRoles.CreatedAt,
			tRoles.Job,
			tRoles.Grade,
		).
		FROM(tRoles).
		WHERE(
			tRoles.Job.EQ(jet.String(job)),
		).
		ORDER_BY(
			tRoles.Job.ASC(),
			tRoles.Grade.ASC(),
		)

	var dest collections.Roles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to get job roles for job %s. %w", job, err)
		}
	}

	return dest, nil
}

func (p *Perms) GetJobRolesUpTo(ctx context.Context, job string, grade int32) (collections.Roles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.ID,
			tRoles.CreatedAt,
			tRoles.Job,
			tRoles.Grade,
		).
		FROM(tRoles).
		WHERE(jet.AND(
			tRoles.Job.EQ(jet.String(job)),
			tRoles.Grade.LT_EQ(jet.Int32(grade)),
		)).
		ORDER_BY(
			tRoles.Job.ASC(),
			tRoles.Grade.ASC(),
		)

	var dest collections.Roles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to get job roles up to grade %d for job %s. %w", grade, job, err)
		}
	}

	return dest, nil
}

func (p *Perms) GetRoles(ctx context.Context, excludeSystem bool) (collections.Roles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.ID,
			tRoles.CreatedAt,
			tRoles.Job,
			tRoles.Grade,
		).
		FROM(tRoles).
		ORDER_BY(
			tRoles.Job.ASC(),
			tRoles.Grade.ASC(),
		)

	if excludeSystem {
		stmt = stmt.WHERE(tRoles.Job.NOT_EQ(jet.String(DefaultRoleJob)))
	}

	var dest collections.Roles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to get roles. %w", err)
		}
	}

	return dest, nil
}

func (p *Perms) GetClosestJobRole(ctx context.Context, job string, grade int32) (*model.FivenetRbacRoles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.ID,
			tRoles.CreatedAt,
			tRoles.Job,
			tRoles.Grade,
		).
		FROM(tRoles).
		WHERE(jet.AND(
			tRoles.Job.EQ(jet.String(job)),
			tRoles.Grade.LT_EQ(jet.Int32(grade)),
		)).
		LIMIT(1)

	var dest model.FivenetRbacRoles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to get closest job role for job %s and grade %d. %w", job, grade, err)
		}
	}

	return &dest, nil
}

func (p *Perms) CountRolesForJob(ctx context.Context, job string) (int64, error) {
	stmt := tRoles.
		SELECT(
			jet.COUNT(tRoles.ID).AS("data_count.total"),
		).
		FROM(tRoles).
		WHERE(
			tRoles.Job.EQ(jet.String(job)),
		)

	var dest database.DataCount
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return -1, fmt.Errorf("failed to count roles for job %s. %w", job, err)
		}
	}

	return dest.Total, nil
}

func (p *Perms) GetRole(ctx context.Context, id uint64) (*model.FivenetRbacRoles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.ID,
			tRoles.CreatedAt,
			tRoles.Job,
			tRoles.Grade,
		).
		FROM(tRoles).
		WHERE(
			tRoles.ID.EQ(jet.Uint64(id)),
		)

	var dest model.FivenetRbacRoles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to get role with ID %d. %w", id, err)
	}

	return &dest, nil
}

func (p *Perms) CreateRole(ctx context.Context, job string, grade int32) (*model.FivenetRbacRoles, error) {
	stmt := tRoles.
		INSERT(
			tRoles.Job,
			tRoles.Grade,
		).
		VALUES(
			job,
			grade,
		)

	res, err := stmt.ExecContext(ctx, p.db)
	if err != nil && !dbutils.IsDuplicateError(err) {
		return nil, fmt.Errorf("failed to create role for job %s and grade %d. %w", job, grade, err)
	}

	var role *model.FivenetRbacRoles
	if res != nil {
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve last insert ID for role creation. %w", err)
		}

		role, err = p.GetRole(ctx, uint64(lastId))
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve created role with ID %d. %w", lastId, err)
		}
	} else {
		role, err = p.GetRoleByJobAndGrade(ctx, job, grade)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve existing role for job %s and grade %d. %w", job, grade, err)
		}
	}

	p.permsRoleMap.Store(role.ID, xsync.NewMap[uint64, bool]())

	grades, _ := p.permsJobsRoleMap.LoadOrCompute(role.Job, func() (*xsync.Map[int32, uint64], bool) {
		return xsync.NewMap[int32, uint64](), false
	})
	grades.Store(role.Grade, role.ID)

	p.roleIDToJobMap.Store(role.ID, role.Job)

	if err := p.publishMessage(ctx, RoleCreatedSubject, RoleIDEvent{
		RoleID: role.ID,
		Job:    role.Job,
	}); err != nil {
		return nil, fmt.Errorf("failed to publish role creation message for role ID %d. %w", role.ID, err)
	}

	return role, nil
}

func (p *Perms) DeleteRole(ctx context.Context, id uint64) error {
	role, err := p.GetRole(ctx, id)
	if err != nil {
		// Role not found? It shouldn't exist anymore
		if errors.Is(err, qrm.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("failed to retrieve role for deletion with ID %d. %w", id, err)
	}

	stmt := tRoles.
		DELETE().
		WHERE(
			tRoles.ID.EQ(jet.Uint64(id)),
		).LIMIT(1)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("failed to delete role with ID %d. %w", id, err)
	}

	if err := p.publishMessage(ctx, RoleDeletedSubject, RoleIDEvent{
		RoleID: role.ID,
		Job:    role.Job,
		Grade:  role.Grade,
	}); err != nil {
		return fmt.Errorf("failed to publish role deletion message for role ID %d. %w", role.ID, err)
	}

	p.deleteRole(role.ID, role.Job, role.Grade)

	return nil
}

func (p *Perms) deleteRole(id uint64, job string, grade int32) {
	p.permsRoleMap.Delete(id)

	grades, _ := p.permsJobsRoleMap.LoadOrCompute(job, func() (*xsync.Map[int32, uint64], bool) {
		return xsync.NewMap[int32, uint64](), false
	})
	grades.Delete(grade)

	p.roleIDToJobMap.Delete(id)
}

func (p *Perms) GetRoleByJobAndGrade(ctx context.Context, job string, grade int32) (*model.FivenetRbacRoles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.ID,
			tRoles.CreatedAt,
			tRoles.Job,
			tRoles.Grade,
		).
		FROM(tRoles).
		WHERE(jet.AND(
			tRoles.Job.EQ(jet.String(job)),
			tRoles.Grade.EQ(jet.Int32(grade)),
		)).
		LIMIT(1)

	var dest model.FivenetRbacRoles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to get role for job %s and grade %d. %w", job, grade, err)
		}
	}

	return &dest, nil
}

func (p *Perms) GetRolePermissions(ctx context.Context, id uint64) ([]*permissions.Permission, error) {
	tRolePerms := tRolePerms
	tPerms := tPerms.AS("permission")
	stmt := tRolePerms.
		SELECT(
			tPerms.ID,
			tPerms.CreatedAt,
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
			tRolePerms.Val.AS("permission.val"),
		).
		FROM(
			tRolePerms.
				INNER_JOIN(tPerms,
					tPerms.ID.EQ(tRolePerms.PermissionID),
				),
		).
		WHERE(
			tRolePerms.RoleID.EQ(jet.Uint64(id)),
		).
		ORDER_BY(
			tPerms.ID.ASC(),
		)

	var dest []*permissions.Permission
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to get permissions for role ID %d. %w", id, err)
		}
	}

	return dest, nil
}

func (p *Perms) GetEffectiveRolePermissions(ctx context.Context, roleId uint64) ([]*permissions.Permission, error) {
	defaultRoleId, ok := p.lookupRoleIDForJobAndGrade(DefaultRoleJob, p.startJobGrade)
	if !ok {
		return nil, fmt.Errorf("failed to fallback to default role")
	}

	role, err := p.GetRole(ctx, roleId)
	if err != nil {
		return nil, err
	}

	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(role.Job, role.Grade)
	if !ok {
		// Fallback to default role
		roleIds = []uint64{defaultRoleId}
	} else {
		// Prepend default role to default perms
		roleIds = append([]uint64{defaultRoleId}, roleIds...)
	}

	perms := map[uint64]bool{}
	for i := range slices.Backward(roleIds) {
		permsRoleMap, ok := p.permsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		permsRoleMap.Range(func(key uint64, value bool) bool {
			// Only allow the first perm "value" to be set (because that's how role perms inheritance works)
			if _, ok := perms[key]; !ok {
				perms[key] = value
			}

			return true
		})
	}

	ps := []*permissions.Permission{}
	for i, v := range perms {
		p, ok := p.lookupPermByID(i)
		if !ok {
			continue
		}

		ps = append(ps, &permissions.Permission{
			Id:        p.ID,
			Category:  string(p.Category),
			Name:      string(p.Name),
			GuardName: p.GuardName,
			Val:       v,
			Order:     p.Order,
		})
	}

	// Order by `GuardName` ascending
	slices.SortFunc(ps, func(a, b *permissions.Permission) int {
		return strings.Compare(a.GuardName, b.GuardName)
	})

	return ps, nil
}

func (p *Perms) UpdateRolePermissions(ctx context.Context, roleId uint64, perms ...AddPerm) error {
	var rolePerms []struct {
		RoleID       uint64
		PermissionID uint64
		Val          bool
	}
	for _, perm := range perms {
		rolePerms = append(rolePerms, struct {
			RoleID       uint64
			PermissionID uint64
			Val          bool
		}{
			RoleID:       roleId,
			PermissionID: perm.Id,
			Val:          perm.Val,
		})
	}

	stmt := tRolePerms.
		INSERT(
			tRolePerms.RoleID,
			tRolePerms.PermissionID,
			tRolePerms.Val,
		).
		MODELS(rolePerms).
		ON_DUPLICATE_KEY_UPDATE(
			tRolePerms.Val.SET(jet.BoolExp(jet.Raw("VALUES(`val`)"))),
		)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return fmt.Errorf("failed to update permissions for role ID %d. %w", roleId, err)
		}
	}

	roleCache, _ := p.permsRoleMap.LoadOrCompute(roleId, func() (*xsync.Map[uint64, bool], bool) {
		return xsync.NewMap[uint64, bool](), false
	})
	for _, v := range rolePerms {
		roleCache.Store(v.PermissionID, v.Val)
	}

	if err := p.publishMessage(ctx, RolePermUpdateSubject, RoleIDEvent{
		RoleID: roleId,
	}); err != nil {
		return fmt.Errorf("failed to publish role permission update message for role ID %d. %w", roleId, err)
	}

	return nil
}

func (p *Perms) RemovePermissionsFromRole(ctx context.Context, roleId uint64, perms ...uint64) error {
	ids := make([]jet.Expression, len(perms))
	for i := range perms {
		ids[i] = jet.Uint64(perms[i])
	}

	stmt := tRolePerms.
		DELETE().
		WHERE(jet.AND(
			tRolePerms.RoleID.EQ(jet.Uint64(roleId)),
			tRolePerms.PermissionID.IN(ids...),
		)).
		LIMIT(int64(len(ids)))

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("failed to remove permissions from role ID %d. %w", roleId, err)
	}

	if permsRoleMap, ok := p.permsRoleMap.Load(roleId); ok {
		for _, permId := range perms {
			permsRoleMap.Delete(permId)
		}
	}

	if err := p.publishMessage(ctx, RolePermUpdateSubject, RoleIDEvent{
		RoleID: roleId,
	}); err != nil {
		return fmt.Errorf("failed to publish role permission removal message for role ID %d. %w", roleId, err)
	}

	return nil
}

func (p *Perms) GetJobPermissions(ctx context.Context, job string) ([]*permissions.Permission, error) {
	tPerms := tPerms.AS("permission")
	stmt := tJobPerms.
		SELECT(
			tJobPerms.PermissionID.AS("permission.id"),
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
			tJobPerms.Val.AS("permission.val"),
		).
		FROM(
			tJobPerms.
				INNER_JOIN(
					tPerms,
					tPerms.ID.EQ(tJobPerms.PermissionID),
				),
		).
		WHERE(
			tJobPerms.Job.EQ(jet.String(job)),
		)

	var dest []*permissions.Permission
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to get permissions for job %s. %w", job, err)
		}
	}

	return dest, nil
}

func (p *Perms) UpdateJobPermissions(ctx context.Context, job string, id uint64, val bool) error {
	stmt := tJobPerms.
		INSERT(
			tJobPerms.Job,
			tJobPerms.PermissionID,
			tJobPerms.Val,
		).
		VALUES(
			job,
			id,
			val,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobPerms.Val.SET(jet.RawBool("VALUES(`val`)")),
		)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return fmt.Errorf("failed to update job permissions for job %s and permission ID %d. %w", job, id, err)
		}
	}

	return nil
}

func (p *Perms) ClearJobPermissions(ctx context.Context, job string) error {
	stmt := tJobPerms.
		DELETE().
		WHERE(
			tJobPerms.Job.EQ(jet.String(job)),
		)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return fmt.Errorf("failed to clear job permissions for job %s. %w", job, err)
		}
	}

	if err := p.publishMessage(ctx, JobAttrUpdateSubject, RoleIDEvent{
		Job: job,
	}); err != nil {
		return fmt.Errorf("failed to publish job attribute update message for job %s. %w", job, err)
	}

	return nil
}
