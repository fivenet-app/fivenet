package calendar

import (
	"context"
	"fmt"
	"strings"

	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const birthdayCalendarNameTranslationKey = "components.calendar.birthday_calendar_name"

func birthdayCalendarDisplayName(
	t func(string, map[string]any) string,
	job string,
	jobInfo *jobs.Job,
) string {
	jobLabel := strings.TrimSpace(job)
	if jobInfo != nil && strings.TrimSpace(jobInfo.GetLabel()) != "" {
		jobLabel = strings.TrimSpace(jobInfo.GetLabel())
	}

	if jobLabel == "" {
		jobLabel = job
	}

	return t(birthdayCalendarNameTranslationKey, map[string]any{
		"job": jobLabel,
	})
}

func birthdayCalendarTitle(
	i18n *i18n.I18n,
	appCfg appconfig.IConfig,
	job string,
	jobInfo *jobs.Job,
) string {
	locale := ""
	if cfg := appCfg.Get(); cfg != nil {
		locale = strings.TrimSpace(cfg.GetDefaultLocale())
	}
	if locale == "" {
		locale = "en"
		if i18n != nil {
			locale = i18n.GetFallbackLanguage()
		}
	}

	translator := func(_ string, vars map[string]any) string {
		if jobLabel, ok := vars["job"].(string); ok && strings.TrimSpace(jobLabel) != "" {
			return fmt.Sprintf("%s Birthdays", jobLabel)
		}
		return fmt.Sprintf("%s Birthdays", job)
	}
	if i18n != nil {
		translator = i18n.Translator(locale)
	}

	return birthdayCalendarDisplayName(translator, job, jobInfo)
}

func birthdayCalendarAccessEntries(
	calendarID int64,
	job string,
	jobInfo *jobs.Job,
) []*calendaraccess.CalendarJobAccess {
	minimumGrade := int32(0)
	highestGrade := int32(0)

	if jobInfo != nil && len(jobInfo.GetGrades()) > 0 {
		minimumGrade = jobInfo.GetGrades()[0].GetGrade()
		highestGrade = jobInfo.GetGrades()[len(jobInfo.GetGrades())-1].GetGrade()
	}

	entries := []*calendaraccess.CalendarJobAccess{
		{
			TargetId:     calendarID,
			Job:          job,
			MinimumGrade: minimumGrade,
			Access:       calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		},
	}

	if highestGrade > minimumGrade {
		entries = append(entries, &calendaraccess.CalendarJobAccess{
			TargetId:     calendarID,
			Job:          job,
			MinimumGrade: highestGrade,
			Access:       calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT,
		})
		return entries
	}

	entries[0].Access = calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT
	return entries
}

func ensureBirthdayCalendarAccess(
	ctx context.Context,
	q qrm.Executable,
	calendarID int64,
	job string,
	jobInfo *jobs.Job,
) error {
	jobAccess := birthdayCalendarAccessEntries(calendarID, job, jobInfo)

	tAccess := table.FivenetCalendarAccess
	if _, err := tAccess.
		DELETE().
		WHERE(tAccess.TargetID.EQ(mysql.Int64(calendarID))).
		ExecContext(ctx, q); err != nil {
		return err
	}

	for i := range jobAccess {
		if _, err := tAccess.
			INSERT(
				tAccess.TargetID,
				tAccess.Access,
				tAccess.Job,
				tAccess.MinimumGrade,
			).
			VALUES(
				jobAccess[i].GetTargetId(),
				jobAccess[i].GetAccess(),
				jobAccess[i].GetJob(),
				jobAccess[i].GetMinimumGrade(),
			).
			ExecContext(ctx, q); err != nil {
			return err
		}
	}

	return nil
}
