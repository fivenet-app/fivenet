package notificator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	notifications "github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	tNotifications = table.FivenetNotifications
	tUsers         = table.Users.AS("user")
	tJobs          = table.Jobs
	tJobGrades     = table.JobGrades
	tJobProps      = table.FivenetJobProps.AS("jobprops")
)

var (
	ErrFailedRequest = status.Error(codes.InvalidArgument, "errors.NotificatorService.ErrFailedRequest")
	ErrFailedStream  = status.Error(codes.InvalidArgument, "errors.NotificatorService.ErrFailedStream")
)

type Server struct {
	NotificatorServiceServer

	logger   *zap.Logger
	db       *sql.DB
	p        perms.Permissions
	tm       *auth.TokenMgr
	ui       userinfo.UserInfoRetriever
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	appCfg   appconfig.IConfig
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	DB        *sql.DB
	Perms     perms.Permissions
	TM        *auth.TokenMgr
	UI        userinfo.UserInfoRetriever
	JS        *events.JSWrapper
	Enricher  *mstlystcdata.Enricher
	AppConfig appconfig.IConfig
}

func NewServer(p Params) *Server {
	s := &Server{
		logger:   p.Logger,
		db:       p.DB,
		p:        p.Perms,
		tm:       p.TM,
		ui:       p.UI,
		js:       p.JS,
		enricher: p.Enricher,
		appCfg:   p.AppConfig,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterNotificatorServiceServer(srv, s)
}

func (s *Server) GetNotifications(ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tNotifications := tNotifications.AS("notification")
	condition := tNotifications.UserID.EQ(jet.Int32(userInfo.UserId))
	if req.IncludeRead != nil && !*req.IncludeRead {
		condition = condition.AND(tNotifications.ReadAt.IS_NULL())
	}

	if len(req.Categories) > 0 {
		categoryIds := make([]jet.Expression, len(req.Categories))
		for i := 0; i < len(req.Categories); i++ {
			categoryIds[i] = jet.Int16(int16(req.Categories[i]))
		}

		condition = condition.AND(tNotifications.Category.IN(categoryIds...))
	}

	countStmt := tNotifications.
		SELECT(
			jet.COUNT(tNotifications.ID).AS("datacount.totalcount"),
		).
		FROM(tNotifications).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrFailedRequest)
		}
	}

	pag, limit := req.Pagination.GetResponse(count.TotalCount)
	resp := &GetNotificationsResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tNotifications.
		SELECT(
			tNotifications.ID,
			tNotifications.CreatedAt,
			tNotifications.ReadAt,
			tNotifications.UserID,
			tNotifications.Title,
			tNotifications.Type,
			tNotifications.Content,
			tNotifications.Category,
			tNotifications.Data,
		).
		FROM(tNotifications).
		WHERE(
			condition,
		).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(tNotifications.ID.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Notifications); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrFailedRequest)
		}
	}

	resp.Pagination.Update(len(resp.Notifications))

	return resp, nil
}

func (s *Server) MarkNotifications(ctx context.Context, req *MarkNotificationsRequest) (*MarkNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tNotifications.UserID.EQ(
		jet.Int32(userInfo.UserId)).AND(
		tNotifications.ReadAt.IS_NULL(),
	)
	// If not all
	if len(req.Ids) > 0 {
		ids := make([]jet.Expression, len(req.Ids))
		for i := 0; i < len(req.Ids); i++ {
			ids[i] = jet.Uint64(req.Ids[i])
		}
		condition = condition.AND(tNotifications.ID.IN(ids...))
	} else if req.All == nil || !*req.All {
		return &MarkNotificationsResponse{}, nil
	}

	stmt := tNotifications.
		UPDATE(
			tNotifications.ReadAt,
		).
		SET(
			tNotifications.ReadAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(condition)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	return &MarkNotificationsResponse{
		Updated: uint64(rows),
	}, nil
}

func (s *Server) getNotificationCount(ctx context.Context, userId int32) (int32, error) {
	stmt := tNotifications.
		SELECT(
			jet.COUNT(tNotifications.ID).AS("count"),
		).
		FROM(tNotifications).
		WHERE(jet.AND(
			tNotifications.UserID.EQ(jet.Int32(userId)),
			tNotifications.ReadAt.IS_NULL(),
		)).
		ORDER_BY(tNotifications.ID.DESC())

	var dest struct {
		Count int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, errswrap.NewError(err, ErrFailedStream)
		}
	}

	return dest.Count, nil
}

func (s *Server) Stream(req *StreamRequest, srv NotificatorService_StreamServer) error {
	userInfo, ok := auth.GetUserInfoFromContext(srv.Context())
	if !ok {
		return ErrFailedStream
	}

	// Track changes to user info, so we can send an updated user info to the user
	currentUserInfo := userInfo.Clone()

	// Setup consumer
	c, err := s.js.CreateConsumer(srv.Context(), notifi.StreamName, jetstream.ConsumerConfig{
		FilterSubjects: []string{
			fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserTopic, currentUserInfo.UserId),
			fmt.Sprintf("%s.%s.%s", notifi.BaseSubject, notifi.JobTopic, currentUserInfo.Job),
			fmt.Sprintf("%s.%s", notifi.BaseSubject, notifi.SystemTopic),
		},
		DeliverPolicy: jetstream.DeliverNewPolicy,
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
				break
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

			_, topic, _ := notifi.SplitSubject(msg.Subject())
			switch topic {
			case notifi.UserTopic:
				var dest notifications.UserEvent
				if err := proto.Unmarshal(msg.Data(), &dest); err != nil {
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
				if err := proto.Unmarshal(msg.Data(), &dest); err != nil {
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

			case notifi.SystemTopic:
				// No events yet...
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

func (s *Server) getCharacter(ctx context.Context, charId int32) (*users.User, *users.JobProps, string, error) {
	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUsers.Group.AS("group"),
			tJobs.Label.AS("user.job_label"),
			tJobGrades.Label.AS("user.job_grade_label"),
			tJobProps.Theme,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.LogoURL,
		).
		FROM(
			tUsers.
				LEFT_JOIN(tJobs,
					tJobs.Name.EQ(tUsers.Job),
				).
				LEFT_JOIN(tJobGrades,
					jet.AND(
						tJobGrades.Grade.EQ(tUsers.JobGrade),
						tJobGrades.JobName.EQ(tUsers.Job),
					),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tJobs.Name),
				),
		).
		WHERE(
			tUsers.ID.EQ(jet.Int32(charId)),
		).
		LIMIT(1)

	var dest struct {
		users.User
		Group    string
		JobProps *users.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, nil, "", err
	}

	if dest.JobProps != nil {
		s.enricher.EnrichJobName(dest.JobProps)
	}

	return &dest.User, dest.JobProps, dest.Group, nil
}
