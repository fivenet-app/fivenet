package rector

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/laws"
	rector "github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tLawBooks = table.FivenetLawbooks
	tLaws     = table.FivenetLawbooksLaws
)

func (s *Server) CreateOrUpdateLawBook(ctx context.Context, req *CreateOrUpdateLawBookRequest) (*CreateOrUpdateLawBookResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateLawBook",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if req.LawBook.Id <= 0 {
		stmt := tLawBooks.
			INSERT(
				tLawBooks.Name,
				tLawBooks.Description,
			).
			VALUES(
				req.LawBook.Name,
				req.LawBook.Description,
			)

		result, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(ErrFailedQuery, err)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(ErrFailedQuery, err)
		}

		req.LawBook.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		stmt := tLawBooks.
			UPDATE(
				tLawBooks.Name,
				tLawBooks.Description,
			).
			SET(
				req.LawBook.Name,
				req.LawBook.Description,
			).
			WHERE(jet.AND(
				tLawBooks.ID.EQ(jet.Uint64(req.LawBook.Id)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(ErrFailedQuery, err)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	lawBook, err := s.getLawBook(ctx, req.LawBook.Id)
	if err != nil {
		return nil, errswrap.NewError(ErrFailedQuery, err)
	}

	s.cache.RefreshLaws(ctx, lawBook.Id)

	return &CreateOrUpdateLawBookResponse{
		LawBook: lawBook,
	}, nil
}

func (s *Server) DeleteLawBook(ctx context.Context, req *DeleteLawBookRequest) (*DeleteLawBookResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "DeleteLawBook",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	lawBook, err := s.getLawBook(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(ErrFailedQuery, err)
	}

	stmt := tLawBooks.
		DELETE().
		WHERE(
			tLawBooks.ID.EQ(jet.Uint64(req.Id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(ErrFailedQuery, err)
	}

	if err := s.cache.RefreshLaws(ctx, lawBook.Id); err != nil {
		return nil, errswrap.NewError(ErrFailedQuery, err)
	}

	return &DeleteLawBookResponse{}, nil
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

func (s *Server) CreateOrUpdateLaw(ctx context.Context, req *CreateOrUpdateLawRequest) (*CreateOrUpdateLawResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateLaw",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if req.Law.Id <= 0 {
		stmt := tLaws.
			INSERT(
				tLaws.LawbookID,
				tLaws.Name,
				tLaws.Description,
				tLaws.Fine,
				tLaws.DetentionTime,
				tLaws.StvoPoints,
			).
			VALUES(
				req.Law.LawbookId,
				req.Law.Name,
				req.Law.Description,
				req.Law.Fine,
				req.Law.DetentionTime,
				req.Law.StvoPoints,
			)

		result, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(ErrFailedQuery, err)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(ErrFailedQuery, err)
		}

		req.Law.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		stmt := tLaws.
			UPDATE(
				tLaws.LawbookID,
				tLaws.Name,
				tLaws.Description,
				tLaws.Fine,
				tLaws.DetentionTime,
				tLaws.StvoPoints,
			).
			SET(
				req.Law.LawbookId,
				req.Law.Name,
				req.Law.Description,
				req.Law.Fine,
				req.Law.DetentionTime,
				req.Law.StvoPoints,
			).
			WHERE(jet.AND(
				tLaws.ID.EQ(jet.Uint64(req.Law.Id)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(ErrFailedQuery, err)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	law, err := s.getLaw(ctx, req.Law.Id)
	if err != nil {
		return nil, errswrap.NewError(ErrFailedQuery, err)
	}

	if err := s.cache.RefreshLaws(ctx, req.Law.LawbookId); err != nil {
		return nil, errswrap.NewError(ErrFailedQuery, err)
	}

	return &CreateOrUpdateLawResponse{
		Law: law,
	}, nil
}

func (s *Server) DeleteLaw(ctx context.Context, req *DeleteLawRequest) (*DeleteLawResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "DeleteLaw",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	law, err := s.getLaw(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(ErrFailedQuery, err)
	}

	stmt := tLaws.
		DELETE().
		WHERE(
			tLaws.ID.EQ(jet.Uint64(req.Id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(ErrFailedQuery, err)
	}

	s.cache.RefreshLaws(ctx, law.LawbookId)

	return &DeleteLawResponse{}, nil
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
