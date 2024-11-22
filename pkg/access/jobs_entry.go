package access

import (
	"context"
	"errors"

	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (a *Jobs[U, T, AccessLevel]) GetEntry(ctx context.Context, tx qrm.DB, id uint64) (T, error) {
	stmt := a.selectTable.
		SELECT(
			a.selectColumns.ID,
			a.selectColumns.CreatedAt,
			a.selectColumns.TargetID,
			a.selectColumns.Access,
			a.selectColumns.Job,
			a.selectColumns.MinimumGrade,
		).
		FROM(a.selectTable).
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

func (a *Jobs[U, T, AccessLevel]) CreateEntry(ctx context.Context, tx qrm.DB, targetId uint64, entry T) error {
	stmt := a.table.
		INSERT(
			a.columns.TargetID,
			a.columns.Access,
			a.columns.Job,
			a.columns.MinimumGrade,
		).
		VALUES(
			targetId,
			entry.GetAccess(),
			entry.GetJob(),
			entry.GetMinimumGrade(),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (a *Jobs[U, T, AccessLevel]) UpdateEntry(ctx context.Context, tx qrm.DB, targetId uint64, entry T) error {
	stmt := a.table.
		UPDATE(
			a.columns.Access,
			a.columns.Job,
			a.columns.MinimumGrade,
		).
		SET(
			entry.GetAccess(),
			entry.GetJob(),
			entry.GetMinimumGrade(),
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

func (a *Jobs[U, T, AccessLevel]) DeleteEntry(ctx context.Context, tx qrm.DB, targetId uint64, id uint64) error {
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
