package notifications

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	notifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	pbnotifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/notifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	pbmailer "github.com/fivenet-app/fivenet/v2025/services/mailer"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const feedFetch = 8

func (s *Server) buildSubjects(ctx context.Context, userInfo userinfo.UserInfo) ([]string, error) {
	subjects := []string{
		fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, userInfo.UserId),
		fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, userInfo.Job),
		fmt.Sprintf("%s.%s.%s.>", notifi.BaseSubject, notifi.JobGradeTopic, userInfo.Job),
		fmt.Sprintf("%s.%s", notifi.BaseSubject, notifi.SystemTopic),
	}

	// Clone user info and disable superuser
	userInfo.Superuser = false
	emails, err := pbmailer.ListUserEmails(ctx, s.db, &userInfo, nil, false)
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedStream)
	}

	for _, email := range emails {
		subjects = append(subjects, fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.MailerTopic, email.Id))
	}

	return subjects, nil
}

func (s *Server) Stream(srv pbnotifications.NotificationsService_StreamServer) error {
	ctx := srv.Context()
	userInfo, ok := auth.GetUserInfoFromContext(ctx)
	if !ok {
		return ErrFailedStream
	}

	// Track changes to user info, so we can send an updated user info to the user
	currentUserInfo := userInfo.Clone()

	// User info update ticker
	updateTicker := time.NewTicker(15 * time.Second)
	defer updateTicker.Stop()

	// Check user token validity and update if necessary
	data, stop, err := s.checkUser(ctx, currentUserInfo)
	if err != nil {
		return err
	}

	notificationCount, err := s.getNotificationCount(ctx, userInfo.UserId)
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}

	if err := srv.Send(&pbnotifications.StreamResponse{
		NotificationCount: notificationCount,
		Data:              data,
		Restart:           &stop,
	}); err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}

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

			// TODO handle client state changes (update consumer subjects, etc.)
			_ = msg
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
		subjects, err := s.buildSubjects(ctx, currentUserInfo)
		if err != nil {
			return err
		}

		// Create pull consumer with multi-filter
		consCfg := jetstream.ConsumerConfig{
			FilterSubjects: subjects,
			DeliverPolicy:  jetstream.DeliverNewPolicy,
			AckPolicy:      jetstream.AckNonePolicy,
			MaxWaiting:     8,
		}
		cons, err := s.js.CreateConsumer(gctx, notifi.StreamName, consCfg)
		if err != nil {
			return fmt.Errorf("failed to create consumer. %w", err)
		}

		for {
			batch, err := cons.Fetch(feedFetch,
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
					return nil
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
					}

					resp := &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_UserEvent{
							UserEvent: &dest,
						},
					}

					if err := srv.Send(resp); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

				case notifi.JobTopic:
					var dest notifications.JobEvent
					if err := protoutils.UnmarshalPartialPJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					resp := &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_JobEvent{
							JobEvent: &dest,
						},
					}

					if err := srv.Send(resp); err != nil {
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

					resp := &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_JobGradeEvent{
							JobGradeEvent: &dest,
						},
					}

					if err := srv.Send(resp); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

				case notifi.SystemTopic:
					var dest notifications.SystemEvent
					if err := protoutils.UnmarshalPartialPJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					resp := &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_SystemEvent{
							SystemEvent: &dest,
						},
					}

					if err := srv.Send(resp); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

				case notifi.MailerTopic:
					var dest mailer.MailerEvent
					if err := protoutils.UnmarshalPartialPJSON(m.Data(), &dest); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}

					resp := &pbnotifications.StreamResponse{
						NotificationCount: notificationCount,
						Data: &pbnotifications.StreamResponse_MailerEvent{
							MailerEvent: &dest,
						},
					}

					if err := srv.Send(resp); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}
				}
			}
		}
	})

	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return nil

			case <-updateTicker.C:
				// Check user token validity
				data, stop, err := s.checkUser(ctx, currentUserInfo)
				if err != nil {
					return err
				}
				if data != nil {
					resp := &pbnotifications.StreamResponse{
						Data: data,
					}
					if err := srv.Send(resp); err != nil {
						return errswrap.NewError(err, ErrFailedStream)
					}
				}

				if stop {
					// End stream if we should "stop"
					// TODO this doesn't end the stream properly, we need to return an error
					return nil
				}

				// Make sure the notification is in sync (again)
				notificationCount, err = s.getNotificationCount(ctx, userInfo.UserId)
				if err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}
			}
		}
	})

	return g.Wait()
}

func (s *Server) checkUser(ctx context.Context, currentUserInfo userinfo.UserInfo) (*pbnotifications.StreamResponse_UserEvent, bool, error) {
	newUserInfo, err := s.ui.GetUserInfo(ctx, currentUserInfo.UserId, currentUserInfo.AccountId)
	if err != nil {
		return nil, true, errswrap.NewError(err, ErrFailedStream)
	}

	if currentUserInfo.LastChar != nil && *newUserInfo.LastChar != currentUserInfo.UserId && s.appCfg.Get().Auth.LastCharLock {
		if !currentUserInfo.CanBeSuperuser && !currentUserInfo.Superuser {
			return nil, true, auth.ErrCharLock
		}
	}

	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, true, auth.ErrInvalidToken
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, true, auth.ErrInvalidToken
	}

	// Either token should be renewed or new user info is not equal
	if time.Until(claims.ExpiresAt.Time) <= auth.TokenRenewalTime || !currentUserInfo.Equal(newUserInfo) {
		// Cause client to refresh token
		return &pbnotifications.StreamResponse_UserEvent{UserEvent: &notifications.UserEvent{
			Data: &notifications.UserEvent_RefreshToken{
				RefreshToken: true,
			},
		}}, true, nil
	}

	return nil, false, nil
}
