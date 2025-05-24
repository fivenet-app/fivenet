package mailer

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const MessagesDefaultPageSize = 20

var tMessages = table.FivenetMailerMessages.AS("message")

func (s *Server) ListThreadMessages(ctx context.Context, req *pbmailer.ListThreadMessagesRequest) (*pbmailer.ListThreadMessagesResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.mailer.thread.id", int64(req.ThreadId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(ctx, userInfo, req.ThreadId, req.EmailId, mailer.AccessLevel_ACCESS_LEVEL_READ); err != nil {
		return nil, err
	}

	countStmt := tMessages.
		SELECT(
			jet.COUNT(jet.DISTINCT(tMessages.ID)).AS("data_count.total"),
		).
		FROM(tMessages).
		WHERE(jet.AND(
			tMessages.DeletedAt.IS_NULL(),
			tMessages.ThreadID.EQ(jet.Uint64(req.ThreadId)),
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, MessagesDefaultPageSize)
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
		WHERE(jet.AND(
			tMessages.DeletedAt.IS_NULL(),
			tMessages.ThreadID.EQ(jet.Uint64(req.ThreadId)),
			tEmails.DeletedAt.IS_NULL(),
		)).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(tMessages.CreatedAt.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Messages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Messages {
		if resp.Messages[i].Sender != nil {
			jobInfoFn(resp.Messages[i].Sender)
		}
	}

	resp.Pagination.Update(len(resp.Messages))

	return resp, nil
}

func (s *Server) getMessage(ctx context.Context, messageId uint64) (*mailer.Message, error) {
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
			tMessages.ID.EQ(jet.Uint64(messageId)),
		).
		LIMIT(1)

	message := &mailer.Message{}
	if err := stmt.QueryContext(ctx, s.db, message); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if message.Id == 0 {
		return nil, nil
	}

	return message, nil
}

func (s *Server) PostMessage(ctx context.Context, req *pbmailer.PostMessageRequest) (*pbmailer.PostMessageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "PostMessage",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if err := s.checkIfEmailPartOfThread(ctx, userInfo, req.Message.ThreadId, req.Message.SenderId, mailer.AccessLevel_ACCESS_LEVEL_WRITE); err != nil {
		return nil, err
	}

	senderEmail, err := s.getEmail(ctx, req.Message.SenderId, false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Prevent disabled emails from sending messages
	if !userInfo.Superuser && senderEmail.Deactivated {
		return nil, errorsmailer.ErrEmailDisabled
	}

	var emails []*mailer.ThreadRecipientEmail
	if len(req.Recipients) > 0 {
		emails, err = s.retrieveRecipientsToEmails(ctx, senderEmail, req.Recipients)
		if err != nil {
			return nil, err
		}
	}

	req.Message.Sender = senderEmail
	req.Message.CreatorId = &userInfo.UserId
	req.Message.CreatorJob = &userInfo.Job

	// Remove titles from attached documents
	for _, attachment := range req.Message.Data.Attachments {
		if a, ok := attachment.Data.(*mailer.MessageAttachment_Document); ok {
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

	lastId, err := s.createMessage(ctx, tx, req.Message)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	req.Message.Id = uint64(lastId)

	if len(emails) > 0 {
		if err := s.handleRecipientsChanges(ctx, tx, req.Message.ThreadId, emails); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	if err := s.updateThreadTime(ctx, tx, req.Message.ThreadId); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	recipients, err := s.getThreadRecipients(ctx, tx, req.Message.ThreadId)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	emailIds := []uint64{}
	for _, ua := range recipients {
		// Skip sender email id
		if ua.EmailId == senderEmail.Id {
			continue
		}

		emailIds = append(emailIds, ua.EmailId)
	}

	if err := s.setUnreadState(ctx, tx, req.Message.ThreadId, senderEmail.Id, emailIds); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	message, err := s.getMessage(ctx, req.Message.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailer.MailerEvent{
		Data: &mailer.MailerEvent_MessageUpdate{
			MessageUpdate: message,
		},
	}, emailIds...)

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return &pbmailer.PostMessageResponse{
		Message: message,
	}, nil
}

func (s *Server) createMessage(ctx context.Context, tx qrm.DB, msg *mailer.Message) (uint64, error) {
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
			msg.ThreadId,
			msg.SenderId,
			msg.Title,
			msg.Content,
			msg.Data,
			msg.CreatorId,
			msg.CreatorJob,
			msg.Sender.Email,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return uint64(lastId), nil
}

func (s *Server) DeleteMessage(ctx context.Context, req *pbmailer.DeleteMessageRequest) (*pbmailer.DeleteMessageResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.mailer.message.id", int64(req.MessageId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "DeleteMessage",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if !userInfo.Superuser {
		return nil, errorsmailer.ErrFailedQuery
	}

	message, err := s.getMessage(ctx, req.MessageId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if message != nil && message.DeletedAt != nil && userInfo.Superuser {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	tMessages := table.FivenetMailerMessages
	stmt := tMessages.
		UPDATE(
			tMessages.DeletedAt,
		).
		SET(
			tMessages.DeletedAt.SET(deletedAtTime),
		).
		WHERE(jet.AND(
			tMessages.ThreadID.EQ(jet.Uint64(req.ThreadId)),
			tMessages.ID.EQ(jet.Uint64(req.MessageId)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	thread, err := s.getThread(ctx, req.ThreadId, req.EmailId, userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if thread != nil && thread.Recipients != nil && len(thread.Recipients) > 0 {
		emailIds := []uint64{}
		for _, ua := range thread.Recipients {
			emailIds = append(emailIds, ua.EmailId)
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_MessageDelete{
				MessageDelete: req.MessageId,
			},
		}, emailIds...)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbmailer.DeleteMessageResponse{}, nil
}

func (s *Server) SearchThreads(ctx context.Context, req *pbmailer.SearchThreadsRequest) (*pbmailer.SearchThreadsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "SearchThreads",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	req.Search = strings.TrimRight(req.Search, "*")
	if req.Search == "" {
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
	if listEmailsResp == nil || len(listEmailsResp.Emails) == 0 {
		return &pbmailer.SearchThreadsResponse{
			Pagination: &database.PaginationResponse{},
		}, nil
	}

	ids := []jet.Expression{}
	for _, email := range listEmailsResp.Emails {
		ids = append(ids, jet.Uint64(email.Id))
	}

	// Get Thread ids via threads recipients list
	condition := tMessages.DeletedAt.IS_NULL().
		AND(tMessages.ThreadID.IN(
			tThreadsRecipients.
				SELECT(
					jet.DISTINCT(tThreadsRecipients.ThreadID),
				).
				FROM(tThreadsRecipients).
				WHERE(
					tThreadsRecipients.EmailID.IN(ids...),
				),
		)).
		AND(jet.OR(
			jet.BoolExp(
				jet.Raw("MATCH(`title`) AGAINST ($search IN BOOLEAN MODE)",
					jet.RawArgs{"$search": req.Search}),
			),
			jet.BoolExp(
				jet.Raw("MATCH(`content`) AGAINST ($search IN BOOLEAN MODE)",
					jet.RawArgs{"$search": req.Search}),
			),
		))

	countStmt := tMessages.
		SELECT(
			jet.COUNT(jet.DISTINCT(tMessages.ID)).AS("data_count.total"),
		).
		FROM(tMessages).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, 15)
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
		OFFSET(req.Pagination.Offset).
		ORDER_BY(tMessages.CreatedAt.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Messages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Messages {
		if resp.Messages[i].Sender != nil {
			jobInfoFn(resp.Messages[i].Sender)
		}
	}

	resp.Pagination.Update(len(resp.Messages))

	return resp, nil
}
