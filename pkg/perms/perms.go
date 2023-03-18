package perms

import (
	"context"
	"database/sql"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/query/arpanet/table"
)

var (
	ap  = table.ArpanetPermissions
	ar  = table.ArpanetRoles
	arp = table.ArpanetRolePermissions
	aup = table.ArpanetUserPermissions
	aur = table.ArpanetUserRoles
)

type Permissions interface {
	GetAllPermissions() (collections.Permissions, error)
	CreatePermission(name string, description string) error

	GetAllPermissionsOfUser(userId int32) (collections.Permissions, error)
	GetAllPermissionsByPrefixOfUser(userId int32, prefix string) (collections.Permissions, error)
	GetSuffixOfPermissionsByPrefixOfUser(userId int32, prefix string) ([]string, error)

	GetRoles(prefix string) (collections.Roles, error)

	GetUserRoles(userId int32) (collections.Roles, error)
	AddUserRoles(userId int32, roles ...string) error
	RemoveUserRoles(userId int32, roles ...string) error

	Can(userId int32, perm ...string) bool
}

type Perms struct {
	db *sql.DB

	ctx    context.Context
	cancel context.CancelFunc

	canCacheTTL time.Duration
	canCache    *cache.Cache[string, bool]

	permsCacheTTL time.Duration
	permsCache    *cache.Cache[int32, collections.Permissions]
}

func New(db *sql.DB) *Perms {
	ctx, cancel := context.WithCancel(context.Background())

	canCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, bool](lru.WithCapacity(128)),
		cache.WithJanitorInterval[string, bool](15*time.Second),
	)
	permsCache := cache.NewContext(
		ctx,
		cache.AsLRU[int32, collections.Permissions](lru.WithCapacity(128)),
		cache.WithJanitorInterval[int32, collections.Permissions](15*time.Second),
	)

	return &Perms{
		db: db,

		ctx:    ctx,
		cancel: cancel,

		canCacheTTL: 2 * time.Minute,
		canCache:    canCache,

		permsCacheTTL: 2 * time.Minute,
		permsCache:    permsCache,
	}
}

func (p *Perms) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
}
