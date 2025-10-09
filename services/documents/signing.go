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

func (s *Server) getSignaturePolicy(
	ctx context.Context, tx qrm.DB, condition mysql.BoolExpression,
) (*documents.SignaturePolicy, error) {
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")

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
			tSignaturePolicy.DeletedAt,
		).
		FROM(tSignaturePolicy).
		WHERE(condition).
		ORDER_BY(tSignaturePolicy.ID.ASC()).
		LIMIT(1)

	policy := &documents.SignaturePolicy{}
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

func (s *Server) getOrCreateSignaturePolicy(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	polId int64,
) (*documents.SignaturePolicy, error) {
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies

	// Attempt to get any policy if no ID is specified
	condition := tSignaturePolicy.AS("signature_policy").DocumentID.EQ(mysql.Int64(documentId))
	if polId != 0 {
		condition = condition.AND(tSignaturePolicy.AS("signature_policy").ID.EQ(mysql.Int64(polId)))
	}

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
	allowedTypesMask := &documents.SignatureTypes{
		Types: []documents.SignatureType{
			documents.SignatureType_SIGNATURE_TYPE_FREEHAND,
			documents.SignatureType_SIGNATURE_TYPE_TYPED,
			documents.SignatureType_SIGNATURE_TYPE_STAMP,
		},
	}

	res, err := tSignaturePolicy.
		INSERT(
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
			tSignaturePolicy.Label,
			tSignaturePolicy.Required,
			tSignaturePolicy.BindingMode,
			tSignaturePolicy.AllowedTypesMask,
		).
		VALUES(
			documentId,
			mysql.CURRENT_TIMESTAMP(),
			"",
			false,
			int32(documents.SignatureBindingMode_SIGNATURE_BINDING_MODE_NONBINDING),
			allowedTypesMask,
		).
		ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pol, err = s.getSignaturePolicy(
		ctx,
		tx,
		tSignaturePolicy.AS("signature_policy").ID.EQ(mysql.Int64(lastId)),
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
	p := req.GetPolicy()

	tSigPolicy := table.FivenetDocumentsSignaturePolicies

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
			mysql.CURRENT_TIMESTAMP(),
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
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")

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
					tSignatureTasks.PolicyID.EQ(mysql.Int64(pol.GetId())),
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
					tSignatureTasks.DocumentID, tSignatureTasks.SnapshotDate, tSignatureTasks.PolicyID,
					tSignatureTasks.AssigneeKind, tSignatureTasks.UserID, tSignatureTasks.SlotNo,
					tSignatureTasks.Status, tSignatureTasks.Comment, tSignatureTasks.DueAt,
					tSignatureTasks.CreatorID, tSignatureTasks.CreatorJob,
				).
				VALUES(
					pol.GetDocumentId(), snap, pol.GetId(),
					int32(documents.SignatureAssigneeKind_SIGNATURE_ASSIGNEE_KIND_USER), seed.GetUserId(), 1,
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
				tSignatureTasks.PolicyID.EQ(mysql.Int64(pol.GetId())),
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
				tSignatureTasks.PolicyID,
				tSignatureTasks.AssigneeKind,
				tSignatureTasks.Job,
				tSignatureTasks.MinimumGrade,
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
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")
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
		AND(tTasks.SnapshotDate.EQ(mysql.TimestampT(snap)))

	// Delete all pending?
	if req.GetDeleteAllPending() {
		condition = condition.
			AND(
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

	tSignatures := table.FivenetDocumentsSignatures.AS("signature")
	tSignaturePolicy := table.FivenetDocumentsSignaturePolicies.AS("signature_policy")

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

	tSignatures = table.FivenetDocumentsSignatures
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

	// Resolve policy, doc, snapshot
	pol, err := s.getOrCreateSignaturePolicy(
		ctx,
		s.db,
		req.GetDocumentId(),
		req.GetPolicyId(),
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
		decidedTask = &documents.SignatureTask{}
		if err := tSignatureTasks.
			SELECT(
				tSignatureTasks.ID,
				tSignatureTasks.DocumentID,
				tSignatureTasks.PolicyID,
				tSignatureTasks.SnapshotDate,
				tSignatureTasks.AssigneeKind,
				tSignatureTasks.UserID,
				tSignatureTasks.Job,
				tSignatureTasks.MinimumGrade,
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
			WHERE(tSignatureTasks.ID.EQ(mysql.Int64(req.GetTaskId()))).
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

	// Insert/Upsert artifact (unique by policy_id + snapshot_date + user_id)
	newSig := &documents.Signature{}
	ins := tSignatures.
		INSERT(
			tSignatures.DocumentID,
			tSignatures.SnapshotDate,
			tSignatures.PolicyID,
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
			pol.GetDocumentId(), snap, pol.GetId(),
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
			tSignatures.PolicyID,
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
			tSignatures.PolicyID.EQ(mysql.Int64(pol.GetId())),
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
	if err := s.recomputeSignaturePolicyTx(ctx, tx, pol.GetDocumentId(), pol.GetId(), snap); err != nil {
		return nil, err
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
	if err := tSignaturePolicy.
		SELECT(
			tSignaturePolicy.ID,
			tSignaturePolicy.DocumentID,
			tSignaturePolicy.SnapshotDate,
		).
		FROM(tSignaturePolicy).
		WHERE(tSignaturePolicy.ID.EQ(mysql.Int64(sig.GetPolicyId()))).
		LIMIT(1).
		QueryContext(ctx, s.db, &pol); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol.Id == 0 {
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

	// Required total
	var totalReq int64
	if err := tSigPolicy.
		SELECT(mysql.COUNT(tSigPolicy.ID)).
		FROM(tSigPolicy).
		WHERE(mysql.AND(
			tSigPolicy.DocumentID.EQ(mysql.Int(req.GetDocumentId())),
			// .AND(tSigPolicy.SnapshotDate.EQ(mysql.TimestampT(req.GetSnapshotDate().AsTime()))).
			tSigPolicy.Required.EQ(mysql.Bool(true)),
		)).
		QueryContext(ctx, s.db, &totalReq); err != nil {
		return nil, err
	}

	// Collected valid
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

	// TODO update signature policy and meta table

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

	// Required total: count of required policies for this doc & snapshot
	// (If you later support multiple signature policies per doc, you may want this across all policies.)
	var reqTotal struct {
		Required int32
		Count    int32
	}
	if err := tSignaturePolicy.
		SELECT(
			mysql.Raw("SUM(required = 1)").AS("required"),
			mysql.COUNT(tSignaturePolicy.ID).AS("count"),
		).
		FROM(tSignaturePolicy).
		WHERE(mysql.AND(
			tSignaturePolicy.DocumentID.EQ(mysql.Int64(documentID)),
			tSignaturePolicy.SnapshotDate.EQ(mysql.TimestampT(snap)),
		)).
		QueryContext(ctx, tx, &reqTotal); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Collected valid for this policy+snapshot
	var collected struct {
		Count int32
	}
	if err := tSignatures.
		SELECT(
			mysql.COUNT(tSignatures.ID).AS("count"),
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

	requiredRemaining := max(reqTotal.Required-collected.Count, 0)
	docSigned := (reqTotal.Required > 0 && collected.Count >= reqTotal.Required)

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
			documentID,
			mysql.CURRENT_TIMESTAMP(),
			docSigned,
			reqTotal.Required,
			collected.Count,
			requiredRemaining,
			reqTotal.Count,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tDocumentsMeta.RecomputedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tDocumentsMeta.Signed.SET(mysql.Bool(docSigned)),
			tDocumentsMeta.SigRequiredTotal.SET(mysql.Int32(reqTotal.Required)),
			tDocumentsMeta.SigCollectedValid.SET(mysql.Int32(collected.Count)),
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
