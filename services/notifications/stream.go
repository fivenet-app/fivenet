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
	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const feedFetch = 8

var (
	// metricActiveUserSessions tracks the number of active user sessions for Prometheus monitoring.
	metricActiveUserSessions = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "user",
		Name:      "active_session_count",
		Help:      "Number of active user sessions.",
	})

	metricLastUserSession = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "user",
		Name:      "last_session_time",
		Help:      "Timestamp of the last started user session.",
	})
)

func (s *Server) buildSubjects(
	ctx context.Context,
	userInfo *pbuserinfo.UserInfo,
) ([]string, []string, error) {
	baseSubjects := []string{
		fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.AccountTopic, userInfo.GetAccountId()),
		fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, userInfo.GetUserId()),
		fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, userInfo.GetJob()),
		fmt.Sprintf("%s.%s.%s.>", notifi.BaseSubject, notifi.JobGradeTopic, userInfo.GetJob()),
		fmt.Sprintf("%s.%s", notifi.BaseSubject, notifi.SystemTopic),
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

	if _, err := s.js.PublishAsyncProto(ctx, userinfo.PollSubject, &pbuserinfo.PollReq{
		AccountId: currentUserInfo.GetAccountId(),
		UserId:    currentUserInfo.GetUserId(),
	}); err != nil {
		s.logger.Error(
			"failed to publish userinfo.poll.request",
			zap.Int32("user_id", currentUserInfo.GetUserId()),
			zap.Error(err),
		)
	}

	subjectsMu := &sync.Mutex{}
	baseSubjects, additionalSubjects, err := s.buildSubjects(ctx, currentUserInfo)
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}
	clientViewSubjects := []string{}
	var currentClientView *notificationsclientview.ClientView

	notificationCount, err := s.store.CountUnread(ctx, userInfo.GetUserId())
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}

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
		return fmt.Errorf("failed to create consumer. %w", err)
	}
	defer s.js.DeleteConsumer(s.ctx, notifi.StreamName, consCfg.Durable)

	// Central pipe: all feeds push messages into outCh
	outCh := make(chan *pbnotifications.StreamResponse, 256)
	defer close(outCh)
	g, gctx := errgroup.WithContext(ctx)

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

	g.Go(func() error {
		// Update metrics for active user sessions in first goroutine
		metricLastUserSession.SetToCurrentTime()
		metricActiveUserSessions.Inc()
		defer metricActiveUserSessions.Dec()

		for {
			msg, err := srv.Recv()
			if errors.Is(err, io.EOF) {
				return err
			}
			if err != nil {
				return err
			}
			if msg == nil {
				continue // Skip nil messages
			}

			switch d := msg.GetData().(type) {
			case *pbnotifications.StreamRequest_Clientview:
				clientView := d.Clientview
				if clientView == nil {
					continue // Skip nil client view
				}

				newClientViewSubjects, err := s.buildClientViewSubjects(
					gctx,
					currentUserInfo,
					clientView,
				)
				if err != nil {
					return err
				}

				subjectsMu.Lock()
				currentClientView = clientView
				clientViewSubjects = newClientViewSubjects
				subjectsMu.Unlock()

				if err := refreshConsumerSubjects(); err != nil {
					return err
				}
			}
		}
	})

	// Writer goroutine - single gRPC send loop
	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return nil

			case msg := <-outCh:
				if msg == nil {
					continue
				}

				if err := srv.Send(msg); err != nil {
					return err
				}
			}
		}
	})

	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return nil

			default:
			}

			batch, err := consumer.Fetch(feedFetch,
				jetstream.FetchMaxWait(3*time.Second))
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) ||
					errors.Is(err, jetstream.ErrNoMessages) {
					continue // keep polling
				}
				return err
			}

			for m := range batch.Messages() {
				// Publish notifications sent directly to user via the message queue
				if m == nil {
					s.logger.Warn(
						"nil notification message received via message queue",
						zap.Int32("user_id", currentUserInfo.GetUserId()),
					)
					continue
				}

				if err := m.Ack(); err != nil {
					s.logger.Error("failed to ack notification message", zap.Error(err))
				}

				topic, parts := notifi.SplitSubject(m.Subject())
				switch topic {
				case notifi.UserTopic, notifi.AccountTopic:
					var dest notificationsevents.UserEvent
					if err := protoutils.UnmarshalPartialJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					needsSubjectRefresh := false
					switch d := dest.GetData().(type) {
					case *notificationsevents.UserEvent_Notification:
						if topic == notifi.UserTopic {
							notificationCount++
						}

					case *notificationsevents.UserEvent_NotificationsReadCount:
						if topic == notifi.UserTopic {
							if notificationCount-d.NotificationsReadCount <= 0 {
								notificationCount = 0
							} else {
								notificationCount -= d.NotificationsReadCount
							}
						}

					case *notificationsevents.UserEvent_UserInfoChanged:
						if topic == notifi.UserTopic {
							applyUserInfoChanged(currentUserInfo, d.UserInfoChanged)
							needsSubjectRefresh = true
						}

					case *notificationsevents.UserEvent_AccountGroupsChanged:
						applyAccountGroupsChanged(currentUserInfo, d.AccountGroupsChanged)
						needsSubjectRefresh = true
					}

					if needsSubjectRefresh {
						if err := rebuildAndRefreshSubjects(); err != nil {
							return err
						}
					}

					outCh <- &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_UserEvent{
							UserEvent: &dest,
						},
					}

				case notifi.JobTopic:
					var dest notificationsevents.JobEvent
					if err := protoutils.UnmarshalPartialJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					outCh <- &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_JobEvent{
							JobEvent: &dest,
						},
					}

				case notifi.JobGradeTopic:
					// Make sure the job grade is included
					if len(parts) < 2 {
						continue
					}
					grade, err := strconv.ParseInt(parts[1], 10, 32)
					if err != nil {
						continue
					}
					if currentUserInfo.GetJobGrade() < int32(grade) {
						continue
					}
					var dest notificationsevents.JobGradeEvent
					if err := protoutils.UnmarshalPartialJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					outCh <- &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_JobGradeEvent{
							JobGradeEvent: &dest,
						},
					}

				case notifi.SystemTopic:
					var dest notificationsevents.SystemEvent
					if err := protoutils.UnmarshalPartialJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					outCh <- &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_SystemEvent{
							SystemEvent: &dest,
						},
					}

				case notifi.ObjectTopic:
					var dest notificationsclientview.ObjectEvent
					if err := protoutils.UnmarshalPartialJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					if !s.shouldDeliverObjectEvent(&dest, currentUserInfo) {
						continue
					}

					outCh <- &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_ObjectEvent{
							ObjectEvent: &dest,
						},
					}

				case notifi.MailerTopic:
					var dest mailerevents.MailerEvent
					if err := protoutils.UnmarshalPartialJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					outCh <- &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_MailerEvent{
							MailerEvent: &dest,
						},
					}
				}
			}
		}
	})

	return g.Wait()
}
