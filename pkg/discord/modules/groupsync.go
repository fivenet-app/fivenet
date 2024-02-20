package modules

import (
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/config"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

const GroupSyncDefaultRoleColor = "#9B59B6"

type GroupSync struct {
	*BaseModule

	groupsToSync map[string]config.DiscordGroupRole

	roles map[string]*discordgo.Role
}

func init() {
	Modules["groupsync"] = NewGroupSync
}

func NewGroupSync(base *BaseModule) (Module, error) {
	return &GroupSync{
		BaseModule:   base,
		groupsToSync: base.cfg.GroupSync.Mapping,
		roles:        map[string]*discordgo.Role{},
	}, nil
}

func (g *GroupSync) Run(settings *users.DiscordSyncSettings) ([]*discordgo.MessageEmbed, error) {
	if err := g.createGroupRoles(); err != nil {
		return nil, err
	}

	return nil, g.syncGroups()
}

func (g *GroupSync) createGroupRoles() error {
	guildRoles, err := g.discord.GuildRoles(g.guild.ID)
	if err != nil {
		return err
	}

	dcRoles := map[string]config.DiscordGroupRole{}
	for _, dcRole := range g.groupsToSync {
		if _, ok := dcRoles[dcRole.RoleName]; !ok {
			dcRoles[dcRole.RoleName] = dcRole
		}
	}

	for _, dcRole := range dcRoles {
		if slices.ContainsFunc(guildRoles, func(in *discordgo.Role) bool {
			if strings.EqualFold(in.Name, dcRole.RoleName) {
				g.roles[dcRole.RoleName] = in
				return true
			}
			return false
		}) {
			// Role permissions are the same no need to edit/update
			if dcRole.Permissions != nil && *dcRole.Permissions == g.roles[dcRole.RoleName].Permissions {
				continue
			}

			var color string
			if dcRole.Color != "" {
				color = dcRole.Color
			}
			color = strings.TrimLeft(color, "#")

			n := new(big.Int)
			n.SetString(color, 16)
			colorDec := int(n.Int64())

			g.logger.Debug("updating group role", zap.String("group_name", dcRole.RoleName))
			role, err := g.discord.GuildRoleEdit(g.guild.ID, g.roles[dcRole.RoleName].ID, &discordgo.RoleParams{
				Name:        dcRole.RoleName,
				Permissions: dcRole.Permissions,
				Color:       &colorDec,
			})
			if err != nil {
				return fmt.Errorf("failed to edit role %s permissions: %w", g.roles[dcRole.RoleName].ID, err)
			}

			g.roles[dcRole.RoleName] = role
			continue
		}

		if _, ok := g.roles[dcRole.RoleName]; ok {
			continue
		}

		g.logger.Debug("creating group role", zap.String("group_name", dcRole.RoleName))
		role, err := g.discord.GuildRoleCreate(g.guild.ID, &discordgo.RoleParams{
			Name:        dcRole.RoleName,
			Permissions: dcRole.Permissions,
		})
		if err != nil {
			return fmt.Errorf("failed to create role %s (%s): %w", dcRole.RoleName, dcRole.RoleName, err)
		}

		g.roles[dcRole.RoleName] = role
	}

	g.logger.Debug("created or updated group roles")

	return nil
}

type GroupSyncUser struct {
	ExternalID string `alias:"external_id"`
	Group      string `alias:"group"`
	License    string `alias:"license"`
	SameJob    bool
}

func (g *GroupSync) syncGroups() error {
	serverGroups := []jet.Expression{}
	for sGroup := range g.groupsToSync {
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

	var dest []*GroupSyncUser
	if err := stmt.QueryContext(g.ctx, g.db, &dest); err != nil {
		return err
	}

	for _, user := range dest {
		// Get group config based on users group
		groupCfg, ok := g.groupsToSync[user.Group]
		if !ok {
			continue
		}

		if groupCfg.NotSameJob {
			has, err := g.checkIfUserHasCharInJob(user.License, g.job)
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

		member, err := g.discord.GuildMember(g.guild.ID, user.ExternalID)
		if err != nil {
			if restErr, ok := err.(*discordgo.RESTError); ok {
				if restErr.Response.StatusCode == http.StatusNotFound {
					continue
				}
			}

			g.logger.Error(fmt.Sprintf("failed to get external guild member %s", user.ExternalID), zap.Error(err))
			continue
		}

		if err := g.setUserGroups(member, user.Group); err != nil {
			g.logger.Error("failed to set user groups", zap.Error(err))
			continue
		}
	}

	return g.cleanupUserGroupMembers(dest)
}

func (g *GroupSync) checkIfUserHasCharInJob(identifier string, job string) (bool, error) {
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
	if err := stmt.QueryContext(g.ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(dest) > 0, nil
}

func (g *GroupSync) setUserGroups(member *discordgo.Member, group string) error {
	dcRole, ok := g.groupsToSync[group]
	if !ok {
		return fmt.Errorf("no discord role mapping found for server group %s", group)
	}

	role, ok := g.roles[dcRole.RoleName]
	if !ok {
		return fmt.Errorf("no discord role found for server group %s", group)
	}

	if !slices.Contains(member.Roles, role.ID) {
		if err := g.discord.GuildMemberRoleAdd(g.guild.ID, member.User.ID, role.ID); err != nil {
			return fmt.Errorf("failed to add member to %s (%s) role: %w", role.Name, role.ID, err)
		}
	}

	return nil
}

func (g *GroupSync) cleanupUserGroupMembers(users []*GroupSyncUser) error {
	guild, err := g.discord.State.Guild(g.guild.ID)
	if err != nil {
		return err
	}

	for i := 0; i < len(guild.Members); i++ {
		for _, role := range g.roles {
			// If user isn't in one of the synced groups, continue
			if !slices.Contains(guild.Members[i].Roles, role.ID) {
				continue
			}

			// If user is in the servergroup, all is okay, otherwise remove user from role
			user := guild.Members[i].User
			if slices.ContainsFunc(users, func(in *GroupSyncUser) bool {
				return in.ExternalID == user.ID && g.groupsToSync[in.Group].RoleName == role.Name && !in.SameJob
			}) {
				continue
			}

			if err := g.discord.GuildMemberRoleRemove(g.guild.ID, user.ID, role.ID); err != nil {
				return fmt.Errorf("failed to remove member from role %s (%s): %w", role.Name, role.ID, err)
			}
		}
	}

	return nil
}
