package modules

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	jet "github.com/go-jet/jet/v2/mysql"
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

func (g *GroupSync) Plan(ctx context.Context) (*types.State, []discord.Embed, error) {
	roles := g.planRoles()

	users, logs, err := g.planUsers(ctx, roles)
	if err != nil {
		return nil, logs, err
	}

	return &types.State{
		Roles: roles,
		Users: users,
	}, logs, nil
}

func (g *GroupSync) planRoles() []*types.Role {
	dcRoles := map[string]config.DiscordGroupRole{}
	for _, dcRole := range g.cfg.GroupSync.Mapping {
		if _, ok := dcRoles[dcRole.RoleName]; !ok {
			dcRoles[dcRole.RoleName] = dcRole
		}
	}

	roles := types.Roles{}
	for _, dcRole := range dcRoles {
		color := defaultGroupSyncRoleColor
		if dcRole.Color != "" {
			color = dcRole.Color
		}
		color = strings.TrimLeft(color, "#")

		n := new(big.Int)
		n.SetString(color, 16)
		colorDec := int32(n.Int64())
		dcColor := discord.Color(colorDec)

		r := &types.Role{
			Name:  dcRole.RoleName,
			Color: &dcColor,

			Module: "GroupSync",
		}
		if dcRole.Permissions != nil {
			ps := discord.Permissions(*dcRole.Permissions)
			r.Permissions = &ps
		}

		roles = append(roles, r)
	}

	return roles
}

func (g *GroupSync) planUsers(ctx context.Context, roles types.Roles) (types.Users, []discord.Embed, error) {
	users := types.Users{}
	logs := []discord.Embed{}

	serverGroups := []jet.Expression{}
	for sGroup := range g.cfg.GroupSync.Mapping {
		serverGroups = append(serverGroups, jet.String(sGroup))
	}

	tUsers := tables.Users().AS("users")

	stmt := tOauth2Accs.
		SELECT(
			tOauth2Accs.ExternalID.AS("groupsyncuser.external_id"),
			tUsers.Group.AS("groupsyncuser.group"),
			tAccs.License.AS("groupsyncuser.license"),
		).
		FROM(
			tOauth2Accs.
				INNER_JOIN(tAccs,
					tAccs.ID.EQ(tOauth2Accs.AccountID),
				).
				INNER_JOIN(tUsers,
					tUsers.Identifier.LIKE(jet.CONCAT(jet.String("%"), tAccs.License)),
				),
		).
		WHERE(jet.AND(
			tOauth2Accs.Provider.EQ(jet.String("discord")),
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
				g.logger.Error(fmt.Sprintf("failed to check if user has char in job %s", user.ExternalID), zap.String("group", user.Group), zap.Error(err))
				continue
			}
			if has {
				g.logger.Debug(fmt.Sprintf("member %s is part of same job, not setting to group", user.ExternalID), zap.String("group", user.Group))
				user.SameJob = true
				continue
			}
		}

		idx := slices.IndexFunc(roles, func(role *types.Role) bool {
			return role.Name == groupCfg.RoleName
		})
		if idx == -1 {
			logs = append(logs, discord.Embed{
				Title:       fmt.Sprintf("Group Sync: Failed to find dc role for group %s", groupCfg.RoleName),
				Description: fmt.Sprintf("For DC ID %s", user.ExternalID),
				Author:      embeds.EmbedAuthor,
				Color:       embeds.ColorInfo,
			})
			continue
		}

		externalId, err := strconv.ParseUint(user.ExternalID, 10, 64)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to parse oauth2 external id %d. %w", externalId, err))
			continue
		}

		users.Add(&types.User{
			ID: discord.UserID(externalId),
			Roles: &types.UserRoles{
				Sum: []*types.Role{roles[idx]},
			},
		})
	}

	return users, logs, errs
}

func (g *GroupSync) checkIfUserIsPartOfJob(ctx context.Context, identifier string, job string) (bool, error) {
	tUsers := tables.Users()

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("id"),
		).
		FROM(tUsers).
		WHERE(jet.AND(
			tUsers.Identifier.LIKE(jet.CONCAT(jet.String("%"), jet.String(identifier))),
			tUsers.Job.EQ(jet.String(job)),
		))

	var dest []int32
	if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(dest) > 0, nil
}
