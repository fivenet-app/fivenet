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
)

type Permissions interface {
	GetAllPermissions() ([]*permissions.Permission, error)
	GetPermissionsByIDs(ids ...uint64) ([]*permissions.Permission, error)
	CreatePermission(category Category, name Name) (uint64, error)
	GetPermissionsOfUser(userInfo *userinfo.UserInfo) (collections.Permissions, error)

	GetJobRoles(job string) (collections.Roles, error)
	GetJobRolesUpTo(job string, grade int32) (collections.Roles, error)
	GetClosesJobRole(job string, grade int32) (*model.FivenetRoles, error)
	CountRolesForJob(prefix string) (int64, error)

	GetRole(id uint64) (*model.FivenetRoles, error)
	GetRoleByJobAndGrade(job string, grade int32) (*model.FivenetRoles, error)
	GetRolePermissions(id uint64) ([]*permissions.Permission, error)

	CreateRole(job string, grade int32) (*model.FivenetRoles, error)
	DeleteRole(id uint64) error
	UpdateRolePermissions(id uint64, perms ...AddPerm) error
	RemovePermissionsFromRole(id uint64, perms ...uint64) error

	Can(userInfo *userinfo.UserInfo, category Category, name Name) bool

	GetAttribute(category Category, name Name, key Key) (*permissions.RoleAttribute, error)
	CreateAttribute(permId uint64, key Key, aType AttributeTypes, validValues any) (uint64, error)
	UpdateAttribute(attributeId uint64, permId uint64, key Key, aType AttributeTypes, validValues any) error
	GetRoleAttributes(job string, grade int32) ([]*permissions.RoleAttribute, error)
	GetAllAttributes(job string) ([]*permissions.RoleAttribute, error)
	AddOrUpdateAttributesToRole(attrs ...*permissions.RoleAttribute) error
	UpdateRoleAttributes(attrs ...*permissions.RoleAttribute) error
	RemoveAttributesFromRole(roleId uint64, attrs ...*permissions.RoleAttribute) error

	Attr(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (any, error)
}

type Perms struct {
	db *sql.DB

	ctx context.Context

	// Guard name to permission ID
	guardToPermIDMap syncx.Map[string, uint64]
	// Job name to map of grade numbers to role ID
	jobsToRoleIDMap syncx.Map[string, map[int32]uint64]
	// Role ID to map of permissions ID and result
	rolePermsMap syncx.Map[uint64, map[uint64]bool]

	// Perm ID to map Key -> cached attribute
	permIDToAttrsMap syncx.Map[uint64, map[Key]cacheAttr]
	// Role ID to map of Key -> cached role attribute
	roleIDToAttrMap syncx.Map[uint64, map[Key]cacheRoleAttr]

	userCanCacheTTL time.Duration
	userCanCache    *cache.Cache[int32, map[uint64]bool]
}

func New(ctx context.Context, db *sql.DB) *Perms {
	userCanCache := cache.NewContext(
		ctx,
		cache.AsLRU[int32, map[uint64]bool](lru.WithCapacity(128)),
		cache.WithJanitorInterval[int32, map[uint64]bool](15*time.Second),
	)

	p := &Perms{
		db: db,

		ctx: ctx,

		guardToPermIDMap: syncx.Map[string, uint64]{},
		jobsToRoleIDMap:  syncx.Map[string, map[int32]uint64]{},
		rolePermsMap:     syncx.Map[uint64, map[uint64]bool]{},

		permIDToAttrsMap: syncx.Map[uint64, map[Key]cacheAttr]{},
		roleIDToAttrMap:  syncx.Map[uint64, map[Key]cacheRoleAttr]{},

		userCanCacheTTL: 30 * time.Second,
		userCanCache:    userCanCache,
	}

	p.load()

	return p
}

type cacheAttr struct {
	ID           uint64
	PermissionID uint64
	Key          Key
	Type         AttributeTypes
	ValidValues  *permissions.AttributeValues
}

type cacheRoleAttr struct {
	Type  AttributeTypes
	Value *permissions.AttributeValues
}

func (p *Perms) load() error {
	if err := p.loadRoleIDs(); err != nil {
		return err
	}

	if err := p.loadRolePermissions(); err != nil {
		return err
	}

	if err := p.loadRoleAttributes(); err != nil {
		return err
	}

	return nil
}

func (p *Perms) loadRoleIDs() error {
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
	if err := stmt.Query(p.db, &dest); err != nil {
		return err
	}

	for _, v := range dest {
		grades, loaded := p.jobsToRoleIDMap.LoadOrStore(v.Job, map[int32]uint64{
			v.Grade: v.ID,
		})
		if loaded {
			grades[v.Grade] = v.ID
		}
	}

	return nil
}

func (p *Perms) loadRolePermissions() error {
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
	if err := stmt.Query(p.db, &dest); err != nil {
		return err
	}

	for _, v := range dest {
		perms, loaded := p.rolePermsMap.LoadOrStore(v.RoleID, map[uint64]bool{
			v.ID: v.Val,
		})
		if loaded {
			perms[v.ID] = v.Val
		}
	}

	return nil
}

func (p *Perms) loadRoleAttributes() error {
	stmt := tRoleAttrs.
		SELECT(
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
		RoleID      uint64
		Key         Key
		Type        AttributeTypes
		Value       string
		ValidValues string
	}

	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return err
		}
	}

	for _, v := range dest {
		if err := p.addOrUpdateRoleAttributeInMap(v.RoleID, v.Key, v.Type, v.Value, v.ValidValues); err != nil {
			return err
		}
	}

	return nil
}

func (p *Perms) addOrUpdateRoleAttributeInMap(roleId uint64, key Key, aType AttributeTypes, value string, validValues string) error {
	val := &permissions.AttributeValues{}
	if err := p.convertRawValue(val, value, aType); err != nil {
		return err
	}

	validVals := &permissions.AttributeValues{}
	if err := p.convertRawValue(validVals, validValues, aType); err != nil {
		return err
	}

	p.updateRoleAttributeInMap(roleId, key, aType, val)

	return nil
}

func (p *Perms) updateRoleAttributeInMap(roleId uint64, key Key, aType AttributeTypes, value *permissions.AttributeValues) {
	attrMap, ok := p.roleIDToAttrMap.Load(roleId)
	if !ok || attrMap == nil {
		attrMap = map[Key]cacheRoleAttr{}
	}

	attrMap[key] = cacheRoleAttr{
		Type:  aType,
		Value: value,
	}

	p.roleIDToAttrMap.Store(roleId, attrMap)
}

func (p *Perms) removeRoleAttributeFromMap(roleId uint64, key Key) {
	attrMap, ok := p.roleIDToAttrMap.Load(roleId)
	if !ok || attrMap == nil {
		return
	}

	if _, ok := attrMap[key]; ok {
		delete(attrMap, key)
		p.roleIDToAttrMap.Store(roleId, attrMap)
	}
}
