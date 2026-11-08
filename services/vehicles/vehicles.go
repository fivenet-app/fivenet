package vehicles

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	permsvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	errorsvehicles "github.com/fivenet-app/fivenet/v2026/services/vehicles/errors"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListVehicles(
	ctx context.Context,
	req *pbvehicles.ListVehiclesRequest,
) (*pbvehicles.ListVehiclesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	logRequest := req.GetLicensePlate() != ""

	if req.Model != nil && req.GetModel() != "" {
		logRequest = true
	}

	if len(req.GetUserIds()) > 0 {
		logRequest = true
	} else if req.Job != nil && req.GetJob() != "" {
		logRequest = true
	}

	propertyList, err := s.perms.AttrPropertyListContext(ctx,
		userInfo,
		permsvehicles.VehiclesService.ListVehicles.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	definitionKeys := append([]string(nil), propertyList.GetKeys()...)
	definitionKeys = append(definitionKeys, requestedVehicleFilterPropertyKeys(req)...)
	definitions, err := s.store.ResolvePropertyDefinitions(ctx, definitionKeys)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	propertyAccess := propertyregistry.NewEffectivePropertyAccess(
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		propertyList,
		definitions,
		userInfo.GetJobAdmin(),
	)
	canAccessWanted := propertyAccess.CanView("vehicle.wanted")
	if req.Wanted != nil && req.GetWanted() && !canAccessWanted {
		return nil, errorsvehicles.ErrPropsPropertiesDenied
	}
	if canAccessWanted {
		if req.Wanted != nil && req.GetWanted() {
			logRequest = true
		}
	}
	propertyFilters := normalizeVehiclePropertyFilters(req, canAccessWanted)
	if len(propertyFilters) > 0 {
		if err := requireVehiclePropertyKeys(
			userInfo,
			propertyList,
			definitions,
			propertyFilters,
		); err != nil {
			return nil, err
		}
		logRequest = true
	}

	if !logRequest {
		grpc_audit.Skip(ctx)
	}

	query := vehiclesstore.ListQuery{
		UserInfo:            userInfo,
		LicensePlate:        req.GetLicensePlate(),
		Model:               req.GetModel(),
		UserIDs:             req.GetUserIds(),
		Job:                 req.GetJob(),
		PropertyFilters:     propertyFilters,
		CanFilterProperties: len(propertyFilters) > 0,
		IncludeWantedFields: canAccessWanted,
		Sort:                req.GetSort(),
	}

	total, err := s.store.Count(ctx, query)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(total, 20)
	resp := &pbvehicles.ListVehiclesResponse{
		Pagination: pag,
	}
	if total <= 0 {
		return resp, nil
	}

	query.Offset = req.GetPagination().GetOffset()
	query.Limit = limit
	resp.Vehicles, err = s.store.List(ctx, query)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	propertyAdmin := userInfo.GetJobAdmin()
	if propertyAdmin || propertyList.Len() > 0 {
		plates := make([]string, 0, len(resp.GetVehicles()))
		for _, vehicle := range resp.GetVehicles() {
			plates = append(plates, vehicle.GetPlate())
		}
		valuesByPlate, err := s.store.ListProperties(ctx, plates)
		if err != nil {
			return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
		}
		keys := append([]string(nil), propertyList.GetKeys()...)
		for _, values := range valuesByPlate {
			for _, value := range values {
				if propertyAdmin || propertyList.Contains(value.GetKey()) {
					keys = append(keys, value.GetKey())
				}
			}
		}
		keys, err = s.appendAdminVehiclePropertyKeys(ctx, propertyAdmin, keys)
		if err != nil {
			return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
		}
		catalog, err := s.store.ResolvePropertyCatalog(ctx, keys)
		if err != nil {
			return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
		}
		groups := catalog.Groups
		for _, vehicle := range resp.GetVehicles() {
			values := valuesByPlate[vehicle.GetPlate()]
			if !propertyAdmin {
				values = filterVehicleProperties(values, propertyList)
			}
			vehicle.Properties = propertyregistry.ResolvePropertiesResourceWithGroups(
				values,
				catalog.Definitions,
				groups,
			)
		}
	}

	for _, vehicle := range resp.GetVehicles() {
		if vehicle.Job != nil && vehicle.GetJob() != "" {
			s.enricher.EnrichJobName(vehicle)
		}
	}

	return resp, nil
}

func normalizeVehiclePropertyFilters(
	req *pbvehicles.ListVehiclesRequest,
	includeWanted bool,
) []*properties.PropertyFilter {
	filters := append([]*properties.PropertyFilter(nil), req.GetPropertyFilters()...)
	if includeWanted && req.Wanted != nil && req.GetWanted() {
		value := &properties.PropertyValue{}
		value.SetBoolValue(true)
		filters = append(filters, &properties.PropertyFilter{
			Key:       "vehicle.wanted",
			Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_EQUALS,
			Value:     value,
		})
	}
	return filters
}

func (s *Server) SetVehicleProps(
	ctx context.Context,
	req *pbvehicles.SetVehiclePropsRequest,
) (*pbvehicles.SetVehiclePropsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.vehicles.plate", req.GetProps().GetPlate()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	if req.GetReason() == "" && len(req.GetPropertyMutations()) == 0 {
		return nil, errorsvehicles.ErrReasonRequired
	}

	// Field Permission Check
	setPropertyList, err := s.perms.AttrPropertyList(
		userInfo,
		permsvehicles.VehiclesService.SetVehicleProps.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	viewPropertyList, err := s.perms.AttrPropertyListContext(ctx,
		userInfo,
		permsvehicles.VehiclesService.ListVehicles.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	propertyAdmin := userInfo.GetJobAdmin()
	wantedMutation := false

	// Wanted is a managed scalar. Public callers must use PropertyMutation;
	// these projection fields are populated only after that mutation is checked.
	if req.GetProps().Wanted != nil || req.GetProps().WantedAt != nil ||
		req.GetProps().WantedTill != nil || req.GetProps().WantedReason != nil {
		return nil, errorsvehicles.ErrInvalidProperty
	}
	game := s.appCfg.Get().GetGame()
	if maxDuration := game.GetMaxWantedDurationVehicle(); game.GetMaxWantedDurationVehicleEnabled() &&
		maxDuration != nil {
		propertyregistry.ClampWantedExpiry(
			properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			req.GetPropertyMutations(),
			time.Now().UTC().Add(maxDuration.AsDuration()),
		)
	}
	for _, mutation := range req.GetPropertyMutations() {
		managed, ok := propertyregistry.ManagedMutation(
			properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			mutation.GetKey(),
		)
		if !ok {
			continue
		}
		wantedMutation = true
		wantedEnabled, err := s.store.IsSystemPropertyEnabled(ctx, managed.Key)
		if err != nil {
			return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
		}
		if !wantedEnabled {
			return nil, errorsvehicles.ErrSystemPropertyDisabled
		}
		if !setPropertyList.Contains(managed.Key) || !viewPropertyList.Contains(managed.Key) {
			return nil, errorsvehicles.ErrPropsWantedDenied
		}
		if err := propertyregistry.ValidateMutationShape(mutation); err != nil {
			return nil, errorsvehicles.ErrInvalidProperty
		}
		definition, _ := propertyregistry.SystemPropertyDefinition(
			properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			managed.Key,
		)
		if err := propertyregistry.ValidateMutationState(definition, mutation); err != nil {
			return nil, errorsvehicles.ErrInvalidProperty
		}
		if err := propertyregistry.ValidateMutationReason(definition, mutation); err != nil {
			return nil, errorsvehicles.ErrInvalidProperty
		}
		if err := managed.Apply(req.Props, mutation); err != nil {
			return nil, errorsvehicles.ErrInvalidProperty
		}
	}
	if err := propertyregistry.ValidateMutationKeys(req.GetPropertyMutations()); err != nil {
		return nil, errorsvehicles.ErrInvalidProperty
	}
	mutationChanges, err := customMutationChanges(
		req.GetPropertyMutations(),
		setPropertyList,
		viewPropertyList,
		propertyAdmin,
	)
	if err != nil {
		return nil, err
	}
	propertyChanges := mutationChanges
	if err := s.validateCustomProperties(
		ctx,
		propertyChanges,
		setPropertyList,
		viewPropertyList,
		propertyAdmin,
	); err != nil {
		return nil, err
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	creatorID := userInfo.GetUserId()
	props, err := s.store.UpdateProps(ctx, vehiclesstore.PropsUpdate{
		Props:                  req.GetProps(),
		CreatorID:              &creatorID,
		CreatorJob:             userInfo.GetJob(),
		Reason:                 req.GetReason(),
		WantedStateReason:      vehicleWantedStateReason(req.GetPropertyMutations()),
		WantedStateClearReason: vehicleWantedStateClearReason(req.GetPropertyMutations()),
		PropertyChanges:        propertyChanges,
		PropertyMutations:      req.GetPropertyMutations(),
		PropertyActivities:     vehicleMutationActivities(req.GetPropertyMutations()),
		UpdateWanted:           wantedMutation,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	values, err := s.store.GetProperties(ctx, req.GetProps().GetPlate())
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	keys := append([]string(nil), viewPropertyList.GetKeys()...)
	for _, value := range values {
		keys = append(keys, value.GetKey())
	}
	keys, err = s.appendAdminVehiclePropertyKeys(ctx, propertyAdmin, keys)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	catalog, err := s.store.ResolvePropertyCatalog(ctx, keys)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	props.Properties = propertyregistry.ResolvePropertiesResourceWithGroups(
		filterVehiclePropertiesForAdmin(
			values,
			viewPropertyList,
			propertyAdmin,
		),
		catalog.Definitions,
		catalog.Groups,
	)

	return &pbvehicles.SetVehiclePropsResponse{Props: props}, nil
}

func vehicleWantedStateReason(mutations []*properties.PropertyMutation) *string {
	for _, mutation := range mutations {
		if mutation.GetKey() == "vehicle.wanted" && mutation.Reason != nil {
			return mutation.Reason
		}
	}
	return nil
}

func vehicleWantedStateClearReason(mutations []*properties.PropertyMutation) bool {
	for _, mutation := range mutations {
		if mutation.GetKey() == "vehicle.wanted" && mutation.GetClearReason() {
			return true
		}
	}
	return false
}

func customMutationChanges(
	mutations []*properties.PropertyMutation,
	settable *permissionsattributes.PropertyList,
	visible *permissionsattributes.PropertyList,
	propertyAdmin bool,
) ([]*propertyregistry.ScalarChange, error) {
	changes := make([]*propertyregistry.ScalarChange, 0, len(mutations))
	seen := make(map[string]struct{}, len(mutations))
	for _, mutation := range mutations {
		if err := propertyregistry.ValidateMutationShape(mutation); err != nil {
			return nil, errorsvehicles.ErrInvalidProperty
		}
		isMath := mutation.GetIntegerMath() != nil
		if !propertyregistry.IsCustomPropertyKey(
			properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			mutation.GetKey(),
		) && !isMath {
			if _, managed := propertyregistry.ManagedMutation(
				properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
				mutation.GetKey(),
			); managed {
				continue
			}
			return nil, errorsvehicles.ErrPropsPropertiesDenied
		}
		if !vehiclePropertyCanBeChanged(mutation.GetKey(), settable, visible, propertyAdmin) {
			return nil, errorsvehicles.ErrPropsPropertiesDenied
		}
		if _, ok := seen[mutation.GetKey()]; ok {
			return nil, errorsvehicles.ErrInvalidProperty
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
		} else if isMath {
			continue
		} else if mutation.GetState() == nil && mutation.Reason == nil && mutation.ClearReason == nil {
			return nil, errorsvehicles.ErrInvalidProperty
		}
	}
	return changes, nil
}

func (s *Server) validateCustomProperties(
	ctx context.Context,
	values []*propertyregistry.ScalarChange,
	settable *permissionsattributes.PropertyList,
	visible *permissionsattributes.PropertyList,
	propertyAdmin bool,
) error {
	err := propertyregistry.ValidateCustomPropertyChanges(
		ctx,
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		values,
		func(key string) bool { return vehiclePropertyCanBeChanged(key, settable, visible, propertyAdmin) },
		s.store.ResolvePropertyDefinitions,
	)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, propertyregistry.ErrInvalidPropertyChange):
		return errorsvehicles.ErrInvalidProperty
	case errors.Is(err, propertyregistry.ErrPropertyChangeDenied):
		return errorsvehicles.ErrPropsPropertiesDenied
	default:
		return errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
}

func vehicleMutationActivities(
	mutations []*properties.PropertyMutation,
) map[string]*vehiclesactivity.PropertyChange {
	activities := make(map[string]*vehiclesactivity.PropertyChange, len(mutations))
	for _, mutation := range mutations {
		if mutation == nil || (mutation.GetSetValue() == nil && mutation.GetClearValue() == nil) {
			continue
		}
		activities[mutation.GetKey()] = &vehiclesactivity.PropertyChange{
			Key:       mutation.GetKey(),
			Lifecycle: vehicleMutationLifecycle(mutation),
			Cause:     vehiclesactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_MANUAL,
		}
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

func vehiclePropertyCanBeChanged(
	key string,
	settable *permissionsattributes.PropertyList,
	visible *permissionsattributes.PropertyList,
	propertyAdmin bool,
) bool {
	return propertyAdmin || (settable.Contains(key) && visible.Contains(key))
}

func (s *Server) appendAdminVehiclePropertyKeys(
	ctx context.Context,
	propertyAdmin bool,
	keys []string,
) ([]string, error) {
	if !propertyAdmin {
		return keys, nil
	}

	for _, definition := range propertyregistry.SystemPropertyDefinitionsForScope(
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
	) {
		keys = append(keys, definition.GetKey())
	}
	customKeys, err := s.store.ListCustomPropertyKeys(ctx)
	if err != nil {
		return nil, err
	}
	keys = append(keys, customKeys...)
	slices.Sort(keys)
	return slices.Compact(keys), nil
}

func filterVehiclePropertiesForAdmin(
	values []*properties.Property,
	allowed *permissionsattributes.PropertyList,
	propertyAdmin bool,
) []*properties.Property {
	if propertyAdmin {
		return values
	}
	return filterVehicleProperties(values, allowed)
}

func requireVehiclePropertyKeys(
	userInfo *userinfo.UserInfo,
	allowed *permissionsattributes.PropertyList,
	definitions map[string]*properties.PropertyDefinition,
	filters []*properties.PropertyFilter,
) error {
	access := propertyregistry.NewEffectivePropertyAccess(
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		allowed,
		definitions,
		userInfo.GetJobAdmin(),
	)
	for _, filter := range filters {
		if errors.Is(
			access.ValidateFilter(filter.GetKey()),
			propertyregistry.ErrPropertyFilterNotSearchable,
		) {
			return errorsvehicles.ErrInvalidProperty
		}
		if !access.CanAccess(filter.GetKey()) {
			return errorsvehicles.ErrPropsPropertiesDenied
		}
	}

	return nil
}

func requestedVehicleFilterPropertyKeys(req *pbvehicles.ListVehiclesRequest) []string {
	keys := make([]string, 0, len(req.GetPropertyFilters())+1)
	for _, filter := range req.GetPropertyFilters() {
		keys = append(keys, filter.GetKey())
	}
	if req.GetWanted() {
		keys = append(keys, "vehicle.wanted")
	}
	return keys
}

func filterVehicleProperties(
	values []*properties.Property,
	allowed *permissionsattributes.PropertyList,
) []*properties.Property {
	filtered := make([]*properties.Property, 0, len(values))
	for _, property := range values {
		if allowed.Contains(property.GetKey()) {
			filtered = append(filtered, property)
		}
	}
	return filtered
}
