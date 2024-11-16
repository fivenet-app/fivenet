package mailer

import (
	"context"
	"errors"

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

var (
	tThreads      = table.FivenetMailerThreads.AS("thread")
	tThreadsState = table.FivenetMailerThreadsState.AS("thread_state")

	tThreadsRecipients = table.FivenetMailerThreadsRecipients
)

func (s *Server) ListThreads(ctx context.Context, req *ListThreadsRequest) (*ListThreadsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	emailIds, err := s.access.CanUserAccessTargetIDs(ctx, userInfo, mailer.AccessLevel_ACCESS_LEVEL_READ, req.EmailIds...)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if len(emailIds) == 0 {
		return &ListThreadsResponse{
			Pagination: &database.PaginationResponse{},
		}, nil
	}

	ids := []jet.Expression{}
	for _, emailId := range emailIds {
		ids = append(ids, jet.Uint64(emailId))
	}

	wheres := []jet.BoolExpression{jet.Bool(true)}
	if !userInfo.SuperUser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tThreads.DeletedAt.IS_NULL(),
				tThreadsRecipients.EmailID.IN(ids...),
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
				LEFT_JOIN(tThreadsRecipients,
					tThreadsRecipients.ThreadID.EQ(tThreads.ID),
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
			tThreads.CreatorID,
			tThreadsState.ThreadID,
			tThreadsState.EmailID,
			tThreadsState.LastRead,
			tThreadsState.Important,
			tThreadsState.Favorite,
			tThreadsState.Muted,
			tThreadsState.Archived,
		).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsState,
					tThreadsState.ThreadID.EQ(tThreads.ID),
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

func (s *Server) getThread(ctx context.Context, threadId uint64, userInfo *userinfo.UserInfo, withRecipients bool) (*mailer.Thread, error) {
	stmt := tThreads.
		SELECT(
			tThreads.ID,
			tThreads.CreatedAt,
			tThreads.UpdatedAt,
			tThreads.DeletedAt,
			tThreads.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
			tThreadsState.ThreadID,
			tThreadsState.EmailID,
			tThreadsState.LastRead,
			tThreadsState.Important,
			tThreadsState.Favorite,
			tThreadsState.Muted,
			tThreadsState.Archived,
		).
		FROM(
			tThreads.
				LEFT_JOIN(tCreator,
					tThreads.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tThreads.CreatorID),
				).
				LEFT_JOIN(tThreadsState,
					tThreadsState.ThreadID.EQ(tThreads.ID).
						AND(tThreadsState.EmailID.EQ(jet.Int32(userInfo.UserId))),
				),
		).
		WHERE(jet.AND(
			tThreads.ID.EQ(jet.Uint64(threadId)),
		)).
		GROUP_BY(tThreads.ID)

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

	if withRecipients {
		recipients, err := s.getThreadRecipients(ctx, threadId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
		thread.Recipients = recipients
	}

	return &thread, nil
}

func (s *Server) getThreadRecipients(ctx context.Context, threadId uint64) ([]*mailer.ThreadRecipientEmail, error) {
	recipients := []*mailer.ThreadRecipientEmail{}

	emailsStmt := tThreadsRecipients.
		SELECT(
			tThreadsRecipients.ID,
			tThreadsRecipients.ThreadID,
			tThreadsRecipients.EmailID,
			tEmails.Email,
		).
		FROM(
			tThreadsRecipients.
				INNER_JOIN(tEmails,
					tEmails.ID.EQ(tThreadsRecipients.EmailID),
				),
		).
		WHERE(
			tThreadsRecipients.ThreadID.EQ(jet.Uint64(threadId)),
		)

	if err := emailsStmt.QueryContext(ctx, s.db, &recipients); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return recipients, nil
}

func (s *Server) GetThread(ctx context.Context, req *GetThreadRequest) (*GetThreadResponse, error) {
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

	tThreads := table.FivenetMailerThreads
	stmt := tThreads.
		INSERT(
			tThreads.CreatorEmailID,
			tThreads.CreatorID,
		).
		VALUES(
			req.Thread.CreatorEmailId,
			userInfo.UserId,
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

	/*
		accessChanges, err := s.access.HandleAccessChanges(ctx, tx, req.Thread.Id, nil, req.Thread.Recipients)
		if err != nil {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	*/

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	thread, err := s.getThread(ctx, req.Thread.Id, userInfo, true)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	/*
		    TODO
				if accessChanges != nil && len(accessChanges.Users.ToDelete) > 0 {
					emailIds := []int32{}
					for _, ua := range accessChanges.Users.ToDelete {
						emailIds = append(emailIds, ua.UserId)
					}

					s.sendUpdate(ctx, &mailer.MailerEvent{
						Data: &mailer.MailerEvent_ThreadDelete{
							ThreadDelete: thread.Id,
						},
					}, emailIds)
				}
	*/

	if len(thread.Recipients) > 0 {
		emailIds := []uint64{}
		if thread != nil && thread.CreatorId != nil {
			emailIds = append(emailIds, thread.CreatorEmailId)
		}
		for _, ua := range thread.Recipients {
			emailIds = append(emailIds, ua.EmailId)
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_ThreadUpdate{
				ThreadUpdate: thread,
			},
		}, emailIds)
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

	check, err := s.access.CanUserAccessTarget(ctx, req.ThreadId, userInfo, mailer.AccessLevel_ACCESS_LEVEL_VIEW)
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

	if thread != nil && thread.Recipients != nil && len(thread.Recipients) > 0 {
		emailids := []uint64{}
		if thread.CreatorId != nil {
			emailids = append(emailids, thread.CreatorEmailId)
		}
		for _, ua := range thread.Recipients {
			emailids = append(emailids, ua.EmailId)
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_ThreadDelete{
				ThreadDelete: req.ThreadId,
			},
		}, emailids)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteThreadResponse{}, nil
}
