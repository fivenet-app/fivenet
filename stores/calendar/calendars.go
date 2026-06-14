package calendarstore

import (
	"context"

	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CreateCalendar(
	ctx context.Context,
	tx qrm.DB,
	cal *calendarresource.Calendar,
	userInfo *userinfo.UserInfo,
	discordSettingsJSON *string,
) (int64, error) {
	stmt := tCalendar.
		INSERT(
			tCalendar.Job,
			tCalendar.DiscordSettings,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendar.CreatorID,
			tCalendar.CreatorJob,
		).
		VALUES(
			cal.GetJob(),
			discordSettingsJSON,
			cal.GetName(),
			cal.GetDescription(),
			cal.GetPublic(),
			cal.GetClosed(),
			cal.GetColor(),
			userInfo.GetUserId(),
			userInfo.GetJob(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCalendar.DiscordSettings.SET(mysql.RawString("VALUES(`discord_settings`)")),
			tCalendar.Name.SET(mysql.String(cal.GetName())),
			tCalendar.Description.SET(mysql.String("VALUES(`description`)")),
			tCalendar.Public.SET(mysql.Bool(cal.GetPublic())),
			tCalendar.Closed.SET(mysql.Bool(cal.GetClosed())),
			tCalendar.Color.SET(mysql.String(cal.GetColor())),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	if cal.GetId() == 0 {
		lastID, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}
		return lastID, nil
	}

	return cal.GetId(), nil
}

func (s *Store) UpdateCalendar(
	ctx context.Context,
	tx qrm.DB,
	cal *calendarresource.Calendar,
	discordSettingsJSON *string,
) error {
	discordSettingsValue := mysql.StringExp(mysql.NULL)
	if discordSettingsJSON != nil {
		discordSettingsValue = mysql.String(*discordSettingsJSON)
	}

	stmt := tCalendar.
		UPDATE(
			tCalendar.DiscordSettings,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
		).
		SET(
			discordSettingsValue,
			cal.GetName(),
			cal.GetDescription(),
			cal.GetPublic(),
			cal.GetClosed(),
			cal.GetColor(),
		).
		WHERE(mysql.AND(
			tCalendar.ID.EQ(mysql.Int64(cal.GetId())),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) DeleteCalendar(
	ctx context.Context,
	tx qrm.DB,
	calendarID int64,
	deletedAt *timestamp.Timestamp,
) error {
	stmt := tCalendar.
		UPDATE(
			tCalendar.DeletedAt,
		).
		SET(
			tCalendar.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAt)),
		).
		WHERE(tCalendar.ID.EQ(mysql.Int64(calendarID))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
