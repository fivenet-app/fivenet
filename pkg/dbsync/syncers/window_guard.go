package syncers

import (
	"fmt"
	"strings"
	"time"
)

func updatedAtUpperBoundCondition(
	updatedTimeColumn *string,
	windowEnd *time.Time,
) string {
	if windowEnd == nil || updatedTimeColumn == nil || strings.TrimSpace(*updatedTimeColumn) == "" {
		return ""
	}

	updatedAtCol := *updatedTimeColumn
	if !strings.Contains(updatedAtCol, "`") {
		updatedAtCol = fmt.Sprintf("%#q", updatedAtCol)
	}

	return fmt.Sprintf(
		"%s <= '%s'",
		updatedAtCol,
		windowEnd.Truncate(time.Millisecond).Format("2006-01-02 15:04:05.000"),
	)
}
