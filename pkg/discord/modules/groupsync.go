package modules

import (
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/utils"
	jet "github.com/go-jet/jet/v2/mysql"
)

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

func (g *GroupSync) Run() error {
	if err := g.createGroupRoles(); err != nil {
		return err
	}

	return g.syncGroups()
}

type GroupSyncUser struct {
	ExternalID string `alias:"external_id"`
	Group      string `alias:"group"`
}

func (g *GroupSync) syncGroups() error {
	groups := []jet.Expression{}
	for serverGroup := range g.groupsToSync {
		groups = append(groups, jet.String(serverGroup))
	}

	stmt := tOauth2Accs.
		SELECT(
			tOauth2Accs.ExternalID.AS("groupsyncuser.external_id"),
			tUsers.Group.AS("groupsyncuser.group"),
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
			tUsers.Group.IN(groups...),
		))

	var dest []*GroupSyncUser
	if err := stmt.QueryContext(g.ctx, g.db, &dest); err != nil {
		return err
	}

	for _, user := range dest {
		member, err := g.discord.GuildMember(g.guild.ID, user.ExternalID)
		if err != nil {
			if restErr, ok := err.(*discordgo.RESTError); ok {
				if restErr.Response.StatusCode == http.StatusNotFound {
					continue
				}
			}
			return err
		}

		if err := g.setUserGroups(member, user.Group); err != nil {
			return err
		}
	}

	return g.cleanupUserGroupMembers(dest)
}

func (g *GroupSync) createGroupRoles() error {
	dcRoles := map[string]config.DiscordGroupRole{}
	for _, dcRole := range g.groupsToSync {
		if _, ok := dcRoles[dcRole.Name]; !ok {
			dcRoles[dcRole.Name] = dcRole
		}
	}

	for _, dcRole := range dcRoles {
		if utils.InSliceFunc(g.guild.Roles, func(in *discordgo.Role) bool {
			if in.Name == dcRole.Name {
				g.roles[dcRole.Name] = in
				return true
			}
			return false
		}) {
			// Role permissions are the same no need to edit/update
			if dcRole.Permissions != nil && *dcRole.Permissions == g.roles[dcRole.Name].Permissions {
				continue
			}

			role, err := g.discord.GuildRoleEdit(g.guild.ID, g.roles[dcRole.Name].ID, &discordgo.RoleParams{
				Name:        dcRole.Name,
				Permissions: dcRole.Permissions,
			})
			if err != nil {
				return err
			}

			g.roles[dcRole.Name] = role
		}

		role, err := g.discord.GuildRoleCreate(g.guild.ID, &discordgo.RoleParams{
			Name:        dcRole.Name,
			Permissions: dcRole.Permissions,
		})
		if err != nil {
			return err
		}

		g.roles[dcRole.Name] = role
	}

	return nil
}

func (g *GroupSync) setUserGroups(member *discordgo.Member, group string) error {
	dcRole, ok := g.groupsToSync[group]
	if !ok {
		return fmt.Errorf("no discord role mapping found for server group %s", group)
	}

	role, ok := g.roles[dcRole.Name]
	if !ok {
		return fmt.Errorf("no discord role found for server group %s", group)
	}

	if !utils.InSlice(member.Roles, role.ID) {
		if err := g.discord.GuildMemberRoleAdd(g.guild.ID, member.User.ID, role.ID); err != nil {
			return err
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
			if !utils.InSlice(guild.Members[i].Roles, role.ID) {
				continue
			}

			// If user is in the servergroup, all is okay, otherwise remove user from role
			if utils.InSliceFunc(users, func(in *GroupSyncUser) bool {
				return in.ExternalID == guild.Members[i].User.ID && g.groupsToSync[in.Group].Name == role.Name
			}) {
				continue
			}

			if err := g.discord.GuildMemberRoleRemove(g.guild.ID, guild.Members[i].User.ID, role.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
