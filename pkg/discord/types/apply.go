package types

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/multierr"
)

func (p *Plan) Apply(ctx context.Context, dc *discordgo.Session) ([]*discordgo.MessageEmbed, error) {
	logs := []*discordgo.MessageEmbed{}

	if p.DryRun {
		return logs, nil
	}

	if err := p.applyRoles(ctx, dc, p.GuildID); err != nil {
		return logs, err
	}

	if err := p.applyUsers(ctx, dc); err != nil {
		return logs, err
	}

	return logs, nil
}

func (p *Plan) applyRoles(ctx context.Context, dc *discordgo.Session, guildId string) error {
	errs := multierr.Combine()

	for _, role := range p.Roles.ToCreate {
		res, err := dc.GuildRoleCreate(guildId, &discordgo.RoleParams{
			Name:        role.Name,
			Color:       role.Color,
			Permissions: role.Permissions,
		}, discordgo.WithContext(ctx))
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to create role %s. %w", role.Name, err))
			continue
		}

		role.ID = res.ID
	}

	for _, role := range p.Roles.ToUpdate {
		_, err := dc.GuildRoleEdit(guildId, role.ID, &discordgo.RoleParams{
			Name:        role.Name,
			Color:       role.Color,
			Permissions: role.Permissions,
		}, discordgo.WithContext(ctx))
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to update role %s. %w", role.Name, err))
			continue
		}
	}

	return errs
}

func (p *Plan) applyUsers(ctx context.Context, dc *discordgo.Session) error {
	errs := multierr.Combine()

	for _, user := range p.Users {
		if user.Kick != nil && *user.Kick {
			if user.KickReason == "" {
				user.KickReason = "FiveNet Bot - Auto Kick (No reason given)"
			}

			if err := dc.GuildMemberDeleteWithReason(p.GuildID, user.ID, user.KickReason,
				discordgo.WithContext(ctx)); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to kick user %s (reason: %q). %w", user.ID, user.KickReason, err))
				continue
			}
			continue
		}

		if user.Nickname != nil && *user.Nickname == "" {
			if err := dc.GuildMemberNickname(p.GuildID, user.ID, *user.Nickname, discordgo.WithContext(ctx)); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to set user %s nickname (%q). %w", user.ID, *user.Nickname, err))
				continue
			}
		}

		for _, role := range user.Roles.ToRemove {
			if err := dc.GuildMemberRoleRemove(p.GuildID, user.ID, role.ID, discordgo.WithContext(ctx)); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to remove user %s from role %s (%s). %w", user.ID, role.ID, role.Name, err))
				continue
			}
		}

		for _, role := range user.Roles.ToAdd {
			if err := dc.GuildMemberRoleAdd(p.GuildID, user.ID, role.ID, discordgo.WithContext(ctx)); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to add user %s to role %s (%s). %w", user.ID, role.ID, role.Name, err))
				continue
			}
		}
	}

	return errs
}
