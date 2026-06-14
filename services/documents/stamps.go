package documents

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsstamps "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/stamps"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const stampLimit = 5

var stampSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_USE),
		int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE),
	},
}

func (s *Server) ListUsableStamps(
	ctx context.Context,
	req *pbdocuments.ListUsableStampsRequest,
) (*pbdocuments.ListUsableStampsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tStamp := table.FivenetDocumentsStamps.AS("stamp")

	deletedAtCond := mysql.Bool(true)
	if !userInfo.GetSuperuser() {
		deletedAtCond = tStamp.DeletedAt.IS_NULL()
	}

	var existsAccess mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		existsAccess = s.signingStampAccess.ACLAccessExistsCondition(
			tStamp.ID,
			userInfo,
			int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_USE),
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

func (s *Server) GetStamp(
	ctx context.Context,
	req *pbdocuments.GetStampRequest,
) (*pbdocuments.GetStampResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.signingStampAccess.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_USE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrPermissionDenied
	}

	stamp, err := s.getStamp(ctx, req.GetId(), true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.GetStampResponse{
		Stamp: stamp,
	}, nil
}

func (s *Server) getStamp(
	ctx context.Context,
	stampId int64,
	withAccess bool,
) (*documentsstamps.Stamp, error) {
	tStamp := table.FivenetDocumentsStamps.AS("stamp")

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
			tStamp.ID.EQ(mysql.Int64(stampId)),
		)).
		LIMIT(1)

	var stamp documentsstamps.Stamp
	if err := stmt.QueryContext(ctx, s.db, &stamp); err != nil {
		return nil, err
	}

	if stamp.Id == 0 {
		return nil, nil
	}

	if withAccess {
		access, err := s.signingStampAccess.ListTargetAccess(
			ctx,
			s.db,
			stampId,
			stampSubjectAccessOptions,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		stamp.Access = access
	}

	return &stamp, nil
}

func (s *Server) checkJobStampCount(ctx context.Context, job string) (int64, error) {
	tStamp := table.FivenetDocumentsStamps.AS("stamp")

	countStmt := tStamp.
		SELECT(
			mysql.COUNT(tStamp.ID).AS("data_count.total"),
		).
		FROM(tStamp).
		WHERE(tStamp.Name.EQ(mysql.String(job)))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Server) UpsertStamp(
	ctx context.Context,
	req *pbdocuments.UpsertStampRequest,
) (*pbdocuments.UpsertStampResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tStamp := table.FivenetDocumentsStamps

	st := req.GetStamp()

	// Stamps are job only and are currently limited to 5!
	if st.Access == nil {
		st.Access = &documentsstamps.StampAccess{}
	}
	if len(st.Access.Jobs) == 0 {
		// Add minimum access for the creator's job
		st.Access.Jobs = append(st.Access.Jobs, &documentsstamps.StampJobAccess{
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
			Access:       int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE),
		})
	}

	// Check if stamp count for the job exceeds the limit
	if count, err := s.checkJobStampCount(ctx, userInfo.GetJob()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	} else if count >= stampLimit && st.GetId() == 0 {
		return nil, errorsdocuments.ErrStampLimitReached
	}

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
			int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE),
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
				st.GetName(),
				st.GetSvgTemplate(),
				nil,
			).
			WHERE(tStamp.ID.EQ(mysql.Int64(st.GetId()))).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
				st.GetName(),
				st.GetSvgTemplate(),
				nil,
			)
		res, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		req.GetStamp().SetId(lastId)
	}

	if _, err := s.signingStampAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.subjectResolver,
		st.GetId(),
		st.GetAccess(),
		stampSubjectAccessOptions,
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	stamp, err := s.getStamp(ctx, st.GetId(), true)
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
		int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrPermissionDenied
	}

	tStamp := table.FivenetDocumentsStamps.AS("stamp")

	stmt := tStamp.
		DELETE().
		WHERE(tStamp.ID.EQ(mysql.Int64(req.GetStampId()))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteStampResponse{}, nil
}
