package perms

import (
	"context"
	"errors"

	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-sql-driver/mysql"
)

func (p *Perms) CreateRole(name string, description string) error {
	stmt := ar.INSERT(
		ar.Name,
		ar.GuardName,
		ar.Description,
	).
		VALUES(name, helpers.Guard(name), description)

	_, err := stmt.ExecContext(context.TODO(), p.db)
	return err
}

func (p *Perms) DeleteRole(name string) error {
	_, err := ar.DELETE().
		WHERE(ar.GuardName.EQ(jet.String(helpers.Guard(name)))).
		ExecContext(context.TODO(), p.db)
	return err
}

func (p *Perms) AddPermissionsToRole(name string, perms collections.Permissions) error {
	var roleIDs []uint
	if err := ar.SELECT(ar.ID).
		FROM(ar).
		WHERE(ar.GuardName.EQ(jet.String(helpers.Guard(name)))).
		LIMIT(1).
		QueryContext(context.TODO(), p.db, &roleIDs); err != nil {
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

	_, err := arp.INSERT(
		arp.PermissionID,
		arp.RoleID,
	).
		MODELS(rolePerms).
		ExecContext(context.TODO(), p.db)

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number != 1062 {
		return err
	}

	return nil
}
