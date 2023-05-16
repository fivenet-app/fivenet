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
	tUser         = table.Users.AS("user")
	tUserLicenses = table.UserLicenses
	tLicenses     = table.Licenses

	tUserProps    = table.FivenetUserProps
	tUserActivity = table.FivenetUserActivity
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
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	selectors := jet.ProjectionList{
		tUser.ID,
		tUser.Identifier,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.Sex,
		tUser.Height,
		tUser.Visum,
		tUserProps.UserID,
	}

	condition := jet.Bool(true)
	// Field Permission Check
	fieldsAttr, err := s.p.Attr(userInfo, CitizenStoreServicePerm, CitizenStoreServiceListCitizensPerm, CitizenStoreServiceListCitizensFieldsPermField)
	if err != nil {
		return nil, FailedQueryErr
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	for _, field := range fields {
		switch field {
		case "PhoneNumber":
			selectors = append(selectors, tUser.PhoneNumber)
			if req.PhoneNumber != "" {
				phoneNumber := strings.ReplaceAll(strings.ReplaceAll(req.PhoneNumber, "%", ""), " ", "") + "%"
				condition = condition.AND(tUser.PhoneNumber.LIKE(jet.String(phoneNumber)))
			}
		case "UserProps.Wanted":
			selectors = append(selectors, tUserProps.Wanted)
			if req.Wanted {
				condition = condition.AND(tUserProps.Wanted.IS_TRUE())
			}
		case "UserProps.Job":
			selectors = append(selectors, tUserProps.Job)
		}
	}

	req.SearchName = strings.TrimSpace(req.SearchName)
	req.SearchName = strings.ReplaceAll(req.SearchName, "%", "")
	req.SearchName = strings.ReplaceAll(req.SearchName, " ", "%")
	if req.SearchName != "" {
		req.SearchName = "%" + req.SearchName + "%"
		condition = condition.AND(
			jet.CONCAT(tUser.Firstname, jet.String(" "), tUser.Lastname).
				LIKE(jet.String(req.SearchName)),
		)
	}

	// Get total count of values
	countStmt := tUser.
		SELECT(
			jet.COUNT(tUser.ID).AS("datacount.totalcount"),
		).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
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

	stmt := tUser.
		SELECT(
			selectors[0], selectors[1:]...,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname")).
		FROM(tUser.
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(
			tUser.Firstname.ASC(),
			tUser.Lastname.ASC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		return nil, FailedQueryErr
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Users))

	jobGradesAttr, err := s.p.Attr(userInfo, CitizenStoreServicePerm, CitizenStoreServiceGetUserPerm, CitizenStoreServiceGetUserJobsPermField)
	if err != nil {
		return nil, FailedQueryErr
	}
	var jobGrades perms.JobGradeList
	if jobGradesAttr != nil {
		jobGrades = jobGradesAttr.(map[string]int32)
	}

	for i := 0; i < len(resp.Users); i++ {
		if utils.InStringSlice(config.C.Game.PublicJobs, resp.Users[i].Job) {
			// Make sure user has permission to see that grade, otherwise "hide" the user's job
			grade, ok := jobGrades[resp.Users[i].Job]
			if !ok || grade > resp.Users[i].JobGrade {
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
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service:      CitizenStoreService_ServiceDesc.ServiceName,
		Method:       "GetUser",
		UserID:       userInfo.UserId,
		UserJob:      userInfo.Job,
		TargetUserID: &req.UserId,
		State:        int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	selectors := jet.ProjectionList{
		tUser.ID,
		tUser.Identifier,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.Sex,
		tUser.Height,
		tUser.Visum,
		tUserProps.UserID,
	}

	// Field Permission Check
	fieldsAttr, err := s.p.Attr(userInfo, CitizenStoreServicePerm, CitizenStoreServiceListCitizensPerm, CitizenStoreServiceListCitizensFieldsPermField)
	if err != nil {
		return nil, FailedQueryErr
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	for _, field := range fields {
		switch field {
		case "PhoneNumber":
			selectors = append(selectors, tUser.PhoneNumber)
		case "UserProps.Wanted":
			selectors = append(selectors, tUserProps.Wanted)
		case "UserProps.Job":
			selectors = append(selectors, tUserProps.Job)
		}
	}

	resp := &GetUserResponse{
		User: &users.User{},
	}
	stmt := tUser.
		SELECT(
			selectors[0], selectors[1:]...,
		).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				),
		).
		WHERE(tUser.ID.EQ(jet.Int32(req.UserId))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, resp.User); err != nil {
		return nil, FailedQueryErr
	}

	if resp.User != nil {
		if utils.InStringSlice(config.C.Game.PublicJobs, resp.User.Job) {
			// Make sure user has permission to see that grade
			jobGradesAttr, err := s.p.Attr(userInfo, CitizenStoreServicePerm, CitizenStoreServiceListCitizensPerm, CitizenStoreServiceGetUserJobsPermField)
			if err != nil {
				return nil, FailedQueryErr
			}
			var jobGrades perms.JobGradeList
			if jobGradesAttr != nil {
				jobGrades = jobGradesAttr.(map[string]int32)
			}

			if len(jobGrades) == 0 {
				return nil, JobGradeNoPermissionErr
			}

			grade, ok := jobGrades[resp.User.Job]
			if ok && grade > resp.User.JobGrade {
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
	if utils.InStringSlice(fields, "Licenses") {
		stmt := tUser.
			SELECT(
				tUserLicenses.Type.AS("license.type"),
				tLicenses.Label.AS("license.label"),
			).
			FROM(
				tUserLicenses.
					INNER_JOIN(tUser,
						tUserLicenses.Owner.EQ(tUser.Identifier),
					).
					LEFT_JOIN(tLicenses,
						tLicenses.Type.EQ(tUserLicenses.Type)),
			).
			WHERE(tUser.ID.EQ(jet.Int32(req.UserId))).
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
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &ListUserActivityResponse{}
	// An user can never see their own activity on their "profile"
	if userInfo.UserId == req.UserId {
		return resp, nil
	}

	uTarget := tUser.AS("target_user")
	uSource := tUser.AS("source_user")
	stmt := tUserActivity.
		SELECT(
			tUserActivity.AllColumns,
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
			tUserActivity.
				LEFT_JOIN(uTarget,
					uTarget.ID.EQ(tUserActivity.TargetUserID),
				).
				LEFT_JOIN(uSource,
					uSource.ID.EQ(tUserActivity.SourceUserID),
				),
		).
		WHERE(
			tUserActivity.TargetUserID.EQ(jet.Int32(req.UserId)),
		).
		ORDER_BY(
			tUserActivity.CreatedAt.DESC(),
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
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service:      CitizenStoreService_ServiceDesc.ServiceName,
		Method:       "SetUserProps",
		UserID:       userInfo.UserId,
		UserJob:      userInfo.Job,
		TargetUserID: &req.Props.UserId,
		State:        int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	resp := &SetUserPropsResponse{
		Props: &users.UserProps{},
	}

	if req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "Must give a reason!")
	}

	// Get current user props to be able to compare
	props, err := s.getUserProps(ctx, req.Props.UserId)
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
	fieldsAttr, err := s.p.Attr(userInfo, CitizenStoreServicePerm, CitizenStoreServiceSetUserPropsPerm, CitizenStoreServiceSetUserPropsFieldsPermField)
	if err != nil {
		return nil, FailedQueryErr
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	if req.Props.Wanted != nil {
		if !utils.InStringSlice(fields, "Wanted") {
			return nil, status.Error(codes.PermissionDenied, "You are not allowed to set a user wanted status!")
		}

		updateSets = append(updateSets, tUserProps.Wanted.SET(jet.Bool(*req.Props.Wanted)))
	} else {
		var current bool
		if props.Wanted != nil {
			current = !*props.Wanted
		}
		req.Props.Wanted = &current
	}
	if req.Props.JobName != nil {
		if !utils.InStringSlice(fields, "Job") {
			return nil, status.Error(codes.PermissionDenied, "You are not allowed to set a user job!")
		}

		if utils.InStringSlice(config.C.Game.PublicJobs, *req.Props.JobName) {
			return nil, status.Error(codes.InvalidArgument, "You can't set a state job!")
		}

		resp.Props.Job = s.c.GetJobByName(*req.Props.JobName)
		if resp.Props.Job == nil {
			return nil, status.Error(codes.PermissionDenied, "Invalid job set!")
		}

		updateSets = append(updateSets, tUserProps.Job.SET(jet.String(*req.Props.JobName)))
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

	stmt := tUserProps.
		INSERT(
			tUserProps.UserID,
			tUserProps.Wanted,
			tUserProps.Job,
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
	if *req.Props.Wanted != *props.Wanted {
		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Props.UserId, users.USER_ACTIVITY_TYPE_CHANGED, "UserProps.Wanted", strconv.FormatBool(*props.Wanted), strconv.FormatBool(*req.Props.Wanted), req.Reason); err != nil {
			return nil, FailedQueryErr
		}
	}
	if *req.Props.JobName != *props.JobName {
		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Props.UserId, users.USER_ACTIVITY_TYPE_CHANGED, "UserProps.Job", *props.JobName, *req.Props.JobName, req.Reason); err != nil {
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
	tUserProps := tUserProps.AS("userprops")
	stmt := tUserProps.
		SELECT(
			tUserProps.UserID,
			tUserProps.Wanted,
			tUserProps.Job,
		).
		FROM(tUserProps).
		WHERE(
			tUserProps.UserID.EQ(jet.Int32(userId)),
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

func (s *Server) addUserActivity(ctx context.Context, tx *sql.Tx, userId int32, targetUserId int32, activityType users.USER_ACTIVITY_TYPE, key string, oldValue string, newValue string, reason string) error {
	stmt := tUserActivity.
		INSERT(
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Key,
			tUserActivity.OldValue,
			tUserActivity.NewValue,
			tUserActivity.Reason,
		).
		MODEL(&model.FivenetUserActivity{
			SourceUserID: userId,
			TargetUserID: targetUserId,
			Type:         int16(activityType),
			Key:          key,
			OldValue:     &oldValue,
			NewValue:     &newValue,
			Reason:       &reason,
		})

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
