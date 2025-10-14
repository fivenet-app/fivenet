package documents

import (
	"context"
	"errors"
	"time"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
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

// ListSignatureTasksInbox lists actionable signature tasks for the current user.
func (s *Server) ListSignatureTasksInbox(
	ctx context.Context,
	req *pbdocuments.ListSignatureTasksInboxRequest,
) (*pbdocuments.ListSignatureTasksInboxResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tSignatureTasks := table.FivenetDocumentsSignatureTasks.AS("sig_task")
	tSignatures := table.FivenetDocumentsSignatures

	var existsAccess mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		existsAccess = mysql.EXISTS(
			mysql.SELECT(mysql.Int(1)).
				FROM(tDAccess).
				WHERE(mysql.AND(
					tDAccess.TargetID.EQ(tSignatureTasks.DocumentID),
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

	// Eligibility for this task (user or job/grade)
	eligible := mysql.OR(
		mysql.AND(
			tSignatureTasks.AssigneeKind.EQ(
				mysql.Int32(int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_USER)),
			),
			tSignatureTasks.UserID.EQ(mysql.Int32(int32(userInfo.GetUserId()))),
		),
		mysql.AND(
			tSignatureTasks.AssigneeKind.EQ(
				mysql.Int32(
					int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_JOB_GRADE),
				),
			),
			tSignatureTasks.Job.EQ(mysql.String(userInfo.GetJob())),
			tSignatureTasks.MinimumGrade.LT_EQ(mysql.Int32(int32(userInfo.GetJobGrade()))),
		),
	)

	// NOT already signed in this round (treat VALID as “acted”; include DECLINED if you support it)
	notAlreadyActed := mysql.NOT(
		mysql.EXISTS(
			mysql.SELECT(mysql.Int(1)).
				FROM(tSignatures).
				WHERE(mysql.AND(
					tSignatures.DocumentID.EQ(tSignatureTasks.DocumentID),
					tSignatures.SnapshotDate.EQ(tSignatureTasks.SnapshotDate),
					tSignatures.UserID.EQ(mysql.Int32(int32(userInfo.GetUserId()))),
					tSignatures.Status.IN(
						mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID)),
						mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_DECLINED)),
					),
				)),
		),
	)

	// For JOB groups: only show the smallest slot_no pending row
	t2 := tSignatureTasks.AS("t2")
	minPendingSlotThisGroup := t2.
		SELECT(mysql.MIN(t2.SlotNo)).
		FROM(t2).
		WHERE(mysql.AND(
			t2.DocumentID.EQ(tSignatureTasks.DocumentID),
			t2.SnapshotDate.EQ(tSignatureTasks.SnapshotDate),
			t2.AssigneeKind.EQ(tSignatureTasks.AssigneeKind),
			mysql.IntExp(mysql.COALESCE(t2.UserID, mysql.Int32(0))).
				EQ(mysql.IntExp(mysql.COALESCE(tSignatureTasks.UserID, mysql.Int32(0)))),
			mysql.StringExp(mysql.COALESCE(t2.Job, mysql.String(""))).
				EQ(mysql.StringExp(mysql.COALESCE(tSignatureTasks.Job, mysql.String("")))),
			mysql.IntExp(mysql.COALESCE(t2.MinimumGrade, mysql.Int32(-1))).
				EQ(mysql.IntExp(mysql.COALESCE(tSignatureTasks.MinimumGrade, mysql.Int32(-1)))),
			t2.Status.EQ(
				mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING)),
			),
		)).
		LIMIT(1)

	onlyFirstSlot := mysql.OR(
		tSignatureTasks.AssigneeKind.EQ(
			mysql.Int32(int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_USER)),
		),
		minPendingSlotThisGroup.IN(tSignatureTasks.SlotNo),
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
		condition = condition.AND(tSignatureTasks.Status.IN(vals...))
	} else {
		condition = condition.AND(
			tSignatureTasks.Status.EQ(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING))),
		)
	}
	if req.OnlyDrafts != nil {
		condition = condition.AND(tDocumentShort.Draft.EQ(mysql.Bool(req.GetOnlyDrafts())))
	}

	// Count
	countStmt := tSignatureTasks.
		SELECT(mysql.COUNT(tSignatureTasks.ID).AS("data_count.total")).
		FROM(
			tSignatureTasks.
				INNER_JOIN(tDocumentShort, tDocumentShort.ID.EQ(tSignatureTasks.DocumentID)),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil &&
		!errors.Is(err, qrm.ErrNoRows) {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	resp := &pbdocuments.ListSignatureTasksInboxResponse{
		Tasks: []*documents.SignatureTask{},
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
			// Task
			tSignatureTasks.ID,
			tSignatureTasks.DocumentID,
			tSignatureTasks.SnapshotDate,
			tSignatureTasks.AssigneeKind,
			tSignatureTasks.UserID,
			tSignatureTasks.Job,
			tSignatureTasks.MinimumGrade,
			tSignatureTasks.Label,
			tSignatureTasks.SlotNo,
			tSignatureTasks.Status,
			tSignatureTasks.Comment,
			tSignatureTasks.CreatedAt,
			tSignatureTasks.CompletedAt,
			tSignatureTasks.DueAt,
			tSignatureTasks.CreatorID,
			tSignatureTasks.CreatorJob,
			// Requester (for USER tasks)
			tUser.ID,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.Job,
			tUser.JobGrade,
			// Document preview/meta
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
			tDMeta.Signed,
			tDMeta.SigRequiredTotal,
			tDMeta.SigCollectedValid,
			tDMeta.SigRequiredRemaining,
			tDMeta.SigDeclinedCount,
			tDMeta.SigPendingCount,
			tDMeta.SigAnyDeclined,
			tDMeta.SigPoliciesActive,
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
			tSignatureTasks.
				INNER_JOIN(tDocumentShort, tDocumentShort.ID.EQ(tSignatureTasks.DocumentID)).
				LEFT_JOIN(tUser, tUser.ID.EQ(tSignatureTasks.UserID)).
				LEFT_JOIN(tCreator, tCreator.ID.EQ(tSignatureTasks.CreatorID)).
				LEFT_JOIN(tDMeta, tDMeta.DocumentID.EQ(tDocumentShort.ID)),
		).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(tSignatureTasks.CreatedAt.ASC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Tasks); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Enrich labels/grades for display
	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for _, t := range resp.Tasks {
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

// ListSignaturePolicies.
func (s *Server) ListSignaturePolicies(
	ctx context.Context,
	req *pbdocuments.ListSignaturePoliciesRequest,
) (*pbdocuments.ListSignaturePoliciesResponse, error) {
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

	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies

	resp := &pbdocuments.ListSignaturePoliciesResponse{
		Policy: &documents.SignaturePolicy{},
	}

	resp.Policy, err = s.getSignaturePolicy(
		ctx,
		s.db,
		tSignaturePolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId())),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) getSignaturePolicy(
	ctx context.Context, tx qrm.DB, condition mysql.BoolExpression,
) (*documents.SignaturePolicy, error) {
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")

	stmt := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
			tSignaturePolicy.RuleKind,
			tSignaturePolicy.RequiredCount,
			tSignaturePolicy.BindingMode,
			tSignaturePolicy.AllowedTypesMask,
			tSignaturePolicy.AssignedCount,
			tSignaturePolicy.ApprovedCount,
			tSignaturePolicy.DeclinedCount,
			tSignaturePolicy.PendingCount,
			tSignaturePolicy.DueAt,
			tSignaturePolicy.StartedAt,
			tSignaturePolicy.CompletedAt,
			tSignaturePolicy.CreatedAt,
			tSignaturePolicy.UpdatedAt,
			tSignaturePolicy.DeletedAt,
		).
		FROM(tSignaturePolicy).
		WHERE(condition).
		LIMIT(1)

	pol := &documents.SignaturePolicy{}
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

func (s *Server) getOrCreateSignaturePolicy(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
) (*documents.SignaturePolicy, error) {
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies

	// Attempt to get any policy if no ID is specified
	condition := tSignaturePolicy.AS("signature_policy").DocumentID.EQ(mysql.Int64(documentId))

	pol, err := s.getSignaturePolicy(
		ctx,
		tx,
		condition,
	)
	if err != nil {
		return nil, err
	}
	if pol != nil {
		return pol, nil
	}

	// Create a new policy if it doesn't exist
	if _, err := tSignaturePolicy.
		INSERT(
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
			tSignaturePolicy.RuleKind,
			tSignaturePolicy.RequiredCount,
			tSignaturePolicy.BindingMode,
			tSignaturePolicy.AllowedTypesMask,
		).
		VALUES(
			documentId,
			mysql.CURRENT_TIMESTAMP(),
			int32(0), // TODO create new enum for signature rule kind
			false,
			int32(documents.SignatureBindingMode_SIGNATURE_BINDING_MODE_NONBINDING),
			"[]",
		).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pol, err = s.getSignaturePolicy(
		ctx,
		tx,
		condition,
	)
	if err != nil {
		return nil, err
	}
	if pol == nil {
		return nil, errorsdocuments.ErrFailedQuery
	}

	return pol, nil
}

// UpsertSignaturePolicy.
func (s *Server) UpsertSignaturePolicy(
	ctx context.Context,
	req *pbdocuments.UpsertSignaturePolicyRequest,
) (*pbdocuments.UpsertSignaturePolicyResponse, error) {
	pol := req.GetPolicy()

	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.documents.id", pol.GetDocumentId()},
	)

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

	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	if pol.GetDocumentId() == 0 {
		stmt := tSignaturePolicy.
			INSERT(
				tSignaturePolicy.DocumentID,
				tSignaturePolicy.SnapshotDate,
				tSignaturePolicy.RuleKind,
				tSignaturePolicy.RequiredCount,
				tSignaturePolicy.BindingMode,
				tSignaturePolicy.AllowedTypesMask,
				tSignaturePolicy.DueAt,
			).
			VALUES(
				pol.GetDocumentId(),
				mysql.CURRENT_TIMESTAMP(),
				pol.GetRequiredCount(),
				int32(pol.GetBindingMode()),
				pol.GetAllowedTypes(),
				pol.GetDueAt(),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		stmt := tSignaturePolicy.
			UPDATE(
				tSignaturePolicy.DocumentID,
				tSignaturePolicy.SnapshotDate,
				tSignaturePolicy.RuleKind,
				tSignaturePolicy.RequiredCount,
				tSignaturePolicy.BindingMode,
				tSignaturePolicy.AllowedTypesMask,
				tSignaturePolicy.DueAt,
			).
			SET(
				pol.GetDocumentId(),
				mysql.CURRENT_TIMESTAMP(),
				pol.GetRuleKind(),
				pol.GetRequiredCount(),
				int32(pol.GetBindingMode()),
				pol.GetAllowedTypes(),
				pol.GetDueAt(),
			).
			WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int64(pol.GetDocumentId()))).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	tSignaturePolicy = tSignaturePolicy.AS("signature_policy")

	policy, err := s.getSignaturePolicy(
		ctx,
		tx,
		tSignaturePolicy.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, policy.GetDocumentId(), policy.GetSnapshotDate().AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	policy, err = s.getSignaturePolicy(
		ctx,
		s.db,
		tSignaturePolicy.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.UpsertSignaturePolicyResponse{
		Policy: pol,
	}, nil
}

// DeleteSignaturePolicy.
func (s *Server) DeleteSignaturePolicy(
	ctx context.Context,
	req *pbdocuments.DeleteSignaturePolicyRequest,
) (*pbdocuments.DeleteSignaturePolicyResponse, error) {
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies
	if _, err := tSignaturePolicy.
		DELETE().
		WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int(req.GetDocumentId()))).
		LIMIT(1).
		ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteSignaturePolicyResponse{}, nil
}

// ListSignatureTasks.
func (s *Server) ListSignatureTasks(
	ctx context.Context,
	req *pbdocuments.ListSignatureTasksRequest,
) (*pbdocuments.ListSignatureTasksResponse, error) {
	tSignatureTasks := table.FivenetDocumentsSignatureTasks

	// Build WHERE
	condition := tSignatureTasks.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))

	if req.GetSnapshotDate() != nil {
		snap := req.GetSnapshotDate().AsTime()
		condition = condition.AND(tSignatureTasks.SnapshotDate.EQ(mysql.TimestampT(snap)))
	}

	// Optional status filter
	if len(req.GetStatuses()) > 0 {
		vals := make([]mysql.Expression, 0, len(req.GetStatuses()))
		for _, st := range req.GetStatuses() {
			vals = append(vals, mysql.Int32(int32(st.Number())))
		}
		condition = condition.AND(tSignatureTasks.Status.IN(vals...))
	}

	resp := &pbdocuments.ListSignatureTasksResponse{
		Signatures: []*documents.Signature{},
	}

	stmt := tSignatureTasks.
		SELECT(
			tSignatureTasks.ID,
			tSignatureTasks.DocumentID,
			tSignatureTasks.SnapshotDate,
			tSignatureTasks.AssigneeKind,
			tSignatureTasks.UserID,
			tSignatureTasks.Job,
			tSignatureTasks.MinimumGrade,
			tSignatureTasks.Label,
			tSignatureTasks.SlotNo,
			tSignatureTasks.Status,
			tSignatureTasks.Comment,
			tSignatureTasks.CreatedAt,
			tSignatureTasks.CompletedAt,
			tSignatureTasks.DueAt,
			tSignatureTasks.CreatorID,
			tSignatureTasks.CreatorJob,
			tSignatureTasks.SignatureID,
		).
		FROM(tSignatureTasks).
		WHERE(condition).
		ORDER_BY(tSignatureTasks.CreatedAt.ASC()).
		LIMIT(15)

	if err := stmt.QueryContext(ctx, s.db, &resp.Signatures); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) getSignatureTask(
	ctx context.Context,
	tx qrm.DB,
	taskId int64,
) (*documents.SignatureTask, error) {
	tSignatureTasks := table.FivenetDocumentsSignatureTasks.AS("signature_task")

	stmt := tSignatureTasks.
		SELECT(
			tSignatureTasks.ID,
			tSignatureTasks.DocumentID,
			tSignatureTasks.SnapshotDate,
			tSignatureTasks.AssigneeKind,
			tSignatureTasks.UserID,
			tSignatureTasks.Job,
			tSignatureTasks.MinimumGrade,
			tSignatureTasks.Label,
			tSignatureTasks.SlotNo,
			tSignatureTasks.Status,
			tSignatureTasks.Comment,
			tSignatureTasks.CreatedAt,
			tSignatureTasks.CompletedAt,
			tSignatureTasks.DueAt,
			tSignatureTasks.CreatorID,
			tSignatureTasks.CreatorJob,
			tSignatureTasks.SignatureID,
		).
		FROM(tSignatureTasks).
		WHERE(tSignatureTasks.ID.EQ(mysql.Int64(taskId))).
		LIMIT(1)

	var task documents.SignatureTask
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

// UpsertSignatureTasks.
func (s *Server) UpsertSignatureTasks(
	ctx context.Context,
	req *pbdocuments.UpsertSignatureTasksRequest,
) (*pbdocuments.UpsertSignatureTasksResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.documents.sig.document_id", req.GetDocumentId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Resolve policy & snapshot
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")

	var pol documents.SignaturePolicy
	if err := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.GetDocumentId() == 0 {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
	}
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

	tSignatureTasks := table.FivenetDocumentsSignatureTasks

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	created := int32(0)
	ensured := int32(0)

	for _, seed := range req.GetSeeds() {
		isUser := seed.GetUserId() != 0
		if isUser {
			// ensure one USER task exists
			var cnt struct{ C int32 }
			if err := mysql.
				SELECT(mysql.COUNT(tSignatureTasks.ID).AS("C")).
				FROM(tSignatureTasks).
				WHERE(mysql.AND(
					tSignatureTasks.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
					tSignatureTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
					tSignatureTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_USER))),
					tSignatureTasks.UserID.EQ(mysql.Int32(seed.GetUserId())),
				)).
				LIMIT(1).
				QueryContext(ctx, tx, &cnt); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
			if cnt.C > 0 {
				ensured++
				continue
			}

			// insert USER task with slot_no=1
			if _, err := tSignatureTasks.
				INSERT(
					tSignatureTasks.DocumentID, tSignatureTasks.SnapshotDate,
					tSignatureTasks.AssigneeKind, tSignatureTasks.UserID, tSignatureTasks.Label,
					tSignatureTasks.SlotNo, tSignatureTasks.Status, tSignatureTasks.Comment,
					tSignatureTasks.DueAt, tSignatureTasks.CreatorID, tSignatureTasks.CreatorJob,
				).
				VALUES(
					pol.GetDocumentId(), snap,
					int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_USER), seed.GetUserId(), seed.GetLabel(), 1,
					int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING), seed.GetComment(), dbutils.TimestampToMySQL(seed.GetDueAt()),
					userInfo.GetUserId(), userInfo.GetJob(),
				).
				ExecContext(ctx, tx); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
			SELECT(mysql.COUNT(tSignatureTasks.ID).AS("C")).
			FROM(tSignatureTasks).
			WHERE(mysql.AND(
				tSignatureTasks.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
				tSignatureTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
				tSignatureTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_JOB_GRADE))),
				tSignatureTasks.Job.EQ(mysql.String(seed.GetJob())),
				tSignatureTasks.MinimumGrade.EQ(mysql.Int32(seed.GetMinimumGrade())),
				tSignatureTasks.Status.EQ(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING))),
			)).
			LIMIT(1).
			QueryContext(ctx, tx, &have); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if have.C >= int32(slots) {
			ensured++
			continue
		}

		// insert missing [have.C+1 .. slots]
		ins := tSignatureTasks.
			INSERT(
				tSignatureTasks.DocumentID,
				tSignatureTasks.SnapshotDate,
				tSignatureTasks.AssigneeKind,
				tSignatureTasks.Job,
				tSignatureTasks.MinimumGrade,
				tSignatureTasks.Label,
				tSignatureTasks.SlotNo,
				tSignatureTasks.Status,
				tSignatureTasks.Comment,
				tSignatureTasks.DueAt,
				tSignatureTasks.CreatorID,
				tSignatureTasks.CreatorJob,
			)
		for slot := have.C + 1; slot <= slots; slot++ {
			ins = ins.
				VALUES(
					pol.GetDocumentId(),
					snap,

					int32(
						documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_JOB_GRADE,
					),
					seed.GetJob(),
					seed.GetMinimumGrade(),
					seed.GetLabel(),
					slot,
					int32(
						documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING,
					),
					seed.GetComment(),
					dbutils.TimestampToMySQL(seed.GetDueAt()),
					userInfo.GetUserId(),
					userInfo.GetJob(),
				)
		}
		if _, err := ins.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		created += int32(slots) - have.C
	}

	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, err
	}

	// Create document activity entry
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   pol.GetDocumentId(),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_SIGNING_ASSIGNED,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.UpsertSignatureTasksResponse{
		TasksCreated: created,
		TasksEnsured: ensured,
		Policy:       &pol,
	}, nil
}

// DeleteSignatureTasks.
func (s *Server) DeleteSignatureTasks(
	ctx context.Context,
	req *pbdocuments.DeleteSignatureTasksRequest,
) (*pbdocuments.DeleteSignatureTasksResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.documents.sig.document_id", req.GetDocumentId()},
	)

	// Resolve policy & snapshot
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")
	var pol documents.SignaturePolicy
	if err := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.GetDocumentId() == 0 {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
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

	tTasks := table.FivenetDocumentsSignatureTasks
	condition := mysql.AND(
		tTasks.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
		tTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
	)

	// Delete all pending?
	if req.GetDeleteAllPending() {
		condition = condition.AND(
			tTasks.Status.EQ(
				mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING)),
			),
		)
	} else if len(req.TaskIds) > 0 {
		ids := make([]mysql.Expression, 0, len(req.TaskIds))
		for _, id := range req.TaskIds {
			ids = append(ids, mysql.Int64(id))
		}
		condition = condition.AND(tTasks.ID.IN(ids...))
	} else {
		return &pbdocuments.DeleteSignatureTasksResponse{}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	if _, err := tTasks.
		DELETE().
		WHERE(condition).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteSignatureTasksResponse{}, nil
}

func (s *Server) ListSignatures(
	ctx context.Context,
	req *pbdocuments.ListSignaturesRequest,
) (*pbdocuments.ListSignaturesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")
	tSignatures := table.FivenetDocumentsSignatures.AS("signature")

	// Resolve policy & ensure access
	var pol documents.SignaturePolicy
	if err := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.GetDocumentId() == 0 {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		pol.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	snap := pol.GetSnapshotDate().AsTime()
	if req.GetSnapshotDate() != nil {
		snap = req.GetSnapshotDate().AsTime()
	}

	condition := mysql.AND(
		tSignatures.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
		tSignatures.SnapshotDate.EQ(mysql.TimestampT(snap)),
	)

	if req.GetStatus() != 0 {
		condition = condition.AND(
			tSignatures.Status.EQ(mysql.Int32(int32(req.GetStatus().Number()))),
		)
	}
	if req.GetUserId() != 0 {
		condition = condition.AND(tSignatures.UserID.EQ(mysql.Int32(req.GetUserId())))
	}
	if req.GetTaskId() != 0 {
		condition = condition.AND(tSignatures.TaskID.EQ(mysql.Int64(req.GetTaskId())))
	}

	// Count
	countStmt := tSignatures.
		SELECT(mysql.COUNT(tSignatures.ID)).
		FROM(tSignatures).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	resp := &pbdocuments.ListSignaturesResponse{
		Signatures: []*documents.Signature{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Page fetch
	stmt := tSignatures.
		SELECT(
			tSignatures.ID, tSignatures.DocumentID, tSignatures.SnapshotDate,
			tSignatures.UserID, tSignatures.UserJob, tSignatures.UserJobGrade,
			tSignatures.Type, tSignatures.PayloadSvg, tSignatures.StampID,
			tSignatures.Status, tSignatures.Comment, tSignatures.TaskID,
			tSignatures.CreatedAt, tSignatures.RevokedAt,
		).
		FROM(tSignatures).
		WHERE(condition).
		ORDER_BY(tSignatures.CreatedAt.DESC()).
		LIMIT(15)

	if err := stmt.QueryContext(ctx, s.db, &resp.Signatures); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return resp, nil
}

func (s *Server) RevokeSignature(
	ctx context.Context,
	req *pbdocuments.RevokeSignatureRequest,
) (*pbdocuments.RevokeSignatureResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tSignatures := table.FivenetDocumentsSignatures.AS("signature")
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")

	// Load the signature (and related policy/doc)
	var sig documents.Signature
	if err := tSignatures.
		SELECT(
			tSignatures.ID, tSignatures.DocumentID, tSignatures.SnapshotDate,
			tSignatures.UserID, tSignatures.UserJob, tSignatures.UserJobGrade,
			tSignatures.Type, tSignatures.PayloadSvg, tSignatures.StampID,
			tSignatures.Status, tSignatures.Comment, tSignatures.TaskID,
			tSignatures.CreatedAt, tSignatures.RevokedAt,
		).
		FROM(tSignatures).
		WHERE(tSignatures.ID.EQ(mysql.Int64(req.GetSignatureId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &sig); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Check permission on the document
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		sig.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	tSignatures = table.FivenetDocumentsSignatures
	tSignatureTasks := table.FivenetDocumentsSignatureTasks

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	// Mark revoked
	if _, err := tSignatures.
		UPDATE().
		SET(
			tSignatures.Status.SET(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_REVOKED))),
			tSignatures.RevokedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tSignatures.Comment.SET(mysql.String(req.GetComment())),
		).
		WHERE(tSignatures.ID.EQ(mysql.Int64(req.GetSignatureId()))).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Reload artifact for response
	tSignatures = table.FivenetDocumentsSignatures
	if err := tSignatures.
		SELECT(
			tSignatures.ID, tSignatures.DocumentID, tSignatures.SnapshotDate,
			tSignatures.UserID, tSignatures.UserJob, tSignatures.UserJobGrade,
			tSignatures.Type, tSignatures.PayloadSvg, tSignatures.StampID,
			tSignatures.Status, tSignatures.Comment, tSignatures.TaskID,
			tSignatures.CreatedAt, tSignatures.RevokedAt,
		).
		FROM(tSignatures).
		WHERE(tSignatures.ID.EQ(mysql.Int64(req.GetSignatureId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &sig); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Resolve policy for recompute
	var pol documents.SignaturePolicy
	if err := tSignaturePolicy.
		SELECT(tSignaturePolicy.DocumentID, tSignaturePolicy.SnapshotDate).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int64(sig.GetDocumentId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if sig.GetTaskId() > 0 {
		// Set PENDING & clear decider snapshot
		if _, err = tSignatureTasks.
			UPDATE().
			SET(
				tSignatureTasks.Status.SET(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING))),
				tSignatureTasks.CompletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			).
			WHERE(
				tSignatureTasks.ID.EQ(mysql.Int64(sig.GetTaskId())),
			).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Recompute rollups for tasks, policy and document
	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, err
	}

	// Create document activity entry
	comment := req.GetComment()
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   pol.GetDocumentId(),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_SIGNING_REVOKED,
		Reason:       &comment,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.RevokeSignatureResponse{Signature: &sig}, nil
}

// DecideSignature creates/upserts a signature artifact.
// If task_id is provided, it will atomically consume that task (when PENDING),
// link the artifact to the task, and set task.signature_id + completed_at.
func (s *Server) DecideSignature(
	ctx context.Context,
	req *pbdocuments.DecideSignatureRequest,
) (*pbdocuments.DecideSignatureResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Resolve policy, doc, snapshot
	pol, err := s.getOrCreateSignaturePolicy(
		ctx,
		s.db,
		req.GetDocumentId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Access check: signer must be able to view/act on document (stronger checks happen when consuming a task)
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		pol.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	snap := pol.GetSnapshotDate().AsTime()

	tSignatureTasks := table.FivenetDocumentsSignatureTasks
	tSignatures := table.FivenetDocumentsSignatures

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	var decidedTask *documents.SignatureTask
	now := time.Now()

	// If task_id provided, validate and mark decided later
	if req.GetTaskId() > 0 {
		// Load the task row into the protobuf message directly
		decidedTask, err = s.getSignatureTask(ctx, tx, req.GetTaskId())
		if err != nil {
			return nil, err
		}
		if decidedTask.DocumentId == 0 {
			return nil, errorsdocuments.ErrNotFoundOrNoPerms
		}

		// Task must belong to same policy/snapshot and be PENDING
		if decidedTask.GetDocumentId() != pol.GetDocumentId() ||
			decidedTask.GetSnapshotDate().AsTime() != snap {
			return nil, errorsdocuments.ErrDocAccessViewDenied
		}
		if decidedTask.GetStatus() != documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING {
			return nil, errorsdocuments.ErrSigningTaskAlreadyHandled
		}

		if decidedTask.GetAssigneeKind() == documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_USER &&
			decidedTask.GetUserId() != userInfo.GetUserId() {
			return nil, errorsdocuments.ErrDocAccessViewDenied
		} else if decidedTask.GetAssigneeKind() == documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_JOB_GRADE {
			if decidedTask.GetJob() != userInfo.GetJob() ||
				decidedTask.GetMinimumGrade() >= userInfo.GetJobGrade() {
				return nil, errorsdocuments.ErrDocAccessViewDenied
			}
		} else {
			return nil, errorsdocuments.ErrDocAccessViewDenied
		}
	}

	var existing struct {
		ID     int64 `alias:"id"`
		Status int32 `alias:"status"`
	}
	if err := tSignatures.
		SELECT(
			tSignatures.ID.AS("id"),
			tSignatures.Status.AS("status"),
		).
		FROM(tSignatures).
		WHERE(mysql.AND(
			tSignatures.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
			tSignatures.SnapshotDate.EQ(mysql.TimestampT(snap)),
			tSignatures.UserID.EQ(mysql.Int32(int32(userInfo.GetUserId()))),
		)).
		LIMIT(1).
		QueryContext(ctx, tx, &existing); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if existing.ID != 0 &&
		(existing.Status == int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID) ||
			existing.Status == int32(documents.SignatureStatus_SIGNATURE_STATUS_DECLINED)) {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrApprovalTaskAlreadyHandled)
	}

	// Insert/Upsert artifact (unique by document_id + snapshot_date + user_id)
	newSig := &documents.Signature{}
	ins := tSignatures.
		INSERT(
			tSignatures.DocumentID,
			tSignatures.SnapshotDate,

			tSignatures.UserID,
			tSignatures.UserJob,
			tSignatures.UserJobGrade,
			tSignatures.Type,
			tSignatures.PayloadSvg,
			tSignatures.StampID,
			tSignatures.Status,
			tSignatures.Comment,
			tSignatures.TaskID,
			tSignatures.CreatedAt,
		).
		VALUES(
			pol.GetDocumentId(), snap,
			userInfo.GetUserId(), userInfo.GetJob(), mysql.Int32(int32(userInfo.GetJobGrade())),
			int32(req.GetType()), req.GetPayloadSvg(), nilOrInt64(req.GetStampId()),
			int32(
				documents.SignatureStatus_SIGNATURE_STATUS_VALID,
			), req.GetComment(), nilOrInt64(req.GetTaskId()), now,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tSignatures.Type.SET(mysql.Int32(int32(req.GetType()))),
			tSignatures.PayloadSvg.SET(mysql.String(req.GetPayloadSvg())),
			tSignatures.StampID.SET(nilOrInt64(req.GetStampId())),
			tSignatures.Status.SET(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID))),
			tSignatures.Comment.SET(mysql.String(req.GetComment())),
			tSignatures.TaskID.SET(nilOrInt64(req.GetTaskId())),
		)
	if _, err := ins.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Read back the artifact row (protobuf hydration)
	tSignatures = table.FivenetDocumentsSignatures
	if err := tSignatures.
		SELECT(
			tSignatures.ID,
			tSignatures.DocumentID,

			tSignatures.SnapshotDate,
			tSignatures.UserID,
			tSignatures.UserJob,
			tSignatures.UserJobGrade,
			tSignatures.Type,
			tSignatures.PayloadSvg,
			tSignatures.StampID,
			tSignatures.Status,
			tSignatures.Comment,
			tSignatures.TaskID,
			tSignatures.CreatedAt,
			tSignatures.RevokedAt,
		).
		FROM(tSignatures).
		WHERE(mysql.AND(
			tSignatures.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
			tSignatures.SnapshotDate.EQ(mysql.TimestampT(snap)),
			tSignatures.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
		)).
		LIMIT(1).
		QueryContext(ctx, tx, newSig); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// If task_id provided, mark the task completed and link the signature id
	if req.GetTaskId() > 0 && decidedTask != nil {
		tSignatureTasks = table.FivenetDocumentsSignatureTasks
		if _, err := tSignatureTasks.
			UPDATE().
			SET(
				tSignatureTasks.Status.SET(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_SIGNED))),
				tSignatureTasks.CompletedAt.SET(mysql.TimestampT(now)),
				tSignatureTasks.SignatureID.SET(mysql.Int64(newSig.GetId())),
			).
			WHERE(tSignatureTasks.ID.EQ(mysql.Int64(req.GetTaskId()))).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		// reflect new state in response
		decidedTask.Status = documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_SIGNED
		decidedTask.CompletedAt = timestamp.New(now)
		decidedTask.SignatureId = &newSig.Id
	}

	// Recompute policy/meta from artifacts for this (policy, snapshot)
	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), snap); err != nil {
		return nil, err
	}

	// Create document activity entry
	comment := req.GetComment()
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   pol.GetDocumentId(),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_SIGNING_SIGNED,
		Reason:       &comment,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DecideSignatureResponse{
		Signature: newSig,
		Task:      decidedTask,
		Policy:    pol,
	}, nil
}

// ReopenSignature.
func (s *Server) ReopenSignature(
	ctx context.Context,
	req *pbdocuments.ReopenSignatureRequest,
) (*pbdocuments.ReopenSignatureResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tSignatures := table.FivenetDocumentsSignatures.AS("signature")
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")

	// Load signature
	var sig documents.Signature
	if err := tSignatures.
		SELECT(
			tSignatures.ID, tSignatures.DocumentID, tSignatures.SnapshotDate,
			tSignatures.UserID, tSignatures.UserJob, tSignatures.UserJobGrade,
			tSignatures.Type, tSignatures.PayloadSvg, tSignatures.StampID,
			tSignatures.Status, tSignatures.Comment, tSignatures.TaskID,
			tSignatures.CreatedAt, tSignatures.RevokedAt,
		).
		FROM(tSignatures).
		WHERE(tSignatures.ID.EQ(mysql.Int64(req.GetSignatureId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &sig); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Permission: require EDIT on document to reopen artifact
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		sig.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	// Resolve policy for recompute
	var pol documents.SignaturePolicy
	if err := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int64(sig.GetDocumentId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.GetDocumentId() == 0 {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	tSignatures = table.FivenetDocumentsSignatures

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	// Set back to VALID, clear revoked_at; optional note into reason/comment
	if _, err := tSignatures.
		UPDATE().
		SET(
			tSignatures.Status.SET(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID))),
			tSignatures.RevokedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tSignatures.Comment.SET(mysql.String(req.GetComment())),
		).
		WHERE(tSignatures.ID.EQ(mysql.Int64(req.GetSignatureId()))).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Reload for response
	tSignatures = table.FivenetDocumentsSignatures
	if err := tSignatures.
		SELECT(
			tSignatures.ID, tSignatures.DocumentID, tSignatures.SnapshotDate,
			tSignatures.UserID, tSignatures.UserJob, tSignatures.UserJobGrade,
			tSignatures.Type, tSignatures.PayloadSvg, tSignatures.StampID,
			tSignatures.Status, tSignatures.Comment, tSignatures.TaskID,
			tSignatures.CreatedAt, tSignatures.RevokedAt,
		).
		FROM(tSignatures).
		WHERE(tSignatures.ID.EQ(mysql.Int64(req.GetSignatureId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &sig); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Recompute signatures for this policy+snapshot
	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, err
	}

	// Create document activity entry
	comment := req.GetComment()
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   pol.GetDocumentId(),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_SIGNING_REVOKED,
		Reason:       &comment,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.ReopenSignatureResponse{Signature: &sig}, nil
}

// RecomputeSignatureStatus.
func (s *Server) RecomputeSignatureStatus(
	ctx context.Context,
	req *pbdocuments.RecomputeSignatureStatusRequest,
) (*pbdocuments.RecomputeSignatureStatusResponse, error) {
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

	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")

	pol, err := s.getSignaturePolicy(
		ctx,
		s.db,
		tSignaturePolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId())),
	)
	if err != nil {
		return nil, err
	}

	if err := s.recomputeSignaturePolicyTx(ctx, s.db, req.GetDocumentId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pol, err = s.getSignaturePolicy(
		ctx,
		s.db,
		tSignaturePolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.RecomputeSignatureStatusResponse{
		Policy: pol,
	}, nil
}

// recomputeSignaturePolicyTx recalculates signature counters for a policy+snapshot
// and updates fivenet_documents_meta for list-time flags.
func (s *Server) recomputeSignaturePolicyTx(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	snap time.Time,
) error {
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies
	tSignatures := table.FivenetDocumentsSignatures
	tSignatureTasks := table.FivenetDocumentsApprovalTasks
	tDocumentsMeta := table.FivenetDocumentsMeta

	// Load policy (required_count, etc.) if given and exists
	var pol documents.SignaturePolicy
	if err := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
			tSignaturePolicy.RequiredCount,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int64(documentID))).
		LIMIT(1).
		QueryContext(ctx, tx, &pol); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	var agg struct {
		Valid    int32
		Declined int32
	}
	if err := tSignatures.
		SELECT(
			mysql.SUM(mysql.CASE().
				WHEN(tSignatures.Status.EQ(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID)))).
				THEN(mysql.Int(1)).
				ELSE(mysql.Int(0))).AS("valid"),
			mysql.SUM(mysql.CASE().
				WHEN(tSignatures.Status.EQ(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_DECLINED)))).
				THEN(mysql.Int(1)).
				ELSE(mysql.Int(0))).AS("declined"),
		).
		FROM(tSignatures).
		WHERE(mysql.AND(
			tSignatures.DocumentID.EQ(mysql.Int(documentID)),
		)).
		QueryContext(ctx, tx, &agg); err != nil {
		return err
	}

	// Pending tasks
	var aggTasks struct {
		Pending  int32
		Assigned int32
	}
	if err := tSignatureTasks.
		SELECT(
			mysql.COUNT(tSignatureTasks.ID).AS("assigned"),
			mysql.SUM(mysql.CASE().
				WHEN(tSignatureTasks.Status.EQ(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING)))).
				THEN(mysql.Int(1)).
				ELSE(mysql.Int(0))).AS("pending"),
		).
		FROM(tSignatureTasks).
		WHERE(mysql.AND(
			tSignatureTasks.DocumentID.EQ(mysql.Int64(documentID)),
			tSignatureTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
		)).
		QueryContext(ctx, tx, &aggTasks); err != nil {
		return err
	}

	requiredTotal := pol.GetRequiredCount()
	requiredRemaining := max(requiredTotal-agg.Valid, 0)

	anyDeclined := agg.Declined > 0
	// Doc is signed when we have enough valid signatures and no declines
	// (if required_total=0, any declines block approval because all are required)
	docSigned := (requiredTotal > 0 && agg.Valid >= requiredTotal) ||
		(requiredTotal == 0 && !anyDeclined)
	var sigPoliciesActive int32
	if pol.DocumentId > 0 {
		sigPoliciesActive = 1
	}

	// update meta rollups (document-level)
	if _, err := tDocumentsMeta.
		INSERT(
			tDocumentsMeta.DocumentID,
			tDocumentsMeta.RecomputedAt,
			tDocumentsMeta.Signed,
			tDocumentsMeta.SigRequiredTotal,
			tDocumentsMeta.SigCollectedValid,
			tDocumentsMeta.SigRequiredRemaining,
			tDocumentsMeta.SigDeclinedCount,
			tDocumentsMeta.SigPendingCount,
			tDocumentsMeta.SigAnyDeclined,
			tDocumentsMeta.SigPoliciesActive,
		).
		VALUES(
			documentID,
			mysql.CURRENT_TIMESTAMP(),
			docSigned,
			requiredTotal,
			agg.Valid,
			requiredRemaining,
			agg.Declined,
			aggTasks.Pending,
			anyDeclined,
			sigPoliciesActive,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tDocumentsMeta.RecomputedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tDocumentsMeta.Signed.SET(mysql.Bool(docSigned)),
			tDocumentsMeta.SigRequiredTotal.SET(mysql.Int32(requiredTotal)),
			tDocumentsMeta.SigCollectedValid.SET(mysql.Int32(agg.Valid)),
			tDocumentsMeta.SigRequiredRemaining.SET(mysql.Int32(requiredRemaining)),
			tDocumentsMeta.SigDeclinedCount.SET(mysql.Int32(agg.Declined)),
			tDocumentsMeta.SigPendingCount.SET(mysql.Int32(aggTasks.Pending)),
			tDocumentsMeta.SigAnyDeclined.SET(mysql.Bool(anyDeclined)),
			tDocumentsMeta.SigPoliciesActive.SET(mysql.Int32(sigPoliciesActive)),
		).
		ExecContext(ctx, tx); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := tSignaturePolicy.
		UPDATE().
		SET(
			tSignaturePolicy.AssignedCount.SET(mysql.Int32(aggTasks.Assigned)),
			tSignaturePolicy.ApprovedCount.SET(mysql.Int32(agg.Valid)),
			tSignaturePolicy.DeclinedCount.SET(mysql.Int32(agg.Declined)),
			tSignaturePolicy.PendingCount.SET(mysql.Int32(aggTasks.Pending)),
			tSignaturePolicy.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
		).
		WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int(documentID))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

// nilOrInt64 helpers for nullable ints (IDs).
func nilOrInt64(v int64) mysql.IntegerExpression {
	if v == 0 {
		return mysql.IntExp(mysql.NULL)
	}
	return mysql.Int64(v)
}

// handleSignatureBindingMode checks the document's signature policies and,
// if configured, revokes the matching signatures and reopens tasks.
func (s *Server) handleSignatureBindingMode(
	ctx context.Context,
	tx qrm.DB,
	doc *documents.Document,
) error {
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")

	pol, err := s.getSignaturePolicy(
		ctx,
		tx,
		tSignaturePolicy.DocumentID.EQ(mysql.Int64(doc.GetId())),
	)
	if err != nil {
		return err
	}

	tSignatures := table.FivenetDocumentsSignatures
	tSignatureTasks := table.FivenetDocumentsSignatureTasks

	if pol.BindingMode <= documents.SignatureBindingMode_SIGNATURE_BINDING_MODE_NONBINDING {
		// Nothing to do
		return nil
	}

	if _, err := tSignatures.
		UPDATE().
		SET(
			tSignatures.Status.SET(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_REVOKED))),
			tSignatures.RevokedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(tSignatures.DocumentID.EQ(mysql.Int64(doc.GetId()))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	if _, err := tSignatureTasks.
		UPDATE().
		SET(
			tSignatureTasks.Status.SET(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING))),
			tSignatureTasks.CompletedAt.SET(mysql.TimestampExp(mysql.NULL)),
		).
		WHERE(tSignatureTasks.DocumentID.EQ(mysql.Int64(doc.GetId()))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	if err := s.recomputeSignaturePolicyTx(ctx, tx, doc.GetId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return err
	}

	return nil
}

func (s *Server) expireSignatureTasks(ctx context.Context) (int64, error) {
	tSignatureTasks := table.FivenetDocumentsSignatureTasks
	now := time.Now()

	stmt := tSignatureTasks.
		UPDATE().
		SET(
			tSignatureTasks.Status.SET(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_EXPIRED))),
		).
		WHERE(mysql.AND(
			tSignatureTasks.DueAt.LT_EQ(mysql.TimestampT(now)),
			tSignatureTasks.Status.EQ(
				mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING)),
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
