package perms

import (
	"github.com/galexrt/fivenet/pkg/dbutils"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/pkg/perms/helpers"
	"github.com/galexrt/fivenet/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (p *Perms) GetRoles(prefix string) (collections.Roles, error) {
	prefix = helpers.Guard(prefix)

	stmt := ar.
		SELECT(
			ar.AllColumns,
		).
		FROM(ar).
		WHERE(
			ar.GuardName.LIKE(jet.String(prefix+"%")),
		).
		ORDER_BY(
			jet.LENGTH(ar.GuardName),
			ar.GuardName.ASC(),
		)

	var dest collections.Roles
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (p *Perms) GetRole(id uint64) (*model.FivenetRoles, error) {
	stmt := ar.
		SELECT(
			ar.AllColumns,
		).
		FROM(ar).
		WHERE(
			ar.ID.EQ(jet.Uint64(id)),
		)

	var dest model.FivenetRoles
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (p *Perms) GetRoleByGuardName(name string) (*model.FivenetRoles, error) {
	stmt := ar.
		SELECT(
			ar.AllColumns,
		).
		FROM(ar).
		WHERE(
			ar.GuardName.EQ(jet.String(helpers.Guard(name))),
		)

	var dest model.FivenetRoles
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (p *Perms) GetRolePermissions(id uint64) (collections.Permissions, error) {
	stmt := arp.
		SELECT(
			ap.AllColumns,
		).
		FROM(
			arp.
				INNER_JOIN(ap,
					ap.ID.EQ(arp.PermissionID),
				),
		).
		WHERE(
			arp.RoleID.EQ(jet.Uint64(id)),
		).
		ORDER_BY(
			ap.Name.ASC(),
			ap.ID.ASC(),
		)

	var dest collections.Permissions
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (p *Perms) CreateRoleWithGuard(name string, guard string, description string) (*model.FivenetRoles, error) {
	stmt := ar.
		INSERT(
			ar.Name,
			ar.GuardName,
			ar.Description,
		).
		VALUES(name, guard, description)

	res, err := stmt.ExecContext(p.ctx, p.db)
	if !dbutils.IsDuplicateError(err) {
		return nil, err
	}

	var role *model.FivenetRoles
	if res != nil {
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		role, err = p.GetRole(uint64(lastId))
		if err != nil {
			return nil, err
		}
	} else {
		role, err = p.GetRoleByGuardName(guard)
		if err != nil {
			return nil, err
		}
	}

	return role, nil
}

func (p *Perms) CreateRole(name string, description string) (*model.FivenetRoles, error) {
	return p.CreateRoleWithGuard(name, helpers.Guard(name), description)
}

func (p *Perms) DeleteRole(id uint64) error {
	_, err := ar.
		DELETE().
		WHERE(
			ar.ID.EQ(jet.Uint64(id)),
		).
		ExecContext(p.ctx, p.db)
	return err
}

func (p *Perms) AddPermissionsToRole(id uint64, perms []uint64) error {
	role, err := p.GetRole(id)
	if err != nil {
		return err
	}

	var rolePerms []struct {
		PermissionID uint64
		RoleID       uint64
	}
	for _, permID := range perms {
		rolePerms = append(rolePerms, struct {
			PermissionID uint64
			RoleID       uint64
		}{
			PermissionID: permID,
			RoleID:       role.ID,
		})
	}

	_, err = arp.
		INSERT(
			arp.PermissionID,
			arp.RoleID,
		).
		MODELS(rolePerms).
		ExecContext(p.ctx, p.db)

	if !dbutils.IsDuplicateError(err) {
		return err
	}

	return nil
}

func (p *Perms) RemovePermissionsFromRole(id uint64, perms []uint64) error {
	ids := make([]jet.Expression, len(perms))
	for i := 0; i < len(perms); i++ {
		ids[i] = jet.Uint64(perms[i])
	}

	stmt := arp.
		DELETE().
		WHERE(jet.AND(
			arp.RoleID.EQ(jet.Uint64(id)),
			arp.PermissionID.IN(ids...),
		))

	_, err := stmt.ExecContext(p.ctx, p.db)

	return err
}
