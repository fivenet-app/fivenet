package jobsstore

import (
	"context"
	"errors"

	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) GetMOTD(ctx context.Context, db qrm.DB, job string) (string, error) {
	stmt := tJobProps.
		SELECT(
			tJobProps.Motd.AS("get_motd_response.motd"),
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(mysql.String(job)),
		).
		LIMIT(1)

	resp := &pbjobs.GetMOTDResponse{}
	if err := stmt.QueryContext(ctx, db, resp); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return "", err
		}
	}

	return resp.GetMotd(), nil
}

func (s *Store) SetMOTD(ctx context.Context, db qrm.DB, job string, motd string) error {
	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.Motd,
		).
		VALUES(
			job,
			motd,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.Motd.SET(mysql.String(motd)),
		)

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (s *Store) GetJobProps(
	ctx context.Context,
	db qrm.DB,
	job string,
) (*jobsprops.JobProps, error) {
	tJobProps := table.FivenetJobProps.AS("job_props")

	stmt := tJobProps.
		SELECT(
			tJobProps.Job,
			tJobProps.UpdatedAt,
			tJobProps.DeletedAt,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.DiscordGuildID,
			tJobProps.DiscordLastSync,
			tJobProps.DiscordSyncSettings,
			tJobProps.DiscordSyncChanges,
			tJobProps.LogoFileID,
			tJobProps.Settings,
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(mysql.String(job)),
		).
		LIMIT(1)

	dest := &jobsprops.JobProps{}
	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	dest.Default(job)

	return dest, nil
}
