package citizens

import (
	context "context"
	"errors"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbcitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2025/services/citizens/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var (
	tUserProps = table.FivenetUserProps

	tFiles = table.FivenetFiles.AS("mugshot")
)

var ZeroTrafficInfractionPoints uint32 = 0

func (s *Server) ListCitizens(
	ctx context.Context,
	req *pbcitizens.ListCitizensRequest,
) (*pbcitizens.ListCitizensResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tUser := tables.User().AS("user")

	selectors := dbutils.Columns{
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
	orderBys := []mysql.OrderByClause{}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permscitizens.CitizensServicePerm,
		permscitizens.CitizensServiceListCitizensPerm,
		permscitizens.CitizensServiceListCitizensFieldsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	for _, field := range fields.GetStrings() {
		switch field {
		case "PhoneNumber":
			selectors = append(selectors, tUser.PhoneNumber)

			if req.PhoneNumber != nil && req.GetPhoneNumber() != "" {
				phoneNumber := strings.ReplaceAll(
					strings.ReplaceAll(req.GetPhoneNumber(), "%", ""),
					" ",
					"",
				) + "%"
				condition = condition.AND(tUser.PhoneNumber.LIKE(mysql.String(phoneNumber)))
			}

		case "UserProps.Wanted":
			selectors = append(selectors, tUserProps.Wanted)

			if req.Wanted != nil && req.GetWanted() {
				condition = condition.AND(tUserProps.Wanted.IS_TRUE())

				orderBys = append(orderBys, tUserProps.UpdatedAt.DESC())
			}

		case "UserProps.Job":
			selectors = append(selectors, tUserProps.Job, tUserProps.JobGrade)

		case "UserProps.TrafficInfractionPoints":
			selectors = append(selectors, tUserProps.TrafficInfractionPoints)

			if req.TrafficInfractionPoints != nil && req.GetTrafficInfractionPoints() > 0 {
				condition = condition.AND(
					tUserProps.TrafficInfractionPoints.GT_EQ(
						mysql.Uint32(req.GetTrafficInfractionPoints()),
					),
				)
			}

		case "UserProps.OpenFines":
			selectors = append(selectors, tUserProps.OpenFines)

			if req.OpenFines != nil && req.GetOpenFines() > 0 {
				condition = condition.AND(
					tUserProps.OpenFines.GT_EQ(mysql.Int64(req.GetOpenFines())),
				)
			}

		case "UserProps.BloodType":
			selectors = append(selectors, tUserProps.BloodType)

		case "UserProps.Mugshot":
			selectors = append(selectors,
				tUserProps.MugshotFileID,
				tFiles.ID,
				tFiles.FilePath,
			)

		case "UserProps.Email":
			selectors = append(selectors, tUserProps.Email)
		}
	}

	req.Search = strings.TrimSpace(req.GetSearch())
	req.Search = strings.ReplaceAll(req.GetSearch(), "%", "")
	req.Search = strings.ReplaceAll(req.GetSearch(), " ", "%")
	req.Search = strings.ReplaceAll(req.GetSearch(), "\t", " ")
	if req.GetSearch() != "" {
		req.Search = "%" + req.GetSearch() + "%"
		condition = condition.AND(
			mysql.CONCAT(tUser.Firstname, mysql.String(" "), tUser.Lastname).
				LIKE(mysql.String(req.GetSearch())),
		)
	}

	if req.Dateofbirth != nil && req.GetDateofbirth() != "" {
		condition = condition.AND(
			tUser.Dateofbirth.LIKE(
				mysql.String(strings.ReplaceAll(req.GetDateofbirth(), "%", " ") + "%"),
			),
		)
	}

	// Get total count of values
	countStmt := tUser.
		SELECT(
			mysql.COUNT(tUser.ID).AS("data_count.total"),
		).
		OPTIMIZER_HINTS(mysql.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponse(count.Total)
	resp := &pbcitizens.ListCitizensResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	if req.GetSort() != nil {
		var column mysql.Column
		switch req.GetSort().GetColumn() {
		case "trafficInfractionPoints":
			if fields.Contains("UserProps.TrafficInfractionPoints") {
				column = tUserProps.TrafficInfractionPoints
			}
		case "openFines":
			if fields.Contains("UserProps.OpenFines") {
				column = tUserProps.OpenFines
			}
		case "name":
			fallthrough
		default:
			column = tUser.Firstname
		}
		if column == nil {
			column = tUser.Firstname
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
			orderBys = append(orderBys,
				column.ASC(),
				tUser.Lastname.ASC(),
			)
		} else {
			orderBys = append(orderBys,
				column.DESC(),
				tUser.Lastname.DESC(),
			)
		}
	} else {
		orderBys = append(orderBys,
			tUser.Firstname.ASC(),
			tUser.Lastname.ASC(),
		)
	}

	stmt := tUser.
		SELECT(
			tUser.ID,
			selectors.Get()...,
		).
		OPTIMIZER_HINTS(mysql.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(tUser.
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUser.ID),
			).
			LEFT_JOIN(tFiles,
				tFiles.ID.EQ(tUserProps.MugshotFileID),
			),
		).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	resp.GetPagination().Update(len(resp.GetUsers()))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetUsers() {
		if resp.GetUsers()[i].GetProps() != nil && resp.Users[i].Props.JobName != nil {
			resp.Users[i].Job = resp.GetUsers()[i].GetProps().GetJobName()
			if resp.Users[i].Props.JobGradeNumber != nil {
				resp.Users[i].JobGrade = resp.GetUsers()[i].GetProps().GetJobGradeNumber()
			} else {
				resp.Users[i].JobGrade = 0
			}

			s.enricher.EnrichJobInfo(resp.GetUsers()[i])
		} else {
			jobInfoFn(resp.GetUsers()[i])
		}
	}

	return resp, nil
}

func (s *Server) GetUser(
	ctx context.Context,
	req *pbcitizens.GetUserRequest,
) (*pbcitizens.GetUserResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.citizens.user_id", req.GetUserId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service:      pbcitizens.CitizensService_ServiceDesc.ServiceName,
		Method:       "GetUser",
		UserId:       userInfo.GetUserId(),
		UserJob:      userInfo.GetJob(),
		TargetUserId: &req.UserId,
		State:        audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	tUser := tables.User().AS("user")

	selectors := dbutils.Columns{
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

	infoOnly := req.InfoOnly != nil && req.GetInfoOnly()

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permscitizens.CitizensServicePerm,
		permscitizens.CitizensServiceListCitizensPerm,
		permscitizens.CitizensServiceListCitizensFieldsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if fields.Strings != nil {
		selectors = append(selectors, tUserProps.UpdatedAt)
	}

	for _, field := range fields.GetStrings() {
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
		case "UserProps.Mugshot":
			selectors = append(selectors,
				tUserProps.MugshotFileID,
				tFiles.ID,
				tFiles.FilePath,
			)
		case "UserProps.Email":
			selectors = append(selectors, tUserProps.Email)
		}
	}

	resp := &pbcitizens.GetUserResponse{
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
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tUserProps.MugshotFileID),
				),
		).
		WHERE(tUser.ID.EQ(mysql.Int32(req.GetUserId()))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, resp.GetUser()); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if resp.GetUser() == nil || resp.GetUser().GetUserId() <= 0 {
		return nil, errorscitizens.ErrJobGradeNoPermission
	}

	auditEntry.TargetUserJob = &resp.User.Job

	if slices.Contains(s.appCfg.Get().JobInfo.GetPublicJobs(), resp.GetUser().GetJob()) ||
		slices.Contains(s.appCfg.Get().JobInfo.GetHiddenJobs(), resp.GetUser().GetJob()) {
		// Make sure user has permission to see that grade
		check, err := s.checkIfUserCanAccess(
			userInfo,
			resp.GetUser().GetJob(),
			resp.GetUser().GetJobGrade(),
		)
		if err != nil {
			return nil, err
		}
		if !check {
			return nil, errorscitizens.ErrJobGradeNoPermission
		}
	}

	// Only let user props override the job if the person isn't in a public job
	if resp.GetUser().GetProps() != nil && resp.User.Props.JobName != nil &&
		!slices.Contains(s.appCfg.Get().JobInfo.GetPublicJobs(), resp.GetUser().GetJob()) {
		resp.User.Job = resp.GetUser().GetProps().GetJobName()
		if resp.User.Props.JobGradeNumber != nil {
			resp.User.JobGrade = resp.GetUser().GetProps().GetJobGradeNumber()
		} else {
			resp.User.JobGrade = 0
		}

		s.enricher.EnrichJobInfo(resp.GetUser())
	} else {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.GetUser())
	}

	if resp.GetUser().GetProps() == nil {
		resp.User.Props = &users.UserProps{
			UserId: resp.GetUser().GetUserId(),
		}
	}

	// Check if user can see licenses and fetch them
	if !infoOnly && fields.Contains("Licenses") {
		tLicenses := tables.Licenses()
		tCitizensLicenses := tables.UserLicenses()

		stmt := tUser.
			SELECT(
				tCitizensLicenses.Type.AS("license.type"),
				tLicenses.Label.AS("license.label"),
			).
			FROM(
				tCitizensLicenses.
					INNER_JOIN(tUser,
						tCitizensLicenses.Owner.EQ(tUser.Identifier),
					).
					LEFT_JOIN(tLicenses,
						tLicenses.Type.EQ(tCitizensLicenses.Type)),
			).
			WHERE(tUser.ID.EQ(mysql.Int32(req.GetUserId()))).
			LIMIT(15)

		if err := stmt.QueryContext(ctx, s.db, &resp.User.Licenses); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
			}
		}
	}

	if fields.Contains("UserProps.Labels") {
		attributes, err := s.getUserLabels(ctx, userInfo, req.GetUserId())
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		resp.User.Props.Labels = attributes
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) SetUserProps(
	ctx context.Context,
	req *pbcitizens.SetUserPropsRequest,
) (*pbcitizens.SetUserPropsResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.citizens.user_id", req.GetProps().GetUserId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service:      pbcitizens.CitizensService_ServiceDesc.ServiceName,
		Method:       "SetUserProps",
		UserId:       userInfo.GetUserId(),
		UserJob:      userInfo.GetJob(),
		TargetUserId: &req.Props.UserId,
		State:        audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if req.GetReason() == "" {
		return nil, errorscitizens.ErrReasonRequired
	}

	// Get current user props to be able to compare
	props, err := s.getUserProps(ctx, userInfo, req.GetProps().GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if props.Wanted == nil {
		wanted := false
		props.Wanted = &wanted
	}
	unemployedJob := s.appCfg.Get().JobInfo.GetUnemployedJob()
	if props.JobName == nil {
		props.JobName = &unemployedJob.Name
	}
	if props.JobGradeNumber == nil {
		props.JobGradeNumber = &unemployedJob.Grade
	}
	if props.TrafficInfractionPoints == nil {
		props.TrafficInfractionPoints = &ZeroTrafficInfractionPoints
	}
	if props.GetLabels() == nil {
		props.Labels = &users.Labels{
			List: []*users.Label{},
		}
	}

	props.Job, props.JobGrade = s.enricher.GetJobGrade(
		props.GetJobName(),
		props.GetJobGradeNumber(),
	)
	// Make sure a job is set
	if props.GetJob() == nil {
		props.Job, props.JobGrade = s.enricher.GetJobGrade(
			unemployedJob.GetName(),
			unemployedJob.GetGrade(),
		)
	}

	resp := &pbcitizens.SetUserPropsResponse{
		Props: &users.UserProps{},
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permscitizens.CitizensServicePerm,
		permscitizens.CitizensServiceSetUserPropsPerm,
		permscitizens.CitizensServiceSetUserPropsFieldsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	tUser := tables.User().AS("user")

	stmt := tUser.
		SELECT(
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
		).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tUserProps.MugshotFileID),
				),
		).
		WHERE(tUser.ID.EQ(mysql.Int32(req.GetProps().GetUserId()))).
		LIMIT(1)

	u := &users.User{}
	if err := stmt.QueryContext(ctx, s.db, u); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if u.GetUserId() <= 0 {
		return nil, errorscitizens.ErrJobGradeNoPermission
	}

	check, err := s.checkIfUserCanAccess(userInfo, u.GetJob(), u.GetJobGrade())
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errorscitizens.ErrJobGradeNoPermission
	}

	// Generate the update sets
	if req.Props.Wanted != nil {
		if !fields.Contains("Wanted") {
			return nil, errorscitizens.ErrPropsWantedDenied
		}
	}

	if req.Props.JobName != nil {
		if !fields.Contains("Job") {
			return nil, errorscitizens.ErrPropsJobDenied
		}

		if slices.Contains(s.appCfg.Get().JobInfo.GetPublicJobs(), req.GetProps().GetJobName()) {
			return nil, errorscitizens.ErrPropsJobPublic
		}

		if req.Props.JobGradeNumber == nil {
			grade := s.cfg.Game.StartJobGrade
			req.Props.JobGradeNumber = &grade
		}

		req.Props.Job, req.Props.JobGrade = s.enricher.GetJobGrade(
			req.GetProps().GetJobName(),
			req.GetProps().GetJobGradeNumber(),
		)
		if req.GetProps().GetJob() == nil || req.GetProps().GetJobGrade() == nil {
			return nil, errorscitizens.ErrPropsJobInvalid
		}
	}

	if req.Props.TrafficInfractionPoints != nil {
		if !fields.Contains("TrafficInfractionPoints") {
			return nil, errorscitizens.ErrPropsTrafficPointsDenied
		}
	}

	// Users aren't allowed to set certain props, unset them so they are set to the db state
	req.Props.OpenFines = nil
	req.Props.BloodType = nil
	req.Props.Email = nil

	if req.GetProps().GetLabels() != nil {
		if !fields.Contains("Labels") {
			return nil, errorscitizens.ErrPropsLabelsDenied
		}

		if req.Props.Labels.List == nil {
			req.Props.Labels.List = []*users.Label{}
		}

		slices.SortFunc(req.GetProps().GetLabels().GetList(), func(a, b *users.Label) int {
			return strings.Compare(a.GetName(), b.GetName())
		})

		added, _ := utils.SlicesDifferenceFunc(
			props.GetLabels().GetList(),
			req.GetProps().GetLabels().GetList(),
			func(in *users.Label) int64 {
				return in.GetId()
			},
		)

		valid, err := s.validateLabels(ctx, userInfo, added)
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		if !valid {
			return nil, errorscitizens.ErrPropsLabelsDenied
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	activities, err := props.HandleChanges(
		ctx,
		tx,
		req.GetProps(),
		&userInfo.UserId,
		req.GetReason(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if err := users.CreateUserActivities(ctx, tx, activities...); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	// Get and return new user props
	user, err := s.GetUser(ctx, &pbcitizens.GetUserRequest{
		UserId: req.GetProps().GetUserId(),
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	resp.Props = user.GetUser().GetProps()
	// Set Job info if set
	if resp.GetProps() != nil && resp.Props.JobName != nil {
		grade := s.cfg.Game.StartJobGrade
		if resp.Props.JobGradeNumber != nil {
			grade = resp.GetProps().GetJobGradeNumber()
		}

		resp.Props.Job, resp.Props.JobGrade = s.enricher.GetJobGrade(
			resp.GetProps().GetJobName(),
			grade,
		)
	}

	userId := int64(user.GetUser().GetUserId())
	s.notifi.SendObjectEvent(ctx, &notifications.ObjectEvent{
		Type:      notifications.ObjectType_OBJECT_TYPE_CITIZEN,
		Id:        &userId,
		EventType: notifications.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return resp, nil
}

func (s *Server) getUserProps(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	userId int32,
) (*users.UserProps, error) {
	tUserProps := tUserProps.AS("user_props")
	tFiles := table.FivenetFiles.AS("mugshot")

	stmt := tUserProps.
		SELECT(
			tUserProps.UserID,
			tUserProps.UpdatedAt,
			tUserProps.Wanted,
			tUserProps.Job,
			tUserProps.JobGrade,
			tUserProps.TrafficInfractionPoints,
			tUserProps.TrafficInfractionPointsUpdatedAt,
			tUserProps.MugshotFileID,
			tFiles.ID,
			tFiles.FilePath,
		).
		FROM(
			tUserProps.
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tUserProps.MugshotFileID),
				),
		).
		WHERE(
			tUserProps.UserID.EQ(mysql.Int32(userId)),
		).
		LIMIT(1)

	var dest users.UserProps
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	dest.UserId = userId

	attributes, err := s.getUserLabels(ctx, userInfo, userId)
	if err != nil {
		return nil, err
	}
	dest.Labels = attributes

	return &dest, nil
}

func (s *Server) checkIfUserCanAccess(
	userInfo *userinfo.UserInfo,
	targetUserJob string,
	targetUserGrade int32,
) (bool, error) {
	// Skip if user is job unemployed
	unemployedJob := s.appCfg.Get().JobInfo.GetUnemployedJob()
	if unemployedJob.GetName() == targetUserJob {
		return true, nil
	}

	// If the user is not part of public or hidden jobs (e.g., police, medics), allow access
	if !slices.Contains(s.appCfg.Get().JobInfo.GetPublicJobs(), targetUserJob) &&
		!slices.Contains(s.appCfg.Get().JobInfo.GetHiddenJobs(), targetUserJob) {
		return true, nil
	}

	jobGrades, err := s.ps.AttrJobGradeList(
		userInfo,
		permscitizens.CitizensServicePerm,
		permscitizens.CitizensServiceGetUserPerm,
		permscitizens.CitizensServiceGetUserJobsPermField,
	)
	if err != nil {
		return false, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if jobGrades.Len() == 0 && !userInfo.GetSuperuser() {
		return false, errorscitizens.ErrJobGradeNoPermission
	}

	// Make sure user has permission to see that job's grade, otherwise deny access to the user
	if ok := jobGrades.HasJobGrade(targetUserJob, targetUserGrade); !ok &&
		!userInfo.GetSuperuser() {
		return false, errorscitizens.ErrJobGradeNoPermission
	}

	return true, nil
}
