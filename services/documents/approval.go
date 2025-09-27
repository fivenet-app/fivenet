package documents

import (
	"context"
	"errors"
	"time"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

// Errors you might return.
var (
	errNotFound       = errors.New("not found")
	errAlreadyHandled = errors.New("already handled")
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

	if req.GetSnapshotDate() != nil {
		snap := req.GetSnapshotDate().AsTime()
		condition = condition.AND(tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)))
	}

	if len(req.GetStatuses()) > 0 {
		vals := make([]mysql.Expression, 0, len(req.GetStatuses()))
		for _, st := range req.GetStatuses() {
			vals = append(vals, mysql.Int32(int32(st)))
		}
		condition = condition.AND(tApprovalTasks.Status.IN(vals...))
	}
	if req.GetUserId() != 0 {
		condition = condition.AND(tApprovalTasks.UserID.EQ(mysql.Int32(req.GetUserId())))
	}
	if job := req.GetJob(); job != "" {
		condition = condition.AND(tApprovalTasks.Job.EQ(mysql.String(job)))
	}
	if grade := req.GetMinimumGrade(); grade != 0 {
		condition = condition.AND(tApprovalTasks.MinimumGrade.GT_EQ(mysql.Int32(grade)))
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
