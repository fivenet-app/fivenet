package perms

import (
	"errors"

	"github.com/galexrt/fivenet/pkg/dbutils"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/pkg/perms/helpers"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (p *Perms) GetUserRoles(userId int32) (collections.Roles, error) {
	var dest collections.Roles

	stmt := aur.
		SELECT(
			ar.ID,
			ar.CreatedAt,
			ar.UpdatedAt,
			ar.Name,
			ar.GuardName,
			ar.Description,
		).
		FROM(aur.
			INNER_JOIN(ar,
				ar.ID.EQ(aur.RoleID)),
		).
		WHERE(
			aur.UserID.EQ(jet.Int32(userId)),
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

	stmt := aur.
		INSERT(
			aur.UserID,
			aur.RoleID,
		).
		QUERY(
			ar.
				SELECT(
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

	roleGuards := make([]jet.Expression, len(roles))
	for i := 0; i < len(roles); i++ {
		roleGuards[i] = jet.String(helpers.Guard(roles[i]))
	}

	stmt := aur.
		DELETE().
		USING(aur.
			INNER_JOIN(ar,
				ar.GuardName.IN(roleGuards...),
			),
		).
		WHERE(jet.AND(
			aur.UserID.EQ(jet.Int32(userId)),
			aur.RoleID.EQ(ar.ID),
		))

	if _, err := stmt.ExecContext(p.ctx, p.db); err != nil {
		return err
	}

	return nil
}
