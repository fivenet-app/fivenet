package documents

import (
	"context"
	"time"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) ListSignatures(
	ctx context.Context,
	req *pbdocuments.ListSignaturesRequest,
) (*pbdocuments.ListSignaturesResponse, error) {
	tSignatures := table.FivenetDocumentsSignatures

	// Build WHERE
	condition := tSignatures.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))

	if req.GetSnapshotDate() != nil {
		snap := req.GetSnapshotDate().AsTime()
		condition = condition.AND(tSignatures.SnapshotDate.EQ(mysql.TimestampT(snap)))
	}

	// Optional status filter
	if len(req.GetStatuses()) > 0 {
		vals := make([]mysql.Expression, 0, len(req.GetStatuses()))
		for _, st := range req.GetStatuses() {
			vals = append(vals, mysql.Int32(int32(st.Number())))
		}
		condition = condition.AND(tSignatures.Status.IN(vals...))
	}

	countStmt := tSignatures.
		SELECT(mysql.COUNT(tSignatures.ID)).
		FROM(tSignatures).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	pag, limit := req.GetPagination().GetResponse(DocsDefaultPageSize)
	resp := &pbdocuments.ListSignaturesResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	stmt := tSignatures.
		SELECT(
			tSignatures.ID,
			tSignatures.DocumentID,
			tSignatures.SnapshotDate,
			tSignatures.RequirementID,
			tSignatures.UserID,
			tSignatures.UserJob,
			tSignatures.Type,
			tSignatures.PayloadJSON,
			tSignatures.StampID,
			tSignatures.Status,
			tSignatures.Reason,
			tSignatures.CreatedAt,
			tSignatures.RevokedAt,
		).
		FROM(tSignatures).
		WHERE(condition).
		ORDER_BY(tSignatures.CreatedAt.ASC()).
		OFFSET(pag.GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Signatures); err != nil {
		return nil, err
	}

	return resp, nil
}

// UpsertRequirement.
func (s *Server) UpsertRequirement(
	ctx context.Context,
	req *pbdocuments.UpsertRequirementRequest,
) (*pbdocuments.UpsertRequirementResponse, error) {
	tSigReq := table.FivenetDocumentsSignatureRequirements
	now := time.Now()
	r := req.GetRequirement()

	stmt := tSigReq.
		INSERT(
			tSigReq.DocumentID,
			tSigReq.SnapshotDate,
			tSigReq.Label,
			tSigReq.Required,
			tSigReq.BindingMode,
			tSigReq.AllowedTypesMask,
			tSigReq.CreatedAt,
		).
		VALUES(
			r.GetDocumentId(),
			now, // TODO: use r.SnapshotDate
			r.GetLabel(),
			r.GetRequired(),
			int32(r.GetBindingMode()),
			r.GetAllowedTypes(),
			now,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tSigReq.Label.SET(mysql.String(r.GetLabel())),
			tSigReq.Required.SET(mysql.Bool(r.GetRequired())),
			tSigReq.BindingMode.SET(mysql.Int32(int32(r.GetBindingMode()))),
			tSigReq.AllowedTypesMask.SET(mysql.RawString("VALUES(`allowed_types_mask`)")),
			tSigReq.UpdatedAt.SET(mysql.TimestampT(now)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &pbdocuments.UpsertRequirementResponse{Requirement: r}, nil
}

// DeleteRequirement.
func (s *Server) DeleteRequirement(
	ctx context.Context,
	req *pbdocuments.DeleteRequirementRequest,
) (*pbdocuments.DeleteRequirementResponse, error) {
	tSigReq := table.FivenetDocumentsSignatureRequirements
	if _, err := tSigReq.
		DELETE().
		WHERE(tSigReq.ID.EQ(mysql.Int(req.GetRequirementId()))).
		LIMIT(1).
		ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &pbdocuments.DeleteRequirementResponse{}, nil
}

// ListRequirementAccess.
func (s *Server) ListRequirementAccess(
	ctx context.Context,
	req *pbdocuments.ListRequirementAccessRequest,
) (*pbdocuments.ListRequirementAccessResponse, error) {
	tSigReqAccess := table.FivenetDocumentsSignatureRequirementsAccess

	cond := tSigReqAccess.TargetID.EQ(mysql.Int(req.GetRequirementId()))

	stmt := tSigReqAccess.
		SELECT(
			tSigReqAccess.ID,
			tSigReqAccess.TargetID,
			tSigReqAccess.Job,
			tSigReqAccess.MinimumGrade,
			tSigReqAccess.Access,
		).
		FROM(tSigReqAccess).
		WHERE(cond).
		ORDER_BY(tSigReqAccess.ID.ASC()).
		LIMIT(40)

	var dest *documents.SignatureAccess
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &pbdocuments.ListRequirementAccessResponse{
		Access: dest,
	}, nil
}

// UpsertRequirementAccess.
func (s *Server) UpsertRequirementAccess(
	ctx context.Context,
	req *pbdocuments.UpsertRequirementAccessRequest,
) (*pbdocuments.UpsertRequirementAccessResponse, error) {
	access := req.GetAccess()

	if _, err := s.signatureAccess.HandleAccessChanges(
		ctx,
		s.db,
		req.GetRequirementId(),
		access.GetJobs(),
		access.GetUsers(),
		nil,
	); err != nil {
		return nil, err
	}

	jobs, err := s.signatureAccess.Jobs.List(ctx, s.db, req.GetRequirementId())
	if err != nil {
		return nil, err
	}

	users, err := s.signatureAccess.Users.List(ctx, s.db, req.GetRequirementId())
	if err != nil {
		return nil, err
	}

	return &pbdocuments.UpsertRequirementAccessResponse{
		Access: &documents.SignatureAccess{
			Jobs:  jobs,
			Users: users,
		},
	}, nil
}

// DeleteRequirementAccess.
func (s *Server) DeleteRequirementAccess(
	ctx context.Context,
	req *pbdocuments.DeleteRequirementAccessRequest,
) (*pbdocuments.DeleteRequirementAccessResponse, error) {
	tSigReqAccess := table.FivenetDocumentsSignatureRequirementsAccess
	if _, err := tSigReqAccess.
		DELETE().
		WHERE(tSigReqAccess.ID.EQ(mysql.Int(req.GetId()))).
		LIMIT(1).
		ExecContext(ctx, s.db); err != nil {
		return nil, err
	}
	return &pbdocuments.DeleteRequirementAccessResponse{}, nil
}

// ApplySignature.
func (s *Server) ApplySignature(
	ctx context.Context,
	req *pbdocuments.ApplySignatureRequest,
) (*pbdocuments.ApplySignatureResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tSignatures := table.FivenetDocumentsSignatures
	tSigReq := table.FivenetDocumentsSignatureRequirements

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	now := time.Now()
	var requirementId *int64
	if req.GetRequirementId() > 0 {
		reqId := req.GetRequirementId()
		requirementId = &reqId
	}
	var stampId *int64
	if req.GetStampId() > 0 {
		id := req.GetStampId()
		stampId = &id
	}

	ins := tSignatures.
		INSERT(
			tSignatures.DocumentID,
			tSignatures.SnapshotDate,
			tSignatures.RequirementID,
			tSignatures.UserID,
			tSignatures.UserJob,
			tSignatures.Type,
			tSignatures.PayloadJSON,
			tSignatures.StampID,
			tSignatures.Status,
			tSignatures.CreatedAt,
		).
		VALUES(
			req.GetDocumentId(),
			now, // TODO: req.SnapshotDate
			requirementId,
			userInfo.GetUserId(),
			userInfo.GetJob(),
			int32(req.GetType()),
			req.GetPayloadJson(),
			stampId,
			int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID),
			now,
		)

	if _, err = ins.ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	// recompute signed state (required reqs vs valid signatures)
	var totalReq int64
	if err = mysql.SELECT(mysql.COUNT(tSigReq.ID)).
		FROM(tSigReq).
		WHERE(
			tSigReq.DocumentID.EQ(mysql.Int64(req.GetDocumentId())).
				// .AND(tSigReq.SnapshotDate.EQ(mysql.Timestamp(req.GetSnapshotDate().AsTime())))
				AND(tSigReq.Required.EQ(mysql.Bool(true))),
		).QueryContext(ctx, tx, &totalReq); err != nil {
		return nil, err
	}

	var haveValid int64
	if err = mysql.SELECT(mysql.COUNT(tSignatures.ID)).
		FROM(tSignatures).
		WHERE(
			tSignatures.DocumentID.EQ(mysql.Int64(req.GetDocumentId())).
				// .AND(tSignatures.SnapshotDate.EQ(mysql.Timestamp(req.GetSnapshotDate().AsTime()))).
				AND(tSignatures.Status.EQ(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID)))),
		).QueryContext(ctx, tx, &haveValid); err != nil {
		return nil, err
	}

	docSigned := totalReq > 0 && haveValid >= totalReq

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &pbdocuments.ApplySignatureResponse{
		Signature: &documents.Signature{
			DocumentId: req.GetDocumentId(),
			Type:       req.GetType(),
		},
		DocumentSigned: docSigned,
	}, nil
}

// RevokeSignature.
func (s *Server) RevokeSignature(
	ctx context.Context,
	req *pbdocuments.RevokeSignatureRequest,
) (*pbdocuments.RevokeSignatureResponse, error) {
	tSignatures := table.FivenetDocumentsSignatures
	now := time.Now()

	if _, err := tSignatures.
		UPDATE().
		SET(
			tSignatures.Status.SET(mysql.Int32(int32(documents.SignatureStatus_SIGNATURE_STATUS_REVOKED))),
			tSignatures.RevokedAt.SET(mysql.TimestampT(now)),
			tSignatures.Reason.SET(mysql.String(req.GetReason())),
		).
		WHERE(tSignatures.ID.EQ(mysql.Int64(req.GetSignatureId()))).
		LIMIT(1).
		ExecContext(ctx, s.db); err != nil {
		return nil, err
	}
	return &pbdocuments.RevokeSignatureResponse{
		Signature: &documents.Signature{
			Id:     req.GetSignatureId(),
			Status: documents.SignatureStatus_SIGNATURE_STATUS_REVOKED,
		},
	}, nil
}

// RecomputeSignatureStatus.
func (s *Server) RecomputeSignatureStatus(
	ctx context.Context,
	req *pbdocuments.RecomputeSignatureStatusRequest,
) (*pbdocuments.RecomputeSignatureStatusResponse, error) {
	tSignatures := table.FivenetDocumentsSignatures
	tSigReq := table.FivenetDocumentsSignatureRequirements

	// required total
	var totalReq int64
	if err := tSigReq.
		SELECT(mysql.COUNT(tSigReq.ID)).
		FROM(tSigReq).
		WHERE(
			tSigReq.DocumentID.EQ(mysql.Int(req.GetDocumentId())).
				// .AND(tSigReq.SnapshotDate.EQ(mysql.TimestampT(req.GetSnapshotDate().AsTime()))).
				AND(tSigReq.Required.EQ(mysql.Bool(true))),
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
		)).QueryContext(ctx, s.db, &collectedValid); err != nil {
		return nil, err
	}

	requiredRemaining := totalReq - collectedValid
	if requiredRemaining < 0 {
		requiredRemaining = 0
	}

	return &pbdocuments.RecomputeSignatureStatusResponse{
		DocumentSigned:    totalReq > 0 && collectedValid >= totalReq,
		RequiredTotal:     int32(totalReq),
		RequiredRemaining: int32(requiredRemaining),
		CollectedValid:    int32(collectedValid),
	}, nil
}
