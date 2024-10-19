package types

import (
	"context"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"go.uber.org/multierr"
)

func (p *Plan) Apply(ctx context.Context, dc *state.State) ([]discord.Embed, error) {
	logs := []discord.Embed{}

	if p.DryRun {
		return logs, nil
	}

	dc = dc.WithContext(ctx)

	if err := p.applyRoles(dc); err != nil {
		return logs, err
	}

	uLogs, err := p.applyUsers(dc)
	if err != nil {
		return logs, err
	}
	logs = append(logs, uLogs...)

	return logs, nil
}

func (p *Plan) applyRoles(dc *state.State) error {
	errs := multierr.Combine()

	for _, role := range p.Roles.ToCreate {
		var ps discord.Permissions
		if role.Permissions != nil {
			ps = *role.Permissions
		}

		roleData := api.CreateRoleData{
			Name:        role.Name,
			Permissions: ps,
		}
		if role.Color != nil {
			roleData.Color = *role.Color
		}
		res, err := dc.CreateRole(p.GuildID, roleData)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to create role %s. %w", role.Name, err))
			continue
		}

		role.ID = res.ID
	}

	for _, role := range p.Roles.ToUpdate {
		roleData := api.ModifyRoleData{
			Name:        option.NewNullableString(role.Name),
			Permissions: role.Permissions,
		}
		if role.Color != nil {
			roleData.Color = *role.Color
		}

		if _, err := dc.ModifyRole(p.GuildID, role.ID, roleData); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to update role %s. %w", role.Name, err))
			continue
		}
	}

	return errs
}

func (p *Plan) applyUsers(dc *state.State) ([]discord.Embed, error) {
	logs := []discord.Embed{}
	errs := multierr.Combine()

	for _, user := range p.Users {
		if user.Kick != nil && *user.Kick {
			if user.KickReason == "" {
				user.KickReason = "FiveNet Bot - Kick (No reason given)"
			}

			if err := dc.Kick(p.GuildID, user.ID, api.AuditLogReason(user.KickReason)); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to kick user %s (reason: %q). %w", user.ID, user.KickReason, err))
				continue
			}
			continue
		}

		if user.Nickname != nil && *user.Nickname == "" {
			if err := dc.ModifyMember(p.GuildID, user.ID, api.ModifyMemberData{
				Nick: user.Nickname,
			}); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to set user %s nickname (%q). %w", user.ID, *user.Nickname, err))
				continue
			}
		}

		if len(user.Roles.ToRemove) > 0 {
			out := ""
			for _, role := range user.Roles.ToRemove.ToSlice() {
				out += "Role Name: " + role.Name + " (ID: " + role.ID.String() + ", " + role.Module + "); "
			}
			log.Printf("removing roles from user %d - roles, toRemove: %s", user.ID, out)
		}

		for _, role := range user.Roles.ToRemove {
			if err := dc.RemoveRole(p.GuildID, user.ID, role.ID, api.AuditLogReason(role.Module)); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to remove user %s from role %s (%s). %w", user.ID, role.Name, role.ID, err))
				continue
			}
		}

		for _, role := range user.Roles.ToAdd {
			if err := dc.AddRole(p.GuildID, user.ID, role.ID, api.AddRoleData{
				AuditLogReason: api.AuditLogReason(role.Module),
			}); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to add user %s to role %s (%s). %w", user.ID, role.Name, role.ID, err))
				continue
			}
		}
	}

	return logs, errs
}
