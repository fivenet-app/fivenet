package rector

import (
	"context"

	database "github.com/galexrt/fivenet/proto/resources/common/database"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	audit = table.FivenetAuditLog.AS("auditentry")
)

func (s *Server) ViewAuditLog(ctx context.Context, req *ViewAuditLogRequest) (*ViewAuditLogResponse, error) {
	countStmt := audit.
		SELECT(
			jet.COUNT(audit.ID).AS("datacount.totalcount"),
		).
		FROM(
			audit,
		)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	resp := &ViewAuditLogResponse{
		Pagination: database.EmptyPaginationResponse(req.Pagination.Offset),
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := audit.
		SELECT(
			audit.AllColumns,
		).
		FROM(
			audit,
		).
		ORDER_BY(
			audit.CreatedAt.ASC(),
		).
		OFFSET(req.Pagination.Offset).
		LIMIT(database.DefaultPageLimit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Logs); err != nil {
		return nil, err
	}

	database.PaginationHelper(resp.Pagination,
		count.TotalCount,
		req.Pagination.Offset,
		len(resp.Logs))

	return resp, nil
}
