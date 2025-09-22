package perms

import (
	"context"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/pkg/errors"
	"github.com/puzpuzpuz/xsync/v4"
)

func (p *Perms) loadData(ctx context.Context) error {
	ctx, span := p.tracer.Start(ctx, "perms.load")
	defer span.End()

	if err := p.loadPermissions(ctx); err != nil {
		return fmt.Errorf("failed to load permissions. %w", err)
	}

	if err := p.loadAttributes(ctx); err != nil {
		return fmt.Errorf("failed to load attributes. %w", err)
	}

	if err := p.loadRoles(ctx, 0); err != nil {
		return fmt.Errorf("failed to load roles. %w", err)
	}

	if err := p.loadRolePermissions(ctx, 0); err != nil {
		return fmt.Errorf("failed to load role permissions. %w", err)
	}

	if err := p.loadRoleAttributes(ctx, 0); err != nil {
		return fmt.Errorf("failed to load role attributes. %w", err)
	}

	return nil
}

func (p *Perms) loadPermissions(ctx context.Context) error {
	tPerms := tPerms.AS("cache_perm")
	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
			tPerms.Order,
		).
		FROM(tPerms)

	var dest []*cachePerm
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query permissions. %w", err)
		}
	}

	for _, perm := range dest {
		p.permsMap.Store(perm.ID, &cachePerm{
			ID:        perm.ID,
			Category:  perm.Category,
			Name:      perm.Name,
			GuardName: BuildGuard(perm.Category, perm.Name),
			Order:     perm.Order,
		})
		p.permsGuardToIDMap.Store(BuildGuard(perm.Category, perm.Name), perm.ID)
	}

	return nil
}

func (p *Perms) loadPermissionByGuard(ctx context.Context, guard string) (int64, error) {
	tPerms := tPerms.AS("cache_perm")
	stmt := tPerms.
		SELECT(
			tPerms.ID,
			tPerms.Category,
			tPerms.Name,
			tPerms.GuardName,
			tPerms.Order,
		).
		FROM(tPerms)

	if guard != "" {
		stmt = stmt.
			WHERE(tPerms.GuardName.EQ(mysql.String(guard))).
			LIMIT(1)
	}

	perm := cachePerm{}
	if err := stmt.QueryContext(ctx, p.db, &perm); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, fmt.Errorf("failed to query permissions. %w", err)
		}
	}

	if perm.ID != 0 {
		p.permsMap.Store(perm.ID, &cachePerm{
			ID:        perm.ID,
			Category:  perm.Category,
			Name:      perm.Name,
			GuardName: BuildGuard(perm.Category, perm.Name),
			Order:     perm.Order,
		})
		p.permsGuardToIDMap.Store(BuildGuard(perm.Category, perm.Name), perm.ID)
	}

	return perm.ID, nil
}

func (p *Perms) loadAttributes(ctx context.Context) error {
	tAttrs := table.FivenetRbacAttrs.AS("role_attribute")
	stmt := tAttrs.
		SELECT(
			tAttrs.ID.AS("role_attribute.attr_id"),
			tAttrs.PermissionID,
			tAttrs.Key,
			tAttrs.Type,
			tAttrs.ValidValues,
		).
		FROM(tAttrs)

	var dest []*permissions.RoleAttribute
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query attributes. %w", err)
		}
	}

	for _, attr := range dest {
		if err := p.addOrUpdateAttributeInMap(attr.GetPermissionId(), attr.GetAttrId(), Key(attr.GetKey()), permissions.AttributeTypes(attr.GetType()), attr.GetValidValues()); err != nil {
			return fmt.Errorf("failed to add/update attribute in map. %w", err)
		}
	}

	return nil
}

func (p *Perms) loadRoles(ctx context.Context, id int64) error {
	stmt := tRoles.
		SELECT(
			tRoles.ID,
			tRoles.Job,
			tRoles.Grade,
		).
		FROM(tRoles)

	if id != 0 {
		stmt = stmt.
			WHERE(tRoles.ID.EQ(mysql.Int64(id)))
	}

	var dest []*permissions.Role
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query roles. %w", err)
		}
	}

	for _, role := range dest {
		grades, _ := p.permsJobsRoleMap.LoadOrCompute(
			role.GetJob(),
			func() (*xsync.Map[int32, int64], bool) {
				return xsync.NewMap[int32, int64](), false
			},
		)
		grades.Store(role.GetGrade(), role.GetId())

		p.roleIDToJobMap.Store(role.GetId(), role.GetJob())
	}

	return nil
}

func (p *Perms) loadRolePermissions(ctx context.Context, roleId int64) error {
	stmt := tRolePerms.
		SELECT(
			tRolePerms.RoleID.AS("role_id"),
			tRolePerms.PermissionID.AS("id"),
			tRolePerms.Val.AS("val"),
		).
		FROM(tRolePerms.
			INNER_JOIN(tRoles,
				tRoles.ID.EQ(tRolePerms.RoleID),
			),
		).
		ORDER_BY(
			tRoles.Job.ASC(),
			tRoles.Grade.DESC(),
		)

	if roleId != 0 {
		stmt = stmt.WHERE(
			tRoles.ID.EQ(mysql.Int64(roleId)),
		)
	}

	var dest []struct {
		RoleID int64
		ID     int64
		Val    bool
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query role permissions. %w", err)
		}
	}

	found := map[int64][]int64{}
	for _, rolePerms := range dest {
		perms, _ := p.permsRoleMap.LoadOrCompute(
			rolePerms.RoleID,
			func() (*xsync.Map[int64, bool], bool) {
				return xsync.NewMap[int64, bool](), false
			},
		)
		perms.Store(rolePerms.ID, rolePerms.Val)

		if _, ok := found[rolePerms.RoleID]; !ok {
			found[rolePerms.RoleID] = []int64{}
		}

		found[rolePerms.RoleID] = append(found[rolePerms.RoleID], rolePerms.ID)
	}

	// Check if any role perms don't exist anymore in the db and need to be deleted
	for rId, list := range found {
		perms, ok := p.permsRoleMap.Load(rId)
		if !ok {
			continue
		}

		perms.Range(func(permId int64, _ bool) bool {
			if !slices.Contains(list, permId) {
				perms.Delete(permId)
			}
			return true
		})
	}

	return nil
}

func (p *Perms) loadRoleAttributes(ctx context.Context, roleId int64) error {
	tRoleAttrs := table.FivenetRbacRolesAttrs.AS("role_attribute")
	stmt := tRoleAttrs.
		SELECT(
			tRoleAttrs.AttrID,
			tAttrs.PermissionID.AS("role_attribute.permission_id"),
			tRoleAttrs.RoleID,
			tAttrs.Key.AS("role_attribute.key"),
			tAttrs.Type.AS("role_attribute.type"),
			tRoleAttrs.Value.AS("role_attribute.value"),
		).
		FROM(
			tRoleAttrs.
				INNER_JOIN(tAttrs,
					tAttrs.ID.EQ(tRoleAttrs.AttrID),
				),
		)

	if roleId != 0 {
		stmt = stmt.WHERE(
			tRoleAttrs.RoleID.EQ(mysql.Int64(roleId)),
		)
	}

	var dest []*permissions.RoleAttribute
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query role attributes. %w", err)
		}
	}

	found := map[int64][]int64{}
	for _, ra := range dest {
		ra.GetValue().Default(permissions.AttributeTypes(ra.GetType()))
		p.updateRoleAttributeInMap(
			ra.GetRoleId(),
			ra.GetPermissionId(),
			ra.GetAttrId(),
			Key(ra.GetKey()),
			permissions.AttributeTypes(ra.GetType()),
			ra.GetValue(),
		)

		if _, ok := found[ra.GetRoleId()]; !ok {
			found[ra.GetRoleId()] = []int64{}
		}
		found[ra.GetRoleId()] = append(found[ra.GetRoleId()], ra.GetAttrId())
	}

	// Check if any role attrs that don't exist anymore in the db and need to be deleted
	for rId, list := range found {
		attrRoleMap, ok := p.attrsRoleMap.Load(rId)
		if !ok {
			continue
		}

		attrRoleMap.Range(func(attrId int64, _ *cacheRoleAttr) bool {
			if !slices.Contains(list, attrId) {
				attrRoleMap.Delete(attrId)
			}
			return true
		})
	}

	return nil
}

func (p *Perms) loadJobRoles(ctx context.Context, job string) error {
	if job == "" {
		if err := p.loadRolePermissions(ctx, 0); err != nil {
			return fmt.Errorf("failed to load roles permissions for job %s. %w", job, err)
		}

		if err := p.loadRoleAttributes(ctx, 0); err != nil {
			return fmt.Errorf("failed to load role attributes for job %s. %w", job, err)
		}

		return nil
	}

	roles, err := p.GetJobRoles(ctx, job)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if err := p.loadRolePermissions(ctx, role.GetId()); err != nil {
			return fmt.Errorf(
				"failed to load role permissions for job %s, role %d. %w",
				job,
				role.GetId(),
				err,
			)
		}

		if err := p.loadRoleAttributes(ctx, role.GetId()); err != nil {
			return fmt.Errorf(
				"failed to load role attributes for job %s, role %d. %w",
				job,
				role.GetId(),
				err,
			)
		}
	}

	return nil
}
