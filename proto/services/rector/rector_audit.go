package rector

import (
	"context"

	database "github.com/galexrt/fivenet/proto/resources/common/database"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	audit = table.FivenetAuditLog.AS("auditentry")
	user  = table.Users.AS("usershort")
)

func (s *Server) ViewAuditLog(ctx context.Context, req *ViewAuditLogRequest) (*ViewAuditLogResponse, error) {
	// TODO add time and user filter conditions
	var condition jet.BoolExpression

	countStmt := audit.
		SELECT(
			jet.COUNT(audit.ID).AS("datacount.totalcount"),
		).
		FROM(
			audit,
		).
		WHERE(condition)

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
			user.ID,
			user.Identifier,
			user.Job,
			user.JobGrade,
			user.Firstname,
			user.Lastname,
		).
		FROM(
			audit.
				LEFT_JOIN(user,
					user.ID.EQ(audit.UserID),
				),
		).
		WHERE(condition).
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
