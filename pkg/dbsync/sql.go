package dbsync

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func getWhereCondition(
	table DBSyncTable,
	state *TableSyncState,
) string {
	// Add "updatedAt" column condition if available
	if state == nil ||
		(table.UpdatedTimeColumn == nil || (state.LastCheck == nil || state.LastCheck.IsZero())) {
		return ""
	}

	return fmt.Sprintf("`%s` >= '%s'\n",
		*table.UpdatedTimeColumn,
		state.LastCheck.Format("2006-01-02 15:04:05"),
	)
}

func prepareStringQuery(
	q string,
	table DBSyncTable,
	state *TableSyncState,
	offset int64,
	limit int64,
) string {
	offsetStr := strconv.FormatInt(offset, 10)
	limitStr := strconv.FormatInt(limit, 10)

	q = strings.ReplaceAll(q, "$offset", offsetStr)
	q = strings.ReplaceAll(q, "$limit", limitStr)

	where := getWhereCondition(table, state)
	if where != "" {
		// Prepend "WHERE " if there is a condition
		where = "WHERE " + where
	}
	q = strings.ReplaceAll(q, "$whereCondition", where)

	return q
}

func buildQueryFromColumns(
	tableName string,
	columns map[string]string,
	whereCondition []string,
	offset int64,
	limit int64,
) string {
	columnsList := []string{}
	keys := make([]string, 0, len(columns))
	for k := range columns {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, alias := range keys {
		column := columns[alias]
		if column == "" {
			continue
		}

		columnsList = append(columnsList, fmt.Sprintf("`%s` AS `%s`", column, alias))
	}

	q := fmt.Sprintf("SELECT %s\nFROM `%s`\n",
		strings.Join(columnsList, ", "),
		tableName,
	)
	if len(whereCondition) > 0 {
		q += "WHERE " + strings.Join(whereCondition, " AND ")
		q += "\n"
	}

	q += fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	return q
}
