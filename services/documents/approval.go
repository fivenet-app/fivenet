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
	if err := stmt.QueryContext(ctx, tx, policy); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if policy.Id == 0 {
		return nil, nil
	}

	return policy, nil
}

func (s *Server) getOrCreateApprovalPolicy(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
) (*documents.ApprovalPolicy, error) {
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies

	pol, err := s.getApprovalPolicy(
		ctx,
		tx,
		tApprovalPolicy.AS("approval_policy").DocumentID.EQ(mysql.Int64(documentId)),
	)
	if err != nil {
		return nil, err
	}
	if pol != nil {
		return pol, nil
	}

	// Create a new policy if it doesn't exist
	res, err := tApprovalPolicy.
		INSERT(
			tApprovalPolicy.DocumentID,
			tApprovalPolicy.RuleKind,
			tApprovalPolicy.RequiredCount,
			tApprovalPolicy.OnEditBehavior,
			tApprovalPolicy.SnapshotDate,
		).
		VALUES(
			documentId,
			int32(documents.ApprovalRuleKind_APPROVAL_RULE_KIND_REQUIRE_ALL),
			1,
			int32(documents.OnEditBehavior_ON_EDIT_BEHAVIOR_KEEP_PROGRESS),
			mysql.CURRENT_TIMESTAMP(),
		).
		ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pol, err = s.getApprovalPolicy(
		ctx,
		tx,
		tApprovalPolicy.AS("approval_policy").ID.EQ(mysql.Int64(lastId)),
	)
	if err != nil {
		return nil, err
	}
	if pol == nil {
		return nil, errorsdocuments.ErrFailedQuery
	}

	return pol, nil
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

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

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	tApprovalPolicy = tApprovalPolicy.AS("approval_policy")

	policy, err := s.getApprovalPolicy(
		ctx,
		tx,
		tApprovalPolicy.DocumentID.EQ(mysql.Int64(p.GetDocumentId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.recomputeApprovalPolicyTx(ctx, tx, policy.GetDocumentId(), policy.GetId(), policy.GetSnapshotDate().AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	policy, err = s.getApprovalPolicy(
		ctx,
		s.db,
		tApprovalPolicy.AS("approval_policy").DocumentID.EQ(mysql.Int64(p.GetDocumentId())),
	)
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
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

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

	tApprovalTasks := table.FivenetDocumentsApprovalTasks

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
			// Ensure one User task exists
			var cnt struct{ C int32 }
			if err := mysql.
				SELECT(mysql.COUNT(tApprovalTasks.ID).AS("C")).
				FROM(tApprovalTasks).
				WHERE(mysql.AND(
					tApprovalTasks.PolicyID.EQ(mysql.Int64(pol.GetId())),
					tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
					tApprovalTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER))),
					tApprovalTasks.UserID.EQ(mysql.Int32(seed.GetUserId())),
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
			if _, err := tApprovalTasks.
				INSERT(
					tApprovalTasks.DocumentID,
					tApprovalTasks.SnapshotDate,
					tApprovalTasks.PolicyID,
					tApprovalTasks.AssigneeKind,
					tApprovalTasks.UserID,
					tApprovalTasks.SlotNo,
					tApprovalTasks.Status,
					tApprovalTasks.Comment,
					tApprovalTasks.DueAt,
					tApprovalTasks.CreatorID,
					tApprovalTasks.CreatorJob,
				).
				VALUES(
					pol.GetDocumentId(),
					snap,
					pol.GetId(),
					int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER),
					seed.GetUserId(),
					1,
					int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING),
					seed.GetComment(),
					dbutils.TimestampToMySQL(seed.GetDueAt()),
					userInfo.GetUserId(),
					userInfo.GetJob(),
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
			SELECT(mysql.COUNT(tApprovalTasks.ID).AS("C")).
			FROM(tApprovalTasks).
			WHERE(mysql.AND(
				tApprovalTasks.PolicyID.EQ(mysql.Int64(pol.GetId())),
				tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
				tApprovalTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE))),
				tApprovalTasks.Job.EQ(mysql.String(seed.GetJob())),
				tApprovalTasks.MinimumGrade.EQ(mysql.Int32(seed.GetMinimumGrade())),
				tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
			)).
			LIMIT(1).
			QueryContext(ctx, tx, &have); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
				tApprovalTasks.PolicyID,
				tApprovalTasks.AssigneeKind,
				tApprovalTasks.Job,
				tApprovalTasks.MinimumGrade,
				tApprovalTasks.SlotNo,
				tApprovalTasks.Status,
				tApprovalTasks.Comment,
				tApprovalTasks.DueAt,
				tApprovalTasks.CreatorID,
				tApprovalTasks.CreatorJob,
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

	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	// Resolve policy & snapshot
	var pol documents.ApprovalPolicy
	if err := tApprovalPolicy.
		SELECT(
			tApprovalPolicy.ID,
			tApprovalPolicy.DocumentID,
			tApprovalPolicy.SnapshotDate,
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
	)

	// Delete all pending?
	if req.GetDeleteAllPending() {
		condition = condition.
			AND(tApprovalTasks.Status.EQ(
				mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)),
			))
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
	documentId := req.GetDocumentId()
	var pol documents.ApprovalPolicy
	if req.GetPolicyId() > 0 {
		if err := tApprovalPolicy.
			SELECT(
				tApprovalPolicy.ID,
				tApprovalPolicy.DocumentID,
				tApprovalPolicy.SnapshotDate,
			).
			FROM(tApprovalPolicy).
			WHERE(tApprovalPolicy.ID.EQ(mysql.Int64(req.GetPolicyId()))).
			LIMIT(1).
			QueryContext(ctx, s.db, &pol); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		if pol.Id != 0 {
			documentId = pol.GetDocumentId()
		}
	}

	// Permission: viewer can list approval artifacts
	ok, err := s.access.CanUserAccessTarget(
		ctx,
		documentId,
		user,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil || !ok {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocAccessViewDenied)
	}

	// Build filters
	condition := tApprovals.DocumentID.EQ(mysql.Int64(documentId))
	if pol.Id != 0 {
		condition = condition.AND(tApprovals.PolicyID.EQ(mysql.Int64(pol.GetId())))
	}

	if req.GetSnapshotDate() != nil {
		condition = condition.AND(
			tApprovals.SnapshotDate.EQ(mysql.TimestampT(req.GetSnapshotDate().AsTime())),
		)
	}
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
	stmt := tApprovals.
		SELECT(
			tApprovals.ID, tApprovals.DocumentID, tApprovals.PolicyID, tApprovals.SnapshotDate,
			tApprovals.UserID, tApprovals.UserJob, tApprovals.UserJobGrade,
			tApprovals.Status, tApprovals.Comment, tApprovals.TaskID,
			tApprovals.CreatedAt, tApprovals.RevokedAt,
		).
		FROM(tApprovals).
		WHERE(condition).
		ORDER_BY(tApprovals.Status.ASC(), tApprovals.CreatedAt.DESC()).
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

	tApprovals := table.FivenetDocumentsApprovals.AS("approval")
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")

	// Load the approval artifact
	var apr documents.Approval
	if err := tApprovals.
		SELECT(
			tApprovals.ID, tApprovals.DocumentID, tApprovals.PolicyID, tApprovals.SnapshotDate,
			tApprovals.UserID, tApprovals.UserJob, tApprovals.UserJobGrade,
			tApprovals.Status, tApprovals.Comment, tApprovals.TaskID,
			tApprovals.CreatedAt, tApprovals.RevokedAt,
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
		if err := tApprovalPolicy.
			SELECT(
				tApprovalPolicy.ID,
				tApprovalPolicy.DocumentID,
				tApprovalPolicy.SnapshotDate,
				tApprovalPolicy.RequiredCount,
			).
			FROM(tApprovalPolicy).
			WHERE(tApprovalPolicy.ID.EQ(mysql.Int64(policyID))).
			LIMIT(1).
			QueryContext(ctx, s.db, &pol); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		// Find by document_id + snapshot_date
		if err := tApprovalPolicy.
			SELECT(
				tApprovalPolicy.ID,
				tApprovalPolicy.DocumentID,
				tApprovalPolicy.SnapshotDate,
				tApprovalPolicy.RequiredCount,
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

	// Mark revoked
	tApprovals = table.FivenetDocumentsApprovals

	if _, err := tApprovals.
		UPDATE().
		SET(
			tApprovals.Status.SET(mysql.Int32(int32(documents.ApprovalStatus_APPROVAL_STATUS_REVOKED))),
			tApprovals.RevokedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tApprovals.Comment.SET(mysql.String(req.GetComment())),
		).
		WHERE(tApprovals.ID.EQ(mysql.Int64(req.GetApprovalId()))).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Reload artifact for response
	if err := tApprovals.
		SELECT(
			tApprovals.ID, tApprovals.DocumentID, tApprovals.PolicyID, tApprovals.SnapshotDate,
			tApprovals.UserID, tApprovals.UserJob, tApprovals.UserJobGrade,
			tApprovals.Status, tApprovals.Comment, tApprovals.TaskID,
			tApprovals.CreatedAt, tApprovals.RevokedAt,
		).
		FROM(tApprovals).
		WHERE(tApprovals.ID.EQ(mysql.Int64(req.GetApprovalId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &apr); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Recompute approvals for this policy+snapshot
	if err := s.recomputeApprovalPolicyTx(ctx, tx, pol.GetDocumentId(), policyID, pol.GetSnapshotDate().AsTime()); err != nil {
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
		"fivenet.documents.approval.policy_id", req.GetPolicyId(),
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

	tApprovalTasks := table.FivenetDocumentsApprovalTasks.AS("approval_task")
	tApprovals := table.FivenetDocumentsApprovals

	// Resolve policy, doc, snapshot
	pol, err := s.getOrCreateApprovalPolicy(
		ctx,
		s.db,
		req.GetDocumentId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	snap := pol.GetSnapshotDate().AsTime()

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
		// Additional per-task access check (user or eligible job/min_grade)
		if err := s.canUserAccessTask(ctx, s.db, req.GetTaskId(), userInfo); err != nil {
			return nil, err
		}

		// Load the task row
		decidedTask = &documents.ApprovalTask{}
		if err := tApprovalTasks.
			SELECT(
				tApprovalTasks.ID, tApprovalTasks.DocumentID, tApprovalTasks.PolicyID, tApprovalTasks.SnapshotDate,
				tApprovalTasks.AssigneeKind, tApprovalTasks.UserID, tApprovalTasks.Job, tApprovalTasks.MinimumGrade, tApprovalTasks.SlotNo,
				tApprovalTasks.Status, tApprovalTasks.Comment, tApprovalTasks.CreatedAt, tApprovalTasks.DecidedAt, tApprovalTasks.DueAt,
				tApprovalTasks.DecisionCount, tApprovalTasks.CreatorID, tApprovalTasks.CreatorJob, tApprovalTasks.ApprovalID,
			).
			FROM(tApprovalTasks).
			WHERE(tApprovalTasks.ID.EQ(mysql.Int64(req.GetTaskId()))).
			LIMIT(1).
			QueryContext(ctx, tx, decidedTask); err != nil {
			if errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
			}
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		// Must be same policy/snapshot and pending
		if decidedTask.GetPolicyId() != pol.GetId() ||
			decidedTask.GetSnapshotDate().AsTime() != snap {
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
			UPDATE().
			SET(
				tApprovalTasks.Status.SET(mysql.Int32(int32(req.GetNewStatus()))),
				tApprovalTasks.DecisionCount.SET(tApprovalTasks.DecisionCount.ADD(mysql.Int32(1))),
				tApprovalTasks.DecidedAt.SET(mysql.TimestampT(now)),
				tApprovalTasks.Comment.SET(mysql.String(req.GetComment())),
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
				tApprovalTasks.ID, tApprovalTasks.DocumentID, tApprovalTasks.PolicyID, tApprovalTasks.SnapshotDate,
				tApprovalTasks.AssigneeKind, tApprovalTasks.UserID, tApprovalTasks.Job, tApprovalTasks.MinimumGrade, tApprovalTasks.SlotNo,
				tApprovalTasks.Status, tApprovalTasks.Comment, tApprovalTasks.CreatedAt, tApprovalTasks.DecidedAt, tApprovalTasks.DueAt,
				tApprovalTasks.DecisionCount, tApprovalTasks.CreatorID, tApprovalTasks.CreatorJob, tApprovalTasks.ApprovalID,
			).
			FROM(tApprovalTasks).
			WHERE(mysql.AND(
				tApprovalTasks.PolicyID.EQ(mysql.Int64(pol.GetId())),
				tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
				tApprovalTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_USER))),
				tApprovalTasks.UserID.EQ(mysql.Int32(int32(userInfo.GetUserId()))),
				tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
			)).
			ORDER_BY(tApprovalTasks.SlotNo.ASC(), tApprovalTasks.CreatedAt.ASC()).
			LIMIT(1).
			QueryContext(ctx, tx, &candidate)

		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		useCandidate := (err == nil && candidate.Id != 0)

		// If no USER task, try JOB target where user is eligible
		if !useCandidate {
			err = tApprovalTasks.
				SELECT(
					tApprovalTasks.ID, tApprovalTasks.DocumentID, tApprovalTasks.PolicyID, tApprovalTasks.SnapshotDate,
					tApprovalTasks.AssigneeKind, tApprovalTasks.UserID, tApprovalTasks.Job, tApprovalTasks.MinimumGrade, tApprovalTasks.SlotNo,
					tApprovalTasks.Status, tApprovalTasks.Comment, tApprovalTasks.CreatedAt, tApprovalTasks.DecidedAt, tApprovalTasks.DueAt,
					tApprovalTasks.DecisionCount, tApprovalTasks.CreatorID, tApprovalTasks.CreatorJob, tApprovalTasks.ApprovalID,
				).
				FROM(tApprovalTasks).
				WHERE(mysql.AND(
					tApprovalTasks.PolicyID.EQ(mysql.Int64(pol.GetId())),
					tApprovalTasks.SnapshotDate.EQ(mysql.TimestampT(snap)),
					tApprovalTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.ApprovalAssigneeKind_APPROVAL_ASSIGNEE_KIND_JOB_GRADE))),
					tApprovalTasks.Job.EQ(mysql.String(userInfo.GetJob())),
					tApprovalTasks.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					tApprovalTasks.Status.EQ(mysql.Int32(int32(documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING))),
				)).
				ORDER_BY(tApprovalTasks.SlotNo.ASC(), tApprovalTasks.CreatedAt.ASC()).
				LIMIT(1).
				QueryContext(ctx, tx, &candidate)

			if err != nil && !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

			useCandidate = (err == nil && candidate.Id != 0)
		}

		// If a task matched, decide it; otherwise proceed as true ad-hoc (no task)
		if useCandidate {
			tApprovalTasks = table.FivenetDocumentsApprovalTasks

			if _, err := tApprovalTasks.
				UPDATE().
				SET(
					tApprovalTasks.Status.SET(mysql.Int32(int32(req.GetNewStatus()))),
					tApprovalTasks.DecisionCount.SET(tApprovalTasks.DecisionCount.ADD(mysql.Int32(1))),
					tApprovalTasks.DecidedAt.SET(mysql.TimestampT(now)),
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

	// Write/UPSERT approval artifact (unique per (policy_id, snapshot_date, user_id))
	// Map task statuses APPROVED/DECLINED -> artifact status
	artifactStatus := documents.ApprovalStatus_APPROVAL_STATUS_APPROVED
	if req.GetNewStatus() == documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_DECLINED ||
		req.GetNewStatus() == documents.ApprovalTaskStatus_APPROVAL_TASK_STATUS_CANCELLED {
		artifactStatus = documents.ApprovalStatus_APPROVAL_STATUS_DECLINED
	}

	// Insert or update the approval artifact for this user in the current round
	if _, err := tApprovals.
		INSERT(
			tApprovals.DocumentID, tApprovals.SnapshotDate, tApprovals.PolicyID,
			tApprovals.UserID, tApprovals.UserJob, tApprovals.UserJobGrade,
			tApprovals.Status, tApprovals.Comment, tApprovals.TaskID,
		).
		VALUES(
			pol.GetDocumentId(), snap, pol.GetId(),
			int32(userInfo.GetUserId()), userInfo.GetJob(), mysql.Int32(userInfo.GetJobGrade()),
			artifactStatus, req.GetComment(), nilOrInt64(taskIDForArtifact),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tApprovals.Status.SET(mysql.Int32(int32(artifactStatus))),
			tApprovals.Comment.SET(mysql.String(req.GetComment())),
			tApprovals.TaskID.SET(nilOrInt64(taskIDForArtifact)),
		).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// If we decided a task, set its approval_id backlink
	tApprovals = table.FivenetDocumentsApprovals.AS("approval")

	var artifact documents.Approval
	if err := tApprovals.
		SELECT(
			tApprovals.ID, tApprovals.DocumentID, tApprovals.PolicyID, tApprovals.SnapshotDate,
			tApprovals.UserID, tApprovals.UserJob, tApprovals.UserJobGrade,
			tApprovals.Status, tApprovals.Comment, tApprovals.TaskID, tApprovals.CreatedAt, tApprovals.RevokedAt,
		).
		FROM(tApprovals).
		WHERE(mysql.AND(
			tApprovals.PolicyID.EQ(mysql.Int64(pol.GetId())),
			tApprovals.SnapshotDate.EQ(mysql.TimestampT(snap)),
			tApprovals.UserID.EQ(mysql.Int32(int32(userInfo.GetUserId()))),
		)).
		LIMIT(1).
		QueryContext(ctx, tx, &artifact); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if decidedTask != nil && artifact.Id != 0 {
		if _, err := tApprovalTasks.
			UPDATE().
			SET(
				tApprovalTasks.ApprovalID.SET(mysql.Int64(artifact.GetId())),
			).
			WHERE(tApprovalTasks.ID.EQ(mysql.Int64(decidedTask.GetId()))).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		decidedTask.ApprovalId = &artifact.Id
		// reflect status in the in-memory task too
		decidedTask.Status = req.GetNewStatus()
		decidedTask.DecidedAt = timestamp.New(now)
		comment := req.GetComment()
		decidedTask.Comment = &comment
	}

	// Recompute rollups for (policy, snapshot)
	if err := s.recomputeApprovalPolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetId(), snap); err != nil {
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

	if err := s.canUserAccessTask(ctx, s.db, req.GetTaskId(), userInfo); err != nil {
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

	pol, err := s.getApprovalPolicy(
		ctx,
		s.db,
		tApprovalPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId())),
	)
	if err != nil {
		return nil, err
	}

	if err := s.recomputeApprovalPolicyTx(ctx, s.db, req.GetDocumentId(), pol.Id, pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pol, err = s.getApprovalPolicy(
		ctx,
		s.db,
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
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")
	tApprovals := table.FivenetDocumentsApprovals
	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	tDocumentsMeta := table.FivenetDocumentsMeta

	// Load policy (required_count, etc.) if given and exists
	var pol documents.ApprovalPolicy
	if policyID > 0 {
		if err := tApprovalPolicy.
			SELECT(
				tApprovalPolicy.ID,
				tApprovalPolicy.DocumentID,
				tApprovalPolicy.SnapshotDate,
				tApprovalPolicy.RequiredCount,
			).
			FROM(tApprovalPolicy).
			WHERE(tApprovalPolicy.ID.EQ(mysql.Int64(policyID))).
			LIMIT(1).
			QueryContext(ctx, tx, &pol); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return err
			}
		}
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
		WHERE(mysql.AND(
			tApprovals.DocumentID.EQ(mysql.Int(documentID)),
		)).
		QueryContext(ctx, tx, &agg); err != nil {
		return err
	}

	// Pending tasks
	var aggTasks struct {
		Pending  int32
		Assigned int32
	}
	if err := tApprovalTasks.
		SELECT(
			mysql.COUNT(tApprovalTasks.ID).AS("assigned"),
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

	requiredTotal := int32(0)
	if pol.RequiredCount != nil {
		requiredTotal = pol.GetRequiredCount()
	}

	requiredRemaining := max(requiredTotal-agg.Approved, 0)

	anyDeclined := agg.Declined > 0
	// Doc is approved when we have enough approved and no declines
	// (if required_total=0, any declines block approval because all are required)
	docApproved := (requiredTotal > 0 && agg.Approved >= requiredTotal || !anyDeclined)
	var apPoliciesActive int32
	if pol.Id > 0 {
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
			tApprovalPolicy.AssignedCount.SET(mysql.Int32(aggTasks.Assigned)),
			tApprovalPolicy.ApprovedCount.SET(mysql.Int32(agg.Approved)),
			tApprovalPolicy.DeclinedCount.SET(mysql.Int32(agg.Declined)),
			tApprovalPolicy.PendingCount.SET(mysql.Int32(aggTasks.Pending)),
			tApprovalPolicy.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
		).
		WHERE(mysql.AND(
			tApprovalPolicy.DocumentID.EQ(mysql.Int(documentID)),
		)).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
