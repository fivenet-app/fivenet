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
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/puzpuzpuz/xsync/v4"
)

const (
	DefaultRoleJob = "__default__"
)

var (
	tRoles     = table.FivenetRbacRoles.AS("role")
	tRolePerms = table.FivenetRbacRolesPermissions
	tJobPerms  = table.FivenetRbacJobPermissions
)

type AddPerm struct {
	Id  int64
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

func (p *Perms) GetJobRolesUpTo(
	ctx context.Context,
	job string,
	grade int32,
) (collections.Roles, error) {
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
			return nil, fmt.Errorf(
				"failed to get job roles up to grade %d for job %s. %w",
				grade,
				job,
				err,
			)
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

func (p *Perms) GetClosestJobRole(
	ctx context.Context,
	job string,
	grade int32,
) (*permissions.Role, error) {
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

	var dest permissions.Role
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf(
				"failed to get closest job role for job %s and grade %d. %w",
				job,
				grade,
				err,
			)
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

func (p *Perms) GetRole(ctx context.Context, id int64) (*permissions.Role, error) {
	stmt := tRoles.
		SELECT(
			tRoles.ID,
			tRoles.CreatedAt,
			tRoles.Job,
			tRoles.Grade,
		).
		FROM(tRoles).
		WHERE(
			tRoles.ID.EQ(jet.Int64(id)),
		).
		LIMIT(1)

	var dest permissions.Role
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to get role with ID %d. %w", id, err)
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return &dest, nil
}

func (p *Perms) CreateRole(
	ctx context.Context,
	job string,
	grade int32,
) (*permissions.Role, error) {
	tRoles := table.FivenetRbacRoles
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

	var role *permissions.Role
	if res != nil {
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve last insert ID for role creation. %w", err)
		}

		role, err = p.GetRole(ctx, lastId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve created role with ID %d. %w", lastId, err)
		}
	} else {
		role, err = p.GetRoleByJobAndGrade(ctx, job, grade)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve existing role for job %s and grade %d. %w", job, grade, err)
		}
	}

	p.permsRoleMap.Store(role.GetId(), xsync.NewMap[int64, bool]())

	grades, _ := p.permsJobsRoleMap.LoadOrCompute(
		role.GetJob(),
		func() (*xsync.Map[int32, int64], bool) {
			return xsync.NewMap[int32, int64](), false
		},
	)
	grades.Store(role.GetGrade(), role.GetId())

	p.roleIDToJobMap.Store(role.GetId(), role.GetJob())

	if err := p.publishMessage(ctx, RoleCreatedSubject, &permissions.RoleIDEvent{
		RoleId: role.GetId(),
		Job:    role.GetJob(),
		Grade:  role.GetGrade(),
	}); err != nil {
		return nil, fmt.Errorf(
			"failed to publish role creation message for role ID %d. %w",
			role.GetId(),
			err,
		)
	}

	return role, nil
}

func (p *Perms) DeleteRole(ctx context.Context, id int64) error {
	role, err := p.GetRole(ctx, id)
	if err != nil {
		// Role not found? It shouldn't exist anymore
		if errors.Is(err, qrm.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("failed to retrieve role for deletion with ID %d. %w", id, err)
	}

	tRoles := table.FivenetRbacRoles

	stmt := tRoles.
		DELETE().
		WHERE(
			tRoles.ID.EQ(jet.Int64(id)),
		).LIMIT(1)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("failed to delete role with ID %d. %w", id, err)
	}

	if err := p.publishMessage(ctx, RoleDeletedSubject, &permissions.RoleIDEvent{
		RoleId: role.GetId(),
		Job:    role.GetJob(),
		Grade:  role.GetGrade(),
	}); err != nil {
		return fmt.Errorf(
			"failed to publish role deletion message for role ID %d. %w",
			role.GetId(),
			err,
		)
	}

	p.deleteRole(role.GetId(), role.GetJob(), role.GetGrade())

	return nil
}

func (p *Perms) deleteRole(id int64, job string, grade int32) {
	p.permsRoleMap.Delete(id)

	grades, _ := p.permsJobsRoleMap.LoadOrCompute(job, func() (*xsync.Map[int32, int64], bool) {
		return xsync.NewMap[int32, int64](), false
	})
	grades.Delete(grade)

	p.roleIDToJobMap.Delete(id)
}

func (p *Perms) GetRoleByJobAndGrade(
	ctx context.Context,
	job string,
	grade int32,
) (*permissions.Role, error) {
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

	var dest permissions.Role
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to get role for job %s and grade %d. %w", job, grade, err)
		}
	}

	return &dest, nil
}

func (p *Perms) GetRolePermissions(
	ctx context.Context,
	id int64,
) ([]*permissions.Permission, error) {
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
			tRolePerms.RoleID.EQ(jet.Int64(id)),
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

func (p *Perms) GetEffectiveRolePermissions(
	ctx context.Context,
	roleId int64,
) ([]*permissions.Permission, error) {
	defaultRoleId, ok := p.lookupRoleIDForJobAndGrade(DefaultRoleJob, p.startJobGrade)
	if !ok {
		return nil, errors.New("failed to fallback to default role")
	}

	role, err := p.GetRole(ctx, roleId)
	if err != nil {
		return nil, err
	}

	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(role.GetJob(), role.GetGrade())
	if !ok {
		// Fallback to default role
		roleIds = []int64{defaultRoleId}
	} else {
		// Prepend default role to default perms
		roleIds = append([]int64{defaultRoleId}, roleIds...)
	}

	perms := map[int64]bool{}
	for i := range slices.Backward(roleIds) {
		permsRoleMap, ok := p.permsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		permsRoleMap.Range(func(key int64, value bool) bool {
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
		return strings.Compare(a.GetGuardName(), b.GetGuardName())
	})

	return ps, nil
}

func (p *Perms) UpdateRolePermissions(ctx context.Context, roleId int64, perms ...AddPerm) error {
	rolePerms := make([]struct {
		RoleID       int64
		PermissionID int64
		Val          bool
	}, len(perms))
	for i, perm := range perms {
		rolePerms[i] = struct {
			RoleID       int64
			PermissionID int64
			Val          bool
		}{
			RoleID:       roleId,
			PermissionID: perm.Id,
			Val:          perm.Val,
		}
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

	roleCache, _ := p.permsRoleMap.LoadOrCompute(roleId, func() (*xsync.Map[int64, bool], bool) {
		return xsync.NewMap[int64, bool](), false
	})
	for _, v := range rolePerms {
		roleCache.Store(v.PermissionID, v.Val)
	}

	if err := p.publishMessage(ctx, RolePermUpdateSubject, &permissions.RoleIDEvent{
		RoleId: roleId,
	}); err != nil {
		return fmt.Errorf(
			"failed to publish role permission update message for role ID %d. %w",
			roleId,
			err,
		)
	}

	return nil
}

func (p *Perms) RemovePermissionsFromRole(
	ctx context.Context,
	roleId int64,
	perms ...int64,
) error {
	ids := make([]jet.Expression, len(perms))
	for i := range perms {
		ids[i] = jet.Int64(perms[i])
	}

	stmt := tRolePerms.
		DELETE().
		WHERE(jet.AND(
			tRolePerms.RoleID.EQ(jet.Int64(roleId)),
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

	if err := p.publishMessage(ctx, RolePermUpdateSubject, &permissions.RoleIDEvent{
		RoleId: roleId,
	}); err != nil {
		return fmt.Errorf(
			"failed to publish role permission removal message for role ID %d. %w",
			roleId,
			err,
		)
	}

	return nil
}

func (p *Perms) GetJobPermissions(
	ctx context.Context,
	job string,
) ([]*permissions.Permission, error) {
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

func (p *Perms) UpdateJobPermissions(
	ctx context.Context,
	job string,
	perms ...*permissions.PermItem,
) error {
	for _, ps := range perms {
		stmt := tJobPerms.
			INSERT(
				tJobPerms.Job,
				tJobPerms.PermissionID,
				tJobPerms.Val,
			).
			VALUES(
				job,
				ps.GetId(),
				ps.GetVal(),
			).
			ON_DUPLICATE_KEY_UPDATE(
				tJobPerms.Val.SET(jet.RawBool("VALUES(`val`)")),
			)

		if _, err := stmt.ExecContext(ctx, p.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return fmt.Errorf(
					"failed to update job permissions for job %s and permission ID %d. %w",
					job,
					ps.GetId(),
					err,
				)
			}
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

	return nil
}
