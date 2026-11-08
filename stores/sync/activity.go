package syncstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	activity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) AddActivity(
	ctx context.Context,
	req *pbsync.AddActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	resp := &pbsync.AddActivityResponse{}
	if req == nil {
		return resp, nil
	}

	switch d := req.GetActivity().(type) {
	case *pbsync.AddActivityRequest_UserOauth2:
		if err := s.handleUserOauth2(ctx, d.UserOauth2); err != nil {
			return nil, fmt.Errorf("failed to handle UserOauth2 activity. %w", err)
		}

	case *pbsync.AddActivityRequest_Dispatch:
		if s.dispatches != nil {
			if _, err := s.dispatches.Create(ctx, d.Dispatch); err != nil {
				return nil, fmt.Errorf("failed to create dispatch. %w", err)
			}
		}

	case *pbsync.AddActivityRequest_UserActivity:
		if err := s.createUserActivities(ctx, s.db, d.UserActivity); err != nil {
			return nil, fmt.Errorf("failed to create user activities. %w", err)
		}

	case *pbsync.AddActivityRequest_UserProps:
		if err := s.handleUserProps(ctx, d.UserProps, nil); err != nil {
			return nil, fmt.Errorf("failed to handle UserProps activity. %w", err)
		}

	case *pbsync.AddActivityRequest_ColleagueActivity:
		if err := s.createColleagueActivity(ctx, s.db, d.ColleagueActivity); err != nil {
			return nil, fmt.Errorf("failed to create jobs user activities. %w", err)
		}

	case *pbsync.AddActivityRequest_ColleagueProps:
		if err := s.handleColleagueProps(ctx, d.ColleagueProps, nil); err != nil {
			return nil, fmt.Errorf("failed to handle ColleagueProps activity. %w", err)
		}

	case *pbsync.AddActivityRequest_JobTimeclock:
		if err := s.handleTimeclockEntry(ctx, d.JobTimeclock); err != nil {
			return nil, fmt.Errorf("failed to handle JobTimeclock activity. %w", err)
		}

	case *pbsync.AddActivityRequest_UserUpdate:
		if err := s.handleUserUpdate(ctx, d.UserUpdate); err != nil {
			return nil, fmt.Errorf("failed to handle UserUpdate activity. %w", err)
		}

	case *pbsync.AddActivityRequest_AccountUpdate:
		if err := s.handleAccountUpdate(ctx, d.AccountUpdate); err != nil {
			return nil, fmt.Errorf("failed to handle AccountUpdate activity. %w", err)
		}
	}

	return resp, nil
}

func (s *Store) AddUserActivity(
	ctx context.Context,
	req *pbsync.AddUserActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	if activity := req.GetUserActivity(); activity != nil && req.GetSourceUser() != nil &&
		req.GetSourceUser().UserId != nil {
		sourceUserID := req.GetSourceUser().GetUserId()
		activity.SourceUserId = &sourceUserID
	}
	if err := s.createUserActivities(ctx, s.db, req.GetUserActivity()); err != nil {
		return nil, fmt.Errorf("failed to create user activities. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddUserUpdate(
	ctx context.Context,
	req *pbsync.AddUserUpdateRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.handleUserUpdate(ctx, req.GetUserUpdate()); err != nil {
		return nil, fmt.Errorf("failed to handle user update data. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

// SetUserProperties applies managed and generic property changes from a
// trusted Sync client.
func (s *Store) SetUserProperties(
	ctx context.Context,
	req *pbsync.SetUserPropertiesRequest,
) (*pbsync.SetUserPropertiesResponse, error) {
	s.normalizeSyncWantedExpiries(
		properties.PropertyScope_PROPERTY_SCOPE_USER,
		req.GetPropertyMutations(),
	)
	if err := s.validateSyncManagedProperties(
		ctx,
		properties.PropertyScope_PROPERTY_SCOPE_USER,
		req.GetPropertyMutations(),
	); err != nil {
		return nil, err
	}
	managedProps, propertyChanges, err := splitUserPropertyMutations(req)
	if err != nil {
		return nil, err
	}
	if err := propertyregistry.ValidateCustomPropertyChanges(
		ctx,
		properties.PropertyScope_PROPERTY_SCOPE_USER,
		propertyChanges,
		func(string) bool { return true },
		s.citizensStore.ResolvePropertyDefinitions,
	); err != nil {
		return nil, fmt.Errorf("failed to validate user property changes: %w", err)
	}

	changes := make([]citizensstore.PropertyChange, 0, len(propertyChanges))
	for _, change := range propertyChanges {
		mutation := findPropertyMutation(req.GetPropertyMutations(), change.Key)
		lifecycle, cause := userMutationActivityMetadata(mutation)
		if change.Value == nil {
			changes = append(changes, citizensstore.PropertyChange{
				Property:  &properties.Property{Key: change.Key},
				Delete:    true,
				Lifecycle: lifecycle,
				Cause:     cause,
			})
			continue
		}
		changes = append(changes, citizensstore.PropertyChange{
			Property:  &properties.Property{Key: change.Key, Value: change.Value},
			Lifecycle: lifecycle,
			Cause:     cause,
		})
	}

	var sourceUserID *int32
	if req.GetSourceUser() != nil && req.GetSourceUser().UserId != nil {
		id := req.GetSourceUser().GetUserId()
		sourceUserID = &id
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin user property update: %w", err)
	}
	defer tx.Rollback()

	activities := make([]*usersactivity.UserActivity, 0, len(changes))
	if managedProps != nil {
		current, err := s.citizensStore.GetUserProps(ctx, tx, req.GetUserId())
		if err != nil {
			return nil, fmt.Errorf("failed to get current user props: %w", err)
		}
		managedActivities, err := s.citizensStore.HandleUserPropsChanges(
			ctx, tx, current, managedProps, sourceUserID, req.GetReason(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to apply managed user property changes: %w", err)
		}
		activities = append(activities, managedActivities...)
		activities = append(activities, citizensstore.UserWantedMutationActivities(
			current, managedProps, req.GetPropertyMutations(), sourceUserID, req.GetReason(),
		)...)
		if err := s.citizensStore.ApplyWantedMutationState(
			ctx,
			tx,
			req.GetUserId(),
			req.GetPropertyMutations(),
		); err != nil {
			return nil, err
		}
	}
	propertyActivities, err := s.citizensStore.ApplyPropertyMutationsWithActivity(
		ctx,
		tx,
		req.GetUserId(),
		changes,
		req.GetPropertyMutations(),
		sourceUserID,
		req.GetReason(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to apply user property mutations: %w", err)
	}
	activities = append(activities, propertyActivities...)
	if err := usersactivity.CreateUserActivities(ctx, tx, activities...); err != nil {
		return nil, fmt.Errorf("failed to create user property activities: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit user property update: %w", err)
	}

	props, err := s.citizensStore.GetUserProps(ctx, s.db, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user props: %w", err)
	}
	if err := s.attachUserProperties(ctx, props); err != nil {
		return nil, err
	}
	return &pbsync.SetUserPropertiesResponse{UserProps: props}, nil
}

func findPropertyMutation(
	mutations []*properties.PropertyMutation,
	key string,
) *properties.PropertyMutation {
	for _, mutation := range mutations {
		if mutation != nil && mutation.GetKey() == key {
			return mutation
		}
	}
	return nil
}

func userMutationActivityMetadata(
	mutation *properties.PropertyMutation,
) (*usersactivity.PropertyLifecycleChange, usersactivity.PropertyChangeCause) {
	if mutation == nil || (mutation.GetSetValue() == nil && mutation.GetClearValue() == nil) {
		return nil, usersactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_UNSPECIFIED
	}
	lifecycle := propertyregistry.MutationLifecycleForMutation(mutation)
	if lifecycle == nil {
		return nil, usersactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_MANUAL
	}
	return &usersactivity.PropertyLifecycleChange{
		ExpiresAt: lifecycle.ExpiresAt, ExpiryCleared: lifecycle.ExpiryCleared,
		StateReason: lifecycle.StateReason, StateReasonCleared: lifecycle.StateReasonCleared,
	}, usersactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_MANUAL
}

func vehicleMutationActivities(
	mutations []*properties.PropertyMutation,
) map[string]*vehiclesactivity.PropertyChange {
	activities := make(map[string]*vehiclesactivity.PropertyChange, len(mutations))
	for _, mutation := range mutations {
		if mutation == nil || (mutation.GetSetValue() == nil && mutation.GetClearValue() == nil) ||
			isManagedMutation(properties.PropertyScope_PROPERTY_SCOPE_VEHICLE, mutation.GetKey()) {
			continue
		}
		activity := &vehiclesactivity.PropertyChange{
			Key:       mutation.GetKey(),
			Lifecycle: vehicleMutationLifecycle(mutation),
			Cause:     vehiclesactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_MANUAL,
		}
		activities[mutation.GetKey()] = activity
	}
	return activities
}

func vehicleMutationLifecycle(
	mutation *properties.PropertyMutation,
) *vehiclesactivity.PropertyLifecycleChange {
	lifecycle := propertyregistry.MutationLifecycleForMutation(mutation)
	if lifecycle == nil {
		return nil
	}
	return &vehiclesactivity.PropertyLifecycleChange{
		ExpiresAt:          lifecycle.ExpiresAt,
		ExpiryCleared:      lifecycle.ExpiryCleared,
		StateReason:        lifecycle.StateReason,
		StateReasonCleared: lifecycle.StateReasonCleared,
	}
}

func (s *Store) GetUserProps(
	ctx context.Context,
	req *pbsync.GetUserPropsRequest,
) (*pbsync.GetUserPropsResponse, error) {
	props, err := s.citizensStore.GetUserProps(ctx, s.db, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to get user props. %w", err)
	}
	if err := s.attachUserProperties(ctx, props); err != nil {
		return nil, err
	}
	return &pbsync.GetUserPropsResponse{
		UserProps: props,
	}, nil
}

func (s *Store) GetVehicleProps(
	ctx context.Context,
	req *pbsync.GetVehiclePropsRequest,
) (*pbsync.GetVehiclePropsResponse, error) {
	vehicleProps, err := s.vehiclesStore.GetProps(ctx, s.db, req.GetPlate())
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle props. %w", err)
	}
	if err := s.attachVehicleProperties(ctx, vehicleProps); err != nil {
		return nil, err
	}
	return &pbsync.GetVehiclePropsResponse{
		VehicleProps: vehicleProps,
	}, nil
}

// SetVehicleProperties applies managed and generic property changes from a
// trusted Sync client.
func (s *Store) SetVehicleProperties(
	ctx context.Context,
	req *pbsync.SetVehiclePropertiesRequest,
) (*pbsync.SetVehiclePropertiesResponse, error) {
	s.normalizeSyncWantedExpiries(
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		req.GetPropertyMutations(),
	)
	if err := s.validateSyncManagedProperties(
		ctx,
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		req.GetPropertyMutations(),
	); err != nil {
		return nil, err
	}
	managedProps, propertyChanges, updateWanted, err := splitVehiclePropertyMutations(req)
	if err != nil {
		return nil, err
	}
	if err := propertyregistry.ValidateCustomPropertyChanges(
		ctx,
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		propertyChanges,
		func(string) bool { return true },
		s.vehiclesStore.ResolvePropertyDefinitions,
	); err != nil {
		return nil, fmt.Errorf("failed to validate vehicle property changes: %w", err)
	}

	var creatorID *int32
	if req.GetSourceUser() != nil && req.GetSourceUser().UserId != nil {
		id := req.GetSourceUser().GetUserId()
		creatorID = &id
	}

	sourceJob := req.GetSourceUser().GetJob()
	if sourceJob == "" && creatorID != nil {
		sourceUser, err := s.citizensStore.GetUser(
			ctx,
			&pbcitizens.GetUserRequest{UserId: *creatorID},
			citizensstore.GetUserOptions{IncludeJob: true},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get source user: %w", err)
		}
		if sourceUser == nil || sourceUser.GetUser() == nil {
			return nil, fmt.Errorf("source user %d not found", *creatorID)
		}
		sourceJob = sourceUser.GetUser().GetJob()
	}

	props, err := s.vehiclesStore.UpdateProps(ctx, vehiclesstore.PropsUpdate{
		Props:             managedProps,
		CreatorID:         creatorID,
		CreatorJob:        sourceJob,
		Reason:            req.GetReason(),
		WantedStateReason: wantedStateReason(req.GetPropertyMutations()),
		WantedStateClearReason: wantedStateClearReason(
			req.GetPropertyMutations(),
		),
		PropertyChanges:    propertyChanges,
		PropertyMutations:  req.GetPropertyMutations(),
		PropertyActivities: vehicleMutationActivities(req.GetPropertyMutations()),
		UpdateWanted:       updateWanted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to apply vehicle property changes: %w", err)
	}
	if err := s.attachVehicleProperties(ctx, props); err != nil {
		return nil, err
	}
	return &pbsync.SetVehiclePropertiesResponse{VehicleProps: props}, nil
}

// normalizeSyncWantedExpiries gives explicit Sync wanted expiries the same
// configured upper bound as the authenticated Citizen and Vehicle APIs.
// clear_expiry intentionally remains a clear: the dedicated wanted
// housekeeper still enforces the maximum duration for no-expiry records.
func (s *Store) normalizeSyncWantedExpiries(
	scope properties.PropertyScope,
	mutations []*properties.PropertyMutation,
) {
	game := s.appCfg.Get().GetGame()
	var maxDuration time.Duration
	switch scope {
	case properties.PropertyScope_PROPERTY_SCOPE_USER:
		if !game.GetMaxWantedDurationUserEnabled() || game.GetMaxWantedDurationUser() == nil {
			return
		}
		maxDuration = game.GetMaxWantedDurationUser().AsDuration()
	case properties.PropertyScope_PROPERTY_SCOPE_VEHICLE:
		if !game.GetMaxWantedDurationVehicleEnabled() || game.GetMaxWantedDurationVehicle() == nil {
			return
		}
		maxDuration = game.GetMaxWantedDurationVehicle().AsDuration()
	default:
		return
	}

	propertyregistry.ClampWantedExpiry(scope, mutations, time.Now().UTC().Add(maxDuration))
}

func (s *Store) validateSyncManagedProperties(
	ctx context.Context,
	scope properties.PropertyScope,
	mutations []*properties.PropertyMutation,
) error {
	for _, mutation := range mutations {
		if mutation == nil || !propertyregistry.IsSystemProperty(scope, mutation.GetKey()) {
			continue
		}
		if _, ok := propertyregistry.ManagedMutation(scope, mutation.GetKey()); !ok &&
			mutation.GetIntegerMath() == nil {
			return fmt.Errorf("system property %q is not Sync-mutable", mutation.GetKey())
		}
		var (
			enabled bool
			err     error
		)
		switch scope {
		case properties.PropertyScope_PROPERTY_SCOPE_USER:
			enabled, err = s.citizensStore.IsSystemPropertyEnabled(ctx, mutation.GetKey())
		case properties.PropertyScope_PROPERTY_SCOPE_VEHICLE:
			enabled, err = s.vehiclesStore.IsSystemPropertyEnabled(ctx, mutation.GetKey())
		}
		if err != nil {
			return err
		}
		if !enabled {
			return fmt.Errorf("system property %q is disabled", mutation.GetKey())
		}
	}
	return nil
}

func wantedStateReason(mutations []*properties.PropertyMutation) *string {
	for _, mutation := range mutations {
		if mutation.GetKey() == "vehicle.wanted" && mutation.Reason != nil {
			return mutation.Reason
		}
	}
	return nil
}

func wantedStateClearReason(mutations []*properties.PropertyMutation) bool {
	for _, mutation := range mutations {
		if mutation.GetKey() == "vehicle.wanted" && mutation.GetClearReason() {
			return true
		}
	}
	return false
}

// Sync clients are trusted but still receive the same fully resolved
// properties resource as application clients. This is particularly important
// for lifecycle state, which is projected on Property rather than UserProps or
// VehicleProps' legacy fields.
func (s *Store) attachUserProperties(ctx context.Context, props *usersprops.UserProps) error {
	if props == nil {
		return nil
	}
	values, err := s.citizensStore.GetProperties(ctx, props.GetUserId())
	if err != nil {
		return fmt.Errorf("failed to get user properties: %w", err)
	}
	keys := make([]string, 0, len(values))
	for _, value := range values {
		keys = append(keys, value.GetKey())
	}
	catalog, err := s.citizensStore.ResolvePropertyCatalog(ctx, keys)
	if err != nil {
		return fmt.Errorf("failed to resolve user property definitions: %w", err)
	}
	props.Properties = propertyregistry.ResolvePropertiesResourceWithGroups(
		values,
		catalog.Definitions,
		catalog.Groups,
	)
	return nil
}

func (s *Store) attachVehicleProperties(
	ctx context.Context,
	props *vehiclesprops.VehicleProps,
) error {
	if props == nil {
		return nil
	}
	values, err := s.vehiclesStore.GetProperties(ctx, props.GetPlate())
	if err != nil {
		return fmt.Errorf("failed to get vehicle properties: %w", err)
	}
	keys := make([]string, 0, len(values))
	for _, value := range values {
		keys = append(keys, value.GetKey())
	}
	catalog, err := s.vehiclesStore.ResolvePropertyCatalog(ctx, keys)
	if err != nil {
		return fmt.Errorf("failed to resolve vehicle property definitions: %w", err)
	}
	props.Properties = propertyregistry.ResolvePropertiesResourceWithGroups(
		values,
		catalog.Definitions,
		catalog.Groups,
	)
	return nil
}

func splitUserPropertyMutations(
	req *pbsync.SetUserPropertiesRequest,
) (*usersprops.UserProps, []*propertyregistry.ScalarChange, error) {
	managed := &usersprops.UserProps{UserId: req.GetUserId()}
	var hasManaged bool
	changes := make([]*propertyregistry.ScalarChange, 0, len(req.GetPropertyMutations()))
	seen := make(map[string]struct{}, len(req.GetPropertyMutations()))
	for _, mutation := range req.GetPropertyMutations() {
		if err := propertyregistry.ValidateMutationShape(mutation); err != nil {
			return nil, nil, err
		}
		if _, ok := seen[mutation.GetKey()]; ok {
			return nil, nil, fmt.Errorf("duplicate property mutation %q", mutation.GetKey())
		}
		seen[mutation.GetKey()] = struct{}{}
		if managedDefinition, ok := propertyregistry.ManagedMutation(
			properties.PropertyScope_PROPERTY_SCOPE_USER,
			mutation.GetKey(),
		); ok {
			definition, _ := propertyregistry.SystemPropertyDefinition(
				properties.PropertyScope_PROPERTY_SCOPE_USER,
				mutation.GetKey(),
			)
			if err := propertyregistry.ValidateMutationState(definition, mutation); err != nil {
				return nil, nil, fmt.Errorf(
					"invalid state for property %q: %w",
					mutation.GetKey(),
					err,
				)
			}
			if err := managedDefinition.Apply(managed, mutation); err != nil {
				return nil, nil, err
			}
			hasManaged = true
			continue
		}
		if value := mutation.GetSetValue(); value != nil {
			changes = append(
				changes,
				&propertyregistry.ScalarChange{Key: mutation.GetKey(), Value: value},
			)
		} else if mutation.GetClearValue() != nil {
			changes = append(changes, &propertyregistry.ScalarChange{Key: mutation.GetKey()})
		} else if mutation.GetIntegerMath() != nil {
			// Resolved by the citizen store while holding the value row lock.
			continue
		} else if mutation.GetState() == nil && mutation.Reason == nil && mutation.ClearReason == nil {
			return nil, nil, fmt.Errorf("property %q has no mutation", mutation.GetKey())
		}
	}
	if !hasManaged {
		managed = nil
	}
	return managed, changes, nil
}

func splitVehiclePropertyMutations(
	req *pbsync.SetVehiclePropertiesRequest,
) (*vehiclesprops.VehicleProps, []*propertyregistry.ScalarChange, bool, error) {
	managed := &vehiclesprops.VehicleProps{Plate: req.GetPlate()}
	var hasManaged bool
	changes := make([]*propertyregistry.ScalarChange, 0, len(req.GetPropertyMutations()))
	seen := make(map[string]struct{}, len(req.GetPropertyMutations()))
	for _, mutation := range req.GetPropertyMutations() {
		if err := propertyregistry.ValidateMutationShape(mutation); err != nil {
			return nil, nil, false, err
		}
		if _, ok := seen[mutation.GetKey()]; ok {
			return nil, nil, false, fmt.Errorf("duplicate property mutation %q", mutation.GetKey())
		}
		seen[mutation.GetKey()] = struct{}{}
		if managedDefinition, ok := propertyregistry.ManagedMutation(
			properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			mutation.GetKey(),
		); ok {
			definition, _ := propertyregistry.SystemPropertyDefinition(
				properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
				mutation.GetKey(),
			)
			if err := propertyregistry.ValidateMutationState(definition, mutation); err != nil {
				return nil, nil, false, fmt.Errorf(
					"invalid state for property %q: %w",
					mutation.GetKey(),
					err,
				)
			}
			if err := managedDefinition.Apply(managed, mutation); err != nil {
				return nil, nil, false, err
			}
			hasManaged = true
			continue
		}
		if value := mutation.GetSetValue(); value != nil {
			changes = append(
				changes,
				&propertyregistry.ScalarChange{Key: mutation.GetKey(), Value: value},
			)
		} else if mutation.GetClearValue() != nil {
			changes = append(changes, &propertyregistry.ScalarChange{Key: mutation.GetKey()})
		} else if mutation.GetIntegerMath() != nil {
			continue
		} else if mutation.GetState() == nil && mutation.Reason == nil && mutation.ClearReason == nil {
			return nil, nil, false, fmt.Errorf("property %q has no mutation", mutation.GetKey())
		}
	}
	if !hasManaged {
		managed = &vehiclesprops.VehicleProps{Plate: req.GetPlate()}
	}
	return managed, changes, hasManaged, nil
}

func isManagedMutation(scope properties.PropertyScope, key string) bool {
	_, ok := propertyregistry.ManagedMutation(scope, key)
	return ok
}

func (s *Store) AddColleagueActivity(
	ctx context.Context,
	req *pbsync.AddColleagueActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	if activity := req.GetColleagueActivity(); activity != nil && req.GetSourceUser() != nil &&
		req.GetSourceUser().UserId != nil {
		sourceUserID := req.GetSourceUser().GetUserId()
		activity.SourceUserId = &sourceUserID
	}
	if err := s.createColleagueActivity(ctx, s.db, req.GetColleagueActivity()); err != nil {
		return nil, fmt.Errorf("failed to create jobs user activities. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddColleagueProps(
	ctx context.Context,
	req *pbsync.AddColleaguePropsRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.handleColleagueProps(
		ctx,
		req.GetColleagueProps(),
		req.GetSourceUser(),
	); err != nil {
		return nil, fmt.Errorf("failed to handle ColleagueProps activity. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddJobTimeclock(
	ctx context.Context,
	req *pbsync.AddJobTimeclockRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.handleTimeclockEntry(ctx, req.GetJobTimeclock()); err != nil {
		return nil, fmt.Errorf("failed to handle JobTimeclock activity. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddDispatch(
	ctx context.Context,
	req *pbsync.AddDispatchRequest,
) (*pbsync.AddActivityResponse, error) {
	if s.dispatches == nil {
		return &pbsync.AddActivityResponse{}, nil
	}

	dsp, err := s.dispatches.Create(ctx, req.GetDispatch())
	if err != nil {
		return nil, fmt.Errorf("failed to create dispatch. %w", err)
	}

	return &pbsync.AddActivityResponse{
		Id:        new(dsp.GetId()),
		CreatedAt: dsp.GetCreatedAt(),
	}, nil
}

func (s *Store) AddMarker(
	ctx context.Context,
	req *pbsync.AddMarkerRequest,
) (*pbsync.AddActivityResponse, error) {
	if s.livemapStore == nil {
		return &pbsync.AddActivityResponse{}, nil
	}

	marker := req.GetMarker()
	// Markers are temporary by default, so we set an expiration time of 24 hours if not provided.
	if marker.GetExpiresAt() == nil {
		marker.ExpiresAt = timestamp.New(time.Now().Add(24 * time.Hour))
	}

	markerId, err := s.livemapStore.CreateMarker(
		ctx,
		req.GetMarker(),
		req.GetMarker().CreatorId,
		req.GetMarker().GetJob(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create marker. %w", err)
	}

	marker, err = s.livemapStore.GetMarker(ctx, markerId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created marker (%d). %w", markerId, err)
	}

	return &pbsync.AddActivityResponse{
		Id:        new(markerId),
		CreatedAt: marker.GetCreatedAt(),
	}, nil
}

func (s *Store) DeleteMarker(
	ctx context.Context,
	req *pbsync.DeleteMarkerRequest,
) (*pbsync.DeleteDataResponse, error) {
	if s.livemapStore == nil {
		return &pbsync.DeleteDataResponse{}, nil
	}

	if err := s.livemapStore.DeleteMarker(ctx, req.GetId(), timestamp.Now()); err != nil {
		return nil, fmt.Errorf("failed to delete marker. %w", err)
	}

	return &pbsync.DeleteDataResponse{
		RowsAffected: 1,
	}, nil
}

func (s *Store) createUserActivities(
	ctx context.Context,
	tx qrm.DB,
	activities ...*usersactivity.UserActivity,
) error {
	if len(activities) == 0 {
		return nil
	}

	tUserActivity := table.FivenetUserActivity
	stmt := tUserActivity.
		INSERT(
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Reason,
			tUserActivity.Data,
		)
	for _, activity := range activities {
		stmt = stmt.VALUES(
			activity.SourceUserId,
			activity.GetTargetUserId(),
			activity.GetType(),
			activity.GetReason(),
			activity.GetData(),
		)
	}

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) createColleagueActivity(
	ctx context.Context,
	tx qrm.DB,
	activities ...*colleaguesactivity.ColleagueActivity,
) error {
	if len(activities) == 0 {
		return nil
	}

	tJobColleagueActivity := table.FivenetJobColleagueActivity
	stmt := tJobColleagueActivity.
		INSERT(
			tJobColleagueActivity.Job,
			tJobColleagueActivity.SourceUserID,
			tJobColleagueActivity.TargetUserID,
			tJobColleagueActivity.ActivityType,
			tJobColleagueActivity.Reason,
			tJobColleagueActivity.Data,
		)
	for _, activity := range activities {
		stmt = stmt.VALUES(
			activity.GetJob(),
			activity.SourceUserId,
			activity.GetTargetUserId(),
			activity.GetActivityType(),
			activity.GetReason(),
			activity.GetData(),
		)
	}

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) handleUserProps(
	ctx context.Context,
	data *activity.UserProps,
	sourceUser *pbsync.SourceUser,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction. %w", err)
	}
	defer tx.Rollback()

	reqP := data.GetProps()
	props, err := s.citizensStore.GetUserProps(ctx, tx, reqP.GetUserId())
	if err != nil {
		return fmt.Errorf("failed to get user props. %w", err)
	}

	if reqP.HasLabels() {
		labels, err := s.citizensStore.GetUserLabels(
			ctx,
			tx,
			mysql.AND(table.FivenetUserLabels.UserID.EQ(mysql.Int32(reqP.GetUserId()))),
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to get user props labels. %w", err)
		}

		props.Labels = labels
	}

	reason := ""
	if data.Reason != nil {
		reason = data.GetReason()
	}

	var sourceUserID *int32
	if sourceUser != nil && sourceUser.UserId != nil {
		value := sourceUser.GetUserId()
		sourceUserID = &value
	} else {
		sourceUserID = &reqP.UserId
	}

	activities, err := s.citizensStore.HandleUserPropsChanges(
		ctx,
		tx,
		props,
		reqP,
		sourceUserID,
		reason,
	)
	if err != nil {
		return fmt.Errorf("failed to handle user props changes. %w", err)
	}

	if data.Reason != nil && data.GetReason() != "" {
		if err := s.createUserActivities(ctx, tx, activities...); err != nil {
			return fmt.Errorf("failed to create user activities. %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. %w", err)
	}

	return nil
}

func (s *Store) handleColleagueProps(
	ctx context.Context,
	data *activity.ColleagueProps,
	sourceUser *pbsync.SourceUser,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction. %w", err)
	}
	defer tx.Rollback()

	input := data.GetProps()
	props, err := s.jobsStore.GetColleagueProps(
		ctx,
		tx,
		input.GetJob(),
		input.GetUserId(),
		nil,
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to get jobs user props. %w", err)
	}

	reason := ""
	if data.Reason != nil {
		reason = data.GetReason()
	}

	var sourceUserID *int32
	if sourceUser != nil && sourceUser.UserId != nil {
		value := sourceUser.GetUserId()
		sourceUserID = &value
	} else {
		sourceUserID = &input.UserId
	}

	activities, err := s.jobsStore.HandleColleaguePropsChanges(
		ctx,
		tx,
		props,
		input,
		input.GetJob(),
		sourceUserID,
		reason,
	)
	if err != nil {
		return fmt.Errorf("failed to handle jobs user props changes. %w", err)
	}

	if data.Reason != nil && data.GetReason() != "" {
		if err := s.createColleagueActivity(ctx, tx, activities...); err != nil {
			return fmt.Errorf("failed to create jobs user activities. %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. %w", err)
	}

	return nil
}

func (s *Store) handleTimeclockEntry(ctx context.Context, data *activity.TimeclockUpdate) error {
	tTimeClock := table.FivenetJobTimeclock

	if data.GetStart() {
		if err := s.startTimeclockEntry(ctx, data); err != nil {
			return err
		}
	} else {
		stmt := tTimeClock.
			UPDATE().
			SET(
				tTimeClock.SpentTime.SET(mysql.RawFloat("`spent_time` + CAST((TIMESTAMPDIFF(SECOND, `start_time`, `end_time`) / 3600) AS DECIMAL(10,2))")),
				tTimeClock.EndTime.SET(mysql.CURRENT_TIMESTAMP()),
			).
			WHERE(mysql.AND(
				tTimeClock.Job.EQ(mysql.String(data.GetJob())),
				tTimeClock.UserID.EQ(mysql.Int32(data.GetUserId())),
				tTimeClock.StartTime.IS_NOT_NULL(),
				tTimeClock.EndTime.IS_NULL()),
			).
			LIMIT(1)
		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return fmt.Errorf("failed to update timeclock entry. %w", err)
		}
	}

	return nil
}

func (s *Store) startTimeclockEntry(ctx context.Context, data *activity.TimeclockUpdate) error {
	tTimeClock := table.FivenetJobTimeclock

	stmt := tTimeClock.
		SELECT(tTimeClock.UserID, tTimeClock.Date, tTimeClock.EndTime).
		FROM(tTimeClock).
		WHERE(mysql.AND(tTimeClock.Job.EQ(mysql.String(data.GetJob())), tTimeClock.UserID.EQ(mysql.Int32(data.GetUserId())))).
		ORDER_BY(tTimeClock.Date.DESC()).
		LIMIT(1)

	dest := &jobstimeclock.TimeclockEntry{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query timeclock entry. %w", err)
		}
	}

	if dest.GetUserId() > 0 {
		return nil
	}

	updateStmt := tTimeClock.
		INSERT(
			tTimeClock.Job,
			tTimeClock.UserID,
			tTimeClock.Date,
		).
		VALUES(
			data.GetJob(),
			data.GetUserId(),
			mysql.CURRENT_DATE(),
		)
	if _, err := updateStmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return fmt.Errorf("failed to insert timeclock entry. %w", err)
		}
	}

	return nil
}

func (s *Store) handleUserUpdate(ctx context.Context, data *activity.UserUpdate) error {
	tUser := table.FivenetUser
	selectStmt := tUser.
		SELECT(tUser.ID).
		FROM(tUser).
		WHERE(tUser.ID.EQ(mysql.Int32(data.GetUserId()))).
		LIMIT(1)

	user := &users.User{}
	if err := selectStmt.QueryContext(ctx, s.db, user); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query user by ID. %w", err)
		}
		return nil
	}

	updateSets := []mysql.ColumnAssigment{}
	if data.Job != nil {
		updateSets = append(updateSets, tUser.Job.SET(mysql.String(data.GetJob())))
	}
	if data.JobGrade != nil {
		updateSets = append(updateSets, tUser.JobGrade.SET(mysql.Int32(data.GetJobGrade())))
	}
	if data.Firstname != nil {
		updateSets = append(updateSets, tUser.Firstname.SET(mysql.String(data.GetFirstname())))
	}
	if data.Lastname != nil {
		updateSets = append(updateSets, tUser.Lastname.SET(mysql.String(data.GetLastname())))
	}

	if len(updateSets) > 0 {
		stmt := tUser.
			UPDATE()
		if len(updateSets) > 1 {
			stmt = stmt.SET(updateSets[0], updateSets[1:])
		} else {
			stmt = stmt.SET(updateSets[0])
		}
		stmt = stmt.
			WHERE(
				tUser.ID.EQ(mysql.Int32(data.GetUserId())),
			).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return fmt.Errorf("failed to update user. %w", err)
		}
	}

	return nil
}
