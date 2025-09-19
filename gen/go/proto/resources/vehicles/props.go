package vehicles

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (x *VehicleProps) HandleChanges(ctx context.Context, tx *sql.Tx, in *VehicleProps) error {
	tUserProps := table.FivenetVehiclesProps

	updateSets := []mysql.ColumnAssigment{}

	// Generate the update sets
	if in.Wanted != nil {
		updateSets = append(updateSets, tUserProps.Wanted.SET(mysql.Bool(in.GetWanted())))
		updateSets = append(
			updateSets,
			tUserProps.WantedReason.SET(mysql.String(in.GetWantedReason())),
		)
	} else {
		in.Wanted = x.Wanted
		in.WantedReason = x.WantedReason
	}

	if len(updateSets) > 0 {
		stmt := tUserProps.
			INSERT(
				tUserProps.Plate,
				tUserProps.UpdatedAt,
				tUserProps.Wanted,
				tUserProps.WantedReason,
			).
			VALUES(
				in.GetPlate(),
				mysql.CURRENT_TIMESTAMP(),
				in.Wanted,
				in.WantedReason,
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
