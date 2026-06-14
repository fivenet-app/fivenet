package mailerstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	mailerthreads "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/threads"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tThreads           = table.FivenetMailerThreads.AS("thread")
	tThreadsState      = table.FivenetMailerThreadsState.AS("thread_state")
	tThreadsRecipients = table.FivenetMailerThreadsRecipients
)

type ThreadListQuery struct {
	EmailIDs  []int64
	Unread    *bool
	Archived  *bool
	Superuser bool
	Offset    int64
	Limit     int64
}

func (s *Store) CountThreads(ctx context.Context, q qrm.DB, in ThreadListQuery) (int64, error) {
	if len(in.EmailIDs) == 0 {
		return 0, nil
	}

	condition, recipientExists := buildThreadFilters(in)

	countStmt := tThreads.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tThreads.ID)).AS("data_count.total"),
		).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsState,
					mysql.AND(
						tThreadsState.ThreadID.EQ(tThreads.ID),
						tThreadsState.EmailID.EQ(mysql.Int64(in.EmailIDs[0])),
					),
				),
		).
		WHERE(mysql.AND(condition, recipientExists))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.dbOr(q), &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListThreads(
	ctx context.Context,
	q qrm.DB,
	in ThreadListQuery,
) ([]*mailerthreads.Thread, error) {
	if len(in.EmailIDs) == 0 {
		return []*mailerthreads.Thread{}, nil
	}

	condition, recipientExists := buildThreadFilters(in)

	page := tThreads.
		SELECT(tThreads.ID.AS("id")).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsState,
					mysql.AND(
						tThreadsState.ThreadID.EQ(tThreads.ID),
						tThreadsState.EmailID.EQ(mysql.Int64(in.EmailIDs[0])),
					),
				),
		).
		WHERE(mysql.AND(condition, recipientExists)).
		ORDER_BY(
			mysql.COALESCE(tThreads.UpdatedAt, tThreads.CreatedAt).DESC(),
			tThreads.ID.DESC(),
		).
		OFFSET(in.Offset).
		LIMIT(in.Limit).
		AsTable("page")

	stmt := page.
		SELECT(
			tThreads.ID,
			tThreads.CreatedAt,
			tThreads.UpdatedAt,
			tThreads.DeletedAt,
			tThreads.Title,
			tThreads.CreatorEmailID,
			tThreads.CreatorEmail,
			tThreads.CreatorEmailID.AS("email.id"),
			tThreads.CreatorEmail.AS("email.email"),
			tThreads.CreatorID,
			tThreadsState.ThreadID,
			tThreadsState.EmailID,
			tThreadsState.Unread,
			tThreadsState.LastRead,
			tThreadsState.Important,
			tThreadsState.Favorite,
			tThreadsState.Muted,
			tThreadsState.Archived,
		).
		FROM(
			page.
				INNER_JOIN(tThreads,
					tThreads.ID.EQ(mysql.RawInt("page.id")),
				).
				LEFT_JOIN(tThreadsState,
					mysql.AND(
						tThreadsState.ThreadID.EQ(tThreads.ID),
						tThreadsState.EmailID.EQ(mysql.Int64(in.EmailIDs[0])),
					),
				),
		).
		ORDER_BY(
			mysql.COALESCE(tThreads.UpdatedAt, tThreads.CreatedAt).DESC(),
			tThreads.ID.DESC(),
		)

	threads := []*mailerthreads.Thread{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), &threads); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return threads, nil
}

func (s *Store) GetThread(
	ctx context.Context,
	q qrm.DB,
	threadID int64,
	emailID int64,
) (*mailerthreads.Thread, error) {
	stmt := tThreads.
		SELECT(
			tThreads.ID,
			tThreads.CreatedAt,
			tThreads.UpdatedAt,
			tThreads.DeletedAt,
			tThreads.Title,
			tThreads.CreatorEmailID,
			tThreads.CreatorID,
			tThreads.CreatorEmailID.AS("email.id"),
			tThreads.CreatorEmail.AS("email.email"),
			tThreadsState.ThreadID,
			tThreadsState.EmailID,
			tThreadsState.Unread,
			tThreadsState.LastRead,
			tThreadsState.Important,
			tThreadsState.Favorite,
			tThreadsState.Muted,
			tThreadsState.Archived,
		).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsState,
					mysql.AND(
						tThreadsState.ThreadID.EQ(tThreads.ID),
						tThreadsState.EmailID.EQ(mysql.Int64(emailID)),
					),
				),
		).
		WHERE(tThreads.ID.EQ(mysql.Int64(threadID))).
		LIMIT(1)

	var thread mailerthreads.Thread
	if err := stmt.QueryContext(ctx, s.dbOr(q), &thread); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if thread.GetId() == 0 {
		return nil, nil
	}

	return &thread, nil
}

func (s *Store) UpdateThreadTime(ctx context.Context, q qrm.DB, threadID int64) error {
	stmt := tThreads.
		UPDATE(
			tThreads.UpdatedAt,
		).
		SET(
			tThreads.UpdatedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(tThreads.ID.EQ(mysql.Int64(threadID))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.dbOr(q)); err != nil {
		return err
	}

	return nil
}

func (s *Store) AddThreadRecipients(
	ctx context.Context,
	q qrm.DB,
	threadID int64,
	recipients []*mailerthreads.ThreadRecipientEmail,
) error {
	if len(recipients) == 0 {
		return nil
	}

	stmt := tThreadsRecipients.
		INSERT(
			tThreadsRecipients.ThreadID,
			tThreadsRecipients.EmailID,
			tThreadsRecipients.Email,
		)

	for _, recipient := range recipients {
		stmt = stmt.VALUES(threadID, recipient.GetEmailId(), recipient.GetEmail().GetEmail())
	}

	if _, err := stmt.ExecContext(ctx, s.dbOr(q)); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (s *Store) ListThreadRecipients(
	ctx context.Context,
	q qrm.DB,
	threadID int64,
) ([]*mailerthreads.ThreadRecipientEmail, error) {
	tRecipients := tThreadsRecipients.AS("thread_recipient_email")
	stmt := tRecipients.
		SELECT(
			tRecipients.ID,
			tRecipients.ThreadID,
			tRecipients.EmailID,
			tEmails.ID,
			tRecipients.Email.AS("email.email"),
		).
		FROM(
			tRecipients.
				INNER_JOIN(tEmails,
					mysql.AND(
						tRecipients.EmailID.EQ(tEmails.ID),
						tEmails.Deactivated.IS_FALSE(),
					),
				),
		).
		WHERE(mysql.AND(
			tRecipients.ThreadID.EQ(mysql.Int64(threadID)),
			tEmails.DeletedAt.IS_NULL(),
		))

	recipients := []*mailerthreads.ThreadRecipientEmail{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), &recipients); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return recipients, nil
}

func buildThreadFilters(in ThreadListQuery) (mysql.BoolExpression, mysql.BoolExpression) {
	if len(in.EmailIDs) == 0 {
		return mysql.Bool(false), mysql.Bool(false)
	}

	condition := mysql.Bool(true)
	if !in.Superuser {
		condition = tThreads.DeletedAt.IS_NULL()
	}

	if in.Unread != nil {
		condition = condition.AND(
			mysql.AND(
				tThreadsState.Unread.IS_NOT_NULL(),
				tThreadsState.Unread.EQ(mysql.Bool(*in.Unread)),
			),
		)
	}
	if in.Archived != nil {
		condition = condition.AND(
			mysql.AND(
				tThreadsState.Archived.IS_NOT_NULL(),
				tThreadsState.Archived.EQ(mysql.Bool(*in.Archived)),
			),
		)
	} else {
		condition = condition.AND(
			tThreadsState.Archived.IS_NULL().OR(tThreadsState.Archived.EQ(mysql.Bool(false))),
		)
	}

	emailIDExprs := make([]mysql.Expression, 0, len(in.EmailIDs))
	for _, emailID := range in.EmailIDs {
		emailIDExprs = append(emailIDExprs, mysql.Int64(emailID))
	}

	recipientExists := mysql.EXISTS(
		mysql.
			SELECT(mysql.Int(1)).
			FROM(tThreadsRecipients).
			WHERE(mysql.AND(
				tThreadsRecipients.ThreadID.EQ(tThreads.ID),
				tThreadsRecipients.EmailID.IN(emailIDExprs...),
			)),
	)

	return condition, recipientExists
}
