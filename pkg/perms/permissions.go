package perms

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/pkg/errors"
)

var tPerms = table.FivenetRbacPermissions

func (p *Perms) CreatePermission(
	ctx context.Context,
	category Category,
	name Name,
) (int64, error) {
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

	res, err := stmt.ExecContext(ctx, p.db)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, fmt.Errorf("failed to execute insert statement. %w", err)
		}

		permId, ok := p.lookupPermIDByGuard(guard)
		if !ok {
			permId, err = p.loadPermissionByGuard(ctx, guard) // try to load it from DB again
			if err != nil || permId == 0 {
				return 0, fmt.Errorf("failed to query permission after duplicate error. %w", err)
			}

			return permId, nil
		}

		return permId, nil
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID. %w", err)
	}

	return lastId, nil
}

func (p *Perms) loadPermissionFromDatabaseByCategoryName(
	ctx context.Context,
	category Category,
	name Name,
) (*permissions.Permission, error) {
	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.CreatedAt,
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
		).
		FROM(tPerms).
		WHERE(mysql.AND(
			tPerms.Category.EQ(mysql.String(string(category))),
			tPerms.Name.EQ(mysql.String(string(name))),
		)).
		LIMIT(1)

	dest := &permissions.Permission{}

	if err := stmt.QueryContext(ctx, p.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query permission by guard. %w", err)
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (p *Perms) UpdatePermission(
	ctx context.Context,
	id int64,
	category Category,
	name Name,
) error {
	guard := Guard(string(category) + "-" + string(name))

	stmt := tPerms.
		UPDATE(
			tPerms.Name,
			tPerms.Category,
			tPerms.GuardName,
		).
		SET(
			tPerms.Category.SET(mysql.String(string(category))),
			tPerms.Name.SET(mysql.String(string(name))),
			tPerms.GuardName.SET(mysql.String(guard)),
		).
		WHERE(
			tPerms.ID.EQ(mysql.Int64(id)),
		)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("failed to execute update statement. %w", err)
	}

	return nil
}

func (p *Perms) GetAllPermissions(ctx context.Context) ([]*permissions.Permission, error) {
	tPerms := tPerms.AS("permission")

	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.CreatedAt,
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
		).
		FROM(tPerms).
		ORDER_BY(
			tPerms.GuardName.ASC(),
		)

	var dest []*permissions.Permission
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query all permissions. %w", err)
		}
	}

	return dest, nil
}

func (p *Perms) RemovePermissionsByIDs(ctx context.Context, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}

	wIds := make([]mysql.Expression, len(ids))
	for i := range ids {
		wIds[i] = mysql.Int64(ids[i])
	}

	stmt := tPerms.
		DELETE().
		WHERE(
			tPerms.ID.IN(wIds...),
		).
		LIMIT(int64(len(wIds)))

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("failed to execute delete statement. %w", err)
	}

	return nil
}

func (p *Perms) GetPermissionsByIDs(
	ctx context.Context,
	ids ...int64,
) ([]*permissions.Permission, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	wIds := make([]mysql.Expression, len(ids))
	for i := range ids {
		wIds[i] = mysql.Int64(ids[i])
	}

	tPerms := tPerms.AS("permission")

	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.CreatedAt,
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
		).
		FROM(tPerms).
		WHERE(
			tPerms.ID.IN(wIds...),
		).
		ORDER_BY(
			tPerms.GuardName.ASC(),
		)

	var dest []*permissions.Permission
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query permissions by IDs. %w", err)
		}
	}

	return dest, nil
}

func (p *Perms) GetPermission(
	ctx context.Context,
	category Category,
	name Name,
) (*permissions.Permission, error) {
	tPerms := tPerms.AS("permission")

	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.CreatedAt,
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
		).
		FROM(tPerms).
		WHERE(mysql.AND(
			tPerms.Category.EQ(mysql.String(string(category))),
			tPerms.Name.EQ(mysql.String(string(name))),
		)).
		LIMIT(1)

	var dest permissions.Permission
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query permission by category and name. %w", err)
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return &dest, nil
}
