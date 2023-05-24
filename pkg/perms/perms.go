package perms

import (
	"context"
	"database/sql"
	"errors"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Permissions interface {
	GetAllPermissions(ctx context.Context) ([]*permissions.Permission, error)
	GetPermissionsByIDs(ctx context.Context, ids ...uint64) ([]*permissions.Permission, error)
	CreatePermission(ctx context.Context, category Category, name Name) (uint64, error)
	GetPermissionsOfUser(userInfo *userinfo.UserInfo) (collections.Permissions, error)

	GetJobRoles(ctx context.Context, job string) (collections.Roles, error)
	GetJobRolesUpTo(ctx context.Context, job string, grade int32) (collections.Roles, error)
	GetClosestJobRole(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error)
	CountRolesForJob(ctx context.Context, prefix string) (int64, error)

	GetRole(ctx context.Context, id uint64) (*model.FivenetRoles, error)
	GetRoleByJobAndGrade(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error)
	GetRolePermissions(ctx context.Context, id uint64) ([]*permissions.Permission, error)

	CreateRole(ctx context.Context, job string, grade int32) (*model.FivenetRoles, error)
	DeleteRole(ctx context.Context, id uint64) error
	UpdateRolePermissions(ctx context.Context, id uint64, perms ...AddPerm) error
	RemovePermissionsFromRole(ctx context.Context, id uint64, perms ...uint64) error

	Can(userInfo *userinfo.UserInfo, category Category, name Name) bool

	GetAttribute(category Category, name Name, key Key) (*permissions.RoleAttribute, error)
	GetAttributeByIDs(ctx context.Context, ids ...uint64) ([]*permissions.RoleAttribute, error)
	CreateAttribute(ctx context.Context, permId uint64, key Key, aType AttributeTypes, validValues any) (uint64, error)
	UpdateAttribute(ctx context.Context, attributeId uint64, permId uint64, key Key, aType AttributeTypes, validValues any) error
	GetRoleAttributes(job string, grade int32) ([]*permissions.RoleAttribute, error)
	FlattenRoleAttributes(job string, grade int32) ([]string, error)
	GetAllAttributes(ctx context.Context, job string) ([]*permissions.RoleAttribute, error)
	AddOrUpdateAttributesToRole(ctx context.Context, attrs ...*permissions.RoleAttribute) error
	UpdateRoleAttributes(ctx context.Context, attrs ...*permissions.RoleAttribute) error
	RemoveAttributesFromRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error

	Attr(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (any, error)
}

type userCacheKey struct {
	userId int32
	permId uint64
}

type Perms struct {
	db *sql.DB

	tracer trace.Tracer
	ctx    context.Context

	permsMap syncx.Map[uint64, *cachePerm]
	// Guard name to permission ID
	permsGuardToIDMap syncx.Map[string, uint64]
	// Job name to map of grade numbers to role ID
	permsJobsRoleMap syncx.Map[string, syncx.Map[int32, uint64]]
	// Role ID to map of permissions ID and result
	permsRoleMap syncx.Map[uint64, syncx.Map[uint64, bool]]

	// Attribute map (key: ID of attribute)
	attrsMap syncx.Map[uint64, *cacheAttr]
	// Role ID to map of role attributes
	attrsRoleMap syncx.Map[uint64, syncx.Map[uint64, *cacheRoleAttr]]
	// Perm ID to map Key -> cached attribute
	attrsPermsMap syncx.Map[uint64, syncx.Map[Key, uint64]]

	userCanCacheTTL time.Duration
	userCanCache    *cache.Cache[userCacheKey, bool]
}

func New(ctx context.Context, tp *tracesdk.TracerProvider, db *sql.DB) *Perms {
	userCanCache := cache.NewContext(
		ctx,
		cache.AsLRU[userCacheKey, bool](lru.WithCapacity(128)),
		cache.WithJanitorInterval[userCacheKey, bool](15*time.Second),
	)

	p := &Perms{
		db: db,

		tracer: tp.Tracer("perms"),
		ctx:    ctx,

		permsMap:          syncx.Map[uint64, *cachePerm]{},
		permsGuardToIDMap: syncx.Map[string, uint64]{},
		permsJobsRoleMap:  syncx.Map[string, syncx.Map[int32, uint64]]{},
		permsRoleMap:      syncx.Map[uint64, syncx.Map[uint64, bool]]{},

		attrsMap:      syncx.Map[uint64, *cacheAttr]{},
		attrsRoleMap:  syncx.Map[uint64, syncx.Map[uint64, *cacheRoleAttr]]{},
		attrsPermsMap: syncx.Map[uint64, syncx.Map[Key, uint64]]{},

		userCanCacheTTL: 30 * time.Second,
		userCanCache:    userCanCache,
	}

	p.load()

	return p
}

type cachePerm struct {
	ID        uint64
	Category  Category
	Name      Name
	GuardName string
}

type cacheAttr struct {
	ID           uint64
	PermissionID uint64
	Category     Category
	Name         Name
	Key          Key
	Type         AttributeTypes
	ValidValues  *permissions.AttributeValues
}

type cacheRoleAttr struct {
	AttrID uint64
	Key    Key
	Type   AttributeTypes
	Value  *permissions.AttributeValues
}

func (p *Perms) load() error {
	ctx, span := p.tracer.Start(p.ctx, "perms-load")
	defer span.End()

	if err := p.loadRoleIDs(ctx); err != nil {
		return err
	}

	if err := p.loadRolePermissions(ctx); err != nil {
		return err
	}

	if err := p.loadRoleAttributes(ctx); err != nil {
		return err
	}

	return nil
}

func (p *Perms) loadRoleIDs(ctx context.Context) error {
	stmt := tRoles.
		SELECT(
			tRoles.ID.AS("id"),
			tRoles.Job.AS("job"),
			tRoles.Grade.AS("grade"),
		).
		FROM(tRoles)

	var dest []struct {
		ID    uint64
		Job   string
		Grade int32
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return err
	}

	for _, role := range dest {
		grades, _ := p.permsJobsRoleMap.LoadOrStore(role.Job, syncx.Map[int32, uint64]{})
		grades.Store(role.Grade, role.ID)

		p.permsJobsRoleMap.Store(role.Job, grades)
	}

	return nil
}

func (p *Perms) loadRolePermissions(ctx context.Context) error {
	stmt := tRolePerms.
		SELECT(
			tRolePerms.RoleID.AS("role_id"),
			tRolePerms.PermissionID.AS("id"),
			tRolePerms.Val.AS("val"),
		).
		FROM(tRolePerms.
			INNER_JOIN(tRoles,
				tRoles.ID.EQ(tRolePerms.RoleID)),
		).
		ORDER_BY(
			tRoles.Job.ASC(),
			tRoles.Grade.DESC(),
		)

	var dest []struct {
		RoleID uint64
		ID     uint64
		Val    bool
	}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		return err
	}

	for _, rPerms := range dest {
		perms, _ := p.permsRoleMap.LoadOrStore(rPerms.RoleID, syncx.Map[uint64, bool]{})
		perms.Store(rPerms.ID, rPerms.Val)

		p.permsRoleMap.Store(rPerms.RoleID, perms)
	}

	return nil
}

func (p *Perms) loadRoleAttributes(ctx context.Context) error {
	stmt := tRoleAttrs.
		SELECT(
			tRoleAttrs.AttrID.AS("attr_id"),
			tAttrs.PermissionID.AS("permission_id"),
			tRoleAttrs.RoleID.AS("role_id"),
			tAttrs.Key.AS("key"),
			tAttrs.Type.AS("type"),
			tRoleAttrs.Value.AS("value"),
			tAttrs.ValidValues.AS("valid_values"),
		).
		FROM(
			tRoleAttrs.
				INNER_JOIN(tAttrs,
					tAttrs.ID.EQ(tRoleAttrs.AttrID),
				),
		)

	var dest []struct {
		AttrID       uint64
		PermissionID uint64
		RoleID       uint64
		Key          Key
		Type         AttributeTypes
		Value        string
		ValidValues  string
	}

	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return err
		}
	}

	for _, v := range dest {
		if err := p.addOrUpdateRoleAttributeInMap(v.RoleID, v.PermissionID, v.AttrID, v.Key, v.Type, v.Value, v.ValidValues); err != nil {
			return err
		}
	}

	return nil
}

func (p *Perms) addOrUpdateRoleAttributeInMap(roleId uint64, permId uint64, attrId uint64, key Key, aType AttributeTypes, value string, validValues string) error {
	val := &permissions.AttributeValues{}
	if err := p.convertRawValue(val, value, aType); err != nil {
		return err
	}

	validVals := &permissions.AttributeValues{}
	if err := p.convertRawValue(validVals, validValues, aType); err != nil {
		return err
	}

	p.updateRoleAttributeInMap(roleId, permId, attrId, key, aType, val)

	return nil
}

func (p *Perms) updateRoleAttributeInMap(roleId uint64, permId uint64, attrId uint64, key Key, aType AttributeTypes, value *permissions.AttributeValues) {
	attrMap, _ := p.attrsRoleMap.LoadOrStore(roleId, syncx.Map[uint64, *cacheRoleAttr]{})

	attrMap.Store(attrId, &cacheRoleAttr{
		AttrID: attrId,
		Key:    key,
		Type:   aType,
		Value:  value,
	})

	p.attrsRoleMap.Store(roleId, attrMap)
}

func (p *Perms) removeRoleAttributeFromMap(roleId uint64, attrId uint64) {
	attrMap, ok := p.attrsRoleMap.Load(roleId)
	if !ok {
		return
	}

	attrMap.Delete(attrId)
	p.attrsRoleMap.Store(roleId, attrMap)
}
