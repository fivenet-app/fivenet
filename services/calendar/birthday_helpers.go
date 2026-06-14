package calendar

import (
	"fmt"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
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
	i18n i18n.Ii18n,
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
