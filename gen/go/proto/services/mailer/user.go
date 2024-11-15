package mailer

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	errorsmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var tSettingsBlocks = table.FivenetMailerSettingsBlocked

func (s *Server) SetThreadState(ctx context.Context, req *SetThreadStateRequest) (*SetThreadStateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.State.EmailId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_READ)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmailer.ErrFailedQuery
		}
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
			tThreadsState.LastRead.SET(jet.RawTimestamp("VALUES(`last_read`)")),
			tThreadsState.Important.SET(jet.Bool(req.State.Important)),
			tThreadsState.Favorite.SET(jet.Bool(req.State.Favorite)),
			tThreadsState.Muted.SET(jet.Bool(req.State.Muted)),
			tThreadsState.Archived.SET(jet.Bool(req.State.Archived)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &SetThreadStateResponse{}, nil
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
