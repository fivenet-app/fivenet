package types

import (
	"context"
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
)

func (s *State) Calculate(ctx context.Context, dc *discordgo.Session) (*Plan, []*discordgo.MessageEmbed, error) {
	plan := NewPlan(s.GuildID, true)
	logs := []*discordgo.MessageEmbed{}

	roles, ls, err := s.calculateRoles(ctx, dc)
	logs = append(logs, ls...)
	if err != nil {
		return plan, logs, err
	}
	plan.Roles = roles

	guild, err := dc.State.Guild(s.GuildID)
	if err != nil {
		return plan, logs, err
	}

	for _, member := range guild.Members {
		if member.User.Bot {
			continue
		}

		u, _ := s.Users[member.User.ID]
		pu, err := s.calculateUserUpdates(ctx, member, u)
		if err != nil {
			return plan, logs, err
		}

		plan.Users = append(plan.Users, pu)
	}

	plan.Users = slices.DeleteFunc(plan.Users, func(u *User) bool {
		return u.ID == "" || ((u.Kick == nil || *u.Kick == false) && (u.Roles == nil || (len(u.Roles.ToAdd) == 0 || len(u.Roles.ToRemove) == 0)))
	})

	return plan, logs, nil
}

func (s *State) calculateUserUpdates(ctx context.Context, member *discordgo.Member, user *User) (*User, error) {
	if user == nil {
		user = &User{
			ID:    member.User.ID,
			Roles: &UserRoles{},
		}
	}
	for _, fn := range s.UserProcessors {
		fn(ctx, s.GuildID, member, user)
	}

	for _, userRole := range user.Roles.Sum {
		if !slices.Contains(member.Roles, userRole.ID) {
			user.Roles.ToAdd = append(user.Roles.ToAdd, userRole)
		}
	}

	for _, role := range member.Roles {
		// If the role is bot managed, and the user is not assigned to the role, remove the role
		if idx := slices.IndexFunc(s.Roles, func(r *Role) bool {
			return r.ID == role
		}); idx > -1 && !slices.ContainsFunc(user.Roles.Sum, func(r *Role) bool {
			return r.ID == role
		}) {
			if s.Roles[idx].Job != "" && s.Roles[idx].Job == user.Job {
				continue
			}
			user.Roles.ToRemove = append(user.Roles.ToRemove, s.Roles[idx])
		}
	}

	return user, nil
}

func (s *State) calculateRoles(ctx context.Context, dc *discordgo.Session) (*PlanRoles, []*discordgo.MessageEmbed, error) {
	logs := []*discordgo.MessageEmbed{}

	roles, err := dc.GuildRoles(s.GuildID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, logs, err
	}
	var botRole *discordgo.Role
	slices.ContainsFunc(roles, func(item *discordgo.Role) bool {
		if item.Name == "FiveNet" && item.Managed {
			botRole = item
			return true
		}
		return false
	})

	pr := &PlanRoles{}
	for _, role := range s.Roles {
		idx := slices.IndexFunc(roles, func(a *discordgo.Role) bool {
			return a.Name == role.Name
		})
		if idx == -1 {
			pr.ToCreate = append(pr.ToCreate, role)
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

			role.ID = roles[idx].ID

			if (role.Color == nil || *role.Color == dcRole.Color) && (role.Permissions == nil || *role.Permissions == dcRole.Permissions) {
				continue
			}

			pr.ToUpdate = append(pr.ToUpdate, role)
		}
	}

	return pr, logs, nil
}
