package documents

import (
	"context"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
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

	countStmt := mysql.
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

	stmt := mysql.
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
