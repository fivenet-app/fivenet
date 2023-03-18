package perms

import (
	"github.com/galexrt/arpanet/pkg/dbutils"
	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
)

func (p *Perms) CreatePermission(name string, description string) error {
	stmt := ap.INSERT(
		ap.Name,
		ap.GuardName,
		ap.Description,
	).
		VALUES(
			name,
			helpers.Guard(name),
			description,
		)

	_, err := stmt.ExecContext(p.ctx, p.db)

	if !dbutils.IsDuplicateError(err) {
		return err
	}

	return nil
}

func (p *Perms) GetAllPermissions() (collections.Permissions, error) {
	stmt := ap.SELECT(
		ap.AllColumns,
	).FROM(ap)

	var dest collections.Permissions
	err := stmt.QueryContext(p.ctx, p.db, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}
