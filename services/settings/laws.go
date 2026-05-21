package settings

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/laws"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorssettings "github.com/fivenet-app/fivenet/v2026/services/settings/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetLawbooks,
		IDColumn:        table.FivenetLawbooks.ID,
		DeletedAtColumn: table.FivenetLawbooks.DeletedAt,

		MinDays: 30,

		DependantTables: []*housekeeper.Table{
			{
				Table:           table.FivenetLawbooksLaws,
				IDColumn:        table.FivenetLawbooksLaws.ID,
				ForeignKey:      table.FivenetLawbooksLaws.LawbookID,
				DeletedAtColumn: table.FivenetLawbooksLaws.DeletedAt,

				MinDays: 30,
			},
		},
	})
}

func (s *Server) ListLawBooks(
	ctx context.Context,
	req *pbsettings.ListLawBooksRequest,
) (*pbsettings.ListLawBooksResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tLawBooks := table.FivenetLawbooks.AS("lawbook")
	tLaws := table.FivenetLawbooksLaws.AS("law")

	columns := mysql.ProjectionList{
		tLawBooks.ID,
		tLawBooks.CreatedAt,
		tLawBooks.UpdatedAt,
		tLawBooks.UpdatedAt,
		tLawBooks.Name,
		tLawBooks.Description,
		tLaws.ID,
		tLaws.LawbookID,
		tLaws.CreatedAt,
		tLaws.UpdatedAt,
		tLaws.Name,
		tLaws.Description,
		tLaws.Hint,
		tLaws.Fine,
		tLaws.DetentionTime,
		tLaws.StvoPoints,
	}

	condition := mysql.Bool(true)

	if !userInfo.GetSuperuser() {
		condition = mysql.AND(
			tLawBooks.DeletedAt.IS_NULL(),
			tLaws.DeletedAt.IS_NULL(),
		)
	} else {
		columns = append(columns,
			tLawBooks.DeletedAt,
			tLaws.DeletedAt,
		)
	}

	stmt := tLawBooks.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		FROM(tLawBooks.
			LEFT_JOIN(tLaws,
				tLaws.LawbookID.EQ(tLawBooks.ID),
			),
		).
		ORDER_BY(
			tLawBooks.DeletedAt.NULLS_FIRST(),
			tLawBooks.SortKey.ASC(),
			tLaws.DeletedAt.NULLS_FIRST(),
			tLaws.SortKey.ASC(),
		).
		WHERE(condition).
		LIMIT(1000)

	resp := &pbsettings.ListLawBooksResponse{
		Books: []*laws.LawBook{},
	}

	if err := stmt.QueryContext(ctx, s.db, &resp.Books); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateLawBook(
	ctx context.Context,
	req *pbsettings.CreateOrUpdateLawBookRequest,
) (*pbsettings.CreateOrUpdateLawBookResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tLawBooks := table.FivenetLawbooks
	if req.GetLawBook().GetId() <= 0 {
		stmt := tLawBooks.
			INSERT(
				tLawBooks.Name,
				tLawBooks.Description,
			).
			VALUES(
				req.GetLawBook().GetName(),
				dbutils.StringEmpty(req.GetLawBook().GetDescription()),
			).
			ON_DUPLICATE_KEY_UPDATE(
				tLawBooks.Name.SET(mysql.RawString("VALUES(`name`)")),
				tLawBooks.Description.SET(mysql.RawString("VALUES(`description`)")),
				tLawBooks.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			)

		result, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		req.LawBook.Id = lastId

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		stmt := tLawBooks.
			UPDATE(
				tLawBooks.Name,
				tLawBooks.Description,
				tLawBooks.DeletedAt,
			).
			SET(
				req.GetLawBook().GetName(),
				dbutils.StringEmpty(req.GetLawBook().GetDescription()),
				mysql.NULL,
			).
			WHERE(mysql.AND(
				tLawBooks.ID.EQ(mysql.Int64(req.GetLawBook().GetId())),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	lawBook, err := s.getLawBook(ctx, req.GetLawBook().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, lawBook.GetId()); err != nil {
		s.logger.Error(
			"failed to refresh law book",
			zap.Int64("law_book_id", lawBook.GetId()),
			zap.Error(err),
		)
	}

	if !userInfo.GetSuperuser() {
		lawBook.DeletedAt = nil
	}

	return &pbsettings.CreateOrUpdateLawBookResponse{
		LawBook: lawBook,
	}, nil
}

func (s *Server) DeleteLawBook(
	ctx context.Context,
	req *pbsettings.DeleteLawBookRequest,
) (*pbsettings.DeleteLawBookResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	lawBook, err := s.getLawBook(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	var deletedAtTime *timestamp.Timestamp
	if lawBook.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
	}

	tLawBooks := table.FivenetLawbooks
	stmt := tLawBooks.
		UPDATE().
		SET(
			tLawBooks.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAtTime)),
		).
		WHERE(
			tLawBooks.ID.EQ(mysql.Int64(req.GetId())),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, lawBook.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return &pbsettings.DeleteLawBookResponse{
		DeletedAt: deletedAtTime,
	}, nil
}

func (s *Server) getLawBook(ctx context.Context, lawbookId int64) (*laws.LawBook, error) {
	tLawBooks := table.FivenetLawbooks.AS("law_book")

	stmt := tLawBooks.
		SELECT(
			tLawBooks.ID,
			tLawBooks.CreatedAt,
			tLawBooks.UpdatedAt,
			tLawBooks.DeletedAt,
			tLawBooks.Name,
			tLawBooks.Description,
		).
		FROM(tLawBooks).
		WHERE(
			tLawBooks.ID.EQ(mysql.Int64(lawbookId)),
		).
		LIMIT(1)

	var dest laws.LawBook
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Server) CreateOrUpdateLaw(
	ctx context.Context,
	req *pbsettings.CreateOrUpdateLawRequest,
) (*pbsettings.CreateOrUpdateLawResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tLaws := table.FivenetLawbooksLaws
	if req.GetLaw().GetId() <= 0 {
		stmt := tLaws.
			INSERT(
				tLaws.LawbookID,
				tLaws.Name,
				tLaws.Description,
				tLaws.Hint,
				tLaws.Fine,
				tLaws.DetentionTime,
				tLaws.StvoPoints,
			).
			VALUES(
				req.GetLaw().GetLawbookId(),
				req.GetLaw().GetName(),
				dbutils.StringEmpty(req.GetLaw().GetDescription()),
				dbutils.StringEmpty(req.GetLaw().GetHint()),
				req.GetLaw().Fine,
				req.GetLaw().DetentionTime,
				req.GetLaw().StvoPoints,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tLaws.LawbookID.SET(mysql.RawInt("VALUES(`lawbook_id`)")),
				tLaws.Name.SET(mysql.RawString("VALUES(`name`)")),
				tLaws.Description.SET(mysql.RawString("VALUES(`description`)")),
				tLaws.Hint.SET(mysql.RawString("VALUES(`hint`)")),
				tLaws.Fine.SET(mysql.RawInt("VALUES(`fine`)")),
				tLaws.DetentionTime.SET(mysql.RawInt("VALUES(`detention_time`)")),
				tLaws.StvoPoints.SET(mysql.RawInt("VALUES(`stvo_points`)")),
				tLaws.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			)

		result, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		req.Law.Id = lastId

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		stmt := tLaws.
			UPDATE(
				tLaws.LawbookID,
				tLaws.Name,
				tLaws.Description,
				tLaws.Hint,
				tLaws.Fine,
				tLaws.DetentionTime,
				tLaws.StvoPoints,
				tLaws.DeletedAt,
			).
			SET(
				req.GetLaw().GetLawbookId(),
				req.GetLaw().GetName(),
				req.GetLaw().Description,
				req.GetLaw().Hint,
				req.GetLaw().Fine,
				req.GetLaw().DetentionTime,
				req.GetLaw().StvoPoints,
				mysql.NULL,
			).
			WHERE(mysql.AND(
				tLaws.ID.EQ(mysql.Int64(req.GetLaw().GetId())),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	law, err := s.getLaw(ctx, req.GetLaw().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, req.GetLaw().GetLawbookId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if !userInfo.GetSuperuser() {
		law.DeletedAt = nil
	}

	return &pbsettings.CreateOrUpdateLawResponse{
		Law: law,
	}, nil
}

func (s *Server) DeleteLaw(
	ctx context.Context,
	req *pbsettings.DeleteLawRequest,
) (*pbsettings.DeleteLawResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	law, err := s.getLaw(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	var deletedAtTime *timestamp.Timestamp
	if law.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
	}

	tLaws := table.FivenetLawbooksLaws
	stmt := tLaws.
		UPDATE().
		SET(
			tLaws.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAtTime)),
		).
		WHERE(
			tLaws.ID.EQ(mysql.Int64(req.GetId())),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, law.GetLawbookId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return &pbsettings.DeleteLawResponse{
		DeletedAt: deletedAtTime,
	}, nil
}

func (s *Server) getLaw(ctx context.Context, lawId int64) (*laws.Law, error) {
	tLaws := table.FivenetLawbooksLaws.AS("law")

	stmt := tLaws.
		SELECT(
			tLaws.ID,
			tLaws.CreatedAt,
			tLaws.UpdatedAt,
			tLaws.DeletedAt,
			tLaws.LawbookID,
			tLaws.Name,
			tLaws.Description,
			tLaws.Hint,
			tLaws.Fine,
			tLaws.DetentionTime,
			tLaws.StvoPoints,
		).
		FROM(tLaws).
		WHERE(
			tLaws.ID.EQ(mysql.Int64(lawId)),
		).
		LIMIT(1)

	var dest laws.Law
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}
