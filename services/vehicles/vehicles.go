package vehicles

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
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

	wantedEnabled, err := s.store.IsSystemPropertyEnabled(ctx, "vehicle.wanted")
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	propertyList, err := s.perms.AttrPropertyListContext(ctx,
		userInfo,
		permsvehicles.VehiclesService.ListVehicles.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	canAccessWanted := wantedEnabled &&
		propertyList.Contains("vehicle.wanted")
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
		IncludePropsUpdated: propertyList.Len() > 0,
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
	if req.GetReason() == "" {
		return nil, errorsvehicles.ErrReasonRequired
	}

	// Field Permission Check
	setPropertyList, err := s.perms.AttrPropertyListContext(ctx,
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

	// Generate the update sets
	if req.GetProps() != nil && req.GetProps().Wanted != nil {
		wantedEnabled, err := s.store.IsSystemPropertyEnabled(ctx, "vehicle.wanted")
		if err != nil {
			return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
		}
		if !wantedEnabled {
			return nil, errorsvehicles.ErrSystemPropertyDisabled
		}
		if !setPropertyList.Contains("vehicle.wanted") ||
			!viewPropertyList.Contains("vehicle.wanted") {
			return nil, errorsvehicles.ErrPropsWantedDenied
		}
	}
	if err := s.validateCustomProperties(
		ctx,
		req.GetPropertyChanges(),
		setPropertyList,
		viewPropertyList,
		propertyAdmin,
	); err != nil {
		return nil, err
	}
	// A protobuf epoch timestamp is the client-side marker for “permanent”.
	// Resolve it against the configured safety maximum before persisting it.
	if req.GetProps().GetWanted() && req.GetProps().GetWantedTill() != nil {
		wantedTill := req.GetProps().GetWantedTill()
		if wantedTill.GetTimestamp() != nil && wantedTill.GetTimestamp().GetSeconds() == 0 &&
			wantedTill.GetTimestamp().GetNanos() == 0 {
			game := s.appCfg.Get().GetGame()
			if maxDuration := game.GetMaxWantedDurationVehicle(); maxDuration != nil &&
				game.GetMaxWantedDurationVehicleEnabled() {
				req.Props.SetWantedTill(
					timestamp.New(time.Now().UTC().Add(maxDuration.AsDuration())),
				)
			}
		} else if wantedTill.GetTimestamp() != nil {
			game := s.appCfg.Get().GetGame()
			if maxDuration := game.GetMaxWantedDurationVehicle(); maxDuration != nil &&
				game.GetMaxWantedDurationVehicleEnabled() {
				maxTill := time.Now().UTC().Add(maxDuration.AsDuration())
				if wantedTill.AsTime().After(maxTill) {
					req.Props.SetWantedTill(timestamp.New(maxTill))
				}
			}
		}
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	creatorID := userInfo.GetUserId()
	props, err := s.store.UpdateProps(ctx, vehiclesstore.PropsUpdate{
		Props:           req.GetProps(),
		CreatorID:       &creatorID,
		CreatorJob:      userInfo.GetJob(),
		Reason:          req.GetReason(),
		PropertyChanges: req.GetPropertyChanges(),
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

func (s *Server) validateCustomProperties(
	ctx context.Context,
	values []*properties.PropertyChange,
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
	filters []*properties.PropertyFilter,
) error {
	for _, filter := range filters {
		if !propertyregistry.IsSearchablePropertyKey(
			properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			filter.GetKey(),
		) {
			return errorsvehicles.ErrInvalidProperty
		}
		if !userInfo.GetJobAdmin() && !allowed.Contains(filter.GetKey()) {
			return errorsvehicles.ErrPropsPropertiesDenied
		}
	}

	return nil
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
