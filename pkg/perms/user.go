package perms

import (
	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (p *perms) GetAllPermissionsOfUser(userID int32) (collections.Permissions, error) {
	ap := table.ArpanetPermissions
	aup := table.ArpanetUserPermissions
	arp := table.ArpanetRolePermissions
	aur := table.ArpanetUserRoles

	stmt := ap.SELECT(ap.AllColumns).
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
	if err := stmt.Query(query.DB, &perms); err != nil {
		return nil, err
	}

	return perms, nil
}

func (p *perms) GetAllPermissionsByPrefixOfUser(userID int32, prefix string) (collections.Permissions, error) {
	// TODO

	return nil, nil
}
