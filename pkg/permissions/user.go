package permissions

import (
	"time"

	"github.com/galexrt/arpanet/pkg/permissions/collections"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

type ArpanetPermissions struct {
	ID          uint64 `sql:"primary_key"`
	Name        string
	GuardName   string
	Description *string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func GetAllPermissionsOfUser(userID int32) (collections.Permissions, error) {
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

func GetAllPermissionsByPrefixOfUser(userID int32, prefix string) (collections.Permissions, error) {
	// TODO

	return nil, nil
}
