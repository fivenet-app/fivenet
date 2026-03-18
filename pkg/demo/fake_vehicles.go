package demo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

	insertVehicles := tOwnedVehicles.
		INSERT(
			tOwnedVehicles.UserID,
			tOwnedVehicles.Plate,
			tOwnedVehicles.Model,
			tOwnedVehicles.Type,
			tOwnedVehicles.Data,
		)

	insertProps := tVehicleProps.
		INSERT(
			tVehicleProps.Plate,
			tVehicleProps.Wanted,
			tVehicleProps.WantedAt,
			tVehicleProps.WantedTill,
			tVehicleProps.WantedReason,
		)

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
			insertProps = insertProps.VALUES(
				plate,
				wanted,
				wantedAt,
				wantedTill,
				wantedReason,
			)
			count++
		}
	}

	insertVehicles = insertVehicles.ON_DUPLICATE_KEY_UPDATE(
		tOwnedVehicles.UserID.SET(mysql.RawInt("VALUES(`user_id`)")),
		tOwnedVehicles.Model.SET(mysql.RawString("VALUES(`model`)")),
		tOwnedVehicles.Type.SET(mysql.RawString("VALUES(`type`)")),
		tOwnedVehicles.Data.SET(mysql.RawString("VALUES(`data`)")),
	)

	if _, err := insertVehicles.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to upsert demo vehicles. %w", err)
	}

	insertProps = insertProps.ON_DUPLICATE_KEY_UPDATE(
		tVehicleProps.Wanted.SET(mysql.RawBool("VALUES(`wanted`)")),
		tVehicleProps.WantedAt.SET(mysql.RawTimestamp("VALUES(`wanted_at`)")),
		tVehicleProps.WantedTill.SET(mysql.RawTimestamp("VALUES(`wanted_till`)")),
		tVehicleProps.WantedReason.SET(mysql.RawString("VALUES(`wanted_reason`)")),
	)

	if _, err := insertProps.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to upsert demo vehicle props. %w", err)
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
	propsStmt := tVehicleProps.
		DELETE().
		WHERE(tVehicleProps.Plate.LIKE(mysql.String(demoVehiclePlatePrefix + "%")))
	if _, err := propsStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to clear demo vehicle props. %w", err)
	}

	vehiclesStmt := tOwnedVehicles.
		DELETE().
		WHERE(tOwnedVehicles.Plate.LIKE(mysql.String(demoVehiclePlatePrefix + "%")))
	if _, err := vehiclesStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to clear demo vehicles. %w", err)
	}

	return nil
}

func (d *Demo) demoVehicleCountForOwner(owner demoVehicleOwner) int {
	return d.fake.Number(1, 2)
}

func (d *Demo) demoVehiclePlate(owner demoVehicleOwner, slot int) string {
	token := strings.ToUpper(d.fake.Lexify("????"))
	nr := d.fake.Numerify("###")
	return fmt.Sprintf("%s%s-%s", demoVehiclePlatePrefix, nr, token)
}

func (d *Demo) demoVehicleJob(owner demoVehicleOwner) string {
	if owner.Job == "" || owner.Job == demoUnemployedJobName {
		return ""
	}
	return owner.Job
}

func (d *Demo) demoVehicleModel(owner demoVehicleOwner, slot int) string {
	return strings.ToLower(strings.ReplaceAll(d.fake.CarModel(), " ", "_"))
}

func (d *Demo) demoVehicleType(owner demoVehicleOwner, slot int) string {
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
