package modules

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

const DefaultNicknameRegex = `^(?P<prefix>\[\S+][ ]*)?(?P<name>[^\[]+)(?P<suffix>[ ]*\[\S+])?`

type UserInfo struct {
	*BaseModule

	nicknameRegex *regexp.Regexp
	roleFormat    string
	jobRoles      map[int32]*discordgo.Role
}

func init() {
	Modules["userinfo"] = NewUserInfo
}

func NewUserInfo(base *BaseModule) (Module, error) {
	nicknameRegex, err := regexp.Compile(base.cfg.UserInfoSync.NicknameRegex)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		BaseModule:    base,
		nicknameRegex: nicknameRegex,
		roleFormat:    base.cfg.UserInfoSync.RoleFormat,
		jobRoles:      map[int32]*discordgo.Role{},
	}, nil
}

func (g *UserInfo) Run() error {
	if err := g.createJobRoles(); err != nil {
		return err
	}

	return g.syncUserInfo()
}

func (g *UserInfo) syncUserInfo() error {
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
			tUsers.Job.EQ(jet.String(g.job)),
		))

	var dest []struct {
		ExternalID string
		JobGrade   int32
		Firstname  string
		Lastname   string
	}
	if err := stmt.QueryContext(g.ctx, g.db, &dest); err != nil {
		return err
	}

	for _, user := range dest {
		member, err := g.discord.GuildMember(g.guild.ID, user.ExternalID)
		if err != nil {
			if restErr, ok := err.(*discordgo.RESTError); ok {
				if restErr.Response.StatusCode == http.StatusNotFound {
					g.logger.Warn("user not found on job discord server",
						zap.String("job", g.job), zap.String("user", fmt.Sprintf("%s, %s (%d)", user.Firstname, user.Lastname, user.JobGrade)))
					continue
				}
			}
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

func (g *UserInfo) setUserNickName(member *discordgo.Member, firstname string, lastname string) error {
	if g.guild.OwnerID == member.User.ID {
		return nil
	}

	fullName := firstname + " " + lastname

	var nickname string
	if member.Nick != "" {
		match := g.nicknameRegex.FindStringSubmatch(member.Nick)
		result := make(map[string]string)
		for i, name := range g.nicknameRegex.SubexpNames() {
			if i != 0 && name != "" {
				if len(match) >= i {
					result[name] = match[i]
				}
			}
		}

		var ok bool
		nickname, ok = result["name"]
		if !ok {
			g.logger.Warn("failed to extract name from discord nickname", zap.String("dc_nick", member.Nick))
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
	fullName = g.nicknameRegex.ReplaceAllString(member.Nick, "${prefix}"+fullName+" ${suffix}")
	fullName = strings.TrimSpace(fullName)

	if err := g.discord.GuildMemberNickname(g.guild.ID, member.User.ID, fullName); err != nil {
		return fmt.Errorf("failed to update user %s (%s) nickname: %w", fullName, member.User.ID, err)
	}

	return nil
}

func (g *UserInfo) createJobRoles() error {
	guild, err := g.discord.Guild(g.guild.ID)
	if err != nil {
		return err
	}

	job := g.enricher.GetJobByName(g.job)
	if job == nil {
		g.logger.Error("unknown job for discord guild, skipping")
		return nil
	}

	for i := len(job.Grades) - 1; i >= 0; i-- {
		grade := job.Grades[i]
		name := fmt.Sprintf(g.roleFormat, grade.Grade, grade.Label)

		if utils.InSliceFunc(guild.Roles, func(in *discordgo.Role) bool {
			if in.Name == name {
				g.jobRoles[grade.Grade] = in
				return true
			}
			return false
		}) {
			continue
		}

		if _, ok := g.jobRoles[grade.Grade]; ok {
			continue
		}

		role, err := g.discord.GuildRoleCreate(guild.ID, &discordgo.RoleParams{
			Name: name,
		})
		if err != nil {
			return fmt.Errorf("failed to create role %s (Grade: %s): %w", name, grade, err)
		}

		g.jobRoles[grade.Grade] = role
	}

	return nil
}

func (g *UserInfo) setUserJobRole(member *discordgo.Member, grade int32) error {
	r, ok := g.jobRoles[grade]
	if !ok {
		return fmt.Errorf("no role for user's job and grade %d found", grade)
	}

	found := false
	removeRoles := []*discordgo.Role{}
	for _, mr := range member.Roles {
		role, ok := g.findGradeRoleByID(mr)
		if ok {
			if r.ID == role.ID {
				found = true
				continue
			} else {
				removeRoles = append(removeRoles, role)
			}
		}
	}

	if false {
		for _, role := range removeRoles {
			if err := g.discord.GuildMemberRoleRemove(g.guild.ID, member.User.ID, role.ID); err != nil {
				return fmt.Errorf("failed to remove role %s (%s) member %s: %w", role.Name, role.ID, member.User.ID, err)
			}
		}
	}

	// Only add user to the rank role if user isn't in it already
	if !found {
		if err := g.discord.GuildMemberRoleAdd(g.guild.ID, member.User.ID, r.ID); err != nil {
			return fmt.Errorf("failed to add role %s (%s) member %s: %w", r.Name, r.ID, member.User.ID, err)
		}
	}

	return nil
}

func (g *UserInfo) findGradeRoleByID(id string) (*discordgo.Role, bool) {
	for _, j := range g.jobRoles {
		if j.ID == id {
			return j, true
		}
	}

	return nil, false
}
