package settings

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorssettings "github.com/fivenet-app/fivenet/v2026/services/settings/errors"
	"go.uber.org/zap"
)

func (s *Server) ListLawBooks(
	ctx context.Context,
	req *pbsettings.ListLawBooksRequest,
) (*pbsettings.ListLawBooksResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp, err := s.store.ListLawBooks(ctx, userInfo.GetSuperuser())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateLawBook(
	ctx context.Context,
	req *pbsettings.CreateOrUpdateLawBookRequest,
) (*pbsettings.CreateOrUpdateLawBookResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	lawBook, err := s.store.CreateOrUpdateLawBook(ctx, req, userInfo.GetSuperuser())
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

	if req.GetLawBook().GetId() <= 0 {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	return &pbsettings.CreateOrUpdateLawBookResponse{LawBook: lawBook}, nil
}

func (s *Server) DeleteLawBook(
	ctx context.Context,
	req *pbsettings.DeleteLawBookRequest,
) (*pbsettings.DeleteLawBookResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	lawBook, err := s.store.GetLawBook(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	var deletedAtTime *timestamp.Timestamp
	if lawBook.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
	}

	if err := s.store.DeleteLawBook(ctx, req.GetId(), deletedAtTime); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, lawBook.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return &pbsettings.DeleteLawBookResponse{DeletedAt: deletedAtTime}, nil
}

func (s *Server) ReorderLawBooks(
	ctx context.Context,
	req *pbsettings.ReorderLawBooksRequest,
) (*pbsettings.ReorderLawBooksResponse, error) {
	if err := s.store.ReorderLawBooks(ctx, req); err != nil {
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

func (s *Server) CreateOrUpdateLaw(
	ctx context.Context,
	req *pbsettings.CreateOrUpdateLawRequest,
) (*pbsettings.CreateOrUpdateLawResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	law, refreshLawBookIDs, err := s.store.CreateOrUpdateLaw(ctx, req, userInfo.GetSuperuser())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	for _, lawBookID := range refreshLawBookIDs {
		if err := s.laws.Refresh(ctx, lawBookID); err != nil {
			s.logger.Error(
				"failed to refresh law book",
				zap.Int64("law_book_id", lawBookID),
				zap.Error(err),
			)
		}
	}

	if req.GetLaw().GetId() <= 0 {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	return &pbsettings.CreateOrUpdateLawResponse{Law: law}, nil
}

func (s *Server) DeleteLaw(
	ctx context.Context,
	req *pbsettings.DeleteLawRequest,
) (*pbsettings.DeleteLawResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	law, err := s.store.GetLaw(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	var deletedAtTime *timestamp.Timestamp
	if law.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
	}

	if err := s.store.DeleteLaw(ctx, req.GetId(), deletedAtTime); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.laws.Refresh(ctx, law.GetLawbookId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return &pbsettings.DeleteLawResponse{DeletedAt: deletedAtTime}, nil
}

func (s *Server) ReorderLaws(
	ctx context.Context,
	req *pbsettings.ReorderLawsRequest,
) (*pbsettings.ReorderLawsResponse, error) {
	if err := s.store.ReorderLaws(ctx, req); err != nil {
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
