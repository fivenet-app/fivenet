package notificator

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	timestamp "github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	tNotifications = table.FivenetNotifications
	tUsers         = table.Users.AS("user")
	tJobs          = table.Jobs
	tJobGrades     = table.JobGrades
	tJobProps      = table.FivenetJobProps.AS("jobprops")
)

var (
	InvalidRequestErr = status.Error(codes.InvalidArgument, "Invalid notificator stream request!")
)

const StartWaitTicks = 3

type Server struct {
	NotificatorServiceServer

	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
	tm     *auth.TokenMgr
	ui     userinfo.UserInfoRetriever
}

func NewServer(logger *zap.Logger, db *sql.DB, p perms.Permissions, tm *auth.TokenMgr, ui userinfo.UserInfoRetriever) *Server {
	return &Server{
		logger: logger,
		db:     db,
		p:      p,
		tm:     tm,
		ui:     ui,
	}
}

func (s *Server) PermissionStreamFuncOverride(ctx context.Context, srv interface{}, info *grpc.StreamServerInfo) (context.Context, error) {
	// Skip permission check for the notificator services
	return ctx, nil
}

func (s *Server) GetNotifications(ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tNotifications.UserID.EQ(jet.Int32(userInfo.UserId))
	if req.IncludeRead {
		condition = jet.AND(
			condition,
			tNotifications.ReadAt.IS_NOT_NULL(),
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
		return nil, err
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
			tNotifications.AllColumns,
		).
		FROM(tNotifications).
		WHERE(
			condition,
		).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, resp.Notifications); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Notifications))

	return resp, nil
}

func (s *Server) ReadNotifications(ctx context.Context, req *ReadNotificationsRequest) (*ReadNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	ids := make([]jet.Expression, len(req.Ids))
	for i := 0; i < len(req.Ids); i++ {
		ids[i] = jet.Uint64(req.Ids[i])
	}

	stmt := tNotifications.
		UPDATE(
			tNotifications.ReadAt,
		).
		SET(
			tNotifications.ReadAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			jet.AND(
				tNotifications.UserID.EQ(jet.Int32(userInfo.UserId)),
				tNotifications.ID.IN(ids...),
			),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &ReadNotificationsResponse{}, nil
}

func (s *Server) Stream(req *StreamRequest, srv NotificatorService_StreamServer) error {
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
		ORDER_BY(nots.ID.DESC()).
		LIMIT(10)

	userInfo, ok := auth.GetUserInfoFromContext(srv.Context())
	if !ok {
		return InvalidRequestErr
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

	waitTicks := StartWaitTicks
	for {
		resp := &StreamResponse{
			LastId: req.LastId,
			Token:  &TokenUpdate{},
		}

		q := stmt.
			WHERE(
				jet.AND(
					nots.UserID.EQ(jet.Int32(userInfo.UserId)),
					nots.ID.GT(jet.Uint64(req.LastId)),
				),
			)

		if err := q.QueryContext(srv.Context(), s.db, &resp.Notifications); err != nil {
			return err
		}

		// Update last notification id of user
		if len(resp.Notifications) > 0 {
			req.LastId = resp.Notifications[0].Id
			resp.LastId = resp.Notifications[0].Id
		}

		claims, restart, err := s.checkAndUpdateToken(srv.Context(), resp.Token)
		if err != nil {
			return err
		}
		if restart {
			resp.RestartStream = true
		}

		if waitTicks <= 0 {
			if claims.CharID > 0 {
				if err := s.checkAndUpdateUserInfo(srv.Context(), resp.Token, &currentUserInfo); err != nil {
					return err
				}
			}
			waitTicks = StartWaitTicks
		}

		if err := srv.Send(resp); err != nil {
			return err
		}

		resp.Notifications = nil

		select {
		case <-srv.Context().Done():
			return nil
		case <-time.After(20 * time.Second):
			waitTicks--
		}
	}
}

func (s *Server) checkAndUpdateToken(ctx context.Context, tu *TokenUpdate) (*auth.CitizenInfoClaims, bool, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, true, auth.InvalidTokenErr
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, true, auth.InvalidTokenErr
	}

	if time.Until(claims.ExpiresAt.Time) <= auth.TokenRenewalTime {
		if claims.RenewedCount >= auth.TokenMaxRenews {
			return nil, true, auth.InvalidTokenErr
		}

		// Increase re-newed count
		claims.RenewedCount++

		auth.SetTokenClaimsTimes(claims)
		newToken, err := s.tm.NewWithClaims(claims)
		if err != nil {
			return nil, true, auth.CheckTokenErr
		}

		tu.NewToken = &newToken
		tu.Expires = timestamp.New(claims.ExpiresAt.Time)

		return claims, true, nil
	}

	return claims, false, nil
}

func (s *Server) checkAndUpdateUserInfo(ctx context.Context, tu *TokenUpdate, currentUserInfo *userinfo.UserInfo) error {
	userInfo, err := s.ui.GetUserInfo(ctx, currentUserInfo.UserId, currentUserInfo.AccId)
	if err != nil {
		return err
	}

	// If the user is logged into a character, update user info and load permissions of user
	if !currentUserInfo.Equal(userInfo) {
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
			return auth.NoPermsErr
		}
		tu.Permissions = ps.GuardNames()

		if userInfo.SuperUser {
			tu.Permissions = append(tu.Permissions, common.SuperuserPermission)
		}

		attrs, err := s.p.FlattenRoleAttributes(userInfo.Job, userInfo.JobGrade)
		if err != nil {
			return auth.NoPermsErr
		}
		tu.Permissions = append(tu.Permissions, attrs...)
	}

	return nil
}

func (s *Server) getCharacter(ctx context.Context, charId int32) (*users.User, *jobs.JobProps, string, error) {
	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
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
		JobProps jobs.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, nil, "", err
	}

	return &dest.User, &dest.JobProps, dest.Group, nil
}
