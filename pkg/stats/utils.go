package stats

import (
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	pbstats "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/stats"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
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

func periodStartDateExpr(period pbstats.StatsPeriod) mysql.DateExpression {
	tRollup := table.FivenetStatsDailyRollup

	switch period {
	case pbstats.StatsPeriod_STATS_PERIOD_MONTHLY:
		return mysql.DateExp(
			tRollup.Day.SUB(
				mysql.INTERVAL(
					mysql.Raw("DAYOFMONTH(`fivenet_stats_daily_rollup`.`day`) - 1"),
					mysql.DAY,
				),
			),
		)

	case pbstats.StatsPeriod_STATS_PERIOD_WEEKLY:
		return mysql.DateExp(tRollup.Day.SUB(mysql.INTERVAL(1, mysql.DAY)))

	case pbstats.StatsPeriod_STATS_PERIOD_DAILY:
		fallthrough

	default:
		return tRollup.Day
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
