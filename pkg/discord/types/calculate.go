package types

import (
	"context"
	"fmt"
	"slices"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
)

func (s *State) Calculate(ctx context.Context, dc *state.State, dryRun bool) (*Plan, []discord.Embed, error) {
	plan := NewPlan(s.GuildID, dryRun)
	logs := []discord.Embed{}

	dc = dc.WithContext(ctx)

	roles, ls, err := s.calculateRoles(dc)
	logs = append(logs, ls...)
	if err != nil {
		return plan, logs, err
	}
	plan.Roles = roles

	members, err := dc.Members(s.GuildID)
	if err != nil {
		return plan, logs, fmt.Errorf("failed to get guild members. %w", err)
	}

	for _, member := range members {
		if member.User.Bot {
			continue
		}

		u := s.Users[member.User.ID]
		pu, err := s.calculateUserUpdates(ctx, member, u)
		if err != nil {
			return plan, logs, err
		}

		plan.Users = append(plan.Users, pu)
	}

	plan.Users = slices.DeleteFunc(plan.Users, func(u *User) bool {
		return u.ID == discord.NullUserID || ((u.Kick == nil || !*u.Kick) && (u.Roles == nil || (len(u.Roles.ToAdd) == 0 && len(u.Roles.ToRemove) == 0)))
	})

	return plan, logs, nil
}

func (s *State) calculateRoles(dc *state.State) (*PlanRoles, []discord.Embed, error) {
	logs := []discord.Embed{}

	roles, err := dc.Roles(s.GuildID)
	if err != nil {
		return nil, logs, fmt.Errorf("failed to get guild roles. %w", err)
	}
	var botRole discord.Role
	if idx := slices.IndexFunc(roles, func(item discord.Role) bool {
		return item.Name == "FiveNet" && item.Managed
	}); idx > -1 {
		botRole = roles[idx]
	}

	pr := &PlanRoles{}
	for _, role := range s.Roles {
		idx := slices.IndexFunc(roles, func(a discord.Role) bool {
			return a.Name == role.Name
		})
		if idx == -1 {
			pr.ToCreate = append(pr.ToCreate, role)
		} else {
			dcRole := roles[idx]

			role.ID = roles[idx].ID

			if botRole.ID != discord.NullRoleID && dcRole.Position > botRole.Position {
				logs = append(logs, discord.Embed{
					Title:       fmt.Sprintf("Roles: Role %s (%s) can't be updated", dcRole.Name, dcRole.ID),
					Description: "FiveNet bot role is not high enough to update the role.",
					Author:      embeds.EmbedAuthor,
					Color:       embeds.ColorWarn,
				})
				continue
			}

			if role.Color == dcRole.Color && role.Permissions == dcRole.Permissions {
				continue
			}

			pr.ToUpdate = append(pr.ToUpdate, role)
		}
	}

	return pr, logs, nil
}

func (s *State) calculateUserUpdates(ctx context.Context, member discord.Member, user *User) (*User, error) {
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
		if !slices.Contains(member.RoleIDs, userRole.ID) {
			user.Roles.ToAdd = append(user.Roles.ToAdd, userRole)
		}
	}

	for _, role := range member.RoleIDs {
		// If the role is bot managed, and the user is not assigned to the role, remove the role
		if idx := slices.IndexFunc(s.Roles, func(r *Role) bool {
			return r.ID == role
		}); idx > -1 {
			if !slices.ContainsFunc(user.Roles.Sum, func(r *Role) bool {
				return r.ID == role
			}) {
				if s.Roles[idx].Job != "" && s.Roles[idx].Job == user.Job {
					continue
				}

				user.Roles.ToRemove = append(user.Roles.ToRemove, s.Roles[idx])
			}
		}
	}

	return user, nil
}
