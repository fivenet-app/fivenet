package mailer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const MessagesDefaultPageSize = 20

var tMessages = table.FivenetMailerMessages.AS("message")

func (s *Server) ListThreadMessages(ctx context.Context, req *ListThreadMessagesRequest) (*ListThreadMessagesResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.mailer.thread.id", int64(req.ThreadId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(ctx, userInfo, req.ThreadId, req.EmailId, mailer.AccessLevel_ACCESS_LEVEL_READ); err != nil {
		return nil, err
	}

	countStmt := tMessages.
		SELECT(
			jet.COUNT(jet.DISTINCT(tMessages.ID)).AS("datacount.totalcount"),
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

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, MessagesDefaultPageSize)
	resp := &ListThreadMessagesResponse{
		Pagination: pag,
		Messages:   []*mailer.Message{},
	}
	if count.TotalCount <= 0 {
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
	for i := 0; i < len(resp.Messages); i++ {
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

func (s *Server) PostMessage(ctx context.Context, req *PostMessageRequest) (*PostMessageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
		Method:  "PostMessage",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
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
	if !userInfo.SuperUser && senderEmail.Deactivated {
		return nil, errorsmailer.ErrEmailDisabled
	}

	var emails []*mailer.ThreadRecipientEmail
	if len(req.Recipients) > 0 {
		emails, err = s.resolveRecipientsToEmails(ctx, senderEmail, req.Recipients)
		if err != nil {
			return nil, err
		}
	}

	req.Message.Sender = senderEmail
	req.Message.CreatorId = &userInfo.UserId
	req.Message.CreatorJob = &userInfo.Job

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

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	message, err := s.getMessage(ctx, req.Message.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	recipients, err := s.getThreadRecipients(ctx, req.Message.ThreadId)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	if len(recipients) > 0 {
		emailIds := []uint64{}
		for _, ua := range recipients {
			emailIds = append(emailIds, ua.EmailId)
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_MessageUpdate{
				MessageUpdate: message,
			},
		}, emailIds...)

		if err := s.setUnreadState(ctx, message.ThreadId, emailIds); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	return &PostMessageResponse{
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

func (s *Server) DeleteMessage(ctx context.Context, req *DeleteMessageRequest) (*DeleteMessageResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.mailer.message.id", int64(req.MessageId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
		Method:  "DeleteMessage",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if !userInfo.SuperUser {
		return nil, errorsmailer.ErrFailedQuery
	}

	stmt := tMessages.
		DELETE().
		WHERE(jet.AND(
			tMessages.ThreadID.EQ(jet.Uint64(req.ThreadId)),
			tMessages.ID.EQ(jet.Uint64(req.MessageId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	thread, err := s.getThread(ctx, req.ThreadId, req.EmailId, userInfo, true)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
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

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteMessageResponse{}, nil
}

func (s *Server) SearchThreads(ctx context.Context, req *SearchThreadsRequest) (*SearchThreadsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
		Method:  "SearchThreads",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	req.Search = strings.TrimRight(req.Search, "*")
	if req.Search == "" {
		return &SearchThreadsResponse{
			Pagination: &database.PaginationResponse{
				PageSize: 15,
			},
		}, nil
	}
	req.Search += "*"

	all := true
	listEmailsResp, err := s.ListEmails(ctx, &ListEmailsRequest{
		Pagination: &database.PaginationRequest{},
		All:        &all,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if listEmailsResp == nil || len(listEmailsResp.Emails) == 0 {
		return &SearchThreadsResponse{
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
			jet.COUNT(jet.DISTINCT(tMessages.ID)).AS("datacount.totalcount"),
		).
		FROM(tMessages).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 15)
	resp := &SearchThreadsResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
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

	fmt.Println(stmt.DebugSql())

	if err := stmt.QueryContext(ctx, s.db, &resp.Messages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Messages); i++ {
		if resp.Messages[i].Sender != nil {
			jobInfoFn(resp.Messages[i].Sender)
		}
	}

	resp.Pagination.Update(len(resp.Messages))

	return resp, nil
}
