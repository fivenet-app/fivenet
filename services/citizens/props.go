package citizens

import (
	context "context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	notificationsclientview "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/clientview"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
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

	if req.GetReason() == "" && len(req.GetPropertyMutations()) == 0 {
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
	propertyList, err := s.perms.AttrPropertyList(
		userInfo,
		permscitizens.CitizensService.SetUserProps.Properties,
	)
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

	// Wanted is a managed scalar. Do not permit the legacy projection to become
	// a second write path; map its PropertyMutation to the projection internally.
	if req.GetProps().Wanted != nil || req.GetProps().WantedAt != nil ||
		req.GetProps().WantedTill != nil {
		return nil, errorscitizens.ErrInvalidProperty
	}
	game := s.appCfg.Get().GetGame()
	if maxDuration := game.GetMaxWantedDurationUser(); game.GetMaxWantedDurationUserEnabled() &&
		maxDuration != nil {
		propertyregistry.ClampWantedExpiry(
			properties.PropertyScope_PROPERTY_SCOPE_USER,
			req.GetPropertyMutations(),
			time.Now().UTC().Add(maxDuration.AsDuration()),
		)
	}
	for _, mutation := range req.GetPropertyMutations() {
		managed, ok := propertyregistry.ManagedMutation(
			properties.PropertyScope_PROPERTY_SCOPE_USER,
			mutation.GetKey(),
		)
		if !ok {
			continue
		}
		wantedEnabled, err := s.store.IsSystemPropertyEnabled(ctx, managed.Key)
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		if !wantedEnabled {
			return nil, errorscitizens.ErrSystemPropertyDisabled
		}
		if !propertyList.Contains(managed.Key) {
			return nil, errorscitizens.ErrPropsWantedDenied
		}
		if err := propertyregistry.ValidateMutationShape(mutation); err != nil {
			return nil, errorscitizens.ErrInvalidProperty
		}
		definition, _ := propertyregistry.SystemPropertyDefinition(
			properties.PropertyScope_PROPERTY_SCOPE_USER,
			managed.Key,
		)
		if err := propertyregistry.ValidateMutationState(definition, mutation); err != nil {
			return nil, errorscitizens.ErrInvalidProperty
		}
		if err := propertyregistry.ValidateMutationReason(definition, mutation); err != nil {
			return nil, errorscitizens.ErrInvalidProperty
		}
		if err := managed.Apply(req.Props, mutation); err != nil {
			return nil, errorscitizens.ErrInvalidProperty
		}
	}

	if req.Props.JobName != nil {
		if !propertyList.Contains("user.job") {
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
		if !propertyList.Contains("user.traffic_infraction_points") {
			return nil, errorscitizens.ErrPropsTrafficPointsDenied
		}
	}

	if err := propertyregistry.ValidateMutationKeys(req.GetPropertyMutations()); err != nil {
		return nil, errorscitizens.ErrInvalidProperty
	}
	mutationChanges, err := customMutationChanges(req.GetPropertyMutations(), propertyList)
	if err != nil {
		return nil, err
	}
	propertyChanges, err := s.preparePropertyChanges(ctx, mutationChanges, propertyList)
	if err != nil {
		return nil, err
	}
	decorateMutationPropertyChanges(propertyChanges, req.GetPropertyMutations())

	// Users aren't allowed to set certain props, unset them so they are set to the db state
	req.Props.OpenFines = nil
	req.Props.BloodType = nil
	req.Props.Email = nil

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()
	if req.GetProps().GetLabels() != nil {
		if !propertyList.Contains("user.labels") {
			return nil, errorscitizens.ErrPropsLabelsDenied
		}

		slices.SortFunc(req.GetProps().GetLabels().GetList(), func(a, b *citizenslabels.Label) int {
			return strings.Compare(a.GetName(), b.GetName())
		})

		currentLabels, effectiveLabels, added, err := s.store.PrepareUserPropsLabelsForUpdate(
			ctx,
			tx,
			userInfo,
			req.GetProps().GetUserId(),
			req.GetProps().GetLabels(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}

		props.Labels = currentLabels
		req.Props.Labels = effectiveLabels

		valid, err := s.store.ValidateLabels(ctx, tx, userInfo.GetJob(), added)
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		if !valid {
			return nil, errorscitizens.ErrPropsLabelsDenied
		}
	}

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
	activities = append(activities, citizensstore.UserWantedMutationActivities(
		props, req.GetProps(), req.GetPropertyMutations(), &userInfo.UserId, req.GetReason(),
	)...)
	if err := s.store.ApplyWantedMutationState(
		ctx,
		tx,
		req.GetProps().GetUserId(),
		req.GetPropertyMutations(),
	); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	propertyActivities, err := s.store.ApplyPropertyMutationsWithActivity(
		ctx,
		tx,
		req.GetProps().GetUserId(),
		propertyChanges,
		req.GetPropertyMutations(),
		&userInfo.UserId,
		req.GetReason(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrInvalidProperty)
	}
	activities = append(activities, propertyActivities...)

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
	if resp.GetProps() != nil {
		resp.Props.Properties = user.GetUser().GetProperties()
	}

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

func decorateMutationPropertyChanges(
	changes []citizensstore.PropertyChange,
	mutations []*properties.PropertyMutation,
) {
	byKey := make(map[string]*properties.PropertyMutation, len(mutations))
	for _, mutation := range mutations {
		if mutation != nil && (mutation.GetSetValue() != nil || mutation.GetClearValue() != nil) {
			byKey[mutation.GetKey()] = mutation
		}
	}
	for index := range changes {
		mutation := byKey[changes[index].Property.GetKey()]
		if mutation == nil {
			continue
		}
		changes[index].Cause = usersactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_MANUAL
		changes[index].Lifecycle = userMutationLifecycle(mutation)
	}
}

func userMutationLifecycle(
	mutation *properties.PropertyMutation,
) *usersactivity.PropertyLifecycleChange {
	lifecycle := propertyregistry.MutationLifecycleForMutation(mutation)
	if lifecycle == nil {
		return nil
	}
	return &usersactivity.PropertyLifecycleChange{
		ExpiresAt: lifecycle.ExpiresAt, ExpiryCleared: lifecycle.ExpiryCleared,
		StateReason: lifecycle.StateReason, StateReasonCleared: lifecycle.StateReasonCleared,
	}
}

// customMutationChanges separates a mutation's optional scalar change from
// its lifecycle patch. State-only mutations deliberately produce no scalar
// change, but still require the caller's property permission.
func customMutationChanges(
	mutations []*properties.PropertyMutation,
	allowed *permissionsattributes.PropertyList,
) ([]*propertyregistry.ScalarChange, error) {
	changes := make([]*propertyregistry.ScalarChange, 0, len(mutations))
	seen := make(map[string]struct{}, len(mutations))
	for _, mutation := range mutations {
		if err := propertyregistry.ValidateMutationShape(mutation); err != nil {
			return nil, errorscitizens.ErrInvalidProperty
		}
		if !propertyregistry.IsCustomPropertyKey(
			properties.PropertyScope_PROPERTY_SCOPE_USER,
			mutation.GetKey(),
		) {
			if _, managed := propertyregistry.ManagedMutation(
				properties.PropertyScope_PROPERTY_SCOPE_USER,
				mutation.GetKey(),
			); managed {
				continue
			}
			return nil, errorscitizens.ErrPropsPropertiesDenied
		}
		if !allowed.Contains(mutation.GetKey()) {
			return nil, errorscitizens.ErrPropsPropertiesDenied
		}
		if _, ok := seen[mutation.GetKey()]; ok {
			return nil, errorscitizens.ErrInvalidProperty
		}
		seen[mutation.GetKey()] = struct{}{}
		if mutation.GetSetValue() != nil {
			changes = append(
				changes,
				&propertyregistry.ScalarChange{
					Key:   mutation.GetKey(),
					Value: mutation.GetSetValue(),
				},
			)
		} else if mutation.GetClearValue() != nil {
			changes = append(changes, &propertyregistry.ScalarChange{Key: mutation.GetKey()})
		} else if mutation.GetState() == nil && mutation.Reason == nil && mutation.ClearReason == nil {
			return nil, errorscitizens.ErrInvalidProperty
		}
	}
	return changes, nil
}

func (s *Server) preparePropertyChanges(
	ctx context.Context,
	values []*propertyregistry.ScalarChange,
	allowed *permissionsattributes.PropertyList,
) ([]citizensstore.PropertyChange, error) {
	err := propertyregistry.ValidateCustomPropertyChanges(
		ctx,
		properties.PropertyScope_PROPERTY_SCOPE_USER,
		values,
		func(key string) bool { return allowed.Contains(key) },
		s.store.ResolvePropertyDefinitions,
	)
	if err != nil {
		switch {
		case errors.Is(err, propertyregistry.ErrInvalidPropertyChange):
			return nil, errorscitizens.ErrInvalidProperty
		case errors.Is(err, propertyregistry.ErrPropertyChangeDenied):
			return nil, errorscitizens.ErrPropsPropertiesDenied
		default:
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	changes := make([]citizensstore.PropertyChange, 0, len(values))
	for _, value := range values {
		if value.Value == nil {
			changes = append(
				changes,
				citizensstore.PropertyChange{
					Property: &properties.Property{Key: value.Key},
					Delete:   true,
				},
			)
			continue
		}
		changes = append(
			changes,
			citizensstore.PropertyChange{
				Property: &properties.Property{Key: value.Key, Value: value.Value},
			},
		)
	}
	return changes, nil
}
