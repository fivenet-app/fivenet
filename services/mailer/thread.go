package mailer

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	maileraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/access"
	mailerevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/events"
	mailerthreads "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/threads"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbmailer "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsmailer "github.com/fivenet-app/fivenet/v2026/services/mailer/errors"
	mailerstore "github.com/fivenet-app/fivenet/v2026/stores/mailer"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	ThreadsDefaultPageSize = 16
)

var tThreadsRecipients = table.FivenetMailerThreadsRecipients

func (s *Server) ListThreads(
	ctx context.Context,
	req *pbmailer.ListThreadsRequest,
) (*pbmailer.ListThreadsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	emailIds, err := s.access.CanUserAccessTargetIDs(
		ctx,
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
		req.GetEmailIds()...)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if len(emailIds) == 0 {
		return &pbmailer.ListThreadsResponse{
			Pagination: &database.PaginationResponse{},
		}, nil
	}
	query := mailerstore.ThreadListQuery{
		EmailIDs:  emailIds,
		Unread:    req.Unread,
		Archived:  req.Archived,
		Superuser: userInfo.GetJobAdmin(),
	}

	count, err := s.store.CountThreads(ctx, s.db, query)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, ThreadsDefaultPageSize)
	resp := &pbmailer.ListThreadsResponse{
		Pagination: pag,
	}
	if count <= 0 {
		return resp, nil
	}
	query.Offset = req.GetPagination().GetOffset()
	query.Limit = limit

	threads, err := s.store.ListThreads(ctx, s.db, query)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	resp.Threads = threads

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetThreads() {
		if resp.GetThreads()[i].GetCreator() != nil {
			jobInfoFn(resp.GetThreads()[i].GetCreator())
		}
	}

	return resp, nil
}

func (s *Server) getThread(
	ctx context.Context,
	threadId int64,
	emailId int64,
	userInfo *userinfo.UserInfo,
) (*mailerthreads.Thread, error) {
	thread, err := s.store.GetThread(
		ctx,
		s.db,
		threadId,
		emailId,
		userInfo != nil && userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if thread == nil {
		return nil, nil
	}

	if thread.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, thread.GetCreator())
	}

	recipients, err := s.store.ListThreadRecipients(
		ctx,
		s.db,
		threadId,
		userInfo != nil && userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	thread.Recipients = recipients

	return thread, nil
}

func (s *Server) GetThread(
	ctx context.Context,
	req *pbmailer.GetThreadRequest,
) (*pbmailer.GetThreadResponse, error) {
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

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetThread().GetCreatorEmailId(),
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_WRITE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check {
		return nil, errorsmailer.ErrNoPerms
	}

	senderEmail, err := s.getEmail(ctx, userInfo, req.GetThread().GetCreatorEmailId(), false, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Prevent disabled emails from creating threads
	if !userInfo.GetJobAdmin() && senderEmail.GetDeactivated() {
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
	if _, err := s.store.CreateMessage(ctx, tx, req.GetMessage()); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Add creator of email to recipients
	emails = append(emails, &mailerthreads.ThreadRecipientEmail{
		EmailId: senderEmail.GetId(),
		Email:   senderEmail,
	})
	if err := s.store.AddThreadRecipients(ctx, tx, req.GetThread().GetId(), emails); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	emailIds := []int64{}
	for _, recipient := range emails {
		emailIds = append(emailIds, recipient.GetEmailId())
	}
	if err := s.store.SetUnreadState(
		ctx,
		tx,
		req.GetThread().GetId(),
		senderEmail.GetId(),
		emailIds,
	); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

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
	thread.State = &mailerthreads.ThreadState{
		ThreadId: thread.GetId(),
		Unread:   &boolTrue,
	}

	if len(thread.GetRecipients()) > 0 {
		if thread != nil && thread.CreatorId != nil {
			s.sendUpdate(ctx, &mailerevents.MailerEvent{
				Data: &mailerevents.MailerEvent_ThreadUpdate{
					ThreadUpdate: thread,
				},
			}, thread.GetCreatorEmailId())
		}

		emailIds := []int64{}
		for _, ua := range thread.GetRecipients() {
			emailIds = append(emailIds, ua.GetEmailId())
		}

		s.sendUpdate(ctx, &mailerevents.MailerEvent{
			Data: &mailerevents.MailerEvent_ThreadUpdate{
				ThreadUpdate: thread,
			},
		}, emailIds...)
	}

	return &pbmailer.CreateThreadResponse{
		Thread: thread,
	}, nil
}

func (s *Server) DeleteThread(
	ctx context.Context,
	req *pbmailer.DeleteThreadRequest,
) (*pbmailer.DeleteThreadResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.mailer.thread_id", req.GetThreadId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if !userInfo.GetJobAdmin() {
		return nil, errorsmailer.ErrFailedQuery
	}

	thread, err := s.getThread(ctx, req.GetThreadId(), req.GetEmailId(), userInfo)
	if err != nil {
		return nil, err
	}

	var deletedAtTime *timestamp.Timestamp
	if thread == nil || thread.GetDeletedAt() == nil || !userInfo.GetJobAdmin() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteThread(ctx, s.db, req.GetThreadId(), deletedAtTime); err != nil {
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

		s.sendUpdate(ctx, &mailerevents.MailerEvent{
			Data: &mailerevents.MailerEvent_ThreadDelete{
				ThreadDelete: req.GetThreadId(),
			},
		}, emailIds...)
	}

	return &pbmailer.DeleteThreadResponse{}, nil
}
