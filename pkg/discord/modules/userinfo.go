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
		BaseModule:    base,
		nicknameRegex: nicknameRegex,
	}, nil
}

func (g *UserInfo) Plan(ctx context.Context) (*types.Plan, []*discordgo.MessageEmbed, error) {
	job := g.enricher.GetJobByName(g.job)
	if job == nil {
		g.logger.Error("unknown job for discord guild, skipping")
		return nil, nil, nil
	}

	roles, err := g.planRoles(job)
	if err != nil {
		return nil, nil, err
	}

	handlers := []types.NotPartOfFactionHandler{}
	for _, role := range roles {
		if role.Module == userInfoRoleModuleUnemployed {
			handlers = append(handlers, func(ctx context.Context, session *discordgo.Session, guildId string, member *discordgo.Member) ([]*discordgo.MessageEmbed, error) {
				if g.settings.DryRun {
					return nil, nil
				}

				switch g.settings.UserInfoSyncSettings.UnemployedMode {
				case pbusers.UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE:
					// Skip if user is already part of unemployed role
					if slices.Contains(member.Roles, role.ID) {
						break
					}

					if err := g.discord.GuildMemberRoleAdd(g.guild.ID, member.User.ID, role.ID, discordgo.WithContext(ctx)); err != nil {
						return nil, fmt.Errorf("failed to add member to unemployed role %s: %w", role.ID, err)
					}

				case pbusers.UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_KICK:
					if err := g.discord.GuildMemberDeleteWithReason(g.guild.ID, member.User.ID,
						fmt.Sprintf("no longer an employee of %s job (unemployed mode: kick)", g.job),
						discordgo.WithContext(ctx)); err != nil {
						return nil, fmt.Errorf("failed to kick unemployed member %s from guild: %w", member.User.ID, err)
					}
				}

				return nil, nil
			})
			break
		}
	}

	users, logs, err := g.planUsers(ctx, roles)
	if err != nil {
		return nil, logs, err
	}

	return &types.Plan{
		Roles:                    roles,
		Users:                    users,
		NotPartOfFactionHandlers: handlers,
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
		})
	}

	if g.settings.UserInfoSyncSettings.UnemployedEnabled {
		roles = append(roles, &types.Role{
			Name:   g.settings.UserInfoSyncSettings.UnemployedRoleName,
			Module: userInfoRoleModuleUnemployed,
		})
	}

	if g.settings.JobsAbsence {
		roles = append(roles, &types.Role{
			Name:   g.settings.JobsAbsenceSettings.AbsenceRole,
			Module: userInfoRoleModuleAbsence,
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
		})

		jobRoles[grade.Grade] = nil
	}

	for i, mapping := range g.settings.UserInfoSyncSettings.GroupMapping {
		roles = append(roles, &types.Role{
			Name:   mapping.Name,
			Module: fmt.Sprintf(userInfoRoleModuleGroupMappingPrefix+"%d", i),
		})
	}

	return roles, nil
}

func (g *UserInfo) planUsers(ctx context.Context, roles types.Roles) (types.Users, []*discordgo.MessageEmbed, error) {
	users := types.Users{}
	logs := []*discordgo.MessageEmbed{}

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
	if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return users, logs, err
		}
	}

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

	errs := multierr.Combine()
	for _, u := range dest {
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

		user := &types.User{
			ID: u.ExternalID,
		}

		if g.settings.UserInfoSyncSettings.SyncNicknames {
			name := g.getUserNickname(member, u.Firstname, u.Lastname)
			if name != "" {
				user.Nickname = &name
			}
		}

		user.Roles, err = g.getUserRoles(gradeRoles, u.Job, u.JobGrade)
		if err != nil {
			g.logger.Error(fmt.Sprintf("failed to set user's job roles %s", u.ExternalID), zap.Error(err))
			continue
		}

		for idx, mapping := range g.settings.UserInfoSyncSettings.GroupMapping {
			if mapping.FromGrade < u.JobGrade || mapping.ToGrade > u.JobGrade {
				continue
			}

			role, ok := groupRoles[idx]
			if !ok {
				return nil, logs, fmt.Errorf("failed to find role for group mapping %s", mapping.Name)
			}

			user.Roles = append(user.Roles, role)
		}

		if g.settings.UserInfoSyncSettings.EmployeeRoleEnabled &&
			employeeRole != nil {
			user.Roles = append(user.Roles, employeeRole)
		}

		if g.settings.JobsAbsence && absenceRole != nil &&
			g.isUserAbsent(u.AbsenceBegin, u.AbsenceEnd) {
			user.Roles = append(user.Roles, absenceRole)
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

func (g *UserInfo) getUserRoles(roles map[int32]*types.Role, job string, grade int32) ([]*types.Role, error) {
	userRoles := []*types.Role{}

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
