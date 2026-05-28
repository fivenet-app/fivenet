package citizens

import (
	context "context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tFiles = table.FivenetFiles.AS("mugshot")

var ZeroTrafficInfractionPoints uint32 = 0

func (s *Server) ListCitizens(
	ctx context.Context,
	req *pbcitizens.ListCitizensRequest,
) (*pbcitizens.ListCitizensResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tUser := table.FivenetUser.AS("user")

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
	fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	for _, field := range fields.Values() {
		switch field {
		case permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber:
			selectors = append(selectors, tUser.PhoneNumber)

			if req.GetPhoneNumber() != "" {
				phoneNumber := dbutils.PrepareForLikeSearch(req.GetPhoneNumber())
				condition = condition.AND(tUser.PhoneNumber.LIKE(mysql.String(phoneNumber)))
			}

		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsWanted:
			selectors = append(selectors, tUserProps.Wanted)

			if req.Wanted != nil && req.GetWanted() {
				condition = condition.AND(tUserProps.Wanted.IS_TRUE())

				orderBys = append(orderBys, tUserProps.UpdatedAt.DESC())
			}

		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsJob:
			selectors = append(selectors, tUserProps.Job, tUserProps.JobGrade)

		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsTrafficInfractionPoints:
			selectors = append(selectors, tUserProps.TrafficInfractionPoints)

			if req.TrafficInfractionPoints != nil && req.GetTrafficInfractionPoints() > 0 {
				condition = condition.AND(
					tUserProps.TrafficInfractionPoints.GT_EQ(
						mysql.Uint32(req.GetTrafficInfractionPoints()),
					),
				)
			}

		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsOpenFines:
			selectors = append(selectors, tUserProps.OpenFines)

			if req.OpenFines != nil && req.GetOpenFines() > 0 {
				condition = condition.AND(
					tUserProps.OpenFines.GT_EQ(mysql.Int64(req.GetOpenFines())),
				)
			}

		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsBloodType:
			selectors = append(selectors, tUserProps.BloodType)

		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsMugshot:
			selectors = append(selectors,
				tUserProps.MugshotFileID,
				tFiles.ID,
				tFiles.FilePath,
			)

		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsEmail:
			selectors = append(selectors, tUserProps.Email)
		}
	}

	if search := dbutils.PrepareForLikeSearch(req.GetSearch()); search != "" {
		condition = condition.AND(
			mysql.CONCAT(tUser.Firstname, mysql.String(" "), tUser.Lastname).
				LIKE(mysql.String(search)),
		)
	}

	if req.GetDateofbirth() != "" {
		dateofbirth := dbutils.PrepareForLikeSearch(req.GetDateofbirth())

		condition = condition.AND(
			tUser.Dateofbirth.LIKE(
				mysql.String(dateofbirth),
			),
		)
	}

	if req.GetMinHeight() > 0 {
		condition = condition.AND(
			tUser.Height.GT_EQ(
				mysql.Float(float64(req.GetMinHeight())),
			),
		)
	}
	if req.GetMaxHeight() > 0 {
		condition = condition.AND(
			tUser.Height.LT_EQ(
				mysql.Float(float64(req.GetMaxHeight())),
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
	if req.GetSort() != nil && len(req.GetSort().GetColumns()) > 0 {
		for _, sc := range req.GetSort().GetColumns() {
			var column mysql.Column
			switch sc.GetId() {
			case "trafficInfractionPoints":
				if fields.Contains(
					permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsTrafficInfractionPoints,
				) {
					column = tUserProps.TrafficInfractionPoints
				}
			case "openFines":
				if fields.Contains(
					permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsOpenFines,
				) {
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

			if sc.GetDesc() {
				orderBys = append(orderBys,
					column.DESC(),
					tUser.Lastname.DESC(),
				)
			} else {
				orderBys = append(orderBys,
					column.ASC(),
					tUser.Lastname.ASC(),
				)
			}
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

	grpc_audit.SetTargetUser(ctx, req.GetUserId(), "")

	tUser := table.FivenetUser.AS("user")

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
	fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if fields.Len() > 0 {
		selectors = append(selectors, tUserProps.UpdatedAt)
	}

	for _, field := range fields.Values() {
		switch field {
		case permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber:
			selectors = append(selectors, tUser.PhoneNumber)
		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsWanted:
			selectors = append(selectors, tUserProps.Wanted)
		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsJob:
			selectors = append(selectors, tUserProps.Job, tUserProps.JobGrade)
		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsTrafficInfractionPoints:
			selectors = append(selectors, tUserProps.TrafficInfractionPoints)
		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsOpenFines:
			selectors = append(selectors, tUserProps.OpenFines)
		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsBloodType:
			selectors = append(selectors, tUserProps.BloodType)
		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsMugshot:
			selectors = append(selectors,
				tUserProps.MugshotFileID,
				tFiles.ID,
				tFiles.FilePath,
			)
		case permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsEmail:
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
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorscitizens.ErrCitizenNotFound
		}
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if resp.GetUser() == nil || resp.GetUser().GetUserId() <= 0 {
		return nil, errorscitizens.ErrJobGradeNoPermission
	}

	grpc_audit.SetTargetUser(ctx, resp.GetUser().GetUserId(), resp.GetUser().GetJob())

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
		resp.User.Props = &usersprops.UserProps{
			UserId: resp.GetUser().GetUserId(),
		}
	}

	// Check if user can see licenses and fetch them
	if !infoOnly &&
		fields.Contains(permscitizens.CitizensServiceListCitizensFieldsPermValueLicenses) {
		tCitizenLicenses := table.FivenetUserLicenses
		tLicenses := table.FivenetLicenses

		stmt := tCitizenLicenses.
			SELECT(
				tLicenses.Type.AS("license.type"),
				tLicenses.Label.AS("license.label"),
			).
			FROM(
				tCitizenLicenses.
					LEFT_JOIN(tLicenses,
						tCitizenLicenses.Type.EQ(tLicenses.Type)),
			).
			WHERE(tCitizenLicenses.UserID.EQ(mysql.Int32(req.GetUserId()))).
			LIMIT(15)

		if err := stmt.QueryContext(ctx, s.db, &resp.User.Licenses); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
			}
		}
	}

	if fields.Contains(permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsLabels) {
		attributes, err := s.getUserLabels(ctx, userInfo, req.GetUserId())
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		resp.User.Props.Labels = attributes
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

	return resp, nil
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
		permscitizens.CitizensService.GetUser.Jobs,
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
