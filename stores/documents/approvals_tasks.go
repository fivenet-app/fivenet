package documentsstore

import (
	"context"
	"time"

	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CreateApprovalTasks(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
	documentID int64,
	snapDate *timestamp.Timestamp,
	seeds []*pbdocuments.ApprovalTaskSeed,
) (int32, int32, error) {
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	created := int32(0)
	ensured := int32(0)

	for _, seed := range seeds {
		isUser := seed.GetUserId() > 0
		if isUser {
			var cnt struct{ C int32 }
			if err := tApprovalTasks.
				SELECT(mysql.COUNT(tApprovalTasks.ID).AS("C")).
				FROM(tApprovalTasks).
				WHERE(mysql.AND(
					tApprovalTasks.DocumentID.EQ(mysql.Int64(documentID)),
					tApprovalTasks.SnapshotDate.EQ(dbutils.TimestampToMySQLDateTimeSec(snapDate)),
					tApprovalTasks.AssigneeKind.EQ(
						mysql.Int32(
							int32(
								documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER,
							),
						),
					),
					tApprovalTasks.UserID.EQ(mysql.Int32(seed.GetUserId())),
				)).
				LIMIT(1).
				QueryContext(ctx, tx, &cnt); err != nil {
				return 0, 0, err
			}
			if cnt.C > 0 {
				ensured++
				continue
			}

			if _, err := tApprovalTasks.
				INSERT(
					tApprovalTasks.DocumentID,
					tApprovalTasks.SnapshotDate,
					tApprovalTasks.AssigneeKind,
					tApprovalTasks.UserID,
					tApprovalTasks.Label,
					tApprovalTasks.SignatureRequired,
					tApprovalTasks.SlotNo,
					tApprovalTasks.Status,
					tApprovalTasks.Comment,
					tApprovalTasks.DueAt,
					tApprovalTasks.CreatorID,
					tApprovalTasks.CreatorJob,
				).
				VALUES(
					documentID,
					dbutils.TimestampToMySQLDateTimeSec(snapDate),
					int32(documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER),
					seed.GetUserId(),
					seed.GetLabel(),
					seed.GetSignatureRequired(),
					1,
					int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
					seed.GetComment(),
					dbutils.TimestampToMySQLDateTime(seed.GetDueAt()),
					userInfo.GetUserId(),
					userInfo.GetJob(),
				).
				ExecContext(ctx, tx); err != nil {
				return 0, 0, err
			}
			created++
			continue
		}

		slots := seed.GetSlots()
		if slots <= 0 {
			slots = 1
		}

		var have struct{ C int32 }
		if err := mysql.
			SELECT(mysql.COUNT(tApprovalTasks.ID).AS("C")).
			FROM(tApprovalTasks).
			WHERE(mysql.AND(
				tApprovalTasks.DocumentID.EQ(mysql.Int64(documentID)),
				tApprovalTasks.SnapshotDate.EQ(dbutils.TimestampToMySQLDateTimeSec(snapDate)),
				tApprovalTasks.AssigneeKind.EQ(
					mysql.Int32(
						int32(
							documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE,
						),
					),
				),
				tApprovalTasks.Job.EQ(mysql.String(seed.GetJob())),
				tApprovalTasks.MinimumGrade.EQ(mysql.Int32(seed.GetMinimumGrade())),
				tApprovalTasks.Status.EQ(
					mysql.Int32(
						int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
					),
				),
			)).
			LIMIT(1).
			QueryContext(ctx, tx, &have); err != nil {
			return 0, 0, err
		}
		if have.C >= slots {
			ensured++
			continue
		}

		ins := tApprovalTasks.
			INSERT(
				tApprovalTasks.DocumentID,
				tApprovalTasks.SnapshotDate,
				tApprovalTasks.AssigneeKind,
				tApprovalTasks.Job,
				tApprovalTasks.MinimumGrade,
				tApprovalTasks.Label,
				tApprovalTasks.SignatureRequired,
				tApprovalTasks.SlotNo,
				tApprovalTasks.Status,
				tApprovalTasks.Comment,
				tApprovalTasks.DueAt,
				tApprovalTasks.CreatorID,
				tApprovalTasks.CreatorJob,
			)

		for slot := have.C + 1; slot <= slots; slot++ {
			ins = ins.VALUES(
				documentID,
				dbutils.TimestampToMySQLDateTimeSec(snapDate),
				int32(documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE),
				seed.GetJob(),
				seed.GetMinimumGrade(),
				seed.GetLabel(),
				seed.GetSignatureRequired(),
				slot,
				int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
				seed.GetComment(),
				dbutils.TimestampToMySQLDateTime(seed.GetDueAt()),
				userInfo.GetUserId(),
				userInfo.GetJob(),
			)
		}

		ins = ins.
			ON_DUPLICATE_KEY_UPDATE(
				tApprovalTasks.Job.SET(mysql.RawString("VALUES(`job`)")),
				tApprovalTasks.MinimumGrade.SET(mysql.RawInt("VALUES(`minimum_grade`)")),
				tApprovalTasks.Label.SET(mysql.RawString("VALUES(`label`)")),
				tApprovalTasks.SignatureRequired.SET(mysql.RawBool("VALUES(`signature_required`)")),
				tApprovalTasks.SlotNo.SET(mysql.RawInt("VALUES(`slot_no`)")),
				tApprovalTasks.Comment.SET(mysql.RawString("VALUES(`comment`)")),
				tApprovalTasks.DueAt.SET(mysql.RawTimestamp("VALUES(`due_at`)")),
			)

		if _, err := ins.ExecContext(ctx, tx); err != nil {
			return 0, 0, err
		}
		created += slots - have.C
	}

	return created, ensured, nil
}

func (s *Store) DeleteApprovalTasks(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	snapDate *timestamp.Timestamp,
	deleteAllPending bool,
	taskIDs []int64,
	pendingCount int32,
) error {
	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	condition := mysql.AND(
		tApprovalTasks.DocumentID.EQ(mysql.Int64(documentID)),
		tApprovalTasks.SnapshotDate.EQ(dbutils.TimestampToMySQLDateTimeSec(snapDate)),
	)

	var deleteLimit int64
	if deleteAllPending {
		condition = condition.AND(
			tApprovalTasks.Status.EQ(
				mysql.Int32(
					int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
				),
			),
		)
		deleteLimit = int64(max(1, int(pendingCount)))
	} else if len(taskIDs) > 0 {
		ids := make([]mysql.Expression, 0, len(taskIDs))
		for _, id := range taskIDs {
			ids = append(ids, mysql.Int64(id))
		}
		condition = condition.AND(tApprovalTasks.ID.IN(ids...))
		deleteLimit = int64(len(ids))
	} else {
		return nil
	}

	_, err := tApprovalTasks.
		DELETE().
		WHERE(condition).
		LIMIT(deleteLimit).
		ExecContext(ctx, tx)
	return err
}

func (s *Store) ExpireApprovalTasks(ctx context.Context, tx qrm.DB) (int64, error) {
	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	now := time.Now()

	stmt := tApprovalTasks.
		UPDATE().
		SET(
			tApprovalTasks.Status.SET(mysql.Int32(int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_EXPIRED))),
		).
		WHERE(mysql.AND(
			tApprovalTasks.DueAt.LT_EQ(mysql.DateTimeT(now)),
			tApprovalTasks.Status.EQ(
				mysql.Int32(
					int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
				),
			),
		)).
		LIMIT(250)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
