package mailer

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const MessagesDefaultPageSize = 20

var tMessages = table.FivenetMailerMessages.AS("message")

func (s *Server) ListThreadMessages(
	ctx context.Context,
	req *pbmailer.ListThreadMessagesRequest,
) (*pbmailer.ListThreadMessagesResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.mailer.thread.id", req.GetThreadId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(ctx, userInfo, req.GetThreadId(), req.GetEmailId(), mailer.AccessLevel_ACCESS_LEVEL_READ); err != nil {
		return nil, err
	}

	countStmt := tMessages.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tMessages.ID)).AS("data_count.total"),
		).
		FROM(tMessages).
		WHERE(mysql.AND(
			tMessages.DeletedAt.IS_NULL(),
			tMessages.ThreadID.EQ(mysql.Int64(req.GetThreadId())),
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, MessagesDefaultPageSize)
	resp := &pbmailer.ListThreadMessagesResponse{
		Pagination: pag,
		Messages:   []*mailer.Message{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

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
			tMessages.ThreadID.EQ(mysql.Int64(req.GetThreadId())),
			tEmails.DeletedAt.IS_NULL(),
		)).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(tMessages.CreatedAt.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Messages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetMessages() {
		if resp.GetMessages()[i].GetSender() != nil {
			jobInfoFn(resp.GetMessages()[i].GetSender())
		}
	}

	return resp, nil
}

func (s *Server) getMessage(ctx context.Context, messageId int64) (*mailer.Message, error) {
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
		WHERE(
			tMessages.ID.EQ(mysql.Int64(messageId)),
		).
		LIMIT(1)

	message := &mailer.Message{}
	if err := stmt.QueryContext(ctx, s.db, message); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if message.GetId() == 0 {
		return nil, nil
	}

	return message, nil
}

func (s *Server) PostMessage(
	ctx context.Context,
	req *pbmailer.PostMessageRequest,
) (*pbmailer.PostMessageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(ctx, userInfo, req.GetMessage().GetThreadId(), req.GetMessage().GetSenderId(), mailer.AccessLevel_ACCESS_LEVEL_WRITE); err != nil {
		return nil, err
	}

	senderEmail, err := s.getEmail(ctx, req.GetMessage().GetSenderId(), false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Prevent disabled emails from sending messages
	if !userInfo.GetSuperuser() && senderEmail.GetDeactivated() {
		return nil, errorsmailer.ErrEmailDisabled
	}

	var emails []*mailer.ThreadRecipientEmail
	if len(req.GetRecipients()) > 0 {
		emails, err = s.retrieveRecipientsToEmails(ctx, senderEmail, req.GetRecipients())
		if err != nil {
			return nil, err
		}
	}

	req.Message.Sender = senderEmail
	req.Message.CreatorId = &userInfo.UserId
	req.Message.CreatorJob = &userInfo.Job

	// Remove titles from attached documents
	for _, attachment := range req.GetMessage().GetData().GetAttachments() {
		if a, ok := attachment.GetData().(*mailer.MessageAttachment_Document); ok {
			a.Document.Title = nil
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	lastId, err := s.createMessage(ctx, tx, req.GetMessage())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	req.Message.Id = lastId

	if len(emails) > 0 {
		if err := s.handleRecipientsChanges(ctx, tx, req.GetMessage().GetThreadId(), emails); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	if err := s.updateThreadTime(ctx, tx, req.GetMessage().GetThreadId()); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	recipients, err := s.getThreadRecipients(ctx, tx, req.GetMessage().GetThreadId())
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	emailIds := []int64{}
	for _, ua := range recipients {
		// Skip sender email id
		if ua.GetEmailId() == senderEmail.GetId() {
			continue
		}

		emailIds = append(emailIds, ua.GetEmailId())
	}

	if err := s.setUnreadState(ctx, tx, req.GetMessage().GetThreadId(), senderEmail.GetId(), emailIds); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	message, err := s.getMessage(ctx, req.GetMessage().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_MessageUpdate{
			MessageUpdate: message,
		},
	}, emailIds...)

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return &pbmailer.PostMessageResponse{
		Message: message,
	}, nil
}

func (s *Server) createMessage(
	ctx context.Context,
	tx qrm.DB,
	msg *mailer.Message,
) (int64, error) {
	tMessages := table.FivenetMailerMessages
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

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return lastId, nil
}

func (s *Server) DeleteMessage(
	ctx context.Context,
	req *pbmailer.DeleteMessageRequest,
) (*pbmailer.DeleteMessageResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.mailer.message_id", req.GetMessageId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if !userInfo.GetSuperuser() {
		return nil, errorsmailer.ErrFailedQuery
	}

	message, err := s.getMessage(ctx, req.GetMessageId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	deletedAtTime := mysql.CURRENT_TIMESTAMP()
	if message != nil && message.GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = mysql.TimestampExp(mysql.NULL)
	}

	tMessages := table.FivenetMailerMessages
	stmt := tMessages.
		UPDATE(
			tMessages.DeletedAt,
		).
		SET(
			tMessages.DeletedAt.SET(deletedAtTime),
		).
		WHERE(mysql.AND(
			tMessages.ThreadID.EQ(mysql.Int64(req.GetThreadId())),
			tMessages.ID.EQ(mysql.Int64(req.GetMessageId())),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	thread, err := s.getThread(ctx, req.GetThreadId(), req.GetEmailId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if thread != nil && thread.Recipients != nil && len(thread.GetRecipients()) > 0 {
		emailIds := []int64{}
		for _, ua := range thread.GetRecipients() {
			emailIds = append(emailIds, ua.GetEmailId())
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_MessageDelete{
				MessageDelete: req.GetMessageId(),
			},
		}, emailIds...)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbmailer.DeleteMessageResponse{}, nil
}

func (s *Server) SearchThreads(
	ctx context.Context,
	req *pbmailer.SearchThreadsRequest,
) (*pbmailer.SearchThreadsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	req.Search = strings.TrimRight(req.GetSearch(), "*")
	if req.GetSearch() == "" {
		return &pbmailer.SearchThreadsResponse{
			Pagination: &database.PaginationResponse{
				PageSize: 15,
			},
		}, nil
	}
	req.Search += "*"

	all := true
	listEmailsResp, err := s.ListEmails(ctx, &pbmailer.ListEmailsRequest{
		Pagination: &database.PaginationRequest{},
		All:        &all,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if listEmailsResp == nil || len(listEmailsResp.GetEmails()) == 0 {
		return &pbmailer.SearchThreadsResponse{
			Pagination: &database.PaginationResponse{},
		}, nil
	}

	ids := []mysql.Expression{}
	for _, email := range listEmailsResp.GetEmails() {
		ids = append(ids, mysql.Int64(email.GetId()))
	}

	// Get Thread ids via threads recipients list
	condition := mysql.AND(
		tMessages.DeletedAt.IS_NULL(),
		tMessages.ThreadID.IN(
			tThreadsRecipients.
				SELECT(
					mysql.DISTINCT(tThreadsRecipients.ThreadID),
				).
				FROM(tThreadsRecipients).
				WHERE(
					tThreadsRecipients.EmailID.IN(ids...),
				),
		),
		mysql.OR(
			dbutils.MATCH(tMessages.Title, mysql.String(req.GetSearch())),
			dbutils.MATCH(tMessages.Content, mysql.String(req.GetSearch())),
		),
	)

	countStmt := tMessages.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tMessages.ID)).AS("data_count.total"),
		).
		FROM(tMessages).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 15)
	resp := &pbmailer.SearchThreadsResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	stmt := tMessages.
		SELECT(
			tMessages.ID,
			tMessages.ThreadID,
			tMessages.SenderID,
			tMessages.Title,
			tMessages.CreatorID,
			tMessages.SenderID.AS("sender.id"),
			tMessages.CreatorEmail.AS("sender.email"),
		).
		FROM(tMessages).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(tMessages.CreatedAt.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Messages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetMessages() {
		if resp.GetMessages()[i].GetSender() != nil {
			jobInfoFn(resp.GetMessages()[i].GetSender())
		}
	}

	return resp, nil
}
