package documents

import (
	"context"
	"errors"
	"fmt"
	"time"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	documentsstore "github.com/fivenet-app/fivenet/v2026/stores/documents"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListApprovalTasksInbox(
	ctx context.Context,
	req *pbdocuments.ListApprovalTasksInboxRequest,
) (*pbdocuments.ListApprovalTasksInboxResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	count, tasks, err := s.store.ListApprovalTasksInbox(
		ctx,
		documentsstore.ListApprovalTasksInboxQuery{
			Pagination:      req.GetPagination(),
			UserInfo:        userInfo,
			Statuses:        req.GetStatuses(),
			NotAlreadyActed: req.GetNotAlreadyActed(),
			OnlyDrafts:      req.OnlyDrafts,
		},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pagReq := req.GetPagination()
	if pagReq == nil {
		pagReq = &database.PaginationRequest{}
	}

	resp := &pbdocuments.ListApprovalTasksInboxResponse{
		Tasks: tasks,
	}

	pag, _ := pagReq.GetResponseWithPageSize(count.Total, 20)
	resp.Pagination = pag
	if count.Total <= 0 {
		return resp, nil
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for _, t := range resp.GetTasks() {
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

func (s *Server) ListApprovalPolicies(
	ctx context.Context,
	req *pbdocuments.ListApprovalPoliciesRequest,
) (*pbdocuments.ListApprovalPoliciesResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
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

	docMeta, err := s.store.GetDocumentMeta(ctx, s.db, req.GetDocumentId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.ListApprovalPoliciesResponse{
		Policy:  policy,
		DocMeta: docMeta,
	}, nil
}

func (s *Server) getApprovalPolicy(
	ctx context.Context, tx qrm.DB, condition mysql.BoolExpression,
) (*documentsapproval.ApprovalPolicy, error) {
	pol, err := s.store.GetApprovalPolicy(ctx, tx, condition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	return pol, nil
}

func (s *Server) getOrCreateApprovalPolicy(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	snapshotDate *timestamp.Timestamp,
) (*documentsapproval.ApprovalPolicy, error) {
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
	if err = s.createApprovalPolicy(ctx, tx, documentId, &documentsapproval.ApprovalPolicy{
		SnapshotDate:       snapshotDate,
		RuleKind:           documentsapproval.ApprovalRuleKind_APPROVAL_RULE_KIND_REQUIRE_ALL,
		RequiredCount:      &requiredCount,
		OnEditBehavior:     documentsapproval.OnEditBehavior_ON_EDIT_BEHAVIOR_KEEP_PROGRESS,
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
	pol *documentsapproval.ApprovalPolicy,
) error {
	return s.store.CreateApprovalPolicy(ctx, tx, documentId, pol)
}

func (s *Server) UpsertApprovalPolicy(
	ctx context.Context,
	req *pbdocuments.UpsertApprovalPolicyRequest,
) (*pbdocuments.UpsertApprovalPolicyResponse, error) {
	pol := req.GetPolicy()

	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", pol.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		pol.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS,
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
			mysql.DateTimeExp(mysql.CURRENT_TIMESTAMP()), // Initialize snapshot_date
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

	if err := s.recomputeApprovalPolicyTx(
		ctx,
		tx,
		policy.GetDocumentId(),
		policy.GetSnapshotDate(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	policy, err = s.getApprovalPolicy(ctx, s.db, condition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	meta, err := s.store.GetDocumentMeta(ctx, s.db, pol.GetDocumentId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.UpsertApprovalPolicyResponse{
		Policy:  policy,
		DocMeta: meta,
	}, nil
}

func (s *Server) ListApprovalTasks(
	ctx context.Context,
	req *pbdocuments.ListApprovalTasksRequest,
) (*pbdocuments.ListApprovalTasksResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
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
		Tasks: []*documentsapproval.ApprovalTask{},
	}

	tUser := table.FivenetUser.AS("usershort")

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

	for i := range resp.Tasks {
		if resp.Tasks[i].GetJob() == "" {
			continue
		}

		s.enricher.EnrichJobInfo(resp.Tasks[i])
	}

	return resp, nil
}

func (s *Server) getApprovalTask(
	ctx context.Context,
	tx qrm.DB,
	taskId int64,
) (*documentsapproval.ApprovalTask, error) {
	task, err := s.store.GetApprovalTask(ctx, tx, taskId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	return task, nil
}

func (s *Server) UpsertApprovalTasks(
	ctx context.Context,
	req *pbdocuments.UpsertApprovalTasksRequest,
) (*pbdocuments.UpsertApprovalTasksResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.document_id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access: must be allowed to edit the document to seed tasks
	ok, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
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

	pol, err := s.getOrCreateApprovalPolicy(ctx, s.db, req.GetDocumentId(), docSnap)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
		docSnap,
		req.GetSeeds(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, req.GetDocumentId(), docSnap); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()
	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   req.GetDocumentId(),
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_ASSIGNED,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.UpsertApprovalTasksResponse{
		TasksCreated: created,
		TasksEnsured: ensured,
		Policy:       pol,
	}, nil
}

func (s *Server) createApprovalTasks(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
	documentId int64,
	snapDate *timestamp.Timestamp,
	seeds []*pbdocuments.ApprovalTaskSeed,
) (int32, int32, error) {
	return s.store.CreateApprovalTasks(ctx, tx, userInfo, documentId, snapDate, seeds)
}

func (s *Server) DeleteApprovalTasks(
	ctx context.Context,
	req *pbdocuments.DeleteApprovalTasksRequest,
) (*pbdocuments.DeleteApprovalTasksResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.document_id", req.GetDocumentId()})

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	// Resolve policy & snapshot
	var pol documentsapproval.ApprovalPolicy
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
	ok, err := s.canUserAccessDocument(
		ctx,
		pol.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.store.DeleteApprovalTasks(
		ctx,
		tx,
		pol.GetDocumentId(),
		pol.GetSnapshotDate(),
		req.GetDeleteAllPending(),
		req.GetTaskIds(),
		pol.GetPendingCount(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeApprovalPolicyTx(
		ctx,
		tx,
		pol.DocumentId,
		pol.GetSnapshotDate(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteApprovalTasksResponse{}, nil
}

func (s *Server) canUserAccessApprovalTask(
	ctx context.Context,
	tx qrm.DB,
	taskId int64,
	userInfo *userinfo.UserInfo,
) error {
	task, err := s.getApprovalTask(ctx, tx, taskId)
	if task == nil || err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	check, err := s.canUserAccessDocument(
		ctx,
		task.DocumentId,
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
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
	ok, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	_, approvals, err := s.store.ListApprovals(ctx, documentsstore.ListApprovalsQuery{
		DocumentID:   req.GetDocumentId(),
		SnapshotDate: req.GetSnapshotDate(),
		Status:       req.GetStatus(),
		UserID:       req.GetUserId(),
		TaskID:       req.GetTaskId(),
		UserInfo:     userInfo,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	resp := &pbdocuments.ListApprovalsResponse{Approvals: approvals}

	// Enrich labels/grades for display
	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for _, t := range resp.GetApprovals() {
		if t.GetUserJob() != "" {
			jobInfoFn(t)
		}

		if t.GetUser() != nil {
			jobInfoFn(t.GetUser())
		}
	}

	return resp, nil
}

func (s *Server) RevokeApproval(
	ctx context.Context,
	req *pbdocuments.RevokeApprovalRequest,
) (*pbdocuments.RevokeApprovalResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Load the approval artifact
	apr, err := s.store.GetApproval(ctx, s.db, req.GetApprovalId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if apr == nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	// Access: require EDIT on the document
	ok, err := s.canUserAccessDocument(
		ctx,
		apr.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	// Resolve policy using document_id + snapshot_date
	pol, err := s.store.GetApprovalPolicy(ctx, s.db, mysql.AND(
		table.FivenetDocumentsApprovalPolicies.AS(
			"approval_policy",
		).DocumentID.EQ(
			mysql.Int64(apr.GetDocumentId()),
		),
		table.FivenetDocumentsApprovalPolicies.AS(
			"approval_policy",
		).SnapshotDate.EQ(
			mysql.DateTimeT(apr.GetSnapshotDate().AsTime()),
		),
	))
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	tApprovals := table.FivenetDocumentsApprovals
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
			tApprovals.Status.SET(mysql.Int32(int32(documentsapproval.ApprovalStatus_APPROVAL_STATUS_REVOKED))),
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
				tApprovalTasks.Status.SET(mysql.Int32(int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
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
	if err := s.recomputeApprovalPolicyTx(
		ctx,
		tx,
		pol.GetDocumentId(),
		pol.GetSnapshotDate(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Create document activity entry
	comment := req.GetComment()
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   pol.GetDocumentId(),
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_REVOKED,
		Reason:       &comment,
		CreatorId:    &userId,
		CreatorJob:   userJob,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	meta, err := s.store.GetDocumentMeta(ctx, tx, pol.GetDocumentId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.RevokeApprovalResponse{
		Approval: apr,
		DocMeta:  meta,
	}, nil
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
	ok, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
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

	// If signature is required, must have either stamp_id or payload_svg. Declining doesn't require signature, so allow empty if declining.
	if pol.GetSignatureRequired() && req.GetStampId() == 0 && req.GetPayloadSvg() == "" &&
		req.GetNewStatus() != documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED {
		return nil, errorsdocuments.ErrApprovalSignatureRequired
	}

	snapDate := pol.GetSnapshotDate()

	tApprovalTasks := table.FivenetDocumentsApprovalTasks.AS("approval_task")
	tApprovals := table.FivenetDocumentsApprovals

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	now := time.Now()

	var decidedTask *documentsapproval.ApprovalTask // may remain nil for ad-hoc
	var taskIDForArtifact int64                     // 0 if ad-hoc

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
			decidedTask.GetSnapshotDate().AsTime() != snapDate.AsTime() {
			return nil, errorsdocuments.ErrDocAccessViewDenied
		}
		if decidedTask.GetStatus() != documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING {
			return nil, errorsdocuments.ErrApprovalTaskAlreadyHandled
		}

		if decidedTask.GetAssigneeKind() == documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER &&
			decidedTask.GetUserId() != userInfo.GetUserId() {
			return nil, errorsdocuments.ErrDocAccessViewDenied
		} else if decidedTask.GetAssigneeKind() == documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE {
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
		var candidate documentsapproval.ApprovalTask
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
			WHERE(mysql.AND(
				tApprovalTasks.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
				tApprovalTasks.SnapshotDate.EQ(dbutils.TimestampToMySQLDateTimeSec(snapDate)),
				tApprovalTasks.AssigneeKind.EQ(
					mysql.Int32(
						int32(documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER),
					),
				),
				tApprovalTasks.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				tApprovalTasks.Status.IN(
					mysql.Int32(
						int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
					),
					mysql.Int32(
						int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED),
					),
					mysql.Int32(
						int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_CANCELLED),
					),
				),
			)).
			ORDER_BY(tApprovalTasks.SlotNo.ASC(), tApprovalTasks.CreatedAt.ASC()).
			LIMIT(1)

		err := stmt.QueryContext(ctx, tx, &candidate)
		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		useCandidate := (err == nil && candidate.Id > 0)

		// If no USER task, try JOB target where user is eligible
		if !useCandidate {
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
				WHERE(mysql.AND(
					tApprovalTasks.DocumentID.EQ(mysql.Int64(pol.GetDocumentId())),
					tApprovalTasks.SnapshotDate.EQ(dbutils.TimestampToMySQLDateTimeSec(snapDate)),
					tApprovalTasks.AssigneeKind.EQ(
						mysql.Int32(
							int32(
								documentsapproval.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE,
							),
						),
					),
					tApprovalTasks.Job.EQ(mysql.String(userInfo.GetJob())),
					tApprovalTasks.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					tApprovalTasks.Status.IN(
						mysql.Int32(
							int32(
								documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING,
							),
						),
						mysql.Int32(
							int32(
								documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED,
							),
						),
						mysql.Int32(
							int32(
								documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_CANCELLED,
							),
						),
					),
				)).
				ORDER_BY(tApprovalTasks.SlotNo.ASC(), tApprovalTasks.CreatedAt.ASC()).
				LIMIT(1)

			err = stmt.QueryContext(ctx, tx, &candidate)
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
					tApprovalTasks.CompletedAt.SET(mysql.DateTimeT(now)),
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

	// If signature is required for the task, must have either stamp_id or payload_svg. Declining doesn't require signature, so allow empty if declining.
	if decidedTask.GetSignatureRequired() && req.GetStampId() == 0 && req.GetPayloadSvg() == "" &&
		req.GetNewStatus() != documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED {
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
			tApprovals.SnapshotDate.EQ(dbutils.TimestampToMySQLDateTimeSec(snapDate)),
			tApprovals.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
		)).
		LIMIT(1).
		QueryContext(ctx, tx, &existing); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Write/UPSERT approval artifact (unique per (document_id, snapshot_date, user_id))
	// Map task statuses APPROVED/DECLINED -> artifact status
	artifactStatus := documentsapproval.ApprovalStatus_APPROVAL_STATUS_APPROVED
	if req.GetNewStatus() == documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED ||
		req.GetNewStatus() == documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_CANCELLED {
		artifactStatus = documentsapproval.ApprovalStatus_APPROVAL_STATUS_DECLINED
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
				tApprovals.RevokedAt,
			).
			VALUES(
				pol.GetDocumentId(),
				dbutils.TimestampToMySQLDateTimeSec(snapDate),
				userInfo.GetUserId(),
				userInfo.GetJob(),
				mysql.Int32(userInfo.GetJobGrade()),
				req.GetPayloadSvg(),
				dbutils.Int64P(req.GetStampId()),
				int32(artifactStatus),
				req.GetComment(),
				dbutils.Int64P(taskIDForArtifact),
				nil,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tApprovals.SnapshotDate.SET(mysql.RawTimestamp("VALUES(`snapshot_date`)")),
				tApprovals.UserID.SET(mysql.Int32(userInfo.GetUserId())),
				tApprovals.UserJob.SET(mysql.String(userInfo.GetJob())),
				tApprovals.UserJobGrade.SET(mysql.Int32(userInfo.GetJobGrade())),
				tApprovals.PayloadSvg.SET(mysql.String(req.GetPayloadSvg())),
				tApprovals.StampID.SET(dbutils.Int64P(req.GetStampId())),
				tApprovals.Status.SET(mysql.Int32(int32(artifactStatus))),
				tApprovals.Comment.SET(mysql.String(req.GetComment())),
				tApprovals.TaskID.SET(dbutils.Int64P(taskIDForArtifact)),
				tApprovals.RevokedAt.SET(mysql.DateTimeExp(mysql.NULL)),
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
				tApprovals.RevokedAt,
			).
			SET(
				snapDate,
				userInfo.GetUserId(),
				userInfo.GetJob(),
				mysql.Int32(userInfo.GetJobGrade()),
				req.GetPayloadSvg(),
				dbutils.Int64P(req.GetStampId()),
				artifactStatus,
				req.GetComment(),
				dbutils.Int64P(taskIDForArtifact),
				nil,
			).
			WHERE(tApprovals.ID.EQ(mysql.Int64(existing.ID))).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	tApprovals = table.FivenetDocumentsApprovals.AS("approval")

	var artifact documentsapproval.Approval
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
			tApprovals.SnapshotDate.EQ(dbutils.TimestampToMySQLDateTimeSec(snapDate)),
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
	if err := s.recomputeApprovalPolicyTx(ctx, tx, pol.GetDocumentId(), snapDate); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Create document activity entry
	comment := req.GetComment()
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	activityType := documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_REJECTED
	if req.GetNewStatus() == documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_APPROVED {
		activityType = documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_APPROVED
	}

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
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

	meta, err := s.store.GetDocumentMeta(ctx, s.db, pol.GetDocumentId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DecideApprovalResponse{
		Approval: &artifact,
		Task:     decidedTask,
		Policy:   pol,
		DocMeta:  meta,
	}, nil
}

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
			tApprovalTasks.Status.SET(mysql.Int32(int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
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
		SnapshotDate *timestamp.Timestamp
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

	if err := s.recomputeApprovalPolicyTx(ctx, tx, k.DocumentID, k.SnapshotDate); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Create document activity entry
	comment := req.GetComment()
	userId := userInfo.GetUserId()
	userJob := userInfo.GetJob()

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   k.DocumentID,
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_APPROVAL_REVOKED,
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
		Task: &documentsapproval.ApprovalTask{
			Id:     req.GetTaskId(),
			Status: documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING,
		},
		Policy: &documentsapproval.ApprovalPolicy{
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

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
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

	if err := s.recomputeApprovalPolicyTx(
		ctx,
		s.db,
		req.GetDocumentId(),
		pol.GetSnapshotDate(),
	); err != nil {
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
	snapDate *timestamp.Timestamp,
) error {
	return s.store.RecomputeApprovalPolicyTx(ctx, tx, documentID, snapDate)
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

	if pol == nil ||
		pol.OnEditBehavior <= documentsapproval.OnEditBehavior_ON_EDIT_BEHAVIOR_KEEP_PROGRESS {
		// Nothing to do
		return nil
	}

	if err := s.store.ResetApprovalProgress(ctx, tx, doc.GetId()); err != nil {
		return err
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, doc.GetId(), pol.GetSnapshotDate()); err != nil {
		return err
	}

	return nil
}

func (s *Server) expireApprovalTasks(ctx context.Context) (int64, error) {
	affected, err := s.store.ExpireApprovalTasks(ctx, s.db)
	if err != nil {
		return 0, err
	}

	return affected, nil
}
