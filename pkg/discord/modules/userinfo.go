package modules

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

const DefaultNicknameRegex = `^(?P<prefix>\[\S+][ ]*)?(?P<name>[^\[]+)(?P<suffix>[ ]*\[\S+])?`

type UserInfo struct {
	*BaseModule

	nicknameRegex      *regexp.Regexp
	roleFormat         string
	employeeRoleFormat string

	jobRoles     map[int32]*discordgo.Role
	employeeRole *discordgo.Role
}

type UserRoleMapping struct {
	AccountID  uint64 `alias:"account_id"`
	ExternalID string `alias:"external_id"`
	JobGrade   int32  `alias:"job_grade"`
	Firstname  string `alias:"firstname"`
	Lastname   string `alias:"lastname"`
	Job        string `alias:"job"`
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
		BaseModule:         base,
		nicknameRegex:      nicknameRegex,
		roleFormat:         base.cfg.UserInfoSync.RoleFormat,
		employeeRoleFormat: base.cfg.UserInfoSync.EmployeeRoleFormat,
		jobRoles:           map[int32]*discordgo.Role{},
	}, nil
}

func (g *UserInfo) Run() error {
	settings, err := g.GetSyncSettings(g.ctx, g.job)
	if err != nil {
		return err
	}
	if settings == nil {
		return nil
	}

	if !settings.UserInfoSync {
		return nil
	}

	if err := g.createJobRoles(); err != nil {
		return err
	}

	return g.syncUserInfo()
}

func (g *UserInfo) syncUserInfo() error {
	stmt := tOauth2Accs.
		SELECT(
			tOauth2Accs.AccountID.AS("userrolemapping.account_id"),
			tOauth2Accs.ExternalID.AS("userrolemapping.external_id"),
			tUsers.JobGrade.AS("userrolemapping.job_grade"),
			tUsers.Firstname.AS("userrolemapping.firstname"),
			tUsers.Lastname.AS("userrolemapping.lastname"),
			tUsers.Job.AS("userrolemapping.job"),
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
		)).
		ORDER_BY(tUsers.ID.ASC())

	var dest []*UserRoleMapping
	if err := stmt.QueryContext(g.ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	for _, user := range dest {
		member, err := g.discord.GuildMember(g.guild.ID, user.ExternalID)
		if err != nil {
			if restErr, ok := err.(*discordgo.RESTError); ok {
				if restErr.Response.StatusCode == http.StatusNotFound {
					g.logger.Warn("user not found on job discord server",
						zap.String("discord_user_id", user.ExternalID), zap.Uint64("account_id", user.AccountID), zap.String("user", fmt.Sprintf("%s, %s", user.Firstname, user.Lastname)), zap.Int32("job_grade", user.JobGrade))
					continue
				}
			}
			return err
		}

		if err := g.setUserNickName(member, user.Firstname, user.Lastname); err != nil {
			g.logger.Error(fmt.Sprintf("failed to set user's nickname %s", user.ExternalID), zap.Error(err))
			continue
		}

		if err := g.setUserJobRole(member, user.Job, user.JobGrade); err != nil {
			g.logger.Error(fmt.Sprintf("failed to set user's job roles %s", user.ExternalID), zap.Error(err))
			continue
		}
	}

	return g.cleanupUserJobRoles(dest)
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

		if slices.ContainsFunc(guild.Roles, func(in *discordgo.Role) bool {
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
			return fmt.Errorf("failed to create job grade role %s (Grade: %s): %w", name, grade, err)
		}

		g.jobRoles[grade.Grade] = role
	}

	employeeRoleName := fmt.Sprintf(g.employeeRoleFormat, job.Label)
	if !slices.ContainsFunc(guild.Roles, func(in *discordgo.Role) bool {
		if in.Name == employeeRoleName {
			g.employeeRole = in
			return true
		}
		return false
	}) {
		role, err := g.discord.GuildRoleCreate(guild.ID, &discordgo.RoleParams{
			Name: employeeRoleName,
		})
		if err != nil {
			return fmt.Errorf("failed to create employee role %s: %w", job.Label, err)
		}
		g.employeeRole = role
	}

	g.logger.Debug("created job employee and rank roles")

	return nil
}

func (g *UserInfo) setUserJobRole(member *discordgo.Member, job string, grade int32) error {
	// Ignore certain jobs when syncing (e.g., "temporary" jobs), example:
	// "ambulance" job Discord, and an user is currently in the ignored "army" job.
	if g.job != job && slices.Contains(g.cfg.UserInfoSync.IgnoreJobs, job) {
		return nil
	}

	r, ok := g.jobRoles[grade]
	if !ok {
		return fmt.Errorf("no role for user's job and grade %d found", grade)
	}

	logger := g.logger.With(zap.String("discord_user_id", member.User.ID), zap.String("discord_nickname", member.Nick))

	hasEmployeeRole := false
	found := false
	removeRoles := []*discordgo.Role{}
	for _, mr := range member.Roles {
		if mr == g.employeeRole.ID {
			hasEmployeeRole = true
			continue
		}

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

	for _, role := range removeRoles {
		logger.Debug("removing role from member", zap.String("discord_role_name", role.Name), zap.String("discord_role_id", role.ID))
		if err := g.discord.GuildMemberRoleRemove(g.guild.ID, member.User.ID, role.ID); err != nil {
			return fmt.Errorf("failed to remove role %s (%s) from member %s: %w", role.Name, role.ID, member.User.ID, err)
		}
	}

	if !hasEmployeeRole {
		logger.Debug("adding employee role to member", zap.String("discord_role_name", g.employeeRole.Name), zap.String("discord_role_id", g.employeeRole.ID))
		if err := g.discord.GuildMemberRoleAdd(g.guild.ID, member.User.ID, g.employeeRole.ID); err != nil {
			return fmt.Errorf("failed to add employee role %s (%s) member %s: %w", r.Name, g.employeeRole.ID, member.User.ID, err)
		}
	}

	// Only add user to the grade role if user isn't in it already
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

func (g *UserInfo) cleanupUserJobRoles(users []*UserRoleMapping) error {
	guild, err := g.discord.State.Guild(g.guild.ID)
	if err != nil {
		return err
	}

outerLoop:
	for i := 0; i < len(guild.Members); i++ {
		member := guild.Members[i]
		for _, role := range g.jobRoles {
			// If user isn't in one of the synced job roles, continue
			if !slices.Contains(guild.Members[i].Roles, role.ID) {
				continue
			}

			// Check if user is suposed to have that job grade role
			if slices.ContainsFunc(users, func(in *UserRoleMapping) bool {
				r, ok := g.findGradeRoleByID(role.ID)
				if in.ExternalID == member.User.ID && (ok && r.ID == role.ID) {
					return true
				}
				return false
			}) {
				continue
			}

			// Lookup user in database
			userMappings, err := g.lookupUsersByDiscordI(member.User.ID)
			if err != nil {
				g.logger.Error("failed to lookup fivenet account via discord id", zap.String("discord_role_name", role.Name), zap.String("discord_role_id", role.ID),
					zap.String("discord_user_id", member.User.ID), zap.String("discord_nickname", member.Nick))
				continue
			}

			for _, userMapping := range userMappings {
				// Ignore certain jobs when syncing (e.g., "temporary" jobs)
				if g.job != userMapping.Job && slices.Contains(g.cfg.UserInfoSync.IgnoreJobs, userMapping.Job) {
					continue outerLoop
				}
			}

			g.logger.Debug("removing job grade role from member", zap.String("discord_role_name", role.Name), zap.String("discord_role_id", role.ID),
				zap.String("discord_user_id", member.User.ID), zap.String("discord_nickname", member.Nick))
			if err := g.discord.GuildMemberRoleRemove(g.guild.ID, member.User.ID, role.ID); err != nil {
				return fmt.Errorf("failed to remove member from role %s (%s): %w", role.Name, role.ID, err)
			}
		}

		role := g.employeeRole
		// If user isn't in the employee role, continue
		if !slices.Contains(guild.Members[i].Roles, role.ID) {
			continue
		}

		// Check if member is suposed to have the employee role
		if slices.ContainsFunc(users, func(in *UserRoleMapping) bool {
			return in.ExternalID == member.User.ID
		}) {
			continue
		}

		g.logger.Debug("removing employee role from member", zap.String("discord_role_name", g.employeeRole.Name), zap.String("discord_role_id", g.employeeRole.ID),
			zap.String("discord_user_id", member.User.ID), zap.String("discord_nickname", member.Nick))
		if err := g.discord.GuildMemberRoleRemove(g.guild.ID, member.User.ID, role.ID); err != nil {
			return fmt.Errorf("failed to remove member from employee job role %s (%s): %w", role.Name, role.ID, err)
		}
	}

	return nil
}

func (g *UserInfo) lookupUsersByDiscordI(externalId string) ([]*UserRoleMapping, error) {
	stmt := tOauth2Accs.
		SELECT(
			tOauth2Accs.AccountID.AS("userrolemapping.account_id"),
			tOauth2Accs.ExternalID.AS("userrolemapping.external_id"),
			tUsers.JobGrade.AS("userrolemapping.job_grade"),
			tUsers.Firstname.AS("userrolemapping.firstname"),
			tUsers.Lastname.AS("userrolemapping.lastname"),
			tUsers.Job.AS("userrolemapping.job"),
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
			tOauth2Accs.ExternalID.EQ(jet.String(externalId)),
		)).
		ORDER_BY(tUsers.ID.ASC())

	var dest []*UserRoleMapping
	if err := stmt.QueryContext(g.ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
