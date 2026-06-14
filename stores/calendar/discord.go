package calendar

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) GetCalendarReminderGuildID(ctx context.Context, job string) (string, error) {
	job = strings.TrimSpace(job)
	if job == "" {
		return "", errors.New("calendar job is required")
	}

	tJobProps := table.FivenetJobProps

	stmt := tJobProps.
		SELECT(tJobProps.DiscordGuildID.AS("discord_guild_id")).
		FROM(tJobProps).
		WHERE(mysql.AND(
			tJobProps.Job.EQ(mysql.String(job)),
			tJobProps.DeletedAt.IS_NULL(),
			tJobProps.DiscordGuildID.IS_NOT_NULL(),
		)).
		LIMIT(1)

	dest := struct {
		DiscordGuildID string
	}{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return "", errorscalendar.ErrNoDiscordGuildID
		}
		return "", err
	}

	if len(dest.DiscordGuildID) == 0 {
		return "", errorscalendar.ErrNoDiscordGuildID
	}

	guildID := strings.TrimSpace(dest.DiscordGuildID)
	if guildID == "" {
		return "", errorscalendar.ErrNoDiscordGuildID
	}

	return guildID, nil
}
