package modules

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord/embeds"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const defaultGroupSyncRoleColor = "#9B59B6"

type GroupSync struct {
	*BaseModule
}

type groupSyncUser struct {
	ExternalID string                  `alias:"external_id"`
	Groups     *accounts.AccountGroups `alias:"groups"`
	UserID     int32                   `alias:"user_id"`
	SameJob    bool
}

func init() {
	Modules["groupsync"] = NewGroupSync
}

func NewGroupSync(base *BaseModule, _ *broker.Broker[any]) (Module, error) {
	return &GroupSync{
		BaseModule: base,
	}, nil
}

func (g *GroupSync) GetName() string {
	return "groupsync"
}

func (g *GroupSync) Plan(ctx context.Context) (*discordtypes.State, []discord.Embed, error) {
	if len(g.cfg.GroupSync.Mapping) == 0 {
		// Nothind to do
		return nil, nil, nil
	}

	roles := g.planRoles()

	users, logs, err := g.planUsers(ctx, roles)
	if err != nil {
		return nil, logs, err
	}

	return &discordtypes.State{
		Roles: roles,
		Users: users,
	}, logs, nil
}

func (g *GroupSync) planRoles() []*discordtypes.Role {
	dcRoles := map[string]config.DiscordGroupRole{}
	for _, dcRole := range g.cfg.GroupSync.Mapping {
		if _, ok := dcRoles[dcRole.RoleName]; !ok {
			dcRoles[dcRole.RoleName] = dcRole
		}
	}

	roles := discordtypes.Roles{}
	for _, dcRole := range dcRoles {
		color := defaultGroupSyncRoleColor
		if dcRole.Color != "" {
			color = dcRole.Color
		}
		color = strings.TrimLeft(color, "#")

		n := new(big.Int)
		n.SetString(color, 16)
		if !n.IsInt64() || n.Int64() > math.MaxInt32 || n.Int64() < math.MinInt32 {
			g.logger.Warn(
				"role color value out of int32 range",
				zap.String("role", dcRole.RoleName),
				zap.String("color", color),
			)
			continue
		}
		//nolint:gosec // We ensure the value is within int32 range above.
		colorDec := int32(n.Int64())
		dcColor := discord.Color(colorDec)

		r := &discordtypes.Role{
			Name:  dcRole.RoleName,
			Color: &dcColor,

			Module: "GroupSync",
		}
		if dcRole.Permissions != nil {
			//nolint:gosec // Permissions are expected to be a valid Discord permissions value, if not at latest the Discord API will complain.
			ps := discord.Permissions(*dcRole.Permissions)
			r.Permissions = &ps
		}

		roles = append(roles, r)
	}

	return roles
}

func (g *GroupSync) planUsers(
	ctx context.Context,
	roles discordtypes.Roles,
) (discordtypes.Users, []discord.Embed, error) {
	users := discordtypes.Users{}
	logs := []discord.Embed{}

	tAccount := table.FivenetAccounts.AS("accounts")
	tUsers := table.FivenetUser.AS("users")

	groupsConditions := []mysql.BoolExpression{}
	for sGroup := range g.cfg.GroupSync.Mapping {
		groupsConditions = append(groupsConditions,
			dbutils.JSON_CONTAINS(tAccount.Groups, mysql.String("\""+sGroup+"\"")),
		)
	}

	condition := mysql.AND(
		tAccsOauth2.Provider.EQ(mysql.String("discord")),
		tAccount.Groups.IS_NOT_NULL(),
		mysql.OR(groupsConditions...),
	)

	stmt := tAccsOauth2.
		SELECT(
			tAccsOauth2.ExternalID.AS("groupsyncuser.external_id"),
			tAccount.Groups.AS("groupsyncuser.groups"),
			tUsers.ID.AS("groupsyncuser.user_id"),
		).
		FROM(
			tAccsOauth2.
				INNER_JOIN(tAccount,
					tAccount.ID.EQ(tAccsOauth2.AccountID),
				).
				INNER_JOIN(tUsers,
					tUsers.ID.EQ(tAccsOauth2.AccountID),
				),
		).
		WHERE(condition)

	var dest []*groupSyncUser
	if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return users, logs, err
		}
	}

	errs := multierr.Combine()
	for _, user := range dest {
		for _, group := range user.Groups.GetGroups() {
			// Get group config based on user's groups
			groupCfg, ok := g.cfg.GroupSync.Mapping[group]
			if !ok {
				continue
			}

			if groupCfg.NotSameJob {
				has, err := g.checkIfUserIsPartOfJob(ctx, user.UserID, g.job)
				if err != nil {
					g.logger.Error(
						fmt.Sprintf("failed to check if user has char in job %s", user.ExternalID),
						zap.String("group", group),
						zap.Error(err),
					)
					continue
				}
				if has {
					g.logger.Debug(
						fmt.Sprintf(
							"member %s is part of same job, not setting to group",
							user.ExternalID,
						),
						zap.String("group", group),
					)
					user.SameJob = true
					continue
				}
			}

			idx := slices.IndexFunc(roles, func(role *discordtypes.Role) bool {
				return role.Name == groupCfg.RoleName
			})
			if idx == -1 {
				logs = append(logs, discord.Embed{
					Title: fmt.Sprintf(
						"Group Sync: Failed to find dc role for group %s",
						groupCfg.RoleName,
					),
					Description: fmt.Sprintf("For DC ID %s", user.ExternalID),
					Author:      embeds.EmbedAuthor,
					Color:       embeds.ColorInfo,
				})
				continue
			}

			externalId, err := strconv.ParseUint(user.ExternalID, 10, 64)
			if err != nil {
				errs = multierr.Append(
					errs,
					fmt.Errorf("failed to parse oauth2 external id %d. %w", externalId, err),
				)
				continue
			}

			users.Add(&discordtypes.User{
				ID: discord.UserID(externalId),
				Roles: &discordtypes.UserRoles{
					Sum: []*discordtypes.Role{roles[idx]},
				},
			})
		}
	}

	return users, logs, errs
}

func (g *GroupSync) checkIfUserIsPartOfJob(
	ctx context.Context,
	userId int32,
	job string,
) (bool, error) {
	tUserJobs := table.FivenetUserJobs

	stmt := tUserJobs.
		SELECT(
			tUserJobs.Job,
			tUserJobs.Grade,
		).
		FROM(tUserJobs).
		WHERE(mysql.AND(
			tUserJobs.UserID.EQ(mysql.Int32(userId)),
			tUserJobs.Job.EQ(mysql.String(job)),
		))

	var dest []*users.UserJob
	if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(dest) > 0, nil
}
