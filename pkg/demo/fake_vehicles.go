package demo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

const (
	demoVehiclePlatePrefix = "DV-"
)

var (
	demoVehicleTypes = []string{
		"aircraft",
		"boat",
		"car",
		"heli",
		"plane",
		"truck",
	}

	demoVehicleWantedReasons = []string{
		"Outstanding warrant linked to vehicle",
		"Suspected involvement in armed robbery",
		"Flagged in ongoing major-case investigation",
		"Vehicle matched BOLO description",
	}
)

type demoVehicleOwner struct {
	UserID     int32  `alias:"user_id"`
	Identifier string `alias:"identifier"`
	Job        string `alias:"job"`
}

func (d *Demo) seedFakeVehicles(ctx context.Context) error {
	owners, err := d.lookupVehicleOwners(ctx, int64(max(200, d.cfg.Demo.FakeUsers.Count*3)))
	if err != nil {
		return err
	}
	if len(owners) == 0 {
		return nil
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := d.clearDemoVehicles(ctx, tx); err != nil {
		return err
	}

	insertVehicles := tVehicles.
		INSERT(
			tVehicles.UserID,
			tVehicles.Plate,
			tVehicles.Model,
			tVehicles.Type,
			tVehicles.Data,
		)

	type wantedState struct {
		plate  string
		wanted bool
		at     any
		till   any
		reason any
	}
	wantedStates := make([]wantedState, 0)

	count := 0
	now := time.Now().UTC()
	for _, owner := range owners {
		vehicleCount := d.demoVehicleCountForOwner(owner)
		for slot := range vehicleCount {
			plate := d.demoVehiclePlate(owner, slot)
			insertVehicles = insertVehicles.VALUES(
				owner.UserID,
				plate,
				d.demoVehicleModel(owner, slot),
				d.demoVehicleType(owner, slot),
				`{"demoSeed":true}`,
			)

			wanted, wantedAt, wantedTill, wantedReason := d.demoVehicleWanted(now)
			wantedStates = append(wantedStates, wantedState{
				plate: plate, wanted: wanted, at: wantedAt, till: wantedTill, reason: wantedReason,
			})
			count++
		}
	}

	insertVehicles = insertVehicles.
		ON_DUPLICATE_KEY_UPDATE(
			tVehicles.UserID.SET(mysql.RawInt("VALUES(`user_id`)")),
			tVehicles.Model.SET(mysql.RawString("VALUES(`model`)")),
			tVehicles.Type.SET(mysql.RawString("VALUES(`type`)")),
			tVehicles.Data.SET(mysql.RawString("VALUES(`data`)")),
		)

	if _, err := insertVehicles.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to upsert demo vehicles. %w", err)
	}

	for _, state := range wantedStates {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO fivenet_vehicle_property_values (plate, property_key, value_type, value_bool)
VALUES (?, 'vehicle.wanted', 4, ?)
ON DUPLICATE KEY UPDATE value_type = VALUES(value_type), value_bool = VALUES(value_bool)`, state.plate, state.wanted); err != nil {
			return fmt.Errorf("failed to upsert demo vehicle wanted property. %w", err)
		}
		if !state.wanted {
			if _, err := tx.ExecContext(ctx, `
DELETE FROM fivenet_vehicle_property_state
WHERE plate = ? AND property_key = 'vehicle.wanted'`, state.plate); err != nil {
				return fmt.Errorf("failed to clear demo vehicle wanted state. %w", err)
			}
			continue
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO fivenet_vehicle_property_state (plate, property_key, reason, active_at, expires_at)
VALUES (?, 'vehicle.wanted', ?, ?, ?)
ON DUPLICATE KEY UPDATE reason = VALUES(reason), active_at = VALUES(active_at), expires_at = VALUES(expires_at)`, state.plate, state.reason, state.at, state.till); err != nil {
			return fmt.Errorf("failed to upsert demo vehicle wanted state. %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	d.logger.Info("completed demo vehicle seeding", zap.Int("count", count))
	return nil
}

func (d *Demo) lookupVehicleOwners(
	ctx context.Context,
	limit int64,
) ([]demoVehicleOwner, error) {
	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("demoVehicleOwner.user_id"),
			tUsers.Identifier.AS("demoVehicleOwner.identifier"),
			tUsers.Job.AS("demoVehicleOwner.job"),
		).
		FROM(tUsers).
		WHERE(tUsers.Identifier.LIKE(mysql.String("char%:%"))).
		ORDER_BY(tUsers.ID.ASC()).
		LIMIT(limit)

	owners := []demoVehicleOwner{}
	if err := stmt.QueryContext(ctx, d.db, &owners); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to lookup demo vehicle owners. %w", err)
		}
	}

	return owners, nil
}

func (d *Demo) clearDemoVehicles(ctx context.Context, tx *sql.Tx) error {
	vehiclesCondition := tVehicles.Plate.LIKE(mysql.String(demoVehiclePlatePrefix + "%"))
	var vehiclesCount database.DataCount
	vehiclesCountStmt := tVehicles.
		SELECT(mysql.COUNT(tVehicles.Plate).AS("data_count.total")).
		FROM(tVehicles).
		WHERE(vehiclesCondition)
	if err := vehiclesCountStmt.QueryContext(ctx, tx, &vehiclesCount); err != nil {
		return fmt.Errorf("failed to count demo vehicles. %w", err)
	}

	vehiclesStmt := tVehicles.
		DELETE().
		WHERE(vehiclesCondition).
		LIMIT(vehiclesCount.Total)
	if _, err := vehiclesStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to clear demo vehicles. %w", err)
	}

	return nil
}

func (d *Demo) demoVehicleCountForOwner(_ demoVehicleOwner) int {
	return d.fake.Number(1, 2)
}

func (d *Demo) demoVehiclePlate(_ demoVehicleOwner, _ int) string {
	token := strings.ToUpper(d.fake.Lexify("????"))
	nr := d.fake.Numerify("###")
	return fmt.Sprintf("%s%s-%s", demoVehiclePlatePrefix, nr, token)
}

func (d *Demo) demoVehicleModel(_ demoVehicleOwner, _ int) string {
	return strings.ToLower(strings.ReplaceAll(d.fake.CarModel(), " ", "_"))
}

func (d *Demo) demoVehicleType(_ demoVehicleOwner, _ int) string {
	return d.fake.RandomString(demoVehicleTypes)
}

func (d *Demo) demoVehicleWanted(
	now time.Time,
) (bool, *time.Time, *time.Time, *string) {
	if d.fake.Number(1, 10) != 1 {
		return false, nil, nil, nil
	}

	reason := d.fake.RandomString(demoVehicleWantedReasons)
	start := now.Add(-time.Duration(d.fake.Number(1, 48)) * time.Hour)
	end := start.Add(time.Duration(d.fake.Number(24, 96)) * time.Hour)
	return true, &start, &end, &reason
}
