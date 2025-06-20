package commands

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbuserinfo "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	permsjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs/perms"
	lang "github.com/fivenet-app/fivenet/v2025/i18n"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/timeutils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tAccsOauth2        = table.FivenetAccountsOauth2
	tAccs              = table.FivenetAccounts
	tColleagueProps    = table.FivenetJobColleagueProps
	tColleagueActivity = table.FivenetJobColleagueActivity
)

const absentDateFormat = "2006-01-02"

type AbsentCommand struct {
	l     *lang.I18n
	db    *sql.DB
	b     types.BotState
	perms perms.Permissions
}

func NewAbsentCommand(p CommandParams) (Command, error) {
	return &AbsentCommand{
		l:     p.L,
		db:    p.DB,
		b:     p.BotState,
		perms: p.Perms,
	}, nil
}

func (c *AbsentCommand) RegisterCommand(router *cmdroute.Router) api.CreateCommandData {
	router.Add("absent", c)

	lEN := c.l.Translator("en")
	lDE := c.l.Translator("de")

	return api.CreateCommandData{
		Type: discord.ChatInputCommand,
		Name: lEN("discord.commands.absent.name", nil),
		NameLocalizations: discord.StringLocales{
			discord.German: lDE("discord.commands.absent.name", nil),
		},
		Description: lEN("discord.commands.absent.desc", nil),
		DescriptionLocalizations: discord.StringLocales{
			discord.German: lDE("discord.commands.absent.desc", nil),
		},
		Options: discord.CommandOptions{
			&discord.StringOption{
				OptionName: lEN("discord.commands.absent.options.start-date.name", nil),
				OptionNameLocalizations: discord.StringLocales{
					discord.German: lDE("discord.commands.absent.options.start-date.name", nil),
				},
				Description: lEN("discord.commands.absent.options.start-date.desc", nil),
				Required:    true,
			},
			&discord.IntegerOption{
				OptionName: lEN("discord.commands.absent.options.days.name", nil),
				OptionNameLocalizations: discord.StringLocales{
					discord.German: lDE("discord.commands.absent.options.days.name", nil),
				},
				Description: lEN("discord.commands.absent.options.days.desc", nil),
				Required:    true,
				Min:         option.NewInt(1),
				Max:         option.NewInt(31),
			},
			&discord.StringOption{
				OptionName: lEN("discord.commands.absent.options.reason.name", nil),
				OptionNameLocalizations: discord.StringLocales{
					discord.German: lDE("discord.commands.absent.options.reason.name", nil),
				},
				Description: lEN("discord.commands.absent.options.reason.desc", nil),
				Required:    true,
				MinLength:   option.NewInt(3),
				MaxLength:   option.NewInt(200),
			},
		},
		DefaultMemberPermissions: discord.NewPermissions(discord.PermissionSendMessages),
	}
}

func (c *AbsentCommand) getBaseResponse() *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Flags: discord.EphemeralMessage,
		Embeds: &[]discord.Embed{
			{
				Type:  discord.LinkEmbed,
				Color: embeds.ColorError,
				Provider: &discord.EmbedProvider{
					Name: "FiveNet",
				},
				Thumbnail: &discord.EmbedThumbnail{
					URL:    "https://cdn.discordapp.com/app-icons/1101207666652618865/94429951df15108c737949ff2770cd8f.png",
					Width:  128,
					Height: 128,
				},
				Footer: embeds.EmbedFooterMadeBy,
			},
		},
	}
}

func (c *AbsentCommand) HandleCommand(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	localizer := c.l.Translator(string(cmd.Event.Locale))
	resp := c.getBaseResponse()

	if cmd.Event.GuildID == discord.NullGuildID {
		return nil
	}

	if cmd.Event.Member == nil {
		(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.wrong_discord.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.wrong_discord.desc", nil)
		(*resp.Embeds)[0].Color = embeds.ColorInfo
		return resp
	}

	job, ok := c.b.GetJobFromGuildID(cmd.Event.GuildID)
	if !ok {
		(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.wrong_discord.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.wrong_discord.desc", nil)
		(*resp.Embeds)[0].Color = embeds.ColorInfo
		return resp
	}

	userId, jobGrade, err := c.getUserIDByJobAndDiscordID(ctx, job, cmd.Event.Member.User.ID)
	if err != nil {
		(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.no_user_found.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.no_user_found.desc", nil)
		return resp
	}
	if userId <= 0 || jobGrade < 0 {
		(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.no_user_found.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.no_user_found.desc", nil)
		return resp
	}

	// For now just check if the user can set
	userInfo := &pbuserinfo.UserInfo{
		UserId:   userId,
		Job:      job,
		JobGrade: jobGrade,
	}
	if !c.perms.Can(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetColleaguePropsPerm) {
		(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.no_perms.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.no_perms.desc",
			map[string]any{"code": "perm"})
		return resp
	}

	startDateOption := cmd.Data.Options.Find("start-date")
	now := timeutils.StartOfDay(time.Now())
	startDate := time.Now()

	startDateValue := strings.ToLower(startDateOption.String())
	if startDateValue != "today" {
		parsed, err := time.Parse(absentDateFormat, startDateValue)
		if err != nil {
			(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.invalid_date.title", nil)
			(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.invalid_date.desc", nil)
			return resp
		}
		startDate = parsed

		if !(now.Equal(startDate) || startDate.After(now)) {
			(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.invalid_date.title", nil)
			(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.invalid_date.desc", nil)
			return resp
		}
	}

	daysOptions := cmd.Data.Options.Find("days")
	days, err := daysOptions.IntValue()
	if err != nil || days <= 0 {
		(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.failed.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.failed.desc",
			map[string]any{"code": "Days wrong"})
		return resp
	}
	endDate := startDate.AddDate(0, 0, int(days))

	reasonOption := cmd.Data.Options.Find("reason")
	reason := strings.TrimSpace(reasonOption.String())
	reason += " (via Discord Bot)"

	check, err := c.createAbsenceForUser(ctx, userId, job, startDate, endDate, reason)
	if err != nil {
		(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.failed.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.failed.desc",
			map[string]any{"code": "Internal Error"},
		)
		return resp
	}

	if !check {
		(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.success.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.success.desc",
			map[string]any{
				"absenceBegin": startDate.Format(absentDateFormat),
				"absenceEnd":   endDate.Format(absentDateFormat),
			},
		)
		(*resp.Embeds)[0].Color = embeds.ColorSuccess
	} else {
		(*resp.Embeds)[0].Title = localizer("discord.commands.absent.results.already_absent.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.absent.results.already_absent.desc",
			map[string]any{
				"absenceBegin": startDate.Format(absentDateFormat),
				"absenceEnd":   endDate.Format(absentDateFormat),
			},
		)
		(*resp.Embeds)[0].Color = embeds.ColorInfo
	}

	return resp
}

func (c *AbsentCommand) getUserIDByJobAndDiscordID(ctx context.Context, job string, discordId discord.UserID) (int32, int32, error) {
	tUsers := tables.User().AS("user")

	stmt := tAccsOauth2.
		SELECT(
			tUsers.ID.AS("user_id"),
			tUsers.JobGrade.AS("job_grade"),
		).
		FROM(
			tAccsOauth2.
				INNER_JOIN(tAccs,
					tAccs.ID.EQ(tAccsOauth2.AccountID),
				).
				INNER_JOIN(tUsers,
					jet.AND(
						tUsers.Identifier.LIKE(jet.CONCAT(jet.String("%"), tAccs.License)),
						tUsers.Job.EQ(jet.String(job)),
					)),
		).
		WHERE(jet.AND(
			tAccsOauth2.Provider.EQ(jet.String("discord")),
			tAccsOauth2.ExternalID.EQ(jet.String(discordId.String())),
		)).
		LIMIT(1)

	var dest struct {
		UserID   int32 `alias:"user_id"`
		JobGrade int32 `alias:"job_grade"`
	}
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		return -1, -1, err
	}

	return dest.UserID, dest.JobGrade, nil
}

func (c *AbsentCommand) createAbsenceForUser(ctx context.Context, charId int32, job string, absenceBegin time.Time, absenceEnd time.Time, reason string) (bool, error) {
	checkStmt := tColleagueProps.
		SELECT(
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
		).
		FROM(tColleagueProps).
		WHERE(jet.AND(
			tColleagueProps.UserID.EQ(jet.Int32(charId)),
			tColleagueProps.Job.EQ(jet.String(job)),
		)).
		LIMIT(1)

	props := jobs.ColleagueProps{}
	if err := checkStmt.QueryContext(ctx, c.db, &props); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	if props.AbsenceBegin != nil && props.AbsenceEnd != nil {
		begin := props.AbsenceBegin.AsTime()
		end := props.AbsenceEnd.AsTime()

		// Check if current absence is equal to the requested one
		if begin.Equal(absenceBegin) && end.Equal(absenceEnd) {
			return true, nil
		}
	}

	// Begin transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	stmt := tColleagueProps.
		INSERT(
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
		).
		VALUES(
			charId,
			job,
			absenceBegin,
			absenceEnd,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tColleagueProps.AbsenceBegin.SET(jet.DateExp(jet.Raw("VALUES(`absence_begin`)"))),
			tColleagueProps.AbsenceEnd.SET(jet.DateExp(jet.Raw("VALUES(`absence_end`)"))),
		)

	if _, err := stmt.ExecContext(ctx, c.db); err != nil {
		return false, err
	}

	activityStmt := tColleagueActivity.
		INSERT(
			tColleagueActivity.Job,
			tColleagueActivity.SourceUserID,
			tColleagueActivity.TargetUserID,
			tColleagueActivity.ActivityType,
			tColleagueActivity.Reason,
			tColleagueActivity.Data,
		).
		VALUES(
			job,
			charId,
			charId,
			jobs.ColleagueActivityType_COLLEAGUE_ACTIVITY_TYPE_ABSENCE_DATE,
			reason,
			&jobs.ColleagueActivityData{
				Data: &jobs.ColleagueActivityData_AbsenceDate{
					AbsenceDate: &jobs.AbsenceDateChange{
						AbsenceBegin: timestamp.New(absenceBegin),
						AbsenceEnd:   timestamp.New(absenceEnd),
					},
				},
			},
		)

	if _, err := activityStmt.ExecContext(ctx, tx); err != nil {
		return false, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return false, err
	}

	return false, nil
}
