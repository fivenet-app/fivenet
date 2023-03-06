package permissions

import (
	"context"
	"strings"

	"github.com/galexrt/arpanet/pkg/permissions/collections"
	"github.com/galexrt/arpanet/pkg/permissions/helpers"
	"github.com/galexrt/arpanet/proto/common"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

func CreatePermission(name string, description string) error {
	r := table.ArpanetPermissions
	stmt := r.INSERT(r.Name,
		r.GuardName,
		r.Description).
		VALUES(name, helpers.Guard(name), description)

	_, err := stmt.ExecContext(context.TODO(), query.DB)
	return err
}

func GetAllPermissions() (collections.Permissions, error) {
	r := table.ArpanetPermissions
	stmt := r.SELECT(r.AllColumns).FROM(table.ArpanetPermissions)

	var dest collections.Permissions
	err := stmt.QueryContext(context.TODO(), query.DB, &dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func Can(user *common.Character, perm ...string) bool {
	return CanID(user.UserID, perm...)
}

func CanID(userID int32, perm ...string) bool {
	return canID(userID, helpers.Guard(strings.Join(perm, ".")))
}

func canID(userID int32, guardName string) bool {
	ap := table.ArpanetPermissions
	aup := table.ArpanetUserPermissions
	aur := table.ArpanetUserRoles
	arp := table.ArpanetRolePermissions

	stmt := ap.
		SELECT(
			ap.ID,
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
	err := stmt.Query(query.DB, &dest)
	if err != nil {
		return false
	}

	return dest.ID > 0
}
