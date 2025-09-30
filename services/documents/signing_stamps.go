package documents

import (
	"context"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) ListUsableStamps(
	ctx context.Context,
	req *pbdocuments.ListUsableStampsRequest,
) (*pbdocuments.ListUsableStampsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tStamp := table.FivenetDocumentsSignaturesStamps.AS("stamp")
	tStampAccess := table.FivenetDocumentsSignaturesStampsAccess.AS("stamp_access")

	var existsAccess mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		existsAccess = mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tStampAccess).
				WHERE(mysql.AND(
					tStampAccess.TargetID.EQ(tStamp.ID),
					// Job + grade access
					mysql.AND(
						tStampAccess.Job.EQ(mysql.String(userInfo.GetJob())),
						tStampAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					),
					tStampAccess.Access.GT_EQ(
						mysql.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)),
					),
				),
				),
		)
	} else {
		existsAccess = mysql.Bool(true)
	}

	condition := mysql.AND(
		tStamp.DeletedAt.IS_NULL(),
		mysql.OR(
			tStamp.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
			existsAccess,
		),
	)

	countStmt := mysql.
		SELECT(mysql.COUNT(tStamp.ID)).
		FROM(tStamp).
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
			tStamp.ID,
			tStamp.Name,
			tStamp.UserID,
			tStamp.SvgTemplate,
			tStamp.VariantsJSON,
			tStamp.CreatedAt,
		).
		FROM(tStamp).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(tStamp.SortKey.ASC(), tStamp.CreatedAt.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Stamps); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) UpsertStamp(
	ctx context.Context,
	req *pbdocuments.UpsertStampRequest,
) (*pbdocuments.UpsertStampResponse, error) {
	// TODO implement

	return nil, nil
}

func (s *Server) DeleteStamp(
	ctx context.Context,
	req *pbdocuments.DeleteStampRequest,
) (*pbdocuments.DeleteStampResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.signatureStampAccess.CanUserAccessTarget(
		ctx,
		req.GetStampId(),
		userInfo,
		documents.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE,
	)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errorsdocuments.ErrPermissionDenied
	}

	tStamp := table.FivenetDocumentsSignaturesStamps.AS("stamp")

	stmt := tStamp.
		DELETE().
		WHERE(tStamp.ID.EQ(mysql.Int64(req.GetStampId()))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteStampResponse{}, nil
}
