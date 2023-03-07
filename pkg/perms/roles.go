package perms

import (
	"context"
	"errors"

	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-sql-driver/mysql"
)

func (p *perms) CreateRole(name string, description string) error {
	r := table.ArpanetRoles
	stmt := r.INSERT(r.Name,
		r.GuardName,
		r.Description).
		VALUES(name, helpers.Guard(name), description)

	_, err := stmt.ExecContext(context.TODO(), query.DB)
	return err
}

func (p *perms) DeleteRole(name string) error {
	r := table.ArpanetRoles
	_, err := r.DELETE().
		WHERE(r.GuardName.EQ(jet.String(helpers.Guard(name)))).
		ExecContext(context.TODO(), query.DB)
	return err
}

func (p *perms) AddPermissionsToRole(name string, perms collections.Permissions) error {
	r := table.ArpanetRoles
	var roleIDs []uint
	if err := r.SELECT(r.ID).
		FROM(r).
		WHERE(r.GuardName.EQ(jet.String(helpers.Guard(name)))).
		LIMIT(1).
		QueryContext(context.TODO(), query.DB, &roleIDs); err != nil {
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

	rp := table.ArpanetRolePermissions
	_, err := rp.INSERT(rp.PermissionID, rp.RoleID).
		MODELS(rolePerms).
		ExecContext(context.TODO(), query.DB)

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number != 1062 {
		return err
	}

	return nil
}
