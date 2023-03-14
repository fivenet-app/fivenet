package perms

import (
	"context"

	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
	"github.com/galexrt/arpanet/query"
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

	_, err := stmt.ExecContext(context.TODO(), query.DB)
	return err
}

func (p *Perms) GetAllPermissions() (collections.Permissions, error) {
	stmt := ap.SELECT(
		ap.AllColumns,
	).FROM(ap)

	var dest collections.Permissions
	err := stmt.QueryContext(context.TODO(), query.DB, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}
