package citizenstore

import (
	context "context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/filestore"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	errorscitizenstore "github.com/galexrt/fivenet/gen/go/proto/services/citizenstore/errors"
	permscitizenstore "github.com/galexrt/fivenet/gen/go/proto/services/citizenstore/perms"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/config/appconfig"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/storage"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

var (
	tUser         = table.Users.AS("user")
	tUserLicenses = table.UserLicenses
	tLicenses     = table.Licenses

	tUserProps    = table.FivenetUserProps
	tUserActivity = table.FivenetUserActivity
)

var ZeroTrafficInfractionPoints uint32 = 0

type Server struct {
	CitizenStoreServiceServer

	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	st       storage.IStorage
	appCfg   appconfig.IConfig

	customDB config.CustomDB
}

type Params struct {
	fx.In

	DB        *sql.DB
	P         perms.Permissions
	Enricher  *mstlystcdata.UserAwareEnricher
	Aud       audit.IAuditer
	Config    *config.Config
	Storage   storage.IStorage
	AppConfig appconfig.IConfig
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		p:        p.P,
		enricher: p.Enricher,
		aud:      p.Aud,
		st:       p.Storage,
		appCfg:   p.AppConfig,

		customDB: p.Config.Database.Custom,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterCitizenStoreServiceServer(srv, s)
}

func (s *Server) ListCitizens(ctx context.Context, req *ListCitizensRequest) (*ListCitizensResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	selectors := dbutils.Columns{
		tUser.Identifier,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.Sex,
		tUser.Height,
		tUserProps.UserID,
		s.customDB.Columns.User.GetVisum(tUser.Alias()),
	}

	condition := s.customDB.Conditions.User.GetFilter(tUser.Alias())
	// Field Permission Check
	fieldsAttr, err := s.p.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
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
			selectors = append(selectors, tUserProps.Job, tUserProps.JobGrade)
		case "UserProps.TrafficInfractionPoints":
			selectors = append(selectors, tUserProps.TrafficInfractionPoints)
			if req.TrafficInfractionPoints != nil && *req.TrafficInfractionPoints > 0 {
				condition = condition.AND(tUserProps.TrafficInfractionPoints.GT_EQ(jet.Uint32(*req.TrafficInfractionPoints)))
			}
		case "UserProps.OpenFines":
			selectors = append(selectors, tUserProps.OpenFines)
			if req.OpenFines != nil && *req.OpenFines > 0 {
				condition = condition.AND(tUserProps.OpenFines.GT_EQ(jet.Uint64(*req.OpenFines)))
			}
		case "UserProps.BloodType":
			selectors = append(selectors, tUserProps.BloodType)
		case "UserProps.MugShot":
			selectors = append(selectors, tUserProps.MugShot)
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
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
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
			tUser.ID,
			selectors.Get()...,
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
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Users))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Users); i++ {
		if resp.Users[i].Props != nil && resp.Users[i].Props.JobName != nil {
			resp.Users[i].Job = *resp.Users[i].Props.JobName
			if resp.Users[i].Props.JobGradeNumber != nil {
				resp.Users[i].JobGrade = *resp.Users[i].Props.JobGradeNumber
			} else {
				resp.Users[i].JobGrade = 0
			}

			s.enricher.EnrichJobInfo(resp.Users[i])
		} else {
			jobInfoFn(resp.Users[i])
		}
	}

	return resp, nil
}

func (s *Server) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.citizenstore.user_id", int64(req.UserId)))

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

	selectors := dbutils.Columns{
		tUser.Identifier,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.Sex,
		tUser.Height,
		tUserProps.UserID,
		s.customDB.Columns.User.GetVisum(tUser.Alias()),
	}

	// Field Permission Check
	fieldsAttr, err := s.p.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
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
			selectors = append(selectors, tUserProps.Job, tUserProps.JobGrade)
		case "UserProps.TrafficInfractionPoints":
			selectors = append(selectors, tUserProps.TrafficInfractionPoints)
		case "UserProps.OpenFines":
			selectors = append(selectors, tUserProps.OpenFines)
		case "UserProps.BloodType":
			selectors = append(selectors, tUserProps.BloodType)
		case "UserProps.MugShot":
			selectors = append(selectors, tUserProps.MugShot)
		}
	}

	resp := &GetUserResponse{
		User: &users.User{},
	}
	stmt := tUser.
		SELECT(
			tUser.ID,
			selectors.Get()...,
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
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}

	if resp.User == nil || resp.User.UserId <= 0 {
		return nil, errorscitizenstore.ErrJobGradeNoPermission
	}

	auditEntry.TargetUserJob = &resp.User.Job

	if slices.Contains(s.appCfg.Get().JobInfo.PublicJobs, resp.User.Job) ||
		slices.Contains(s.appCfg.Get().JobInfo.HiddenJobs, resp.User.Job) {
		// Make sure user has permission to see that grade
		jobGradesAttr, err := s.p.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceGetUserPerm, permscitizenstore.CitizenStoreServiceGetUserJobsPermField)
		if err != nil {
			return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
		}
		var jobGrades perms.JobGradeList
		if jobGradesAttr != nil {
			jobGrades = jobGradesAttr.(map[string]int32)
		}

		if len(jobGrades) == 0 && !userInfo.SuperUser {
			return nil, errorscitizenstore.ErrJobGradeNoPermission
		}

		// Make sure user has permission to see that grade, otherwise "hide" the user's job
		grade, ok := jobGrades[resp.User.Job]
		if !ok || resp.User.JobGrade > grade {
			// Skip for superuser
			if !userInfo.SuperUser {
				return nil, errorscitizenstore.ErrJobGradeNoPermission
			}
		}
	}

	// Only let user props override the job if the person isn't in a public job
	if resp.User.Props != nil && resp.User.Props.JobName != nil &&
		!slices.Contains(s.appCfg.Get().JobInfo.PublicJobs, resp.User.Job) {
		resp.User.Job = *resp.User.Props.JobName
		if resp.User.Props.JobGradeNumber != nil {
			resp.User.JobGrade = *resp.User.Props.JobGradeNumber
		} else {
			resp.User.JobGrade = 0
		}

		s.enricher.EnrichJobInfo(resp.User)
	} else {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.User)
	}

	// Check if user can see licenses and fetch them
	if slices.Contains(fields, "Licenses") {
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
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
			}
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) ListUserActivity(ctx context.Context, req *ListUserActivityRequest) (*ListUserActivityResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.citizenstore.user_id", int64(req.UserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	pag, limit := req.Pagination.GetResponseWithPageSize(15)
	resp := &ListUserActivityResponse{
		Pagination: pag,
		Activity:   []*users.UserActivity{},
	}
	// User can't see their own activities, unless they have "Own" perm attribute, or are a superuser
	fieldsAttr, err := s.p.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListUserActivityPerm, permscitizenstore.CitizenStoreServiceListUserActivityFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	if userInfo.UserId == req.UserId {
		// If isn't superuser or doesn't have 'Own' activity feed access
		if !userInfo.SuperUser || !slices.Contains(fields, "Own") {
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
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
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
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
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
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.citizenstore.user_id", int64(req.Props.UserId)))

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
		return nil, errorscitizenstore.ErrReasonRequired
	}

	// Get current user props to be able to compare
	props, err := s.getUserProps(ctx, req.Props.UserId)
	if err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}
	if props.Wanted == nil {
		wanted := false
		props.Wanted = &wanted
	}
	unemployedJob := s.appCfg.Get().JobInfo.UnemployedJob
	if props.JobName == nil {
		props.JobName = &unemployedJob.Name
	}
	if props.JobGradeNumber == nil {
		props.JobGradeNumber = &unemployedJob.Grade
	}
	if props.TrafficInfractionPoints == nil {
		props.TrafficInfractionPoints = &ZeroTrafficInfractionPoints
	}

	props.Job, props.JobGrade = s.enricher.GetJobGrade(*props.JobName, *props.JobGradeNumber)
	// Make sure a job is set
	if props.Job == nil {
		props.Job, props.JobGrade = s.enricher.GetJobGrade(unemployedJob.Name, unemployedJob.Grade)
	}

	resp := &SetUserPropsResponse{
		Props: &users.UserProps{},
	}

	// Field Permission Check
	fieldsAttr, err := s.p.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceSetUserPropsPerm, permscitizenstore.CitizenStoreServiceSetUserPropsFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	updateSets := []jet.ColumnAssigment{}
	// Generate the update sets
	if req.Props.Wanted != nil {
		if !slices.Contains(fields, "Wanted") {
			return nil, errorscitizenstore.ErrPropsWantedDenied
		}

		updateSets = append(updateSets, tUserProps.Wanted.SET(jet.Bool(*req.Props.Wanted)))
	} else {
		req.Props.Wanted = props.Wanted
	}

	if req.Props.JobName != nil {
		if !slices.Contains(fields, "Job") {
			return nil, errorscitizenstore.ErrPropsJobDenied
		}

		if slices.Contains(s.appCfg.Get().JobInfo.PublicJobs, *req.Props.JobName) {
			return nil, errorscitizenstore.ErrPropsJobPublic
		}

		if req.Props.JobGradeNumber == nil {
			grade := int32(1)
			req.Props.JobGradeNumber = &grade
		}

		req.Props.Job, req.Props.JobGrade = s.enricher.GetJobGrade(*req.Props.JobName, *req.Props.JobGradeNumber)
		if req.Props.Job == nil || req.Props.JobGrade == nil {
			return nil, errorscitizenstore.ErrPropsJobInvalid
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
		if !slices.Contains(fields, "TrafficInfractionPoints") {
			return nil, errorscitizenstore.ErrPropsTrafficPointsDenied
		}

		updateSets = append(updateSets, tUserProps.TrafficInfractionPoints.SET(jet.Uint32(*req.Props.TrafficInfractionPoints)))
	} else {
		req.Props.TrafficInfractionPoints = props.TrafficInfractionPoints
	}

	if req.Props.MugShot != nil {
		// Only update when it has actually changed
		if !slices.Contains(fields, "MugShot") {
			return nil, errorscitizenstore.ErrPropsMugShotDenied
		}

		updateSets = append(updateSets, tUserProps.MugShot.SET(jet.StringExp(jet.Raw("VALUES(`mug_shot`)"))))

		if len(req.Props.MugShot.Data) > 0 {
			if props.MugShot != nil {
				req.Props.MugShot.Url = props.MugShot.Url
			}

			if !req.Props.MugShot.IsImage() {
				return nil, errorscitizenstore.ErrFailedQuery
			}

			if err := req.Props.MugShot.Optimize(ctx); err != nil {
				return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
			}

			if err := req.Props.MugShot.Upload(ctx, s.st, filestore.MugShots, storage.FileNameSplitter(req.Props.MugShot.GetHash())); err != nil {
				return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
			}
		} else {
			// Delete mug shot from store
			if props.MugShot != nil && props.MugShot.Url != nil {
				if err := s.st.Delete(ctx, filestore.StripURLPrefix(*props.MugShot.Url)); err != nil {
					return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
				}
			}
		}
	} else {
		req.Props.MugShot = props.MugShot
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
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
			tUserProps.MugShot,
		).
		VALUES(
			req.Props.UserId,
			req.Props.Wanted,
			req.Props.JobName,
			req.Props.JobGradeNumber,
			req.Props.TrafficInfractionPoints,
			req.Props.MugShot,
		).
		ON_DUPLICATE_KEY_UPDATE(
			updateSets...,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}

	// Create user activity entries
	if *req.Props.Wanted != *props.Wanted {
		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Props.UserId, users.UserActivityType_USER_ACTIVITY_TYPE_CHANGED, "UserProps.Wanted",
			strconv.FormatBool(*props.Wanted), strconv.FormatBool(*req.Props.Wanted), req.Reason); err != nil {
			return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
		}
	}
	if *req.Props.JobName != *props.JobName || *req.Props.JobGradeNumber != *props.JobGradeNumber {
		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Props.UserId, users.UserActivityType_USER_ACTIVITY_TYPE_CHANGED, "UserProps.Job",
			fmt.Sprintf("%s|%s", props.Job.Label, props.JobGrade.Label), fmt.Sprintf("%s|%s", req.Props.Job.Label, req.Props.JobGrade.Label), req.Reason); err != nil {
			return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
		}
	}
	if *req.Props.TrafficInfractionPoints != *props.TrafficInfractionPoints {
		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Props.UserId, users.UserActivityType_USER_ACTIVITY_TYPE_CHANGED, "UserProps.TrafficInfractionPoints",
			strconv.Itoa(int(*props.TrafficInfractionPoints)), strconv.Itoa(int(*req.Props.TrafficInfractionPoints)), req.Reason); err != nil {
			return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
		}
	}
	if req.Props.MugShot != nil && (props.MugShot == nil || req.Props.MugShot.Url != props.MugShot.Url) {
		previousUrl := ""
		if props.MugShot != nil && props.MugShot.Url != nil {
			previousUrl = *props.MugShot.Url
		}
		currentUrl := ""
		if req.Props != nil && req.Props.MugShot != nil && req.Props.MugShot.Url != nil {
			currentUrl = *req.Props.MugShot.Url
		}

		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Props.UserId, users.UserActivityType_USER_ACTIVITY_TYPE_CHANGED, "UserProps.MugShot",
			previousUrl, currentUrl, req.Reason); err != nil {
			return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}

	// Get and return new user props
	user, err := s.GetUser(ctx, &GetUserRequest{
		UserId: req.Props.UserId,
	})
	if err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}

	resp.Props = user.User.Props
	// Set Job info if set
	if resp.Props != nil && resp.Props.JobName != nil {
		grade := int32(1)
		if resp.Props.JobGradeNumber != nil {
			grade = *resp.Props.JobGradeNumber
		}

		resp.Props.Job, resp.Props.JobGrade = s.enricher.GetJobGrade(*resp.Props.JobName, grade)
	}

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
			tUserProps.MugShot,
		).
		FROM(tUserProps).
		WHERE(
			tUserProps.UserID.EQ(jet.Int32(userId)),
		).
		LIMIT(1)

	var dest users.UserProps
	dest.UserId = userId
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) addUserActivity(ctx context.Context, tx qrm.DB, userId int32, targetUserId int32, activityType users.UserActivityType, key string, oldValue string, newValue string, reason string) error {
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

func (s *Server) SetProfilePicture(ctx context.Context, req *SetProfilePictureRequest) (*SetProfilePictureResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.citizenstore.user_id", int64(userInfo.UserId)))

	auditEntry := &model.FivenetAuditLog{
		Service: CitizenStoreService_ServiceDesc.ServiceName,
		Method:  "SetUserProps",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	avatarFile, err := s.getUserAvatar(ctx, userInfo.UserId)
	if err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}

	if len(req.Avatar.Data) > 0 {
		if avatarFile != nil {
			req.Avatar.Url = avatarFile.Url
		}

		if !req.Avatar.IsImage() {
			return nil, errorscitizenstore.ErrFailedQuery
		}

		if err := req.Avatar.Optimize(ctx); err != nil {
			return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
		}

		if err := req.Avatar.Upload(ctx, s.st, filestore.Avatars, storage.FileNameSplitter(req.Avatar.GetHash())); err != nil {
			return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
		}
	} else if req.Avatar.Delete != nil && *req.Avatar.Delete {
		// Delete mug shot from store
		if avatarFile != nil && avatarFile.Url != nil {
			if err := s.st.Delete(ctx, filestore.StripURLPrefix(*avatarFile.Url)); err != nil {
				return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
			}
		}
	}

	stmt := tUserProps.
		INSERT(
			tUserProps.UserID,
			tUserProps.Avatar,
		).
		VALUES(
			userInfo.UserId,
			req.Avatar,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tUserProps.Avatar.SET(jet.StringExp(jet.Raw("VALUES(`avatar`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(errorscitizenstore.ErrFailedQuery, err)
	}

	return &SetProfilePictureResponse{
		Avatar: req.Avatar,
	}, nil
}

func (s *Server) getUserAvatar(ctx context.Context, userId int32) (*filestore.File, error) {
	stmt := tUserProps.
		SELECT(
			tUserProps.Avatar.AS("usershort.avatar"),
		).
		FROM(tUserProps).
		WHERE(
			tUserProps.UserID.EQ(jet.Int32(userId)),
		).
		LIMIT(1)

	var dest users.UserShort
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.Avatar, nil
}
