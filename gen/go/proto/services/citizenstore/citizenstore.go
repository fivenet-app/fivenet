package citizenstore

import (
	context "context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
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
	FailedQueryErr          = status.Error(codes.Internal, "Failed to list/get citizen(s) data!")
	JobGradeNoPermissionErr = status.Error(codes.NotFound, "No permission to access this citizen (based on the citizen's job)")
)

type Server struct {
	CitizenStoreServiceServer

	db *sql.DB
	p  perms.Permissions
	c  *mstlystcdata.Enricher
	a  audit.IAuditer
}

func NewServer(db *sql.DB, p perms.Permissions, c *mstlystcdata.Enricher, aud audit.IAuditer) *Server {
	return &Server{
		db: db,
		p:  p,
		c:  c,
		a:  aud,
	}
}

func (s *Server) ListCitizens(ctx context.Context, req *ListCitizensRequest) (*ListCitizensResponse, error) {
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
	if s.p.Can(userId, CitizenStoreServicePermKey, "ListCitizens", "PhoneNumber") {
		selectors = append(selectors, user.PhoneNumber)
		if req.PhoneNumber != "" {
			phoneNumber := strings.ReplaceAll(strings.ReplaceAll(req.PhoneNumber, "%", ""), " ", "") + "%"
			condition = condition.AND(user.PhoneNumber.LIKE(jet.String(phoneNumber)))
		}
	}
	if s.p.Can(userId, CitizenStoreServicePermKey, "ListCitizens", "UserProps", "Wanted") {
		selectors = append(selectors, userProps.Wanted)
		if req.Wanted {
			condition = condition.AND(userProps.Wanted.IS_TRUE())
		}
	}
	if s.p.Can(userId, CitizenStoreServicePermKey, "ListCitizens", "UserProps", "Job") {
		selectors = append(selectors, userProps.Job.AS("jobname"))
	}

	req.SearchName = strings.TrimSpace(req.SearchName)
	req.SearchName = strings.ReplaceAll(req.SearchName, "%", "")
	req.SearchName = strings.ReplaceAll(req.SearchName, " ", "%")
	if req.SearchName != "" {
		req.SearchName = "%" + req.SearchName + "%"
		condition = condition.AND(
			jet.CONCAT(user.Firstname, jet.String(" "), user.Lastname).
				LIKE(jet.String(req.SearchName)),
		)
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

	pag, limit := req.Pagination.GetResponse()
	resp := &ListCitizensResponse{
		Pagination: pag,
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
		ORDER_BY(
			user.Firstname.ASC(),
			user.Lastname.ASC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		return nil, FailedQueryErr
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Users))

	for i := 0; i < len(resp.Users); i++ {
		if utils.InStringSlice(config.C.Game.PublicJobs, resp.Users[i].Job) {
			// Make sure user has permission to see that grade, otherwise "hide" the user's job
			if !s.p.Can(userId, CitizenStoreServicePermKey, "GetUser", resp.Users[i].Job, strconv.Itoa(int(resp.Users[i].JobGrade))) {
				resp.Users[i].JobGrade = 0
			}
		} else {
			resp.Users[i].Job = config.C.Game.UnemployedJob.Name
			resp.Users[i].JobGrade = config.C.Game.UnemployedJob.Grade
		}

		if resp.Users[i].Props != nil && resp.Users[i].Props.JobName != nil {
			resp.Users[i].Job = *resp.Users[i].Props.JobName
			resp.Users[i].JobGrade = 0
		}

		s.c.EnrichJobInfo(resp.Users[i])
	}

	return resp, nil
}

func (s *Server) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service:      CitizenStoreService_ServiceDesc.ServiceName,
		Method:       "GetUser",
		UserID:       userId,
		UserJob:      job,
		TargetUserID: &req.UserId,
		State:        int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

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
	if s.p.Can(userId, CitizenStoreServicePermKey, "ListCitizens", "UserProps") {
		// Field Permission Check
		if s.p.Can(userId, CitizenStoreServicePermKey, "ListCitizens", "PhoneNumber") {
			selectors = append(selectors, user.PhoneNumber)
		}
		if s.p.Can(userId, CitizenStoreServicePermKey, "ListCitizens", "UserProps", "Wanted") {
			selectors = append(selectors, userProps.Wanted)
		}
		if s.p.Can(userId, CitizenStoreServicePermKey, "ListCitizens", "UserProps", "Job") {
			selectors = append(selectors, userProps.Job.AS("userprops.job_name"))
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

	if resp.User != nil {
		if utils.InStringSlice(config.C.Game.PublicJobs, resp.User.Job) {
			// Make sure user has permission to see that grade
			grades, err := s.p.GetSuffixOfPermissionsByPrefixOfUser(userId, CitizenStoreServicePermKey+".GetUser."+resp.User.Job)
			if err != nil {
				return nil, FailedQueryErr
			}

			if len(grades) == 0 {
				return nil, JobGradeNoPermissionErr
			}

			allowed := false

			for _, grade := range grades {
				i, err := strconv.Atoi(grade)
				if err != nil {
					return nil, FailedQueryErr
				}

				if resp.User.JobGrade <= int32(i) {
					allowed = true
				}
			}

			if !allowed {
				return nil, JobGradeNoPermissionErr
			}
		} else {
			resp.User.Job = config.C.Game.UnemployedJob.Name
			resp.User.JobGrade = config.C.Game.UnemployedJob.Grade
		}

		if resp.User.Props != nil && resp.User.Props.Job != nil {
			resp.User.Job = *resp.User.Props.JobName
			resp.User.JobGrade = 0
		}

		s.c.EnrichJobInfo(resp.User)
	}

	// Check if user can see licenses and fetch them
	if s.p.Can(userId, CitizenStoreServicePermKey, "ListCitizens", "Licenses") {
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

	auditEntry.State = int16(rector.EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) ListUserActivity(ctx context.Context, req *ListUserActivityRequest) (*ListUserActivityResponse, error) {
	userId, _, _ := auth.GetUserInfoFromContext(ctx)

	resp := &ListUserActivityResponse{}
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
		ORDER_BY(
			userAct.CreatedAt.DESC(),
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
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service:      CitizenStoreService_ServiceDesc.ServiceName,
		Method:       "SetUserProps",
		UserID:       userId,
		UserJob:      job,
		TargetUserID: &req.Props.UserId,
		State:        int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	resp := &SetUserPropsResponse{
		Props: &users.UserProps{},
	}

	// Use getUserProps
	props, err := s.getUserProps(ctx, userId)
	if err != nil {
		return nil, FailedQueryErr
	}
	if props.Wanted == nil {
		wanted := false
		props.Wanted = &wanted
	}
	if props.JobName == nil {
		props.JobName = &config.C.Game.UnemployedJob.Name
	}

	updateSets := []jet.ColumnAssigment{}
	// Field Permission Check
	if req.Props.Wanted != nil {
		if !s.p.Can(userId, CitizenStoreServicePermKey, "SetUserProps", "Wanted") {
			return nil, status.Error(codes.PermissionDenied, "You are not allowed to set a user wanted status!")
		}

		updateSets = append(updateSets, userProps.Wanted.SET(jet.Bool(*req.Props.Wanted)))
	} else {
		req.Props.Wanted = props.Wanted
	}
	if req.Props.JobName != nil {
		if !s.p.Can(userId, CitizenStoreServicePermKey, "SetUserProps", "Job") {
			return nil, status.Error(codes.PermissionDenied, "You are not allowed to set a user job!")
		}

		if utils.InStringSlice(config.C.Game.PublicJobs, *req.Props.JobName) {
			return nil, status.Error(codes.InvalidArgument, "You can't set a state job!")
		}

		resp.Props.Job = s.c.GetJobByName(*req.Props.JobName)
		if resp.Props.Job == nil {
			return nil, status.Error(codes.PermissionDenied, "Invalid job set!")
		}

		updateSets = append(updateSets, userProps.Job.SET(jet.String(*req.Props.JobName)))
	} else {
		req.Props.JobName = props.JobName
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
			userProps.UserID,
			userProps.Wanted,
			userProps.Job,
		).
		VALUES(
			req.Props.UserId,
			req.Props.Wanted,
			req.Props.JobName,
		).
		ON_DUPLICATE_KEY_UPDATE(
			updateSets...,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, FailedQueryErr
	}

	// Create user activity
	if req.Props.Wanted != props.Wanted {
		if err := s.addUserAcitvity(ctx, tx,
			userId, req.Props.UserId, users.USER_ACTIVITY_TYPE_CHANGED, "UserProps.Wanted", strconv.FormatBool(!*props.Wanted), strconv.FormatBool(*req.Props.Wanted)); err != nil {
			return nil, FailedQueryErr
		}
	}
	if req.Props.JobName != props.JobName {
		if err := s.addUserAcitvity(ctx, tx,
			userId, req.Props.UserId, users.USER_ACTIVITY_TYPE_CHANGED, "UserProps.Job", *props.JobName, *req.Props.JobName); err != nil {
			return nil, FailedQueryErr
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, FailedQueryErr
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return resp, nil
}

func (s *Server) getUserProps(ctx context.Context, userId int32) (*users.UserProps, error) {
	stmt := userProps.
		SELECT(
			userProps.UserID,
			userProps.Wanted,
			userProps.Job,
		).
		FROM(userProps).
		WHERE(
			userProps.UserID.EQ(jet.Int32(userId)),
		).
		LIMIT(1)

	var dest users.UserProps
	dest.UserId = userId
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) addUserAcitvity(ctx context.Context, tx *sql.Tx, userId int32, targetUserId int32, activityType users.USER_ACTIVITY_TYPE, key string, oldValue string, newValue string) error {
	stmt := userAct.
		INSERT(
			userAct.SourceUserID,
			userAct.TargetUserID,
			userAct.Type,
			userAct.Key,
			userAct.OldValue,
			userAct.NewValue,
		).
		MODEL(&model.FivenetUserActivity{
			SourceUserID: userId,
			TargetUserID: targetUserId,
			Type:         int16(activityType),
			Key:          key,
			OldValue:     &oldValue,
			NewValue:     &newValue,
		})

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}
