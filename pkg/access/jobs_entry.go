package access

import (
	"context"
	"errors"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

// GetEntry retrieves a single job access entry by its ID.
func (a *Jobs[U, T, AccessLevel]) GetEntry(ctx context.Context, tx qrm.DB, id int64) (T, error) {
	stmt := a.selectTable.
		SELECT(
			a.selectColumns.ID,
			a.selectColumns.TargetID,
			a.selectColumns.Access,
			a.selectColumns.Job,
			a.selectColumns.MinimumGrade,
		).
		FROM(a.selectTable).
		WHERE(
			a.selectColumns.ID.EQ(mysql.Int64(id)),
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

// CreateEntry inserts a new job access entry for a given targetId and entry.
func (a *Jobs[U, T, AccessLevel]) CreateEntry(
	ctx context.Context,
	tx qrm.DB,
	targetId int64,
	entry T,
) error {
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

// UpdateEntry updates an existing job access entry for a given targetId and entry.
func (a *Jobs[U, T, AccessLevel]) UpdateEntry(
	ctx context.Context,
	tx qrm.DB,
	targetId int64,
	entry T,
) error {
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
		WHERE(mysql.AND(
			a.columns.ID.EQ(mysql.Int64(entry.GetId())),
			a.columns.TargetID.EQ(mysql.Int64(targetId)),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

// DeleteEntry deletes a job access entry by its ID and targetId.
func (a *Jobs[U, T, AccessLevel]) DeleteEntry(
	ctx context.Context,
	tx qrm.DB,
	targetId int64,
	id int64,
) error {
	stmt := a.table.
		DELETE().
		WHERE(mysql.AND(
			a.columns.ID.EQ(mysql.Int64(id)),
			a.columns.TargetID.EQ(mysql.Int64(targetId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	return nil
}
