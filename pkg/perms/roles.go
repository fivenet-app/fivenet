package perms

import (
	"errors"

	"github.com/galexrt/arpanet/pkg/dbutils"
	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-sql-driver/mysql"
)

func (p *Perms) GetRoles(prefix string) (collections.Roles, error) {
	prefix = helpers.Guard(prefix)

	stmt := ar.
		SELECT(
			ar.AllColumns,
		).
		FROM(ar).
		WHERE(
			ar.GuardName.LIKE(jet.String(prefix + "%")),
		)

	var dest collections.Roles
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (p *Perms) CreateRole(name string, description string) error {
	stmt := ar.
		INSERT(
			ar.Name,
			ar.GuardName,
			ar.Description,
		).
		VALUES(name, helpers.Guard(name), description)

	_, err := stmt.ExecContext(p.ctx, p.db)

	if !dbutils.IsDuplicateError(err) {
		return err
	}

	return nil
}

func (p *Perms) DeleteRole(name string) error {
	_, err := ar.
		DELETE().
		WHERE(
			ar.GuardName.EQ(jet.String(helpers.Guard(name))),
		).
		ExecContext(p.ctx, p.db)
	return err
}

func (p *Perms) AddPermissionsToRole(name string, perms collections.Permissions) error {
	// Make sure the role exists
	rolesStmt := ar.
		SELECT(ar.ID).
		FROM(ar).
		WHERE(
			ar.GuardName.EQ(jet.String(helpers.Guard(name))),
		).
		LIMIT(1)

	var roleIDs []uint
	if err := rolesStmt.QueryContext(p.ctx, p.db, &roleIDs); err != nil {
		return err
	}

	var rolePerms []struct {
		PermissionID uint
		RoleID       uint
	}
	for _, permID := range perms.IDs() {
		rolePerms = append(rolePerms, struct {
			PermissionID uint
			RoleID       uint
		}{
			PermissionID: uint(permID),
			RoleID:       roleIDs[0],
		})
	}

	_, err := arp.
		INSERT(
			arp.PermissionID,
			arp.RoleID,
		).
		MODELS(rolePerms).
		ExecContext(p.ctx, p.db)

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number != 1062 {
		return err
	}

	return nil
}
