package perms

import (
	"context"
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

var P *perms

type perms struct {
	cancel context.CancelFunc

	canCacheTTL time.Duration
	canCache    *cache.Cache[string, bool]

	permsCacheTTL time.Duration
	permsCache    *cache.Cache[int32, collections.Permissions]
}

func Setup() {
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

	P = &perms{
		cancel: cancel,

		canCacheTTL: 2 * time.Minute,
		canCache:    canCache,

		permsCacheTTL: 2 * time.Minute,
		permsCache:    permsCache,
	}
}

func Stop() {
	if P.cancel != nil {
		P.cancel()
	}
}
