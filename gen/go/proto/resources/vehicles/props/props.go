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
		x.NormalizeWantedChange(in, "")
		updateSets = append(
			updateSets,
			tVehicleProps.Wanted.SET(mysql.Bool(in.GetWanted())),
			tVehicleProps.WantedReason.SET(dbutils.StringPP(in.WantedReason)),
			tVehicleProps.WantedAt.SET(dbutils.TimestampToMySQL(in.GetWantedAt())),
			tVehicleProps.WantedTill.SET(dbutils.TimestampToMySQL(in.GetWantedTill())),
		)
	} else {
		in.Wanted = x.Wanted
		in.WantedReason = x.WantedReason
		in.SetWantedAt(x.GetWantedAt())
		in.SetWantedTill(x.GetWantedTill())
	}

	if len(updateSets) > 0 {
		stmt := tVehicleProps.
			INSERT(
				tVehicleProps.Plate,
				tVehicleProps.UpdatedAt,
				tVehicleProps.Wanted,
				tVehicleProps.WantedReason,
				tVehicleProps.WantedAt,
				tVehicleProps.WantedTill,
			).
			VALUES(
				in.GetPlate(),
				mysql.CURRENT_TIMESTAMP(),
				in.Wanted,
				dbutils.StringEmpty(in.GetWantedReason()),
				in.GetWantedAt(),
				in.GetWantedTill(),
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

func (x *VehicleProps) NormalizeWantedChange(in *VehicleProps, reason string) {
	if in == nil {
		return
	}

	if in.Wanted == nil {
		in.SetWanted(x.GetWanted())
		in.SetWantedReason(x.GetWantedReason())
		in.SetWantedAt(x.GetWantedAt())
		in.SetWantedTill(x.GetWantedTill())
		return
	}

	if !in.GetWanted() {
		in.ClearWantedAt()
		in.ClearWantedTill()
		in.ClearWantedReason()
		return
	}

	if reason != "" {
		in.SetWantedReason(reason)
	} else if in.GetWantedReason() == "" {
		in.SetWantedReason(x.GetWantedReason())
	}

	if !x.GetWanted() {
		in.SetWantedAt(timestamp.Now())
	} else {
		in.SetWantedAt(x.GetWantedAt())
	}

	if in.GetWantedTill() == nil {
		in.SetWantedTill(x.GetWantedTill())
	}
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
			tVehicleProps.WantedTill,
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
		x.SetPlate(plate)
	}

	return nil
}
