package documents

import (
	"context"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) ListUsableStamps(
	ctx context.Context,
	req *pbdocuments.ListUsableStampsRequest,
) (*pbdocuments.ListUsableStampsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tStamps := table.FivenetDocumentsSignaturesStamps.AS("stamp")

	condition := mysql.AND(
		tStamps.DeletedAt.IS_NULL(),
		tStamps.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
		// TODO add access check for jobs
	)

	countStmt := mysql.
		SELECT(mysql.COUNT(tStamps.ID)).
		FROM(tStamps).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	pag, limit := req.GetPagination().
		GetResponse(DocsDefaultPageSize)
	resp := &pbdocuments.ListUsableStampsResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	stmt := mysql.
		SELECT(
			tStamps.ID,
			tStamps.Name,
			tStamps.UserID,
			tStamps.SvgTemplate,
			tStamps.VariantsJSON,
			tStamps.CreatedAt,
		).
		FROM(tStamps).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(tStamps.SortKey.ASC(), tStamps.CreatedAt.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Stamps); err != nil {
		return nil, err
	}

	return resp, nil
}
