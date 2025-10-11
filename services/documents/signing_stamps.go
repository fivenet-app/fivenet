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

	deletedAtCond := mysql.Bool(true)
	if !userInfo.GetSuperuser() {
		deletedAtCond = tStamp.DeletedAt.IS_NULL()
	}

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
				)),
		)
	} else {
		existsAccess = mysql.Bool(true)
	}

	condition := mysql.AND(
		deletedAtCond,
		existsAccess,
	)

	countStmt := mysql.
		SELECT(
			mysql.COUNT(tStamp.ID).AS("data_count.total"),
		).
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

func (s *Server) getStamp(ctx context.Context, stampID int64) (*documents.Stamp, error) {
	tStamp := table.FivenetDocumentsSignaturesStamps.AS("stamp")

	stmt := mysql.
		SELECT(
			tStamp.ID,
			tStamp.Name,
			tStamp.SvgTemplate,
			tStamp.VariantsJSON,
			tStamp.CreatedAt,
		).
		FROM(tStamp).
		WHERE(mysql.AND(
			tStamp.ID.EQ(mysql.Int64(stampID)),
		)).
		LIMIT(1)

	var stamp documents.Stamp
	if err := stmt.QueryContext(ctx, s.db, &stamp); err != nil {
		return nil, err
	}

	if stamp.Id == 0 {
		return nil, nil
	}

	return &stamp, nil
}

func (s *Server) UpsertStamp(
	ctx context.Context,
	req *pbdocuments.UpsertStampRequest,
) (*pbdocuments.UpsertStampResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tStamp := table.FivenetDocumentsSignaturesStamps

	st := req.GetStamp()

	// Stamps are job only!
	// TODO Ensure that at least the highest grade in the user's job has edit access

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	if st.GetId() != 0 {
		check, err := s.signingStampAccess.CanUserAccessTarget(
			ctx,
			st.GetId(),
			userInfo,
			documents.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE,
		)
		if err != nil {
			return nil, err
		}
		if !check {
			return nil, errorsdocuments.ErrPermissionDenied
		}

		stmt := tStamp.
			UPDATE(
				tStamp.Name,
				tStamp.SvgTemplate,
				tStamp.VariantsJSON,
			).
			SET(
				mysql.String(st.GetJob()),
				mysql.String(st.GetSvgTemplate()),
				nil,
			).
			WHERE(tStamp.ID.EQ(mysql.Int64(st.GetId()))).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		if _, err := s.signingStampAccess.HandleAccessChanges(ctx, s.db, st.GetId(), st.Access.Jobs, nil, nil); err != nil {
			return nil, err
		}
	} else {
		// Create new stamp in db
		stmt := tStamp.
			INSERT(
				tStamp.Name,
				tStamp.SvgTemplate,
				tStamp.VariantsJSON,
			).
			VALUES(
				st.GetJob(),
				st.GetSvgTemplate(),
				nil,
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		if _, err := s.signingStampAccess.HandleAccessChanges(ctx, s.db, st.GetId(), st.Access.Jobs, nil, nil); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	stamp, err := s.getStamp(ctx, st.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.UpsertStampResponse{
		Stamp: stamp,
	}, nil
}

func (s *Server) DeleteStamp(
	ctx context.Context,
	req *pbdocuments.DeleteStampRequest,
) (*pbdocuments.DeleteStampResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.signingStampAccess.CanUserAccessTarget(
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
