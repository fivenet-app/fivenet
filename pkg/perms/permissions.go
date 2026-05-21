package perms

import (
	"context"
	"fmt"

	permissionspermissions "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/permissions"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/pkg/errors"
)

var tPerms = table.FivenetRbacPermissions

func (ps *Perms) createPermission(
	ctx context.Context,
	namespace Namespace,
	service Service,
	name Name,
	order int32,
	icon *string,
) (int64, error) {
	guard := BuildGuard(namespace, service, name)
	stmt := tPerms.
		INSERT(
			tPerms.Namespace,
			tPerms.Service,
			tPerms.Name,
			tPerms.GuardName,
			tPerms.Order,
			tPerms.Icon,
		).
		VALUES(
			namespace,
			service,
			name,
			guard,
			order,
			icon,
		)

	res, err := stmt.ExecContext(ctx, ps.db)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, fmt.Errorf("failed to execute insert statement. %w", err)
		}

		permId, ok := ps.lookupPermIDByGuard(guard)
		if !ok {
			permId, err = ps.loadPermissionByGuard(ctx, guard) // try to load it from DB again
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

func (ps *Perms) loadPermissionFromDatabaseByCategoryName(
	ctx context.Context,
	namespace Namespace,
	service Service,
	name Name,
) (*permissionspermissions.Permission, error) {
	tPerms := tPerms.AS("permission")

	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.CreatedAt,
			tPerms.Namespace,
			tPerms.Service,
			tPerms.Name,
			tPerms.GuardName,
			tPerms.Order,
			tPerms.Icon,
		).
		FROM(tPerms).
		WHERE(mysql.AND(
			tPerms.Namespace.EQ(mysql.String(string(namespace))),
			tPerms.Service.EQ(mysql.String(string(service))),
			tPerms.Name.EQ(mysql.String(string(name))),
		)).
		LIMIT(1)

	dest := &permissionspermissions.Permission{}
	if err := stmt.QueryContext(ctx, ps.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query permission by guard. %w", err)
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (ps *Perms) updatePermission(
	ctx context.Context,
	id int64,
	namespace Namespace,
	service Service,
	name Name,
	order int32,
	icon *string,
) error {
	guard := guard(string(namespace) + "." + string(service) + "-" + string(name))

	var iconExp mysql.StringExpression
	if icon != nil {
		iconExp = mysql.String(*icon)
	} else {
		iconExp = mysql.StringExp(mysql.NULL)
	}

	stmt := tPerms.
		UPDATE(
			tPerms.Name,
			tPerms.Namespace,
			tPerms.Service,
			tPerms.GuardName,
			tPerms.Order,
			tPerms.Icon,
		).
		SET(
			tPerms.Namespace.SET(mysql.String(string(namespace))),
			tPerms.Service.SET(mysql.String(string(service))),
			tPerms.Name.SET(mysql.String(string(name))),
			tPerms.GuardName.SET(mysql.String(guard)),
			tPerms.Order.SET(mysql.Int32(order)),
			tPerms.Icon.SET(iconExp),
		).
		WHERE(
			tPerms.ID.EQ(mysql.Int64(id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, ps.db); err != nil {
		return fmt.Errorf("failed to execute update statement. %w", err)
	}

	return nil
}

func (ps *Perms) GetAllPermissions(
	ctx context.Context,
) ([]*permissionspermissions.Permission, error) {
	tPerms := tPerms.AS("permission")

	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.CreatedAt,
			tPerms.Namespace,
			tPerms.Service,
			tPerms.Name,
			tPerms.GuardName,
			tPerms.Order,
			tPerms.Icon,
		).
		FROM(tPerms).
		ORDER_BY(
			tPerms.GuardName.ASC(),
		)

	var dest []*permissionspermissions.Permission
	if err := stmt.QueryContext(ctx, ps.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query all permissions. %w", err)
		}
	}

	return dest, nil
}

func (ps *Perms) RemovePermissionsByIDs(ctx context.Context, ids ...int64) error {
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

	if _, err := stmt.ExecContext(ctx, ps.db); err != nil {
		return fmt.Errorf("failed to execute delete statement. %w", err)
	}

	return nil
}

func (ps *Perms) GetPermissionsByIDs(
	ctx context.Context,
	ids ...int64,
) ([]*permissionspermissions.Permission, error) {
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
			tPerms.Namespace,
			tPerms.Service,
			tPerms.Name,
			tPerms.GuardName,
			tPerms.Order,
			tPerms.Icon,
		).
		FROM(tPerms).
		WHERE(
			tPerms.ID.IN(wIds...),
		).
		ORDER_BY(
			tPerms.GuardName.ASC(),
		)

	var dest []*permissionspermissions.Permission
	if err := stmt.QueryContext(ctx, ps.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query permissions by IDs. %w", err)
		}
	}

	return dest, nil
}

func (ps *Perms) GetPermission(
	ctx context.Context,
	namespace Namespace,
	service Service,
	name Name,
) (*permissionspermissions.Permission, error) {
	tPerms := tPerms.AS("permission")

	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.CreatedAt,
			tPerms.Namespace,
			tPerms.Service,
			tPerms.Name,
			tPerms.GuardName,
			tPerms.Order,
			tPerms.Icon,
		).
		FROM(tPerms).
		WHERE(mysql.AND(
			tPerms.Namespace.EQ(mysql.String(string(namespace))),
			tPerms.Service.EQ(mysql.String(string(service))),
			tPerms.Name.EQ(mysql.String(string(name))),
		)).
		LIMIT(1)

	var dest permissionspermissions.Permission
	if err := stmt.QueryContext(ctx, ps.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query permission by category and name. %w", err)
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return &dest, nil
}
