package documents

import (
	"context"
	"database/sql"
	"errors"
	"time"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) ListApprovalAccess(
	ctx context.Context,
	req *pbdocuments.ListApprovalAccessRequest,
) (*pbdocuments.ListApprovalAccessResponse, error) {
	tAccess := table.FivenetDocumentsApprovalAccess

	condition := tAccess.TargetID.EQ(mysql.Int64(req.GetPolicyId()))

	resp := &pbdocuments.ListApprovalAccessResponse{}

	stmt := mysql.
		SELECT(
			tAccess.ID,
			tAccess.TargetID,
			tAccess.Job,
			tAccess.MinimumGrade,
			tAccess.Access,
		).
		FROM(tAccess).
		WHERE(condition).
		ORDER_BY(tAccess.ID.ASC()).
		LIMIT(40)

	if err := stmt.QueryContext(ctx, s.db, &resp.Access); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) ListTasks(
	ctx context.Context,
	req *pbdocuments.ListTasksRequest,
) (*pbdocuments.ListTasksResponse, error) {
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	condition := tApprovalTasks.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))

	if len(req.GetStatuses()) > 0 {
		vals := make([]mysql.Expression, 0, len(req.GetStatuses()))
		for _, st := range req.GetStatuses() {
			vals = append(vals, mysql.Int32(int32(st)))
		}
		condition = condition.AND(tApprovalTasks.Status.IN(vals...))
	}

	countStmt := mysql.
		SELECT(mysql.COUNT(tApprovalTasks.ID)).
		FROM(tApprovalTasks).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	pag, limit := req.GetPagination().GetResponse(DocsDefaultPageSize)
	resp := &pbdocuments.ListTasksResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	stmt := mysql.
		SELECT(
			tApprovalTasks.ID,
			tApprovalTasks.DocumentID,
			tApprovalTasks.SnapshotDate,
			tApprovalTasks.AssigneeKind,
			tApprovalTasks.UserID,
			tApprovalTasks.Job,
			tApprovalTasks.MinimumGrade,
			tApprovalTasks.DecidedByUserID,
			tApprovalTasks.DecidedByJob,
			tApprovalTasks.DecidedByUserGrade,
			tApprovalTasks.Status,
			tApprovalTasks.Comment,
			tApprovalTasks.CreatedAt,
			tApprovalTasks.DecidedAt,
			tApprovalTasks.DueAt,
			tApprovalTasks.DecisionCount,
		).
		FROM(tApprovalTasks).
		WHERE(condition).
		ORDER_BY(tApprovalTasks.CreatedAt.ASC()).
		OFFSET(pag.GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Tasks); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) DecideTask(
	ctx context.Context,
	req *pbdocuments.DecideTaskRequest,
) (*pbdocuments.DecideTaskResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	tPolicies := table.FivenetDocumentsApprovalPolicies

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
		return nil, err
	}

	// 2) Fetch doc/snapshot for aggregation
	type keyRow struct {
		DocumentID   int64
		SnapshotDate time.Time
	}
	var k keyRow
	if err = tApprovalTasks.
		SELECT(
			tApprovalTasks.DocumentID, tApprovalTasks.SnapshotDate,
		).
		FROM(tApprovalTasks).
		WHERE(tApprovalTasks.ID.EQ(mysql.Int64(req.GetTaskId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &k); err != nil {
		return nil, err
	}

	// 3) Aggregate counters
	var agg struct {
		Assigned int64
		Approved int64
		Declined int64
		Pending  int64
	}
	if err = mysql.SELECT(
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
		WHERE(
			tApprovalTasks.DocumentID.EQ(mysql.Int64(k.DocumentID)).
				AND(tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(k.SnapshotDate))),
		).
		QueryContext(ctx, tx, &agg); err != nil {
		return nil, err
	}

	// 4) Update policy counters
	if _, err = tPolicies.
		UPDATE().
		SET(
			tPolicies.AssignedCount.SET(mysql.Int32(int32(agg.Assigned))),
			tPolicies.ApprovedCount.SET(mysql.Int32(int32(agg.Approved))),
			tPolicies.DeclinedCount.SET(mysql.Int32(int32(agg.Declined))),
			tPolicies.PendingCount.SET(mysql.Int32(int32(agg.Pending))),
			tPolicies.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
			tPolicies.UpdatedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(tPolicies.DocumentID.EQ(mysql.Int64(k.DocumentID))).
		ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Minimal response; fetchers can be added if you want full objects
	return &pbdocuments.DecideTaskResponse{
		Task:   &documents.ApprovalTask{Id: int64(req.GetTaskId()), Status: req.GetNewStatus()},
		Policy: &documents.ApprovalPolicy{DocumentId: int64(k.DocumentID)},
	}, nil
}

// GetPolicy.
func (s *Server) GetPolicy(
	ctx context.Context,
	req *pbdocuments.GetPolicyRequest,
) (*pbdocuments.GetPolicyResponse, error) {
	tApprovalPolicies := table.FivenetDocumentsApprovalPolicies

	stmt := tApprovalPolicies.
		SELECT(
			tApprovalPolicies.ID,
			tApprovalPolicies.DocumentID,
			tApprovalPolicies.RuleKind,
			tApprovalPolicies.RequiredCount,
			tApprovalPolicies.ActiveSnapshotDate,
			tApprovalPolicies.DueAt,
			tApprovalPolicies.StartedAt,
			tApprovalPolicies.CompletedAt,
			tApprovalPolicies.OnEditBehavior,
			tApprovalPolicies.AssignedCount,
			tApprovalPolicies.ApprovedCount,
			tApprovalPolicies.DeclinedCount,
			tApprovalPolicies.PendingCount,
			tApprovalPolicies.AnyDeclined,
			tApprovalPolicies.CreatedAt,
			tApprovalPolicies.UpdatedAt,
		).
		FROM(tApprovalPolicies).
		WHERE(tApprovalPolicies.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		LIMIT(1)

	type row struct {
		ID, DocumentID                uint64
		RuleKind, RequiredCount       int32
		ActiveSnapshotDate            time.Time
		DueAt, StartedAt, CompletedAt *time.Time
		OnEditBehavior                int32
		AssignedCount, ApprovedCount  int32
		DeclinedCount, PendingCount   int32
		AnyDeclined                   bool
		CreatedAt, UpdatedAt          time.Time
	}
	var r row
	if err := stmt.QueryContext(ctx, s.db, &r); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pbdocuments.GetPolicyResponse{}, nil
		}
		return nil, err
	}

	out := &documents.ApprovalPolicy{
		Id:            int64(r.ID),
		DocumentId:    int64(r.DocumentID),
		RuleKind:      documents.ApprovalRuleKind(r.RuleKind),
		RequiredCount: r.RequiredCount,
		// TODO: map timestamps if needed
		OnEditBehavior: documents.OnEditBehavior(r.OnEditBehavior),
		AssignedCount:  r.AssignedCount,
		ApprovedCount:  r.ApprovedCount,
		DeclinedCount:  r.DeclinedCount,
		PendingCount:   r.PendingCount,
		AnyDeclined:    r.AnyDeclined,
	}
	return &pbdocuments.GetPolicyResponse{Policy: out}, nil
}

// UpsertPolicy.
func (s *Server) UpsertPolicy(
	ctx context.Context,
	req *pbdocuments.UpsertPolicyRequest,
) (*pbdocuments.UpsertPolicyResponse, error) {
	tApprovalPolicies := table.FivenetDocumentsApprovalPolicies
	now := time.Now()

	stmt := tApprovalPolicies.
		INSERT().
		VALUES(
			tApprovalPolicies.DocumentID,
			tApprovalPolicies.RuleKind,
			tApprovalPolicies.RequiredCount,
			tApprovalPolicies.DueAt,
			tApprovalPolicies.OnEditBehavior,
			tApprovalPolicies.ActiveSnapshotDate,
			tApprovalPolicies.CreatedAt,
			tApprovalPolicies.UpdatedAt,
		).
		VALUES(
			req.GetDocumentId(),
			int32(req.GetRuleKind()),
			req.GetRequiredCount(),
			// TODO: convert req.DueAt to time.Time if present
			nil,
			int32(req.GetOnEditBehavior()),
			now, // initialize active_snapshot_date if you want
			now,
			now,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tApprovalPolicies.RuleKind.SET(mysql.Int32(int32(req.GetRuleKind()))),
			tApprovalPolicies.RequiredCount.SET(mysql.Int32(req.GetRequiredCount())),
			tApprovalPolicies.DueAt.SET(mysql.TimestampT(req.DueAt.AsTime())),
			tApprovalPolicies.OnEditBehavior.SET(mysql.Int32(int32(req.GetOnEditBehavior()))),
			tApprovalPolicies.UpdatedAt.SET(mysql.TimestampT(now)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	policy, err := s.GetPolicy(ctx, &pbdocuments.GetPolicyRequest{DocumentId: req.GetDocumentId()})
	if err != nil {
		return nil, err
	}

	return &pbdocuments.UpsertPolicyResponse{
		Policy: policy.Policy,
	}, nil
}

func (s *Server) StartApprovalRound(
	ctx context.Context,
	req *pbdocuments.StartApprovalRoundRequest,
) (*pbdocuments.StartApprovalRoundResponse, error) {
	tApprovalPolicies := table.FivenetDocumentsApprovalPolicies
	tApprovalAccess := table.FivenetDocumentsApprovalAccess
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	snap := time.Now()

	// 1) Reset policy counters for new round
	if _, err = tApprovalPolicies.UPDATE().
		SET(
			tApprovalPolicies.ActiveSnapshotDate.SET(mysql.TimestampT(snap)),
			tApprovalPolicies.StartedAt.SET(mysql.TimestampT(time.Now())),
			tApprovalPolicies.CompletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tApprovalPolicies.AssignedCount.SET(mysql.Int(0)),
			tApprovalPolicies.ApprovedCount.SET(mysql.Int(0)),
			tApprovalPolicies.DeclinedCount.SET(mysql.Int(0)),
			tApprovalPolicies.PendingCount.SET(mysql.Int(0)),
			tApprovalPolicies.AnyDeclined.SET(mysql.Bool(false)),
		).
		WHERE(tApprovalPolicies.DocumentID.EQ(mysql.Int(req.GetDocumentId()))).ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	// 2) regenerate tasks (example: from ACL USE entries)
	if req.GetRegenTasks() {
		// delete pending for same (doc,snap) if any
		_, _ = tApprovalTasks.
			DELETE().
			WHERE(
				tApprovalTasks.DocumentID.EQ(mysql.Int(req.GetDocumentId())).
					AND(tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap))).
					AND(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))),
			).
			ExecContext(ctx, tx)

		// fetch ACL entries
		type accRow struct {
			Job          *string
			MinimumGrade *int32
		}
		var accRows []accRow
		if err = tApprovalAccess.
			SELECT(
				tApprovalAccess.Job,
				tApprovalAccess.MinimumGrade,
			).
			FROM(tApprovalAccess.
				INNER_JOIN(tApprovalPolicies,
					tApprovalPolicies.ID.EQ(tApprovalAccess.TargetID),
				),
			).
			WHERE(
				tApprovalPolicies.DocumentID.EQ(mysql.Int64(req.GetDocumentId())).
					AND(tApprovalAccess.Access.EQ(mysql.Int32(int32(documents.ApprovalAccessLevel_APPROVAL_ACCESS_LEVEL_USE)))),
			).QueryContext(ctx, tx, &accRows); err != nil {
			return nil, err
		}

		ins := tApprovalTasks.
			INSERT().
			VALUES(
				tApprovalTasks.DocumentID,
				tApprovalTasks.SnapshotDate,
				tApprovalTasks.AssigneeKind,
				tApprovalTasks.Job,
				tApprovalTasks.MinimumGrade,
				tApprovalTasks.Status,
				tApprovalTasks.CreatedAt,
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
			ins = ins.VALUES(
				req.GetDocumentId(),
				snap,
				int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE),
				job,
				mg,
				int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
				time.Now(),
			)
		}
		if len(accRows) > 0 {
			if _, err = ins.ExecContext(ctx, tx); err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// minimal response (fetch full tasks if you like)
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
	tApprovalPolicies := table.FivenetDocumentsApprovalPolicies
	now := time.Now()

	if _, err := tApprovalPolicies.UPDATE().
		SET(
			tApprovalPolicies.CompletedAt.SET(mysql.TimestampT(now)),
			tApprovalPolicies.UpdatedAt.SET(mysql.TimestampT(now)),
		).
		WHERE(tApprovalPolicies.DocumentID.EQ(mysql.Int(req.GetDocumentId()))).ExecContext(ctx, s.db); err != nil {
		return nil, err
	}
	return &pbdocuments.CompleteApprovalRoundResponse{
		Policy: &documents.ApprovalPolicy{DocumentId: req.GetDocumentId()},
	}, nil
}

// RecomputePolicyCounters (aggregate from tasks).
func (s *Server) RecomputePolicyCounters(
	ctx context.Context,
	req *pbdocuments.RecomputePolicyCountersRequest,
) (*pbdocuments.RecomputePolicyCountersResponse, error) {
	tApprovalPolicies := table.FivenetDocumentsApprovalPolicies
	tApprovalTasks := table.FivenetDocumentsApprovalTasks

	// get active snapshot
	var snap time.Time
	if err := tApprovalPolicies.SELECT(
		tApprovalPolicies.ActiveSnapshotDate,
	).
		FROM(tApprovalPolicies).
		WHERE(tApprovalPolicies.DocumentID.EQ(mysql.Int(req.GetDocumentId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &snap); err != nil {
		return nil, err
	}

	var agg struct {
		Assigned int64
		Approved int64
		Declined int64
		Pending  int64
	}
	if err := mysql.SELECT(
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
	).FROM(tApprovalTasks).
		WHERE(
			tApprovalTasks.DocumentID.EQ(mysql.Int(req.GetDocumentId())).
				AND(tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap))),
		).QueryContext(ctx, s.db, &agg); err != nil {
		return nil, err
	}

	if _, err := tApprovalPolicies.UPDATE().
		SET(
			tApprovalPolicies.AssignedCount.SET(mysql.Int32(int32(agg.Assigned))),
			tApprovalPolicies.ApprovedCount.SET(mysql.Int32(int32(agg.Approved))),
			tApprovalPolicies.DeclinedCount.SET(mysql.Int32(int32(agg.Declined))),
			tApprovalPolicies.PendingCount.SET(mysql.Int32(int32(agg.Pending))),
			tApprovalPolicies.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
			tApprovalPolicies.UpdatedAt.SET(mysql.TimestampT(time.Now())),
		).
		WHERE(tApprovalPolicies.DocumentID.EQ(mysql.Int(req.GetDocumentId()))).ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &pbdocuments.RecomputePolicyCountersResponse{
		Policy: &documents.ApprovalPolicy{DocumentId: req.GetDocumentId()},
	}, nil
}

// ReopenTask.
func (s *Server) ReopenTask(
	ctx context.Context,
	req *pbdocuments.ReopenTaskRequest,
) (*pbdocuments.ReopenTaskResponse, error) {
	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	tApprovalPolicies := table.FivenetDocumentsApprovalPolicies

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// set PENDING & clear decider snapshot
	if _, err = tApprovalTasks.
		UPDATE().
		SET(
			tApprovalTasks.Status.SET(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
			tApprovalTasks.DecidedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tApprovalTasks.DecidedByUserID.SET(mysql.IntExp(mysql.NULL)),
			tApprovalTasks.DecidedByJob.SET(mysql.StringExp(mysql.NULL)),
			tApprovalTasks.DecidedByUserGrade.SET(mysql.IntExp(mysql.NULL)),
		).
		WHERE(
			tApprovalTasks.ID.EQ(mysql.Int64(req.GetTaskId())),
		).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	// fetch key
	var k struct {
		DocumentID   uint64
		SnapshotDate time.Time
	}
	if err = mysql.SELECT(tApprovalTasks.DocumentID, tApprovalTasks.SnapshotDate).
		FROM(tApprovalTasks).
		WHERE(tApprovalTasks.ID.EQ(mysql.Int(req.GetTaskId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &k); err != nil {
		return nil, err
	}

	// recompute counters
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
		return nil, err
	}

	if _, err = tApprovalPolicies.UPDATE().
		SET(
			tApprovalPolicies.AssignedCount.SET(mysql.Int32(int32(agg.Assigned))),
			tApprovalPolicies.ApprovedCount.SET(mysql.Int32(int32(agg.Approved))),
			tApprovalPolicies.DeclinedCount.SET(mysql.Int32(int32(agg.Declined))),
			tApprovalPolicies.PendingCount.SET(mysql.Int32(int32(agg.Pending))),
			tApprovalPolicies.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
			tApprovalPolicies.UpdatedAt.SET(mysql.TimestampT(time.Now())),
		).WHERE(tApprovalPolicies.DocumentID.EQ(mysql.Int(int64(k.DocumentID)))).ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &pbdocuments.ReopenTaskResponse{
		Task: &documents.ApprovalTask{
			Id:     req.GetTaskId(),
			Status: documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING,
		},
		Policy: &documents.ApprovalPolicy{DocumentId: int64(k.DocumentID)},
	}, nil
}
