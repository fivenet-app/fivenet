package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

const DefaultNicknameRegex = `^(?P<prefix>\[\S+][ ]*)?(?P<name>[^\[]+)(?P<suffix>[ ]*\[\S+])?`

type Guild struct {
	Job string `alias:"job"`
	ID  string `alias:"id"`

	bot      *Bot
	guild    *discordgo.Guild
	jobRoles map[int32]*discordgo.Role
}

func (g *Guild) init(b *Bot, guild *discordgo.Guild) {
	g.bot = b
	g.guild = guild
	g.jobRoles = map[int32]*discordgo.Role{}
}

func (g *Guild) createJobRoles() error {
	guild, err := g.bot.discord.Guild(g.ID)
	if err != nil {
		return err
	}

	job := g.bot.enricher.GetJobByName(g.Job)
	if job == nil {
		g.bot.logger.Error("unknown job for discord guild, skipping")
		return nil
	}

	for i := len(job.Grades) - 1; i >= 0; i-- {
		grade := job.Grades[i]
		name := fmt.Sprintf(RoleFormat, grade.Grade, grade.Label)

		if utils.InSliceFunc(guild.Roles, func(in *discordgo.Role) bool {
			if in.Name == name {
				g.jobRoles[grade.Grade] = in
				return true
			}
			return false
		}) {
			continue
		}

		role, err := g.bot.discord.GuildRoleCreate(guild.ID, &discordgo.RoleParams{
			Name: name,
		})
		if err != nil {
			return err
		}
		g.jobRoles[grade.Grade] = role
	}

	return nil
}

func (g *Guild) syncUserInfo() error {
	stmt := tOauth2Accs.
		SELECT(
			tOauth2Accs.ExternalID.AS("external_id"),
			tUsers.JobGrade.AS("job_grade"),
			tUsers.Firstname.AS("firstname"),
			tUsers.Lastname.AS("lastname"),
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
			tUsers.Job.EQ(jet.String(g.Job)),
		))

	var dest []struct {
		ExternalID string
		JobGrade   int32
		Firstname  string
		Lastname   string
	}
	if err := stmt.QueryContext(g.bot.ctx, g.bot.db, &dest); err != nil {
		return err
	}

	for _, user := range dest {
		member, err := g.bot.discord.GuildMember(g.ID, user.ExternalID)
		if err != nil {
			return err
		}

		if err := g.setUserNickName(member, user.Firstname, user.Lastname); err != nil {
			return err
		}

		if err := g.setUserJobRole(member, user.JobGrade); err != nil {
			return err
		}
	}

	return nil
}

func (g *Guild) setUserNickName(member *discordgo.Member, firstname string, lastname string) error {
	if g.guild.OwnerID == member.User.ID {
		return nil
	}

	fullName := firstname + " " + lastname

	var nickname string
	if member.Nick != "" {
		match := g.bot.nicknameRegex.FindStringSubmatch(member.Nick)
		result := make(map[string]string)
		for i, name := range g.bot.nicknameRegex.SubexpNames() {
			if i != 0 && name != "" {
				if len(match) >= i {
					result[name] = match[i]
				}
			}
		}

		var ok bool
		nickname, ok = result["name"]
		if !ok {
			g.bot.logger.Warn("failed to extract name from discord nickname", zap.String("dc_nick", member.Nick))
			nickname = member.Nick
		}
	} else {
		nickname = member.User.Username
	}

	extractedName := strings.TrimSpace(nickname)
	if extractedName == fullName {
		return nil
	}

	// Last space on the name is lost due to the space trimming combined with the regex capture
	fullName = g.bot.nicknameRegex.ReplaceAllString(member.Nick, "${prefix}"+fullName+" ${suffix}")
	fullName = strings.TrimSpace(fullName)

	if err := g.bot.discord.GuildMemberNickname(g.ID, member.User.ID, fullName); err != nil {
		return err
	}

	return nil
}

func (g *Guild) setUserJobRole(member *discordgo.Member, grade int32) error {
	r, ok := g.jobRoles[grade]
	if !ok {
		return fmt.Errorf("no role for user's job and grade found")
	}

	found := false
	removeRoles := []string{}
	for _, mr := range member.Roles {
		role, ok := g.findGradeRoleByID(mr)
		if ok {
			if r.ID == role.ID {
				found = true
				continue
			} else {
				removeRoles = append(removeRoles, role.ID)
			}
		}
	}

	if false {
		for _, r := range removeRoles {
			if err := g.bot.discord.GuildMemberRoleRemove(g.ID, member.User.ID, r); err != nil {
				return err
			}
		}
	}

	// Only add user to the rank role if user isn't in it already
	if !found {
		if err := g.bot.discord.GuildMemberRoleAdd(g.ID, member.User.ID, r.ID); err != nil {
			return err
		}
	}

	return nil
}

func (g *Guild) findGradeRoleByID(id string) (*discordgo.Role, bool) {
	for _, j := range g.jobRoles {
		if j.ID == id {
			return j, true
		}
	}

	return nil, false
}
