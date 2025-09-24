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
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/embeds"
	discordtypes "github.com/fivenet-app/fivenet/v2025/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2025/services/auth"
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
	ExternalID string `alias:"external_id"`
	Group      string `alias:"group"`
	License    string `alias:"license"`
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

	serverGroups := []mysql.Expression{}
	for sGroup := range g.cfg.GroupSync.Mapping {
		serverGroups = append(serverGroups, mysql.String(sGroup))
	}

	tUsers := tables.User().AS("users")

	stmt := tAccsOauth2.
		SELECT(
			tAccsOauth2.ExternalID.AS("groupsyncuser.external_id"),
			tUsers.Group.AS("groupsyncuser.group"),
			tAccs.License.AS("groupsyncuser.license"),
		).
		FROM(
			tAccsOauth2.
				INNER_JOIN(tAccs,
					tAccs.ID.EQ(tAccsOauth2.AccountID),
				).
				INNER_JOIN(tUsers,
					tUsers.Identifier.LIKE(mysql.CONCAT(mysql.String("%"), tAccs.License)),
				),
		).
		WHERE(mysql.AND(
			tAccsOauth2.Provider.EQ(mysql.String("discord")),
			tUsers.Group.IN(serverGroups...),
		))

	var dest []*groupSyncUser
	if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return users, logs, err
		}
	}

	errs := multierr.Combine()
	for _, user := range dest {
		// Get group config based on users group
		groupCfg, ok := g.cfg.GroupSync.Mapping[user.Group]
		if !ok {
			continue
		}

		if groupCfg.NotSameJob {
			has, err := g.checkIfUserIsPartOfJob(ctx, user.License, g.job)
			if err != nil {
				g.logger.Error(
					fmt.Sprintf("failed to check if user has char in job %s", user.ExternalID),
					zap.String("group", user.Group),
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
					zap.String("group", user.Group),
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

	return users, logs, errs
}

func (g *GroupSync) checkIfUserIsPartOfJob(
	ctx context.Context,
	identifier string,
	job string,
) (bool, error) {
	tUsers := tables.User()

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("id"),
		).
		FROM(tUsers).
		WHERE(mysql.AND(
			tUsers.Identifier.LIKE(mysql.String(auth.BuildCharSearchIdentifier(identifier))),
			tUsers.Job.EQ(mysql.String(job)),
		))

	var dest []int32
	if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(dest) > 0, nil
}
