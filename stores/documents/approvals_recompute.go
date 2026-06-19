package documentsstore

import (
	"context"

	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) RecomputeApprovalPolicyTx(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	snapDate *timestamp.Timestamp,
) error {
	tApprovalPolicy := table.FivenetDocumentsApprovalPolicies.AS("approval_policy")
	tApprovals := table.FivenetDocumentsApprovals
	tApprovalTasks := table.FivenetDocumentsApprovalTasks
	tDocumentsMeta := table.FivenetDocumentsMeta

	pol, err := s.GetApprovalPolicy(ctx, tx, tApprovalPolicy.DocumentID.EQ(mysql.Int64(documentID)))
	if err != nil {
		return err
	}
	if pol == nil {
		pol = &documentsapproval.ApprovalPolicy{}
	}
	pol.Default()

	var docCreatorId int32
	if !pol.SelfApproveAllowed {
		var docCreator struct {
			CreatorId int32 `alias:"creator_id"`
		}
		tDocuments := table.FivenetDocuments
		if err := tDocuments.
			SELECT(tDocuments.CreatorID.AS("creator_id")).
			FROM(tDocuments).
			WHERE(tDocuments.ID.EQ(mysql.Int64(documentID))).
			LIMIT(1).
			QueryContext(ctx, tx, &docCreator); err != nil {
			return err
		}
		if docCreator.CreatorId > 0 {
			docCreatorId = docCreator.CreatorId
		}
	}

	approvalCondition := tApprovals.DocumentID.EQ(mysql.Int64(documentID))
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
			mysql.SUM(mysql.CASE().WHEN(mysql.AND(
				tApprovals.RevokedAt.IS_NULL(),
				tApprovals.Status.EQ(
					mysql.Int32(int32(documentsapproval.ApprovalStatus_APPROVAL_STATUS_APPROVED)),
				),
			)).THEN(mysql.Int(1)).ELSE(mysql.Int(0))).AS("approved"),
			mysql.SUM(mysql.CASE().WHEN(mysql.AND(
				tApprovals.RevokedAt.IS_NULL(),
				tApprovals.Status.EQ(
					mysql.Int32(int32(documentsapproval.ApprovalStatus_APPROVAL_STATUS_DECLINED)),
				),
			)).THEN(mysql.Int(1)).ELSE(mysql.Int(0))).AS("declined"),
		).
		FROM(tApprovals).
		WHERE(approvalCondition).
		QueryContext(ctx, tx, &agg); err != nil {
		return err
	}

	var aggTasks struct {
		Total   int32
		Pending int32
	}
	if err := tApprovalTasks.
		SELECT(
			mysql.COUNT(tApprovalTasks.ID).AS("total"),
			mysql.SUM(mysql.CASE().WHEN(tApprovalTasks.Status.EQ(mysql.Int32(int32(documentsapproval.ApprovalTaskStatus_APPROVAL_TASK_STATUS_PENDING)))).THEN(mysql.Int(1)).ELSE(mysql.Int(0))).
				AS("pending"),
		).
		FROM(tApprovalTasks).
		WHERE(mysql.AND(
			tApprovalTasks.DocumentID.EQ(mysql.Int64(documentID)),
			tApprovalTasks.SnapshotDate.EQ(dbutils.TimestampToMySQLDateTimeSec(snapDate)),
		)).
		QueryContext(ctx, tx, &aggTasks); err != nil {
		return err
	}

	requiredTotal := pol.GetRequiredCount()
	requiredRemaining := max(requiredTotal-agg.Approved, 0)

	anyDeclined := agg.Declined > 0
	var docApproved bool
	if pol.GetRuleKind() == documentsapproval.ApprovalRuleKind_APPROVAL_RULE_KIND_REQUIRE_ALL {
		docApproved = !anyDeclined && agg.Approved > 0 && (agg.Approved >= aggTasks.Total)
	} else if pol.GetRuleKind() == documentsapproval.ApprovalRuleKind_APPROVAL_RULE_KIND_QUORUM_ANY {
		docApproved = (agg.Approved >= requiredTotal)
	}

	var apPoliciesActive int32
	if pol.DocumentId > 0 {
		apPoliciesActive = 1
	}

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

	tApprovalPolicy = table.FivenetDocumentsApprovalPolicies
	if _, err := tApprovalPolicy.
		UPDATE().
		SET(
			tApprovalPolicy.AssignedCount.SET(mysql.Int32(aggTasks.Total)),
			tApprovalPolicy.ApprovedCount.SET(mysql.Int32(agg.Approved)),
			tApprovalPolicy.DeclinedCount.SET(mysql.Int32(agg.Declined)),
			tApprovalPolicy.PendingCount.SET(mysql.Int32(aggTasks.Pending)),
			tApprovalPolicy.AnyDeclined.SET(mysql.Bool(agg.Declined > 0)),
		).
		WHERE(tApprovalPolicy.DocumentID.EQ(mysql.Int64(documentID))).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
