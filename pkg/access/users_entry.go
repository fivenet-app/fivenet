package access

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

// GetEntry retrieves a single user access entry by its ID, joining with user_short for additional user info.
func (a *Users[U, T, AccessLevel]) GetEntry(ctx context.Context, tx qrm.DB, id uint64) (T, error) {
	tUsers := tables.User().AS("user_short")

	stmt := a.selectTable.
		SELECT(
			a.selectColumns.ID,
			a.selectColumns.TargetID,
			a.selectColumns.Access,
			a.selectColumns.UserId,
			tUsers.ID,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			a.table.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(a.selectColumns.UserId),
				),
		).
		WHERE(
			a.selectColumns.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	var dest T
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

// CreateEntry inserts a new user access entry for a given targetId and entry.
func (a *Users[U, T, AccessLevel]) CreateEntry(
	ctx context.Context,
	tx qrm.DB,
	targetId uint64,
	entry T,
) error {
	stmt := a.table.
		INSERT(
			a.columns.TargetID,
			a.columns.Access,
			a.columns.UserId,
		).
		VALUES(
			targetId,
			entry.GetAccess(),
			entry.GetUserId(),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

// UpdateEntry updates an existing user access entry for a given targetId and entry.
func (a *Users[U, T, AccessLevel]) UpdateEntry(
	ctx context.Context,
	tx qrm.DB,
	targetId uint64,
	entry T,
) error {
	stmt := a.table.
		UPDATE(
			a.columns.Access,
			a.columns.UserId,
		).
		SET(
			entry.GetAccess(),
			entry.GetUserId(),
		).
		WHERE(jet.AND(
			a.columns.ID.EQ(jet.Uint64(entry.GetId())),
			a.columns.TargetID.EQ(jet.Uint64(targetId)),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

// DeleteEntry deletes a user access entry by its ID and targetId.
func (a *Users[U, T, AccessLevel]) DeleteEntry(
	ctx context.Context,
	tx qrm.DB,
	targetId uint64,
	id uint64,
) error {
	stmt := a.table.
		DELETE().
		WHERE(jet.AND(
			a.columns.ID.EQ(jet.Uint64(id)),
			a.columns.TargetID.EQ(jet.Uint64(targetId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	return nil
}

// DeleteEntryWithCondition deletes a user access entry matching a custom condition and targetId.
func (a *Users[U, T, AccessLevel]) DeleteEntryWithCondition(
	ctx context.Context,
	tx qrm.DB,
	condition jet.BoolExpression,
	targetId uint64,
) error {
	stmt := a.table.
		DELETE().
		WHERE(jet.AND(
			condition,
			a.columns.TargetID.EQ(jet.Uint64(targetId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	return nil
}
