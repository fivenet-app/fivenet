package perms

import (
	"context"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
	"github.com/galexrt/arpanet/proto/common"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var P *perms

type perms struct {
	canCache   *cache.Cache[string, bool]
	permsCache *cache.Cache[string, []string]
}

func Setup() {
	canCache := cache.New(
		cache.AsLRU[string, bool](lru.WithCapacity(128)),
		cache.WithJanitorInterval[string, bool](15*time.Second),
	)
	permsCache := cache.New(
		cache.AsLRU[string, []string](lru.WithCapacity(128)),
		cache.WithJanitorInterval[string, []string](15*time.Second),
	)

	P = &perms{
		canCache:   canCache,
		permsCache: permsCache,
	}
}

// TODO Use https://github.com/Code-Hex/go-generics-cache as a cache

func (p *perms) CreatePermission(name string, description string) error {
	r := table.ArpanetPermissions
	stmt := r.INSERT(r.Name,
		r.GuardName,
		r.Description).
		VALUES(name, helpers.Guard(name), description)

	_, err := stmt.ExecContext(context.TODO(), query.DB)
	return err
}

func (p *perms) GetAllPermissions() (collections.Permissions, error) {
	r := table.ArpanetPermissions
	stmt := r.SELECT(r.AllColumns).FROM(table.ArpanetPermissions)

	var dest collections.Permissions
	err := stmt.QueryContext(context.TODO(), query.DB, &dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (p *perms) Can(user common.IGetUserID, perm ...string) bool {
	return p.CanID(user.GetUserID(), perm...)
}

func (p *perms) CanID(userID int32, perm ...string) bool {
	return p.canID(userID, helpers.Guard(strings.Join(perm, ".")))
}

func (p *perms) canID(userID int32, guardName string) bool {
	ap := table.ArpanetPermissions
	aup := table.ArpanetUserPermissions
	aur := table.ArpanetUserRoles
	arp := table.ArpanetRolePermissions

	stmt := ap.
		SELECT(
			ap.ID.AS("id"),
		).
		FROM(
			ap.LEFT_JOIN(aup,
				aup.PermissionID.EQ(ap.ID),
			).
				LEFT_JOIN(aur, aur.UserID.EQ(jet.Int32(userID))).
				LEFT_JOIN(arp, arp.PermissionID.EQ(ap.ID).
					AND(
						arp.RoleID.EQ(aur.RoleID),
					)),
		).
		WHERE(ap.GuardName.EQ(jet.String(guardName))).LIMIT(1)

	var dest struct {
		ID int32
	}
	if err := stmt.Query(query.DB, &dest); err != nil {
		return false
	}

	return dest.ID > 0
}
