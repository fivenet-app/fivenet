package mailer

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	ThreadsDefaultPageSize = 16
)

var (
	tThreads      = table.FivenetMailerThreads.AS("thread")
	tThreadsState = table.FivenetMailerThreadsState.AS("thread_state")

	tThreadsRecipients = table.FivenetMailerThreadsRecipients
)

func (s *Server) ListThreads(
	ctx context.Context,
	req *pbmailer.ListThreadsRequest,
) (*pbmailer.ListThreadsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	emailIds, err := s.access.CanUserAccessTargetIDs(
		ctx,
		userInfo,
		mailer.AccessLevel_ACCESS_LEVEL_READ,
		req.GetEmailIds()...)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if len(emailIds) == 0 {
		return &pbmailer.ListThreadsResponse{
			Pagination: &database.PaginationResponse{},
		}, nil
	}

	ids := []jet.Expression{}
	for _, emailId := range emailIds {
		ids = append(ids, jet.Int64(emailId))
	}

	wheres := []jet.BoolExpression{jet.Bool(true)}
	if !userInfo.GetSuperuser() {
		wheres = []jet.BoolExpression{
			jet.AND(
				tThreads.DeletedAt.IS_NULL(),
				tThreadsRecipients.EmailID.IN(ids...),
			),
		}
	}

	if req.Unread != nil {
		wheres = append(
			wheres,
			tThreadsState.Unread.IS_NOT_NULL().
				AND(tThreadsState.Unread.EQ(jet.Bool(req.GetUnread()))),
		)
	}
	if req.Archived != nil {
		wheres = append(
			wheres,
			tThreadsState.Archived.IS_NOT_NULL().
				AND(tThreadsState.Archived.EQ(jet.Bool(req.GetArchived()))),
		)
	} else {
		// Skip archived emails by default
		wheres = append(wheres, tThreadsState.Archived.IS_NULL().OR(tThreadsState.Archived.EQ(jet.Bool(false))))
	}

	countStmt := tThreads.
		SELECT(
			jet.COUNT(jet.DISTINCT(tThreads.ID)).AS("data_count.total"),
		).
		FROM(
			tThreads.
				INNER_JOIN(tThreadsRecipients,
					tThreadsRecipients.ThreadID.EQ(tThreads.ID),
				).
				LEFT_JOIN(tThreadsState,
					tThreadsState.ThreadID.EQ(tThreads.ID).
						AND(tThreadsState.EmailID.EQ(jet.Int64(req.GetEmailIds()[0]))),
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

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, ThreadsDefaultPageSize)
	resp := &pbmailer.ListThreadsResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	stmt := tThreads.
		SELECT(
			tThreads.ID,
			tThreads.CreatedAt,
			tThreads.UpdatedAt,
			tThreads.DeletedAt,
			tThreads.Title,
			tThreads.CreatorEmailID,
			tThreads.CreatorEmail,
			tThreads.CreatorEmailID.AS("email.id"),
			tThreads.CreatorEmail.AS("email.email"),
			tThreads.CreatorID,
			tThreadsState.ThreadID,
			tThreadsState.EmailID,
			tThreadsState.Unread,
			tThreadsState.LastRead,
			tThreadsState.Important,
			tThreadsState.Favorite,
			tThreadsState.Muted,
			tThreadsState.Archived,
		).
		FROM(
			tThreads.
				INNER_JOIN(tThreadsRecipients,
					tThreadsRecipients.ThreadID.EQ(tThreads.ID),
				).
				LEFT_JOIN(tThreadsState,
					tThreadsState.ThreadID.EQ(tThreads.ID).
						AND(tThreadsState.EmailID.EQ(jet.Int64(req.GetEmailIds()[0]))),
				),
		).
		WHERE(jet.AND(
			wheres...,
		)).
		OFFSET(req.GetPagination().GetOffset()).
		GROUP_BY(tThreads.ID).
		ORDER_BY(jet.COALESCE(tThreads.UpdatedAt, tThreads.CreatedAt).DESC(), tThreads.ID.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Threads); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetThreads() {
		if resp.GetThreads()[i].GetCreator() != nil {
			jobInfoFn(resp.GetThreads()[i].GetCreator())
		}
	}

	resp.GetPagination().Update(len(resp.GetThreads()))

	return resp, nil
}

func (s *Server) getThread(
	ctx context.Context,
	threadId int64,
	emailId int64,
	userInfo *userinfo.UserInfo,
) (*mailer.Thread, error) {
	stmt := tThreads.
		SELECT(
			tThreads.ID,
			tThreads.CreatedAt,
			tThreads.UpdatedAt,
			tThreads.DeletedAt,
			tThreads.Title,
			tThreads.CreatorEmailID,
			tThreads.CreatorID,
			tThreads.CreatorEmailID.AS("email.id"),
			tThreads.CreatorEmail.AS("email.email"),
			tThreadsState.ThreadID,
			tThreadsState.EmailID,
			tThreadsState.Unread,
			tThreadsState.LastRead,
			tThreadsState.Important,
			tThreadsState.Favorite,
			tThreadsState.Muted,
			tThreadsState.Archived,
		).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsState,
					tThreadsState.ThreadID.EQ(tThreads.ID).
						AND(tThreadsState.EmailID.EQ(jet.Int64(emailId))),
				),
		).
		WHERE(jet.AND(
			tThreads.ID.EQ(jet.Int64(threadId)),
		)).
		GROUP_BY(tThreads.ID).
		LIMIT(1)

	var thread mailer.Thread
	if err := stmt.QueryContext(ctx, s.db, &thread); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}
	}

	if thread.GetId() == 0 {
		return nil, nil
	}

	if thread.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, thread.GetCreator())
	}

	recipients, err := s.getThreadRecipients(ctx, s.db, threadId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	thread.Recipients = recipients

	return &thread, nil
}

func (s *Server) GetThread(
	ctx context.Context,
	req *pbmailer.GetThreadRequest,
) (*pbmailer.GetThreadResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.checkIfEmailPartOfThread(ctx, userInfo, req.GetThreadId(), req.GetEmailId(), mailer.AccessLevel_ACCESS_LEVEL_READ); err != nil {
		return nil, err
	}

	thread, err := s.getThread(ctx, req.GetThreadId(), req.GetEmailId(), userInfo)
	if err != nil {
		return nil, err
	}

	return &pbmailer.GetThreadResponse{
		Thread: thread,
	}, nil
}

func (s *Server) CreateThread(
	ctx context.Context,
	req *pbmailer.CreateThreadRequest,
) (*pbmailer.CreateThreadResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.mailer.thread_id", req.GetThread().GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "CreateThread",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetThread().GetCreatorEmailId(),
		userInfo,
		mailer.AccessLevel_ACCESS_LEVEL_WRITE,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	senderEmail, err := s.getEmail(ctx, req.GetThread().GetCreatorEmailId(), false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Prevent disabled emails from creating threads
	if !userInfo.GetSuperuser() && senderEmail.GetDeactivated() {
		return nil, errorsmailer.ErrEmailDisabled
	}

	emails, err := s.retrieveRecipientsToEmails(ctx, senderEmail, req.GetRecipients())
	if err != nil {
		return nil, err
	}

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
			tThreads.Title,
			tThreads.CreatorEmailID,
			tThreads.CreatorID,
			tThreads.CreatorEmail,
		).
		VALUES(
			req.GetThread().GetTitle(),
			req.GetThread().GetCreatorEmailId(),
			userInfo.GetUserId(),
			senderEmail.GetEmail(),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errorsmailer.ErrFailedQuery
	}

	req.Thread.Id = lastId

	req.Message.ThreadId = req.GetThread().GetId()
	req.Message.Sender = senderEmail
	req.Message.CreatorId = &userInfo.UserId
	req.Message.CreatorJob = &userInfo.Job
	if _, err := s.createMessage(ctx, tx, req.GetMessage()); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Add creator of email to recipients
	emails = append(emails, &mailer.ThreadRecipientEmail{
		EmailId: senderEmail.GetId(),
		Email:   senderEmail,
	})
	if err := s.handleRecipientsChanges(ctx, tx, req.GetThread().GetId(), emails); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	emailIds := []int64{}
	for _, recipient := range emails {
		emailIds = append(emailIds, recipient.GetEmailId())
	}
	if err := s.setUnreadState(ctx, tx, req.GetThread().GetId(), senderEmail.GetId(), emailIds); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	thread, err := s.getThread(
		ctx,
		req.GetThread().GetId(),
		req.GetThread().GetCreatorEmailId(),
		userInfo,
	)
	if err != nil {
		return nil, err
	}

	// Set dummy thread state to make client-side handling easier
	boolTrue := true
	thread.State = &mailer.ThreadState{
		ThreadId: thread.GetId(),
		Unread:   &boolTrue,
	}

	if len(thread.GetRecipients()) > 0 {
		if thread != nil && thread.CreatorId != nil {
			s.sendUpdate(ctx, &mailer.MailerEvent{
				Data: &mailer.MailerEvent_ThreadUpdate{
					ThreadUpdate: thread,
				},
			}, thread.GetCreatorEmailId())
		}

		emailIds := []int64{}
		for _, ua := range thread.GetRecipients() {
			emailIds = append(emailIds, ua.GetEmailId())
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_ThreadUpdate{
				ThreadUpdate: thread,
			},
		}, emailIds...)
	}

	return &pbmailer.CreateThreadResponse{
		Thread: thread,
	}, nil
}

func (s *Server) updateThreadTime(ctx context.Context, tx qrm.DB, threadId int64) error {
	stmt := tThreads.
		UPDATE(
			tThreads.UpdatedAt,
		).
		SET(
			tThreads.UpdatedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tThreads.ID.EQ(jet.Int64(threadId)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Server) DeleteThread(
	ctx context.Context,
	req *pbmailer.DeleteThreadRequest,
) (*pbmailer.DeleteThreadResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.mailer.thread_id", req.GetThreadId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbmailer.MailerService_ServiceDesc.ServiceName,
		Method:  "DeleteThread",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if !userInfo.GetSuperuser() {
		return nil, errorsmailer.ErrFailedQuery
	}

	thread, err := s.getThread(ctx, req.GetThreadId(), req.GetEmailId(), userInfo)
	if err != nil {
		return nil, err
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if thread != nil && thread.GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	tThreads := table.FivenetMailerThreads
	stmt := tThreads.
		UPDATE(
			tThreads.DeletedAt,
		).
		SET(
			tThreads.DeletedAt.SET(deletedAtTime),
		).
		WHERE(
			tThreads.ID.EQ(jet.Int64(req.GetThreadId())),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	if thread != nil && thread.Recipients != nil && len(thread.GetRecipients()) > 0 {
		emailIds := []int64{}
		if thread.CreatorId != nil {
			emailIds = append(emailIds, thread.GetCreatorEmailId())
		}
		for _, ua := range thread.GetRecipients() {
			emailIds = append(emailIds, ua.GetEmailId())
		}

		s.sendUpdate(ctx, &mailer.MailerEvent{
			Data: &mailer.MailerEvent_ThreadDelete{
				ThreadDelete: req.GetThreadId(),
			},
		}, emailIds...)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbmailer.DeleteThreadResponse{}, nil
}
