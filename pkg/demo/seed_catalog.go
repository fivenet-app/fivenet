package demo

import (
	"context"
	"fmt"
	"strconv"

	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	"github.com/go-jet/jet/v2/mysql"
)

func (d *Demo) seedDemoCatalog(ctx context.Context) error {
	if err := d.upsertDemoJobs(ctx); err != nil {
		return err
	}
	if err := d.upsertDemoJobGrades(ctx); err != nil {
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

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
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

func (d *Demo) upsertDemoJobs(ctx context.Context) error {
	if len(demoSeedJobs) == 0 {
		return nil
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
		tJobs.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
	)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo jobs. %w", err)
	}

	return nil
}

func (d *Demo) upsertDemoJobGrades(ctx context.Context) error {
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
		tJobsGrades.Label.SET(mysql.RawString("VALUES(`label`)")),
	)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
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

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
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

	stmt := tLawbooks.
		INSERT(
			tLawbooks.ID,
			tLawbooks.Name,
			tLawbooks.Description,
		)

	for _, lawbook := range demoSeedLawbooks {
		stmt = stmt.VALUES(
			lawbook.ID,
			lawbook.Name,
			lawbook.Description,
		)
	}

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tLawbooks.Name.SET(mysql.RawString("VALUES(`name`)")),
		tLawbooks.Description.SET(mysql.RawString("VALUES(`description`)")),
	)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo lawbooks. %w", err)
	}

	return nil
}

func (d *Demo) upsertDemoLaws(ctx context.Context) error {
	if len(demoSeedLaws) == 0 {
		return nil
	}

	stmt := tLawbooksLaws.
		INSERT(
			tLawbooksLaws.ID,
			tLawbooksLaws.LawbookID,
			tLawbooksLaws.Name,
			tLawbooksLaws.Description,
			tLawbooksLaws.Hint,
			tLawbooksLaws.Fine,
			tLawbooksLaws.DetentionTime,
			tLawbooksLaws.StvoPoints,
		)

	for _, law := range demoSeedLaws {
		stmt = stmt.VALUES(
			law.ID,
			law.LawbookID,
			law.Name,
			law.Description,
			law.Hint,
			law.Fine,
			law.DetentionTime,
			law.StvoPoints,
		)
	}

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tLawbooksLaws.LawbookID.SET(mysql.RawInt("VALUES(`lawbook_id`)")),
		tLawbooksLaws.Name.SET(mysql.RawString("VALUES(`name`)")),
		tLawbooksLaws.Description.SET(mysql.RawString("VALUES(`description`)")),
		tLawbooksLaws.Hint.SET(mysql.RawString("VALUES(`hint`)")),
		tLawbooksLaws.Fine.SET(mysql.RawInt("VALUES(`fine`)")),
		tLawbooksLaws.DetentionTime.SET(mysql.RawInt("VALUES(`detention_time`)")),
		tLawbooksLaws.StvoPoints.SET(mysql.RawInt("VALUES(`stvo_points`)")),
	)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo laws. %w", err)
	}

	return nil
}
