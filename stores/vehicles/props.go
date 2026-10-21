package vehiclesstore

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	vehiclesprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/props"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const vehicleWantedPropertyKey = "vehicle.wanted"

func handlePropsChanges(
	ctx context.Context,
	tx qrm.Executable,
	current *vehiclesprops.VehicleProps,
	in *vehiclesprops.VehicleProps,
) error {
	if in.Wanted == nil {
		in.Wanted = current.Wanted
		in.WantedReason = current.WantedReason
		in.SetWantedAt(current.GetWantedAt())
		in.SetWantedTill(current.GetWantedTill())
		return nil
	}

	normalizeWantedChange(current, in, "")
	tValues := table.FivenetVehiclePropertyValues
	if _, err := tValues.
		INSERT(tValues.Plate, tValues.PropertyKey, tValues.ValueType, tValues.ValueBool).
		VALUES(in.GetPlate(), vehicleWantedPropertyKey, propertyValueTypeBool, in.GetWanted()).
		ON_DUPLICATE_KEY_UPDATE(
			tValues.ValueType.SET(mysql.IntExp(mysql.RawString("VALUES(`value_type`)"))),
			tValues.ValueBool.SET(mysql.BoolExp(mysql.RawString("VALUES(`value_bool`)"))),
		).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	tState := table.FivenetVehiclePropertyState
	if !in.GetWanted() {
		_, err := tState.DELETE().WHERE(
			tState.Plate.EQ(mysql.String(in.GetPlate())).AND(
				tState.PropertyKey.EQ(mysql.String(vehicleWantedPropertyKey)),
			),
		).ExecContext(ctx, tx)
		return err
	}

	_, err := tState.
		INSERT(tState.Plate, tState.PropertyKey, tState.Reason, tState.ActiveAt, tState.ExpiresAt).
		VALUES(
			in.GetPlate(), vehicleWantedPropertyKey, in.WantedReason,
			timestampTime(in.GetWantedAt()), timestampTime(in.GetWantedTill()),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tState.Reason.SET(mysql.RawString("VALUES(`reason`)")),
			tState.ActiveAt.SET(mysql.TimestampExp(mysql.RawString("VALUES(`active_at`)"))),
			tState.ExpiresAt.SET(mysql.TimestampExp(mysql.RawString("VALUES(`expires_at`)"))),
		).
		ExecContext(ctx, tx)
	return err
}

func timestampTime(value *timestamp.Timestamp) any {
	if value == nil || value.GetTimestamp() == nil {
		return nil
	}
	return value.AsTime().UTC().Truncate(time.Millisecond)
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

	if isZeroTimestamp(in.GetWantedTill()) {
		in.ClearWantedTill()
		return
	}

	if in.GetWantedTill() == nil {
		in.SetWantedTill(current.GetWantedTill())
	}
}

func isZeroTimestamp(ts *timestamp.Timestamp) bool {
	return ts != nil && ts.GetTimestamp() != nil &&
		ts.GetTimestamp().GetSeconds() == 0 && ts.GetTimestamp().GetNanos() == 0
}
