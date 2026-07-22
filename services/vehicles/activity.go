package vehicles

import (
	context "context"

	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	permsvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsvehicles "github.com/fivenet-app/fivenet/v2026/services/vehicles/errors"
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

	fields, err := permsvehicles.VehiclesService.ListVehicleActivity.FieldsTyped.Get(s.ps, userInfo)
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

	queryOpts := vehiclesstore.CountVehicleActivityOptions{
		VehicleActivityOptions: vehiclesstore.VehicleActivityOptions{
			Plate: req.GetPlate(),
			Types: req.GetTypes(),
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
		VehicleActivityOptions: vehiclesstore.VehicleActivityOptions{
			Plate: req.GetPlate(),
			Types: req.GetTypes(),
		},
		Sort:   req.GetSort(),
		Offset: req.GetPagination().GetOffset(),
		Limit:  limit,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	resp.Activity = activity

	canViewCreator := fields.Contains(
		permsvehicles.VehiclesServiceListVehicleActivityFieldsPermValueCreator,
	) || userInfo.GetJobAdmin()
	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetActivity() {
		if canViewCreator {
			if resp.GetActivity()[i].GetCreator() != nil {
				jobInfoFn(resp.GetActivity()[i].GetCreator())
			}
		} else {
			resp.Activity[i].CreatorId = nil
			resp.Activity[i].Creator = nil
		}
	}

	return resp, nil
}
