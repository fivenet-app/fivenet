package perms

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/pkg/perms/collections"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/puzpuzpuz/xsync/v3"
)

const (
	DefaultRoleJob      = "__default__"
	DefaultRoleJobGrade = int32(1)
)

var (
	tRoles     = table.FivenetRoles
	tRolePerms = table.FivenetRolePermissions
	tJobPerms  = table.FivenetJobPermissions
)

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
			return nil, err
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
			return nil, err
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
			return nil, err
		}
	}

	return dest, nil
}

func (p *Perms) GetClosestJobRole(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error) {
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

	var dest model.FivenetRoles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (p *Perms) CountRolesForJob(ctx context.Context, job string) (int64, error) {
	stmt := tRoles.
		SELECT(
			jet.COUNT(tRoles.ID).AS("datacount.totalcount"),
		).
		FROM(tRoles).
		WHERE(
			tRoles.Job.EQ(jet.String(job)),
		)

	var dest database.DataCount
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return -1, err
		}
	}

	return dest.TotalCount, nil
}

func (p *Perms) GetRole(ctx context.Context, id uint64) (*model.FivenetRoles, error) {
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

	var dest model.FivenetRoles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (p *Perms) CreateRole(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error) {
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
		return nil, err
	}

	var role *model.FivenetRoles
	if res != nil {
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		role, err = p.GetRole(ctx, uint64(lastId))
		if err != nil {
			return nil, err
		}
	} else {
		role, err = p.GetRoleByJobAndGrade(ctx, job, grade)
		if err != nil {
			return nil, err
		}
	}

	p.permsRoleMap.Store(role.ID, xsync.NewMapOf[uint64, bool]())

	grades, _ := p.permsJobsRoleMap.LoadOrCompute(role.Job, func() *xsync.MapOf[int32, uint64] {
		return xsync.NewMapOf[int32, uint64]()
	})
	grades.Store(role.Grade, role.ID)

	p.roleIDToJobMap.Store(role.ID, role.Job)

	if err := p.publishMessage(ctx, RoleCreatedSubject, RoleIDEvent{
		RoleID: role.ID,
	}); err != nil {
		return nil, err
	}

	return role, nil
}

func (p *Perms) DeleteRole(ctx context.Context, id uint64) error {
	role, err := p.GetRole(ctx, id)
	if err != nil {
		return err
	}

	p.deleteRole(role)

	if err := p.publishMessage(ctx, RoleDeletedSubject, RoleIDEvent{
		RoleID: role.ID,
	}); err != nil {
		return err
	}

	stmt := tRoles.
		DELETE().
		WHERE(
			tRoles.ID.EQ(jet.Uint64(id)),
		).LIMIT(1)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return err
	}

	return nil
}

func (p *Perms) deleteRole(role *model.FivenetRoles) {
	p.permsRoleMap.Delete(role.ID)

	grades, _ := p.permsJobsRoleMap.LoadOrCompute(role.Job, func() *xsync.MapOf[int32, uint64] {
		return xsync.NewMapOf[int32, uint64]()
	})
	grades.Delete(role.Grade)

	p.roleIDToJobMap.Delete(role.ID)
}

func (p *Perms) GetRoleByJobAndGrade(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error) {
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

	var dest model.FivenetRoles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return nil, err
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
			return nil, err
		}
	}

	return dest, nil
}

type AddPerm struct {
	Id  uint64
	Val bool
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
			return err
		}
	}

	roleCache, _ := p.permsRoleMap.LoadOrCompute(roleId, func() *xsync.MapOf[uint64, bool] {
		return xsync.NewMapOf[uint64, bool]()
	})
	for _, v := range rolePerms {
		roleCache.Store(v.PermissionID, v.Val)
	}

	if err := p.publishMessage(ctx, RolePermUpdateSubject, RoleIDEvent{
		RoleID: roleId,
	}); err != nil {
		return err
	}

	return nil
}

func (p *Perms) RemovePermissionsFromRole(ctx context.Context, roleId uint64, perms ...uint64) error {
	ids := make([]jet.Expression, len(perms))
	for i := 0; i < len(perms); i++ {
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
		return err
	}

	if permsRoleMap, ok := p.permsRoleMap.Load(roleId); ok {
		for _, permId := range perms {
			permsRoleMap.Delete(permId)
		}
	}

	if err := p.publishMessage(ctx, RolePermUpdateSubject, RoleIDEvent{
		RoleID: roleId,
	}); err != nil {
		return err
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
			return nil, err
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
			return err
		}
	}

	return nil
}
