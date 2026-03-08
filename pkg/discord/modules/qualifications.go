package modules

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord/embeds"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
)

const (
	qualificationsRoleModulePrefix = "Qualifications-"

	qualificationsQueryLimit = 500
)

var (
	tQualifications        = table.FivenetQualifications
	tQualificationsResults = table.FivenetQualificationsResults
)

type QualificationsSync struct {
	*BaseModule
}

type qualificationsEntry struct {
	ID                 int64  `alias:"qualifications_entry.id"           sql:"primary_key"`
	Abbreviation       string `alias:"qualifications_entry.abbreviation"`
	QualificationTitle string `alias:"qualifications_entry.title"`

	DiscordSettings *qualifications.QualificationDiscordSettings `alias:"qualifications_entry.discord_settings"`
	RoleName        string
}

type qualificationUserMapping struct {
	UserID     int32            `alias:"qualification_user_mapping.user_id" sql:"primary_key"`
	ExternalID string           `alias:"external_id"                        sql:"primary_key"`
	Jobs       []*users.UserJob `alias:"jobs"`
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
		WHERE(mysql.AND(
			tQualifications.DeletedAt.IS_NULL(),
			tQualifications.CreatorJob.EQ(mysql.String(g.job)),
			tQualifications.DiscordSyncEnabled.IS_TRUE(),
		))

	var qualifications []*qualificationsEntry
	if err := stmt.QueryContext(ctx, g.db, &qualifications); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			errs = multierr.Append(errs, err)
			return nil, nil, errs
		}
	}

	rolesMap, logs, err := g.planRoles(qualifications)
	if err != nil {
		errs = multierr.Append(errs, err)
		return nil, logs, errs
	}

	users, uLogs, err := g.planUsers(ctx, rolesMap)
	logs = append(logs, uLogs...)
	if err != nil {
		errs = multierr.Append(errs, err)
		return nil, logs, errs
	}

	roles := make(discordtypes.Roles, 0, len(rolesMap))
	for _, role := range rolesMap {
		roles = append(roles, role)
	}

	return &discordtypes.State{
		Roles: roles,
		Users: users,
	}, logs, nil
}

func (g *QualificationsSync) planRoles(
	qualifications []*qualificationsEntry,
) (map[int64]*discordtypes.Role, []discord.Embed, error) {
	logs := []discord.Embed{}
	roles := map[int64]*discordtypes.Role{}

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

		roles[entry.ID] = &discordtypes.Role{
			Name:   entry.RoleName,
			Module: fmt.Sprintf(qualificationsRoleModulePrefix+"%d", entry.ID),
			Job:    g.job,
		}
	}

	return roles, logs, errs
}

func (g *QualificationsSync) planUsers(
	ctx context.Context,
	roles map[int64]*discordtypes.Role,
) (discordtypes.Users, []discord.Embed, error) {
	logs := []discord.Embed{}

	errs := multierr.Combine()

	users := discordtypes.Users{}
	for qualificationId, role := range roles {
		if err := g.queryAndPlanUsersForQualification(ctx, qualificationId, role, &users); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
	}

	return users, logs, errs
}

func (g *QualificationsSync) queryAndPlanUsersForQualification(
	ctx context.Context,
	qualificationId int64,
	role *discordtypes.Role,
	users *discordtypes.Users,
) error {
	offset := int64(0)

	var errs error
	for {
		userMappings, err := g.queryUsers(ctx, qualificationId, offset)
		if err != nil {
			errs = multierr.Append(errs, err)
			break
		}

		for _, u := range userMappings {
			externalId, err := strconv.ParseUint(u.ExternalID, 10, 64)
			if err != nil {
				errs = multierr.Append(
					errs,
					fmt.Errorf("failed to parse user oauth2 external id %d. %w", externalId, err),
				)
				continue
			}

			users.Add(&discordtypes.User{
				ID: discord.UserID(externalId),
				Roles: &discordtypes.UserRoles{
					Sum: []*discordtypes.Role{role},
				},
				Jobs: u.Jobs,
			})
		}

		count := int64(len(userMappings))
		if count < qualificationsQueryLimit {
			break
		}
		offset += count
	}

	return errs
}

func (g *QualificationsSync) queryUsers(
	ctx context.Context,
	qualificationId int64,
	offset int64,
) ([]*qualificationUserMapping, error) {
	tAccs := table.FivenetAccounts
	tUsers := table.FivenetUser.AS("users")
	tUserJobs := table.FivenetUserJobs.AS("user_jobs")

	stmt := tAccsOauth2.
		SELECT(
			tAccsOauth2.ExternalID.AS("qualification_user_mapping.external_id"),
			tUsers.ID.AS("qualification_user_mapping.user_id"),
			// User's jobs
			tUserJobs.Job.AS("jobs.job"),
			tUserJobs.Grade.AS("jobs.grade"),
		).
		FROM(
			tQualificationsResults.
				INNER_JOIN(tQualifications,
					tQualifications.ID.EQ(tQualificationsResults.QualificationID),
				).
				INNER_JOIN(tUsers,
					tUsers.ID.EQ(tQualificationsResults.UserID),
				).
				INNER_JOIN(tAccs,
					mysql.OR(
						tAccs.ID.EQ(tUsers.AccountID),
						tAccs.License.EQ(tUsers.License),
					),
				).
				INNER_JOIN(tAccsOauth2,
					tAccsOauth2.AccountID.EQ(tAccs.ID),
				).
				INNER_JOIN(tUserJobs,
					mysql.AND(
						tUserJobs.UserID.EQ(tUsers.ID),
						tUserJobs.Job.EQ(mysql.String(g.job)),
					),
				),
		).
		WHERE(mysql.AND(
			tAccsOauth2.Provider.EQ(mysql.String("discord")),
			tQualifications.Job.EQ(mysql.String(g.job)),
			tQualifications.DeletedAt.IS_NULL(),
			tQualificationsResults.QualificationID.EQ(mysql.Int64(qualificationId)),
			tQualificationsResults.DeletedAt.IS_NULL(),
			tQualificationsResults.Status.EQ(
				mysql.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)),
			),
		)).
		OFFSET(offset).
		LIMIT(qualificationsQueryLimit)

	var userMappings []*qualificationUserMapping
	if err := stmt.QueryContext(ctx, g.db, &userMappings); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query qualification user mappings. %w", err)
		}
	}

	return userMappings, nil
}
