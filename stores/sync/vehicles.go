package syncstore

import (
	"context"
	"fmt"

	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Store) SendVehicles(
	ctx context.Context,
	req *pbsync.SendVehiclesRequest,
) (*pbsync.SendDataResponse, error) {
	rowsAffected, err := s.handleVehiclesData(ctx, req.GetVehicles())
	if err != nil {
		return nil, fmt.Errorf("failed to handle vehicles data. %w", err)
	}

	return &pbsync.SendDataResponse{RowsAffected: rowsAffected}, nil
}

func (s *Store) DeleteVehicles(
	ctx context.Context,
	plates []string,
) (*pbsync.DeleteDataResponse, error) {
	if len(plates) == 0 {
		return &pbsync.DeleteDataResponse{}, nil
	}

	plateExprs := make([]mysql.Expression, 0, len(plates))
	for _, plate := range plates {
		plateExprs = append(plateExprs, mysql.String(plate))
	}

	tVehicles := table.FivenetOwnedVehicles
	delStmt := tVehicles.
		DELETE().
		WHERE(tVehicles.Plate.IN(plateExprs...)).
		LIMIT(int64(len(plates)))

	res, err := delStmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to execute vehicles delete statement. %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rows affected for vehicles delete. %w", err)
	}

	return &pbsync.DeleteDataResponse{RowsAffected: rows}, nil
}

func (s *Store) handleVehiclesData(
	ctx context.Context,
	data []*resourcesvehicles.Vehicle,
) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	tVehicles := table.FivenetOwnedVehicles
	tUsers := table.FivenetUser

	stmt := tVehicles.
		INSERT(
			tVehicles.UserID,
			tVehicles.Job,
			tVehicles.Plate,
			tVehicles.Model,
			tVehicles.Type,
			tVehicles.Data,
		)

	for _, vehicle := range data {
		var ownerId mysql.Expression
		if vehicle.OwnerId != nil && vehicle.GetOwnerId() != 0 {
			ownerId = mysql.Int32(vehicle.GetOwnerId())
		} else if vehicle.OwnerIdentifier != nil && vehicle.GetOwnerIdentifier() != "" {
			ownerId = tUsers.
				SELECT(tUsers.ID).
				WHERE(tUsers.Identifier.EQ(mysql.String(vehicle.GetOwnerIdentifier())))
		}

		stmt = stmt.VALUES(
			ownerId,
			vehicle.Job,
			vehicle.GetPlate(),
			vehicle.Model,
			vehicle.GetType(),
			mysql.NULL,
		)
	}

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tVehicles.UserID.SET(mysql.RawInt("VALUES(`user_id`)")),
		tVehicles.Job.SET(mysql.RawString("VALUES(`job`)")),
		tVehicles.Plate.SET(mysql.RawString("VALUES(`plate`)")),
		tVehicles.Model.SET(mysql.RawString("VALUES(`model`)")),
		tVehicles.Type.SET(mysql.RawString("VALUES(`type`)")),
		tVehicles.Data.SET(mysql.RawString("VALUES(`data`)")),
	)

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
