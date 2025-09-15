package vehicles

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/vehicles"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	pbvehicles "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/vehicles"
	permsvehicles "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/vehicles/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsvehicles "github.com/fivenet-app/fivenet/v2025/services/vehicles/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListVehicles(
	ctx context.Context,
	req *pbvehicles.ListVehiclesRequest,
) (*pbvehicles.ListVehiclesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	logRequest := false

	tVehicles := tables.OwnedVehicles().AS("vehicle")
	tVehicleProps := table.FivenetVehiclesProps.AS("vehicle_props")
	tUsers := tables.User().AS("user_short")

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsvehicles.VehiclesServicePerm,
		permsvehicles.VehiclesServiceSetVehiclePropsPerm,
		permsvehicles.VehiclesServiceSetVehiclePropsFieldsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	condition := mysql.Bool(true)
	userCondition := tUsers.Identifier.EQ(tVehicles.Owner)
	if req.LicensePlate != nil && req.GetLicensePlate() != "" {
		logRequest = true
		condition = mysql.AND(condition, tVehicles.Plate.LIKE(mysql.String(
			strings.ReplaceAll(req.GetLicensePlate(), "%", "")+"%",
		)))
	}

	// Make sure the model column is available
	modelColumn := s.customDB.Columns.Vehicle.GetModel(tVehicles.Alias())
	if modelColumn != nil && req.Model != nil && req.GetModel() != "" {
		logRequest = true
		condition = mysql.AND(condition, tVehicles.Model.LIKE(mysql.String(
			strings.ReplaceAll(req.GetModel(), "%", "")+"%",
		)))
	}

	if len(req.GetUserIds()) > 0 {
		logRequest = true
		userIds := []mysql.Expression{}
		for _, v := range req.GetUserIds() {
			userIds = append(userIds, mysql.Int32(v))
		}

		condition = mysql.AND(condition,
			tUsers.Identifier.EQ(tVehicles.Owner),
			tUsers.ID.IN(userIds...),
		)
		userCondition = mysql.AND(userCondition, tUsers.ID.IN(userIds...))
	} else if req.Job != nil && req.GetJob() != "" && !tables.IsESXCompatEnabled() {
		logRequest = true
		condition = mysql.AND(condition,
			tVehicles.Job.EQ(mysql.String(req.GetJob())),
		)
	}

	if fields.Contains("Wanted") || userInfo.GetSuperuser() {
		if req.Wanted != nil && req.GetWanted() {
			logRequest = true
			condition = mysql.AND(condition,
				tVehicleProps.Wanted.EQ(mysql.Bool(req.GetWanted())),
			)
		}
	}

	if logRequest {
		defer s.aud.Log(&audit.AuditEntry{
			Service: pbvehicles.VehiclesService_ServiceDesc.ServiceName,
			Method:  "ListVehicles",
			UserId:  userInfo.GetUserId(),
			UserJob: userInfo.GetJob(),
			State:   audit.EventType_EVENT_TYPE_VIEWED,
		}, req)
	}

	countStmt := tVehicles.
		SELECT(
			mysql.COUNT(tVehicles.Owner).AS("data_count.total"),
		).
		FROM(
			tVehicles.
				LEFT_JOIN(tUsers,
					userCondition,
				).
				LEFT_JOIN(tVehicleProps,
					tVehicleProps.Plate.EQ(tVehicles.Plate),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 20)
	resp := &pbvehicles.ListVehiclesResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{
		tVehicles.Type.ASC(),
	}
	if req.GetSort() != nil {
		var column mysql.Column
		switch req.GetSort().GetColumn() {
		case "model":
			column = tVehicles.Model
		case "plate":
			fallthrough
		default:
			column = tVehicles.Plate
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tVehicles.Plate.ASC())
	}

	columns := dbutils.Columns{
		modelColumn,
		mysql.REPLACE(tVehicles.Type, mysql.String("_"), mysql.String(" ")).AS("vehicle.type"),
		tUsers.ID.AS("vehicle.owner_id"),
		tUsers.ID,
		tUsers.Firstname,
		tUsers.Lastname,
		tUsers.Dateofbirth,
		tVehicleProps.Plate,
	}

	if !tables.IsESXCompatEnabled() {
		columns = append(columns,
			tVehicles.Job,
			tVehicles.Data,
		)
	}

	// Field Permission Check
	userFields, err := s.ps.AttrStringList(
		userInfo,
		permscitizens.CitizensServicePerm,
		permscitizens.CitizensServiceListCitizensPerm,
		permscitizens.CitizensServiceListCitizensFieldsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	if userFields.Contains("PhoneNumber") {
		columns = append(columns, tUsers.PhoneNumber)
	}

	if fields.Len() > 0 {
		columns = append(columns, tVehicleProps.UpdatedAt)
	}
	if fields.Contains("Wanted") || userInfo.GetSuperuser() {
		columns = append(columns,
			tVehicleProps.Wanted,
			tVehicleProps.WantedReason,
		)
	}

	stmt := tVehicles.
		SELECT(
			tVehicles.Plate,
			columns.Get()...,
		).
		FROM(
			tVehicles.
				LEFT_JOIN(tUsers,
					userCondition,
				).
				LEFT_JOIN(tVehicleProps,
					tVehicleProps.Plate.EQ(tVehicles.Plate),
				),
		).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Vehicles); err != nil {
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

	auditEntry := &audit.AuditEntry{
		Service: pbvehicles.VehiclesService_ServiceDesc.ServiceName,
		Method:  "SetVehicleProps",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Get current vehicle props to be able to compare
	props, err := s.getVehicleProps(ctx, req.GetProps().GetPlate())
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	if props.Wanted == nil {
		wanted := false
		props.Wanted = &wanted
	}

	resp := &pbvehicles.SetVehiclePropsResponse{
		Props: &vehicles.VehicleProps{},
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsvehicles.VehiclesServicePerm,
		permsvehicles.VehiclesServiceSetVehiclePropsPerm,
		permsvehicles.VehiclesServiceSetVehiclePropsFieldsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	// Generate the update sets
	if req.Props.Wanted != nil {
		if !fields.Contains("Wanted") && !userInfo.GetSuperuser() {
			return nil, errorsvehicles.ErrPropsWantedDenied
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := props.HandleChanges(ctx, tx, req.GetProps()); err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	resp.Props, err = s.getVehicleProps(ctx, req.GetProps().GetPlate())
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	return resp, nil
}

func (s *Server) getVehicleProps(
	ctx context.Context,
	plate string,
) (*vehicles.VehicleProps, error) {
	tVehicleProps := table.FivenetVehiclesProps.AS("vehicle_props")

	stmt := tVehicleProps.
		SELECT(
			tVehicleProps.Plate,
			tVehicleProps.UpdatedAt,
			tVehicleProps.Wanted,
			tVehicleProps.WantedReason,
		).
		FROM(tVehicleProps).
		WHERE(
			tVehicleProps.Plate.EQ(mysql.String(plate)),
		).
		LIMIT(1)

	var dest vehicles.VehicleProps
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetPlate() == "" {
		dest.Plate = plate
	}

	return &dest, nil
}
