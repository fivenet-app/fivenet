package rector

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/laws"
	rector "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbrector "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/rector"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsrector "github.com/fivenet-app/fivenet/v2025/services/rector/errors"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tLawBooks = table.FivenetLawbooks
	tLaws     = table.FivenetLawbooksLaws
)

func (s *Server) CreateOrUpdateLawBook(ctx context.Context, req *pbrector.CreateOrUpdateLawBookRequest) (*pbrector.CreateOrUpdateLawBookResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorLawsService_ServiceDesc.ServiceName,
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
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
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
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	lawBook, err := s.getLawBook(ctx, req.LawBook.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	s.laws.Refresh(ctx, lawBook.Id)

	return &pbrector.CreateOrUpdateLawBookResponse{
		LawBook: lawBook,
	}, nil
}

func (s *Server) DeleteLawBook(ctx context.Context, req *pbrector.DeleteLawBookRequest) (*pbrector.DeleteLawBookResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorLawsService_ServiceDesc.ServiceName,
		Method:  "DeleteLawBook",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	lawBook, err := s.getLawBook(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	stmt := tLawBooks.
		DELETE().
		WHERE(
			tLawBooks.ID.EQ(jet.Uint64(req.Id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, lawBook.Id); err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	return &pbrector.DeleteLawBookResponse{}, nil
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

func (s *Server) CreateOrUpdateLaw(ctx context.Context, req *pbrector.CreateOrUpdateLawRequest) (*pbrector.CreateOrUpdateLawResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorLawsService_ServiceDesc.ServiceName,
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
				tLaws.Hint,
				tLaws.Fine,
				tLaws.DetentionTime,
				tLaws.StvoPoints,
			).
			VALUES(
				req.Law.LawbookId,
				req.Law.Name,
				req.Law.Description,
				req.Law.Hint,
				req.Law.Fine,
				req.Law.DetentionTime,
				req.Law.StvoPoints,
			)

		result, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}

		req.Law.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
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
				req.Law.LawbookId,
				req.Law.Name,
				req.Law.Description,
				req.Law.Hint,
				req.Law.Fine,
				req.Law.DetentionTime,
				req.Law.StvoPoints,
			).
			WHERE(jet.AND(
				tLaws.ID.EQ(jet.Uint64(req.Law.Id)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	law, err := s.getLaw(ctx, req.Law.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, req.Law.LawbookId); err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	return &pbrector.CreateOrUpdateLawResponse{
		Law: law,
	}, nil
}

func (s *Server) DeleteLaw(ctx context.Context, req *pbrector.DeleteLawRequest) (*pbrector.DeleteLawResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorLawsService_ServiceDesc.ServiceName,
		Method:  "DeleteLaw",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	law, err := s.getLaw(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	stmt := tLaws.
		DELETE().
		WHERE(
			tLaws.ID.EQ(jet.Uint64(req.Id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	s.laws.Refresh(ctx, law.LawbookId)

	return &pbrector.DeleteLawResponse{}, nil
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
