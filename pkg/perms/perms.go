package perms

import (
	"context"
	"database/sql"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
	"github.com/galexrt/fivenet/query/fivenet/model"
)

type Permissions interface {
	GetAllPermissions() ([]*permissions.Permission, error)
	GetPermissionsByIDs(ids ...uint64) ([]*permissions.Permission, error)
	CreatePermission(category Category, name Name) (uint64, error)
	GetPermissionsOfUser(userId int32, job string, grade int32) (collections.Permissions, error)

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

	Can(userId int32, job string, grade int32, category Category, name Name) bool

	GetAttribute(permId uint64, key Key) (*model.FivenetAttrs, error)
	CreateAttribute(permId uint64, key Key, aType AttributeTypes, validValues string) (uint64, error)
	UpdateAttribute(attributeId uint64, permId uint64, key Key, aType AttributeTypes, validValues string) error
	GetRoleAttributes(job string, grade int32) ([]*permissions.RoleAttribute, error)
	GetAllAttributes(job string) ([]*permissions.RoleAttribute, error)
	AddAttributesToRole(roleId uint64, attrs ...*permissions.RoleAttribute) error
	UpdateRoleAttributes(roleId uint64, attrs ...*permissions.RoleAttribute) error
	RemoveAttributesFromRole(roleId uint64, attrs ...*permissions.RoleAttribute) error

	Attr(userId int32, job string, grade int32, category Category, name Name, key Key) (any, error)
}

type Perms struct {
	db *sql.DB

	ctx context.Context

	guardToPermIDMap syncx.Map[string, uint64]
	permIdToAttrsMap syncx.Map[uint64, map[Key]Attr]
	jobsToRoleIDMap  syncx.Map[string, map[int32]uint64]
	rolePermsMap     syncx.Map[uint64, map[uint64]bool]

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
		permIdToAttrsMap: syncx.Map[uint64, map[Key]Attr]{},
		jobsToRoleIDMap:  syncx.Map[string, map[int32]uint64]{},
		rolePermsMap:     syncx.Map[uint64, map[uint64]bool]{},

		userCanCacheTTL: 30 * time.Second,
		userCanCache:    userCanCache,
	}

	p.load()

	return p
}

func (p *Perms) load() error {
	if err := p.loadRoleIDs(); err != nil {
		return err
	}

	if err := p.loadRolePermissions(); err != nil {
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
