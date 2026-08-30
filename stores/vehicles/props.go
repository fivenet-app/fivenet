package vehiclesstore

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func handlePropsChanges(
	ctx context.Context,
	tx qrm.Executable,
	current *vehiclesprops.VehicleProps,
	in *vehiclesprops.VehicleProps,
) error {
	tVehicleProps := table.FivenetVehiclesProps
	updateSets := []mysql.ColumnAssigment{}

	if in.Wanted != nil {
		normalizeWantedChange(current, in, "")
		updateSets = append(updateSets,
			tVehicleProps.Wanted.SET(mysql.Bool(in.GetWanted())),
			tVehicleProps.WantedReason.SET(dbutils.StringPP(in.WantedReason)),
			tVehicleProps.WantedAt.SET(dbutils.TimestampToMySQL(in.GetWantedAt())),
			tVehicleProps.WantedTill.SET(dbutils.TimestampToMySQL(in.GetWantedTill())),
		)
	} else {
		in.Wanted = current.Wanted
		in.WantedReason = current.WantedReason
		in.SetWantedAt(current.GetWantedAt())
		in.SetWantedTill(current.GetWantedTill())
	}

	if len(updateSets) == 0 {
		return nil
	}

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
		ON_DUPLICATE_KEY_UPDATE(updateSets...)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func normalizeWantedChange(
	current *vehiclesprops.VehicleProps,
	in *vehiclesprops.VehicleProps,
	reason string,
) {
	if in == nil {
		return
	}

	if in.Wanted == nil {
		in.SetWanted(current.GetWanted())
		in.SetWantedReason(current.GetWantedReason())
		in.SetWantedAt(current.GetWantedAt())
		in.SetWantedTill(current.GetWantedTill())
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
		in.SetWantedReason(current.GetWantedReason())
	}

	if !current.GetWanted() {
		in.SetWantedAt(timestamp.Now())
	} else {
		in.SetWantedAt(current.GetWantedAt())
	}

	if in.GetWantedTill() == nil {
		in.SetWantedTill(current.GetWantedTill())
	}
}
