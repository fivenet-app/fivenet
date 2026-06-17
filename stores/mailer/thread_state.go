package mailerstore

import (
	"context"
	"errors"

	mailerthreads "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/threads"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) GetThreadState(
	ctx context.Context,
	q qrm.DB,
	threadID int64,
	emailID int64,
) (*mailerthreads.ThreadState, error) {
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
		WHERE(mysql.AND(
			tThreadsState.ThreadID.EQ(mysql.Int64(threadID)),
			tThreadsState.EmailID.EQ(mysql.Int64(emailID)),
		))

	dest := &mailerthreads.ThreadState{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetThreadId() == 0 || dest.GetEmailId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Store) SetThreadState(
	ctx context.Context,
	q qrm.DB,
	state *mailerthreads.ThreadState,
) error {
	updateSets := []mysql.ColumnAssigment{}
	if state.Unread != nil {
		updateSets = append(updateSets, tThreadsState.Unread.SET(mysql.RawBool("VALUES(`unread`)")))
	}
	if state.GetLastRead() != nil {
		updateSets = append(
			updateSets,
			tThreadsState.LastRead.SET(mysql.RawTimestamp("VALUES(`last_read`)")),
		)
	}
	if state.Important != nil {
		updateSets = append(
			updateSets,
			tThreadsState.Important.SET(mysql.RawBool("VALUES(`important`)")),
		)
	}
	if state.Favorite != nil {
		updateSets = append(
			updateSets,
			tThreadsState.Favorite.SET(mysql.RawBool("VALUES(`favorite`)")),
		)
	}
	if state.Muted != nil {
		updateSets = append(updateSets, tThreadsState.Muted.SET(mysql.RawBool("VALUES(`muted`)")))
	}
	if state.Archived != nil {
		updateSets = append(
			updateSets,
			tThreadsState.Archived.SET(mysql.RawBool("VALUES(`archived`)")),
		)
	}

	if len(updateSets) == 0 {
		return nil
	}

	tThreadsStateTbl := table.FivenetMailerThreadsState
	stmt := tThreadsStateTbl.
		INSERT(
			tThreadsStateTbl.ThreadID,
			tThreadsStateTbl.EmailID,
			tThreadsStateTbl.Unread,
			tThreadsStateTbl.LastRead,
			tThreadsStateTbl.Important,
			tThreadsStateTbl.Favorite,
			tThreadsStateTbl.Muted,
			tThreadsStateTbl.Archived,
		).
		VALUES(
			state.GetThreadId(),
			state.GetEmailId(),
			state.Unread,
			state.GetLastRead(),
			state.Important,
			state.Favorite,
			state.Muted,
			state.Archived,
		).
		ON_DUPLICATE_KEY_UPDATE(updateSets...)

	if _, err := stmt.ExecContext(ctx, s.dbOr(q)); err != nil {
		return err
	}

	return nil
}

func (s *Store) SetUnreadState(
	ctx context.Context,
	q qrm.DB,
	threadID int64,
	senderID int64,
	emailIDs []int64,
) error {
	if len(emailIDs) == 0 {
		return nil
	}

	stmt := table.FivenetMailerThreadsState.
		INSERT(
			table.FivenetMailerThreadsState.ThreadID,
			table.FivenetMailerThreadsState.EmailID,
			table.FivenetMailerThreadsState.Unread,
		)

	for _, emailID := range emailIDs {
		stmt = stmt.VALUES(threadID, emailID, emailID != senderID)
	}

	stmt = stmt.
		ON_DUPLICATE_KEY_UPDATE(
			table.FivenetMailerThreadsState.Unread.SET(mysql.RawBool("VALUES(`unread`)")),
		)

	if _, err := stmt.ExecContext(ctx, s.dbOr(q)); err != nil {
		return err
	}

	return nil
}
