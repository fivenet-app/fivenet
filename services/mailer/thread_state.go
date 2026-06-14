package mailer

import (
	"context"

	maileraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/access"
	mailerevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/events"
	pbmailer "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsmailer "github.com/fivenet-app/fivenet/v2026/services/mailer/errors"
)

func (s *Server) GetThreadState(
	ctx context.Context,
	req *pbmailer.GetThreadStateRequest,
) (*pbmailer.GetThreadStateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(
		ctx,
		userInfo,
		req.GetThreadId(),
		req.GetEmailId(),
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
	); err != nil {
		return nil, err
	}

	state, err := s.store.GetThreadState(ctx, s.db, req.GetThreadId(), req.GetEmailId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &pbmailer.GetThreadStateResponse{
		State: state,
	}, nil
}

func (s *Server) SetThreadState(
	ctx context.Context,
	req *pbmailer.SetThreadStateRequest,
) (*pbmailer.SetThreadStateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(
		ctx,
		userInfo,
		req.GetState().GetThreadId(),
		req.GetState().GetEmailId(),
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_WRITE),
	); err != nil {
		return nil, err
	}

	if err := s.store.SetThreadState(ctx, s.db, req.GetState()); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	state, err := s.store.GetThreadState(
		ctx,
		s.db,
		req.GetState().GetThreadId(),
		req.GetState().GetEmailId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailerevents.MailerEvent{
		Data: &mailerevents.MailerEvent_ThreadStateUpdate{
			ThreadStateUpdate: state,
		},
	}, req.GetState().GetEmailId())

	return &pbmailer.SetThreadStateResponse{
		State: state,
	}, nil
}
