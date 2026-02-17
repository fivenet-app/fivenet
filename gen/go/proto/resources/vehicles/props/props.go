package vehiclesprops

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (x *VehicleProps) HandleChanges(
	ctx context.Context,
	tx qrm.Executable,
	in *VehicleProps,
) error {
	tVehicleProps := table.FivenetVehiclesProps

	updateSets := []mysql.ColumnAssigment{}

	// Generate the update sets
	if in.Wanted != nil {
		updateSets = append(updateSets, tVehicleProps.Wanted.SET(mysql.Bool(in.GetWanted())))
		updateSets = append(
			updateSets,
			tVehicleProps.WantedReason.SET(mysql.String(in.GetWantedReason())),
		)
		if in.GetWanted() {
			in.WantedAt = timestamp.Now()
		} else {
			in.WantedAt = nil
		}
		updateSets = append(
			updateSets,
			tVehicleProps.WantedAt.SET(dbutils.TimestampToMySQL(in.WantedAt)),
		)
	} else {
		in.Wanted = x.Wanted
		in.WantedReason = x.WantedReason
		in.WantedAt = x.WantedAt
	}

	if len(updateSets) > 0 {
		stmt := tVehicleProps.
			INSERT(
				tVehicleProps.Plate,
				tVehicleProps.UpdatedAt,
				tVehicleProps.Wanted,
				tVehicleProps.WantedReason,
				tVehicleProps.WantedAt,
			).
			VALUES(
				in.GetPlate(),
				mysql.CURRENT_TIMESTAMP(),
				in.Wanted,
				in.WantedReason,
				in.WantedAt,
			).
			ON_DUPLICATE_KEY_UPDATE(
				updateSets...,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (x *VehicleProps) LoadFromDB(
	ctx context.Context,
	tx qrm.DB,
	plate string,
) error {
	tVehicleProps := table.FivenetVehiclesProps.AS("vehicle_props")

	stmt := tVehicleProps.
		SELECT(
			tVehicleProps.Plate,
			tVehicleProps.UpdatedAt,
			tVehicleProps.Wanted,
			tVehicleProps.WantedReason,
			tVehicleProps.WantedAt,
		).
		FROM(tVehicleProps).
		WHERE(
			tVehicleProps.Plate.EQ(mysql.String(plate)),
		).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, tx, x); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	if x.GetPlate() == "" {
		x.Plate = plate
	}

	return nil
}
