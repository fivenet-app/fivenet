package mailer

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) ListThreadMessages(ctx context.Context, req *ListThreadMessagesRequest) (*ListThreadMessagesResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.mailer.thread.id", int64(req.ThreadId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.ThreadId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmailer.ErrFailedQuery
		}
	}

	resp := &ListThreadMessagesResponse{
		Messages: []*mailer.Message{},
	}

	stmt := tMessages.
		SELECT(
			tMessages.ID,
			tMessages.ThreadID,
			tMessages.SenderID,
			tMessages.CreatedAt,
			tMessages.UpdatedAt,
			tMessages.DeletedAt,
			tMessages.Content,
			tMessages.Data,
			tMessages.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("sender_user.avatar"),
		).
		FROM(
			tMessages.
				LEFT_JOIN(tEmails,
					tEmails.ID.EQ(tMessages.SenderID),
				),
		).
		LIMIT(20)

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

	return resp, nil
}

func (s *Server) getMessage(ctx context.Context, messageId uint64, userInfo *userinfo.UserInfo) (*mailer.Message, error) {
	stmt := tMessages.
		SELECT(
			tMessages.ID,
			tMessages.ThreadID,
			tMessages.CreatedAt,
			tMessages.UpdatedAt,
			tMessages.DeletedAt,
			tMessages.Content,
			tMessages.Data,
			tMessages.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
		).
		FROM(
			tMessages.
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tMessages.SenderID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCreator.ID),
				),
		).
		WHERE(tMessages.ID.EQ(jet.Uint64(messageId))).
		LIMIT(1)

	var message mailer.Message
	if err := stmt.QueryContext(ctx, s.db, &message); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if message.Id == 0 {
		return nil, nil
	}

	if message.Sender != nil {
		// TODO
		// s.enricher.EnrichJobInfoSafe(userInfo, message.Sender)
	}

	return &message, nil
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

	check, err := s.access.CanUserAccessTarget(ctx, req.Message.ThreadId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmailer.ErrFailedQuery
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

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	message, err := s.getMessage(ctx, req.Message.Id, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	thread, err := s.getThread(ctx, req.Message.ThreadId, userInfo, true)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	if thread != nil && thread.Recipients != nil && len(thread.Recipients) > 0 {
		emailIds := []uint64{}
		for _, ua := range thread.Recipients {
			emailIds = append(emailIds, ua.EmailId)
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_MessageUpdate{
				MessageUpdate: message,
			},
		}, emailIds)

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
		).
		VALUES(
			msg.ThreadId,
			msg.SenderId,
			msg.Title,
			msg.Content,
			msg.Data,
			msg.CreatorId,
			msg.CreatorJob,
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

	check, err := s.access.CanUserAccessTarget(ctx, req.ThreadId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmailer.ErrFailedQuery
		}
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

	thread, err := s.getThread(ctx, req.ThreadId, userInfo, true)
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
		}, emailIds)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteMessageResponse{}, nil
}
