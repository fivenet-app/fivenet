package vehicles

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	permsvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsvehicles "github.com/fivenet-app/fivenet/v2026/services/vehicles/errors"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListVehicles(
	ctx context.Context,
	req *pbvehicles.ListVehiclesRequest,
) (*pbvehicles.ListVehiclesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	logRequest := false

	// Field Permission Check
	fields, err := permsvehicles.VehiclesService.SetVehicleProps.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	if req.GetLicensePlate() != "" {
		logRequest = true
	}

	if req.Model != nil && req.GetModel() != "" {
		logRequest = true
	}

	if len(req.GetUserIds()) > 0 {
		logRequest = true
	} else if req.Job != nil && req.GetJob() != "" {
		logRequest = true
	}

	canAccessWanted := fields.Contains(
		permsvehicles.VehiclesServiceSetVehiclePropsFieldsPermValueWanted,
	) ||
		userInfo.GetSuperuser()
	if canAccessWanted {
		if req.Wanted != nil && req.GetWanted() {
			logRequest = true
		}
	}

	if !logRequest {
		grpc_audit.Skip(ctx)
	}

	// Field Permission Check
	userFields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	query := vehiclesstore.ListQuery{
		LicensePlate:    req.GetLicensePlate(),
		Model:           req.GetModel(),
		UserIDs:         req.GetUserIds(),
		Job:             req.GetJob(),
		Wanted:          req.Wanted,
		CanFilterWanted: canAccessWanted,
		IncludePhoneNumber: userFields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		),
		IncludePropsUpdated: fields.Len() > 0,
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

	for i := range resp.GetVehicles() {
		if resp.Vehicles[i].Job != nil && resp.GetVehicles()[i].GetJob() != "" {
			s.enricher.EnrichJobName(resp.GetVehicles()[i])
		}
	}

	return resp, nil
}

func (s *Server) SetVehicleProps(
	ctx context.Context,
	req *pbvehicles.SetVehiclePropsRequest,
) (*pbvehicles.SetVehiclePropsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.vehicles.plate", req.GetProps().GetPlate()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Field Permission Check
	fields, err := permsvehicles.VehiclesService.SetVehicleProps.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	// Generate the update sets
	if req.GetProps() != nil && req.GetProps().Wanted != nil {
		if !fields.Contains(permsvehicles.VehiclesServiceSetVehiclePropsFieldsPermValueWanted) &&
			!userInfo.GetSuperuser() {
			return nil, errorsvehicles.ErrPropsWantedDenied
		}
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	props, err := s.store.UpdateProps(ctx, req.GetProps())
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	return &pbvehicles.SetVehiclePropsResponse{Props: props}, nil
}
