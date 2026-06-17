package documentsstore

import (
	"context"
	"errors"

	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListApprovalTasksInbox(
	ctx context.Context,
	q ListApprovalTasksInboxQuery,
) (resourcesdatabase.DataCount, []*documentsapproval.ApprovalTask, error) {
	tApprovalTasks := table.FivenetDocumentsApprovalTasks.AS("approval_task")
	tApprovals := table.FivenetDocumentsApprovals
	tDocumentShort := table.FivenetDocuments.AS("document_short")
	tDMeta := table.FivenetDocumentsMeta.AS("meta")
	tUser := table.FivenetUser.AS("requester")
	tCreator := table.FivenetUser.AS("creator")
	pagination := q.Pagination
	if pagination == nil {
		pagination = &resourcesdatabase.PaginationRequest{}
	}

	documentCondition := buildApprovalInboxDocumentCondition(tDocumentShort, q)

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
				mysql.
					SELECT(mysql.Int(1)).
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

	taskCondition := mysql.AND(
		eligible,
		notAlreadyActed,
		onlyFirstSlot,
		tApprovalTasks.Status.IN(statuses...),
	)

	var count resourcesdatabase.DataCount
	var countStmt mysql.Statement
	var stmt mysql.Statement
	var ctes []mysql.CommonTableExpression

	if q.UserInfo.GetSuperuser() {
		baseFrom := tApprovalTasks.
			INNER_JOIN(tDocumentShort, tDocumentShort.ID.EQ(tApprovalTasks.DocumentID))
		countStmt = tApprovalTasks.
			SELECT(
				mysql.COUNT(tApprovalTasks.ID).AS("data_count.total"),
			).
			FROM(baseFrom).
			WHERE(mysql.AND(taskCondition, documentCondition))

		stmt = mysql.
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
				baseFrom.
					LEFT_JOIN(tUser, tUser.ID.EQ(tApprovalTasks.UserID)).
					LEFT_JOIN(tCreator, tCreator.ID.EQ(tApprovalTasks.CreatorID)).
					LEFT_JOIN(tDMeta, tDMeta.DocumentID.EQ(tDocumentShort.ID)),
			).
			WHERE(mysql.AND(taskCondition, documentCondition)).
			OFFSET(pagination.GetOffset()).
			ORDER_BY(tApprovalTasks.CreatedAt.ASC()).
			LIMIT(20)
	} else {
		visibleIDs := s.subjectAccess.VisibleIDsByConditionQuery(
			q.UserInfo,
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			false,
			documentCondition,
		)
		ctes = visibleIDs.CTEs
		visibleDocID := mysql.IntegerColumn("id").From(visibleIDs.Table)
		baseFrom := visibleIDs.Table.
			INNER_JOIN(tDocumentShort, tDocumentShort.ID.EQ(visibleDocID)).
			INNER_JOIN(tApprovalTasks, tApprovalTasks.DocumentID.EQ(tDocumentShort.ID))

		countStmt = tApprovalTasks.
			SELECT(
				mysql.COUNT(tApprovalTasks.ID).AS("data_count.total"),
			).
			FROM(baseFrom).
			WHERE(taskCondition)

		stmt = mysql.
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
				baseFrom.
					LEFT_JOIN(tUser, tUser.ID.EQ(tApprovalTasks.UserID)).
					LEFT_JOIN(tCreator, tCreator.ID.EQ(tApprovalTasks.CreatorID)).
					LEFT_JOIN(tDMeta, tDMeta.DocumentID.EQ(tDocumentShort.ID)),
			).
			WHERE(taskCondition).
			OFFSET(pagination.GetOffset()).
			ORDER_BY(tApprovalTasks.CreatedAt.ASC()).
			LIMIT(20)
	}

	if len(ctes) > 0 {
		countStmt = mysql.WITH(ctes...)(countStmt)
		stmt = mysql.WITH(ctes...)(stmt)
	}

	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	if count.Total <= 0 {
		return count, []*documentsapproval.ApprovalTask{}, nil
	}

	var tasks []*documentsapproval.ApprovalTask
	if err := stmt.QueryContext(ctx, s.db, &tasks); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	return count, tasks, nil
}

func buildApprovalInboxDocumentCondition(
	document *table.FivenetDocumentsTable,
	q ListApprovalTasksInboxQuery,
) mysql.BoolExpression {
	condition := mysql.Bool(true)
	if q.OnlyDrafts != nil {
		condition = condition.AND(document.Draft.EQ(mysql.Bool(*q.OnlyDrafts)))
	}

	return condition
}
