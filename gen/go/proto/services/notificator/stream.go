package notificator

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	notifications "github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	pbmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *Server) Stream(req *StreamRequest, srv NotificatorService_StreamServer) error {
	userInfo, ok := auth.GetUserInfoFromContext(srv.Context())
	if !ok {
		return ErrFailedStream
	}

	// Track changes to user info, so we can send an updated user info to the user
	currentUserInfo := userInfo.Clone()

	subjects := []string{
		fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, currentUserInfo.UserId),
		fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, currentUserInfo.Job),
		fmt.Sprintf("%s.%s.%s.>", notifi.BaseSubject, notifi.JobGradeTopic, currentUserInfo.Job),
		fmt.Sprintf("%s.%s", notifi.BaseSubject, notifi.SystemTopic),
	}

	// Clone user info and disable superuser
	cloned := currentUserInfo.Clone()
	cloned.SuperUser = false
	emails, err := pbmailer.ListUserEmails(srv.Context(), s.db, &cloned, nil, false)
	if err != nil {
		return ErrFailedStream
	}

	for _, email := range emails {
		subjects = append(subjects, fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.MailerTopic, email.Id))
	}

	// Setup consumer
	c, err := s.js.CreateConsumer(srv.Context(), notifi.StreamName, jetstream.ConsumerConfig{
		FilterSubjects: subjects,
		DeliverPolicy:  jetstream.DeliverNewPolicy,
	})
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}

	cons, err := c.Messages()
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}
	defer cons.Stop()

	msgCh := make(chan jetstream.Msg, 8)
	go func() {
		for {
			msg, err := cons.Next()
			if err != nil {
				close(msgCh)
				return
			}

			msgCh <- msg
		}
	}()

	// Update Ticker
	updateTicker := time.NewTicker(35 * time.Second)
	defer updateTicker.Stop()

	// Check user token validity and update if necessary
	data, stop, err := s.checkUser(srv.Context(), currentUserInfo)
	if err != nil {
		return err
	}

	notsCount, err := s.getNotificationCount(srv.Context(), userInfo.UserId)
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}

	if err := srv.Send(&StreamResponse{
		NotificationCount: notsCount,
		Data:              data,
		Restart:           &stop,
	}); err != nil {
		return errswrap.NewError(err, ErrFailedStream)
	}

	for {
		select {
		case <-srv.Context().Done():
			return nil

		case <-updateTicker.C:
			// Check user token validity
			data, stop, err := s.checkUser(srv.Context(), currentUserInfo)
			if err != nil {
				return err
			}
			if data != nil {
				resp := &StreamResponse{
					Data: data,
				}
				if err := srv.Send(resp); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}
			}

			if stop {
				// End stream if we should "stop"
				return nil
			}

			// Make sure the notification is in sync (again)
			notsCount, err = s.getNotificationCount(srv.Context(), userInfo.UserId)
			if err != nil {
				return errswrap.NewError(err, ErrFailedStream)
			}

		case msg := <-msgCh:
			// Publish notifications sent directly to user via the message queue
			if msg == nil {
				s.logger.Warn("nil notification message received via message queue", zap.Int32("user_id", currentUserInfo.UserId))
				return nil
			}

			if err := msg.Ack(); err != nil {
				s.logger.Error("failed to ack notification message", zap.Error(err))
			}

			_, topic, parts := notifi.SplitSubject(msg.Subject())
			switch topic {
			case notifi.UserTopic:
				var dest notifications.UserEvent
				if err := protojson.Unmarshal(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				switch dest.Data.(type) {
				case *notifications.UserEvent_Notification:
					notsCount++
				}

				resp := &StreamResponse{
					NotificationCount: notsCount,
					Data: &StreamResponse_UserEvent{
						UserEvent: &dest,
					},
				}

				if err := srv.Send(resp); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

			case notifi.JobTopic:
				var dest notifications.JobEvent
				if err := protojson.Unmarshal(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				resp := &StreamResponse{
					NotificationCount: notsCount,
					Data: &StreamResponse_JobEvent{
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
				grade, err := strconv.Atoi(parts[1])
				if err != nil {
					continue
				}
				if currentUserInfo.JobGrade < int32(grade) {
					continue
				}
				var dest notifications.JobGradeEvent
				if err := protojson.Unmarshal(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				resp := &StreamResponse{
					NotificationCount: notsCount,
					Data: &StreamResponse_JobGradeEvent{
						JobGradeEvent: &dest,
					},
				}

				if err := srv.Send(resp); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

			case notifi.SystemTopic:
				// No events yet...

			case notifi.MailerTopic:
				var dest mailer.MailerEvent
				if err := protojson.Unmarshal(msg.Data(), &dest); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}

				resp := &StreamResponse{
					NotificationCount: notsCount,
					Data: &StreamResponse_MailerEvent{
						MailerEvent: &dest,
					},
				}

				if err := srv.Send(resp); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}
			}
		}
	}
}

func (s *Server) checkUser(ctx context.Context, currentUserInfo userinfo.UserInfo) (isStreamResponse_Data, bool, error) {
	newUserInfo, err := s.ui.GetUserInfo(ctx, currentUserInfo.UserId, currentUserInfo.AccountId)
	if err != nil {
		return nil, true, errswrap.NewError(err, ErrFailedStream)
	}

	if currentUserInfo.LastChar != nil && *newUserInfo.LastChar != currentUserInfo.UserId && s.appCfg.Get().Auth.LastCharLock {
		if !currentUserInfo.CanBeSuper && !currentUserInfo.SuperUser {
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
		return &StreamResponse_UserEvent{UserEvent: &notifications.UserEvent{
			Data: &notifications.UserEvent_RefreshToken{
				RefreshToken: true,
			},
		}}, true, nil
	}

	return nil, false, nil
}
