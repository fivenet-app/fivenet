package messenger

import (
	"context"
	"errors"
	"slices"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
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
						tThreadsUserAccess.Access.GT_EQ(jet.Int32(int32(messenger.AccessLevel_ACCESS_LEVEL_BLOCKED))),
					),
					jet.AND(
						tThreadsUserAccess.Access.IS_NULL(),
						tThreadsJobAccess.Access.IS_NOT_NULL(),
						tThreadsJobAccess.Access.GT_EQ(jet.Int32(int32(messenger.AccessLevel_ACCESS_LEVEL_BLOCKED))),
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
				).
				LEFT_JOIN(tThreadsJobAccess,
					tThreadsJobAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tThreadsJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(jet.AND(
			wheres...,
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
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
			tThreads.Archived,
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
		).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsUserAccess,
					tThreadsUserAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				).
				LEFT_JOIN(tThreadsJobAccess,
					tThreadsJobAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tThreadsJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
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
			return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
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

func (s *Server) getThread(ctx context.Context, threadId uint64, userInfo *userinfo.UserInfo, withAccess bool) (*messenger.Thread, error) {
	stmt := tThreads.
		SELECT(
			tThreads.ID,
			tThreads.CreatedAt,
			tThreads.UpdatedAt,
			tThreads.DeletedAt,
			tThreads.Title,
			tThreads.Archived,
			tThreads.CreatorID,
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
			tThreads.
				LEFT_JOIN(tThreadsUserAccess,
					tThreadsUserAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				).
				LEFT_JOIN(tThreadsJobAccess,
					tThreadsJobAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tThreadsJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				).
				LEFT_JOIN(tCreator,
					tThreads.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tThreads.CreatorID),
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
					tThreadsUserAccess.Access.GT_EQ(jet.Int32(int32(messenger.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
				jet.AND(
					tThreadsUserAccess.Access.IS_NULL(),
					tThreadsJobAccess.Access.IS_NOT_NULL(),
					tThreadsJobAccess.Access.GT_EQ(jet.Int32(int32(messenger.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
			),
		)).
		LIMIT(1)

	var thread messenger.Thread
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
		access, err := s.getThreadAccess(ctx, threadId)
		if err != nil {
			return nil, err
		}
		thread.Access = access
	}

	return &thread, nil
}

func (s *Server) GetThread(ctx context.Context, req *GetThreadRequest) (*GetThreadResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.checkIfUserHasAccessToThread(ctx, req.ThreadId, userInfo, messenger.AccessLevel_ACCESS_LEVEL_ADMIN)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsmessenger.ErrFailedQuery
		}
	}

	thread, err := s.getThread(ctx, req.ThreadId, userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}

	return &GetThreadResponse{
		Thread: thread,
	}, nil
}

func (s *Server) CreateOrUpdateThread(ctx context.Context, req *CreateOrUpdateThreadRequest) (*CreateOrUpdateThreadResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.messenger.thread.id", int64(req.Thread.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MessengerService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateThread",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if req.Thread.Id == 0 {
		tThreads := table.FivenetMsgsThreads
		stmt := tThreads.
			INSERT(
				tThreads.Title,
				tThreads.Archived,
				tThreads.CreatorID,
				tThreads.CreatorJob,
			).
			VALUES(
				req.Thread.Title,
				req.Thread.Archived,
				userInfo.UserId,
				userInfo.Job,
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errorsmessenger.ErrFailedQuery
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errorsmessenger.ErrFailedQuery
		}

		req.Thread.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		check, err := s.checkIfUserHasAccessToThread(ctx, req.Thread.Id, userInfo, messenger.AccessLevel_ACCESS_LEVEL_ADMIN)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
		}
		if !check && !userInfo.SuperUser {
			if !userInfo.SuperUser {
				return nil, errorsmessenger.ErrFailedQuery
			}
		}

		tThreads := table.FivenetMsgsThreads
		stmt := tThreads.
			UPDATE(
				tThreads.Title,
				tThreads.Archived,
			).
			SET(
				tThreads.Title.SET(jet.String(req.Thread.Title)),
				tThreads.Archived.SET(jet.Bool(req.Thread.Archived)),
			).
			WHERE(tThreads.ID.EQ(jet.Uint64(req.Thread.Id)))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errorsmessenger.ErrFailedQuery
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	accessToDelete, err := s.handleThreadAccessChanges(ctx, tx, messenger.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE, req.Thread.Id, req.Thread.Access)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}

	thread, err := s.getThread(ctx, req.Thread.Id, userInfo, true)
	if err != nil {
		return nil, errorsmessenger.ErrFailedQuery
	}

	if accessToDelete != nil && len(accessToDelete.Users) > 0 {
		userIds := []int32{}
		for _, ua := range accessToDelete.Users {
			userIds = append(userIds, ua.UserId)
		}

		s.sendUpdate(ctx, &messenger.MessengerEvent{
			Data: &messenger.MessengerEvent_ThreadDelete{
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

		s.sendUpdate(ctx, &messenger.MessengerEvent{
			Data: &messenger.MessengerEvent_ThreadUpdate{
				ThreadUpdate: thread,
			},
		}, userIds)
	}

	return &CreateOrUpdateThreadResponse{
		Thread: thread,
	}, nil
}

func (s *Server) DeleteThread(ctx context.Context, req *DeleteThreadRequest) (*DeleteThreadResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.messenger.thread.id", int64(req.ThreadId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MessengerService_ServiceDesc.ServiceName,
		Method:  "DeleteThread",
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

	thread, err := s.getThread(ctx, req.ThreadId, userInfo, true)
	if err != nil {
		return nil, errorsmessenger.ErrFailedQuery
	}

	stmt := tThreads.
		DELETE().
		WHERE(tThreads.ID.EQ(jet.Uint64(req.ThreadId))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
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
			Data: &messenger.MessengerEvent_ThreadDelete{
				ThreadDelete: req.ThreadId,
			},
		}, userIds)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteThreadResponse{}, nil
}

func (s *Server) LeaveThread(ctx context.Context, req *LeaveThreadRequest) (*LeaveThreadResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.messenger.thread.id", int64(req.ThreadId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: MessengerService_ServiceDesc.ServiceName,
		Method:  "DeleteThread",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	resp := &LeaveThreadResponse{}

	check, err := s.checkIfUserHasAccessToThread(ctx, req.ThreadId, userInfo, messenger.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return resp, nil
	}

	thread, err := s.getThread(ctx, req.ThreadId, userInfo, true)
	if err != nil {
		return nil, errorsmessenger.ErrFailedQuery
	}

	if thread == nil {
		return resp, nil
	}

	if thread.Access != nil && len(thread.Access.Users) > 0 {
		idx := slices.IndexFunc(thread.Access.Users, func(ua *messenger.ThreadUserAccess) bool {
			return ua.UserId == userInfo.UserId
		})
		if idx == -1 {
			return resp, nil
		}

		if err := s.deleteThreadAccess(ctx, s.db, thread.Id, &messenger.ThreadAccess{
			Users: []*messenger.ThreadUserAccess{
				thread.Access.Users[idx],
			},
		}); err != nil {
			return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
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

		s.sendUpdate(ctx, &messenger.MessengerEvent{
			Data: &messenger.MessengerEvent_ThreadUpdate{
				ThreadUpdate: thread,
			},
		}, userIds)

		if err := s.setUnreadState(ctx, thread.Id, userIds); err != nil {
			return nil, errswrap.NewError(err, errorsmessenger.ErrFailedQuery)
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return resp, nil
}
