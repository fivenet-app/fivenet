package modules

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	pbusers "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type UserInfo struct {
	*BaseModule

	nicknameRegex   *regexp.Regexp
	gradeRoleFormat string

	employeeRoleEnabled bool
	employeeRoleFormat  string
	unemployedRoleName  string
	unemployedMode      pbusers.UserInfoSyncUnemployedMode

	jobsAbsenceEnabled  bool
	jobsAbsenceRoleName string

	jobRoles        map[int32]*discordgo.Role
	employeeRole    *discordgo.Role
	unemployedRole  *discordgo.Role
	jobsAbsenceRole *discordgo.Role
	ignoredJobs     []string
}

type UserRoleMapping struct {
	AccountID    uint64               `alias:"account_id"`
	ExternalID   string               `alias:"external_id"`
	JobGrade     int32                `alias:"job_grade"`
	Firstname    string               `alias:"firstname"`
	Lastname     string               `alias:"lastname"`
	Job          string               `alias:"job"`
	AbsenceBegin *timestamp.Timestamp `alias:"absence_begin"`
	AbsenceEnd   *timestamp.Timestamp `alias:"absence_end"`
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
		BaseModule:          base,
		nicknameRegex:       nicknameRegex,
		employeeRoleEnabled: true,
		employeeRoleFormat:  base.cfg.UserInfoSync.EmployeeRoleFormat,
		gradeRoleFormat:     base.cfg.UserInfoSync.GradeRoleFormat,
		unemployedRoleName:  base.cfg.UserInfoSync.UnemployedRoleName,
		ignoredJobs:         []string{},

		jobRoles: map[int32]*discordgo.Role{},
	}, nil
}

func (g *UserInfo) Run(settings *pbusers.DiscordSyncSettings) ([]*discordgo.MessageEmbed, error) {
	if !settings.UserInfoSync {
		return nil, nil
	}

	g.employeeRoleEnabled = settings.UserInfoSyncSettings.EmployeeRoleEnabled
	g.employeeRoleFormat = settings.UserInfoSyncSettings.EmployeeRoleFormat
	g.gradeRoleFormat = settings.UserInfoSyncSettings.GradeRoleFormat
	g.unemployedRoleName = settings.UserInfoSyncSettings.UnemployedRoleName
	g.unemployedMode = settings.UserInfoSyncSettings.UnemployedMode

	g.jobsAbsenceEnabled = settings.JobsAbsence
	g.jobsAbsenceRoleName = settings.JobsAbsenceSettings.AbsenceRole

	guildRoles, err := g.discord.GuildRoles(g.guild.ID)
	if err != nil {
		return nil, err
	}

	logs, err := g.createJobRoles(guildRoles, settings.UserInfoSyncSettings.UnemployedEnabled)
	if err != nil {
		return nil, err
	}

	return g.syncUserInfo(logs)
}

func (g *UserInfo) syncUserInfo(logs []*discordgo.MessageEmbed) ([]*discordgo.MessageEmbed, error) {
	stmt := tOauth2Accs.
		SELECT(
			tOauth2Accs.AccountID.AS("userrolemapping.account_id"),
			tOauth2Accs.ExternalID.AS("userrolemapping.external_id"),
			tUsers.JobGrade.AS("userrolemapping.job_grade"),
			tUsers.Firstname.AS("userrolemapping.firstname"),
			tUsers.Lastname.AS("userrolemapping.lastname"),
			tUsers.Job.AS("userrolemapping.job"),
			tJobsUserProps.AbsenceBegin.AS("userrolemapping.absence_begin"),
			tJobsUserProps.AbsenceEnd.AS("userrolemapping.absence_end"),
		).
		FROM(
			tOauth2Accs.
				INNER_JOIN(tAccs,
					tAccs.ID.EQ(tOauth2Accs.AccountID),
				).
				INNER_JOIN(tUsers,
					tUsers.Identifier.LIKE(jet.CONCAT(jet.String("char%:"), tAccs.License)),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUsers.ID),
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
			return logs, err
		}
	}

	errs := multierr.Combine()
	for _, user := range dest {
		member, err := g.discord.GuildMember(g.guild.ID, user.ExternalID)
		if err != nil {
			if restErr, ok := err.(*discordgo.RESTError); ok {
				if restErr.Response.StatusCode == http.StatusNotFound {
					g.logger.Warn("user not found on job discord server",
						zap.String("discord_user_id", user.ExternalID), zap.Uint64("account_id", user.AccountID), zap.String("user", fmt.Sprintf("%s, %s", user.Firstname, user.Lastname)), zap.Int32("job_grade", user.JobGrade))

					// Add log about user not being on discord
					logs = append(logs, &discordgo.MessageEmbed{
						Type:        discordgo.EmbedTypeRich,
						Title:       fmt.Sprintf("UserInfo: Employee not found on Discord: %s, %s", user.Firstname, user.Lastname),
						Description: fmt.Sprintf("Discord ID: %s, Rank: %d", user.ExternalID, user.JobGrade),
						Author:      embeds.EmbedAuthor,
						Color:       embeds.ColorWarn,
					})
					continue
				}
			}

			return logs, err
		}

		if err := g.setUserNickname(member, user.Firstname, user.Lastname); err != nil {
			g.logger.Error(fmt.Sprintf("failed to set user's nickname %s", user.ExternalID), zap.Error(err))
			continue
		}

		if err := g.setUserJobRole(member, user.Job, user.JobGrade); err != nil {
			g.logger.Error(fmt.Sprintf("failed to set user's job roles %s", user.ExternalID), zap.Error(err))
			continue
		}

		if g.jobsAbsenceRole != nil {
			if err := g.setJobsAbsenceRole(member, user.AbsenceBegin, user.AbsenceEnd); err != nil {
				g.logger.Error(fmt.Sprintf("failed to set user's jobs absence roles %s", user.ExternalID), zap.Error(err))
				continue
			}
		}
	}

	if err := g.cleanupUserJobRoles(g.guild, dest); err != nil {
		errs = multierr.Append(errs, err)
	}

	return logs, errs
}

func (g *UserInfo) setUserNickname(member *discordgo.Member, firstname string, lastname string) error {
	if g.guild.OwnerID == member.User.ID {
		return nil
	}

	fullName := strings.TrimSpace(firstname + " " + lastname)

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

	if strings.TrimSpace(nickname) == fullName {
		return nil
	}

	// Last space on the name is lost due to the space trimming combined with the regex capture
	fullName = g.nicknameRegex.ReplaceAllString(member.Nick, "${prefix}"+fullName+" ${suffix}")
	fullName = strings.TrimSpace(fullName)

	if err := g.discord.GuildMemberNickname(g.guild.ID, member.User.ID, fullName); err != nil {
		return fmt.Errorf("failed to update user %s (%s) nickname: %w", fullName, member.User.ID, err)
	}

	g.logger.Debug("update discord nickname", zap.String("dc_current_nick", member.Nick), zap.String("dc_user_id", member.User.ID), zap.String("dc_new_nick", fullName))

	return nil
}

func (g *UserInfo) createJobRoles(roles []*discordgo.Role, unemployedRoleEnabled bool) ([]*discordgo.MessageEmbed, error) {
	logs := []*discordgo.MessageEmbed{}

	job := g.enricher.GetJobByName(g.job)
	if job == nil {
		g.logger.Error("unknown job for discord guild, skipping")
		return logs, nil
	}

	for i := len(job.Grades) - 1; i >= 0; i-- {
		grade := job.Grades[i]
		name := strings.ReplaceAll(g.gradeRoleFormat, "%grade_label%", grade.Label)
		name = strings.ReplaceAll(name, "%grade%", fmt.Sprintf("%02d", grade.Grade))

		if slices.ContainsFunc(roles, func(in *discordgo.Role) bool {
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

		role, err := g.discord.GuildRoleCreate(g.guild.ID, &discordgo.RoleParams{
			Name: name,
		})
		if err != nil {
			return logs, fmt.Errorf("failed to create job grade role %s (Grade: %s): %w", name, grade, err)
		}

		// Add log about user not being on discord
		logs = append(logs, &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       fmt.Sprintf("UserInfo: Role %s created", role.Name),
			Description: fmt.Sprintf("Role %s created (%s)", role.Name, role.ID),
			Author:      embeds.EmbedAuthor,
			Color:       embeds.ColorInfo,
		})

		g.jobRoles[grade.Grade] = role
	}

	if g.employeeRoleEnabled {
		employeeRoleName := strings.ReplaceAll(g.employeeRoleFormat, "%s", job.Label)
		if !slices.ContainsFunc(roles, func(in *discordgo.Role) bool {
			if in.Name == employeeRoleName {
				g.employeeRole = in
				return true
			}
			return false
		}) {
			role, err := g.discord.GuildRoleCreate(g.guild.ID, &discordgo.RoleParams{
				Name: employeeRoleName,
			})
			if err != nil {
				return logs, fmt.Errorf("failed to create employee role %s: %w", job.Label, err)
			}

			// Add log about user not being on discord
			logs = append(logs, &discordgo.MessageEmbed{
				Type:        discordgo.EmbedTypeRich,
				Title:       fmt.Sprintf("UserInfo: Role %s created", role.Name),
				Description: fmt.Sprintf("Role %s created (%s)", role.Name, role.ID),
				Author:      embeds.EmbedAuthor,
				Color:       embeds.ColorInfo,
			})

			g.employeeRole = role
		}
	} else {
		g.employeeRole = nil
	}

	if g.jobsAbsenceEnabled {
		if !slices.ContainsFunc(roles, func(in *discordgo.Role) bool {
			if in.Name == g.jobsAbsenceRoleName {
				g.jobsAbsenceRole = in
				return true
			}
			return false
		}) {
			role, err := g.discord.GuildRoleCreate(g.guild.ID, &discordgo.RoleParams{
				Name: g.jobsAbsenceRoleName,
			})
			if err != nil {
				return logs, fmt.Errorf("failed to create jobs absence role %s: %w", job.Label, err)
			}

			// Add log about user not being on discord
			logs = append(logs, &discordgo.MessageEmbed{
				Type:        discordgo.EmbedTypeRich,
				Title:       fmt.Sprintf("UserInfo: Role %s created", role.Name),
				Description: fmt.Sprintf("Role %s created (%s)", role.Name, role.ID),
				Author:      embeds.EmbedAuthor,
				Color:       embeds.ColorInfo,
			})

			g.jobsAbsenceRole = role
		}
	} else {
		g.jobsAbsenceRole = nil
	}

	if unemployedRoleEnabled {
		if err := g.createUnemployedRole(roles); err != nil {
			return nil, err
		}
	}

	g.logger.Debug("created job employee and rank roles")

	return logs, nil
}

func (g *UserInfo) createUnemployedRole(roles []*discordgo.Role) error {
	if !slices.ContainsFunc(roles, func(in *discordgo.Role) bool {
		if in.Name == g.unemployedRoleName {
			g.unemployedRole = in
			return true
		}
		return false
	}) {
		role, err := g.discord.GuildRoleCreate(g.guild.ID, &discordgo.RoleParams{
			Name: g.unemployedRoleName,
		})
		if err != nil {
			return fmt.Errorf("failed to create employee role %s: %w", g.unemployedRoleName, err)
		}

		g.unemployedRole = role
	}

	return nil
}

func (g *UserInfo) setUserJobRole(member *discordgo.Member, job string, grade int32) error {
	// Ignore certain jobs when syncing (e.g., "temporary" jobs), example:
	// "ambulance" job Discord, and an user is currently in the ignored, e.g., "army", jobs.
	if g.job != job && slices.Contains(g.ignoredJobs, job) {
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
		if g.employeeRole == nil || mr == g.employeeRole.ID {
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

	// If unemployed role is enabled and user is an employee, make sure to remove the unemployed role
	if g.unemployedRole != nil && slices.Contains(member.Roles, g.unemployedRole.ID) {
		removeRoles = append(removeRoles, g.unemployedRole)
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

func (g *UserInfo) setJobsAbsenceRole(member *discordgo.Member, beginDate *timestamp.Timestamp, endDate *timestamp.Timestamp) error {
	// Either the user has no dates set or the absence is over (due to dates we have to think end date + 24 hours)
	if (beginDate == nil || endDate == nil) || (time.Since(beginDate.AsTime()) < 0*time.Hour || time.Since(endDate.AsTime()) > 24*time.Hour) {
		if !slices.Contains(member.Roles, g.jobsAbsenceRole.ID) {
			return nil
		}

		if err := g.discord.GuildMemberRoleRemove(g.guild.ID, member.User.ID, g.jobsAbsenceRole.ID); err != nil {
			return fmt.Errorf("failed to remove jobs absence role %s (%s) from member %s: %w", g.jobsAbsenceRole.Name, g.jobsAbsenceRole.ID, member.User.ID, err)
		}
	} else {
		if slices.Contains(member.Roles, g.jobsAbsenceRole.ID) {
			return nil
		}

		if err := g.discord.GuildMemberRoleAdd(g.guild.ID, member.User.ID, g.jobsAbsenceRole.ID); err != nil {
			return fmt.Errorf("failed to add role %s (%s) member %s: %w", g.jobsAbsenceRole.Name, g.jobsAbsenceRole.ID, member.User.ID, err)
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

func (g *UserInfo) cleanupUserJobRoles(guild *discordgo.Guild, users []*UserRoleMapping) error {
	errs := multierr.Combine()

outerLoop:
	for i := 0; i < len(guild.Members); i++ {
		member := guild.Members[i]
		if member.User.Bot {
			continue
		}

		isEmployee := false
		for _, role := range g.jobRoles {
			// If user isn't in one of the synced job roles, continue
			if !slices.Contains(member.Roles, role.ID) {
				continue
			}

			// Check if user is suposed to have that job grade role
			if slices.ContainsFunc(users, func(in *UserRoleMapping) bool {
				r, ok := g.findGradeRoleByID(role.ID)
				if in.ExternalID == member.User.ID && ok && r.ID == role.ID {
					return true
				}
				return false
			}) {
				isEmployee = true
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
				if g.job != userMapping.Job && slices.Contains(g.ignoredJobs, userMapping.Job) {
					continue outerLoop
				}
			}

			g.logger.Debug("removing job grade role from member", zap.String("discord_role_name", role.Name), zap.String("discord_role_id", role.ID),
				zap.String("discord_user_id", member.User.ID), zap.String("discord_nickname", member.Nick))
			if err := g.discord.GuildMemberRoleRemove(g.guild.ID, member.User.ID, role.ID); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to remove member from role %s (%s): %w", role.Name, role.ID, err))
				continue
			}
		}

		if !isEmployee {
			// If employee role is disabled or user isn't in the employee role, continue
			if g.employeeRole != nil && slices.Contains(member.Roles, g.employeeRole.ID) {
				g.logger.Debug("removing employee role from member", zap.String("discord_role_name", g.employeeRole.Name), zap.String("discord_role_id", g.employeeRole.ID),
					zap.String("discord_user_id", member.User.ID), zap.String("discord_nickname", member.Nick))
				if err := g.discord.GuildMemberRoleRemove(g.guild.ID, member.User.ID, g.employeeRole.ID); err != nil {
					errs = multierr.Append(errs, fmt.Errorf("failed to remove member from employee job role %s (%s): %w", g.employeeRole.Name, g.employeeRole.ID, err))
					continue
				}
			}

			// If unemployed role is enabled and user is not an employee, give them the unemployed role
			if g.unemployedRole != nil {
				switch g.unemployedMode {
				case pbusers.UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE:
					// Skip if user is already part of unemployed role
					if slices.Contains(member.Roles, g.unemployedRole.ID) {
						break
					}

					g.logger.Debug("adding unemployed role to member", zap.String("discord_role_name", g.unemployedRole.Name), zap.String("discord_role_id", g.unemployedRole.ID),
						zap.String("discord_user_id", member.User.ID), zap.String("discord_nickname", member.Nick))
					if err := g.discord.GuildMemberRoleAdd(g.guild.ID, member.User.ID, g.unemployedRole.ID); err != nil {
						errs = multierr.Append(errs, fmt.Errorf("failed to add member to unemployed role %s: %w", g.unemployedRole.ID, err))
						continue
					}

				case pbusers.UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_KICK:
					if err := g.discord.GuildMemberDeleteWithReason(g.guild.ID, member.User.ID,
						fmt.Sprintf("no longer an employee of %s job", g.job)); err != nil {
						errs = multierr.Append(errs, fmt.Errorf("failed to kick unemployed member %s from guild: %w", member.User.ID, err))
						continue
					}
				}
			}

			// Remove user from jobs absence role when no longer an employee
			if g.jobsAbsenceRole != nil && slices.Contains(member.Roles, g.jobsAbsenceRole.ID) {
				g.logger.Debug("removing jobs absence role from unemployed member", zap.String("discord_role_name", g.jobsAbsenceRole.Name), zap.String("discord_role_id", g.jobsAbsenceRole.ID),
					zap.String("discord_user_id", member.User.ID), zap.String("discord_nickname", member.Nick))
				if err := g.discord.GuildMemberRoleRemove(g.guild.ID, member.User.ID, g.jobsAbsenceRole.ID); err != nil {
					errs = multierr.Append(errs, fmt.Errorf("failed to remove member from jobs absence role %s (%s): %w", g.jobsAbsenceRole.Name, g.jobsAbsenceRole.ID, err))
					continue
				}
			}
		}
	}

	return errs
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
