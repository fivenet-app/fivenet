package citizenstore

import (
	context "context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/proto/resources/common/database"
	users "github.com/galexrt/fivenet/proto/resources/users"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	user         = table.Users.AS("user")
	userLicenses = table.UserLicenses
	licenses     = table.Licenses

	userProps = table.FivenetUserProps
	userAct   = table.FivenetUserActivity
)

var (
	FailedQueryErr      = status.Error(codes.Internal, "Failed to list/get citizen(s) date!")
	NotFoundOrNoPermErr = status.Error(codes.NotFound, "Citizen not found or no permission to access")
)

type Server struct {
	CitizenStoreServiceServer

	db *sql.DB
	p  perms.Permissions
	c  *mstlystcdata.Enricher
}

func NewServer(db *sql.DB, p perms.Permissions, c *mstlystcdata.Enricher) *Server {
	return &Server{
		db: db,
		p:  p,
		c:  c,
	}
}

func (s *Server) FindUsers(ctx context.Context, req *FindUsersRequest) (*FindUsersResponse, error) {
	userId, _, _ := auth.GetUserInfoFromContext(ctx)

	selectors := jet.ProjectionList{
		user.ID,
		user.Identifier,
		user.Firstname,
		user.Lastname,
		user.Job,
		user.JobGrade,
		user.Dateofbirth,
		user.Sex,
		user.Height,
		user.Visum,
		userProps.UserID,
	}

	condition := jet.Bool(true)
	// Field Permission Check
	if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "UserProps") {
		if req.PhoneNumber != "" {
			if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "PhoneNumber") {
				selectors = append(selectors, user.PhoneNumber)
				phoneNumber := strings.ReplaceAll(strings.ReplaceAll(req.PhoneNumber, "%", ""), " ", "") + "%"
				condition = condition.AND(user.PhoneNumber.LIKE(jet.String(phoneNumber)))
			}
		}
		if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "UserProps", "Wanted") {
			selectors = append(selectors, userProps.Wanted)
			if req.Wanted {
				condition = condition.AND(userProps.Wanted.IS_TRUE())
			}
		}
		if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "UserProps", "Job") {
			selectors = append(selectors, userProps.Job)
		}
	}

	req.SearchName = strings.ReplaceAll(req.SearchName, "%", "")
	if req.SearchName != "" {
		condition = condition.AND(jet.BoolExp(jet.Raw("MATCH(firstname,lastname) AGAINST ($search IN NATURAL LANGUAGE MODE)", jet.RawArgs{"$search": req.SearchName})))
	}

	// Get total count of values
	countStmt := user.
		SELECT(
			jet.COUNT(user.ID).AS("datacount.totalcount"),
		).
		FROM(
			user.
				LEFT_JOIN(userProps,
					userProps.UserID.EQ(user.ID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, FailedQueryErr
	}

	resp := &FindUsersResponse{
		Pagination: database.EmptyPaginationResponse(req.Pagination.Offset),
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := user.
		SELECT(
			selectors[0], selectors[1:]...,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname")).
		FROM(user.
			LEFT_JOIN(userProps,
				userProps.UserID.EQ(user.ID),
			),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(database.DefaultPageLimit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		return nil, FailedQueryErr
	}

	database.PaginationHelper(resp.Pagination,
		count.TotalCount,
		req.Pagination.Offset,
		len(resp.Users))

	for i := 0; i < len(resp.Users); i++ {
		if utils.InStringSlice(config.C.Game.PublicJobs, resp.Users[i].Job) {
			// Make sure user has permission to see that grade, otherwise "hide" the user's job
			if !s.p.Can(userId, CitizenStoreServicePermKey, "GetUser", resp.Users[i].Job, strconv.Itoa(int(resp.Users[i].JobGrade))) {
				resp.Users[i].JobGrade = -1
			}
		} else {
			resp.Users[i].Job = config.C.Game.UnemployedJob.Name
			resp.Users[i].JobGrade = config.C.Game.UnemployedJob.Grade
		}

		if resp.Users[i].Props != nil && resp.Users[i].Props.Job != "" {
			resp.Users[i].Job = resp.Users[i].Props.Job
			resp.Users[i].JobGrade = -1
		}

		s.c.EnrichJobInfo(resp.Users[i])
	}

	return resp, nil
}

func (s *Server) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	userId := auth.GetUserIDFromContext(ctx)

	selectors := jet.ProjectionList{
		user.ID,
		user.Identifier,
		user.Firstname,
		user.Lastname,
		user.Job,
		user.JobGrade,
		user.Dateofbirth,
		user.Sex,
		user.Height,
		user.Visum,
		userProps.UserID,
	}

	// Field Permission Check
	if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "UserProps") {
		// Field Permission Check
		if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "PhoneNumber") {
			selectors = append(selectors, user.PhoneNumber)
		}
		if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "UserProps", "Wanted") {
			selectors = append(selectors, userProps.Wanted)
		}
		if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "UserProps", "Job") {
			selectors = append(selectors, userProps.Job)
		}
	}

	resp := &GetUserResponse{
		User: &users.User{},
	}
	stmt := user.
		SELECT(
			selectors[0], selectors[1:]...,
		).
		FROM(
			user.
				LEFT_JOIN(userProps,
					userProps.UserID.EQ(user.ID),
				),
		).
		WHERE(user.ID.EQ(jet.Int32(req.UserId))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, resp.User); err != nil {
		return nil, FailedQueryErr
	}

	// Check if user can see licenses and fetch them separately
	if s.p.Can(userId, CitizenStoreServicePermKey, "FindUsers", "Licenses") {
		stmt := user.
			SELECT(
				userLicenses.Type.AS("license.type"),
				licenses.Label.AS("license.label"),
			).
			FROM(
				userLicenses.
					INNER_JOIN(user,
						userLicenses.Owner.EQ(user.Identifier),
					).
					LEFT_JOIN(licenses,
						licenses.Type.EQ(userLicenses.Type)),
			).
			WHERE(user.ID.EQ(jet.Int32(req.UserId))).
			LIMIT(15)

		if err := stmt.QueryContext(ctx, s.db, &resp.User.Licenses); err != nil {
			if !errors.Is(qrm.ErrNoRows, err) {
				return nil, FailedQueryErr
			}
		}
	}

	if resp.User != nil {
		if utils.InStringSlice(config.C.Game.PublicJobs, resp.User.Job) {
			// Make sure user has permission to see that grade
			if !s.p.Can(userId, CitizenStoreServicePermKey, "GetUser", resp.User.Job, strconv.Itoa(int(resp.User.JobGrade))) {
				return nil, NotFoundOrNoPermErr
			}
		} else {
			resp.User.Job = config.C.Game.UnemployedJob.Name
			resp.User.JobGrade = config.C.Game.UnemployedJob.Grade
		}

		if resp.User.Props != nil && resp.User.Props.Job != "" {
			resp.User.Job = resp.User.Props.Job
			resp.User.JobGrade = -1
		}

		s.c.EnrichJobInfo(resp.User)
	}

	return resp, nil
}

func (s *Server) GetUserActivity(ctx context.Context, req *GetUserActivityRequest) (*GetUserActivityResponse, error) {
	userId, _, _ := auth.GetUserInfoFromContext(ctx)

	resp := &GetUserActivityResponse{}
	// An user can never see their own activity on their own "profile"
	if userId == req.UserId {
		return resp, nil
	}

	uTarget := user.AS("target_user")
	uSource := user.AS("source_user")
	stmt := userAct.
		SELECT(
			userAct.AllColumns,
			uTarget.ID,
			uTarget.Identifier,
			uTarget.Job,
			uTarget.JobGrade,
			uTarget.Firstname,
			uTarget.Lastname,
			uSource.ID,
			uSource.Identifier,
			uSource.Job,
			uSource.JobGrade,
			uSource.Firstname,
			uSource.Lastname,
		).
		FROM(
			userAct.
				LEFT_JOIN(uTarget,
					uTarget.ID.EQ(userAct.TargetUserID),
				).
				LEFT_JOIN(uSource,
					uSource.ID.EQ(userAct.SourceUserID),
				),
		).
		WHERE(
			userAct.TargetUserID.EQ(jet.Int32(req.UserId)),
		).
		LIMIT(12)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, FailedQueryErr
		}
	}

	for i := 0; i < len(resp.Activity); i++ {
		s.c.EnrichJobInfo(resp.Activity[i].SourceUser)
		s.c.EnrichJobInfo(resp.Activity[i].TargetUser)
	}

	return resp, nil
}

func (s *Server) SetUserProps(ctx context.Context, req *SetUserPropsRequest) (*SetUserPropsResponse, error) {
	userId := auth.GetUserIDFromContext(ctx)

	// Field Permission Check
	if !s.p.Can(userId, CitizenStoreServicePermKey, "SetUserProps", "Wanted") {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to set user wanted status!")
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, FailedQueryErr
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	stmt := userProps.
		INSERT(
			userProps.AllColumns,
		).
		VALUES(
			req.Props.UserId,
			req.Props.Wanted,
		).
		ON_DUPLICATE_KEY_UPDATE(
			userProps.Wanted.SET(jet.Bool(req.Props.Wanted)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, FailedQueryErr
	}

	// Create user activity
	key := "UserProps.Wanted"
	newValue := strconv.FormatBool(req.Props.Wanted)
	oldValue := strconv.FormatBool(!req.Props.Wanted)
	if err := s.addUserAcitvity(ctx, tx,
		&model.FivenetUserActivity{
			SourceUserID: userId,
			TargetUserID: req.Props.UserId,
			Type:         int16(users.USER_ACTIVITY_TYPE_CHANGED),
			Key:          key,
			OldValue:     &oldValue,
			NewValue:     &newValue,
		}); err != nil {
		return nil, FailedQueryErr
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, FailedQueryErr
	}

	return &SetUserPropsResponse{}, nil
}

func (s *Server) addUserAcitvity(ctx context.Context, tx *sql.Tx, activity *model.FivenetUserActivity) error {
	stmt := userAct.
		INSERT(
			userAct.SourceUserID,
			userAct.TargetUserID,
			userAct.Type,
			userAct.Key,
			userAct.OldValue,
			userAct.NewValue,
		).
		MODEL(activity)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}
