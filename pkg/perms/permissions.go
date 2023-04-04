package perms

import (
	"github.com/galexrt/fivenet/pkg/dbutils"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/pkg/perms/helpers"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (p *Perms) CreatePermission(name string, description string) error {
	stmt := ap.
		INSERT(
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

func (p *Perms) UpdatePermission(id uint64, name string, description string) error {
	stmt := ap.
		UPDATE(
			ap.Name,
			ap.GuardName,
			ap.Description,
		).
		SET(
			ap.Name.SET(jet.String(name)),
			ap.GuardName.SET(jet.String(helpers.Guard(name))),
			ap.Description.SET(jet.String(description)),
		).
		WHERE(
			ap.ID.EQ(jet.Uint64(id)),
		)

	_, err := stmt.ExecContext(p.ctx, p.db)
	return err
}

func (p *Perms) GetAllPermissions() (collections.Permissions, error) {
	stmt := ap.
		SELECT(
			ap.AllColumns,
		).
		FROM(ap)

	var dest collections.Permissions
	err := stmt.QueryContext(p.ctx, p.db, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (p *Perms) GetPermissionsByIDs(ids ...uint64) (collections.Permissions, error) {
	wIds := make([]jet.Expression, len(ids))
	for i := 0; i < len(ids); i++ {
		wIds[i] = jet.Uint64(ids[i])
	}

	stmt := ap.
		SELECT(
			ap.AllColumns,
		).
		FROM(ap).
		WHERE(
			ap.ID.IN(wIds...),
		)

	var dest collections.Permissions
	err := stmt.QueryContext(p.ctx, p.db, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}
