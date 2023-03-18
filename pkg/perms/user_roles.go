package perms

import (
	"errors"
	"fmt"

	"github.com/galexrt/arpanet/pkg/dbutils"
	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (p *Perms) GetUserRoles(userId int32) (collections.Roles, error) {
	var dest collections.Roles

	stmt := aur.SELECT(
		aur.RoleID,
	).
		FROM(aur.
			INNER_JOIN(ar,
				ar.ID.EQ(aur.RoleID)),
		)

	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return dest, nil
}

func (p *Perms) AddUserRoles(userId int32, roles ...string) error {
	if len(roles) == 0 {
		return nil
	}

	roleGuards := make([]jet.Expression, len(roles))
	for i := 0; i < len(roles); i++ {
		roleGuards[i] = jet.String(helpers.Guard(roles[i]))
	}

	stmt := aur.INSERT(
		aur.UserID,
		aur.RoleID,
	).
		QUERY(
			ar.SELECT(
				jet.Int32(userId),
				ar.ID,
			).
				FROM(ar).
				WHERE(ar.GuardName.IN(roleGuards...)),
		)

	if _, err := stmt.ExecContext(p.ctx, p.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (p *Perms) RemoveUserRoles(userId int32, roles ...string) error {
	if len(roles) == 0 {
		return nil
	}

	roleGuards := []jet.Expression{}
	for i := 0; i < len(roles); i++ {
		roleGuards[i] = jet.String(helpers.Guard(roles[i]))
	}

	stmt := aur.DELETE().
		USING(ar).
		WHERE(jet.AND(
			aur.UserID.EQ(jet.Int32(userId)),
			ar.GuardName.IN(roleGuards...),
			aur.RoleID.EQ(ar.ID),
		))

	fmt.Println(stmt.DebugSql())

	if _, err := stmt.ExecContext(p.ctx, p.db); err != nil {
		return err
	}

	return nil
}
