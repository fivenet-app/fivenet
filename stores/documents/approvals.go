package documentsstore

import (
	"context"
	"errors"

	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) GetApprovalPolicy(
	ctx context.Context,
	db qrm.DB,
	condition mysql.BoolExpression,
) (*documentsapproval.ApprovalPolicy, error) {
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
			tApprovalPolicy.CreatedAt,
			tApprovalPolicy.UpdatedAt,
			tApprovalPolicy.DeletedAt,
		).
		FROM(tApprovalPolicy).
		WHERE(condition).
		LIMIT(1)

	pol := &documentsapproval.ApprovalPolicy{}
	if err := stmt.QueryContext(ctx, db, pol); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if pol.GetDocumentId() == 0 {
		return nil, nil
	}

	return pol, nil
}

func (s *Store) CreateApprovalPolicy(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	pol *documentsapproval.ApprovalPolicy,
) error {
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies
	_, err := tApprovalPolicy.
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
			documentID,
			dbutils.TimestampToMySQLDateTimeSec(pol.GetSnapshotDate()),
			int32(pol.GetRuleKind()),
			pol.GetRequiredCount(),
			pol.GetOnEditBehavior(),
			pol.GetSignatureRequired(),
			pol.GetSelfApproveAllowed(),
		).
		ExecContext(ctx, tx)
	return err
}

func (s *Store) GetApprovalTask(
	ctx context.Context,
	db qrm.DB,
	taskID int64,
) (*documentsapproval.ApprovalTask, error) {
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
		WHERE(tApprovalTasks.ID.EQ(mysql.Int64(taskID))).
		LIMIT(1)

	var task documentsapproval.ApprovalTask
	if err := stmt.QueryContext(ctx, db, &task); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if task.Id == 0 {
		return nil, nil
	}

	return &task, nil
}
