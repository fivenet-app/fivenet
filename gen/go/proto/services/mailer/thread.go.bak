package mailer

import (
	"context"
	"errors"
	"slices"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
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

const (
	ThreadsDefaultPageSize = 16
)

func (s *Server) ListThreads(ctx context.Context, req *ListThreadsRequest) (*ListThreadsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	wheres := []jet.BoolExpression{jet.Bool(true)}
	if !userInfo.SuperUser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tThreads.DeletedAt.IS_NULL(),
				jet.OR(
					jet.AND(
						tThreads.CreatorID.EQ(jet.Int32(userInfo.UserId)),
						tThreads.CreatorJob.EQ(jet.String(userInfo.Job)),
					),
					jet.AND(
						tThreadsUserAccess.Access.IS_NOT_NULL(),
						tThreadsUserAccess.Access.GT_EQ(jet.Int32(int32(mailer.AccessLevel_ACCESS_LEVEL_BLOCKED))),
					),
				),
			),
		}
	}

	if req.After != nil {
		wheres = append(wheres, tThreads.UpdatedAt.GT_EQ(jet.TimestampT(req.After.AsTime())))
	}

	countStmt := tThreads.
		SELECT(
			jet.COUNT(jet.DISTINCT(tThreads.ID)).AS("datacount.totalcount"),
		).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsUserAccess,
					tThreadsUserAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				),
		).
		WHERE(jet.AND(
			wheres...,
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, ThreadsDefaultPageSize)
	resp := &ListThreadsResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tThreads.
		SELECT(
			tThreads.ID,
			tThreads.CreatedAt,
			tThreads.UpdatedAt,
			tThreads.DeletedAt,
			tThreads.Title,
			tThreads.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
			tThreadsUserState.ThreadID,
			tThreadsUserState.UserID,
			tThreadsUserState.LastRead,
			tThreadsUserState.Important,
			tThreadsUserState.Favorite,
			tThreadsUserState.Muted,
			tThreadsUserState.Archived,
		).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsUserAccess,
					tThreadsUserAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				).
				LEFT_JOIN(tCreator,
					tThreads.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tThreadsUserState,
					tThreadsUserState.ThreadID.EQ(tThreads.ID).
						AND(tThreadsUserState.UserID.EQ(jet.Int32(userInfo.UserId))),
				),
		).
		WHERE(jet.AND(
			wheres...,
		)).
		OFFSET(req.Pagination.Offset).
		GROUP_BY(tThreads.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Threads); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Threads); i++ {
		if resp.Threads[i].Creator != nil {
			jobInfoFn(resp.Threads[i].Creator)
		}
	}

	resp.Pagination.Update(len(resp.Threads))

	return resp, nil
}

func (s *Server) getThread(ctx context.Context, threadId uint64, userInfo *userinfo.UserInfo, withAccess bool) (*mailer.Thread, error) {
	stmt := tThreads.
		SELECT(
			tThreads.ID,
			tThreads.CreatedAt,
			tThreads.UpdatedAt,
			tThreads.DeletedAt,
			tThreads.Title,
			tThreads.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
			tThreadsUserState.ThreadID,
			tThreadsUserState.UserID,
			tThreadsUserState.LastRead,
			tThreadsUserState.Important,
			tThreadsUserState.Favorite,
			tThreadsUserState.Muted,
			tThreadsUserState.Archived,
		).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsUserAccess,
					tThreadsUserAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				).
				LEFT_JOIN(tCreator,
					tThreads.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tThreads.CreatorID),
				).
				LEFT_JOIN(tThreadsUserState,
					tThreadsUserState.ThreadID.EQ(tThreads.ID).
						AND(tThreadsUserState.UserID.EQ(jet.Int32(userInfo.UserId))),
				),
		).
		WHERE(jet.AND(
			tThreads.ID.EQ(jet.Uint64(threadId)),
			jet.OR(
				jet.AND(
					tThreads.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					tThreads.CreatorJob.EQ(jet.String(userInfo.Job)),
				),
				jet.AND(
					tThreadsUserAccess.Access.IS_NOT_NULL(),
					tThreadsUserAccess.Access.GT_EQ(jet.Int32(int32(mailer.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
			),
		)).
		LIMIT(1)

	var thread mailer.Thread
	if err := stmt.QueryContext(ctx, s.db, &thread); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if thread.Id == 0 {
		return nil, nil
	}

	if thread.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, thread.Creator)
	}

	if withAccess {
		userAccess, err := s.access.Users.List(ctx, s.db, threadId)
		if err != nil {
			return nil, err
		}
		thread.Access = &mailer.ThreadAccess{
			Users: userAccess,
		}
	}

	return &thread, nil
}

func (s *Server) GetThread(ctx context.Context, req *GetThreadRequest) (*GetThreadResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.ThreadId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_PARTICIPANT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmailer.ErrFailedQuery
		}
	}

	thread, err := s.getThread(ctx, req.ThreadId, userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	return &GetThreadResponse{
		Thread: thread,
	}, nil
}

func (s *Server) CreateThread(ctx context.Context, req *CreateThreadRequest) (*CreateThreadResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.mailer.thread.id", int64(req.Thread.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
		Method:  "CreateThread",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tThreads := table.FivenetMsgsThreads
	stmt := tThreads.
		INSERT(
			tThreads.Title,
			tThreads.CreatorID,
			tThreads.CreatorJob,
		).
		VALUES(
			req.Thread.Title,
			userInfo.UserId,
			userInfo.Job,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	req.Thread.Id = uint64(lastId)

	req.Message.ThreadId = req.Thread.Id
	req.Message.CreatorId = &userInfo.UserId
	if _, err := s.createMessage(ctx, tx, req.Message); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	accessChanges, err := s.access.HandleAccessChanges(ctx, tx, req.Thread.Id, nil, req.Thread.Access.Users)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	thread, err := s.getThread(ctx, req.Thread.Id, userInfo, true)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	if accessChanges != nil && len(accessChanges.Users.ToDelete) > 0 {
		userIds := []int32{}
		for _, ua := range accessChanges.Users.ToDelete {
			userIds = append(userIds, ua.UserId)
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_ThreadDelete{
				ThreadDelete: thread.Id,
			},
		}, userIds)
	}

	if len(thread.Access.Users) > 0 {
		userIds := []int32{userInfo.UserId}
		if thread != nil && thread.CreatorId != nil {
			userIds = append(userIds, *thread.CreatorId)
		}
		for _, ua := range thread.Access.Users {
			userIds = append(userIds, ua.UserId)
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_ThreadUpdate{
				ThreadUpdate: thread,
			},
		}, userIds)
	}

	return &CreateThreadResponse{
		Thread: thread,
	}, nil
}

func (s *Server) DeleteThread(ctx context.Context, req *DeleteThreadRequest) (*DeleteThreadResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.mailer.thread.id", int64(req.ThreadId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
		Method:  "DeleteThread",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.ThreadId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_PARTICIPANT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmailer.ErrFailedQuery
		}
	}

	thread, err := s.getThread(ctx, req.ThreadId, userInfo, true)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	stmt := tThreads.
		DELETE().
		WHERE(tThreads.ID.EQ(jet.Uint64(req.ThreadId))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if thread != nil && thread.Access != nil && len(thread.Access.Users) > 0 {
		userIds := []int32{userInfo.UserId}
		if thread.CreatorId != nil {
			userIds = append(userIds, *thread.CreatorId)
		}
		for _, ua := range thread.Access.Users {
			userIds = append(userIds, ua.UserId)
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_ThreadDelete{
				ThreadDelete: req.ThreadId,
			},
		}, userIds)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteThreadResponse{}, nil
}

func (s *Server) LeaveThread(ctx context.Context, req *LeaveThreadRequest) (*LeaveThreadResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.mailer.thread.id", int64(req.ThreadId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MailerService_ServiceDesc.ServiceName,
		Method:  "DeleteThread",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	resp := &LeaveThreadResponse{}

	check, err := s.access.CanUserAccessTarget(ctx, req.ThreadId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_PARTICIPANT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return resp, nil
	}

	thread, err := s.getThread(ctx, req.ThreadId, userInfo, true)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	if thread == nil {
		return resp, nil
	}

	if thread.Access != nil && len(thread.Access.Users) > 0 {
		idx := slices.IndexFunc(thread.Access.Users, func(ua *mailer.ThreadUserAccess) bool {
			return ua.UserId == userInfo.UserId
		})
		if idx == -1 {
			return resp, nil
		}

		if err := s.access.Users.DeleteEntryWithCondition(ctx, s.db, table.FivenetMsgsThreadsUserAccess.UserID.EQ(jet.Int32(userInfo.UserId)), thread.Id); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		thread.Access.Users = slices.Delete(thread.Access.Users, idx, 1)
	}

	if thread.Access != nil && len(thread.Access.Users) > 0 {
		userIds := []int32{userInfo.UserId}
		if thread.CreatorId != nil {
			userIds = append(userIds, *thread.CreatorId)
		}
		for _, ua := range thread.Access.Users {
			userIds = append(userIds, ua.UserId)
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_ThreadUpdate{
				ThreadUpdate: thread,
			},
		}, userIds)

		if err := s.setUnreadState(ctx, thread.Id, userIds); err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return resp, nil
}
