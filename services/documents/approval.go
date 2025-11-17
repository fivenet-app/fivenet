package documents

import (
	"context"
	"errors"
	"fmt"
	"time"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// ListApprovalTasksInbox.
func (s *Server) ListApprovalTasksInbox(
	ctx context.Context,
	req *pbdocuments.ListApprovalTasksInboxRequest,
) (*pbdocuments.ListApprovalTasksInboxResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tApprovalTasks := table.FivenetDocumentsApprovalTasks.AS("approval_task")
	tApprovals := table.FivenetDocumentsApprovals

	var existsAccess mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		existsAccess = mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tDAccess).
				WHERE(mysql.AND(
					tDAccess.TargetID.EQ(tApprovalTasks.DocumentID),
					mysql.OR(
						// Direct user access
						tDAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						// or job + grade access
						mysql.AND(
							tDAccess.Job.EQ(mysql.String(userInfo.GetJob())),
							tDAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
						),
					),
					tDAccess.Access.GT_EQ(
						mysql.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)),
					),
				)),
		)
	} else {
		existsAccess = mysql.Bool(true)
	}

	// Eligibility for this task
	eligible := mysql.OR(
		// USER slot
		mysql.AND(
			tApprovalTasks.AssigneeKind.EQ(
				mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER)),
			),
			tApprovalTasks.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
		),
		// JOB slot
		mysql.AND(
			tApprovalTasks.AssigneeKind.EQ(
				mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE)),
			),
			tApprovalTasks.Job.EQ(mysql.String(userInfo.GetJob())),
			tApprovalTasks.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
		),
	)

	// NOT already approved/declined in this round
	notAlreadyActed := mysql.NOT(
		mysql.EXISTS(
			mysql.SELECT(mysql.Int(1)).
				FROM(tApprovals).
				WHERE(mysql.AND(
					tApprovals.DocumentID.EQ(tApprovalTasks.DocumentID),
					tApprovals.SnapshotDate.EQ(tApprovalTasks.SnapshotDate),
					tApprovals.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tApprovals.Status.IN(
						mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_APPROVED)),
						mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_DECLINED)),
					),
				)),
		),
	)

	// For JOB groups: only the smallest slot_no
	t2 := tApprovalTasks.AS("t2")

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
			t2.Status.EQ(
				mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)),
			),
		)).
		LIMIT(1)

	// onlyFirstSlot is true if:
	//  - USER task (there’s only one anyway), or
	//  - JOB task with slot_no = MIN(slot_no) among pending in its group
	onlyFirstSlot := mysql.OR(
		tApprovalTasks.AssigneeKind.EQ(
			mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER)),
		),
		maxSlotThisGroup.IN(tApprovalTasks.SlotNo),
	)

	condition := mysql.AND(
		existsAccess,
		tDocumentShort.DeletedAt.IS_NULL(),
		eligible,
		notAlreadyActed,
		onlyFirstSlot,
	)
	if len(req.GetStatuses()) > 0 {
		vals := make([]mysql.Expression, 0, len(req.GetStatuses()))
		for _, st := range req.GetStatuses() {
			vals = append(vals, mysql.Int32(int32(st)))
		}
		condition = condition.AND(tApprovalTasks.Status.IN(vals...))
	} else {
		condition = condition.AND(
			tApprovalTasks.Status.EQ(
				mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)),
			),
		)
	}

	if req.OnlyDrafts != nil {
		condition = condition.AND(tDocumentShort.Draft.EQ(mysql.Bool(req.GetOnlyDrafts())))
	}

	// Get total count of values
	countStmt := tApprovalTasks.
		SELECT(
			mysql.COUNT(tApprovalTasks.ID).AS("data_count.total"),
		).
		FROM(
			tApprovalTasks.
				INNER_JOIN(tDocumentShort,
					tDocumentShort.ID.EQ(tApprovalTasks.DocumentID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	resp := &pbdocuments.ListApprovalTasksInboxResponse{
		Tasks: []*documents.ApprovalTask{},
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 20)
	resp.Pagination = pag
	if count.Total <= 0 {
		return resp, nil
	}

	tUser := tables.User().AS("requester")
	tCreator := tables.User().AS("creator")

	stmt := mysql.
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
		).
		FROM(
			tApprovalTasks.
				INNER_JOIN(tDocumentShort,
					tDocumentShort.ID.EQ(tApprovalTasks.DocumentID),
				).
				LEFT_JOIN(tUser,
					tUser.ID.EQ(tApprovalTasks.UserID),
				).
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tApprovalTasks.CreatorID),
				).
				LEFT_JOIN(tDMeta,
					tDMeta.DocumentID.EQ(tDocumentShort.ID),
				),
		).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(tApprovalTasks.CreatedAt.ASC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Tasks); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	tasks := resp.GetTasks()
	for _, t := range tasks {
		if t.GetJob() != "" {
			jobInfoFn(t)
		}

		if t.GetCreator() != nil {
			jobInfoFn(t.GetCreator())
		}

		if t.GetDocument() != nil {
			doc := t.GetDocument()
			if doc.GetCreator() != nil {
				jobInfoFn(doc.GetCreator())
			}

			if job := s.enricher.GetJobByName(doc.GetCreatorJob()); job != nil {
				doc.CreatorJobLabel = &job.Label
			}
		}
	}

	return resp, nil
}

// ListApprovalPolicies.
func (s *Server) ListApprovalPolicies(
	ctx context.Context,
	req *pbdocuments.ListApprovalPoliciesRequest,
) (*pbdocuments.ListApprovalPoliciesResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	condition := tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))
	if !userInfo.GetSuperuser() {
		condition = condition.AND(tApprovalPolicy.DeletedAt.IS_NULL())
	}

	policy, err := s.getApprovalPolicy(ctx, s.db, condition)
	if err != nil {
		return nil, err
	}

	return &pbdocuments.ListApprovalPoliciesResponse{
		Policy: policy,
	}, nil
}

func (s *Server) getApprovalPolicy(
	ctx context.Context, tx qrm.DB, condition mysql.BoolExpression,
) (*documents.ApprovalPolicy, error) {
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	stmt := tApprovalPolicy.
		SELECT(
			tApprovalPolicy.DocumentID,
			tApprovalPolicy.RuleKind,
			tApprovalPolicy.RequiredCount,
			tApprovalPolicy.SignatureRequired,
			tApprovalPolicy.SnapshotDate,
			tApprovalPolicy.StartedAt,
			tApprovalPolicy.CompletedAt,
			tApprovalPolicy.OnEditBehavior,
			tApprovalPolicy.SelfApproveAllowed,
			tApprovalPolicy.AssignedCount,
			tApprovalPolicy.ApprovedCount,
			tApprovalPolicy.DeclinedCount,
			tApprovalPolicy.PendingCount,
			tApprovalPolicy.AnyDeclined,
			tApprovalPolicy.StartedAt,
			tApprovalPolicy.CompletedAt,
			tApprovalPolicy.CreatedAt,
			tApprovalPolicy.UpdatedAt,
			tApprovalPolicy.DeletedAt,
		).
		FROM(tApprovalPolicy).
		WHERE(condition).
		LIMIT(1)

	pol := &documents.ApprovalPolicy{}
	if err := stmt.QueryContext(ctx, tx, pol); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if pol.GetDocumentId() == 0 {
		return nil, nil
	}

	return pol, nil
}

func (s *Server) getOrCreateApprovalPolicy(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	snapshotDate *timestamp.Timestamp,
) (*documents.ApprovalPolicy, error) {
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies

	condition := tApprovalPolicy.AS("approval_policy").DocumentID.EQ(mysql.Int64(documentId))

	pol, err := s.getApprovalPolicy(ctx, tx, condition)
	if err != nil {
		return nil, err
	}
	if pol != nil {
		return pol, nil
	}

	requiredCount := int32(1)
	if err = s.createApprovalPolicy(ctx, tx, documentId, &documents.ApprovalPolicy{
		SnapshotDate:       snapshotDate,
		RuleKind:           documents.ApprovalRuleKind_APPROVAL_RULE_KIND_REQUIRE_ALL,
		RequiredCount:      &requiredCount,
		OnEditBehavior:     documents.OnEditBehavior_ON_EDIT_BEHAVIOR_KEEP_PROGRESS,
		SignatureRequired:  false,
		SelfApproveAllowed: false,
	}); err != nil {
		return nil, err
	}

	pol, err = s.getApprovalPolicy(ctx, tx, condition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol == nil {
		return nil, errorsdocuments.ErrFailedQuery
	}

	return pol, nil
}

func (s *Server) createApprovalPolicy(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	pol *documents.ApprovalPolicy,
) error {
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies
	// Create a new policy if it doesn't exist
	if _, err := tApprovalPolicy.
		INSERT(
			tApprovalPolicy.DocumentID,
			tApprovalPolicy.SnapshotDate,
			tApprovalPolicy.RuleKind,
			tApprovalPolicy.RequiredCount,
			tApprovalPolicy.OnEditBehavior,
			tApprovalPolicy.SignatureRequired,
			tApprovalPolicy.SelfApproveAllowed,
		).
		VALUES(
			documentId,
			pol.GetSnapshotDate(),
			int32(pol.GetRuleKind()),
			pol.GetRequiredCount(),
			pol.GetOnEditBehavior(),
			pol.GetSignatureRequired(),
			pol.GetSelfApproveAllowed(),
		).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

// UpsertApprovalPolicy.
func (s *Server) UpsertApprovalPolicy(
	ctx context.Context,
	req *pbdocuments.UpsertApprovalPolicyRequest,
) (*pbdocuments.UpsertApprovalPolicyResponse, error) {
	pol := req.GetPolicy()

	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", pol.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		pol.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_STATUS,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	stmt := tApprovalPolicy.
		INSERT(
			tApprovalPolicy.DocumentID,
			tApprovalPolicy.SnapshotDate,
			tApprovalPolicy.OnEditBehavior,
			tApprovalPolicy.RuleKind,
			tApprovalPolicy.RequiredCount,
			tApprovalPolicy.SignatureRequired,
			tApprovalPolicy.SelfApproveAllowed,
		).
		VALUES(
			pol.GetDocumentId(),
			mysql.CURRENT_TIMESTAMP(), // Initialize snapshot_date
			int32(pol.GetOnEditBehavior()),
			int32(pol.GetRuleKind()),
			pol.GetRequiredCount(),
			pol.GetSignatureRequired(),
			pol.GetSelfApproveAllowed(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tApprovalPolicy.OnEditBehavior.SET(mysql.Int32(int32(pol.GetOnEditBehavior()))),
			tApprovalPolicy.RuleKind.SET(mysql.Int32(int32(pol.GetRuleKind()))),
			tApprovalPolicy.RequiredCount.SET(mysql.Int32(pol.GetRequiredCount())),
			tApprovalPolicy.SignatureRequired.SET(mysql.Bool(pol.GetSignatureRequired())),
			tApprovalPolicy.SelfApproveAllowed.SET(mysql.Bool(pol.GetSelfApproveAllowed())),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	tApprovalPolicy = tApprovalPolicy.AS("approval_policy")

	condition := tApprovalPolicy.DocumentID.EQ(mysql.Int64(pol.GetDocumentId()))
	policy, err := s.getApprovalPolicy(ctx, tx, condition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, policy.GetDocumentId(), policy.GetSnapshotDate().AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	policy, err = s.getApprovalPolicy(ctx, s.db, condition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.UpsertApprovalPolicyResponse{
		Policy: policy,
	}, nil
}

// ListApprovalTasks.
func (s *Server) ListApprovalTasks(
	ctx context.Context,
	req *pbdocuments.ListApprovalTasksRequest,
) (*pbdocuments.ListApprovalTasksResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	tApprovalTasks := table.FivenetDocumentsApprovalTasks.AS("approval_task")

	condition := tApprovalTasks.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))

	if len(req.GetStatuses()) > 0 {
		vals := make([]mysql.Expression, 0, len(req.GetStatuses()))
		for _, st := range req.GetStatuses() {
			vals = append(vals, mysql.Int32(int32(st)))
		}
		condition = condition.AND(tApprovalTasks.Status.IN(vals...))
	}

	resp := &pbdocuments.ListApprovalTasksResponse{
		Tasks: []*documents.ApprovalTask{},
	}

	tUser := tables.User().AS("usershort")

	stmt := mysql.
		SELECT(
			tApprovalTasks.ID,
			tApprovalTasks.DocumentID,
			tApprovalTasks.SnapshotDate,
			tApprovalTasks.AssigneeKind,
			tApprovalTasks.UserID,
			tApprovalTasks.Job,
			tApprovalTasks.MinimumGrade,
			tApprovalTasks.Label,
			tApprovalTasks.SlotNo,
			tApprovalTasks.Status,
			tApprovalTasks.Comment,
			tApprovalTasks.CreatedAt,
			tApprovalTasks.CompletedAt,
			tApprovalTasks.DueAt,
			tApprovalTasks.DecisionCount,
			tApprovalTasks.ApprovalID,
			tApprovalTasks.CreatorID,
			tApprovalTasks.CreatorJob,
			tUser.ID,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.Job,
			tUser.JobGrade,
		).
		FROM(
			tApprovalTasks.
				LEFT_JOIN(tUser,
					tUser.ID.EQ(tApprovalTasks.UserID),
				),
		).
		WHERE(condition).
		ORDER_BY(tApprovalTasks.CreatedAt.ASC()).
		LIMIT(15)

	if err := stmt.QueryContext(ctx, s.db, &resp.Tasks); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Tasks {
		if resp.Tasks[i].GetJob() == "" {
			continue
		}

		jobInfoFn(resp.Tasks[i])
	}

	return resp, nil
}

func (s *Server) getApprovalTask(
	ctx context.Context,
	tx qrm.Queryable,
	taskId int64,
) (*documents.ApprovalTask, error) {
	tApprovalTasks := table.FivenetDocumentsApprovalTasks.AS("approval_task")

	stmt := tApprovalTasks.
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
			tApprovalTasks.ApprovalID,
		).
		FROM(tApprovalTasks).
		WHERE(tApprovalTasks.ID.EQ(mysql.Int64(taskId))).
		LIMIT(1)

	var task documents.ApprovalTask
	if err := stmt.QueryContext(ctx, tx, &task); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if task.Id == 0 {
		return nil, nil
	}

	return &task, nil
}

// UpsertApprovalTasks.
func (s *Server) UpsertApprovalTasks(
	ctx context.Context,
	req *pbdocuments.UpsertApprovalTasksRequest,
) (*pbdocuments.UpsertApprovalTasksResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.document_id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Resolve policy & snapshot
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	var pol documents.ApprovalPolicy
	if err := tApprovalPolicy.
		SELECT(
			tApprovalPolicy.DocumentID, tApprovalPolicy.SnapshotDate,
		).
		FROM(tApprovalPolicy).
		WHERE(tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.GetDocumentId() == 0 {
		return nil, errorsdocuments.ErrNotFoundOrNoPerms
	}
	// Access: must be allowed to edit the document to seed tasks
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		pol.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	snap := pol.GetSnapshotDate().AsTime()
	if req.GetSnapshotDate() != nil {
		snap = req.GetSnapshotDate().AsTime()
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	created, ensured, err := s.createApprovalTasks(
		ctx,
		tx,
		userInfo,
		req.GetDocumentId(),
		snap,
		req.GetSeeds(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.UpsertApprovalTasksResponse{
		TasksCreated: created,
		TasksEnsured: ensured,
		Policy:       &pol,
	}, nil
}

func (s *Server) createApprovalTasks(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
	documentId int64,
	snap time.Time,
	seeds []*pbdocuments.ApprovalTaskSeed,
) (int32, int32, error) {
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	created := int32(0)
	ensured := int32(0)

	for _, seed := range seeds {
		isUser := seed.GetUserId() > 0
		if isUser {
			// Ensure one USER task exists
			var cnt struct{ C int32 }
			if err := tApprovalTasks.
				SELECT(mysql.COUNT(tApprovalTasks.ID).AS("C")).
				FROM(tApprovalTasks).
				WHERE(mysql.AND(
					tApprovalTasks.DocumentID.EQ(mysql.Int64(documentId)),
					tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
					tApprovalTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER))),
					tApprovalTasks.UserID.EQ(mysql.Int32(seed.GetUserId())),
				)).
				LIMIT(1).
				QueryContext(ctx, tx, &cnt); err != nil {
				return 0, 0, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
			if cnt.C > 0 {
				ensured++
				continue
			}

			// Insert USER task with slot_no=1
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
					documentId,
					snap,
					int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER),
					seed.GetUserId(),
					seed.GetLabel(),
					seed.SignatureRequired,
					1,
					int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
					seed.GetComment(),
					dbutils.TimestampToMySQL(seed.GetDueAt()),
					userInfo.GetUserId(),
					userInfo.GetJob(),
				).
				ExecContext(ctx, tx); err != nil {
				return 0, 0, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
			created++
			continue
		}

		// JOB target with N slots
		slots := seed.GetSlots()
		if slots <= 0 {
			slots = 1
		}
		// Count existing PENDING slots for this target
		var have struct{ C int32 }
		if err := mysql.
			SELECT(mysql.COUNT(tApprovalTasks.ID).AS("C")).
			FROM(tApprovalTasks).
			WHERE(mysql.AND(
				tApprovalTasks.DocumentID.EQ(mysql.Int64(documentId)),
				tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
				tApprovalTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE))),
				tApprovalTasks.Job.EQ(mysql.String(seed.GetJob())),
				tApprovalTasks.MinimumGrade.EQ(mysql.Int32(seed.GetMinimumGrade())),
				tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
			)).
			LIMIT(1).
			QueryContext(ctx, tx, &have); err != nil {
			return 0, 0, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if have.C >= int32(slots) {
			ensured++
			continue
		}

		// Insert missing [have.C+1 .. slots] rows
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
				documentId,
				snap,
				int32(
					documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE,
				),
				seed.GetJob(),
				seed.GetMinimumGrade(),
				seed.GetLabel(),
				seed.GetSignatureRequired(),
				slot,
				int32(
					documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING,
				),
				seed.GetComment(),
				dbutils.TimestampToMySQL(seed.GetDueAt()),
				userInfo.GetUserId(),
				userInfo.GetJob(),
			)
		}
		if _, err := ins.ExecContext(ctx, tx); err != nil {
			return 0, 0, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		created += int32(slots) - have.C
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, documentId, snap); err != nil {
		return 0, 0, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Create document activity entry
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   documentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_ASSIGNED,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return 0, 0, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return created, ensured, nil
}

// DeleteApprovalTasks.
func (s *Server) DeleteApprovalTasks(
	ctx context.Context,
	req *pbdocuments.DeleteApprovalTasksRequest,
) (*pbdocuments.DeleteApprovalTasksResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.document_id", req.GetDocumentId()})

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	// Resolve policy & snapshot
	var pol documents.ApprovalPolicy
	if err := tApprovalPolicy.
		SELECT(
			tApprovalPolicy.DocumentID,
			tApprovalPolicy.SnapshotDate,
		).
		FROM(tApprovalPolicy).
		WHERE(tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.GetDocumentId() == 0 {
		return nil, errorsdocuments.ErrNotFoundOrNoPerms
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		pol.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	snap := pol.GetSnapshotDate().AsTime()

	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	condition := mysql.AND(
		tApprovalTasks.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
		tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
	)

	// Delete all pending?
	if req.GetDeleteAllPending() {
		condition = condition.AND(
			tApprovalTasks.Status.EQ(
				mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)),
			),
		)
	} else if len(req.GetTaskIds()) > 0 {
		ids := make([]mysql.Expression, 0, len(req.GetTaskIds()))
		for _, id := range req.GetTaskIds() {
			ids = append(ids, mysql.Int64(id))
		}

		condition = condition.AND(tApprovalTasks.ID.IN(ids...))
	} else {
		return &pbdocuments.DeleteApprovalTasksResponse{}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	if _, err := tApprovalTasks.
		DELETE().
		WHERE(condition).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, pol.DocumentId, pol.SnapshotDate.AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteApprovalTasksResponse{}, nil
}

func (s *Server) canUserAccessApprovalTask(
	ctx context.Context,
	tx qrm.Queryable,
	taskId int64,
	userInfo *userinfo.UserInfo,
) error {
	task, err := s.getApprovalTask(ctx, tx, taskId)
	if task == nil || err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	check, err := s.access.CanUserAccessTarget(
		ctx,
		task.DocumentId,
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check {
		return errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	return nil
}

func (s *Server) ListApprovals(
	ctx context.Context,
	req *pbdocuments.ListApprovalsRequest,
) (*pbdocuments.ListApprovalsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Permission: viewer can list approval artifacts
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	tApprovals := table.FivenetDocumentsApprovals.AS("approval")

	// Build filters
	condition := tApprovals.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))

	if req.GetSnapshotDate() != nil {
		condition = condition.AND(
			tApprovals.SnapshotDate.EQ(mysql.TimestampT(req.GetSnapshotDate().AsTime())),
		)
	}
	if req.GetStatus() > 0 {
		condition = condition.AND(
			tApprovals.Status.EQ(mysql.Int32(int32(req.GetStatus().Number()))),
		)
	}
	if req.GetUserId() > 0 {
		condition = condition.AND(tApprovals.UserID.EQ(mysql.Int32(req.GetUserId())))
	}
	if req.GetTaskId() > 0 {
		condition = condition.AND(tApprovals.TaskID.EQ(mysql.Int64(req.GetTaskId())))
	}

	// Count
	countStmt := mysql.
		SELECT(
			mysql.COUNT(tApprovals.ID).AS("data_count.total"),
		).
		FROM(tApprovals).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	resp := &pbdocuments.ListApprovalsResponse{
		Approvals: []*documents.Approval{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Page fetch
	tUser := tables.User().AS("usershort")
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

			// User info
			tUser.ID,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.Job,
			tUser.JobGrade,
		).
		FROM(
			tApprovals.
				LEFT_JOIN(tUser,
					tUser.ID.EQ(tApprovals.UserID),
				).
				LEFT_JOIN(tStamp,
					tApprovals.ID.EQ(tStamp.ID),
				),
		).
		WHERE(condition).
		ORDER_BY(tApprovals.Status.ASC(), tApprovals.CreatedAt.DESC()).
		LIMIT(15)

	if err := stmt.QueryContext(ctx, s.db, &resp.Approvals); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Enrich labels/grades for display
	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for _, t := range resp.Approvals {
		if t.GetUserJob() != "" {
			jobInfoFn(t)
		}

		if t.GetUser() != nil {
			jobInfoFn(t.GetUser())
		}
	}

	return resp, nil
}

// RevokeApproval.
func (s *Server) RevokeApproval(
	ctx context.Context,
	req *pbdocuments.RevokeApprovalRequest,
) (*pbdocuments.RevokeApprovalResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tApprovals := table.FivenetDocumentsApprovals.AS("approval")
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	// Load the approval artifact
	var apr documents.Approval
	if err := tApprovals.
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
		WHERE(tApprovals.ID.EQ(mysql.Int64(req.GetApprovalId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &apr); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Access: require EDIT on the document
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		apr.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	// Resolve policy using document_id + snapshot_date
	var pol documents.ApprovalPolicy
	if err := tApprovalPolicy.
		SELECT(
			tApprovalPolicy.DocumentID,
			tApprovalPolicy.SnapshotDate,
		).
		FROM(tApprovalPolicy).
		WHERE(mysql.AND(
			tApprovalPolicy.DocumentID.EQ(mysql.Int64(apr.GetDocumentId())),
			tApprovalPolicy.SnapshotDate.EQ(mysql.TimestampT(apr.GetSnapshotDate().AsTime())),
		)).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	tApprovals = table.FivenetDocumentsApprovals
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	// Mark revoked
	if _, err := tApprovals.
		UPDATE().
		SET(
			tApprovals.Status.SET(mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_REVOKED))),
			tApprovals.RevokedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tApprovals.Comment.SET(mysql.String(req.GetComment())),
		).
		WHERE(mysql.AND(
			tApprovals.DocumentID.EQ(mysql.Int64(apr.GetDocumentId())),
			tApprovals.ID.EQ(mysql.Int64(req.GetApprovalId())),
		)).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Reload artifact for response
	if err := tApprovals.
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
		WHERE(tApprovals.ID.EQ(mysql.Int64(req.GetApprovalId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &apr); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if apr.GetTaskId() > 0 {
		// Set PENDING & clear decider snapshot
		if _, err = tApprovalTasks.
			UPDATE().
			SET(
				tApprovalTasks.Status.SET(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
				tApprovalTasks.CompletedAt.SET(mysql.TimestampExp(mysql.NULL)),
				tApprovalTasks.DecisionCount.SET(mysql.Int32(0)),
			).
			WHERE(
				tApprovalTasks.ID.EQ(mysql.Int64(apr.GetTaskId())),
			).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Recompute approvals for this policy+snapshot
	if err := s.recomputeApprovalPolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Create document activity entry
	comment := req.GetComment()
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   pol.GetDocumentId(),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_REVOKED,
		Reason:       &comment,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.RevokeApprovalResponse{Approval: &apr}, nil
}

// DecideApproval supports both task-based and ad-hoc approvals.
//   - If task_id is provided: decide that task (PENDING -> APPROVED/DECLINED), create/UPSERT artifact, link both ways.
//   - If task_id is empty: try to auto-match a pending task by (user_id) or (job,min_grade<=user.grade).
//     If found -> decide that task; else -> create an ad-hoc artifact (no task).
func (s *Server) DecideApproval(
	ctx context.Context,
	req *pbdocuments.DecideApprovalRequest,
) (*pbdocuments.DecideApprovalResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		"fivenet.documents.approval.document_id", req.GetDocumentId(),
		"fivenet.documents.approval.task_id", req.GetTaskId(),
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access: must be able to VIEW the document to decide (tighten if you want)
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	doc, err := s.getDocument(
		ctx,
		tDocument.ID.EQ(mysql.Int64(req.GetDocumentId())),
		userInfo,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if doc == nil || doc.GetMeta() == nil || doc.GetMeta().GetDraft() {
		return nil, errorsdocuments.ErrApprovalDocIsDraft
	}

	docSnap := doc.GetUpdatedAt()
	if docSnap == nil {
		docSnap = timestamp.Now()
	}
	// Resolve policy, doc, snapshot
	pol, err := s.getOrCreateApprovalPolicy(
		ctx,
		s.db,
		req.GetDocumentId(),
		docSnap,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if pol.GetSignatureRequired() && req.GetStampId() == 0 && req.GetPayloadSvg() == "" {
		return nil, errorsdocuments.ErrApprovalSignatureRequired
	}

	snapTime := pol.GetSnapshotDate().AsTime()

	tApprovalTasks := table.FivenetDocumentsApprovalTasks.AS("approval_task")
	tApprovals := table.FivenetDocumentsApprovals

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	now := time.Now()

	var decidedTask *documents.ApprovalTask // may remain nil for ad-hoc
	var taskIDForArtifact int64             // 0 if ad-hoc

	// Path A: task_id provided -> validate and mark decided
	if req.GetTaskId() > 0 {
		// Load the task row
		decidedTask, err = s.getApprovalTask(ctx, tx, req.GetTaskId())
		if err != nil {
			return nil, err
		}
		if decidedTask.DocumentId == 0 {
			return nil, errorsdocuments.ErrNotFoundOrNoPerms
		}

		// Must be same policy/snapshot and pending
		if decidedTask.GetDocumentId() != pol.GetDocumentId() ||
			decidedTask.GetSnapshotDate().AsTime() != snapTime {
			return nil, errorsdocuments.ErrDocAccessViewDenied
		}
		if decidedTask.GetStatus() != documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING {
			return nil, errorsdocuments.ErrApprovalTaskAlreadyHandled
		}

		if decidedTask.GetAssigneeKind() == documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER &&
			decidedTask.GetUserId() != userInfo.GetUserId() {
			return nil, errorsdocuments.ErrDocAccessViewDenied
		} else if decidedTask.GetAssigneeKind() == documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE {
			if decidedTask.GetJob() != userInfo.GetJob() ||
				decidedTask.GetMinimumGrade() >= userInfo.GetJobGrade() {
				return nil, errorsdocuments.ErrDocAccessViewDenied
			}
		} else {
			return nil, errorsdocuments.ErrDocAccessViewDenied
		}

		// Update task
		tApprovalTasks = table.FivenetDocumentsApprovalTasks

		if _, err := tApprovalTasks.
			UPDATE(
				tApprovalTasks.Status,
				tApprovalTasks.DecisionCount,
				tApprovalTasks.CompletedAt,
				tApprovalTasks.Comment,
			).
			SET(
				int32(req.GetNewStatus()),
				tApprovalTasks.DecisionCount.ADD(mysql.Int32(1)),
				now,
				req.GetComment(),
			).
			WHERE(tApprovalTasks.ID.EQ(mysql.Int64(req.GetTaskId()))).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		taskIDForArtifact = req.GetTaskId()
	} else {
		// Path B: no task_id -> try to auto-match a pending task for this user

		// First: exact USER task
		var candidate documents.ApprovalTask
		err := tApprovalTasks.
			SELECT(
				tApprovalTasks.ID, tApprovalTasks.DocumentID, tApprovalTasks.SnapshotDate, tApprovalTasks.AssigneeKind,
				tApprovalTasks.UserID, tApprovalTasks.Job, tApprovalTasks.MinimumGrade, tApprovalTasks.Label, tApprovalTasks.SignatureRequired, tApprovalTasks.SlotNo,
				tApprovalTasks.Status, tApprovalTasks.Comment, tApprovalTasks.CreatedAt, tApprovalTasks.CompletedAt, tApprovalTasks.DueAt,
				tApprovalTasks.DecisionCount, tApprovalTasks.CreatorID, tApprovalTasks.CreatorJob, tApprovalTasks.ApprovalID,
			).
			FROM(tApprovalTasks).
			WHERE(mysql.AND(
				tApprovalTasks.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
				tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snapTime)),
				tApprovalTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER))),
				tApprovalTasks.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				tApprovalTasks.Status.IN(
					mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)),
					mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED)),
					mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_CANCELLED)),
				),
			)).
			ORDER_BY(tApprovalTasks.SlotNo.ASC(), tApprovalTasks.CreatedAt.ASC()).
			LIMIT(1).
			QueryContext(ctx, tx, &candidate)

		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		useCandidate := (err == nil && candidate.Id > 0)

		// If no USER task, try JOB target where user is eligible
		if !useCandidate {
			err = tApprovalTasks.
				SELECT(
					tApprovalTasks.ID, tApprovalTasks.DocumentID, tApprovalTasks.SnapshotDate, tApprovalTasks.AssigneeKind,
					tApprovalTasks.UserID, tApprovalTasks.Job, tApprovalTasks.MinimumGrade, tApprovalTasks.Label, tApprovalTasks.SignatureRequired, tApprovalTasks.SlotNo,
					tApprovalTasks.Status, tApprovalTasks.Comment, tApprovalTasks.CreatedAt, tApprovalTasks.CompletedAt, tApprovalTasks.DueAt,
					tApprovalTasks.DecisionCount, tApprovalTasks.CreatorID, tApprovalTasks.CreatorJob, tApprovalTasks.ApprovalID,
				).
				FROM(tApprovalTasks).
				WHERE(mysql.AND(
					tApprovalTasks.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
					tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snapTime)),
					tApprovalTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE))),
					tApprovalTasks.Job.EQ(mysql.String(userInfo.GetJob())),
					tApprovalTasks.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					tApprovalTasks.Status.IN(
						mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)),
						mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED)),
						mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_CANCELLED)),
					),
				)).
				ORDER_BY(tApprovalTasks.SlotNo.ASC(), tApprovalTasks.CreatedAt.ASC()).
				LIMIT(1).
				QueryContext(ctx, tx, &candidate)

			if err != nil && !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

			useCandidate = (err == nil && candidate.Id > 0)
		}

		// If a task matched, decide it; otherwise proceed as true ad-hoc (no task)
		if useCandidate {
			tApprovalTasks = table.FivenetDocumentsApprovalTasks

			if _, err := tApprovalTasks.
				UPDATE().
				SET(
					tApprovalTasks.Status.SET(mysql.Int32(int32(req.GetNewStatus()))),
					tApprovalTasks.DecisionCount.SET(tApprovalTasks.DecisionCount.ADD(mysql.Int32(1))),
					tApprovalTasks.CompletedAt.SET(mysql.TimestampT(now)),
					tApprovalTasks.Comment.SET(mysql.String(req.GetComment())),
				).
				WHERE(tApprovalTasks.ID.EQ(mysql.Int64(candidate.GetId()))).
				LIMIT(1).
				ExecContext(ctx, tx); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

			decidedTask = &candidate
			taskIDForArtifact = candidate.GetId()
		}
	}

	if decidedTask.GetSignatureRequired() && req.GetStampId() == 0 && req.GetPayloadSvg() == "" {
		return nil, errorsdocuments.ErrApprovalSignatureRequired
	}

	var existing struct {
		ID int64 `alias:"id"`
	}
	if err := tApprovals.
		SELECT(
			tApprovals.ID.AS("id"),
		).
		FROM(tApprovals).
		WHERE(mysql.AND(
			tApprovals.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
			tApprovals.SnapshotDate.EQ(mysql.TimestampT(snapTime)),
			tApprovals.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
		)).
		LIMIT(1).
		QueryContext(ctx, tx, &existing); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Write/UPSERT approval artifact (unique per (document_id, snapshot_date, user_id))
	// Map task statuses APPROVED/DECLINED -> artifact status
	artifactStatus := documents.ApprovalStatus_APPROVAL_STATUS_APPROVED
	if req.GetNewStatus() == documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED ||
		req.GetNewStatus() == documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_CANCELLED {
		artifactStatus = documents.ApprovalStatus_APPROVAL_STATUS_DECLINED
	}

	// Insert or update the approval artifact for this user
	if existing.ID <= 0 {
		if _, err := tApprovals.
			INSERT(
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
			).
			VALUES(
				pol.GetDocumentId(),
				dbutils.TimestampToMySQL(pol.GetSnapshotDate()),
				userInfo.GetUserId(),
				userInfo.GetJob(),
				mysql.Int32(userInfo.GetJobGrade()),
				req.GetPayloadSvg(),
				dbutils.Int64P(req.GetStampId()),
				int32(artifactStatus),
				req.GetComment(),
				dbutils.Int64P(taskIDForArtifact),
			).
			ON_DUPLICATE_KEY_UPDATE(
				tApprovals.SnapshotDate.SET(mysql.DateTimeExp(mysql.RawDate("VALUES(`snapshot_date`)"))),
				tApprovals.UserID.SET(mysql.Int32(userInfo.GetUserId())),
				tApprovals.UserJob.SET(mysql.String(userInfo.GetJob())),
				tApprovals.UserJobGrade.SET(mysql.Int32(userInfo.GetJobGrade())),
				tApprovals.PayloadSvg.SET(mysql.String(req.GetPayloadSvg())),
				tApprovals.StampID.SET(dbutils.Int64P(req.GetStampId())),
				tApprovals.Status.SET(mysql.Int32(int32(artifactStatus))),
				tApprovals.Comment.SET(mysql.String(req.GetComment())),
				tApprovals.TaskID.SET(dbutils.Int64P(taskIDForArtifact)),
			).
			ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		if _, err := tApprovals.
			UPDATE(
				tApprovals.SnapshotDate,
				tApprovals.UserID,
				tApprovals.UserJob,
				tApprovals.UserJobGrade,
				tApprovals.PayloadSvg,
				tApprovals.StampID,
				tApprovals.Status,
				tApprovals.Comment,
				tApprovals.TaskID,
			).
			SET(
				snapTime,
				userInfo.GetUserId(),
				userInfo.GetJob(),
				mysql.Int32(userInfo.GetJobGrade()),
				req.GetPayloadSvg(),
				dbutils.Int64P(req.GetStampId()),
				artifactStatus,
				req.GetComment(),
				dbutils.Int64P(taskIDForArtifact),
			).
			WHERE(tApprovals.ID.EQ(mysql.Int64(existing.ID))).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	tApprovals = table.FivenetDocumentsApprovals.AS("approval")

	var artifact documents.Approval
	if err := tApprovals.
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
		WHERE(mysql.AND(
			tApprovals.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
			tApprovals.SnapshotDate.EQ(mysql.TimestampT(snapTime)),
			tApprovals.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
		)).
		LIMIT(1).
		QueryContext(ctx, tx, &artifact); err != nil {
		return nil, errswrap.NewError(
			fmt.Errorf("failed to get approval artifact. %w", err),
			errorsdocuments.ErrFailedQuery,
		)
	}

	// If we decided a task, set its approval_id backlink
	if decidedTask != nil && artifact.Id > 0 {
		if _, err := tApprovalTasks.
			UPDATE(tApprovalTasks.ApprovalID).
			SET(artifact.GetId()).
			WHERE(tApprovalTasks.ID.EQ(mysql.Int64(decidedTask.GetId()))).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		decidedTask.ApprovalId = &artifact.Id
		// reflect status in the in-memory task too
		decidedTask.Status = req.GetNewStatus()
		decidedTask.CompletedAt = timestamp.New(now)
		comment := req.GetComment()
		decidedTask.Comment = &comment
	}

	// Recompute rollups for tasks, policy and document
	if err := s.recomputeApprovalPolicyTx(ctx, tx, pol.GetDocumentId(), snapTime); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Create document activity entry
	comment := req.GetComment()
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	activityType := documents.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_REJECTED
	if req.GetNewStatus() == documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_APPROVED {
		activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_APPROVED
	}

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   pol.GetDocumentId(),
		ActivityType: activityType,
		Reason:       &comment,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DecideApprovalResponse{
		Approval: &artifact,
		Task:     decidedTask,
		Policy:   pol,
	}, nil
}

// ReopenApprovalTask.
func (s *Server) ReopenApprovalTask(
	ctx context.Context,
	req *pbdocuments.ReopenApprovalTaskRequest,
) (*pbdocuments.ReopenApprovalTaskResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.task_id", req.GetTaskId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.canUserAccessApprovalTask(ctx, s.db, req.GetTaskId(), userInfo); err != nil {
		return nil, err
	}

	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	// Set PENDING & clear decider snapshot
	if _, err = tApprovalTasks.
		UPDATE().
		SET(
			tApprovalTasks.Status.SET(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
			tApprovalTasks.CompletedAt.SET(mysql.TimestampExp(mysql.NULL)),
		).
		WHERE(
			tApprovalTasks.ID.EQ(mysql.Int64(req.GetTaskId())),
		).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Fetch key
	var k struct {
		DocumentID   int64
		SnapshotDate timestamp.Timestamp
	}
	if err = mysql.
		SELECT(
			tApprovalTasks.DocumentID,
			tApprovalTasks.SnapshotDate,
		).
		FROM(tApprovalTasks).
		WHERE(tApprovalTasks.ID.EQ(mysql.Int(req.GetTaskId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &k); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, k.DocumentID, k.SnapshotDate.AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Create document activity entry
	comment := req.GetComment()
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   k.DocumentID,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_REVOKED,
		Reason:       &comment,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err = tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.ReopenApprovalTaskResponse{
		Task: &documents.ApprovalTask{
			Id:     req.GetTaskId(),
			Status: documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING,
		},
		Policy: &documents.ApprovalPolicy{
			DocumentId: k.DocumentID,
		},
	}, nil
}

// RecomputeApprovalPolicyCounters (aggregate from tasks).
func (s *Server) RecomputeApprovalPolicyCounters(
	ctx context.Context,
	req *pbdocuments.RecomputeApprovalPolicyCountersRequest,
) (*pbdocuments.RecomputeApprovalPolicyCountersResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	condition := tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))

	pol, err := s.getApprovalPolicy(ctx, s.db, condition)
	if err != nil {
		return nil, err
	}

	if err := s.recomputeApprovalPolicyTx(ctx, s.db, req.GetDocumentId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pol, err = s.getApprovalPolicy(ctx, s.db, condition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.RecomputeApprovalPolicyCountersResponse{
		Policy: pol,
	}, nil
}

// recomputeApprovalPolicyTx recalculates approval counters for a policy+snapshot
// and updates fivenet_documents_meta for list-time flags.
func (s *Server) recomputeApprovalPolicyTx(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	snap time.Time,
) error {
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")
	tApprovals := table.FivenetDocumentsApprovals
	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	tDocumentsMeta := table.FivenetDocumentsMeta

	// Load policy if exists
	pol, err := s.getApprovalPolicy(
		ctx,
		tx,
		tApprovalPolicy.AS("approval_policy").DocumentID.EQ(mysql.Int64(documentID)),
	)
	if err != nil {
		return err
	}
	if pol == nil {
		pol = &documents.ApprovalPolicy{}
	}
	pol.Default()

	var docCreatorId int32
	if !pol.SelfApproveAllowed {
		// Get document creator ID
		var docCreator struct {
			CreatorId int32 `alias:"creator_id"`
		}
		tDocuments := table.FivenetDocuments
		if err := tDocuments.
			SELECT(
				tDocuments.CreatorID.AS("creator_id"),
			).
			FROM(tDocuments).
			WHERE(tDocuments.ID.EQ(mysql.Int64(documentID))).
			LIMIT(1).
			QueryContext(ctx, tx, &docCreator); err != nil {
			return err
		}

		// In case the document has no creator (anymore).
		if docCreator.CreatorId > 0 {
			docCreatorId = docCreator.CreatorId
		}
	}

	approvalCondition := tApprovals.DocumentID.EQ(mysql.Int(documentID))
	if !pol.SelfApproveAllowed && docCreatorId > 0 {
		approvalCondition = mysql.AND(
			approvalCondition,
			tApprovals.UserID.NOT_EQ(mysql.Int32(docCreatorId)),
		)
	}

	var agg struct {
		Approved int32
		Declined int32
	}
	if err := tApprovals.
		SELECT(
			mysql.SUM(mysql.CASE().
				WHEN(tApprovals.Status.EQ(mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_APPROVED)))).
				THEN(mysql.Int(1)).
				ELSE(mysql.Int(0))).AS("approved"),
			mysql.SUM(mysql.CASE().
				WHEN(tApprovals.Status.EQ(mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_DECLINED)))).
				THEN(mysql.Int(1)).
				ELSE(mysql.Int(0))).AS("declined"),
		).
		FROM(tApprovals).
		WHERE(approvalCondition).
		QueryContext(ctx, tx, &agg); err != nil {
		return err
	}

	// Pending tasks
	var aggTasks struct {
		Total   int32
		Pending int32
	}
	if err := tApprovalTasks.
		SELECT(
			mysql.COUNT(tApprovalTasks.ID).AS("total"),
			mysql.SUM(mysql.CASE().
				WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))).
				THEN(mysql.Int(1)).
				ELSE(mysql.Int(0))).AS("pending"),
		).
		FROM(tApprovalTasks).
		WHERE(mysql.AND(
			tApprovalTasks.DocumentID.EQ(mysql.Int64(documentID)),
			tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
		)).
		QueryContext(ctx, tx, &aggTasks); err != nil {
		return err
	}

	requiredTotal := pol.GetRequiredCount()
	requiredRemaining := max(requiredTotal-agg.Approved, 0)

	anyDeclined := agg.Declined > 0
	// Doc is approved when we have enough approved and no declines
	// (if required_total=0, any declines block approval because all are required)
	var docApproved bool
	if pol.GetRuleKind() == documents.ApprovalRuleKind_APPROVAL_RULE_KIND_REQUIRE_ALL {
		// Ensure that there is at least one approval if all are required
		docApproved = !anyDeclined && agg.Approved > 0 && (agg.Approved >= aggTasks.Total)
	} else if pol.GetRuleKind() == documents.ApprovalRuleKind_APPROVAL_RULE_KIND_QUORUM_ANY {
		// Quorum-any: enough approvals, regardless of declines
		docApproved = (agg.Approved >= requiredTotal)
	}

	var apPoliciesActive int32
	if pol.DocumentId > 0 {
		apPoliciesActive = 1
	}

	// Update document meta rollups (document-level)
	if _, err := tDocumentsMeta.
		INSERT(
			tDocumentsMeta.DocumentID,
			tDocumentsMeta.RecomputedAt,
			tDocumentsMeta.Approved,
			tDocumentsMeta.ApRequiredTotal,
			tDocumentsMeta.ApCollectedApproved,
			tDocumentsMeta.ApRequiredRemaining,
			tDocumentsMeta.ApDeclinedCount,
			tDocumentsMeta.ApPendingCount,
			tDocumentsMeta.ApAnyDeclined,
			tDocumentsMeta.ApPoliciesActive,
		).
		VALUES(
			documentID,
			mysql.CURRENT_TIMESTAMP(),
			docApproved,
			requiredTotal,
			agg.Approved,
			requiredRemaining,
			agg.Declined,
			aggTasks.Pending,
			anyDeclined,
			apPoliciesActive,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tDocumentsMeta.RecomputedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tDocumentsMeta.Approved.SET(mysql.Bool(docApproved)),
			tDocumentsMeta.ApRequiredTotal.SET(mysql.Int32(requiredTotal)),
			tDocumentsMeta.ApCollectedApproved.SET(mysql.Int32(agg.Approved)),
			tDocumentsMeta.ApRequiredRemaining.SET(mysql.Int32(requiredRemaining)),
			tDocumentsMeta.ApDeclinedCount.SET(mysql.Int32(agg.Declined)),
			tDocumentsMeta.ApPendingCount.SET(mysql.Int32(aggTasks.Pending)),
			tDocumentsMeta.ApAnyDeclined.SET(mysql.Bool(anyDeclined)),
			tDocumentsMeta.ApPoliciesActive.SET(mysql.Int32(apPoliciesActive)),
		).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	if _, err := tApprovalPolicy.
		UPDATE().
		SET(
			tApprovalPolicy.AssignedCount.SET(mysql.Int32(aggTasks.Total)),
			tApprovalPolicy.ApprovedCount.SET(mysql.Int32(agg.Approved)),
			tApprovalPolicy.DeclinedCount.SET(mysql.Int32(agg.Declined)),
			tApprovalPolicy.PendingCount.SET(mysql.Int32(aggTasks.Pending)),
			tApprovalPolicy.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
		).
		WHERE(mysql.AND(
			tApprovalPolicy.DocumentID.EQ(mysql.Int(documentID)),
		)).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

// handleApprovalOnEditBehaviors checks the document's approval policy and,
// if configured, revokes existing approvals and reopens tasks.
func (s *Server) handleApprovalOnEditBehaviors(
	ctx context.Context,
	tx qrm.DB,
	doc *documents.Document,
) error {
	pol, err := s.getApprovalPolicy(
		ctx,
		tx,
		table.FivenetDocumentsApprovalPolicies.AS(
			"approval_policy",
		).DocumentID.EQ(
			mysql.Int64(doc.GetId()),
		),
	)
	if err != nil {
		return err
	}

	if pol == nil || pol.OnEditBehavior <= documents.OnEditBehavior_ON_EDIT_BEHAVIOR_KEEP_PROGRESS {
		// Nothing to do
		return nil
	}

	tApprovals := table.FivenetDocumentsApprovals
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	if _, err := tApprovals.
		UPDATE().
		SET(
			tApprovals.Status.SET(mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_REVOKED))),
			tApprovals.RevokedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(tApprovals.DocumentID.EQ(mysql.Int64(doc.GetId()))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	if _, err := tApprovalTasks.
		UPDATE().
		SET(
			tApprovalTasks.Status.SET(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
			tApprovalTasks.CompletedAt.SET(mysql.TimestampExp(mysql.NULL)),
		).
		WHERE(tApprovalTasks.DocumentID.EQ(mysql.Int64(doc.GetId()))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, doc.GetId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return err
	}

	return nil
}

func (s *Server) expireApprovalTasks(ctx context.Context) (int64, error) {
	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	now := time.Now()

	stmt := tApprovalTasks.
		UPDATE().
		SET(
			tApprovalTasks.Status.SET(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_EXPIRED))),
		).
		WHERE(mysql.AND(
			tApprovalTasks.DueAt.LT_EQ(mysql.TimestampT(now)),
			tApprovalTasks.Status.EQ(
				mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)),
			),
		)).
		LIMIT(250)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}
