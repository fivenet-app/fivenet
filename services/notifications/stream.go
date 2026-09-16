package notifications

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"sync"
	"time"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	mailerevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/events"
	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	notificationsclientview "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/clientview"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/notifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/grpcws"
	natsutils "github.com/fivenet-app/fivenet/v2026/pkg/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const feedFetch = 8

var errAccountGroupsChangesResync = errors.New("account group changes subscriber requires resync")

func (s *Server) buildSubjects(
	ctx context.Context,
	userInfo *pbuserinfo.UserInfo,
) ([]string, []string, error) {
	baseSubjects := []string{
		fmt.Sprintf("%s.%s", notifi.BaseSubject, notifi.SystemTopic),
	}

	if userInfo.GetUserId() == 0 {
		return baseSubjects, nil, nil
	}

	if userInfo.GetUserId() > 0 {
		baseSubjects = append(
			baseSubjects,
			fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, userInfo.GetUserId()),
		)
	}
	if userInfo.GetJob() != "" {
		baseSubjects = append(baseSubjects,
			fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, userInfo.GetJob()),
			fmt.Sprintf("%s.%s.%s.>", notifi.BaseSubject, notifi.JobGradeTopic, userInfo.GetJob()),
		)
	}

	// Clone user info and disable superuser (so a superuser doesn't receive notifications for "all" emails..)
	clonedUserInfo := userInfo.Clone()
	clonedUserInfo.Superuser = false
	emails, err := s.mailerStore.ListUserEmails(ctx, s.db, clonedUserInfo, nil, false, false)
	if err != nil {
		return baseSubjects, nil, errswrap.NewError(err, ErrFailedStream)
	}

	additionalSubjects := []string{}
	for _, email := range emails {
		additionalSubjects = append(
			additionalSubjects,
			fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.MailerTopic, email.GetId()),
		)
	}

	return baseSubjects, additionalSubjects, nil
}

func (s *Server) buildClientViewSubjects(
	ctx context.Context,
	userInfo *pbuserinfo.UserInfo,
	clientView *notificationsclientview.ClientView,
) ([]string, error) {
	spec, ok := clientView.GetType().Spec()
	if !ok || clientView.GetType() == notificationsclientview.ObjectType_OBJECT_TYPE_UNSPECIFIED {
		return nil, nil
	}

	if clientView.Id == nil {
		return nil, nil
	}

	switch spec.Visibility {
	case notificationsclientview.VisibilityTargetAccess:
		gAccess, ok := access.GetAccess(spec.AccessRegistryKey)
		if !ok {
			return nil, nil
		}

		check, err := gAccess.CanUserAccessTarget(
			ctx,
			clientView.GetId(),
			userInfo,
			2,
		)
		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrFailedStream)
		}

		if !check {
			s.logger.Warn("user does not have access to the object",
				zap.Int32("user_id", userInfo.GetUserId()),
				zap.String("object_type", clientView.GetType().String()),
				zap.Int64("object_id", clientView.GetId()),
			)
			return nil, nil
		}

	case notificationsclientview.VisibilityJobScoped:
	}

	return []string{
		fmt.Sprintf(
			"%s.%s.%s.%d",
			notifi.BaseSubject,
			notifi.ObjectTopic,
			spec.NatsKey,
			clientView.GetId(),
		),
	}, nil
}

func (s *Server) shouldDeliverObjectEvent(
	dest *notificationsclientview.ObjectEvent,
	userInfo *pbuserinfo.UserInfo,
) bool {
	if dest.UserId != nil && dest.GetUserId() == userInfo.GetUserId() {
		return false
	}

	spec, ok := dest.GetType().Spec()
	if !ok {
		return false
	}

	if spec.Visibility == notificationsclientview.VisibilityJobScoped {
		if dest.Job == nil || userInfo.GetJob() != dest.GetJob() {
			return false
		}
	}

	return true
}

func applyUserInfoChanged(currentUserInfo *pbuserinfo.UserInfo, event *pbuserinfo.UserInfoChanged) {
	if currentUserInfo == nil || event == nil {
		return
	}

	currentUserInfo.Job = event.GetNewJob()
	currentUserInfo.JobGrade = event.GetNewJobGrade()
}

func applyAccountGroupsChanged(
	currentUserInfo *pbuserinfo.UserInfo,
	event *pbuserinfo.AccountGroupsChanged,
) {
	if currentUserInfo == nil || event == nil {
		return
	}

	wasSuperuser := currentUserInfo.GetSuperuser()
	currentUserInfo.CanBeSuperuser = event.GetCanBeSuperuser()
	currentUserInfo.CanBeConfigAdmin = event.GetCanBeConfigAdmin()
	if event.GetNewGroups() == nil {
		currentUserInfo.Groups = nil
	} else {
		currentUserInfo.Groups = &accounts.AccountGroups{
			Groups: slices.Clone(event.GetNewGroups().GetGroups()),
		}
	}

	if !currentUserInfo.GetCanBeSuperuser() {
		currentUserInfo.Superuser = false
	}

	if wasSuperuser && !currentUserInfo.GetCanBeSuperuser() &&
		currentUserInfo.GetOriginalJob() != nil {
		currentUserInfo.Job = currentUserInfo.GetOriginalJob().GetJob()
		currentUserInfo.JobGrade = currentUserInfo.GetOriginalJob().GetJobGrade()
	}
}

func (s *Server) Stream(srv pbnotifications.NotificationsService_StreamServer) error {
	ctx := srv.Context()
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Track changes to user info, so we can send an updated user info to the user
	currentUserInfo := userInfo.Clone()
	accountOnly := currentUserInfo.GetUserId() == 0
	accountGroupsChanges := s.userinfoChanges.SubscribeAccountGroupsChanges()
	defer s.userinfoChanges.UnsubscribeAccountGroupsChanges(accountGroupsChanges)

	subjectsMu := &sync.Mutex{}
	stateMu := &sync.RWMutex{}
	baseSubjects, additionalSubjects, err := s.buildSubjects(ctx, currentUserInfo)
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}
	clientViewSubjects := []string{}
	var currentClientView *notificationsclientview.ClientView

	meta := metadata.ExtractIncoming(ctx)
	connId := meta.Get(grpcws.ConnectionIdHeader)

	// Create durable pull consumer with multi-filter, required to update filter subjects dynamically
	consCfg := jetstream.ConsumerConfig{
		Durable: natsutils.GenerateConsumerName(
			userInfo.GetAccountId(),
			userInfo.GetUserId(),
			connId,
		),
		FilterSubjects:    append(baseSubjects, additionalSubjects...),
		DeliverPolicy:     jetstream.DeliverNewPolicy,
		AckPolicy:         jetstream.AckNonePolicy,
		MaxWaiting:        8,
		InactiveThreshold: 15 * time.Second,
	}
	consumer, err := s.js.CreateOrUpdateConsumer(ctx, notifi.StreamName, consCfg)
	if err != nil {
		return fmt.Errorf("failed to create consumer (%v). %w", consCfg.FilterSubjects, err)
	}

	// Count only after the consumer boundary exists. Durable notification events
	// carry an updated unread count, so this is the stream's only count lookup.
	var notificationCount int64
	var notificationCountMu sync.RWMutex
	if !accountOnly {
		notificationCount, err = s.store.CountUnread(ctx, userInfo.GetUserId())
		if err != nil {
			return errswrap.NewError(err, ErrFailedStream)
		}
	}
	notificationCountSnapshot := func() int64 {
		notificationCountMu.RLock()
		defer notificationCountMu.RUnlock()
		return notificationCount
	}
	setNotificationCount := func(count int64) {
		notificationCountMu.Lock()
		defer notificationCountMu.Unlock()
		notificationCount = count
	}
	decrementNotificationCount := func(count int64) {
		notificationCountMu.Lock()
		defer notificationCountMu.Unlock()
		if notificationCount-count <= 0 {
			notificationCount = 0
			return
		}
		notificationCount -= count
	}
	currentUserInfoSnapshot := func() *pbuserinfo.UserInfo {
		stateMu.RLock()
		defer stateMu.RUnlock()
		return currentUserInfo.Clone()
	}
	// Keep the durable consumer alive across websocket stream restarts on the
	// same connection. The consumer will expire via InactiveThreshold once the
	// connection is truly idle, which avoids an old handler deleting the fresh
	// consumer created by a restarted stream.

	// Central pipe: all feeds push messages into outCh
	outCh := make(chan *pbnotifications.StreamResponse, 256)
	defer close(outCh)
	g, gctx := errgroup.WithContext(ctx)

	// A consumer uses DeliverNew, so no notification event is guaranteed to be
	// delivered on connection. Send the snapshot explicitly before live events.
	outCh <- &pbnotifications.StreamResponse{
		NotificationCount: notificationCountSnapshot(),
		Data: &pbnotifications.StreamResponse_NotificationState{
			NotificationState: true,
		},
	}

	refreshConsumerSubjects := func() error {
		info, err := consumer.Info(gctx)
		if err != nil {
			return errswrap.NewError(err, ErrFailedStream)
		}

		cfg := info.Config
		subjectsMu.Lock()
		filterSubjects := make(
			[]string,
			0,
			len(baseSubjects)+len(additionalSubjects)+len(clientViewSubjects),
		)
		filterSubjects = append(filterSubjects, baseSubjects...)
		filterSubjects = append(filterSubjects, additionalSubjects...)
		filterSubjects = append(filterSubjects, clientViewSubjects...)
		subjectsMu.Unlock()

		cfg.FilterSubjects = filterSubjects

		if _, err := s.js.UpdateConsumer(gctx, notifi.StreamName, cfg); err != nil {
			return errswrap.NewError(
				fmt.Errorf("failed to update consumer. %w", err),
				ErrFailedStream,
			)
		}

		return nil
	}

	rebuildAndRefreshSubjects := func() error {
		newBaseSubjects, newAdditionalSubjects, err := s.buildSubjects(gctx, currentUserInfo)
		if err != nil {
			return errswrap.NewError(err, ErrFailedStream)
		}

		subjectsMu.Lock()
		baseSubjects = newBaseSubjects
		additionalSubjects = newAdditionalSubjects
		clientView := currentClientView
		subjectsMu.Unlock()

		if clientView != nil {
			newClientViewSubjects, err := s.buildClientViewSubjects(
				gctx,
				currentUserInfo,
				clientView,
			)
			if err != nil {
				return err
			}
			subjectsMu.Lock()
			clientViewSubjects = newClientViewSubjects
			subjectsMu.Unlock()
		}

		return refreshConsumerSubjects()
	}

	applyAndRefreshAccountGroups := func(change *pbuserinfo.AccountGroupsChanged) error {
		stateMu.Lock()
		defer stateMu.Unlock()

		applyAccountGroupsChanged(currentUserInfo, change)
		return rebuildAndRefreshSubjects()
	}

	applyAndRefreshUserInfo := func(change *pbuserinfo.UserInfoChanged) error {
		stateMu.Lock()
		defer stateMu.Unlock()

		applyUserInfoChanged(currentUserInfo, change)
		return rebuildAndRefreshSubjects()
	}

	g.Go(func() error {
		// Update metrics for active user sessions in first goroutine
		s.metrics.lastSession.SetToCurrentTime()
		s.metrics.activeSessions.Inc()
		defer s.metrics.activeSessions.Dec()

		for {
			msg, err := srv.Recv()
			if errors.Is(err, io.EOF) || protoutils.IsContextCanceled(err) {
				return nil
			}
			if err != nil {
				return err
			}
			if msg == nil {
				continue // Skip nil messages
			}

			switch d := msg.GetData().(type) {
			case *pbnotifications.StreamRequest_Clientview:
				if accountOnly {
					continue
				}

				clientView := d.Clientview
				if clientView == nil {
					continue // Skip nil client view
				}

				stateMu.Lock()
				newClientViewSubjects, err := s.buildClientViewSubjects(
					gctx,
					currentUserInfo,
					clientView,
				)
				if err != nil {
					stateMu.Unlock()
					return err
				}

				subjectsMu.Lock()
				currentClientView = clientView
				clientViewSubjects = newClientViewSubjects
				subjectsMu.Unlock()

				if err := refreshConsumerSubjects(); err != nil {
					stateMu.Unlock()
					return err
				}
				stateMu.Unlock()
			}
		}
	})

	// Writer goroutine - single gRPC send loop
	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return nil

			case msg, ok := <-outCh:
				if !ok {
					return nil
				}
				if msg == nil {
					continue
				}

				if err := srv.Send(msg); err != nil {
					if protoutils.IsContextCanceled(err) {
						return nil
					}
					return err
				}
			}
		}
	})

	// Canonical account changes are delivered through pkg/userinfo and adapted
	// to the existing browser UserEvent contract here.
	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return nil

			case change, ok := <-accountGroupsChanges:
				if !ok {
					return errAccountGroupsChangesResync
				}
				if change == nil || change.GetAccountId() != userInfo.GetAccountId() {
					continue
				}

				if err := applyAndRefreshAccountGroups(change); err != nil {
					return err
				}

				dest := &notificationsevents.UserEvent{
					Data: &notificationsevents.UserEvent_AccountGroupsChanged{
						AccountGroupsChanged: change,
					},
				}
				select {
				case <-gctx.Done():
					return nil
				case outCh <- &pbnotifications.StreamResponse{
					NotificationCount: notificationCountSnapshot(),
					Data: &pbnotifications.StreamResponse_UserEvent{
						UserEvent: dest,
					},
				}:
				}
			}
		}
	})

	g.Go(func() error {
		msgs, err := consumer.Messages(
			jetstream.PullMaxMessages(feedFetch),
			jetstream.WithMessagesErrOnMissingHeartbeat(false),
		)
		if err != nil {
			return err
		}
		defer msgs.Stop()

		for {
			msg, err := msgs.Next(jetstream.NextContext(gctx))
			if err != nil {
				if protoutils.IsContextCanceled(err) ||
					errors.Is(err, jetstream.ErrMsgIteratorClosed) {
					return nil
				}
				return err
			}

			// Publish notifications sent directly to user via the message queue
			if msg == nil {
				s.logger.Warn(
					"nil notification message received via message queue",
					zap.Int32("user_id", userInfo.GetUserId()),
				)
				continue
			}

			topic, parts := notifi.SplitSubject(msg.Subject())
			switch topic {
			case notifi.UserTopic:
				if accountOnly {
					continue
				}

				var dest notificationsevents.UserEvent
				if err := protoutils.UnmarshalPartialJSON(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				switch d := dest.GetData().(type) {
				case *notificationsevents.UserEvent_Notification:
					if err := s.hydrateNotifications(
						gctx,
						currentUserInfoSnapshot(),
						[]*resourcesnotifications.Notification{d.Notification},
					); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}
					if topic == notifi.UserTopic && d.Notification.GetId() > 0 &&
						dest.UnreadCount != nil {
						// This is intentionally eventually consistent: concurrent
						// publishers can observe and deliver counts out of order.
						setNotificationCount(*dest.UnreadCount)
					}

				case *notificationsevents.UserEvent_NotificationsReadCount:
					if topic == notifi.UserTopic {
						decrementNotificationCount(d.NotificationsReadCount)
					}

				case *notificationsevents.UserEvent_NotificationsUnreadCount:
					if topic == notifi.UserTopic {
						setNotificationCount(d.NotificationsUnreadCount)
					}

				case *notificationsevents.UserEvent_UserInfoChanged:
					if err := applyAndRefreshUserInfo(d.UserInfoChanged); err != nil {
						return err
					}
				}

				select {
				case <-gctx.Done():
					return nil

				case outCh <- &pbnotifications.StreamResponse{
					NotificationCount: notificationCountSnapshot(),
					Data: &pbnotifications.StreamResponse_UserEvent{
						UserEvent: &dest,
					},
				}:
				}

			case notifi.JobTopic:
				if accountOnly {
					continue
				}

				var dest notificationsevents.JobEvent
				if err := protoutils.UnmarshalPartialJSON(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				select {
				case <-gctx.Done():
					return nil

				case outCh <- &pbnotifications.StreamResponse{
					NotificationCount: notificationCountSnapshot(),
					Data: &pbnotifications.StreamResponse_JobEvent{
						JobEvent: &dest,
					},
				}:
				}

			case notifi.JobGradeTopic:
				if accountOnly {
					continue
				}

				// Make sure the job grade is included
				if len(parts) < 2 {
					continue
				}
				grade, err := strconv.ParseInt(parts[1], 10, 32)
				if err != nil {
					continue
				}
				if currentUserInfoSnapshot().GetJobGrade() < int32(grade) {
					continue
				}
				var dest notificationsevents.JobGradeEvent
				if err := protoutils.UnmarshalPartialJSON(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				select {
				case <-gctx.Done():
					return nil
				case outCh <- &pbnotifications.StreamResponse{
					NotificationCount: notificationCountSnapshot(),
					Data: &pbnotifications.StreamResponse_JobGradeEvent{
						JobGradeEvent: &dest,
					},
				}:
				}

			case notifi.SystemTopic:
				var dest notificationsevents.SystemEvent
				if err := protoutils.UnmarshalPartialJSON(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				select {
				case <-gctx.Done():
					return nil

				case outCh <- &pbnotifications.StreamResponse{
					NotificationCount: notificationCountSnapshot(),
					Data: &pbnotifications.StreamResponse_SystemEvent{
						SystemEvent: &dest,
					},
				}:
				}

			case notifi.ObjectTopic:
				if accountOnly {
					continue
				}

				var dest notificationsclientview.ObjectEvent
				if err := protoutils.UnmarshalPartialJSON(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				if !s.shouldDeliverObjectEvent(&dest, currentUserInfoSnapshot()) {
					continue
				}

				select {
				case <-gctx.Done():
					return nil

				case outCh <- &pbnotifications.StreamResponse{
					NotificationCount: notificationCountSnapshot(),
					Data: &pbnotifications.StreamResponse_ObjectEvent{
						ObjectEvent: &dest,
					},
				}:
				}

			case notifi.MailerTopic:
				if accountOnly {
					continue
				}

				var dest mailerevents.MailerEvent
				if err := protoutils.UnmarshalPartialJSON(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				select {
				case <-gctx.Done():
					return nil

				case outCh <- &pbnotifications.StreamResponse{
					NotificationCount: notificationCountSnapshot(),
					Data: &pbnotifications.StreamResponse_MailerEvent{
						MailerEvent: &dest,
					},
				}:
				}
			}
		}
	})

	return g.Wait()
}
