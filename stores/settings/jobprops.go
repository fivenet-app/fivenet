package settingsstore

import (
	"context"

	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

var tJobProps = table.FivenetJobProps

func (s *Store) SetJobProps(ctx context.Context, props *jobsprops.JobProps) error {
	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.DiscordGuildID,
			tJobProps.DiscordSyncSettings,
			tJobProps.Settings,
		).
		VALUES(
			props.GetJob(),
			props.GetLivemapMarkerColor(),
			dbutils.StringEmpty(props.GetRadioFrequency()),
			props.GetQuickButtons(),
			dbutils.StringEmpty(props.GetDiscordGuildId()),
			props.GetDiscordSyncSettings(),
			props.GetSettings(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.LivemapMarkerColor.SET(mysql.RawString("VALUES(`livemap_marker_color`)")),
			tJobProps.RadioFrequency.SET(mysql.RawString("VALUES(`radio_frequency`)")),
			tJobProps.QuickButtons.SET(mysql.RawString("VALUES(`quick_buttons`)")),
			tJobProps.DiscordGuildID.SET(mysql.RawString("VALUES(`discord_guild_id`)")),
			tJobProps.DiscordSyncSettings.SET(mysql.RawString("VALUES(`discord_sync_settings`)")),
			tJobProps.Settings.SET(mysql.RawString("VALUES(`settings`)")),
		)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}

func (s *Store) DeleteJobProps(ctx context.Context, job string) error {
	stmt := tJobProps.
		UPDATE().
		SET(
			tJobProps.DeletedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tJobProps.Job.EQ(mysql.String(job)),
		).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}
