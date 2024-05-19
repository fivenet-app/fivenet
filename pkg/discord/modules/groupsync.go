package modules

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
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

func NewGroupSync(base *BaseModule) (Module, error) {
	return &GroupSync{
		BaseModule: base,
	}, nil
}

func (g *GroupSync) Plan(ctx context.Context) (*types.Plan, []*discordgo.MessageEmbed, error) {
	roles := g.planRoles()

	users, logs, err := g.planUsers(ctx, roles)
	if err != nil {
		return nil, logs, err
	}

	return &types.Plan{
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

	i := 0
	for _, dcRole := range dcRoles {
		color := defaultGroupSyncRoleColor
		if dcRole.Color != "" {
			color = dcRole.Color
		}
		color = strings.TrimLeft(color, "#")

		n := new(big.Int)
		n.SetString(color, 16)
		colorDec := int(n.Int64())

		roles = append(roles, &types.Role{
			Name:        dcRole.RoleName,
			Permissions: dcRole.Permissions,
			Color:       &colorDec,

			Module: "GroupSync",
		})

		i++
	}

	return roles
}

func (g *GroupSync) planUsers(ctx context.Context, roles types.Roles) (types.Users, []*discordgo.MessageEmbed, error) {
	users := types.Users{}
	logs := []*discordgo.MessageEmbed{}

	serverGroups := []jet.Expression{}
	for sGroup := range g.cfg.GroupSync.Mapping {
		serverGroups = append(serverGroups, jet.String(sGroup))
	}

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
					tUsers.Identifier.LIKE(jet.CONCAT(jet.String("char%:"), tAccs.License)),
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
			has, err := g.checkIfUserHasCharInJob(ctx, user.License, g.job)
			if err != nil {
				g.logger.Error(fmt.Sprintf("failed to check if user has char in job %s", user.ExternalID), zap.Error(err))
				continue
			}
			if has {
				g.logger.Debug(fmt.Sprintf("member is part of same job, not setting group %s", user.ExternalID))
				user.SameJob = true
				continue
			}
		}

		idx := slices.IndexFunc(roles, func(role *types.Role) bool {
			return role.Name == groupCfg.RoleName
		})
		if idx == -1 {
			logs = append(logs, &discordgo.MessageEmbed{
				Type:        discordgo.EmbedTypeRich,
				Title:       fmt.Sprintf("Group Sync: Failed to find dc role for group %s", groupCfg.RoleName),
				Description: fmt.Sprintf("For DC ID %s", user.ExternalID),
				Author:      embeds.EmbedAuthor,
				Color:       embeds.ColorInfo,
			})
			continue
		}

		users.Add(&types.User{
			ID:            user.ExternalID,
			Roles:         []*types.Role{roles[idx]},
			ForceNotFound: true,
		})
	}

	return users, logs, errs
}

func (g *GroupSync) checkIfUserHasCharInJob(ctx context.Context, identifier string, job string) (bool, error) {
	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("id"),
		).
		FROM(tUsers).
		WHERE(jet.AND(
			tUsers.Identifier.LIKE(jet.String("char%:"+identifier)),
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
