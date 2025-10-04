package documents

import (
	"context"
	"errors"
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

	policy, err := s.getApprovalPolicy(ctx, condition)
	if err != nil {
		return nil, err
	}

	return &pbdocuments.ListApprovalPoliciesResponse{
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
	p := req.GetPolicy()

	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", p.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		p.GetDocumentId(),
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
			tApprovalPolicy.OnEditBehavior,
			tApprovalPolicy.SnapshotDate,
		).
		VALUES(
			p.GetDocumentId(),
			int32(p.GetRuleKind()),
			p.GetRequiredCount(),
			int32(p.GetOnEditBehavior()),
			mysql.CURRENT_TIMESTAMP(), // Initialize snapshot_date
		).
		ON_DUPLICATE_KEY_UPDATE(
			tApprovalPolicy.RuleKind.SET(mysql.Int32(int32(p.GetRuleKind()))),
			tApprovalPolicy.RequiredCount.SET(mysql.Int32(p.GetRequiredCount())),
			tApprovalPolicy.OnEditBehavior.SET(mysql.Int32(int32(p.GetOnEditBehavior()))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	policy, err := s.getApprovalPolicy(
		ctx,
		tApprovalPolicy.AS("approval_policy").DocumentID.EQ(mysql.Int64(p.GetDocumentId())),
	)
	if err != nil {
		return nil, err
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
			tApprovalTasks.PolicyID,
			tApprovalTasks.SnapshotDate,
			tApprovalTasks.AssigneeKind,
			tApprovalTasks.UserID,
			tApprovalTasks.Job,
			tApprovalTasks.MinimumGrade,
			tApprovalTasks.SlotNo,
			tApprovalTasks.Status,
			tApprovalTasks.Comment,
			tApprovalTasks.CreatedAt,
			tApprovalTasks.DecidedAt,
			tApprovalTasks.DueAt,
			tApprovalTasks.DecisionCount,
			tApprovalTasks.CreatorID,
			tApprovalTasks.CreatorJob,
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

// UpsertApprovalTasks.
func (s *Server) UpsertApprovalTasks(
	ctx context.Context,
	req *pbdocuments.UpsertApprovalTasksRequest,
) (*pbdocuments.UpsertApprovalTasksResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.policy_id", req.GetPolicyId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Resolve policy & snapshot
	tPol := table.FivenetDocumentsApprovalPolicies
	var pol documents.ApprovalPolicy
	if err := tPol.
		SELECT(
			tPol.ID, tPol.DocumentID, tPol.SnapshotDate,
		).
		FROM(tPol).
		WHERE(tPol.ID.EQ(mysql.Int64(req.GetPolicyId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.Id == 0 {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
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

	tTasks := table.FivenetDocumentsApprovalTasks

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	created := int32(0)
	ensured := int32(0)

	now := time.Now().UTC()

	for _, seed := range req.GetSeeds() {
		isUser := seed.GetUserId() != 0
		if isUser {
			// Ensure one User task exists
			var cnt struct{ C int32 }
			if err := mysql.
				SELECT(mysql.COUNT(tTasks.ID).AS("C")).
				FROM(tTasks).
				WHERE(mysql.AND(
					tTasks.PolicyID.EQ(mysql.Int64(pol.GetId())),
					tTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
					tTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER))),
					tTasks.UserID.EQ(mysql.Int32(seed.GetUserId())),
				)).
				LIMIT(1).
				QueryContext(ctx, tx, &cnt); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
			if cnt.C > 0 {
				ensured++
				continue
			}

			// Insert USER task with slot_no=1
			if _, err := tTasks.
				INSERT(
					tTasks.DocumentID, tTasks.SnapshotDate, tTasks.PolicyID,
					tTasks.AssigneeKind, tTasks.UserID, tTasks.SlotNo,
					tTasks.Status, tTasks.Comment, tTasks.CreatedAt, tTasks.DueAt,
					tTasks.CreatorID, tTasks.CreatorJob,
				).
				VALUES(
					pol.GetDocumentId(), snap, pol.GetId(),
					int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER), seed.GetUserId(), 1,
					int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING), seed.GetComment(), now, dbutils.TimestampToMySQL(seed.GetDueAt()),
					int32(userInfo.GetUserId()), userInfo.GetJob(),
				).
				ExecContext(ctx, tx); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
			created++
			continue
		}

		// Job target with slots
		slots := seed.GetSlots()
		if slots <= 0 {
			slots = 1
		}
		// Count existing PENDING slots for this target
		var have struct{ C int32 }
		if err := mysql.
			SELECT(mysql.COUNT(tTasks.ID).AS("C")).
			FROM(tTasks).
			WHERE(
				tTasks.PolicyID.EQ(mysql.Int64(pol.GetId())).
					AND(tTasks.SnapshotDate.EQ(mysql.TimestampT(snap))).
					AND(tTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE)))).
					AND(tTasks.Job.EQ(mysql.String(seed.GetJob()))).
					AND(tTasks.MinimumGrade.EQ(mysql.Int32(seed.GetMinimumGrade()))).
					AND(tTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))),
			).
			LIMIT(1).
			QueryContext(ctx, tx, &have); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if have.C >= int32(slots) {
			ensured++
			continue
		}

		// Insert missing [have.C+1 .. slots] rows
		ins := tTasks.
			INSERT(
				tTasks.DocumentID, tTasks.SnapshotDate, tTasks.PolicyID,
				tTasks.AssigneeKind, tTasks.Job, tTasks.MinimumGrade, tTasks.SlotNo,
				tTasks.Status, tTasks.Comment, tTasks.CreatedAt, tTasks.DueAt,
				tTasks.CreatorID, tTasks.CreatorJob,
			)
		for slot := have.C + 1; slot <= slots; slot++ {
			ins = ins.VALUES(
				pol.GetDocumentId(),
				snap,
				pol.GetId(),
				int32(
					documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE,
				),
				seed.GetJob(),
				seed.GetMinimumGrade(),
				slot,
				int32(
					documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING,
				),
				seed.GetComment(),
				now,
				dbutils.TimestampToMySQL(seed.GetDueAt()),
				int32(userInfo.GetUserId()),
				userInfo.GetJob(),
			)
		}
		if _, err := ins.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		created += int32(slots) - have.C
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, pol.DocumentId, pol.Id, pol.SnapshotDate.AsTime()); err != nil {
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

// DeleteApprovalTasks.
func (s *Server) DeleteApprovalTasks(
	ctx context.Context,
	req *pbdocuments.DeleteApprovalTasksRequest,
) (*pbdocuments.DeleteApprovalTasksResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.policy_id", req.GetPolicyId()})

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	// Resolve policy & snapshot
	tPol := table.FivenetDocumentsApprovalPolicies
	var pol documents.ApprovalPolicy
	if err := tPol.
		SELECT(
			tPol.ID,
			tPol.DocumentID,
			tPol.SnapshotDate,
		).
		FROM(tPol).
		WHERE(tPol.ID.EQ(mysql.Int64(req.GetPolicyId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.Id == 0 {
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

	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	condition := mysql.AND(
		tApprovalTasks.PolicyID.EQ(mysql.Int64(pol.GetId())),
		tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
		tApprovalTasks.Status.EQ(
			mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)),
		),
	)

	// Delete all pending?
	if req.GetDeleteAllPending() {
		condition = tApprovalTasks.PolicyID.EQ(mysql.Int64(pol.GetId())).
			AND(tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap))).
			AND(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))))
	}

	if _, err := tApprovalTasks.
		DELETE().
		WHERE(condition).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, pol.DocumentId, pol.Id, pol.SnapshotDate.AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteApprovalTasksResponse{}, nil
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
			tApprovalTasks.PolicyID,
			tApprovalTasks.SnapshotDate,
			tApprovalTasks.AssigneeKind,
			tApprovalTasks.UserID,
			tApprovalTasks.Job,
			tApprovalTasks.MinimumGrade,
			tApprovalTasks.SlotNo,
			tApprovalTasks.Status,
			tApprovalTasks.Comment,
			tApprovalTasks.CreatedAt,
			tApprovalTasks.DecidedAt,
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

	// Optional: also verify the actor is eligible for this task
	// (user_id match OR job/minimum_grade eligibility). Your existing ACL util
	// can be invoked here if desired.

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
	user := auth.MustGetUserInfoFromContext(ctx)

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")
	tApprovals := table.FivenetDocumentsApprovals.AS("approval")

	// Resolve policy & base snapshot
	var pol documents.ApprovalPolicy
	if err := tApprovalPolicy.
		SELECT(
			tApprovalPolicy.ID, tApprovalPolicy.DocumentID, tApprovalPolicy.SnapshotDate,
		).
		FROM(tApprovalPolicy).
		WHERE(tApprovalPolicy.ID.EQ(mysql.Int64(req.GetPolicyId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.Id == 0 {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	// Permission: viewer can list approval artifacts
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		pol.GetDocumentId(),
		user,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	snap := pol.GetSnapshotDate().AsTime()
	if req.GetSnapshotDate() != nil {
		snap = req.GetSnapshotDate().AsTime()
	}

	// Build filters
	condition := tApprovals.PolicyID.EQ(mysql.Int64(pol.GetId())).
		AND(tApprovals.SnapshotDate.EQ(mysql.TimestampT(snap)))

	if req.GetStatus() != 0 {
		condition = condition.AND(
			tApprovals.Status.EQ(mysql.Int32(int32(req.GetStatus().Number()))),
		)
	}
	if req.GetUserId() != 0 {
		condition = condition.AND(tApprovals.UserID.EQ(mysql.Int32(req.GetUserId())))
	}
	if req.GetTaskId() != 0 {
		condition = condition.AND(tApprovals.TaskID.EQ(mysql.Int64(req.GetTaskId())))
	}

	// Count
	countStmt := mysql.
		SELECT(mysql.COUNT(tApprovals.ID)).
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
	stmt := tApprovals.
		SELECT(
			tApprovals.ID, tApprovals.DocumentID, tApprovals.PolicyID, tApprovals.SnapshotDate,
			tApprovals.UserID, tApprovals.UserJob, tApprovals.UserJobGrade,
			tApprovals.Status, tApprovals.Comment, tApprovals.TaskID,
			tApprovals.CreatedAt, tApprovals.RevokedAt,
		).
		FROM(tApprovals).
		WHERE(condition).
		ORDER_BY(tApprovals.CreatedAt.DESC()).
		LIMIT(15)

	if err := stmt.QueryContext(ctx, s.db, &resp.Approvals); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return resp, nil
}

// RevokeApproval.
func (s *Server) RevokeApproval(
	ctx context.Context,
	req *pbdocuments.RevokeApprovalRequest,
) (*pbdocuments.RevokeApprovalResponse, error) {
	user := auth.MustGetUserInfoFromContext(ctx)

	tApr := table.FivenetDocumentsApprovals
	tPol := table.FivenetDocumentsApprovalPolicies

	// Load the approval artifact
	var apr documents.Approval
	if err := tApr.SELECT(
		tApr.ID, tApr.DocumentID, tApr.PolicyID, tApr.SnapshotDate,
		tApr.UserID, tApr.UserJob, tApr.UserJobGrade,
		tApr.Status, tApr.Comment, tApr.TaskID,
		tApr.CreatedAt, tApr.RevokedAt,
	).FROM(tApr).
		WHERE(tApr.ID.EQ(mysql.Int64(req.GetApprovalId()))).
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
		user,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	// Resolve policy (prefer artifact.policy_id; if missing, fall back to doc+snapshot)
	policyID := apr.GetPolicyId()
	var pol documents.ApprovalPolicy
	if policyID != 0 {
		if err := tPol.SELECT(
			tPol.ID, tPol.DocumentID, tPol.SnapshotDate, tPol.RequiredCount,
		).FROM(tPol).
			WHERE(tPol.ID.EQ(mysql.Int64(policyID))).
			LIMIT(1).
			QueryContext(ctx, s.db, &pol); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		// compatibility: find by document_id + snapshot_date
		if err := tPol.SELECT(
			tPol.ID, tPol.DocumentID, tPol.SnapshotDate, tPol.RequiredCount,
		).FROM(tPol).
			WHERE(
				tPol.DocumentID.EQ(mysql.Int64(apr.GetDocumentId())).
					AND(tPol.SnapshotDate.EQ(mysql.TimestampT(apr.GetSnapshotDate().AsTime()))),
			).
			LIMIT(1).
			QueryContext(ctx, s.db, &pol); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		policyID = pol.GetId()
	}
	if pol.Id == 0 {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	// Mark revoked
	if _, err := tApr.UPDATE().SET(
		tApr.Status.SET(mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_REVOKED))),
		tApr.RevokedAt.SET(mysql.TimestampT(now)),
		tApr.Comment.SET(mysql.String(req.GetComment())),
	).WHERE(tApr.ID.EQ(mysql.Int64(req.GetApprovalId()))).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Reload artifact for response
	if err := tApr.SELECT(
		tApr.ID, tApr.DocumentID, tApr.PolicyID, tApr.SnapshotDate,
		tApr.UserID, tApr.UserJob, tApr.UserJobGrade,
		tApr.Status, tApr.Comment, tApr.TaskID,
		tApr.CreatedAt, tApr.RevokedAt,
	).FROM(tApr).
		WHERE(tApr.ID.EQ(mysql.Int64(req.GetApprovalId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &apr); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Recompute approvals for this policy+snapshot
	if err := s.recomputeApprovalPolicyTx(ctx, tx, pol.GetDocumentId(), policyID, pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.RevokeApprovalResponse{Approval: &apr}, nil
}

// DecideApproval.
func (s *Server) DecideApproval(
	ctx context.Context,
	req *pbdocuments.DecideApprovalRequest,
) (*pbdocuments.DecideApprovalResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.task_id", req.GetTaskId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.canUserAccessTask(ctx, s.db, req.GetTaskId(), userInfo); err != nil {
		return nil, err
	}

	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	tDApprovalPolicy := table.FivenetDocumentsApprovalPolicies

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

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
		SnapshotDate timestamp.Timestamp
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
			tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(k.SnapshotDate.AsTime())),
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

	if err := s.recomputeApprovalPolicyTx(ctx, tx, k.DocumentID, req.PolicyId, k.SnapshotDate.AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err = tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Minimal response; fetchers can be added if you want full objects
	return &pbdocuments.DecideApprovalResponse{
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
	tDApprovalPolicy := table.FivenetDocumentsApprovalPolicies

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
		DocumentID   int64
		SnapshotDate timestamp.Timestamp
		PolicyID     int64
	}
	if err = mysql.
		SELECT(
			tApprovalTasks.DocumentID,
			tApprovalTasks.SnapshotDate,
			tApprovalTasks.PolicyID,
		).
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
			mysql.SUM(mysql.CASE().
				WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_APPROVED)))).
				THEN(mysql.Int(1)).
				ELSE(mysql.Int(0))),
			mysql.SUM(mysql.CASE().
				WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED)))).
				THEN(mysql.Int(1)).
				ELSE(mysql.Int(0))),
			mysql.SUM(mysql.CASE().
				WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))).
				THEN(mysql.Int(1)).
				ELSE(mysql.Int(0))),
		).FROM(tApprovalTasks).
		WHERE(
			tApprovalTasks.DocumentID.EQ(mysql.Int(int64(k.DocumentID))).
				AND(tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(k.SnapshotDate.AsTime()))),
		).QueryContext(ctx, tx, &agg); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err = tDApprovalPolicy.
		UPDATE().
		SET(
			tDApprovalPolicy.AssignedCount.SET(mysql.Int32(int32(agg.Assigned))),
			tDApprovalPolicy.ApprovedCount.SET(mysql.Int32(int32(agg.Approved))),
			tDApprovalPolicy.DeclinedCount.SET(mysql.Int32(int32(agg.Declined))),
			tDApprovalPolicy.PendingCount.SET(mysql.Int32(int32(agg.Pending))),
			tDApprovalPolicy.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
		).
		WHERE(tDApprovalPolicy.DocumentID.EQ(mysql.Int(int64(k.DocumentID)))).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, k.DocumentID, k.PolicyID, k.SnapshotDate.AsTime()); err != nil {
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

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	pol, err := s.getApprovalPolicy(
		ctx,
		tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId())),
	)
	if err != nil {
		return nil, err
	}

	if err := s.recomputeApprovalPolicyTx(ctx, s.db, req.GetDocumentId(), pol.Id, pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, err
	}

	pol, err = s.getApprovalPolicy(
		ctx,
		tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId())),
	)
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
	policyID int64,
	snap time.Time,
) error {
	tDApprovalPolicy := table.FivenetDocumentsApprovalPolicies
	tApprovals := table.FivenetDocumentsApprovals
	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	tDocumentsMeta := table.FivenetDocumentsMeta

	// Load policy (required_count, etc.)
	var pol documents.ApprovalPolicy
	if err := tDApprovalPolicy.
		SELECT(
			tDApprovalPolicy.ID,
			tDApprovalPolicy.DocumentID,
			tDApprovalPolicy.SnapshotDate,
			tDApprovalPolicy.RequiredCount,
		).
		FROM(tDApprovalPolicy).
		WHERE(tDApprovalPolicy.ID.EQ(mysql.Int64(policyID))).
		LIMIT(1).
		QueryContext(ctx, tx, &pol); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.Id == 0 {
		return errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	// Collected approvals (valid=not revoked & APPROVED)
	var aprApproved struct{ N int32 }
	if err := tApprovals.
		SELECT(
			mysql.COUNT(tApprovals.ID).AS("N"),
		).
		FROM(tApprovals).
		WHERE(
			tApprovals.PolicyID.EQ(mysql.Int64(policyID)).
				AND(tApprovals.SnapshotDate.EQ(mysql.TimestampT(snap))).
				AND(tApprovals.Status.EQ(mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_APPROVED)))),
		).
		QueryContext(ctx, tx, &aprApproved); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Declined (current snapshot)
	var aprDeclined struct{ N int32 }
	if err := tApprovals.
		SELECT(
			mysql.COUNT(tApprovals.ID).AS("N"),
		).
		FROM(tApprovals).
		WHERE(
			tApprovals.PolicyID.EQ(mysql.Int64(policyID)).
				AND(tApprovals.SnapshotDate.EQ(mysql.TimestampT(snap))).
				AND(tApprovals.Status.EQ(mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_DECLINED)))),
		).
		QueryContext(ctx, tx, &aprDeclined); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Pending tasks
	var pendingTasks struct{ N int32 }
	if err := tApprovalTasks.
		SELECT(
			mysql.COUNT(tApprovalTasks.ID).AS("N"),
		).
		FROM(tApprovalTasks).
		WHERE(
			tApprovalTasks.PolicyID.EQ(mysql.Int64(policyID)).
				AND(tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap))).
				AND(tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))),
		).
		QueryContext(ctx, tx, &pendingTasks); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	requiredTotal := int32(0)
	if pol.RequiredCount != nil {
		requiredTotal = pol.GetRequiredCount()
	}

	requiredRemaining := max(requiredTotal-aprApproved.N, 0)

	anyDeclined := aprDeclined.N > 0
	docApproved := (requiredTotal > 0 && aprApproved.N >= requiredTotal && !anyDeclined)

	// Update document meta rollups (document-level)
	if _, err := tDocumentsMeta.
		INSERT(
			tDocumentsMeta.DocumentID, tDocumentsMeta.RecomputedAt,
			tDocumentsMeta.Approved,
			tDocumentsMeta.ApRequiredTotal, tDocumentsMeta.ApCollectedApproved, tDocumentsMeta.ApRequiredRemaining,
			tDocumentsMeta.ApDeclinedCount, tDocumentsMeta.ApPendingCount, tDocumentsMeta.ApAnyDeclined,
			tDocumentsMeta.ApPoliciesActive,
		).
		VALUES(
			documentID, time.Now().UTC(),
			docApproved,
			requiredTotal, aprApproved.N, requiredRemaining,
			aprDeclined.N, pendingTasks.N, anyDeclined,
			1, // if later multiple approval policies are supported, compute actual count
			snap,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tDocumentsMeta.RecomputedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tDocumentsMeta.Approved.SET(mysql.Bool(docApproved)),
			tDocumentsMeta.ApRequiredTotal.SET(mysql.Int32(requiredTotal)),
			tDocumentsMeta.ApCollectedApproved.SET(mysql.Int32(aprApproved.N)),
			tDocumentsMeta.ApRequiredRemaining.SET(mysql.Int32(requiredRemaining)),
			tDocumentsMeta.ApDeclinedCount.SET(mysql.Int32(aprDeclined.N)),
			tDocumentsMeta.ApPendingCount.SET(mysql.Int32(pendingTasks.N)),
			tDocumentsMeta.ApAnyDeclined.SET(mysql.Bool(anyDeclined)),
			tDocumentsMeta.ApPoliciesActive.SET(mysql.Int32(1)),
		).
		ExecContext(ctx, tx); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return nil
}
