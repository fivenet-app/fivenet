package perms

import (
	"context"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	DefaultRoleJob      = "__default__"
	DefaultRoleJobGrade = int32(1)
)

var (
	tRoles     = table.FivenetRoles
	tRolePerms = table.FivenetRolePermissions
)

func (p *Perms) GetJobRoles(ctx context.Context, job string) (collections.Roles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.AllColumns,
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
		return nil, err
	}

	return dest, nil
}

func (p *Perms) GetJobRolesUpTo(ctx context.Context, job string, grade int32) (collections.Roles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.AllColumns,
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
		return nil, err
	}

	return dest, nil
}

func (p *Perms) getRoles(ctx context.Context) (collections.Roles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.AllColumns,
		).
		FROM(tRoles).
		ORDER_BY(
			tRoles.Job.ASC(),
			tRoles.Grade.ASC(),
		)

	var dest collections.Roles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (p *Perms) GetClosestJobRole(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.AllColumns,
		).
		FROM(tRoles).
		WHERE(jet.AND(
			tRoles.Job.EQ(jet.String(job)),
			tRoles.Grade.LT_EQ(jet.Int32(grade)),
		)).
		LIMIT(1)

	var dest model.FivenetRoles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return nil, err
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
		return -1, err
	}

	return dest.TotalCount, nil
}

func (p *Perms) GetRole(ctx context.Context, id uint64) (*model.FivenetRoles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.AllColumns,
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

	p.permsRoleMap.Store(role.ID, &syncx.Map[uint64, bool]{})

	return role, nil
}

func (p *Perms) DeleteRole(ctx context.Context, id uint64) error {
	_, err := tRoles.
		DELETE().
		WHERE(
			tRoles.ID.EQ(jet.Uint64(id)),
		).
		ExecContext(ctx, p.db)
	if err != nil {
		return err
	}

	p.permsRoleMap.Delete(id)

	return nil
}

func (p *Perms) GetRoleByJobAndGrade(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error) {
	stmt := tRoles.
		SELECT(
			tRoles.AllColumns,
		).
		FROM(tRoles).
		WHERE(jet.AND(
			tRoles.Job.EQ(jet.String(job)),
			tRoles.Grade.EQ(jet.Int32(grade)),
		))

	var dest model.FivenetRoles
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if errors.Is(qrm.ErrNoRows, err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &dest, nil
}

func (p *Perms) GetRolePermissions(ctx context.Context, id uint64) ([]*permissions.Permission, error) {
	tRolePerms := tRolePerms
	tPerms := tPerms.AS("permission")
	stmt := tRolePerms.
		SELECT(
			tPerms.AllColumns,
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
			tPerms.Name.ASC(),
			tPerms.ID.ASC(),
		)

	var dest []*permissions.Permission
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return nil, err
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
			tRolePerms.Val.SET(jet.BoolExp(jet.Raw("values(`val`)"))),
		)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil && !dbutils.IsDuplicateError(err) {
		return err
	}

	roleCache, ok := p.permsRoleMap.Load(roleId)
	if !ok {
		return nil
	}

	for _, v := range rolePerms {
		roleCache.Store(v.PermissionID, v.Val)
	}

	if err := p.publishMessage(RolePermUpdateSubject, RolePermUpdateEvent{
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
		))

	_, err := stmt.ExecContext(ctx, p.db)
	if err != nil {
		return err
	}

	roleCache, ok := p.permsRoleMap.Load(roleId)
	if !ok {
		return nil
	}

	for _, permId := range perms {
		roleCache.Delete(permId)
	}

	if err := p.publishMessage(RolePermUpdateSubject, RolePermUpdateEvent{
		RoleID: roleId,
	}); err != nil {
		return err
	}

	return nil
}
