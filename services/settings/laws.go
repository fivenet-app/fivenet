package settings

import (
	"context"
	"errors"
	"math"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/laws"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
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

type sortOrderResult struct {
	SortOrder int32 `alias:"sort_order"`
}

func (s *Server) nextLawBookSortOrder(ctx context.Context, q qrm.Queryable) (int32, error) {
	tLawBooks := table.FivenetLawbooks.AS("lawbook")

	stmt := tLawBooks.
		SELECT(
			mysql.COALESCE(mysql.MAX(tLawBooks.SortOrder), mysql.Int32(-1)).AS("sort_order"),
		).
		FROM(tLawBooks)

	var dest sortOrderResult
	if err := stmt.QueryContext(ctx, q, &dest); err != nil {
		return 0, err
	}

	return dest.SortOrder + 1, nil
}

func (s *Server) nextLawSortOrder(
	ctx context.Context,
	q qrm.Queryable,
	lawbookID int64,
) (int32, error) {
	tLaws := table.FivenetLawbooksLaws.AS("law")

	stmt := tLaws.
		SELECT(
			mysql.COALESCE(mysql.MAX(tLaws.SortOrder), mysql.Int32(-1)).AS("sort_order"),
		).
		FROM(tLaws).
		WHERE(
			tLaws.LawbookID.EQ(mysql.Int64(lawbookID)),
		)

	var dest sortOrderResult
	if err := stmt.QueryContext(ctx, q, &dest); err != nil {
		return 0, err
	}

	if dest.SortOrder == math.MaxInt32 {
		return 0, errors.New("law sort order overflow")
	}

	return dest.SortOrder + 1, nil
}

func (s *Server) ensureActiveLawBook(ctx context.Context, q qrm.Queryable, lawbookID int64) error {
	tLawBooks := table.FivenetLawbooks.AS("lawbook")

	stmt := tLawBooks.
		SELECT(tLawBooks.ID.AS("id")).
		FROM(tLawBooks).
		WHERE(mysql.AND(
			tLawBooks.ID.EQ(mysql.Int64(lawbookID)),
			tLawBooks.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	var dest struct {
		ID int64 `alias:"id"`
	}
	if err := stmt.QueryContext(ctx, q, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return errors.New("invalid lawbook")
		}
		return err
	}

	return nil
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
		tLawBooks.SortOrder,
		tLawBooks.Name,
		tLawBooks.Description,
		tLaws.ID,
		tLaws.LawbookID,
		tLaws.CreatedAt,
		tLaws.UpdatedAt,
		tLaws.SortOrder,
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
			tLawBooks.SortOrder.ASC(),
			tLawBooks.SortKey.ASC(),
			tLaws.DeletedAt.NULLS_FIRST(),
			tLaws.SortOrder.ASC(),
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	defer tx.Rollback()

	if req.GetLawBook().GetId() <= 0 {
		sortOrder, err := s.nextLawBookSortOrder(ctx, tx)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		stmt := tLawBooks.
			INSERT(
				tLawBooks.Name,
				tLawBooks.Description,
				tLawBooks.SortOrder,
			).
			VALUES(
				req.GetLawBook().GetName(),
				dbutils.StringEmpty(req.GetLawBook().GetDescription()),
				mysql.Int32(sortOrder),
			).
			ON_DUPLICATE_KEY_UPDATE(
				tLawBooks.ID.SET(mysql.RawInt("LAST_INSERT_ID(`id`)")),
				tLawBooks.Name.SET(mysql.RawString("VALUES(`name`)")),
				tLawBooks.Description.SET(mysql.RawString("VALUES(`description`)")),
				tLawBooks.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			)

		result, err := stmt.ExecContext(ctx, tx)
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

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
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

func (s *Server) ReorderLawBooks(
	ctx context.Context,
	req *pbsettings.ReorderLawBooksRequest,
) (*pbsettings.ReorderLawBooksResponse, error) {
	lawBookIds := utils.RemoveSliceDuplicates(req.GetLawBookIds())

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	defer tx.Rollback()

	tLawBooks := table.FivenetLawbooks
	stmt := tLawBooks.
		SELECT(tLawBooks.ID).
		FROM(tLawBooks).
		WHERE(tLawBooks.DeletedAt.IS_NULL()).
		LIMIT(int64(len(lawBookIds)))

	var dest []int64
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	existing := make(map[int64]struct{}, len(lawBookIds))
	for _, lawBookID := range dest {
		existing[lawBookID] = struct{}{}
	}

	if len(existing) != len(lawBookIds) {
		return nil, errswrap.NewError(
			errors.New("invalid lawbook reorder"),
			errorssettings.ErrFailedQuery,
		)
	}

	for _, lawBookID := range lawBookIds {
		if _, ok := existing[lawBookID]; !ok {
			return nil, errswrap.NewError(
				errors.New("invalid lawbook reorder"),
				errorssettings.ErrFailedQuery,
			)
		}
	}

	for idx, lawBookID := range lawBookIds {
		if idx > math.MaxInt32 {
			return nil, errswrap.NewError(
				errors.New("invalid lawbook reorder"),
				errorssettings.ErrFailedQuery,
			)
		}

		_, err := tLawBooks.
			UPDATE().
			SET(
				tLawBooks.SortOrder.SET(mysql.Int32(int32(idx))),
			).
			WHERE(mysql.AND(
				tLawBooks.ID.EQ(mysql.Int64(lawBookID)),
				tLawBooks.DeletedAt.IS_NULL(),
			)).
			LIMIT(1).
			ExecContext(ctx, tx)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, 0); err != nil {
		s.logger.Error(
			"failed to refresh law books",
			zap.Error(err),
		)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbsettings.ReorderLawBooksResponse{}, nil
}

func (s *Server) getLawBook(ctx context.Context, lawbookId int64) (*laws.LawBook, error) {
	tLawBooks := table.FivenetLawbooks.AS("law_book")

	stmt := tLawBooks.
		SELECT(
			tLawBooks.ID,
			tLawBooks.CreatedAt,
			tLawBooks.UpdatedAt,
			tLawBooks.SortOrder,
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	defer tx.Rollback()

	refreshLawBookIDs := map[int64]struct{}{}

	if req.GetLaw().GetId() <= 0 {
		if err := s.ensureActiveLawBook(ctx, tx, req.GetLaw().GetLawbookId()); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		sortOrder, err := s.nextLawSortOrder(ctx, tx, req.GetLaw().GetLawbookId())
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		stmt := tLaws.
			INSERT(
				tLaws.LawbookID,
				tLaws.SortOrder,
				tLaws.Name,
				tLaws.Description,
				tLaws.Hint,
				tLaws.Fine,
				tLaws.DetentionTime,
				tLaws.StvoPoints,
			).
			VALUES(
				req.GetLaw().GetLawbookId(),
				mysql.Int32(sortOrder),
				req.GetLaw().GetName(),
				dbutils.StringEmpty(req.GetLaw().GetDescription()),
				dbutils.StringEmpty(req.GetLaw().GetHint()),
				req.GetLaw().Fine,
				req.GetLaw().DetentionTime,
				req.GetLaw().StvoPoints,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tLaws.ID.SET(mysql.RawInt("LAST_INSERT_ID(`id`)")),
				tLaws.LawbookID.SET(mysql.RawInt("VALUES(`lawbook_id`)")),
				tLaws.Name.SET(mysql.RawString("VALUES(`name`)")),
				tLaws.Description.SET(mysql.RawString("VALUES(`description`)")),
				tLaws.Hint.SET(mysql.RawString("VALUES(`hint`)")),
				tLaws.Fine.SET(mysql.RawInt("VALUES(`fine`)")),
				tLaws.DetentionTime.SET(mysql.RawInt("VALUES(`detention_time`)")),
				tLaws.StvoPoints.SET(mysql.RawInt("VALUES(`stvo_points`)")),
				tLaws.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			)

		result, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		req.Law.Id = lastId
		refreshLawBookIDs[req.GetLaw().GetLawbookId()] = struct{}{}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		existingLaw, err := s.getLaw(ctx, req.GetLaw().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		newLawBookID := req.GetLaw().GetLawbookId()
		sortOrder := existingLaw.GetSortOrder()
		if existingLaw.GetLawbookId() != newLawBookID {
			if err := s.ensureActiveLawBook(ctx, tx, newLawBookID); err != nil {
				return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
			}

			sortOrder, err = s.nextLawSortOrder(ctx, tx, newLawBookID)
			if err != nil {
				return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
			}
			refreshLawBookIDs[existingLaw.GetLawbookId()] = struct{}{}
		}
		refreshLawBookIDs[newLawBookID] = struct{}{}

		stmt := tLaws.
			UPDATE(
				tLaws.LawbookID,
				tLaws.SortOrder,
				tLaws.Name,
				tLaws.Description,
				tLaws.Hint,
				tLaws.Fine,
				tLaws.DetentionTime,
				tLaws.StvoPoints,
				tLaws.DeletedAt,
			).
			SET(
				newLawBookID,
				mysql.Int32(sortOrder),
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

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	law, err := s.getLaw(ctx, req.GetLaw().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	for lawBookID := range refreshLawBookIDs {
		if err := s.laws.Refresh(ctx, lawBookID); err != nil {
			s.logger.Error(
				"failed to refresh law book",
				zap.Int64("law_book_id", lawBookID),
				zap.Error(err),
			)
		}
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

func (s *Server) ReorderLaws(
	ctx context.Context,
	req *pbsettings.ReorderLawsRequest,
) (*pbsettings.ReorderLawsResponse, error) {
	lawIds := utils.RemoveSliceDuplicates(req.GetLawIds())

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.ensureActiveLawBook(ctx, tx, req.GetLawBookId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	tLaws := table.FivenetLawbooksLaws
	stmt := tLaws.
		SELECT(tLaws.ID).
		FROM(tLaws).
		WHERE(mysql.AND(
			tLaws.LawbookID.EQ(mysql.Int64(req.GetLawBookId())),
			tLaws.DeletedAt.IS_NULL(),
		)).
		LIMIT(int64(len(lawIds)))

	var dest []int64
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	existing := make(map[int64]struct{}, len(lawIds))
	for _, lawID := range dest {
		existing[lawID] = struct{}{}
	}

	if len(existing) != len(lawIds) {
		return nil, errswrap.NewError(
			errors.New("invalid law reorder"),
			errorssettings.ErrFailedQuery,
		)
	}

	for _, lawID := range lawIds {
		if _, ok := existing[lawID]; !ok {
			return nil, errswrap.NewError(
				errors.New("invalid law reorder"),
				errorssettings.ErrFailedQuery,
			)
		}
	}

	for idx, lawID := range lawIds {
		if idx > math.MaxInt32 {
			return nil, errswrap.NewError(
				errors.New("invalid law reorder"),
				errorssettings.ErrFailedQuery,
			)
		}

		_, err := tLaws.
			UPDATE().
			SET(
				tLaws.SortOrder.SET(mysql.Int32(int32(idx))),
			).
			WHERE(mysql.AND(
				tLaws.ID.EQ(mysql.Int64(lawID)),
				tLaws.LawbookID.EQ(mysql.Int64(req.GetLawBookId())),
				tLaws.DeletedAt.IS_NULL(),
			)).
			LIMIT(1).
			ExecContext(ctx, tx)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, req.GetLawBookId()); err != nil {
		s.logger.Error(
			"failed to refresh law book",
			zap.Int64("law_book_id", req.GetLawBookId()),
			zap.Error(err),
		)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbsettings.ReorderLawsResponse{}, nil
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
			tLaws.SortOrder,
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
