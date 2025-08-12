package settings

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/laws"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

var (
	tLawBooks = table.FivenetLawbooks
	tLaws     = table.FivenetLawbooksLaws
)

func (s *Server) CreateOrUpdateLawBook(
	ctx context.Context,
	req *pbsettings.CreateOrUpdateLawBookRequest,
) (*pbsettings.CreateOrUpdateLawBookResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.LawsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateLawBook",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if req.GetLawBook().GetId() <= 0 {
		stmt := tLawBooks.
			INSERT(
				tLawBooks.Name,
				tLawBooks.Description,
			).
			VALUES(
				req.GetLawBook().GetName(),
				req.GetLawBook().GetDescription(),
			)

		result, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		req.LawBook.Id = uint64(lastId)

		auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
	} else {
		stmt := tLawBooks.
			UPDATE(
				tLawBooks.Name,
				tLawBooks.Description,
			).
			SET(
				req.GetLawBook().GetName(),
				req.GetLawBook().GetDescription(),
			).
			WHERE(jet.AND(
				tLawBooks.ID.EQ(jet.Uint64(req.GetLawBook().GetId())),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED
	}

	lawBook, err := s.getLawBook(ctx, req.GetLawBook().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, lawBook.GetId()); err != nil {
		s.logger.Error(
			"failed to refresh law book",
			zap.Uint64("law_book_id", lawBook.GetId()),
			zap.Error(err),
		)
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

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.LawsService_ServiceDesc.ServiceName,
		Method:  "DeleteLawBook",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	lawBook, err := s.getLawBook(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	tLawBooks := table.FivenetLawbooks

	stmt := tLawBooks.
		DELETE().
		WHERE(
			tLawBooks.ID.EQ(jet.Uint64(req.GetId())),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, lawBook.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return &pbsettings.DeleteLawBookResponse{}, nil
}

func (s *Server) getLawBook(ctx context.Context, lawbookId uint64) (*laws.LawBook, error) {
	tLawBooks := tLawBooks.AS("law_book")
	stmt := tLawBooks.
		SELECT(
			tLawBooks.ID,
			tLawBooks.CreatedAt,
			tLawBooks.UpdatedAt,
			tLawBooks.Name,
			tLawBooks.Description,
		).
		FROM(tLawBooks).
		WHERE(
			tLawBooks.ID.EQ(jet.Uint64(lawbookId)),
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

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.LawsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateLaw",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

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
				req.GetLaw().GetDescription(),
				req.GetLaw().GetHint(),
				req.GetLaw().GetFine(),
				req.GetLaw().GetDetentionTime(),
				req.GetLaw().GetStvoPoints(),
			)

		result, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		req.Law.Id = uint64(lastId)

		auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
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
			).
			SET(
				req.GetLaw().GetLawbookId(),
				req.GetLaw().GetName(),
				req.GetLaw().GetDescription(),
				req.GetLaw().GetHint(),
				req.GetLaw().GetFine(),
				req.GetLaw().GetDetentionTime(),
				req.GetLaw().GetStvoPoints(),
			).
			WHERE(jet.AND(
				tLaws.ID.EQ(jet.Uint64(req.GetLaw().GetId())),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED
	}

	law, err := s.getLaw(ctx, req.GetLaw().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, req.GetLaw().GetLawbookId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
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

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.LawsService_ServiceDesc.ServiceName,
		Method:  "DeleteLaw",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	law, err := s.getLaw(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	tLaws := table.FivenetLawbooksLaws

	stmt := tLaws.
		DELETE().
		WHERE(
			tLaws.ID.EQ(jet.Uint64(req.GetId())),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	s.laws.Refresh(ctx, law.GetLawbookId())

	return &pbsettings.DeleteLawResponse{}, nil
}

func (s *Server) getLaw(ctx context.Context, lawId uint64) (*laws.Law, error) {
	tLaws := tLaws.AS("law")
	stmt := tLaws.
		SELECT(
			tLaws.ID,
			tLaws.CreatedAt,
			tLaws.UpdatedAt,
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
			tLaws.ID.EQ(jet.Uint64(lawId)),
		).
		LIMIT(1)

	var dest laws.Law
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}
