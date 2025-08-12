package access

import (
	"context"
	"errors"

	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

// GetEntry retrieves a single qualification access entry by its ID, joining with the qualifications table for additional info.
func (a *Qualifications[U, T, AccessLevel]) GetEntry(
	ctx context.Context,
	tx qrm.DB,
	id uint64,
) (T, error) {
	stmt := a.selectTable.
		SELECT(
			a.selectColumns.ID,
			a.selectColumns.TargetID,
			a.selectColumns.Access,
			a.selectColumns.QualificationId,
			tQualifications.ID,
			tQualifications.Job,
			tQualifications.Abbreviation,
			tQualifications.Title,
		).
		FROM(
			a.table.
				INNER_JOIN(tQualifications,
					tQualifications.ID.EQ(a.selectColumns.QualificationId),
				),
		).
		WHERE(jet.AND(
			a.selectColumns.ID.EQ(jet.Uint64(id)),
			tQualifications.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	var dest T
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

// CreateEntry inserts a new qualification access entry for a given targetId and entry.
func (a *Qualifications[U, T, AccessLevel]) CreateEntry(
	ctx context.Context,
	tx qrm.DB,
	targetId uint64,
	entry T,
) error {
	stmt := a.table.
		INSERT(
			a.columns.TargetID,
			a.columns.Access,
			a.columns.QualificationId,
		).
		VALUES(
			targetId,
			entry.GetAccess(),
			entry.GetQualificationId(),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

// UpdateEntry updates an existing qualification access entry for a given targetId and entry.
func (a *Qualifications[U, T, AccessLevel]) UpdateEntry(
	ctx context.Context,
	tx qrm.DB,
	targetId uint64,
	entry T,
) error {
	stmt := a.table.
		UPDATE(
			a.columns.Access,
			a.columns.QualificationId,
		).
		SET(
			entry.GetAccess(),
			entry.GetQualificationId(),
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

// DeleteEntry deletes a qualification access entry by its ID and targetId.
func (a *Qualifications[U, T, AccessLevel]) DeleteEntry(
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

// DeleteEntryWithCondition deletes a qualification access entry matching a custom condition and targetId.
func (a *Qualifications[U, T, AccessLevel]) DeleteEntryWithCondition(
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
