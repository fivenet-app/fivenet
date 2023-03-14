package perms

import (
	"strconv"
	"strings"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (p *Perms) GetAllPermissionsOfUser(userID int32) (collections.Permissions, error) {
	if cached, ok := p.permsCache.Get(userID); ok {
		return cached, nil
	}

	stmt := ap.SELECT(
		ap.AllColumns,
	).
		FROM(ap).
		WHERE(
			ap.ID.IN(
				aup.
					SELECT(
						aup.PermissionID,
					).
					FROM(
						aup,
					).WHERE(aup.UserID.EQ(jet.Int32(userID))).
					UNION(
						aur.SELECT(arp.PermissionID).
							FROM(aur.INNER_JOIN(arp, arp.RoleID.EQ(aur.RoleID))).
							WHERE(
								aur.UserID.EQ(jet.Int32(userID)),
							),
					),
			),
		)

	var perms collections.Permissions
	if err := stmt.Query(p.db, &perms); err != nil {
		return nil, err
	}

	p.permsCache.Set(userID, perms)

	return perms, nil
}

func (p *Perms) GetAllPermissionsByPrefixOfUser(userID int32, prefix string) (collections.Permissions, error) {
	prefix = helpers.Guard(prefix)

	return p.getAllPermissionsByPrefixOfUser(userID, prefix)
}

func (p *Perms) getAllPermissionsByPrefixOfUser(userID int32, prefix string) (collections.Permissions, error) {
	if cached, ok := p.permsCache.Get(userID); ok {
		return cached.HasPrefix(prefix), nil
	}

	stmt := ap.SELECT(
		ap.AllColumns,
	).
		FROM(ap).
		WHERE(
			jet.AND(
				ap.GuardName.LIKE(jet.String(prefix+"%")),
				ap.ID.IN(
					aup.
						SELECT(
							aup.PermissionID,
						).
						FROM(
							aup,
						).WHERE(aup.UserID.EQ(jet.Int32(userID))).
						UNION(
							aur.SELECT(arp.PermissionID).
								FROM(aur.INNER_JOIN(arp, arp.RoleID.EQ(aur.RoleID))).
								WHERE(
									aur.UserID.EQ(jet.Int32(userID)),
								),
						),
				),
			),
		)

	var perms collections.Permissions
	if err := stmt.Query(p.db, &perms); err != nil {
		return nil, err
	}

	return perms, nil
}

func (p *Perms) GetSuffixOfPermissionsByPrefixOfUser(userID int32, prefix string) ([]string, error) {
	prefix = helpers.Guard(prefix) + "-"

	perms, err := p.getAllPermissionsByPrefixOfUser(userID, prefix)
	if err != nil {
		return nil, err
	}

	suffixes := []string{}
	for _, perm := range perms {
		suffixes = append(suffixes, strings.TrimPrefix(perm.GuardName, prefix))
	}

	return suffixes, nil
}

func (p *Perms) CanID(userID int32, perm ...string) bool {
	return p.canID(userID, helpers.Guard(strings.Join(perm, ".")))
}

func (p *Perms) canID(userID int32, guardName string) bool {
	cacheKey := buildCanCacheKey(userID, guardName)
	if cached, ok := p.canCache.Get(cacheKey); ok {
		return cached
	}

	stmt := ap.
		SELECT(
			ap.ID.AS("id"),
		).
		FROM(
			ap.LEFT_JOIN(aup,
				aup.PermissionID.EQ(ap.ID),
			).
				LEFT_JOIN(aur,
					aur.UserID.EQ(jet.Int32(userID)),
				).
				LEFT_JOIN(arp,
					arp.PermissionID.EQ(ap.ID).
						AND(
							arp.RoleID.EQ(aur.RoleID),
						),
				),
		).
		WHERE(ap.GuardName.EQ(jet.String(guardName))).
		LIMIT(1)

	var dest struct {
		ID int32
	}
	if err := stmt.Query(p.db, &dest); err != nil {
		return false
	}

	result := dest.ID > 0
	p.canCache.Set(cacheKey, result, cache.WithExpiration(p.canCacheTTL))

	return result
}

func buildCanCacheKey(userID int32, guardName string) string {
	return strconv.Itoa(int(userID)) + "-" + guardName
}
