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

func (s *Server) SetThreadUserState(ctx context.Context, req *SetThreadUserStateRequest) (*SetThreadUserStateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.checkIfUserHasAccessToThread(ctx, req.State.ThreadId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_ADMIN)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmailer.ErrFailedQuery
		}
	}

	tThreadsUserState := table.FivenetMsgsThreadsUserState
	stmt := tThreadsUserState.
		INSERT(
			tThreadsUserState.ThreadID,
			tThreadsUserState.UserID,
			tThreadsUserState.Unread,
			tThreadsUserState.LastRead,
			tThreadsUserState.Important,
			tThreadsUserState.Favorite,
			tThreadsUserState.Muted,
		).
		VALUES(
			req.State.ThreadId,
			userInfo.UserId,
			req.State.Unread,
			req.State.LastRead,
			req.State.Important,
			req.State.Favorite,
			req.State.Muted,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tThreadsUserState.LastRead.SET(jet.RawTimestamp("VALUES(`last_read`)")),
			tThreadsUserState.Important.SET(jet.Bool(req.State.Important)),
			tThreadsUserState.Favorite.SET(jet.Bool(req.State.Favorite)),
			tThreadsUserState.Muted.SET(jet.Bool(req.State.Muted)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &SetThreadUserStateResponse{}, nil
}

func (s *Server) setUnreadState(ctx context.Context, threadId uint64, userIds []int32) error {
	if len(userIds) == 0 {
		return nil
	}

	tThreadsUserState := table.FivenetMsgsThreadsUserState
	stmt := tThreadsUserState.
		INSERT(
			tThreadsUserState.ThreadID,
			tThreadsUserState.UserID,
			tThreadsUserState.Unread,
		)

	for _, userId := range userIds {
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
