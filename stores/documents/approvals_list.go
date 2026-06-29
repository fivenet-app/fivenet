package documentsstore

import (
	"context"
	"errors"

	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListApprovals(
	ctx context.Context,
	q ListApprovalsQuery,
) (resourcesdatabase.DataCount, []*documentsapproval.Approval, error) {
	tApprovals := table.FivenetDocumentsApprovals.AS("approval")

	condition := tApprovals.DocumentID.EQ(mysql.Int64(q.DocumentID))
	if q.SnapshotDate != nil {
		condition = condition.AND(
			tApprovals.SnapshotDate.EQ(mysql.DateTimeT(q.SnapshotDate.AsTime())),
		)
	}
	if q.Status > 0 {
		condition = condition.AND(tApprovals.Status.EQ(mysql.Int32(int32(q.Status))))
	}
	if q.UserID > 0 {
		condition = condition.AND(tApprovals.UserID.EQ(mysql.Int32(q.UserID)))
	}
	if q.TaskID > 0 {
		condition = condition.AND(tApprovals.TaskID.EQ(mysql.Int64(q.TaskID)))
	}

	countStmt := mysql.
		SELECT(mysql.COUNT(tApprovals.ID).AS("data_count.total")).
		FROM(tApprovals).
		WHERE(condition)

	var count resourcesdatabase.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	if count.Total <= 0 {
		return count, []*documentsapproval.Approval{}, nil
	}

	tUser := table.FivenetUser.AS("usershort")
	tStamp := table.FivenetDocumentsStamps.AS("stamp")

	stmt := tApprovals.
		SELECT(
			tApprovals.ID,
			tApprovals.DocumentID,
			tApprovals.SnapshotDate,
			tApprovals.UserID,
			tApprovals.UserJob,
			tApprovals.UserJobGrade,
			tApprovals.PayloadSvg,
			tApprovals.StampID,
			tApprovals.Status,
			tApprovals.Comment,
			tApprovals.TaskID,
			tApprovals.CreatedAt,
			tApprovals.RevokedAt,

			tStamp.ID,
			tStamp.SvgTemplate,

			tUser.ID,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.Job,
			tUser.JobGrade,
		).
		FROM(
			tApprovals.
				LEFT_JOIN(tUser, tUser.ID.EQ(tApprovals.UserID)).
				LEFT_JOIN(tStamp, tApprovals.ID.EQ(tStamp.ID)),
		).
		WHERE(condition).
		ORDER_BY(tApprovals.Status.ASC(), tApprovals.CreatedAt.DESC()).
		LIMIT(15)

	var approvals []*documentsapproval.Approval
	if err := stmt.QueryContext(ctx, s.db, &approvals); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	return count, approvals, nil
}

func (s *Store) GetApproval(
	ctx context.Context,
	db qrm.DB,
	approvalID int64,
) (*documentsapproval.Approval, error) {
	tApprovals := table.FivenetDocumentsApprovals.AS("approval")

	stmt := tApprovals.
		SELECT(
			tApprovals.ID,
			tApprovals.DocumentID,
			tApprovals.SnapshotDate,
			tApprovals.UserID,
			tApprovals.UserJob,
			tApprovals.UserJobGrade,
			tApprovals.PayloadSvg,
			tApprovals.StampID,
			tApprovals.Status,
			tApprovals.Comment,
			tApprovals.TaskID,
			tApprovals.CreatedAt,
			tApprovals.RevokedAt,
		).
		FROM(tApprovals).
		WHERE(tApprovals.ID.EQ(mysql.Int64(approvalID))).
		LIMIT(1)

	var approval documentsapproval.Approval
	if err := stmt.QueryContext(ctx, db, &approval); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if approval.GetId() == 0 {
		return nil, nil
	}

	return &approval, nil
}
