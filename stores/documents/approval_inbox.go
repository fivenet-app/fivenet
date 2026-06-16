package documentsstore

import (
	"context"
	"errors"

	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListApprovalTasksInbox(
	ctx context.Context,
	q ListApprovalTasksInboxQuery,
) (resourcesdatabase.DataCount, []*documentsapproval.ApprovalTask, error) {
	if q.Pagination == nil {
		q.Pagination = &resourcesdatabase.PaginationRequest{}
	}
	if q.UserInfo == nil {
		q.UserInfo = &userinfo.UserInfo{}
	}

	tApprovalTasks := table.FivenetDocumentsApprovalTasks.AS("approval_task")
	tApprovals := table.FivenetDocumentsApprovals
	tDocumentShort := table.FivenetDocuments.AS("document_short")
	tDMeta := table.FivenetDocumentsMeta.AS("meta")

	visibleQuery := s.subjectAccess.VisibleIDsByConditionQuery(
		q.UserInfo,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		table.FivenetDocuments.DeletedAt.IS_NULL(),
	)
	visibleIDs := mysql.
		SELECT(mysql.IntegerColumn("id").From(visibleQuery.Table)).
		FROM(visibleQuery.Table)

	eligible := mysql.OR(
		mysql.AND(
			tApprovalTasks.AssigneeKind.EQ(
				mysql.Int32(
					int32(documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER),
				),
			),
			tApprovalTasks.UserID.EQ(mysql.Int32(q.UserInfo.GetUserId())),
		),
		mysql.AND(
			tApprovalTasks.AssigneeKind.EQ(
				mysql.Int32(
					int32(documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE),
				),
			),
			tApprovalTasks.Job.EQ(mysql.String(q.UserInfo.GetJob())),
			tApprovalTasks.MinimumGrade.LT_EQ(mysql.Int32(q.UserInfo.GetJobGrade())),
		),
	)

	var notAlreadyActed mysql.BoolExpression
	if q.NotAlreadyActed {
		notAlreadyActed = mysql.NOT(
			mysql.EXISTS(
				mysql.SELECT(mysql.Int(1)).
					FROM(tApprovals).
					WHERE(mysql.AND(
						tApprovals.DocumentID.EQ(tApprovalTasks.DocumentID),
						tApprovals.SnapshotDate.EQ(tApprovalTasks.SnapshotDate),
						tApprovals.UserID.EQ(mysql.Int32(q.UserInfo.GetUserId())),
						tApprovals.Status.IN(
							mysql.Int32(
								int32(documentsapproval.ApprovalStatus_APPROVAL_STATUS_APPROVED),
							),
							mysql.Int32(
								int32(documentsapproval.ApprovalStatus_APPROVAL_STATUS_DECLINED),
							),
						),
					)),
			),
		)
	} else {
		notAlreadyActed = mysql.Bool(true)
	}

	t2 := tApprovalTasks.AS("t2")
	statuses := make([]mysql.Expression, 0, len(q.Statuses))
	if len(q.Statuses) > 0 {
		for _, st := range q.Statuses {
			statuses = append(statuses, mysql.Int32(int32(st)))
		}
	} else {
		statuses = append(
			statuses,
			mysql.Int32(int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)),
		)
	}

	maxSlotThisGroup := t2.
		SELECT(mysql.MAX(t2.SlotNo)).
		FROM(t2).
		WHERE(mysql.AND(
			t2.DocumentID.EQ(tApprovalTasks.DocumentID),
			t2.SnapshotDate.EQ(tApprovalTasks.SnapshotDate),
			t2.AssigneeKind.EQ(tApprovalTasks.AssigneeKind),
			mysql.IntExp(mysql.COALESCE(t2.UserID, mysql.Int32(0))).
				EQ(mysql.IntExp(mysql.COALESCE(tApprovalTasks.UserID, mysql.Int32(0)))),
			mysql.StringExp(mysql.COALESCE(t2.Job, mysql.String(""))).
				EQ(mysql.StringExp(mysql.COALESCE(tApprovalTasks.Job, mysql.String("")))),
			mysql.IntExp(mysql.COALESCE(t2.MinimumGrade, mysql.Int32(-1))).
				EQ(mysql.IntExp(mysql.COALESCE(tApprovalTasks.MinimumGrade, mysql.Int32(-1)))),
			t2.Status.IN(statuses...),
		)).
		LIMIT(1)

	onlyFirstSlot := mysql.OR(
		tApprovalTasks.AssigneeKind.EQ(
			mysql.Int32(int32(documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER)),
		),
		maxSlotThisGroup.IN(tApprovalTasks.SlotNo),
	)

	condition := mysql.AND(
		tDocumentShort.DeletedAt.IS_NULL(),
		tApprovalTasks.DocumentID.IN(visibleIDs),
		eligible,
		notAlreadyActed,
		onlyFirstSlot,
		tApprovalTasks.Status.IN(statuses...),
	)

	if q.OnlyDrafts != nil {
		condition = condition.AND(tDocumentShort.Draft.EQ(mysql.Bool(*q.OnlyDrafts)))
	}

	var countStmt mysql.Statement = tApprovalTasks.
		SELECT(mysql.COUNT(tApprovalTasks.ID).AS("data_count.total")).
		FROM(
			tApprovalTasks.
				INNER_JOIN(tDocumentShort, tDocumentShort.ID.EQ(tApprovalTasks.DocumentID)),
		).
		WHERE(condition)
	if len(visibleQuery.CTEs) > 0 {
		countStmt = mysql.WITH(visibleQuery.CTEs...)(countStmt)
	}

	var count resourcesdatabase.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	if count.Total <= 0 {
		return count, []*documentsapproval.ApprovalTask{}, nil
	}

	tUser := table.FivenetUser.AS("requester")
	tCreator := table.FivenetUser.AS("creator")
	var stmt mysql.Statement = mysql.
		SELECT(
			tApprovalTasks.ID,
			tApprovalTasks.DocumentID,
			tApprovalTasks.SnapshotDate,
			tApprovalTasks.AssigneeKind,
			tApprovalTasks.UserID,
			tApprovalTasks.Job,
			tApprovalTasks.MinimumGrade,
			tApprovalTasks.Label,
			tApprovalTasks.SignatureRequired,
			tApprovalTasks.SlotNo,
			tApprovalTasks.Status,
			tApprovalTasks.Comment,
			tApprovalTasks.CreatedAt,
			tApprovalTasks.CompletedAt,
			tApprovalTasks.DueAt,
			tApprovalTasks.DecisionCount,
			tApprovalTasks.CreatorID,
			tApprovalTasks.CreatorJob,
			tUser.ID,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.Job,
			tUser.JobGrade,
			tDocumentShort.Title,
			tDocumentShort.ContentType,
			tDocumentShort.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tDocumentShort.CreatorJob,
			tDocumentShort.State.AS("meta.state"),
			tDocumentShort.Closed.AS("meta.closed"),
			tDocumentShort.Draft.AS("meta.draft"),
			tDocumentShort.Public.AS("meta.public"),
			tDMeta.DocumentID,
			tDMeta.Approved,
			tDMeta.ApRequiredTotal,
			tDMeta.ApCollectedApproved,
			tDMeta.ApRequiredRemaining,
			tDMeta.ApDeclinedCount,
			tDMeta.ApPendingCount,
			tDMeta.ApAnyDeclined,
			tDMeta.ApPoliciesActive,
			tDMeta.CommentCount,
		).
		FROM(
			tApprovalTasks.
				INNER_JOIN(tDocumentShort, tDocumentShort.ID.EQ(tApprovalTasks.DocumentID)).
				LEFT_JOIN(tUser, tUser.ID.EQ(tApprovalTasks.UserID)).
				LEFT_JOIN(tCreator, tCreator.ID.EQ(tApprovalTasks.CreatorID)).
				LEFT_JOIN(tDMeta, tDMeta.DocumentID.EQ(tDocumentShort.ID)),
		).
		WHERE(condition).
		OFFSET(q.Pagination.GetOffset()).
		ORDER_BY(tApprovalTasks.CreatedAt.ASC()).
		LIMIT(20)
	if len(visibleQuery.CTEs) > 0 {
		stmt = mysql.WITH(visibleQuery.CTEs...)(stmt)
	}

	var tasks []*documentsapproval.ApprovalTask
	if err := stmt.QueryContext(ctx, s.db, &tasks); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	return count, tasks, nil
}
