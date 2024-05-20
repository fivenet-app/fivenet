package modules

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbusers "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type UserInfo struct {
	*BaseModule

	nicknameRegex *regexp.Regexp
	ignoredJobs   []string
}

type userRoleMapping struct {
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
		BaseModule:    base,
		nicknameRegex: nicknameRegex,
		ignoredJobs:   base.appCfg.Get().Discord.IgnoredJobs,
	}, nil
}

func (g *UserInfo) Plan(ctx context.Context) (*types.State, []*discordgo.MessageEmbed, error) {
	job := g.enricher.GetJobByName(g.job)
	if job == nil {
		g.logger.Error("unknown job for discord guild, skipping")
		return nil, nil, nil
	}

	roles, err := g.planRoles(job)
	if err != nil {
		return nil, nil, err
	}

	handlers := []types.UserProcessorHandler{}
	for _, role := range roles {
		if role.Module == userInfoRoleModuleUnemployed {
			handlers = append(handlers, func(ctx context.Context, guildId string, member *discordgo.Member, user *types.User) (*types.User, []*discordgo.MessageEmbed, error) {
				if user.Job == g.job {
					return user, nil, nil
				}

				if slices.Contains(g.ignoredJobs, user.Job) {
					user.Job = g.job
					return user, nil, nil
				}

				switch g.settings.UserInfoSyncSettings.UnemployedMode {
				case pbusers.UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE:
					user.Roles.Sum = append(user.Roles.Sum, role)

				case pbusers.UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_KICK:
					kick := true
					user.Kick = &kick
					user.KickReason = fmt.Sprintf("no longer an employee of %s job (unemployed mode: kick)", g.job)
				}

				return user, nil, nil
			})
			break
		}
	}

	users, logs, err := g.planUsers(ctx, roles)
	if err != nil {
		return nil, logs, err
	}

	return &types.State{
		Roles: roles,
		Users: users,

		UserProcessors: handlers,
	}, logs, err
}

const (
	userInfoRoleModuleEmployee           = "UserInfo-Employees"
	userInfoRoleModuleUnemployed         = "UserInfo-Unemployed"
	userInfoRoleModuleAbsence            = "UserInfo-Absence"
	userInfoRoleModuleJobGradePrefix     = "UserInfo-Grade-"
	userInfoRoleModuleGroupMappingPrefix = "UserInfo-GroupMapping-"
)

func (g *UserInfo) planRoles(job *users.Job) (types.Roles, error) {
	roles := types.Roles{}

	if g.settings.UserInfoSyncSettings.EmployeeRoleEnabled {
		roles = append(roles, &types.Role{
			Name:   strings.ReplaceAll(g.settings.UserInfoSyncSettings.EmployeeRoleFormat, "%s", job.Label),
			Module: userInfoRoleModuleEmployee,
			Job:    g.job,
		})
	}

	if g.settings.UserInfoSyncSettings.UnemployedEnabled {
		roles = append(roles, &types.Role{
			Name:   g.settings.UserInfoSyncSettings.UnemployedRoleName,
			Module: userInfoRoleModuleUnemployed,
			Job:    g.job,
		})
	}

	if g.settings.JobsAbsence {
		roles = append(roles, &types.Role{
			Name:   g.settings.JobsAbsenceSettings.AbsenceRole,
			Module: userInfoRoleModuleAbsence,
			Job:    g.job,
		})
	}

	jobRoles := map[int32]interface{}{}
	for i := len(job.Grades) - 1; i >= 0; i-- {
		grade := job.Grades[i]
		name := strings.ReplaceAll(g.settings.UserInfoSyncSettings.GradeRoleFormat, "%grade_label%", grade.Label)
		name = strings.ReplaceAll(name, "%grade%", fmt.Sprintf("%02d", grade.Grade))

		if _, ok := jobRoles[grade.Grade]; ok {
			continue
		}

		roles = append(roles, &types.Role{
			Name:   name,
			Module: fmt.Sprintf(userInfoRoleModuleJobGradePrefix+"%d", grade.Grade),
			Job:    g.job,
		})

		jobRoles[grade.Grade] = nil
	}

	for i, mapping := range g.settings.UserInfoSyncSettings.GroupMapping {
		roles = append(roles, &types.Role{
			Name:   mapping.Name,
			Module: fmt.Sprintf(userInfoRoleModuleGroupMappingPrefix+"%d", i),
			Job:    g.job,
		})
	}

	return roles, nil
}

func (g *UserInfo) planUsers(ctx context.Context, roles types.Roles) (types.Users, []*discordgo.MessageEmbed, error) {
	users := types.Users{}
	logs := []*discordgo.MessageEmbed{}

	var employeeRole *types.Role
	var absenceRole *types.Role
	gradeRoles := map[int32]*types.Role{}
	groupRoles := map[int]*types.Role{}
	for _, role := range roles {
		if role.Module == userInfoRoleModuleEmployee {
			employeeRole = role
		} else if role.Module == userInfoRoleModuleAbsence {
			absenceRole = role
		} else if strings.HasPrefix(role.Module, userInfoRoleModuleJobGradePrefix) {
			sGrade, found := strings.CutPrefix(role.Module, userInfoRoleModuleJobGradePrefix)
			if !found {
				continue
			}
			grade, err := strconv.Atoi(sGrade)
			if err != nil {
				return nil, nil, err
			}
			gradeRoles[int32(grade)] = role
		} else if strings.HasPrefix(role.Module, userInfoRoleModuleGroupMappingPrefix) {
			sGroup, found := strings.CutPrefix(role.Module, userInfoRoleModuleGroupMappingPrefix)
			if !found {
				continue
			}
			index, err := strconv.Atoi(sGroup)
			if err != nil {
				return nil, nil, err
			}
			groupRoles[index] = role
		}
	}

	jobs := []jet.Expression{jet.String(g.job)}
	for _, job := range g.ignoredJobs {
		jobs = append(jobs, jet.String(job))
	}

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
			tUsers.Job.IN(jobs...),
		)).
		ORDER_BY(tUsers.ID.ASC())

	var dest []*userRoleMapping
	if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return users, logs, err
		}
	}

	errs := multierr.Combine()
	for _, u := range dest {
		user := &types.User{
			ID:    u.ExternalID,
			Roles: &types.UserRoles{},
			Job:   u.Job,
		}

		if u.Job != g.job {
			users.Add(user)
			continue
		}

		member, err := g.discord.GuildMember(g.guild.ID, u.ExternalID)
		if err != nil {
			if restErr, ok := err.(*discordgo.RESTError); ok {
				if restErr.Response.StatusCode == http.StatusNotFound {

					// Add log about employee not being on discord
					logs = append(logs, &discordgo.MessageEmbed{
						Type:        discordgo.EmbedTypeRich,
						Title:       fmt.Sprintf("UserInfo: Employee not found on Discord: %s %s", u.Firstname, u.Lastname),
						Description: fmt.Sprintf("Discord ID: %s, Rank: %d", u.ExternalID, u.JobGrade),
						Author:      embeds.EmbedAuthor,
						Color:       embeds.ColorWarn,
					})
					continue
				}
			}

			errs = multierr.Append(errs, err)
			continue
		}

		if g.settings.UserInfoSyncSettings.SyncNicknames {
			name := g.getUserNickname(member, u.Firstname, u.Lastname)
			if name != "" {
				user.Nickname = &name
			}
		}

		user.Roles.Sum, err = g.getUserRoles(gradeRoles, u.Job, u.JobGrade)
		if err != nil {
			g.logger.Error(fmt.Sprintf("failed to set user's job roles %s", u.ExternalID), zap.Error(err))
			continue
		}

		for idx, mapping := range g.settings.UserInfoSyncSettings.GroupMapping {
			if u.JobGrade < mapping.FromGrade || u.JobGrade > mapping.ToGrade {
				continue
			}

			role, ok := groupRoles[idx]
			if !ok {
				return nil, logs, fmt.Errorf("failed to find role for group mapping %s", mapping.Name)
			}

			user.Roles.Sum = append(user.Roles.Sum, role)
		}

		if g.settings.UserInfoSyncSettings.EmployeeRoleEnabled &&
			employeeRole != nil {
			user.Roles.Sum = append(user.Roles.Sum, employeeRole)
		}

		if g.settings.JobsAbsence && absenceRole != nil &&
			g.isUserAbsent(u.AbsenceBegin, u.AbsenceEnd) {
			user.Roles.Sum = append(user.Roles.Sum, absenceRole)
		}

		users.Add(user)
	}

	return users, logs, errs
}

func (g *UserInfo) getUserNickname(member *discordgo.Member, firstname string, lastname string) string {
	if g.guild.OwnerID == member.User.ID {
		return ""
	}

	fullName := strings.TrimSpace(firstname + " " + lastname)

	nickname := fullName
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
	}

	if strings.TrimSpace(nickname) == fullName {
		return fullName
	}

	// Last space on the name is lost due to the space trimming combined with the regex capture
	fullName = g.nicknameRegex.ReplaceAllString(member.Nick, "${prefix}"+fullName+" ${suffix}")
	fullName = strings.TrimSpace(fullName)

	return fullName
}

func (g *UserInfo) getUserRoles(roles map[int32]*types.Role, job string, grade int32) (types.Roles, error) {
	userRoles := types.Roles{}

	// Ignore certain jobs when syncing (e.g., "temporary" jobs), example:
	// "ambulance" job Discord, and an user is currently in the ignored, e.g., "army", jobs.
	ignoredJobs := g.appCfg.Get().Discord.IgnoredJobs
	if g.job != job && slices.Contains(ignoredJobs, job) {
		return userRoles, nil
	}

	role, ok := roles[grade]
	if !ok {
		return nil, fmt.Errorf("failed to find role for job %s grade %d", job, grade)
	}
	userRoles = append(userRoles, role)

	return userRoles, nil
}

func (g *UserInfo) isUserAbsent(beginDate *timestamp.Timestamp, endDate *timestamp.Timestamp) bool {
	// Either the user has no dates set or the absence is over (due to dates we have to think end date + 24 hours)
	return !((beginDate == nil || endDate == nil) || (time.Since(beginDate.AsTime()) < 0*time.Hour || time.Since(endDate.AsTime()) > 24*time.Hour))
}

func (g *UserInfo) lookupUsersByDiscordID(ctx context.Context, externalId string) ([]*userRoleMapping, error) {
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

	var dest []*userRoleMapping
	if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
