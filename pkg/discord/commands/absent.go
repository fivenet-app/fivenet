package commands

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	permsjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/fivenet-app/fivenet/pkg/perms"
	timeutils "github.com/fivenet-app/fivenet/pkg/utils/time"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	tOauth2Accs       = table.FivenetOauth2Accounts
	tAccs             = table.FivenetAccounts
	tUsers            = table.Users.AS("user")
	tJobsUserProps    = table.FivenetJobsUserProps
	tJobsUserActivity = table.FivenetJobsUserActivity
)

func init() {
	CommandsFactories["absent"] = NewAbsentCommand
}

type AbsentCommand struct {
	l     *lang.I18n
	db    *sql.DB
	b     types.BotState
	perms perms.Permissions
}

func NewAbsentCommand(router *cmdroute.Router, cfg *config.Config, p CommandParams) (api.CreateCommandData, error) {
	lEN := p.L.I18n("en")
	lDE := p.L.I18n("de")

	router.Add("absent", &AbsentCommand{
		l:     p.L,
		db:    p.DB,
		b:     p.BotState,
		perms: p.Perms,
	})

	return api.CreateCommandData{
			Type: discord.ChatInputCommand,
			Name: lEN.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.name"}),
			NameLocalizations: discord.StringLocales{
				discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.name"}),
			},
			Description: lEN.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.desc"}),
			DescriptionLocalizations: discord.StringLocales{
				discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.desc"}),
			},
			Options: discord.CommandOptions{
				&discord.StringOption{
					OptionName: lEN.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.options.start-date.name"}),
					OptionNameLocalizations: discord.StringLocales{
						discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.options.start-date.name"}),
					},
					Description: lEN.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.options.start-date.desc"}),
					Required:    true,
				},
				&discord.IntegerOption{
					OptionName: lEN.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.options.days.name"}),
					OptionNameLocalizations: discord.StringLocales{
						discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.options.days.name"}),
					},
					Description: lEN.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.options.days.desc"}),
					Required:    true,
					Min:         option.NewInt(1),
					Max:         option.NewInt(31),
				},
				&discord.StringOption{
					OptionName: lEN.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.options.reason.name"}),
					OptionNameLocalizations: discord.StringLocales{
						discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.options.reason.name"}),
					},
					Description: lEN.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.options.reason.desc"}),
					Required:    true,
					MinLength:   option.NewInt(3),
					MaxLength:   option.NewInt(200),
				},
			},
		},
		nil
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
	localizer := c.l.I18n(string(cmd.Event.Locale))
	resp := c.getBaseResponse()

	if cmd.Event.GuildID == discord.NullGuildID {
		return nil
	}

	if cmd.Event.Member == nil {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.wrong_discord.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.wrong_discord.desc"})
		(*resp.Embeds)[0].Color = embeds.ColorInfo
		return resp
	}

	job, ok := c.b.GetJobFromGuildID(cmd.Event.GuildID)
	if !ok {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.wrong_discord.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.wrong_discord.desc"})
		(*resp.Embeds)[0].Color = embeds.ColorInfo
		return resp
	}

	userId, jobGrade, err := c.getUserIDByJobAndDiscordID(ctx, job, cmd.Event.Member.User.ID)
	if err != nil {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.no_user_found.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.no_user_found.desc"})
		return resp
	}
	if userId <= 0 || jobGrade < 0 {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.no_user_found.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.no_user_found.desc"})
		return resp
	}

	// For now just check if the user can set
	userInfo := &userinfo.UserInfo{
		UserId:   userId,
		Job:      job,
		JobGrade: jobGrade,
	}
	if !c.perms.Can(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetJobsUserPropsPerm) {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.no_perms.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "discord.commands.absent.results.no_perms.desc",
			TemplateData: map[string]string{
				"Code": "perm",
			},
		})
		return resp
	}

	startDateOption := cmd.Data.Options.Find("start-date")
	startDate := time.Now()

	startDateValue := strings.ToLower(startDateOption.String())
	if startDateValue != "today" {
		parsed, err := time.Parse("2006-01-02", startDateValue)
		if err != nil {
			(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.invalid_date.title"})
			(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.invalid_date.desc"})
			return resp
		}
		startDate = parsed

		now := timeutils.TruncateToDay(time.Now())
		if !startDate.After(now) {
			(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.invalid_date.title"})
			(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.invalid_date.desc"})
			return resp
		}
	}

	daysOptions := cmd.Data.Options.Find("days")
	days, err := daysOptions.IntValue()
	if err != nil {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.failed.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.failed.desc"})
		return resp
	}
	endDate := time.Now().AddDate(0, 0, int(days))

	reasonOption := cmd.Data.Options.Find("reason")
	reason := reasonOption.String()
	reason += " (via Discord Bot)"

	if err := c.createAbsenceForUser(ctx, userId, job, timestamp.New(startDate), timestamp.New(endDate), reason); err != nil {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.failed.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.failed.desc"})
		return resp
	}

	(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.absent.results.success.title"})
	(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "discord.commands.absent.results.success.desc",
		TemplateData: map[string]string{
			"AbsenceBegin": startDate.Format("2006-01-02"),
			"AbsenceEnd":   endDate.Format("2006-01-02"),
		},
	})
	(*resp.Embeds)[0].Color = embeds.ColorSuccess
	return resp
}

func (c *AbsentCommand) getUserIDByJobAndDiscordID(ctx context.Context, job string, discordId discord.UserID) (int32, int32, error) {
	stmt := tOauth2Accs.
		SELECT(
			tUsers.ID.AS("user_id"),
			tUsers.JobGrade.AS("job_grade"),
		).
		FROM(
			tOauth2Accs.
				INNER_JOIN(tAccs,
					tAccs.ID.EQ(tOauth2Accs.AccountID),
				).
				INNER_JOIN(tUsers,
					jet.AND(
						tUsers.Identifier.LIKE(jet.CONCAT(jet.String("%"), tAccs.License)),
						tUsers.Job.EQ(jet.String(job)),
					)),
		).
		WHERE(jet.AND(
			tOauth2Accs.Provider.EQ(jet.String("discord")),
			tOauth2Accs.ExternalID.EQ(jet.String(discordId.String())),
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

func (c *AbsentCommand) createAbsenceForUser(ctx context.Context, charId int32, job string, absenceBegin *timestamp.Timestamp, absenceEnd *timestamp.Timestamp, reason string) error {
	// Begin transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	updateSets := []jet.ColumnAssigment{
		tJobsUserProps.AbsenceBegin.SET(jet.DateExp(jet.Raw("VALUES(`absence_begin`)"))),
		tJobsUserProps.AbsenceEnd.SET(jet.DateExp(jet.Raw("VALUES(`absence_end`)"))),
	}
	stmt := tJobsUserProps.
		INSERT(
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.AbsenceBegin,
			tJobsUserProps.AbsenceEnd,
		).
		VALUES(
			charId,
			job,
			absenceBegin,
			absenceEnd,
		).
		ON_DUPLICATE_KEY_UPDATE(
			updateSets...,
		)

	if _, err := stmt.ExecContext(ctx, c.db); err != nil {
		return err
	}

	activityStmt := tJobsUserActivity.
		INSERT(
			tJobsUserActivity.Job,
			tJobsUserActivity.SourceUserID,
			tJobsUserActivity.TargetUserID,
			tJobsUserActivity.ActivityType,
			tJobsUserActivity.Reason,
			tJobsUserActivity.Data,
		).
		VALUES(
			job,
			charId,
			charId,
			jobs.JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_ABSENCE_DATE,
			reason,
			&jobs.JobsUserActivityData{
				Data: &jobs.JobsUserActivityData_AbsenceDate{
					AbsenceDate: &jobs.ColleagueAbsenceDate{
						AbsenceBegin: absenceBegin,
						AbsenceEnd:   absenceEnd,
					},
				},
			},
		)

	if _, err := activityStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
