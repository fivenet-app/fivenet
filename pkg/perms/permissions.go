package perms

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/pkg/errors"
)

var tPerms = table.FivenetRbacPermissions

func (p *Perms) CreatePermission(ctx context.Context, category Category, name Name) (uint64, error) {
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
			return 0, fmt.Errorf("created permission not found in our cache")
		}

		return permId, nil
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID. %w", err)
	}

	return uint64(lastId), nil
}

func (p *Perms) loadPermissionFromDatabaseByCategoryName(ctx context.Context, category Category, name Name) (*model.FivenetRbacPermissions, error) {
	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.CreatedAt,
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
		).
		FROM(tPerms).
		WHERE(jet.AND(
			tPerms.Category.EQ(jet.String(string(category))),
			tPerms.Name.EQ(jet.String(string(name))),
		)).
		LIMIT(1)

	dest := &model.FivenetRbacPermissions{}

	if err := stmt.QueryContext(ctx, p.db, dest); err != nil {
		return nil, fmt.Errorf("failed to query permission by guard. %w", err)
	}

	return dest, nil
}

func (p *Perms) UpdatePermission(ctx context.Context, id uint64, category Category, name Name) error {
	guard := Guard(string(category) + "-" + string(name))

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

func (p *Perms) RemovePermissionsByIDs(ctx context.Context, ids ...uint64) error {
	if len(ids) == 0 {
		return nil
	}

	wIds := make([]jet.Expression, len(ids))
	for i := range ids {
		wIds[i] = jet.Uint64(ids[i])
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

func (p *Perms) GetPermissionsByIDs(ctx context.Context, ids ...uint64) ([]*permissions.Permission, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	wIds := make([]jet.Expression, len(ids))
	for i := range ids {
		wIds[i] = jet.Uint64(ids[i])
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
		return nil, fmt.Errorf("failed to query permissions by IDs. %w", err)
	}

	return dest, nil
}

func (p *Perms) GetPermission(ctx context.Context, category Category, name Name) (*permissions.Permission, error) {
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
		WHERE(jet.AND(
			tPerms.Category.EQ(jet.String(string(category))),
			tPerms.Name.EQ(jet.String(string(name))),
		)).
		LIMIT(1)

	var dest permissions.Permission
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to query permission by category and name. %w", err)
	}

	return &dest, nil
}
