package documentsstore

import (
	"context"

	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ResetApprovalProgress(ctx context.Context, tx qrm.DB, documentID int64) error {
	tApprovals := table.FivenetDocumentsApprovals
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	if _, err := tApprovals.
		UPDATE().
		SET(
			tApprovals.Status.SET(mysql.Int32(int32(documentsapproval.ApprovalStatus_APPROVAL_STATUS_REVOKED))),
			tApprovals.RevokedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(tApprovals.DocumentID.EQ(mysql.Int64(documentID))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	if _, err := tApprovalTasks.
		UPDATE().
		SET(
			tApprovalTasks.Status.SET(mysql.Int32(int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
			tApprovalTasks.CompletedAt.SET(mysql.TimestampExp(mysql.NULL)),
		).
		WHERE(tApprovalTasks.DocumentID.EQ(mysql.Int64(documentID))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
