package citizenstore

import (
	context "context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
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
	ErrFailedQuery          = status.Error(codes.Internal, "errors.CitizenStoreService.ErrFailedQuery")
	ErrJobGradeNoPermission = status.Error(codes.NotFound, "errors.CitizenStoreService.ErrJobGradeNoPermission")
)

var ZeroTrafficInfractionPoints uint64 = 0

type Server struct {
	CitizenStoreServiceServer

	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.Enricher
	aud      audit.IAuditer

	publicJobs         []string
	unemployedJob      string
	unemployedJobGrade int32
}

func NewServer(db *sql.DB, p perms.Permissions, enricher *mstlystcdata.Enricher, aud audit.IAuditer, cfg *config.Config) *Server {
	return &Server{
		db:       db,
		p:        p,
		enricher: enricher,
		aud:      aud,

		publicJobs:         cfg.Game.PublicJobs,
		unemployedJob:      cfg.Game.UnemployedJob.Name,
		unemployedJobGrade: cfg.Game.UnemployedJob.Grade,
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
		return nil, ErrFailedQuery
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	for _, field := range fields {
		switch field {
		case "PhoneNumber":
			selectors = append(selectors, tUser.PhoneNumber)
			if req.PhoneNumber != nil && *req.PhoneNumber != "" {
				phoneNumber := strings.ReplaceAll(strings.ReplaceAll(*req.PhoneNumber, "%", ""), " ", "") + "%"
				condition = condition.AND(tUser.PhoneNumber.LIKE(jet.String(phoneNumber)))
			}
		case "UserProps.Wanted":
			selectors = append(selectors, tUserProps.Wanted)
			if req.Wanted != nil && *req.Wanted {
				condition = condition.AND(tUserProps.Wanted.IS_TRUE())
			}
		case "UserProps.Job":
			selectors = append(selectors, tUserProps.Job)
		case "UserProps.TrafficInfractionPoints":
			selectors = append(selectors, tUserProps.TrafficInfractionPoints)
			if req.TrafficPoints != nil && *req.TrafficPoints > 0 {
				condition = condition.AND(tUserProps.TrafficInfractionPoints.GT_EQ(jet.Uint64(*req.TrafficPoints)))
			}
		case "UserProps.OpenFines":
			selectors = append(selectors, tUserProps.OpenFines)
			if req.OpenFines != nil && *req.OpenFines > 0 {
				condition = condition.AND(tUserProps.OpenFines.GT_EQ(jet.Uint64(*req.OpenFines)))
			}
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

	if req.Dateofbirth != nil && *req.Dateofbirth != "" {
		condition = condition.AND(tUser.Dateofbirth.LIKE(jet.String(strings.ReplaceAll(*req.Dateofbirth, "%", " ") + "%")))
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
		return nil, ErrFailedQuery
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
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext")).
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
		return nil, ErrFailedQuery
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Users))

	jobInfoFn := s.enricher.EnrichJobInfoFunc(userInfo)
	for i := 0; i < len(resp.Users); i++ {
		jobInfoFn(resp.Users[i])

		if resp.Users[i].Job == s.unemployedJob {
			// Only set the user's job from the user props when it isn't a public job
			if resp.Users[i].Props != nil && resp.Users[i].Props.JobName != nil {
				resp.Users[i].Job = *resp.Users[i].Props.JobName
				if resp.Users[i].Props.JobGradeNumber != nil {
					resp.Users[i].JobGrade = *resp.Users[i].Props.JobGradeNumber
				} else {
					resp.Users[i].JobGrade = 0
				}
			}
		}
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
		State:        int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

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
		return nil, ErrFailedQuery
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
			selectors = append(selectors, tUserProps.JobGrade)
		case "UserProps.TrafficInfractionPoints":
			selectors = append(selectors, tUserProps.TrafficInfractionPoints)
		case "UserProps.OpenFines":
			selectors = append(selectors, tUserProps.OpenFines)
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
		return nil, ErrFailedQuery
	}

	if resp.User == nil || resp.User.UserId <= 0 {
		return nil, ErrJobGradeNoPermission
	}

	if utils.InSlice(s.publicJobs, resp.User.Job) {
		// Make sure user has permission to see that grade
		jobGradesAttr, err := s.p.Attr(userInfo, CitizenStoreServicePerm, CitizenStoreServiceGetUserPerm, CitizenStoreServiceGetUserJobsPermField)
		if err != nil {
			return nil, ErrFailedQuery
		}
		var jobGrades perms.JobGradeList
		if jobGradesAttr != nil {
			jobGrades = jobGradesAttr.(map[string]int32)
		}

		if len(jobGrades) == 0 && !userInfo.SuperUser {
			return nil, ErrJobGradeNoPermission
		}

		// Make sure user has permission to see that grade, otherwise "hide" the user's job
		grade, ok := jobGrades[resp.User.Job]
		if !ok || resp.User.JobGrade > grade {
			// Skip for superuser
			if !userInfo.SuperUser {
				return nil, ErrJobGradeNoPermission
			}
		}
	} else {
		resp.User.Job = s.unemployedJob
		resp.User.JobGrade = s.unemployedJobGrade
	}

	if resp.User.Props != nil {
		if resp.User.Props.JobName != nil {
			resp.User.Job = *resp.User.Props.JobName
		}

		if resp.User.Props.JobGradeNumber != nil {
			resp.User.JobGrade = *resp.User.Props.JobGradeNumber
		}
	}

	s.enricher.EnrichJobInfoSafe(userInfo, resp.User)

	// Check if user can see licenses and fetch them
	if utils.InSlice(fields, "Licenses") {
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
				return nil, ErrFailedQuery
			}
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) ListUserActivity(ctx context.Context, req *ListUserActivityRequest) (*ListUserActivityResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	pag, limit := req.Pagination.GetResponseWithPageSize(15)
	resp := &ListUserActivityResponse{
		Pagination: pag,
		Activity:   []*users.UserActivity{},
	}
	// User can't see their own activities, unless they have "Own" perm attribute, or are a superuser
	fieldsAttr, err := s.p.Attr(userInfo, CitizenStoreServicePerm, CitizenStoreServiceListUserActivityPerm, CitizenStoreServiceListUserActivityFieldsPermField)
	if err != nil {
		return nil, ErrFailedQuery
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	if userInfo.UserId == req.UserId {
		// If isn't superuser or doesn't have 'Own' activity feed access
		if !userInfo.SuperUser || !utils.InSlice(fields, "Own") {
			return resp, nil
		}
	}

	condition := tUserActivity.TargetUserID.EQ(jet.Int32(req.UserId))

	// Get total count of values
	countStmt := tUserActivity.
		SELECT(
			jet.COUNT(tUserActivity.ID).AS("datacount.totalcount"),
		).
		FROM(tUserActivity).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, ErrFailedQuery
	}

	if count.TotalCount <= 0 {
		return resp, nil
	}

	tUTarget := tUser.AS("target_user")
	tUSource := tUser.AS("source_user")
	stmt := tUserActivity.
		SELECT(
			tUserActivity.ID,
			tUserActivity.CreatedAt,
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Key,
			tUserActivity.OldValue,
			tUserActivity.NewValue,
			tUserActivity.Reason,
			tUTarget.ID,
			tUTarget.Identifier,
			tUTarget.Job,
			tUTarget.JobGrade,
			tUTarget.Firstname,
			tUTarget.Lastname,
			tUSource.ID,
			tUSource.Identifier,
			tUSource.Job,
			tUSource.JobGrade,
			tUSource.Firstname,
			tUSource.Lastname,
		).
		FROM(
			tUserActivity.
				LEFT_JOIN(tUTarget,
					tUTarget.ID.EQ(tUserActivity.TargetUserID),
				).
				LEFT_JOIN(tUSource,
					tUSource.ID.EQ(tUserActivity.SourceUserID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(
			tUserActivity.CreatedAt.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrFailedQuery
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoFunc(userInfo)
	for i := 0; i < len(resp.Activity); i++ {
		if resp.Activity[i].SourceUser != nil {
			jobInfoFn(resp.Activity[i].SourceUser)
		}
		if resp.Activity[i].TargetUser != nil {
			jobInfoFn(resp.Activity[i].TargetUser)
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Activity))

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
		State:        int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "Must give a reason!")
	}

	// Get current user props to be able to compare
	props, err := s.getUserProps(ctx, req.Props.UserId)
	if err != nil {
		return nil, ErrFailedQuery
	}
	if props.Wanted == nil {
		wanted := false
		props.Wanted = &wanted
	}
	if props.JobName == nil {
		props.JobName = &s.unemployedJob
	}
	if props.JobGradeNumber == nil {
		props.JobGradeNumber = &s.unemployedJobGrade
	}
	if props.TrafficInfractionPoints == nil {
		props.TrafficInfractionPoints = &ZeroTrafficInfractionPoints
	}

	props.Job, props.JobGrade = s.enricher.GetJobGrade(*props.JobName, *props.JobGradeNumber)
	// Make sure a job is set
	if props.Job == nil {
		props.Job, props.JobGrade = s.enricher.GetJobGrade(s.unemployedJob, s.unemployedJobGrade)
	}

	resp := &SetUserPropsResponse{
		Props: &users.UserProps{},
	}

	// Field Permission Check
	fieldsAttr, err := s.p.Attr(userInfo, CitizenStoreServicePerm, CitizenStoreServiceSetUserPropsPerm, CitizenStoreServiceSetUserPropsFieldsPermField)
	if err != nil {
		return nil, ErrFailedQuery
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	updateSets := []jet.ColumnAssigment{}
	if req.Props.Wanted != nil {
		if !utils.InSlice(fields, "Wanted") {
			return nil, status.Error(codes.PermissionDenied, "You are not allowed to set a user wanted status!")
		}

		updateSets = append(updateSets, tUserProps.Wanted.SET(jet.Bool(*req.Props.Wanted)))
	} else {
		req.Props.Wanted = props.Wanted
	}
	if req.Props.JobName != nil {
		if !utils.InSlice(fields, "Job") {
			return nil, status.Error(codes.PermissionDenied, "You are not allowed to set a user job!")
		}

		if utils.InSlice(s.publicJobs, *req.Props.JobName) {
			return nil, status.Error(codes.InvalidArgument, "You can't set a state job!")
		}

		req.Props.Job, req.Props.JobGrade = s.enricher.GetJobGrade(*req.Props.JobName, *req.Props.JobGradeNumber)
		if req.Props.Job == nil || req.Props.JobGrade == nil {
			return nil, status.Error(codes.PermissionDenied, "Invalid job or job rank set!")
		}

		updateSets = append(updateSets, tUserProps.Job.SET(jet.String(*req.Props.JobName)))
		updateSets = append(updateSets, tUserProps.JobGrade.SET(jet.Int32(*req.Props.JobGradeNumber)))
	} else {
		req.Props.JobName = props.JobName
		req.Props.Job = props.Job
		req.Props.JobGradeNumber = props.JobGradeNumber
		req.Props.JobGrade = props.JobGrade
	}
	if req.Props.TrafficInfractionPoints != nil {
		// Only update when it has actually changed
		if !utils.InSlice(fields, "TrafficInfractionPoints") {
			return nil, status.Error(codes.PermissionDenied, "You are not allowed to set a user's traffic infraction points!")
		}

		updateSets = append(updateSets, tUserProps.TrafficInfractionPoints.SET(jet.Uint64(*req.Props.TrafficInfractionPoints)))
	} else {
		req.Props.TrafficInfractionPoints = props.TrafficInfractionPoints
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	stmt := tUserProps.
		INSERT(
			tUserProps.UserID,
			tUserProps.Wanted,
			tUserProps.Job,
			tUserProps.JobGrade,
			tUserProps.TrafficInfractionPoints,
		).
		VALUES(
			req.Props.UserId,
			req.Props.Wanted,
			req.Props.JobName,
			req.Props.JobGradeNumber,
			req.Props.TrafficInfractionPoints,
		).
		ON_DUPLICATE_KEY_UPDATE(
			updateSets...,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, ErrFailedQuery
	}

	// Create user activity entries
	if *req.Props.Wanted != *props.Wanted {
		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Props.UserId, users.UserActivityType_USER_ACTIVITY_TYPE_CHANGED, "UserProps.Wanted",
			strconv.FormatBool(*props.Wanted), strconv.FormatBool(*req.Props.Wanted), req.Reason); err != nil {
			return nil, ErrFailedQuery
		}
	}
	if *req.Props.JobName != *props.JobName || *req.Props.JobGradeNumber != *props.JobGradeNumber {
		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Props.UserId, users.UserActivityType_USER_ACTIVITY_TYPE_CHANGED, "UserProps.Job",
			fmt.Sprintf("%s|%s", props.Job.Label, props.JobGrade.Label), fmt.Sprintf("%s|%s", req.Props.Job.Label, req.Props.JobGrade.Label), req.Reason); err != nil {
			return nil, ErrFailedQuery
		}
	}
	if *req.Props.TrafficInfractionPoints != *props.TrafficInfractionPoints {
		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Props.UserId, users.UserActivityType_USER_ACTIVITY_TYPE_CHANGED, "UserProps.TrafficInfractionPoints",
			strconv.Itoa(int(*props.TrafficInfractionPoints)), strconv.Itoa(int(*req.Props.TrafficInfractionPoints)), req.Reason); err != nil {
			return nil, ErrFailedQuery
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, ErrFailedQuery
	}

	// Set Job info
	resp.Props = req.Props
	resp.Props.Job, resp.Props.JobGrade = s.enricher.GetJobGrade(*resp.Props.JobName, *resp.Props.JobGradeNumber)

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return resp, nil
}

func (s *Server) getUserProps(ctx context.Context, userId int32) (*users.UserProps, error) {
	tUserProps := tUserProps.AS("userprops")
	stmt := tUserProps.
		SELECT(
			tUserProps.UserID,
			tUserProps.Wanted,
			tUserProps.Job,
			tUserProps.JobGrade,
			tUserProps.TrafficInfractionPoints,
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

func (s *Server) addUserActivity(ctx context.Context, tx *sql.Tx, userId int32, targetUserId int32, activityType users.UserActivityType, key string, oldValue string, newValue string, reason string) error {
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
			SourceUserID: &userId,
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
