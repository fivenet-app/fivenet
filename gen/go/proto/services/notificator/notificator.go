package notificator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/notifications"
	timestamp "github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/notifi"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/table"
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
	js       jetstream.JetStream
	enricher *mstlystcdata.Enricher
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	DB       *sql.DB
	Perms    perms.Permissions
	TM       *auth.TokenMgr
	UI       userinfo.UserInfoRetriever
	JS       jetstream.JetStream
	Enricher *mstlystcdata.Enricher
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
		condition = jet.AND(
			condition,
			tNotifications.ReadAt.IS_NULL(),
		)
	}

	countStmt := tNotifications.
		SELECT(
			jet.COUNT(tNotifications.ID).AS("datacount.totalcount"),
		).
		FROM(tNotifications).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
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
		FilterSubject: fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserNotification, currentUserInfo.UserId),
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
	updateTicker := time.NewTicker(45 * time.Second)
	defer updateTicker.Stop()

	// Check user token validity
	data, stop, err := s.checkUser(srv.Context(), currentUserInfo)
	if err != nil {
		return errswrap.NewError(err, ErrFailedStream)
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
		resp := &StreamResponse{
			NotificationCount: notsCount,
		}

		select {
		case <-srv.Context().Done():
			return nil

		case <-updateTicker.C:
			// Check user token validity
			data, stop, err := s.checkUser(srv.Context(), currentUserInfo)
			if err != nil {
				return errswrap.NewError(err, ErrFailedStream)
			}
			if data != nil {
				resp.Data = data

				if err := srv.Send(resp); err != nil {
					return errswrap.NewError(err, ErrFailedStream)
				}
			}

			if stop {
				// End stream if we should "stop"
				return nil
			}

			// Make sure the notification is in sync (again)
			resp.NotificationCount, err = s.getNotificationCount(srv.Context(), userInfo.UserId)
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

			var dest notifications.Notification
			if err := proto.Unmarshal(msg.Data(), &dest); err != nil {
				return errswrap.NewError(err, ErrFailedStream)
			}

			resp.Data = &StreamResponse_Notification{
				Notification: &dest,
			}
			resp.NotificationCount++

			if err := srv.Send(resp); err != nil {
				return errswrap.NewError(err, ErrFailedStream)
			}
		}
	}
}

func (s *Server) checkUser(ctx context.Context, userInfo userinfo.UserInfo) (isStreamResponse_Data, bool, error) {
	claims, restart, tu, err := s.checkAndUpdateToken(ctx)
	if err != nil {
		return nil, false, err
	}

	if tu != nil && claims.CharID > 0 {
		if err := s.checkAndUpdateUserInfo(ctx, tu, userInfo); err != nil {
			return nil, false, err
		}

		return &StreamResponse_Token{
			Token: tu,
		}, true, nil
	}

	return nil, restart, nil
}

func (s *Server) checkAndUpdateToken(ctx context.Context) (*auth.CitizenInfoClaims, bool, *TokenUpdate, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, true, nil, auth.ErrInvalidToken
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, true, nil, auth.ErrInvalidToken
	}

	if time.Until(claims.ExpiresAt.Time) <= auth.TokenRenewalTime {
		if claims.RenewedCount >= auth.TokenMaxRenews {
			return nil, true, nil, auth.ErrInvalidToken
		}

		// Increase re-newed count
		claims.RenewedCount++

		auth.SetTokenClaimsTimes(claims)
		newToken, err := s.tm.NewWithClaims(claims)
		if err != nil {
			return nil, true, nil, auth.ErrCheckToken
		}

		tu := &TokenUpdate{
			NewToken: &newToken,
			Expires:  timestamp.New(claims.ExpiresAt.Time),
		}

		return claims, true, tu, nil
	}

	return claims, false, nil, nil
}

func (s *Server) checkAndUpdateUserInfo(ctx context.Context, tu *TokenUpdate, currentUserInfo userinfo.UserInfo) error {
	userInfo, err := s.ui.GetUserInfo(ctx, currentUserInfo.UserId, currentUserInfo.AccountId)
	if err != nil {
		return err
	}

	// If the user is logged into a character, update user info and load permissions of user
	if currentUserInfo.Equal(userInfo) {
		return nil
	}

	char, jobProps, group, err := s.getCharacter(ctx, userInfo.UserId)
	if err != nil {
		return err
	}
	tu.UserInfo = char
	tu.JobProps = jobProps

	// Update current user info with new data from database
	currentUserInfo.UserId = char.UserId
	currentUserInfo.Job = char.Job
	currentUserInfo.JobGrade = char.JobGrade
	currentUserInfo.Group = group

	ps, err := s.p.GetPermissionsOfUser(&userinfo.UserInfo{
		UserId:   userInfo.UserId,
		Job:      userInfo.Job,
		JobGrade: userInfo.JobGrade,
	})
	if err != nil {
		return auth.ErrUserNoPerms
	}
	tu.Permissions = ps.GuardNames()

	if userInfo.SuperUser {
		tu.Permissions = append(tu.Permissions, common.SuperuserPermission)
	}

	attrs, err := s.p.FlattenRoleAttributes(userInfo.Job, userInfo.JobGrade)
	if err != nil {
		return auth.ErrUserNoPerms
	}
	tu.Permissions = append(tu.Permissions, attrs...)

	return nil
}

func (s *Server) getCharacter(ctx context.Context, charId int32) (*users.User, *users.JobProps, string, error) {
	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Identifier,
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
