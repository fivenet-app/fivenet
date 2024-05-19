package types

import (
	"context"
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/utils"
)

type NotPartOfFactionHandler func(ctx context.Context, session *discordgo.Session, guildId string, member *discordgo.Member) ([]*discordgo.MessageEmbed, error)

type Plan struct {
	DryRun bool

	Roles Roles
	Users Users

	NotPartOfFactionHandlers []NotPartOfFactionHandler `yaml:"-"`
}

func (p *Plan) Merge(plan *Plan) {
	p.Roles = append(p.Roles, plan.Roles...)

	for _, user := range plan.Users {
		p.Users.Add(user)
	}

	p.NotPartOfFactionHandlers = append(p.NotPartOfFactionHandlers, plan.NotPartOfFactionHandlers...)
}

func (p *Plan) Apply(ctx context.Context, dc *discordgo.Session, guildId string) ([]*discordgo.MessageEmbed, error) {
	logs := []*discordgo.MessageEmbed{}

	ls, err := p.Roles.Apply(ctx, dc, guildId, p.DryRun)
	logs = append(logs, ls...)
	if err != nil {
		return logs, err
	}

	guild, err := dc.State.Guild(guildId)
	if err != nil {
		return logs, err
	}

	for _, member := range guild.Members {
		// Skip users that have an assigment
		u, ok := p.Users[member.User.ID]
		if ok {
			if err := p.applyUser(ctx, dc, guildId, member, u); err != nil {
				return logs, err
			}

			if !u.ForceNotFound {
				continue
			}
		}

		for _, handler := range p.NotPartOfFactionHandlers {
			ls, err := handler(ctx, dc, guildId, member)
			logs = append(logs, ls...)
			if err != nil {
				return logs, err
			}
		}
	}

	return logs, nil
}

func (p *Plan) applyUser(ctx context.Context, dc *discordgo.Session, guildId string, member *discordgo.Member, user *User) error {
	if p.DryRun {
		return nil
	}

	addRoles := []string{}
	removeRoles := []string{}
	for _, userRole := range user.Roles {
		if !slices.Contains(member.Roles, userRole.ID) {
			addRoles = append(addRoles, userRole.ID)
		}
	}

	for _, role := range member.Roles {
		// If the role is bot managed, and the user is not assigned to the role, remove the role
		if slices.ContainsFunc(p.Roles, func(r *Role) bool {
			return r.ID == role
		}) && !slices.ContainsFunc(user.Roles, func(r *Role) bool {
			return r.ID == role
		}) {
			removeRoles = append(removeRoles, role)
		}
	}

	for _, roleId := range removeRoles {
		if err := dc.GuildMemberRoleRemove(guildId, member.User.ID, roleId, discordgo.WithContext(ctx)); err != nil {
			return err
		}
	}

	for _, roleId := range addRoles {
		if err := dc.GuildMemberRoleAdd(guildId, member.User.ID, roleId, discordgo.WithContext(ctx)); err != nil {
			return err
		}
	}

	if user.Nickname != nil && member.Nick != *user.Nickname {
		if err := dc.GuildMemberNickname(guildId, member.User.ID, *user.Nickname, discordgo.WithContext(ctx)); err != nil {
			return err
		}
	}

	return nil
}

type Roles []*Role

func (r Roles) Apply(ctx context.Context, dc *discordgo.Session, guildId string, dryRun bool) ([]*discordgo.MessageEmbed, error) {
	logs := []*discordgo.MessageEmbed{}
	if dryRun {
		return logs, nil
	}

	guild, err := dc.State.Guild(guildId)
	if err != nil {
		return logs, err
	}
	var botRole *discordgo.Role
	slices.ContainsFunc(guild.Roles, func(item *discordgo.Role) bool {
		if item.Name == "FiveNet" && item.Managed {
			botRole = item
			return true
		}
		return false
	})

	roles, err := dc.GuildRoles(guildId, discordgo.WithContext(ctx))
	if err != nil {
		return logs, err
	}

	for _, role := range r {
		idx := slices.IndexFunc(roles, func(a *discordgo.Role) bool {
			return a.Name == role.Name
		})
		if idx == -1 {
			res, err := dc.GuildRoleCreate(guildId, &discordgo.RoleParams{
				Name:        role.Name,
				Color:       role.Color,
				Permissions: role.Permissions,
			}, discordgo.WithContext(ctx))
			if err != nil {
				return logs, err
			}

			role.ID = res.ID
		} else {
			dcRole := roles[idx]

			if botRole != nil && dcRole.Position > botRole.Position {
				logs = append(logs, &discordgo.MessageEmbed{
					Type:        discordgo.EmbedTypeRich,
					Title:       fmt.Sprintf("Roles: Role %s (%s) can't be updated", dcRole.Name, dcRole.ID),
					Description: "FiveNet bot role is not high enough to update the role.",
					Author:      embeds.EmbedAuthor,
					Color:       embeds.ColorWarn,
				})
				continue
			}

			res, err := dc.GuildRoleEdit(guildId, dcRole.ID, &discordgo.RoleParams{
				Name:        role.Name,
				Color:       role.Color,
				Permissions: role.Permissions,
			}, discordgo.WithContext(ctx))
			if err != nil {
				return logs, err
			}

			role.ID = res.ID
		}
	}

	return logs, nil
}

type Role struct {
	ID          string
	Name        string
	Color       *int   `yaml:",omitempty"`
	Permissions *int64 `yaml:",omitempty"`

	Module string
}

type Users map[string]*User

func (u Users) Add(user *User) {
	existing, ok := u[user.ID]
	if !ok {
		u[user.ID] = user
	} else {
		existing.Merge(user)
	}
}

type User struct {
	ID       string
	Nickname *string `yaml:",omitempty"`

	Roles []*Role

	ForceNotFound bool `yaml:"-"`
}

func (u *User) Merge(user *User) {
	if u.ID != user.ID {
		return
	}

	if user.Nickname != nil {
		u.Nickname = user.Nickname
	}

	if len(user.Roles) > 0 {
		u.Roles = append(u.Roles, user.Roles...)
		u.Roles = utils.RemoveSliceDuplicates(u.Roles)
	}

	if u.ForceNotFound && !user.ForceNotFound {
		u.ForceNotFound = false
	}
}
