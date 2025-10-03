package documents

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListSignaturePolicies(
	ctx context.Context,
	req *pbdocuments.ListSignaturePoliciesRequest,
) (*pbdocuments.ListSignaturePoliciesResponse, error) {
	// TODO

	return nil, nil
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
			now, // TODO: use r.SnapshotDate
			p.GetLabel(),
			p.GetRequired(),
			int32(p.GetBindingMode()),
			p.GetAllowedTypes(),
			now,
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

	resp := &pbdocuments.ListSignaturesResponse{
		Signatures: []*documents.Signature{},
	}

	stmt := tSignatures.
		SELECT(
			tSignatures.ID,
			tSignatures.DocumentID,
			tSignatures.SnapshotDate,
			tSignatures.TaskID,
			tSignatures.UserID,
			tSignatures.UserJob,
			tSignatures.Type,
			tSignatures.PayloadSvg,
			tSignatures.StampID,
			tSignatures.Status,
			tSignatures.Reason,
			tSignatures.CreatedAt,
			tSignatures.RevokedAt,
		).
		FROM(tSignatures).
		WHERE(condition).
		ORDER_BY(tSignatures.CreatedAt.ASC()).
		LIMIT(15)

	if err := stmt.QueryContext(ctx, s.db, &resp.Signatures); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return resp, nil
}

// ApplySignature.
func (s *Server) ApplySignature(
	ctx context.Context,
	req *pbdocuments.ApplySignatureRequest,
) (*pbdocuments.ApplySignatureResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tSignatures := table.FivenetDocumentsSignatures
	tSigPolicy := table.FivenetDocumentsSignaturePolicies

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	var policyId *int64
	if req.GetPolicyId() > 0 {
		polId := req.GetPolicyId()
		policyId = &polId
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
			tSignatures.TaskID,
			tSignatures.UserID,
			tSignatures.UserJob,
			tSignatures.Type,
			tSignatures.PayloadSvg,
			tSignatures.StampID,
			tSignatures.Status,
		).
		VALUES(
			req.GetDocumentId(),
			now, // TODO: req.SnapshotDate
			policyId,
			userInfo.GetUserId(),
			userInfo.GetJob(),
			int32(req.GetType()),
			req.GetPayloadSvg(),
			stampId,
			int32(documents.SignatureStatus_SIGNATURE_STATUS_VALID),
		)

	if _, err = ins.ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	// Recompute signed state (required reqs vs valid signatures)
	var totalReq int64
	if err = mysql.
		SELECT(mysql.COUNT(tSigPolicy.ID)).
		FROM(tSigPolicy).
		WHERE(
			tSigPolicy.DocumentID.EQ(mysql.Int64(req.GetDocumentId())).
				// .AND(tSigPolicy.SnapshotDate.EQ(mysql.Timestamp(req.GetSnapshotDate().AsTime())))
				AND(tSigPolicy.Required.EQ(mysql.Bool(true))),
		).QueryContext(ctx, tx, &totalReq); err != nil {
		return nil, err
	}

	var haveValid int64
	if err = mysql.
		SELECT(mysql.COUNT(tSignatures.ID)).
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
