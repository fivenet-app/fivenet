package vehicles

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	pbvehicles "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/vehicles"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorsvehicles "github.com/fivenet-app/fivenet/v2025/services/vehicles/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListVehicles(ctx context.Context, req *pbvehicles.ListVehiclesRequest) (*pbvehicles.ListVehiclesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	logRequest := false

	tVehicles := tables.OwnedVehicles().AS("vehicle")
	tUsers := tables.User().AS("user_short")

	condition := jet.Bool(true)
	userCondition := tUsers.Identifier.EQ(tVehicles.Owner)
	if req.LicensePlate != nil && *req.LicensePlate != "" {
		logRequest = true
		condition = jet.AND(condition, tVehicles.Plate.LIKE(jet.String(
			strings.ReplaceAll(*req.LicensePlate, "%", "")+"%",
		)))
	}

	// Make sure the model column is available
	modelColumn := s.customDB.Columns.Vehicle.GetModel(tVehicles.Alias())
	if modelColumn != nil && req.Model != nil && *req.Model != "" {
		logRequest = true
		condition = jet.AND(condition, tVehicles.Model.LIKE(jet.String(
			strings.ReplaceAll(*req.Model, "%", "")+"%",
		)))
	}

	if len(req.UserIds) > 0 {
		logRequest = true
		userIds := []jet.Expression{}
		for _, v := range req.UserIds {
			userIds = append(userIds, jet.Int32(v))
		}

		condition = jet.AND(condition,
			tUsers.Identifier.EQ(tVehicles.Owner),
			tUsers.ID.IN(userIds...),
		)
		userCondition = jet.AND(userCondition, tUsers.ID.IN(userIds...))
	} else if req.Job != nil && *req.Job != "" && !tables.ESXCompatEnabled {
		logRequest = true
		condition = jet.AND(condition,
			tVehicles.Job.EQ(jet.String(*req.Job)),
		)
	}

	if logRequest {
		defer s.aud.Log(&audit.AuditEntry{
			Service: pbvehicles.VehiclesService_ServiceDesc.ServiceName,
			Method:  "ListVehicles",
			UserId:  userInfo.UserId,
			UserJob: userInfo.Job,
			State:   audit.EventType_EVENT_TYPE_VIEWED,
		}, req)
	}

	countStmt := tVehicles.
		SELECT(
			jet.COUNT(tVehicles.Owner).AS("data_count.total"),
		).
		FROM(
			tVehicles.
				LEFT_JOIN(tUsers,
					userCondition,
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, 16)
	resp := &pbvehicles.ListVehiclesResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{
		tVehicles.Type.ASC(),
	}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "model":
			column = tVehicles.Model
		case "plate":
			fallthrough
		default:
			column = tVehicles.Plate
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tVehicles.Plate.ASC())
	}

	columns := dbutils.Columns{
		modelColumn,
		jet.REPLACE(tVehicles.Type, jet.String("_"), jet.String(" ")).AS("vehicle.type"),
		tUsers.ID.AS("vehicle.owner_id"),
		tUsers.ID,
		tUsers.Firstname,
		tUsers.Lastname,
		tUsers.Dateofbirth,
	}

	if !tables.ESXCompatEnabled {
		columns = append(columns,
			tVehicles.Job,
			tVehicles.Data,
		)
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permscitizens.CitizensServicePerm, permscitizens.CitizensServiceListCitizensPerm, permscitizens.CitizensServiceListCitizensFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	if fields.Contains("PhoneNumber") {
		columns = append(columns, tUsers.PhoneNumber)
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
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Vehicles); err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}

	for i := range resp.Vehicles {
		if resp.Vehicles[i].Job != nil && *resp.Vehicles[i].Job != "" {
			s.enricher.EnrichJobName(resp.Vehicles[i])
		}
	}

	resp.Pagination.Update(len(resp.Vehicles))

	return resp, nil
}
