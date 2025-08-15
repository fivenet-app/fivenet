package modules

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/embeds"
	discordtypes "github.com/fivenet-app/fivenet/v2025/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
)

const (
	qualificationsRoleModulePrefix = "Qualifications-"
)

var (
	tQualifications        = table.FivenetQualifications
	tQualificationsResults = table.FivenetQualificationsResults
)

type QualificationsSync struct {
	*BaseModule
}

type qualificationsEntry struct {
	ID                 int64  `alias:"qualifications_entry.id"`
	Abbreviation       string `alias:"qualifications_entry.abbreviation"`
	QualificationTitle string `alias:"qualifications_entry.title"`

	DiscordSettings *qualifications.QualificationDiscordSettings `alias:"qualifications_entry.discord_settings"`
	RoleName        string
}

type qualificationUserMapping struct {
	AccountID  int64  `alias:"account_id"`
	ExternalID string `alias:"external_id"`
	Job        string `alias:"job"`
}

func init() {
	Modules["qualifications"] = NewQualifications
}

func NewQualifications(base *BaseModule, _ *broker.Broker[any]) (Module, error) {
	return &QualificationsSync{
		BaseModule: base,
	}, nil
}

func (g *QualificationsSync) GetName() string {
	return "qualifications"
}

func (g *QualificationsSync) Plan(
	ctx context.Context,
) (*discordtypes.State, []discord.Embed, error) {
	errs := multierr.Combine()

	stmt := tQualifications.
		SELECT(
			tQualifications.ID.AS("qualifications_entry.id"),
			tQualifications.Abbreviation.AS("qualifications_entry.abbreviation"),
			tQualifications.Title.AS("qualifications_entry.title"),
			tQualifications.DiscordSettings.AS("qualifications_entry.discord_settings"),
		).
		FROM(tQualifications).
		WHERE(jet.AND(
			tQualifications.CreatorJob.EQ(jet.String(g.job)),
			tQualifications.DiscordSyncEnabled.IS_TRUE(),
		))

	var qualifications []*qualificationsEntry
	if err := stmt.QueryContext(ctx, g.db, &qualifications); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			errs = multierr.Append(errs, err)
			return nil, nil, errs
		}
	}

	roles, logs, err := g.planRoles(qualifications)
	if err != nil {
		errs = multierr.Append(errs, err)
		return nil, logs, errs
	}

	users, uLogs, err := g.planUsers(ctx, roles)
	logs = append(logs, uLogs...)
	if err != nil {
		errs = multierr.Append(errs, err)
		return nil, logs, errs
	}

	return &discordtypes.State{
		Roles: roles,
		Users: users,
	}, logs, nil
}

func (g *QualificationsSync) planRoles(
	qualifications []*qualificationsEntry,
) ([]*discordtypes.Role, []discord.Embed, error) {
	logs := []discord.Embed{}
	roles := discordtypes.Roles{}

	syncSettings := g.settings.Load()

	errs := multierr.Combine()
	for _, entry := range qualifications {
		if entry.DiscordSettings.RoleName == nil ||
			strings.TrimSpace(entry.DiscordSettings.GetRoleName()) == "" {
			logs = append(logs, discord.Embed{
				Title: fmt.Sprintf(
					"Qualifications: Empty role name in qualification's discord sync settings \"%s: %s\" (ID: %d)",
					entry.Abbreviation,
					entry.QualificationTitle,
					entry.ID,
				),
				Description: fmt.Sprintf("Qualification ID: %d", entry.ID),
				Author:      embeds.EmbedAuthor,
				Color:       embeds.ColorWarn,
			})
			continue
		}

		roleFormat := strings.TrimSpace(syncSettings.GetQualificationsRoleFormat())
		roleName := strings.TrimSpace(entry.DiscordSettings.GetRoleName())
		if entry.DiscordSettings.RoleFormat != nil &&
			strings.TrimSpace(entry.DiscordSettings.GetRoleFormat()) != "" {
			rf := strings.TrimSpace(entry.DiscordSettings.GetRoleFormat())
			if strings.Contains(roleFormat, "%abbr%") || strings.Contains(roleFormat, "%name%") {
				roleFormat = rf
			}
		}

		entry.RoleName = strings.ReplaceAll(roleFormat, "%abbr%", entry.Abbreviation)
		entry.RoleName = strings.ReplaceAll(entry.RoleName, "%name%", roleName)

		roles = append(roles, &discordtypes.Role{
			Name:   entry.RoleName,
			Module: fmt.Sprintf(qualificationsRoleModulePrefix+"%d", entry.ID),
			Job:    g.job,
		})
	}

	return roles, logs, errs
}

func (g *QualificationsSync) planUsers(
	ctx context.Context,
	roles discordtypes.Roles,
) (discordtypes.Users, []discord.Embed, error) {
	logs := []discord.Embed{}

	qualificationRoles := map[int64]*discordtypes.Role{}
	for _, role := range roles {
		if strings.HasPrefix(role.Module, qualificationsRoleModulePrefix) {
			sGroup, found := strings.CutPrefix(role.Module, qualificationsRoleModulePrefix)
			if !found {
				continue
			}
			index, err := strconv.ParseInt(sGroup, 10, 64)
			if err != nil {
				return nil, logs, err
			}
			qualificationRoles[index] = role
		}
	}

	errs := multierr.Combine()

	jobs := []jet.Expression{jet.String(g.job)}
	for _, job := range g.appCfg.Get().Discord.GetIgnoredJobs() {
		jobs = append(jobs, jet.String(job))
	}

	tUsers := tables.User().AS("users")

	users := discordtypes.Users{}
	for qualificationId, role := range qualificationRoles {
		stmt := tAccsOauth2.
			SELECT(
				tAccsOauth2.AccountID.AS("qualification_user_mapping.account_id"),
				tAccsOauth2.ExternalID.AS("qualification_user_mapping.external_id"),
				tUsers.Job.AS("qualification_user_mapping.job"),
			).
			FROM(
				tQualificationsResults.
					INNER_JOIN(tQualifications,
						tQualifications.ID.EQ(tQualificationsResults.QualificationID).
							AND(tQualifications.DeletedAt.IS_NULL()),
					).
					INNER_JOIN(tUsers,
						tUsers.ID.EQ(tQualificationsResults.UserID),
					).
					INNER_JOIN(tAccs,
						tAccs.License.LIKE(jet.RawString("SUBSTRING_INDEX(`users`.`identifier`, ':', -1)")),
					).
					INNER_JOIN(tAccsOauth2,
						tAccsOauth2.AccountID.EQ(tAccs.ID),
					),
			).
			WHERE(jet.AND(
				tQualificationsResults.QualificationID.EQ(jet.Int64(qualificationId)),
				tQualificationsResults.DeletedAt.IS_NULL(),
				tQualificationsResults.Status.EQ(
					jet.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)),
				),
				tQualifications.Job.IN(jobs...),
				tAccsOauth2.Provider.EQ(jet.String("discord")),
			))

		var dest []*qualificationUserMapping
		if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, logs, err
			}
		}

		for _, u := range dest {
			externalId, err := strconv.ParseUint(u.ExternalID, 10, 64)
			if err != nil {
				errs = multierr.Append(
					errs,
					fmt.Errorf("failed to parse user oauth2 external id %d. %w", externalId, err),
				)
				continue
			}

			user := &discordtypes.User{
				ID:    discord.UserID(externalId),
				Roles: &discordtypes.UserRoles{},
				Job:   u.Job,
			}

			if u.Job != g.job {
				users.Add(user)
				continue
			}

			user.Roles.Sum = append(user.Roles.Sum, role)

			users.Add(user)
		}
	}

	return users, logs, errs
}
