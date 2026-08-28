package demo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (d *Demo) seedDemoCatalog(ctx context.Context) error {
	if err := d.upsertDemoJobsAndGrades(ctx); err != nil {
		return err
	}

	if err := d.upsertDemoLicenses(ctx); err != nil {
		return err
	}

	if err := d.upsertDemoJobProps(ctx); err != nil {
		return err
	}

	if err := d.upsertDemoCentrumSettings(ctx); err != nil {
		return err
	}
	if err := d.upsertDemoCentrumUnits(ctx); err != nil {
		return err
	}

	if err := d.upsertDemoLawbooks(ctx); err != nil {
		return err
	}
	if err := d.upsertDemoLaws(ctx); err != nil {
		return err
	}

	if err := d.upsertDemoTargetJobHighestGradeRolePerms(ctx); err != nil {
		return err
	}

	if d.perms == nil {
		return errors.New("failed to reload demo RBAC perms: permissions service is not available")
	}
	if err := d.perms.ReloadJob(ctx, d.targetJobName()); err != nil {
		return fmt.Errorf("failed to reload demo RBAC perms for job %s. %w", d.targetJobName(), err)
	}

	return nil
}

func (d *Demo) upsertDemoJobProps(ctx context.Context) error {
	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.QuickButtons,
			tJobProps.RadioFrequency,
			tJobProps.Motd,
		).
		VALUES(
			d.targetJobName(),
			`{"penaltyCalculator":true}`,
			d.randomDemoRadioFrequency(),
			d.randomDemoMotd(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tJobProps.QuickButtons.SET(mysql.RawString("VALUES(`quick_buttons`)")),
			tJobProps.RadioFrequency.SET(mysql.RawString("VALUES(`radio_frequency`)")),
			tJobProps.Motd.SET(mysql.RawString("VALUES(`motd`)")),
		)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo job props. %w", err)
	}

	return nil
}

func (d *Demo) upsertDemoCentrumSettings(ctx context.Context) error {
	stmt := tCentrumSettings.
		INSERT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Type,
			tCentrumSettings.Public,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
		).
		VALUES(
			d.targetJobName(),
			true,
			centrumsettings.CentrumType_CENTRUM_TYPE_DISPATCH,
			true,
			centrumsettings.CentrumMode_CENTRUM_MODE_MANUAL,
			centrumsettings.CentrumMode_CENTRUM_MODE_AUTO_ROUND_ROBIN,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCentrumSettings.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tCentrumSettings.Enabled.SET(mysql.RawBool("VALUES(`enabled`)")),
			tCentrumSettings.Type.SET(mysql.RawInt("VALUES(`type`)")),
			tCentrumSettings.Public.SET(mysql.RawBool("VALUES(`public`)")),
			tCentrumSettings.Mode.SET(mysql.RawInt("VALUES(`mode`)")),
			tCentrumSettings.FallbackMode.SET(mysql.RawInt("VALUES(`fallback_mode`)")),
		)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo centrum settings. %w", err)
	}

	return nil
}

func (d *Demo) upsertDemoCentrumUnits(ctx context.Context) error {
	stmt := tCentrumUnits.
		INSERT(
			tCentrumUnits.Job,
			tCentrumUnits.Name,
			tCentrumUnits.Initials,
			tCentrumUnits.Color,
			tCentrumUnits.Icon,
			tCentrumUnits.Description,
		)

	for _, unit := range demoSeedCentrumUnits {
		stmt = stmt.VALUES(
			d.targetJobName(),
			unit.Name,
			unit.Initials,
			unit.Color,
			unit.Icon,
			unit.Description,
		)
	}

	stmt = stmt.
		ON_DUPLICATE_KEY_UPDATE(
			tCentrumUnits.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tCentrumUnits.Name.SET(mysql.RawString("VALUES(`name`)")),
			tCentrumUnits.Initials.SET(mysql.RawString("VALUES(`initials`)")),
			tCentrumUnits.Color.SET(mysql.RawString("VALUES(`color`)")),
			tCentrumUnits.Icon.SET(mysql.RawString("VALUES(`icon`)")),
			tCentrumUnits.Description.SET(mysql.RawString("VALUES(`description`)")),
		)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo centrum units. %w", err)
	}

	return nil
}

func (d *Demo) randomDemoRadioFrequency() string {
	// 2-3 digit frequency range for simple demo radios.
	return strconv.Itoa(d.randIntN(980) + 20)
}

func (d *Demo) randomDemoMotd() string {
	motds := []string{
		"Stay sharp and keep comms clear.",
		"Report status changes on radio.",
		"Team first, paperwork right after.",
		"Check equipment before every shift.",
		"Safety and professionalism come first.",
		"Log incidents completely and on time.",
		"Treat civilians with respect.",
		"Keep channels clean and concise.",
	}

	return motds[d.randIntN(len(motds))]
}

func (d *Demo) upsertDemoJobsAndGrades(ctx context.Context) error {
	if len(demoSeedJobs) == 0 {
		return nil
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	delStmt := tJobs.
		DELETE().
		WHERE(mysql.AND(
			tJobs.DeletedAt.IS_NOT_NULL(),
			tJobs.DeletedAt.IS_NULL(),
		))
	if _, err := delStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to remove jobs before demo jobs upsert. %w", err)
	}

	stmt := tJobs.
		INSERT(
			tJobs.Name,
			tJobs.Label,
		)
	for _, job := range demoSeedJobs {
		stmt = stmt.VALUES(job.Name, job.Label)
	}
	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tJobs.Label.SET(mysql.RawString("VALUES(`label`)")),
	)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to upsert demo jobs. %w", err)
	}

	if err := d.upsertDemoJobGrades(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit upsert jobs transaction. %w", err)
	}

	return nil
}

func (d *Demo) upsertDemoJobGrades(ctx context.Context, tx *sql.Tx) error {
	if len(demoSeedJobGrades) == 0 {
		return nil
	}

	stmt := tJobsGrades.
		INSERT(
			tJobsGrades.JobName,
			tJobsGrades.Grade,
			tJobsGrades.Label,
		)
	for _, grade := range demoSeedJobGrades {
		stmt = stmt.VALUES(grade.JobName, grade.Grade, grade.Label)
	}
	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tJobsGrades.Grade.SET(mysql.RawInt("VALUES(`grade`)")),
		tJobsGrades.Label.SET(mysql.RawString("VALUES(`label`)")),
	)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to upsert demo job grades. %w", err)
	}

	return nil
}

func (d *Demo) upsertDemoLicenses(ctx context.Context) error {
	if len(demoSeedLicenses) == 0 {
		return nil
	}

	stmt := tLicenses.
		INSERT(
			tLicenses.Type,
			tLicenses.Label,
		)

	for _, license := range demoSeedLicenses {
		stmt = stmt.VALUES(license.Type, license.Label)
	}

	stmt = stmt.
		ON_DUPLICATE_KEY_UPDATE(
			tLicenses.Label.SET(mysql.RawString("VALUES(`label`)")),
		)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo licenses. %w", err)
	}

	return nil
}

func (d *Demo) upsertDemoLawbooks(ctx context.Context) error {
	if len(demoSeedLawbooks) == 0 {
		return nil
	}

	for _, lawbook := range demoSeedLawbooks {
		var existing struct{ ID int64 }
		if err := tLawbooks.
			SELECT(tLawbooks.ID.AS("id")).
			FROM(tLawbooks).
			WHERE(mysql.OR(
				tLawbooks.ID.EQ(mysql.Int32(lawbook.ID)),
				tLawbooks.Name.EQ(mysql.String(lawbook.Name)),
			)).
			LIMIT(1).
			QueryContext(ctx, d.db, &existing); err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to find demo lawbook. %w", err)
		}
		if existing.ID > 0 {
			if _, err := tLawbooks.UPDATE(
				tLawbooks.Name, tLawbooks.Description, tLawbooks.DeletedAt,
			).SET(lawbook.Name, lawbook.Description, mysql.NULL).
				WHERE(tLawbooks.ID.EQ(mysql.Int64(existing.ID))).LIMIT(1).ExecContext(ctx, d.db); err != nil {
				return fmt.Errorf("failed to update demo lawbook. %w", err)
			}
			continue
		}
		if _, err := tLawbooks.INSERT(tLawbooks.ID, tLawbooks.Name, tLawbooks.Description).
			VALUES(lawbook.ID, lawbook.Name, lawbook.Description).ExecContext(ctx, d.db); err != nil {
			return fmt.Errorf("failed to insert demo lawbook. %w", err)
		}
	}

	return nil
}

func (d *Demo) upsertDemoLaws(ctx context.Context) error {
	if len(demoSeedLaws) == 0 {
		return nil
	}

	for _, law := range demoSeedLaws {
		var existing struct{ ID int64 }
		if err := tLawbooksLaws.
			SELECT(tLawbooksLaws.ID.AS("id")).
			FROM(tLawbooksLaws).
			WHERE(mysql.OR(
				tLawbooksLaws.ID.EQ(mysql.Int32(law.ID)),
				mysql.AND(
					tLawbooksLaws.LawbookID.EQ(mysql.Int32(law.LawbookID)),
					tLawbooksLaws.Name.EQ(mysql.String(law.Name)),
				),
			)).
			LIMIT(1).
			QueryContext(ctx, d.db, &existing); err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to find demo law. %w", err)
		}
		if existing.ID > 0 {
			if _, err := tLawbooksLaws.UPDATE(
				tLawbooksLaws.LawbookID, tLawbooksLaws.Name, tLawbooksLaws.Description,
				tLawbooksLaws.Hint, tLawbooksLaws.Fine, tLawbooksLaws.DetentionTime,
				tLawbooksLaws.StvoPoints, tLawbooksLaws.DeletedAt,
			).SET(
				law.LawbookID, law.Name, law.Description, law.Hint, law.Fine,
				law.DetentionTime, law.StvoPoints, mysql.NULL,
			).WHERE(tLawbooksLaws.ID.EQ(mysql.Int64(existing.ID))).LIMIT(1).ExecContext(ctx, d.db); err != nil {
				return fmt.Errorf("failed to update demo law. %w", err)
			}
			continue
		}
		if _, err := tLawbooksLaws.INSERT(
			tLawbooksLaws.ID, tLawbooksLaws.LawbookID, tLawbooksLaws.Name, tLawbooksLaws.Description,
			tLawbooksLaws.Hint, tLawbooksLaws.Fine, tLawbooksLaws.DetentionTime, tLawbooksLaws.StvoPoints,
		).VALUES(
			law.ID, law.LawbookID, law.Name, law.Description, law.Hint, law.Fine,
			law.DetentionTime, law.StvoPoints,
		).ExecContext(ctx, d.db); err != nil {
			return fmt.Errorf("failed to insert demo law. %w", err)
		}
	}

	return nil
}
