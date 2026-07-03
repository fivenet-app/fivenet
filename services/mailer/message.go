package mailer

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	maileraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/access"
	mailerevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/events"
	mailermessages "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/messages"
	mailerthreads "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/threads"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbmailer "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2026/services/mailer/errors"
	mailerstore "github.com/fivenet-app/fivenet/v2026/stores/mailer"
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

	if err := s.checkIfEmailPartOfThread(
		ctx,
		userInfo,
		req.GetThreadId(),
		req.GetEmailId(),
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
	); err != nil {
		return nil, err
	}

	count, err := s.store.CountThreadMessages(ctx, s.db, req.GetThreadId(), userInfo.GetJobAdmin())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, MessagesDefaultPageSize)
	resp := &pbmailer.ListThreadMessagesResponse{
		Pagination: pag,
		Messages:   []*mailermessages.Message{},
	}
	if count <= 0 {
		return resp, nil
	}
	messages, err := s.store.ListThreadMessages(ctx, s.db, mailerstore.MessageListQuery{
		ThreadID: req.GetThreadId(),
		Offset:   req.GetPagination().GetOffset(),
		Limit:    limit,
	}, userInfo.GetJobAdmin())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	resp.Messages = messages

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetMessages() {
		if resp.GetMessages()[i].GetSender() != nil {
			jobInfoFn(resp.GetMessages()[i].GetSender())
		}
	}

	return resp, nil
}

func (s *Server) PostMessage(
	ctx context.Context,
	req *pbmailer.PostMessageRequest,
) (*pbmailer.PostMessageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(
		ctx,
		userInfo,
		req.GetMessage().GetThreadId(),
		req.GetMessage().GetSenderId(),
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_WRITE),
	); err != nil {
		return nil, err
	}

	senderEmail, err := s.getEmail(ctx, userInfo, req.GetMessage().GetSenderId(), false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Prevent disabled emails from sending messages
	if !userInfo.GetJobAdmin() && senderEmail.GetDeactivated() {
		return nil, errorsmailer.ErrEmailDisabled
	}

	var emails []*mailerthreads.ThreadRecipientEmail
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
		if a, ok := attachment.GetData().(*mailermessages.MessageAttachment_Document); ok {
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

	lastId, err := s.store.CreateMessage(ctx, tx, req.GetMessage())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	req.Message.Id = lastId

	if len(emails) > 0 {
		if err := s.store.AddThreadRecipients(
			ctx,
			tx,
			req.GetMessage().GetThreadId(),
			emails,
		); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	if err := s.store.UpdateThreadTime(ctx, tx, req.GetMessage().GetThreadId()); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	recipients, err := s.store.ListThreadRecipients(
		ctx,
		tx,
		req.GetMessage().GetThreadId(),
		userInfo.GetJobAdmin(),
	)
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

	if err := s.store.SetUnreadState(
		ctx,
		tx,
		req.GetMessage().GetThreadId(),
		senderEmail.GetId(),
		emailIds,
	); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	message, err := s.store.GetMessage(ctx, s.db, req.GetMessage().GetId(), userInfo.GetJobAdmin())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	s.sendUpdate(ctx, &mailerevents.MailerEvent{
		Data: &mailerevents.MailerEvent_MessageUpdate{
			MessageUpdate: message,
		},
	}, emailIds...)

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return &pbmailer.PostMessageResponse{
		Message: message,
	}, nil
}

func (s *Server) DeleteMessage(
	ctx context.Context,
	req *pbmailer.DeleteMessageRequest,
) (*pbmailer.DeleteMessageResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.mailer.message_id", req.GetMessageId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if !userInfo.GetJobAdmin() {
		return nil, errorsmailer.ErrFailedQuery
	}

	message, err := s.store.GetMessage(ctx, s.db, req.GetMessageId(), userInfo.GetJobAdmin())
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	var deletedAtTime *timestamp.Timestamp
	if message == nil || message.GetDeletedAt() == nil || !userInfo.GetJobAdmin() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteMessage(
		ctx,
		s.db,
		req.GetThreadId(),
		req.GetMessageId(),
		deletedAtTime,
	); err != nil {
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

		s.sendUpdate(ctx, &mailerevents.MailerEvent{
			Data: &mailerevents.MailerEvent_MessageDelete{
				MessageDelete: req.GetMessageId(),
			},
		}, emailIds...)
	}

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
	includeDeleted := userInfo.GetJobAdmin()

	// Get Thread ids via threads recipients list
	condition := mysql.AND(
		mysql.OR(
			mysql.Bool(includeDeleted),
			tMessages.DeletedAt.IS_NULL(),
		),
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
