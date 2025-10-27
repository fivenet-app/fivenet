package notifications

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	notifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	pbuserinfo "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbnotifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/notifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/grpcws"
	natsutils "github.com/fivenet-app/fivenet/v2025/pkg/nats"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2025/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	pbmailer "github.com/fivenet-app/fivenet/v2025/services/mailer"
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

	// Clone user info and disable superuser
	userInfo.Superuser = false
	emails, err := pbmailer.ListUserEmails(ctx, s.db, userInfo, nil, false)
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
	clientViewSubject := []string{}

	notificationCount, err := s.getNotificationCount(ctx, userInfo.GetUserId())
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
			case *pbnotifications.StreamRequest_ClientView:
				clientView := d.ClientView
				if clientView == nil {
					continue // Skip nil client view
				}

				info, err := consumer.Info(gctx)
				if err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}
				cfg := info.Config

				// If client view is not "unspecified", add specific subject for it
				if clientView.Id != nil && clientView.GetType() > notifications.ObjectType_OBJECT_TYPE_UNSPECIFIED {
					gAccess := access.GetAccess(clientView.GetType().ToAccessKey())
					if gAccess != nil {
						check, err := gAccess.CanUserAccessTarget(gctx, clientView.GetId(), userInfo, 2)
						if err != nil {
							if !errors.Is(err, qrm.ErrNoRows) {
								return errswrap.NewError(err, ErrFailedStream)
							}
						}

						if !check {
							s.logger.Warn("user does not have access to the object",
								zap.Int32("user_id", userInfo.GetUserId()),
								zap.String("object_type", clientView.GetType().String()),
								zap.Int64("object_id", clientView.GetId()),
							)
							continue
						}
					}

					// Generate subject for the client view
					clientViewSubject = []string{
						fmt.Sprintf("%s.%s.%s.%d", notifi.BaseSubject, notifi.ObjectTopic, clientView.GetType().ToNatsKey(), clientView.GetId()),
					}
				}

				subjectsMu.Lock()
				cfg.FilterSubjects = []string{}
				cfg.FilterSubjects = append(cfg.FilterSubjects, baseSubjects...)
				cfg.FilterSubjects = append(cfg.FilterSubjects, additionalSubjects...)
				cfg.FilterSubjects = append(cfg.FilterSubjects, clientViewSubject...)
				subjectsMu.Unlock()

				// Update consumer
				if _, err := s.js.UpdateConsumer(gctx, notifi.StreamName, cfg); err != nil {
					return fmt.Errorf("failed to update consumer. %w", err)
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
				case notifi.UserTopic:
					var dest notifications.UserEvent
					if err := protoutils.UnmarshalPartialJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					switch d := dest.GetData().(type) {
					case *notifications.UserEvent_Notification:
						notificationCount++

					case *notifications.UserEvent_NotificationsReadCount:
						if notificationCount-d.NotificationsReadCount <= 0 {
							notificationCount = 0
						} else {
							notificationCount -= d.NotificationsReadCount
						}

					case *notifications.UserEvent_UserInfoChanged:
						currentUserInfo.Job = d.UserInfoChanged.GetNewJob()
						currentUserInfo.JobGrade = d.UserInfoChanged.GetNewJobGrade()

						baseSubjects, additionalSubjects, err = s.buildSubjects(gctx, currentUserInfo)
						if err != nil {
							return errswrap.NewError(err, ErrFailedStream)
						}

						info, err := consumer.Info(gctx)
						if err != nil {
							return errswrap.NewError(err, ErrFailedStream)
						}

						cfg := info.Config
						subjectsMu.Lock()
						cfg.FilterSubjects = []string{}
						// Rebuild filter subjects with the new user info
						cfg.FilterSubjects = append(cfg.FilterSubjects, baseSubjects...)
						cfg.FilterSubjects = append(cfg.FilterSubjects, additionalSubjects...)
						cfg.FilterSubjects = append(cfg.FilterSubjects, clientViewSubject...)
						subjectsMu.Unlock()

						// Update consumer subjects
						if _, err := s.js.UpdateConsumer(gctx, notifi.StreamName, cfg); err != nil {
							return fmt.Errorf("failed to update consumer. %w", err)
						}
					}

					outCh <- &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_UserEvent{
							UserEvent: &dest,
						},
					}

				case notifi.JobTopic:
					var dest notifications.JobEvent
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
					var dest notifications.JobGradeEvent
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
					var dest notifications.SystemEvent
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
					var dest notifications.ObjectEvent
					if err := protoutils.UnmarshalPartialJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					// Skip if the object event is from the current user
					if dest.UserId != nil && dest.GetUserId() == currentUserInfo.GetUserId() {
						continue
					}

					// Check if the user has access to the object for job specific objects
					if dest.GetType() != notifications.ObjectType_OBJECT_TYPE_UNSPECIFIED &&
						dest.GetType() != notifications.ObjectType_OBJECT_TYPE_DOCUMENT &&
						dest.GetType() != notifications.ObjectType_OBJECT_TYPE_WIKI_PAGE {
						// No job specified or job doesn't match the user's job
						if dest.Job == nil || userInfo.GetJob() != dest.GetJob() {
							continue
						}
					}

					outCh <- &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_ObjectEvent{
							ObjectEvent: &dest,
						},
					}

				case notifi.MailerTopic:
					var dest mailer.MailerEvent
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
