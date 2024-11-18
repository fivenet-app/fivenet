package mailer

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	errorsmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) SetThreadState(ctx context.Context, req *SetThreadStateRequest) (*SetThreadStateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(ctx, userInfo, req.State.ThreadId, req.State.EmailId, mailer.AccessLevel_ACCESS_LEVEL_READ); err != nil {
		return nil, err
	}

	tThreadsState := table.FivenetMailerThreadsState
	stmt := tThreadsState.
		INSERT(
			tThreadsState.ThreadID,
			tThreadsState.EmailID,
			tThreadsState.Unread,
			tThreadsState.LastRead,
			tThreadsState.Important,
			tThreadsState.Favorite,
			tThreadsState.Muted,
			tThreadsState.Archived,
		).
		VALUES(
			req.State.ThreadId,
			req.State.EmailId,
			req.State.Unread,
			req.State.LastRead,
			req.State.Important,
			req.State.Favorite,
			req.State.Muted,
			req.State.Archived,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tThreadsState.Unread.SET(jet.Bool(req.State.Unread)),
			tThreadsState.LastRead.SET(jet.RawTimestamp("VALUES(`last_read`)")),
			tThreadsState.Important.SET(jet.Bool(req.State.Important)),
			tThreadsState.Favorite.SET(jet.Bool(req.State.Favorite)),
			tThreadsState.Muted.SET(jet.Bool(req.State.Muted)),
			tThreadsState.Archived.SET(jet.Bool(req.State.Archived)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	state, err := s.getThreadState(ctx, req.State.ThreadId, req.State.EmailId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_ThreadStateUpdate{
			ThreadStateUpdate: state,
		},
	}, req.State.EmailId)

	return &SetThreadStateResponse{
		State: state,
	}, nil
}

func (s *Server) getThreadState(ctx context.Context, threadId uint64, emaildId uint64) (*mailer.ThreadState, error) {
	stmt := tThreadsState.
		SELECT(
			tThreadsState.ThreadID,
			tThreadsState.EmailID,
			tThreadsState.Unread,
			tThreadsState.LastRead,
			tThreadsState.Important,
			tThreadsState.Favorite,
			tThreadsState.Muted,
			tThreadsState.Archived,
		).
		FROM(tThreadsState).
		WHERE(jet.AND(
			tThreadsState.ThreadID.EQ(jet.Uint64(threadId)),
			tThreadsState.EmailID.EQ(jet.Uint64(emaildId)),
		))

	dest := &mailer.ThreadState{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.ThreadId == 0 || dest.EmailId == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Server) setUnreadState(ctx context.Context, threadId uint64, emailIds []uint64) error {
	if len(emailIds) == 0 {
		return nil
	}

	tThreadsUserState := table.FivenetMailerThreadsState
	stmt := tThreadsUserState.
		INSERT(
			tThreadsUserState.ThreadID,
			tThreadsUserState.EmailID,
			tThreadsUserState.Unread,
		)

	for _, userId := range emailIds {
		stmt = stmt.VALUES(
			threadId,
			userId,
			true,
		)
	}

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tThreadsUserState.Unread.SET(jet.RawBool("VALUES(`unread`)")),
	)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
