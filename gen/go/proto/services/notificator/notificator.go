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
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/notifi"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

const pingTickerTime = 35 * time.Second

type Server struct {
	NotificatorServiceServer

	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
	tm     *auth.TokenMgr
	ui     userinfo.UserInfoRetriever
	events *events.Eventus
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	Perms  perms.Permissions
	TM     *auth.TokenMgr
	UI     userinfo.UserInfoRetriever
	Events *events.Eventus
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger,
		db:     p.DB,
		p:      p.Perms,
		tm:     p.TM,
		ui:     p.UI,
		events: p.Events,
	}

	return s
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
		return nil, ErrFailedRequest
	}

	pag, limit := req.Pagination.GetResponse()
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
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrFailedRequest
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Notifications))

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
		return nil, ErrFailedRequest
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, ErrFailedRequest
	}

	return &MarkNotificationsResponse{
		Updated: uint64(rows),
	}, nil
}

func (s *Server) getNotifications(ctx context.Context, userId int32, lastId uint64) ([]*notifications.Notification, error) {
	nots := tNotifications.AS("notification")

	stmt := nots.
		SELECT(
			nots.ID,
			nots.Title,
			nots.Type,
			nots.Content,
			nots.Data,
		).
		FROM(nots).
		WHERE(
			jet.AND(
				nots.UserID.EQ(jet.Int32(userId)),
				nots.ID.GT(jet.Uint64(lastId)),
				nots.ReadAt.IS_NULL(),
			),
		).
		ORDER_BY(nots.ID.DESC()).
		LIMIT(10)

	var dest []*notifications.Notification
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, ErrFailedStream
	}

	return dest, nil
}

func (s *Server) Stream(req *StreamRequest, srv NotificatorService_StreamServer) error {
	userInfo, ok := auth.GetUserInfoFromContext(srv.Context())
	if !ok {
		return ErrFailedStream
	}

	// Track changes to user info, so we can send an updated user info to the user
	currentUserInfo := userinfo.UserInfo{
		AccId:        userInfo.AccId,
		UserId:       userInfo.UserId,
		Job:          userInfo.Job,
		JobGrade:     userInfo.JobGrade,
		OrigJob:      userInfo.OrigJob,
		OrigJobGrade: userInfo.OrigJobGrade,
		Group:        userInfo.Group,
		SuperUser:    userInfo.SuperUser,
	}

	msgCh := make(chan *nats.Msg, 8)
	sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.%d", notifi.BaseSubject, notifi.UserNotification, currentUserInfo.UserId), msgCh, nats.DeliverNew())
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	// Ping pingTicker to ensure better stream quality
	pingTicker := time.NewTicker(pingTickerTime * 2)
	defer pingTicker.Stop()

	// Update Ticker
	updateTicker := time.NewTicker(40 * time.Second)
	defer updateTicker.Stop()

	// Check user token validity
	if stop, err := s.checkUser(srv, req, &currentUserInfo); err != nil {
		return err
	} else if stop {
		// End stream if we should "stop"
		return nil
	}

	// Check if user has any (unsent/unread) notifications
	if err := s.checkNotifications(srv, req, &currentUserInfo); err != nil {
		return err
	}

	for {
		resp := &StreamResponse{
			LastId: req.LastId,
		}

		select {
		case <-srv.Context().Done():
			return nil

		case t := <-pingTicker.C:
			resp.Data = &StreamResponse_Ping{
				Ping: t.String(),
			}

		case <-updateTicker.C:
			// Check for new user notifications
			if err := s.checkNotifications(srv, req, &currentUserInfo); err != nil {
				return err
			}

			// Check user token validity
			if stop, err := s.checkUser(srv, req, &currentUserInfo); err != nil {
				return err
			} else if stop {
				// End stream if we should "stop"
				return nil
			}

			// Make sure message queue subscription is still valid, otherwise restart stream
			if !sub.IsValid() {
				restart := true
				if err := srv.Send(&StreamResponse{
					LastId:  req.LastId,
					Restart: &restart,
				}); err != nil {
					return ErrFailedStream
				}
			}

		case msg := <-msgCh:
			// Publish notifications sent directly to user via the message queue
			msg.Ack()

			var dest notifications.Notification
			if err := proto.Unmarshal(msg.Data, &dest); err != nil {
				return err
			}

			resp.Data = &StreamResponse_Notifications{
				Notifications: &StreamNotifications{
					Notifications: []*notifications.Notification{&dest},
				},
			}

			if err := srv.Send(resp); err != nil {
				return ErrFailedStream
			}
		}
	}
}

func (s *Server) checkNotifications(srv NotificatorService_StreamServer, req *StreamRequest, userInfo *userinfo.UserInfo) error {
	resp := &StreamResponse{
		LastId: req.LastId,
	}

	notifications, err := s.getNotifications(srv.Context(), userInfo.UserId, req.LastId)
	if err != nil {
		return ErrFailedStream
	}
	// Update last notification id of the user and send notifications
	if len(notifications) > 0 {
		req.LastId = notifications[0].Id
		resp.LastId = notifications[0].Id

		resp.Data = &StreamResponse_Notifications{
			Notifications: &StreamNotifications{
				Notifications: notifications,
			},
		}

		if err := srv.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) checkUser(srv NotificatorService_StreamServer, req *StreamRequest, userInfo *userinfo.UserInfo) (bool, error) {
	claims, restart, tu, err := s.checkAndUpdateToken(srv.Context())
	if err != nil {
		return false, ErrFailedStream
	}

	resp := &StreamResponse{
		LastId: req.LastId,
	}
	if tu != nil && claims.CharID > 0 {
		if err := s.checkAndUpdateUserInfo(srv.Context(), tu, userInfo); err != nil {
			return false, ErrFailedStream
		}

		resp.Data = &StreamResponse_Token{
			Token: tu,
		}
	}

	if restart {
		resp.Restart = &restart

		if err := srv.Send(resp); err != nil {
			return false, ErrFailedStream
		}
	}

	return false, nil
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

func (s *Server) checkAndUpdateUserInfo(ctx context.Context, tu *TokenUpdate, currentUserInfo *userinfo.UserInfo) error {
	userInfo, err := s.ui.GetUserInfo(ctx, currentUserInfo.UserId, currentUserInfo.AccId)
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
			tJobProps.QuickButtons,
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
		JobProps users.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, nil, "", err
	}

	return &dest.User, &dest.JobProps, dest.Group, nil
}
