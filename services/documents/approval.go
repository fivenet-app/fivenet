package documents

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// GetApprovalPolicy.
func (s *Server) GetApprovalPolicy(
	ctx context.Context,
	req *pbdocuments.GetApprovalPolicyRequest,
) (*pbdocuments.GetApprovalPolicyResponse, error) {
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

	policy, err := s.getApprovalPolicy(ctx, condition)
	if err != nil {
		return nil, err
	}

	return &pbdocuments.GetApprovalPolicyResponse{
		Policy: policy,
	}, nil
}

func (s *Server) getApprovalPolicy(
	ctx context.Context, condition mysql.BoolExpression,
) (*documents.ApprovalPolicy, error) {
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	stmt := tApprovalPolicy.
		SELECT(
			tApprovalPolicy.ID,
			tApprovalPolicy.DocumentID,
			tApprovalPolicy.RuleKind,
			tApprovalPolicy.RequiredCount,
			tApprovalPolicy.SnapshotDate,
			tApprovalPolicy.DueAt,
			tApprovalPolicy.StartedAt,
			tApprovalPolicy.CompletedAt,
			tApprovalPolicy.OnEditBehavior,
			tApprovalPolicy.AssignedCount,
			tApprovalPolicy.ApprovedCount,
			tApprovalPolicy.DeclinedCount,
			tApprovalPolicy.PendingCount,
			tApprovalPolicy.AnyDeclined,
			tApprovalPolicy.CreatedAt,
			tApprovalPolicy.UpdatedAt,
			tApprovalPolicy.DeletedAt,
		).
		FROM(tApprovalPolicy).
		WHERE(condition).
		LIMIT(1)

	policy := &documents.ApprovalPolicy{}
	if err := stmt.QueryContext(ctx, s.db, policy); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if policy.Id == 0 {
		return nil, nil
	}

	return policy, nil
}

// UpsertApprovalPolicy.
func (s *Server) UpsertApprovalPolicy(
	ctx context.Context,
	req *pbdocuments.UpsertApprovalPolicyRequest,
) (*pbdocuments.UpsertApprovalPolicyResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
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

	stmt := tApprovalPolicy.
		INSERT(
			tApprovalPolicy.DocumentID,
			tApprovalPolicy.RuleKind,
			tApprovalPolicy.RequiredCount,
			tApprovalPolicy.DueAt,
			tApprovalPolicy.OnEditBehavior,
			tApprovalPolicy.SnapshotDate,
		).
		VALUES(
			req.GetDocumentId(),
			int32(req.GetRuleKind()),
			req.GetRequiredCount(),
			req.DueAt,
			int32(req.GetOnEditBehavior()),
			mysql.CURRENT_TIMESTAMP(), // Initialize snapshot_date
		).
		ON_DUPLICATE_KEY_UPDATE(
			tApprovalPolicy.RuleKind.SET(mysql.Int32(int32(req.GetRuleKind()))),
			tApprovalPolicy.RequiredCount.SET(mysql.Int32(req.GetRequiredCount())),
			tApprovalPolicy.DueAt.SET(mysql.TimestampT(req.DueAt.AsTime())),
			tApprovalPolicy.OnEditBehavior.SET(mysql.Int32(int32(req.GetOnEditBehavior()))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	policy, err := s.getApprovalPolicy(
		ctx,
		tApprovalPolicy.AS("approval_policy").DocumentID.EQ(mysql.Int64(req.GetDocumentId())),
	)
	if err != nil {
		return nil, err
	}

	return &pbdocuments.UpsertApprovalPolicyResponse{
		Policy: policy,
	}, nil
}

func (s *Server) StartApprovalRound(
	ctx context.Context,
	req *pbdocuments.StartApprovalRoundRequest,
) (*pbdocuments.StartApprovalRoundResponse, error) {
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

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	snap := time.Now()

	// Reset policy counters for new round
	if _, err = tApprovalPolicy.
		UPDATE().
		SET(
			tApprovalPolicy.SnapshotDate.SET(mysql.TimestampT(snap)),
			tApprovalPolicy.StartedAt.SET(mysql.TimestampT(time.Now())),
			tApprovalPolicy.CompletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tApprovalPolicy.AssignedCount.SET(mysql.Int(0)),
			tApprovalPolicy.ApprovedCount.SET(mysql.Int(0)),
			tApprovalPolicy.DeclinedCount.SET(mysql.Int(0)),
			tApprovalPolicy.PendingCount.SET(mysql.Int(0)),
			tApprovalPolicy.AnyDeclined.SET(mysql.Bool(false)),
		).
		WHERE(tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Regenerate tasks - Delete pending for same (doc, snap) if any
	_, _ = tApprovalTasks.
		DELETE().
		WHERE(
			tApprovalTasks.DocumentID.EQ(mysql.Int64(req.GetDocumentId())).
				AND(tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap))).
				AND(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))),
		).
		ExecContext(ctx, tx)

	// Fetch ACL entries
	type accRow struct {
		Job          *string
		MinimumGrade *int32
	}
	var accRows []accRow
	if err = tApprovalTasks.
		SELECT(
			tApprovalTasks.Job,
			tApprovalTasks.MinimumGrade,
		).
		FROM(tApprovalTasks.
			INNER_JOIN(tApprovalPolicy,
				tApprovalPolicy.ID.EQ(tApprovalTasks.DocumentID),
			),
		).
		WHERE(
			tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId())).
				AND(tApprovalTasks.Status.LT_EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))),
		).
		QueryContext(ctx, tx, &accRows); err != nil {
		return nil, err
	}

	ins := tApprovalTasks.
		INSERT(
			tApprovalTasks.DocumentID,
			tApprovalTasks.SnapshotDate,
			tApprovalTasks.AssigneeKind,
			tApprovalTasks.Job,
			tApprovalTasks.MinimumGrade,
			tApprovalTasks.Status,
		)

	for _, a := range accRows {
		job := ""
		if a.Job != nil {
			job = *a.Job
		}

		mg := int32(-1)
		if a.MinimumGrade != nil {
			mg = *a.MinimumGrade
		}

		ins = ins.
			VALUES(
				req.GetDocumentId(),
				snap,
				int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE),
				job,
				mg,
				int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
			)
	}

	if len(accRows) > 0 {
		if _, err = ins.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Minimal response (fetch full tasks if you like)
	return &pbdocuments.StartApprovalRoundResponse{
		Policy: &documents.ApprovalPolicy{
			DocumentId: req.GetDocumentId(),
		},
		Tasks: nil,
	}, nil
}

// CompleteApprovalRound.
func (s *Server) CompleteApprovalRound(
	ctx context.Context,
	req *pbdocuments.CompleteApprovalRoundRequest,
) (*pbdocuments.CompleteApprovalRoundResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies
	now := time.Now()

	if _, err := tApprovalPolicy.
		UPDATE().
		SET(
			tApprovalPolicy.CompletedAt.SET(mysql.TimestampT(now)),
		).
		WHERE(tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.CompleteApprovalRoundResponse{
		Policy: &documents.ApprovalPolicy{DocumentId: req.GetDocumentId()},
	}, nil
}

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
			tApprovalTasks.Status,
			tApprovalTasks.Comment,
			tApprovalTasks.CreatedAt,
			tApprovalTasks.DecidedAt,
			tApprovalTasks.DueAt,
			tApprovalTasks.DecisionCount,
			tUser.ID,
			tUser.Firstname,
			tUser.Lastname,
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
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
			tApprovalTasks.Status,
			tApprovalTasks.Comment,
			tApprovalTasks.CreatedAt,
			tApprovalTasks.DecidedAt,
			tApprovalTasks.DueAt,
			tApprovalTasks.DecisionCount,
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

func (s *Server) canUserAccessTask(
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

func (s *Server) DecideApprovalTask(
	ctx context.Context,
	req *pbdocuments.DecideApprovalTaskRequest,
) (*pbdocuments.DecideApprovalTaskResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.task_id", req.GetTaskId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.canUserAccessTask(ctx, s.db, req.GetTaskId(), userInfo); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	tDApprovalPolicy := table.FivenetDocumentsApprovalPolicies

	// 1) Update task atomically (allow reconsideration: +1 decision_count)
	up := tApprovalTasks.
		UPDATE().
		SET(
			tApprovalTasks.Status.SET(mysql.Int32(int32(req.GetNewStatus()))),
			tApprovalTasks.DecisionCount.SET(tApprovalTasks.DecisionCount.ADD(mysql.Int32(1))),
			tApprovalTasks.DecidedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tApprovalTasks.Comment.SET(mysql.String(req.GetComment())),
		).
		WHERE(tApprovalTasks.ID.EQ(mysql.Int64(req.GetTaskId()))).
		LIMIT(1)

	if _, err = up.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// 2) Fetch doc/snapshot for aggregation
	type keyRow struct {
		DocumentID   int64
		SnapshotDate time.Time
	}
	var k keyRow
	if err = tApprovalTasks.
		SELECT(
			tApprovalTasks.DocumentID,
			tApprovalTasks.SnapshotDate,
		).
		FROM(tApprovalTasks).
		WHERE(tApprovalTasks.ID.EQ(mysql.Int64(req.GetTaskId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &k); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// 3) Aggregate counters
	var agg struct {
		Assigned int64
		Approved int64
		Declined int64
		Pending  int64
	}
	if err = mysql.
		SELECT(
			mysql.COUNT(tApprovalTasks.ID),
			mysql.SUM(
				mysql.CASE().
					WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_APPROVED)))).
					THEN(mysql.Int32(1)).ELSE(mysql.Int32(0)),
			),
			mysql.SUM(
				mysql.CASE().
					WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED)))).
					THEN(mysql.Int32(1)).ELSE(mysql.Int32(0)),
			),
			mysql.SUM(
				mysql.CASE().
					WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))).
					THEN(mysql.Int32(1)).ELSE(mysql.Int32(0)),
			),
		).
		FROM(tApprovalTasks).
		WHERE(mysql.AND(
			tApprovalTasks.DocumentID.EQ(mysql.Int64(k.DocumentID)),
			tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(k.SnapshotDate)),
		)).
		QueryContext(ctx, tx, &agg); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// 4) Update policy counters
	if _, err = tDApprovalPolicy.
		UPDATE().
		SET(
			tDApprovalPolicy.AssignedCount.SET(mysql.Int32(int32(agg.Assigned))),
			tDApprovalPolicy.ApprovedCount.SET(mysql.Int32(int32(agg.Approved))),
			tDApprovalPolicy.DeclinedCount.SET(mysql.Int32(int32(agg.Declined))),
			tDApprovalPolicy.PendingCount.SET(mysql.Int32(int32(agg.Pending))),
			tDApprovalPolicy.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
		).
		WHERE(tDApprovalPolicy.DocumentID.EQ(mysql.Int64(k.DocumentID))).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err = tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Minimal response; fetchers can be added if you want full objects
	return &pbdocuments.DecideApprovalTaskResponse{
		Task: &documents.ApprovalTask{
			Id:     int64(req.GetTaskId()),
			Status: req.GetNewStatus(),
		},
		Policy: &documents.ApprovalPolicy{
			DocumentId: int64(k.DocumentID),
		},
	}, nil
}

// ReopenApprovalTask.
func (s *Server) ReopenApprovalTask(
	ctx context.Context,
	req *pbdocuments.ReopenApprovalTaskRequest,
) (*pbdocuments.ReopenApprovalTaskResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.task_id", req.GetTaskId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.canUserAccessTask(ctx, s.db, req.GetTaskId(), userInfo); err != nil {
		return nil, err
	}

	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies

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
			tApprovalTasks.DecidedAt.SET(mysql.TimestampExp(mysql.NULL)),
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
		DocumentID   uint64
		SnapshotDate time.Time
	}
	if err = mysql.
		SELECT(tApprovalTasks.DocumentID, tApprovalTasks.SnapshotDate).
		FROM(tApprovalTasks).
		WHERE(tApprovalTasks.ID.EQ(mysql.Int(req.GetTaskId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &k); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Recompute counters
	var agg struct{ Assigned, Approved, Declined, Pending int64 }
	if err = tApprovalTasks.
		SELECT(
			mysql.COUNT(tApprovalTasks.ID),
			mysql.SUM(mysql.CASE().WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_APPROVED)))).THEN(mysql.Int(1)).ELSE(mysql.Int(0))),
			mysql.SUM(mysql.CASE().WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED)))).THEN(mysql.Int(1)).ELSE(mysql.Int(0))),
			mysql.SUM(mysql.CASE().WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))).THEN(mysql.Int(1)).ELSE(mysql.Int(0))),
		).FROM(tApprovalTasks).
		WHERE(
			tApprovalTasks.DocumentID.EQ(mysql.Int(int64(k.DocumentID))).
				AND(tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(k.SnapshotDate))),
		).QueryContext(ctx, tx, &agg); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err = tApprovalPolicy.
		UPDATE().
		SET(
			tApprovalPolicy.AssignedCount.SET(mysql.Int32(int32(agg.Assigned))),
			tApprovalPolicy.ApprovedCount.SET(mysql.Int32(int32(agg.Approved))),
			tApprovalPolicy.DeclinedCount.SET(mysql.Int32(int32(agg.Declined))),
			tApprovalPolicy.PendingCount.SET(mysql.Int32(int32(agg.Pending))),
			tApprovalPolicy.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
		).
		WHERE(tApprovalPolicy.DocumentID.EQ(mysql.Int(int64(k.DocumentID)))).
		ExecContext(ctx, tx); err != nil {
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
		Policy: &documents.ApprovalPolicy{DocumentId: int64(k.DocumentID)},
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

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	// Get active snapshot
	var snap time.Time
	if err := tApprovalPolicy.
		SELECT(
			tApprovalPolicy.SnapshotDate,
		).
		FROM(tApprovalPolicy).
		WHERE(tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &snap); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	var agg struct {
		Assigned int64
		Approved int64
		Declined int64
		Pending  int64
	}
	if err := mysql.
		SELECT(
			mysql.COUNT(tApprovalTasks.ID),
			mysql.SUM(
				mysql.CASE().WHEN(
					tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_APPROVED)))).
					THEN(mysql.Int(1)).
					ELSE(mysql.Int(0)),
			),
			mysql.SUM(
				mysql.CASE().WHEN(
					tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED)))).
					THEN(mysql.Int(1)).
					ELSE(mysql.Int(0)),
			),
			mysql.SUM(
				mysql.CASE().WHEN(
					tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))).
					THEN(mysql.Int(1)).
					ELSE(mysql.Int(0)),
			),
		).
		FROM(tApprovalTasks).
		WHERE(mysql.AND(
			tApprovalTasks.DocumentID.EQ(mysql.Int64(req.GetDocumentId())),
			tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
		)).
		QueryContext(ctx, s.db, &agg); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := tApprovalPolicy.
		UPDATE().
		SET(
			tApprovalPolicy.AssignedCount.SET(mysql.Int32(int32(agg.Assigned))),
			tApprovalPolicy.ApprovedCount.SET(mysql.Int32(int32(agg.Approved))),
			tApprovalPolicy.DeclinedCount.SET(mysql.Int32(int32(agg.Declined))),
			tApprovalPolicy.PendingCount.SET(mysql.Int32(int32(agg.Pending))),
			tApprovalPolicy.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
		).
		WHERE(tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.RecomputeApprovalPolicyCountersResponse{
		Policy: &documents.ApprovalPolicy{
			DocumentId: req.GetDocumentId(),
		},
	}, nil
}
