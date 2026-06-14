package mailer

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	mailermessages "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/messages"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tMessages = table.FivenetMailerMessages.AS("message")

type MessageListQuery struct {
	ThreadID int64
	Offset   int64
	Limit    int64
}

func (s *Store) CountThreadMessages(ctx context.Context, q qrm.DB, threadID int64) (int64, error) {
	countStmt := tMessages.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tMessages.ID)).AS("data_count.total"),
		).
		FROM(tMessages).
		WHERE(mysql.AND(
			tMessages.DeletedAt.IS_NULL(),
			tMessages.ThreadID.EQ(mysql.Int64(threadID)),
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.dbOr(q), &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListThreadMessages(
	ctx context.Context,
	q qrm.DB,
	in MessageListQuery,
) ([]*mailermessages.Message, error) {
	stmt := tMessages.
		SELECT(
			tMessages.ID,
			tMessages.ThreadID,
			tMessages.SenderID,
			tMessages.CreatedAt,
			tMessages.UpdatedAt,
			tMessages.DeletedAt,
			tMessages.Title,
			tMessages.Content,
			tMessages.Data,
			tMessages.CreatorID,
			tMessages.SenderID.AS("sender.id"),
			tMessages.CreatorEmail.AS("sender.email"),
		).
		FROM(
			tMessages.
				LEFT_JOIN(tEmails,
					tEmails.ID.EQ(tMessages.SenderID),
				),
		).
		WHERE(mysql.AND(
			tMessages.DeletedAt.IS_NULL(),
			tMessages.ThreadID.EQ(mysql.Int64(in.ThreadID)),
			tEmails.DeletedAt.IS_NULL(),
		)).
		OFFSET(in.Offset).
		ORDER_BY(tMessages.CreatedAt.DESC()).
		LIMIT(in.Limit)

	messages := []*mailermessages.Message{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), &messages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return messages, nil
}

func (s *Store) GetMessage(
	ctx context.Context,
	q qrm.DB,
	messageID int64,
) (*mailermessages.Message, error) {
	stmt := tMessages.
		SELECT(
			tMessages.ID,
			tMessages.ThreadID,
			tMessages.SenderID,
			tMessages.CreatedAt,
			tMessages.UpdatedAt,
			tMessages.DeletedAt,
			tMessages.Title,
			tMessages.Content,
			tMessages.Data,
			tMessages.CreatorID,
			tMessages.SenderID.AS("sender.id"),
			tMessages.CreatorEmail.AS("sender.email"),
		).
		FROM(tMessages).
		WHERE(tMessages.ID.EQ(mysql.Int64(messageID))).
		LIMIT(1)

	message := &mailermessages.Message{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), message); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if message.GetId() == 0 {
		return nil, nil
	}

	return message, nil
}

func (s *Store) CreateMessage(
	ctx context.Context,
	q qrm.DB,
	msg *mailermessages.Message,
) (int64, error) {
	stmt := tMessages.
		INSERT(
			tMessages.ThreadID,
			tMessages.SenderID,
			tMessages.Title,
			tMessages.Content,
			tMessages.Data,
			tMessages.CreatorID,
			tMessages.CreatorJob,
			tMessages.CreatorEmail,
		).
		VALUES(
			msg.GetThreadId(),
			msg.GetSenderId(),
			msg.GetTitle(),
			msg.GetContent(),
			msg.GetData(),
			msg.CreatorId,
			msg.CreatorJob,
			msg.GetSender().GetEmail(),
		)

	res, err := stmt.ExecContext(ctx, s.dbOr(q))
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}
