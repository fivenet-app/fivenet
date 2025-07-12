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

func (s *Server) buildSubjects(ctx context.Context, userInfo *pbuserinfo.UserInfo) ([]string, []string, error) {
	baseSubjects := []string{
		fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, userInfo.UserId),
		fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, userInfo.Job),
		fmt.Sprintf("%s.%s.%s.>", notifi.BaseSubject, notifi.JobGradeTopic, userInfo.Job),
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
		additionalSubjects = append(additionalSubjects, fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.MailerTopic, email.Id))
	}

	return baseSubjects, additionalSubjects, nil
}

func (s *Server) Stream(srv pbnotifications.NotificationsService_StreamServer) error {
	ctx := srv.Context()
	userInfo, ok := auth.GetUserInfoFromContext(ctx)
	if !ok {
		return ErrFailedStream
	}

	// Track changes to user info, so we can send an updated user info to the user
	currentUserInfo := userInfo.Clone()

	if _, err := s.js.PublishAsyncProto(ctx, userinfo.PollSubject, &pbuserinfo.PollReq{
		AccountId: currentUserInfo.AccountId,
		UserId:    currentUserInfo.UserId,
	}); err != nil {
		s.logger.Error("failed to publish userinfo.poll.request", zap.Int32("user_id", currentUserInfo.UserId), zap.Error(err))
	}

	subjectsMu := &sync.Mutex{}
	baseSubjects, additionalSubjects, err := s.buildSubjects(ctx, currentUserInfo)
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}
	clientViewSubject := []string{}

	notificationCount, err := s.getNotificationCount(ctx, userInfo.UserId)
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}

	meta := metadata.ExtractIncoming(ctx)
	connId := meta.Get(grpcws.ConnectionIdHeader)

	// Create durable pull consumer with multi-filter, required to update filter subjects dynamically
	consCfg := jetstream.ConsumerConfig{
		Durable:           natsutils.GenerateConsumerName(userInfo.AccountId, userInfo.UserId, connId),
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
	defer s.js.DeleteConsumer(ctx, notifi.StreamName, consCfg.Durable)

	metricLastUserSession.SetToCurrentTime()

	metricActiveUserSessions.Inc()
	defer metricActiveUserSessions.Dec()

	// Central pipe: all feeds push messages into outCh
	outCh := make(chan *pbnotifications.StreamResponse, 256)
	defer close(outCh)
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			msg, err := srv.Recv()
			if err == io.EOF {
				return err
			}
			if err != nil {
				return err
			}
			if msg == nil {
				continue // Skip nil messages
			}

			switch d := msg.Data.(type) {
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
				if clientView.Id != nil && clientView.Type > notifications.ObjectType_OBJECT_TYPE_UNSPECIFIED {
					gAccess := access.GetAccess(clientView.Type.ToAccessKey())
					if gAccess != nil {
						check, err := gAccess.CanUserAccessTarget(gctx, *clientView.Id, userInfo, 2)
						if err != nil {
							if !errors.Is(err, qrm.ErrNoRows) {
								return errswrap.NewError(err, ErrFailedStream)
							}
						}

						if !check {
							s.logger.Warn("user does not have access to the object", zap.Int32("user_id", userInfo.UserId), zap.String("object_type", clientView.Type.String()), zap.Uint64p("object_id", clientView.Id))
							continue
						}
					}

					// Generate subject for the client view
					clientViewSubject = []string{
						fmt.Sprintf("%s.%s.%s.%d", notifi.BaseSubject, notifi.ObjectTopic, clientView.Type.ToNatsKey(), *clientView.Id),
					}
				}

				subjectsMu.Lock()
				cfg.FilterSubjects = append(baseSubjects, additionalSubjects...)
				cfg.FilterSubjects = append(cfg.FilterSubjects, clientViewSubject...)
				subjectsMu.Unlock()

				// Update consumer
				if _, err := s.js.UpdateConsumer(ctx, notifi.StreamName, cfg); err != nil {
					return fmt.Errorf("failed to update consumer. %w", err)
				}
			}
		}
	})

	// Writer goroutine – single gRPC send loop
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
			batch, err := consumer.Fetch(feedFetch,
				jetstream.FetchMaxWait(2*time.Second))
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
					s.logger.Warn("nil notification message received via message queue", zap.Int32("user_id", currentUserInfo.UserId))
					continue
				}

				if err := m.Ack(); err != nil {
					s.logger.Error("failed to ack notification message", zap.Error(err))
				}

				topic, parts := notifi.SplitSubject(m.Subject())
				switch topic {
				case notifi.UserTopic:
					var dest notifications.UserEvent
					if err := protoutils.UnmarshalPartialPJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					switch d := dest.Data.(type) {
					case *notifications.UserEvent_Notification:
						notificationCount++

					case *notifications.UserEvent_NotificationsReadCount:
						if notificationCount-d.NotificationsReadCount <= 0 {
							notificationCount = 0
						} else {
							notificationCount -= d.NotificationsReadCount
						}

					case *notifications.UserEvent_UserInfoChanged:
						currentUserInfo.Job = d.UserInfoChanged.NewJob
						currentUserInfo.JobGrade = d.UserInfoChanged.NewJobGrade

						baseSubjects, additionalSubjects, err = s.buildSubjects(ctx, currentUserInfo)
						if err != nil {
							return errswrap.NewError(err, ErrFailedStream)
						}

						info, err := consumer.Info(gctx)
						if err != nil {
							return errswrap.NewError(err, ErrFailedStream)
						}

						cfg := info.Config
						subjectsMu.Lock()
						cfg.FilterSubjects = append(baseSubjects, additionalSubjects...)
						cfg.FilterSubjects = append(cfg.FilterSubjects, clientViewSubject...)
						subjectsMu.Unlock()

						// Update consumer subjects
						if _, err := s.js.UpdateConsumer(ctx, notifi.StreamName, cfg); err != nil {
							return fmt.Errorf("failed to update consumer. %w", err)
						}
					}

					if err := srv.Send(&pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_UserEvent{
							UserEvent: &dest,
						},
					}); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

				case notifi.JobTopic:
					var dest notifications.JobEvent
					if err := protoutils.UnmarshalPartialPJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					if err := srv.Send(&pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_JobEvent{
							JobEvent: &dest,
						},
					}); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
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
					if currentUserInfo.JobGrade < int32(grade) {
						continue
					}
					var dest notifications.JobGradeEvent
					if err := protoutils.UnmarshalPartialPJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					if err := srv.Send(&pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_JobGradeEvent{
							JobGradeEvent: &dest,
						},
					}); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

				case notifi.SystemTopic:
					var dest notifications.SystemEvent
					if err := protoutils.UnmarshalPartialPJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					if err := srv.Send(&pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_SystemEvent{
							SystemEvent: &dest,
						},
					}); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

				case notifi.ObjectTopic:
					var dest notifications.ObjectEvent
					if err := protoutils.UnmarshalPartialPJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					// Skip if the object event is from the current user
					if dest.UserId != nil && *dest.UserId == currentUserInfo.UserId {
						continue
					}

					// Check if the user has access to the object for job specific objects
					if dest.Type != notifications.ObjectType_OBJECT_TYPE_UNSPECIFIED && dest.Type != notifications.ObjectType_OBJECT_TYPE_DOCUMENT && dest.Type != notifications.ObjectType_OBJECT_TYPE_WIKI_PAGE {
						if dest.Job == nil {
							continue
						}
						// Job doesn't match the user's job
						if userInfo.Job != *dest.Job {
							continue
						}
					}

					if err := srv.Send(&pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_ObjectEvent{
							ObjectEvent: &dest,
						},
					}); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

				case notifi.MailerTopic:
					var dest mailer.MailerEvent
					if err := protoutils.UnmarshalPartialPJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					if err := srv.Send(&pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_MailerEvent{
							MailerEvent: &dest,
						},
					}); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}
				}
			}
		}
	})

	return g.Wait()
}
