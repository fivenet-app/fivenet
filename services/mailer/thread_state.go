package mailer

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) GetThreadState(
	ctx context.Context,
	req *pbmailer.GetThreadStateRequest,
) (*pbmailer.GetThreadStateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(ctx, userInfo, req.GetThreadId(), req.GetEmailId(), mailer.AccessLevel_ACCESS_LEVEL_READ); err != nil {
		return nil, err
	}

	state, err := s.getThreadState(ctx, req.GetThreadId(), req.GetEmailId())
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

	if err := s.checkIfEmailPartOfThread(ctx, userInfo, req.GetState().GetThreadId(), req.GetState().GetEmailId(), mailer.AccessLevel_ACCESS_LEVEL_WRITE); err != nil {
		return nil, err
	}

	updateSets := []jet.ColumnAssigment{}
	if req.State.Unread != nil {
		updateSets = append(updateSets, tThreadsState.Unread.SET(jet.RawBool("VALUES(`unread`)")))
	}
	if req.GetState().GetLastRead() != nil {
		updateSets = append(
			updateSets,
			tThreadsState.LastRead.SET(jet.RawTimestamp("VALUES(`last_read`)")),
		)
	}
	if req.State.Important != nil {
		updateSets = append(
			updateSets,
			tThreadsState.Important.SET(jet.RawBool("VALUES(`important`)")),
		)
	}
	if req.State.Favorite != nil {
		updateSets = append(
			updateSets,
			tThreadsState.Favorite.SET(jet.RawBool("VALUES(`favorite`)")),
		)
	}
	if req.State.Muted != nil {
		updateSets = append(updateSets, tThreadsState.Muted.SET(jet.RawBool("VALUES(`muted`)")))
	}
	if req.State.Archived != nil {
		updateSets = append(
			updateSets,
			tThreadsState.Archived.SET(jet.RawBool("VALUES(`archived`)")),
		)
	}

	if len(updateSets) > 0 {
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
				req.GetState().GetThreadId(),
				req.GetState().GetEmailId(),
				req.GetState().GetUnread(),
				req.GetState().GetLastRead(),
				req.GetState().GetImportant(),
				req.GetState().GetFavorite(),
				req.GetState().GetMuted(),
				req.GetState().GetArchived(),
			).
			ON_DUPLICATE_KEY_UPDATE(updateSets...)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	state, err := s.getThreadState(ctx, req.GetState().GetThreadId(), req.GetState().GetEmailId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_ThreadStateUpdate{
			ThreadStateUpdate: state,
		},
	}, req.GetState().GetEmailId())

	return &pbmailer.SetThreadStateResponse{
		State: state,
	}, nil
}

func (s *Server) getThreadState(
	ctx context.Context,
	threadId uint64,
	emaildId uint64,
) (*mailer.ThreadState, error) {
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

	if dest.GetThreadId() == 0 || dest.GetEmailId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Server) setUnreadState(
	ctx context.Context,
	tx qrm.DB,
	threadId uint64,
	senderId uint64,
	emailIds []uint64,
) error {
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

	for _, emailId := range emailIds {
		stmt = stmt.VALUES(
			threadId,
			emailId,
			emailId != senderId,
		)
	}

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tThreadsUserState.Unread.SET(jet.RawBool("VALUES(`unread`)")),
	)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
