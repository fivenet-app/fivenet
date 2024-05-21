package messenger

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/messenger"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsmessenger "github.com/fivenet-app/fivenet/gen/go/proto/services/messenger/errors"
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

func (s *Server) GetThreadMessages(ctx context.Context, req *GetThreadMessagesRequest) (*GetThreadMessagesResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.messenger.thread.id", int64(req.ThreadId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.checkIfUserHasAccessToThread(ctx, req.ThreadId, userInfo, messenger.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmessenger.ErrFailedQuery
		}
	}

	resp := &GetThreadMessagesResponse{
		Messages: []*messenger.Message{},
	}

	stmt := tMessages.
		SELECT(
			tMessages.ID,
			tMessages.ThreadID,
			tMessages.CreatedAt,
			tMessages.UpdatedAt,
			tMessages.DeletedAt,
			tMessages.Message,
			tMessages.Data,
			tMessages.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
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
					tMessages.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tMessages.CreatorID),
				),
		).
		LIMIT(20)

	if err := stmt.QueryContext(ctx, s.db, &resp.Messages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Messages); i++ {
		if resp.Messages[i].Creator != nil {
			jobInfoFn(resp.Messages[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) getMessage(ctx context.Context, messageId uint64, userInfo *userinfo.UserInfo) (*messenger.Message, error) {
	stmt := tMessages.
		SELECT(
			tMessages.ID,
			tMessages.ThreadID,
			tMessages.CreatedAt,
			tMessages.UpdatedAt,
			tMessages.DeletedAt,
			tMessages.Message,
			tMessages.Data,
			tMessages.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
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
					tCreator.ID.EQ(tMessages.CreatorID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCreator.ID),
				),
		).
		WHERE(tMessages.ID.EQ(jet.Uint64(messageId))).
		LIMIT(1)

	var message messenger.Message
	if err := stmt.QueryContext(ctx, s.db, &message); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if message.Id == 0 {
		return nil, nil
	}

	if message.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, message.Creator)
	}

	return &message, nil
}

func (s *Server) PostMessage(ctx context.Context, req *PostMessageRequest) (*PostMessageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MessengerService_ServiceDesc.ServiceName,
		Method:  "PostMessage",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToThread(ctx, req.Message.ThreadId, userInfo, messenger.AccessLevel_ACCESS_LEVEL_MESSAGE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmessenger.ErrFailedQuery
		}
	}

	tMessages := table.FivenetMsgsMessages
	stmt := tMessages.
		INSERT(
			tMessages.ThreadID,
			tMessages.Message,
			tMessages.Data,
			tMessages.CreatorID,
		).
		VALUES(
			req.Message.ThreadId,
			req.Message.Message,
			req.Message.Data,
			userInfo.UserId,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}
	req.Message.Id = uint64(lastId)

	message, err := s.getMessage(ctx, req.Message.Id, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	thread, err := s.getThread(ctx, message.ThreadId, userInfo, true)
	if err != nil {
		return nil, errorsmessenger.ErrFailedQuery
	}

	if thread != nil && thread.Access != nil && len(thread.Access.Users) > 0 {
		userIds := []int32{userInfo.UserId}
		if thread.CreatorId != nil {
			userIds = append(userIds, *thread.CreatorId)
		}
		for _, ua := range thread.Access.Users {
			userIds = append(userIds, ua.UserId)
		}

		s.sendUpdate(ctx, &messenger.MessengerEvent{
			Data: &messenger.MessengerEvent_MessageUpdate{
				MessageUpdate: message,
			},
		}, userIds)

		if err := s.setUnreadState(ctx, message.ThreadId, userIds); err != nil {
			return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
		}
	}

	return &PostMessageResponse{
		Message: message,
	}, nil
}

func (s *Server) DeleteMessage(ctx context.Context, req *DeleteMessageRequest) (*DeleteMessageResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.messenger.message.id", int64(req.MessageId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MessengerService_ServiceDesc.ServiceName,
		Method:  "DeleteMessage",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToThread(ctx, req.ThreadId, userInfo, messenger.AccessLevel_ACCESS_LEVEL_ADMIN)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmessenger.ErrFailedQuery
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
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}

	thread, err := s.getThread(ctx, req.ThreadId, userInfo, true)
	if err != nil {
		return nil, errorsmessenger.ErrFailedQuery
	}

	if thread != nil && thread.Access != nil && len(thread.Access.Users) > 0 {
		userIds := []int32{userInfo.UserId}
		for _, ua := range thread.Access.Users {
			userIds = append(userIds, ua.UserId)
		}

		s.sendUpdate(ctx, &messenger.MessengerEvent{
			Data: &messenger.MessengerEvent_MessageDelete{
				MessageDelete: req.MessageId,
			},
		}, userIds)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteMessageResponse{}, nil
}
