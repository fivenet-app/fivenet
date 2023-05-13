package perms

import (
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/perms/helpers"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var tPerms = table.FivenetPermissions

func (p *Perms) CreatePermission(category Category, name Name) (uint64, error) {
	guard := BuildGuard(category, name)
	stmt := tPerms.
		INSERT(
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
		).
		VALUES(
			category,
			name,
			guard,
		)

	res, err := stmt.ExecContext(p.ctx, p.db)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, err
		}

		perm, err := p.getPermissionByGuard(guard)
		if err != nil {
			return 0, err
		}

		return perm.ID, nil
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastId), nil
}

func (p *Perms) UpdatePermission(id uint64, category Category, name Name) error {
	guard := helpers.Guard(fmt.Sprintf("%s-%s", category, name))
	stmt := tPerms.
		UPDATE(
			tPerms.Name,
			tPerms.Category,
			tPerms.GuardName,
		).
		SET(
			tPerms.Category.SET(jet.String(string(category))),
			tPerms.Name.SET(jet.String(string(name))),
			tPerms.GuardName.SET(jet.String(guard)),
		).
		WHERE(
			tPerms.ID.EQ(jet.Uint64(id)),
		)

	_, err := stmt.ExecContext(p.ctx, p.db)
	if err != nil {
		return err
	}

	return nil
}

func (p *Perms) GetAllPermissions() ([]*permissions.Permission, error) {
	tPerms := tPerms.AS("permission")

	stmt := tPerms.
		SELECT(
			tPerms.AllColumns,
		).
		FROM(tPerms)

	var dest []*permissions.Permission
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (p *Perms) RemovePermissionsByIDs(ids ...uint64) error {
	wIds := make([]jet.Expression, len(ids))
	for i := 0; i < len(ids); i++ {
		wIds[i] = jet.Uint64(ids[i])
	}

	stmt := tPerms.
		DELETE().
		WHERE(
			tPerms.ID.IN(wIds...),
		)

	_, err := stmt.ExecContext(p.ctx, p.db)
	if err != nil {
		return err
	}

	return nil
}

func (p *Perms) GetPermissionsByIDs(ids ...uint64) ([]*permissions.Permission, error) {
	wIds := make([]jet.Expression, len(ids))
	for i := 0; i < len(ids); i++ {
		wIds[i] = jet.Uint64(ids[i])
	}

	tPerms := tPerms.AS("permission")

	stmt := tPerms.
		SELECT(
			tPerms.AllColumns,
		).
		FROM(tPerms).
		WHERE(
			tPerms.ID.IN(wIds...),
		)

	var dest []*permissions.Permission
	err := stmt.QueryContext(p.ctx, p.db, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (p *Perms) getPermissionByGuard(name string) (*model.FivenetPermissions, error) {
	guard := helpers.Guard(name)

	stmt := tPerms.
		SELECT(
			tPerms.AllColumns,
		).
		FROM(tPerms).
		WHERE(
			tPerms.GuardName.EQ(jet.String(guard)),
		).
		LIMIT(1)

	var dest model.FivenetPermissions
	err := stmt.QueryContext(p.ctx, p.db, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}
