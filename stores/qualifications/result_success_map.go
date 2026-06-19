package qualificationsstore

import (
	"context"
	"errors"

	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) syncQualificationResultSuccessMap(
	ctx context.Context,
	tx qrm.DB,
	resultId int64,
	qualificationId int64,
	userId int32,
	successful bool,
) error {
	if err := s.deleteQualificationResultSuccessMapByResultID(ctx, tx, resultId); err != nil {
		return err
	}

	if !successful {
		return nil
	}

	return s.upsertQualificationResultSuccessMap(ctx, tx, resultId, qualificationId, userId)
}

func (s *Store) upsertQualificationResultSuccessMap(
	ctx context.Context,
	tx qrm.DB,
	resultId int64,
	qualificationId int64,
	userId int32,
) error {
	stmt := tQualiResultSuccess.
		INSERT(
			tQualiResultSuccess.QualificationID,
			tQualiResultSuccess.UserID,
			tQualiResultSuccess.ResultID,
		).
		VALUES(
			qualificationId,
			userId,
			resultId,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tQualiResultSuccess.ResultID.SET(mysql.RawInt("VALUES(`result_id`)")),
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) deleteQualificationResultSuccessMapByResultID(
	ctx context.Context,
	tx qrm.DB,
	resultId int64,
) error {
	stmt := tQualiResultSuccess.
		DELETE().
		WHERE(tQualiResultSuccess.ResultID.EQ(mysql.Int64(resultId))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) deleteQualificationResultSuccessMapByQualificationID(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
) error {
	stmt := tQualiResultSuccess.
		DELETE().
		WHERE(tQualiResultSuccess.QualificationID.EQ(mysql.Int64(qualificationId)))

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) rebuildQualificationResultSuccessMapByQualificationID(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
) error {
	if err := s.deleteQualificationResultSuccessMapByQualificationID(
		ctx,
		tx,
		qualificationId,
	); err != nil {
		return err
	}

	stmt := tQualiResult.
		SELECT(
			tQualiResult.QualificationID.AS("qualification_id"),
			tQualiResult.UserID.AS("user_id"),
			mysql.MAX(tQualiResult.ID).AS("result_id"),
		).
		FROM(tQualiResult).
		WHERE(mysql.AND(
			tQualiResult.QualificationID.EQ(mysql.Int64(qualificationId)),
			tQualiResult.DeletedAt.IS_NULL(),
			tQualiResult.Status.EQ(
				mysql.Int32(int32(resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)),
			),
		)).
		GROUP_BY(tQualiResult.QualificationID, tQualiResult.UserID)

	var rows []struct {
		QualificationID int64 `alias:"qualification_id"`
		UserID          int32 `alias:"user_id"`
		ResultID        int64 `alias:"result_id"`
	}
	if err := stmt.QueryContext(ctx, tx, &rows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	for _, row := range rows {
		if err := s.upsertQualificationResultSuccessMap(
			ctx,
			tx,
			row.ResultID,
			row.QualificationID,
			row.UserID,
		); err != nil {
			return err
		}
	}

	return nil
}
