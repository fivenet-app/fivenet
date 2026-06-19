package citizens

import (
	context "context"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	notificationsclientview "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/clientview"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) SetUserProps(
	ctx context.Context,
	req *pbcitizens.SetUserPropsRequest,
) (*pbcitizens.SetUserPropsResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{citizenIDLogFieldKey, req.GetProps().GetUserId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	grpc_audit.SetTargetUser(ctx, req.GetProps().GetUserId(), "")

	if req.GetReason() == "" {
		return nil, errorscitizens.ErrReasonRequired
	}

	// Get current user props to be able to compare
	props, err := s.store.GetUserProps(ctx, s.db, req.GetProps().GetUserId())
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
		props.Labels = &citizenslabels.Labels{
			List: []*citizenslabels.Label{},
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
		Props: &usersprops.UserProps{},
	}

	// Field Permission Check
	fields, err := permscitizens.CitizensService.SetUserProps.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	u, err := s.store.GetUserAccess(ctx, req.GetProps().GetUserId())
	if err != nil {
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
		if !fields.Contains(permscitizens.CitizensServiceSetUserPropsFieldsPermValueWanted) {
			return nil, errorscitizens.ErrPropsWantedDenied
		}
	}

	if req.Props.JobName != nil {
		if !fields.Contains(permscitizens.CitizensServiceSetUserPropsFieldsPermValueJob) {
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
		if !fields.Contains(
			permscitizens.CitizensServiceSetUserPropsFieldsPermValueTrafficInfractionPoints,
		) {
			return nil, errorscitizens.ErrPropsTrafficPointsDenied
		}
	}

	// Users aren't allowed to set certain props, unset them so they are set to the db state
	req.Props.OpenFines = nil
	req.Props.BloodType = nil
	req.Props.Email = nil

	if req.GetProps().GetLabels() != nil {
		if !fields.Contains(permscitizens.CitizensServiceSetUserPropsFieldsPermValueLabels) {
			return nil, errorscitizens.ErrPropsLabelsDenied
		}

		if req.GetProps().GetLabels().GetList() == nil {
			req.Props.Labels.List = []*citizenslabels.Label{}
		}

		slices.SortFunc(req.GetProps().GetLabels().GetList(), func(a, b *citizenslabels.Label) int {
			return strings.Compare(a.GetName(), b.GetName())
		})

		added, _ := utils.SlicesDifferenceFunc(
			props.GetLabels().GetList(),
			req.GetProps().GetLabels().GetList(),
			func(in *citizenslabels.Label) int64 {
				return in.GetId()
			},
		)

		valid, err := s.store.ValidateLabels(ctx, userInfo.GetJob(), added)
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

	activities, err := s.store.HandleUserPropsChanges(
		ctx,
		tx,
		props,
		req.GetProps(),
		&userInfo.UserId,
		req.GetReason(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if err := usersactivity.CreateUserActivities(ctx, tx, activities...); err != nil {
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
	s.notifi.SendObjectEvent(ctx, &notificationsclientview.ObjectEvent{
		Type:      notificationsclientview.ObjectType_OBJECT_TYPE_CITIZEN,
		Id:        &userId,
		EventType: notificationsclientview.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return resp, nil
}
