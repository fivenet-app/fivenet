package vehicles

import (
	context "context"
	"slices"

	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	permsvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	errorsvehicles "github.com/fivenet-app/fivenet/v2026/services/vehicles/errors"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListVehicleActivity(
	ctx context.Context,
	req *pbvehicles.ListVehicleActivityRequest,
) (*pbvehicles.ListVehicleActivityResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.vehicles.plate", req.GetPlate()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbvehicles.ListVehicleActivityResponse{
		Activity: []*vehiclesactivity.VehicleActivity{},
	}

	fields, err := permsvehicles.VehiclesService.ListVehicleActivity.FieldsTyped.Get(
		s.perms,
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	if !userInfo.GetJobAdmin() &&
		!fields.Contains(permsvehicles.VehiclesServiceListVehicleActivityFieldsPermValueOwn) {
		isOwner, err := s.store.IsVehicleOwner(ctx, req.GetPlate(), userInfo.GetUserId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
		}
		if isOwner {
			return resp, nil
		}
	}
	propertyList, err := s.perms.AttrPropertyListContext(ctx,
		userInfo,
		permsvehicles.VehiclesService.ListVehicles.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	visiblePropertyKeys, err := s.effectiveActivityPropertyKeys(
		ctx,
		propertyList.GetKeys(),
		userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	activityTypes := effectiveVehicleActivityTypes(
		req.GetTypes(),
		visiblePropertyKeys,
	)

	queryOpts := vehiclesstore.CountVehicleActivityOptions{
		VehicleActivityOptions: vehiclesstore.VehicleActivityOptions{
			Plate:               req.GetPlate(),
			Types:               activityTypes,
			PropertyKeys:        req.GetPropertyKeys(),
			VisiblePropertyKeys: visiblePropertyKeys,
			From:                req.GetFrom(),
			To:                  req.GetTo(),
		},
	}
	count, err := s.store.CountVehicleActivity(ctx, queryOpts)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, 20)
	resp.Pagination = pag
	if count <= 0 {
		return resp, nil
	}

	activity, err := s.store.ListVehicleActivity(ctx, vehiclesstore.ListVehicleActivityOptions{
		VehicleActivityOptions: queryOpts.VehicleActivityOptions,
		Sort:                   req.GetSort(),
		Offset:                 req.GetPagination().GetOffset(),
		Limit:                  limit,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	resp.Activity = activity

	canViewCreator := fields.Contains(
		permsvehicles.VehiclesServiceListVehicleActivityFieldsPermValueCreator,
	) || userInfo.GetJobAdmin()
	targets := make([]citizenshydrator.BasicTarget, 0, len(resp.GetActivity()))
	for i, entry := range resp.GetActivity() {
		if canViewCreator {
			targets = append(targets, citizenshydrator.BasicTarget{
				UserID: entry.GetCreatorId(),
				Set:    resp.Activity[i].SetCreator,
			})
		} else {
			entry.ClearCreatorId()
			entry.ClearCreator()
		}
	}

	if canViewCreator {
		hydrateShort := s.hydrator.HydrateBasicTargetsSafeFunc(userInfo)
		if err := hydrateShort(ctx, nil, targets); err != nil {
			return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
		}
	}

	return resp, nil
}

func effectiveVehicleActivityTypes(
	requested []vehiclesactivity.VehicleActivityType,
	visiblePropertyKeys []string,
) []vehiclesactivity.VehicleActivityType {
	activityTypes := []vehiclesactivity.VehicleActivityType{
		vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_WANTED,
		vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_FILE,
		vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_PROPERTY,
	}
	if len(requested) > 0 {
		activityTypes = requested
	}

	allowed := make([]vehiclesactivity.VehicleActivityType, 0, len(activityTypes))
	for _, activityType := range activityTypes {
		if activityType == vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_WANTED &&
			!slices.Contains(visiblePropertyKeys, "vehicle.wanted") {
			continue
		}
		if activityType == vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_FILE &&
			!slices.Contains(visiblePropertyKeys, vehicleFilesPropertyKey) {
			continue
		}
		allowed = append(allowed, activityType)
	}
	if len(allowed) == 0 {
		return []vehiclesactivity.VehicleActivityType{
			vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_UNSPECIFIED,
		}
	}
	return allowed
}

func (s *Server) effectiveActivityPropertyKeys(
	ctx context.Context,
	allowedKeys []string,
	propertyAdmin bool,
) ([]string, error) {
	keys, err := s.appendAdminVehiclePropertyKeys(
		ctx,
		propertyAdmin,
		append([]string(nil), allowedKeys...),
	)
	if err != nil {
		return nil, err
	}
	definitions, err := s.store.ResolvePropertyDefinitions(ctx, keys)
	if err != nil {
		return nil, err
	}
	return propertyregistry.EnabledPropertyKeys(keys, definitions), nil
}
