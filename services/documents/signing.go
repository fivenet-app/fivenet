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
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// ListSignaturePolicies.
func (s *Server) ListSignaturePolicies(
	ctx context.Context,
	req *pbdocuments.ListSignaturePoliciesRequest,
) (*pbdocuments.ListSignaturePoliciesResponse, error) {
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies

	resp := &pbdocuments.ListSignaturePoliciesResponse{
		Policies: []*documents.SignaturePolicy{},
	}

	stmt := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.ID,
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
			tSignaturePolicy.Label,
			tSignaturePolicy.Required,
			tSignaturePolicy.BindingMode,
			tSignaturePolicy.AllowedTypesMask,
			tSignaturePolicy.CreatedAt,
			tSignaturePolicy.UpdatedAt,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))).
		ORDER_BY(tSignaturePolicy.CreatedAt.ASC())

	if err := stmt.QueryContext(ctx, s.db, &resp.Policies); err != nil &&
		!errors.Is(err, qrm.ErrNoRows) {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return resp, nil
}

// UpsertSignaturePolicy.
func (s *Server) UpsertSignaturePolicy(
	ctx context.Context,
	req *pbdocuments.UpsertSignaturePolicyRequest,
) (*pbdocuments.UpsertSignaturePolicyResponse, error) {
	tSigPolicy := table.FivenetDocumentsSignaturePolicies
	now := time.Now()
	p := req.GetPolicy()

	stmt := tSigPolicy.
		INSERT(
			tSigPolicy.DocumentID,
			tSigPolicy.SnapshotDate,
			tSigPolicy.Label,
			tSigPolicy.Required,
			tSigPolicy.BindingMode,
			tSigPolicy.AllowedTypesMask,
		).
		VALUES(
			p.GetDocumentId(),
			now,
			p.GetLabel(),
			p.GetRequired(),
			int32(p.GetBindingMode()),
			p.GetAllowedTypes(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tSigPolicy.Label.SET(mysql.String(p.GetLabel())),
			tSigPolicy.Required.SET(mysql.Bool(p.GetRequired())),
			tSigPolicy.BindingMode.SET(mysql.Int32(int32(p.GetBindingMode()))),
			tSigPolicy.AllowedTypesMask.SET(mysql.RawString("VALUES(`allowed_types_mask`)")),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.UpsertSignaturePolicyResponse{
		Policy: p,
	}, nil
}

// DeleteSignaturePolicy.
func (s *Server) DeleteSignaturePolicy(
	ctx context.Context,
	req *pbdocuments.DeleteSignaturePolicyRequest,
) (*pbdocuments.DeleteSignaturePolicyResponse, error) {
	tSigPolicy := table.FivenetDocumentsSignaturePolicies
	if _, err := tSigPolicy.
		DELETE().
		WHERE(tSigPolicy.ID.EQ(mysql.Int(req.GetPolicyId()))).
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
	tTasks := table.FivenetDocumentsSignatureTasks

	// Build WHERE
	condition := tTasks.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))

	if req.GetSnapshotDate() != nil {
		snap := req.GetSnapshotDate().AsTime()
		condition = condition.AND(tTasks.SnapshotDate.EQ(mysql.TimestampT(snap)))
	}

	// Optional status filter
	if len(req.GetStatuses()) > 0 {
		vals := make([]mysql.Expression, 0, len(req.GetStatuses()))
		for _, st := range req.GetStatuses() {
			vals = append(vals, mysql.Int32(int32(st.Number())))
		}
		condition = condition.AND(tTasks.Status.IN(vals...))
	}

	resp := &pbdocuments.ListSignatureTasksResponse{
		Signatures: []*documents.Signature{},
	}

	stmt := tTasks.
		SELECT(
			tTasks.ID, tTasks.DocumentID, tTasks.PolicyID, tTasks.SnapshotDate,
			tTasks.AssigneeKind, tTasks.UserID, tTasks.Job, tTasks.MinimumGrade, tTasks.SlotNo,
			tTasks.Status, tTasks.Comment, tTasks.CreatedAt, tTasks.CompletedAt, tTasks.DueAt,
			tTasks.CreatorID, tTasks.CreatorJob, tTasks.SignatureID,
		).
		FROM(tTasks).
		WHERE(condition).
		ORDER_BY(tTasks.CreatedAt.ASC()).
		LIMIT(15)

	if err := stmt.QueryContext(ctx, s.db, &resp.Signatures); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return resp, nil
}

// UpsertSignatureTasks.
func (s *Server) UpsertSignatureTasks(
	ctx context.Context,
	req *pbdocuments.UpsertSignatureTasksRequest,
) (*pbdocuments.UpsertSignatureTasksResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.sig.policy_id", req.GetPolicyId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Resolve policy & snapshot
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies

	var pol documents.SignaturePolicy
	if err := tSignaturePolicy.
		SELECT(tSignaturePolicy.ID, tSignaturePolicy.DocumentID, tSignaturePolicy.SnapshotDate).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.ID.EQ(mysql.Int64(req.GetPolicyId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.Id == 0 {
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

	tTasks := table.FivenetDocumentsSignatureTasks

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
			// ensure one USER task exists
			var cnt struct{ C int32 }
			if err := mysql.
				SELECT(mysql.COUNT(tTasks.ID).AS("C")).
				FROM(tTasks).
				WHERE(
					tTasks.PolicyID.EQ(mysql.Int64(pol.GetId())).
						AND(tTasks.SnapshotDate.EQ(mysql.TimestampT(snap))).
						AND(tTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_USER)))).
						AND(tTasks.UserID.EQ(mysql.Int32(seed.GetUserId()))),
				).
				LIMIT(1).
				QueryContext(ctx, tx, &cnt); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
			if cnt.C > 0 {
				ensured++
				continue
			}
			// insert USER task with slot_no=1
			if _, err := tTasks.INSERT(
				tTasks.DocumentID, tTasks.SnapshotDate, tTasks.PolicyID,
				tTasks.AssigneeKind, tTasks.UserID, tTasks.SlotNo,
				tTasks.Status, tTasks.Comment, tTasks.CreatedAt, tTasks.DueAt,
				tTasks.CreatorID, tTasks.CreatorJob,
			).VALUES(
				pol.GetDocumentId(), snap, pol.GetId(),
				int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_USER), seed.GetUserId(), 1,
				int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING), seed.GetComment(), now, dbutils.TimestampToMySQL(seed.GetDueAt()),
				int32(userInfo.GetUserId()), userInfo.GetJob(),
			).ExecContext(ctx, tx); err != nil {
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
		// count existing PENDING slots for this target
		var have struct{ C int32 }
		if err := mysql.
			SELECT(mysql.COUNT(tTasks.ID).AS("C")).
			FROM(tTasks).
			WHERE(
				tTasks.PolicyID.EQ(mysql.Int64(pol.GetId())).
					AND(tTasks.SnapshotDate.EQ(mysql.TimestampT(snap))).
					AND(tTasks.AssigneeKind.EQ(mysql.Int32(int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_JOB_GRADE)))).
					AND(tTasks.Job.EQ(mysql.String(seed.GetJob()))).
					AND(tTasks.MinimumGrade.EQ(mysql.Int32(seed.GetMinimumGrade()))).
					AND(tTasks.Status.EQ(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING)))),
			).
			LIMIT(1).
			QueryContext(ctx, tx, &have); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if have.C >= int32(slots) {
			ensured++
			continue
		}
		// insert missing [have.C+1 .. slots]
		ins := tTasks.
			INSERT(
				tTasks.DocumentID, tTasks.SnapshotDate, tTasks.PolicyID,
				tTasks.AssigneeKind, tTasks.Job, tTasks.MinimumGrade, tTasks.SlotNo,
				tTasks.Status, tTasks.Comment, tTasks.CreatedAt, tTasks.DueAt,
				tTasks.CreatorID, tTasks.CreatorJob,
			)
		for slot := have.C + 1; slot <= slots; slot++ {
			ins = ins.
				VALUES(
					pol.GetDocumentId(),
					snap,
					pol.GetId(),
					int32(
						documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_JOB_GRADE,
					),
					seed.GetJob(),
					seed.GetMinimumGrade(),
					slot,
					int32(
						documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING,
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

	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, err
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
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.sig.policy_id", req.GetPolicyId()})

	// Resolve policy & snapshot
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies
	var pol documents.SignaturePolicy
	if err := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.ID,
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.ID.EQ(mysql.Int64(req.GetPolicyId()))).
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

	tTasks := table.FivenetDocumentsSignatureTasks
	condition := tTasks.PolicyID.EQ(mysql.Int64(pol.GetId())).
		AND(tTasks.SnapshotDate.EQ(mysql.TimestampT(snap))).
		AND(tTasks.Status.EQ(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING))))

	// Delete all pending?
	if req.GetDeleteAllPending() {
		condition = tTasks.PolicyID.EQ(mysql.Int64(pol.GetId())).
			AND(tTasks.SnapshotDate.EQ(mysql.TimestampT(snap))).
			AND(tTasks.Status.EQ(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING))))
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

	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetId(), pol.GetSnapshotDate().AsTime()); err != nil {
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
			tSignaturePolicy.ID,
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.ID.EQ(mysql.Int64(req.GetPolicyId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.Id == 0 {
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

	condition := tSignatures.PolicyID.EQ(mysql.Int64(pol.GetId())).
		AND(tSignatures.SnapshotDate.EQ(mysql.TimestampT(snap)))

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
			tSignatures.ID, tSignatures.DocumentID, tSignatures.PolicyID, tSignatures.SnapshotDate,
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

	tSignatures := table.FivenetDocumentsSignatures
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies

	// Load the signature (and related policy/doc)
	var sig documents.Signature
	if err := tSignatures.
		SELECT(
			tSignatures.ID, tSignatures.DocumentID, tSignatures.PolicyID, tSignatures.SnapshotDate,
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

	// Tx: mark revoked and recompute
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	if _, err := tSignatures.
		UPDATE().
		SET(
			tSignatures.Status.SET(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_REVOKED))),
			tSignatures.RevokedAt.SET(mysql.TimestampT(now)),
			tSignatures.Comment.SET(mysql.String(req.GetComment())),
		).
		WHERE(tSignatures.ID.EQ(mysql.Int64(req.GetSignatureId()))).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// reload artifact for response
	if err := tSignatures.
		SELECT(
			tSignatures.ID, tSignatures.DocumentID, tSignatures.PolicyID, tSignatures.SnapshotDate,
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
		SELECT(tSignaturePolicy.ID, tSignaturePolicy.DocumentID, tSignaturePolicy.SnapshotDate).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.ID.EQ(mysql.Int64(sig.GetPolicyId()))).
		LIMIT(1).
		QueryContext(ctx, tx, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Recompute (signature side)
	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, err
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

	// Resolve policy & base snapshot
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies
	var pol documents.SignaturePolicy
	if err := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.ID,
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
			tSignaturePolicy.Required,
			tSignaturePolicy.AllowedTypesMask,
			tSignaturePolicy.BindingMode,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.ID.EQ(mysql.Int64(req.GetPolicyId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.Id == 0 {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
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

	tTasks := table.FivenetDocumentsSignatureTasks
	tSigs := table.FivenetDocumentsSignatures

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	var decidedTask *documents.SignatureTask
	now := time.Now().UTC()

	// If task_id provided, validate and mark decided later
	if req.GetTaskId() > 0 {
		// Load the task row into the protobuf message directly
		decidedTask = &documents.SignatureTask{}
		if err := tTasks.
			SELECT(
				tTasks.ID,
				tTasks.DocumentID,
				tTasks.PolicyID,
				tTasks.SnapshotDate,
				tTasks.AssigneeKind,
				tTasks.UserID,
				tTasks.Job,
				tTasks.MinimumGrade,
				tTasks.SlotNo,
				tTasks.Status,
				tTasks.Comment,
				tTasks.CreatedAt,
				tTasks.CompletedAt,
				tTasks.DueAt,
				tTasks.CreatorID,
				tTasks.CreatorJob,
				tTasks.SignatureID,
			).
			FROM(tTasks).
			WHERE(tTasks.ID.EQ(mysql.Int64(req.GetTaskId()))).
			LIMIT(1).
			QueryContext(ctx, tx, decidedTask); err != nil {
			if errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
			}
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		// Task must belong to same policy/snapshot and be PENDING
		if decidedTask.GetPolicyId() != pol.GetId() ||
			decidedTask.GetSnapshotDate().AsTime() != snap {
			return nil, errorsdocuments.ErrDocAccessViewDenied
		}
		if decidedTask.GetStatus() != documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_PENDING {
			// return nil, errorsdocuments.ErrTaskAlreadyHandled
			return nil, errorsdocuments.ErrFailedQuery
		}

		// TODO enforce eligibility (user_id match OR job/min_grade >=)
	}

	// Insert/Upsert artifact (unique by policy_id + snapshot_date + user_id)
	newSig := &documents.Signature{}
	ins := tSigs.
		INSERT(
			tSigs.DocumentID,
			tSigs.SnapshotDate,
			tSigs.PolicyID,
			tSigs.UserID,
			tSigs.UserJob,
			tSigs.UserJobGrade,
			tSigs.Type,
			tSigs.PayloadSvg,
			tSigs.StampID,
			tSigs.Status,
			tSigs.Comment,
			tSigs.TaskID,
			tSigs.CreatedAt,
		).
		VALUES(
			pol.GetDocumentId(), snap, pol.GetId(),
			userInfo.GetUserId(), userInfo.GetJob(), mysql.Int32(int32(userInfo.GetJobGrade())),
			int32(req.GetType()), req.GetPayloadSvg(), nilOrInt64(req.GetStampId()),
			int32(
				documents.SignatureStatus_SIGNATURE_STATUS_VALID,
			), req.GetComment(), nilOrInt64(req.GetTaskId()), now,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tSigs.Type.SET(mysql.Int32(int32(req.GetType()))),
			tSigs.PayloadSvg.SET(mysql.String(req.GetPayloadSvg())),
			tSigs.StampID.SET(nilOrInt64(req.GetStampId())),
			tSigs.Status.SET(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID))),
			tSigs.Comment.SET(mysql.String(req.GetComment())),
			tSigs.TaskID.SET(nilOrInt64(req.GetTaskId())),
		)
	if _, err := ins.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Read back the artifact row (protobuf hydration)
	if err := tSigs.
		SELECT(
			tSigs.ID,
			tSigs.DocumentID,
			tSigs.PolicyID,
			tSigs.SnapshotDate,
			tSigs.UserID,
			tSigs.UserJob,
			tSigs.UserJobGrade,
			tSigs.Type,
			tSigs.PayloadSvg,
			tSigs.StampID,
			tSigs.Status,
			tSigs.Comment,
			tSigs.TaskID,
			tSigs.CreatedAt,
			tSigs.RevokedAt,
		).
		FROM(tSigs).
		WHERE(
			tSigs.PolicyID.EQ(mysql.Int64(pol.GetId())).
				AND(tSigs.SnapshotDate.EQ(mysql.TimestampT(snap))).
				AND(tSigs.UserID.EQ(mysql.Int32(int32(userInfo.GetUserId())))),
		).
		LIMIT(1).
		QueryContext(ctx, tx, newSig); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// If task_id provided, mark the task completed and link the signature id
	if req.GetTaskId() > 0 && decidedTask != nil {
		if _, err := tTasks.
			UPDATE().
			SET(
				tTasks.Status.SET(mysql.Int32(int32(documents.SignatureTaskStatus_SIGNATURE_TASK_STATUS_SIGNED))),
				tTasks.CompletedAt.SET(mysql.TimestampT(now)),
				tTasks.SignatureID.SET(mysql.Int64(newSig.GetId())),
			).
			WHERE(tTasks.ID.EQ(mysql.Int64(req.GetTaskId()))).
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
	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetId(), snap); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DecideSignatureResponse{
		Signature: newSig,
		Task:      decidedTask,
		Policy:    &pol,
	}, nil
}

// ReopenSignature.
func (s *Server) ReopenSignature(
	ctx context.Context,
	req *pbdocuments.ReopenSignatureRequest,
) (*pbdocuments.ReopenSignatureResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tSignatures := table.FivenetDocumentsSignatures
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies

	// Load signature
	var sig documents.Signature
	if err := tSignatures.
		SELECT(
			tSignatures.ID, tSignatures.DocumentID, tSignatures.PolicyID, tSignatures.SnapshotDate,
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
	if err := tSignaturePolicy.SELECT(tSignaturePolicy.ID, tSignaturePolicy.DocumentID, tSignaturePolicy.SnapshotDate).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.ID.EQ(mysql.Int64(sig.GetPolicyId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.Id == 0 {
		return nil, errswrap.NewError(nil, errorsdocuments.ErrNotFoundOrNoPerms)
	}

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
	if err := tSignatures.
		SELECT(
			tSignatures.ID, tSignatures.DocumentID, tSignatures.PolicyID, tSignatures.SnapshotDate,
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
	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetId(), pol.GetSnapshotDate().AsTime()); err != nil {
		return nil, err
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
	tSignatures := table.FivenetDocumentsSignatures
	tSigPolicy := table.FivenetDocumentsSignaturePolicies

	// required total
	var totalReq int64
	if err := tSigPolicy.
		SELECT(mysql.COUNT(tSigPolicy.ID)).
		FROM(tSigPolicy).
		WHERE(
			tSigPolicy.DocumentID.EQ(mysql.Int(req.GetDocumentId())).
				// .AND(tSigPolicy.SnapshotDate.EQ(mysql.TimestampT(req.GetSnapshotDate().AsTime()))).
				AND(tSigPolicy.Required.EQ(mysql.Bool(true))),
		).
		QueryContext(ctx, s.db, &totalReq); err != nil {
		return nil, err
	}

	// collected valid
	var collectedValid int64
	if err := tSignatures.
		SELECT(mysql.COUNT(tSignatures.ID)).
		FROM(tSignatures).
		WHERE(mysql.AND(
			tSignatures.DocumentID.EQ(mysql.Int(req.GetDocumentId())),
			// tSignatures.SnapshotDate.EQ(mysql.TimestampT(req.GetSnapshotDate().AsTime())).
			tSignatures.Status.EQ(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID))),
		)).
		QueryContext(ctx, s.db, &collectedValid); err != nil {
		return nil, err
	}

	requiredRemaining := max(totalReq-collectedValid, 0)

	return &pbdocuments.RecomputeSignatureStatusResponse{
		DocumentSigned:    totalReq > 0 && collectedValid >= totalReq,
		RequiredTotal:     int32(totalReq),
		RequiredRemaining: int32(requiredRemaining),
		CollectedValid:    int32(collectedValid),
	}, nil
}

// recomputeSignaturePolicyTx recalculates signature counters for a policy+snapshot
// and updates fivenet_documents_meta for list-time flags.
func (s *Server) recomputeSignaturePolicyTx(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	policyID int64,
	snap time.Time,
) error {
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies
	tSignatures := table.FivenetDocumentsSignatures
	tDocumentsMeta := table.FivenetDocumentsMeta

	// required total: count of required policies for this doc & snapshot
	// (If you later support multiple signature policies per doc, you may want this across all policies.)
	var reqTotal struct{ N int32 }
	if err := tSignaturePolicy.
		SELECT(
			mysql.Raw("SUM(required = 1) AS N"),
		).
		FROM(tSignaturePolicy).
		WHERE(
			tSignaturePolicy.DocumentID.EQ(mysql.Int64(documentID)).
				AND(tSignaturePolicy.SnapshotDate.EQ(mysql.TimestampT(snap))),
		).
		QueryContext(ctx, tx, &reqTotal); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// collected valid for this policy+snapshot
	var collected struct{ N int32 }
	if err := tSignatures.
		SELECT(
			mysql.COUNT(tSignatures.ID).AS("N"),
		).
		FROM(tSignatures).
		WHERE(mysql.AND(
			tSignatures.PolicyID.EQ(mysql.Int64(policyID)),
			tSignatures.SnapshotDate.EQ(mysql.TimestampT(snap)),
			tSignatures.Status.EQ(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID))),
		)).
		QueryContext(ctx, tx, &collected); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	requiredRemaining := max(reqTotal.N-collected.N, 0)
	docSigned := (reqTotal.N > 0 && collected.N >= reqTotal.N)

	// update meta rollups (document-level)
	if _, err := tDocumentsMeta.
		INSERT(
			tDocumentsMeta.DocumentID,
			tDocumentsMeta.RecomputedAt,
			tDocumentsMeta.Signed,
			tDocumentsMeta.SigRequiredTotal,
			tDocumentsMeta.SigCollectedValid,
			tDocumentsMeta.SigRequiredRemaining,
			tDocumentsMeta.SigPoliciesActive,
		).
		VALUES(
			documentID, time.Now().UTC(), docSigned,
			reqTotal.N, collected.N, requiredRemaining,
			1, // if you later have multiple policies, compute actual active count
			snap,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tDocumentsMeta.RecomputedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tDocumentsMeta.Signed.SET(mysql.Bool(docSigned)),
			tDocumentsMeta.SigRequiredTotal.SET(mysql.Int32(reqTotal.N)),
			tDocumentsMeta.SigCollectedValid.SET(mysql.Int32(collected.N)),
			tDocumentsMeta.SigRequiredRemaining.SET(mysql.Int32(requiredRemaining)),
			tDocumentsMeta.SigPoliciesActive.SET(mysql.Int32(1)),
		).
		ExecContext(ctx, tx); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
