package sync

import (
	"context"
	"fmt"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) handleVehiclesData(
	ctx context.Context,
	data *pbsync.SendDataRequest_Vehicles,
) (int64, error) {
	if len(data.Vehicles.GetVehicles()) == 0 {
		return 0, nil
	}

	tVehicles := table.FivenetOwnedVehicles

	stmt := tVehicles.
		INSERT(
			tVehicles.UserID,
			tVehicles.Plate,
			tVehicles.Model,
			tVehicles.Type,
			tVehicles.Job,
			tVehicles.Data,
		)

	for _, vehicle := range data.Vehicles.GetVehicles() {
		var ownerId mysql.Expression
		if vehicle.OwnerIdentifier != nil && vehicle.GetOwnerIdentifier() != "" {
			ownerId = mysql.String(vehicle.GetOwnerIdentifier())
		} else if vehicle.OwnerId != nil {
			ownerId = mysql.Int32(vehicle.GetOwnerId())
		}

		stmt = stmt.VALUES(
			ownerId,
			vehicle.GetPlate(),
			vehicle.Model,
			vehicle.GetType(),
			vehicle.Job,
			mysql.NULL,
		)
	}

	assignments := []mysql.ColumnAssigment{
		tVehicles.UserID.SET(mysql.IntExp(mysql.Raw("VALUES(`user_id`)"))),
		tVehicles.Job.SET(mysql.StringExp(mysql.Raw("VALUES(`job`)"))),
		tVehicles.Plate.SET(mysql.StringExp(mysql.Raw("VALUES(`plate`)"))),
		tVehicles.Model.SET(mysql.StringExp(mysql.Raw("VALUES(`model`)"))),
		tVehicles.Type.SET(mysql.StringExp(mysql.Raw("VALUES(`type`)"))),
		tVehicles.Data.SET(mysql.StringExp(mysql.Raw("VALUES(`data`)"))),
	}

	stmt = stmt.
		ON_DUPLICATE_KEY_UPDATE(assignments...)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, fmt.Errorf("failed to execute vehicles insert statement. %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected for vehicles insert. %w", err)
	}

	return rowsAffected, nil
}
