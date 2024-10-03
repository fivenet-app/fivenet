package modules

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
)

var (
	tQualifications        = table.FivenetQualifications
	tQualificationsResults = table.FivenetQualificationsResults
)

type QualificationsSync struct {
	*BaseModule

	ignoredJobs []string
}

type qualificationsEntry struct {
	ID              uint64                                       `alias:"qualifications_entry.id"`
	DiscordSettings *qualifications.QualificationDiscordSettings `alias:"qualifications_entry.discord_settings"`
	RoleName        string
}

type qualificationUserMapping struct {
	AccountID  uint64 `alias:"account_id"`
	ExternalID string `alias:"external_id"`
	Job        string `alias:"job"`
}

func init() {
	Modules["qualifications"] = NewQualifications
}

func NewQualifications(base *BaseModule) (Module, error) {
	return &QualificationsSync{
		BaseModule:  base,
		ignoredJobs: base.appCfg.Get().Discord.IgnoredJobs,
	}, nil
}

func (g *QualificationsSync) Plan(ctx context.Context) (*types.State, []discord.Embed, error) {
	errs := multierr.Combine()

	stmt := tQualifications.
		SELECT(
			tQualifications.ID.AS("qualificationsentry.id"),
			tQualifications.DiscordSettings.AS("qualificationsentry.discord_settings"),
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

	return &types.State{
		Roles: roles,
		Users: users,
	}, logs, nil
}

const (
	qualificationsRoleModulePrefix = "Qualifications-"
)

func (g *QualificationsSync) planRoles(qualifications []*qualificationsEntry) ([]*types.Role, []discord.Embed, error) {
	logs := []discord.Embed{}
	roles := types.Roles{}

	errs := multierr.Combine()

	for _, entry := range qualifications {
		if entry.DiscordSettings.RoleName == nil || *entry.DiscordSettings.RoleName == "" {
			logs = append(logs, discord.Embed{
				Title:       fmt.Sprintf("Qualifications: Empty role name in qualification Discord settings %d", entry.ID),
				Description: fmt.Sprintf("Qualification ID: %d", entry.ID),
				Author:      embeds.EmbedAuthor,
				Color:       embeds.ColorWarn,
			})
			continue
		}

		entry.RoleName = strings.ReplaceAll(g.settings.QualificationsRoleFormat, "%name%", *entry.DiscordSettings.RoleName)

		roles = append(roles, &types.Role{
			Name:   entry.RoleName,
			Module: fmt.Sprintf(qualificationsRoleModulePrefix+"%d", entry.ID),
			Job:    g.job,
		})
	}

	return roles, logs, errs
}

func (g *QualificationsSync) planUsers(ctx context.Context, roles types.Roles) (types.Users, []discord.Embed, error) {
	logs := []discord.Embed{}

	qualificationRoles := map[uint64]*types.Role{}
	for _, role := range roles {
		if strings.HasPrefix(role.Module, qualificationsRoleModulePrefix) {
			sGroup, found := strings.CutPrefix(role.Module, qualificationsRoleModulePrefix)
			if !found {
				continue
			}
			index, err := strconv.Atoi(sGroup)
			if err != nil {
				return nil, logs, err
			}
			qualificationRoles[uint64(index)] = role
		}
	}

	errs := multierr.Combine()

	jobs := []jet.Expression{jet.String(g.job)}
	for _, job := range g.ignoredJobs {
		jobs = append(jobs, jet.String(job))
	}

	users := types.Users{}
	for qualificationId, role := range qualificationRoles {
		stmt := tOauth2Accs.
			SELECT(
				tOauth2Accs.AccountID.AS("qualificationusermapping.account_id"),
				tOauth2Accs.ExternalID.AS("qualificationusermapping.external_id"),
				tUsers.Job.AS("qualificationusermapping.job"),
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
					INNER_JOIN(tOauth2Accs,
						tOauth2Accs.AccountID.EQ(tAccs.ID),
					),
			).
			WHERE(jet.AND(
				tQualificationsResults.QualificationID.EQ(jet.Uint64(qualificationId)),
				tQualificationsResults.Status.EQ(jet.Int16(int16(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL))),
				tQualifications.Job.IN(jobs...),
				tOauth2Accs.Provider.EQ(jet.String("discord")),
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
				errs = multierr.Append(errs, fmt.Errorf("failed to parse user oauth2 external id %d. %w", externalId, err))
				continue
			}

			user := &types.User{
				ID:    discord.UserID(externalId),
				Roles: &types.UserRoles{},
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
