package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) SendVehicles(
	ctx context.Context,
	req *pbsync.SendVehiclesRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())

	rowsAffected, err := s.handleVehiclesData(ctx, req.GetVehicles())
	if err != nil {
		return nil, fmt.Errorf("failed to handle vehicles data. %w", err)
	}

	return &pbsync.SendDataResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (s *Server) handleVehiclesData(
	ctx context.Context,
	data []*vehicles.Vehicle,
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
			// Use a subquery to find the user ID based on the owner identifier
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

	stmt = stmt.
		ON_DUPLICATE_KEY_UPDATE(
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

func (s *Server) DeleteVehicles(
	ctx context.Context,
	req *pbsync.DeleteVehiclesRequest,
) (*pbsync.DeleteDataResponse, error) {
	if len(req.GetPlates()) == 0 {
		return &pbsync.DeleteDataResponse{}, nil
	}

	plates := []mysql.Expression{}
	for _, plate := range req.GetPlates() {
		plates = append(plates, mysql.String(plate))
	}

	tVehicles := table.FivenetOwnedVehicles

	delStmt := tVehicles.
		DELETE().
		WHERE(tVehicles.Plate.IN(plates...)).
		LIMIT(int64(len(req.GetPlates())))

	res, err := delStmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to execute vehicles delete statement. %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rows affected for vehicles delete. %w", err)
	}

	return &pbsync.DeleteDataResponse{
		RowsAffected: rows,
	}, nil
}
