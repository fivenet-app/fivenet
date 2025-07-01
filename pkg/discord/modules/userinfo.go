package modules

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	userInfoRoleModuleEmployee           = "UserInfo-Employees"
	userInfoRoleModuleUnemployed         = "UserInfo-Unemployed"
	userInfoRoleModuleAbsence            = "UserInfo-Absence"
	userInfoRoleModuleJobGradePrefix     = "UserInfo-Grade-"
	userInfoRoleModuleGroupMappingPrefix = "UserInfo-GroupMapping-"

	DiscordNicknameMaxLength = 32
)

type UserInfo struct {
	*BaseModule

	mu sync.Mutex

	employeeRole   *types.Role
	unemployedRole *types.Role
	absenceRole    *types.Role
	jobGradeRoles  map[int32]*types.Role
	groupRoles     map[string]*types.Role
}

type userRoleMapping struct {
	AccountID  uint64 `alias:"account_id"`
	ExternalID string `alias:"external_id"`
	UserID     int32  `alias:"user_id"`
	Job        string `alias:"job"`
	JobGrade   int32  `alias:"job_grade"`
	Firstname  string `alias:"firstname"`
	Lastname   string `alias:"lastname"`

	// Job Props
	NamePrefix   string               `alias:"name_prefix"`
	NameSuffix   string               `alias:"name_suffix"`
	AbsenceBegin *timestamp.Timestamp `alias:"absence_begin"`
	AbsenceEnd   *timestamp.Timestamp `alias:"absence_end"`
}

func init() {
	Modules["userinfo"] = NewUserInfo
}

func NewUserInfo(base *BaseModule, events *broker.Broker[any]) (Module, error) {
	ui := &UserInfo{
		BaseModule: base,

		mu: sync.Mutex{},

		jobGradeRoles: map[int32]*types.Role{},
		groupRoles:    map[string]*types.Role{},
	}

	eventsCh := events.Subscribe()
	go func() {
		select {
		case <-base.ctx.Done():
			events.Unsubscribe(eventsCh)
			return

		case ev := <-eventsCh:
			ui.watchEvents(ev)
		}
	}()

	return ui, nil
}

func (g *UserInfo) GetName() string {
	return "userinfo"
}

func (g *UserInfo) Plan(ctx context.Context) (*types.State, []discord.Embed, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	job := g.enricher.GetJobByName(g.job)
	if job == nil {
		g.logger.Warn("unknown job for discord guild, skipping")
		return nil, nil, nil
	}

	roles, err := g.planRoles(job)
	if err != nil {
		return nil, nil, err
	}

	settings := g.settings.Load()
	handlers := []types.UserProcessorHandler{}
	if settings.UserInfoSyncSettings.UnemployedEnabled {
		handlers = append(handlers, func(ctx context.Context, guildId discord.GuildID, member discord.Member, user *types.User) ([]discord.Embed, error) {
			if user.Job == g.job {
				return nil, nil
			}

			if g.checkIfJobIgnored(user.Job) {
				user.Job = g.job
				return nil, nil
			}

			switch settings.UserInfoSyncSettings.UnemployedMode {
			case jobs.UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE:
				user.Roles.Sum = append(user.Roles.Sum, g.unemployedRole)

			case jobs.UserInfoSyncUnemployedMode_USER_INFO_SYNC_UNEMPLOYED_MODE_KICK:
				kick := true
				user.Kick = &kick
				user.KickReason = fmt.Sprintf("no longer an employee of %s job (unemployed mode: kick)", g.job)
			}

			return nil, nil
		})
	}

	users, logs, err := g.planUsers(ctx)
	if err != nil {
		return nil, logs, err
	}

	return &types.State{
		Roles: roles,
		Users: users,

		UserProcessors: handlers,
	}, logs, err
}

func (g *UserInfo) planRoles(job *jobs.Job) (types.Roles, error) {
	roles := types.Roles{}
	settings := g.settings.Load()

	if settings.UserInfoSyncSettings.EmployeeRoleEnabled {
		g.employeeRole = &types.Role{
			Name:   strings.ReplaceAll(strings.ReplaceAll(settings.UserInfoSyncSettings.EmployeeRoleFormat, "%job%", job.Label), "%s", job.Label),
			Module: userInfoRoleModuleEmployee,
			Job:    g.job,
		}
		roles = append(roles, g.employeeRole)
	} else {
		g.employeeRole = nil
	}

	if settings.UserInfoSyncSettings.UnemployedEnabled {
		g.unemployedRole = &types.Role{
			Name:   settings.UserInfoSyncSettings.UnemployedRoleName,
			Module: userInfoRoleModuleUnemployed,
			Job:    g.job,

			KeepIfJobDifferent: true,
		}
		roles = append(roles, g.unemployedRole)
	} else {
		g.unemployedRole = nil
	}

	if settings.JobsAbsence {
		g.absenceRole = &types.Role{
			Name:   settings.JobsAbsenceSettings.AbsenceRole,
			Module: userInfoRoleModuleAbsence,
			Job:    g.job,
		}
		roles = append(roles, g.absenceRole)
	} else {
		g.absenceRole = nil
	}

	g.jobGradeRoles = make(map[int32]*types.Role, len(job.Grades))
	for _, grade := range slices.Backward(job.Grades) {
		name := strings.ReplaceAll(settings.UserInfoSyncSettings.GradeRoleFormat, "%grade_label%", grade.Label)
		name = strings.ReplaceAll(name, "%grade%", fmt.Sprintf("%02d", grade.Grade))
		name = strings.ReplaceAll(name, "%grade_single%", fmt.Sprintf("%d", grade.Grade))

		role := &types.Role{
			Name:   name,
			Module: fmt.Sprintf(userInfoRoleModuleJobGradePrefix+"%d", grade.Grade),
			Job:    g.job,
		}
		g.jobGradeRoles[grade.Grade] = role
		roles = append(roles, role)
	}

	g.groupRoles = make(map[string]*types.Role, len(settings.UserInfoSyncSettings.GroupMapping))
	for i, mapping := range settings.UserInfoSyncSettings.GroupMapping {
		role := &types.Role{
			Name:   mapping.Name,
			Module: fmt.Sprintf(userInfoRoleModuleGroupMappingPrefix+"%d", i),
			Job:    g.job,
		}
		g.groupRoles[mapping.Name] = role
		roles = append(roles, role)
	}

	return roles, nil
}

func (g *UserInfo) planUsers(ctx context.Context) (types.Users, []discord.Embed, error) {
	users := types.Users{}
	logs := []discord.Embed{}
	settings := g.settings.Load()

	jobs := []jet.Expression{jet.String(g.job)}
	for _, job := range g.BaseModule.appCfg.Get().Discord.IgnoredJobs {
		jobs = append(jobs, jet.String(job))
	}

	tUsers := tables.User().AS("users")

	stmt := tAccsOauth2.
		SELECT(
			tAccsOauth2.AccountID.AS("userrolemapping.account_id"),
			tAccsOauth2.ExternalID.AS("userrolemapping.external_id"),
			tUsers.ID.AS("userrolemapping.user_id"),
			tUsers.Job.AS("userrolemapping.job"),
			tUsers.JobGrade.AS("userrolemapping.job_grade"),
			tUsers.Firstname.AS("userrolemapping.firstname"),
			tUsers.Lastname.AS("userrolemapping.lastname"),
			// Job Props
			tColleagueProps.NamePrefix.AS("userrolemapping.name_prefix"),
			tColleagueProps.NameSuffix.AS("userrolemapping.name_suffix"),
			tColleagueProps.AbsenceBegin.AS("userrolemapping.absence_begin"),
			tColleagueProps.AbsenceEnd.AS("userrolemapping.absence_end"),
		).
		FROM(
			tAccsOauth2.
				INNER_JOIN(tAccs,
					tAccs.ID.EQ(tAccsOauth2.AccountID),
				).
				INNER_JOIN(tUsers,
					tUsers.Identifier.LIKE(jet.CONCAT(jet.String("%"), tAccs.License)),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tUsers.ID).
						AND(tColleagueProps.Job.EQ(jet.String(g.job))),
				),
		).
		WHERE(jet.AND(
			tAccsOauth2.Provider.EQ(jet.String("discord")),
			tUsers.Job.IN(jobs...),
		)).
		ORDER_BY(tUsers.ID.ASC())

	var dest []*userRoleMapping
	if err := stmt.QueryContext(ctx, g.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return users, logs, fmt.Errorf("failed to get user info from database. %w", err)
		}
	}

	errs := multierr.Combine()
	for _, u := range dest {
		externalId, err := strconv.ParseUint(u.ExternalID, 10, 64)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to parse user oauth2 external id %d. %w", externalId, err))
			continue
		}

		user := &types.User{
			ID:    discord.UserID(externalId),
			Roles: &types.UserRoles{},
			Job:   u.Job,
		}

		if u.Job != g.job {
			users.Add(user)
			continue
		}

		member, err := g.discord.Member(g.guild.ID, discord.UserID(externalId))
		if err != nil {
			if restErr, ok := err.(*httputil.HTTPError); ok {
				if restErr.Status == http.StatusNotFound {
					// Add log about employee not being on discord
					logs = append(logs, discord.Embed{
						Title:       fmt.Sprintf("UserInfo: Employee not found on Discord: %s %s", u.Firstname, u.Lastname),
						Description: fmt.Sprintf("Discord ID: %d, Rank: %d", externalId, u.JobGrade),
						Author:      embeds.EmbedAuthor,
						Color:       embeds.ColorWarn,
					})
					continue
				}
			}

			errs = multierr.Append(errs, fmt.Errorf("discord API returned an error. %w", err))
			continue
		}

		if settings.UserInfoSyncSettings.SyncNicknames {
			if name := g.getUserNickname(member, strings.TrimSpace(u.Firstname), strings.TrimSpace(u.Lastname), u.NamePrefix, u.NameSuffix); name != nil {
				user.Nickname = name
			}
		}

		user.Roles.Sum, err = g.getUserRoles(u.Job, u.JobGrade)
		if err != nil {
			g.logger.Warn(fmt.Sprintf("failed to get user's job roles %d", externalId), zap.Int32("user_id", u.UserID), zap.Error(err))
			continue
		}

		for _, mapping := range settings.UserInfoSyncSettings.GroupMapping {
			if u.JobGrade < mapping.FromGrade || u.JobGrade > mapping.ToGrade {
				continue
			}

			role, ok := g.groupRoles[mapping.Name]
			if !ok {
				return nil, logs, fmt.Errorf("failed to find role for group mapping %s", mapping.Name)
			}

			user.Roles.Sum = append(user.Roles.Sum, role)
		}

		if settings.UserInfoSyncSettings.EmployeeRoleEnabled &&
			g.employeeRole != nil {
			user.Roles.Sum = append(user.Roles.Sum, g.employeeRole)
		}

		if settings.JobsAbsence && g.absenceRole != nil &&
			g.isUserAbsent(u.AbsenceBegin, u.AbsenceEnd) {
			user.Roles.Sum = append(user.Roles.Sum, g.absenceRole)
		}

		users.Add(user)
	}

	return users, logs, errs
}

func (g *UserInfo) getUserNickname(member *discord.Member, firstname string, lastname string, prefix string, suffix string) *string {
	if g.guild.OwnerID == member.User.ID {
		return nil // Can't change owner's nickname
	}

	targetNickname := g.constructUserNickname(firstname, lastname, prefix, suffix)

	// No need to set nickname when they are already equal
	if member.Nick == targetNickname {
		return nil
	}

	return &targetNickname
}

// constructUserNickname constructs a user's nickname based on the provided input and ensures that it isn't longer than what Discord allows (32 chars)
// prefix and suffix are optional and will be trimmed of spaces, max length is 12/24
func (g *UserInfo) constructUserNickname(firstname string, lastname string, prefix string, suffix string) string {
	maxLength := DiscordNicknameMaxLength

	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)

	if prefix != "" {
		prefix = strings.TrimSpace(prefix) + " "
	}
	if suffix != "" {
		suffix = " " + strings.TrimSpace(suffix)
	}

	last := strings.TrimSpace(lastname + suffix)
	fullName := strings.TrimSpace(prefix + firstname)
	if last != "" {
		fullName += " " + last
	}
	fullName = strings.TrimSpace(fullName)

	// If within limit, return as is
	if len(fullName) <= maxLength {
		return fullName
	}

	// Truncate the firstname progressively
	firstParts := strings.Split(firstname, " ")
	truncatedFirst := ""
	for i := range firstParts {
		if i == len(firstParts)-1 {
			truncatedFirst += string(firstParts[i][0]) + "."
		} else {
			truncatedFirst += firstParts[i] + " "
		}
	}
	truncatedFirst = strings.TrimSpace(truncatedFirst)

	fullName = fmt.Sprintf("%s%s %s%s", prefix, truncatedFirst, lastname, suffix)
	fullName = strings.TrimSpace(fullName)

	if len(fullName) <= maxLength {
		return fullName
	}

	// As a last resort, truncate the last name
	availableLength := maxLength - len(prefix) - len(suffix) // Ensure spacing
	if availableLength > 0 {
		truncatedBase := truncatedFirst + " " + lastname[:availableLength-len(truncatedFirst)-1]
		return fmt.Sprintf("%s%s%s", prefix, truncatedBase, suffix)
	}

	// If even the prefix and suffix alone exceed the limit, just truncate everything
	result := (prefix + suffix)[:maxLength]
	return strings.TrimSpace(result)
}

func (g *UserInfo) getUserRoles(job string, grade int32) (types.Roles, error) {
	userRoles := types.Roles{}

	if g.checkIfJobIgnored(job) {
		return userRoles, nil
	}

	role, ok := g.jobGradeRoles[grade]
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

func (g *UserInfo) watchEvents(e any) {
	switch ev := e.(type) {
	case *gateway.GuildMemberAddEvent:
		settings := g.settings.Load()
		if settings == nil || !settings.UserInfoSyncSettings.UnemployedEnabled {
			return
		}

		func() {
			g.mu.Lock()
			defer g.mu.Unlock()

			if g.unemployedRole == nil {
				return
			}

			if err := g.discord.AddRole(g.guild.ID, ev.User.ID, g.unemployedRole.ID, api.AddRoleData{
				AuditLogReason: api.AuditLogReason("On Join (Unemployed Role)"),
			}); err != nil {
				g.logger.Error("failed to add unemployed role to user on join", zap.Error(err))
				return
			}
		}()

	default:
		g.logger.Warn("unknown event received", zap.Reflect("dc_event_type", e))
	}
}
