package citizens

import (
	context "context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var ZeroTrafficInfractionPoints uint32 = 0

func (s *Server) ListCitizens(
	ctx context.Context,
	req *pbcitizens.ListCitizensRequest,
) (*pbcitizens.ListCitizensResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Field Permission Check
	fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	resp, err := s.store.ListCitizens(ctx, req, citizensstore.ListCitizensOptions{
		IncludePhoneNumber: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		),
		IncludeWanted: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsWanted,
		),
		IncludeJob: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsJob,
		),
		IncludeTrafficInfractionPoints: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsTrafficInfractionPoints,
		),
		IncludeOpenFines: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsOpenFines,
		),
		IncludeBloodType: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsBloodType,
		),
		IncludeMugshot: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsMugshot,
		),
		IncludeEmail: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsEmail,
		),
	})
	if err != nil {
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

	infoOnly := req.InfoOnly != nil && req.GetInfoOnly()

	// Field Permission Check
	fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	resp, err := s.store.GetUser(ctx, req, citizensstore.GetUserOptions{
		IncludePropsUpdated: fields.Len() > 0,
		IncludePhoneNumber: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		),
		IncludeWanted: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsWanted,
		),
		IncludeJob: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsJob,
		),
		IncludeTrafficInfractionPoints: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsTrafficInfractionPoints,
		),
		IncludeOpenFines: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsOpenFines,
		),
		IncludeBloodType: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsBloodType,
		),
		IncludeMugshot: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsMugshot,
		),
		IncludeEmail: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsEmail,
		),
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if resp == nil || resp.GetUser() == nil || resp.GetUser().GetUserId() <= 0 {
		return nil, errorscitizens.ErrCitizenNotFound
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
		resp.User.Props = &usersprops.UserProps{UserId: resp.GetUser().GetUserId()}
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
		condition := mysql.AND(
			tCitizensLabelsJob.DeletedAt.IS_NULL(),
			tCitizenLabels.UserID.EQ(mysql.Int32(req.GetUserId())),
		)
		if !userInfo.GetSuperuser() {
			jobAccessExists := s.labelsAccess.ACLAccessExistsCondition(
				tCitizensLabelsJob.ID,
				userInfo,
				int32(citizenslabels.AccessLevel_ACCESS_LEVEL_VIEW),
			)

			condition = condition.AND(jobAccessExists)
		}

		labels, err := s.store.GetUserLabels(ctx, s.db, condition)
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		resp.User.Props.Labels = labels
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
