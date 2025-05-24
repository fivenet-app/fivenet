package perms

import (
	"context"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/pkg/errors"
	"github.com/puzpuzpuz/xsync/v4"
)

func (p *Perms) loadData(ctx context.Context) error {
	ctx, span := p.tracer.Start(ctx, "perms-load")
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

	if err := p.loadJobAttrs(ctx, ""); err != nil {
		return fmt.Errorf("failed to load job attributes. %w", err)
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

func (p *Perms) loadAttributes(ctx context.Context) error {
	stmt := tAttrs.
		SELECT(
			tAttrs.ID.AS("id"),
			tAttrs.PermissionID.AS("permission_id"),
			tAttrs.Key.AS("key"),
			tAttrs.Type.AS("type"),
			tAttrs.ValidValues.AS("valid_values"),
		).
		FROM(tAttrs)

	var dest []struct {
		ID           uint64
		PermissionID uint64
		Key          Key
		Type         permissions.AttributeTypes
		ValidValues  *permissions.AttributeValues
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query attributes. %w", err)
		}
	}

	for _, attr := range dest {
		if err := p.addOrUpdateAttributeInMap(attr.PermissionID, attr.ID, attr.Key, attr.Type, attr.ValidValues); err != nil {
			return fmt.Errorf("failed to add/update attribute in map. %w", err)
		}
	}

	return nil
}

func (p *Perms) loadRoles(ctx context.Context, id uint64) error {
	stmt := tRoles.
		SELECT(
			tRoles.ID.AS("id"),
			tRoles.Job.AS("job"),
			tRoles.Grade.AS("grade"),
		).
		FROM(tRoles)

	if id != 0 {
		stmt = stmt.
			WHERE(tRoles.ID.EQ(jet.Uint64(id)))
	}

	var dest []struct {
		ID    uint64
		Job   string
		Grade int32
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query roles. %w", err)
		}
	}

	for _, role := range dest {
		grades, _ := p.permsJobsRoleMap.LoadOrCompute(role.Job, func() (*xsync.Map[int32, uint64], bool) {
			return xsync.NewMap[int32, uint64](), false
		})
		grades.Store(role.Grade, role.ID)

		p.roleIDToJobMap.Store(role.ID, role.Job)
	}

	return nil
}

func (p *Perms) loadRolePermissions(ctx context.Context, roleId uint64) error {
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
			tRoles.ID.EQ(jet.Uint64(roleId)),
		)
	}

	var dest []struct {
		RoleID uint64
		ID     uint64
		Val    bool
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query role permissions. %w", err)
		}
	}

	found := map[uint64][]uint64{}
	for _, rolePerms := range dest {
		perms, _ := p.permsRoleMap.LoadOrCompute(rolePerms.RoleID, func() (*xsync.Map[uint64, bool], bool) {
			return xsync.NewMap[uint64, bool](), false
		})
		perms.Store(rolePerms.ID, rolePerms.Val)

		if _, ok := found[rolePerms.RoleID]; !ok {
			found[rolePerms.RoleID] = []uint64{}
		}

		found[rolePerms.RoleID] = append(found[rolePerms.RoleID], rolePerms.ID)
	}

	// Check if any role perms don't exist anymore in the db and need to be deleted
	for roleId, list := range found {
		perms, ok := p.permsRoleMap.Load(roleId)
		if !ok {
			continue
		}

		perms.Range(func(permId uint64, _ bool) bool {
			if !slices.Contains(list, permId) {
				perms.Delete(permId)
			}
			return true
		})
	}

	return nil
}

func (p *Perms) loadJobAttrs(ctx context.Context, job string) error {
	stmt := tJobAttrs.
		SELECT(
			tJobAttrs.Job.AS("job"),
			tJobAttrs.AttrID.AS("attr_id"),
			tJobAttrs.MaxValues.AS("max_values"),
		).
		FROM(tJobAttrs).
		ORDER_BY(
			tJobAttrs.Job.ASC(),
		)

	if job != "" {
		stmt = stmt.WHERE(
			tJobAttrs.Job.EQ(jet.String(job)),
		)
	}

	var dest []struct {
		Job       string
		AttrID    uint64
		MaxValues *permissions.AttributeValues
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query job attributes. %w", err)
		}
	}

	// No attributes? Delete cached data
	if len(dest) == 0 {
		if job != "" {
			p.attrsJobMaxValuesMap.Delete(job)
		} else {
			p.attrsJobMaxValuesMap.Clear()
		}
	} else {
		found := map[string][]uint64{}
		for _, jobAttrs := range dest {
			attrs, _ := p.attrsJobMaxValuesMap.LoadOrCompute(jobAttrs.Job, func() (*xsync.Map[uint64, *permissions.AttributeValues], bool) {
				return xsync.NewMap[uint64, *permissions.AttributeValues](), false
			})
			attrs.Store(jobAttrs.AttrID, jobAttrs.MaxValues)

			if _, ok := found[jobAttrs.Job]; !ok {
				found[jobAttrs.Job] = []uint64{}
			}
			found[jobAttrs.Job] = append(found[jobAttrs.Job], jobAttrs.AttrID)
		}

		// Check if any job attrs don't exist anymore in the db and need to be deleted
		for job, list := range found {
			attrs, ok := p.attrsJobMaxValuesMap.Load(job)
			if !ok {
				continue
			}

			attrs.Range(func(attrId uint64, _ *permissions.AttributeValues) bool {
				if !slices.Contains(list, attrId) {
					attrs.Delete(attrId)
				}
				return true
			})
		}
	}

	return nil
}

func (p *Perms) loadRoleAttributes(ctx context.Context, roleId uint64) error {
	stmt := tRoleAttrs.
		SELECT(
			tRoleAttrs.AttrID.AS("attr_id"),
			tAttrs.PermissionID.AS("permission_id"),
			tRoleAttrs.RoleID.AS("role_id"),
			tAttrs.Key.AS("key"),
			tAttrs.Type.AS("type"),
			tRoleAttrs.Value.AS("value"),
		).
		FROM(
			tRoleAttrs.
				INNER_JOIN(tAttrs,
					tAttrs.ID.EQ(tRoleAttrs.AttrID),
				),
		)

	if roleId != 0 {
		stmt = stmt.WHERE(
			tRoleAttrs.RoleID.EQ(jet.Uint64(roleId)),
		)
	}

	var dest []struct {
		AttrID       uint64
		PermissionID uint64
		RoleID       uint64
		Key          Key
		Type         permissions.AttributeTypes
		Value        *permissions.AttributeValues
	}

	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query role attributes. %w", err)
		}
	}

	found := map[uint64][]uint64{}
	for _, ra := range dest {
		ra.Value.Default(ra.Type)
		p.updateRoleAttributeInMap(ra.RoleID, ra.PermissionID, ra.AttrID, ra.Key, ra.Type, ra.Value)

		if _, ok := found[ra.RoleID]; !ok {
			found[ra.RoleID] = []uint64{}
		}
		found[ra.RoleID] = append(found[ra.RoleID], ra.AttrID)
	}

	// Check if any role attrs that don't exist anymore in the db and need to be deleted
	for roleId, list := range found {
		attrRoleMap, ok := p.attrsRoleMap.Load(roleId)
		if !ok {
			continue
		}

		attrRoleMap.Range(func(attrId uint64, _ *cacheRoleAttr) bool {
			if !slices.Contains(list, attrId) {
				attrRoleMap.Delete(attrId)
			}
			return true
		})
	}

	return nil
}
