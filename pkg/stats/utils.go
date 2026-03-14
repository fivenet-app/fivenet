package stats

import (
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	pbstats "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/stats"
)

func isPublishedDocument(doc *documents.Document) bool {
	if doc == nil {
		return false
	}
	if doc.GetDeletedAt() != nil {
		return false
	}
	if strings.TrimSpace(doc.GetCreatorJob()) == "" {
		return false
	}
	if doc.GetMeta() == nil {
		return false
	}

	return !doc.GetMeta().GetDraft()
}

func periodStartExpr(period pbstats.StatsPeriod) string {
	switch period {
	case pbstats.StatsPeriod_STATS_PERIOD_MONTHLY:
		return "DATE_SUB(day, INTERVAL DAYOFMONTH(day) - 1 DAY)"

	case pbstats.StatsPeriod_STATS_PERIOD_WEEKLY:
		return "DATE_SUB(day, INTERVAL WEEKDAY(day) DAY)"

	case pbstats.StatsPeriod_STATS_PERIOD_DAILY:
		fallthrough
	default:
		return "day"
	}
}

func dateRangeArgs(startDay, endDay time.Time, repeats int) []any {
	args := make([]any, 0, repeats*2)
	start := startDay.Format(time.DateOnly)
	end := endDay.Format(time.DateOnly)

	for range repeats {
		args = append(args, start, end)
	}

	return args
}
